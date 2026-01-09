package rest

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"

	"github.com/nmn3m/pulsar/backend/internal/domain"
	"github.com/nmn3m/pulsar/backend/internal/middleware"
	"github.com/nmn3m/pulsar/backend/internal/service"
)

const (
	// Time allowed to write a message to the peer
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer
	pongWait = 60 * time.Second

	// Send pings to peer with this period (must be less than pongWait)
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// TODO: In production, validate origin properly
		return true
	},
}

type WebSocketHandler struct {
	wsService *service.WebSocketService
	logger    *zap.Logger
}

func NewWebSocketHandler(wsService *service.WebSocketService, logger *zap.Logger) *WebSocketHandler {
	return &WebSocketHandler{
		wsService: wsService,
		logger:    logger,
	}
}

// HandleWebSocket handles WebSocket connections
func (h *WebSocketHandler) HandleWebSocket(c *gin.Context) {
	// Get user and organization from context (set by auth middleware)
	userID, _ := middleware.GetUserID(c)
	orgID, _ := middleware.GetOrganizationID(c)

	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		h.logger.Error("Failed to upgrade to WebSocket", zap.Error(err))
		return
	}

	// Create new client
	client := &domain.WSClient{
		ID:             uuid.New(),
		UserID:         userID,
		OrganizationID: orgID,
		Send:           make(chan *domain.WSMessage, 256),
	}

	// Register client with hub
	h.wsService.GetHub().Register <- client

	// Start goroutines for reading and writing
	go h.writePump(conn, client)
	go h.readPump(conn, client)
}

// readPump pumps messages from the WebSocket connection to the hub
func (h *WebSocketHandler) readPump(conn *websocket.Conn, client *domain.WSClient) {
	defer func() {
		h.wsService.GetHub().Unregister <- client
		conn.Close()
	}()

	conn.SetReadLimit(maxMessageSize)
	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				h.logger.Error("WebSocket read error", zap.Error(err))
			}
			break
		}

		// Handle incoming messages if needed
		h.handleIncomingMessage(client, message)
	}
}

// writePump pumps messages from the hub to the WebSocket connection
func (h *WebSocketHandler) writePump(conn *websocket.Conn, client *domain.WSClient) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		conn.Close()
	}()

	for {
		select {
		case message, ok := <-client.Send:
			conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// Channel was closed
				conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// Write message as JSON
			if err := conn.WriteJSON(message); err != nil {
				h.logger.Error("Failed to write WebSocket message", zap.Error(err))
				return
			}

		case <-ticker.C:
			conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// handleIncomingMessage handles messages received from clients
func (h *WebSocketHandler) handleIncomingMessage(client *domain.WSClient, data []byte) {
	var msg map[string]interface{}
	if err := json.Unmarshal(data, &msg); err != nil {
		h.logger.Error("Failed to unmarshal incoming message", zap.Error(err))
		return
	}

	// Handle pong messages
	if msgType, ok := msg["type"].(string); ok {
		if msgType == string(domain.WSEventPong) {
			h.logger.Debug("Received pong from client",
				zap.String("client_id", client.ID.String()),
			)
		}
	}

	// Add more message type handling here if needed
}

// GetStats returns WebSocket connection statistics
func (h *WebSocketHandler) GetStats(c *gin.Context) {
	orgID, _ := middleware.GetOrganizationID(c)

	stats := gin.H{
		"organization_client_count": h.wsService.GetClientCount(orgID),
		"total_client_count":        h.wsService.GetTotalClientCount(),
	}

	c.JSON(http.StatusOK, stats)
}
