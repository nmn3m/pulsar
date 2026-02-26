package repository

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/nmn3m/pulsar/backend/internal/domain"
)

type EscalationPolicyRepository interface {
	Create(ctx context.Context, policy *domain.EscalationPolicy) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.EscalationPolicy, error)
	Update(ctx context.Context, policy *domain.EscalationPolicy) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, orgID uuid.UUID, limit, offset int) ([]*domain.EscalationPolicy, error)
	GetWithRules(ctx context.Context, id uuid.UUID) (*domain.EscalationPolicyWithRules, error)
	CreateRule(ctx context.Context, rule *domain.EscalationRule) error
	GetRule(ctx context.Context, id uuid.UUID) (*domain.EscalationRule, error)
	UpdateRule(ctx context.Context, rule *domain.EscalationRule) error
	DeleteRule(ctx context.Context, id uuid.UUID) error
	ListRules(ctx context.Context, policyID uuid.UUID) ([]*domain.EscalationRule, error)
	AddTarget(ctx context.Context, target *domain.EscalationTarget) error
	RemoveTarget(ctx context.Context, id uuid.UUID) error
	ListTargets(ctx context.Context, ruleID uuid.UUID) ([]*domain.EscalationTarget, error)
	CreateEvent(ctx context.Context, event *domain.AlertEscalationEvent) error
	GetLatestEvent(ctx context.Context, alertID uuid.UUID) (*domain.AlertEscalationEvent, error)
	UpdateEvent(ctx context.Context, event *domain.AlertEscalationEvent) error
	ListPendingEscalations(ctx context.Context, before time.Time) ([]*domain.AlertEscalationEvent, error)
}
