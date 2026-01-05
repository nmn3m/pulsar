package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nmn3m/pulsar/backend/internal/domain"
	"github.com/nmn3m/pulsar/backend/internal/repository"
)

type IncidentService struct {
	incidentRepo repository.IncidentRepository
}

func NewIncidentService(incidentRepo repository.IncidentRepository) *IncidentService {
	return &IncidentService{
		incidentRepo: incidentRepo,
	}
}

// Request/Response types

type CreateIncidentRequest struct {
	Title            string     `json:"title" binding:"required"`
	Description      *string    `json:"description"`
	Severity         string     `json:"severity" binding:"required"`
	Priority         string     `json:"priority" binding:"required"`
	AssignedToTeamID *uuid.UUID `json:"assigned_to_team_id"`
}

type UpdateIncidentRequest struct {
	Title            *string    `json:"title"`
	Description      *string    `json:"description"`
	Severity         *string    `json:"severity"`
	Status           *string    `json:"status"`
	Priority         *string    `json:"priority"`
	AssignedToTeamID *uuid.UUID `json:"assigned_to_team_id"`
}

type AddResponderRequest struct {
	UserID uuid.UUID `json:"user_id" binding:"required"`
	Role   string    `json:"role" binding:"required"`
}

type UpdateResponderRoleRequest struct {
	Role string `json:"role" binding:"required"`
}

type AddNoteRequest struct {
	Note string `json:"note" binding:"required"`
}

type LinkAlertRequest struct {
	AlertID uuid.UUID `json:"alert_id" binding:"required"`
}

type ListIncidentsRequest struct {
	Status           []string   `form:"status"`
	Severity         []string   `form:"severity"`
	AssignedToTeamID *uuid.UUID `form:"assigned_to_team_id"`
	Search           *string    `form:"search"`
	Page             int        `form:"page"`
	PageSize         int        `form:"page_size"`
}

type ListIncidentsResponse struct {
	Incidents []*domain.Incident `json:"incidents"`
	Total     int                `json:"total"`
	Page      int                `json:"page"`
	PageSize  int                `json:"page_size"`
}

// Incident CRUD

func (s *IncidentService) CreateIncident(ctx context.Context, orgID, userID uuid.UUID, req *CreateIncidentRequest) (*domain.Incident, error) {
	// Validate severity
	severity := domain.IncidentSeverity(req.Severity)
	if !severity.IsValid() {
		return nil, fmt.Errorf("invalid severity: %s", req.Severity)
	}

	// Validate priority
	priority := domain.AlertPriority(req.Priority)
	if !priority.IsValid() {
		return nil, fmt.Errorf("invalid priority: %s", req.Priority)
	}

	incident := &domain.Incident{
		ID:               uuid.New(),
		OrganizationID:   orgID,
		Title:            req.Title,
		Description:      req.Description,
		Severity:         severity,
		Status:           domain.IncidentStatusInvestigating,
		Priority:         priority,
		CreatedByUserID:  userID,
		AssignedToTeamID: req.AssignedToTeamID,
		StartedAt:        time.Now(),
	}

	if err := s.incidentRepo.Create(ctx, incident); err != nil {
		return nil, fmt.Errorf("failed to create incident: %w", err)
	}

	// Add timeline event for creation
	timelineEvent := &domain.IncidentTimelineEvent{
		ID:          uuid.New(),
		IncidentID:  incident.ID,
		EventType:   domain.TimelineEventCreated,
		UserID:      &userID,
		Description: fmt.Sprintf("Incident created with severity %s", severity),
		Metadata:    make(map[string]interface{}),
	}

	if err := s.incidentRepo.AddTimelineEvent(ctx, timelineEvent); err != nil {
		// Log error but don't fail the incident creation
		fmt.Printf("Failed to add timeline event: %v\n", err)
	}

	return incident, nil
}

func (s *IncidentService) GetIncident(ctx context.Context, id uuid.UUID) (*domain.Incident, error) {
	incident, err := s.incidentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get incident: %w", err)
	}

	return incident, nil
}

func (s *IncidentService) GetIncidentWithDetails(ctx context.Context, id uuid.UUID) (*domain.IncidentWithDetails, error) {
	incident, err := s.incidentRepo.GetWithDetails(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get incident with details: %w", err)
	}

	return incident, nil
}

func (s *IncidentService) UpdateIncident(ctx context.Context, id, userID uuid.UUID, req *UpdateIncidentRequest) (*domain.Incident, error) {
	incident, err := s.incidentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get incident: %w", err)
	}

	// Track changes for timeline
	var changes []string

	// Update fields if provided
	if req.Title != nil {
		incident.Title = *req.Title
		changes = append(changes, "title")
	}

	if req.Description != nil {
		incident.Description = req.Description
		changes = append(changes, "description")
	}

	if req.Severity != nil {
		severity := domain.IncidentSeverity(*req.Severity)
		if !severity.IsValid() {
			return nil, fmt.Errorf("invalid severity: %s", *req.Severity)
		}
		oldSeverity := incident.Severity
		incident.Severity = severity

		// Add timeline event for severity change
		if oldSeverity != severity {
			timelineEvent := &domain.IncidentTimelineEvent{
				ID:          uuid.New(),
				IncidentID:  incident.ID,
				EventType:   domain.TimelineEventSeverityChanged,
				UserID:      &userID,
				Description: fmt.Sprintf("Severity changed from %s to %s", oldSeverity, severity),
				Metadata: map[string]interface{}{
					"old_severity": oldSeverity,
					"new_severity": severity,
				},
			}
			s.incidentRepo.AddTimelineEvent(ctx, timelineEvent)
		}
	}

	if req.Status != nil {
		status := domain.IncidentStatus(*req.Status)
		if !status.IsValid() {
			return nil, fmt.Errorf("invalid status: %s", *req.Status)
		}
		oldStatus := incident.Status
		incident.Status = status

		// If resolved, set resolved_at
		if status == domain.IncidentStatusResolved && oldStatus != domain.IncidentStatusResolved {
			now := time.Now()
			incident.ResolvedAt = &now

			// Add timeline event for resolution
			timelineEvent := &domain.IncidentTimelineEvent{
				ID:          uuid.New(),
				IncidentID:  incident.ID,
				EventType:   domain.TimelineEventResolved,
				UserID:      &userID,
				Description: "Incident resolved",
				Metadata:    make(map[string]interface{}),
			}
			s.incidentRepo.AddTimelineEvent(ctx, timelineEvent)
		}

		// Add timeline event for status change
		if oldStatus != status {
			timelineEvent := &domain.IncidentTimelineEvent{
				ID:          uuid.New(),
				IncidentID:  incident.ID,
				EventType:   domain.TimelineEventStatusChanged,
				UserID:      &userID,
				Description: fmt.Sprintf("Status changed from %s to %s", oldStatus, status),
				Metadata: map[string]interface{}{
					"old_status": oldStatus,
					"new_status": status,
				},
			}
			s.incidentRepo.AddTimelineEvent(ctx, timelineEvent)
		}
	}

	if req.Priority != nil {
		priority := domain.AlertPriority(*req.Priority)
		if !priority.IsValid() {
			return nil, fmt.Errorf("invalid priority: %s", *req.Priority)
		}
		incident.Priority = priority
	}

	if req.AssignedToTeamID != nil {
		incident.AssignedToTeamID = req.AssignedToTeamID
	}

	if err := s.incidentRepo.Update(ctx, incident); err != nil {
		return nil, fmt.Errorf("failed to update incident: %w", err)
	}

	return incident, nil
}

func (s *IncidentService) DeleteIncident(ctx context.Context, id uuid.UUID) error {
	if err := s.incidentRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete incident: %w", err)
	}

	return nil
}

func (s *IncidentService) ListIncidents(ctx context.Context, orgID uuid.UUID, req *ListIncidentsRequest) (*ListIncidentsResponse, error) {
	// Parse status filters
	var statuses []domain.IncidentStatus
	for _, statusStr := range req.Status {
		status := domain.IncidentStatus(statusStr)
		if status.IsValid() {
			statuses = append(statuses, status)
		}
	}

	// Parse severity filters
	var severities []domain.IncidentSeverity
	for _, severityStr := range req.Severity {
		severity := domain.IncidentSeverity(severityStr)
		if severity.IsValid() {
			severities = append(severities, severity)
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

	filter := &domain.IncidentFilter{
		OrganizationID:   orgID,
		Status:           statuses,
		Severity:         severities,
		AssignedToTeamID: req.AssignedToTeamID,
		Search:           req.Search,
		Limit:            pageSize,
		Offset:           offset,
	}

	incidents, total, err := s.incidentRepo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list incidents: %w", err)
	}

	return &ListIncidentsResponse{
		Incidents: incidents,
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
	}, nil
}

// Responder management

func (s *IncidentService) AddResponder(ctx context.Context, incidentID uuid.UUID, userID uuid.UUID, req *AddResponderRequest) (*domain.IncidentResponder, error) {
	role := domain.ResponderRole(req.Role)
	if !role.IsValid() {
		return nil, fmt.Errorf("invalid responder role: %s", req.Role)
	}

	responder := &domain.IncidentResponder{
		ID:         uuid.New(),
		IncidentID: incidentID,
		UserID:     req.UserID,
		Role:       role,
	}

	if err := s.incidentRepo.AddResponder(ctx, responder); err != nil {
		return nil, fmt.Errorf("failed to add responder: %w", err)
	}

	// Add timeline event
	timelineEvent := &domain.IncidentTimelineEvent{
		ID:          uuid.New(),
		IncidentID:  incidentID,
		EventType:   domain.TimelineEventResponderAdded,
		UserID:      &userID,
		Description: fmt.Sprintf("Responder added with role %s", role),
		Metadata: map[string]interface{}{
			"responder_user_id": req.UserID.String(),
			"role":              role,
		},
	}

	if err := s.incidentRepo.AddTimelineEvent(ctx, timelineEvent); err != nil {
		fmt.Printf("Failed to add timeline event: %v\n", err)
	}

	return responder, nil
}

func (s *IncidentService) RemoveResponder(ctx context.Context, incidentID, responderUserID, actionUserID uuid.UUID) error {
	if err := s.incidentRepo.RemoveResponder(ctx, incidentID, responderUserID); err != nil {
		return fmt.Errorf("failed to remove responder: %w", err)
	}

	// Add timeline event
	timelineEvent := &domain.IncidentTimelineEvent{
		ID:          uuid.New(),
		IncidentID:  incidentID,
		EventType:   domain.TimelineEventResponderRemoved,
		UserID:      &actionUserID,
		Description: "Responder removed",
		Metadata: map[string]interface{}{
			"responder_user_id": responderUserID.String(),
		},
	}

	if err := s.incidentRepo.AddTimelineEvent(ctx, timelineEvent); err != nil {
		fmt.Printf("Failed to add timeline event: %v\n", err)
	}

	return nil
}

func (s *IncidentService) UpdateResponderRole(ctx context.Context, incidentID, responderUserID uuid.UUID, req *UpdateResponderRoleRequest) error {
	role := domain.ResponderRole(req.Role)
	if !role.IsValid() {
		return fmt.Errorf("invalid responder role: %s", req.Role)
	}

	if err := s.incidentRepo.UpdateResponderRole(ctx, incidentID, responderUserID, role); err != nil {
		return fmt.Errorf("failed to update responder role: %w", err)
	}

	return nil
}

func (s *IncidentService) ListResponders(ctx context.Context, incidentID uuid.UUID) ([]*domain.ResponderWithUser, error) {
	responders, err := s.incidentRepo.ListResponders(ctx, incidentID)
	if err != nil {
		return nil, fmt.Errorf("failed to list responders: %w", err)
	}

	return responders, nil
}

// Timeline management

func (s *IncidentService) AddNote(ctx context.Context, incidentID, userID uuid.UUID, req *AddNoteRequest) (*domain.IncidentTimelineEvent, error) {
	event := &domain.IncidentTimelineEvent{
		ID:          uuid.New(),
		IncidentID:  incidentID,
		EventType:   domain.TimelineEventNoteAdded,
		UserID:      &userID,
		Description: req.Note,
		Metadata:    make(map[string]interface{}),
	}

	if err := s.incidentRepo.AddTimelineEvent(ctx, event); err != nil {
		return nil, fmt.Errorf("failed to add note: %w", err)
	}

	return event, nil
}

func (s *IncidentService) GetTimeline(ctx context.Context, incidentID uuid.UUID) ([]*domain.TimelineEventWithUser, error) {
	timeline, err := s.incidentRepo.GetTimeline(ctx, incidentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get timeline: %w", err)
	}

	return timeline, nil
}

// Alert linking

func (s *IncidentService) LinkAlert(ctx context.Context, incidentID, userID uuid.UUID, req *LinkAlertRequest) (*domain.IncidentAlert, error) {
	link := &domain.IncidentAlert{
		ID:             uuid.New(),
		IncidentID:     incidentID,
		AlertID:        req.AlertID,
		LinkedByUserID: &userID,
	}

	if err := s.incidentRepo.LinkAlert(ctx, link); err != nil {
		return nil, fmt.Errorf("failed to link alert: %w", err)
	}

	// Add timeline event
	timelineEvent := &domain.IncidentTimelineEvent{
		ID:          uuid.New(),
		IncidentID:  incidentID,
		EventType:   domain.TimelineEventAlertLinked,
		UserID:      &userID,
		Description: "Alert linked to incident",
		Metadata: map[string]interface{}{
			"alert_id": req.AlertID.String(),
		},
	}

	if err := s.incidentRepo.AddTimelineEvent(ctx, timelineEvent); err != nil {
		fmt.Printf("Failed to add timeline event: %v\n", err)
	}

	return link, nil
}

func (s *IncidentService) UnlinkAlert(ctx context.Context, incidentID, alertID, userID uuid.UUID) error {
	if err := s.incidentRepo.UnlinkAlert(ctx, incidentID, alertID); err != nil {
		return fmt.Errorf("failed to unlink alert: %w", err)
	}

	// Add timeline event
	timelineEvent := &domain.IncidentTimelineEvent{
		ID:          uuid.New(),
		IncidentID:  incidentID,
		EventType:   domain.TimelineEventAlertUnlinked,
		UserID:      &userID,
		Description: "Alert unlinked from incident",
		Metadata: map[string]interface{}{
			"alert_id": alertID.String(),
		},
	}

	if err := s.incidentRepo.AddTimelineEvent(ctx, timelineEvent); err != nil {
		fmt.Printf("Failed to add timeline event: %v\n", err)
	}

	return nil
}

func (s *IncidentService) ListAlerts(ctx context.Context, incidentID uuid.UUID) ([]*domain.IncidentAlertWithDetails, error) {
	alerts, err := s.incidentRepo.ListAlerts(ctx, incidentID)
	if err != nil {
		return nil, fmt.Errorf("failed to list alerts: %w", err)
	}

	return alerts, nil
}
