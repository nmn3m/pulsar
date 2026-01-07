package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nmn3m/pulsar/backend/internal/domain"
)

type EmailVerificationRepository struct {
	db *DB
}

func NewEmailVerificationRepository(db *DB) *EmailVerificationRepository {
	return &EmailVerificationRepository{db: db}
}

func (r *EmailVerificationRepository) Create(ctx context.Context, verification *domain.EmailVerification) error {
	// Delete any existing unverified OTPs for this user first
	deleteQuery := `DELETE FROM email_verifications WHERE user_id = $1 AND verified = FALSE`
	_, err := r.db.ExecContext(ctx, deleteQuery, verification.UserID)
	if err != nil {
		return fmt.Errorf("failed to delete existing verifications: %w", err)
	}

	query := `
		INSERT INTO email_verifications (id, user_id, email, otp, expires_at, verified)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at
	`

	err = r.db.QueryRowContext(
		ctx,
		query,
		verification.ID,
		verification.UserID,
		verification.Email,
		verification.OTP,
		verification.ExpiresAt,
		verification.Verified,
	).Scan(&verification.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create email verification: %w", err)
	}

	return nil
}

func (r *EmailVerificationRepository) GetByEmail(ctx context.Context, email string) (*domain.EmailVerification, error) {
	query := `
		SELECT id, user_id, email, otp, expires_at, verified, created_at
		FROM email_verifications
		WHERE email = $1 AND verified = FALSE
		ORDER BY created_at DESC
		LIMIT 1
	`

	var verification domain.EmailVerification
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&verification.ID,
		&verification.UserID,
		&verification.Email,
		&verification.OTP,
		&verification.ExpiresAt,
		&verification.Verified,
		&verification.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("verification not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get verification: %w", err)
	}

	return &verification, nil
}

func (r *EmailVerificationRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.EmailVerification, error) {
	query := `
		SELECT id, user_id, email, otp, expires_at, verified, created_at
		FROM email_verifications
		WHERE user_id = $1 AND verified = FALSE
		ORDER BY created_at DESC
		LIMIT 1
	`

	var verification domain.EmailVerification
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&verification.ID,
		&verification.UserID,
		&verification.Email,
		&verification.OTP,
		&verification.ExpiresAt,
		&verification.Verified,
		&verification.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("verification not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get verification: %w", err)
	}

	return &verification, nil
}

func (r *EmailVerificationRepository) MarkVerified(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE email_verifications SET verified = TRUE WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to mark verification as verified: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("verification not found")
	}

	return nil
}

func (r *EmailVerificationRepository) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	query := `DELETE FROM email_verifications WHERE user_id = $1`

	_, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete verifications: %w", err)
	}

	return nil
}

func (r *EmailVerificationRepository) DeleteExpired(ctx context.Context) error {
	query := `DELETE FROM email_verifications WHERE expires_at < $1`

	_, err := r.db.ExecContext(ctx, query, time.Now())
	if err != nil {
		return fmt.Errorf("failed to delete expired verifications: %w", err)
	}

	return nil
}
