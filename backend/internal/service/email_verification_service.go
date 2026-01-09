package service

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/google/uuid"

	"github.com/nmn3m/pulsar/backend/internal/domain"
	"github.com/nmn3m/pulsar/backend/internal/repository"
)

const (
	OTPLength     = 6
	OTPExpiration = 10 * time.Minute
)

type EmailVerificationService struct {
	verificationRepo repository.EmailVerificationRepository
	userRepo         repository.UserRepository
	emailService     *EmailService
}

func NewEmailVerificationService(
	verificationRepo repository.EmailVerificationRepository,
	userRepo repository.UserRepository,
	emailService *EmailService,
) *EmailVerificationService {
	return &EmailVerificationService{
		verificationRepo: verificationRepo,
		userRepo:         userRepo,
		emailService:     emailService,
	}
}

func (s *EmailVerificationService) IsEmailServiceConfigured() bool {
	return s.emailService != nil && s.emailService.IsConfigured()
}

// GenerateOTP generates a random 6-digit OTP
func (s *EmailVerificationService) GenerateOTP() (string, error) {
	const digits = "0123456789"
	otp := make([]byte, OTPLength)

	for i := range otp {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", fmt.Errorf("failed to generate OTP: %w", err)
		}
		otp[i] = digits[num.Int64()]
	}

	return string(otp), nil
}

// CreateAndSendOTP creates a new OTP and sends it via email
func (s *EmailVerificationService) CreateAndSendOTP(ctx context.Context, userID uuid.UUID, email, username string) error {
	// Generate OTP
	otp, err := s.GenerateOTP()
	if err != nil {
		return fmt.Errorf("failed to generate OTP: %w", err)
	}

	// Create verification record
	verification := &domain.EmailVerification{
		ID:        uuid.New(),
		UserID:    userID,
		Email:     email,
		OTP:       otp,
		ExpiresAt: time.Now().Add(OTPExpiration),
		Verified:  false,
	}

	if err := s.verificationRepo.Create(ctx, verification); err != nil {
		return fmt.Errorf("failed to create verification record: %w", err)
	}

	// Send email if configured
	if s.IsEmailServiceConfigured() {
		if err := s.emailService.SendOTPEmail(email, otp, username); err != nil {
			return fmt.Errorf("failed to send OTP email: %w", err)
		}
	}

	return nil
}

// VerifyOTP verifies the OTP and marks the email as verified
func (s *EmailVerificationService) VerifyOTP(ctx context.Context, email, otp string) error {
	// Get the verification record
	verification, err := s.verificationRepo.GetByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("invalid or expired OTP")
	}

	// Check if OTP has expired
	if time.Now().After(verification.ExpiresAt) {
		return fmt.Errorf("OTP has expired")
	}

	// Check if OTP matches
	if verification.OTP != otp {
		return fmt.Errorf("invalid OTP")
	}

	// Mark verification as verified
	if err := s.verificationRepo.MarkVerified(ctx, verification.ID); err != nil {
		return fmt.Errorf("failed to mark verification: %w", err)
	}

	// Update user's email_verified status
	user, err := s.userRepo.GetByID(ctx, verification.UserID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	user.EmailVerified = true
	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// ResendOTP creates a new OTP and sends it
func (s *EmailVerificationService) ResendOTP(ctx context.Context, email string) error {
	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	// Check if already verified
	if user.EmailVerified {
		return fmt.Errorf("email already verified")
	}

	// Create and send new OTP
	return s.CreateAndSendOTP(ctx, user.ID, email, user.Username)
}

// GetPendingVerification gets the pending verification for a user
func (s *EmailVerificationService) GetPendingVerification(ctx context.Context, userID uuid.UUID) (*domain.EmailVerification, error) {
	return s.verificationRepo.GetByUserID(ctx, userID)
}

// CleanupExpired removes expired verification records
func (s *EmailVerificationService) CleanupExpired(ctx context.Context) error {
	return s.verificationRepo.DeleteExpired(ctx)
}
