<script lang="ts">
  import { onMount } from 'svelte';
  import { api } from '$lib/api/client';
  import Button from '$lib/components/ui/Button.svelte';
  import type { WebhookEndpoint, CreateWebhookEndpointRequest } from '$lib/types/webhook';

  let endpoints: WebhookEndpoint[] = [];
  let isLoading = true;
  let error: string | null = null;

  // Create form state
  let showCreateForm = false;
  let name = '';
  let url = '';
  let enabled = true;
  let alertCreated = false;
  let alertUpdated = false;
  let alertAcknowledged = false;
  let alertClosed = false;
  let alertEscalated = false;
  let incidentCreated = false;
  let incidentUpdated = false;
  let incidentResolved = false;
  let customHeaders: Record<string, string> = {};
  let headerKey = '';
  let headerValue = '';
  let timeoutSeconds = 30;
  let maxRetries = 3;
  let retryDelaySeconds = 60;
  let createError = '';
  let creating = false;

  onMount(async () => {
    await loadEndpoints();
  });

  async function loadEndpoints() {
    isLoading = true;
    error = null;

    try {
      endpoints = (await api.listWebhookEndpoints()) || [];
    } catch (err) {
      error = err instanceof Error ? err.message : 'Failed to load webhook endpoints';
    } finally {
      isLoading = false;
    }
  }

  async function handleCreate() {
    createError = '';
    creating = true;

    try {
      const data: CreateWebhookEndpointRequest = {
        name: name.trim(),
        url: url.trim(),
        enabled,
        alert_created: alertCreated,
        alert_updated: alertUpdated,
        alert_acknowledged: alertAcknowledged,
        alert_closed: alertClosed,
        alert_escalated: alertEscalated,
        incident_created: incidentCreated,
        incident_updated: incidentUpdated,
        incident_resolved: incidentResolved,
        headers: customHeaders,
        timeout_seconds: timeoutSeconds,
        max_retries: maxRetries,
        retry_delay_seconds: retryDelaySeconds,
      };

      await api.createWebhookEndpoint(data);
      await loadEndpoints();

      // Reset form
      name = '';
      url = '';
      enabled = true;
      alertCreated = false;
      alertUpdated = false;
      alertAcknowledged = false;
      alertClosed = false;
      alertEscalated = false;
      incidentCreated = false;
      incidentUpdated = false;
      incidentResolved = false;
      customHeaders = {};
      timeoutSeconds = 30;
      maxRetries = 3;
      retryDelaySeconds = 60;
      showCreateForm = false;
    } catch (err) {
      createError = err instanceof Error ? err.message : 'Failed to create webhook endpoint';
    } finally {
      creating = false;
    }
  }

  async function handleDelete(id: string) {
    if (!confirm('Are you sure you want to delete this webhook endpoint?')) return;

    try {
      await api.deleteWebhookEndpoint(id);
      await loadEndpoints();
    } catch (err) {
      error = err instanceof Error ? err.message : 'Failed to delete webhook endpoint';
    }
  }

  async function toggleEnabled(endpoint: WebhookEndpoint) {
    try {
      await api.updateWebhookEndpoint(endpoint.id, { enabled: !endpoint.enabled });
      await loadEndpoints();
    } catch (err) {
      error = err instanceof Error ? err.message : 'Failed to update webhook endpoint';
    }
  }

  function addHeader() {
    if (headerKey.trim() && headerValue.trim()) {
      customHeaders[headerKey.trim()] = headerValue.trim();
      customHeaders = { ...customHeaders }; // Trigger reactivity
      headerKey = '';
      headerValue = '';
    }
  }

  function removeHeader(key: string) {
    delete customHeaders[key];
    customHeaders = { ...customHeaders }; // Trigger reactivity
  }

  function getEventFilters(endpoint: WebhookEndpoint): string[] {
    const filters: string[] = [];
    if (endpoint.alert_created) filters.push('Alert Created');
    if (endpoint.alert_updated) filters.push('Alert Updated');
    if (endpoint.alert_acknowledged) filters.push('Alert Acknowledged');
    if (endpoint.alert_closed) filters.push('Alert Closed');
    if (endpoint.alert_escalated) filters.push('Alert Escalated');
    if (endpoint.incident_created) filters.push('Incident Created');
    if (endpoint.incident_updated) filters.push('Incident Updated');
    if (endpoint.incident_resolved) filters.push('Incident Resolved');
    return filters;
  }
</script>

<svelte:head>
  <title>Webhook Endpoints - Pulsar</title>
</svelte:head>

<div class="space-y-6">
  <div class="flex justify-between items-center">
    <div>
      <h2 class="text-3xl font-bold text-gray-900">Webhook Endpoints</h2>
      <p class="mt-2 text-gray-500">Send events to external services via HTTP webhooks</p>
    </div>
    <Button variant="primary" on:click={() => (showCreateForm = !showCreateForm)}>
      {showCreateForm ? 'Cancel' : 'Create Endpoint'}
    </Button>
  </div>

  {#if error}
    <div class="bg-red-50 border border-red-200 text-red-600 px-4 py-3 rounded-lg">
      {error}
    </div>
  {/if}

  {#if showCreateForm}
    <div class="bg-white backdrop-blur-sm p-6 rounded-xl border border-primary-200 shadow-sm">
      <h3 class="text-lg font-semibold mb-4 text-gray-900">Create Webhook Endpoint</h3>
      <form on:submit|preventDefault={handleCreate} class="space-y-4">
        <div class="grid grid-cols-2 gap-4">
          <div>
            <label for="name" class="block text-sm font-medium text-gray-600 mb-1"> Name * </label>
            <input
              id="name"
              type="text"
              bind:value={name}
              required
              class="w-full px-3 py-2 bg-white border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500 text-gray-900 placeholder-gray-400"
              placeholder="My Webhook"
            />
          </div>

          <div>
            <label for="url" class="block text-sm font-medium text-gray-600 mb-1"> URL * </label>
            <input
              id="url"
              type="url"
              bind:value={url}
              required
              class="w-full px-3 py-2 bg-white border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500 text-gray-900 placeholder-gray-400"
              placeholder="https://example.com/webhook"
            />
          </div>
        </div>

        <div>
          <label class="flex items-center space-x-2">
            <input
              type="checkbox"
              bind:checked={enabled}
              class="rounded bg-white border-gray-300 text-primary-500 focus:ring-primary-500"
            />
            <span class="text-sm font-medium text-gray-600">Enabled</span>
          </label>
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-600 mb-2"> Event Triggers </label>
          <div class="grid grid-cols-2 gap-2">
            <label class="flex items-center space-x-2">
              <input
                type="checkbox"
                bind:checked={alertCreated}
                class="rounded bg-white border-gray-300 text-primary-500 focus:ring-primary-500"
              />
              <span class="text-sm text-gray-600">Alert Created</span>
            </label>
            <label class="flex items-center space-x-2">
              <input
                type="checkbox"
                bind:checked={alertUpdated}
                class="rounded bg-white border-gray-300 text-primary-500 focus:ring-primary-500"
              />
              <span class="text-sm text-gray-600">Alert Updated</span>
            </label>
            <label class="flex items-center space-x-2">
              <input
                type="checkbox"
                bind:checked={alertAcknowledged}
                class="rounded bg-white border-gray-300 text-primary-500 focus:ring-primary-500"
              />
              <span class="text-sm text-gray-600">Alert Acknowledged</span>
            </label>
            <label class="flex items-center space-x-2">
              <input
                type="checkbox"
                bind:checked={alertClosed}
                class="rounded bg-white border-gray-300 text-primary-500 focus:ring-primary-500"
              />
              <span class="text-sm text-gray-600">Alert Closed</span>
            </label>
            <label class="flex items-center space-x-2">
              <input
                type="checkbox"
                bind:checked={alertEscalated}
                class="rounded bg-white border-gray-300 text-primary-500 focus:ring-primary-500"
              />
              <span class="text-sm text-gray-600">Alert Escalated</span>
            </label>
            <label class="flex items-center space-x-2">
              <input
                type="checkbox"
                bind:checked={incidentCreated}
                class="rounded bg-white border-gray-300 text-primary-500 focus:ring-primary-500"
              />
              <span class="text-sm text-gray-600">Incident Created</span>
            </label>
            <label class="flex items-center space-x-2">
              <input
                type="checkbox"
                bind:checked={incidentUpdated}
                class="rounded bg-white border-gray-300 text-primary-500 focus:ring-primary-500"
              />
              <span class="text-sm text-gray-600">Incident Updated</span>
            </label>
            <label class="flex items-center space-x-2">
              <input
                type="checkbox"
                bind:checked={incidentResolved}
                class="rounded bg-white border-gray-300 text-primary-500 focus:ring-primary-500"
              />
              <span class="text-sm text-gray-600">Incident Resolved</span>
            </label>
          </div>
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-600 mb-2"> Custom Headers </label>
          <div class="flex gap-2 mb-2">
            <input
              type="text"
              bind:value={headerKey}
              placeholder="Header name"
              class="flex-1 px-3 py-2 bg-white border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500 text-gray-900 placeholder-gray-400"
            />
            <input
              type="text"
              bind:value={headerValue}
              placeholder="Header value"
              class="flex-1 px-3 py-2 bg-white border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500 text-gray-900 placeholder-gray-400"
            />
            <Button type="button" variant="secondary" on:click={addHeader}>Add</Button>
          </div>
          {#if Object.keys(customHeaders).length > 0}
            <div class="space-y-1">
              {#each Object.entries(customHeaders) as [key, value]}
                <div
                  class="flex items-center justify-between bg-gray-100 px-3 py-2 rounded-lg border border-gray-200"
                >
                  <span class="text-sm text-gray-700"
                    ><strong class="text-primary-600">{key}:</strong>
                    {value}</span
                  >
                  <button
                    type="button"
                    on:click={() => removeHeader(key)}
                    class="text-error-dark hover:text-error text-sm"
                  >
                    Remove
                  </button>
                </div>
              {/each}
            </div>
          {/if}
        </div>

        <div class="grid grid-cols-3 gap-4">
          <div>
            <label for="timeout" class="block text-sm font-medium text-gray-600 mb-1">
              Timeout (seconds)
            </label>
            <input
              id="timeout"
              type="number"
              bind:value={timeoutSeconds}
              min="1"
              max="300"
              class="w-full px-3 py-2 bg-white border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500 text-gray-900"
            />
          </div>

          <div>
            <label for="retries" class="block text-sm font-medium text-gray-600 mb-1">
              Max Retries
            </label>
            <input
              id="retries"
              type="number"
              bind:value={maxRetries}
              min="0"
              max="10"
              class="w-full px-3 py-2 bg-white border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500 text-gray-900"
            />
          </div>

          <div>
            <label for="delay" class="block text-sm font-medium text-gray-600 mb-1">
              Retry Delay (seconds)
            </label>
            <input
              id="delay"
              type="number"
              bind:value={retryDelaySeconds}
              min="1"
              max="3600"
              class="w-full px-3 py-2 bg-white border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500 text-gray-900"
            />
          </div>
        </div>

        {#if createError}
          <div class="bg-red-50 border border-red-200 text-red-600 px-4 py-3 rounded-lg">
            {createError}
          </div>
        {/if}

        <div class="flex gap-2">
          <Button type="submit" variant="primary" disabled={creating}>
            {creating ? 'Creating...' : 'Create Endpoint'}
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
      <p class="mt-2 text-gray-500">Loading webhook endpoints...</p>
    </div>
  {:else if endpoints.length === 0}
    <div
      class="text-center py-12 bg-white backdrop-blur-sm rounded-xl border border-gray-200 shadow-sm"
    >
      <p class="text-gray-600">No webhook endpoints configured</p>
      <p class="text-sm text-gray-400 mt-2">
        Create your first webhook to send events to external services
      </p>
    </div>
  {:else}
    <div class="grid grid-cols-1 gap-4">
      {#each endpoints as endpoint (endpoint.id)}
        <div
          class="bg-white backdrop-blur-sm p-6 rounded-xl border border-gray-200 hover:border-primary-400 transition-all duration-300 shadow-sm hover:shadow-lg"
        >
          <div class="flex items-start justify-between">
            <div class="flex-1">
              <div class="flex items-center gap-3 mb-2">
                <h3 class="text-lg font-semibold text-gray-900">
                  {endpoint.name}
                </h3>
                <button
                  type="button"
                  on:click={() => toggleEnabled(endpoint)}
                  class="px-2 py-1 text-xs rounded border {endpoint.enabled
                    ? 'bg-green-100 text-green-700 border-green-300'
                    : 'bg-gray-100 text-gray-500 border-gray-300'}"
                >
                  {endpoint.enabled ? 'Enabled' : 'Disabled'}
                </button>
              </div>
              <p class="text-sm text-primary-600 mb-3 font-mono break-all">
                {endpoint.url}
              </p>

              <div class="space-y-2">
                <div class="flex flex-wrap gap-1">
                  <span class="text-xs font-medium text-gray-500">Events:</span>
                  {#each getEventFilters(endpoint) as filter}
                    <span
                      class="px-2 py-1 bg-primary-100 text-primary-700 text-xs rounded border border-primary-200"
                      >{filter}</span
                    >
                  {/each}
                </div>

                <div class="text-xs text-gray-400">
                  Timeout: {endpoint.timeout_seconds}s | Max Retries: {endpoint.max_retries} | Retry Delay:
                  {endpoint.retry_delay_seconds}s
                </div>

                {#if Object.keys(endpoint.headers).length > 0}
                  <div class="text-xs text-gray-400">
                    Custom Headers: {Object.keys(endpoint.headers).length}
                  </div>
                {/if}
              </div>
            </div>

            <div class="flex gap-2">
              <Button variant="danger" size="sm" on:click={() => handleDelete(endpoint.id)}>
                Delete
              </Button>
            </div>
          </div>
        </div>
      {/each}
    </div>
  {/if}
</div>
