import { writable } from 'svelte/store';
import { api } from '$lib/api/client';
import type { Schedule } from '$lib/types/schedule';

interface SchedulesState {
	schedules: Schedule[];
	isLoading: boolean;
	error: string | null;
}

function createSchedulesStore() {
	const { subscribe, set, update } = writable<SchedulesState>({
		schedules: [],
		isLoading: false,
		error: null
	});

	return {
		subscribe,

		async load() {
			update((state) => ({ ...state, isLoading: true, error: null }));

			try {
				const response = await api.listSchedules();
				update((state) => ({
					...state,
					schedules: response.schedules,
					isLoading: false
				}));
			} catch (err) {
				update((state) => ({
					...state,
					error: err instanceof Error ? err.message : 'Failed to load schedules',
					isLoading: false
				}));
			}
		},

		async create(data: { name: string; description?: string; timezone?: string; team_id?: string }) {
			try {
				const schedule = await api.createSchedule(data);
				update((state) => ({
					...state,
					schedules: [schedule, ...state.schedules]
				}));
			} catch (err) {
				throw err;
			}
		},

		async update(id: string, data: { name?: string; description?: string; timezone?: string; team_id?: string }) {
			try {
				const schedule = await api.updateSchedule(id, data);
				update((state) => ({
					...state,
					schedules: state.schedules.map((s) => (s.id === id ? schedule : s))
				}));
			} catch (err) {
				throw err;
			}
		},

		async delete(id: string) {
			try {
				await api.deleteSchedule(id);
				update((state) => ({
					...state,
					schedules: state.schedules.filter((s) => s.id !== id)
				}));
			} catch (err) {
				throw err;
			}
		}
	};
}

export const schedulesStore = createSchedulesStore();
