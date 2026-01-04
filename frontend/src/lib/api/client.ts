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
}

export const api = new APIClient(API_URL);
