package domain

import (
	"time"

	"github.com/google/uuid"
)

type EmailVerification struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Email     string
	OTP       string
	ExpiresAt time.Time
	Verified  bool
	CreatedAt time.Time
}
