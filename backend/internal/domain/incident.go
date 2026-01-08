package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// IncidentSeverity represents the severity level of an incident
type IncidentSeverity string

const (
	IncidentSeverityCritical IncidentSeverity = "critical"
	IncidentSeverityHigh     IncidentSeverity = "high"
	IncidentSeverityMedium   IncidentSeverity = "medium"
	IncidentSeverityLow      IncidentSeverity = "low"
)

// IsValid checks if the severity is valid
func (s IncidentSeverity) IsValid() bool {
	switch s {
	case IncidentSeverityCritical, IncidentSeverityHigh, IncidentSeverityMedium, IncidentSeverityLow:
		return true
	}
	return false
}

func (s IncidentSeverity) String() string {
	return string(s)
}

// IncidentStatus represents the current status of an incident
type IncidentStatus string

const (
	IncidentStatusInvestigating IncidentStatus = "investigating"
	IncidentStatusIdentified    IncidentStatus = "identified"
	IncidentStatusMonitoring    IncidentStatus = "monitoring"
	IncidentStatusResolved      IncidentStatus = "resolved"
)

// IsValid checks if the status is valid
func (s IncidentStatus) IsValid() bool {
	switch s {
	case IncidentStatusInvestigating, IncidentStatusIdentified, IncidentStatusMonitoring, IncidentStatusResolved:
		return true
	}
	return false
}

func (s IncidentStatus) String() string {
	return string(s)
}

// Incident represents an incident
type Incident struct {
	ID               uuid.UUID        `json:"id" db:"id"`
	OrganizationID   uuid.UUID        `json:"organization_id" db:"organization_id"`
	Title            string           `json:"title" db:"title"`
	Description      *string          `json:"description" db:"description"`
	Severity         IncidentSeverity `json:"severity" db:"severity"`
	Status           IncidentStatus   `json:"status" db:"status"`
	Priority         AlertPriority    `json:"priority" db:"priority"`
	CreatedByUserID  uuid.UUID        `json:"created_by_user_id" db:"created_by_user_id"`
	AssignedToTeamID *uuid.UUID       `json:"assigned_to_team_id" db:"assigned_to_team_id"`
	StartedAt        time.Time        `json:"started_at" db:"started_at"`
	ResolvedAt       *time.Time       `json:"resolved_at" db:"resolved_at"`
	CreatedAt        time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time        `json:"updated_at" db:"updated_at"`
}

// ResponderRole represents the role of an incident responder
type ResponderRole string

const (
	ResponderRoleIncidentCommander ResponderRole = "incident_commander"
	ResponderRoleResponder         ResponderRole = "responder"
)

// IsValid checks if the responder role is valid
func (r ResponderRole) IsValid() bool {
	switch r {
	case ResponderRoleIncidentCommander, ResponderRoleResponder:
		return true
	}
	return false
}

func (r ResponderRole) String() string {
	return string(r)
}

// IncidentResponder represents a user assigned to an incident
type IncidentResponder struct {
	ID         uuid.UUID     `json:"id" db:"id"`
	IncidentID uuid.UUID     `json:"incident_id" db:"incident_id"`
	UserID     uuid.UUID     `json:"user_id" db:"user_id"`
	Role       ResponderRole `json:"role" db:"role"`
	AddedAt    time.Time     `json:"added_at" db:"added_at"`
}

// ResponderWithUser extends IncidentResponder with user details
type ResponderWithUser struct {
	IncidentResponder
	User *User `json:"user"`
}

// TimelineEventType represents the type of timeline event
type TimelineEventType string

const (
	TimelineEventCreated          TimelineEventType = "created"
	TimelineEventStatusChanged    TimelineEventType = "status_changed"
	TimelineEventSeverityChanged  TimelineEventType = "severity_changed"
	TimelineEventResponderAdded   TimelineEventType = "responder_added"
	TimelineEventResponderRemoved TimelineEventType = "responder_removed"
	TimelineEventNoteAdded        TimelineEventType = "note_added"
	TimelineEventAlertLinked      TimelineEventType = "alert_linked"
	TimelineEventAlertUnlinked    TimelineEventType = "alert_unlinked"
	TimelineEventResolved         TimelineEventType = "resolved"
)

// IsValid checks if the timeline event type is valid
func (t TimelineEventType) IsValid() bool {
	switch t {
	case TimelineEventCreated, TimelineEventStatusChanged, TimelineEventSeverityChanged,
		TimelineEventResponderAdded, TimelineEventResponderRemoved, TimelineEventNoteAdded,
		TimelineEventAlertLinked, TimelineEventAlertUnlinked, TimelineEventResolved:
		return true
	}
	return false
}

func (t TimelineEventType) String() string {
	return string(t)
}

// IncidentTimelineEvent represents an event in the incident timeline
type IncidentTimelineEvent struct {
	ID          uuid.UUID              `json:"id" db:"id"`
	IncidentID  uuid.UUID              `json:"incident_id" db:"incident_id"`
	EventType   TimelineEventType      `json:"event_type" db:"event_type"`
	UserID      *uuid.UUID             `json:"user_id" db:"user_id"`
	Description string                 `json:"description" db:"description"`
	Metadata    map[string]interface{} `json:"metadata" db:"metadata"`
	CreatedAt   time.Time              `json:"created_at" db:"created_at"`
}

// TimelineEventWithUser extends IncidentTimelineEvent with user details
type TimelineEventWithUser struct {
	IncidentTimelineEvent
	User *User `json:"user,omitempty"`
}

// IncidentAlert represents a link between an incident and an alert
type IncidentAlert struct {
	ID             uuid.UUID  `json:"id" db:"id"`
	IncidentID     uuid.UUID  `json:"incident_id" db:"incident_id"`
	AlertID        uuid.UUID  `json:"alert_id" db:"alert_id"`
	LinkedAt       time.Time  `json:"linked_at" db:"linked_at"`
	LinkedByUserID *uuid.UUID `json:"linked_by_user_id" db:"linked_by_user_id"`
}

// IncidentAlertWithDetails extends IncidentAlert with alert details
type IncidentAlertWithDetails struct {
	IncidentAlert
	Alert *Alert `json:"alert"`
}

// IncidentFilter represents filters for listing incidents
type IncidentFilter struct {
	OrganizationID   uuid.UUID
	Status           []IncidentStatus
	Severity         []IncidentSeverity
	AssignedToTeamID *uuid.UUID
	Search           *string
	Limit            int
	Offset           int
}

// Validate validates the incident filter
func (f *IncidentFilter) Validate() error {
	if f.Limit < 0 {
		return fmt.Errorf("limit must be non-negative")
	}
	if f.Offset < 0 {
		return fmt.Errorf("offset must be non-negative")
	}
	return nil
}

// IncidentWithDetails extends Incident with related data
type IncidentWithDetails struct {
	Incident
	Responders []*ResponderWithUser        `json:"responders,omitempty"`
	Alerts     []*IncidentAlertWithDetails `json:"alerts,omitempty"`
	Timeline   []*TimelineEventWithUser    `json:"timeline,omitempty"`
}
