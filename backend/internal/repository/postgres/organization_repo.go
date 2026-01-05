package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/nmn3m/pulsar/backend/internal/domain"
)

type OrganizationRepository struct {
	db *DB
}

func NewOrganizationRepository(db *DB) *OrganizationRepository {
	return &OrganizationRepository{db: db}
}

func (r *OrganizationRepository) Create(ctx context.Context, org *domain.Organization) error {
	query := `
		INSERT INTO organizations (id, name, slug, plan, settings)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING created_at, updated_at
	`

	settings, err := json.Marshal(org.Settings)
	if err != nil {
		return fmt.Errorf("failed to marshal settings: %w", err)
	}

	err = r.db.QueryRowContext(
		ctx,
		query,
		org.ID,
		org.Name,
		org.Slug,
		org.Plan,
		settings,
	).Scan(&org.CreatedAt, &org.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create organization: %w", err)
	}

	return nil
}

func (r *OrganizationRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Organization, error) {
	query := `
		SELECT id, name, slug, plan, settings, created_at, updated_at
		FROM organizations
		WHERE id = $1
	`

	var org domain.Organization
	var settingsJSON []byte

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&org.ID,
		&org.Name,
		&org.Slug,
		&org.Plan,
		&settingsJSON,
		&org.CreatedAt,
		&org.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("organization not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get organization: %w", err)
	}

	if err := json.Unmarshal(settingsJSON, &org.Settings); err != nil {
		return nil, fmt.Errorf("failed to unmarshal settings: %w", err)
	}

	return &org, nil
}

func (r *OrganizationRepository) GetBySlug(ctx context.Context, slug string) (*domain.Organization, error) {
	query := `
		SELECT id, name, slug, plan, settings, created_at, updated_at
		FROM organizations
		WHERE slug = $1
	`

	var org domain.Organization
	var settingsJSON []byte

	err := r.db.QueryRowContext(ctx, query, slug).Scan(
		&org.ID,
		&org.Name,
		&org.Slug,
		&org.Plan,
		&settingsJSON,
		&org.CreatedAt,
		&org.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("organization not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get organization: %w", err)
	}

	if err := json.Unmarshal(settingsJSON, &org.Settings); err != nil {
		return nil, fmt.Errorf("failed to unmarshal settings: %w", err)
	}

	return &org, nil
}

func (r *OrganizationRepository) Update(ctx context.Context, org *domain.Organization) error {
	query := `
		UPDATE organizations
		SET name = $2, slug = $3, plan = $4, settings = $5
		WHERE id = $1
		RETURNING updated_at
	`

	settings, err := json.Marshal(org.Settings)
	if err != nil {
		return fmt.Errorf("failed to marshal settings: %w", err)
	}

	err = r.db.QueryRowContext(
		ctx,
		query,
		org.ID,
		org.Name,
		org.Slug,
		org.Plan,
		settings,
	).Scan(&org.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to update organization: %w", err)
	}

	return nil
}

func (r *OrganizationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM organizations WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete organization: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("organization not found")
	}

	return nil
}

func (r *OrganizationRepository) List(ctx context.Context, limit, offset int) ([]*domain.Organization, error) {
	query := `
		SELECT id, name, slug, plan, settings, created_at, updated_at
		FROM organizations
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list organizations: %w", err)
	}
	defer rows.Close()

	var orgs []*domain.Organization
	for rows.Next() {
		var org domain.Organization
		var settingsJSON []byte

		err := rows.Scan(
			&org.ID,
			&org.Name,
			&org.Slug,
			&org.Plan,
			&settingsJSON,
			&org.CreatedAt,
			&org.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan organization: %w", err)
		}

		if err := json.Unmarshal(settingsJSON, &org.Settings); err != nil {
			return nil, fmt.Errorf("failed to unmarshal settings: %w", err)
		}

		orgs = append(orgs, &org)
	}

	return orgs, nil
}

func (r *OrganizationRepository) AddUser(ctx context.Context, orgID, userID uuid.UUID, role domain.UserRole) error {
	query := `
		INSERT INTO organization_users (organization_id, user_id, role)
		VALUES ($1, $2, $3)
	`

	_, err := r.db.ExecContext(ctx, query, orgID, userID, role.String())
	if err != nil {
		return fmt.Errorf("failed to add user to organization: %w", err)
	}

	return nil
}

func (r *OrganizationRepository) RemoveUser(ctx context.Context, orgID, userID uuid.UUID) error {
	query := `DELETE FROM organization_users WHERE organization_id = $1 AND user_id = $2`

	result, err := r.db.ExecContext(ctx, query, orgID, userID)
	if err != nil {
		return fmt.Errorf("failed to remove user from organization: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("user not found in organization")
	}

	return nil
}

func (r *OrganizationRepository) GetUserRole(ctx context.Context, orgID, userID uuid.UUID) (domain.UserRole, error) {
	query := `SELECT role FROM organization_users WHERE organization_id = $1 AND user_id = $2`

	var role string
	err := r.db.QueryRowContext(ctx, query, orgID, userID).Scan(&role)
	if err == sql.ErrNoRows {
		return "", fmt.Errorf("user not found in organization")
	}
	if err != nil {
		return "", fmt.Errorf("failed to get user role: %w", err)
	}

	return domain.UserRole(role), nil
}

func (r *OrganizationRepository) UpdateUserRole(ctx context.Context, orgID, userID uuid.UUID, role domain.UserRole) error {
	query := `
		UPDATE organization_users
		SET role = $3
		WHERE organization_id = $1 AND user_id = $2
	`

	result, err := r.db.ExecContext(ctx, query, orgID, userID, role.String())
	if err != nil {
		return fmt.Errorf("failed to update user role: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("user not found in organization")
	}

	return nil
}

func (r *OrganizationRepository) ListUsers(ctx context.Context, orgID uuid.UUID) ([]*domain.UserWithOrganization, error) {
	query := `
		SELECT u.id, u.email, u.username, u.full_name, u.phone, u.timezone,
		       u.notification_preferences, u.is_active, u.created_at, u.updated_at,
		       ou.organization_id, ou.role
		FROM users u
		JOIN organization_users ou ON u.id = ou.user_id
		WHERE ou.organization_id = $1
		ORDER BY u.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, orgID)
	if err != nil {
		return nil, fmt.Errorf("failed to list organization users: %w", err)
	}
	defer rows.Close()

	var users []*domain.UserWithOrganization
	for rows.Next() {
		var user domain.UserWithOrganization
		var prefsJSON []byte
		var role string

		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.Username,
			&user.FullName,
			&user.Phone,
			&user.Timezone,
			&prefsJSON,
			&user.IsActive,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.OrganizationID,
			&role,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}

		if err := json.Unmarshal(prefsJSON, &user.NotificationPreferences); err != nil {
			return nil, fmt.Errorf("failed to unmarshal notification preferences: %w", err)
		}

		user.Role = domain.UserRole(role)
		users = append(users, &user)
	}

	return users, nil
}

func (r *OrganizationRepository) ListUserOrganizations(ctx context.Context, userID uuid.UUID) ([]*domain.Organization, error) {
	query := `
		SELECT o.id, o.name, o.slug, o.plan, o.settings, o.created_at, o.updated_at
		FROM organizations o
		JOIN organization_users ou ON o.id = ou.organization_id
		WHERE ou.user_id = $1
		ORDER BY o.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list user organizations: %w", err)
	}
	defer rows.Close()

	var orgs []*domain.Organization
	for rows.Next() {
		var org domain.Organization
		var settingsJSON []byte

		err := rows.Scan(
			&org.ID,
			&org.Name,
			&org.Slug,
			&org.Plan,
			&settingsJSON,
			&org.CreatedAt,
			&org.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan organization: %w", err)
		}

		if err := json.Unmarshal(settingsJSON, &org.Settings); err != nil {
			return nil, fmt.Errorf("failed to unmarshal settings: %w", err)
		}

		orgs = append(orgs, &org)
	}

	return orgs, nil
}
