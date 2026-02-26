package domain

import (
	"time"

	"github.com/google/uuid"
)

type Schedule struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	TeamID         *uuid.UUID
	Name           string
	Description    *string
	Timezone       string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type ScheduleRotation struct {
	ID             uuid.UUID
	ScheduleID     uuid.UUID
	Name           string
	RotationType   RotationType
	RotationLength int
	StartDate      time.Time
	StartTime      time.Time
	EndTime        *time.Time
	HandoffDay     *int
	HandoffTime    time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type ScheduleRotationParticipant struct {
	ID         uuid.UUID
	RotationID uuid.UUID
	UserID     uuid.UUID
	Position   int
	CreatedAt  time.Time
}

type ScheduleOverride struct {
	ID         uuid.UUID
	ScheduleID uuid.UUID
	UserID     uuid.UUID
	StartTime  time.Time
	EndTime    time.Time
	Note       *string
	CreatedAt  time.Time
	UpdatedAt  time.Time
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
	Rotations []*ScheduleRotation
}

// RotationWithParticipants includes the rotation and its participants
type RotationWithParticipants struct {
	ScheduleRotation
	Participants []*ParticipantWithUser
}

// ParticipantWithUser includes participant info and user details
type ParticipantWithUser struct {
	ScheduleRotationParticipant
	User User
}

// OnCallUser represents who is on-call at a specific time
type OnCallUser struct {
	UserID     uuid.UUID
	User       *User
	ScheduleID uuid.UUID
	StartTime  time.Time
	EndTime    time.Time
	IsOverride bool
}
