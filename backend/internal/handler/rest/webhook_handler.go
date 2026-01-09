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

// CreateEndpoint godoc
// @Summary      Create a webhook endpoint
// @Description  Create a new outgoing webhook endpoint
// @Tags         Webhooks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body domain.CreateWebhookEndpointRequest true "Create endpoint request"
// @Success      201 {object} domain.WebhookEndpoint
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /webhooks/endpoints [post]
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

// GetEndpoint godoc
// @Summary      Get a webhook endpoint
// @Description  Get a webhook endpoint by ID
// @Tags         Webhooks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Endpoint ID" format(uuid)
// @Success      200 {object} domain.WebhookEndpoint
// @Failure      400 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /webhooks/endpoints/{id} [get]
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

// ListEndpoints godoc
// @Summary      List webhook endpoints
// @Description  List all outgoing webhook endpoints for the organization
// @Tags         Webhooks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {array} domain.WebhookEndpoint
// @Failure      500 {object} map[string]string
// @Router       /webhooks/endpoints [get]
func (h *WebhookHandler) ListEndpoints(c *gin.Context) {
	orgID, _ := middleware.GetOrganizationID(c)

	endpoints, err := h.webhookService.ListEndpoints(c.Request.Context(), orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list webhook endpoints"})
		return
	}

	// Ensure we return an empty array instead of null
	if endpoints == nil {
		endpoints = []*domain.WebhookEndpoint{}
	}

	c.JSON(http.StatusOK, endpoints)
}

// UpdateEndpoint godoc
// @Summary      Update a webhook endpoint
// @Description  Update a webhook endpoint by ID
// @Tags         Webhooks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Endpoint ID" format(uuid)
// @Param        request body domain.UpdateWebhookEndpointRequest true "Update endpoint request"
// @Success      200 {object} domain.WebhookEndpoint
// @Failure      400 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /webhooks/endpoints/{id} [patch]
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

// DeleteEndpoint godoc
// @Summary      Delete a webhook endpoint
// @Description  Delete a webhook endpoint by ID
// @Tags         Webhooks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Endpoint ID" format(uuid)
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /webhooks/endpoints/{id} [delete]
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

// ListDeliveries godoc
// @Summary      List webhook deliveries
// @Description  List all webhook delivery attempts for the organization
// @Tags         Webhooks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        limit query int false "Limit" default(20)
// @Param        offset query int false "Offset" default(0)
// @Success      200 {object} map[string]interface{}
// @Failure      500 {object} map[string]string
// @Router       /webhooks/deliveries [get]
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

// CreateIncomingToken godoc
// @Summary      Create an incoming webhook token
// @Description  Create a new incoming webhook token for receiving webhooks from external sources
// @Tags         Webhooks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body domain.CreateIncomingWebhookTokenRequest true "Create token request"
// @Success      201 {object} domain.IncomingWebhookToken
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /webhooks/incoming [post]
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

// ListIncomingTokens godoc
// @Summary      List incoming webhook tokens
// @Description  List all incoming webhook tokens for the organization
// @Tags         Webhooks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {array} domain.IncomingWebhookToken
// @Failure      500 {object} map[string]string
// @Router       /webhooks/incoming [get]
func (h *WebhookHandler) ListIncomingTokens(c *gin.Context) {
	orgID, _ := middleware.GetOrganizationID(c)

	tokens, err := h.webhookService.ListIncomingTokens(c.Request.Context(), orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list incoming webhook tokens"})
		return
	}

	// Ensure we return an empty array instead of null
	if tokens == nil {
		tokens = []*domain.IncomingWebhookToken{}
	}

	c.JSON(http.StatusOK, tokens)
}

// DeleteIncomingToken godoc
// @Summary      Delete an incoming webhook token
// @Description  Delete an incoming webhook token by ID
// @Tags         Webhooks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Token ID" format(uuid)
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /webhooks/incoming/{id} [delete]
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
