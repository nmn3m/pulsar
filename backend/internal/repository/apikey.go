package repository

import (
	"context"

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
