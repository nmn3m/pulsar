package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// UserDNDSettings represents a user's Do Not Disturb configuration
type UserDNDSettings struct {
	ID              uuid.UUID       `json:"id" db:"id"`
	UserID          uuid.UUID       `json:"user_id" db:"user_id"`
	Enabled         bool            `json:"enabled" db:"enabled"`
	Schedule        json.RawMessage `json:"schedule" db:"schedule"`
	Overrides       json.RawMessage `json:"overrides" db:"overrides"`
	AllowP1Override bool            `json:"allow_p1_override" db:"allow_p1_override"`
	CreatedAt       time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at" db:"updated_at"`
}

// DNDSchedule represents the weekly schedule for DND
type DNDSchedule struct {
	Weekly   []DNDTimeSlot `json:"weekly"`
	Timezone string        `json:"timezone"` // IANA timezone string
}

// DNDTimeSlot represents a single time slot in the DND schedule
type DNDTimeSlot struct {
	Day   string `json:"day"`   // monday, tuesday, wednesday, thursday, friday, saturday, sunday
	Start string `json:"start"` // HH:MM format (24-hour)
	End   string `json:"end"`   // HH:MM format (24-hour)
}

// DNDOverride represents a one-time DND override period
type DNDOverride struct {
	Start  time.Time `json:"start"`
	End    time.Time `json:"end"`
	Reason string    `json:"reason,omitempty"`
}

// CreateDNDSettingsRequest is the request to create/update DND settings
type CreateDNDSettingsRequest struct {
	Enabled         *bool           `json:"enabled,omitempty"`
	Schedule        json.RawMessage `json:"schedule,omitempty"`
	Overrides       json.RawMessage `json:"overrides,omitempty"`
	AllowP1Override *bool           `json:"allow_p1_override,omitempty"`
}

// UpdateDNDSettingsRequest is the request to update DND settings
type UpdateDNDSettingsRequest struct {
	Enabled         *bool           `json:"enabled,omitempty"`
	Schedule        json.RawMessage `json:"schedule,omitempty"`
	Overrides       json.RawMessage `json:"overrides,omitempty"`
	AllowP1Override *bool           `json:"allow_p1_override,omitempty"`
}

// AddDNDOverrideRequest is the request to add a temporary DND override
type AddDNDOverrideRequest struct {
	Start  time.Time `json:"start" binding:"required"`
	End    time.Time `json:"end" binding:"required"`
	Reason string    `json:"reason,omitempty"`
}

// ParseSchedule parses the raw JSON schedule into a structured format
func (s *UserDNDSettings) ParseSchedule() (*DNDSchedule, error) {
	if len(s.Schedule) == 0 {
		return &DNDSchedule{}, nil
	}
	var schedule DNDSchedule
	if err := json.Unmarshal(s.Schedule, &schedule); err != nil {
		return nil, err
	}
	return &schedule, nil
}

// ParseOverrides parses the raw JSON overrides into a structured format
func (s *UserDNDSettings) ParseOverrides() ([]DNDOverride, error) {
	if len(s.Overrides) == 0 {
		return []DNDOverride{}, nil
	}
	var overrides []DNDOverride
	if err := json.Unmarshal(s.Overrides, &overrides); err != nil {
		return nil, err
	}
	return overrides, nil
}

// Valid day names
var validDays = map[string]int{
	"sunday":    0,
	"monday":    1,
	"tuesday":   2,
	"wednesday": 3,
	"thursday":  4,
	"friday":    5,
	"saturday":  6,
}

// IsValidDay checks if a day name is valid
func IsValidDay(day string) bool {
	_, ok := validDays[day]
	return ok
}

// GetDayIndex returns the weekday index for a day name (0 = Sunday)
func GetDayIndex(day string) int {
	return validDays[day]
}
