package rest

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/nmn3m/pulsar/backend/internal/domain"
	"github.com/nmn3m/pulsar/backend/internal/service"
)

type DNDHandler struct {
	dndService *service.DNDService
}

func NewDNDHandler(dndService *service.DNDService) *DNDHandler {
	return &DNDHandler{dndService: dndService}
}

// GetDNDSettings godoc
// @Summary Get user's DND settings
// @Description Get the current user's Do Not Disturb settings
// @Tags dnd
// @Accept json
// @Produce json
// @Success 200 {object} domain.UserDNDSettings
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/me/dnd [get]
// @Security BearerAuth
func (h *DNDHandler) GetDNDSettings(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	settings, err := h.dndService.GetSettings(c.Request.Context(), userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, settings)
}

// UpdateDNDSettings godoc
// @Summary Update user's DND settings
// @Description Update the current user's Do Not Disturb settings
// @Tags dnd
// @Accept json
// @Produce json
// @Param request body domain.UpdateDNDSettingsRequest true "DND settings"
// @Success 200 {object} domain.UserDNDSettings
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/me/dnd [put]
// @Security BearerAuth
func (h *DNDHandler) UpdateDNDSettings(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req domain.UpdateDNDSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	settings, err := h.dndService.UpdateSettings(c.Request.Context(), userID.(uuid.UUID), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, settings)
}

// AddDNDOverride godoc
// @Summary Add a temporary DND override
// @Description Add a temporary Do Not Disturb period (e.g., vacation)
// @Tags dnd
// @Accept json
// @Produce json
// @Param request body domain.AddDNDOverrideRequest true "Override details"
// @Success 200 {object} domain.UserDNDSettings
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/me/dnd/overrides [post]
// @Security BearerAuth
func (h *DNDHandler) AddDNDOverride(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req domain.AddDNDOverrideRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate that end is after start
	if !req.End.After(req.Start) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "end time must be after start time"})
		return
	}

	settings, err := h.dndService.AddOverride(c.Request.Context(), userID.(uuid.UUID), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, settings)
}

// RemoveDNDOverride godoc
// @Summary Remove a DND override
// @Description Remove a temporary Do Not Disturb period by index
// @Tags dnd
// @Accept json
// @Produce json
// @Param index path int true "Override index"
// @Success 200 {object} domain.UserDNDSettings
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/me/dnd/overrides/{index} [delete]
// @Security BearerAuth
func (h *DNDHandler) RemoveDNDOverride(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	indexStr := c.Param("index")
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid index"})
		return
	}

	settings, err := h.dndService.RemoveOverride(c.Request.Context(), userID.(uuid.UUID), index)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, settings)
}

// DeleteDNDSettings godoc
// @Summary Delete user's DND settings
// @Description Delete all DND settings for the current user
// @Tags dnd
// @Accept json
// @Produce json
// @Success 204 "No Content"
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/me/dnd [delete]
// @Security BearerAuth
func (h *DNDHandler) DeleteDNDSettings(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	err := h.dndService.DeleteSettings(c.Request.Context(), userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// CheckDNDStatus godoc
// @Summary Check if user is currently in DND mode
// @Description Check if the current user is in Do Not Disturb mode
// @Tags dnd
// @Accept json
// @Produce json
// @Param priority query string false "Alert priority to check against (P1, P2, P3, P4, P5)"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/me/dnd/status [get]
// @Security BearerAuth
func (h *DNDHandler) CheckDNDStatus(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Default to P3 if no priority specified
	priority := domain.PriorityP3
	if p := c.Query("priority"); p != "" {
		priority = domain.AlertPriority(p)
	}

	inDND, err := h.dndService.IsInDNDMode(c.Request.Context(), userID.(uuid.UUID), priority)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"in_dnd_mode": inDND,
		"priority":    priority,
	})
}
