<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { incidentsStore } from '$lib/stores/incidents';
	import Button from '$lib/components/ui/Button.svelte';
	import type { CreateIncidentRequest, IncidentSeverity, IncidentStatus } from '$lib/types/incident';

	let showCreateForm = false;
	let title = '';
	let description = '';
	let severity: IncidentSeverity = 'medium';
	let priority = 'P3';

	// Filters
	let selectedStatuses: IncidentStatus[] = [];
	let selectedSeverities: IncidentSeverity[] = [];
	let searchQuery = '';

	$: ({ incidents, isLoading, error, total } = $incidentsStore);

	onMount(() => {
		loadIncidents();
	});

	function loadIncidents() {
		incidentsStore.load({
			status: selectedStatuses.length > 0 ? selectedStatuses : undefined,
			severity: selectedSeverities.length > 0 ? selectedSeverities : undefined,
			search: searchQuery || undefined
		});
	}

	function applyFilters() {
		loadIncidents();
	}

	async function handleCreate() {
		if (!title.trim()) return;

		const data: CreateIncidentRequest = {
			title: title.trim(),
			description: description.trim() || undefined,
			severity,
			priority: priority as any
		};

		try {
			await incidentsStore.create(data);
			title = '';
			description = '';
			severity = 'medium';
			priority = 'P3';
			showCreateForm = false;
		} catch (err) {
			console.error('Failed to create incident:', err);
		}
	}

	function viewIncident(id: string) {
		goto(`/incidents/${id}`);
	}

	function getSeverityColor(sev: IncidentSeverity): string {
		switch (sev) {
			case 'critical':
				return 'bg-red-100 text-red-800';
			case 'high':
				return 'bg-orange-100 text-orange-800';
			case 'medium':
				return 'bg-yellow-100 text-yellow-800';
			case 'low':
				return 'bg-blue-100 text-blue-800';
			default:
				return 'bg-gray-100 text-gray-800';
		}
	}

	function getStatusColor(status: IncidentStatus): string {
		switch (status) {
			case 'investigating':
				return 'bg-yellow-100 text-yellow-800';
			case 'identified':
				return 'bg-blue-100 text-blue-800';
			case 'monitoring':
				return 'bg-purple-100 text-purple-800';
			case 'resolved':
				return 'bg-green-100 text-green-800';
			default:
				return 'bg-gray-100 text-gray-800';
		}
	}

	function formatDate(dateStr: string): string {
		const date = new Date(dateStr);
		const now = new Date();
		const diff = now.getTime() - date.getTime();
		const seconds = Math.floor(diff / 1000);
		const minutes = Math.floor(seconds / 60);
		const hours = Math.floor(minutes / 60);
		const days = Math.floor(hours / 24);

		if (days > 0) return `${days} day${days > 1 ? 's' : ''} ago`;
		if (hours > 0) return `${hours} hour${hours > 1 ? 's' : ''} ago`;
		if (minutes > 0) return `${minutes} minute${minutes > 1 ? 's' : ''} ago`;
		return 'Just now';
	}

	function toggleStatusFilter(status: IncidentStatus) {
		if (selectedStatuses.includes(status)) {
			selectedStatuses = selectedStatuses.filter((s) => s !== status);
		} else {
			selectedStatuses = [...selectedStatuses, status];
		}
	}

	function toggleSeverityFilter(sev: IncidentSeverity) {
		if (selectedSeverities.includes(sev)) {
			selectedSeverities = selectedSeverities.filter((s) => s !== sev);
		} else {
			selectedSeverities = [...selectedSeverities, sev];
		}
	}
</script>

<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
	<div class="mb-8">
		<h1 class="text-3xl font-bold text-gray-900">Incidents</h1>
		<p class="mt-2 text-sm text-gray-600">Manage and track incidents across your organization</p>
	</div>

	<!-- Create Incident Button -->
	<div class="mb-6">
		<Button on:click={() => (showCreateForm = !showCreateForm)}>
			{showCreateForm ? 'Cancel' : 'Create Incident'}
		</Button>
	</div>

	<!-- Create Incident Form -->
	{#if showCreateForm}
		<div class="bg-white shadow rounded-lg p-6 mb-6">
			<h2 class="text-xl font-semibold mb-4">Create New Incident</h2>

			<div class="space-y-4">
				<div>
					<label for="title" class="block text-sm font-medium text-gray-700">Title *</label>
					<input
						id="title"
						type="text"
						bind:value={title}
						class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm px-3 py-2 border"
						placeholder="Brief description of the incident"
					/>
				</div>

				<div>
					<label for="description" class="block text-sm font-medium text-gray-700"
						>Description</label
					>
					<textarea
						id="description"
						bind:value={description}
						rows="3"
						class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm px-3 py-2 border"
						placeholder="Detailed description of the incident"
					/>
				</div>

				<div class="grid grid-cols-2 gap-4">
					<div>
						<label for="severity" class="block text-sm font-medium text-gray-700">Severity *</label>
						<select
							id="severity"
							bind:value={severity}
							class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm px-3 py-2 border"
						>
							<option value="critical">Critical</option>
							<option value="high">High</option>
							<option value="medium">Medium</option>
							<option value="low">Low</option>
						</select>
					</div>

					<div>
						<label for="priority" class="block text-sm font-medium text-gray-700">Priority *</label>
						<select
							id="priority"
							bind:value={priority}
							class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm px-3 py-2 border"
						>
							<option value="P1">P1 - Critical</option>
							<option value="P2">P2 - High</option>
							<option value="P3">P3 - Medium</option>
							<option value="P4">P4 - Low</option>
							<option value="P5">P5 - Info</option>
						</select>
					</div>
				</div>

				<div class="flex gap-2">
					<Button on:click={handleCreate} disabled={!title.trim() || isLoading}>
						Create Incident
					</Button>
					<Button variant="secondary" on:click={() => (showCreateForm = false)}>Cancel</Button>
				</div>
			</div>
		</div>
	{/if}

	<!-- Filters -->
	<div class="bg-white shadow rounded-lg p-6 mb-6">
		<h3 class="text-lg font-medium mb-4">Filters</h3>

		<div class="space-y-4">
			<!-- Status Filter -->
			<div>
				<label class="block text-sm font-medium text-gray-700 mb-2">Status</label>
				<div class="flex flex-wrap gap-2">
					{#each ['investigating', 'identified', 'monitoring', 'resolved'] as status}
						<button
							type="button"
							on:click={() => toggleStatusFilter(status)}
							class="px-3 py-1 rounded-full text-sm font-medium transition-colors {selectedStatuses.includes(
								status
							)
								? getStatusColor(status)
								: 'bg-gray-100 text-gray-700 hover:bg-gray-200'}"
						>
							{status.charAt(0).toUpperCase() + status.slice(1)}
						</button>
					{/each}
				</div>
			</div>

			<!-- Severity Filter -->
			<div>
				<label class="block text-sm font-medium text-gray-700 mb-2">Severity</label>
				<div class="flex flex-wrap gap-2">
					{#each ['critical', 'high', 'medium', 'low'] as sev}
						<button
							type="button"
							on:click={() => toggleSeverityFilter(sev)}
							class="px-3 py-1 rounded-full text-sm font-medium transition-colors {selectedSeverities.includes(
								sev
							)
								? getSeverityColor(sev)
								: 'bg-gray-100 text-gray-700 hover:bg-gray-200'}"
						>
							{sev.charAt(0).toUpperCase() + sev.slice(1)}
						</button>
					{/each}
				</div>
			</div>

			<!-- Search -->
			<div>
				<label for="search" class="block text-sm font-medium text-gray-700 mb-2">Search</label>
				<input
					id="search"
					type="text"
					bind:value={searchQuery}
					placeholder="Search incidents..."
					class="w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm px-3 py-2 border"
				/>
			</div>

			<Button on:click={applyFilters}>Apply Filters</Button>
		</div>
	</div>

	<!-- Error Display -->
	{#if error}
		<div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded mb-4">
			{error}
		</div>
	{/if}

	<!-- Loading State -->
	{#if isLoading}
		<div class="text-center py-12">
			<p class="text-gray-500">Loading incidents...</p>
		</div>
	{:else if incidents.length === 0}
		<!-- Empty State -->
		<div class="text-center py-12">
			<p class="text-gray-500">No incidents found. Create your first incident to get started.</p>
		</div>
	{:else}
		<!-- Incidents List -->
		<div class="bg-white shadow rounded-lg overflow-hidden">
			<div class="px-6 py-3 bg-gray-50 border-b border-gray-200">
				<p class="text-sm text-gray-700">
					Showing {incidents.length} of {total} incident{total !== 1 ? 's' : ''}
				</p>
			</div>

			<ul class="divide-y divide-gray-200">
				{#each incidents as incident (incident.id)}
					<li class="hover:bg-gray-50 transition-colors">
						<button
							type="button"
							on:click={() => viewIncident(incident.id)}
							class="w-full text-left px-6 py-4"
						>
							<div class="flex items-start justify-between">
								<div class="flex-1 min-w-0">
									<div class="flex items-center gap-2 mb-2">
										<span
											class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium {getSeverityColor(
												incident.severity
											)}"
										>
											{incident.severity.toUpperCase()}
										</span>
										<span
											class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium {getStatusColor(
												incident.status
											)}"
										>
											{incident.status.charAt(0).toUpperCase() + incident.status.slice(1)}
										</span>
										<span
											class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-gray-100 text-gray-800"
										>
											{incident.priority}
										</span>
									</div>
									<p class="text-sm font-medium text-gray-900">{incident.title}</p>
									{#if incident.description}
										<p class="mt-1 text-sm text-gray-600 line-clamp-2">{incident.description}</p>
									{/if}
									<p class="mt-2 text-xs text-gray-500">
										Started {formatDate(incident.started_at)}
										{#if incident.resolved_at}
											• Resolved {formatDate(incident.resolved_at)}
										{/if}
									</p>
								</div>
								<div class="ml-4 flex-shrink-0">
									<svg
										class="h-5 w-5 text-gray-400"
										xmlns="http://www.w3.org/2000/svg"
										viewBox="0 0 20 20"
										fill="currentColor"
									>
										<path
											fill-rule="evenodd"
											d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z"
											clip-rule="evenodd"
										/>
									</svg>
								</div>
							</div>
						</button>
					</li>
				{/each}
			</ul>
		</div>
	{/if}
</div>
