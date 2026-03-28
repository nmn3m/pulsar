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
	ID               uuid.UUID
	OrganizationID   uuid.UUID
	Title            string
	Description      *string
	Severity         IncidentSeverity
	Status           IncidentStatus
	Priority         AlertPriority
	CreatedByUserID  uuid.UUID
	AssignedToTeamID *uuid.UUID
	StartedAt        time.Time
	ResolvedAt       *time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
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
	ID         uuid.UUID
	IncidentID uuid.UUID
	UserID     uuid.UUID
	Role       ResponderRole
	AddedAt    time.Time
}

// ResponderWithUser extends IncidentResponder with user details
type ResponderWithUser struct {
	IncidentResponder
	User *User
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
	ID          uuid.UUID
	IncidentID  uuid.UUID
	EventType   TimelineEventType
	UserID      *uuid.UUID
	Description string
	Metadata    map[string]interface{}
	CreatedAt   time.Time
}

// TimelineEventWithUser extends IncidentTimelineEvent with user details
type TimelineEventWithUser struct {
	IncidentTimelineEvent
	User *User
}

// IncidentAlert represents a link between an incident and an alert
type IncidentAlert struct {
	ID             uuid.UUID
	IncidentID     uuid.UUID
	AlertID        uuid.UUID
	LinkedAt       time.Time
	LinkedByUserID *uuid.UUID
}

// IncidentAlertWithDetails extends IncidentAlert with alert details
type IncidentAlertWithDetails struct {
	IncidentAlert
	Alert *Alert
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
	Responders []*ResponderWithUser
	Alerts     []*IncidentAlertWithDetails
	Timeline   []*TimelineEventWithUser
}
