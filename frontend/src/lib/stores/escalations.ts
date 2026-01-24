import { writable } from 'svelte/store';
import { api } from '$lib/api/client';
import type { EscalationPolicy } from '$lib/types/escalation';

interface EscalationPoliciesState {
  policies: EscalationPolicy[];
  isLoading: boolean;
  error: string | null;
}

function createEscalationPoliciesStore() {
  const { subscribe, update } = writable<EscalationPoliciesState>({
    policies: [],
    isLoading: false,
    error: null,
  });

  return {
    subscribe,

    async load() {
      update((state) => ({ ...state, isLoading: true, error: null }));

      try {
        const response = await api.listEscalationPolicies();
        update((state) => ({
          ...state,
          policies: response.policies || [],
          isLoading: false,
        }));
      } catch (err) {
        update((state) => ({
          ...state,
          error: err instanceof Error ? err.message : 'Failed to load escalation policies',
          isLoading: false,
        }));
      }
    },

    async create(data: {
      name: string;
      description?: string;
      repeat_enabled?: boolean;
      repeat_count?: number;
    }) {
      const policy = await api.createEscalationPolicy(data);
      update((state) => ({
        ...state,
        policies: [policy, ...state.policies],
      }));
    },

    async update(
      id: string,
      data: {
        name?: string;
        description?: string;
        repeat_enabled?: boolean;
        repeat_count?: number;
      }
    ) {
      const policy = await api.updateEscalationPolicy(id, data);
      update((state) => ({
        ...state,
        policies: state.policies.map((p) => (p.id === id ? policy : p)),
      }));
    },

    async delete(id: string) {
      await api.deleteEscalationPolicy(id);
      update((state) => ({
        ...state,
        policies: state.policies.filter((p) => p.id !== id),
      }));
    },
  };
}

export const escalationPoliciesStore = createEscalationPoliciesStore();
