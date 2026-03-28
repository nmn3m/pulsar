package outbound

import (
	"context"

	"github.com/google/uuid"

	"github.com/nmn3m/pulsar/backend/internal/core/domain"
)

type RoutingRuleRepository interface {
	Create(ctx context.Context, rule *domain.AlertRoutingRule) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.AlertRoutingRule, error)
	Update(ctx context.Context, rule *domain.AlertRoutingRule) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, orgID uuid.UUID, limit, offset int) ([]*domain.AlertRoutingRule, error)
	ListEnabled(ctx context.Context, orgID uuid.UUID) ([]*domain.AlertRoutingRule, error)
	Reorder(ctx context.Context, orgID uuid.UUID, ruleIDs []uuid.UUID) error
}
