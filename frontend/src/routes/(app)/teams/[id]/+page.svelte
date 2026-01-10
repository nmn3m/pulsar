<script lang="ts">
  import { onMount } from 'svelte';
  import { page } from '$app/stores';
  import { goto } from '$app/navigation';
  import { api } from '$lib/api/client';
  import type { TeamWithMembers, TeamRole, TeamInvitation } from '$lib/types/team';
  import type { User } from '$lib/types/user';
  import Button from '$lib/components/ui/Button.svelte';
  import Input from '$lib/components/ui/Input.svelte';

  let teamId = $page.params.id!;
  let team: TeamWithMembers | null = null;
  let organizationUsers: User[] = [];
  let invitations: TeamInvitation[] = [];
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
  let addMemberMode: 'select' | 'email' = 'email';
  let selectedUserId = '';
  let inviteEmail = '';
  let selectedUserRole: TeamRole = 'member';
  let addMemberError = '';
  let addMemberSuccess = '';
  let isAddingMember = false;

  onMount(async () => {
    await loadTeam();
    await loadOrganizationUsers();
    await loadInvitations();
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

  async function loadInvitations() {
    try {
      const response = await api.listTeamInvitations(teamId);
      invitations = response.invitations || [];
    } catch (err) {
      console.error('Failed to load invitations:', err);
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
    if (!team) return;

    addMemberError = '';
    addMemberSuccess = '';
    isAddingMember = true;

    try {
      if (addMemberMode === 'select' && selectedUserId) {
        await api.addTeamMember(team.id, {
          user_id: selectedUserId,
          role: selectedUserRole,
        });
        addMemberSuccess = 'Member added successfully';
      } else if (addMemberMode === 'email' && inviteEmail) {
        const result = await api.inviteTeamMember(team.id, {
          email: inviteEmail,
          role: selectedUserRole,
        });
        addMemberSuccess = result.message;
        if (result.invited) {
          await loadInvitations();
        }
      } else {
        addMemberError = 'Please select a user or enter an email';
        return;
      }

      await loadTeam();
      selectedUserId = '';
      inviteEmail = '';
      selectedUserRole = 'member';

      // Keep form open briefly to show success message
      setTimeout(() => {
        showAddMemberForm = false;
        addMemberSuccess = '';
      }, 2000);
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

  async function handleUpdateMemberRole(userId: string, role: string) {
    if (!team) return;

    try {
      await api.updateTeamMemberRole(team.id, userId, { role: role as TeamRole });
      await loadTeam();
    } catch (err) {
      window.alert(err instanceof Error ? err.message : 'Failed to update member role');
    }
  }

  async function handleCancelInvitation(invitationId: string) {
    if (!team) return;

    if (!confirm('Cancel this invitation?')) {
      return;
    }

    try {
      await api.cancelTeamInvitation(team.id, invitationId);
      await loadInvitations();
    } catch (err) {
      alert(err instanceof Error ? err.message : 'Failed to cancel invitation');
    }
  }

  async function handleResendInvitation(invitationId: string) {
    if (!team) return;

    try {
      await api.resendTeamInvitation(team.id, invitationId);
      alert('Invitation resent successfully');
    } catch (err) {
      alert(err instanceof Error ? err.message : 'Failed to resend invitation');
    }
  }

  function formatDate(dateStr: string): string {
    return new Date(dateStr).toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      year: 'numeric',
    });
  }

  $: availableUsers = organizationUsers.filter(
    (user) => !team?.members?.some((member) => member.id === user.id)
  );

  $: pendingInvitations = invitations.filter((inv) => inv.status === 'pending');
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
        <h3 class="text-lg font-semibold mb-4 text-gray-900">Edit Team</h3>
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
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 bg-white text-gray-900"
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
        <h3 class="text-lg font-semibold text-gray-900">
          Team Members ({team.members?.length || 0})
        </h3>
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
          <h4 class="text-sm font-semibold mb-3 text-gray-900">Add New Member</h4>

          <!-- Mode Toggle -->
          <div class="flex gap-2 mb-4">
            <button
              type="button"
              class="px-3 py-1.5 text-sm rounded-lg transition-colors {addMemberMode === 'email'
                ? 'bg-primary-600 text-white'
                : 'bg-gray-200 text-gray-700'}"
              on:click={() => (addMemberMode = 'email')}
            >
              Invite by Email
            </button>
            <button
              type="button"
              class="px-3 py-1.5 text-sm rounded-lg transition-colors {addMemberMode === 'select'
                ? 'bg-primary-600 text-white'
                : 'bg-gray-200 text-gray-700'}"
              on:click={() => (addMemberMode = 'select')}
            >
              Select Existing User
            </button>
          </div>

          <form on:submit|preventDefault={handleAddMember} class="space-y-3">
            {#if addMemberMode === 'email'}
              <div>
                <label for="invite-email" class="block text-sm font-medium text-gray-700 mb-1">
                  Email Address
                </label>
                <input
                  type="email"
                  id="invite-email"
                  bind:value={inviteEmail}
                  class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 bg-white text-gray-900"
                  placeholder="user@example.com"
                  required
                />
                <p class="mt-1 text-xs text-gray-500">
                  If the user exists, they'll be added directly. Otherwise, an invitation email will
                  be sent.
                </p>
              </div>
            {:else}
              <div>
                <label for="user-select" class="block text-sm font-medium text-gray-700 mb-1">
                  Select User
                </label>
                <select
                  id="user-select"
                  bind:value={selectedUserId}
                  class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 bg-white text-gray-900"
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
            {/if}

            <div>
              <label for="role-select" class="block text-sm font-medium text-gray-700 mb-1">
                Role
              </label>
              <select
                id="role-select"
                bind:value={selectedUserRole}
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 bg-white text-gray-900"
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

            {#if addMemberSuccess}
              <div
                class="bg-green-50 border border-green-200 text-green-700 px-3 py-2 rounded text-sm"
              >
                {addMemberSuccess}
              </div>
            {/if}

            <div class="flex gap-2">
              <Button type="submit" variant="primary" size="sm" disabled={isAddingMember}>
                {isAddingMember
                  ? 'Adding...'
                  : addMemberMode === 'email'
                    ? 'Send Invitation'
                    : 'Add Member'}
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
                  class="px-3 py-1.5 text-sm border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 bg-white text-gray-900"
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

    <!-- Pending Invitations Section -->
    {#if pendingInvitations.length > 0}
      <div class="bg-white p-6 rounded-lg shadow">
        <h3 class="text-lg font-semibold mb-4 text-gray-900">
          Pending Invitations ({pendingInvitations.length})
        </h3>
        <div class="space-y-3">
          {#each pendingInvitations as invitation (invitation.id)}
            <div
              class="flex items-center justify-between p-4 border border-yellow-200 bg-yellow-50 rounded-lg"
            >
              <div class="flex-1">
                <div class="font-medium text-gray-900">{invitation.email}</div>
                <div class="text-sm text-gray-600">
                  Role: {invitation.role} | Expires: {formatDate(invitation.expires_at)}
                </div>
              </div>
              <div class="flex items-center gap-2">
                <button
                  on:click={() => handleResendInvitation(invitation.id)}
                  class="text-primary-600 hover:text-primary-800 text-sm font-medium px-3 py-1.5"
                >
                  Resend
                </button>
                <button
                  on:click={() => handleCancelInvitation(invitation.id)}
                  class="text-red-600 hover:text-red-800 text-sm font-medium px-3 py-1.5"
                >
                  Cancel
                </button>
              </div>
            </div>
          {/each}
        </div>
      </div>
    {/if}
  {/if}
</div>
