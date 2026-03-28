package domain

import (
	"time"

	"github.com/google/uuid"
)

type Alert struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	Source         string
	SourceID       *string
	Priority       AlertPriority
	Status         AlertStatus
	Message        string
	Description    *string
	Tags           []string
	CustomFields   map[string]interface{}

	// Assignment
	AssignedToUserID *uuid.UUID
	AssignedToTeamID *uuid.UUID

	// Acknowledgment
	AcknowledgedBy *uuid.UUID
	AcknowledgedAt *time.Time

	// Closure
	ClosedBy    *uuid.UUID
	ClosedAt    *time.Time
	CloseReason *string

	// Snooze
	SnoozedUntil *time.Time

	// Escalation
	EscalationPolicyID *uuid.UUID
	EscalationLevel    int
	LastEscalatedAt    *time.Time

	// Deduplication
	DedupKey          *string
	DedupCount        int
	FirstOccurrenceAt *time.Time
	LastOccurrenceAt  *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
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
	AssignedUser *User
	AssignedTeam *Team
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
