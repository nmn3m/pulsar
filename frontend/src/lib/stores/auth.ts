import { writable } from 'svelte/store';
import { browser } from '$app/environment';
import { api } from '$lib/api/client';
import type { User, Organization, LoginRequest, RegisterRequest, VerifyEmailRequest, ResendOTPRequest } from '$lib/types/user';

interface AuthState {
	user: User | null;
	organization: Organization | null;
	isAuthenticated: boolean;
	isLoading: boolean;
	pendingVerificationEmail: string | null;
}

function createAuthStore() {
	const { subscribe, set, update } = writable<AuthState>({
		user: null,
		organization: null,
		isAuthenticated: false,
		isLoading: true,
		pendingVerificationEmail: null
	});

	return {
		subscribe,

		async init() {
			if (!browser) return;

			const token = api.getAccessToken();

			if (!token) {
				update(state => ({ ...state, isLoading: false }));
				return;
			}

			try {
				const user = await api.getMe();
				update(state => ({
					...state,
					user,
					isAuthenticated: true,
					isLoading: false
				}));
			} catch (error) {
				// Token might be expired, try to refresh
				try {
					const response = await api.refreshToken();
					update(state => ({
						...state,
						user: response.user,
						organization: response.organization,
						isAuthenticated: true,
						isLoading: false
					}));
				} catch {
					// Refresh failed, clear auth
					api.setAccessToken(null);
					update(state => ({
						...state,
						user: null,
						organization: null,
						isAuthenticated: false,
						isLoading: false
					}));
				}
			}
		},

		async register(data: RegisterRequest) {
			try {
				const response = await api.register(data);
				if (response.requires_email_verification) {
					update(state => ({
						...state,
						pendingVerificationEmail: data.email,
						isLoading: false
					}));
				} else {
					set({
						user: response.user,
						organization: response.organization,
						isAuthenticated: true,
						isLoading: false,
						pendingVerificationEmail: null
					});
				}
				return response;
			} catch (error) {
				throw error;
			}
		},

		async login(data: LoginRequest) {
			try {
				const response = await api.login(data);
				if (response.requires_email_verification) {
					update(state => ({
						...state,
						pendingVerificationEmail: data.email,
						isLoading: false
					}));
				} else {
					set({
						user: response.user,
						organization: response.organization,
						isAuthenticated: true,
						isLoading: false,
						pendingVerificationEmail: null
					});
				}
				return response;
			} catch (error) {
				throw error;
			}
		},

		async verifyEmail(data: VerifyEmailRequest) {
			try {
				await api.verifyEmail(data);
				update(state => ({
					...state,
					pendingVerificationEmail: null
				}));
			} catch (error) {
				throw error;
			}
		},

		async resendOTP(data: ResendOTPRequest) {
			try {
				await api.resendOTP(data);
			} catch (error) {
				throw error;
			}
		},

		setPendingVerificationEmail(email: string | null) {
			update(state => ({
				...state,
				pendingVerificationEmail: email
			}));
		},

		async logout() {
			try {
				await api.logout();
			} finally {
				set({
					user: null,
					organization: null,
					isAuthenticated: false,
					isLoading: false,
					pendingVerificationEmail: null
				});
			}
		}
	};
}

export const authStore = createAuthStore();
