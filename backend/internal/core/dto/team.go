package dto

import (
	"github.com/google/uuid"

	"github.com/nmn3m/pulsar/backend/internal/core/domain"
)

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
