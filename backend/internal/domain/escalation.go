package domain

import (
	"time"

	"github.com/google/uuid"
)

type EscalationPolicy struct {
	ID             uuid.UUID `json:"id" db:"id"`
	OrganizationID uuid.UUID `json:"organization_id" db:"organization_id"`
	Name           string    `json:"name" db:"name"`
	Description    *string   `json:"description,omitempty" db:"description"`
	RepeatEnabled  bool      `json:"repeat_enabled" db:"repeat_enabled"`
	RepeatCount    *int      `json:"repeat_count,omitempty" db:"repeat_count"` // NULL = infinite
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

type EscalationRule struct {
	ID              uuid.UUID `json:"id" db:"id"`
	PolicyID        uuid.UUID `json:"policy_id" db:"policy_id"`
	Position        int       `json:"position" db:"position"`
	EscalationDelay int       `json:"escalation_delay" db:"escalation_delay"` // minutes
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

type EscalationTarget struct {
	ID         uuid.UUID            `json:"id" db:"id"`
	RuleID     uuid.UUID            `json:"rule_id" db:"rule_id"`
	TargetType EscalationTargetType `json:"target_type" db:"target_type"`
	TargetID   uuid.UUID            `json:"target_id" db:"target_id"`
	CreatedAt  time.Time            `json:"created_at" db:"created_at"`
}

type AlertEscalationEvent struct {
	ID               uuid.UUID           `json:"id" db:"id"`
	AlertID          uuid.UUID           `json:"alert_id" db:"alert_id"`
	PolicyID         uuid.UUID           `json:"policy_id" db:"policy_id"`
	RuleID           *uuid.UUID          `json:"rule_id,omitempty" db:"rule_id"`
	EventType        EscalationEventType `json:"event_type" db:"event_type"`
	CurrentLevel     int                 `json:"current_level" db:"current_level"`
	RepeatCount      int                 `json:"repeat_count" db:"repeat_count"`
	NextEscalationAt *time.Time          `json:"next_escalation_at,omitempty" db:"next_escalation_at"`
	CreatedAt        time.Time           `json:"created_at" db:"created_at"`
}

type EscalationTargetType string

const (
	EscalationTargetTypeUser     EscalationTargetType = "user"
	EscalationTargetTypeTeam     EscalationTargetType = "team"
	EscalationTargetTypeSchedule EscalationTargetType = "schedule"
)

func (t EscalationTargetType) String() string {
	return string(t)
}

func (t EscalationTargetType) Validate() error {
	switch t {
	case EscalationTargetTypeUser, EscalationTargetTypeTeam, EscalationTargetTypeSchedule:
		return nil
	default:
		return ErrInvalidEscalationTarget
	}
}

type EscalationEventType string

const (
	EscalationEventTriggered    EscalationEventType = "triggered"
	EscalationEventAcknowledged EscalationEventType = "acknowledged"
	EscalationEventCompleted    EscalationEventType = "completed"
	EscalationEventStopped      EscalationEventType = "stopped"
)

func (t EscalationEventType) String() string {
	return string(t)
}

// EscalationPolicyWithRules includes the policy and its rules with targets
type EscalationPolicyWithRules struct {
	EscalationPolicy
	Rules []*EscalationRuleWithTargets `json:"rules"`
}

// EscalationRuleWithTargets includes the rule and its targets
type EscalationRuleWithTargets struct {
	EscalationRule
	Targets []*EscalationTarget `json:"targets"`
}
