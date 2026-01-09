package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// AlertRoutingRule defines how alerts should be automatically routed
type AlertRoutingRule struct {
	ID             uuid.UUID       `json:"id" db:"id"`
	OrganizationID uuid.UUID       `json:"organization_id" db:"organization_id"`
	Name           string          `json:"name" db:"name"`
	Description    *string         `json:"description,omitempty" db:"description"`
	Priority       int             `json:"priority" db:"priority"` // Lower number = higher priority
	Conditions     json.RawMessage `json:"conditions" db:"conditions"`
	Actions        json.RawMessage `json:"actions" db:"actions"`
	Enabled        bool            `json:"enabled" db:"enabled"`
	CreatedAt      time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at" db:"updated_at"`
}

// RoutingConditions represents the condition configuration for a routing rule
type RoutingConditions struct {
	Match      string             `json:"match"` // "all" or "any"
	Conditions []RoutingCondition `json:"conditions"`
}

// RoutingCondition represents a single condition to evaluate
type RoutingCondition struct {
	Field    string `json:"field"`    // source, priority, tags, message
	Operator string `json:"operator"` // equals, not_equals, contains, not_contains, regex, gte, lte
	Value    string `json:"value"`
}

// RoutingActions represents the actions to take when a rule matches
type RoutingActions struct {
	AssignTeamID             *uuid.UUID `json:"assign_team_id,omitempty"`
	AssignUserID             *uuid.UUID `json:"assign_user_id,omitempty"`
	AssignEscalationPolicyID *uuid.UUID `json:"assign_escalation_policy_id,omitempty"`
	SetPriority              *string    `json:"set_priority,omitempty"`
	AddTags                  []string   `json:"add_tags,omitempty"`
	Suppress                 bool       `json:"suppress"`
}

// CreateRoutingRuleRequest is the request to create a routing rule
type CreateRoutingRuleRequest struct {
	Name        string          `json:"name" binding:"required"`
	Description *string         `json:"description,omitempty"`
	Priority    int             `json:"priority"`
	Conditions  json.RawMessage `json:"conditions" binding:"required"`
	Actions     json.RawMessage `json:"actions" binding:"required"`
	Enabled     *bool           `json:"enabled,omitempty"`
}

// UpdateRoutingRuleRequest is the request to update a routing rule
type UpdateRoutingRuleRequest struct {
	Name        *string         `json:"name,omitempty"`
	Description *string         `json:"description,omitempty"`
	Priority    *int            `json:"priority,omitempty"`
	Conditions  json.RawMessage `json:"conditions,omitempty"`
	Actions     json.RawMessage `json:"actions,omitempty"`
	Enabled     *bool           `json:"enabled,omitempty"`
}

// ReorderRoutingRulesRequest is the request to reorder routing rules
type ReorderRoutingRulesRequest struct {
	RuleIDs []uuid.UUID `json:"rule_ids" binding:"required"`
}

// ParseConditions parses the raw JSON conditions into a structured format
func (r *AlertRoutingRule) ParseConditions() (*RoutingConditions, error) {
	var conditions RoutingConditions
	if err := json.Unmarshal(r.Conditions, &conditions); err != nil {
		return nil, err
	}
	return &conditions, nil
}

// ParseActions parses the raw JSON actions into a structured format
func (r *AlertRoutingRule) ParseActions() (*RoutingActions, error) {
	var actions RoutingActions
	if err := json.Unmarshal(r.Actions, &actions); err != nil {
		return nil, err
	}
	return &actions, nil
}
