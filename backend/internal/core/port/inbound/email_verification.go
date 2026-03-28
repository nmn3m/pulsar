package inbound

import (
	"context"

	"github.com/google/uuid"

	"github.com/nmn3m/pulsar/backend/internal/core/domain"
)

type EmailVerificationService interface {
	IsEmailServiceConfigured() bool
	GenerateOTP() (string, error)
	CreateAndSendOTP(ctx context.Context, userID uuid.UUID, email, username string) error
	VerifyOTP(ctx context.Context, email, otp string) error
	ResendOTP(ctx context.Context, email string) error
	GetPendingVerification(ctx context.Context, userID uuid.UUID) (*domain.EmailVerification, error)
	CleanupExpired(ctx context.Context) error
}
