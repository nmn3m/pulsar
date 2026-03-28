package outbound

import (
	"context"

	"github.com/google/uuid"
)

type WebhookDispatcher interface {
	TriggerWebhooks(ctx context.Context, orgID uuid.UUID, eventType string, data map[string]interface{})
}
