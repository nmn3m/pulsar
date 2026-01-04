import { writable } from 'svelte/store';
import { browser } from '$app/environment';
import { api } from '$lib/api/client';
import type { User, Organization, LoginRequest, RegisterRequest } from '$lib/types/user';

interface AuthState {
	user: User | null;
	organization: Organization | null;
	isAuthenticated: boolean;
	isLoading: boolean;
}

function createAuthStore() {
	const { subscribe, set, update } = writable<AuthState>({
		user: null,
		organization: null,
		isAuthenticated: false,
		isLoading: true
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
				set({
					user: response.user,
					organization: response.organization,
					isAuthenticated: true,
					isLoading: false
				});
				return response;
			} catch (error) {
				throw error;
			}
		},

		async login(data: LoginRequest) {
			try {
				const response = await api.login(data);
				set({
					user: response.user,
					organization: response.organization,
					isAuthenticated: true,
					isLoading: false
				});
				return response;
			} catch (error) {
				throw error;
			}
		},

		async logout() {
			try {
				await api.logout();
			} finally {
				set({
					user: null,
					organization: null,
					isAuthenticated: false,
					isLoading: false
				});
			}
		}
	};
}

export const authStore = createAuthStore();
