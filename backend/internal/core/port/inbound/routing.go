package inbound

import (
	"context"

	"github.com/google/uuid"

	"github.com/nmn3m/pulsar/backend/internal/core/domain"
	"github.com/nmn3m/pulsar/backend/internal/core/dto"
)

type RoutingService interface {
	CreateRule(ctx context.Context, orgID uuid.UUID, req *dto.CreateRoutingRuleRequest) (*domain.AlertRoutingRule, error)
	GetRule(ctx context.Context, id uuid.UUID) (*domain.AlertRoutingRule, error)
	UpdateRule(ctx context.Context, id uuid.UUID, req *dto.UpdateRoutingRuleRequest) (*domain.AlertRoutingRule, error)
	DeleteRule(ctx context.Context, id uuid.UUID) error
	ListRules(ctx context.Context, orgID uuid.UUID, page, pageSize int) ([]*domain.AlertRoutingRule, error)
	ReorderRules(ctx context.Context, orgID uuid.UUID, req *dto.ReorderRoutingRulesRequest) error
	ApplyRouting(ctx context.Context, orgID uuid.UUID, alert *domain.Alert) (*domain.RoutingActions, error)
}
