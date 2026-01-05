package rest

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nmn3m/pulsar/backend/internal/domain"
	"github.com/nmn3m/pulsar/backend/internal/middleware"
	"github.com/nmn3m/pulsar/backend/internal/service"
)

type WebhookHandler struct {
	webhookService *service.WebhookService
}

func NewWebhookHandler(webhookService *service.WebhookService) *WebhookHandler {
	return &WebhookHandler{
		webhookService: webhookService,
	}
}

// Outgoing Webhook Endpoints

func (h *WebhookHandler) CreateEndpoint(c *gin.Context) {
	orgID, _ := middleware.GetOrganizationID(c)

	var req domain.CreateWebhookEndpointRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	endpoint, err := h.webhookService.CreateEndpoint(c.Request.Context(), orgID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create webhook endpoint"})
		return
	}

	c.JSON(http.StatusCreated, endpoint)
}

func (h *WebhookHandler) GetEndpoint(c *gin.Context) {
	orgID, _ := middleware.GetOrganizationID(c)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid endpoint ID"})
		return
	}

	endpoint, err := h.webhookService.GetEndpoint(c.Request.Context(), id)
	if err != nil {
		if err == domain.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Webhook endpoint not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get webhook endpoint"})
		return
	}

	if endpoint.OrganizationID != orgID {
		c.JSON(http.StatusNotFound, gin.H{"error": "Webhook endpoint not found"})
		return
	}

	c.JSON(http.StatusOK, endpoint)
}

func (h *WebhookHandler) ListEndpoints(c *gin.Context) {
	orgID, _ := middleware.GetOrganizationID(c)

	endpoints, err := h.webhookService.ListEndpoints(c.Request.Context(), orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list webhook endpoints"})
		return
	}

	c.JSON(http.StatusOK, endpoints)
}

func (h *WebhookHandler) UpdateEndpoint(c *gin.Context) {
	orgID, _ := middleware.GetOrganizationID(c)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid endpoint ID"})
		return
	}

	var req domain.UpdateWebhookEndpointRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	endpoint, err := h.webhookService.UpdateEndpoint(c.Request.Context(), id, orgID, &req)
	if err != nil {
		if err == domain.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Webhook endpoint not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update webhook endpoint"})
		return
	}

	c.JSON(http.StatusOK, endpoint)
}

func (h *WebhookHandler) DeleteEndpoint(c *gin.Context) {
	orgID, _ := middleware.GetOrganizationID(c)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid endpoint ID"})
		return
	}

	if err := h.webhookService.DeleteEndpoint(c.Request.Context(), id, orgID); err != nil {
		if err == domain.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Webhook endpoint not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete webhook endpoint"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Webhook endpoint deleted successfully"})
}

// Webhook Deliveries

func (h *WebhookHandler) ListDeliveries(c *gin.Context) {
	orgID, _ := middleware.GetOrganizationID(c)

	limit := 20
	offset := 0

	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil {
			limit = parsedLimit
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil {
			offset = parsedOffset
		}
	}

	deliveries, err := h.webhookService.ListDeliveries(c.Request.Context(), orgID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list webhook deliveries"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"deliveries": deliveries,
		"limit":      limit,
		"offset":     offset,
	})
}

// Incoming Webhook Tokens

func (h *WebhookHandler) CreateIncomingToken(c *gin.Context) {
	orgID, _ := middleware.GetOrganizationID(c)

	var req domain.CreateIncomingWebhookTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.webhookService.CreateIncomingToken(c.Request.Context(), orgID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create incoming webhook token"})
		return
	}

	c.JSON(http.StatusCreated, token)
}

func (h *WebhookHandler) ListIncomingTokens(c *gin.Context) {
	orgID, _ := middleware.GetOrganizationID(c)

	tokens, err := h.webhookService.ListIncomingTokens(c.Request.Context(), orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list incoming webhook tokens"})
		return
	}

	c.JSON(http.StatusOK, tokens)
}

func (h *WebhookHandler) DeleteIncomingToken(c *gin.Context) {
	orgID, _ := middleware.GetOrganizationID(c)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token ID"})
		return
	}

	if err := h.webhookService.DeleteIncomingToken(c.Request.Context(), id, orgID); err != nil {
		if err == domain.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Incoming webhook token not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete incoming webhook token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Incoming webhook token deleted successfully"})
}
