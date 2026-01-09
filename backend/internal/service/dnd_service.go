package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/nmn3m/pulsar/backend/internal/domain"
	"github.com/nmn3m/pulsar/backend/internal/repository"
)

type DNDService struct {
	dndRepo repository.DNDSettingsRepository
}

func NewDNDService(dndRepo repository.DNDSettingsRepository) *DNDService {
	return &DNDService{dndRepo: dndRepo}
}

// GetSettings retrieves DND settings for a user
func (s *DNDService) GetSettings(ctx context.Context, userID uuid.UUID) (*domain.UserDNDSettings, error) {
	settings, err := s.dndRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get DND settings: %w", err)
	}

	// If no settings exist, return default settings
	if settings == nil {
		return &domain.UserDNDSettings{
			UserID:          userID,
			Enabled:         false,
			Schedule:        json.RawMessage("{}"),
			Overrides:       json.RawMessage("[]"),
			AllowP1Override: true,
		}, nil
	}

	return settings, nil
}

// UpdateSettings updates or creates DND settings for a user
func (s *DNDService) UpdateSettings(ctx context.Context, userID uuid.UUID, req *domain.UpdateDNDSettingsRequest) (*domain.UserDNDSettings, error) {
	// Get existing settings or create new
	settings, err := s.dndRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing DND settings: %w", err)
	}

	if settings == nil {
		// Create new settings
		settings = &domain.UserDNDSettings{
			ID:              uuid.New(),
			UserID:          userID,
			Enabled:         false,
			Schedule:        json.RawMessage("{}"),
			Overrides:       json.RawMessage("[]"),
			AllowP1Override: true,
		}
	}

	// Apply updates
	if req.Enabled != nil {
		settings.Enabled = *req.Enabled
	}
	if req.Schedule != nil {
		// Validate schedule format
		var schedule domain.DNDSchedule
		if err := json.Unmarshal(req.Schedule, &schedule); err != nil {
			return nil, fmt.Errorf("invalid schedule format: %w", err)
		}
		// Validate timezone
		if schedule.Timezone != "" {
			if _, err := time.LoadLocation(schedule.Timezone); err != nil {
				return nil, fmt.Errorf("invalid timezone: %w", err)
			}
		}
		// Validate days
		for _, slot := range schedule.Weekly {
			if !domain.IsValidDay(slot.Day) {
				return nil, fmt.Errorf("invalid day: %s", slot.Day)
			}
		}
		settings.Schedule = req.Schedule
	}
	if req.Overrides != nil {
		// Validate overrides format
		var overrides []domain.DNDOverride
		if err := json.Unmarshal(req.Overrides, &overrides); err != nil {
			return nil, fmt.Errorf("invalid overrides format: %w", err)
		}
		settings.Overrides = req.Overrides
	}
	if req.AllowP1Override != nil {
		settings.AllowP1Override = *req.AllowP1Override
	}

	// Upsert the settings
	if err := s.dndRepo.Upsert(ctx, settings); err != nil {
		return nil, fmt.Errorf("failed to save DND settings: %w", err)
	}

	return settings, nil
}

// AddOverride adds a temporary DND override
func (s *DNDService) AddOverride(ctx context.Context, userID uuid.UUID, req *domain.AddDNDOverrideRequest) (*domain.UserDNDSettings, error) {
	settings, err := s.dndRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get DND settings: %w", err)
	}

	if settings == nil {
		// Create settings with the override
		settings = &domain.UserDNDSettings{
			ID:              uuid.New(),
			UserID:          userID,
			Enabled:         true, // Enable DND when adding override
			Schedule:        json.RawMessage("{}"),
			AllowP1Override: true,
		}
	}

	// Parse existing overrides
	overrides, err := settings.ParseOverrides()
	if err != nil {
		return nil, fmt.Errorf("failed to parse existing overrides: %w", err)
	}

	// Add new override
	newOverride := domain.DNDOverride{
		Start:  req.Start,
		End:    req.End,
		Reason: req.Reason,
	}
	overrides = append(overrides, newOverride)

	// Marshal back to JSON
	overridesJSON, err := json.Marshal(overrides)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal overrides: %w", err)
	}
	settings.Overrides = overridesJSON

	// Save
	if err := s.dndRepo.Upsert(ctx, settings); err != nil {
		return nil, fmt.Errorf("failed to save DND settings: %w", err)
	}

	return settings, nil
}

// RemoveOverride removes a DND override by index
func (s *DNDService) RemoveOverride(ctx context.Context, userID uuid.UUID, index int) (*domain.UserDNDSettings, error) {
	settings, err := s.dndRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get DND settings: %w", err)
	}

	if settings == nil {
		return nil, fmt.Errorf("no DND settings found")
	}

	// Parse existing overrides
	overrides, err := settings.ParseOverrides()
	if err != nil {
		return nil, fmt.Errorf("failed to parse existing overrides: %w", err)
	}

	if index < 0 || index >= len(overrides) {
		return nil, fmt.Errorf("invalid override index")
	}

	// Remove override at index
	overrides = append(overrides[:index], overrides[index+1:]...)

	// Marshal back to JSON
	overridesJSON, err := json.Marshal(overrides)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal overrides: %w", err)
	}
	settings.Overrides = overridesJSON

	// Save
	if err := s.dndRepo.Update(ctx, settings); err != nil {
		return nil, fmt.Errorf("failed to save DND settings: %w", err)
	}

	return settings, nil
}

// IsInDNDMode checks if a user is currently in DND mode
// Returns true if the user should not be notified
// If priority is P1 and AllowP1Override is true, returns false (allow notification)
func (s *DNDService) IsInDNDMode(ctx context.Context, userID uuid.UUID, priority domain.AlertPriority) (bool, error) {
	settings, err := s.dndRepo.GetByUserID(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("failed to get DND settings: %w", err)
	}

	// No settings or DND disabled
	if settings == nil || !settings.Enabled {
		return false, nil
	}

	// Check if P1 alerts bypass DND
	if priority == domain.PriorityP1 && settings.AllowP1Override {
		return false, nil
	}

	now := time.Now()

	// First check overrides (temporary DND periods)
	overrides, err := settings.ParseOverrides()
	if err != nil {
		return false, fmt.Errorf("failed to parse overrides: %w", err)
	}

	for _, override := range overrides {
		if now.After(override.Start) && now.Before(override.End) {
			return true, nil // Currently in an override period
		}
	}

	// Check weekly schedule
	schedule, err := settings.ParseSchedule()
	if err != nil {
		return false, fmt.Errorf("failed to parse schedule: %w", err)
	}

	if len(schedule.Weekly) == 0 {
		return false, nil // No schedule defined
	}

	// Load timezone
	loc := time.UTC
	if schedule.Timezone != "" {
		var err error
		loc, err = time.LoadLocation(schedule.Timezone)
		if err != nil {
			return false, fmt.Errorf("invalid timezone: %w", err)
		}
	}

	// Get current time in the user's timezone
	localNow := now.In(loc)
	currentDay := getDayName(localNow.Weekday())
	currentTimeStr := localNow.Format("15:04")

	// Check if current time falls within any DND slot
	for _, slot := range schedule.Weekly {
		if slot.Day != currentDay {
			continue
		}

		if isTimeInRange(currentTimeStr, slot.Start, slot.End) {
			return true, nil
		}
	}

	return false, nil
}

// CleanExpiredOverrides removes overrides that have ended
func (s *DNDService) CleanExpiredOverrides(ctx context.Context, userID uuid.UUID) error {
	settings, err := s.dndRepo.GetByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get DND settings: %w", err)
	}

	if settings == nil {
		return nil
	}

	overrides, err := settings.ParseOverrides()
	if err != nil {
		return fmt.Errorf("failed to parse overrides: %w", err)
	}

	now := time.Now()
	var activeOverrides []domain.DNDOverride
	for _, override := range overrides {
		if override.End.After(now) {
			activeOverrides = append(activeOverrides, override)
		}
	}

	if len(activeOverrides) == len(overrides) {
		return nil // No changes needed
	}

	overridesJSON, err := json.Marshal(activeOverrides)
	if err != nil {
		return fmt.Errorf("failed to marshal overrides: %w", err)
	}
	settings.Overrides = overridesJSON

	return s.dndRepo.Update(ctx, settings)
}

// DeleteSettings removes all DND settings for a user
func (s *DNDService) DeleteSettings(ctx context.Context, userID uuid.UUID) error {
	return s.dndRepo.Delete(ctx, userID)
}

// Helper function to get day name from weekday
func getDayName(weekday time.Weekday) string {
	days := []string{"sunday", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday"}
	return days[weekday]
}

// Helper function to check if a time is within a range
// Handles overnight ranges (e.g., 22:00 to 08:00)
func isTimeInRange(current, start, end string) bool {
	// Simple string comparison for HH:MM format
	if start <= end {
		// Normal range (e.g., 09:00 to 17:00)
		return current >= start && current <= end
	}
	// Overnight range (e.g., 22:00 to 08:00)
	return current >= start || current <= end
}
