package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/nmn3m/pulsar/backend/internal/domain"
	"github.com/nmn3m/pulsar/backend/internal/middleware"
	"github.com/nmn3m/pulsar/backend/internal/service"
)

type APIKeyHandler struct {
	apiKeyService *service.APIKeyService
}

func NewAPIKeyHandler(apiKeyService *service.APIKeyService) *APIKeyHandler {
	return &APIKeyHandler{
		apiKeyService: apiKeyService,
	}
}

// Create godoc
// @Summary      Create a new API key
// @Description  Create a new API key for programmatic access. The raw key is only shown once.
// @Tags         API Keys
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body domain.CreateAPIKeyRequest true "Create API key request"
// @Success      201 {object} domain.APIKeyResponse
// @Failure      400 {object} map[string]string
// @Failure      401 {object} map[string]string
// @Router       /api-keys [post]
func (h *APIKeyHandler) Create(c *gin.Context) {
	orgID, ok := middleware.GetOrganizationID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req domain.CreateAPIKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.apiKeyService.CreateAPIKey(c.Request.Context(), orgID, userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// List godoc
// @Summary      List API keys
// @Description  List all API keys for the current user
// @Tags         API Keys
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} map[string][]domain.APIKey
// @Failure      401 {object} map[string]string
// @Router       /api-keys [get]
func (h *APIKeyHandler) List(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	keys, err := h.apiKeyService.ListUserAPIKeys(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"api_keys": keys})
}

// ListAll godoc
// @Summary      List all API keys in organization
// @Description  List all API keys for the organization (admin only)
// @Tags         API Keys
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} map[string][]domain.APIKey
// @Failure      401 {object} map[string]string
// @Router       /api-keys/all [get]
func (h *APIKeyHandler) ListAll(c *gin.Context) {
	orgID, ok := middleware.GetOrganizationID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	keys, err := h.apiKeyService.ListAPIKeys(c.Request.Context(), orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"api_keys": keys})
}

// Get godoc
// @Summary      Get an API key
// @Description  Get an API key by ID
// @Tags         API Keys
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "API Key ID" format(uuid)
// @Success      200 {object} domain.APIKey
// @Failure      400 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Router       /api-keys/{id} [get]
func (h *APIKeyHandler) Get(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid API key ID"})
		return
	}

	key, err := h.apiKeyService.GetAPIKey(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "API key not found"})
		return
	}

	// Ensure user owns this key
	if key.UserID != userID {
		c.JSON(http.StatusNotFound, gin.H{"error": "API key not found"})
		return
	}

	c.JSON(http.StatusOK, key)
}

// Update godoc
// @Summary      Update an API key
// @Description  Update an API key by ID
// @Tags         API Keys
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "API Key ID" format(uuid)
// @Param        request body domain.UpdateAPIKeyRequest true "Update API key request"
// @Success      200 {object} domain.APIKey
// @Failure      400 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Router       /api-keys/{id} [patch]
func (h *APIKeyHandler) Update(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid API key ID"})
		return
	}

	// Check ownership first
	existingKey, err := h.apiKeyService.GetAPIKey(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "API key not found"})
		return
	}

	if existingKey.UserID != userID {
		c.JSON(http.StatusNotFound, gin.H{"error": "API key not found"})
		return
	}

	var req domain.UpdateAPIKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	key, err := h.apiKeyService.UpdateAPIKey(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, key)
}

// Revoke godoc
// @Summary      Revoke an API key
// @Description  Revoke (deactivate) an API key by ID
// @Tags         API Keys
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "API Key ID" format(uuid)
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Router       /api-keys/{id}/revoke [post]
func (h *APIKeyHandler) Revoke(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid API key ID"})
		return
	}

	// Check ownership first
	existingKey, err := h.apiKeyService.GetAPIKey(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "API key not found"})
		return
	}

	if existingKey.UserID != userID {
		c.JSON(http.StatusNotFound, gin.H{"error": "API key not found"})
		return
	}

	if err := h.apiKeyService.RevokeAPIKey(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "API key revoked"})
}

// Delete godoc
// @Summary      Delete an API key
// @Description  Permanently delete an API key by ID
// @Tags         API Keys
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "API Key ID" format(uuid)
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Router       /api-keys/{id} [delete]
func (h *APIKeyHandler) Delete(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid API key ID"})
		return
	}

	// Check ownership first
	existingKey, err := h.apiKeyService.GetAPIKey(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "API key not found"})
		return
	}

	if existingKey.UserID != userID {
		c.JSON(http.StatusNotFound, gin.H{"error": "API key not found"})
		return
	}

	if err := h.apiKeyService.DeleteAPIKey(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "API key deleted"})
}

// GetScopes godoc
// @Summary      Get available scopes
// @Description  Get a list of all available API key scopes
// @Tags         API Keys
// @Accept       json
// @Produce      json
// @Success      200 {object} map[string][]string
// @Router       /api-keys/scopes [get]
func (h *APIKeyHandler) GetScopes(c *gin.Context) {
	scopes := domain.ValidScopes()
	scopeStrings := make([]string, len(scopes))
	for i, s := range scopes {
		scopeStrings[i] = string(s)
	}
	c.JSON(http.StatusOK, gin.H{"scopes": scopeStrings})
}
