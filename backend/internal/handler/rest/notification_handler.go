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

type NotificationHandler struct {
	notificationService *service.NotificationService
}

func NewNotificationHandler(notificationService *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{
		notificationService: notificationService,
	}
}

// ==================== Notification Channels ====================

// CreateChannel godoc
// @Summary      Create a notification channel
// @Description  Creates a new notification channel for the organization
// @Tags         Notifications
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body domain.CreateNotificationChannelRequest true "Channel creation request"
// @Success      201 {object} domain.NotificationChannel
// @Failure      400 {object} map[string]string
// @Failure      401 {object} map[string]string
// @Router       /notifications/channels [post]
func (h *NotificationHandler) CreateChannel(c *gin.Context) {
	orgID, ok := middleware.GetOrganizationID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req domain.CreateNotificationChannelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	channel, err := h.notificationService.CreateChannel(c.Request.Context(), orgID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, channel)
}

// GetChannel godoc
// @Summary      Get a notification channel
// @Description  Retrieves a notification channel by its ID
// @Tags         Notifications
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Channel ID" format(uuid)
// @Success      200 {object} domain.NotificationChannel
// @Failure      400 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Router       /notifications/channels/{id} [get]
func (h *NotificationHandler) GetChannel(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid channel ID"})
		return
	}

	channel, err := h.notificationService.GetChannel(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, channel)
}

// ListChannels godoc
// @Summary      List notification channels
// @Description  Lists all notification channels for the organization
// @Tags         Notifications
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} map[string]interface{}
// @Failure      401 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /notifications/channels [get]
func (h *NotificationHandler) ListChannels(c *gin.Context) {
	orgID, ok := middleware.GetOrganizationID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	channels, err := h.notificationService.ListChannels(c.Request.Context(), orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"channels": channels,
		"total":    len(channels),
	})
}

// UpdateChannel godoc
// @Summary      Update a notification channel
// @Description  Updates an existing notification channel by its ID
// @Tags         Notifications
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Channel ID" format(uuid)
// @Param        request body domain.UpdateNotificationChannelRequest true "Channel update request"
// @Success      200 {object} domain.NotificationChannel
// @Failure      400 {object} map[string]string
// @Router       /notifications/channels/{id} [patch]
func (h *NotificationHandler) UpdateChannel(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid channel ID"})
		return
	}

	var req domain.UpdateNotificationChannelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	channel, err := h.notificationService.UpdateChannel(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, channel)
}

// DeleteChannel godoc
// @Summary      Delete a notification channel
// @Description  Deletes a notification channel by its ID
// @Tags         Notifications
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Channel ID" format(uuid)
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Router       /notifications/channels/{id} [delete]
func (h *NotificationHandler) DeleteChannel(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid channel ID"})
		return
	}

	if err := h.notificationService.DeleteChannel(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "notification channel deleted successfully"})
}

// ==================== User Notification Preferences ====================

// CreatePreference godoc
// @Summary      Create a user notification preference
// @Description  Creates a new notification preference for the authenticated user
// @Tags         Notifications
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body domain.CreateUserNotificationPreferenceRequest true "Preference creation request"
// @Success      201 {object} domain.UserNotificationPreference
// @Failure      400 {object} map[string]string
// @Failure      401 {object} map[string]string
// @Router       /notifications/preferences [post]
func (h *NotificationHandler) CreatePreference(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req domain.CreateUserNotificationPreferenceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pref, err := h.notificationService.CreatePreference(c.Request.Context(), userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, pref)
}

// GetPreference godoc
// @Summary      Get a user notification preference
// @Description  Retrieves a notification preference by its ID
// @Tags         Notifications
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Preference ID" format(uuid)
// @Success      200 {object} domain.UserNotificationPreference
// @Failure      400 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Router       /notifications/preferences/{id} [get]
func (h *NotificationHandler) GetPreference(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid preference ID"})
		return
	}

	pref, err := h.notificationService.GetPreference(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pref)
}

// ListUserPreferences godoc
// @Summary      List user notification preferences
// @Description  Lists all notification preferences for the authenticated user
// @Tags         Notifications
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} map[string]interface{}
// @Failure      401 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /notifications/preferences [get]
func (h *NotificationHandler) ListUserPreferences(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	prefs, err := h.notificationService.ListUserPreferences(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"preferences": prefs,
		"total":       len(prefs),
	})
}

// UpdatePreference godoc
// @Summary      Update a user notification preference
// @Description  Updates an existing notification preference by its ID
// @Tags         Notifications
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Preference ID" format(uuid)
// @Param        request body domain.UpdateUserNotificationPreferenceRequest true "Preference update request"
// @Success      200 {object} domain.UserNotificationPreference
// @Failure      400 {object} map[string]string
// @Router       /notifications/preferences/{id} [patch]
func (h *NotificationHandler) UpdatePreference(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid preference ID"})
		return
	}

	var req domain.UpdateUserNotificationPreferenceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pref, err := h.notificationService.UpdatePreference(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pref)
}

// DeletePreference godoc
// @Summary      Delete a user notification preference
// @Description  Deletes a notification preference by its ID
// @Tags         Notifications
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Preference ID" format(uuid)
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Router       /notifications/preferences/{id} [delete]
func (h *NotificationHandler) DeletePreference(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid preference ID"})
		return
	}

	if err := h.notificationService.DeletePreference(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "notification preference deleted successfully"})
}

// ==================== Sending Notifications ====================

// SendNotification godoc
// @Summary      Send a notification
// @Description  Sends a notification through the specified channel
// @Tags         Notifications
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body domain.SendNotificationRequest true "Notification send request"
// @Success      201 {object} domain.NotificationLog
// @Failure      400 {object} map[string]string
// @Failure      401 {object} map[string]string
// @Router       /notifications/send [post]
func (h *NotificationHandler) SendNotification(c *gin.Context) {
	orgID, ok := middleware.GetOrganizationID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req domain.SendNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log, err := h.notificationService.SendNotification(c.Request.Context(), orgID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, log)
}

// ==================== Notification Logs ====================

// GetLog godoc
// @Summary      Get a notification log
// @Description  Retrieves a notification log entry by its ID
// @Tags         Notifications
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Log ID" format(uuid)
// @Success      200 {object} domain.NotificationLog
// @Failure      400 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Router       /notifications/logs/{id} [get]
func (h *NotificationHandler) GetLog(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid log ID"})
		return
	}

	log, err := h.notificationService.GetLog(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, log)
}

// ListLogs godoc
// @Summary      List notification logs
// @Description  Lists all notification logs for the organization with pagination
// @Tags         Notifications
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        limit query int false "Number of logs to return" default(50)
// @Param        offset query int false "Number of logs to skip" default(0)
// @Success      200 {object} map[string]interface{}
// @Failure      401 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /notifications/logs [get]
func (h *NotificationHandler) ListLogs(c *gin.Context) {
	orgID, ok := middleware.GetOrganizationID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	logs, err := h.notificationService.ListLogs(c.Request.Context(), orgID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"logs":   logs,
		"total":  len(logs),
		"limit":  limit,
		"offset": offset,
	})
}

// ListLogsByAlert godoc
// @Summary      List notification logs by alert
// @Description  Lists all notification logs associated with a specific alert
// @Tags         Notifications
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        alertId path string true "Alert ID" format(uuid)
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /notifications/logs/alert/{alertId} [get]
func (h *NotificationHandler) ListLogsByAlert(c *gin.Context) {
	alertIDStr := c.Param("alertId")
	alertID, err := uuid.Parse(alertIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid alert ID"})
		return
	}

	logs, err := h.notificationService.ListLogsByAlert(c.Request.Context(), alertID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"logs":  logs,
		"total": len(logs),
	})
}

// ListLogsByUser godoc
// @Summary      List notification logs for the current user
// @Description  Lists all notification logs for the authenticated user with pagination
// @Tags         Notifications
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        limit query int false "Number of logs to return" default(50)
// @Param        offset query int false "Number of logs to skip" default(0)
// @Success      200 {object} map[string]interface{}
// @Failure      401 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /notifications/logs/user/me [get]
func (h *NotificationHandler) ListLogsByUser(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	logs, err := h.notificationService.ListLogsByUser(c.Request.Context(), userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"logs":   logs,
		"total":  len(logs),
		"limit":  limit,
		"offset": offset,
	})
}
