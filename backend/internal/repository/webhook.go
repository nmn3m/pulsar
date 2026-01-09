package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/nmn3m/pulsar/backend/internal/domain"
)

type WebhookRepository interface {
	// Webhook Endpoints
	CreateEndpoint(ctx context.Context, endpoint *domain.WebhookEndpoint) error
	GetEndpointByID(ctx context.Context, id uuid.UUID) (*domain.WebhookEndpoint, error)
	ListEndpoints(ctx context.Context, orgID uuid.UUID) ([]*domain.WebhookEndpoint, error)
	UpdateEndpoint(ctx context.Context, endpoint *domain.WebhookEndpoint) error
	DeleteEndpoint(ctx context.Context, id, orgID uuid.UUID) error

	// Webhook Deliveries
	CreateDelivery(ctx context.Context, delivery *domain.WebhookDelivery) error
	UpdateDelivery(ctx context.Context, delivery *domain.WebhookDelivery) error
	GetPendingDeliveries(ctx context.Context, limit int) ([]*domain.WebhookDelivery, error)
	ListDeliveries(ctx context.Context, orgID uuid.UUID, limit, offset int) ([]*domain.WebhookDelivery, error)

	// Incoming Webhook Tokens
	CreateIncomingToken(ctx context.Context, token *domain.IncomingWebhookToken) error
	GetIncomingTokenByToken(ctx context.Context, token string) (*domain.IncomingWebhookToken, error)
	ListIncomingTokens(ctx context.Context, orgID uuid.UUID) ([]*domain.IncomingWebhookToken, error)
	UpdateIncomingTokenUsage(ctx context.Context, id uuid.UUID) error
	DeleteIncomingToken(ctx context.Context, id, orgID uuid.UUID) error
}
