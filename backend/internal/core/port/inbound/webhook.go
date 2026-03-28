package inbound

import (
	"context"

	"github.com/google/uuid"

	"github.com/nmn3m/pulsar/backend/internal/core/domain"
	"github.com/nmn3m/pulsar/backend/internal/core/dto"
)

type WebhookService interface {
	CreateEndpoint(ctx context.Context, orgID uuid.UUID, req *dto.CreateWebhookEndpointRequest) (*domain.WebhookEndpoint, error)
	GetEndpoint(ctx context.Context, id uuid.UUID) (*domain.WebhookEndpoint, error)
	ListEndpoints(ctx context.Context, orgID uuid.UUID) ([]*domain.WebhookEndpoint, error)
	UpdateEndpoint(ctx context.Context, id, orgID uuid.UUID, req *dto.UpdateWebhookEndpointRequest) (*domain.WebhookEndpoint, error)
	DeleteEndpoint(ctx context.Context, id, orgID uuid.UUID) error
	TriggerWebhooks(ctx context.Context, orgID uuid.UUID, eventType string, data map[string]interface{})
	ProcessPendingDeliveries(ctx context.Context) error
	ListDeliveries(ctx context.Context, orgID uuid.UUID, limit, offset int) ([]*domain.WebhookDelivery, error)
	CreateIncomingToken(ctx context.Context, orgID uuid.UUID, req *dto.CreateIncomingWebhookTokenRequest) (*domain.IncomingWebhookToken, error)
	GetIncomingTokenByToken(ctx context.Context, token string) (*domain.IncomingWebhookToken, error)
	ListIncomingTokens(ctx context.Context, orgID uuid.UUID) ([]*domain.IncomingWebhookToken, error)
	UpdateIncomingTokenUsage(ctx context.Context, tokenID uuid.UUID) error
	DeleteIncomingToken(ctx context.Context, id, orgID uuid.UUID) error
}
