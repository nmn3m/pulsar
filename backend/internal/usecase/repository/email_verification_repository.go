package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/nmn3m/pulsar/backend/internal/domain"
)

type EmailVerificationRepository interface {
	Create(ctx context.Context, verification *domain.EmailVerification) error
	GetByEmail(ctx context.Context, email string) (*domain.EmailVerification, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.EmailVerification, error)
	MarkVerified(ctx context.Context, id uuid.UUID) error
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
	DeleteExpired(ctx context.Context) error
}
