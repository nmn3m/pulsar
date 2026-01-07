package rest

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nmn3m/pulsar/backend/internal/domain"
	"github.com/nmn3m/pulsar/backend/internal/service"
)

type MetricsHandler struct {
	metricsService *service.MetricsService
}

func NewMetricsHandler(metricsService *service.MetricsService) *MetricsHandler {
	return &MetricsHandler{
		metricsService: metricsService,
	}
}

// parseMetricsFilter extracts filter parameters from query string
func parseMetricsFilter(c *gin.Context) *domain.MetricsFilter {
	filter := &domain.MetricsFilter{}

	if startStr := c.Query("start_time"); startStr != "" {
		if t, err := time.Parse(time.RFC3339, startStr); err == nil {
			filter.StartTime = &t
		}
	}

	if endStr := c.Query("end_time"); endStr != "" {
		if t, err := time.Parse(time.RFC3339, endStr); err == nil {
			filter.EndTime = &t
		}
	}

	if teamID := c.Query("team_id"); teamID != "" {
		filter.TeamID = &teamID
	}

	if period := c.Query("period"); period != "" {
		filter.Period = period
	}

	return filter
}

// GetDashboard godoc
// @Summary Get dashboard metrics
// @Description Get aggregated metrics for the dashboard including alerts, incidents, notifications, and trends
// @Tags Metrics
// @Accept json
// @Produce json
// @Param start_time query string false "Start time (RFC3339 format)"
// @Param end_time query string false "End time (RFC3339 format)"
// @Param period query string false "Trend period: hourly, daily, weekly" default(daily)
// @Success 200 {object} domain.DashboardMetrics
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /metrics/dashboard [get]
func (h *MetricsHandler) GetDashboard(c *gin.Context) {
	orgID, exists := c.Get("organization_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Organization not found"})
		return
	}

	filter := parseMetricsFilter(c)

	metrics, err := h.metricsService.GetDashboardMetrics(c.Request.Context(), orgID.(uuid.UUID), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	metrics.UpdatedAt = time.Now()
	c.JSON(http.StatusOK, metrics)
}

// GetAlertMetrics godoc
// @Summary Get alert metrics
// @Description Get detailed alert metrics including counts by status, priority, source, and response times
// @Tags Metrics
// @Accept json
// @Produce json
// @Param start_time query string false "Start time (RFC3339 format)"
// @Param end_time query string false "End time (RFC3339 format)"
// @Success 200 {object} domain.AlertMetrics
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /metrics/alerts [get]
func (h *MetricsHandler) GetAlertMetrics(c *gin.Context) {
	orgID, exists := c.Get("organization_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Organization not found"})
		return
	}

	filter := parseMetricsFilter(c)

	metrics, err := h.metricsService.GetAlertMetrics(c.Request.Context(), orgID.(uuid.UUID), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, metrics)
}

// GetIncidentMetrics godoc
// @Summary Get incident metrics
// @Description Get detailed incident metrics including counts by status, severity, and resolution times
// @Tags Metrics
// @Accept json
// @Produce json
// @Param start_time query string false "Start time (RFC3339 format)"
// @Param end_time query string false "End time (RFC3339 format)"
// @Success 200 {object} domain.IncidentMetrics
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /metrics/incidents [get]
func (h *MetricsHandler) GetIncidentMetrics(c *gin.Context) {
	orgID, exists := c.Get("organization_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Organization not found"})
		return
	}

	filter := parseMetricsFilter(c)

	metrics, err := h.metricsService.GetIncidentMetrics(c.Request.Context(), orgID.(uuid.UUID), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, metrics)
}

// GetNotificationMetrics godoc
// @Summary Get notification metrics
// @Description Get notification delivery metrics including counts by status and channel
// @Tags Metrics
// @Accept json
// @Produce json
// @Param start_time query string false "Start time (RFC3339 format)"
// @Param end_time query string false "End time (RFC3339 format)"
// @Success 200 {object} domain.NotificationMetrics
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /metrics/notifications [get]
func (h *MetricsHandler) GetNotificationMetrics(c *gin.Context) {
	orgID, exists := c.Get("organization_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Organization not found"})
		return
	}

	filter := parseMetricsFilter(c)

	metrics, err := h.metricsService.GetNotificationMetrics(c.Request.Context(), orgID.(uuid.UUID), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, metrics)
}

// GetAlertTrend godoc
// @Summary Get alert trend data
// @Description Get time-series data for alert creation and closure
// @Tags Metrics
// @Accept json
// @Produce json
// @Param start_time query string false "Start time (RFC3339 format)"
// @Param end_time query string false "End time (RFC3339 format)"
// @Param period query string false "Trend period: hourly, daily, weekly" default(daily)
// @Success 200 {object} domain.AlertTrend
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /metrics/alerts/trend [get]
func (h *MetricsHandler) GetAlertTrend(c *gin.Context) {
	orgID, exists := c.Get("organization_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Organization not found"})
		return
	}

	filter := parseMetricsFilter(c)

	trend, err := h.metricsService.GetAlertTrend(c.Request.Context(), orgID.(uuid.UUID), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, trend)
}

// GetTeamMetrics godoc
// @Summary Get team performance metrics
// @Description Get performance metrics for all teams including alert counts and response times
// @Tags Metrics
// @Accept json
// @Produce json
// @Param start_time query string false "Start time (RFC3339 format)"
// @Param end_time query string false "End time (RFC3339 format)"
// @Success 200 {object} []domain.TeamMetrics
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /metrics/teams [get]
func (h *MetricsHandler) GetTeamMetrics(c *gin.Context) {
	orgID, exists := c.Get("organization_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Organization not found"})
		return
	}

	filter := parseMetricsFilter(c)

	metrics, err := h.metricsService.GetTeamMetrics(c.Request.Context(), orgID.(uuid.UUID), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"teams": metrics})
}
