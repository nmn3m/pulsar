package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/nmn3m/pulsar/backend/internal/domain"
)

type UserRepository struct {
	db *DB
}

func NewUserRepository(db *DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (id, email, username, password_hash, full_name, phone, timezone, notification_preferences, is_active, email_verified)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING created_at, updated_at
	`

	prefs, err := json.Marshal(user.NotificationPreferences)
	if err != nil {
		return fmt.Errorf("failed to marshal notification preferences: %w", err)
	}

	err = r.db.QueryRowContext(
		ctx,
		query,
		user.ID,
		user.Email,
		user.Username,
		user.PasswordHash,
		user.FullName,
		user.Phone,
		user.Timezone,
		prefs,
		user.IsActive,
		user.EmailVerified,
	).Scan(&user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	query := `
		SELECT id, email, username, password_hash, full_name, phone, timezone,
		       notification_preferences, is_active, email_verified, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user domain.User
	var prefsJSON []byte

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.PasswordHash,
		&user.FullName,
		&user.Phone,
		&user.Timezone,
		&prefsJSON,
		&user.IsActive,
		&user.EmailVerified,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if err := json.Unmarshal(prefsJSON, &user.NotificationPreferences); err != nil {
		return nil, fmt.Errorf("failed to unmarshal notification preferences: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, email, username, password_hash, full_name, phone, timezone,
		       notification_preferences, is_active, email_verified, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	var user domain.User
	var prefsJSON []byte

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.PasswordHash,
		&user.FullName,
		&user.Phone,
		&user.Timezone,
		&prefsJSON,
		&user.IsActive,
		&user.EmailVerified,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if err := json.Unmarshal(prefsJSON, &user.NotificationPreferences); err != nil {
		return nil, fmt.Errorf("failed to unmarshal notification preferences: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	query := `
		SELECT id, email, username, password_hash, full_name, phone, timezone,
		       notification_preferences, is_active, email_verified, created_at, updated_at
		FROM users
		WHERE username = $1
	`

	var user domain.User
	var prefsJSON []byte

	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.PasswordHash,
		&user.FullName,
		&user.Phone,
		&user.Timezone,
		&prefsJSON,
		&user.IsActive,
		&user.EmailVerified,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if err := json.Unmarshal(prefsJSON, &user.NotificationPreferences); err != nil {
		return nil, fmt.Errorf("failed to unmarshal notification preferences: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users
		SET email = $2, username = $3, full_name = $4, phone = $5, timezone = $6,
		    notification_preferences = $7, is_active = $8, email_verified = $9
		WHERE id = $1
		RETURNING updated_at
	`

	prefs, err := json.Marshal(user.NotificationPreferences)
	if err != nil {
		return fmt.Errorf("failed to marshal notification preferences: %w", err)
	}

	err = r.db.QueryRowContext(
		ctx,
		query,
		user.ID,
		user.Email,
		user.Username,
		user.FullName,
		user.Phone,
		user.Timezone,
		prefs,
		user.IsActive,
		user.EmailVerified,
	).Scan(&user.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (r *UserRepository) List(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	query := `
		SELECT id, email, username, password_hash, full_name, phone, timezone,
		       notification_preferences, is_active, email_verified, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var user domain.User
		var prefsJSON []byte

		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.Username,
			&user.PasswordHash,
			&user.FullName,
			&user.Phone,
			&user.Timezone,
			&prefsJSON,
			&user.IsActive,
			&user.EmailVerified,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}

		if err := json.Unmarshal(prefsJSON, &user.NotificationPreferences); err != nil {
			return nil, fmt.Errorf("failed to unmarshal notification preferences: %w", err)
		}

		users = append(users, &user)
	}

	return users, nil
}

func (r *UserRepository) ListByTeam(ctx context.Context, teamID uuid.UUID) ([]*domain.User, error) {
	query := `
		SELECT u.id, u.email, u.username, u.password_hash, u.full_name, u.phone, u.timezone,
		       u.notification_preferences, u.is_active, u.email_verified, u.created_at, u.updated_at
		FROM users u
		JOIN team_members tm ON u.id = tm.user_id
		WHERE tm.team_id = $1
		ORDER BY u.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, teamID)
	if err != nil {
		return nil, fmt.Errorf("failed to list users by team: %w", err)
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var user domain.User
		var prefsJSON []byte

		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.Username,
			&user.PasswordHash,
			&user.FullName,
			&user.Phone,
			&user.Timezone,
			&prefsJSON,
			&user.IsActive,
			&user.EmailVerified,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}

		if err := json.Unmarshal(prefsJSON, &user.NotificationPreferences); err != nil {
			return nil, fmt.Errorf("failed to unmarshal notification preferences: %w", err)
		}

		users = append(users, &user)
	}

	return users, nil
}
