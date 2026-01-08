<script lang="ts">
  import { onMount } from 'svelte';
  import { page } from '$app/stores';
  import { goto } from '$app/navigation';
  import { api } from '$lib/api/client';
  import type { TeamWithMembers, TeamRole } from '$lib/types/team';
  import type { User } from '$lib/types/user';
  import Button from '$lib/components/ui/Button.svelte';
  import Input from '$lib/components/ui/Input.svelte';

  let teamId = $page.params.id;
  let team: TeamWithMembers | null = null;
  let organizationUsers: User[] = [];
  let isLoading = true;
  let error = '';

  // Edit team form
  let showEditForm = false;
  let editName = '';
  let editDescription = '';
  let editError = '';
  let isEditing = false;

  // Add member form
  let showAddMemberForm = false;
  let selectedUserId = '';
  let selectedUserRole: TeamRole = 'member';
  let addMemberError = '';
  let isAddingMember = false;

  onMount(async () => {
    await loadTeam();
    await loadOrganizationUsers();
  });

  async function loadTeam() {
    try {
      isLoading = true;
      error = '';
      team = await api.getTeam(teamId);
      editName = team.name;
      editDescription = team.description || '';
    } catch (err) {
      error = err instanceof Error ? err.message : 'Failed to load team';
    } finally {
      isLoading = false;
    }
  }

  async function loadOrganizationUsers() {
    try {
      const response = await api.listUsers();
      organizationUsers = response.users;
    } catch (err) {
      console.error('Failed to load organization users:', err);
    }
  }

  async function handleUpdateTeam() {
    if (!team) return;

    editError = '';
    isEditing = true;

    try {
      await api.updateTeam(team.id, {
        name: editName,
        description: editDescription || undefined,
      });

      await loadTeam();
      showEditForm = false;
    } catch (err) {
      editError = err instanceof Error ? err.message : 'Failed to update team';
    } finally {
      isEditing = false;
    }
  }

  async function handleDeleteTeam() {
    if (!team) return;

    if (!confirm(`Are you sure you want to delete team "${team.name}"?`)) {
      return;
    }

    try {
      await api.deleteTeam(team.id);
      goto('/teams');
    } catch (err) {
      alert(err instanceof Error ? err.message : 'Failed to delete team');
    }
  }

  async function handleAddMember() {
    if (!team || !selectedUserId) return;

    addMemberError = '';
    isAddingMember = true;

    try {
      await api.addTeamMember(team.id, {
        user_id: selectedUserId,
        role: selectedUserRole,
      });

      await loadTeam();
      selectedUserId = '';
      selectedUserRole = 'member';
      showAddMemberForm = false;
    } catch (err) {
      addMemberError = err instanceof Error ? err.message : 'Failed to add member';
    } finally {
      isAddingMember = false;
    }
  }

  async function handleRemoveMember(userId: string, username: string) {
    if (!team) return;

    if (!confirm(`Remove ${username} from this team?`)) {
      return;
    }

    try {
      await api.removeTeamMember(team.id, userId);
      await loadTeam();
    } catch (err) {
      alert(err instanceof Error ? err.message : 'Failed to remove member');
    }
  }

  async function handleUpdateMemberRole(userId: string, role: TeamRole) {
    if (!team) return;

    try {
      await api.updateTeamMemberRole(team.id, userId, { role });
      await loadTeam();
    } catch (err) {
      alert(err instanceof Error ? err.message : 'Failed to update member role');
    }
  }

  $: availableUsers = organizationUsers.filter(
    (user) => !team?.members?.some((member) => member.id === user.id)
  );
</script>

<svelte:head>
  <title>{team?.name || 'Team'} - Pulsar</title>
</svelte:head>

<div class="space-y-6">
  <!-- Header -->
  <div class="flex justify-between items-start">
    <div>
      <div class="flex items-center gap-3 mb-2">
        <button
          on:click={() => goto('/teams')}
          class="text-gray-600 hover:text-gray-900"
          aria-label="Back to teams"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M15 19l-7-7 7-7"
            />
          </svg>
        </button>
        <h2 class="text-3xl font-bold text-gray-900">{team?.name || 'Loading...'}</h2>
      </div>
      {#if team?.description}
        <p class="text-gray-600 ml-8">{team.description}</p>
      {/if}
    </div>
    <div class="flex gap-2">
      <Button variant="secondary" on:click={() => (showEditForm = !showEditForm)}>
        {showEditForm ? 'Cancel' : 'Edit Team'}
      </Button>
      <Button variant="danger" on:click={handleDeleteTeam}>Delete Team</Button>
    </div>
  </div>

  {#if isLoading}
    <div class="text-center py-12">
      <div
        class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600"
      ></div>
      <p class="mt-2 text-gray-600">Loading team...</p>
    </div>
  {:else if error}
    <div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg">
      {error}
    </div>
  {:else if team}
    <!-- Edit Team Form -->
    {#if showEditForm}
      <div class="bg-white p-6 rounded-lg shadow">
        <h3 class="text-lg font-semibold mb-4">Edit Team</h3>
        <form on:submit|preventDefault={handleUpdateTeam} class="space-y-4">
          <Input
            id="edit-name"
            label="Team Name"
            bind:value={editName}
            placeholder="Engineering, DevOps, Support..."
            required
          />

          <div>
            <label for="edit-description" class="block text-sm font-medium text-gray-700 mb-1">
              Description
            </label>
            <textarea
              id="edit-description"
              bind:value={editDescription}
              rows="3"
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500"
              placeholder="Team description..."
            ></textarea>
          </div>

          {#if editError}
            <div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg">
              {editError}
            </div>
          {/if}

          <div class="flex gap-2">
            <Button type="submit" variant="primary" disabled={isEditing}>
              {isEditing ? 'Saving...' : 'Save Changes'}
            </Button>
            <Button type="button" variant="secondary" on:click={() => (showEditForm = false)}>
              Cancel
            </Button>
          </div>
        </form>
      </div>
    {/if}

    <!-- Team Members Section -->
    <div class="bg-white p-6 rounded-lg shadow">
      <div class="flex justify-between items-center mb-4">
        <h3 class="text-lg font-semibold">Team Members ({team.members?.length || 0})</h3>
        <Button
          variant="primary"
          size="sm"
          on:click={() => (showAddMemberForm = !showAddMemberForm)}
        >
          {showAddMemberForm ? 'Cancel' : 'Add Member'}
        </Button>
      </div>

      <!-- Add Member Form -->
      {#if showAddMemberForm}
        <div class="mb-6 p-4 bg-gray-50 rounded-lg">
          <h4 class="text-sm font-semibold mb-3">Add New Member</h4>
          <form on:submit|preventDefault={handleAddMember} class="space-y-3">
            <div>
              <label for="user-select" class="block text-sm font-medium text-gray-700 mb-1">
                Select User
              </label>
              <select
                id="user-select"
                bind:value={selectedUserId}
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500"
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

            <div>
              <label for="role-select" class="block text-sm font-medium text-gray-700 mb-1">
                Role
              </label>
              <select
                id="role-select"
                bind:value={selectedUserRole}
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500"
              >
                <option value="member">Member</option>
                <option value="lead">Lead</option>
              </select>
            </div>

            {#if addMemberError}
              <div class="bg-red-50 border border-red-200 text-red-700 px-3 py-2 rounded text-sm">
                {addMemberError}
              </div>
            {/if}

            <div class="flex gap-2">
              <Button type="submit" variant="primary" size="sm" disabled={isAddingMember}>
                {isAddingMember ? 'Adding...' : 'Add Member'}
              </Button>
              <Button
                type="button"
                variant="secondary"
                size="sm"
                on:click={() => (showAddMemberForm = false)}
              >
                Cancel
              </Button>
            </div>
          </form>
        </div>
      {/if}

      <!-- Members List -->
      {#if team.members && team.members.length > 0}
        <div class="space-y-3">
          {#each team.members as member (member.id)}
            <div class="flex items-center justify-between p-4 border border-gray-200 rounded-lg">
              <div class="flex-1">
                <div class="font-medium text-gray-900">
                  {member.full_name || member.username}
                </div>
                <div class="text-sm text-gray-600">{member.email}</div>
              </div>

              <div class="flex items-center gap-3">
                <select
                  value={member.role || 'member'}
                  on:change={(e) => handleUpdateMemberRole(member.id, e.currentTarget.value)}
                  class="px-3 py-1.5 text-sm border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500"
                >
                  <option value="member">Member</option>
                  <option value="lead">Lead</option>
                </select>

                <button
                  on:click={() => handleRemoveMember(member.id, member.username)}
                  class="text-red-600 hover:text-red-800 text-sm font-medium px-3 py-1.5"
                >
                  Remove
                </button>
              </div>
            </div>
          {/each}
        </div>
      {:else}
        <div class="text-center py-8 text-gray-500">
          <p>No members in this team yet</p>
          <p class="text-sm mt-1">Click "Add Member" to get started</p>
        </div>
      {/if}
    </div>
  {/if}
</div>
