package dto

import (
	"encoding/json"
	"time"
)

// UpdateDNDSettingsRequest represents a request to update DND settings
type UpdateDNDSettingsRequest struct {
	Enabled         *bool           `json:"enabled"`
	Schedule        json.RawMessage `json:"schedule"`
	Overrides       json.RawMessage `json:"overrides"`
	AllowP1Override *bool           `json:"allow_p1_override"`
}

// AddDNDOverrideRequest represents a request to add a DND override
type AddDNDOverrideRequest struct {
	Start  time.Time `json:"start" binding:"required"`
	End    time.Time `json:"end" binding:"required"`
	Reason string    `json:"reason"`
}
