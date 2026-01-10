<script lang="ts">
  import { onMount } from 'svelte';
  import { api } from '$lib/api/client';
  import Button from '$lib/components/ui/Button.svelte';
  import Input from '$lib/components/ui/Input.svelte';
  import type {
    AlertRoutingRule,
    RoutingCondition,
    RoutingConditions,
    RoutingActions,
    CreateRoutingRuleRequest,
  } from '$lib/types/routing';
  import { ROUTING_FIELDS, ROUTING_OPERATORS, PRIORITY_OPTIONS } from '$lib/types/routing';
  import type { Team } from '$lib/types/team';
  import type { EscalationPolicy } from '$lib/types/escalation';
  import type { User } from '$lib/types/user';

  let rules: AlertRoutingRule[] = [];
  let teams: Team[] = [];
  let policies: EscalationPolicy[] = [];
  let users: User[] = [];
  let isLoading = true;
  let error = '';

  // Form state
  let showCreateForm = false;
  let editingRule: AlertRoutingRule | null = null;
  let formError = '';
  let isSaving = false;

  // Form fields
  let name = '';
  let description = '';
  let matchType: 'all' | 'any' = 'all';
  let conditions: RoutingCondition[] = [{ field: 'source', operator: 'equals', value: '' }];
  let actionTeamId = '';
  let actionUserId = '';
  let actionPolicyId = '';
  let actionPriority = '';
  let actionTags = '';
  let actionSuppress = false;
  let enabled = true;

  onMount(async () => {
    await Promise.all([loadRules(), loadTeams(), loadPolicies(), loadUsers()]);
  });

  async function loadRules() {
    isLoading = true;
    error = '';
    try {
      const response = await api.listRoutingRules();
      rules = response.rules || [];
    } catch (err) {
      error = err instanceof Error ? err.message : 'Failed to load routing rules';
    } finally {
      isLoading = false;
    }
  }

  async function loadTeams() {
    try {
      const response = await api.listTeams();
      teams = response.teams || [];
    } catch (err) {
      console.error('Failed to load teams:', err);
    }
  }

  async function loadPolicies() {
    try {
      const response = await api.listEscalationPolicies();
      policies = response.policies || [];
    } catch (err) {
      console.error('Failed to load escalation policies:', err);
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

  function resetForm() {
    name = '';
    description = '';
    matchType = 'all';
    conditions = [{ field: 'source', operator: 'equals', value: '' }];
    actionTeamId = '';
    actionUserId = '';
    actionPolicyId = '';
    actionPriority = '';
    actionTags = '';
    actionSuppress = false;
    enabled = true;
    formError = '';
    editingRule = null;
    showCreateForm = false;
  }

  function editRule(rule: AlertRoutingRule) {
    editingRule = rule;
    showCreateForm = true;
    name = rule.name;
    description = rule.description || '';
    enabled = rule.enabled;

    // Parse conditions
    const parsedConditions = rule.conditions;
    matchType = parsedConditions.match || 'all';
    conditions =
      parsedConditions.conditions?.length > 0
        ? parsedConditions.conditions
        : [{ field: 'source', operator: 'equals', value: '' }];

    // Parse actions
    const actions = rule.actions;
    actionTeamId = actions.assign_team_id || '';
    actionUserId = actions.assign_user_id || '';
    actionPolicyId = actions.assign_escalation_policy_id || '';
    actionPriority = actions.set_priority || '';
    actionTags = actions.add_tags?.join(', ') || '';
    actionSuppress = actions.suppress || false;
  }

  function addCondition() {
    conditions = [...conditions, { field: 'source', operator: 'equals', value: '' }];
  }

  function removeCondition(index: number) {
    if (conditions.length > 1) {
      conditions = conditions.filter((_, i) => i !== index);
    }
  }

  async function handleSubmit() {
    if (!name.trim()) {
      formError = 'Name is required';
      return;
    }

    if (conditions.some((c) => !c.value.trim())) {
      formError = 'All condition values are required';
      return;
    }

    formError = '';
    isSaving = true;

    const conditionsData: RoutingConditions = {
      match: matchType,
      conditions: conditions,
    };

    const actionsData: RoutingActions = {
      suppress: actionSuppress,
    };

    if (actionTeamId) actionsData.assign_team_id = actionTeamId;
    if (actionUserId) actionsData.assign_user_id = actionUserId;
    if (actionPolicyId) actionsData.assign_escalation_policy_id = actionPolicyId;
    if (actionPriority) actionsData.set_priority = actionPriority;
    if (actionTags.trim()) {
      actionsData.add_tags = actionTags
        .split(',')
        .map((t) => t.trim())
        .filter((t) => t);
    }

    try {
      if (editingRule) {
        await api.updateRoutingRule(editingRule.id, {
          name: name.trim(),
          description: description.trim() || undefined,
          conditions: conditionsData,
          actions: actionsData,
          enabled,
        });
      } else {
        const request: CreateRoutingRuleRequest = {
          name: name.trim(),
          description: description.trim() || undefined,
          conditions: conditionsData,
          actions: actionsData,
          enabled,
        };
        await api.createRoutingRule(request);
      }

      resetForm();
      await loadRules();
    } catch (err) {
      formError = err instanceof Error ? err.message : 'Failed to save routing rule';
    } finally {
      isSaving = false;
    }
  }

  async function toggleEnabled(rule: AlertRoutingRule) {
    try {
      await api.updateRoutingRule(rule.id, { enabled: !rule.enabled });
      await loadRules();
    } catch (err) {
      error = err instanceof Error ? err.message : 'Failed to update rule';
    }
  }

  async function deleteRule(rule: AlertRoutingRule) {
    if (!confirm(`Are you sure you want to delete "${rule.name}"?`)) {
      return;
    }

    try {
      await api.deleteRoutingRule(rule.id);
      await loadRules();
    } catch (err) {
      error = err instanceof Error ? err.message : 'Failed to delete rule';
    }
  }

  function getTeamName(id: string): string {
    return teams.find((t) => t.id === id)?.name || 'Unknown Team';
  }

  function getUserName(id: string): string {
    const user = users.find((u) => u.id === id);
    return user?.full_name || user?.username || 'Unknown User';
  }

  function getPolicyName(id: string): string {
    return policies.find((p) => p.id === id)?.name || 'Unknown Policy';
  }

  function formatConditions(rule: AlertRoutingRule): string {
    const conds = rule.conditions;
    if (!conds.conditions || conds.conditions.length === 0) return 'No conditions';

    const parts = conds.conditions.map((c) => `${c.field} ${c.operator} "${c.value}"`);
    return parts.join(conds.match === 'all' ? ' AND ' : ' OR ');
  }

  function formatActions(rule: AlertRoutingRule): string[] {
    const actions: string[] = [];
    const a = rule.actions;

    if (a.assign_team_id) actions.push(`Assign to team: ${getTeamName(a.assign_team_id)}`);
    if (a.assign_user_id) actions.push(`Assign to user: ${getUserName(a.assign_user_id)}`);
    if (a.assign_escalation_policy_id)
      actions.push(`Use policy: ${getPolicyName(a.assign_escalation_policy_id)}`);
    if (a.set_priority) actions.push(`Set priority: ${a.set_priority}`);
    if (a.add_tags?.length) actions.push(`Add tags: ${a.add_tags.join(', ')}`);
    if (a.suppress) actions.push('Suppress alert');

    return actions.length > 0 ? actions : ['No actions configured'];
  }
</script>

<svelte:head>
  <title>Routing Rules - Pulsar</title>
</svelte:head>

<div class="space-y-6">
  <div class="flex justify-between items-center">
    <div>
      <h2 class="text-3xl font-bold text-gray-900">Alert Routing Rules</h2>
      <p class="mt-2 text-gray-500">
        Automatically route alerts to the right team or escalation policy based on conditions
      </p>
    </div>
    <Button
      variant="primary"
      on:click={() => {
        if (showCreateForm) {
          resetForm();
        } else {
          showCreateForm = true;
        }
      }}
    >
      {showCreateForm ? 'Cancel' : 'Create Rule'}
    </Button>
  </div>

  {#if error}
    <div class="bg-red-50 border border-red-200 text-red-600 px-4 py-3 rounded-lg">
      {error}
    </div>
  {/if}

  <!-- Create/Edit Form -->
  {#if showCreateForm}
    <div class="bg-white backdrop-blur-sm p-6 rounded-xl border border-primary-200 shadow-sm">
      <h3 class="text-lg font-semibold mb-4 text-gray-900">
        {editingRule ? 'Edit Routing Rule' : 'Create New Routing Rule'}
      </h3>
      <form on:submit|preventDefault={handleSubmit} class="space-y-6">
        <!-- Basic Info -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <Input
            id="name"
            label="Rule Name"
            bind:value={name}
            placeholder="e.g., Route production alerts"
            required
          />
          <Input
            id="description"
            label="Description (optional)"
            bind:value={description}
            placeholder="Brief description of this rule"
          />
        </div>

        <!-- Conditions -->
        <div class="space-y-4">
          <div class="flex items-center justify-between">
            <h4 class="text-sm font-medium text-gray-700">Conditions</h4>
            <div class="flex items-center gap-4">
              <label class="flex items-center gap-2">
                <input type="radio" bind:group={matchType} value="all" class="text-primary-600" />
                <span class="text-sm text-gray-600">Match ALL conditions</span>
              </label>
              <label class="flex items-center gap-2">
                <input type="radio" bind:group={matchType} value="any" class="text-primary-600" />
                <span class="text-sm text-gray-600">Match ANY condition</span>
              </label>
            </div>
          </div>

          {#each conditions as condition, index}
            <div class="flex items-center gap-2 bg-gray-50 p-3 rounded-lg">
              <select
                bind:value={condition.field}
                class="px-3 py-2 bg-white border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 text-gray-900"
              >
                {#each ROUTING_FIELDS as field}
                  <option value={field.value}>{field.label}</option>
                {/each}
              </select>

              <select
                bind:value={condition.operator}
                class="px-3 py-2 bg-white border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 text-gray-900"
              >
                {#each ROUTING_OPERATORS as op}
                  <option value={op.value}>{op.label}</option>
                {/each}
              </select>

              <input
                type="text"
                bind:value={condition.value}
                placeholder="Value"
                class="flex-1 px-3 py-2 bg-white border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 text-gray-900"
              />

              {#if conditions.length > 1}
                <button
                  type="button"
                  on:click={() => removeCondition(index)}
                  class="p-2 text-red-500 hover:bg-red-50 rounded"
                >
                  <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      stroke-width="2"
                      d="M6 18L18 6M6 6l12 12"
                    />
                  </svg>
                </button>
              {/if}
            </div>
          {/each}

          <Button type="button" variant="secondary" size="sm" on:click={addCondition}>
            + Add Condition
          </Button>
        </div>

        <!-- Actions -->
        <div class="space-y-4">
          <h4 class="text-sm font-medium text-gray-700">Actions (when conditions match)</h4>

          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label for="actionTeam" class="block text-sm font-medium text-gray-600 mb-1">
                Assign to Team
              </label>
              <select
                id="actionTeam"
                bind:value={actionTeamId}
                class="w-full px-3 py-2 bg-white border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 text-gray-900"
              >
                <option value="">-- None --</option>
                {#each teams as team}
                  <option value={team.id}>{team.name}</option>
                {/each}
              </select>
            </div>

            <div>
              <label for="actionUser" class="block text-sm font-medium text-gray-600 mb-1">
                Assign to User
              </label>
              <select
                id="actionUser"
                bind:value={actionUserId}
                class="w-full px-3 py-2 bg-white border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 text-gray-900"
              >
                <option value="">-- None --</option>
                {#each users as user}
                  <option value={user.id}>{user.full_name || user.username}</option>
                {/each}
              </select>
            </div>

            <div>
              <label for="actionPolicy" class="block text-sm font-medium text-gray-600 mb-1">
                Apply Escalation Policy
              </label>
              <select
                id="actionPolicy"
                bind:value={actionPolicyId}
                class="w-full px-3 py-2 bg-white border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 text-gray-900"
              >
                <option value="">-- None --</option>
                {#each policies as policy}
                  <option value={policy.id}>{policy.name}</option>
                {/each}
              </select>
            </div>

            <div>
              <label for="actionPriority" class="block text-sm font-medium text-gray-600 mb-1">
                Set Priority
              </label>
              <select
                id="actionPriority"
                bind:value={actionPriority}
                class="w-full px-3 py-2 bg-white border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 text-gray-900"
              >
                <option value="">-- Don't change --</option>
                {#each PRIORITY_OPTIONS as priority}
                  <option value={priority.value}>{priority.label}</option>
                {/each}
              </select>
            </div>
          </div>

          <Input
            id="actionTags"
            label="Add Tags (comma-separated)"
            bind:value={actionTags}
            placeholder="e.g., auto-routed, production"
          />

          <label class="flex items-center gap-2">
            <input type="checkbox" bind:checked={actionSuppress} class="rounded text-primary-600" />
            <span class="text-sm text-gray-600">Suppress this alert (don't notify)</span>
          </label>
        </div>

        <!-- Enabled toggle -->
        <label class="flex items-center gap-2">
          <input type="checkbox" bind:checked={enabled} class="rounded text-primary-600" />
          <span class="text-sm text-gray-600">Rule is enabled</span>
        </label>

        {#if formError}
          <div class="bg-red-50 border border-red-200 text-red-600 px-4 py-3 rounded-lg">
            {formError}
          </div>
        {/if}

        <div class="flex gap-2">
          <Button type="submit" variant="primary" disabled={isSaving}>
            {isSaving ? 'Saving...' : editingRule ? 'Update Rule' : 'Create Rule'}
          </Button>
          <Button type="button" variant="secondary" on:click={resetForm}>Cancel</Button>
        </div>
      </form>
    </div>
  {/if}

  <!-- Rules List -->
  {#if isLoading}
    <div class="text-center py-12">
      <div
        class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-primary-500"
      ></div>
      <p class="mt-2 text-gray-500">Loading routing rules...</p>
    </div>
  {:else if rules.length === 0}
    <div
      class="text-center py-12 bg-white backdrop-blur-sm rounded-xl border border-gray-200 shadow-sm"
    >
      <svg
        class="mx-auto h-12 w-12 text-gray-400"
        fill="none"
        stroke="currentColor"
        viewBox="0 0 24 24"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"
        />
      </svg>
      <h3 class="mt-2 text-sm font-medium text-gray-900">No routing rules</h3>
      <p class="mt-1 text-sm text-gray-500">Get started by creating a new routing rule.</p>
      <div class="mt-6">
        <Button variant="primary" on:click={() => (showCreateForm = true)}>
          Create Routing Rule
        </Button>
      </div>
    </div>
  {:else}
    <div class="space-y-4">
      <p class="text-sm text-gray-500">
        Rules are evaluated in priority order (lower number = higher priority). First matching rule
        wins.
      </p>

      {#each rules as rule, index}
        <div
          class="bg-white rounded-xl border shadow-sm overflow-hidden {rule.enabled
            ? 'border-gray-200'
            : 'border-gray-300 opacity-60'}"
        >
          <div class="px-6 py-4">
            <div class="flex items-start justify-between">
              <div class="flex-1">
                <div class="flex items-center gap-3 mb-2">
                  <span class="text-sm font-mono text-gray-400">#{rule.priority}</span>
                  <h3 class="text-lg font-semibold text-gray-900">{rule.name}</h3>
                  {#if rule.enabled}
                    <span
                      class="px-2 py-0.5 text-xs font-medium rounded bg-green-100 text-green-700"
                    >
                      Enabled
                    </span>
                  {:else}
                    <span class="px-2 py-0.5 text-xs font-medium rounded bg-gray-100 text-gray-500">
                      Disabled
                    </span>
                  {/if}
                </div>

                {#if rule.description}
                  <p class="text-sm text-gray-500 mb-3">{rule.description}</p>
                {/if}

                <div class="flex flex-col gap-2 text-sm">
                  <div class="flex items-start gap-2">
                    <span class="font-medium text-gray-600 min-w-[80px]">When:</span>
                    <span class="text-gray-800 font-mono text-xs bg-gray-50 px-2 py-1 rounded">
                      {formatConditions(rule)}
                    </span>
                  </div>
                  <div class="flex items-start gap-2">
                    <span class="font-medium text-gray-600 min-w-[80px]">Then:</span>
                    <div class="flex flex-wrap gap-1">
                      {#each formatActions(rule) as action}
                        <span class="bg-primary-50 text-primary-700 px-2 py-1 rounded text-xs">
                          {action}
                        </span>
                      {/each}
                    </div>
                  </div>
                </div>
              </div>

              <div class="flex items-center gap-2">
                <button
                  on:click={() => toggleEnabled(rule)}
                  class="p-2 rounded hover:bg-gray-100"
                  title={rule.enabled ? 'Disable' : 'Enable'}
                >
                  {#if rule.enabled}
                    <svg
                      class="w-5 h-5 text-green-500"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                    >
                      <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
                      />
                      <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"
                      />
                    </svg>
                  {:else}
                    <svg
                      class="w-5 h-5 text-gray-400"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                    >
                      <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21"
                      />
                    </svg>
                  {/if}
                </button>
                <button
                  on:click={() => editRule(rule)}
                  class="p-2 rounded hover:bg-gray-100 text-primary-600"
                  title="Edit"
                >
                  <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      stroke-width="2"
                      d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"
                    />
                  </svg>
                </button>
                <button
                  on:click={() => deleteRule(rule)}
                  class="p-2 rounded hover:bg-red-50 text-red-500"
                  title="Delete"
                >
                  <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      stroke-width="2"
                      d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                    />
                  </svg>
                </button>
              </div>
            </div>
          </div>
        </div>
      {/each}
    </div>
  {/if}

  <!-- Info Box -->
  <div class="bg-blue-50 border border-blue-200 p-4 rounded-xl">
    <h4 class="text-sm font-semibold text-blue-800 mb-2">How Routing Rules Work</h4>
    <ul class="text-sm text-blue-700 space-y-1 list-disc list-inside">
      <li>Rules are evaluated in priority order when a new alert is created</li>
      <li>The first matching rule's actions are applied to the alert</li>
      <li>If no rules match, the alert follows default routing (no automatic assignment)</li>
      <li>Use conditions to match alerts by source, priority, tags, or message content</li>
    </ul>
  </div>
</div>
