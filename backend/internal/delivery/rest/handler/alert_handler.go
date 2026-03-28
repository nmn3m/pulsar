package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/nmn3m/pulsar/backend/internal/delivery/rest/middleware"
	"github.com/nmn3m/pulsar/backend/internal/usecase"
)

type AlertHandler struct {
	alertUsecase *usecase.AlertUsecase
}

func NewAlertHandler(alertUsecase *usecase.AlertUsecase) *AlertHandler {
	return &AlertHandler{
		alertUsecase: alertUsecase,
	}
}

// Create godoc
// @Summary      Create a new alert
// @Description  Create a new alert in the organization
// @Tags         Alerts
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body usecase.CreateAlertRequest true "Create alert request"
// @Success      201 {object} domain.Alert
// @Failure      400 {object} map[string]string
// @Failure      401 {object} map[string]string
// @Router       /alerts [post]
func (h *AlertHandler) Create(c *gin.Context) {
	orgID, ok := middleware.GetOrganizationID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req usecase.CreateAlertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	alert, err := h.alertUsecase.CreateAlert(c.Request.Context(), orgID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, alert)
}

// Get godoc
// @Summary      Get an alert
// @Description  Get an alert by ID
// @Tags         Alerts
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Alert ID" format(uuid)
// @Success      200 {object} domain.Alert
// @Failure      400 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Router       /alerts/{id} [get]
func (h *AlertHandler) Get(c *gin.Context) {
	orgID, ok := middleware.GetOrganizationID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid alert ID"})
		return
	}

	alert, err := h.alertUsecase.GetAlert(c.Request.Context(), id, orgID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, alert)
}

// Update godoc
// @Summary      Update an alert
// @Description  Update an alert by ID
// @Tags         Alerts
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Alert ID" format(uuid)
// @Param        request body usecase.UpdateAlertRequest true "Update alert request"
// @Success      200 {object} domain.Alert
// @Failure      400 {object} map[string]string
// @Router       /alerts/{id} [patch]
func (h *AlertHandler) Update(c *gin.Context) {
	orgID, ok := middleware.GetOrganizationID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid alert ID"})
		return
	}

	var req usecase.UpdateAlertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	alert, err := h.alertUsecase.UpdateAlert(c.Request.Context(), id, orgID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, alert)
}

// Delete godoc
// @Summary      Delete an alert
// @Description  Delete an alert by ID
// @Tags         Alerts
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Alert ID" format(uuid)
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Router       /alerts/{id} [delete]
func (h *AlertHandler) Delete(c *gin.Context) {
	orgID, ok := middleware.GetOrganizationID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid alert ID"})
		return
	}

	if err := h.alertUsecase.DeleteAlert(c.Request.Context(), id, orgID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "alert deleted successfully"})
}

// List godoc
// @Summary      List alerts
// @Description  List alerts with optional filters
// @Tags         Alerts
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        status query []string false "Filter by status" collectionFormat(multi)
// @Param        priority query []string false "Filter by priority" collectionFormat(multi)
// @Param        assigned_to_user query string false "Filter by assigned user ID" format(uuid)
// @Param        assigned_to_team query string false "Filter by assigned team ID" format(uuid)
// @Param        source query string false "Filter by source"
// @Param        search query string false "Search in message and description"
// @Param        page query int false "Page number" default(1)
// @Param        page_size query int false "Page size" default(20)
// @Success      200 {object} usecase.ListAlertsResponse
// @Failure      400 {object} map[string]string
// @Failure      401 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /alerts [get]
func (h *AlertHandler) List(c *gin.Context) {
	orgID, ok := middleware.GetOrganizationID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req usecase.ListAlertsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.alertUsecase.ListAlerts(c.Request.Context(), orgID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// Acknowledge godoc
// @Summary      Acknowledge an alert
// @Description  Acknowledge an alert by ID
// @Tags         Alerts
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Alert ID" format(uuid)
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Failure      401 {object} map[string]string
// @Router       /alerts/{id}/acknowledge [post]
func (h *AlertHandler) Acknowledge(c *gin.Context) {
	orgID, ok := middleware.GetOrganizationID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid alert ID"})
		return
	}

	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := h.alertUsecase.AcknowledgeAlert(c.Request.Context(), id, orgID, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "alert acknowledged successfully"})
}

// Close godoc
// @Summary      Close an alert
// @Description  Close an alert by ID with a reason
// @Tags         Alerts
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Alert ID" format(uuid)
// @Param        request body usecase.CloseAlertRequest true "Close alert request"
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Failure      401 {object} map[string]string
// @Router       /alerts/{id}/close [post]
func (h *AlertHandler) Close(c *gin.Context) {
	orgID, ok := middleware.GetOrganizationID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid alert ID"})
		return
	}

	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req usecase.CloseAlertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.alertUsecase.CloseAlert(c.Request.Context(), id, orgID, userID, req.Reason); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "alert closed successfully"})
}

// Snooze godoc
// @Summary      Snooze an alert
// @Description  Snooze an alert until a specified time
// @Tags         Alerts
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Alert ID" format(uuid)
// @Param        request body usecase.SnoozeAlertRequest true "Snooze alert request"
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Router       /alerts/{id}/snooze [post]
func (h *AlertHandler) Snooze(c *gin.Context) {
	orgID, ok := middleware.GetOrganizationID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid alert ID"})
		return
	}

	var req usecase.SnoozeAlertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.alertUsecase.SnoozeAlert(c.Request.Context(), id, orgID, req.Until); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "alert snoozed successfully"})
}

// Assign godoc
// @Summary      Assign an alert
// @Description  Assign an alert to a user or team
// @Tags         Alerts
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Alert ID" format(uuid)
// @Param        request body usecase.AssignAlertRequest true "Assign alert request"
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Router       /alerts/{id}/assign [post]
func (h *AlertHandler) Assign(c *gin.Context) {
	orgID, ok := middleware.GetOrganizationID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid alert ID"})
		return
	}

	var req usecase.AssignAlertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.alertUsecase.AssignAlert(c.Request.Context(), id, orgID, req.UserID, req.TeamID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "alert assigned successfully"})
}
