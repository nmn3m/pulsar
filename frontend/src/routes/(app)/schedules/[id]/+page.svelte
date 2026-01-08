<script lang="ts">
  import { onMount } from 'svelte';
  import { page } from '$app/stores';
  import { goto } from '$app/navigation';
  import { api } from '$lib/api/client';
  import type { ScheduleWithRotations, OnCallUser, RotationType } from '$lib/types/schedule';
  import type { User } from '$lib/types/user';
  import Button from '$lib/components/ui/Button.svelte';
  import Input from '$lib/components/ui/Input.svelte';
  import dayjs from 'dayjs';

  let scheduleId = $page.params.id;
  let schedule: ScheduleWithRotations | null = null;
  let users: User[] = [];
  let onCallUser: OnCallUser | null = null;
  let isLoading = true;
  let error = '';

  // Create rotation form
  let showCreateRotationForm = false;
  let rotationName = '';
  let rotationType: RotationType = 'weekly';
  let rotationLength = 1;
  let startDate = '';
  let startTime = '00:00';
  let endTime = '';
  let handoffTime = '09:00';
  let rotationError = '';
  let creatingRotation = false;

  onMount(async () => {
    await loadSchedule();
    await loadUsers();
    await loadOnCallUser();
  });

  async function loadSchedule() {
    try {
      isLoading = true;
      error = '';
      schedule = await api.getSchedule(scheduleId);
    } catch (err) {
      error = err instanceof Error ? err.message : 'Failed to load schedule';
    } finally {
      isLoading = false;
    }
  }

  async function loadUsers() {
    try {
      const response = await api.listUsers();
      users = response.users;
    } catch (err) {
      console.error('Failed to load users:', err);
    }
  }

  async function loadOnCallUser() {
    try {
      onCallUser = await api.getOnCallUser(scheduleId);
    } catch (err) {
      console.log('No one currently on-call or error loading:', err);
    }
  }

  async function handleCreateRotation() {
    rotationError = '';
    creatingRotation = true;

    try {
      await api.createRotation(scheduleId, {
        name: rotationName,
        rotation_type: rotationType,
        rotation_length: rotationLength,
        start_date: startDate,
        start_time: startTime,
        end_time: endTime || undefined,
        handoff_time: handoffTime,
      });

      await loadSchedule();

      // Reset form
      rotationName = '';
      rotationType = 'weekly';
      rotationLength = 1;
      startDate = '';
      startTime = '00:00';
      endTime = '';
      handoffTime = '09:00';
      showCreateRotationForm = false;
    } catch (err) {
      rotationError = err instanceof Error ? err.message : 'Failed to create rotation';
    } finally {
      creatingRotation = false;
    }
  }

  async function handleDeleteRotation(rotationId: string, rotationName: string) {
    if (!confirm(`Delete rotation "${rotationName}"?`)) return;

    try {
      await api.deleteRotation(scheduleId, rotationId);
      await loadSchedule();
    } catch (err) {
      alert(err instanceof Error ? err.message : 'Failed to delete rotation');
    }
  }

  function getRotationTypeLabel(type: RotationType): string {
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
</script>

<svelte:head>
  <title>{schedule?.name || 'Schedule'} - Pulsar</title>
</svelte:head>

<div class="space-y-6">
  <!-- Header -->
  <div class="flex items-center gap-3 mb-6">
    <button
      on:click={() => goto('/schedules')}
      class="text-gray-600 hover:text-gray-900"
      aria-label="Back to schedules"
    >
      <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
      </svg>
    </button>
    <div>
      <h2 class="text-3xl font-bold text-gray-900">{schedule?.name || 'Loading...'}</h2>
      {#if schedule?.description}
        <p class="text-gray-600 mt-1">{schedule.description}</p>
      {/if}
    </div>
  </div>

  {#if isLoading}
    <div class="text-center py-12">
      <div
        class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600"
      ></div>
      <p class="mt-2 text-gray-600">Loading schedule...</p>
    </div>
  {:else if error}
    <div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg">
      {error}
    </div>
  {:else if schedule}
    <!-- Currently On-Call -->
    <div class="bg-white p-6 rounded-lg shadow">
      <h3 class="text-lg font-semibold mb-4">Currently On-Call</h3>
      {#if onCallUser}
        <div class="flex items-center gap-4">
          <div class="flex-1">
            <p class="font-medium text-gray-900">
              {onCallUser.user?.full_name || onCallUser.user?.username || 'Unknown User'}
            </p>
            <p class="text-sm text-gray-600">{onCallUser.user?.email}</p>
            <p class="text-xs text-gray-500 mt-1">
              {dayjs(onCallUser.start_time).format('MMM D, YYYY h:mm A')} -
              {dayjs(onCallUser.end_time).format('MMM D, YYYY h:mm A')}
            </p>
            {#if onCallUser.is_override}
              <span
                class="inline-block px-2 py-1 text-xs bg-yellow-100 text-yellow-800 rounded mt-2"
              >
                Override
              </span>
            {/if}
          </div>
        </div>
      {:else}
        <p class="text-gray-600">No one currently on-call</p>
      {/if}
    </div>

    <!-- Rotations Section -->
    <div class="bg-white p-6 rounded-lg shadow">
      <div class="flex justify-between items-center mb-4">
        <h3 class="text-lg font-semibold">Rotations ({schedule.rotations?.length || 0})</h3>
        <Button
          variant="primary"
          size="sm"
          on:click={() => (showCreateRotationForm = !showCreateRotationForm)}
        >
          {showCreateRotationForm ? 'Cancel' : 'Add Rotation'}
        </Button>
      </div>

      {#if showCreateRotationForm}
        <div class="mb-6 p-4 bg-gray-50 rounded-lg">
          <h4 class="text-sm font-semibold mb-3">Create Rotation</h4>
          <form on:submit|preventDefault={handleCreateRotation} class="space-y-3">
            <Input id="rotation-name" label="Rotation Name" bind:value={rotationName} required />

            <div class="grid grid-cols-2 gap-3">
              <div>
                <label for="rotation-type" class="block text-sm font-medium text-gray-700 mb-1">
                  Rotation Type
                </label>
                <select
                  id="rotation-type"
                  bind:value={rotationType}
                  class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500"
                >
                  <option value="daily">Daily</option>
                  <option value="weekly">Weekly</option>
                  <option value="custom">Custom</option>
                </select>
              </div>

              <Input
                id="rotation-length"
                label="Rotation Length (days/weeks)"
                type="number"
                bind:value={rotationLength}
                min="1"
                required
              />
            </div>

            <div class="grid grid-cols-2 gap-3">
              <Input
                id="start-date"
                label="Start Date"
                type="date"
                bind:value={startDate}
                required
              />

              <Input
                id="start-time"
                label="Start Time"
                type="time"
                bind:value={startTime}
                required
              />
            </div>

            <div class="grid grid-cols-2 gap-3">
              <Input id="end-time" label="End Time (optional)" type="time" bind:value={endTime} />

              <Input
                id="handoff-time"
                label="Handoff Time"
                type="time"
                bind:value={handoffTime}
                required
              />
            </div>

            {#if rotationError}
              <div class="bg-red-50 border border-red-200 text-red-700 px-3 py-2 rounded text-sm">
                {rotationError}
              </div>
            {/if}

            <div class="flex gap-2">
              <Button type="submit" variant="primary" size="sm" disabled={creatingRotation}>
                {creatingRotation ? 'Creating...' : 'Create Rotation'}
              </Button>
              <Button
                type="button"
                variant="secondary"
                size="sm"
                on:click={() => (showCreateRotationForm = false)}
              >
                Cancel
              </Button>
            </div>
          </form>
        </div>
      {/if}

      {#if schedule.rotations && schedule.rotations.length > 0}
        <div class="space-y-3">
          {#each schedule.rotations as rotation (rotation.id)}
            <div class="p-4 border border-gray-200 rounded-lg">
              <div class="flex justify-between items-start">
                <div class="flex-1">
                  <h4 class="font-medium text-gray-900">{rotation.name}</h4>
                  <div class="text-sm text-gray-600 mt-1 space-y-1">
                    <p>
                      Type: {getRotationTypeLabel(rotation.rotation_type)}
                      ({rotation.rotation_length}
                      {rotation.rotation_type === 'daily' ? 'day(s)' : 'week(s)'})
                    </p>
                    <p>Start: {dayjs(rotation.start_date).format('MMM D, YYYY')}</p>
                    <p>
                      Time: {dayjs(rotation.start_time).format('HH:mm')}
                      {#if rotation.end_time}
                        - {dayjs(rotation.end_time).format('HH:mm')}
                      {:else}
                        - 24/7
                      {/if}
                    </p>
                    <p>Handoff: {dayjs(rotation.handoff_time).format('HH:mm')}</p>
                  </div>
                </div>

                <div class="flex gap-2">
                  <Button
                    variant="primary"
                    size="sm"
                    on:click={() => goto(`/schedules/${scheduleId}/rotations/${rotation.id}`)}
                  >
                    Manage
                  </Button>
                  <Button
                    variant="danger"
                    size="sm"
                    on:click={() => handleDeleteRotation(rotation.id, rotation.name)}
                  >
                    Delete
                  </Button>
                </div>
              </div>
            </div>
          {/each}
        </div>
      {:else}
        <div class="text-center py-8 text-gray-500">
          <p>No rotations configured</p>
          <p class="text-sm mt-1">Click "Add Rotation" to get started</p>
        </div>
      {/if}
    </div>
  {/if}
</div>
