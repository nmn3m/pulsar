package outbound

import (
	"context"

	"github.com/google/uuid"

	"github.com/nmn3m/pulsar/backend/internal/core/domain"
)

type AlertNotificationSender interface {
	NotifyAlertCreated(ctx context.Context, alert *domain.Alert) error
	NotifyAlertAcknowledged(ctx context.Context, alert *domain.Alert, acknowledgedBy uuid.UUID) error
	NotifyAlertClosed(ctx context.Context, alert *domain.Alert, closedBy uuid.UUID, reason string) error
	NotifyAlertEscalated(ctx context.Context, alert *domain.Alert, escalationRule *domain.EscalationRule, targets []domain.EscalationTarget) error
}
