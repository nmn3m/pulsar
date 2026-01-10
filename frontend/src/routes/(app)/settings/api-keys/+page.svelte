<script lang="ts">
  import { onMount } from 'svelte';
  import { apiKeysStore } from '$lib/stores/apikeys';
  import Button from '$lib/components/ui/Button.svelte';
  import Input from '$lib/components/ui/Input.svelte';

  let showCreateForm = false;
  let name = '';
  let selectedScopes: string[] = [];
  let expiresIn: string = '';
  let createError = '';
  let creating = false;
  let copiedKey = false;

  onMount(() => {
    apiKeysStore.load();
  });

  function getExpirationDate(): string | undefined {
    if (!expiresIn) return undefined;

    const now = new Date();
    switch (expiresIn) {
      case '7d':
        now.setDate(now.getDate() + 7);
        break;
      case '30d':
        now.setDate(now.getDate() + 30);
        break;
      case '90d':
        now.setDate(now.getDate() + 90);
        break;
      case '1y':
        now.setFullYear(now.getFullYear() + 1);
        break;
      default:
        return undefined;
    }
    return now.toISOString();
  }

  async function handleCreateKey() {
    if (!name.trim() || selectedScopes.length === 0) {
      createError = 'Name and at least one scope are required';
      return;
    }

    createError = '';
    creating = true;

    try {
      await apiKeysStore.create({
        name: name.trim(),
        scopes: selectedScopes,
        expires_at: getExpirationDate(),
      });

      name = '';
      selectedScopes = [];
      expiresIn = '';
      showCreateForm = false;
    } catch (err) {
      createError = err instanceof Error ? err.message : 'Failed to create API key';
    } finally {
      creating = false;
    }
  }

  async function handleRevokeKey(id: string, keyName: string) {
    if (!confirm(`Are you sure you want to revoke API key "${keyName}"? This cannot be undone.`)) {
      return;
    }

    try {
      await apiKeysStore.revoke(id);
    } catch (err) {
      alert(err instanceof Error ? err.message : 'Failed to revoke API key');
    }
  }

  async function handleDeleteKey(id: string, keyName: string) {
    if (!confirm(`Are you sure you want to permanently delete API key "${keyName}"?`)) {
      return;
    }

    try {
      await apiKeysStore.delete(id);
    } catch (err) {
      alert(err instanceof Error ? err.message : 'Failed to delete API key');
    }
  }

  function toggleScope(scope: string) {
    if (selectedScopes.includes(scope)) {
      selectedScopes = selectedScopes.filter((s) => s !== scope);
    } else {
      selectedScopes = [...selectedScopes, scope];
    }
  }

  async function copyToClipboard(text: string) {
    try {
      await navigator.clipboard.writeText(text);
      copiedKey = true;
      setTimeout(() => (copiedKey = false), 2000);
    } catch {
      alert('Failed to copy to clipboard');
    }
  }

  function formatDate(dateStr: string): string {
    return new Date(dateStr).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
    });
  }

  function isExpired(dateStr?: string): boolean {
    if (!dateStr) return false;
    return new Date(dateStr) < new Date();
  }
</script>

<svelte:head>
  <title>API Keys - Pulsar</title>
</svelte:head>

<div class="space-y-6">
  <div class="flex justify-between items-center">
    <div>
      <h2 class="text-3xl font-bold text-gray-900">API Keys</h2>
      <p class="mt-2 text-gray-500">Manage API keys for programmatic access to Pulsar</p>
    </div>
    <Button
      variant="primary"
      on:click={() => {
        showCreateForm = !showCreateForm;
        if (!showCreateForm) {
          name = '';
          selectedScopes = [];
          expiresIn = '';
          createError = '';
        }
      }}
    >
      {showCreateForm ? 'Cancel' : 'Create API Key'}
    </Button>
  </div>

  <!-- Newly Created Key Display -->
  {#if $apiKeysStore.newlyCreatedKey}
    <div class="bg-green-50 border border-green-200 p-6 rounded-xl">
      <div class="flex items-start justify-between">
        <div>
          <h3 class="text-lg font-semibold text-green-800">API Key Created Successfully</h3>
          <p class="text-sm text-green-700 mt-1">
            Make sure to copy your API key now. You won't be able to see it again!
          </p>
        </div>
        <button
          on:click={() => apiKeysStore.clearNewlyCreatedKey()}
          class="text-green-600 hover:text-green-800"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M6 18L18 6M6 6l12 12"
            />
          </svg>
        </button>
      </div>
      <div class="mt-4 flex items-center gap-2">
        <code
          class="flex-1 bg-white px-4 py-2 rounded border border-green-300 text-sm font-mono text-gray-900 break-all"
        >
          {$apiKeysStore.newlyCreatedKey.key}
        </code>
        <Button
          variant="primary"
          size="sm"
          on:click={() => copyToClipboard($apiKeysStore.newlyCreatedKey?.key || '')}
        >
          {copiedKey ? 'Copied!' : 'Copy'}
        </Button>
      </div>
    </div>
  {/if}

  <!-- Create Form -->
  {#if showCreateForm}
    <div class="bg-white backdrop-blur-sm p-6 rounded-xl border border-primary-200 shadow-sm">
      <h3 class="text-lg font-semibold mb-4 text-gray-900">Create New API Key</h3>
      <form on:submit|preventDefault={handleCreateKey} class="space-y-4">
        <Input id="name" label="Key Name" bind:value={name} placeholder="My API Key" required />

        <div>
          <label class="block text-sm font-medium text-gray-600 mb-2">
            Scopes (select at least one)
          </label>
          <div class="flex flex-wrap gap-2">
            {#each $apiKeysStore.scopes as scope}
              <button
                type="button"
                on:click={() => toggleScope(scope)}
                class="px-3 py-1.5 rounded-lg text-sm font-medium transition-all duration-200 {selectedScopes.includes(
                  scope
                )
                  ? 'bg-primary-600 text-white'
                  : 'bg-gray-100 text-gray-600 hover:bg-gray-200 border border-gray-300'}"
              >
                {scope}
              </button>
            {/each}
          </div>
          {#if $apiKeysStore.scopes.length === 0}
            <p class="text-sm text-gray-500 mt-2">Loading available scopes...</p>
          {/if}
        </div>

        <div>
          <label for="expires" class="block text-sm font-medium text-gray-600 mb-1">
            Expiration (optional)
          </label>
          <select
            id="expires"
            bind:value={expiresIn}
            class="w-full px-3 py-2 bg-white border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500 text-gray-900"
          >
            <option value="">Never expires</option>
            <option value="7d">7 days</option>
            <option value="30d">30 days</option>
            <option value="90d">90 days</option>
            <option value="1y">1 year</option>
          </select>
        </div>

        {#if createError}
          <div class="bg-red-50 border border-red-200 text-red-600 px-4 py-3 rounded-lg">
            {createError}
          </div>
        {/if}

        <div class="flex gap-2">
          <Button type="submit" variant="primary" disabled={creating}>
            {creating ? 'Creating...' : 'Create API Key'}
          </Button>
          <Button
            type="button"
            variant="secondary"
            on:click={() => {
              showCreateForm = false;
              name = '';
              selectedScopes = [];
              expiresIn = '';
              createError = '';
            }}
          >
            Cancel
          </Button>
        </div>
      </form>
    </div>
  {/if}

  <!-- API Keys List -->
  {#if $apiKeysStore.isLoading}
    <div class="text-center py-12">
      <div
        class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-primary-500"
      ></div>
      <p class="mt-2 text-gray-500">Loading API keys...</p>
    </div>
  {:else if $apiKeysStore.error}
    <div class="bg-red-50 border border-red-200 text-red-600 px-4 py-3 rounded-lg">
      {$apiKeysStore.error}
    </div>
  {:else if $apiKeysStore.apiKeys.length === 0}
    <div
      class="text-center py-12 bg-white backdrop-blur-sm rounded-xl border border-gray-200 shadow-sm"
    >
      <p class="text-gray-600">No API keys found</p>
      <p class="text-sm text-gray-400 mt-2">Create your first API key to get started</p>
    </div>
  {:else}
    <div
      class="bg-white backdrop-blur-sm rounded-xl border border-gray-200 overflow-hidden shadow-sm"
    >
      <div class="px-6 py-3 bg-gray-50 border-b border-gray-200">
        <p class="text-sm text-gray-600">
          {$apiKeysStore.apiKeys.length} API key{$apiKeysStore.apiKeys.length !== 1 ? 's' : ''}
        </p>
      </div>
      <ul class="divide-y divide-gray-200">
        {#each $apiKeysStore.apiKeys as key (key.id)}
          <li class="px-6 py-4">
            <div class="flex items-start justify-between">
              <div class="flex-1 min-w-0">
                <div class="flex items-center gap-2 mb-2">
                  <h4 class="text-sm font-semibold text-gray-900">{key.name}</h4>
                  {#if !key.is_active}
                    <span class="px-2 py-0.5 text-xs font-medium rounded bg-red-100 text-red-700">
                      Revoked
                    </span>
                  {:else if isExpired(key.expires_at)}
                    <span
                      class="px-2 py-0.5 text-xs font-medium rounded bg-yellow-100 text-yellow-700"
                    >
                      Expired
                    </span>
                  {:else}
                    <span
                      class="px-2 py-0.5 text-xs font-medium rounded bg-green-100 text-green-700"
                    >
                      Active
                    </span>
                  {/if}
                </div>

                <div class="flex items-center gap-4 text-sm text-gray-500">
                  <span class="font-mono bg-gray-100 px-2 py-0.5 rounded">
                    {key.key_prefix}...
                  </span>
                  <span>Created {formatDate(key.created_at)}</span>
                  {#if key.expires_at}
                    <span>
                      {isExpired(key.expires_at) ? 'Expired' : 'Expires'}
                      {formatDate(key.expires_at)}
                    </span>
                  {:else}
                    <span>Never expires</span>
                  {/if}
                  {#if key.last_used_at}
                    <span>Last used {formatDate(key.last_used_at)}</span>
                  {:else}
                    <span>Never used</span>
                  {/if}
                </div>

                <div class="flex flex-wrap gap-1 mt-2">
                  {#each key.scopes as scope}
                    <span class="px-2 py-0.5 text-xs font-medium rounded bg-gray-100 text-gray-600">
                      {scope}
                    </span>
                  {/each}
                </div>
              </div>

              <div class="flex items-center gap-2 ml-4">
                {#if key.is_active && !isExpired(key.expires_at)}
                  <Button
                    variant="secondary"
                    size="sm"
                    on:click={() => handleRevokeKey(key.id, key.name)}
                  >
                    Revoke
                  </Button>
                {/if}
                <Button
                  variant="danger"
                  size="sm"
                  on:click={() => handleDeleteKey(key.id, key.name)}
                >
                  Delete
                </Button>
              </div>
            </div>
          </li>
        {/each}
      </ul>
    </div>
  {/if}

  <!-- Usage Info -->
  <div class="bg-blue-50 border border-blue-200 p-4 rounded-xl">
    <h4 class="text-sm font-semibold text-blue-800 mb-2">Using API Keys</h4>
    <p class="text-sm text-blue-700">
      Include your API key in requests using the <code class="bg-blue-100 px-1 rounded"
        >X-API-Key</code
      > header:
    </p>
    <pre
      class="mt-2 bg-blue-100 p-3 rounded text-xs font-mono text-blue-900 overflow-x-auto">curl -H "X-API-Key: your-api-key" https://api.pulsar.example.com/api/v1/alerts</pre>
  </div>
</div>
