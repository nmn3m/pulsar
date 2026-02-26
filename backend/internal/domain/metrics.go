package domain

import "time"

// AlertMetrics contains aggregated alert statistics
type AlertMetrics struct {
	Total             int64
	Open              int64
	Acknowledged      int64
	Closed            int64
	Snoozed           int64
	ByPriority        map[string]int64
	BySource          map[string]int64
	AvgResponseTime   *float64 // Time to acknowledge
	AvgResolutionTime *float64 // Time to close
}

// IncidentMetrics contains aggregated incident statistics
type IncidentMetrics struct {
	Total             int64
	Open              int64
	Investigating     int64
	Identified        int64
	Monitoring        int64
	Resolved          int64
	Closed            int64
	BySeverity        map[string]int64
	AvgResolutionTime *float64
}

// NotificationMetrics contains notification delivery statistics
type NotificationMetrics struct {
	Total     int64
	Sent      int64
	Pending   int64
	Failed    int64
	ByChannel map[string]int64
}

// TeamMetrics contains team performance metrics
type TeamMetrics struct {
	TeamID             string
	TeamName           string
	TotalAlerts        int64
	AcknowledgedAlerts int64
	ClosedAlerts       int64
	AvgResponseTime    *float64
}

// OnCallCoverage contains schedule coverage information
type OnCallCoverage struct {
	ScheduleID      string
	ScheduleName    string
	CoveragePercent float64
	GapsCount       int
	TotalHours      float64
	CoveredHours    float64
}

// TimeSeriesPoint represents a single data point in a time series
type TimeSeriesPoint struct {
	Timestamp time.Time
	Value     int64
}

// AlertTrend contains time-series data for alerts
type AlertTrend struct {
	Period  string // hourly, daily, weekly
	Created []TimeSeriesPoint
	Closed  []TimeSeriesPoint
}

// DashboardMetrics contains all metrics for the dashboard
type DashboardMetrics struct {
	Alerts        *AlertMetrics
	Incidents     *IncidentMetrics
	Notifications *NotificationMetrics
	AlertTrend    *AlertTrend
	UpdatedAt     time.Time
}

// MetricsFilter contains filter options for metrics queries
type MetricsFilter struct {
	StartTime *time.Time
	EndTime   *time.Time
	TeamID    *string
	Period    string // hourly, daily, weekly (for trends)
}
