package dto

type CreateWebhookEndpointRequest struct {
	Name              string            `json:"name" binding:"required"`
	URL               string            `json:"url" binding:"required"`
	Enabled           bool              `json:"enabled"`
	AlertCreated      bool              `json:"alert_created"`
	AlertUpdated      bool              `json:"alert_updated"`
	AlertAcknowledged bool              `json:"alert_acknowledged"`
	AlertClosed       bool              `json:"alert_closed"`
	AlertEscalated    bool              `json:"alert_escalated"`
	IncidentCreated   bool              `json:"incident_created"`
	IncidentUpdated   bool              `json:"incident_updated"`
	IncidentResolved  bool              `json:"incident_resolved"`
	Headers           map[string]string `json:"headers"`
	TimeoutSeconds    *int              `json:"timeout_seconds"`
	MaxRetries        *int              `json:"max_retries"`
	RetryDelaySeconds *int              `json:"retry_delay_seconds"`
}

type UpdateWebhookEndpointRequest struct {
	Name              *string           `json:"name"`
	URL               *string           `json:"url"`
	Enabled           *bool             `json:"enabled"`
	AlertCreated      *bool             `json:"alert_created"`
	AlertUpdated      *bool             `json:"alert_updated"`
	AlertAcknowledged *bool             `json:"alert_acknowledged"`
	AlertClosed       *bool             `json:"alert_closed"`
	AlertEscalated    *bool             `json:"alert_escalated"`
	IncidentCreated   *bool             `json:"incident_created"`
	IncidentUpdated   *bool             `json:"incident_updated"`
	IncidentResolved  *bool             `json:"incident_resolved"`
	Headers           map[string]string `json:"headers"`
	TimeoutSeconds    *int              `json:"timeout_seconds"`
	MaxRetries        *int              `json:"max_retries"`
	RetryDelaySeconds *int              `json:"retry_delay_seconds"`
}

type CreateIncomingWebhookTokenRequest struct {
	Name            string   `json:"name" binding:"required"`
	IntegrationType string   `json:"integration_type" binding:"required"`
	DefaultPriority *string  `json:"default_priority"`
	DefaultTags     []string `json:"default_tags"`
}
