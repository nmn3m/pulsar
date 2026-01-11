<script lang="ts">
  import { onMount } from 'svelte';
  import { page } from '$app/stores';
  import { goto } from '$app/navigation';
  import { api } from '$lib/api/client';
  import type {
    EscalationPolicyWithRules,
    EscalationTargetType,
    TargetNotificationConfig,
  } from '$lib/types/escalation';
  import type { User } from '$lib/types/user';
  import type { Team } from '$lib/types/team';
  import type { Schedule } from '$lib/types/schedule';
  import type { NotificationChannel } from '$lib/types/notification';
  import Button from '$lib/components/ui/Button.svelte';
  import Input from '$lib/components/ui/Input.svelte';

  let policyId = $page.params.id!;
  let policy: EscalationPolicyWithRules | null = null;
  let users: User[] = [];
  let teams: Team[] = [];
  let schedules: Schedule[] = [];
  let notificationChannels: NotificationChannel[] = [];
  let isLoading = true;
  let error = '';

  // Create rule form
  let showCreateRuleForm = false;
  let rulePosition = '1';
  let escalationDelay = '15';
  let ruleError = '';
  let creatingRule = false;

  // Edit policy form
  let showEditPolicyForm = false;
  let editPolicyName = '';
  let editPolicyDescription = '';
  let editRepeatEnabled = false;
  let editRepeatCount: string | undefined = undefined;
  let editPolicyError = '';
  let editingPolicy = false;

  // Edit rule form
  let editingRuleId: string | null = null;
  let editRuleDelay = 0;
  let editRuleError = '';
  let savingRule = false;

  // Add target with channel override
  let showAddTargetModal = false;
  let addTargetRuleId: string | null = null;
  let addTargetType: EscalationTargetType = 'user';
  let addTargetId = '';
  let addTargetChannels: string[] = [];
  let addTargetUrgent = false;
  let addingTarget = false;

  onMount(async () => {
    await loadPolicy();
    await loadResources();
  });

  async function loadPolicy() {
    try {
      isLoading = true;
      error = '';
      policy = await api.getEscalationPolicy(policyId);
    } catch (err) {
      error = err instanceof Error ? err.message : 'Failed to load escalation policy';
    } finally {
      isLoading = false;
    }
  }

  async function loadResources() {
    try {
      const [usersResp, teamsResp, schedulesResp, channelsResp] = await Promise.all([
        api.listUsers(),
        api.listTeams(),
        api.listSchedules(),
        api.listNotificationChannels(),
      ]);
      users = usersResp.users;
      teams = teamsResp.teams;
      schedules = schedulesResp.schedules;
      notificationChannels = channelsResp.channels || [];
    } catch (err) {
      console.error('Failed to load resources:', err);
    }
  }

  async function handleCreateRule() {
    ruleError = '';
    creatingRule = true;

    try {
      await api.createEscalationRule(policyId, {
        position: Number(rulePosition),
        escalation_delay: Number(escalationDelay),
      });

      await loadPolicy();

      // Reset form
      rulePosition = policy ? String(policy.rules.length + 1) : '1';
      escalationDelay = '15';
      showCreateRuleForm = false;
    } catch (err) {
      ruleError = err instanceof Error ? err.message : 'Failed to create rule';
    } finally {
      creatingRule = false;
    }
  }

  async function handleDeleteRule(ruleId: string) {
    if (!confirm('Delete this escalation rule?')) return;

    try {
      await api.deleteEscalationRule(policyId, ruleId);
      await loadPolicy();
    } catch (err) {
      alert(err instanceof Error ? err.message : 'Failed to delete rule');
    }
  }

  function openAddTargetModal(ruleId: string) {
    addTargetRuleId = ruleId;
    addTargetType = 'user';
    addTargetId = '';
    addTargetChannels = [];
    addTargetUrgent = false;
    showAddTargetModal = true;
  }

  function closeAddTargetModal() {
    showAddTargetModal = false;
    addTargetRuleId = null;
  }

  async function handleAddTargetWithChannels() {
    if (!addTargetRuleId || !addTargetId) return;

    addingTarget = true;
    try {
      const notificationConfig: TargetNotificationConfig | undefined =
        addTargetChannels.length > 0
          ? { channels: addTargetChannels, urgent: addTargetUrgent }
          : undefined;

      await api.addEscalationTarget(policyId, addTargetRuleId, {
        target_type: addTargetType,
        target_id: addTargetId,
        notification_channels: notificationConfig,
      });
      await loadPolicy();
      closeAddTargetModal();
    } catch (err) {
      alert(err instanceof Error ? err.message : 'Failed to add target');
    } finally {
      addingTarget = false;
    }
  }

  async function handleQuickAddTarget(
    ruleId: string,
    targetType: EscalationTargetType,
    targetId: string
  ) {
    try {
      await api.addEscalationTarget(policyId, ruleId, {
        target_type: targetType,
        target_id: targetId,
      });
      await loadPolicy();
    } catch (err) {
      alert(err instanceof Error ? err.message : 'Failed to add target');
    }
  }

  function toggleChannel(channelType: string) {
    if (addTargetChannels.includes(channelType)) {
      addTargetChannels = addTargetChannels.filter((c) => c !== channelType);
    } else {
      addTargetChannels = [...addTargetChannels, channelType];
    }
  }

  function getChannelTypeDisplay(channelType: string): string {
    const displayNames: Record<string, string> = {
      email: 'Email',
      slack: 'Slack',
      sms: 'SMS',
      webhook: 'Webhook',
      msteams: 'MS Teams',
    };
    return displayNames[channelType] || channelType;
  }

  function getTargetChannelsDisplay(channels: TargetNotificationConfig | undefined): string {
    if (!channels || channels.channels.length === 0) return '';
    return channels.channels.map(getChannelTypeDisplay).join(', ');
  }

  async function handleRemoveTarget(ruleId: string, targetId: string) {
    if (!confirm('Remove this target?')) return;

    try {
      await api.removeEscalationTarget(policyId, ruleId, targetId);
      await loadPolicy();
    } catch (err) {
      alert(err instanceof Error ? err.message : 'Failed to remove target');
    }
  }

  function getTargetDisplay(targetType: EscalationTargetType, targetId: string): string {
    if (targetType === 'user') {
      const user = users.find((u) => u.id === targetId);
      return user ? `👤 ${user.full_name || user.username}` : 'Unknown User';
    } else if (targetType === 'team') {
      const team = teams.find((t) => t.id === targetId);
      return team ? `👥 ${team.name}` : 'Unknown Team';
    } else if (targetType === 'schedule') {
      const schedule = schedules.find((s) => s.id === targetId);
      return schedule ? `📅 ${schedule.name}` : 'Unknown Schedule';
    }
    return 'Unknown Target';
  }

  function openEditPolicyForm() {
    if (!policy) return;
    editPolicyName = policy.name;
    editPolicyDescription = policy.description || '';
    editRepeatEnabled = policy.repeat_enabled;
    editRepeatCount = policy.repeat_count ? String(policy.repeat_count) : undefined;
    editPolicyError = '';
    showEditPolicyForm = true;
  }

  async function handleUpdatePolicy() {
    if (!policy) return;
    editPolicyError = '';
    editingPolicy = true;

    try {
      await api.updateEscalationPolicy(policyId, {
        name: editPolicyName,
        description: editPolicyDescription || undefined,
        repeat_enabled: editRepeatEnabled,
        repeat_count: editRepeatEnabled && editRepeatCount ? Number(editRepeatCount) : undefined,
      });

      await loadPolicy();
      showEditPolicyForm = false;
    } catch (err) {
      editPolicyError = err instanceof Error ? err.message : 'Failed to update policy';
    } finally {
      editingPolicy = false;
    }
  }

  function startEditingRule(ruleId: string, currentDelay: number) {
    editingRuleId = ruleId;
    editRuleDelay = currentDelay;
    editRuleError = '';
  }

  function cancelEditingRule() {
    editingRuleId = null;
    editRuleError = '';
  }

  async function handleUpdateRule(ruleId: string) {
    editRuleError = '';
    savingRule = true;

    try {
      await api.updateEscalationRule(policyId, ruleId, {
        escalation_delay: Number(editRuleDelay),
      });

      await loadPolicy();
      editingRuleId = null;
    } catch (err) {
      editRuleError = err instanceof Error ? err.message : 'Failed to update rule';
    } finally {
      savingRule = false;
    }
  }
</script>

<svelte:head>
  <title>{policy?.name || 'Escalation Policy'} - Pulsar</title>
</svelte:head>

<div class="space-y-6">
  <!-- Header -->
  <div class="flex items-center gap-3 mb-6">
    <button
      on:click={() => goto('/escalation-policies')}
      class="text-gray-600 hover:text-gray-900"
      aria-label="Back to escalation policies"
    >
      <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
      </svg>
    </button>
    <div>
      <h2 class="text-3xl font-bold text-gray-900">{policy?.name || 'Loading...'}</h2>
      {#if policy?.description}
        <p class="text-gray-600 mt-1">{policy.description}</p>
      {/if}
    </div>
  </div>

  {#if isLoading}
    <div class="text-center py-12">
      <div
        class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600"
      ></div>
      <p class="mt-2 text-gray-600">Loading policy...</p>
    </div>
  {:else if error}
    <div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg">
      {error}
    </div>
  {:else if policy}
    <!-- Policy Info -->
    <div class="bg-white p-6 rounded-lg shadow">
      <div class="flex justify-between items-start mb-4">
        <h3 class="text-lg font-semibold">Policy Settings</h3>
        <Button variant="secondary" size="sm" on:click={openEditPolicyForm}>Edit Policy</Button>
      </div>

      {#if showEditPolicyForm}
        <div class="mb-4 p-4 bg-gray-50 rounded-lg border border-gray-200">
          <h4 class="text-sm font-semibold mb-3">Edit Policy</h4>
          <form on:submit|preventDefault={handleUpdatePolicy} class="space-y-3">
            <Input id="edit-policy-name" label="Policy Name" bind:value={editPolicyName} required />

            <div>
              <label
                for="edit-policy-description"
                class="block text-sm font-medium text-gray-700 mb-1"
              >
                Description
              </label>
              <textarea
                id="edit-policy-description"
                bind:value={editPolicyDescription}
                rows="2"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 text-sm"
                placeholder="Policy description..."
              ></textarea>
            </div>

            <div class="flex items-center gap-2">
              <input
                id="edit-repeat-enabled"
                type="checkbox"
                bind:checked={editRepeatEnabled}
                class="rounded bg-white border-gray-300 text-primary-600 focus:ring-primary-500"
              />
              <label for="edit-repeat-enabled" class="text-sm font-medium text-gray-700">
                Enable repeat escalation
              </label>
            </div>

            {#if editRepeatEnabled}
              <div>
                <label for="edit-repeat-count" class="block text-sm font-medium text-gray-700 mb-1">
                  Maximum Repeat Count (leave empty for infinite)
                </label>
                <input
                  id="edit-repeat-count"
                  type="number"
                  bind:value={editRepeatCount}
                  min="1"
                  placeholder="Leave empty for infinite"
                  class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 text-sm"
                />
              </div>
            {/if}

            {#if editPolicyError}
              <div class="bg-red-50 border border-red-200 text-red-700 px-3 py-2 rounded text-sm">
                {editPolicyError}
              </div>
            {/if}

            <div class="flex gap-2">
              <Button type="submit" variant="primary" size="sm" disabled={editingPolicy}>
                {editingPolicy ? 'Saving...' : 'Save Changes'}
              </Button>
              <Button
                type="button"
                variant="secondary"
                size="sm"
                on:click={() => (showEditPolicyForm = false)}
              >
                Cancel
              </Button>
            </div>
          </form>
        </div>
      {:else}
        <div class="text-sm text-gray-600 space-y-1">
          <p>
            Repeat: {policy.repeat_enabled
              ? policy.repeat_count
                ? `${policy.repeat_count} times`
                : 'Infinite'
              : 'Disabled'}
          </p>
        </div>
      {/if}
    </div>

    <!-- Escalation Rules -->
    <div class="bg-white p-6 rounded-lg shadow">
      <div class="flex justify-between items-center mb-4">
        <h3 class="text-lg font-semibold">Escalation Rules ({policy.rules?.length || 0})</h3>
        <Button
          variant="primary"
          size="sm"
          on:click={() => {
            rulePosition = String((policy?.rules?.length || 0) + 1);
            showCreateRuleForm = !showCreateRuleForm;
          }}
        >
          {showCreateRuleForm ? 'Cancel' : 'Add Rule'}
        </Button>
      </div>

      {#if showCreateRuleForm}
        <div class="mb-6 p-4 bg-gray-50 rounded-lg">
          <h4 class="text-sm font-semibold mb-3">Create Escalation Rule</h4>
          <form on:submit|preventDefault={handleCreateRule} class="space-y-3">
            <div>
              <label for="position" class="block text-sm font-medium text-gray-700 mb-1">
                Position (order)
              </label>
              <input
                id="position"
                type="number"
                bind:value={rulePosition}
                min="1"
                required
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 text-sm"
              />
            </div>

            <div>
              <label for="escalation-delay" class="block text-sm font-medium text-gray-700 mb-1">
                Escalation Delay (minutes)
              </label>
              <input
                id="escalation-delay"
                type="number"
                bind:value={escalationDelay}
                min="0"
                required
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 text-sm"
              />
            </div>

            {#if ruleError}
              <div class="bg-red-50 border border-red-200 text-red-700 px-3 py-2 rounded text-sm">
                {ruleError}
              </div>
            {/if}

            <div class="flex gap-2">
              <Button type="submit" variant="primary" size="sm" disabled={creatingRule}>
                {creatingRule ? 'Creating...' : 'Create Rule'}
              </Button>
              <Button
                type="button"
                variant="secondary"
                size="sm"
                on:click={() => (showCreateRuleForm = false)}
              >
                Cancel
              </Button>
            </div>
          </form>
        </div>
      {/if}

      {#if policy.rules && policy.rules.length > 0}
        <div class="space-y-4">
          {#each policy.rules as rule (rule.id)}
            <div class="border border-gray-200 rounded-lg p-4">
              <div class="flex justify-between items-start mb-3">
                <div class="flex-1">
                  <span class="text-sm font-semibold text-gray-900">
                    Level {rule.position + 1}
                  </span>

                  {#if editingRuleId === rule.id}
                    <div class="mt-2 flex items-center gap-2">
                      <span class="text-sm text-gray-600">Escalate after</span>
                      <input
                        type="number"
                        bind:value={editRuleDelay}
                        min="0"
                        class="w-20 px-2 py-1 text-sm border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-primary-500"
                      />
                      <span class="text-sm text-gray-600">minutes</span>
                      <Button
                        variant="primary"
                        size="sm"
                        on:click={() => handleUpdateRule(rule.id)}
                        disabled={savingRule}
                      >
                        {savingRule ? 'Saving...' : 'Save'}
                      </Button>
                      <Button variant="secondary" size="sm" on:click={cancelEditingRule}>
                        Cancel
                      </Button>
                    </div>
                    {#if editRuleError}
                      <p class="text-xs text-red-600 mt-1">{editRuleError}</p>
                    {/if}
                  {:else}
                    <p class="text-sm text-gray-600 mt-1">
                      Escalate after {rule.escalation_delay} minutes
                      <button
                        on:click={() => startEditingRule(rule.id, rule.escalation_delay)}
                        class="ml-2 text-primary-600 hover:text-primary-800 text-xs underline"
                      >
                        Edit
                      </button>
                    </p>
                  {/if}
                </div>
                <Button variant="danger" size="sm" on:click={() => handleDeleteRule(rule.id)}>
                  Delete
                </Button>
              </div>

              <!-- Targets -->
              <div class="mt-3">
                <p class="text-xs font-semibold text-gray-700 mb-2">Notify:</p>
                {#if rule.targets && rule.targets.length > 0}
                  <div class="space-y-1">
                    {#each rule.targets as target (target.id)}
                      <div class="flex items-center justify-between text-sm bg-gray-50 p-2 rounded">
                        <div>
                          <span>{getTargetDisplay(target.target_type, target.target_id)}</span>
                          {#if target.notification_channels && target.notification_channels.channels?.length > 0}
                            <span class="ml-2 text-xs text-primary-600">
                              via {getTargetChannelsDisplay(target.notification_channels)}
                              {#if target.notification_channels.urgent}
                                <span class="text-red-600">(urgent)</span>
                              {/if}
                            </span>
                          {/if}
                        </div>
                        <button
                          on:click={() => handleRemoveTarget(rule.id, target.id)}
                          class="text-red-600 hover:text-red-800 text-xs"
                        >
                          Remove
                        </button>
                      </div>
                    {/each}
                  </div>
                {:else}
                  <p class="text-xs text-gray-500">No targets added</p>
                {/if}

                <!-- Quick add target buttons -->
                <div class="mt-2 flex gap-2 flex-wrap">
                  <select
                    on:change={(e) => {
                      const value = e.currentTarget.value;
                      if (value) {
                        handleQuickAddTarget(rule.id, 'user', value);
                        e.currentTarget.value = '';
                      }
                    }}
                    class="text-xs px-2 py-1 border border-gray-300 rounded"
                  >
                    <option value="">+ Add User</option>
                    {#each users as user (user.id)}
                      <option value={user.id}>{user.full_name || user.username}</option>
                    {/each}
                  </select>

                  <select
                    on:change={(e) => {
                      const value = e.currentTarget.value;
                      if (value) {
                        handleQuickAddTarget(rule.id, 'team', value);
                        e.currentTarget.value = '';
                      }
                    }}
                    class="text-xs px-2 py-1 border border-gray-300 rounded"
                  >
                    <option value="">+ Add Team</option>
                    {#each teams as team (team.id)}
                      <option value={team.id}>{team.name}</option>
                    {/each}
                  </select>

                  <select
                    on:change={(e) => {
                      const value = e.currentTarget.value;
                      if (value) {
                        handleQuickAddTarget(rule.id, 'schedule', value);
                        e.currentTarget.value = '';
                      }
                    }}
                    class="text-xs px-2 py-1 border border-gray-300 rounded"
                  >
                    <option value="">+ Add Schedule</option>
                    {#each schedules as schedule (schedule.id)}
                      <option value={schedule.id}>{schedule.name}</option>
                    {/each}
                  </select>

                  <button
                    on:click={() => openAddTargetModal(rule.id)}
                    class="text-xs px-2 py-1 text-primary-600 hover:text-primary-800 border border-primary-300 rounded hover:bg-primary-50"
                  >
                    + Advanced
                  </button>
                </div>
              </div>
            </div>
          {/each}
        </div>
      {:else}
        <div class="text-center py-8 text-gray-500">
          <p>No escalation rules configured</p>
          <p class="text-sm mt-1">Click "Add Rule" to create the first escalation level</p>
        </div>
      {/if}
    </div>
  {/if}
</div>

<!-- Add Target Modal -->
{#if showAddTargetModal}
  <div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
    <div class="bg-white rounded-lg shadow-xl max-w-md w-full mx-4 p-6">
      <h3 class="text-lg font-semibold mb-4">Add Target with Channel Override</h3>

      <div class="space-y-4">
        <!-- Target Type -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Target Type</label>
          <select
            bind:value={addTargetType}
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500"
          >
            <option value="user">User</option>
            <option value="team">Team</option>
            <option value="schedule">Schedule</option>
          </select>
        </div>

        <!-- Target Selection -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">
            Select {addTargetType === 'user'
              ? 'User'
              : addTargetType === 'team'
                ? 'Team'
                : 'Schedule'}
          </label>
          <select
            bind:value={addTargetId}
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500"
          >
            <option value="">Select...</option>
            {#if addTargetType === 'user'}
              {#each users as user (user.id)}
                <option value={user.id}>{user.full_name || user.username}</option>
              {/each}
            {:else if addTargetType === 'team'}
              {#each teams as team (team.id)}
                <option value={team.id}>{team.name}</option>
              {/each}
            {:else}
              {#each schedules as schedule (schedule.id)}
                <option value={schedule.id}>{schedule.name}</option>
              {/each}
            {/if}
          </select>
        </div>

        <!-- Notification Channels Override -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">
            Notification Channels (optional override)
          </label>
          <p class="text-xs text-gray-500 mb-2">Leave unchecked to use all available channels</p>
          <div class="space-y-2">
            {#each ['email', 'slack', 'sms', 'webhook', 'msteams'] as channelType}
              <label class="flex items-center gap-2">
                <input
                  type="checkbox"
                  checked={addTargetChannels.includes(channelType)}
                  on:change={() => toggleChannel(channelType)}
                  class="rounded border-gray-300 text-primary-600 focus:ring-primary-500"
                />
                <span class="text-sm">{getChannelTypeDisplay(channelType)}</span>
              </label>
            {/each}
          </div>
        </div>

        <!-- Urgent Flag -->
        {#if addTargetChannels.length > 0}
          <label class="flex items-center gap-2">
            <input
              type="checkbox"
              bind:checked={addTargetUrgent}
              class="rounded border-gray-300 text-primary-600 focus:ring-primary-500"
            />
            <span class="text-sm font-medium">Mark as urgent (high-priority notification)</span>
          </label>
        {/if}
      </div>

      <!-- Modal Actions -->
      <div class="flex justify-end gap-2 mt-6">
        <Button variant="secondary" on:click={closeAddTargetModal}>Cancel</Button>
        <Button
          variant="primary"
          on:click={handleAddTargetWithChannels}
          disabled={!addTargetId || addingTarget}
        >
          {addingTarget ? 'Adding...' : 'Add Target'}
        </Button>
      </div>
    </div>
  </div>
{/if}
