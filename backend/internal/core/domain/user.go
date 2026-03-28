package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                      uuid.UUID
	Email                   string
	Username                string
	PasswordHash            string
	FullName                *string
	Phone                   *string
	Timezone                string
	NotificationPreferences map[string]interface{}
	IsActive                bool
	EmailVerified           bool
	CreatedAt               time.Time
	UpdatedAt               time.Time
}

type OrganizationUser struct {
	OrganizationID uuid.UUID
	UserID         uuid.UUID
	Role           string
	JoinedAt       time.Time
}

// UserRole represents the user's role within an organization
type UserRole string

const (
	RoleOwner  UserRole = "owner"
	RoleAdmin  UserRole = "admin"
	RoleMember UserRole = "member"
	RoleViewer UserRole = "viewer"
)

func (r UserRole) String() string {
	return string(r)
}

func (r UserRole) IsValid() bool {
	switch r {
	case RoleOwner, RoleAdmin, RoleMember, RoleViewer:
		return true
	}
	return false
}

// UserWithOrganization represents a user with their organization role
type UserWithOrganization struct {
	User
	OrganizationID uuid.UUID
	Role           UserRole
}
