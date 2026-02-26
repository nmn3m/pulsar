package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/nmn3m/pulsar/backend/internal/domain"
)

type DNDSettingsRepository interface {
	Create(ctx context.Context, settings *domain.UserDNDSettings) error
	GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.UserDNDSettings, error)
	Update(ctx context.Context, settings *domain.UserDNDSettings) error
	Delete(ctx context.Context, userID uuid.UUID) error
	Upsert(ctx context.Context, settings *domain.UserDNDSettings) error
}
