package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/nmn3m/pulsar/backend/internal/domain"
	"github.com/nmn3m/pulsar/backend/internal/repository"
)

type apiKeyRepository struct {
	db *sqlx.DB
}

func NewAPIKeyRepository(db *sqlx.DB) repository.APIKeyRepository {
	return &apiKeyRepository{db: db}
}

func (r *apiKeyRepository) Create(ctx context.Context, key *domain.APIKey) error {
	query := `
		INSERT INTO api_keys (
			id, organization_id, user_id, name, key_prefix, key_hash, scopes,
			expires_at, is_active, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
		)
	`

	now := time.Now()
	key.CreatedAt = now
	key.UpdatedAt = now

	_, err := r.db.ExecContext(ctx, query,
		key.ID,
		key.OrganizationID,
		key.UserID,
		key.Name,
		key.KeyPrefix,
		key.KeyHash,
		pq.StringArray(key.Scopes),
		key.ExpiresAt,
		key.IsActive,
		key.CreatedAt,
		key.UpdatedAt,
	)

	return err
}

func (r *apiKeyRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.APIKey, error) {
	query := `
		SELECT id, organization_id, user_id, name, key_prefix, key_hash, scopes,
			last_used_at, expires_at, is_active, created_at, updated_at
		FROM api_keys
		WHERE id = $1
	`

	var key domain.APIKey
	var scopes pq.StringArray

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&key.ID,
		&key.OrganizationID,
		&key.UserID,
		&key.Name,
		&key.KeyPrefix,
		&key.KeyHash,
		&scopes,
		&key.LastUsedAt,
		&key.ExpiresAt,
		&key.IsActive,
		&key.CreatedAt,
		&key.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	key.Scopes = scopes
	return &key, nil
}

func (r *apiKeyRepository) GetByHash(ctx context.Context, keyHash string) (*domain.APIKey, error) {
	query := `
		SELECT id, organization_id, user_id, name, key_prefix, key_hash, scopes,
			last_used_at, expires_at, is_active, created_at, updated_at
		FROM api_keys
		WHERE key_hash = $1 AND is_active = true
	`

	var key domain.APIKey
	var scopes pq.StringArray

	err := r.db.QueryRowContext(ctx, query, keyHash).Scan(
		&key.ID,
		&key.OrganizationID,
		&key.UserID,
		&key.Name,
		&key.KeyPrefix,
		&key.KeyHash,
		&scopes,
		&key.LastUsedAt,
		&key.ExpiresAt,
		&key.IsActive,
		&key.CreatedAt,
		&key.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	key.Scopes = scopes
	return &key, nil
}

func (r *apiKeyRepository) ListByOrganization(ctx context.Context, orgID uuid.UUID) ([]domain.APIKey, error) {
	query := `
		SELECT id, organization_id, user_id, name, key_prefix, key_hash, scopes,
			last_used_at, expires_at, is_active, created_at, updated_at
		FROM api_keys
		WHERE organization_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var keys []domain.APIKey
	for rows.Next() {
		var key domain.APIKey
		var scopes pq.StringArray
		if err := rows.Scan(
			&key.ID,
			&key.OrganizationID,
			&key.UserID,
			&key.Name,
			&key.KeyPrefix,
			&key.KeyHash,
			&scopes,
			&key.LastUsedAt,
			&key.ExpiresAt,
			&key.IsActive,
			&key.CreatedAt,
			&key.UpdatedAt,
		); err != nil {
			return nil, err
		}
		key.Scopes = scopes
		keys = append(keys, key)
	}

	return keys, nil
}

func (r *apiKeyRepository) ListByUser(ctx context.Context, userID uuid.UUID) ([]domain.APIKey, error) {
	query := `
		SELECT id, organization_id, user_id, name, key_prefix, key_hash, scopes,
			last_used_at, expires_at, is_active, created_at, updated_at
		FROM api_keys
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var keys []domain.APIKey
	for rows.Next() {
		var key domain.APIKey
		var scopes pq.StringArray
		if err := rows.Scan(
			&key.ID,
			&key.OrganizationID,
			&key.UserID,
			&key.Name,
			&key.KeyPrefix,
			&key.KeyHash,
			&scopes,
			&key.LastUsedAt,
			&key.ExpiresAt,
			&key.IsActive,
			&key.CreatedAt,
			&key.UpdatedAt,
		); err != nil {
			return nil, err
		}
		key.Scopes = scopes
		keys = append(keys, key)
	}

	return keys, nil
}

func (r *apiKeyRepository) Update(ctx context.Context, key *domain.APIKey) error {
	query := `
		UPDATE api_keys
		SET name = $1, scopes = $2, is_active = $3, updated_at = $4
		WHERE id = $5
	`

	key.UpdatedAt = time.Now()

	result, err := r.db.ExecContext(ctx, query,
		key.Name,
		pq.StringArray(key.Scopes),
		key.IsActive,
		key.UpdatedAt,
		key.ID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *apiKeyRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM api_keys WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *apiKeyRepository) UpdateLastUsed(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE api_keys SET last_used_at = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, time.Now(), id)
	return err
}

func (r *apiKeyRepository) RevokeAllByUser(ctx context.Context, userID uuid.UUID) error {
	query := `UPDATE api_keys SET is_active = false, updated_at = $1 WHERE user_id = $2`
	_, err := r.db.ExecContext(ctx, query, time.Now(), userID)
	return err
}
