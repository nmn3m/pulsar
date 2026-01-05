package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/nmn3m/pulsar/backend/internal/domain"
	"github.com/nmn3m/pulsar/backend/internal/repository"
)

type incidentRepository struct {
	db *sqlx.DB
}

// NewIncidentRepository creates a new incident repository
func NewIncidentRepository(db *sqlx.DB) repository.IncidentRepository {
	return &incidentRepository{db: db}
}

// Create creates a new incident
func (r *incidentRepository) Create(ctx context.Context, incident *domain.Incident) error {
	query := `
		INSERT INTO incidents (
			id, organization_id, title, description, severity, status, priority,
			created_by_user_id, assigned_to_team_id, started_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
		)
		RETURNING created_at, updated_at
	`

	return r.db.QueryRowContext(
		ctx, query,
		incident.ID, incident.OrganizationID, incident.Title, incident.Description,
		incident.Severity, incident.Status, incident.Priority, incident.CreatedByUserID,
		incident.AssignedToTeamID, incident.StartedAt,
	).Scan(&incident.CreatedAt, &incident.UpdatedAt)
}

// GetByID retrieves an incident by ID
func (r *incidentRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Incident, error) {
	var incident domain.Incident
	query := `SELECT * FROM incidents WHERE id = $1`

	err := r.db.GetContext(ctx, &incident, query, id)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &incident, nil
}

// Update updates an incident
func (r *incidentRepository) Update(ctx context.Context, incident *domain.Incident) error {
	query := `
		UPDATE incidents SET
			title = $1,
			description = $2,
			severity = $3,
			status = $4,
			priority = $5,
			assigned_to_team_id = $6,
			resolved_at = $7
		WHERE id = $8
		RETURNING updated_at
	`

	return r.db.QueryRowContext(
		ctx, query,
		incident.Title, incident.Description, incident.Severity, incident.Status,
		incident.Priority, incident.AssignedToTeamID, incident.ResolvedAt, incident.ID,
	).Scan(&incident.UpdatedAt)
}

// Delete deletes an incident
func (r *incidentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM incidents WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return domain.ErrNotFound
	}

	return nil
}

// List retrieves incidents with filtering and pagination
func (r *incidentRepository) List(ctx context.Context, filter *domain.IncidentFilter) ([]*domain.Incident, int, error) {
	if err := filter.Validate(); err != nil {
		return nil, 0, err
	}

	where := []string{"organization_id = $1"}
	args := []interface{}{filter.OrganizationID}
	argCount := 1

	// Filter by status
	if len(filter.Status) > 0 {
		argCount++
		placeholders := make([]string, len(filter.Status))
		for i, status := range filter.Status {
			args = append(args, status)
			placeholders[i] = fmt.Sprintf("$%d", argCount+i)
		}
		where = append(where, fmt.Sprintf("status IN (%s)", strings.Join(placeholders, ",")))
		argCount += len(filter.Status) - 1
	}

	// Filter by severity
	if len(filter.Severity) > 0 {
		argCount++
		placeholders := make([]string, len(filter.Severity))
		for i, severity := range filter.Severity {
			args = append(args, severity)
			placeholders[i] = fmt.Sprintf("$%d", argCount+i)
		}
		where = append(where, fmt.Sprintf("severity IN (%s)", strings.Join(placeholders, ",")))
		argCount += len(filter.Severity) - 1
	}

	// Filter by assigned team
	if filter.AssignedToTeamID != nil {
		argCount++
		where = append(where, fmt.Sprintf("assigned_to_team_id = $%d", argCount))
		args = append(args, *filter.AssignedToTeamID)
	}

	// Search in title and description
	if filter.Search != nil && *filter.Search != "" {
		argCount++
		searchPattern := "%" + *filter.Search + "%"
		where = append(where, fmt.Sprintf("(title ILIKE $%d OR description ILIKE $%d)", argCount, argCount))
		args = append(args, searchPattern)
	}

	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM incidents WHERE %s", strings.Join(where, " AND "))
	var total int
	if err := r.db.GetContext(ctx, &total, countQuery, args...); err != nil {
		return nil, 0, err
	}

	// Get incidents
	argCount++
	args = append(args, filter.Limit)
	argCount++
	args = append(args, filter.Offset)

	query := fmt.Sprintf(`
		SELECT * FROM incidents
		WHERE %s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, strings.Join(where, " AND "), argCount-1, argCount)

	var incidents []*domain.Incident
	if err := r.db.SelectContext(ctx, &incidents, query, args...); err != nil {
		return nil, 0, err
	}

	return incidents, total, nil
}

// AddResponder adds a responder to an incident
func (r *incidentRepository) AddResponder(ctx context.Context, responder *domain.IncidentResponder) error {
	query := `
		INSERT INTO incident_responders (id, incident_id, user_id, role)
		VALUES ($1, $2, $3, $4)
		RETURNING added_at
	`

	return r.db.QueryRowContext(
		ctx, query,
		responder.ID, responder.IncidentID, responder.UserID, responder.Role,
	).Scan(&responder.AddedAt)
}

// RemoveResponder removes a responder from an incident
func (r *incidentRepository) RemoveResponder(ctx context.Context, incidentID, userID uuid.UUID) error {
	query := `DELETE FROM incident_responders WHERE incident_id = $1 AND user_id = $2`
	result, err := r.db.ExecContext(ctx, query, incidentID, userID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return domain.ErrNotFound
	}

	return nil
}

// UpdateResponderRole updates a responder's role
func (r *incidentRepository) UpdateResponderRole(ctx context.Context, incidentID, userID uuid.UUID, role domain.ResponderRole) error {
	query := `UPDATE incident_responders SET role = $1 WHERE incident_id = $2 AND user_id = $3`
	result, err := r.db.ExecContext(ctx, query, role, incidentID, userID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return domain.ErrNotFound
	}

	return nil
}

// ListResponders retrieves all responders for an incident with user details
func (r *incidentRepository) ListResponders(ctx context.Context, incidentID uuid.UUID) ([]*domain.ResponderWithUser, error) {
	query := `
		SELECT
			ir.id, ir.incident_id, ir.user_id, ir.role, ir.added_at,
			u.id, u.email, u.username, u.full_name, u.created_at, u.updated_at
		FROM incident_responders ir
		JOIN users u ON ir.user_id = u.id
		WHERE ir.incident_id = $1
		ORDER BY ir.added_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query, incidentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var responders []*domain.ResponderWithUser
	for rows.Next() {
		var r domain.ResponderWithUser
		r.User = &domain.User{}

		err := rows.Scan(
			&r.ID, &r.IncidentID, &r.UserID, &r.Role, &r.AddedAt,
			&r.User.ID, &r.User.Email, &r.User.Username, &r.User.FullName,
			&r.User.CreatedAt, &r.User.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Clear password hash for security
		r.User.PasswordHash = ""

		responders = append(responders, &r)
	}

	return responders, nil
}

// AddTimelineEvent adds an event to the incident timeline
func (r *incidentRepository) AddTimelineEvent(ctx context.Context, event *domain.IncidentTimelineEvent) error {
	metadata, err := json.Marshal(event.Metadata)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO incident_timeline (id, incident_id, event_type, user_id, description, metadata)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at
	`

	return r.db.QueryRowContext(
		ctx, query,
		event.ID, event.IncidentID, event.EventType, event.UserID, event.Description, metadata,
	).Scan(&event.CreatedAt)
}

// GetTimeline retrieves the timeline for an incident with user details
func (r *incidentRepository) GetTimeline(ctx context.Context, incidentID uuid.UUID) ([]*domain.TimelineEventWithUser, error) {
	query := `
		SELECT
			t.id, t.incident_id, t.event_type, t.user_id, t.description, t.metadata, t.created_at,
			u.id, u.email, u.username, u.full_name, u.created_at, u.updated_at
		FROM incident_timeline t
		LEFT JOIN users u ON t.user_id = u.id
		WHERE t.incident_id = $1
		ORDER BY t.created_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query, incidentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var timeline []*domain.TimelineEventWithUser
	for rows.Next() {
		var event domain.TimelineEventWithUser
		var metadataJSON []byte
		var userID sql.NullString
		var userEmail sql.NullString
		var userUsername sql.NullString
		var userFullName sql.NullString
		var userCreatedAt sql.NullTime
		var userUpdatedAt sql.NullTime

		err := rows.Scan(
			&event.ID, &event.IncidentID, &event.EventType, &event.UserID,
			&event.Description, &metadataJSON, &event.CreatedAt,
			&userID, &userEmail, &userUsername, &userFullName,
			&userCreatedAt, &userUpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Parse metadata
		if err := json.Unmarshal(metadataJSON, &event.Metadata); err != nil {
			event.Metadata = make(map[string]interface{})
		}

		// Attach user if present
		if userID.Valid {
			event.User = &domain.User{
				Email:        userEmail.String,
				Username:     userUsername.String,
				FullName:     &userFullName.String,
				CreatedAt:    userCreatedAt.Time,
				UpdatedAt:    userUpdatedAt.Time,
				PasswordHash: "",
			}
			userUUID, err := uuid.Parse(userID.String)
			if err == nil {
				event.User.ID = userUUID
			}
		}

		timeline = append(timeline, &event)
	}

	return timeline, nil
}

// LinkAlert links an alert to an incident
func (r *incidentRepository) LinkAlert(ctx context.Context, link *domain.IncidentAlert) error {
	query := `
		INSERT INTO incident_alerts (id, incident_id, alert_id, linked_by_user_id)
		VALUES ($1, $2, $3, $4)
		RETURNING linked_at
	`

	return r.db.QueryRowContext(
		ctx, query,
		link.ID, link.IncidentID, link.AlertID, link.LinkedByUserID,
	).Scan(&link.LinkedAt)
}

// UnlinkAlert unlinks an alert from an incident
func (r *incidentRepository) UnlinkAlert(ctx context.Context, incidentID, alertID uuid.UUID) error {
	query := `DELETE FROM incident_alerts WHERE incident_id = $1 AND alert_id = $2`
	result, err := r.db.ExecContext(ctx, query, incidentID, alertID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return domain.ErrNotFound
	}

	return nil
}

// ListAlerts retrieves all alerts linked to an incident
func (r *incidentRepository) ListAlerts(ctx context.Context, incidentID uuid.UUID) ([]*domain.IncidentAlertWithDetails, error) {
	query := `
		SELECT
			ia.id, ia.incident_id, ia.alert_id, ia.linked_at, ia.linked_by_user_id,
			a.id, a.organization_id, a.source, a.source_id, a.priority, a.status,
			a.message, a.description, a.tags, a.custom_fields,
			a.assigned_to_user_id, a.assigned_to_team_id, a.escalation_policy_id,
			a.escalation_level, a.acknowledged_at, a.acknowledged_by_user_id,
			a.closed_at, a.closed_by_user_id, a.close_reason, a.snoozed_until,
			a.last_escalated_at, a.created_at, a.updated_at
		FROM incident_alerts ia
		JOIN alerts a ON ia.alert_id = a.id
		WHERE ia.incident_id = $1
		ORDER BY ia.linked_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, incidentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alerts []*domain.IncidentAlertWithDetails
	for rows.Next() {
		var ia domain.IncidentAlertWithDetails
		ia.Alert = &domain.Alert{}

		var tagsJSON []byte
		var customFieldsJSON []byte

		err := rows.Scan(
			&ia.ID, &ia.IncidentID, &ia.AlertID, &ia.LinkedAt, &ia.LinkedByUserID,
			&ia.Alert.ID, &ia.Alert.OrganizationID, &ia.Alert.Source, &ia.Alert.SourceID,
			&ia.Alert.Priority, &ia.Alert.Status, &ia.Alert.Message, &ia.Alert.Description,
			&tagsJSON, &customFieldsJSON,
			&ia.Alert.AssignedToUserID, &ia.Alert.AssignedToTeamID, &ia.Alert.EscalationPolicyID,
			&ia.Alert.EscalationLevel, &ia.Alert.AcknowledgedAt, &ia.Alert.AcknowledgedBy,
			&ia.Alert.ClosedAt, &ia.Alert.ClosedBy, &ia.Alert.CloseReason,
			&ia.Alert.SnoozedUntil, &ia.Alert.LastEscalatedAt, &ia.Alert.CreatedAt, &ia.Alert.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Parse JSON fields
		if err := json.Unmarshal(tagsJSON, &ia.Alert.Tags); err != nil {
			ia.Alert.Tags = []string{}
		}
		if err := json.Unmarshal(customFieldsJSON, &ia.Alert.CustomFields); err != nil {
			ia.Alert.CustomFields = make(map[string]interface{})
		}

		alerts = append(alerts, &ia)
	}

	return alerts, nil
}

// GetWithDetails retrieves an incident with all related data
func (r *incidentRepository) GetWithDetails(ctx context.Context, id uuid.UUID) (*domain.IncidentWithDetails, error) {
	// Get base incident
	incident, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	result := &domain.IncidentWithDetails{
		Incident: *incident,
	}

	// Get responders
	responders, err := r.ListResponders(ctx, id)
	if err != nil {
		return nil, err
	}
	result.Responders = responders

	// Get alerts
	alerts, err := r.ListAlerts(ctx, id)
	if err != nil {
		return nil, err
	}
	result.Alerts = alerts

	// Get timeline
	timeline, err := r.GetTimeline(ctx, id)
	if err != nil {
		return nil, err
	}
	result.Timeline = timeline

	return result, nil
}
