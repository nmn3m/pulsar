package domain

import (
	"time"

	"github.com/google/uuid"
)

type Team struct {
	ID             uuid.UUID `json:"id" db:"id"`
	OrganizationID uuid.UUID `json:"organization_id" db:"organization_id"`
	Name           string    `json:"name" db:"name"`
	Description    *string   `json:"description,omitempty" db:"description"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

type TeamMember struct {
	TeamID   uuid.UUID `json:"team_id" db:"team_id"`
	UserID   uuid.UUID `json:"user_id" db:"user_id"`
	Role     string    `json:"role" db:"role"`
	JoinedAt time.Time `json:"joined_at" db:"joined_at"`
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
	Role     string    `json:"role"`
	JoinedAt time.Time `json:"joined_at"`
}
