import { writable } from 'svelte/store';
import { api } from '$lib/api/client';
import type { Team, CreateTeamRequest, UpdateTeamRequest } from '$lib/types/team';

interface TeamsState {
	teams: Team[];
	isLoading: boolean;
	error: string | null;
}

function createTeamsStore() {
	const { subscribe, set, update } = writable<TeamsState>({
		teams: [],
		isLoading: false,
		error: null
	});

	return {
		subscribe,

		async load() {
			update(state => ({ ...state, isLoading: true, error: null }));

			try {
				const response = await api.listTeams();
				set({
					teams: response.teams || [],
					isLoading: false,
					error: null
				});
			} catch (error) {
				update(state => ({
					...state,
					isLoading: false,
					error: error instanceof Error ? error.message : 'Failed to load teams'
				}));
			}
		},

		async create(data: CreateTeamRequest) {
			try {
				const team = await api.createTeam(data);
				update(state => ({
					...state,
					teams: [team, ...state.teams]
				}));
				return team;
			} catch (error) {
				throw error;
			}
		},

		async update(id: string, data: UpdateTeamRequest) {
			try {
				const updatedTeam = await api.updateTeam(id, data);
				update(state => ({
					...state,
					teams: state.teams.map(t => (t.id === id ? updatedTeam : t))
				}));
				return updatedTeam;
			} catch (error) {
				throw error;
			}
		},

		async delete(id: string) {
			try {
				await api.deleteTeam(id);
				update(state => ({
					...state,
					teams: state.teams.filter(t => t.id !== id)
				}));
			} catch (error) {
				throw error;
			}
		}
	};
}

export const teamsStore = createTeamsStore();
