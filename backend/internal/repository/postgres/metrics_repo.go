package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/nmn3m/pulsar/backend/internal/domain"
	"github.com/nmn3m/pulsar/backend/internal/repository"
)

type metricsRepository struct {
	db *sqlx.DB
}

func NewMetricsRepository(db *sqlx.DB) repository.MetricsRepository {
	return &metricsRepository{db: db}
}

func (r *metricsRepository) GetAlertMetrics(ctx context.Context, orgID uuid.UUID, filter *domain.MetricsFilter) (*domain.AlertMetrics, error) {
	metrics := &domain.AlertMetrics{
		ByPriority: make(map[string]int64),
		BySource:   make(map[string]int64),
	}

	// Build time filter
	startTime := time.Now().AddDate(0, 0, -30) // Default: last 30 days
	endTime := time.Now()
	if filter != nil {
		if filter.StartTime != nil {
			startTime = *filter.StartTime
		}
		if filter.EndTime != nil {
			endTime = *filter.EndTime
		}
	}

	// Get total and status counts
	query := `
		SELECT
			COUNT(*) as total,
			COUNT(*) FILTER (WHERE status = 'open') as open,
			COUNT(*) FILTER (WHERE status = 'acknowledged') as acknowledged,
			COUNT(*) FILTER (WHERE status = 'closed') as closed,
			COUNT(*) FILTER (WHERE status = 'snoozed') as snoozed
		FROM alerts
		WHERE organization_id = $1 AND created_at >= $2 AND created_at <= $3
	`

	err := r.db.QueryRowContext(ctx, query, orgID, startTime, endTime).Scan(
		&metrics.Total,
		&metrics.Open,
		&metrics.Acknowledged,
		&metrics.Closed,
		&metrics.Snoozed,
	)
	if err != nil {
		return nil, err
	}

	// Get counts by priority
	priorityQuery := `
		SELECT priority, COUNT(*) as count
		FROM alerts
		WHERE organization_id = $1 AND created_at >= $2 AND created_at <= $3
		GROUP BY priority
	`
	rows, err := r.db.QueryContext(ctx, priorityQuery, orgID, startTime, endTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var priority string
		var count int64
		if err := rows.Scan(&priority, &count); err != nil {
			return nil, err
		}
		metrics.ByPriority[priority] = count
	}

	// Get counts by source
	sourceQuery := `
		SELECT source, COUNT(*) as count
		FROM alerts
		WHERE organization_id = $1 AND created_at >= $2 AND created_at <= $3
		GROUP BY source
	`
	rows, err = r.db.QueryContext(ctx, sourceQuery, orgID, startTime, endTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var source string
		var count int64
		if err := rows.Scan(&source, &count); err != nil {
			return nil, err
		}
		metrics.BySource[source] = count
	}

	// Get average response time (time to acknowledge)
	avgResponseQuery := `
		SELECT AVG(EXTRACT(EPOCH FROM (acknowledged_at - created_at)))
		FROM alerts
		WHERE organization_id = $1
			AND created_at >= $2 AND created_at <= $3
			AND acknowledged_at IS NOT NULL
	`
	var avgResponse sql.NullFloat64
	if err := r.db.QueryRowContext(ctx, avgResponseQuery, orgID, startTime, endTime).Scan(&avgResponse); err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if avgResponse.Valid {
		metrics.AvgResponseTime = &avgResponse.Float64
	}

	// Get average resolution time (time to close)
	avgResolutionQuery := `
		SELECT AVG(EXTRACT(EPOCH FROM (closed_at - created_at)))
		FROM alerts
		WHERE organization_id = $1
			AND created_at >= $2 AND created_at <= $3
			AND closed_at IS NOT NULL
	`
	var avgResolution sql.NullFloat64
	if err := r.db.QueryRowContext(ctx, avgResolutionQuery, orgID, startTime, endTime).Scan(&avgResolution); err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if avgResolution.Valid {
		metrics.AvgResolutionTime = &avgResolution.Float64
	}

	return metrics, nil
}

func (r *metricsRepository) GetIncidentMetrics(ctx context.Context, orgID uuid.UUID, filter *domain.MetricsFilter) (*domain.IncidentMetrics, error) {
	metrics := &domain.IncidentMetrics{
		BySeverity: make(map[string]int64),
	}

	// Build time filter
	startTime := time.Now().AddDate(0, 0, -30)
	endTime := time.Now()
	if filter != nil {
		if filter.StartTime != nil {
			startTime = *filter.StartTime
		}
		if filter.EndTime != nil {
			endTime = *filter.EndTime
		}
	}

	// Get total and status counts
	query := `
		SELECT
			COUNT(*) as total,
			COUNT(*) FILTER (WHERE status = 'open') as open,
			COUNT(*) FILTER (WHERE status = 'investigating') as investigating,
			COUNT(*) FILTER (WHERE status = 'identified') as identified,
			COUNT(*) FILTER (WHERE status = 'monitoring') as monitoring,
			COUNT(*) FILTER (WHERE status = 'resolved') as resolved,
			COUNT(*) FILTER (WHERE status = 'closed') as closed
		FROM incidents
		WHERE organization_id = $1 AND created_at >= $2 AND created_at <= $3
	`

	err := r.db.QueryRowContext(ctx, query, orgID, startTime, endTime).Scan(
		&metrics.Total,
		&metrics.Open,
		&metrics.Investigating,
		&metrics.Identified,
		&metrics.Monitoring,
		&metrics.Resolved,
		&metrics.Closed,
	)
	if err != nil {
		return nil, err
	}

	// Get counts by severity
	severityQuery := `
		SELECT severity, COUNT(*) as count
		FROM incidents
		WHERE organization_id = $1 AND created_at >= $2 AND created_at <= $3
		GROUP BY severity
	`
	rows, err := r.db.QueryContext(ctx, severityQuery, orgID, startTime, endTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var severity string
		var count int64
		if err := rows.Scan(&severity, &count); err != nil {
			return nil, err
		}
		metrics.BySeverity[severity] = count
	}

	// Get average resolution time
	avgResolutionQuery := `
		SELECT AVG(EXTRACT(EPOCH FROM (resolved_at - created_at)))
		FROM incidents
		WHERE organization_id = $1
			AND created_at >= $2 AND created_at <= $3
			AND resolved_at IS NOT NULL
	`
	var avgResolution sql.NullFloat64
	if err := r.db.QueryRowContext(ctx, avgResolutionQuery, orgID, startTime, endTime).Scan(&avgResolution); err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if avgResolution.Valid {
		metrics.AvgResolutionTime = &avgResolution.Float64
	}

	return metrics, nil
}

func (r *metricsRepository) GetNotificationMetrics(ctx context.Context, orgID uuid.UUID, filter *domain.MetricsFilter) (*domain.NotificationMetrics, error) {
	metrics := &domain.NotificationMetrics{
		ByChannel: make(map[string]int64),
	}

	// Build time filter
	startTime := time.Now().AddDate(0, 0, -30)
	endTime := time.Now()
	if filter != nil {
		if filter.StartTime != nil {
			startTime = *filter.StartTime
		}
		if filter.EndTime != nil {
			endTime = *filter.EndTime
		}
	}

	// Get total and status counts
	query := `
		SELECT
			COUNT(*) as total,
			COUNT(*) FILTER (WHERE status = 'sent') as sent,
			COUNT(*) FILTER (WHERE status = 'pending') as pending,
			COUNT(*) FILTER (WHERE status = 'failed') as failed
		FROM notification_logs
		WHERE organization_id = $1 AND created_at >= $2 AND created_at <= $3
	`

	err := r.db.QueryRowContext(ctx, query, orgID, startTime, endTime).Scan(
		&metrics.Total,
		&metrics.Sent,
		&metrics.Pending,
		&metrics.Failed,
	)
	if err != nil {
		return nil, err
	}

	// Get counts by channel type
	channelQuery := `
		SELECT nc.channel_type, COUNT(*) as count
		FROM notification_logs nl
		JOIN notification_channels nc ON nl.channel_id = nc.id
		WHERE nl.organization_id = $1 AND nl.created_at >= $2 AND nl.created_at <= $3
		GROUP BY nc.channel_type
	`
	rows, err := r.db.QueryContext(ctx, channelQuery, orgID, startTime, endTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var channelType string
		var count int64
		if err := rows.Scan(&channelType, &count); err != nil {
			return nil, err
		}
		metrics.ByChannel[channelType] = count
	}

	return metrics, nil
}

func (r *metricsRepository) GetAlertTrend(ctx context.Context, orgID uuid.UUID, filter *domain.MetricsFilter) (*domain.AlertTrend, error) {
	period := "daily"
	if filter != nil && filter.Period != "" {
		period = filter.Period
	}

	// Build time filter
	startTime := time.Now().AddDate(0, 0, -30)
	endTime := time.Now()
	if filter != nil {
		if filter.StartTime != nil {
			startTime = *filter.StartTime
		}
		if filter.EndTime != nil {
			endTime = *filter.EndTime
		}
	}

	trend := &domain.AlertTrend{
		Period:  period,
		Created: []domain.TimeSeriesPoint{},
		Closed:  []domain.TimeSeriesPoint{},
	}

	// Determine the date_trunc interval
	truncInterval := "day"
	switch period {
	case "hourly":
		truncInterval = "hour"
	case "weekly":
		truncInterval = "week"
	}

	// Get created alerts trend
	createdQuery := `
		SELECT date_trunc($1, created_at) as bucket, COUNT(*) as count
		FROM alerts
		WHERE organization_id = $2 AND created_at >= $3 AND created_at <= $4
		GROUP BY bucket
		ORDER BY bucket
	`
	rows, err := r.db.QueryContext(ctx, createdQuery, truncInterval, orgID, startTime, endTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var point domain.TimeSeriesPoint
		if err := rows.Scan(&point.Timestamp, &point.Value); err != nil {
			return nil, err
		}
		trend.Created = append(trend.Created, point)
	}

	// Get closed alerts trend
	closedQuery := `
		SELECT date_trunc($1, closed_at) as bucket, COUNT(*) as count
		FROM alerts
		WHERE organization_id = $2 AND closed_at >= $3 AND closed_at <= $4 AND closed_at IS NOT NULL
		GROUP BY bucket
		ORDER BY bucket
	`
	rows, err = r.db.QueryContext(ctx, closedQuery, truncInterval, orgID, startTime, endTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var point domain.TimeSeriesPoint
		if err := rows.Scan(&point.Timestamp, &point.Value); err != nil {
			return nil, err
		}
		trend.Closed = append(trend.Closed, point)
	}

	return trend, nil
}

func (r *metricsRepository) GetTeamMetrics(ctx context.Context, orgID uuid.UUID, filter *domain.MetricsFilter) ([]domain.TeamMetrics, error) {
	// Build time filter
	startTime := time.Now().AddDate(0, 0, -30)
	endTime := time.Now()
	if filter != nil {
		if filter.StartTime != nil {
			startTime = *filter.StartTime
		}
		if filter.EndTime != nil {
			endTime = *filter.EndTime
		}
	}

	query := `
		SELECT
			t.id,
			t.name,
			COUNT(a.id) as total_alerts,
			COUNT(a.id) FILTER (WHERE a.acknowledged_at IS NOT NULL) as acknowledged_alerts,
			COUNT(a.id) FILTER (WHERE a.closed_at IS NOT NULL) as closed_alerts,
			AVG(EXTRACT(EPOCH FROM (a.acknowledged_at - a.created_at))) FILTER (WHERE a.acknowledged_at IS NOT NULL) as avg_response_time
		FROM teams t
		LEFT JOIN alerts a ON a.assigned_to_team_id = t.id
			AND a.created_at >= $2 AND a.created_at <= $3
		WHERE t.organization_id = $1
		GROUP BY t.id, t.name
		ORDER BY total_alerts DESC
	`

	rows, err := r.db.QueryContext(ctx, query, orgID, startTime, endTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metrics []domain.TeamMetrics
	for rows.Next() {
		var m domain.TeamMetrics
		var avgResponseTime sql.NullFloat64
		if err := rows.Scan(
			&m.TeamID,
			&m.TeamName,
			&m.TotalAlerts,
			&m.AcknowledgedAlerts,
			&m.ClosedAlerts,
			&avgResponseTime,
		); err != nil {
			return nil, err
		}
		if avgResponseTime.Valid {
			m.AvgResponseTime = &avgResponseTime.Float64
		}
		metrics = append(metrics, m)
	}

	return metrics, nil
}
