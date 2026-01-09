<script lang="ts">
  import { onMount } from 'svelte';
  import { page } from '$app/stores';
  import { goto } from '$app/navigation';
  import { api } from '$lib/api/client';
  import type { EscalationPolicyWithRules, EscalationTargetType } from '$lib/types/escalation';
  import type { User } from '$lib/types/user';
  import type { Team } from '$lib/types/team';
  import type { Schedule } from '$lib/types/schedule';
  import Button from '$lib/components/ui/Button.svelte';
  import Input from '$lib/components/ui/Input.svelte';

  let policyId = $page.params.id;
  let policy: EscalationPolicyWithRules | null = null;
  let users: User[] = [];
  let teams: Team[] = [];
  let schedules: Schedule[] = [];
  let isLoading = true;
  let error = '';

  // Create rule form
  let showCreateRuleForm = false;
  let rulePosition = 0;
  let escalationDelay = 15;
  let ruleError = '';
  let creatingRule = false;

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
      const [usersResp, teamsResp, schedulesResp] = await Promise.all([
        api.listUsers(),
        api.listTeams(),
        api.listSchedules(),
      ]);
      users = usersResp.users;
      teams = teamsResp.teams;
      schedules = schedulesResp.schedules;
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
      rulePosition = policy ? policy.rules.length : 0;
      escalationDelay = 15;
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

  async function handleAddTarget(
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
      <h3 class="text-lg font-semibold mb-2">Policy Settings</h3>
      <div class="text-sm text-gray-600 space-y-1">
        <p>
          Repeat: {policy.repeat_enabled
            ? policy.repeat_count
              ? `${policy.repeat_count} times`
              : 'Infinite'
            : 'Disabled'}
        </p>
      </div>
    </div>

    <!-- Escalation Rules -->
    <div class="bg-white p-6 rounded-lg shadow">
      <div class="flex justify-between items-center mb-4">
        <h3 class="text-lg font-semibold">Escalation Rules ({policy.rules?.length || 0})</h3>
        <Button
          variant="primary"
          size="sm"
          on:click={() => {
            rulePosition = policy?.rules?.length || 0;
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
            <Input
              id="position"
              label="Position (order)"
              type="number"
              bind:value={rulePosition}
              min="0"
              required
            />

            <Input
              id="escalation-delay"
              label="Escalation Delay (minutes)"
              type="number"
              bind:value={escalationDelay}
              min="0"
              required
            />

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
                <div>
                  <span class="text-sm font-semibold text-gray-900">
                    Level {rule.position + 1}
                  </span>
                  <p class="text-sm text-gray-600 mt-1">
                    Escalate after {rule.escalation_delay} minutes
                  </p>
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
                        <span>{getTargetDisplay(target.target_type, target.target_id)}</span>
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
                        handleAddTarget(rule.id, 'user', value);
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
                        handleAddTarget(rule.id, 'team', value);
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
                        handleAddTarget(rule.id, 'schedule', value);
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
