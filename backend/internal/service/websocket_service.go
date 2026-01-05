package service

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/nmn3m/pulsar/backend/internal/domain"
	"go.uber.org/zap"
)

// WebSocketService manages WebSocket connections and message broadcasting
type WebSocketService struct {
	hub    *domain.WSHub
	logger *zap.Logger
	mu     sync.RWMutex
}

// NewWebSocketService creates a new WebSocket service
func NewWebSocketService(logger *zap.Logger) *WebSocketService {
	return &WebSocketService{
		hub:    domain.NewWSHub(),
		logger: logger,
	}
}

// GetHub returns the WebSocket hub
func (s *WebSocketService) GetHub() *domain.WSHub {
	return s.hub
}

// Run starts the WebSocket hub
func (s *WebSocketService) Run() {
	for {
		select {
		case client := <-s.hub.Register:
			s.registerClient(client)

		case client := <-s.hub.Unregister:
			s.unregisterClient(client)

		case message := <-s.hub.Broadcast:
			s.broadcastMessage(message)
		}
	}
}

// registerClient registers a new WebSocket client
func (s *WebSocketService) registerClient(client *domain.WSClient) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.hub.Clients[client.OrganizationID] == nil {
		s.hub.Clients[client.OrganizationID] = make(map[*domain.WSClient]bool)
	}

	s.hub.Clients[client.OrganizationID][client] = true

	s.logger.Info("WebSocket client registered",
		zap.String("client_id", client.ID.String()),
		zap.String("user_id", client.UserID.String()),
		zap.String("org_id", client.OrganizationID.String()),
	)

	// Send welcome message
	welcomeMsg := domain.NewWSMessage(
		domain.WSEventConnected,
		client.OrganizationID,
		map[string]interface{}{
			"client_id": client.ID.String(),
			"user_id":   client.UserID.String(),
			"message":   "Connected to WebSocket server",
		},
	)

	select {
	case client.Send <- welcomeMsg:
	default:
		close(client.Send)
		delete(s.hub.Clients[client.OrganizationID], client)
	}
}

// unregisterClient unregisters a WebSocket client
func (s *WebSocketService) unregisterClient(client *domain.WSClient) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if clients, ok := s.hub.Clients[client.OrganizationID]; ok {
		if _, exists := clients[client]; exists {
			delete(clients, client)
			close(client.Send)

			s.logger.Info("WebSocket client unregistered",
				zap.String("client_id", client.ID.String()),
				zap.String("user_id", client.UserID.String()),
			)

			// Clean up empty organization maps
			if len(clients) == 0 {
				delete(s.hub.Clients, client.OrganizationID)
			}
		}
	}
}

// broadcastMessage broadcasts a message to all clients in an organization
func (s *WebSocketService) broadcastMessage(message *domain.WSMessage) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	clients, ok := s.hub.Clients[message.OrganizationID]
	if !ok {
		s.logger.Debug("No clients to broadcast to",
			zap.String("org_id", message.OrganizationID.String()),
			zap.String("event_type", string(message.Type)),
		)
		return
	}

	messageJSON, err := json.Marshal(message)
	if err != nil {
		s.logger.Error("Failed to marshal WebSocket message",
			zap.Error(err),
			zap.String("event_type", string(message.Type)),
		)
		return
	}

	s.logger.Debug("Broadcasting WebSocket message",
		zap.String("org_id", message.OrganizationID.String()),
		zap.String("event_type", string(message.Type)),
		zap.Int("client_count", len(clients)),
	)

	for client := range clients {
		select {
		case client.Send <- message:
			// Message sent successfully
		default:
			// Client's send channel is full, close and remove
			close(client.Send)
			delete(s.hub.Clients[message.OrganizationID], client)
			s.logger.Warn("Client send channel full, disconnecting",
				zap.String("client_id", client.ID.String()),
			)
		}
	}
}

// BroadcastAlertEvent broadcasts an alert-related event
func (s *WebSocketService) BroadcastAlertEvent(eventType domain.WSEventType, orgID uuid.UUID, alert *domain.Alert) {
	payload := map[string]interface{}{
		"alert_id":  alert.ID.String(),
		"message":   alert.Message,
		"priority":  alert.Priority,
		"status":    alert.Status,
		"source":    alert.Source,
		"created_at": alert.CreatedAt,
	}

	message := domain.NewWSMessage(eventType, orgID, payload)
	s.hub.Broadcast <- message
}

// BroadcastIncidentEvent broadcasts an incident-related event
func (s *WebSocketService) BroadcastIncidentEvent(eventType domain.WSEventType, orgID uuid.UUID, incident *domain.Incident) {
	payload := map[string]interface{}{
		"incident_id": incident.ID.String(),
		"title":       incident.Title,
		"severity":    incident.Severity,
		"status":      incident.Status,
		"priority":    incident.Priority,
		"started_at":  incident.StartedAt,
	}

	message := domain.NewWSMessage(eventType, orgID, payload)
	s.hub.Broadcast <- message
}

// BroadcastIncidentTimelineEvent broadcasts an incident timeline event
func (s *WebSocketService) BroadcastIncidentTimelineEvent(orgID, incidentID uuid.UUID, event *domain.IncidentTimelineEvent) {
	payload := map[string]interface{}{
		"incident_id": incidentID.String(),
		"event_id":    event.ID.String(),
		"event_type":  event.EventType,
		"description": event.Description,
		"created_at":  event.CreatedAt,
	}

	if event.UserID != nil {
		payload["user_id"] = event.UserID.String()
	}

	message := domain.NewWSMessage(domain.WSEventIncidentTimelineAdded, orgID, payload)
	s.hub.Broadcast <- message
}

// GetClientCount returns the number of connected clients for an organization
func (s *WebSocketService) GetClientCount(orgID uuid.UUID) int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if clients, ok := s.hub.Clients[orgID]; ok {
		return len(clients)
	}
	return 0
}

// GetTotalClientCount returns the total number of connected clients
func (s *WebSocketService) GetTotalClientCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	total := 0
	for _, clients := range s.hub.Clients {
		total += len(clients)
	}
	return total
}
