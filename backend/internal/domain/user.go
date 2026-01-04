package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                      uuid.UUID              `json:"id" db:"id"`
	Email                   string                 `json:"email" db:"email"`
	Username                string                 `json:"username" db:"username"`
	PasswordHash            string                 `json:"-" db:"password_hash"`
	FullName                *string                `json:"full_name,omitempty" db:"full_name"`
	Phone                   *string                `json:"phone,omitempty" db:"phone"`
	Timezone                string                 `json:"timezone" db:"timezone"`
	NotificationPreferences map[string]interface{} `json:"notification_preferences" db:"notification_preferences"`
	IsActive                bool                   `json:"is_active" db:"is_active"`
	CreatedAt               time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt               time.Time              `json:"updated_at" db:"updated_at"`
}

type OrganizationUser struct {
	OrganizationID uuid.UUID `json:"organization_id" db:"organization_id"`
	UserID         uuid.UUID `json:"user_id" db:"user_id"`
	Role           string    `json:"role" db:"role"`
	JoinedAt       time.Time `json:"joined_at" db:"joined_at"`
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
	OrganizationID uuid.UUID `json:"organization_id"`
	Role           UserRole  `json:"role"`
}
