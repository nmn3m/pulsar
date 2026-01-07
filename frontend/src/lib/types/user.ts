export interface User {
	id: string;
	email: string;
	username: string;
	full_name?: string;
	phone?: string;
	timezone: string;
	notification_preferences: Record<string, unknown>;
	is_active: boolean;
	email_verified: boolean;
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
	requires_email_verification?: boolean;
}

export interface VerifyEmailRequest {
	email: string;
	otp: string;
}

export interface ResendOTPRequest {
	email: string;
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
