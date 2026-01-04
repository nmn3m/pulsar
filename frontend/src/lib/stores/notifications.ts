import { writable } from 'svelte/store';
import { api } from '$lib/api/client';
import type {
	NotificationChannel,
	UserNotificationPreference,
	CreateNotificationChannelRequest,
	UpdateNotificationChannelRequest,
	CreateUserNotificationPreferenceRequest,
	UpdateUserNotificationPreferenceRequest
} from '$lib/types/notification';

// Notification Channels Store
interface NotificationChannelsState {
	channels: NotificationChannel[];
	isLoading: boolean;
	error: string | null;
}

function createNotificationChannelsStore() {
	const { subscribe, set, update } = writable<NotificationChannelsState>({
		channels: [],
		isLoading: false,
		error: null
	});

	return {
		subscribe,

		async load() {
			update((state) => ({ ...state, isLoading: true, error: null }));

			try {
				const response = await api.listNotificationChannels();
				update((state) => ({
					...state,
					channels: response.channels,
					isLoading: false
				}));
			} catch (err) {
				update((state) => ({
					...state,
					error: err instanceof Error ? err.message : 'Failed to load notification channels',
					isLoading: false
				}));
			}
		},

		async create(data: CreateNotificationChannelRequest) {
			try {
				const channel = await api.createNotificationChannel(data);
				update((state) => ({
					...state,
					channels: [channel, ...state.channels]
				}));
				return channel;
			} catch (err) {
				throw err;
			}
		},

		async update(id: string, data: UpdateNotificationChannelRequest) {
			try {
				const channel = await api.updateNotificationChannel(id, data);
				update((state) => ({
					...state,
					channels: state.channels.map((c) => (c.id === id ? channel : c))
				}));
				return channel;
			} catch (err) {
				throw err;
			}
		},

		async delete(id: string) {
			try {
				await api.deleteNotificationChannel(id);
				update((state) => ({
					...state,
					channels: state.channels.filter((c) => c.id !== id)
				}));
			} catch (err) {
				throw err;
			}
		}
	};
}

// User Notification Preferences Store
interface UserNotificationPreferencesState {
	preferences: UserNotificationPreference[];
	isLoading: boolean;
	error: string | null;
}

function createUserNotificationPreferencesStore() {
	const { subscribe, set, update } = writable<UserNotificationPreferencesState>({
		preferences: [],
		isLoading: false,
		error: null
	});

	return {
		subscribe,

		async load() {
			update((state) => ({ ...state, isLoading: true, error: null }));

			try {
				const response = await api.listUserNotificationPreferences();
				update((state) => ({
					...state,
					preferences: response.preferences,
					isLoading: false
				}));
			} catch (err) {
				update((state) => ({
					...state,
					error:
						err instanceof Error ? err.message : 'Failed to load notification preferences',
					isLoading: false
				}));
			}
		},

		async create(data: CreateUserNotificationPreferenceRequest) {
			try {
				const preference = await api.createUserNotificationPreference(data);
				update((state) => ({
					...state,
					preferences: [preference, ...state.preferences]
				}));
				return preference;
			} catch (err) {
				throw err;
			}
		},

		async update(id: string, data: UpdateUserNotificationPreferenceRequest) {
			try {
				const preference = await api.updateUserNotificationPreference(id, data);
				update((state) => ({
					...state,
					preferences: state.preferences.map((p) => (p.id === id ? preference : p))
				}));
				return preference;
			} catch (err) {
				throw err;
			}
		},

		async delete(id: string) {
			try {
				await api.deleteUserNotificationPreference(id);
				update((state) => ({
					...state,
					preferences: state.preferences.filter((p) => p.id !== id)
				}));
			} catch (err) {
				throw err;
			}
		}
	};
}

export const notificationChannelsStore = createNotificationChannelsStore();
export const userNotificationPreferencesStore = createUserNotificationPreferencesStore();
