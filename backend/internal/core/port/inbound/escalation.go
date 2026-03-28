package inbound

import (
	"context"

	"github.com/google/uuid"

	"github.com/nmn3m/pulsar/backend/internal/core/domain"
	"github.com/nmn3m/pulsar/backend/internal/core/dto"
)

type EscalationService interface {
	CreatePolicy(ctx context.Context, orgID uuid.UUID, req *dto.CreateEscalationPolicyRequest) (*domain.EscalationPolicy, error)
	GetPolicy(ctx context.Context, id uuid.UUID) (*domain.EscalationPolicy, error)
	GetPolicyWithRules(ctx context.Context, id uuid.UUID) (*domain.EscalationPolicyWithRules, error)
	UpdatePolicy(ctx context.Context, id uuid.UUID, req *dto.UpdateEscalationPolicyRequest) (*domain.EscalationPolicy, error)
	DeletePolicy(ctx context.Context, id uuid.UUID) error
	ListPolicies(ctx context.Context, orgID uuid.UUID, page, pageSize int) ([]*domain.EscalationPolicy, error)
	CreateRule(ctx context.Context, policyID uuid.UUID, req *dto.CreateEscalationRuleRequest) (*domain.EscalationRule, error)
	GetRule(ctx context.Context, id uuid.UUID) (*domain.EscalationRule, error)
	UpdateRule(ctx context.Context, id uuid.UUID, req *dto.UpdateEscalationRuleRequest) (*domain.EscalationRule, error)
	DeleteRule(ctx context.Context, id uuid.UUID) error
	ListRules(ctx context.Context, policyID uuid.UUID) ([]*domain.EscalationRule, error)
	AddTarget(ctx context.Context, ruleID uuid.UUID, req *dto.AddEscalationTargetRequest) (*domain.EscalationTarget, error)
	RemoveTarget(ctx context.Context, id uuid.UUID) error
	ListTargets(ctx context.Context, ruleID uuid.UUID) ([]*domain.EscalationTarget, error)
	StartEscalation(ctx context.Context, alertID, orgID uuid.UUID) error
	ProcessPendingEscalations(ctx context.Context) error
	StopEscalation(ctx context.Context, alertID uuid.UUID) error
}
