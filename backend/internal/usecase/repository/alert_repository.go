package repository

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/nmn3m/pulsar/backend/internal/domain"
)

type AlertRepository interface {
	Create(ctx context.Context, alert *domain.Alert) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Alert, error)
	Update(ctx context.Context, alert *domain.Alert) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filter *domain.AlertFilter) ([]*domain.Alert, int, error)
	Acknowledge(ctx context.Context, id, userID uuid.UUID) error
	Close(ctx context.Context, id, userID uuid.UUID, reason string) error
	Snooze(ctx context.Context, id uuid.UUID, until time.Time) error
	Assign(ctx context.Context, id uuid.UUID, userID, teamID *uuid.UUID) error
	FindByDedupKey(ctx context.Context, orgID uuid.UUID, dedupKey string) (*domain.Alert, error)
	IncrementDedupCount(ctx context.Context, id uuid.UUID) error
}
