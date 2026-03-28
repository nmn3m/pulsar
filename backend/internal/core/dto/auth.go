package dto

import (
	"github.com/nmn3m/pulsar/backend/internal/core/domain"
)

type RegisterRequest struct {
	Email            string `json:"email" binding:"required,email"`
	Username         string `json:"username" binding:"required,min=3,max=50"`
	Password         string `json:"password" binding:"required,min=8"`
	FullName         string `json:"full_name"`
	OrganizationName string `json:"organization_name" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	User                      *domain.User         `json:"user"`
	Organization              *domain.Organization `json:"organization"`
	AccessToken               string               `json:"access_token"`
	RefreshToken              string               `json:"refresh_token"`
	RequiresEmailVerification bool                 `json:"requires_email_verification,omitempty"`
}
