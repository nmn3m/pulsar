<script lang="ts">
  import { onMount } from 'svelte';
  import { authStore } from '$lib/stores/auth';
  import { api } from '$lib/api/client';
  import type { DashboardMetrics, TeamMetrics } from '$lib/types/metrics';

  let metrics: DashboardMetrics | null = null;
  let teamMetrics: TeamMetrics[] = [];
  let loading = true;
  let error: string | null = null;
  let selectedPeriod: 'hourly' | 'daily' | 'weekly' = 'daily';

  async function loadMetrics() {
    loading = true;
    error = null;
    try {
      const [dashboardData, teamsData] = await Promise.all([
        api.getDashboardMetrics({ period: selectedPeriod }),
        api.getTeamMetrics(),
      ]);
      metrics = dashboardData;
      teamMetrics = teamsData.teams || [];
    } catch (e) {
      error = e instanceof Error ? e.message : 'Failed to load metrics';
    } finally {
      loading = false;
    }
  }

  function formatDuration(seconds: number | undefined): string {
    if (!seconds) return '-';
    if (seconds < 60) return `${Math.round(seconds)}s`;
    if (seconds < 3600) return `${Math.round(seconds / 60)}m`;
    if (seconds < 86400) return `${Math.round(seconds / 3600)}h`;
    return `${Math.round(seconds / 86400)}d`;
  }

  function getStatusColor(status: string): string {
    switch (status) {
      case 'open':
        return 'text-red-500';
      case 'acknowledged':
        return 'text-yellow-500';
      case 'investigating':
        return 'text-orange-500';
      case 'identified':
        return 'text-blue-500';
      case 'monitoring':
        return 'text-purple-500';
      case 'resolved':
      case 'closed':
        return 'text-green-500';
      default:
        return 'text-gray-500';
    }
  }

  function getPriorityColor(priority: string): string {
    switch (priority) {
      case 'P1':
        return 'bg-red-100 text-red-800';
      case 'P2':
        return 'bg-orange-100 text-orange-800';
      case 'P3':
        return 'bg-yellow-100 text-yellow-800';
      case 'P4':
        return 'bg-blue-100 text-blue-800';
      case 'P5':
        return 'bg-gray-100 text-gray-800';
      default:
        return 'bg-gray-100 text-gray-800';
    }
  }

  onMount(() => {
    loadMetrics();
  });

  $: if (selectedPeriod) {
    loadMetrics();
  }
</script>

<svelte:head>
  <title>Dashboard - Pulsar</title>
</svelte:head>

<div class="space-y-6">
  <div class="flex justify-between items-center">
    <div>
      <h2 class="text-3xl font-bold text-gray-900">Dashboard</h2>
      <p class="mt-2 text-gray-500">Overview of your incident management metrics</p>
    </div>
    <div class="flex items-center gap-2">
      <select
        bind:value={selectedPeriod}
        class="px-3 py-2 bg-white border border-gray-200 rounded-lg text-sm focus:ring-2 focus:ring-primary-500 focus:border-transparent"
      >
        <option value="hourly">Hourly</option>
        <option value="daily">Daily</option>
        <option value="weekly">Weekly</option>
      </select>
      <button
        on:click={loadMetrics}
        class="p-2 text-gray-500 hover:text-gray-700"
        title="Refresh metrics"
      >
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"
          />
        </svg>
      </button>
    </div>
  </div>

  {#if loading}
    <div class="flex justify-center py-12">
      <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-500"></div>
    </div>
  {:else if error}
    <div class="bg-red-50 border border-red-200 rounded-lg p-4 text-red-700">
      {error}
    </div>
  {:else if metrics}
    <!-- Main Metrics Cards -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
      <!-- Active Alerts Card -->
      <div
        class="bg-white backdrop-blur-sm p-6 rounded-xl border border-error/20 hover:border-error/40 transition-all duration-300 hover:shadow-lg shadow-sm"
      >
        <div class="flex items-center justify-between">
          <h3 class="text-sm font-medium text-gray-500">Open Alerts</h3>
          <span class="text-error">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
              />
            </svg>
          </span>
        </div>
        <p class="mt-3 text-4xl font-bold text-gray-900">
          {metrics.alerts.open}
        </p>
        <p class="mt-2 text-sm text-gray-400">
          {metrics.alerts.total} total ({metrics.alerts.acknowledged} ack)
        </p>
      </div>

      <!-- Open Incidents Card -->
      <div
        class="bg-white backdrop-blur-sm p-6 rounded-xl border border-red-200 hover:border-red-400 transition-all duration-300 hover:shadow-lg shadow-sm"
      >
        <div class="flex items-center justify-between">
          <h3 class="text-sm font-medium text-gray-500">Open Incidents</h3>
          <span class="text-red-500">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M13 10V3L4 14h7v7l9-11h-7z"
              />
            </svg>
          </span>
        </div>
        <p class="mt-3 text-4xl font-bold text-gray-900">
          {metrics.incidents.open +
            metrics.incidents.investigating +
            metrics.incidents.identified +
            metrics.incidents.monitoring}
        </p>
        <p class="mt-2 text-sm text-gray-400">
          {metrics.incidents.total} total ({metrics.incidents.resolved + metrics.incidents.closed} resolved)
        </p>
      </div>

      <!-- Avg Response Time Card -->
      <div
        class="bg-white backdrop-blur-sm p-6 rounded-xl border border-primary-200 hover:border-primary-400 transition-all duration-300 hover:shadow-lg shadow-sm"
      >
        <div class="flex items-center justify-between">
          <h3 class="text-sm font-medium text-gray-500">Avg Response Time</h3>
          <span class="text-primary-500">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
              />
            </svg>
          </span>
        </div>
        <p class="mt-3 text-4xl font-bold text-gray-900">
          {formatDuration(metrics.alerts.avg_response_time_seconds)}
        </p>
        <p class="mt-2 text-sm text-gray-400">Time to acknowledge</p>
      </div>

      <!-- Notifications Card -->
      <div
        class="bg-white backdrop-blur-sm p-6 rounded-xl border border-green-200 hover:border-green-400 transition-all duration-300 hover:shadow-lg shadow-sm"
      >
        <div class="flex items-center justify-between">
          <h3 class="text-sm font-medium text-gray-500">Notifications Sent</h3>
          <span class="text-green-500">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9"
              />
            </svg>
          </span>
        </div>
        <p class="mt-3 text-4xl font-bold text-gray-900">
          {metrics.notifications.sent}
        </p>
        <p class="mt-2 text-sm text-gray-400">
          {metrics.notifications.failed > 0
            ? `${metrics.notifications.failed} failed`
            : 'All delivered'}
        </p>
      </div>
    </div>

    <!-- Alerts by Priority and Source -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- By Priority -->
      <div class="bg-white backdrop-blur-sm p-6 rounded-xl border border-gray-200 shadow-sm">
        <h3 class="text-lg font-semibold text-gray-900 mb-4">Alerts by Priority</h3>
        {#if Object.keys(metrics.alerts.by_priority).length > 0}
          <div class="space-y-3">
            {#each Object.entries(metrics.alerts.by_priority).sort( (a, b) => a[0].localeCompare(b[0]) ) as [priority, count]}
              <div class="flex items-center justify-between">
                <span
                  class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium {getPriorityColor(
                    priority
                  )}"
                >
                  {priority}
                </span>
                <div class="flex items-center gap-2">
                  <div class="w-32 bg-gray-200 rounded-full h-2">
                    <div
                      class="h-2 rounded-full {priority === 'P1'
                        ? 'bg-red-500'
                        : priority === 'P2'
                          ? 'bg-orange-500'
                          : priority === 'P3'
                            ? 'bg-yellow-500'
                            : 'bg-blue-500'}"
                      style="width: {Math.min((count / metrics.alerts.total) * 100, 100)}%"
                    ></div>
                  </div>
                  <span class="text-sm text-gray-600 w-8 text-right">{count}</span>
                </div>
              </div>
            {/each}
          </div>
        {:else}
          <p class="text-gray-500 text-sm">No alerts in this period</p>
        {/if}
      </div>

      <!-- By Source -->
      <div class="bg-white backdrop-blur-sm p-6 rounded-xl border border-gray-200 shadow-sm">
        <h3 class="text-lg font-semibold text-gray-900 mb-4">Alerts by Source</h3>
        {#if Object.keys(metrics.alerts.by_source).length > 0}
          <div class="space-y-3">
            {#each Object.entries(metrics.alerts.by_source)
              .sort((a, b) => b[1] - a[1])
              .slice(0, 5) as [source, count]}
              <div class="flex items-center justify-between">
                <span class="text-sm text-gray-700 truncate max-w-[150px]" title={source}>
                  {source}
                </span>
                <div class="flex items-center gap-2">
                  <div class="w-32 bg-gray-200 rounded-full h-2">
                    <div
                      class="bg-primary-500 h-2 rounded-full"
                      style="width: {Math.min((count / metrics.alerts.total) * 100, 100)}%"
                    ></div>
                  </div>
                  <span class="text-sm text-gray-600 w-8 text-right">{count}</span>
                </div>
              </div>
            {/each}
          </div>
        {:else}
          <p class="text-gray-500 text-sm">No alerts in this period</p>
        {/if}
      </div>
    </div>

    <!-- Incident Status Breakdown -->
    <div class="bg-white backdrop-blur-sm p-6 rounded-xl border border-gray-200 shadow-sm">
      <h3 class="text-lg font-semibold text-gray-900 mb-4">Incident Status</h3>
      <div class="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-7 gap-4">
        {#each [{ label: 'Open', value: metrics.incidents.open, color: 'red' }, { label: 'Investigating', value: metrics.incidents.investigating, color: 'orange' }, { label: 'Identified', value: metrics.incidents.identified, color: 'blue' }, { label: 'Monitoring', value: metrics.incidents.monitoring, color: 'purple' }, { label: 'Resolved', value: metrics.incidents.resolved, color: 'green' }, { label: 'Closed', value: metrics.incidents.closed, color: 'gray' }] as stat}
          <div class="text-center p-3 bg-gray-50 rounded-lg">
            <p class="text-2xl font-bold text-{stat.color}-600">
              {stat.value}
            </p>
            <p class="text-xs text-gray-500 mt-1">{stat.label}</p>
          </div>
        {/each}
      </div>
    </div>

    <!-- Team Performance -->
    {#if teamMetrics.length > 0}
      <div class="bg-white backdrop-blur-sm p-6 rounded-xl border border-gray-200 shadow-sm">
        <h3 class="text-lg font-semibold text-gray-900 mb-4">Team Performance</h3>
        <div class="overflow-x-auto">
          <table class="min-w-full">
            <thead>
              <tr class="border-b border-gray-200">
                <th class="text-left py-3 px-4 text-sm font-medium text-gray-500">Team</th>
                <th class="text-right py-3 px-4 text-sm font-medium text-gray-500">Total</th>
                <th class="text-right py-3 px-4 text-sm font-medium text-gray-500">Acknowledged</th>
                <th class="text-right py-3 px-4 text-sm font-medium text-gray-500">Closed</th>
                <th class="text-right py-3 px-4 text-sm font-medium text-gray-500">Avg Response</th>
              </tr>
            </thead>
            <tbody>
              {#each teamMetrics as team}
                <tr class="border-b border-gray-100 hover:bg-gray-50">
                  <td class="py-3 px-4 text-sm text-gray-900">{team.team_name}</td>
                  <td class="py-3 px-4 text-sm text-gray-600 text-right">{team.total_alerts}</td>
                  <td class="py-3 px-4 text-sm text-gray-600 text-right"
                    >{team.acknowledged_alerts}</td
                  >
                  <td class="py-3 px-4 text-sm text-gray-600 text-right">{team.closed_alerts}</td>
                  <td class="py-3 px-4 text-sm text-gray-600 text-right">
                    {formatDuration(team.avg_response_time_seconds)}
                  </td>
                </tr>
              {/each}
            </tbody>
          </table>
        </div>
      </div>
    {/if}

    <!-- User Info Card -->
    <div class="bg-white backdrop-blur-sm p-6 rounded-xl border border-gray-200 shadow-sm">
      <h3 class="text-lg font-semibold text-gray-900">Your Information</h3>
      <dl class="mt-4 space-y-3">
        <div class="flex justify-between items-center py-2 border-b border-gray-100">
          <dt class="text-sm font-medium text-gray-400">Email:</dt>
          <dd class="text-sm text-primary-600">{$authStore.user?.email}</dd>
        </div>
        <div class="flex justify-between items-center py-2 border-b border-gray-100">
          <dt class="text-sm font-medium text-gray-400">Username:</dt>
          <dd class="text-sm text-gray-700">{$authStore.user?.username}</dd>
        </div>
        <div class="flex justify-between items-center py-2 border-b border-gray-100">
          <dt class="text-sm font-medium text-gray-400">Organization:</dt>
          <dd class="text-sm text-gray-700">{$authStore.organization?.name}</dd>
        </div>
        <div class="flex justify-between items-center py-2">
          <dt class="text-sm font-medium text-gray-400">Plan:</dt>
          <dd class="text-sm text-primary-600 capitalize">
            {$authStore.organization?.plan}
          </dd>
        </div>
      </dl>
    </div>
  {/if}
</div>
