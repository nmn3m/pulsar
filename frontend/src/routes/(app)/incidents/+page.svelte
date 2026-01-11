<script lang="ts">
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { incidentsStore } from '$lib/stores/incidents';
  import Button from '$lib/components/ui/Button.svelte';
  import type {
    CreateIncidentRequest,
    IncidentSeverity,
    IncidentStatus,
  } from '$lib/types/incident';

  const statusOptions: IncidentStatus[] = ['investigating', 'identified', 'monitoring', 'resolved'];
  const severityOptions: IncidentSeverity[] = ['critical', 'high', 'medium', 'low'];

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
      search: searchQuery || undefined,
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
      priority: priority as any,
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

<div class="space-y-6">
  <div class="flex justify-between items-center">
    <div>
      <h2 class="text-3xl font-bold text-gray-900">Incidents</h2>
      <p class="mt-2 text-gray-500">Manage and track incidents across your organization</p>
    </div>
    <Button variant="primary" on:click={() => (showCreateForm = !showCreateForm)}>
      {showCreateForm ? 'Cancel' : 'Create Incident'}
    </Button>
  </div>

  <!-- Create Incident Form -->
  {#if showCreateForm}
    <div class="bg-white backdrop-blur-sm p-6 rounded-xl border border-primary-200 shadow-sm">
      <h3 class="text-lg font-semibold mb-4 text-gray-900">Create New Incident</h3>

      <form on:submit|preventDefault={handleCreate} class="space-y-4">
        <div>
          <label for="title" class="block text-sm font-medium text-gray-600 mb-1">Title *</label>
          <input
            id="title"
            type="text"
            bind:value={title}
            class="w-full px-3 py-2 bg-white border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500 text-gray-900 placeholder-gray-400"
            placeholder="Brief description of the incident"
          />
        </div>

        <div>
          <label for="description" class="block text-sm font-medium text-gray-600 mb-1"
            >Description</label
          >
          <textarea
            id="description"
            bind:value={description}
            rows="3"
            class="w-full px-3 py-2 bg-white border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500 text-gray-900 placeholder-gray-400"
            placeholder="Detailed description of the incident"
          ></textarea>
        </div>

        <div class="grid grid-cols-2 gap-4">
          <div>
            <label for="severity" class="block text-sm font-medium text-gray-600 mb-1"
              >Severity *</label
            >
            <select
              id="severity"
              bind:value={severity}
              class="w-full px-3 py-2 bg-white border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500 text-gray-900"
            >
              <option value="critical">Critical</option>
              <option value="high">High</option>
              <option value="medium">Medium</option>
              <option value="low">Low</option>
            </select>
          </div>

          <div>
            <label for="priority" class="block text-sm font-medium text-gray-600 mb-1"
              >Priority *</label
            >
            <select
              id="priority"
              bind:value={priority}
              class="w-full px-3 py-2 bg-white border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500 text-gray-900"
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
          <Button type="submit" variant="primary" disabled={!title.trim() || isLoading}>
            Create Incident
          </Button>
          <Button type="button" variant="secondary" on:click={() => (showCreateForm = false)}
            >Cancel</Button
          >
        </div>
      </form>
    </div>
  {/if}

  <!-- Filters -->
  <div class="bg-white backdrop-blur-sm p-4 rounded-xl border border-gray-200 shadow-sm">
    <div class="space-y-4">
      <!-- Status Filter -->
      <div>
        <label class="block text-sm font-medium text-gray-600 mb-2">Status</label>
        <div class="flex flex-wrap gap-2">
          {#each statusOptions as status}
            <button
              type="button"
              on:click={() => toggleStatusFilter(status)}
              class="px-3 py-1.5 rounded-lg text-sm font-medium transition-all duration-200 {selectedStatuses.includes(
                status
              )
                ? getStatusColor(status)
                : 'bg-gray-100 text-gray-600 hover:bg-gray-200 border border-gray-300 hover:border-primary-500/50'}"
            >
              {status.charAt(0).toUpperCase() + status.slice(1)}
            </button>
          {/each}
        </div>
      </div>

      <!-- Severity Filter -->
      <div>
        <label class="block text-sm font-medium text-gray-600 mb-2">Severity</label>
        <div class="flex flex-wrap gap-2">
          {#each severityOptions as sev}
            <button
              type="button"
              on:click={() => toggleSeverityFilter(sev)}
              class="px-3 py-1.5 rounded-lg text-sm font-medium transition-all duration-200 {selectedSeverities.includes(
                sev
              )
                ? getSeverityColor(sev)
                : 'bg-gray-100 text-gray-600 hover:bg-gray-200 border border-gray-300 hover:border-error/50'}"
            >
              {sev.charAt(0).toUpperCase() + sev.slice(1)}
            </button>
          {/each}
        </div>
      </div>

      <!-- Search -->
      <div>
        <label for="search" class="block text-sm font-medium text-gray-600 mb-2">Search</label>
        <div class="flex gap-2">
          <input
            id="search"
            type="text"
            bind:value={searchQuery}
            placeholder="Search incidents..."
            class="flex-1 px-3 py-2 bg-white border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500 text-gray-900 placeholder-gray-400"
          />
          <Button variant="primary" on:click={applyFilters}>Search</Button>
        </div>
      </div>
    </div>
  </div>

  <!-- Error Display -->
  {#if error}
    <div class="bg-red-50 border border-red-200 text-red-600 px-4 py-3 rounded-lg">
      {error}
    </div>
  {/if}

  <!-- Incidents List -->
  <div class="space-y-4">
    {#if isLoading}
      <div class="text-center py-12">
        <div
          class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-primary-500"
        ></div>
        <p class="mt-2 text-gray-500">Loading incidents...</p>
      </div>
    {:else if incidents.length === 0}
      <div
        class="text-center py-12 bg-white backdrop-blur-sm rounded-xl border border-gray-200 shadow-sm"
      >
        <p class="text-gray-600">No incidents found</p>
        <p class="text-sm text-gray-400 mt-2">Create your first incident to get started</p>
      </div>
    {:else}
      <div
        class="bg-white backdrop-blur-sm rounded-xl border border-gray-200 overflow-hidden shadow-sm"
      >
        <div class="px-6 py-3 bg-gray-50 border-b border-gray-200">
          <p class="text-sm text-gray-600">
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
                    <p class="text-sm font-medium text-gray-900">
                      {incident.title}
                    </p>
                    {#if incident.description}
                      <p class="mt-1 text-sm text-gray-600 line-clamp-2">
                        {incident.description}
                      </p>
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
</div>
