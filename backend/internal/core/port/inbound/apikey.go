package inbound

import (
	"context"

	"github.com/google/uuid"

	"github.com/nmn3m/pulsar/backend/internal/core/domain"
	"github.com/nmn3m/pulsar/backend/internal/core/dto"
)

type APIKeyService interface {
	CreateAPIKey(ctx context.Context, orgID, userID uuid.UUID, req *dto.CreateAPIKeyRequest) (*dto.APIKeyResponse, error)
	ValidateAPIKey(ctx context.Context, rawKey string) (*domain.APIKey, error)
	GetAPIKey(ctx context.Context, id uuid.UUID) (*domain.APIKey, error)
	ListAPIKeys(ctx context.Context, orgID uuid.UUID) ([]domain.APIKey, error)
	ListUserAPIKeys(ctx context.Context, userID uuid.UUID) ([]domain.APIKey, error)
	UpdateAPIKey(ctx context.Context, id uuid.UUID, req *dto.UpdateAPIKeyRequest) (*domain.APIKey, error)
	RevokeAPIKey(ctx context.Context, id uuid.UUID) error
	DeleteAPIKey(ctx context.Context, id uuid.UUID) error
	RevokeAllUserAPIKeys(ctx context.Context, userID uuid.UUID) error
	CheckScope(key *domain.APIKey, scope domain.APIKeyScope) bool
}
