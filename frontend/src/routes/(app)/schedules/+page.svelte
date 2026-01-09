<script lang="ts">
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { schedulesStore } from '$lib/stores/schedules';
  import Button from '$lib/components/ui/Button.svelte';
  import Input from '$lib/components/ui/Input.svelte';

  let showCreateForm = false;
  let name = '';
  let description = '';
  let timezone = 'UTC';
  let createError = '';
  let creating = false;

  onMount(() => {
    schedulesStore.load();
  });

  async function handleCreateSchedule() {
    createError = '';
    creating = true;

    try {
      await schedulesStore.create({
        name,
        description: description || undefined,
        timezone,
      });

      // Reset form
      name = '';
      description = '';
      timezone = 'UTC';
      showCreateForm = false;
    } catch (err) {
      createError = err instanceof Error ? err.message : 'Failed to create schedule';
    } finally {
      creating = false;
    }
  }

  async function handleDeleteSchedule(id: string, scheduleName: string) {
    if (!confirm(`Are you sure you want to delete schedule "${scheduleName}"?`)) {
      return;
    }

    try {
      await schedulesStore.delete(id);
    } catch (err) {
      alert(err instanceof Error ? err.message : 'Failed to delete schedule');
    }
  }
</script>

<svelte:head>
  <title>Schedules - Pulsar</title>
</svelte:head>

<div class="space-y-6">
  <div class="flex justify-between items-center">
    <div>
      <h2 class="text-3xl font-bold text-gray-900">On-Call Schedules</h2>
      <p class="mt-2 text-gray-500">Manage on-call rotations and schedules</p>
    </div>
    <Button variant="primary" on:click={() => (showCreateForm = !showCreateForm)}>
      {showCreateForm ? 'Cancel' : 'Create Schedule'}
    </Button>
  </div>

  {#if showCreateForm}
    <div
      class="bg-white backdrop-blur-sm p-6 rounded-xl border border-primary-200 shadow-sm"
    >
      <h3 class="text-lg font-semibold mb-4 text-gray-900">
        Create New Schedule
      </h3>
      <form on:submit|preventDefault={handleCreateSchedule} class="space-y-4">
        <Input
          id="name"
          label="Schedule Name"
          bind:value={name}
          placeholder="Primary On-Call, Weekend Support..."
          required
        />

        <div>
          <label
            for="description"
            class="block text-sm font-medium text-gray-600 mb-1"
          >
            Description
          </label>
          <textarea
            id="description"
            bind:value={description}
            rows="3"
            class="w-full px-3 py-2 bg-white border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500 text-gray-900 placeholder-gray-400"
            placeholder="Schedule description..."
          ></textarea>
        </div>

        <Input
          id="timezone"
          label="Timezone"
          bind:value={timezone}
          placeholder="UTC, America/New_York..."
          required
        />

        {#if createError}
          <div
            class="bg-red-50 border border-red-200 text-red-600 px-4 py-3 rounded-lg"
          >
            {createError}
          </div>
        {/if}

        <div class="flex gap-2">
          <Button type="submit" variant="primary" disabled={creating}>
            {creating ? 'Creating...' : 'Create Schedule'}
          </Button>
          <Button type="button" variant="secondary" on:click={() => (showCreateForm = false)}>
            Cancel
          </Button>
        </div>
      </form>
    </div>
  {/if}

  {#if $schedulesStore.isLoading}
    <div class="text-center py-12">
      <div
        class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-primary-500"
      ></div>
      <p class="mt-2 text-gray-500">Loading schedules...</p>
    </div>
  {:else if $schedulesStore.error}
    <div
      class="bg-red-50 border border-red-200 text-red-600 px-4 py-3 rounded-lg"
    >
      {$schedulesStore.error}
    </div>
  {:else if $schedulesStore.schedules.length === 0}
    <div
      class="text-center py-12 bg-white backdrop-blur-sm rounded-xl border border-gray-200 shadow-sm"
    >
      <p class="text-gray-600">No schedules found</p>
      <p class="text-sm text-gray-400 mt-2">
        Create your first schedule to get started
      </p>
    </div>
  {:else}
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      {#each $schedulesStore.schedules as schedule (schedule.id)}
        <div
          class="bg-white backdrop-blur-sm p-6 rounded-xl border border-gray-200 hover:border-primary-400 transition-all duration-300 hover:shadow-lg shadow-sm"
        >
          <div class="flex justify-between items-start mb-4">
            <div class="flex-1">
              <h3 class="text-lg font-semibold text-gray-900">
                {schedule.name}
              </h3>
              {#if schedule.description}
                <p class="text-sm text-gray-500 mt-1">{schedule.description}</p>
              {/if}
              <p class="text-xs text-gray-400 mt-2">
                Timezone: {schedule.timezone}
              </p>
            </div>
          </div>

          <div class="flex gap-2 mt-4">
            <Button variant="primary" size="sm" on:click={() => goto(`/schedules/${schedule.id}`)}>
              View Schedule
            </Button>
            <Button
              variant="danger"
              size="sm"
              on:click={() => handleDeleteSchedule(schedule.id, schedule.name)}
            >
              Delete
            </Button>
          </div>
        </div>
      {/each}
    </div>
  {/if}
</div>
