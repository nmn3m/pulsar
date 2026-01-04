<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { escalationPoliciesStore } from '$lib/stores/escalations';
	import Button from '$lib/components/ui/Button.svelte';
	import Input from '$lib/components/ui/Input.svelte';

	let showCreateForm = false;
	let name = '';
	let description = '';
	let repeatEnabled = false;
	let repeatCount: number | undefined = undefined;
	let createError = '';
	let creating = false;

	onMount(() => {
		escalationPoliciesStore.load();
	});

	async function handleCreatePolicy() {
		createError = '';
		creating = true;

		try {
			await escalationPoliciesStore.create({
				name,
				description: description || undefined,
				repeat_enabled: repeatEnabled,
				repeat_count: repeatEnabled && repeatCount ? repeatCount : undefined
			});

			// Reset form
			name = '';
			description = '';
			repeatEnabled = false;
			repeatCount = undefined;
			showCreateForm = false;
		} catch (err) {
			createError = err instanceof Error ? err.message : 'Failed to create escalation policy';
		} finally {
			creating = false;
		}
	}

	async function handleDeletePolicy(id: string, policyName: string) {
		if (!confirm(`Are you sure you want to delete escalation policy "${policyName}"?`)) {
			return;
		}

		try {
			await escalationPoliciesStore.delete(id);
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to delete escalation policy');
		}
	}
</script>

<svelte:head>
	<title>Escalation Policies - Pulsar</title>
</svelte:head>

<div class="space-y-6">
	<div class="flex justify-between items-center">
		<div>
			<h2 class="text-3xl font-bold text-gray-900">Escalation Policies</h2>
			<p class="mt-2 text-gray-600">
				Define how alerts escalate through different notification levels
			</p>
		</div>
		<Button variant="primary" on:click={() => (showCreateForm = !showCreateForm)}>
			{showCreateForm ? 'Cancel' : 'Create Policy'}
		</Button>
	</div>

	{#if showCreateForm}
		<div class="bg-white p-6 rounded-lg shadow">
			<h3 class="text-lg font-semibold mb-4">Create Escalation Policy</h3>
			<form on:submit|preventDefault={handleCreatePolicy} class="space-y-4">
				<Input
					id="name"
					label="Policy Name"
					bind:value={name}
					placeholder="Primary Escalation, Weekend Escalation..."
					required
				/>

				<div>
					<label for="description" class="block text-sm font-medium text-gray-700 mb-1">
						Description
					</label>
					<textarea
						id="description"
						bind:value={description}
						rows="3"
						class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500"
						placeholder="Policy description..."
					></textarea>
				</div>

				<div class="flex items-center gap-2">
					<input
						id="repeat-enabled"
						type="checkbox"
						bind:checked={repeatEnabled}
						class="rounded border-gray-300 text-primary-600 focus:ring-primary-500"
					/>
					<label for="repeat-enabled" class="text-sm font-medium text-gray-700">
						Enable repeat escalation
					</label>
				</div>

				{#if repeatEnabled}
					<Input
						id="repeat-count"
						label="Maximum Repeat Count (leave empty for infinite)"
						type="number"
						bind:value={repeatCount}
						min="1"
						placeholder="Leave empty for infinite repeats"
					/>
				{/if}

				{#if createError}
					<div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg">
						{createError}
					</div>
				{/if}

				<div class="flex gap-2">
					<Button type="submit" variant="primary" disabled={creating}>
						{creating ? 'Creating...' : 'Create Policy'}
					</Button>
					<Button type="button" variant="secondary" on:click={() => (showCreateForm = false)}>
						Cancel
					</Button>
				</div>
			</form>
		</div>
	{/if}

	{#if $escalationPoliciesStore.isLoading}
		<div class="text-center py-12">
			<div
				class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600"
			></div>
			<p class="mt-2 text-gray-600">Loading escalation policies...</p>
		</div>
	{:else if $escalationPoliciesStore.error}
		<div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg">
			{$escalationPoliciesStore.error}
		</div>
	{:else if $escalationPoliciesStore.policies.length === 0}
		<div class="text-center py-12 bg-white rounded-lg shadow">
			<p class="text-gray-600">No escalation policies found</p>
			<p class="text-sm text-gray-500 mt-2">Create your first policy to get started</p>
		</div>
	{:else}
		<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
			{#each $escalationPoliciesStore.policies as policy (policy.id)}
				<div class="bg-white p-6 rounded-lg shadow hover:shadow-md transition-shadow">
					<div class="flex justify-between items-start mb-4">
						<div class="flex-1">
							<h3 class="text-lg font-semibold text-gray-900">{policy.name}</h3>
							{#if policy.description}
								<p class="text-sm text-gray-600 mt-1">{policy.description}</p>
							{/if}
						</div>
					</div>

					<div class="space-y-2 text-sm text-gray-600">
						{#if policy.repeat_enabled}
							<div class="flex items-center gap-2">
								<span class="px-2 py-1 bg-blue-100 text-blue-800 rounded text-xs">
									Repeat: {policy.repeat_count ? `${policy.repeat_count}x` : 'Infinite'}
								</span>
							</div>
						{/if}
					</div>

					<div class="flex gap-2 mt-4">
						<Button
							variant="primary"
							size="sm"
							on:click={() => goto(`/escalation-policies/${policy.id}`)}
						>
							Manage Rules
						</Button>
						<Button
							variant="danger"
							size="sm"
							on:click={() => handleDeletePolicy(policy.id, policy.name)}
						>
							Delete
						</Button>
					</div>
				</div>
			{/each}
		</div>
	{/if}
</div>
