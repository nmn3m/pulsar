package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/pulsar/backend/internal/domain"
	"github.com/pulsar/backend/internal/repository"
)

type AlertService struct {
	alertRepo repository.AlertRepository
}

func NewAlertService(alertRepo repository.AlertRepository) *AlertService {
	return &AlertService{
		alertRepo: alertRepo,
	}
}

type CreateAlertRequest struct {
	Source       string                 `json:"source" binding:"required"`
	SourceID     *string                `json:"source_id"`
	Priority     string                 `json:"priority" binding:"required"`
	Message      string                 `json:"message" binding:"required"`
	Description  *string                `json:"description"`
	Tags         []string               `json:"tags"`
	CustomFields map[string]interface{} `json:"custom_fields"`
}

type UpdateAlertRequest struct {
	Priority    *string                 `json:"priority"`
	Message     *string                 `json:"message"`
	Description *string                 `json:"description"`
	Tags        []string                `json:"tags"`
	CustomFields map[string]interface{} `json:"custom_fields"`
}

type AcknowledgeAlertRequest struct {
	UserID uuid.UUID
}

type CloseAlertRequest struct {
	UserID uuid.UUID
	Reason string `json:"reason"`
}

type SnoozeAlertRequest struct {
	Until time.Time `json:"until" binding:"required"`
}

type AssignAlertRequest struct {
	UserID *uuid.UUID `json:"user_id"`
	TeamID *uuid.UUID `json:"team_id"`
}

type ListAlertsRequest struct {
	Status         []string    `form:"status"`
	Priority       []string    `form:"priority"`
	AssignedToUser *uuid.UUID  `form:"assigned_to_user"`
	AssignedToTeam *uuid.UUID  `form:"assigned_to_team"`
	Source         *string     `form:"source"`
	Search         *string     `form:"search"`
	Page           int         `form:"page"`
	PageSize       int         `form:"page_size"`
}

type ListAlertsResponse struct {
	Alerts []*domain.Alert `json:"alerts"`
	Total  int             `json:"total"`
	Page   int             `json:"page"`
	PageSize int           `json:"page_size"`
}

func (s *AlertService) CreateAlert(ctx context.Context, orgID uuid.UUID, req *CreateAlertRequest) (*domain.Alert, error) {
	// Validate priority
	priority := domain.AlertPriority(req.Priority)
	if !priority.IsValid() {
		return nil, fmt.Errorf("invalid priority: %s", req.Priority)
	}

	// Initialize tags and custom fields if nil
	tags := req.Tags
	if tags == nil {
		tags = []string{}
	}

	customFields := req.CustomFields
	if customFields == nil {
		customFields = make(map[string]interface{})
	}

	alert := &domain.Alert{
		ID:             uuid.New(),
		OrganizationID: orgID,
		Source:         req.Source,
		SourceID:       req.SourceID,
		Priority:       priority,
		Status:         domain.AlertStatusOpen,
		Message:        req.Message,
		Description:    req.Description,
		Tags:           tags,
		CustomFields:   customFields,
		EscalationLevel: 0,
	}

	if err := s.alertRepo.Create(ctx, alert); err != nil {
		return nil, fmt.Errorf("failed to create alert: %w", err)
	}

	return alert, nil
}

func (s *AlertService) GetAlert(ctx context.Context, id uuid.UUID) (*domain.Alert, error) {
	alert, err := s.alertRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get alert: %w", err)
	}

	return alert, nil
}

func (s *AlertService) UpdateAlert(ctx context.Context, id uuid.UUID, req *UpdateAlertRequest) (*domain.Alert, error) {
	alert, err := s.alertRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get alert: %w", err)
	}

	// Update fields if provided
	if req.Priority != nil {
		priority := domain.AlertPriority(*req.Priority)
		if !priority.IsValid() {
			return nil, fmt.Errorf("invalid priority: %s", *req.Priority)
		}
		alert.Priority = priority
	}

	if req.Message != nil {
		alert.Message = *req.Message
	}

	if req.Description != nil {
		alert.Description = req.Description
	}

	if req.Tags != nil {
		alert.Tags = req.Tags
	}

	if req.CustomFields != nil {
		alert.CustomFields = req.CustomFields
	}

	if err := s.alertRepo.Update(ctx, alert); err != nil {
		return nil, fmt.Errorf("failed to update alert: %w", err)
	}

	return alert, nil
}

func (s *AlertService) DeleteAlert(ctx context.Context, id uuid.UUID) error {
	if err := s.alertRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete alert: %w", err)
	}

	return nil
}

func (s *AlertService) ListAlerts(ctx context.Context, orgID uuid.UUID, req *ListAlertsRequest) (*ListAlertsResponse, error) {
	// Parse status filters
	var statuses []domain.AlertStatus
	for _, statusStr := range req.Status {
		status := domain.AlertStatus(statusStr)
		if status.IsValid() {
			statuses = append(statuses, status)
		}
	}

	// Parse priority filters
	var priorities []domain.AlertPriority
	for _, priorityStr := range req.Priority {
		priority := domain.AlertPriority(priorityStr)
		if priority.IsValid() {
			priorities = append(priorities, priority)
		}
	}

	// Set defaults
	page := req.Page
	if page < 1 {
		page = 1
	}

	pageSize := req.PageSize
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize

	filter := &domain.AlertFilter{
		OrganizationID: orgID,
		Status:         statuses,
		Priority:       priorities,
		AssignedToUser: req.AssignedToUser,
		AssignedToTeam: req.AssignedToTeam,
		Source:         req.Source,
		Search:         req.Search,
		Limit:          pageSize,
		Offset:         offset,
	}

	alerts, total, err := s.alertRepo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list alerts: %w", err)
	}

	return &ListAlertsResponse{
		Alerts:   alerts,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

func (s *AlertService) AcknowledgeAlert(ctx context.Context, id, userID uuid.UUID) error {
	if err := s.alertRepo.Acknowledge(ctx, id, userID); err != nil {
		return fmt.Errorf("failed to acknowledge alert: %w", err)
	}

	return nil
}

func (s *AlertService) CloseAlert(ctx context.Context, id, userID uuid.UUID, reason string) error {
	if err := s.alertRepo.Close(ctx, id, userID, reason); err != nil {
		return fmt.Errorf("failed to close alert: %w", err)
	}

	return nil
}

func (s *AlertService) SnoozeAlert(ctx context.Context, id uuid.UUID, until time.Time) error {
	if until.Before(time.Now()) {
		return fmt.Errorf("snooze time must be in the future")
	}

	if err := s.alertRepo.Snooze(ctx, id, until); err != nil {
		return fmt.Errorf("failed to snooze alert: %w", err)
	}

	return nil
}

func (s *AlertService) AssignAlert(ctx context.Context, id uuid.UUID, userID, teamID *uuid.UUID) error {
	if userID == nil && teamID == nil {
		return fmt.Errorf("must assign to either a user or a team")
	}

	if err := s.alertRepo.Assign(ctx, id, userID, teamID); err != nil {
		return fmt.Errorf("failed to assign alert: %w", err)
	}

	return nil
}
