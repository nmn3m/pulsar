package domain

import (
	"time"

	"github.com/google/uuid"
)

type Alert struct {
	ID             uuid.UUID              `json:"id" db:"id"`
	OrganizationID uuid.UUID              `json:"organization_id" db:"organization_id"`
	Source         string                 `json:"source" db:"source"`
	SourceID       *string                `json:"source_id,omitempty" db:"source_id"`
	Priority       AlertPriority          `json:"priority" db:"priority"`
	Status         AlertStatus            `json:"status" db:"status"`
	Message        string                 `json:"message" db:"message"`
	Description    *string                `json:"description,omitempty" db:"description"`
	Tags           []string               `json:"tags" db:"tags"`
	CustomFields   map[string]interface{} `json:"custom_fields" db:"custom_fields"`

	// Assignment
	AssignedToUserID *uuid.UUID `json:"assigned_to_user_id,omitempty" db:"assigned_to_user_id"`
	AssignedToTeamID *uuid.UUID `json:"assigned_to_team_id,omitempty" db:"assigned_to_team_id"`

	// Acknowledgment
	AcknowledgedBy *uuid.UUID `json:"acknowledged_by,omitempty" db:"acknowledged_by"`
	AcknowledgedAt *time.Time `json:"acknowledged_at,omitempty" db:"acknowledged_at"`

	// Closure
	ClosedBy    *uuid.UUID `json:"closed_by,omitempty" db:"closed_by"`
	ClosedAt    *time.Time `json:"closed_at,omitempty" db:"closed_at"`
	CloseReason *string    `json:"close_reason,omitempty" db:"close_reason"`

	// Snooze
	SnoozedUntil *time.Time `json:"snoozed_until,omitempty" db:"snoozed_until"`

	// Escalation
	EscalationPolicyID *uuid.UUID `json:"escalation_policy_id,omitempty" db:"escalation_policy_id"`
	EscalationLevel    int        `json:"escalation_level" db:"escalation_level"`
	LastEscalatedAt    *time.Time `json:"last_escalated_at,omitempty" db:"last_escalated_at"`

	// Deduplication
	DedupKey          *string    `json:"dedup_key,omitempty" db:"dedup_key"`
	DedupCount        int        `json:"dedup_count" db:"dedup_count"`
	FirstOccurrenceAt *time.Time `json:"first_occurrence_at,omitempty" db:"first_occurrence_at"`
	LastOccurrenceAt  *time.Time `json:"last_occurrence_at,omitempty" db:"last_occurrence_at"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type AlertPriority string

const (
	PriorityP1 AlertPriority = "P1" // Critical
	PriorityP2 AlertPriority = "P2" // High
	PriorityP3 AlertPriority = "P3" // Medium
	PriorityP4 AlertPriority = "P4" // Low
	PriorityP5 AlertPriority = "P5" // Informational
)

func (p AlertPriority) String() string {
	return string(p)
}

func (p AlertPriority) IsValid() bool {
	switch p {
	case PriorityP1, PriorityP2, PriorityP3, PriorityP4, PriorityP5:
		return true
	}
	return false
}

type AlertStatus string

const (
	AlertStatusOpen         AlertStatus = "open"
	AlertStatusAcknowledged AlertStatus = "acknowledged"
	AlertStatusClosed       AlertStatus = "closed"
	AlertStatusSnoozed      AlertStatus = "snoozed"
)

func (s AlertStatus) String() string {
	return string(s)
}

func (s AlertStatus) IsValid() bool {
	switch s {
	case AlertStatusOpen, AlertStatusAcknowledged, AlertStatusClosed, AlertStatusSnoozed:
		return true
	}
	return false
}

type AlertSource string

const (
	AlertSourceWebhook     AlertSource = "webhook"
	AlertSourceAPI         AlertSource = "api"
	AlertSourceEmail       AlertSource = "email"
	AlertSourceIntegration AlertSource = "integration"
	AlertSourceManual      AlertSource = "manual"
)

func (s AlertSource) String() string {
	return string(s)
}

// AlertWithAssignees includes assigned user and team information
type AlertWithAssignees struct {
	Alert
	AssignedUser *User `json:"assigned_user,omitempty"`
	AssignedTeam *Team `json:"assigned_team,omitempty"`
}

// AlertFilter for filtering and pagination
type AlertFilter struct {
	OrganizationID uuid.UUID
	Status         []AlertStatus
	Priority       []AlertPriority
	AssignedToUser *uuid.UUID
	AssignedToTeam *uuid.UUID
	Source         *string
	Search         *string // Search in message and description
	Limit          int
	Offset         int
}
