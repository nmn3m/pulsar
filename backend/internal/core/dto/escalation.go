package dto

import (
	"encoding/json"

	"github.com/google/uuid"
)

type CreateEscalationPolicyRequest struct {
	Name          string  `json:"name" binding:"required"`
	Description   *string `json:"description"`
	RepeatEnabled bool    `json:"repeat_enabled"`
	RepeatCount   *int    `json:"repeat_count"`
}

type UpdateEscalationPolicyRequest struct {
	Name          *string `json:"name"`
	Description   *string `json:"description"`
	RepeatEnabled *bool   `json:"repeat_enabled"`
	RepeatCount   *int    `json:"repeat_count"`
}

type CreateEscalationRuleRequest struct {
	Position        int `json:"position" binding:"required"`
	EscalationDelay int `json:"escalation_delay" binding:"required"`
}

type UpdateEscalationRuleRequest struct {
	Position        *int `json:"position"`
	EscalationDelay *int `json:"escalation_delay"`
}

type AddEscalationTargetRequest struct {
	TargetType           string          `json:"target_type" binding:"required"`
	TargetID             uuid.UUID       `json:"target_id" binding:"required"`
	NotificationChannels json.RawMessage `json:"notification_channels,omitempty"` // Optional channel override
}
