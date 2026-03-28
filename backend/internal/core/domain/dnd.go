package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// UserDNDSettings represents a user's Do Not Disturb configuration
type UserDNDSettings struct {
	ID              uuid.UUID
	UserID          uuid.UUID
	Enabled         bool
	Schedule        json.RawMessage
	Overrides       json.RawMessage
	AllowP1Override bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
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
