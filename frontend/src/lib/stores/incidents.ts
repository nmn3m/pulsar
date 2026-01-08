import { writable } from 'svelte/store';
import { api } from '$lib/api/client';
import type {
  Incident,
  IncidentWithDetails,
  CreateIncidentRequest,
  UpdateIncidentRequest,
  ListIncidentsParams,
} from '$lib/types/incident';

interface IncidentsState {
  incidents: Incident[];
  isLoading: boolean;
  error: string | null;
  total: number;
  page: number;
  pageSize: number;
}

function createIncidentsStore() {
  const { subscribe, set, update } = writable<IncidentsState>({
    incidents: [],
    isLoading: false,
    error: null,
    total: 0,
    page: 1,
    pageSize: 20,
  });

  return {
    subscribe,

    async load(params?: ListIncidentsParams) {
      update((state) => ({ ...state, isLoading: true, error: null }));

      try {
        const response = await api.listIncidents(params);
        set({
          incidents: response.incidents || [],
          isLoading: false,
          error: null,
          total: response.total || 0,
          page: response.page || 1,
          pageSize: response.page_size || 20,
        });
      } catch (err) {
        const error = err instanceof Error ? err.message : 'Failed to load incidents';
        update((state) => ({ ...state, isLoading: false, error }));
      }
    },

    async create(data: CreateIncidentRequest) {
      update((state) => ({ ...state, isLoading: true, error: null }));

      try {
        const incident = await api.createIncident(data);
        update((state) => ({
          ...state,
          incidents: [incident, ...state.incidents],
          isLoading: false,
          total: state.total + 1,
        }));
        return incident;
      } catch (err) {
        const error = err instanceof Error ? err.message : 'Failed to create incident';
        update((state) => ({ ...state, isLoading: false, error }));
        throw err;
      }
    },

    async update(id: string, data: UpdateIncidentRequest) {
      update((state) => ({ ...state, isLoading: true, error: null }));

      try {
        const updatedIncident = await api.updateIncident(id, data);
        update((state) => ({
          ...state,
          incidents: state.incidents.map((incident) =>
            incident.id === id ? updatedIncident : incident
          ),
          isLoading: false,
        }));
        return updatedIncident;
      } catch (err) {
        const error = err instanceof Error ? err.message : 'Failed to update incident';
        update((state) => ({ ...state, isLoading: false, error }));
        throw err;
      }
    },

    async delete(id: string) {
      update((state) => ({ ...state, isLoading: true, error: null }));

      try {
        await api.deleteIncident(id);
        update((state) => ({
          ...state,
          incidents: state.incidents.filter((incident) => incident.id !== id),
          isLoading: false,
          total: state.total - 1,
        }));
      } catch (err) {
        const error = err instanceof Error ? err.message : 'Failed to delete incident';
        update((state) => ({ ...state, isLoading: false, error }));
        throw err;
      }
    },

    reset() {
      set({
        incidents: [],
        isLoading: false,
        error: null,
        total: 0,
        page: 1,
        pageSize: 20,
      });
    },
  };
}

export const incidentsStore = createIncidentsStore();
