package outbound

import (
	"context"

	"github.com/google/uuid"

	"github.com/nmn3m/pulsar/backend/internal/core/domain"
)

type OrganizationRepository interface {
	Create(ctx context.Context, org *domain.Organization) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Organization, error)
	GetBySlug(ctx context.Context, slug string) (*domain.Organization, error)
	Update(ctx context.Context, org *domain.Organization) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*domain.Organization, error)
	AddUser(ctx context.Context, orgID, userID uuid.UUID, role domain.UserRole) error
	RemoveUser(ctx context.Context, orgID, userID uuid.UUID) error
	GetUserRole(ctx context.Context, orgID, userID uuid.UUID) (domain.UserRole, error)
	UpdateUserRole(ctx context.Context, orgID, userID uuid.UUID, role domain.UserRole) error
	ListUsers(ctx context.Context, orgID uuid.UUID) ([]*domain.UserWithOrganization, error)
	ListUserOrganizations(ctx context.Context, userID uuid.UUID) ([]*domain.Organization, error)
}
