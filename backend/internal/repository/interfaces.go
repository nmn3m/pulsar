package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/pulsar/backend/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*domain.User, error)
}

type OrganizationRepository interface {
	Create(ctx context.Context, org *domain.Organization) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Organization, error)
	GetBySlug(ctx context.Context, slug string) (*domain.Organization, error)
	Update(ctx context.Context, org *domain.Organization) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*domain.Organization, error)

	// Organization user methods
	AddUser(ctx context.Context, orgID, userID uuid.UUID, role domain.UserRole) error
	RemoveUser(ctx context.Context, orgID, userID uuid.UUID) error
	GetUserRole(ctx context.Context, orgID, userID uuid.UUID) (domain.UserRole, error)
	UpdateUserRole(ctx context.Context, orgID, userID uuid.UUID, role domain.UserRole) error
	ListUsers(ctx context.Context, orgID uuid.UUID) ([]*domain.UserWithOrganization, error)
	ListUserOrganizations(ctx context.Context, userID uuid.UUID) ([]*domain.Organization, error)
}

type AlertRepository interface {
	Create(ctx context.Context, alert *domain.Alert) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Alert, error)
	Update(ctx context.Context, alert *domain.Alert) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filter *domain.AlertFilter) ([]*domain.Alert, int, error)

	// Alert actions
	Acknowledge(ctx context.Context, id, userID uuid.UUID) error
	Close(ctx context.Context, id, userID uuid.UUID, reason string) error
	Snooze(ctx context.Context, id uuid.UUID, until time.Time) error
	Assign(ctx context.Context, id uuid.UUID, userID, teamID *uuid.UUID) error
}

type TeamRepository interface {
	Create(ctx context.Context, team *domain.Team) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Team, error)
	Update(ctx context.Context, team *domain.Team) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, orgID uuid.UUID, limit, offset int) ([]*domain.Team, error)

	// Team member methods
	AddMember(ctx context.Context, teamID, userID uuid.UUID, role domain.TeamRole) error
	RemoveMember(ctx context.Context, teamID, userID uuid.UUID) error
	UpdateMemberRole(ctx context.Context, teamID, userID uuid.UUID, role domain.TeamRole) error
	ListMembers(ctx context.Context, teamID uuid.UUID) ([]*domain.UserWithTeamRole, error)
	ListUserTeams(ctx context.Context, userID uuid.UUID) ([]*domain.Team, error)
}
