package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"

	"github.com/nmn3m/pulsar/backend/internal/domain"
)

type TeamInvitationRepo struct {
	db *DB
}

func NewTeamInvitationRepo(db *DB) *TeamInvitationRepo {
	return &TeamInvitationRepo{db: db}
}

func (r *TeamInvitationRepo) Create(ctx context.Context, invitation *domain.TeamInvitation) error {
	query := `
		INSERT INTO team_invitations (id, team_id, organization_id, email, role, token, status, invited_by_id, expires_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
	`
	_, err := r.db.ExecContext(ctx, query,
		invitation.ID,
		invitation.TeamID,
		invitation.OrganizationID,
		invitation.Email,
		invitation.Role,
		invitation.Token,
		invitation.Status,
		invitation.InvitedByID,
		invitation.ExpiresAt,
	)
	return err
}

func (r *TeamInvitationRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.TeamInvitation, error) {
	var invitation domain.TeamInvitation
	query := `SELECT * FROM team_invitations WHERE id = $1`
	err := r.db.GetContext(ctx, &invitation, query, id)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("invitation not found")
	}
	return &invitation, err
}

func (r *TeamInvitationRepo) GetByToken(ctx context.Context, token string) (*domain.TeamInvitationWithDetails, error) {
	var invitation domain.TeamInvitationWithDetails
	query := `
		SELECT
			ti.*,
			t.name as team_name,
			u.email as invited_by,
			u.username as inviter_name
		FROM team_invitations ti
		JOIN teams t ON t.id = ti.team_id
		JOIN users u ON u.id = ti.invited_by_id
		WHERE ti.token = $1
	`
	err := r.db.GetContext(ctx, &invitation, query, token)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("invitation not found")
	}
	return &invitation, err
}

func (r *TeamInvitationRepo) GetByEmailAndTeam(ctx context.Context, email string, teamID uuid.UUID) (*domain.TeamInvitation, error) {
	var invitation domain.TeamInvitation
	query := `SELECT * FROM team_invitations WHERE email = $1 AND team_id = $2 AND status = 'pending'`
	err := r.db.GetContext(ctx, &invitation, query, email, teamID)
	if err == sql.ErrNoRows {
		return nil, nil // No existing invitation
	}
	return &invitation, err
}

func (r *TeamInvitationRepo) Update(ctx context.Context, invitation *domain.TeamInvitation) error {
	query := `
		UPDATE team_invitations
		SET status = $1, updated_at = NOW()
		WHERE id = $2
	`
	_, err := r.db.ExecContext(ctx, query, invitation.Status, invitation.ID)
	return err
}

func (r *TeamInvitationRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM team_invitations WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *TeamInvitationRepo) ListByTeam(ctx context.Context, teamID uuid.UUID) ([]*domain.TeamInvitation, error) {
	var invitations []*domain.TeamInvitation
	query := `SELECT * FROM team_invitations WHERE team_id = $1 ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &invitations, query, teamID)
	return invitations, err
}

func (r *TeamInvitationRepo) ListByEmail(ctx context.Context, email string) ([]*domain.TeamInvitationWithDetails, error) {
	var invitations []*domain.TeamInvitationWithDetails
	query := `
		SELECT
			ti.*,
			t.name as team_name,
			u.email as invited_by,
			u.username as inviter_name
		FROM team_invitations ti
		JOIN teams t ON t.id = ti.team_id
		JOIN users u ON u.id = ti.invited_by_id
		WHERE ti.email = $1 AND ti.status = 'pending' AND ti.expires_at > NOW()
		ORDER BY ti.created_at DESC
	`
	err := r.db.SelectContext(ctx, &invitations, query, email)
	return invitations, err
}

func (r *TeamInvitationRepo) ListPendingByOrganization(ctx context.Context, orgID uuid.UUID) ([]*domain.TeamInvitationWithDetails, error) {
	var invitations []*domain.TeamInvitationWithDetails
	query := `
		SELECT
			ti.*,
			t.name as team_name,
			u.email as invited_by,
			u.username as inviter_name
		FROM team_invitations ti
		JOIN teams t ON t.id = ti.team_id
		JOIN users u ON u.id = ti.invited_by_id
		WHERE ti.organization_id = $1 AND ti.status = 'pending'
		ORDER BY ti.created_at DESC
	`
	err := r.db.SelectContext(ctx, &invitations, query, orgID)
	return invitations, err
}

func (r *TeamInvitationRepo) ExpireOldInvitations(ctx context.Context) error {
	query := `
		UPDATE team_invitations
		SET status = 'expired', updated_at = NOW()
		WHERE status = 'pending' AND expires_at < NOW()
	`
	_, err := r.db.ExecContext(ctx, query)
	return err
}
