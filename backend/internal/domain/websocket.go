package domain

import (
	"time"

	"github.com/google/uuid"
)

// WSEventType represents the type of WebSocket event
type WSEventType string

const (
	// Alert events
	WSEventAlertCreated      WSEventType = "alert.created"
	WSEventAlertUpdated      WSEventType = "alert.updated"
	WSEventAlertDeleted      WSEventType = "alert.deleted"
	WSEventAlertAcknowledged WSEventType = "alert.acknowledged"
	WSEventAlertClosed       WSEventType = "alert.closed"
	WSEventAlertEscalated    WSEventType = "alert.escalated"

	// Incident events
	WSEventIncidentCreated          WSEventType = "incident.created"
	WSEventIncidentUpdated          WSEventType = "incident.updated"
	WSEventIncidentDeleted          WSEventType = "incident.deleted"
	WSEventIncidentTimelineAdded    WSEventType = "incident.timeline_added"
	WSEventIncidentResponderAdded   WSEventType = "incident.responder_added"
	WSEventIncidentResponderRemoved WSEventType = "incident.responder_removed"
	WSEventIncidentAlertLinked      WSEventType = "incident.alert_linked"
	WSEventIncidentAlertUnlinked    WSEventType = "incident.alert_unlinked"

	// Connection events
	WSEventConnected WSEventType = "connection.connected"
	WSEventError     WSEventType = "connection.error"
	WSEventPing      WSEventType = "connection.ping"
	WSEventPong      WSEventType = "connection.pong"
)

// WSMessage represents a WebSocket message
type WSMessage struct {
	ID             string                 `json:"id"`
	Type           WSEventType            `json:"type"`
	OrganizationID uuid.UUID              `json:"organization_id"`
	Payload        map[string]interface{} `json:"payload"`
	Timestamp      time.Time              `json:"timestamp"`
}

// NewWSMessage creates a new WebSocket message
func NewWSMessage(eventType WSEventType, orgID uuid.UUID, payload map[string]interface{}) *WSMessage {
	return &WSMessage{
		ID:             uuid.New().String(),
		Type:           eventType,
		OrganizationID: orgID,
		Payload:        payload,
		Timestamp:      time.Now(),
	}
}

// WSClient represents a WebSocket client connection
type WSClient struct {
	ID             uuid.UUID
	UserID         uuid.UUID
	OrganizationID uuid.UUID
	Send           chan *WSMessage
}

// WSHub manages WebSocket client connections
type WSHub struct {
	// Registered clients by organization
	Clients map[uuid.UUID]map[*WSClient]bool

	// Inbound messages from clients
	Broadcast chan *WSMessage

	// Register requests from clients
	Register chan *WSClient

	// Unregister requests from clients
	Unregister chan *WSClient
}

// NewWSHub creates a new WebSocket hub
func NewWSHub() *WSHub {
	return &WSHub{
		Clients:    make(map[uuid.UUID]map[*WSClient]bool),
		Broadcast:  make(chan *WSMessage, 256),
		Register:   make(chan *WSClient),
		Unregister: make(chan *WSClient),
	}
}
