package rest

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nmn3m/pulsar/backend/internal/domain"
	"github.com/nmn3m/pulsar/backend/internal/service"
	"go.uber.org/zap"
)

type IncomingWebhookHandler struct {
	webhookService *service.WebhookService
	alertService   *service.AlertService
	logger         *zap.Logger
}

func NewIncomingWebhookHandler(webhookService *service.WebhookService, alertService *service.AlertService, logger *zap.Logger) *IncomingWebhookHandler {
	return &IncomingWebhookHandler{
		webhookService: webhookService,
		alertService:   alertService,
		logger:         logger,
	}
}

// ReceiveWebhook handles incoming webhooks from external sources
func (h *IncomingWebhookHandler) ReceiveWebhook(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token is required"})
		return
	}

	// Get the incoming webhook token
	webhookToken, err := h.webhookService.GetIncomingTokenByToken(c.Request.Context(), token)
	if err != nil {
		if err == domain.ErrNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate token"})
		return
	}

	if !webhookToken.Enabled {
		c.JSON(http.StatusForbidden, gin.H{"error": "This webhook token is disabled"})
		return
	}

	// Update usage stats
	if err := h.webhookService.UpdateIncomingTokenUsage(c.Request.Context(), webhookToken.ID); err != nil {
		h.logger.Warn("Failed to update webhook token usage", zap.Error(err))
	}

	// Read request body
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}

	// Parse based on integration type
	var alerts []*domain.CreateAlertRequest

	switch webhookToken.IntegrationType {
	case domain.IncomingWebhookPrometheus:
		alerts, err = h.parsePrometheusWebhook(body)
	case domain.IncomingWebhookGrafana:
		alerts, err = h.parseGrafanaWebhook(body)
	case domain.IncomingWebhookGeneric:
		alerts, err = h.parseGenericWebhook(body)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported integration type"})
		return
	}

	if err != nil {
		h.logger.Error("Failed to parse webhook payload",
			zap.String("integration_type", string(webhookToken.IntegrationType)),
			zap.Error(err),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Failed to parse webhook: %v", err)})
		return
	}

	// Create alerts
	createdAlerts := []string{}
	for _, alertReq := range alerts {
		// Apply default priority and tags
		if alertReq.Priority == "" {
			alertReq.Priority = domain.AlertPriority(webhookToken.DefaultPriority)
		}

		// Merge default tags
		tagMap := make(map[string]bool)
		for _, tag := range webhookToken.DefaultTags {
			tagMap[tag] = true
		}
		for _, tag := range alertReq.Tags {
			tagMap[tag] = true
		}
		finalTags := make([]string, 0, len(tagMap))
		for tag := range tagMap {
			finalTags = append(finalTags, tag)
		}
		alertReq.Tags = finalTags

		// Create alert
		alert, err := h.alertService.CreateAlert(c.Request.Context(), webhookToken.OrganizationID, alertReq)
		if err != nil {
			h.logger.Error("Failed to create alert from webhook",
				zap.Error(err),
				zap.String("message", alertReq.Message),
			)
			continue
		}

		createdAlerts = append(createdAlerts, alert.ID.String())
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":         "Alerts created successfully",
		"alerts_created":  len(createdAlerts),
		"alerts_received": len(alerts),
		"alert_ids":       createdAlerts,
	})
}

func (h *IncomingWebhookHandler) parsePrometheusWebhook(body []byte) ([]*domain.CreateAlertRequest, error) {
	var payload struct {
		Alerts []struct {
			Status      string            `json:"status"`
			Labels      map[string]string `json:"labels"`
			Annotations map[string]string `json:"annotations"`
			StartsAt    string            `json:"startsAt"`
			EndsAt      string            `json:"endsAt"`
		} `json:"alerts"`
	}

	if err := parseJSON(body, &payload); err != nil {
		return nil, err
	}

	var alerts []*domain.CreateAlertRequest
	for _, prometheusAlert := range payload.Alerts {
		// Skip resolved alerts
		if prometheusAlert.Status == "resolved" {
			continue
		}

		message := prometheusAlert.Annotations["summary"]
		if message == "" {
			message = prometheusAlert.Labels["alertname"]
		}
		if message == "" {
			message = "Alert from Prometheus"
		}

		description := prometheusAlert.Annotations["description"]

		// Determine priority based on severity label
		priority := domain.AlertPriorityP3
		if severity, ok := prometheusAlert.Labels["severity"]; ok {
			switch strings.ToLower(severity) {
			case "critical":
				priority = domain.AlertPriorityP1
			case "error", "high":
				priority = domain.AlertPriorityP2
			case "warning", "medium":
				priority = domain.AlertPriorityP3
			case "info", "low":
				priority = domain.AlertPriorityP4
			}
		}

		// Convert labels to tags
		tags := []string{"prometheus"}
		for key, value := range prometheusAlert.Labels {
			tags = append(tags, fmt.Sprintf("%s:%s", key, value))
		}

		alerts = append(alerts, &domain.CreateAlertRequest{
			Source:      "prometheus",
			Priority:    priority,
			Message:     message,
			Description: &description,
			Tags:        tags,
		})
	}

	return alerts, nil
}

func (h *IncomingWebhookHandler) parseGrafanaWebhook(body []byte) ([]*domain.CreateAlertRequest, error) {
	var payload struct {
		Title   string `json:"title"`
		State   string `json:"state"`
		Message string `json:"message"`
		RuleURL string `json:"ruleUrl"`
		Tags    map[string]string `json:"tags"`
	}

	if err := parseJSON(body, &payload); err != nil {
		return nil, err
	}

	// Skip resolved alerts
	if payload.State == "ok" {
		return []*domain.CreateAlertRequest{}, nil
	}

	message := payload.Title
	if message == "" {
		message = "Alert from Grafana"
	}

	description := payload.Message
	if description == "" && payload.RuleURL != "" {
		description = "Rule: " + payload.RuleURL
	}

	// Determine priority based on state
	priority := domain.AlertPriorityP3
	switch payload.State {
	case "alerting":
		priority = domain.AlertPriorityP2
	case "no_data":
		priority = domain.AlertPriorityP3
	}

	// Convert tags
	tags := []string{"grafana"}
	for key, value := range payload.Tags {
		tags = append(tags, fmt.Sprintf("%s:%s", key, value))
	}

	return []*domain.CreateAlertRequest{
		{
			Source:      "grafana",
			Priority:    priority,
			Message:     message,
			Description: &description,
			Tags:        tags,
		},
	}, nil
}

func (h *IncomingWebhookHandler) parseGenericWebhook(body []byte) ([]*domain.CreateAlertRequest, error) {
	var payload struct {
		Message     string   `json:"message"`
		Description string   `json:"description"`
		Priority    string   `json:"priority"`
		Tags        []string `json:"tags"`
	}

	if err := parseJSON(body, &payload); err != nil {
		return nil, err
	}

	if payload.Message == "" {
		payload.Message = "Alert from generic webhook"
	}

	priority := domain.AlertPriorityP3
	if payload.Priority != "" {
		priority = domain.AlertPriority(payload.Priority)
	}

	var description *string
	if payload.Description != "" {
		description = &payload.Description
	}

	tags := payload.Tags
	if tags == nil {
		tags = []string{}
	}
	tags = append(tags, "webhook")

	return []*domain.CreateAlertRequest{
		{
			Source:      "webhook",
			Priority:    priority,
			Message:     payload.Message,
			Description: description,
			Tags:        tags,
		},
	}, nil
}
