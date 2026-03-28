package outbound

import (
	"github.com/google/uuid"

	"github.com/nmn3m/pulsar/backend/internal/core/domain"
)

type EventBroadcaster interface {
	BroadcastAlertEvent(eventType domain.WSEventType, orgID uuid.UUID, alert *domain.Alert)
	BroadcastIncidentEvent(eventType domain.WSEventType, orgID uuid.UUID, incident *domain.Incident)
	BroadcastIncidentTimelineEvent(orgID, incidentID uuid.UUID, event *domain.IncidentTimelineEvent)
}
