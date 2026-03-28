package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/nmn3m/pulsar/backend/internal/domain"
	"github.com/nmn3m/pulsar/backend/internal/usecase/repository"
)

type AlertUsecase struct {
	alertRepo      repository.AlertRepository
	alertNotifier  *AlertNotifier
	wsUsecase      *WebSocketUsecase
	webhookUsecase *WebhookUsecase
}

func NewAlertUsecase(alertRepo repository.AlertRepository, alertNotifier *AlertNotifier, wsUsecase *WebSocketUsecase, webhookUsecase *WebhookUsecase) *AlertUsecase {
	return &AlertUsecase{
		alertRepo:      alertRepo,
		alertNotifier:  alertNotifier,
		wsUsecase:      wsUsecase,
		webhookUsecase: webhookUsecase,
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
	DedupKey     *string                `json:"dedup_key"` // Optional deduplication key
}

type UpdateAlertRequest struct {
	Priority     *string                `json:"priority"`
	Message      *string                `json:"message"`
	Description  *string                `json:"description"`
	Tags         []string               `json:"tags"`
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
	Status         []string   `form:"status"`
	Priority       []string   `form:"priority"`
	AssignedToUser *uuid.UUID `form:"assigned_to_user"`
	AssignedToTeam *uuid.UUID `form:"assigned_to_team"`
	Source         *string    `form:"source"`
	Search         *string    `form:"search"`
	Page           int        `form:"page"`
	PageSize       int        `form:"page_size"`
}

type ListAlertsResponse struct {
	Alerts   []*domain.Alert `json:"alerts"`
	Total    int             `json:"total"`
	Page     int             `json:"page"`
	PageSize int             `json:"page_size"`
}

func (s *AlertUsecase) CreateAlert(ctx context.Context, orgID uuid.UUID, req *CreateAlertRequest) (*domain.Alert, error) {
	// Validate priority
	priority := domain.AlertPriority(req.Priority)
	if !priority.IsValid() {
		return nil, fmt.Errorf("invalid priority: %s", req.Priority)
	}

	// Check for deduplication
	if req.DedupKey != nil && *req.DedupKey != "" {
		existingAlert, err := s.alertRepo.FindByDedupKey(ctx, orgID, *req.DedupKey)
		if err != nil {
			return nil, fmt.Errorf("failed to check for duplicate alert: %w", err)
		}

		if existingAlert != nil {
			// Increment dedup count instead of creating new alert
			if err := s.alertRepo.IncrementDedupCount(ctx, existingAlert.ID); err != nil {
				return nil, fmt.Errorf("failed to increment dedup count: %w", err)
			}

			// Refresh the alert to get updated values
			updatedAlert, err := s.alertRepo.GetByID(ctx, existingAlert.ID, orgID)
			if err != nil {
				return nil, fmt.Errorf("failed to get updated alert: %w", err)
			}

			// Broadcast WebSocket event for dedup
			if s.wsUsecase != nil {
				s.wsUsecase.BroadcastAlertEvent(domain.WSEventAlertUpdated, orgID, updatedAlert)
			}

			return updatedAlert, nil
		}
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

	now := time.Now()
	alert := &domain.Alert{
		ID:                uuid.New(),
		OrganizationID:    orgID,
		Source:            req.Source,
		SourceID:          req.SourceID,
		Priority:          priority,
		Status:            domain.AlertStatusOpen,
		Message:           req.Message,
		Description:       req.Description,
		Tags:              tags,
		CustomFields:      customFields,
		EscalationLevel:   0,
		DedupKey:          req.DedupKey,
		DedupCount:        1,
		FirstOccurrenceAt: &now,
		LastOccurrenceAt:  &now,
	}

	if err := s.alertRepo.Create(ctx, alert); err != nil {
		return nil, fmt.Errorf("failed to create alert: %w", err)
	}

	// Send notification for new alert (async, don't fail if notification fails)
	if s.alertNotifier != nil {
		go func() {
			if err := s.alertNotifier.NotifyAlertCreated(context.Background(), alert); err != nil {
				// Log error but don't fail alert creation
				fmt.Printf("Failed to send alert creation notification: %v\n", err)
			}
		}()
	}

	// Broadcast WebSocket event
	if s.wsUsecase != nil {
		s.wsUsecase.BroadcastAlertEvent(domain.WSEventAlertCreated, orgID, alert)
	}

	// Trigger webhooks
	if s.webhookUsecase != nil {
		s.webhookUsecase.TriggerWebhooks(ctx, orgID, "alert.created", map[string]interface{}{
			"alert_id":    alert.ID.String(),
			"source":      alert.Source,
			"priority":    string(alert.Priority),
			"status":      string(alert.Status),
			"message":     alert.Message,
			"description": alert.Description,
			"tags":        alert.Tags,
			"created_at":  alert.CreatedAt,
		})
	}

	return alert, nil
}

func (s *AlertUsecase) GetAlert(ctx context.Context, id, orgID uuid.UUID) (*domain.Alert, error) {
	alert, err := s.alertRepo.GetByID(ctx, id, orgID)
	if err != nil {
		return nil, fmt.Errorf("failed to get alert: %w", err)
	}

	return alert, nil
}

func (s *AlertUsecase) UpdateAlert(ctx context.Context, id, orgID uuid.UUID, req *UpdateAlertRequest) (*domain.Alert, error) {
	alert, err := s.alertRepo.GetByID(ctx, id, orgID)
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

	// Broadcast WebSocket event
	if s.wsUsecase != nil {
		s.wsUsecase.BroadcastAlertEvent(domain.WSEventAlertUpdated, alert.OrganizationID, alert)
	}

	// Trigger webhooks
	if s.webhookUsecase != nil {
		s.webhookUsecase.TriggerWebhooks(ctx, alert.OrganizationID, "alert.updated", map[string]interface{}{
			"alert_id":    alert.ID.String(),
			"source":      alert.Source,
			"priority":    string(alert.Priority),
			"status":      string(alert.Status),
			"message":     alert.Message,
			"description": alert.Description,
			"tags":        alert.Tags,
			"updated_at":  alert.UpdatedAt,
		})
	}

	return alert, nil
}

func (s *AlertUsecase) DeleteAlert(ctx context.Context, id, orgID uuid.UUID) error {
	if err := s.alertRepo.Delete(ctx, id, orgID); err != nil {
		return fmt.Errorf("failed to delete alert: %w", err)
	}

	return nil
}

func (s *AlertUsecase) ListAlerts(ctx context.Context, orgID uuid.UUID, req *ListAlertsRequest) (*ListAlertsResponse, error) {
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

func (s *AlertUsecase) AcknowledgeAlert(ctx context.Context, id, orgID, userID uuid.UUID) error {
	if err := s.alertRepo.Acknowledge(ctx, id, orgID, userID); err != nil {
		return fmt.Errorf("failed to acknowledge alert: %w", err)
	}

	// Send notification for acknowledged alert (async)
	if s.alertNotifier != nil {
		go func() {
			alert, err := s.alertRepo.GetByID(context.Background(), id, orgID)
			if err == nil {
				if err := s.alertNotifier.NotifyAlertAcknowledged(context.Background(), alert, userID); err != nil {
					fmt.Printf("Failed to send alert acknowledgment notification: %v\n", err)
				}
			}
		}()
	}

	// Broadcast WebSocket event and trigger webhooks
	if s.wsUsecase != nil || s.webhookUsecase != nil {
		alert, err := s.alertRepo.GetByID(ctx, id, orgID)
		if err == nil {
			if s.wsUsecase != nil {
				s.wsUsecase.BroadcastAlertEvent(domain.WSEventAlertAcknowledged, alert.OrganizationID, alert)
			}
			if s.webhookUsecase != nil {
				s.webhookUsecase.TriggerWebhooks(ctx, alert.OrganizationID, "alert.acknowledged", map[string]interface{}{
					"alert_id":        alert.ID.String(),
					"source":          alert.Source,
					"priority":        string(alert.Priority),
					"status":          string(alert.Status),
					"message":         alert.Message,
					"acknowledged_at": alert.AcknowledgedAt,
					"acknowledged_by": userID.String(),
				})
			}
		}
	}

	return nil
}

func (s *AlertUsecase) CloseAlert(ctx context.Context, id, orgID, userID uuid.UUID, reason string) error {
	if err := s.alertRepo.Close(ctx, id, orgID, userID, reason); err != nil {
		return fmt.Errorf("failed to close alert: %w", err)
	}

	// Send notification for closed alert (async)
	if s.alertNotifier != nil {
		go func() {
			alert, err := s.alertRepo.GetByID(context.Background(), id, orgID)
			if err == nil {
				if err := s.alertNotifier.NotifyAlertClosed(context.Background(), alert, userID, reason); err != nil {
					fmt.Printf("Failed to send alert closure notification: %v\n", err)
				}
			}
		}()
	}

	// Broadcast WebSocket event and trigger webhooks
	if s.wsUsecase != nil || s.webhookUsecase != nil {
		alert, err := s.alertRepo.GetByID(ctx, id, orgID)
		if err == nil {
			if s.wsUsecase != nil {
				s.wsUsecase.BroadcastAlertEvent(domain.WSEventAlertClosed, alert.OrganizationID, alert)
			}
			if s.webhookUsecase != nil {
				s.webhookUsecase.TriggerWebhooks(ctx, alert.OrganizationID, "alert.closed", map[string]interface{}{
					"alert_id":     alert.ID.String(),
					"source":       alert.Source,
					"priority":     string(alert.Priority),
					"status":       string(alert.Status),
					"message":      alert.Message,
					"closed_at":    alert.ClosedAt,
					"closed_by":    userID.String(),
					"close_reason": reason,
				})
			}
		}
	}

	return nil
}

func (s *AlertUsecase) SnoozeAlert(ctx context.Context, id, orgID uuid.UUID, until time.Time) error {
	if until.Before(time.Now()) {
		return fmt.Errorf("snooze time must be in the future")
	}

	if err := s.alertRepo.Snooze(ctx, id, orgID, until); err != nil {
		return fmt.Errorf("failed to snooze alert: %w", err)
	}

	return nil
}

func (s *AlertUsecase) AssignAlert(ctx context.Context, id, orgID uuid.UUID, userID, teamID *uuid.UUID) error {
	if userID == nil && teamID == nil {
		return fmt.Errorf("must assign to either a user or a team")
	}

	if err := s.alertRepo.Assign(ctx, id, orgID, userID, teamID); err != nil {
		return fmt.Errorf("failed to assign alert: %w", err)
	}

	return nil
}
