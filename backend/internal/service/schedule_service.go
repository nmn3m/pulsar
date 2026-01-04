package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/pulsar/backend/internal/domain"
	"github.com/pulsar/backend/internal/repository"
)

type ScheduleService struct {
	scheduleRepo repository.ScheduleRepository
	userRepo     repository.UserRepository
}

func NewScheduleService(scheduleRepo repository.ScheduleRepository, userRepo repository.UserRepository) *ScheduleService {
	return &ScheduleService{
		scheduleRepo: scheduleRepo,
		userRepo:     userRepo,
	}
}

// Request/Response types

type CreateScheduleRequest struct {
	TeamID      *uuid.UUID `json:"team_id"`
	Name        string     `json:"name" binding:"required"`
	Description *string    `json:"description"`
	Timezone    string     `json:"timezone"`
}

type UpdateScheduleRequest struct {
	TeamID      *uuid.UUID `json:"team_id"`
	Name        *string    `json:"name"`
	Description *string    `json:"description"`
	Timezone    *string    `json:"timezone"`
}

type CreateRotationRequest struct {
	Name           string     `json:"name" binding:"required"`
	RotationType   string     `json:"rotation_type" binding:"required"`
	RotationLength int        `json:"rotation_length" binding:"required"`
	StartDate      string     `json:"start_date" binding:"required"`
	StartTime      string     `json:"start_time"`
	EndTime        *string    `json:"end_time"`
	HandoffDay     *int       `json:"handoff_day"`
	HandoffTime    string     `json:"handoff_time"`
}

type UpdateRotationRequest struct {
	Name           *string `json:"name"`
	RotationType   *string `json:"rotation_type"`
	RotationLength *int    `json:"rotation_length"`
	StartDate      *string `json:"start_date"`
	StartTime      *string `json:"start_time"`
	EndTime        *string `json:"end_time"`
	HandoffDay     *int    `json:"handoff_day"`
	HandoffTime    *string `json:"handoff_time"`
}

type AddParticipantRequest struct {
	UserID   uuid.UUID `json:"user_id" binding:"required"`
	Position int       `json:"position" binding:"required"`
}

type ReorderParticipantsRequest struct {
	UserIDs []uuid.UUID `json:"user_ids" binding:"required"`
}

type CreateOverrideRequest struct {
	UserID    uuid.UUID `json:"user_id" binding:"required"`
	StartTime string    `json:"start_time" binding:"required"`
	EndTime   string    `json:"end_time" binding:"required"`
	Note      *string   `json:"note"`
}

type UpdateOverrideRequest struct {
	UserID    *uuid.UUID `json:"user_id"`
	StartTime *string    `json:"start_time"`
	EndTime   *string    `json:"end_time"`
	Note      *string    `json:"note"`
}

// Schedule CRUD

func (s *ScheduleService) CreateSchedule(ctx context.Context, orgID uuid.UUID, req *CreateScheduleRequest) (*domain.Schedule, error) {
	timezone := req.Timezone
	if timezone == "" {
		timezone = "UTC"
	}

	schedule := &domain.Schedule{
		ID:             uuid.New(),
		OrganizationID: orgID,
		TeamID:         req.TeamID,
		Name:           req.Name,
		Description:    req.Description,
		Timezone:       timezone,
	}

	if err := s.scheduleRepo.Create(ctx, schedule); err != nil {
		return nil, fmt.Errorf("failed to create schedule: %w", err)
	}

	return schedule, nil
}

func (s *ScheduleService) GetSchedule(ctx context.Context, id uuid.UUID) (*domain.Schedule, error) {
	schedule, err := s.scheduleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get schedule: %w", err)
	}

	return schedule, nil
}

func (s *ScheduleService) GetScheduleWithRotations(ctx context.Context, id uuid.UUID) (*domain.ScheduleWithRotations, error) {
	schedule, err := s.scheduleRepo.GetWithRotations(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get schedule with rotations: %w", err)
	}

	return schedule, nil
}

func (s *ScheduleService) UpdateSchedule(ctx context.Context, id uuid.UUID, req *UpdateScheduleRequest) (*domain.Schedule, error) {
	schedule, err := s.scheduleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get schedule: %w", err)
	}

	if req.Name != nil {
		schedule.Name = *req.Name
	}
	if req.Description != nil {
		schedule.Description = req.Description
	}
	if req.Timezone != nil {
		schedule.Timezone = *req.Timezone
	}
	if req.TeamID != nil {
		schedule.TeamID = req.TeamID
	}

	if err := s.scheduleRepo.Update(ctx, schedule); err != nil {
		return nil, fmt.Errorf("failed to update schedule: %w", err)
	}

	return schedule, nil
}

func (s *ScheduleService) DeleteSchedule(ctx context.Context, id uuid.UUID) error {
	if err := s.scheduleRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete schedule: %w", err)
	}

	return nil
}

func (s *ScheduleService) ListSchedules(ctx context.Context, orgID uuid.UUID, page, pageSize int) ([]*domain.Schedule, error) {
	if page < 1 {
		page = 1
	}

	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize

	schedules, err := s.scheduleRepo.List(ctx, orgID, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list schedules: %w", err)
	}

	return schedules, nil
}

// Rotation CRUD

func (s *ScheduleService) CreateRotation(ctx context.Context, scheduleID uuid.UUID, req *CreateRotationRequest) (*domain.ScheduleRotation, error) {
	// Validate rotation type
	rotationType := domain.RotationType(req.RotationType)
	if err := rotationType.Validate(); err != nil {
		return nil, err
	}

	// Parse start date
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start_date format: %w", err)
	}

	// Parse times
	startTime, err := time.Parse("15:04", req.StartTime)
	if err != nil {
		startTime, _ = time.Parse("15:04", "00:00")
	}

	handoffTime, err := time.Parse("15:04", req.HandoffTime)
	if err != nil {
		handoffTime, _ = time.Parse("15:04", "09:00")
	}

	var endTime *time.Time
	if req.EndTime != nil {
		t, err := time.Parse("15:04", *req.EndTime)
		if err == nil {
			endTime = &t
		}
	}

	rotation := &domain.ScheduleRotation{
		ID:             uuid.New(),
		ScheduleID:     scheduleID,
		Name:           req.Name,
		RotationType:   rotationType,
		RotationLength: req.RotationLength,
		StartDate:      startDate,
		StartTime:      startTime,
		EndTime:        endTime,
		HandoffDay:     req.HandoffDay,
		HandoffTime:    handoffTime,
	}

	if err := s.scheduleRepo.CreateRotation(ctx, rotation); err != nil {
		return nil, fmt.Errorf("failed to create rotation: %w", err)
	}

	return rotation, nil
}

func (s *ScheduleService) GetRotation(ctx context.Context, id uuid.UUID) (*domain.ScheduleRotation, error) {
	rotation, err := s.scheduleRepo.GetRotation(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get rotation: %w", err)
	}

	return rotation, nil
}

func (s *ScheduleService) UpdateRotation(ctx context.Context, id uuid.UUID, req *UpdateRotationRequest) (*domain.ScheduleRotation, error) {
	rotation, err := s.scheduleRepo.GetRotation(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get rotation: %w", err)
	}

	if req.Name != nil {
		rotation.Name = *req.Name
	}
	if req.RotationType != nil {
		rotationType := domain.RotationType(*req.RotationType)
		if err := rotationType.Validate(); err != nil {
			return nil, err
		}
		rotation.RotationType = rotationType
	}
	if req.RotationLength != nil {
		rotation.RotationLength = *req.RotationLength
	}
	if req.StartDate != nil {
		startDate, err := time.Parse("2006-01-02", *req.StartDate)
		if err != nil {
			return nil, fmt.Errorf("invalid start_date format: %w", err)
		}
		rotation.StartDate = startDate
	}
	if req.StartTime != nil {
		startTime, err := time.Parse("15:04", *req.StartTime)
		if err != nil {
			return nil, fmt.Errorf("invalid start_time format: %w", err)
		}
		rotation.StartTime = startTime
	}
	if req.EndTime != nil {
		endTime, err := time.Parse("15:04", *req.EndTime)
		if err != nil {
			return nil, fmt.Errorf("invalid end_time format: %w", err)
		}
		rotation.EndTime = &endTime
	}
	if req.HandoffDay != nil {
		rotation.HandoffDay = req.HandoffDay
	}
	if req.HandoffTime != nil {
		handoffTime, err := time.Parse("15:04", *req.HandoffTime)
		if err != nil {
			return nil, fmt.Errorf("invalid handoff_time format: %w", err)
		}
		rotation.HandoffTime = handoffTime
	}

	if err := s.scheduleRepo.UpdateRotation(ctx, rotation); err != nil {
		return nil, fmt.Errorf("failed to update rotation: %w", err)
	}

	return rotation, nil
}

func (s *ScheduleService) DeleteRotation(ctx context.Context, id uuid.UUID) error {
	if err := s.scheduleRepo.DeleteRotation(ctx, id); err != nil {
		return fmt.Errorf("failed to delete rotation: %w", err)
	}

	return nil
}

func (s *ScheduleService) ListRotations(ctx context.Context, scheduleID uuid.UUID) ([]*domain.ScheduleRotation, error) {
	rotations, err := s.scheduleRepo.ListRotations(ctx, scheduleID)
	if err != nil {
		return nil, fmt.Errorf("failed to list rotations: %w", err)
	}

	return rotations, nil
}

// Rotation participants

func (s *ScheduleService) AddParticipant(ctx context.Context, rotationID uuid.UUID, req *AddParticipantRequest) (*domain.ScheduleRotationParticipant, error) {
	// Verify user exists
	_, err := s.userRepo.GetByID(ctx, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	participant := &domain.ScheduleRotationParticipant{
		ID:         uuid.New(),
		RotationID: rotationID,
		UserID:     req.UserID,
		Position:   req.Position,
	}

	if err := s.scheduleRepo.AddParticipant(ctx, participant); err != nil {
		return nil, fmt.Errorf("failed to add participant: %w", err)
	}

	return participant, nil
}

func (s *ScheduleService) RemoveParticipant(ctx context.Context, rotationID, userID uuid.UUID) error {
	if err := s.scheduleRepo.RemoveParticipant(ctx, rotationID, userID); err != nil {
		return fmt.Errorf("failed to remove participant: %w", err)
	}

	return nil
}

func (s *ScheduleService) ListParticipants(ctx context.Context, rotationID uuid.UUID) ([]*domain.ParticipantWithUser, error) {
	participants, err := s.scheduleRepo.ListParticipants(ctx, rotationID)
	if err != nil {
		return nil, fmt.Errorf("failed to list participants: %w", err)
	}

	return participants, nil
}

func (s *ScheduleService) ReorderParticipants(ctx context.Context, rotationID uuid.UUID, req *ReorderParticipantsRequest) error {
	if err := s.scheduleRepo.ReorderParticipants(ctx, rotationID, req.UserIDs); err != nil {
		return fmt.Errorf("failed to reorder participants: %w", err)
	}

	return nil
}

// Overrides

func (s *ScheduleService) CreateOverride(ctx context.Context, scheduleID uuid.UUID, req *CreateOverrideRequest) (*domain.ScheduleOverride, error) {
	// Verify user exists
	_, err := s.userRepo.GetByID(ctx, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// Parse times
	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		return nil, fmt.Errorf("invalid start_time format: %w", err)
	}

	endTime, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		return nil, fmt.Errorf("invalid end_time format: %w", err)
	}

	if endTime.Before(startTime) || endTime.Equal(startTime) {
		return nil, fmt.Errorf("end_time must be after start_time")
	}

	override := &domain.ScheduleOverride{
		ID:         uuid.New(),
		ScheduleID: scheduleID,
		UserID:     req.UserID,
		StartTime:  startTime,
		EndTime:    endTime,
		Note:       req.Note,
	}

	if err := s.scheduleRepo.CreateOverride(ctx, override); err != nil {
		return nil, fmt.Errorf("failed to create override: %w", err)
	}

	return override, nil
}

func (s *ScheduleService) GetOverride(ctx context.Context, id uuid.UUID) (*domain.ScheduleOverride, error) {
	override, err := s.scheduleRepo.GetOverride(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get override: %w", err)
	}

	return override, nil
}

func (s *ScheduleService) UpdateOverride(ctx context.Context, id uuid.UUID, req *UpdateOverrideRequest) (*domain.ScheduleOverride, error) {
	override, err := s.scheduleRepo.GetOverride(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get override: %w", err)
	}

	if req.UserID != nil {
		// Verify user exists
		_, err := s.userRepo.GetByID(ctx, *req.UserID)
		if err != nil {
			return nil, fmt.Errorf("user not found")
		}
		override.UserID = *req.UserID
	}

	if req.StartTime != nil {
		startTime, err := time.Parse(time.RFC3339, *req.StartTime)
		if err != nil {
			return nil, fmt.Errorf("invalid start_time format: %w", err)
		}
		override.StartTime = startTime
	}

	if req.EndTime != nil {
		endTime, err := time.Parse(time.RFC3339, *req.EndTime)
		if err != nil {
			return nil, fmt.Errorf("invalid end_time format: %w", err)
		}
		override.EndTime = endTime
	}

	if override.EndTime.Before(override.StartTime) || override.EndTime.Equal(override.StartTime) {
		return nil, fmt.Errorf("end_time must be after start_time")
	}

	if req.Note != nil {
		override.Note = req.Note
	}

	if err := s.scheduleRepo.UpdateOverride(ctx, override); err != nil {
		return nil, fmt.Errorf("failed to update override: %w", err)
	}

	return override, nil
}

func (s *ScheduleService) DeleteOverride(ctx context.Context, id uuid.UUID) error {
	if err := s.scheduleRepo.DeleteOverride(ctx, id); err != nil {
		return fmt.Errorf("failed to delete override: %w", err)
	}

	return nil
}

func (s *ScheduleService) ListOverrides(ctx context.Context, scheduleID uuid.UUID, start, end time.Time) ([]*domain.ScheduleOverride, error) {
	overrides, err := s.scheduleRepo.ListOverrides(ctx, scheduleID, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to list overrides: %w", err)
	}

	return overrides, nil
}

// On-call calculation

func (s *ScheduleService) GetOnCallUser(ctx context.Context, scheduleID uuid.UUID, at time.Time) (*domain.OnCallUser, error) {
	// First, check for overrides
	overrides, err := s.scheduleRepo.ListOverrides(ctx, scheduleID, at, at.Add(1*time.Second))
	if err != nil {
		return nil, fmt.Errorf("failed to check overrides: %w", err)
	}

	if len(overrides) > 0 {
		// Return the most recent override
		override := overrides[0]
		user, _ := s.userRepo.GetByID(ctx, override.UserID)

		return &domain.OnCallUser{
			UserID:     override.UserID,
			User:       user,
			ScheduleID: scheduleID,
			StartTime:  override.StartTime,
			EndTime:    override.EndTime,
			IsOverride: true,
		}, nil
	}

	// No override, calculate from rotation
	rotations, err := s.scheduleRepo.ListRotations(ctx, scheduleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get rotations: %w", err)
	}

	if len(rotations) == 0 {
		return nil, fmt.Errorf("no rotations configured for schedule")
	}

	// Use the first rotation for simplicity (in production, you'd handle multiple rotations)
	rotation := rotations[0]

	participants, err := s.scheduleRepo.ListParticipants(ctx, rotation.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get participants: %w", err)
	}

	if len(participants) == 0 {
		return nil, fmt.Errorf("no participants in rotation")
	}

	// Calculate who is on-call based on rotation type
	onCallUser := s.calculateOnCallFromRotation(rotation, participants, at)
	if onCallUser == nil {
		return nil, fmt.Errorf("could not determine on-call user")
	}

	onCallUser.ScheduleID = scheduleID
	onCallUser.IsOverride = false

	return onCallUser, nil
}

func (s *ScheduleService) calculateOnCallFromRotation(
	rotation *domain.ScheduleRotation,
	participants []*domain.ParticipantWithUser,
	at time.Time,
) *domain.OnCallUser {
	if len(participants) == 0 {
		return nil
	}

	// Calculate days since rotation start
	daysSinceStart := int(at.Sub(rotation.StartDate).Hours() / 24)
	if daysSinceStart < 0 {
		return nil // Before rotation starts
	}

	// Calculate which participant based on rotation type
	var participantIndex int
	switch rotation.RotationType {
	case domain.RotationTypeDaily:
		participantIndex = (daysSinceStart / rotation.RotationLength) % len(participants)
	case domain.RotationTypeWeekly:
		weeksSinceStart := daysSinceStart / 7
		participantIndex = (weeksSinceStart / rotation.RotationLength) % len(participants)
	case domain.RotationTypeCustom:
		participantIndex = (daysSinceStart / rotation.RotationLength) % len(participants)
	default:
		participantIndex = 0
	}

	participant := participants[participantIndex]

	// Calculate the shift start and end times for this participant
	shiftStart := rotation.StartDate.AddDate(0, 0, participantIndex*rotation.RotationLength)
	shiftEnd := shiftStart.AddDate(0, 0, rotation.RotationLength)

	return &domain.OnCallUser{
		UserID:    participant.UserID,
		User:      &participant.User,
		StartTime: shiftStart,
		EndTime:   shiftEnd,
	}
}
