package dto

import (
	"github.com/nmn3m/pulsar/backend/internal/core/domain"
)

type CreateAPIKeyRequest struct {
	Name      string   `json:"name" binding:"required,min=1,max=255"`
	Scopes    []string `json:"scopes" binding:"required,min=1"`
	ExpiresAt *string  `json:"expires_at,omitempty"`
}

type UpdateAPIKeyRequest struct {
	Name     *string  `json:"name,omitempty"`
	Scopes   []string `json:"scopes,omitempty"`
	IsActive *bool    `json:"is_active,omitempty"`
}

type APIKeyResponse struct {
	*domain.APIKey
	RawKey string `json:"key,omitempty"`
}
