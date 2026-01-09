package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/nmn3m/pulsar/backend/internal/domain"
)

type APIKeyRepository interface {
	Create(ctx context.Context, key *domain.APIKey) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.APIKey, error)
	GetByHash(ctx context.Context, keyHash string) (*domain.APIKey, error)
	ListByOrganization(ctx context.Context, orgID uuid.UUID) ([]domain.APIKey, error)
	ListByUser(ctx context.Context, userID uuid.UUID) ([]domain.APIKey, error)
	Update(ctx context.Context, key *domain.APIKey) error
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateLastUsed(ctx context.Context, id uuid.UUID) error
	RevokeAllByUser(ctx context.Context, userID uuid.UUID) error
}

type APIKeyService struct {
	repo APIKeyRepository
}

func NewAPIKeyService(repo APIKeyRepository) *APIKeyService {
	return &APIKeyService{repo: repo}
}

// CreateAPIKey creates a new API key and returns the raw key (only shown once)
func (s *APIKeyService) CreateAPIKey(ctx context.Context, orgID, userID uuid.UUID, req *domain.CreateAPIKeyRequest) (*domain.APIKeyResponse, error) {
	// Validate scopes
	for _, scope := range req.Scopes {
		if !domain.IsValidScope(scope) {
			return nil, fmt.Errorf("invalid scope: %s", scope)
		}
	}

	// Generate the API key
	rawKey, keyPrefix, keyHash, err := domain.GenerateAPIKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate API key: %w", err)
	}

	// Parse expiration if provided
	var expiresAt *time.Time
	if req.ExpiresAt != nil {
		t, err := time.Parse(time.RFC3339, *req.ExpiresAt)
		if err != nil {
			return nil, fmt.Errorf("invalid expires_at format, use RFC3339: %w", err)
		}
		if t.Before(time.Now()) {
			return nil, fmt.Errorf("expires_at must be in the future")
		}
		expiresAt = &t
	}

	key := &domain.APIKey{
		ID:             uuid.New(),
		OrganizationID: orgID,
		UserID:         userID,
		Name:           req.Name,
		KeyPrefix:      keyPrefix,
		KeyHash:        keyHash,
		Scopes:         req.Scopes,
		ExpiresAt:      expiresAt,
		IsActive:       true,
	}

	if err := s.repo.Create(ctx, key); err != nil {
		return nil, fmt.Errorf("failed to create API key: %w", err)
	}

	// Return the response with the raw key (only shown once)
	return &domain.APIKeyResponse{
		APIKey: key,
		RawKey: rawKey,
	}, nil
}

// ValidateAPIKey validates an API key and returns the key details if valid
func (s *APIKeyService) ValidateAPIKey(ctx context.Context, rawKey string) (*domain.APIKey, error) {
	// Hash the provided key
	keyHash := domain.HashAPIKey(rawKey)

	// Look up the key by hash
	key, err := s.repo.GetByHash(ctx, keyHash)
	if err != nil {
		return nil, domain.ErrUnauthorized
	}

	// Check if key is active
	if !key.IsActive {
		return nil, domain.ErrUnauthorized
	}

	// Check if key has expired
	if key.ExpiresAt != nil && key.ExpiresAt.Before(time.Now()) {
		return nil, domain.ErrUnauthorized
	}

	// Update last used timestamp (async, don't block on errors)
	go func() {
		_ = s.repo.UpdateLastUsed(context.Background(), key.ID)
	}()

	return key, nil
}

// GetAPIKey gets an API key by ID
func (s *APIKeyService) GetAPIKey(ctx context.Context, id uuid.UUID) (*domain.APIKey, error) {
	return s.repo.GetByID(ctx, id)
}

// ListAPIKeys lists all API keys for an organization
func (s *APIKeyService) ListAPIKeys(ctx context.Context, orgID uuid.UUID) ([]domain.APIKey, error) {
	return s.repo.ListByOrganization(ctx, orgID)
}

// ListUserAPIKeys lists all API keys for a user
func (s *APIKeyService) ListUserAPIKeys(ctx context.Context, userID uuid.UUID) ([]domain.APIKey, error) {
	return s.repo.ListByUser(ctx, userID)
}

// UpdateAPIKey updates an API key
func (s *APIKeyService) UpdateAPIKey(ctx context.Context, id uuid.UUID, req *domain.UpdateAPIKeyRequest) (*domain.APIKey, error) {
	key, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		key.Name = *req.Name
	}

	if len(req.Scopes) > 0 {
		// Validate scopes
		for _, scope := range req.Scopes {
			if !domain.IsValidScope(scope) {
				return nil, fmt.Errorf("invalid scope: %s", scope)
			}
		}
		key.Scopes = req.Scopes
	}

	if req.IsActive != nil {
		key.IsActive = *req.IsActive
	}

	if err := s.repo.Update(ctx, key); err != nil {
		return nil, err
	}

	return key, nil
}

// RevokeAPIKey revokes (deactivates) an API key
func (s *APIKeyService) RevokeAPIKey(ctx context.Context, id uuid.UUID) error {
	key, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	key.IsActive = false
	return s.repo.Update(ctx, key)
}

// DeleteAPIKey permanently deletes an API key
func (s *APIKeyService) DeleteAPIKey(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

// RevokeAllUserAPIKeys revokes all API keys for a user
func (s *APIKeyService) RevokeAllUserAPIKeys(ctx context.Context, userID uuid.UUID) error {
	return s.repo.RevokeAllByUser(ctx, userID)
}

// CheckScope checks if an API key has a specific scope
func (s *APIKeyService) CheckScope(key *domain.APIKey, scope domain.APIKeyScope) bool {
	return key.HasScope(scope)
}
