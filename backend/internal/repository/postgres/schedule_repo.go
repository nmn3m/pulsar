package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/pulsar/backend/internal/domain"
)

type ScheduleRepository struct {
	db *DB
}

func NewScheduleRepository(db *DB) *ScheduleRepository {
	return &ScheduleRepository{db: db}
}

// Schedule CRUD operations

func (r *ScheduleRepository) Create(ctx context.Context, schedule *domain.Schedule) error {
	query := `
		INSERT INTO schedules (id, organization_id, team_id, name, description, timezone)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at, updated_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		schedule.ID,
		schedule.OrganizationID,
		schedule.TeamID,
		schedule.Name,
		schedule.Description,
		schedule.Timezone,
	).Scan(&schedule.CreatedAt, &schedule.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create schedule: %w", err)
	}

	return nil
}

func (r *ScheduleRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Schedule, error) {
	query := `
		SELECT id, organization_id, team_id, name, description, timezone, created_at, updated_at
		FROM schedules
		WHERE id = $1
	`

	var schedule domain.Schedule
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&schedule.ID,
		&schedule.OrganizationID,
		&schedule.TeamID,
		&schedule.Name,
		&schedule.Description,
		&schedule.Timezone,
		&schedule.CreatedAt,
		&schedule.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("schedule not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get schedule: %w", err)
	}

	return &schedule, nil
}

func (r *ScheduleRepository) Update(ctx context.Context, schedule *domain.Schedule) error {
	query := `
		UPDATE schedules
		SET name = $2, description = $3, timezone = $4, team_id = $5
		WHERE id = $1
		RETURNING updated_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		schedule.ID,
		schedule.Name,
		schedule.Description,
		schedule.Timezone,
		schedule.TeamID,
	).Scan(&schedule.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to update schedule: %w", err)
	}

	return nil
}

func (r *ScheduleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM schedules WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete schedule: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("schedule not found")
	}

	return nil
}

func (r *ScheduleRepository) List(ctx context.Context, orgID uuid.UUID, limit, offset int) ([]*domain.Schedule, error) {
	query := `
		SELECT id, organization_id, team_id, name, description, timezone, created_at, updated_at
		FROM schedules
		WHERE organization_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, orgID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list schedules: %w", err)
	}
	defer rows.Close()

	var schedules []*domain.Schedule
	for rows.Next() {
		var schedule domain.Schedule
		err := rows.Scan(
			&schedule.ID,
			&schedule.OrganizationID,
			&schedule.TeamID,
			&schedule.Name,
			&schedule.Description,
			&schedule.Timezone,
			&schedule.CreatedAt,
			&schedule.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan schedule: %w", err)
		}

		schedules = append(schedules, &schedule)
	}

	return schedules, nil
}

func (r *ScheduleRepository) GetWithRotations(ctx context.Context, id uuid.UUID) (*domain.ScheduleWithRotations, error) {
	schedule, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	rotations, err := r.ListRotations(ctx, id)
	if err != nil {
		return nil, err
	}

	return &domain.ScheduleWithRotations{
		Schedule:  *schedule,
		Rotations: rotations,
	}, nil
}

// Rotation CRUD operations

func (r *ScheduleRepository) CreateRotation(ctx context.Context, rotation *domain.ScheduleRotation) error {
	query := `
		INSERT INTO schedule_rotations (
			id, schedule_id, name, rotation_type, rotation_length,
			start_date, start_time, end_time, handoff_day, handoff_time
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING created_at, updated_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		rotation.ID,
		rotation.ScheduleID,
		rotation.Name,
		rotation.RotationType.String(),
		rotation.RotationLength,
		rotation.StartDate,
		rotation.StartTime,
		rotation.EndTime,
		rotation.HandoffDay,
		rotation.HandoffTime,
	).Scan(&rotation.CreatedAt, &rotation.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create rotation: %w", err)
	}

	return nil
}

func (r *ScheduleRepository) GetRotation(ctx context.Context, id uuid.UUID) (*domain.ScheduleRotation, error) {
	query := `
		SELECT id, schedule_id, name, rotation_type, rotation_length,
		       start_date, start_time, end_time, handoff_day, handoff_time,
		       created_at, updated_at
		FROM schedule_rotations
		WHERE id = $1
	`

	var rotation domain.ScheduleRotation
	var rotationType string

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&rotation.ID,
		&rotation.ScheduleID,
		&rotation.Name,
		&rotationType,
		&rotation.RotationLength,
		&rotation.StartDate,
		&rotation.StartTime,
		&rotation.EndTime,
		&rotation.HandoffDay,
		&rotation.HandoffTime,
		&rotation.CreatedAt,
		&rotation.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("rotation not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get rotation: %w", err)
	}

	rotation.RotationType = domain.RotationType(rotationType)

	return &rotation, nil
}

func (r *ScheduleRepository) UpdateRotation(ctx context.Context, rotation *domain.ScheduleRotation) error {
	query := `
		UPDATE schedule_rotations
		SET name = $2, rotation_type = $3, rotation_length = $4,
		    start_date = $5, start_time = $6, end_time = $7,
		    handoff_day = $8, handoff_time = $9
		WHERE id = $1
		RETURNING updated_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		rotation.ID,
		rotation.Name,
		rotation.RotationType.String(),
		rotation.RotationLength,
		rotation.StartDate,
		rotation.StartTime,
		rotation.EndTime,
		rotation.HandoffDay,
		rotation.HandoffTime,
	).Scan(&rotation.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to update rotation: %w", err)
	}

	return nil
}

func (r *ScheduleRepository) DeleteRotation(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM schedule_rotations WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete rotation: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("rotation not found")
	}

	return nil
}

func (r *ScheduleRepository) ListRotations(ctx context.Context, scheduleID uuid.UUID) ([]*domain.ScheduleRotation, error) {
	query := `
		SELECT id, schedule_id, name, rotation_type, rotation_length,
		       start_date, start_time, end_time, handoff_day, handoff_time,
		       created_at, updated_at
		FROM schedule_rotations
		WHERE schedule_id = $1
		ORDER BY created_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query, scheduleID)
	if err != nil {
		return nil, fmt.Errorf("failed to list rotations: %w", err)
	}
	defer rows.Close()

	var rotations []*domain.ScheduleRotation
	for rows.Next() {
		var rotation domain.ScheduleRotation
		var rotationType string

		err := rows.Scan(
			&rotation.ID,
			&rotation.ScheduleID,
			&rotation.Name,
			&rotationType,
			&rotation.RotationLength,
			&rotation.StartDate,
			&rotation.StartTime,
			&rotation.EndTime,
			&rotation.HandoffDay,
			&rotation.HandoffTime,
			&rotation.CreatedAt,
			&rotation.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan rotation: %w", err)
		}

		rotation.RotationType = domain.RotationType(rotationType)
		rotations = append(rotations, &rotation)
	}

	return rotations, nil
}

// Rotation participant operations

func (r *ScheduleRepository) AddParticipant(ctx context.Context, participant *domain.ScheduleRotationParticipant) error {
	query := `
		INSERT INTO schedule_rotation_participants (id, rotation_id, user_id, position)
		VALUES ($1, $2, $3, $4)
		RETURNING created_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		participant.ID,
		participant.RotationID,
		participant.UserID,
		participant.Position,
	).Scan(&participant.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to add participant: %w", err)
	}

	return nil
}

func (r *ScheduleRepository) RemoveParticipant(ctx context.Context, rotationID, userID uuid.UUID) error {
	query := `DELETE FROM schedule_rotation_participants WHERE rotation_id = $1 AND user_id = $2`

	result, err := r.db.ExecContext(ctx, query, rotationID, userID)
	if err != nil {
		return fmt.Errorf("failed to remove participant: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("participant not found")
	}

	return nil
}

func (r *ScheduleRepository) ListParticipants(ctx context.Context, rotationID uuid.UUID) ([]*domain.ParticipantWithUser, error) {
	query := `
		SELECT p.id, p.rotation_id, p.user_id, p.position, p.created_at,
		       u.id, u.email, u.username, u.full_name, u.phone, u.timezone,
		       u.notification_preferences, u.is_active, u.created_at, u.updated_at
		FROM schedule_rotation_participants p
		JOIN users u ON p.user_id = u.id
		WHERE p.rotation_id = $1
		ORDER BY p.position ASC
	`

	rows, err := r.db.QueryContext(ctx, query, rotationID)
	if err != nil {
		return nil, fmt.Errorf("failed to list participants: %w", err)
	}
	defer rows.Close()

	var participants []*domain.ParticipantWithUser
	for rows.Next() {
		var p domain.ParticipantWithUser
		var prefsJSON []byte

		err := rows.Scan(
			&p.ID,
			&p.RotationID,
			&p.UserID,
			&p.Position,
			&p.CreatedAt,
			&p.User.ID,
			&p.User.Email,
			&p.User.Username,
			&p.User.FullName,
			&p.User.Phone,
			&p.User.Timezone,
			&prefsJSON,
			&p.User.IsActive,
			&p.User.CreatedAt,
			&p.User.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan participant: %w", err)
		}

		if err := json.Unmarshal(prefsJSON, &p.User.NotificationPreferences); err != nil {
			return nil, fmt.Errorf("failed to unmarshal notification preferences: %w", err)
		}

		// Clear password hash
		p.User.PasswordHash = ""

		participants = append(participants, &p)
	}

	return participants, nil
}

func (r *ScheduleRepository) ReorderParticipants(ctx context.Context, rotationID uuid.UUID, userIDs []uuid.UUID) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Update each participant's position
	for i, userID := range userIDs {
		query := `
			UPDATE schedule_rotation_participants
			SET position = $1
			WHERE rotation_id = $2 AND user_id = $3
		`

		_, err := tx.ExecContext(ctx, query, i, rotationID, userID)
		if err != nil {
			return fmt.Errorf("failed to update participant position: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// Override operations

func (r *ScheduleRepository) CreateOverride(ctx context.Context, override *domain.ScheduleOverride) error {
	query := `
		INSERT INTO schedule_overrides (id, schedule_id, user_id, start_time, end_time, note)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at, updated_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		override.ID,
		override.ScheduleID,
		override.UserID,
		override.StartTime,
		override.EndTime,
		override.Note,
	).Scan(&override.CreatedAt, &override.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create override: %w", err)
	}

	return nil
}

func (r *ScheduleRepository) GetOverride(ctx context.Context, id uuid.UUID) (*domain.ScheduleOverride, error) {
	query := `
		SELECT id, schedule_id, user_id, start_time, end_time, note, created_at, updated_at
		FROM schedule_overrides
		WHERE id = $1
	`

	var override domain.ScheduleOverride
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&override.ID,
		&override.ScheduleID,
		&override.UserID,
		&override.StartTime,
		&override.EndTime,
		&override.Note,
		&override.CreatedAt,
		&override.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("override not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get override: %w", err)
	}

	return &override, nil
}

func (r *ScheduleRepository) UpdateOverride(ctx context.Context, override *domain.ScheduleOverride) error {
	query := `
		UPDATE schedule_overrides
		SET user_id = $2, start_time = $3, end_time = $4, note = $5
		WHERE id = $1
		RETURNING updated_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		override.ID,
		override.UserID,
		override.StartTime,
		override.EndTime,
		override.Note,
	).Scan(&override.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to update override: %w", err)
	}

	return nil
}

func (r *ScheduleRepository) DeleteOverride(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM schedule_overrides WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete override: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("override not found")
	}

	return nil
}

func (r *ScheduleRepository) ListOverrides(ctx context.Context, scheduleID uuid.UUID, start, end time.Time) ([]*domain.ScheduleOverride, error) {
	query := `
		SELECT id, schedule_id, user_id, start_time, end_time, note, created_at, updated_at
		FROM schedule_overrides
		WHERE schedule_id = $1
		  AND start_time < $3
		  AND end_time > $2
		ORDER BY start_time ASC
	`

	rows, err := r.db.QueryContext(ctx, query, scheduleID, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to list overrides: %w", err)
	}
	defer rows.Close()

	var overrides []*domain.ScheduleOverride
	for rows.Next() {
		var override domain.ScheduleOverride
		err := rows.Scan(
			&override.ID,
			&override.ScheduleID,
			&override.UserID,
			&override.StartTime,
			&override.EndTime,
			&override.Note,
			&override.CreatedAt,
			&override.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan override: %w", err)
		}

		overrides = append(overrides, &override)
	}

	return overrides, nil
}

// On-call calculation

func (r *ScheduleRepository) GetOnCallUser(ctx context.Context, scheduleID uuid.UUID, at time.Time) (*domain.OnCallUser, error) {
	// First, check for overrides
	query := `
		SELECT user_id, start_time, end_time
		FROM schedule_overrides
		WHERE schedule_id = $1
		  AND start_time <= $2
		  AND end_time > $2
		ORDER BY created_at DESC
		LIMIT 1
	`

	var onCall domain.OnCallUser
	err := r.db.QueryRowContext(ctx, query, scheduleID, at).Scan(
		&onCall.UserID,
		&onCall.StartTime,
		&onCall.EndTime,
	)

	if err == nil {
		// Override found
		onCall.ScheduleID = scheduleID
		onCall.IsOverride = true
		return &onCall, nil
	}

	if err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to check overrides: %w", err)
	}

	// No override, calculate from rotation
	// This is a simplified implementation - in production, you'd need more complex logic
	// to handle different rotation types, timezones, and edge cases

	// For now, return nil to indicate no one is on-call (will be implemented in service layer)
	return nil, fmt.Errorf("rotation-based on-call calculation not yet implemented in repository")
}
