package usecase

import (
	"context"

	"github.com/google/uuid"

	"github.com/nmn3m/pulsar/backend/internal/domain"
	"github.com/nmn3m/pulsar/backend/internal/usecase/repository"
)

type UserUsecase struct {
	orgRepo repository.OrganizationRepository
}

func NewUserUsecase(orgRepo repository.OrganizationRepository) *UserUsecase {
	return &UserUsecase{orgRepo: orgRepo}
}

func (s *UserUsecase) ListOrganizationUsers(ctx context.Context, orgID uuid.UUID) ([]*domain.UserWithOrganization, error) {
	return s.orgRepo.ListUsers(ctx, orgID)
}
