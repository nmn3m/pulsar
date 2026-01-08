package domain

import (
	"time"

	"github.com/google/uuid"
)

type Schedule struct {
	ID             uuid.UUID  `json:"id" db:"id"`
	OrganizationID uuid.UUID  `json:"organization_id" db:"organization_id"`
	TeamID         *uuid.UUID `json:"team_id,omitempty" db:"team_id"`
	Name           string     `json:"name" db:"name"`
	Description    *string    `json:"description,omitempty" db:"description"`
	Timezone       string     `json:"timezone" db:"timezone"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
}

type ScheduleRotation struct {
	ID             uuid.UUID    `json:"id" db:"id"`
	ScheduleID     uuid.UUID    `json:"schedule_id" db:"schedule_id"`
	Name           string       `json:"name" db:"name"`
	RotationType   RotationType `json:"rotation_type" db:"rotation_type"`
	RotationLength int          `json:"rotation_length" db:"rotation_length"`
	StartDate      time.Time    `json:"start_date" db:"start_date"`
	StartTime      time.Time    `json:"start_time" db:"start_time"`
	EndTime        *time.Time   `json:"end_time,omitempty" db:"end_time"`
	HandoffDay     *int         `json:"handoff_day,omitempty" db:"handoff_day"`
	HandoffTime    time.Time    `json:"handoff_time" db:"handoff_time"`
	CreatedAt      time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at" db:"updated_at"`
}

type ScheduleRotationParticipant struct {
	ID         uuid.UUID `json:"id" db:"id"`
	RotationID uuid.UUID `json:"rotation_id" db:"rotation_id"`
	UserID     uuid.UUID `json:"user_id" db:"user_id"`
	Position   int       `json:"position" db:"position"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

type ScheduleOverride struct {
	ID         uuid.UUID `json:"id" db:"id"`
	ScheduleID uuid.UUID `json:"schedule_id" db:"schedule_id"`
	UserID     uuid.UUID `json:"user_id" db:"user_id"`
	StartTime  time.Time `json:"start_time" db:"start_time"`
	EndTime    time.Time `json:"end_time" db:"end_time"`
	Note       *string   `json:"note,omitempty" db:"note"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

type RotationType string

const (
	RotationTypeDaily  RotationType = "daily"
	RotationTypeWeekly RotationType = "weekly"
	RotationTypeCustom RotationType = "custom"
)

func (r RotationType) String() string {
	return string(r)
}

func (r RotationType) Validate() error {
	switch r {
	case RotationTypeDaily, RotationTypeWeekly, RotationTypeCustom:
		return nil
	default:
		return ErrInvalidRotationType
	}
}

// ScheduleWithRotations includes the schedule and its rotations
type ScheduleWithRotations struct {
	Schedule
	Rotations []*ScheduleRotation `json:"rotations"`
}

// RotationWithParticipants includes the rotation and its participants
type RotationWithParticipants struct {
	ScheduleRotation
	Participants []*ParticipantWithUser `json:"participants"`
}

// ParticipantWithUser includes participant info and user details
type ParticipantWithUser struct {
	ScheduleRotationParticipant
	User User `json:"user"`
}

// OnCallUser represents who is on-call at a specific time
type OnCallUser struct {
	UserID     uuid.UUID `json:"user_id"`
	User       *User     `json:"user,omitempty"`
	ScheduleID uuid.UUID `json:"schedule_id"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
	IsOverride bool      `json:"is_override"`
}
