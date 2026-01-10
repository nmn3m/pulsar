package domain

import (
	"time"

	"github.com/google/uuid"
)

// InvitationStatus represents the status of an invitation
type InvitationStatus string

const (
	InvitationStatusPending  InvitationStatus = "pending"
	InvitationStatusAccepted InvitationStatus = "accepted"
	InvitationStatusDeclined InvitationStatus = "declined"
	InvitationStatusExpired  InvitationStatus = "expired"
)

// TeamInvitation represents an invitation to join a team
type TeamInvitation struct {
	ID             uuid.UUID        `json:"id" db:"id"`
	TeamID         uuid.UUID        `json:"team_id" db:"team_id"`
	OrganizationID uuid.UUID        `json:"organization_id" db:"organization_id"`
	Email          string           `json:"email" db:"email"`
	Role           TeamRole         `json:"role" db:"role"`
	Token          string           `json:"-" db:"token"` // Hidden from JSON
	Status         InvitationStatus `json:"status" db:"status"`
	InvitedByID    uuid.UUID        `json:"invited_by_id" db:"invited_by_id"`
	ExpiresAt      time.Time        `json:"expires_at" db:"expires_at"`
	CreatedAt      time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time        `json:"updated_at" db:"updated_at"`
}

// TeamInvitationWithDetails includes team and inviter info
type TeamInvitationWithDetails struct {
	*TeamInvitation
	TeamName    string `json:"team_name" db:"team_name"`
	InvitedBy   string `json:"invited_by" db:"invited_by"` // Inviter's email
	InviterName string `json:"inviter_name,omitempty" db:"inviter_name"`
}

// IsExpired checks if the invitation has expired
func (i *TeamInvitation) IsExpired() bool {
	return time.Now().After(i.ExpiresAt)
}

// IsValid checks if the invitation can still be accepted
func (i *TeamInvitation) IsValid() bool {
	return i.Status == InvitationStatusPending && !i.IsExpired()
}
