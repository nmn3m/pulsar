<script lang="ts">
  import { onMount } from 'svelte';
  import { api } from '$lib/api/client';
  import Button from '$lib/components/ui/Button.svelte';
  import type { WebhookDelivery } from '$lib/types/webhook';

  let deliveries: WebhookDelivery[] = [];
  let isLoading = true;
  let error: string | null = null;
  let limit = 20;
  let offset = 0;

  // Expanded delivery details
  let expandedDeliveryId: string | null = null;

  onMount(async () => {
    await loadDeliveries();
  });

  async function loadDeliveries() {
    isLoading = true;
    error = null;

    try {
      const response = await api.listWebhookDeliveries(limit, offset);
      deliveries = response.deliveries || [];
    } catch (err) {
      error = err instanceof Error ? err.message : 'Failed to load webhook deliveries';
    } finally {
      isLoading = false;
    }
  }

  function toggleExpanded(id: string) {
    expandedDeliveryId = expandedDeliveryId === id ? null : id;
  }

  function getStatusColor(status: string): string {
    switch (status) {
      case 'success':
        return 'bg-green-100 text-green-700 border border-green-300';
      case 'failed':
        return 'bg-red-100 text-red-700 border border-red-300';
      case 'pending':
        return 'bg-yellow-100 text-yellow-700 border border-yellow-300';
      default:
        return 'bg-gray-100 text-gray-700 border border-gray-300';
    }
  }

  function formatDate(dateStr: string): string {
    const date = new Date(dateStr);
    return date.toLocaleString();
  }

  function formatRelativeTime(dateStr?: string): string {
    if (!dateStr) return '';
    const date = new Date(dateStr);
    const now = new Date();
    const diffMs = now.getTime() - date.getTime();
    const diffMins = Math.floor(diffMs / 60000);

    if (diffMins < 1) return 'just now';
    if (diffMins < 60) return `${diffMins}m ago`;

    const diffHours = Math.floor(diffMins / 60);
    if (diffHours < 24) return `${diffHours}h ago`;

    const diffDays = Math.floor(diffHours / 24);
    return `${diffDays}d ago`;
  }

  async function handleNextPage() {
    offset += limit;
    await loadDeliveries();
  }

  async function handlePrevPage() {
    offset = Math.max(0, offset - limit);
    await loadDeliveries();
  }
</script>

<svelte:head>
  <title>Webhook Deliveries - Pulsar</title>
</svelte:head>

<div class="space-y-6">
  <div class="flex justify-between items-center">
    <div>
      <h2 class="text-3xl font-bold text-gray-900">Webhook Deliveries</h2>
      <p class="mt-2 text-gray-500">
        Monitor webhook delivery status and troubleshoot issues
      </p>
    </div>
    <Button variant="secondary" on:click={loadDeliveries}>Refresh</Button>
  </div>

  {#if error}
    <div
      class="bg-red-50 border border-red-200 text-red-600 px-4 py-3 rounded-lg"
    >
      {error}
    </div>
  {/if}

  {#if isLoading}
    <div class="text-center py-12">
      <div
        class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-primary-500"
      ></div>
      <p class="mt-2 text-gray-500">Loading webhook deliveries...</p>
    </div>
  {:else if deliveries.length === 0}
    <div
      class="text-center py-12 bg-white backdrop-blur-sm rounded-xl border border-gray-200 shadow-sm"
    >
      <p class="text-gray-600">No webhook deliveries found</p>
      <p class="text-sm text-gray-400 mt-2">
        Deliveries will appear here when webhooks are triggered
      </p>
    </div>
  {:else}
    <div class="space-y-3">
      {#each deliveries as delivery (delivery.id)}
        <div
          class="bg-white backdrop-blur-sm rounded-xl border border-gray-200 overflow-hidden shadow-sm"
        >
          <button
            type="button"
            on:click={() => toggleExpanded(delivery.id)}
            class="w-full px-6 py-4 text-left hover:bg-gray-50 transition-colors"
          >
            <div class="flex items-center justify-between">
              <div class="flex-1 min-w-0">
                <div class="flex items-center gap-3 mb-2">
                  <span
                    class="px-2 py-1 text-xs rounded font-medium {getStatusColor(delivery.status)}"
                  >
                    {delivery.status.toUpperCase()}
                  </span>
                  <span class="text-sm font-medium text-gray-900"
                    >{delivery.event_type}</span
                  >
                  {#if delivery.response_status_code}
                    <span class="text-xs text-gray-500"
                      >HTTP {delivery.response_status_code}</span
                    >
                  {/if}
                </div>

                <div class="flex items-center gap-4 text-xs text-gray-500">
                  <span>Attempts: {delivery.attempts}</span>
                  <span>{formatDate(delivery.created_at)}</span>
                  {#if delivery.last_attempt_at}
                    <span>Last attempt: {formatRelativeTime(delivery.last_attempt_at)}</span>
                  {/if}
                  {#if delivery.next_retry_at}
                    <span class="text-yellow-600"
                      >Next retry: {formatDate(delivery.next_retry_at)}</span
                    >
                  {/if}
                </div>
              </div>

              <div class="flex items-center gap-2">
                {#if expandedDeliveryId === delivery.id}
                  <span class="text-primary-600">▼</span>
                {:else}
                  <span class="text-gray-400">▶</span>
                {/if}
              </div>
            </div>
          </button>

          {#if expandedDeliveryId === delivery.id}
            <div
              class="border-t border-gray-200 px-6 py-4 bg-gray-50"
            >
              <div class="space-y-4">
                <!-- Payload -->
                <div>
                  <h4 class="text-sm font-medium text-gray-700 mb-2">Payload</h4>
                  <pre
                    class="bg-gray-100 p-3 rounded-lg border border-gray-200 text-xs overflow-x-auto text-gray-700">{JSON.stringify(
                      delivery.payload,
                      null,
                      2
                    )}</pre>
                </div>

                <!-- Response -->
                {#if delivery.response_body}
                  <div>
                    <h4 class="text-sm font-medium text-gray-700 mb-2">
                      Response Body
                    </h4>
                    <pre
                      class="bg-gray-100 p-3 rounded-lg border border-gray-200 text-xs overflow-x-auto max-h-40 text-gray-700">{delivery.response_body}</pre>
                  </div>
                {/if}

                <!-- Error -->
                {#if delivery.error_message}
                  <div>
                    <h4 class="text-sm font-medium text-red-600 mb-2">
                      Error Message
                    </h4>
                    <div
                      class="bg-red-50 border border-red-200 text-red-600 p-3 rounded-lg text-sm"
                    >
                      {delivery.error_message}
                    </div>
                  </div>
                {/if}

                <!-- Metadata -->
                <div class="grid grid-cols-2 gap-4 text-sm">
                  <div>
                    <span class="font-medium text-gray-500">Delivery ID:</span>
                    <span class="text-primary-600 font-mono text-xs ml-2"
                      >{delivery.id}</span
                    >
                  </div>
                  <div>
                    <span class="font-medium text-gray-500">Endpoint ID:</span>
                    <span class="text-primary-600 font-mono text-xs ml-2"
                      >{delivery.webhook_endpoint_id}</span
                    >
                  </div>
                </div>
              </div>
            </div>
          {/if}
        </div>
      {/each}
    </div>

    <!-- Pagination -->
    <div
      class="flex items-center justify-between bg-white backdrop-blur-sm px-6 py-4 rounded-xl border border-gray-200 shadow-sm"
    >
      <div class="text-sm text-gray-500">
        Showing {offset + 1} - {offset + deliveries.length}
      </div>
      <div class="flex gap-2">
        <Button variant="secondary" size="sm" on:click={handlePrevPage} disabled={offset === 0}>
          Previous
        </Button>
        <Button
          variant="secondary"
          size="sm"
          on:click={handleNextPage}
          disabled={deliveries.length < limit}
        >
          Next
        </Button>
      </div>
    </div>
  {/if}
</div>
