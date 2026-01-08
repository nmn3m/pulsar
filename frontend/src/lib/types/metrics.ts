export interface AlertMetrics {
  total: number;
  open: number;
  acknowledged: number;
  closed: number;
  snoozed: number;
  by_priority: Record<string, number>;
  by_source: Record<string, number>;
  avg_response_time_seconds?: number;
  avg_resolution_time_seconds?: number;
}

export interface IncidentMetrics {
  total: number;
  open: number;
  investigating: number;
  identified: number;
  monitoring: number;
  resolved: number;
  closed: number;
  by_severity: Record<string, number>;
  avg_resolution_time_seconds?: number;
}

export interface NotificationMetrics {
  total: number;
  sent: number;
  pending: number;
  failed: number;
  by_channel: Record<string, number>;
}

export interface TimeSeriesPoint {
  timestamp: string;
  value: number;
}

export interface AlertTrend {
  period: string;
  created: TimeSeriesPoint[];
  closed: TimeSeriesPoint[];
}

export interface TeamMetrics {
  team_id: string;
  team_name: string;
  total_alerts: number;
  acknowledged_alerts: number;
  closed_alerts: number;
  avg_response_time_seconds?: number;
}

export interface DashboardMetrics {
  alerts: AlertMetrics;
  incidents: IncidentMetrics;
  notifications: NotificationMetrics;
  alert_trend?: AlertTrend;
  updated_at: string;
}

export interface MetricsFilter {
  start_time?: string;
  end_time?: string;
  period?: 'hourly' | 'daily' | 'weekly';
  team_id?: string;
}
