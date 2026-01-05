package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nmn3m/pulsar/backend/internal/domain"
	"github.com/nmn3m/pulsar/backend/internal/repository"
)

type EscalationService struct {
	escalationRepo repository.EscalationPolicyRepository
	alertRepo      repository.AlertRepository
	alertNotifier  *AlertNotifier
}

func NewEscalationService(
	escalationRepo repository.EscalationPolicyRepository,
	alertRepo repository.AlertRepository,
	alertNotifier *AlertNotifier,
) *EscalationService {
	return &EscalationService{
		escalationRepo: escalationRepo,
		alertRepo:      alertRepo,
		alertNotifier:  alertNotifier,
	}
}

// Request/Response types

type CreateEscalationPolicyRequest struct {
	Name          string  `json:"name" binding:"required"`
	Description   *string `json:"description"`
	RepeatEnabled bool    `json:"repeat_enabled"`
	RepeatCount   *int    `json:"repeat_count"`
}

type UpdateEscalationPolicyRequest struct {
	Name          *string `json:"name"`
	Description   *string `json:"description"`
	RepeatEnabled *bool   `json:"repeat_enabled"`
	RepeatCount   *int    `json:"repeat_count"`
}

type CreateEscalationRuleRequest struct {
	Position        int `json:"position" binding:"required"`
	EscalationDelay int `json:"escalation_delay" binding:"required"`
}

type UpdateEscalationRuleRequest struct {
	Position        *int `json:"position"`
	EscalationDelay *int `json:"escalation_delay"`
}

type AddEscalationTargetRequest struct {
	TargetType string    `json:"target_type" binding:"required"`
	TargetID   uuid.UUID `json:"target_id" binding:"required"`
}

// Policy CRUD

func (s *EscalationService) CreatePolicy(ctx context.Context, orgID uuid.UUID, req *CreateEscalationPolicyRequest) (*domain.EscalationPolicy, error) {
	policy := &domain.EscalationPolicy{
		ID:             uuid.New(),
		OrganizationID: orgID,
		Name:           req.Name,
		Description:    req.Description,
		RepeatEnabled:  req.RepeatEnabled,
		RepeatCount:    req.RepeatCount,
	}

	if err := s.escalationRepo.Create(ctx, policy); err != nil {
		return nil, fmt.Errorf("failed to create escalation policy: %w", err)
	}

	return policy, nil
}

func (s *EscalationService) GetPolicy(ctx context.Context, id uuid.UUID) (*domain.EscalationPolicy, error) {
	policy, err := s.escalationRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get escalation policy: %w", err)
	}

	return policy, nil
}

func (s *EscalationService) GetPolicyWithRules(ctx context.Context, id uuid.UUID) (*domain.EscalationPolicyWithRules, error) {
	policy, err := s.escalationRepo.GetWithRules(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get escalation policy with rules: %w", err)
	}

	return policy, nil
}

func (s *EscalationService) UpdatePolicy(ctx context.Context, id uuid.UUID, req *UpdateEscalationPolicyRequest) (*domain.EscalationPolicy, error) {
	policy, err := s.escalationRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get escalation policy: %w", err)
	}

	if req.Name != nil {
		policy.Name = *req.Name
	}
	if req.Description != nil {
		policy.Description = req.Description
	}
	if req.RepeatEnabled != nil {
		policy.RepeatEnabled = *req.RepeatEnabled
	}
	if req.RepeatCount != nil {
		policy.RepeatCount = req.RepeatCount
	}

	if err := s.escalationRepo.Update(ctx, policy); err != nil {
		return nil, fmt.Errorf("failed to update escalation policy: %w", err)
	}

	return policy, nil
}

func (s *EscalationService) DeletePolicy(ctx context.Context, id uuid.UUID) error {
	if err := s.escalationRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete escalation policy: %w", err)
	}

	return nil
}

func (s *EscalationService) ListPolicies(ctx context.Context, orgID uuid.UUID, page, pageSize int) ([]*domain.EscalationPolicy, error) {
	if page < 1 {
		page = 1
	}

	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize

	policies, err := s.escalationRepo.List(ctx, orgID, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list escalation policies: %w", err)
	}

	return policies, nil
}

// Rule CRUD

func (s *EscalationService) CreateRule(ctx context.Context, policyID uuid.UUID, req *CreateEscalationRuleRequest) (*domain.EscalationRule, error) {
	rule := &domain.EscalationRule{
		ID:              uuid.New(),
		PolicyID:        policyID,
		Position:        req.Position,
		EscalationDelay: req.EscalationDelay,
	}

	if err := s.escalationRepo.CreateRule(ctx, rule); err != nil {
		return nil, fmt.Errorf("failed to create escalation rule: %w", err)
	}

	return rule, nil
}

func (s *EscalationService) GetRule(ctx context.Context, id uuid.UUID) (*domain.EscalationRule, error) {
	rule, err := s.escalationRepo.GetRule(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get escalation rule: %w", err)
	}

	return rule, nil
}

func (s *EscalationService) UpdateRule(ctx context.Context, id uuid.UUID, req *UpdateEscalationRuleRequest) (*domain.EscalationRule, error) {
	rule, err := s.escalationRepo.GetRule(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get escalation rule: %w", err)
	}

	if req.Position != nil {
		rule.Position = *req.Position
	}
	if req.EscalationDelay != nil {
		rule.EscalationDelay = *req.EscalationDelay
	}

	if err := s.escalationRepo.UpdateRule(ctx, rule); err != nil {
		return nil, fmt.Errorf("failed to update escalation rule: %w", err)
	}

	return rule, nil
}

func (s *EscalationService) DeleteRule(ctx context.Context, id uuid.UUID) error {
	if err := s.escalationRepo.DeleteRule(ctx, id); err != nil {
		return fmt.Errorf("failed to delete escalation rule: %w", err)
	}

	return nil
}

func (s *EscalationService) ListRules(ctx context.Context, policyID uuid.UUID) ([]*domain.EscalationRule, error) {
	rules, err := s.escalationRepo.ListRules(ctx, policyID)
	if err != nil {
		return nil, fmt.Errorf("failed to list escalation rules: %w", err)
	}

	return rules, nil
}

// Target CRUD

func (s *EscalationService) AddTarget(ctx context.Context, ruleID uuid.UUID, req *AddEscalationTargetRequest) (*domain.EscalationTarget, error) {
	targetType := domain.EscalationTargetType(req.TargetType)
	if err := targetType.Validate(); err != nil {
		return nil, err
	}

	target := &domain.EscalationTarget{
		ID:         uuid.New(),
		RuleID:     ruleID,
		TargetType: targetType,
		TargetID:   req.TargetID,
	}

	if err := s.escalationRepo.AddTarget(ctx, target); err != nil {
		return nil, fmt.Errorf("failed to add escalation target: %w", err)
	}

	return target, nil
}

func (s *EscalationService) RemoveTarget(ctx context.Context, id uuid.UUID) error {
	if err := s.escalationRepo.RemoveTarget(ctx, id); err != nil {
		return fmt.Errorf("failed to remove escalation target: %w", err)
	}

	return nil
}

func (s *EscalationService) ListTargets(ctx context.Context, ruleID uuid.UUID) ([]*domain.EscalationTarget, error) {
	targets, err := s.escalationRepo.ListTargets(ctx, ruleID)
	if err != nil {
		return nil, fmt.Errorf("failed to list escalation targets: %w", err)
	}

	return targets, nil
}

// Escalation logic

func (s *EscalationService) StartEscalation(ctx context.Context, alertID uuid.UUID) error {
	// Get the alert to find its policy
	alert, err := s.alertRepo.GetByID(ctx, alertID)
	if err != nil {
		return fmt.Errorf("failed to get alert: %w", err)
	}

	if alert.EscalationPolicyID == nil {
		return nil // No escalation policy configured
	}

	// Get policy with rules
	policy, err := s.escalationRepo.GetWithRules(ctx, *alert.EscalationPolicyID)
	if err != nil {
		return fmt.Errorf("failed to get escalation policy: %w", err)
	}

	if len(policy.Rules) == 0 {
		return nil // No rules to escalate
	}

	// Create initial escalation event
	firstRule := policy.Rules[0]
	nextEscalationTime := time.Now().Add(time.Duration(firstRule.EscalationDelay) * time.Minute)

	event := &domain.AlertEscalationEvent{
		ID:               uuid.New(),
		AlertID:          alertID,
		PolicyID:         policy.ID,
		RuleID:           &firstRule.ID,
		EventType:        domain.EscalationEventTriggered,
		CurrentLevel:     0,
		RepeatCount:      0,
		NextEscalationAt: &nextEscalationTime,
	}

	if err := s.escalationRepo.CreateEvent(ctx, event); err != nil {
		return fmt.Errorf("failed to create escalation event: %w", err)
	}

	return nil
}

func (s *EscalationService) ProcessPendingEscalations(ctx context.Context) error {
	// Get all escalations that should be triggered now
	events, err := s.escalationRepo.ListPendingEscalations(ctx, time.Now())
	if err != nil {
		return fmt.Errorf("failed to list pending escalations: %w", err)
	}

	for _, event := range events {
		if err := s.processEscalation(ctx, event); err != nil {
			// Log error but continue processing other escalations
			fmt.Printf("Failed to process escalation for alert %s: %v\n", event.AlertID, err)
		}
	}

	return nil
}

func (s *EscalationService) processEscalation(ctx context.Context, event *domain.AlertEscalationEvent) error {
	// Get the policy with rules
	policy, err := s.escalationRepo.GetWithRules(ctx, event.PolicyID)
	if err != nil {
		return fmt.Errorf("failed to get policy: %w", err)
	}

	// Check if there are more rules to escalate to
	nextLevel := event.CurrentLevel + 1

	if nextLevel < len(policy.Rules) {
		// Move to next rule
		nextRule := policy.Rules[nextLevel]
		nextEscalationTime := time.Now().Add(time.Duration(nextRule.EscalationDelay) * time.Minute)

		event.CurrentLevel = nextLevel
		event.RuleID = &nextRule.ID
		event.NextEscalationAt = &nextEscalationTime

		if err := s.escalationRepo.UpdateEvent(ctx, event); err != nil {
			return fmt.Errorf("failed to update event: %w", err)
		}

		// Trigger notifications to targets in nextRule
		if err := s.sendEscalationNotifications(ctx, event, nextRule); err != nil {
			fmt.Printf("Failed to send escalation notifications for alert %s: %v\n", event.AlertID, err)
		}

	} else if policy.RepeatEnabled {
		// Check if we should repeat
		if policy.RepeatCount == nil || event.RepeatCount < *policy.RepeatCount {
			// Restart from first rule
			firstRule := policy.Rules[0]
			nextEscalationTime := time.Now().Add(time.Duration(firstRule.EscalationDelay) * time.Minute)

			event.CurrentLevel = 0
			event.RuleID = &firstRule.ID
			event.RepeatCount++
			event.NextEscalationAt = &nextEscalationTime

			if err := s.escalationRepo.UpdateEvent(ctx, event); err != nil {
				return fmt.Errorf("failed to update event: %w", err)
			}

			// Trigger notifications to targets in firstRule (repeat cycle)
			if err := s.sendEscalationNotifications(ctx, event, firstRule); err != nil {
				fmt.Printf("Failed to send escalation notifications for alert %s (repeat): %v\n", event.AlertID, err)
			}
		} else {
			// Max repeats reached, mark as completed
			event.EventType = domain.EscalationEventCompleted
			event.NextEscalationAt = nil

			if err := s.escalationRepo.UpdateEvent(ctx, event); err != nil {
				return fmt.Errorf("failed to update event: %w", err)
			}
		}
	} else {
		// No more rules and no repeat, mark as completed
		event.EventType = domain.EscalationEventCompleted
		event.NextEscalationAt = nil

		if err := s.escalationRepo.UpdateEvent(ctx, event); err != nil {
			return fmt.Errorf("failed to update event: %w", err)
		}
	}

	return nil
}

func (s *EscalationService) sendEscalationNotifications(ctx context.Context, event *domain.AlertEscalationEvent, rule domain.EscalationRule) error {
	// Only send notifications if alertNotifier is configured
	if s.alertNotifier == nil {
		return nil
	}

	// Get the alert
	alert, err := s.alertRepo.GetByID(ctx, event.AlertID)
	if err != nil {
		return fmt.Errorf("failed to get alert: %w", err)
	}

	// Get targets for this rule
	targets, err := s.escalationRepo.ListTargets(ctx, rule.ID)
	if err != nil {
		return fmt.Errorf("failed to get targets: %w", err)
	}

	if len(targets) == 0 {
		return nil // No targets configured
	}

	// Update alert escalation level
	alert.EscalationLevel = event.CurrentLevel
	alert.LastEscalatedAt = &event.TriggeredAt

	// Send notifications to all targets
	if err := s.alertNotifier.NotifyAlertEscalated(ctx, alert, &rule, targets); err != nil {
		return fmt.Errorf("failed to send notifications: %w", err)
	}

	return nil
}

func (s *EscalationService) StopEscalation(ctx context.Context, alertID uuid.UUID) error {
	// Get the latest escalation event for this alert
	event, err := s.escalationRepo.GetLatestEvent(ctx, alertID)
	if err != nil {
		return fmt.Errorf("failed to get escalation event: %w", err)
	}

	if event == nil || event.EventType != domain.EscalationEventTriggered {
		return nil // No active escalation
	}

	// Mark as acknowledged/stopped
	event.EventType = domain.EscalationEventAcknowledged
	event.NextEscalationAt = nil

	if err := s.escalationRepo.UpdateEvent(ctx, event); err != nil {
		return fmt.Errorf("failed to stop escalation: %w", err)
	}

	return nil
}
