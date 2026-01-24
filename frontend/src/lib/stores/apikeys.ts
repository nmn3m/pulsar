import { writable } from 'svelte/store';
import { api } from '$lib/api/client';
import type {
  APIKey,
  APIKeyResponse,
  CreateAPIKeyRequest,
  UpdateAPIKeyRequest,
} from '$lib/types/apikey';

interface APIKeysState {
  apiKeys: APIKey[];
  scopes: string[];
  isLoading: boolean;
  error: string | null;
  newlyCreatedKey: APIKeyResponse | null;
}

function createAPIKeysStore() {
  const { subscribe, set, update } = writable<APIKeysState>({
    apiKeys: [],
    scopes: [],
    isLoading: false,
    error: null,
    newlyCreatedKey: null,
  });

  return {
    subscribe,

    async load() {
      update((state) => ({ ...state, isLoading: true, error: null }));

      try {
        const [keysResponse, scopesResponse] = await Promise.all([
          api.listAPIKeys(),
          api.getAPIKeyScopes(),
        ]);
        set({
          apiKeys: keysResponse.api_keys || [],
          scopes: scopesResponse.scopes || [],
          isLoading: false,
          error: null,
          newlyCreatedKey: null,
        });
      } catch (error) {
        update((state) => ({
          ...state,
          isLoading: false,
          error: error instanceof Error ? error.message : 'Failed to load API keys',
        }));
      }
    },

    async create(data: CreateAPIKeyRequest): Promise<APIKeyResponse> {
      const key = await api.createAPIKey(data);
      update((state) => ({
        ...state,
        apiKeys: [key, ...state.apiKeys],
        newlyCreatedKey: key,
      }));
      return key;
    },

    async updateKey(id: string, data: UpdateAPIKeyRequest) {
      const updatedKey = await api.updateAPIKey(id, data);
      update((state) => ({
        ...state,
        apiKeys: state.apiKeys.map((k) => (k.id === id ? updatedKey : k)),
      }));
      return updatedKey;
    },

    async revoke(id: string) {
      await api.revokeAPIKey(id);
      update((state) => ({
        ...state,
        apiKeys: state.apiKeys.map((k) => (k.id === id ? { ...k, is_active: false } : k)),
      }));
    },

    async delete(id: string) {
      await api.deleteAPIKey(id);
      update((state) => ({
        ...state,
        apiKeys: state.apiKeys.filter((k) => k.id !== id),
      }));
    },

    clearNewlyCreatedKey() {
      update((state) => ({
        ...state,
        newlyCreatedKey: null,
      }));
    },
  };
}

export const apiKeysStore = createAPIKeysStore();
