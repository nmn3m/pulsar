import type { User } from './user';
import type { Alert } from './alert';

export type IncidentSeverity = 'critical' | 'high' | 'medium' | 'low';
export type IncidentStatus = 'investigating' | 'identified' | 'monitoring' | 'resolved';
export type AlertPriority = 'P1' | 'P2' | 'P3' | 'P4' | 'P5';
export type ResponderRole = 'incident_commander' | 'responder';
export type TimelineEventType =
	| 'created'
	| 'status_changed'
	| 'severity_changed'
	| 'responder_added'
	| 'responder_removed'
	| 'note_added'
	| 'alert_linked'
	| 'alert_unlinked'
	| 'resolved';

export interface Incident {
	id: string;
	organization_id: string;
	title: string;
	description?: string;
	severity: IncidentSeverity;
	status: IncidentStatus;
	priority: AlertPriority;
	created_by_user_id: string;
	assigned_to_team_id?: string;
	started_at: string;
	resolved_at?: string;
	created_at: string;
	updated_at: string;
}

export interface IncidentResponder {
	id: string;
	incident_id: string;
	user_id: string;
	role: ResponderRole;
	added_at: string;
}

export interface ResponderWithUser extends IncidentResponder {
	user: User;
}

export interface IncidentTimelineEvent {
	id: string;
	incident_id: string;
	event_type: TimelineEventType;
	user_id?: string;
	description: string;
	metadata: Record<string, any>;
	created_at: string;
}

export interface TimelineEventWithUser extends IncidentTimelineEvent {
	user?: User;
}

export interface IncidentAlert {
	id: string;
	incident_id: string;
	alert_id: string;
	linked_at: string;
	linked_by_user_id?: string;
}

export interface IncidentAlertWithDetails extends IncidentAlert {
	alert: Alert;
}

export interface IncidentWithDetails extends Incident {
	responders?: ResponderWithUser[];
	alerts?: IncidentAlertWithDetails[];
	timeline?: TimelineEventWithUser[];
}

// Request types
export interface CreateIncidentRequest {
	title: string;
	description?: string;
	severity: IncidentSeverity;
	priority: AlertPriority;
	assigned_to_team_id?: string;
}

export interface UpdateIncidentRequest {
	title?: string;
	description?: string;
	severity?: IncidentSeverity;
	status?: IncidentStatus;
	priority?: AlertPriority;
	assigned_to_team_id?: string;
}

export interface AddResponderRequest {
	user_id: string;
	role: ResponderRole;
}

export interface UpdateResponderRoleRequest {
	role: ResponderRole;
}

export interface AddNoteRequest {
	note: string;
}

export interface LinkAlertRequest {
	alert_id: string;
}

export interface ListIncidentsParams {
	status?: IncidentStatus[];
	severity?: IncidentSeverity[];
	assigned_to_team_id?: string;
	search?: string;
	page?: number;
	page_size?: number;
}

export interface ListIncidentsResponse {
	incidents: Incident[];
	total: number;
	page: number;
	page_size: number;
}
