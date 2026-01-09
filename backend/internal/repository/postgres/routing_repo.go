package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"

	"github.com/nmn3m/pulsar/backend/internal/domain"
)

type RoutingRuleRepository struct {
	db *DB
}

func NewRoutingRuleRepository(db *DB) *RoutingRuleRepository {
	return &RoutingRuleRepository{db: db}
}

func (r *RoutingRuleRepository) Create(ctx context.Context, rule *domain.AlertRoutingRule) error {
	query := `
		INSERT INTO alert_routing_rules (id, organization_id, name, description, priority, conditions, actions, enabled)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING created_at, updated_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		rule.ID,
		rule.OrganizationID,
		rule.Name,
		rule.Description,
		rule.Priority,
		rule.Conditions,
		rule.Actions,
		rule.Enabled,
	).Scan(&rule.CreatedAt, &rule.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create routing rule: %w", err)
	}

	return nil
}

func (r *RoutingRuleRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.AlertRoutingRule, error) {
	query := `
		SELECT id, organization_id, name, description, priority, conditions, actions, enabled, created_at, updated_at
		FROM alert_routing_rules
		WHERE id = $1
	`

	var rule domain.AlertRoutingRule
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&rule.ID,
		&rule.OrganizationID,
		&rule.Name,
		&rule.Description,
		&rule.Priority,
		&rule.Conditions,
		&rule.Actions,
		&rule.Enabled,
		&rule.CreatedAt,
		&rule.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("routing rule not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get routing rule: %w", err)
	}

	return &rule, nil
}

func (r *RoutingRuleRepository) Update(ctx context.Context, rule *domain.AlertRoutingRule) error {
	query := `
		UPDATE alert_routing_rules
		SET name = $2, description = $3, priority = $4, conditions = $5, actions = $6, enabled = $7, updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		rule.ID,
		rule.Name,
		rule.Description,
		rule.Priority,
		rule.Conditions,
		rule.Actions,
		rule.Enabled,
	).Scan(&rule.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to update routing rule: %w", err)
	}

	return nil
}

func (r *RoutingRuleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM alert_routing_rules WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete routing rule: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("routing rule not found")
	}

	return nil
}

func (r *RoutingRuleRepository) List(ctx context.Context, orgID uuid.UUID, limit, offset int) ([]*domain.AlertRoutingRule, error) {
	query := `
		SELECT id, organization_id, name, description, priority, conditions, actions, enabled, created_at, updated_at
		FROM alert_routing_rules
		WHERE organization_id = $1
		ORDER BY priority ASC, created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, orgID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list routing rules: %w", err)
	}
	defer rows.Close()

	var rules []*domain.AlertRoutingRule
	for rows.Next() {
		var rule domain.AlertRoutingRule
		err := rows.Scan(
			&rule.ID,
			&rule.OrganizationID,
			&rule.Name,
			&rule.Description,
			&rule.Priority,
			&rule.Conditions,
			&rule.Actions,
			&rule.Enabled,
			&rule.CreatedAt,
			&rule.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan routing rule: %w", err)
		}

		rules = append(rules, &rule)
	}

	return rules, nil
}

func (r *RoutingRuleRepository) ListEnabled(ctx context.Context, orgID uuid.UUID) ([]*domain.AlertRoutingRule, error) {
	query := `
		SELECT id, organization_id, name, description, priority, conditions, actions, enabled, created_at, updated_at
		FROM alert_routing_rules
		WHERE organization_id = $1 AND enabled = true
		ORDER BY priority ASC
	`

	rows, err := r.db.QueryContext(ctx, query, orgID)
	if err != nil {
		return nil, fmt.Errorf("failed to list enabled routing rules: %w", err)
	}
	defer rows.Close()

	var rules []*domain.AlertRoutingRule
	for rows.Next() {
		var rule domain.AlertRoutingRule
		err := rows.Scan(
			&rule.ID,
			&rule.OrganizationID,
			&rule.Name,
			&rule.Description,
			&rule.Priority,
			&rule.Conditions,
			&rule.Actions,
			&rule.Enabled,
			&rule.CreatedAt,
			&rule.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan routing rule: %w", err)
		}

		rules = append(rules, &rule)
	}

	return rules, nil
}

func (r *RoutingRuleRepository) Reorder(ctx context.Context, orgID uuid.UUID, ruleIDs []uuid.UUID) error {
	// Start a transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Update priority for each rule based on its position in the array
	query := `UPDATE alert_routing_rules SET priority = $1, updated_at = NOW() WHERE id = $2 AND organization_id = $3`

	for i, ruleID := range ruleIDs {
		result, err := tx.ExecContext(ctx, query, i, ruleID, orgID)
		if err != nil {
			return fmt.Errorf("failed to update rule priority: %w", err)
		}

		rows, err := result.RowsAffected()
		if err != nil {
			return fmt.Errorf("failed to get rows affected: %w", err)
		}

		if rows == 0 {
			return fmt.Errorf("routing rule not found or does not belong to organization")
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
