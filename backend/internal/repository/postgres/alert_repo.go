package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/nmn3m/pulsar/backend/internal/domain"
)

type AlertRepository struct {
	db *DB
}

func NewAlertRepository(db *DB) *AlertRepository {
	return &AlertRepository{db: db}
}

func (r *AlertRepository) Create(ctx context.Context, alert *domain.Alert) error {
	query := `
		INSERT INTO alerts (
			id, organization_id, source, source_id, priority, status,
			message, description, tags, custom_fields,
			assigned_to_user_id, assigned_to_team_id,
			escalation_policy_id, escalation_level,
			dedup_key, dedup_count, first_occurrence_at, last_occurrence_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
		RETURNING created_at, updated_at
	`

	tags, err := json.Marshal(alert.Tags)
	if err != nil {
		return fmt.Errorf("failed to marshal tags: %w", err)
	}

	customFields, err := json.Marshal(alert.CustomFields)
	if err != nil {
		return fmt.Errorf("failed to marshal custom_fields: %w", err)
	}

	err = r.db.QueryRowContext(
		ctx,
		query,
		alert.ID,
		alert.OrganizationID,
		alert.Source,
		alert.SourceID,
		alert.Priority,
		alert.Status,
		alert.Message,
		alert.Description,
		tags,
		customFields,
		alert.AssignedToUserID,
		alert.AssignedToTeamID,
		alert.EscalationPolicyID,
		alert.EscalationLevel,
		alert.DedupKey,
		alert.DedupCount,
		alert.FirstOccurrenceAt,
		alert.LastOccurrenceAt,
	).Scan(&alert.CreatedAt, &alert.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create alert: %w", err)
	}

	return nil
}

func (r *AlertRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Alert, error) {
	query := `
		SELECT
			id, organization_id, source, source_id, priority, status,
			message, description, tags, custom_fields,
			assigned_to_user_id, assigned_to_team_id,
			acknowledged_by, acknowledged_at,
			closed_by, closed_at, close_reason,
			snoozed_until,
			escalation_policy_id, escalation_level, last_escalated_at,
			dedup_key, dedup_count, first_occurrence_at, last_occurrence_at,
			created_at, updated_at
		FROM alerts
		WHERE id = $1
	`

	var alert domain.Alert
	var tagsJSON, customFieldsJSON []byte

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&alert.ID,
		&alert.OrganizationID,
		&alert.Source,
		&alert.SourceID,
		&alert.Priority,
		&alert.Status,
		&alert.Message,
		&alert.Description,
		&tagsJSON,
		&customFieldsJSON,
		&alert.AssignedToUserID,
		&alert.AssignedToTeamID,
		&alert.AcknowledgedBy,
		&alert.AcknowledgedAt,
		&alert.ClosedBy,
		&alert.ClosedAt,
		&alert.CloseReason,
		&alert.SnoozedUntil,
		&alert.EscalationPolicyID,
		&alert.EscalationLevel,
		&alert.LastEscalatedAt,
		&alert.DedupKey,
		&alert.DedupCount,
		&alert.FirstOccurrenceAt,
		&alert.LastOccurrenceAt,
		&alert.CreatedAt,
		&alert.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("alert not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get alert: %w", err)
	}

	if err := json.Unmarshal(tagsJSON, &alert.Tags); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tags: %w", err)
	}

	if err := json.Unmarshal(customFieldsJSON, &alert.CustomFields); err != nil {
		return nil, fmt.Errorf("failed to unmarshal custom_fields: %w", err)
	}

	return &alert, nil
}

func (r *AlertRepository) Update(ctx context.Context, alert *domain.Alert) error {
	query := `
		UPDATE alerts
		SET
			source = $2, source_id = $3, priority = $4, status = $5,
			message = $6, description = $7, tags = $8, custom_fields = $9,
			assigned_to_user_id = $10, assigned_to_team_id = $11,
			acknowledged_by = $12, acknowledged_at = $13,
			closed_by = $14, closed_at = $15, close_reason = $16,
			snoozed_until = $17,
			escalation_policy_id = $18, escalation_level = $19, last_escalated_at = $20
		WHERE id = $1
		RETURNING updated_at
	`

	tags, err := json.Marshal(alert.Tags)
	if err != nil {
		return fmt.Errorf("failed to marshal tags: %w", err)
	}

	customFields, err := json.Marshal(alert.CustomFields)
	if err != nil {
		return fmt.Errorf("failed to marshal custom_fields: %w", err)
	}

	err = r.db.QueryRowContext(
		ctx,
		query,
		alert.ID,
		alert.Source,
		alert.SourceID,
		alert.Priority,
		alert.Status,
		alert.Message,
		alert.Description,
		tags,
		customFields,
		alert.AssignedToUserID,
		alert.AssignedToTeamID,
		alert.AcknowledgedBy,
		alert.AcknowledgedAt,
		alert.ClosedBy,
		alert.ClosedAt,
		alert.CloseReason,
		alert.SnoozedUntil,
		alert.EscalationPolicyID,
		alert.EscalationLevel,
		alert.LastEscalatedAt,
	).Scan(&alert.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to update alert: %w", err)
	}

	return nil
}

func (r *AlertRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM alerts WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete alert: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("alert not found")
	}

	return nil
}

func (r *AlertRepository) List(ctx context.Context, filter *domain.AlertFilter) ([]*domain.Alert, int, error) {
	// Build WHERE clause
	where := []string{"organization_id = $1"}
	args := []interface{}{filter.OrganizationID}
	argCount := 1

	if len(filter.Status) > 0 {
		argCount++
		statusPlaceholders := make([]string, len(filter.Status))
		for i, status := range filter.Status {
			statusPlaceholders[i] = fmt.Sprintf("$%d", argCount)
			args = append(args, status.String())
			argCount++
		}
		where = append(where, fmt.Sprintf("status IN (%s)", strings.Join(statusPlaceholders, ",")))
		argCount--
	}

	if len(filter.Priority) > 0 {
		argCount++
		priorityPlaceholders := make([]string, len(filter.Priority))
		for i, priority := range filter.Priority {
			priorityPlaceholders[i] = fmt.Sprintf("$%d", argCount)
			args = append(args, priority.String())
			argCount++
		}
		where = append(where, fmt.Sprintf("priority IN (%s)", strings.Join(priorityPlaceholders, ",")))
		argCount--
	}

	if filter.AssignedToUser != nil {
		argCount++
		where = append(where, fmt.Sprintf("assigned_to_user_id = $%d", argCount))
		args = append(args, filter.AssignedToUser)
	}

	if filter.AssignedToTeam != nil {
		argCount++
		where = append(where, fmt.Sprintf("assigned_to_team_id = $%d", argCount))
		args = append(args, filter.AssignedToTeam)
	}

	if filter.Source != nil {
		argCount++
		where = append(where, fmt.Sprintf("source = $%d", argCount))
		args = append(args, *filter.Source)
	}

	if filter.Search != nil && *filter.Search != "" {
		argCount++
		where = append(where, fmt.Sprintf("(message ILIKE $%d OR description ILIKE $%d)", argCount, argCount))
		args = append(args, "%"+*filter.Search+"%")
	}

	whereClause := strings.Join(where, " AND ")

	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM alerts WHERE %s", whereClause)
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count alerts: %w", err)
	}

	// Query alerts
	query := fmt.Sprintf(`
		SELECT
			id, organization_id, source, source_id, priority, status,
			message, description, tags, custom_fields,
			assigned_to_user_id, assigned_to_team_id,
			acknowledged_by, acknowledged_at,
			closed_by, closed_at, close_reason,
			snoozed_until,
			escalation_policy_id, escalation_level, last_escalated_at,
			dedup_key, dedup_count, first_occurrence_at, last_occurrence_at,
			created_at, updated_at
		FROM alerts
		WHERE %s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argCount+1, argCount+2)

	args = append(args, filter.Limit, filter.Offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list alerts: %w", err)
	}
	defer rows.Close()

	var alerts []*domain.Alert
	for rows.Next() {
		var alert domain.Alert
		var tagsJSON, customFieldsJSON []byte

		err := rows.Scan(
			&alert.ID,
			&alert.OrganizationID,
			&alert.Source,
			&alert.SourceID,
			&alert.Priority,
			&alert.Status,
			&alert.Message,
			&alert.Description,
			&tagsJSON,
			&customFieldsJSON,
			&alert.AssignedToUserID,
			&alert.AssignedToTeamID,
			&alert.AcknowledgedBy,
			&alert.AcknowledgedAt,
			&alert.ClosedBy,
			&alert.ClosedAt,
			&alert.CloseReason,
			&alert.SnoozedUntil,
			&alert.EscalationPolicyID,
			&alert.EscalationLevel,
			&alert.LastEscalatedAt,
			&alert.DedupKey,
			&alert.DedupCount,
			&alert.FirstOccurrenceAt,
			&alert.LastOccurrenceAt,
			&alert.CreatedAt,
			&alert.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan alert: %w", err)
		}

		if err := json.Unmarshal(tagsJSON, &alert.Tags); err != nil {
			return nil, 0, fmt.Errorf("failed to unmarshal tags: %w", err)
		}

		if err := json.Unmarshal(customFieldsJSON, &alert.CustomFields); err != nil {
			return nil, 0, fmt.Errorf("failed to unmarshal custom_fields: %w", err)
		}

		alerts = append(alerts, &alert)
	}

	return alerts, total, nil
}

func (r *AlertRepository) Acknowledge(ctx context.Context, id, userID uuid.UUID) error {
	query := `
		UPDATE alerts
		SET
			status = $2,
			acknowledged_by = $3,
			acknowledged_at = $4
		WHERE id = $1 AND status = 'open'
		RETURNING updated_at
	`

	var updatedAt time.Time
	err := r.db.QueryRowContext(
		ctx,
		query,
		id,
		domain.AlertStatusAcknowledged.String(),
		userID,
		time.Now(),
	).Scan(&updatedAt)

	if err == sql.ErrNoRows {
		return fmt.Errorf("alert not found or already acknowledged")
	}
	if err != nil {
		return fmt.Errorf("failed to acknowledge alert: %w", err)
	}

	return nil
}

func (r *AlertRepository) Close(ctx context.Context, id, userID uuid.UUID, reason string) error {
	query := `
		UPDATE alerts
		SET
			status = $2,
			closed_by = $3,
			closed_at = $4,
			close_reason = $5
		WHERE id = $1 AND status != 'closed'
		RETURNING updated_at
	`

	var updatedAt time.Time
	err := r.db.QueryRowContext(
		ctx,
		query,
		id,
		domain.AlertStatusClosed.String(),
		userID,
		time.Now(),
		reason,
	).Scan(&updatedAt)

	if err == sql.ErrNoRows {
		return fmt.Errorf("alert not found or already closed")
	}
	if err != nil {
		return fmt.Errorf("failed to close alert: %w", err)
	}

	return nil
}

func (r *AlertRepository) Snooze(ctx context.Context, id uuid.UUID, until time.Time) error {
	query := `
		UPDATE alerts
		SET
			status = $2,
			snoozed_until = $3
		WHERE id = $1 AND status != 'closed'
		RETURNING updated_at
	`

	var updatedAt time.Time
	err := r.db.QueryRowContext(
		ctx,
		query,
		id,
		domain.AlertStatusSnoozed.String(),
		until,
	).Scan(&updatedAt)

	if err == sql.ErrNoRows {
		return fmt.Errorf("alert not found or already closed")
	}
	if err != nil {
		return fmt.Errorf("failed to snooze alert: %w", err)
	}

	return nil
}

func (r *AlertRepository) Assign(ctx context.Context, id uuid.UUID, userID, teamID *uuid.UUID) error {
	query := `
		UPDATE alerts
		SET
			assigned_to_user_id = $2,
			assigned_to_team_id = $3
		WHERE id = $1
		RETURNING updated_at
	`

	var updatedAt time.Time
	err := r.db.QueryRowContext(
		ctx,
		query,
		id,
		userID,
		teamID,
	).Scan(&updatedAt)

	if err == sql.ErrNoRows {
		return fmt.Errorf("alert not found")
	}
	if err != nil {
		return fmt.Errorf("failed to assign alert: %w", err)
	}

	return nil
}

// FindByDedupKey finds an open alert with the given dedup key
func (r *AlertRepository) FindByDedupKey(ctx context.Context, orgID uuid.UUID, dedupKey string) (*domain.Alert, error) {
	query := `
		SELECT
			id, organization_id, source, source_id, priority, status,
			message, description, tags, custom_fields,
			assigned_to_user_id, assigned_to_team_id,
			acknowledged_by, acknowledged_at,
			closed_by, closed_at, close_reason,
			snoozed_until,
			escalation_policy_id, escalation_level, last_escalated_at,
			dedup_key, dedup_count, first_occurrence_at, last_occurrence_at,
			created_at, updated_at
		FROM alerts
		WHERE organization_id = $1 AND dedup_key = $2 AND status != 'closed'
		ORDER BY created_at DESC
		LIMIT 1
	`

	var alert domain.Alert
	var tagsJSON, customFieldsJSON []byte

	err := r.db.QueryRowContext(ctx, query, orgID, dedupKey).Scan(
		&alert.ID,
		&alert.OrganizationID,
		&alert.Source,
		&alert.SourceID,
		&alert.Priority,
		&alert.Status,
		&alert.Message,
		&alert.Description,
		&tagsJSON,
		&customFieldsJSON,
		&alert.AssignedToUserID,
		&alert.AssignedToTeamID,
		&alert.AcknowledgedBy,
		&alert.AcknowledgedAt,
		&alert.ClosedBy,
		&alert.ClosedAt,
		&alert.CloseReason,
		&alert.SnoozedUntil,
		&alert.EscalationPolicyID,
		&alert.EscalationLevel,
		&alert.LastEscalatedAt,
		&alert.DedupKey,
		&alert.DedupCount,
		&alert.FirstOccurrenceAt,
		&alert.LastOccurrenceAt,
		&alert.CreatedAt,
		&alert.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil // No matching alert found (not an error)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find alert by dedup key: %w", err)
	}

	if err := json.Unmarshal(tagsJSON, &alert.Tags); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tags: %w", err)
	}

	if err := json.Unmarshal(customFieldsJSON, &alert.CustomFields); err != nil {
		return nil, fmt.Errorf("failed to unmarshal custom_fields: %w", err)
	}

	return &alert, nil
}

// IncrementDedupCount increments the dedup count and updates last_occurrence_at
func (r *AlertRepository) IncrementDedupCount(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE alerts
		SET
			dedup_count = dedup_count + 1,
			last_occurrence_at = NOW()
		WHERE id = $1
		RETURNING dedup_count
	`

	var dedupCount int
	err := r.db.QueryRowContext(ctx, query, id).Scan(&dedupCount)

	if err == sql.ErrNoRows {
		return fmt.Errorf("alert not found")
	}
	if err != nil {
		return fmt.Errorf("failed to increment dedup count: %w", err)
	}

	return nil
}
