import { writable } from 'svelte/store';
import { api } from '$lib/api/client';
import type { Schedule } from '$lib/types/schedule';

interface SchedulesState {
  schedules: Schedule[];
  isLoading: boolean;
  error: string | null;
}

function createSchedulesStore() {
  const { subscribe, update } = writable<SchedulesState>({
    schedules: [],
    isLoading: false,
    error: null,
  });

  return {
    subscribe,

    async load() {
      update((state) => ({ ...state, isLoading: true, error: null }));

      try {
        const response = await api.listSchedules();
        update((state) => ({
          ...state,
          schedules: response.schedules || [],
          isLoading: false,
        }));
      } catch (err) {
        update((state) => ({
          ...state,
          error: err instanceof Error ? err.message : 'Failed to load schedules',
          isLoading: false,
        }));
      }
    },

    async create(data: {
      name: string;
      description?: string;
      timezone?: string;
      team_id?: string;
    }) {
      const schedule = await api.createSchedule(data);
      update((state) => ({
        ...state,
        schedules: [schedule, ...state.schedules],
      }));
    },

    async update(
      id: string,
      data: { name?: string; description?: string; timezone?: string; team_id?: string }
    ) {
      const schedule = await api.updateSchedule(id, data);
      update((state) => ({
        ...state,
        schedules: state.schedules.map((s) => (s.id === id ? schedule : s)),
      }));
    },

    async delete(id: string) {
      await api.deleteSchedule(id);
      update((state) => ({
        ...state,
        schedules: state.schedules.filter((s) => s.id !== id),
      }));
    },
  };
}

export const schedulesStore = createSchedulesStore();
