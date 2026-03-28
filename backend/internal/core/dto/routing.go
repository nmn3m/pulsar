package dto

import (
	"encoding/json"

	"github.com/google/uuid"
)

type CreateRoutingRuleRequest struct {
	Name        string          `json:"name" binding:"required"`
	Description *string         `json:"description"`
	Priority    int             `json:"priority"`
	Conditions  json.RawMessage `json:"conditions" binding:"required"`
	Actions     json.RawMessage `json:"actions" binding:"required"`
	Enabled     *bool           `json:"enabled"`
}

type UpdateRoutingRuleRequest struct {
	Name        *string         `json:"name"`
	Description *string         `json:"description"`
	Priority    *int            `json:"priority"`
	Conditions  json.RawMessage `json:"conditions"`
	Actions     json.RawMessage `json:"actions"`
	Enabled     *bool           `json:"enabled"`
}

type ReorderRoutingRulesRequest struct {
	RuleIDs []uuid.UUID `json:"rule_ids" binding:"required"`
}
