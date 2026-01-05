package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/nmn3m/pulsar/backend/internal/domain"
)

type NotificationRepository struct {
	db *sqlx.DB
}

func NewNotificationRepository(db *sqlx.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

// ==================== Notification Channels ====================

func (r *NotificationRepository) CreateChannel(ctx context.Context, channel *domain.NotificationChannel) error {
	query := `
		INSERT INTO notification_channels (organization_id, name, channel_type, is_enabled, config)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRowContext(
		ctx,
		query,
		channel.OrganizationID,
		channel.Name,
		channel.ChannelType,
		channel.IsEnabled,
		channel.Config,
	).Scan(&channel.ID, &channel.CreatedAt, &channel.UpdatedAt)
}

func (r *NotificationRepository) GetChannelByID(ctx context.Context, id uuid.UUID) (*domain.NotificationChannel, error) {
	var channel domain.NotificationChannel
	query := `SELECT * FROM notification_channels WHERE id = $1`

	err := r.db.GetContext(ctx, &channel, query, id)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &channel, nil
}

func (r *NotificationRepository) ListChannels(ctx context.Context, orgID uuid.UUID) ([]domain.NotificationChannel, error) {
	var channels []domain.NotificationChannel
	query := `
		SELECT * FROM notification_channels
		WHERE organization_id = $1
		ORDER BY name ASC
	`

	err := r.db.SelectContext(ctx, &channels, query, orgID)
	if err != nil {
		return nil, err
	}

	if channels == nil {
		channels = []domain.NotificationChannel{}
	}

	return channels, nil
}

func (r *NotificationRepository) UpdateChannel(ctx context.Context, channel *domain.NotificationChannel) error {
	query := `
		UPDATE notification_channels
		SET name = $1, channel_type = $2, is_enabled = $3, config = $4
		WHERE id = $5
		RETURNING updated_at
	`
	err := r.db.QueryRowContext(
		ctx,
		query,
		channel.Name,
		channel.ChannelType,
		channel.IsEnabled,
		channel.Config,
		channel.ID,
	).Scan(&channel.UpdatedAt)

	if err == sql.ErrNoRows {
		return domain.ErrNotFound
	}

	return err
}

func (r *NotificationRepository) DeleteChannel(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM notification_channels WHERE id = $1`
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

// ==================== User Notification Preferences ====================

func (r *NotificationRepository) CreatePreference(ctx context.Context, pref *domain.UserNotificationPreference) error {
	query := `
		INSERT INTO user_notification_preferences
		(user_id, channel_id, is_enabled, dnd_enabled, dnd_start_time, dnd_end_time, min_priority)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRowContext(
		ctx,
		query,
		pref.UserID,
		pref.ChannelID,
		pref.IsEnabled,
		pref.DNDEnabled,
		pref.DNDStartTime,
		pref.DNDEndTime,
		pref.MinPriority,
	).Scan(&pref.ID, &pref.CreatedAt, &pref.UpdatedAt)
}

func (r *NotificationRepository) GetPreferenceByID(ctx context.Context, id uuid.UUID) (*domain.UserNotificationPreference, error) {
	var pref domain.UserNotificationPreference
	query := `SELECT * FROM user_notification_preferences WHERE id = $1`

	err := r.db.GetContext(ctx, &pref, query, id)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &pref, nil
}

func (r *NotificationRepository) GetPreferenceByUserAndChannel(ctx context.Context, userID, channelID uuid.UUID) (*domain.UserNotificationPreference, error) {
	var pref domain.UserNotificationPreference
	query := `SELECT * FROM user_notification_preferences WHERE user_id = $1 AND channel_id = $2`

	err := r.db.GetContext(ctx, &pref, query, userID, channelID)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &pref, nil
}

func (r *NotificationRepository) ListPreferencesByUser(ctx context.Context, userID uuid.UUID) ([]domain.UserNotificationPreference, error) {
	var prefs []domain.UserNotificationPreference
	query := `
		SELECT * FROM user_notification_preferences
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	err := r.db.SelectContext(ctx, &prefs, query, userID)
	if err != nil {
		return nil, err
	}

	if prefs == nil {
		prefs = []domain.UserNotificationPreference{}
	}

	return prefs, nil
}

func (r *NotificationRepository) UpdatePreference(ctx context.Context, pref *domain.UserNotificationPreference) error {
	query := `
		UPDATE user_notification_preferences
		SET is_enabled = $1, dnd_enabled = $2, dnd_start_time = $3, dnd_end_time = $4, min_priority = $5
		WHERE id = $6
		RETURNING updated_at
	`
	err := r.db.QueryRowContext(
		ctx,
		query,
		pref.IsEnabled,
		pref.DNDEnabled,
		pref.DNDStartTime,
		pref.DNDEndTime,
		pref.MinPriority,
		pref.ID,
	).Scan(&pref.UpdatedAt)

	if err == sql.ErrNoRows {
		return domain.ErrNotFound
	}

	return err
}

func (r *NotificationRepository) DeletePreference(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM user_notification_preferences WHERE id = $1`
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

// ==================== Notification Logs ====================

func (r *NotificationRepository) CreateLog(ctx context.Context, log *domain.NotificationLog) error {
	query := `
		INSERT INTO notification_logs
		(organization_id, channel_id, user_id, alert_id, recipient, subject, message, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at
	`
	return r.db.QueryRowContext(
		ctx,
		query,
		log.OrganizationID,
		log.ChannelID,
		log.UserID,
		log.AlertID,
		log.Recipient,
		log.Subject,
		log.Message,
		log.Status,
	).Scan(&log.ID, &log.CreatedAt)
}

func (r *NotificationRepository) GetLogByID(ctx context.Context, id uuid.UUID) (*domain.NotificationLog, error) {
	var log domain.NotificationLog
	query := `SELECT * FROM notification_logs WHERE id = $1`

	err := r.db.GetContext(ctx, &log, query, id)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &log, nil
}

func (r *NotificationRepository) ListLogs(ctx context.Context, orgID uuid.UUID, limit, offset int) ([]domain.NotificationLog, error) {
	var logs []domain.NotificationLog
	query := `
		SELECT * FROM notification_logs
		WHERE organization_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	err := r.db.SelectContext(ctx, &logs, query, orgID, limit, offset)
	if err != nil {
		return nil, err
	}

	if logs == nil {
		logs = []domain.NotificationLog{}
	}

	return logs, nil
}

func (r *NotificationRepository) ListLogsByAlert(ctx context.Context, alertID uuid.UUID) ([]domain.NotificationLog, error) {
	var logs []domain.NotificationLog
	query := `
		SELECT * FROM notification_logs
		WHERE alert_id = $1
		ORDER BY created_at DESC
	`

	err := r.db.SelectContext(ctx, &logs, query, alertID)
	if err != nil {
		return nil, err
	}

	if logs == nil {
		logs = []domain.NotificationLog{}
	}

	return logs, nil
}

func (r *NotificationRepository) ListLogsByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]domain.NotificationLog, error) {
	var logs []domain.NotificationLog
	query := `
		SELECT * FROM notification_logs
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	err := r.db.SelectContext(ctx, &logs, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}

	if logs == nil {
		logs = []domain.NotificationLog{}
	}

	return logs, nil
}

func (r *NotificationRepository) GetPendingNotifications(ctx context.Context, limit int) ([]domain.NotificationLog, error) {
	var logs []domain.NotificationLog
	query := `
		SELECT * FROM notification_logs
		WHERE status = $1
		ORDER BY created_at ASC
		LIMIT $2
	`

	err := r.db.SelectContext(ctx, &logs, query, domain.NotificationStatusPending, limit)
	if err != nil {
		return nil, err
	}

	if logs == nil {
		logs = []domain.NotificationLog{}
	}

	return logs, nil
}

func (r *NotificationRepository) UpdateLogStatus(ctx context.Context, id uuid.UUID, status domain.NotificationStatus, errorMsg *string) error {
	var query string
	var err error

	if status == domain.NotificationStatusSent {
		query = `
			UPDATE notification_logs
			SET status = $1, sent_at = NOW(), error_message = NULL
			WHERE id = $2
		`
		_, err = r.db.ExecContext(ctx, query, status, id)
	} else {
		query = `
			UPDATE notification_logs
			SET status = $1, error_message = $2
			WHERE id = $3
		`
		_, err = r.db.ExecContext(ctx, query, status, errorMsg, id)
	}

	if err == sql.ErrNoRows {
		return domain.ErrNotFound
	}

	return err
}

func (r *NotificationRepository) CountLogsByStatus(ctx context.Context, orgID uuid.UUID, status domain.NotificationStatus) (int, error) {
	var count int
	query := `
		SELECT COUNT(*) FROM notification_logs
		WHERE organization_id = $1 AND status = $2
	`

	err := r.db.GetContext(ctx, &count, query, orgID, status)
	return count, err
}

// Helper function to check if a user is in DND mode at a specific time
func (r *NotificationRepository) IsUserInDND(ctx context.Context, userID, channelID uuid.UUID) (bool, error) {
	query := `
		SELECT
			CASE
				WHEN dnd_enabled = true
					AND dnd_start_time IS NOT NULL
					AND dnd_end_time IS NOT NULL
					AND CURRENT_TIME BETWEEN dnd_start_time AND dnd_end_time
				THEN true
				ELSE false
			END as in_dnd
		FROM user_notification_preferences
		WHERE user_id = $1 AND channel_id = $2 AND is_enabled = true
	`

	var inDND bool
	err := r.db.GetContext(ctx, &inDND, query, userID, channelID)
	if err == sql.ErrNoRows {
		return false, nil // No preference set, not in DND
	}
	if err != nil {
		return false, fmt.Errorf("failed to check DND status: %w", err)
	}

	return inDND, nil
}
