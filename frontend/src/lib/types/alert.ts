export type AlertPriority = 'P1' | 'P2' | 'P3' | 'P4' | 'P5';
export type AlertStatus = 'open' | 'acknowledged' | 'closed' | 'snoozed';
export type AlertSource = 'webhook' | 'api' | 'email' | 'integration' | 'manual';

export interface Alert {
	id: string;
	organization_id: string;
	source: AlertSource;
	source_id?: string;
	priority: AlertPriority;
	status: AlertStatus;
	message: string;
	description?: string;
	tags: string[];
	custom_fields: Record<string, unknown>;

	// Assignment
	assigned_to_user_id?: string;
	assigned_to_team_id?: string;

	// Acknowledgment
	acknowledged_by?: string;
	acknowledged_at?: string;

	// Closure
	closed_by?: string;
	closed_at?: string;
	close_reason?: string;

	// Snooze
	snoozed_until?: string;

	// Escalation
	escalation_policy_id?: string;
	escalation_level: number;
	last_escalated_at?: string;

	created_at: string;
	updated_at: string;
}

export interface CreateAlertRequest {
	source: string;
	source_id?: string;
	priority: AlertPriority;
	message: string;
	description?: string;
	tags?: string[];
	custom_fields?: Record<string, unknown>;
}

export interface UpdateAlertRequest {
	priority?: AlertPriority;
	message?: string;
	description?: string;
	tags?: string[];
	custom_fields?: Record<string, unknown>;
}

export interface CloseAlertRequest {
	reason: string;
}

export interface SnoozeAlertRequest {
	until: string; // ISO 8601 date string
}

export interface AssignAlertRequest {
	user_id?: string;
	team_id?: string;
}

export interface ListAlertsParams {
	status?: AlertStatus[];
	priority?: AlertPriority[];
	assigned_to_user?: string;
	assigned_to_team?: string;
	source?: string;
	search?: string;
	page?: number;
	page_size?: number;
}

export interface ListAlertsResponse {
	alerts: Alert[];
	total: number;
	page: number;
	page_size: number;
}
