package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/nmn3m/pulsar/backend/internal/domain"
)

type IncidentRepository interface {
	Create(ctx context.Context, incident *domain.Incident) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Incident, error)
	Update(ctx context.Context, incident *domain.Incident) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filter *domain.IncidentFilter) ([]*domain.Incident, int, error)
	AddResponder(ctx context.Context, responder *domain.IncidentResponder) error
	RemoveResponder(ctx context.Context, incidentID, userID uuid.UUID) error
	UpdateResponderRole(ctx context.Context, incidentID, userID uuid.UUID, role domain.ResponderRole) error
	ListResponders(ctx context.Context, incidentID uuid.UUID) ([]*domain.ResponderWithUser, error)
	AddTimelineEvent(ctx context.Context, event *domain.IncidentTimelineEvent) error
	GetTimeline(ctx context.Context, incidentID uuid.UUID) ([]*domain.TimelineEventWithUser, error)
	LinkAlert(ctx context.Context, link *domain.IncidentAlert) error
	UnlinkAlert(ctx context.Context, incidentID, alertID uuid.UUID) error
	ListAlerts(ctx context.Context, incidentID uuid.UUID) ([]*domain.IncidentAlertWithDetails, error)
	GetWithDetails(ctx context.Context, id uuid.UUID) (*domain.IncidentWithDetails, error)
}
