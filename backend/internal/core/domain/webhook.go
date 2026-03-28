package domain

import (
	"time"

	"github.com/google/uuid"
)

// WebhookEndpoint represents an outgoing webhook configuration
type WebhookEndpoint struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	Name           string
	URL            string
	Secret         string // Never expose in JSON
	Enabled        bool

	// Event filters
	AlertCreated      bool
	AlertUpdated      bool
	AlertAcknowledged bool
	AlertClosed       bool
	AlertEscalated    bool
	IncidentCreated   bool
	IncidentUpdated   bool
	IncidentResolved  bool

	// HTTP configuration
	Headers        map[string]string
	TimeoutSeconds int

	// Retry configuration
	MaxRetries        int
	RetryDelaySeconds int

	CreatedAt time.Time
	UpdatedAt time.Time
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
	ID                uuid.UUID
	WebhookEndpointID uuid.UUID
	OrganizationID    uuid.UUID
	EventType         string
	Payload           map[string]interface{}
	Status            WebhookDeliveryStatus
	Attempts          int
	LastAttemptAt     *time.Time
	NextRetryAt       *time.Time
	ResponseStatus    *int
	ResponseBody      *string
	ErrorMessage      *string
	CreatedAt         time.Time
	UpdatedAt         time.Time
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
	ID              uuid.UUID
	OrganizationID  uuid.UUID
	Name            string
	Token           string
	Enabled         bool
	IntegrationType IncomingWebhookIntegrationType
	DefaultPriority string
	DefaultTags     []string
	LastUsedAt      *time.Time
	RequestCount    int
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// WebhookPayload represents the payload sent in outgoing webhooks
type WebhookPayload struct {
	EventType      string                 `json:"event_type"`
	EventID        string                 `json:"event_id"`
	OrganizationID string                 `json:"organization_id"`
	Timestamp      time.Time              `json:"timestamp"`
	Data           map[string]interface{} `json:"data"`
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
