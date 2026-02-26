package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/nmn3m/pulsar/backend/internal/domain"
)

type TeamInvitationRepository interface {
	Create(ctx context.Context, invitation *domain.TeamInvitation) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.TeamInvitation, error)
	GetByToken(ctx context.Context, token string) (*domain.TeamInvitationWithDetails, error)
	GetByEmailAndTeam(ctx context.Context, email string, teamID uuid.UUID) (*domain.TeamInvitation, error)
	Update(ctx context.Context, invitation *domain.TeamInvitation) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListByTeam(ctx context.Context, teamID uuid.UUID) ([]*domain.TeamInvitation, error)
	ListByEmail(ctx context.Context, email string) ([]*domain.TeamInvitationWithDetails, error)
	ListPendingByOrganization(ctx context.Context, orgID uuid.UUID) ([]*domain.TeamInvitationWithDetails, error)
	ExpireOldInvitations(ctx context.Context) error
}
