package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/nmn3m/pulsar/backend/internal/domain"
)

type TeamRepository interface {
	Create(ctx context.Context, team *domain.Team) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Team, error)
	Update(ctx context.Context, team *domain.Team) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, orgID uuid.UUID, limit, offset int) ([]*domain.Team, error)
	AddMember(ctx context.Context, teamID, userID uuid.UUID, role domain.TeamRole) error
	RemoveMember(ctx context.Context, teamID, userID uuid.UUID) error
	UpdateMemberRole(ctx context.Context, teamID, userID uuid.UUID, role domain.TeamRole) error
	ListMembers(ctx context.Context, teamID uuid.UUID) ([]*domain.UserWithTeamRole, error)
	ListUserTeams(ctx context.Context, userID uuid.UUID) ([]*domain.Team, error)
}
