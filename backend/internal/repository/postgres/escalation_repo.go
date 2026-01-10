package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/nmn3m/pulsar/backend/internal/domain"
)

type EscalationPolicyRepository struct {
	db *DB
}

func NewEscalationPolicyRepository(db *DB) *EscalationPolicyRepository {
	return &EscalationPolicyRepository{db: db}
}

// Policy CRUD operations

func (r *EscalationPolicyRepository) Create(ctx context.Context, policy *domain.EscalationPolicy) error {
	query := `
		INSERT INTO escalation_policies (id, organization_id, name, description, repeat_enabled, repeat_count)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at, updated_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		policy.ID,
		policy.OrganizationID,
		policy.Name,
		policy.Description,
		policy.RepeatEnabled,
		policy.RepeatCount,
	).Scan(&policy.CreatedAt, &policy.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create escalation policy: %w", err)
	}

	return nil
}

func (r *EscalationPolicyRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.EscalationPolicy, error) {
	query := `
		SELECT id, organization_id, name, description, repeat_enabled, repeat_count, created_at, updated_at
		FROM escalation_policies
		WHERE id = $1
	`

	var policy domain.EscalationPolicy
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&policy.ID,
		&policy.OrganizationID,
		&policy.Name,
		&policy.Description,
		&policy.RepeatEnabled,
		&policy.RepeatCount,
		&policy.CreatedAt,
		&policy.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("escalation policy not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get escalation policy: %w", err)
	}

	return &policy, nil
}

func (r *EscalationPolicyRepository) Update(ctx context.Context, policy *domain.EscalationPolicy) error {
	query := `
		UPDATE escalation_policies
		SET name = $2, description = $3, repeat_enabled = $4, repeat_count = $5
		WHERE id = $1
		RETURNING updated_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		policy.ID,
		policy.Name,
		policy.Description,
		policy.RepeatEnabled,
		policy.RepeatCount,
	).Scan(&policy.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to update escalation policy: %w", err)
	}

	return nil
}

func (r *EscalationPolicyRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM escalation_policies WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete escalation policy: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("escalation policy not found")
	}

	return nil
}

func (r *EscalationPolicyRepository) List(ctx context.Context, orgID uuid.UUID, limit, offset int) ([]*domain.EscalationPolicy, error) {
	query := `
		SELECT id, organization_id, name, description, repeat_enabled, repeat_count, created_at, updated_at
		FROM escalation_policies
		WHERE organization_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, orgID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list escalation policies: %w", err)
	}
	defer rows.Close()

	var policies []*domain.EscalationPolicy
	for rows.Next() {
		var policy domain.EscalationPolicy
		err := rows.Scan(
			&policy.ID,
			&policy.OrganizationID,
			&policy.Name,
			&policy.Description,
			&policy.RepeatEnabled,
			&policy.RepeatCount,
			&policy.CreatedAt,
			&policy.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan escalation policy: %w", err)
		}

		policies = append(policies, &policy)
	}

	return policies, nil
}

func (r *EscalationPolicyRepository) GetWithRules(ctx context.Context, id uuid.UUID) (*domain.EscalationPolicyWithRules, error) {
	policy, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	rules, err := r.ListRules(ctx, id)
	if err != nil {
		return nil, err
	}

	// Load targets for each rule
	rulesWithTargets := make([]*domain.EscalationRuleWithTargets, 0, len(rules))
	for _, rule := range rules {
		targets, err := r.ListTargets(ctx, rule.ID)
		if err != nil {
			return nil, err
		}

		rulesWithTargets = append(rulesWithTargets, &domain.EscalationRuleWithTargets{
			EscalationRule: *rule,
			Targets:        targets,
		})
	}

	return &domain.EscalationPolicyWithRules{
		EscalationPolicy: *policy,
		Rules:            rulesWithTargets,
	}, nil
}

// Rule CRUD operations

func (r *EscalationPolicyRepository) CreateRule(ctx context.Context, rule *domain.EscalationRule) error {
	query := `
		INSERT INTO escalation_rules (id, policy_id, position, escalation_delay)
		VALUES ($1, $2, $3, $4)
		RETURNING created_at, updated_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		rule.ID,
		rule.PolicyID,
		rule.Position,
		rule.EscalationDelay,
	).Scan(&rule.CreatedAt, &rule.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create escalation rule: %w", err)
	}

	return nil
}

func (r *EscalationPolicyRepository) GetRule(ctx context.Context, id uuid.UUID) (*domain.EscalationRule, error) {
	query := `
		SELECT id, policy_id, position, escalation_delay, created_at, updated_at
		FROM escalation_rules
		WHERE id = $1
	`

	var rule domain.EscalationRule
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&rule.ID,
		&rule.PolicyID,
		&rule.Position,
		&rule.EscalationDelay,
		&rule.CreatedAt,
		&rule.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("escalation rule not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get escalation rule: %w", err)
	}

	return &rule, nil
}

func (r *EscalationPolicyRepository) UpdateRule(ctx context.Context, rule *domain.EscalationRule) error {
	query := `
		UPDATE escalation_rules
		SET position = $2, escalation_delay = $3
		WHERE id = $1
		RETURNING updated_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		rule.ID,
		rule.Position,
		rule.EscalationDelay,
	).Scan(&rule.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to update escalation rule: %w", err)
	}

	return nil
}

func (r *EscalationPolicyRepository) DeleteRule(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM escalation_rules WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete escalation rule: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("escalation rule not found")
	}

	return nil
}

func (r *EscalationPolicyRepository) ListRules(ctx context.Context, policyID uuid.UUID) ([]*domain.EscalationRule, error) {
	query := `
		SELECT id, policy_id, position, escalation_delay, created_at, updated_at
		FROM escalation_rules
		WHERE policy_id = $1
		ORDER BY position ASC
	`

	rows, err := r.db.QueryContext(ctx, query, policyID)
	if err != nil {
		return nil, fmt.Errorf("failed to list escalation rules: %w", err)
	}
	defer rows.Close()

	var rules []*domain.EscalationRule
	for rows.Next() {
		var rule domain.EscalationRule
		err := rows.Scan(
			&rule.ID,
			&rule.PolicyID,
			&rule.Position,
			&rule.EscalationDelay,
			&rule.CreatedAt,
			&rule.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan escalation rule: %w", err)
		}

		rules = append(rules, &rule)
	}

	return rules, nil
}

// Target CRUD operations

func (r *EscalationPolicyRepository) AddTarget(ctx context.Context, target *domain.EscalationTarget) error {
	query := `
		INSERT INTO escalation_targets (id, rule_id, target_type, target_id, notification_channels)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING created_at
	`

	// Handle nil or empty notification channels - use nil for NULL in PostgreSQL
	var notificationChannels interface{}
	if len(target.NotificationChannels) > 0 {
		notificationChannels = target.NotificationChannels
	}

	err := r.db.QueryRowContext(
		ctx,
		query,
		target.ID,
		target.RuleID,
		target.TargetType.String(),
		target.TargetID,
		notificationChannels,
	).Scan(&target.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to add escalation target: %w", err)
	}

	return nil
}

func (r *EscalationPolicyRepository) RemoveTarget(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM escalation_targets WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to remove escalation target: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("escalation target not found")
	}

	return nil
}

func (r *EscalationPolicyRepository) ListTargets(ctx context.Context, ruleID uuid.UUID) ([]*domain.EscalationTarget, error) {
	query := `
		SELECT id, rule_id, target_type, target_id, COALESCE(notification_channels, 'null'::jsonb), created_at
		FROM escalation_targets
		WHERE rule_id = $1
		ORDER BY created_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query, ruleID)
	if err != nil {
		return nil, fmt.Errorf("failed to list escalation targets: %w", err)
	}
	defer rows.Close()

	var targets []*domain.EscalationTarget
	for rows.Next() {
		var target domain.EscalationTarget
		var targetType string

		err := rows.Scan(
			&target.ID,
			&target.RuleID,
			&targetType,
			&target.TargetID,
			&target.NotificationChannels,
			&target.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan escalation target: %w", err)
		}

		target.TargetType = domain.EscalationTargetType(targetType)
		// Clear notification channels if it's just "null"
		if string(target.NotificationChannels) == "null" {
			target.NotificationChannels = nil
		}
		targets = append(targets, &target)
	}

	return targets, nil
}

// Escalation event operations

func (r *EscalationPolicyRepository) CreateEvent(ctx context.Context, event *domain.AlertEscalationEvent) error {
	query := `
		INSERT INTO alert_escalation_events (id, alert_id, policy_id, rule_id, event_type, current_level, repeat_count, next_escalation_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING created_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		event.ID,
		event.AlertID,
		event.PolicyID,
		event.RuleID,
		event.EventType.String(),
		event.CurrentLevel,
		event.RepeatCount,
		event.NextEscalationAt,
	).Scan(&event.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create escalation event: %w", err)
	}

	return nil
}

func (r *EscalationPolicyRepository) GetLatestEvent(ctx context.Context, alertID uuid.UUID) (*domain.AlertEscalationEvent, error) {
	query := `
		SELECT id, alert_id, policy_id, rule_id, event_type, current_level, repeat_count, next_escalation_at, created_at
		FROM alert_escalation_events
		WHERE alert_id = $1
		ORDER BY created_at DESC
		LIMIT 1
	`

	var event domain.AlertEscalationEvent
	var eventType string

	err := r.db.QueryRowContext(ctx, query, alertID).Scan(
		&event.ID,
		&event.AlertID,
		&event.PolicyID,
		&event.RuleID,
		&eventType,
		&event.CurrentLevel,
		&event.RepeatCount,
		&event.NextEscalationAt,
		&event.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil // No event found
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get latest escalation event: %w", err)
	}

	event.EventType = domain.EscalationEventType(eventType)

	return &event, nil
}

func (r *EscalationPolicyRepository) UpdateEvent(ctx context.Context, event *domain.AlertEscalationEvent) error {
	query := `
		UPDATE alert_escalation_events
		SET rule_id = $2, event_type = $3, current_level = $4, repeat_count = $5, next_escalation_at = $6
		WHERE id = $1
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		event.ID,
		event.RuleID,
		event.EventType.String(),
		event.CurrentLevel,
		event.RepeatCount,
		event.NextEscalationAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update escalation event: %w", err)
	}

	return nil
}

func (r *EscalationPolicyRepository) ListPendingEscalations(ctx context.Context, before time.Time) ([]*domain.AlertEscalationEvent, error) {
	query := `
		SELECT id, alert_id, policy_id, rule_id, event_type, current_level, repeat_count, next_escalation_at, created_at
		FROM alert_escalation_events
		WHERE next_escalation_at IS NOT NULL
		  AND next_escalation_at <= $1
		  AND event_type = 'triggered'
		ORDER BY next_escalation_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query, before)
	if err != nil {
		return nil, fmt.Errorf("failed to list pending escalations: %w", err)
	}
	defer rows.Close()

	var events []*domain.AlertEscalationEvent
	for rows.Next() {
		var event domain.AlertEscalationEvent
		var eventType string

		err := rows.Scan(
			&event.ID,
			&event.AlertID,
			&event.PolicyID,
			&event.RuleID,
			&eventType,
			&event.CurrentLevel,
			&event.RepeatCount,
			&event.NextEscalationAt,
			&event.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan escalation event: %w", err)
		}

		event.EventType = domain.EscalationEventType(eventType)
		events = append(events, &event)
	}

	return events, nil
}
