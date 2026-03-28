package inbound

import (
	"context"

	"github.com/google/uuid"

	"github.com/nmn3m/pulsar/backend/internal/core/domain"
	"github.com/nmn3m/pulsar/backend/internal/core/dto"
)

type TeamService interface {
	CreateTeam(ctx context.Context, orgID uuid.UUID, req *dto.CreateTeamRequest) (*domain.Team, error)
	GetTeam(ctx context.Context, id uuid.UUID) (*domain.Team, error)
	GetTeamWithMembers(ctx context.Context, id uuid.UUID) (*dto.TeamWithMembers, error)
	UpdateTeam(ctx context.Context, id uuid.UUID, req *dto.UpdateTeamRequest) (*domain.Team, error)
	DeleteTeam(ctx context.Context, id uuid.UUID) error
	ListTeams(ctx context.Context, orgID uuid.UUID, page, pageSize int) ([]*domain.Team, error)
	AddMember(ctx context.Context, teamID uuid.UUID, req *dto.AddTeamMemberRequest) error
	AddMemberOrInvite(ctx context.Context, teamID, orgID, inviterID uuid.UUID, req *dto.InviteMemberRequest) (*dto.InvitationResponse, error)
	AcceptInvitation(ctx context.Context, token string, userID uuid.UUID) error
	DeclineInvitation(ctx context.Context, token string) error
	GetPendingInvitations(ctx context.Context, email string) ([]*domain.TeamInvitationWithDetails, error)
	ListTeamInvitations(ctx context.Context, teamID uuid.UUID) ([]*domain.TeamInvitation, error)
	CancelInvitation(ctx context.Context, invitationID uuid.UUID) error
	ResendInvitation(ctx context.Context, invitationID uuid.UUID) error
	RemoveMember(ctx context.Context, teamID, userID uuid.UUID) error
	UpdateMemberRole(ctx context.Context, teamID, userID uuid.UUID, req *dto.UpdateTeamMemberRoleRequest) error
	ListMembers(ctx context.Context, teamID uuid.UUID) ([]*domain.UserWithTeamRole, error)
	ListUserTeams(ctx context.Context, userID uuid.UUID) ([]*domain.Team, error)
}
