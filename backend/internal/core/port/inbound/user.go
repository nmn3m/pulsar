package inbound

import (
	"context"

	"github.com/google/uuid"

	"github.com/nmn3m/pulsar/backend/internal/core/domain"
	"github.com/nmn3m/pulsar/backend/internal/core/dto"
)

type UserService interface {
	ListOrganizationUsers(ctx context.Context, orgID uuid.UUID) ([]*domain.UserWithOrganization, error)
	UpdateProfile(ctx context.Context, userID uuid.UUID, req *dto.UpdateProfileRequest) (*domain.User, error)
}
