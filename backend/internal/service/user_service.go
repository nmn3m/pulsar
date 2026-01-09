package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/nmn3m/pulsar/backend/internal/domain"
	"github.com/nmn3m/pulsar/backend/internal/repository"
)

type UserService struct {
	orgRepo repository.OrganizationRepository
}

func NewUserService(orgRepo repository.OrganizationRepository) *UserService {
	return &UserService{
		orgRepo: orgRepo,
	}
}

func (s *UserService) ListOrganizationUsers(ctx context.Context, orgID uuid.UUID) ([]*domain.UserWithOrganization, error) {
	return s.orgRepo.ListUsers(ctx, orgID)
}
