package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/nmn3m/pulsar/backend/internal/domain"
	"github.com/nmn3m/pulsar/backend/internal/usecase/repository"
)

type UpdateProfileRequest struct {
	FullName *string `json:"full_name,omitempty"`
	Phone    *string `json:"phone,omitempty"`
	Timezone *string `json:"timezone,omitempty"`
}

type UserUsecase struct {
	orgRepo  repository.OrganizationRepository
	userRepo repository.UserRepository
}

func NewUserUsecase(orgRepo repository.OrganizationRepository, userRepo repository.UserRepository) *UserUsecase {
	return &UserUsecase{orgRepo: orgRepo, userRepo: userRepo}
}

func (s *UserUsecase) ListOrganizationUsers(ctx context.Context, orgID uuid.UUID) ([]*domain.UserWithOrganization, error) {
	return s.orgRepo.ListUsers(ctx, orgID)
}

func (s *UserUsecase) UpdateProfile(ctx context.Context, userID uuid.UUID, req *UpdateProfileRequest) (*domain.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	if req.FullName != nil {
		user.FullName = req.FullName
	}
	if req.Phone != nil {
		user.Phone = req.Phone
	}
	if req.Timezone != nil {
		user.Timezone = *req.Timezone
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}

	return user, nil
}
