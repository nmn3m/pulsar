package inbound

import (
	"context"

	"github.com/google/uuid"

	"github.com/nmn3m/pulsar/backend/internal/core/domain"
	"github.com/nmn3m/pulsar/backend/internal/core/dto"
)

type IncidentService interface {
	CreateIncident(ctx context.Context, orgID, userID uuid.UUID, req *dto.CreateIncidentRequest) (*domain.Incident, error)
	GetIncident(ctx context.Context, id, orgID uuid.UUID) (*domain.Incident, error)
	GetIncidentWithDetails(ctx context.Context, id, orgID uuid.UUID) (*domain.IncidentWithDetails, error)
	UpdateIncident(ctx context.Context, id, orgID, userID uuid.UUID, req *dto.UpdateIncidentRequest) (*domain.Incident, error)
	DeleteIncident(ctx context.Context, id, orgID uuid.UUID) error
	ListIncidents(ctx context.Context, orgID uuid.UUID, req *dto.ListIncidentsRequest) (*dto.ListIncidentsResponse, error)
	AddResponder(ctx context.Context, incidentID uuid.UUID, userID uuid.UUID, req *dto.AddResponderRequest) (*domain.IncidentResponder, error)
	RemoveResponder(ctx context.Context, incidentID, orgID, responderUserID, actionUserID uuid.UUID) error
	UpdateResponderRole(ctx context.Context, incidentID, orgID, responderUserID uuid.UUID, req *dto.UpdateResponderRoleRequest) error
	ListResponders(ctx context.Context, incidentID, orgID uuid.UUID) ([]*domain.ResponderWithUser, error)
	AddNote(ctx context.Context, incidentID, orgID, userID uuid.UUID, req *dto.AddNoteRequest) (*domain.IncidentTimelineEvent, error)
	GetTimeline(ctx context.Context, incidentID, orgID uuid.UUID) ([]*domain.TimelineEventWithUser, error)
	LinkAlert(ctx context.Context, incidentID, orgID, userID uuid.UUID, req *dto.LinkAlertRequest) (*domain.IncidentAlert, error)
	UnlinkAlert(ctx context.Context, incidentID, orgID, alertID, userID uuid.UUID) error
	ListAlerts(ctx context.Context, incidentID, orgID uuid.UUID) ([]*domain.IncidentAlertWithDetails, error)
}
