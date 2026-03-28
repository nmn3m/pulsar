package domain

import (
	"time"

	"github.com/google/uuid"
)

type Team struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	Name           string
	Description    *string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type TeamMember struct {
	TeamID   uuid.UUID
	UserID   uuid.UUID
	Role     string
	JoinedAt time.Time
}

type TeamRole string

const (
	TeamRoleLead   TeamRole = "lead"
	TeamRoleMember TeamRole = "member"
)

func (r TeamRole) String() string {
	return string(r)
}

type UserWithTeamRole struct {
	User
	Role     string
	JoinedAt time.Time
}
