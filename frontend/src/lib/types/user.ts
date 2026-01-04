export interface User {
	id: string;
	email: string;
	username: string;
	full_name?: string;
	phone?: string;
	timezone: string;
	notification_preferences: Record<string, unknown>;
	is_active: boolean;
	created_at: string;
	updated_at: string;
}

export interface Organization {
	id: string;
	name: string;
	slug: string;
	plan: string;
	settings: Record<string, unknown>;
	created_at: string;
	updated_at: string;
}

export interface AuthResponse {
	user: User;
	organization: Organization;
	access_token: string;
	refresh_token: string;
}

export interface RegisterRequest {
	email: string;
	username: string;
	password: string;
	full_name?: string;
	organization_name: string;
}

export interface LoginRequest {
	email: string;
	password: string;
}
