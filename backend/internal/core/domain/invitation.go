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
	ID             uuid.UUID
	TeamID         uuid.UUID
	OrganizationID uuid.UUID
	Email          string
	Role           TeamRole
	Token          string // Hidden from JSON
	Status         InvitationStatus
	InvitedByID    uuid.UUID
	ExpiresAt      time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// TeamInvitationWithDetails includes team and inviter info
type TeamInvitationWithDetails struct {
	*TeamInvitation
	TeamName    string
	InvitedBy   string // Inviter's email
	InviterName string
}

// IsExpired checks if the invitation has expired
func (i *TeamInvitation) IsExpired() bool {
	return time.Now().After(i.ExpiresAt)
}

// IsValid checks if the invitation can still be accepted
func (i *TeamInvitation) IsValid() bool {
	return i.Status == InvitationStatusPending && !i.IsExpired()
}
