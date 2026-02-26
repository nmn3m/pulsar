package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type EscalationPolicy struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	Name           string
	Description    *string
	RepeatEnabled  bool
	RepeatCount    *int // NULL = infinite
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type EscalationRule struct {
	ID              uuid.UUID
	PolicyID        uuid.UUID
	Position        int
	EscalationDelay int // minutes
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type EscalationTarget struct {
	ID                   uuid.UUID
	RuleID               uuid.UUID
	TargetType           EscalationTargetType
	TargetID             uuid.UUID
	NotificationChannels json.RawMessage
	CreatedAt            time.Time
}

// TargetNotificationConfig represents the notification channel override for a target
type TargetNotificationConfig struct {
	Channels []string `json:"channels"` // e.g., ["email", "slack", "sms", "webhook"]
	Urgent   bool     `json:"urgent"`   // If true, use urgent/high-priority notification
}

// ParseNotificationChannels parses the notification channels config
func (t *EscalationTarget) ParseNotificationChannels() (*TargetNotificationConfig, error) {
	if len(t.NotificationChannels) == 0 {
		return nil, nil // No override configured
	}
	var config TargetNotificationConfig
	if err := json.Unmarshal(t.NotificationChannels, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

// HasNotificationOverride returns true if this target has custom notification channels
func (t *EscalationTarget) HasNotificationOverride() bool {
	return len(t.NotificationChannels) > 0
}

type AlertEscalationEvent struct {
	ID               uuid.UUID
	AlertID          uuid.UUID
	PolicyID         uuid.UUID
	RuleID           *uuid.UUID
	EventType        EscalationEventType
	CurrentLevel     int
	RepeatCount      int
	NextEscalationAt *time.Time
	CreatedAt        time.Time
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
	Rules []*EscalationRuleWithTargets
}

// EscalationRuleWithTargets includes the rule and its targets
type EscalationRuleWithTargets struct {
	EscalationRule
	Targets []*EscalationTarget
}
