package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// AlertRoutingRule defines how alerts should be automatically routed
type AlertRoutingRule struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	Name           string
	Description    *string
	Priority       int // Lower number = higher priority
	Conditions     json.RawMessage
	Actions        json.RawMessage
	Enabled        bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
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
