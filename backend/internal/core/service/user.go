package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/nmn3m/pulsar/backend/internal/core/domain"
	"github.com/nmn3m/pulsar/backend/internal/core/dto"
	"github.com/nmn3m/pulsar/backend/internal/core/port/outbound"
)

type UserService struct {
	orgRepo  outbound.OrganizationRepository
	userRepo outbound.UserRepository
}

func NewUserService(orgRepo outbound.OrganizationRepository, userRepo outbound.UserRepository) *UserService {
	return &UserService{orgRepo: orgRepo, userRepo: userRepo}
}

func (s *UserService) ListOrganizationUsers(ctx context.Context, orgID uuid.UUID) ([]*domain.UserWithOrganization, error) {
	return s.orgRepo.ListUsers(ctx, orgID)
}

func (s *UserService) UpdateProfile(ctx context.Context, userID uuid.UUID, req *dto.UpdateProfileRequest) (*domain.User, error) {
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
