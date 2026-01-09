package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"

	"github.com/nmn3m/pulsar/backend/internal/domain"
)

type DNDSettingsRepository struct {
	db *DB
}

func NewDNDSettingsRepository(db *DB) *DNDSettingsRepository {
	return &DNDSettingsRepository{db: db}
}

func (r *DNDSettingsRepository) Create(ctx context.Context, settings *domain.UserDNDSettings) error {
	query := `
		INSERT INTO user_dnd_settings (id, user_id, enabled, schedule, overrides, allow_p1_override)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at, updated_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		settings.ID,
		settings.UserID,
		settings.Enabled,
		settings.Schedule,
		settings.Overrides,
		settings.AllowP1Override,
	).Scan(&settings.CreatedAt, &settings.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create DND settings: %w", err)
	}

	return nil
}

func (r *DNDSettingsRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.UserDNDSettings, error) {
	query := `
		SELECT id, user_id, enabled, schedule, overrides, allow_p1_override, created_at, updated_at
		FROM user_dnd_settings
		WHERE user_id = $1
	`

	var settings domain.UserDNDSettings
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&settings.ID,
		&settings.UserID,
		&settings.Enabled,
		&settings.Schedule,
		&settings.Overrides,
		&settings.AllowP1Override,
		&settings.CreatedAt,
		&settings.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil // No settings found (not an error)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get DND settings: %w", err)
	}

	return &settings, nil
}

func (r *DNDSettingsRepository) Update(ctx context.Context, settings *domain.UserDNDSettings) error {
	query := `
		UPDATE user_dnd_settings
		SET enabled = $2, schedule = $3, overrides = $4, allow_p1_override = $5, updated_at = NOW()
		WHERE user_id = $1
		RETURNING updated_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		settings.UserID,
		settings.Enabled,
		settings.Schedule,
		settings.Overrides,
		settings.AllowP1Override,
	).Scan(&settings.UpdatedAt)

	if err == sql.ErrNoRows {
		return fmt.Errorf("DND settings not found")
	}
	if err != nil {
		return fmt.Errorf("failed to update DND settings: %w", err)
	}

	return nil
}

func (r *DNDSettingsRepository) Delete(ctx context.Context, userID uuid.UUID) error {
	query := `DELETE FROM user_dnd_settings WHERE user_id = $1`

	result, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete DND settings: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("DND settings not found")
	}

	return nil
}

func (r *DNDSettingsRepository) Upsert(ctx context.Context, settings *domain.UserDNDSettings) error {
	query := `
		INSERT INTO user_dnd_settings (id, user_id, enabled, schedule, overrides, allow_p1_override)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (user_id) DO UPDATE SET
			enabled = EXCLUDED.enabled,
			schedule = EXCLUDED.schedule,
			overrides = EXCLUDED.overrides,
			allow_p1_override = EXCLUDED.allow_p1_override,
			updated_at = NOW()
		RETURNING created_at, updated_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		settings.ID,
		settings.UserID,
		settings.Enabled,
		settings.Schedule,
		settings.Overrides,
		settings.AllowP1Override,
	).Scan(&settings.CreatedAt, &settings.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to upsert DND settings: %w", err)
	}

	return nil
}
