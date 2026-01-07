package domain

import "time"

// AlertMetrics contains aggregated alert statistics
type AlertMetrics struct {
	Total          int64 `json:"total"`
	Open           int64 `json:"open"`
	Acknowledged   int64 `json:"acknowledged"`
	Closed         int64 `json:"closed"`
	Snoozed        int64 `json:"snoozed"`
	ByPriority     map[string]int64 `json:"by_priority"`
	BySource       map[string]int64 `json:"by_source"`
	AvgResponseTime *float64 `json:"avg_response_time_seconds,omitempty"` // Time to acknowledge
	AvgResolutionTime *float64 `json:"avg_resolution_time_seconds,omitempty"` // Time to close
}

// IncidentMetrics contains aggregated incident statistics
type IncidentMetrics struct {
	Total              int64 `json:"total"`
	Open               int64 `json:"open"`
	Investigating      int64 `json:"investigating"`
	Identified         int64 `json:"identified"`
	Monitoring         int64 `json:"monitoring"`
	Resolved           int64 `json:"resolved"`
	Closed             int64 `json:"closed"`
	BySeverity         map[string]int64 `json:"by_severity"`
	AvgResolutionTime  *float64 `json:"avg_resolution_time_seconds,omitempty"`
}

// NotificationMetrics contains notification delivery statistics
type NotificationMetrics struct {
	Total     int64 `json:"total"`
	Sent      int64 `json:"sent"`
	Pending   int64 `json:"pending"`
	Failed    int64 `json:"failed"`
	ByChannel map[string]int64 `json:"by_channel"`
}

// TeamMetrics contains team performance metrics
type TeamMetrics struct {
	TeamID          string `json:"team_id"`
	TeamName        string `json:"team_name"`
	TotalAlerts     int64  `json:"total_alerts"`
	AcknowledgedAlerts int64 `json:"acknowledged_alerts"`
	ClosedAlerts    int64 `json:"closed_alerts"`
	AvgResponseTime *float64 `json:"avg_response_time_seconds,omitempty"`
}

// OnCallCoverage contains schedule coverage information
type OnCallCoverage struct {
	ScheduleID   string `json:"schedule_id"`
	ScheduleName string `json:"schedule_name"`
	CoveragePercent float64 `json:"coverage_percent"`
	GapsCount    int    `json:"gaps_count"`
	TotalHours   float64 `json:"total_hours"`
	CoveredHours float64 `json:"covered_hours"`
}

// TimeSeriesPoint represents a single data point in a time series
type TimeSeriesPoint struct {
	Timestamp time.Time `json:"timestamp"`
	Value     int64     `json:"value"`
}

// AlertTrend contains time-series data for alerts
type AlertTrend struct {
	Period    string            `json:"period"` // hourly, daily, weekly
	Created   []TimeSeriesPoint `json:"created"`
	Closed    []TimeSeriesPoint `json:"closed"`
}

// DashboardMetrics contains all metrics for the dashboard
type DashboardMetrics struct {
	Alerts        *AlertMetrics        `json:"alerts"`
	Incidents     *IncidentMetrics     `json:"incidents"`
	Notifications *NotificationMetrics `json:"notifications"`
	AlertTrend    *AlertTrend          `json:"alert_trend,omitempty"`
	UpdatedAt     time.Time            `json:"updated_at"`
}

// MetricsFilter contains filter options for metrics queries
type MetricsFilter struct {
	StartTime *time.Time `json:"start_time,omitempty"`
	EndTime   *time.Time `json:"end_time,omitempty"`
	TeamID    *string    `json:"team_id,omitempty"`
	Period    string     `json:"period,omitempty"` // hourly, daily, weekly (for trends)
}
