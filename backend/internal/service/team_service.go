package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/nmn3m/pulsar/backend/internal/domain"
	"github.com/nmn3m/pulsar/backend/internal/repository"
)

type TeamService struct {
	teamRepo       repository.TeamRepository
	userRepo       repository.UserRepository
	invitationRepo repository.TeamInvitationRepository
	emailService   EmailServiceInterface
}

// EmailServiceInterface defines the interface for sending emails
type EmailServiceInterface interface {
	SendTeamInvitation(ctx context.Context, toEmail, teamName, inviterName, inviteToken string) error
}

func NewTeamService(teamRepo repository.TeamRepository, userRepo repository.UserRepository) *TeamService {
	return &TeamService{
		teamRepo: teamRepo,
		userRepo: userRepo,
	}
}

// SetInvitationRepo sets the invitation repository (optional dependency)
func (s *TeamService) SetInvitationRepo(repo repository.TeamInvitationRepository) {
	s.invitationRepo = repo
}

// SetEmailService sets the email service (optional dependency)
func (s *TeamService) SetEmailService(emailSvc EmailServiceInterface) {
	s.emailService = emailSvc
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
	UserID *uuid.UUID `json:"user_id"` // Optional: if provided, add existing user
	Email  string     `json:"email"`   // Optional: if provided without user_id, find or invite
	Role   string     `json:"role"`
}

type InviteMemberRequest struct {
	Email string `json:"email" binding:"required,email"`
	Role  string `json:"role"`
}

type InvitationResponse struct {
	UserAdded  bool                   `json:"user_added"` // True if user was directly added
	Invited    bool                   `json:"invited"`    // True if invitation was sent
	Invitation *domain.TeamInvitation `json:"invitation,omitempty"`
	Message    string                 `json:"message"`
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
	// Parse role
	role := domain.TeamRoleMember
	if req.Role != "" {
		role = domain.TeamRole(req.Role)
	}

	// If user_id is provided, add directly
	if req.UserID != nil {
		_, err := s.userRepo.GetByID(ctx, *req.UserID)
		if err != nil {
			return fmt.Errorf("user not found")
		}
		if err := s.teamRepo.AddMember(ctx, teamID, *req.UserID, role); err != nil {
			return fmt.Errorf("failed to add team member: %w", err)
		}
		return nil
	}

	// If email is provided, try to find the user
	if req.Email != "" {
		user, err := s.userRepo.GetByEmail(ctx, req.Email)
		if err == nil && user != nil {
			// User exists, add them directly
			if err := s.teamRepo.AddMember(ctx, teamID, user.ID, role); err != nil {
				return fmt.Errorf("failed to add team member: %w", err)
			}
			return nil
		}
		// User not found - return error suggesting to use invite endpoint
		return fmt.Errorf("user not found with email %s, use invite endpoint to send invitation", req.Email)
	}

	return fmt.Errorf("either user_id or email is required")
}

// AddMemberOrInvite adds an existing user or sends an invitation
func (s *TeamService) AddMemberOrInvite(ctx context.Context, teamID, orgID, inviterID uuid.UUID, req *InviteMemberRequest) (*InvitationResponse, error) {
	// Parse role
	role := domain.TeamRoleMember
	if req.Role != "" {
		role = domain.TeamRole(req.Role)
	}

	// First check if user already exists
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err == nil && user != nil {
		// User exists, add them directly
		if err := s.teamRepo.AddMember(ctx, teamID, user.ID, role); err != nil {
			return nil, fmt.Errorf("failed to add team member: %w", err)
		}
		return &InvitationResponse{
			UserAdded: true,
			Invited:   false,
			Message:   "User added to team successfully",
		}, nil
	}

	// User doesn't exist, create an invitation
	if s.invitationRepo == nil {
		return nil, fmt.Errorf("invitation feature not configured")
	}

	// Check for existing pending invitation
	existingInvite, _ := s.invitationRepo.GetByEmailAndTeam(ctx, req.Email, teamID)
	if existingInvite != nil && existingInvite.IsValid() {
		return &InvitationResponse{
			UserAdded:  false,
			Invited:    true,
			Invitation: existingInvite,
			Message:    "Invitation already sent to this email",
		}, nil
	}

	// Generate invitation token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return nil, fmt.Errorf("failed to generate invitation token: %w", err)
	}
	token := hex.EncodeToString(tokenBytes)

	// Create invitation
	invitation := &domain.TeamInvitation{
		ID:             uuid.New(),
		TeamID:         teamID,
		OrganizationID: orgID,
		Email:          req.Email,
		Role:           role,
		Token:          token,
		Status:         domain.InvitationStatusPending,
		InvitedByID:    inviterID,
		ExpiresAt:      time.Now().Add(7 * 24 * time.Hour), // 7 days expiry
	}

	if err := s.invitationRepo.Create(ctx, invitation); err != nil {
		return nil, fmt.Errorf("failed to create invitation: %w", err)
	}

	// Send invitation email
	if s.emailService != nil {
		team, _ := s.teamRepo.GetByID(ctx, teamID)
		inviter, _ := s.userRepo.GetByID(ctx, inviterID)
		teamName := "the team"
		inviterName := "A team member"
		if team != nil {
			teamName = team.Name
		}
		if inviter != nil {
			inviterName = inviter.Username
		}
		// Send email asynchronously (don't block on email failure)
		go func() {
			_ = s.emailService.SendTeamInvitation(context.Background(), req.Email, teamName, inviterName, token)
		}()
	}

	return &InvitationResponse{
		UserAdded:  false,
		Invited:    true,
		Invitation: invitation,
		Message:    "Invitation sent successfully",
	}, nil
}

// AcceptInvitation accepts a team invitation
func (s *TeamService) AcceptInvitation(ctx context.Context, token string, userID uuid.UUID) error {
	if s.invitationRepo == nil {
		return fmt.Errorf("invitation feature not configured")
	}

	invitation, err := s.invitationRepo.GetByToken(ctx, token)
	if err != nil {
		return fmt.Errorf("invitation not found")
	}

	if !invitation.IsValid() {
		if invitation.IsExpired() {
			return fmt.Errorf("invitation has expired")
		}
		return fmt.Errorf("invitation is no longer valid")
	}

	// Add user to team
	if err := s.teamRepo.AddMember(ctx, invitation.TeamID, userID, invitation.Role); err != nil {
		return fmt.Errorf("failed to add team member: %w", err)
	}

	// Update invitation status
	invitation.Status = domain.InvitationStatusAccepted
	if err := s.invitationRepo.Update(ctx, invitation.TeamInvitation); err != nil {
		return fmt.Errorf("failed to update invitation: %w", err)
	}

	return nil
}

// DeclineInvitation declines a team invitation
func (s *TeamService) DeclineInvitation(ctx context.Context, token string) error {
	if s.invitationRepo == nil {
		return fmt.Errorf("invitation feature not configured")
	}

	invitation, err := s.invitationRepo.GetByToken(ctx, token)
	if err != nil {
		return fmt.Errorf("invitation not found")
	}

	invitation.Status = domain.InvitationStatusDeclined
	if err := s.invitationRepo.Update(ctx, invitation.TeamInvitation); err != nil {
		return fmt.Errorf("failed to update invitation: %w", err)
	}

	return nil
}

// GetPendingInvitations returns pending invitations for a user's email
func (s *TeamService) GetPendingInvitations(ctx context.Context, email string) ([]*domain.TeamInvitationWithDetails, error) {
	if s.invitationRepo == nil {
		return nil, fmt.Errorf("invitation feature not configured")
	}
	return s.invitationRepo.ListByEmail(ctx, email)
}

// ListTeamInvitations returns all invitations for a team
func (s *TeamService) ListTeamInvitations(ctx context.Context, teamID uuid.UUID) ([]*domain.TeamInvitation, error) {
	if s.invitationRepo == nil {
		return nil, fmt.Errorf("invitation feature not configured")
	}
	return s.invitationRepo.ListByTeam(ctx, teamID)
}

// CancelInvitation cancels a pending invitation
func (s *TeamService) CancelInvitation(ctx context.Context, invitationID uuid.UUID) error {
	if s.invitationRepo == nil {
		return fmt.Errorf("invitation feature not configured")
	}
	return s.invitationRepo.Delete(ctx, invitationID)
}

// ResendInvitation resends an invitation email
func (s *TeamService) ResendInvitation(ctx context.Context, invitationID uuid.UUID) error {
	if s.invitationRepo == nil {
		return fmt.Errorf("invitation feature not configured")
	}

	invitation, err := s.invitationRepo.GetByID(ctx, invitationID)
	if err != nil {
		return fmt.Errorf("invitation not found")
	}

	if invitation.Status != domain.InvitationStatusPending {
		return fmt.Errorf("can only resend pending invitations")
	}

	// Send invitation email
	if s.emailService != nil {
		team, _ := s.teamRepo.GetByID(ctx, invitation.TeamID)
		inviter, _ := s.userRepo.GetByID(ctx, invitation.InvitedByID)
		teamName := "the team"
		inviterName := "A team member"
		if team != nil {
			teamName = team.Name
		}
		if inviter != nil {
			inviterName = inviter.Username
		}
		if err := s.emailService.SendTeamInvitation(ctx, invitation.Email, teamName, inviterName, invitation.Token); err != nil {
			return fmt.Errorf("failed to send invitation email: %w", err)
		}
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
