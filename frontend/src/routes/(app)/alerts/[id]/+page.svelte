<script lang="ts">
  import { onMount } from 'svelte';
  import { page } from '$app/stores';
  import { goto } from '$app/navigation';
  import { api } from '$lib/api/client';
  import type { Alert } from '$lib/types/alert';
  import type { User } from '$lib/types/user';
  import type { Team } from '$lib/types/team';
  import Button from '$lib/components/ui/Button.svelte';
  import dayjs from 'dayjs';
  import relativeTime from 'dayjs/plugin/relativeTime';

  dayjs.extend(relativeTime);

  let alertId = $page.params.id!;
  let alertData: Alert | null = null;
  let users: User[] = [];
  let teams: Team[] = [];
  let isLoading = true;
  let error = '';

  // Assignment form
  let showAssignForm = false;
  let assignmentType: 'user' | 'team' = 'user';
  let selectedUserId = '';
  let selectedTeamId = '';
  let assignError = '';
  let isAssigning = false;

  // Snooze form
  let showSnoozeForm = false;
  let snoozeMinutes = 30;
  let isSnoozing = false;

  const priorityColors = {
    P1: 'bg-red-100 text-red-800 border-red-200',
    P2: 'bg-orange-100 text-orange-800 border-orange-200',
    P3: 'bg-yellow-100 text-yellow-800 border-yellow-200',
    P4: 'bg-blue-100 text-blue-800 border-blue-200',
    P5: 'bg-gray-100 text-gray-800 border-gray-200',
  };

  const statusColors = {
    open: 'bg-red-100 text-red-800',
    acknowledged: 'bg-yellow-100 text-yellow-800',
    closed: 'bg-green-100 text-green-800',
    snoozed: 'bg-purple-100 text-purple-800',
  };

  onMount(async () => {
    await loadAlert();
    await loadUsersAndTeams();
  });

  async function loadAlert() {
    try {
      isLoading = true;
      error = '';
      alertData = await api.getAlert(alertId);
    } catch (err) {
      error = err instanceof Error ? err.message : 'Failed to load alert';
    } finally {
      isLoading = false;
    }
  }

  async function loadUsersAndTeams() {
    try {
      const [usersResp, teamsResp] = await Promise.all([api.listUsers(), api.listTeams()]);
      users = usersResp.users;
      teams = teamsResp.teams;
    } catch (err) {
      console.error('Failed to load users/teams:', err);
    }
  }

  async function handleAcknowledge() {
    if (!alertData) return;

    try {
      await api.acknowledgeAlert(alertData.id);
      await loadAlert();
    } catch (err) {
      alert(err instanceof Error ? err.message : 'Failed to acknowledge alert');
    }
  }

  async function handleClose() {
    if (!alertData) return;

    const reason = prompt('Reason for closing (optional):');
    if (reason === null) return;

    try {
      await api.closeAlert(alertData.id, { reason: reason || 'Closed manually' });
      await loadAlert();
    } catch (err) {
      alert(err instanceof Error ? err.message : 'Failed to close alert');
    }
  }

  async function handleSnooze() {
    if (!alertData) return;

    isSnoozing = true;
    try {
      const until = new Date();
      until.setMinutes(until.getMinutes() + snoozeMinutes);

      await api.snoozeAlert(alertData.id, { until: until.toISOString() });
      await loadAlert();
      showSnoozeForm = false;
    } catch (err) {
      alert(err instanceof Error ? err.message : 'Failed to snooze alert');
    } finally {
      isSnoozing = false;
    }
  }

  async function handleAssign() {
    if (!alertData) return;

    assignError = '';
    isAssigning = true;

    try {
      await api.assignAlert(alertData.id, {
        user_id: assignmentType === 'user' ? selectedUserId : undefined,
        team_id: assignmentType === 'team' ? selectedTeamId : undefined,
      });

      await loadAlert();
      showAssignForm = false;
      selectedUserId = '';
      selectedTeamId = '';
    } catch (err) {
      assignError = err instanceof Error ? err.message : 'Failed to assign alert';
    } finally {
      isAssigning = false;
    }
  }

  async function handleDelete() {
    if (!alertData) return;

    if (!confirm('Are you sure you want to delete this alert?')) {
      return;
    }

    try {
      await api.deleteAlert(alertData.id);
      goto('/alerts');
    } catch (err) {
      alert(err instanceof Error ? err.message : 'Failed to delete alert');
    }
  }

  function getAssignedToDisplay(): string {
    if (!alertData) return 'Unassigned';
    const currentAlert = alertData;
    if (currentAlert.assigned_to_user_id) {
      const user = users.find((u) => u.id === currentAlert.assigned_to_user_id);
      return user ? `👤 ${user.full_name || user.username}` : 'Assigned to user';
    }
    if (currentAlert.assigned_to_team_id) {
      const team = teams.find((t) => t.id === currentAlert.assigned_to_team_id);
      return team ? `👥 ${team.name}` : 'Assigned to team';
    }
    return 'Unassigned';
  }
</script>

<svelte:head>
  <title>{alertData?.message || 'Alert'} - Pulsar</title>
</svelte:head>

<div class="space-y-6">
  <!-- Header -->
  <div class="flex items-center gap-3 mb-6">
    <button
      on:click={() => goto('/alerts')}
      class="text-gray-600 hover:text-gray-900"
      aria-label="Back to alerts"
    >
      <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
      </svg>
    </button>
    <h2 class="text-3xl font-bold text-gray-900">Alert Details</h2>
  </div>

  {#if isLoading}
    <div class="text-center py-12">
      <div
        class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600"
      ></div>
      <p class="mt-2 text-gray-600">Loading alert...</p>
    </div>
  {:else if error}
    <div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg">
      {error}
    </div>
  {:else if alertData}
    <!-- Alert Info Card -->
    <div class="bg-white p-6 rounded-lg shadow border-l-4 {priorityColors[alertData.priority]}">
      <div class="flex items-start justify-between mb-4">
        <div class="flex items-center gap-2">
          <span class="px-3 py-1 text-sm font-semibold rounded {priorityColors[alertData.priority]}">
            {alertData.priority}
          </span>
          <span class="px-3 py-1 text-sm font-semibold rounded {statusColors[alertData.status]}">
            {alertData.status}
          </span>
          <span class="text-sm text-gray-500">{alertData.source}</span>
        </div>
        <div class="flex gap-2">
          {#if alertData.status === 'open'}
            <Button variant="primary" size="sm" on:click={handleAcknowledge}>Acknowledge</Button>
          {/if}
          {#if alertData.status !== 'closed'}
            <Button
              variant="secondary"
              size="sm"
              on:click={() => (showSnoozeForm = !showSnoozeForm)}
            >
              Snooze
            </Button>
            <Button variant="secondary" size="sm" on:click={handleClose}>Close</Button>
          {/if}
          <Button variant="danger" size="sm" on:click={handleDelete}>Delete</Button>
        </div>
      </div>

      <h3 class="text-2xl font-semibold text-gray-900 mb-3">{alertData.message}</h3>

      {#if alertData.description}
        <p class="text-gray-700 mb-4">{alertData.description}</p>
      {/if}

      {#if alertData.tags && alertData.tags.length > 0}
        <div class="flex gap-2 flex-wrap mb-4">
          {#each alertData.tags as tag}
            <span class="px-3 py-1 text-sm bg-gray-100 text-gray-700 rounded">{tag}</span>
          {/each}
        </div>
      {/if}

      <div class="grid grid-cols-2 gap-4 pt-4 border-t border-gray-200">
        <div>
          <div class="text-sm text-gray-600">Created</div>
          <div class="font-medium">{dayjs(alertData.created_at).format('MMM D, YYYY h:mm A')}</div>
          <div class="text-xs text-gray-500">{dayjs(alertData.created_at).fromNow()}</div>
        </div>

        <div>
          <div class="text-sm text-gray-600">Assigned To</div>
          <div class="font-medium">{getAssignedToDisplay()}</div>
          <button
            on:click={() => (showAssignForm = !showAssignForm)}
            class="text-xs text-primary-600 hover:text-primary-800 mt-1"
          >
            {showAssignForm ? 'Cancel' : 'Change assignment'}
          </button>
        </div>

        {#if alertData.acknowledged_at}
          <div>
            <div class="text-sm text-gray-600">Acknowledged</div>
            <div class="font-medium">
              {dayjs(alertData.acknowledged_at).format('MMM D, YYYY h:mm A')}
            </div>
            <div class="text-xs text-gray-500">{dayjs(alertData.acknowledged_at).fromNow()}</div>
          </div>
        {/if}

        {#if alertData.closed_at}
          <div>
            <div class="text-sm text-gray-600">Closed</div>
            <div class="font-medium">{dayjs(alertData.closed_at).format('MMM D, YYYY h:mm A')}</div>
            {#if alertData.close_reason}
              <div class="text-xs text-gray-500">Reason: {alertData.close_reason}</div>
            {/if}
          </div>
        {/if}

        {#if alertData.snoozed_until}
          <div>
            <div class="text-sm text-gray-600">Snoozed Until</div>
            <div class="font-medium">{dayjs(alertData.snoozed_until).format('MMM D, YYYY h:mm A')}</div>
          </div>
        {/if}
      </div>
    </div>

    <!-- Snooze Form -->
    {#if showSnoozeForm}
      <div class="bg-white p-6 rounded-lg shadow">
        <h3 class="text-lg font-semibold mb-4">Snooze Alert</h3>
        <form on:submit|preventDefault={handleSnooze} class="space-y-4">
          <div>
            <label for="snooze-duration" class="block text-sm font-medium text-gray-700 mb-1">
              Snooze for (minutes)
            </label>
            <input
              id="snooze-duration"
              type="number"
              bind:value={snoozeMinutes}
              min="5"
              max="1440"
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500"
            />
            <p class="text-xs text-gray-500 mt-1">Max 24 hours (1440 minutes)</p>
          </div>

          <div class="flex gap-2">
            <Button type="submit" variant="primary" disabled={isSnoozing}>
              {isSnoozing ? 'Snoozing...' : 'Snooze'}
            </Button>
            <Button type="button" variant="secondary" on:click={() => (showSnoozeForm = false)}>
              Cancel
            </Button>
          </div>
        </form>
      </div>
    {/if}

    <!-- Assignment Form -->
    {#if showAssignForm}
      <div class="bg-white p-6 rounded-lg shadow">
        <h3 class="text-lg font-semibold mb-4">Assign Alert</h3>
        <form on:submit|preventDefault={handleAssign} class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Assign to</label>
            <div class="flex gap-4 mb-3">
              <label class="flex items-center">
                <input type="radio" bind:group={assignmentType} value="user" class="mr-2" />
                <span class="text-sm">User</span>
              </label>
              <label class="flex items-center">
                <input type="radio" bind:group={assignmentType} value="team" class="mr-2" />
                <span class="text-sm">Team</span>
              </label>
            </div>

            {#if assignmentType === 'user'}
              <select
                bind:value={selectedUserId}
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500"
                required
              >
                <option value="">Select a user...</option>
                {#each users as user (user.id)}
                  <option value={user.id}>
                    {user.full_name || user.username} ({user.email})
                  </option>
                {/each}
              </select>
            {:else}
              <select
                bind:value={selectedTeamId}
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500"
                required
              >
                <option value="">Select a team...</option>
                {#each teams as team (team.id)}
                  <option value={team.id}>
                    {team.name}
                    {#if team.description}- {team.description}{/if}
                  </option>
                {/each}
              </select>
            {/if}
          </div>

          {#if assignError}
            <div class="bg-red-50 border border-red-200 text-red-700 px-3 py-2 rounded text-sm">
              {assignError}
            </div>
          {/if}

          <div class="flex gap-2">
            <Button type="submit" variant="primary" disabled={isAssigning}>
              {isAssigning ? 'Assigning...' : 'Assign'}
            </Button>
            <Button type="button" variant="secondary" on:click={() => (showAssignForm = false)}>
              Cancel
            </Button>
          </div>
        </form>
      </div>
    {/if}
  {/if}
</div>
