<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { goto } from '$app/navigation';
	import { alertsStore } from '$lib/stores/alerts';
	import { wsStore } from '$lib/stores/websocket';
	import AlertCard from '$lib/components/alerts/AlertCard.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import type { AlertStatus, AlertPriority } from '$lib/types/alert';

	let showCreateForm = false;
	let selectedStatus: AlertStatus[] = [];
	let selectedPriority: AlertPriority[] = [];
	let searchQuery = '';

	// Create alert form
	let message = '';
	let description = '';
	let priority: AlertPriority = 'P3';
	let tags = '';
	let createError = '';
	let creatingAlert = false;

	let unsubscribeWS: (() => void)[] = [];

	onMount(() => {
		loadAlerts();

		// Listen for WebSocket alert events
		unsubscribeWS.push(
			wsStore.on('alert.created', () => {
				loadAlerts(); // Refresh alerts list
			})
		);
		unsubscribeWS.push(
			wsStore.on('alert.updated', () => {
				loadAlerts();
			})
		);
		unsubscribeWS.push(
			wsStore.on('alert.acknowledged', () => {
				loadAlerts();
			})
		);
		unsubscribeWS.push(
			wsStore.on('alert.closed', () => {
				loadAlerts();
			})
		);
		unsubscribeWS.push(
			wsStore.on('alert.deleted', () => {
				loadAlerts();
			})
		);
	});

	onDestroy(() => {
		// Unsubscribe from WebSocket events
		unsubscribeWS.forEach((unsub) => unsub());
	});

	function loadAlerts() {
		const params = {
			status: selectedStatus.length > 0 ? selectedStatus : undefined,
			priority: selectedPriority.length > 0 ? selectedPriority : undefined,
			search: searchQuery || undefined
		};
		alertsStore.load(params);
	}

	async function handleCreateAlert() {
		createError = '';
		creatingAlert = true;

		try {
			await alertsStore.create({
				source: 'manual',
				priority,
				message,
				description: description || undefined,
				tags: tags ? tags.split(',').map(t => t.trim()) : []
			});

			// Reset form
			message = '';
			description = '';
			priority = 'P3';
			tags = '';
			showCreateForm = false;
		} catch (err) {
			createError = err instanceof Error ? err.message : 'Failed to create alert';
		} finally {
			creatingAlert = false;
		}
	}

	async function handleAcknowledge(id: string) {
		try {
			await alertsStore.acknowledge(id);
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to acknowledge alert');
		}
	}

	async function handleClose(id: string) {
		const reason = prompt('Reason for closing (optional):');
		if (reason === null) return;

		try {
			await alertsStore.close(id, { reason: reason || 'Closed manually' });
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to close alert');
		}
	}

	function handleStatusFilter(status: AlertStatus) {
		if (selectedStatus.includes(status)) {
			selectedStatus = selectedStatus.filter(s => s !== status);
		} else {
			selectedStatus = [...selectedStatus, status];
		}
		loadAlerts();
	}

	function handlePriorityFilter(prio: AlertPriority) {
		if (selectedPriority.includes(prio)) {
			selectedPriority = selectedPriority.filter(p => p !== prio);
		} else {
			selectedPriority = [...selectedPriority, prio];
		}
		loadAlerts();
	}
</script>

<svelte:head>
	<title>Alerts - Pulsar</title>
</svelte:head>

<div class="space-y-6">
	<div class="flex justify-between items-center">
		<div>
			<h2 class="text-3xl font-bold text-gray-100">Alerts</h2>
			<p class="mt-2 text-gray-400">Manage and respond to alerts</p>
		</div>
		<Button variant="primary" on:click={() => (showCreateForm = !showCreateForm)}>
			{showCreateForm ? 'Cancel' : 'Create Alert'}
		</Button>
	</div>

	{#if showCreateForm}
		<div class="bg-space-800/50 backdrop-blur-sm p-6 rounded-xl border border-primary-500/30">
			<h3 class="text-lg font-semibold mb-4 text-gray-100">Create New Alert</h3>
			<form on:submit|preventDefault={handleCreateAlert} class="space-y-4">
				<div>
					<label for="message" class="block text-sm font-medium text-gray-300 mb-1">
						Message *
					</label>
					<input
						id="message"
						type="text"
						bind:value={message}
						required
						class="w-full px-3 py-2 bg-space-700 border border-space-500 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500 text-gray-100 placeholder-gray-500"
						placeholder="Brief description of the alert"
					/>
				</div>

				<div>
					<label for="description" class="block text-sm font-medium text-gray-300 mb-1">
						Description
					</label>
					<textarea
						id="description"
						bind:value={description}
						rows="3"
						class="w-full px-3 py-2 bg-space-700 border border-space-500 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500 text-gray-100 placeholder-gray-500"
						placeholder="Additional details..."
					></textarea>
				</div>

				<div>
					<label for="priority" class="block text-sm font-medium text-gray-300 mb-1">
						Priority *
					</label>
					<select
						id="priority"
						bind:value={priority}
						class="w-full px-3 py-2 bg-space-700 border border-space-500 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500 text-gray-100"
					>
						<option value="P1">P1 - Critical</option>
						<option value="P2">P2 - High</option>
						<option value="P3">P3 - Medium</option>
						<option value="P4">P4 - Low</option>
						<option value="P5">P5 - Informational</option>
					</select>
				</div>

				<div>
					<label for="tags" class="block text-sm font-medium text-gray-300 mb-1">
						Tags (comma-separated)
					</label>
					<input
						id="tags"
						type="text"
						bind:value={tags}
						class="w-full px-3 py-2 bg-space-700 border border-space-500 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500 text-gray-100 placeholder-gray-500"
						placeholder="production, database, critical"
					/>
				</div>

				{#if createError}
					<div class="bg-accent-900/30 border border-accent-500/50 text-accent-300 px-4 py-3 rounded-lg">
						{createError}
					</div>
				{/if}

				<div class="flex gap-2">
					<Button type="submit" variant="primary" disabled={creatingAlert}>
						{creatingAlert ? 'Creating...' : 'Create Alert'}
					</Button>
					<Button type="button" variant="secondary" on:click={() => (showCreateForm = false)}>
						Cancel
					</Button>
				</div>
			</form>
		</div>
	{/if}

	<!-- Filters -->
	<div class="bg-space-800/50 backdrop-blur-sm p-4 rounded-xl border border-space-600">
		<div class="space-y-4">
			<div>
				<label class="block text-sm font-medium text-gray-300 mb-2">Status</label>
				<div class="flex gap-2 flex-wrap">
					{#each ['open', 'acknowledged', 'closed', 'snoozed'] as status}
						<button
							type="button"
							class="px-3 py-1.5 rounded-lg text-sm font-medium transition-all duration-200 {selectedStatus.includes(status)
								? 'bg-primary-500 text-space-900 shadow-neon-cyan'
								: 'bg-space-700 text-gray-300 hover:bg-space-600 border border-space-500 hover:border-primary-500/50'}"
							on:click={() => handleStatusFilter(status)}
						>
							{status}
						</button>
					{/each}
				</div>
			</div>

			<div>
				<label class="block text-sm font-medium text-gray-300 mb-2">Priority</label>
				<div class="flex gap-2 flex-wrap">
					{#each ['P1', 'P2', 'P3', 'P4', 'P5'] as prio}
						<button
							type="button"
							class="px-3 py-1.5 rounded-lg text-sm font-medium transition-all duration-200 {selectedPriority.includes(prio)
								? 'bg-accent-500 text-white shadow-neon-pink'
								: 'bg-space-700 text-gray-300 hover:bg-space-600 border border-space-500 hover:border-accent-500/50'}"
							on:click={() => handlePriorityFilter(prio)}
						>
							{prio}
						</button>
					{/each}
				</div>
			</div>

			<div>
				<label for="search" class="block text-sm font-medium text-gray-300 mb-2">Search</label>
				<div class="flex gap-2">
					<input
						id="search"
						type="text"
						bind:value={searchQuery}
						class="flex-1 px-3 py-2 bg-space-700 border border-space-500 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500 text-gray-100 placeholder-gray-500"
						placeholder="Search alerts..."
					/>
					<Button variant="primary" on:click={loadAlerts}>Search</Button>
				</div>
			</div>
		</div>
	</div>

	<!-- Alerts List -->
	<div class="space-y-4">
		{#if $alertsStore.isLoading}
			<div class="text-center py-12">
				<div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-primary-500"></div>
				<p class="mt-2 text-gray-400">Loading alerts...</p>
			</div>
		{:else if $alertsStore.error}
			<div class="bg-accent-900/30 border border-accent-500/50 text-accent-300 px-4 py-3 rounded-lg">
				{$alertsStore.error}
			</div>
		{:else if $alertsStore.alerts.length === 0}
			<div class="text-center py-12 bg-space-800/50 backdrop-blur-sm rounded-xl border border-space-600">
				<p class="text-gray-300">No alerts found</p>
				<p class="text-sm text-gray-500 mt-2">Create your first alert to get started</p>
			</div>
		{:else}
			{#each $alertsStore.alerts as alert (alert.id)}
				<AlertCard
					{alert}
					onAcknowledge={handleAcknowledge}
					onClose={handleClose}
					onClick={(id) => goto(`/alerts/${id}`)}
				/>
			{/each}

			<div class="text-sm text-gray-500 text-center py-4">
				Showing {$alertsStore.alerts.length} of {$alertsStore.total} alerts
			</div>
		{/if}
	</div>
</div>
