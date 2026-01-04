import { browser } from '$app/environment';
import type { AuthResponse, LoginRequest, RegisterRequest, User } from '$lib/types/user';
import type {
	Alert,
	AssignAlertRequest,
	CloseAlertRequest,
	CreateAlertRequest,
	ListAlertsParams,
	ListAlertsResponse,
	SnoozeAlertRequest,
	UpdateAlertRequest
} from '$lib/types/alert';
import type {
	Team,
	TeamWithMembers,
	CreateTeamRequest,
	UpdateTeamRequest,
	AddTeamMemberRequest,
	UpdateTeamMemberRoleRequest
} from '$lib/types/team';
import type {
	Schedule,
	ScheduleWithRotations,
	ScheduleRotation,
	ScheduleOverride,
	OnCallUser,
	ParticipantWithUser,
	CreateScheduleRequest,
	UpdateScheduleRequest,
	CreateRotationRequest,
	UpdateRotationRequest,
	AddParticipantRequest,
	ReorderParticipantsRequest,
	CreateOverrideRequest,
	UpdateOverrideRequest,
	ListSchedulesResponse,
	ListRotationsResponse,
	ListParticipantsResponse,
	ListOverridesResponse
} from '$lib/types/schedule';
import type {
	EscalationPolicy,
	EscalationPolicyWithRules,
	EscalationRule,
	EscalationTarget,
	CreateEscalationPolicyRequest,
	UpdateEscalationPolicyRequest,
	CreateEscalationRuleRequest,
	UpdateEscalationRuleRequest,
	AddEscalationTargetRequest,
	ListEscalationPoliciesResponse,
	ListEscalationRulesResponse,
	ListEscalationTargetsResponse
} from '$lib/types/escalation';
import type {
	NotificationChannel,
	UserNotificationPreference,
	NotificationLog,
	CreateNotificationChannelRequest,
	UpdateNotificationChannelRequest,
	CreateUserNotificationPreferenceRequest,
	UpdateUserNotificationPreferenceRequest,
	SendNotificationRequest,
	ListNotificationChannelsResponse,
	ListUserNotificationPreferencesResponse,
	ListNotificationLogsResponse
} from '$lib/types/notification';

const API_URL = browser ? import.meta.env.VITE_API_URL || 'http://localhost:8080' : 'http://backend:8080';

class APIClient {
	private baseURL: string;
	private accessToken: string | null = null;

	constructor(baseURL: string) {
		this.baseURL = baseURL;

		// Load token from localStorage if in browser
		if (browser) {
			this.accessToken = localStorage.getItem('access_token');
		}
	}

	setAccessToken(token: string | null) {
		this.accessToken = token;
		if (browser) {
			if (token) {
				localStorage.setItem('access_token', token);
			} else {
				localStorage.removeItem('access_token');
			}
		}
	}

	getAccessToken(): string | null {
		return this.accessToken;
	}

	private async request<T>(
		endpoint: string,
		options: RequestInit = {}
	): Promise<T> {
		const headers: HeadersInit = {
			'Content-Type': 'application/json',
			...options.headers
		};

		if (this.accessToken) {
			headers['Authorization'] = `Bearer ${this.accessToken}`;
		}

		const response = await fetch(`${this.baseURL}${endpoint}`, {
			...options,
			headers
		});

		if (!response.ok) {
			const error = await response.json().catch(() => ({ error: 'An error occurred' }));
			throw new Error(error.error || response.statusText);
		}

		return response.json();
	}

	// Auth endpoints
	async register(data: RegisterRequest): Promise<AuthResponse> {
		const response = await this.request<AuthResponse>('/api/v1/auth/register', {
			method: 'POST',
			body: JSON.stringify(data)
		});

		this.setAccessToken(response.access_token);
		if (browser) {
			localStorage.setItem('refresh_token', response.refresh_token);
		}

		return response;
	}

	async login(data: LoginRequest): Promise<AuthResponse> {
		const response = await this.request<AuthResponse>('/api/v1/auth/login', {
			method: 'POST',
			body: JSON.stringify(data)
		});

		this.setAccessToken(response.access_token);
		if (browser) {
			localStorage.setItem('refresh_token', response.refresh_token);
		}

		return response;
	}

	async logout(): Promise<void> {
		try {
			await this.request('/api/v1/auth/logout', {
				method: 'POST'
			});
		} finally {
			this.setAccessToken(null);
			if (browser) {
				localStorage.removeItem('refresh_token');
			}
		}
	}

	async refreshToken(): Promise<AuthResponse> {
		const refreshToken = browser ? localStorage.getItem('refresh_token') : null;

		if (!refreshToken) {
			throw new Error('No refresh token available');
		}

		const response = await this.request<AuthResponse>('/api/v1/auth/refresh', {
			method: 'POST',
			body: JSON.stringify({ refresh_token: refreshToken })
		});

		this.setAccessToken(response.access_token);
		if (browser) {
			localStorage.setItem('refresh_token', response.refresh_token);
		}

		return response;
	}

	async getMe(): Promise<User> {
		return this.request<User>('/api/v1/auth/me');
	}

	// User endpoints
	async listUsers(): Promise<{ users: User[] }> {
		return this.request<{ users: User[] }>('/api/v1/users');
	}

	// Alert endpoints
	async listAlerts(params?: ListAlertsParams): Promise<ListAlertsResponse> {
		const queryParams = new URLSearchParams();

		if (params?.status) {
			params.status.forEach(s => queryParams.append('status', s));
		}
		if (params?.priority) {
			params.priority.forEach(p => queryParams.append('priority', p));
		}
		if (params?.assigned_to_user) {
			queryParams.append('assigned_to_user', params.assigned_to_user);
		}
		if (params?.assigned_to_team) {
			queryParams.append('assigned_to_team', params.assigned_to_team);
		}
		if (params?.source) {
			queryParams.append('source', params.source);
		}
		if (params?.search) {
			queryParams.append('search', params.search);
		}
		if (params?.page) {
			queryParams.append('page', params.page.toString());
		}
		if (params?.page_size) {
			queryParams.append('page_size', params.page_size.toString());
		}

		const query = queryParams.toString();
		const url = query ? `/api/v1/alerts?${query}` : '/api/v1/alerts';

		return this.request<ListAlertsResponse>(url);
	}

	async createAlert(data: CreateAlertRequest): Promise<Alert> {
		return this.request<Alert>('/api/v1/alerts', {
			method: 'POST',
			body: JSON.stringify(data)
		});
	}

	async getAlert(id: string): Promise<Alert> {
		return this.request<Alert>(`/api/v1/alerts/${id}`);
	}

	async updateAlert(id: string, data: UpdateAlertRequest): Promise<Alert> {
		return this.request<Alert>(`/api/v1/alerts/${id}`, {
			method: 'PATCH',
			body: JSON.stringify(data)
		});
	}

	async deleteAlert(id: string): Promise<void> {
		await this.request(`/api/v1/alerts/${id}`, {
			method: 'DELETE'
		});
	}

	async acknowledgeAlert(id: string): Promise<void> {
		await this.request(`/api/v1/alerts/${id}/acknowledge`, {
			method: 'POST'
		});
	}

	async closeAlert(id: string, data: CloseAlertRequest): Promise<void> {
		await this.request(`/api/v1/alerts/${id}/close`, {
			method: 'POST',
			body: JSON.stringify(data)
		});
	}

	async snoozeAlert(id: string, data: SnoozeAlertRequest): Promise<void> {
		await this.request(`/api/v1/alerts/${id}/snooze`, {
			method: 'POST',
			body: JSON.stringify(data)
		});
	}

	async assignAlert(id: string, data: AssignAlertRequest): Promise<void> {
		await this.request(`/api/v1/alerts/${id}/assign`, {
			method: 'POST',
			body: JSON.stringify(data)
		});
	}

	// Team endpoints
	async listTeams(page = 1, pageSize = 20): Promise<{ teams: Team[] }> {
		return this.request<{ teams: Team[] }>(
			`/api/v1/teams?page=${page}&page_size=${pageSize}`
		);
	}

	async createTeam(data: CreateTeamRequest): Promise<Team> {
		return this.request<Team>('/api/v1/teams', {
			method: 'POST',
			body: JSON.stringify(data)
		});
	}

	async getTeam(id: string): Promise<TeamWithMembers> {
		return this.request<TeamWithMembers>(`/api/v1/teams/${id}`);
	}

	async updateTeam(id: string, data: UpdateTeamRequest): Promise<Team> {
		return this.request<Team>(`/api/v1/teams/${id}`, {
			method: 'PATCH',
			body: JSON.stringify(data)
		});
	}

	async deleteTeam(id: string): Promise<void> {
		await this.request(`/api/v1/teams/${id}`, {
			method: 'DELETE'
		});
	}

	async addTeamMember(teamId: string, data: AddTeamMemberRequest): Promise<void> {
		await this.request(`/api/v1/teams/${teamId}/members`, {
			method: 'POST',
			body: JSON.stringify(data)
		});
	}

	async removeTeamMember(teamId: string, userId: string): Promise<void> {
		await this.request(`/api/v1/teams/${teamId}/members/${userId}`, {
			method: 'DELETE'
		});
	}

	async updateTeamMemberRole(
		teamId: string,
		userId: string,
		data: UpdateTeamMemberRoleRequest
	): Promise<void> {
		await this.request(`/api/v1/teams/${teamId}/members/${userId}`, {
			method: 'PATCH',
			body: JSON.stringify(data)
		});
	}

	async listTeamMembers(teamId: string): Promise<{ members: User[] }> {
		return this.request<{ members: User[] }>(`/api/v1/teams/${teamId}/members`);
	}

	// Schedule endpoints
	async listSchedules(page = 1, pageSize = 20): Promise<ListSchedulesResponse> {
		return this.request<ListSchedulesResponse>(
			`/api/v1/schedules?page=${page}&page_size=${pageSize}`
		);
	}

	async createSchedule(data: CreateScheduleRequest): Promise<Schedule> {
		return this.request<Schedule>('/api/v1/schedules', {
			method: 'POST',
			body: JSON.stringify(data)
		});
	}

	async getSchedule(id: string): Promise<ScheduleWithRotations> {
		return this.request<ScheduleWithRotations>(`/api/v1/schedules/${id}`);
	}

	async updateSchedule(id: string, data: UpdateScheduleRequest): Promise<Schedule> {
		return this.request<Schedule>(`/api/v1/schedules/${id}`, {
			method: 'PATCH',
			body: JSON.stringify(data)
		});
	}

	async deleteSchedule(id: string): Promise<void> {
		await this.request(`/api/v1/schedules/${id}`, {
			method: 'DELETE'
		});
	}

	async getOnCallUser(scheduleId: string, at?: string): Promise<OnCallUser> {
		const params = at ? `?at=${encodeURIComponent(at)}` : '';
		return this.request<OnCallUser>(`/api/v1/schedules/${scheduleId}/oncall${params}`);
	}

	// Rotation endpoints
	async listRotations(scheduleId: string): Promise<ListRotationsResponse> {
		return this.request<ListRotationsResponse>(`/api/v1/schedules/${scheduleId}/rotations`);
	}

	async createRotation(scheduleId: string, data: CreateRotationRequest): Promise<ScheduleRotation> {
		return this.request<ScheduleRotation>(`/api/v1/schedules/${scheduleId}/rotations`, {
			method: 'POST',
			body: JSON.stringify(data)
		});
	}

	async getRotation(scheduleId: string, rotationId: string): Promise<ScheduleRotation> {
		return this.request<ScheduleRotation>(
			`/api/v1/schedules/${scheduleId}/rotations/${rotationId}`
		);
	}

	async updateRotation(
		scheduleId: string,
		rotationId: string,
		data: UpdateRotationRequest
	): Promise<ScheduleRotation> {
		return this.request<ScheduleRotation>(
			`/api/v1/schedules/${scheduleId}/rotations/${rotationId}`,
			{
				method: 'PATCH',
				body: JSON.stringify(data)
			}
		);
	}

	async deleteRotation(scheduleId: string, rotationId: string): Promise<void> {
		await this.request(`/api/v1/schedules/${scheduleId}/rotations/${rotationId}`, {
			method: 'DELETE'
		});
	}

	// Participant endpoints
	async listParticipants(scheduleId: string, rotationId: string): Promise<ListParticipantsResponse> {
		return this.request<ListParticipantsResponse>(
			`/api/v1/schedules/${scheduleId}/rotations/${rotationId}/participants`
		);
	}

	async addParticipant(
		scheduleId: string,
		rotationId: string,
		data: AddParticipantRequest
	): Promise<void> {
		await this.request(`/api/v1/schedules/${scheduleId}/rotations/${rotationId}/participants`, {
			method: 'POST',
			body: JSON.stringify(data)
		});
	}

	async removeParticipant(
		scheduleId: string,
		rotationId: string,
		userId: string
	): Promise<void> {
		await this.request(
			`/api/v1/schedules/${scheduleId}/rotations/${rotationId}/participants/${userId}`,
			{
				method: 'DELETE'
			}
		);
	}

	async reorderParticipants(
		scheduleId: string,
		rotationId: string,
		data: ReorderParticipantsRequest
	): Promise<void> {
		await this.request(
			`/api/v1/schedules/${scheduleId}/rotations/${rotationId}/participants/reorder`,
			{
				method: 'PUT',
				body: JSON.stringify(data)
			}
		);
	}

	// Override endpoints
	async listOverrides(scheduleId: string, start?: string, end?: string): Promise<ListOverridesResponse> {
		const params = new URLSearchParams();
		if (start) params.append('start', start);
		if (end) params.append('end', end);
		const queryString = params.toString();
		return this.request<ListOverridesResponse>(
			`/api/v1/schedules/${scheduleId}/overrides${queryString ? '?' + queryString : ''}`
		);
	}

	async createOverride(scheduleId: string, data: CreateOverrideRequest): Promise<ScheduleOverride> {
		return this.request<ScheduleOverride>(`/api/v1/schedules/${scheduleId}/overrides`, {
			method: 'POST',
			body: JSON.stringify(data)
		});
	}

	async getOverride(scheduleId: string, overrideId: string): Promise<ScheduleOverride> {
		return this.request<ScheduleOverride>(
			`/api/v1/schedules/${scheduleId}/overrides/${overrideId}`
		);
	}

	async updateOverride(
		scheduleId: string,
		overrideId: string,
		data: UpdateOverrideRequest
	): Promise<ScheduleOverride> {
		return this.request<ScheduleOverride>(
			`/api/v1/schedules/${scheduleId}/overrides/${overrideId}`,
			{
				method: 'PATCH',
				body: JSON.stringify(data)
			}
		);
	}

	async deleteOverride(scheduleId: string, overrideId: string): Promise<void> {
		await this.request(`/api/v1/schedules/${scheduleId}/overrides/${overrideId}`, {
			method: 'DELETE'
		});
	}

	// Escalation policy endpoints
	async listEscalationPolicies(page = 1, pageSize = 20): Promise<ListEscalationPoliciesResponse> {
		return this.request<ListEscalationPoliciesResponse>(
			`/api/v1/escalation-policies?page=${page}&page_size=${pageSize}`
		);
	}

	async createEscalationPolicy(data: CreateEscalationPolicyRequest): Promise<EscalationPolicy> {
		return this.request<EscalationPolicy>('/api/v1/escalation-policies', {
			method: 'POST',
			body: JSON.stringify(data)
		});
	}

	async getEscalationPolicy(id: string): Promise<EscalationPolicyWithRules> {
		return this.request<EscalationPolicyWithRules>(`/api/v1/escalation-policies/${id}`);
	}

	async updateEscalationPolicy(
		id: string,
		data: UpdateEscalationPolicyRequest
	): Promise<EscalationPolicy> {
		return this.request<EscalationPolicy>(`/api/v1/escalation-policies/${id}`, {
			method: 'PATCH',
			body: JSON.stringify(data)
		});
	}

	async deleteEscalationPolicy(id: string): Promise<void> {
		await this.request(`/api/v1/escalation-policies/${id}`, {
			method: 'DELETE'
		});
	}

	// Escalation rule endpoints
	async listEscalationRules(policyId: string): Promise<ListEscalationRulesResponse> {
		return this.request<ListEscalationRulesResponse>(
			`/api/v1/escalation-policies/${policyId}/rules`
		);
	}

	async createEscalationRule(
		policyId: string,
		data: CreateEscalationRuleRequest
	): Promise<EscalationRule> {
		return this.request<EscalationRule>(`/api/v1/escalation-policies/${policyId}/rules`, {
			method: 'POST',
			body: JSON.stringify(data)
		});
	}

	async getEscalationRule(policyId: string, ruleId: string): Promise<EscalationRule> {
		return this.request<EscalationRule>(
			`/api/v1/escalation-policies/${policyId}/rules/${ruleId}`
		);
	}

	async updateEscalationRule(
		policyId: string,
		ruleId: string,
		data: UpdateEscalationRuleRequest
	): Promise<EscalationRule> {
		return this.request<EscalationRule>(
			`/api/v1/escalation-policies/${policyId}/rules/${ruleId}`,
			{
				method: 'PATCH',
				body: JSON.stringify(data)
			}
		);
	}

	async deleteEscalationRule(policyId: string, ruleId: string): Promise<void> {
		await this.request(`/api/v1/escalation-policies/${policyId}/rules/${ruleId}`, {
			method: 'DELETE'
		});
	}

	// Escalation target endpoints
	async listEscalationTargets(
		policyId: string,
		ruleId: string
	): Promise<ListEscalationTargetsResponse> {
		return this.request<ListEscalationTargetsResponse>(
			`/api/v1/escalation-policies/${policyId}/rules/${ruleId}/targets`
		);
	}

	async addEscalationTarget(
		policyId: string,
		ruleId: string,
		data: AddEscalationTargetRequest
	): Promise<EscalationTarget> {
		return this.request<EscalationTarget>(
			`/api/v1/escalation-policies/${policyId}/rules/${ruleId}/targets`,
			{
				method: 'POST',
				body: JSON.stringify(data)
			}
		);
	}

	async removeEscalationTarget(policyId: string, ruleId: string, targetId: string): Promise<void> {
		await this.request(
			`/api/v1/escalation-policies/${policyId}/rules/${ruleId}/targets/${targetId}`,
			{
				method: 'DELETE'
			}
		);
	}

	// ==================== Notification Channels ====================

	async listNotificationChannels(): Promise<ListNotificationChannelsResponse> {
		return this.request<ListNotificationChannelsResponse>('/api/v1/notifications/channels');
	}

	async createNotificationChannel(
		data: CreateNotificationChannelRequest
	): Promise<NotificationChannel> {
		return this.request<NotificationChannel>('/api/v1/notifications/channels', {
			method: 'POST',
			body: JSON.stringify(data)
		});
	}

	async getNotificationChannel(id: string): Promise<NotificationChannel> {
		return this.request<NotificationChannel>(`/api/v1/notifications/channels/${id}`);
	}

	async updateNotificationChannel(
		id: string,
		data: UpdateNotificationChannelRequest
	): Promise<NotificationChannel> {
		return this.request<NotificationChannel>(`/api/v1/notifications/channels/${id}`, {
			method: 'PATCH',
			body: JSON.stringify(data)
		});
	}

	async deleteNotificationChannel(id: string): Promise<void> {
		await this.request(`/api/v1/notifications/channels/${id}`, {
			method: 'DELETE'
		});
	}

	// ==================== User Notification Preferences ====================

	async listUserNotificationPreferences(): Promise<ListUserNotificationPreferencesResponse> {
		return this.request<ListUserNotificationPreferencesResponse>(
			'/api/v1/notifications/preferences'
		);
	}

	async createUserNotificationPreference(
		data: CreateUserNotificationPreferenceRequest
	): Promise<UserNotificationPreference> {
		return this.request<UserNotificationPreference>('/api/v1/notifications/preferences', {
			method: 'POST',
			body: JSON.stringify(data)
		});
	}

	async getUserNotificationPreference(id: string): Promise<UserNotificationPreference> {
		return this.request<UserNotificationPreference>(`/api/v1/notifications/preferences/${id}`);
	}

	async updateUserNotificationPreference(
		id: string,
		data: UpdateUserNotificationPreferenceRequest
	): Promise<UserNotificationPreference> {
		return this.request<UserNotificationPreference>(
			`/api/v1/notifications/preferences/${id}`,
			{
				method: 'PATCH',
				body: JSON.stringify(data)
			}
		);
	}

	async deleteUserNotificationPreference(id: string): Promise<void> {
		await this.request(`/api/v1/notifications/preferences/${id}`, {
			method: 'DELETE'
		});
	}

	// ==================== Sending Notifications ====================

	async sendNotification(data: SendNotificationRequest): Promise<NotificationLog> {
		return this.request<NotificationLog>('/api/v1/notifications/send', {
			method: 'POST',
			body: JSON.stringify(data)
		});
	}

	// ==================== Notification Logs ====================

	async listNotificationLogs(limit?: number, offset?: number): Promise<ListNotificationLogsResponse> {
		const params = new URLSearchParams();
		if (limit) params.append('limit', limit.toString());
		if (offset) params.append('offset', offset.toString());

		const queryString = params.toString();
		const endpoint = queryString ? `/api/v1/notifications/logs?${queryString}` : '/api/v1/notifications/logs';

		return this.request<ListNotificationLogsResponse>(endpoint);
	}

	async getNotificationLog(id: string): Promise<NotificationLog> {
		return this.request<NotificationLog>(`/api/v1/notifications/logs/${id}`);
	}

	async listNotificationLogsByUser(
		limit?: number,
		offset?: number
	): Promise<ListNotificationLogsResponse> {
		const params = new URLSearchParams();
		if (limit) params.append('limit', limit.toString());
		if (offset) params.append('offset', offset.toString());

		const queryString = params.toString();
		const endpoint = queryString
			? `/api/v1/notifications/logs/user/me?${queryString}`
			: '/api/v1/notifications/logs/user/me';

		return this.request<ListNotificationLogsResponse>(endpoint);
	}

	async listNotificationLogsByAlert(alertId: string): Promise<ListNotificationLogsResponse> {
		return this.request<ListNotificationLogsResponse>(
			`/api/v1/notifications/logs/alert/${alertId}`
		);
	}
}

export const api = new APIClient(API_URL);
