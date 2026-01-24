import { writable } from 'svelte/store';
import { api } from '$lib/api/client';
import type {
  Alert,
  AssignAlertRequest,
  CloseAlertRequest,
  CreateAlertRequest,
  ListAlertsParams,
  SnoozeAlertRequest,
  UpdateAlertRequest,
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
    error: null,
  });

  return {
    subscribe,

    async load(params?: ListAlertsParams) {
      update((state) => ({ ...state, isLoading: true, error: null }));

      try {
        const response = await api.listAlerts(params);
        set({
          alerts: response.alerts || [],
          total: response.total || 0,
          page: response.page || 1,
          pageSize: response.page_size || 20,
          isLoading: false,
          error: null,
        });
      } catch (error) {
        update((state) => ({
          ...state,
          isLoading: false,
          error: error instanceof Error ? error.message : 'Failed to load alerts',
        }));
      }
    },

    async create(data: CreateAlertRequest) {
      const alert = await api.createAlert(data);
      update((state) => ({
        ...state,
        alerts: [alert, ...state.alerts],
        total: state.total + 1,
      }));
      return alert;
    },

    async update(id: string, data: UpdateAlertRequest) {
      const updatedAlert = await api.updateAlert(id, data);
      update((state) => ({
        ...state,
        alerts: state.alerts.map((a) => (a.id === id ? updatedAlert : a)),
      }));
      return updatedAlert;
    },

    async delete(id: string) {
      await api.deleteAlert(id);
      update((state) => ({
        ...state,
        alerts: state.alerts.filter((a) => a.id !== id),
        total: state.total - 1,
      }));
    },

    async acknowledge(id: string) {
      await api.acknowledgeAlert(id);
      // Reload to get updated alert
      const alert = await api.getAlert(id);
      update((state) => ({
        ...state,
        alerts: state.alerts.map((a) => (a.id === id ? alert : a)),
      }));
    },

    async close(id: string, data: CloseAlertRequest) {
      await api.closeAlert(id, data);
      // Reload to get updated alert
      const alert = await api.getAlert(id);
      update((state) => ({
        ...state,
        alerts: state.alerts.map((a) => (a.id === id ? alert : a)),
      }));
    },

    async snooze(id: string, data: SnoozeAlertRequest) {
      await api.snoozeAlert(id, data);
      // Reload to get updated alert
      const alert = await api.getAlert(id);
      update((state) => ({
        ...state,
        alerts: state.alerts.map((a) => (a.id === id ? alert : a)),
      }));
    },

    async assign(id: string, data: AssignAlertRequest) {
      await api.assignAlert(id, data);
      // Reload to get updated alert
      const alert = await api.getAlert(id);
      update((state) => ({
        ...state,
        alerts: state.alerts.map((a) => (a.id === id ? alert : a)),
      }));
    },
  };
}

export const alertsStore = createAlertsStore();
