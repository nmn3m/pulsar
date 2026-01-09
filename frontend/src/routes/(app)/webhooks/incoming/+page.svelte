<script lang="ts">
  import { onMount } from 'svelte';
  import { api } from '$lib/api/client';
  import Button from '$lib/components/ui/Button.svelte';
  import type {
    IncomingWebhookToken,
    CreateIncomingWebhookTokenRequest,
    IncomingWebhookIntegrationType,
  } from '$lib/types/webhook';

  let tokens: IncomingWebhookToken[] = [];
  let isLoading = true;
  let error: string | null = null;

  // Create form state
  let showCreateForm = false;
  let name = '';
  let integrationType: IncomingWebhookIntegrationType = 'generic';
  let defaultPriority = 'P3';
  let defaultTags = '';
  let createError = '';
  let creating = false;

  onMount(async () => {
    await loadTokens();
  });

  async function loadTokens() {
    isLoading = true;
    error = null;

    try {
      tokens = (await api.listIncomingWebhookTokens()) || [];
    } catch (err) {
      error = err instanceof Error ? err.message : 'Failed to load incoming webhook tokens';
    } finally {
      isLoading = false;
    }
  }

  async function handleCreate() {
    createError = '';
    creating = true;

    try {
      const data: CreateIncomingWebhookTokenRequest = {
        name: name.trim(),
        integration_type: integrationType,
        default_priority: defaultPriority,
        default_tags: defaultTags ? defaultTags.split(',').map((t) => t.trim()) : [],
      };

      await api.createIncomingWebhookToken(data);
      await loadTokens();

      // Reset form
      name = '';
      integrationType = 'generic';
      defaultPriority = 'P3';
      defaultTags = '';
      showCreateForm = false;
    } catch (err) {
      createError = err instanceof Error ? err.message : 'Failed to create incoming webhook token';
    } finally {
      creating = false;
    }
  }

  async function handleDelete(id: string) {
    if (!confirm('Are you sure you want to delete this incoming webhook token?')) return;

    try {
      await api.deleteIncomingWebhookToken(id);
      await loadTokens();
    } catch (err) {
      error = err instanceof Error ? err.message : 'Failed to delete incoming webhook token';
    }
  }

  function copyToClipboard(text: string) {
    navigator.clipboard.writeText(text);
    alert('Copied to clipboard!');
  }

  function getWebhookURL(token: string): string {
    const baseUrl = window.location.origin;
    return `${baseUrl}/api/v1/webhook/${token}`;
  }

  function formatDate(dateStr?: string): string {
    if (!dateStr) return 'Never';
    const date = new Date(dateStr);
    return date.toLocaleString();
  }
</script>

<svelte:head>
  <title>Incoming Webhooks - Pulsar</title>
</svelte:head>

<div class="space-y-6">
  <div class="flex justify-between items-center">
    <div>
      <h2 class="text-3xl font-bold text-gray-900">Incoming Webhooks</h2>
      <p class="mt-2 text-gray-500">
        Receive alerts from external monitoring tools
      </p>
    </div>
    <Button variant="primary" on:click={() => (showCreateForm = !showCreateForm)}>
      {showCreateForm ? 'Cancel' : 'Create Token'}
    </Button>
  </div>

  {#if error}
    <div
      class="bg-red-50 border border-red-200 text-red-600 px-4 py-3 rounded-lg"
    >
      {error}
    </div>
  {/if}

  {#if showCreateForm}
    <div
      class="bg-white backdrop-blur-sm p-6 rounded-xl border border-primary-200 shadow-sm"
    >
      <h3 class="text-lg font-semibold mb-4 text-gray-900">
        Create Incoming Webhook Token
      </h3>
      <form on:submit|preventDefault={handleCreate} class="space-y-4">
        <div>
          <label for="name" class="block text-sm font-medium text-gray-600 mb-1">
            Name *
          </label>
          <input
            id="name"
            type="text"
            bind:value={name}
            required
            class="w-full px-3 py-2 bg-white border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500 text-gray-900 placeholder-gray-400"
            placeholder="Prometheus Production"
          />
        </div>

        <div>
          <label
            for="integration-type"
            class="block text-sm font-medium text-gray-600 mb-1"
          >
            Integration Type *
          </label>
          <select
            id="integration-type"
            bind:value={integrationType}
            class="w-full px-3 py-2 bg-white border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500 text-gray-900"
          >
            <option value="generic">Generic</option>
            <option value="prometheus">Prometheus Alertmanager</option>
            <option value="grafana">Grafana</option>
            <option value="datadog">Datadog</option>
          </select>
        </div>

        <div>
          <label
            for="priority"
            class="block text-sm font-medium text-gray-600 mb-1"
          >
            Default Priority
          </label>
          <select
            id="priority"
            bind:value={defaultPriority}
            class="w-full px-3 py-2 bg-white border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500 text-gray-900"
          >
            <option value="P1">P1 - Critical</option>
            <option value="P2">P2 - High</option>
            <option value="P3">P3 - Medium</option>
            <option value="P4">P4 - Low</option>
            <option value="P5">P5 - Info</option>
          </select>
        </div>

        <div>
          <label for="tags" class="block text-sm font-medium text-gray-600 mb-1">
            Default Tags (comma-separated)
          </label>
          <input
            id="tags"
            type="text"
            bind:value={defaultTags}
            class="w-full px-3 py-2 bg-white border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500 text-gray-900 placeholder-gray-400"
            placeholder="production, monitoring"
          />
        </div>

        {#if createError}
          <div
            class="bg-red-50 border border-red-200 text-red-600 px-4 py-3 rounded-lg"
          >
            {createError}
          </div>
        {/if}

        <div class="flex gap-2">
          <Button type="submit" variant="primary" disabled={creating}>
            {creating ? 'Creating...' : 'Create Token'}
          </Button>
          <Button type="button" variant="secondary" on:click={() => (showCreateForm = false)}>
            Cancel
          </Button>
        </div>
      </form>
    </div>
  {/if}

  {#if isLoading}
    <div class="text-center py-12">
      <div
        class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-primary-500"
      ></div>
      <p class="mt-2 text-gray-500">Loading incoming webhook tokens...</p>
    </div>
  {:else if tokens.length === 0}
    <div
      class="text-center py-12 bg-white backdrop-blur-sm rounded-xl border border-gray-200 shadow-sm"
    >
      <p class="text-gray-600">No incoming webhook tokens configured</p>
      <p class="text-sm text-gray-400 mt-2">
        Create a token to receive alerts from external monitoring tools
      </p>
    </div>
  {:else}
    <div class="grid grid-cols-1 gap-4">
      {#each tokens as token (token.id)}
        <div
          class="bg-white backdrop-blur-sm p-6 rounded-xl border border-gray-200 hover:border-primary-400 transition-all duration-300 shadow-sm hover:shadow-lg"
        >
          <div class="flex items-start justify-between mb-4">
            <div class="flex-1">
              <div class="flex items-center gap-3 mb-2">
                <h3 class="text-lg font-semibold text-gray-900">{token.name}</h3>
                <span
                  class="px-2 py-1 text-xs rounded bg-primary-100 text-primary-700 border border-primary-200"
                >
                  {token.integration_type}
                </span>
                <span
                  class="px-2 py-1 text-xs rounded border {token.enabled
                    ? 'bg-green-100 text-green-700 border-green-300'
                    : 'bg-gray-100 text-gray-500 border-gray-300'}"
                >
                  {token.enabled ? 'Enabled' : 'Disabled'}
                </span>
              </div>

              <div class="space-y-2 text-sm text-gray-500">
                <div>
                  <span class="font-medium text-gray-700">Default Priority:</span>
                  {token.default_priority}
                </div>
                {#if token.default_tags.length > 0}
                  <div class="flex flex-wrap gap-1">
                    <span class="font-medium text-gray-700">Default Tags:</span>
                    {#each token.default_tags as tag}
                      <span
                        class="px-2 py-0.5 bg-gray-100 text-gray-700 text-xs rounded border border-gray-200"
                        >{tag}</span
                      >
                    {/each}
                  </div>
                {/if}
                <div>
                  <span class="font-medium text-gray-700">Requests:</span>
                  {token.request_count}
                </div>
                <div>
                  <span class="font-medium text-gray-700">Last Used:</span>
                  {formatDate(token.last_used_at)}
                </div>
              </div>
            </div>

            <div class="flex gap-2">
              <Button variant="danger" size="sm" on:click={() => handleDelete(token.id)}>
                Delete
              </Button>
            </div>
          </div>

          <div class="border-t border-gray-200 pt-4">
            <label class="block text-sm font-medium text-gray-700 mb-2">
              Webhook URL
            </label>
            <div class="flex gap-2">
              <input
                type="text"
                value={getWebhookURL(token.token)}
                readonly
                class="flex-1 px-3 py-2 bg-gray-100 border border-gray-200 rounded-lg font-mono text-sm text-primary-600"
              />
              <Button
                variant="secondary"
                size="sm"
                on:click={() => copyToClipboard(getWebhookURL(token.token))}
              >
                Copy
              </Button>
            </div>
            <p class="mt-2 text-xs text-gray-400">
              Use this URL to send webhooks from {token.integration_type === 'generic'
                ? 'your service'
                : token.integration_type}
            </p>
          </div>

          {#if token.integration_type === 'prometheus'}
            <div
              class="mt-4 p-3 bg-gray-100 rounded-lg border border-gray-200 text-xs"
            >
              <p class="font-medium mb-1 text-gray-700">
                Prometheus Alertmanager Configuration:
              </p>
              <pre class="overflow-x-auto text-primary-600"><code
                  >receivers:
  - name: 'pulsar'
    webhook_configs:
      - url: '{getWebhookURL(token.token)}'
        send_resolved: false</code
                ></pre>
            </div>
          {:else if token.integration_type === 'grafana'}
            <div
              class="mt-4 p-3 bg-gray-100 rounded-lg border border-gray-200 text-xs"
            >
              <p class="font-medium mb-1 text-gray-700">
                Grafana Webhook Configuration:
              </p>
              <p class="text-gray-500">
                Add this URL as a webhook notification channel in Grafana.
              </p>
            </div>
          {/if}
        </div>
      {/each}
    </div>
  {/if}
</div>
