package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	"github.com/nmn3m/pulsar/backend/internal/domain"
)

type TeamRepository struct {
	db *DB
}

func NewTeamRepository(db *DB) *TeamRepository {
	return &TeamRepository{db: db}
}

func (r *TeamRepository) Create(ctx context.Context, team *domain.Team) error {
	query := `
		INSERT INTO teams (id, organization_id, name, description)
		VALUES ($1, $2, $3, $4)
		RETURNING created_at, updated_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		team.ID,
		team.OrganizationID,
		team.Name,
		team.Description,
	).Scan(&team.CreatedAt, &team.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create team: %w", err)
	}

	return nil
}

func (r *TeamRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Team, error) {
	query := `
		SELECT id, organization_id, name, description, created_at, updated_at
		FROM teams
		WHERE id = $1
	`

	var team domain.Team
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&team.ID,
		&team.OrganizationID,
		&team.Name,
		&team.Description,
		&team.CreatedAt,
		&team.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("team not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get team: %w", err)
	}

	return &team, nil
}

func (r *TeamRepository) Update(ctx context.Context, team *domain.Team) error {
	query := `
		UPDATE teams
		SET name = $2, description = $3
		WHERE id = $1
		RETURNING updated_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		team.ID,
		team.Name,
		team.Description,
	).Scan(&team.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to update team: %w", err)
	}

	return nil
}

func (r *TeamRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM teams WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete team: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("team not found")
	}

	return nil
}

func (r *TeamRepository) List(ctx context.Context, orgID uuid.UUID, limit, offset int) ([]*domain.Team, error) {
	query := `
		SELECT id, organization_id, name, description, created_at, updated_at
		FROM teams
		WHERE organization_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, orgID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list teams: %w", err)
	}
	defer rows.Close()

	var teams []*domain.Team
	for rows.Next() {
		var team domain.Team
		err := rows.Scan(
			&team.ID,
			&team.OrganizationID,
			&team.Name,
			&team.Description,
			&team.CreatedAt,
			&team.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan team: %w", err)
		}

		teams = append(teams, &team)
	}

	return teams, nil
}

func (r *TeamRepository) AddMember(ctx context.Context, teamID, userID uuid.UUID, role domain.TeamRole) error {
	query := `
		INSERT INTO team_members (team_id, user_id, role)
		VALUES ($1, $2, $3)
	`

	_, err := r.db.ExecContext(ctx, query, teamID, userID, role.String())
	if err != nil {
		return fmt.Errorf("failed to add team member: %w", err)
	}

	return nil
}

func (r *TeamRepository) RemoveMember(ctx context.Context, teamID, userID uuid.UUID) error {
	query := `DELETE FROM team_members WHERE team_id = $1 AND user_id = $2`

	result, err := r.db.ExecContext(ctx, query, teamID, userID)
	if err != nil {
		return fmt.Errorf("failed to remove team member: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("team member not found")
	}

	return nil
}

func (r *TeamRepository) UpdateMemberRole(ctx context.Context, teamID, userID uuid.UUID, role domain.TeamRole) error {
	query := `
		UPDATE team_members
		SET role = $3
		WHERE team_id = $1 AND user_id = $2
	`

	result, err := r.db.ExecContext(ctx, query, teamID, userID, role.String())
	if err != nil {
		return fmt.Errorf("failed to update team member role: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("team member not found")
	}

	return nil
}

func (r *TeamRepository) ListMembers(ctx context.Context, teamID uuid.UUID) ([]*domain.UserWithTeamRole, error) {
	query := `
		SELECT u.id, u.email, u.username, u.full_name, u.phone, u.timezone,
		       u.notification_preferences, u.is_active, u.created_at, u.updated_at,
		       tm.role, tm.joined_at
		FROM users u
		JOIN team_members tm ON u.id = tm.user_id
		WHERE tm.team_id = $1
		ORDER BY tm.joined_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, teamID)
	if err != nil {
		return nil, fmt.Errorf("failed to list team members: %w", err)
	}
	defer rows.Close()

	var users []*domain.UserWithTeamRole
	for rows.Next() {
		var user domain.UserWithTeamRole
		var prefsJSON []byte

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
			&user.Role,
			&user.JoinedAt,
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

func (r *TeamRepository) ListUserTeams(ctx context.Context, userID uuid.UUID) ([]*domain.Team, error) {
	query := `
		SELECT t.id, t.organization_id, t.name, t.description, t.created_at, t.updated_at
		FROM teams t
		JOIN team_members tm ON t.id = tm.team_id
		WHERE tm.user_id = $1
		ORDER BY t.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list user teams: %w", err)
	}
	defer rows.Close()

	var teams []*domain.Team
	for rows.Next() {
		var team domain.Team
		err := rows.Scan(
			&team.ID,
			&team.OrganizationID,
			&team.Name,
			&team.Description,
			&team.CreatedAt,
			&team.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan team: %w", err)
		}

		teams = append(teams, &team)
	}

	return teams, nil
}
