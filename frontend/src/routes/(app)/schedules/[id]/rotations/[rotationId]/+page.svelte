<script lang="ts">
  import { onMount } from 'svelte';
  import { page } from '$app/stores';
  import { goto } from '$app/navigation';
  import { api } from '$lib/api/client';
  import type { ScheduleRotation, ParticipantWithUser } from '$lib/types/schedule';
  import type { User } from '$lib/types/user';
  import Button from '$lib/components/ui/Button.svelte';
  import dayjs from 'dayjs';

  let scheduleId = $page.params.id;
  let rotationId = $page.params.rotationId;

  let rotation: ScheduleRotation | null = null;
  let participants: ParticipantWithUser[] = [];
  let users: User[] = [];
  let isLoading = true;
  let error = '';

  // Add participant form
  let showAddParticipantForm = false;
  let selectedUserId = '';
  let addingParticipant = false;
  let addError = '';

  onMount(async () => {
    await Promise.all([loadRotation(), loadParticipants(), loadUsers()]);
  });

  async function loadRotation() {
    try {
      isLoading = true;
      error = '';
      rotation = await api.getRotation(scheduleId, rotationId);
    } catch (err) {
      error = err instanceof Error ? err.message : 'Failed to load rotation';
    } finally {
      isLoading = false;
    }
  }

  async function loadParticipants() {
    try {
      const response = await api.listParticipants(scheduleId, rotationId);
      participants = response.participants || [];
    } catch (err) {
      console.error('Failed to load participants:', err);
    }
  }

  async function loadUsers() {
    try {
      const response = await api.listUsers();
      users = response.users || [];
    } catch (err) {
      console.error('Failed to load users:', err);
    }
  }

  async function handleAddParticipant() {
    if (!selectedUserId) return;

    addError = '';
    addingParticipant = true;

    try {
      // Position is the next available slot (0-indexed)
      const position = participants.length;
      await api.addParticipant(scheduleId, rotationId, { user_id: selectedUserId, position });
      await loadParticipants();
      selectedUserId = '';
      showAddParticipantForm = false;
    } catch (err) {
      addError = err instanceof Error ? err.message : 'Failed to add participant';
    } finally {
      addingParticipant = false;
    }
  }

  async function handleRemoveParticipant(userId: string, userName: string) {
    if (!confirm(`Remove ${userName} from rotation?`)) return;

    try {
      await api.removeParticipant(scheduleId, rotationId, userId);
      await loadParticipants();
    } catch (err) {
      alert(err instanceof Error ? err.message : 'Failed to remove participant');
    }
  }

  function getRotationTypeLabel(type: string): string {
    switch (type) {
      case 'daily':
        return 'Daily';
      case 'weekly':
        return 'Weekly';
      case 'custom':
        return 'Custom';
      default:
        return type;
    }
  }

  // Filter out users already in the rotation
  $: availableUsers = users.filter((user) => !participants.some((p) => p.user_id === user.id));
</script>

<svelte:head>
  <title>{rotation?.name || 'Rotation'} - Pulsar</title>
</svelte:head>

<div class="space-y-6">
  <!-- Header -->
  <div class="flex items-center gap-3 mb-6">
    <button
      on:click={() => goto(`/schedules/${scheduleId}`)}
      class="text-gray-600 hover:text-gray-900 dark:text-gray-400 dark:hover:text-gray-200"
      aria-label="Back to schedule"
    >
      <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
      </svg>
    </button>
    <div>
      <h2 class="text-3xl font-bold text-gray-900 dark:text-white">
        {rotation?.name || 'Loading...'}
      </h2>
      <p class="text-gray-600 dark:text-gray-400 mt-1">Manage rotation participants</p>
    </div>
  </div>

  {#if isLoading}
    <div class="text-center py-12">
      <div
        class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600"
      ></div>
      <p class="mt-2 text-gray-600 dark:text-gray-400">Loading rotation...</p>
    </div>
  {:else if error}
    <div
      class="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 text-red-700 dark:text-red-400 px-4 py-3 rounded-lg"
    >
      {error}
    </div>
  {:else if rotation}
    <!-- Rotation Details -->
    <div class="bg-white dark:bg-gray-800 p-6 rounded-lg shadow">
      <h3 class="text-lg font-semibold mb-4 text-gray-900 dark:text-white">Rotation Details</h3>
      <div class="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm">
        <div>
          <span class="text-gray-500 dark:text-gray-400">Type</span>
          <p class="font-medium text-gray-900 dark:text-white">
            {getRotationTypeLabel(rotation.rotation_type)}
          </p>
        </div>
        <div>
          <span class="text-gray-500 dark:text-gray-400">Length</span>
          <p class="font-medium text-gray-900 dark:text-white">
            {rotation.rotation_length}
            {rotation.rotation_type === 'daily' ? 'day(s)' : 'week(s)'}
          </p>
        </div>
        <div>
          <span class="text-gray-500 dark:text-gray-400">Start Date</span>
          <p class="font-medium text-gray-900 dark:text-white">
            {dayjs(rotation.start_date).format('MMM D, YYYY')}
          </p>
        </div>
        <div>
          <span class="text-gray-500 dark:text-gray-400">Handoff Time</span>
          <p class="font-medium text-gray-900 dark:text-white">
            {dayjs(rotation.handoff_time).format('HH:mm')}
          </p>
        </div>
      </div>
    </div>

    <!-- Participants Section -->
    <div class="bg-white dark:bg-gray-800 p-6 rounded-lg shadow">
      <div class="flex justify-between items-center mb-4">
        <h3 class="text-lg font-semibold text-gray-900 dark:text-white">
          Participants ({participants.length})
        </h3>
        <Button
          variant="primary"
          size="sm"
          on:click={() => (showAddParticipantForm = !showAddParticipantForm)}
        >
          {showAddParticipantForm ? 'Cancel' : 'Add Participant'}
        </Button>
      </div>

      {#if showAddParticipantForm}
        <div class="mb-6 p-4 bg-gray-50 dark:bg-gray-700 rounded-lg">
          <h4 class="text-sm font-semibold mb-3 text-gray-900 dark:text-white">Add Participant</h4>
          <form on:submit|preventDefault={handleAddParticipant} class="space-y-3">
            <div>
              <label
                for="user-select"
                class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
              >
                Select User
              </label>
              <select
                id="user-select"
                bind:value={selectedUserId}
                class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 bg-white dark:bg-gray-800 text-gray-900 dark:text-white"
                required
              >
                <option value="">Choose a user...</option>
                {#each availableUsers as user (user.id)}
                  <option value={user.id}>
                    {user.full_name || user.username} ({user.email})
                  </option>
                {/each}
              </select>
            </div>

            {#if addError}
              <div
                class="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 text-red-700 dark:text-red-400 px-3 py-2 rounded text-sm"
              >
                {addError}
              </div>
            {/if}

            {#if availableUsers.length === 0}
              <p class="text-sm text-gray-500 dark:text-gray-400">
                All users are already in this rotation.
              </p>
            {/if}

            <div class="flex gap-2">
              <Button
                type="submit"
                variant="primary"
                size="sm"
                disabled={addingParticipant || !selectedUserId}
              >
                {addingParticipant ? 'Adding...' : 'Add Participant'}
              </Button>
              <Button
                type="button"
                variant="secondary"
                size="sm"
                on:click={() => (showAddParticipantForm = false)}
              >
                Cancel
              </Button>
            </div>
          </form>
        </div>
      {/if}

      {#if participants.length > 0}
        <div class="space-y-2">
          {#each participants as participant, index (participant.user_id)}
            <div
              class="flex items-center justify-between p-3 border border-gray-200 dark:border-gray-700 rounded-lg"
            >
              <div class="flex items-center gap-3">
                <span
                  class="w-6 h-6 flex items-center justify-center bg-primary-100 dark:bg-primary-900 text-primary-700 dark:text-primary-300 text-sm font-medium rounded-full"
                >
                  {index + 1}
                </span>
                <div>
                  <p class="font-medium text-gray-900 dark:text-white">
                    {participant.user?.full_name || participant.user?.username || 'Unknown User'}
                  </p>
                  <p class="text-sm text-gray-500 dark:text-gray-400">
                    {participant.user?.email}
                  </p>
                </div>
              </div>
              <Button
                variant="danger"
                size="sm"
                on:click={() =>
                  handleRemoveParticipant(
                    participant.user_id,
                    participant.user?.full_name || participant.user?.username || 'User'
                  )}
              >
                Remove
              </Button>
            </div>
          {/each}
        </div>
        <p class="text-sm text-gray-500 dark:text-gray-400 mt-4">
          Participants rotate in the order shown above.
        </p>
      {:else}
        <div class="text-center py-8 text-gray-500 dark:text-gray-400">
          <p>No participants in this rotation</p>
          <p class="text-sm mt-1">Click "Add Participant" to add users to this rotation</p>
        </div>
      {/if}
    </div>
  {/if}
</div>
