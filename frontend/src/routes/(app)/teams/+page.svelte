<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { teamsStore } from '$lib/stores/teams';
	import Button from '$lib/components/ui/Button.svelte';
	import Input from '$lib/components/ui/Input.svelte';

	let showCreateForm = false;
	let name = '';
	let description = '';
	let createError = '';
	let creatingTeam = false;

	onMount(() => {
		teamsStore.load();
	});

	async function handleCreateTeam() {
		createError = '';
		creatingTeam = true;

		try {
			await teamsStore.create({
				name,
				description: description || undefined
			});

			// Reset form
			name = '';
			description = '';
			showCreateForm = false;
		} catch (err) {
			createError = err instanceof Error ? err.message : 'Failed to create team';
		} finally {
			creatingTeam = false;
		}
	}

	async function handleDeleteTeam(id: string, teamName: string) {
		if (!confirm(`Are you sure you want to delete team "${teamName}"?`)) {
			return;
		}

		try {
			await teamsStore.delete(id);
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to delete team');
		}
	}
</script>

<svelte:head>
	<title>Teams - Pulsar</title>
</svelte:head>

<div class="space-y-6">
	<div class="flex justify-between items-center">
		<div>
			<h2 class="text-3xl font-bold text-gray-900">Teams</h2>
			<p class="mt-2 text-gray-600">Manage your organization's teams</p>
		</div>
		<Button variant="primary" on:click={() => (showCreateForm = !showCreateForm)}>
			{showCreateForm ? 'Cancel' : 'Create Team'}
		</Button>
	</div>

	{#if showCreateForm}
		<div class="bg-white p-6 rounded-lg shadow">
			<h3 class="text-lg font-semibold mb-4">Create New Team</h3>
			<form on:submit|preventDefault={handleCreateTeam} class="space-y-4">
				<Input
					id="name"
					label="Team Name"
					bind:value={name}
					placeholder="Engineering, DevOps, Support..."
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
						placeholder="Team description..."
					></textarea>
				</div>

				{#if createError}
					<div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg">
						{createError}
					</div>
				{/if}

				<div class="flex gap-2">
					<Button type="submit" variant="primary" disabled={creatingTeam}>
						{creatingTeam ? 'Creating...' : 'Create Team'}
					</Button>
					<Button type="button" variant="secondary" on:click={() => (showCreateForm = false)}>
						Cancel
					</Button>
				</div>
			</form>
		</div>
	{/if}

	<!-- Teams List -->
	<div class="space-y-4">
		{#if $teamsStore.isLoading}
			<div class="text-center py-12">
				<div
					class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600"
				></div>
				<p class="mt-2 text-gray-600">Loading teams...</p>
			</div>
		{:else if $teamsStore.error}
			<div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg">
				{$teamsStore.error}
			</div>
		{:else if $teamsStore.teams.length === 0}
			<div class="text-center py-12 bg-white rounded-lg shadow">
				<p class="text-gray-600">No teams found</p>
				<p class="text-sm text-gray-500 mt-2">Create your first team to get started</p>
			</div>
		{:else}
			<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
				{#each $teamsStore.teams as team (team.id)}
					<div class="bg-white p-6 rounded-lg shadow hover:shadow-md transition-shadow">
						<div class="flex justify-between items-start mb-3">
							<h3 class="text-lg font-semibold text-gray-900">{team.name}</h3>
							<button
								class="text-red-600 hover:text-red-800 text-sm"
								on:click={() => handleDeleteTeam(team.id, team.name)}
							>
								Delete
							</button>
						</div>

						{#if team.description}
							<p class="text-sm text-gray-600 mb-4">{team.description}</p>
						{/if}

						<Button variant="secondary" fullWidth on:click={() => goto(`/teams/${team.id}`)}>
							View Team
						</Button>
					</div>
				{/each}
			</div>
		{/if}
	</div>
</div>
