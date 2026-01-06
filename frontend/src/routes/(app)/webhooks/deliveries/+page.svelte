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
				return 'bg-green-100 dark:bg-green-900/50 text-green-700 dark:text-green-300 border border-green-300 dark:border-green-500/30';
			case 'failed':
				return 'bg-red-100 dark:bg-accent-900/50 text-red-700 dark:text-accent-300 border border-red-300 dark:border-accent-500/30';
			case 'pending':
				return 'bg-yellow-100 dark:bg-yellow-900/50 text-yellow-700 dark:text-yellow-300 border border-yellow-300 dark:border-yellow-500/30';
			default:
				return 'bg-gray-100 dark:bg-space-700 text-gray-700 dark:text-gray-300 border border-gray-300 dark:border-space-500';
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
			<h2 class="text-3xl font-bold text-gray-900 dark:text-gray-100">Webhook Deliveries</h2>
			<p class="mt-2 text-gray-500 dark:text-gray-400">Monitor webhook delivery status and troubleshoot issues</p>
		</div>
		<Button variant="secondary" on:click={loadDeliveries}>
			Refresh
		</Button>
	</div>

	{#if error}
		<div class="bg-red-50 dark:bg-accent-900/30 border border-red-200 dark:border-accent-500/50 text-red-600 dark:text-accent-300 px-4 py-3 rounded-lg">
			{error}
		</div>
	{/if}

	{#if isLoading}
		<div class="text-center py-12">
			<div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-primary-500"></div>
			<p class="mt-2 text-gray-500 dark:text-gray-400">Loading webhook deliveries...</p>
		</div>
	{:else if deliveries.length === 0}
		<div class="text-center py-12 bg-white dark:bg-space-800/50 backdrop-blur-sm rounded-xl border border-gray-200 dark:border-space-600 shadow-sm">
			<p class="text-gray-600 dark:text-gray-300">No webhook deliveries found</p>
			<p class="text-sm text-gray-400 dark:text-gray-500 mt-2">Deliveries will appear here when webhooks are triggered</p>
		</div>
	{:else}
		<div class="space-y-3">
			{#each deliveries as delivery (delivery.id)}
				<div class="bg-white dark:bg-space-800/50 backdrop-blur-sm rounded-xl border border-gray-200 dark:border-space-600 overflow-hidden shadow-sm">
					<button
						type="button"
						on:click={() => toggleExpanded(delivery.id)}
						class="w-full px-6 py-4 text-left hover:bg-gray-50 dark:hover:bg-space-700/50 transition-colors"
					>
						<div class="flex items-center justify-between">
							<div class="flex-1 min-w-0">
								<div class="flex items-center gap-3 mb-2">
									<span class="px-2 py-1 text-xs rounded font-medium {getStatusColor(delivery.status)}">
										{delivery.status.toUpperCase()}
									</span>
									<span class="text-sm font-medium text-gray-900 dark:text-gray-100">{delivery.event_type}</span>
									{#if delivery.response_status_code}
										<span class="text-xs text-gray-500 dark:text-gray-400">HTTP {delivery.response_status_code}</span>
									{/if}
								</div>

								<div class="flex items-center gap-4 text-xs text-gray-500 dark:text-gray-400">
									<span>Attempts: {delivery.attempts}</span>
									<span>{formatDate(delivery.created_at)}</span>
									{#if delivery.last_attempt_at}
										<span>Last attempt: {formatRelativeTime(delivery.last_attempt_at)}</span>
									{/if}
									{#if delivery.next_retry_at}
										<span class="text-yellow-600 dark:text-yellow-400">Next retry: {formatDate(delivery.next_retry_at)}</span>
									{/if}
								</div>
							</div>

							<div class="flex items-center gap-2">
								{#if expandedDeliveryId === delivery.id}
									<span class="text-primary-600 dark:text-primary-400">▼</span>
								{:else}
									<span class="text-gray-400 dark:text-gray-500">▶</span>
								{/if}
							</div>
						</div>
					</button>

					{#if expandedDeliveryId === delivery.id}
						<div class="border-t border-gray-200 dark:border-space-600 px-6 py-4 bg-gray-50 dark:bg-space-900/50">
							<div class="space-y-4">
								<!-- Payload -->
								<div>
									<h4 class="text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Payload</h4>
									<pre class="bg-gray-100 dark:bg-space-700 p-3 rounded-lg border border-gray-200 dark:border-space-500 text-xs overflow-x-auto text-gray-700 dark:text-gray-300">{JSON.stringify(delivery.payload, null, 2)}</pre>
								</div>

								<!-- Response -->
								{#if delivery.response_body}
									<div>
										<h4 class="text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Response Body</h4>
										<pre class="bg-gray-100 dark:bg-space-700 p-3 rounded-lg border border-gray-200 dark:border-space-500 text-xs overflow-x-auto max-h-40 text-gray-700 dark:text-gray-300">{delivery.response_body}</pre>
									</div>
								{/if}

								<!-- Error -->
								{#if delivery.error_message}
									<div>
										<h4 class="text-sm font-medium text-red-600 dark:text-accent-400 mb-2">Error Message</h4>
										<div class="bg-red-50 dark:bg-accent-900/30 border border-red-200 dark:border-accent-500/50 text-red-600 dark:text-accent-300 p-3 rounded-lg text-sm">
											{delivery.error_message}
										</div>
									</div>
								{/if}

								<!-- Metadata -->
								<div class="grid grid-cols-2 gap-4 text-sm">
									<div>
										<span class="font-medium text-gray-500 dark:text-gray-400">Delivery ID:</span>
										<span class="text-primary-600 dark:text-primary-300 font-mono text-xs ml-2">{delivery.id}</span>
									</div>
									<div>
										<span class="font-medium text-gray-500 dark:text-gray-400">Endpoint ID:</span>
										<span class="text-primary-600 dark:text-primary-300 font-mono text-xs ml-2">{delivery.webhook_endpoint_id}</span>
									</div>
								</div>
							</div>
						</div>
					{/if}
				</div>
			{/each}
		</div>

		<!-- Pagination -->
		<div class="flex items-center justify-between bg-white dark:bg-space-800/50 backdrop-blur-sm px-6 py-4 rounded-xl border border-gray-200 dark:border-space-600 shadow-sm">
			<div class="text-sm text-gray-500 dark:text-gray-400">
				Showing {offset + 1} - {offset + deliveries.length}
			</div>
			<div class="flex gap-2">
				<Button variant="secondary" size="sm" on:click={handlePrevPage} disabled={offset === 0}>
					Previous
				</Button>
				<Button variant="secondary" size="sm" on:click={handleNextPage} disabled={deliveries.length < limit}>
					Next
				</Button>
			</div>
		</div>
	{/if}
</div>
