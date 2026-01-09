package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/nmn3m/pulsar/backend/internal/domain"
	"github.com/nmn3m/pulsar/backend/internal/repository"
)

type TeamService struct {
	teamRepo repository.TeamRepository
	userRepo repository.UserRepository
}

func NewTeamService(teamRepo repository.TeamRepository, userRepo repository.UserRepository) *TeamService {
	return &TeamService{
		teamRepo: teamRepo,
		userRepo: userRepo,
	}
}

type CreateTeamRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description"`
}

type UpdateTeamRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

type AddTeamMemberRequest struct {
	UserID uuid.UUID `json:"user_id" binding:"required"`
	Role   string    `json:"role"`
}

type UpdateTeamMemberRoleRequest struct {
	Role string `json:"role" binding:"required"`
}

type TeamWithMembers struct {
	*domain.Team
	Members []*domain.UserWithTeamRole `json:"members"`
}

func (s *TeamService) CreateTeam(ctx context.Context, orgID uuid.UUID, req *CreateTeamRequest) (*domain.Team, error) {
	team := &domain.Team{
		ID:             uuid.New(),
		OrganizationID: orgID,
		Name:           req.Name,
		Description:    req.Description,
	}

	if err := s.teamRepo.Create(ctx, team); err != nil {
		return nil, fmt.Errorf("failed to create team: %w", err)
	}

	return team, nil
}

func (s *TeamService) GetTeam(ctx context.Context, id uuid.UUID) (*domain.Team, error) {
	team, err := s.teamRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get team: %w", err)
	}

	return team, nil
}

func (s *TeamService) GetTeamWithMembers(ctx context.Context, id uuid.UUID) (*TeamWithMembers, error) {
	team, err := s.teamRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get team: %w", err)
	}

	members, err := s.teamRepo.ListMembers(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get team members: %w", err)
	}

	return &TeamWithMembers{
		Team:    team,
		Members: members,
	}, nil
}

func (s *TeamService) UpdateTeam(ctx context.Context, id uuid.UUID, req *UpdateTeamRequest) (*domain.Team, error) {
	team, err := s.teamRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get team: %w", err)
	}

	if req.Name != nil {
		team.Name = *req.Name
	}

	if req.Description != nil {
		team.Description = req.Description
	}

	if err := s.teamRepo.Update(ctx, team); err != nil {
		return nil, fmt.Errorf("failed to update team: %w", err)
	}

	return team, nil
}

func (s *TeamService) DeleteTeam(ctx context.Context, id uuid.UUID) error {
	if err := s.teamRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete team: %w", err)
	}

	return nil
}

func (s *TeamService) ListTeams(ctx context.Context, orgID uuid.UUID, page, pageSize int) ([]*domain.Team, error) {
	if page < 1 {
		page = 1
	}

	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize

	teams, err := s.teamRepo.List(ctx, orgID, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list teams: %w", err)
	}

	return teams, nil
}

func (s *TeamService) AddMember(ctx context.Context, teamID uuid.UUID, req *AddTeamMemberRequest) error {
	// Check if user exists
	_, err := s.userRepo.GetByID(ctx, req.UserID)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	// Parse role
	role := domain.TeamRoleMember
	if req.Role != "" {
		role = domain.TeamRole(req.Role)
	}

	if err := s.teamRepo.AddMember(ctx, teamID, req.UserID, role); err != nil {
		return fmt.Errorf("failed to add team member: %w", err)
	}

	return nil
}

func (s *TeamService) RemoveMember(ctx context.Context, teamID, userID uuid.UUID) error {
	if err := s.teamRepo.RemoveMember(ctx, teamID, userID); err != nil {
		return fmt.Errorf("failed to remove team member: %w", err)
	}

	return nil
}

func (s *TeamService) UpdateMemberRole(ctx context.Context, teamID, userID uuid.UUID, req *UpdateTeamMemberRoleRequest) error {
	role := domain.TeamRole(req.Role)

	if err := s.teamRepo.UpdateMemberRole(ctx, teamID, userID, role); err != nil {
		return fmt.Errorf("failed to update team member role: %w", err)
	}

	return nil
}

func (s *TeamService) ListMembers(ctx context.Context, teamID uuid.UUID) ([]*domain.UserWithTeamRole, error) {
	members, err := s.teamRepo.ListMembers(ctx, teamID)
	if err != nil {
		return nil, fmt.Errorf("failed to list team members: %w", err)
	}

	// Clear password hashes
	for _, member := range members {
		member.PasswordHash = ""
	}

	return members, nil
}

func (s *TeamService) ListUserTeams(ctx context.Context, userID uuid.UUID) ([]*domain.Team, error) {
	teams, err := s.teamRepo.ListUserTeams(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list user teams: %w", err)
	}

	return teams, nil
}
