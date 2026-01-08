package domain

import (
	"time"

	"github.com/google/uuid"
)

// WebhookEndpoint represents an outgoing webhook configuration
type WebhookEndpoint struct {
	ID             uuid.UUID `json:"id" db:"id"`
	OrganizationID uuid.UUID `json:"organization_id" db:"organization_id"`
	Name           string    `json:"name" db:"name"`
	URL            string    `json:"url" db:"url"`
	Secret         string    `json:"-" db:"secret"` // Never expose in JSON
	Enabled        bool      `json:"enabled" db:"enabled"`

	// Event filters
	AlertCreated      bool `json:"alert_created" db:"alert_created"`
	AlertUpdated      bool `json:"alert_updated" db:"alert_updated"`
	AlertAcknowledged bool `json:"alert_acknowledged" db:"alert_acknowledged"`
	AlertClosed       bool `json:"alert_closed" db:"alert_closed"`
	AlertEscalated    bool `json:"alert_escalated" db:"alert_escalated"`
	IncidentCreated   bool `json:"incident_created" db:"incident_created"`
	IncidentUpdated   bool `json:"incident_updated" db:"incident_updated"`
	IncidentResolved  bool `json:"incident_resolved" db:"incident_resolved"`

	// HTTP configuration
	Headers        map[string]string `json:"headers" db:"headers"`
	TimeoutSeconds int               `json:"timeout_seconds" db:"timeout_seconds"`

	// Retry configuration
	MaxRetries        int `json:"max_retries" db:"max_retries"`
	RetryDelaySeconds int `json:"retry_delay_seconds" db:"retry_delay_seconds"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// WebhookDeliveryStatus represents the delivery status of a webhook
type WebhookDeliveryStatus string

const (
	WebhookDeliveryPending WebhookDeliveryStatus = "pending"
	WebhookDeliverySuccess WebhookDeliveryStatus = "success"
	WebhookDeliveryFailed  WebhookDeliveryStatus = "failed"
)

// WebhookDelivery represents a webhook delivery attempt
type WebhookDelivery struct {
	ID                uuid.UUID              `json:"id" db:"id"`
	WebhookEndpointID uuid.UUID              `json:"webhook_endpoint_id" db:"webhook_endpoint_id"`
	OrganizationID    uuid.UUID              `json:"organization_id" db:"organization_id"`
	EventType         string                 `json:"event_type" db:"event_type"`
	Payload           map[string]interface{} `json:"payload" db:"payload"`
	Status            WebhookDeliveryStatus  `json:"status" db:"status"`
	Attempts          int                    `json:"attempts" db:"attempts"`
	LastAttemptAt     *time.Time             `json:"last_attempt_at,omitempty" db:"last_attempt_at"`
	NextRetryAt       *time.Time             `json:"next_retry_at,omitempty" db:"next_retry_at"`
	ResponseStatus    *int                   `json:"response_status_code,omitempty" db:"response_status_code"`
	ResponseBody      *string                `json:"response_body,omitempty" db:"response_body"`
	ErrorMessage      *string                `json:"error_message,omitempty" db:"error_message"`
	CreatedAt         time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time              `json:"updated_at" db:"updated_at"`
}

// IncomingWebhookIntegrationType represents the type of incoming webhook integration
type IncomingWebhookIntegrationType string

const (
	IncomingWebhookGeneric    IncomingWebhookIntegrationType = "generic"
	IncomingWebhookPrometheus IncomingWebhookIntegrationType = "prometheus"
	IncomingWebhookGrafana    IncomingWebhookIntegrationType = "grafana"
	IncomingWebhookDatadog    IncomingWebhookIntegrationType = "datadog"
)

// IncomingWebhookToken represents a token for receiving webhooks from external sources
type IncomingWebhookToken struct {
	ID              uuid.UUID                      `json:"id" db:"id"`
	OrganizationID  uuid.UUID                      `json:"organization_id" db:"organization_id"`
	Name            string                         `json:"name" db:"name"`
	Token           string                         `json:"token" db:"token"`
	Enabled         bool                           `json:"enabled" db:"enabled"`
	IntegrationType IncomingWebhookIntegrationType `json:"integration_type" db:"integration_type"`
	DefaultPriority string                         `json:"default_priority" db:"default_priority"`
	DefaultTags     []string                       `json:"default_tags" db:"default_tags"`
	LastUsedAt      *time.Time                     `json:"last_used_at,omitempty" db:"last_used_at"`
	RequestCount    int                            `json:"request_count" db:"request_count"`
	CreatedAt       time.Time                      `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time                      `json:"updated_at" db:"updated_at"`
}

// WebhookPayload represents the payload sent in outgoing webhooks
type WebhookPayload struct {
	EventType      string                 `json:"event_type"`
	EventID        string                 `json:"event_id"`
	OrganizationID string                 `json:"organization_id"`
	Timestamp      time.Time              `json:"timestamp"`
	Data           map[string]interface{} `json:"data"`
}

// CreateWebhookEndpointRequest represents the request to create a webhook endpoint
type CreateWebhookEndpointRequest struct {
	Name              string            `json:"name" binding:"required"`
	URL               string            `json:"url" binding:"required,url"`
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

// UpdateWebhookEndpointRequest represents the request to update a webhook endpoint
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

// CreateIncomingWebhookTokenRequest represents the request to create an incoming webhook token
type CreateIncomingWebhookTokenRequest struct {
	Name            string                         `json:"name" binding:"required"`
	IntegrationType IncomingWebhookIntegrationType `json:"integration_type" binding:"required"`
	DefaultPriority *string                        `json:"default_priority"`
	DefaultTags     []string                       `json:"default_tags"`
}

// ShouldTriggerEvent checks if the webhook should trigger for the given event type
func (w *WebhookEndpoint) ShouldTriggerEvent(eventType string) bool {
	switch eventType {
	case "alert.created":
		return w.AlertCreated
	case "alert.updated":
		return w.AlertUpdated
	case "alert.acknowledged":
		return w.AlertAcknowledged
	case "alert.closed":
		return w.AlertClosed
	case "alert.escalated":
		return w.AlertEscalated
	case "incident.created":
		return w.IncidentCreated
	case "incident.updated":
		return w.IncidentUpdated
	case "incident.resolved":
		return w.IncidentResolved
	default:
		return false
	}
}
