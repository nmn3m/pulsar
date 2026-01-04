import { writable } from 'svelte/store';
import { api } from '$lib/api/client';
import type {
	Alert,
	AssignAlertRequest,
	CloseAlertRequest,
	CreateAlertRequest,
	ListAlertsParams,
	SnoozeAlertRequest,
	UpdateAlertRequest
} from '$lib/types/alert';

interface AlertsState {
	alerts: Alert[];
	total: number;
	page: number;
	pageSize: number;
	isLoading: boolean;
	error: string | null;
}

function createAlertsStore() {
	const { subscribe, set, update } = writable<AlertsState>({
		alerts: [],
		total: 0,
		page: 1,
		pageSize: 20,
		isLoading: false,
		error: null
	});

	return {
		subscribe,

		async load(params?: ListAlertsParams) {
			update(state => ({ ...state, isLoading: true, error: null }));

			try {
				const response = await api.listAlerts(params);
				set({
					alerts: response.alerts,
					total: response.total,
					page: response.page,
					pageSize: response.page_size,
					isLoading: false,
					error: null
				});
			} catch (error) {
				update(state => ({
					...state,
					isLoading: false,
					error: error instanceof Error ? error.message : 'Failed to load alerts'
				}));
			}
		},

		async create(data: CreateAlertRequest) {
			try {
				const alert = await api.createAlert(data);
				update(state => ({
					...state,
					alerts: [alert, ...state.alerts],
					total: state.total + 1
				}));
				return alert;
			} catch (error) {
				throw error;
			}
		},

		async update(id: string, data: UpdateAlertRequest) {
			try {
				const updatedAlert = await api.updateAlert(id, data);
				update(state => ({
					...state,
					alerts: state.alerts.map(a => (a.id === id ? updatedAlert : a))
				}));
				return updatedAlert;
			} catch (error) {
				throw error;
			}
		},

		async delete(id: string) {
			try {
				await api.deleteAlert(id);
				update(state => ({
					...state,
					alerts: state.alerts.filter(a => a.id !== id),
					total: state.total - 1
				}));
			} catch (error) {
				throw error;
			}
		},

		async acknowledge(id: string) {
			try {
				await api.acknowledgeAlert(id);
				// Reload to get updated alert
				const alert = await api.getAlert(id);
				update(state => ({
					...state,
					alerts: state.alerts.map(a => (a.id === id ? alert : a))
				}));
			} catch (error) {
				throw error;
			}
		},

		async close(id: string, data: CloseAlertRequest) {
			try {
				await api.closeAlert(id, data);
				// Reload to get updated alert
				const alert = await api.getAlert(id);
				update(state => ({
					...state,
					alerts: state.alerts.map(a => (a.id === id ? alert : a))
				}));
			} catch (error) {
				throw error;
			}
		},

		async snooze(id: string, data: SnoozeAlertRequest) {
			try {
				await api.snoozeAlert(id, data);
				// Reload to get updated alert
				const alert = await api.getAlert(id);
				update(state => ({
					...state,
					alerts: state.alerts.map(a => (a.id === id ? alert : a))
				}));
			} catch (error) {
				throw error;
			}
		},

		async assign(id: string, data: AssignAlertRequest) {
			try {
				await api.assignAlert(id, data);
				// Reload to get updated alert
				const alert = await api.getAlert(id);
				update(state => ({
					...state,
					alerts: state.alerts.map(a => (a.id === id ? alert : a))
				}));
			} catch (error) {
				throw error;
			}
		}
	};
}

export const alertsStore = createAlertsStore();
