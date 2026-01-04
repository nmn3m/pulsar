export type ChannelType = 'email' | 'slack' | 'teams' | 'webhook';
export type NotificationStatus = 'pending' | 'sent' | 'failed';

export interface NotificationChannel {
	id: string;
	organization_id: string;
	name: string;
	channel_type: ChannelType;
	is_enabled: boolean;
	config: Record<string, unknown>;
	created_at: string;
	updated_at: string;
}

export interface UserNotificationPreference {
	id: string;
	user_id: string;
	channel_id: string;
	is_enabled: boolean;
	dnd_enabled: boolean;
	dnd_start_time?: string;
	dnd_end_time?: string;
	min_priority?: string;
	created_at: string;
	updated_at: string;
}

export interface NotificationLog {
	id: string;
	organization_id: string;
	channel_id: string;
	user_id?: string;
	alert_id?: string;
	recipient: string;
	subject?: string;
	message: string;
	status: NotificationStatus;
	error_message?: string;
	sent_at?: string;
	created_at: string;
}

// Request types
export interface CreateNotificationChannelRequest {
	name: string;
	channel_type: ChannelType;
	is_enabled?: boolean;
	config: Record<string, unknown>;
}

export interface UpdateNotificationChannelRequest {
	name?: string;
	channel_type?: ChannelType;
	is_enabled?: boolean;
	config?: Record<string, unknown>;
}

export interface CreateUserNotificationPreferenceRequest {
	channel_id: string;
	is_enabled?: boolean;
	dnd_enabled?: boolean;
	dnd_start_time?: string;
	dnd_end_time?: string;
	min_priority?: string;
}

export interface UpdateUserNotificationPreferenceRequest {
	is_enabled?: boolean;
	dnd_enabled?: boolean;
	dnd_start_time?: string;
	dnd_end_time?: string;
	min_priority?: string;
}

export interface SendNotificationRequest {
	channel_id: string;
	user_id?: string;
	alert_id?: string;
	recipient: string;
	subject?: string;
	message: string;
}

// Response types
export interface ListNotificationChannelsResponse {
	channels: NotificationChannel[];
	total: number;
}

export interface ListUserNotificationPreferencesResponse {
	preferences: UserNotificationPreference[];
	total: number;
}

export interface ListNotificationLogsResponse {
	logs: NotificationLog[];
	total: number;
	limit: number;
	offset: number;
}

// Provider-specific config types
export interface EmailConfig {
	smtp_host: string;
	smtp_port: number;
	smtp_username: string;
	smtp_password: string;
	from_address: string;
	from_name?: string;
	use_tls?: boolean;
}

export interface SlackConfig {
	webhook_url: string;
	channel?: string;
	username?: string;
	icon_emoji?: string;
}

export interface TeamsConfig {
	webhook_url: string;
	theme_color?: string;
}

export interface WebhookConfig {
	url: string;
	method?: string;
	headers?: Record<string, string>;
	timeout?: number;
}
