package inbound

import (
	"context"

	"github.com/google/uuid"

	"github.com/nmn3m/pulsar/backend/internal/core/domain"
	"github.com/nmn3m/pulsar/backend/internal/core/dto"
)

type DNDService interface {
	GetSettings(ctx context.Context, userID uuid.UUID) (*domain.UserDNDSettings, error)
	UpdateSettings(ctx context.Context, userID uuid.UUID, req *dto.UpdateDNDSettingsRequest) (*domain.UserDNDSettings, error)
	AddOverride(ctx context.Context, userID uuid.UUID, req *dto.AddDNDOverrideRequest) (*domain.UserDNDSettings, error)
	RemoveOverride(ctx context.Context, userID uuid.UUID, index int) (*domain.UserDNDSettings, error)
	IsInDNDMode(ctx context.Context, userID uuid.UUID, priority domain.AlertPriority) (bool, error)
	CleanExpiredOverrides(ctx context.Context, userID uuid.UUID) error
	DeleteSettings(ctx context.Context, userID uuid.UUID) error
}
