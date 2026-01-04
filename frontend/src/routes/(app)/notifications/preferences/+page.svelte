<script lang="ts">
	import { onMount } from 'svelte';
	import {
		notificationChannelsStore,
		userNotificationPreferencesStore
	} from '$lib/stores/notifications';
	import Button from '$lib/components/ui/Button.svelte';
	import Input from '$lib/components/ui/Input.svelte';

	let showCreateForm = false;
	let selectedChannelId = '';
	let isEnabled = true;
	let dndEnabled = false;
	let dndStartTime = '22:00:00';
	let dndEndTime = '08:00:00';
	let minPriority = '';

	let createError = '';
	let creating = false;

	onMount(async () => {
		await Promise.all([
			notificationChannelsStore.load(),
			userNotificationPreferencesStore.load()
		]);
	});

	async function handleCreatePreference() {
		createError = '';
		creating = true;

		try {
			await userNotificationPreferencesStore.create({
				channel_id: selectedChannelId,
				is_enabled: isEnabled,
				dnd_enabled: dndEnabled,
				dnd_start_time: dndEnabled ? dndStartTime : undefined,
				dnd_end_time: dndEnabled ? dndEndTime : undefined,
				min_priority: minPriority || undefined
			});

			// Reset form
			selectedChannelId = '';
			isEnabled = true;
			dndEnabled = false;
			dndStartTime = '22:00:00';
			dndEndTime = '08:00:00';
			minPriority = '';
			showCreateForm = false;
		} catch (err) {
			createError = err instanceof Error ? err.message : 'Failed to create notification preference';
		} finally {
			creating = false;
		}
	}

	async function handleTogglePreference(id: string, currentlyEnabled: boolean) {
		try {
			await userNotificationPreferencesStore.update(id, {
				is_enabled: !currentlyEnabled
			});
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to update preference');
		}
	}

	async function handleDeletePreference(id: string) {
		if (!confirm('Delete this notification preference?')) return;

		try {
			await userNotificationPreferencesStore.delete(id);
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to delete preference');
		}
	}

	function getChannelName(channelId: string): string {
		const channel = $notificationChannelsStore.channels.find((c) => c.id === channelId);
		return channel ? channel.name : 'Unknown Channel';
	}

	function getChannelType(channelId: string): string {
		const channel = $notificationChannelsStore.channels.find((c) => c.id === channelId);
		return channel ? channel.channel_type : '';
	}

	function getAvailableChannels() {
		const preferenceChannelIds = new Set(
			$userNotificationPreferencesStore.preferences.map((p) => p.channel_id)
		);
		return $notificationChannelsStore.channels.filter((c) => !preferenceChannelIds.has(c.id));
	}
</script>

<svelte:head>
	<title>Notification Preferences - Pulsar</title>
</svelte:head>

<div class="space-y-6">
	<div class="flex justify-between items-center">
		<div>
			<h2 class="text-3xl font-bold text-gray-900">Notification Preferences</h2>
			<p class="mt-2 text-gray-600">Configure how you receive notifications for each channel</p>
		</div>
		<Button variant="primary" on:click={() => (showCreateForm = !showCreateForm)}>
			{showCreateForm ? 'Cancel' : 'Add Preference'}
		</Button>
	</div>

	{#if showCreateForm}
		<div class="bg-white p-6 rounded-lg shadow">
			<h3 class="text-lg font-semibold mb-4">Create Notification Preference</h3>
			<form on:submit|preventDefault={handleCreatePreference} class="space-y-4">
				<div>
					<label for="channel" class="block text-sm font-medium text-gray-700 mb-1">
						Notification Channel
					</label>
					<select
						id="channel"
						bind:value={selectedChannelId}
						class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500"
						required
					>
						<option value="">Select a channel...</option>
						{#each getAvailableChannels() as channel (channel.id)}
							<option value={channel.id}>{channel.name} ({channel.channel_type})</option>
						{/each}
					</select>
				</div>

				<div class="flex items-center gap-2">
					<input
						id="is-enabled"
						type="checkbox"
						bind:checked={isEnabled}
						class="rounded border-gray-300 text-primary-600 focus:ring-primary-500"
					/>
					<label for="is-enabled" class="text-sm font-medium text-gray-700">
						Enable notifications for this channel
					</label>
				</div>

				<div class="flex items-center gap-2">
					<input
						id="dnd-enabled"
						type="checkbox"
						bind:checked={dndEnabled}
						class="rounded border-gray-300 text-primary-600 focus:ring-primary-500"
					/>
					<label for="dnd-enabled" class="text-sm font-medium text-gray-700">
						Enable Do Not Disturb
					</label>
				</div>

				{#if dndEnabled}
					<div class="grid grid-cols-2 gap-4 p-4 bg-gray-50 rounded-lg">
						<div>
							<label for="dnd-start" class="block text-sm font-medium text-gray-700 mb-1">
								DND Start Time
							</label>
							<input
								id="dnd-start"
								type="time"
								bind:value={dndStartTime}
								class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500"
								step="1"
							/>
						</div>
						<div>
							<label for="dnd-end" class="block text-sm font-medium text-gray-700 mb-1">
								DND End Time
							</label>
							<input
								id="dnd-end"
								type="time"
								bind:value={dndEndTime}
								class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500"
								step="1"
							/>
						</div>
					</div>
				{/if}

				<div>
					<label for="min-priority" class="block text-sm font-medium text-gray-700 mb-1">
						Minimum Priority (optional)
					</label>
					<select
						id="min-priority"
						bind:value={minPriority}
						class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500"
					>
						<option value="">All priorities</option>
						<option value="P1">P1 only</option>
						<option value="P2">P2 and above</option>
						<option value="P3">P3 and above</option>
						<option value="P4">P4 and above</option>
						<option value="P5">P5 and above</option>
					</select>
					<p class="mt-1 text-xs text-gray-500">
						Only receive notifications for alerts at or above this priority level
					</p>
				</div>

				{#if createError}
					<div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg">
						{createError}
					</div>
				{/if}

				<div class="flex gap-2">
					<Button type="submit" variant="primary" disabled={creating}>
						{creating ? 'Creating...' : 'Create Preference'}
					</Button>
					<Button type="button" variant="secondary" on:click={() => (showCreateForm = false)}>
						Cancel
					</Button>
				</div>
			</form>
		</div>
	{/if}

	{#if $userNotificationPreferencesStore.isLoading || $notificationChannelsStore.isLoading}
		<div class="text-center py-12">
			<div
				class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600"
			></div>
			<p class="mt-2 text-gray-600">Loading preferences...</p>
		</div>
	{:else if $userNotificationPreferencesStore.error}
		<div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg">
			{$userNotificationPreferencesStore.error}
		</div>
	{:else if $userNotificationPreferencesStore.preferences.length === 0}
		<div class="text-center py-12 bg-white rounded-lg shadow">
			<p class="text-gray-600">No notification preferences configured</p>
			<p class="text-sm text-gray-500 mt-2">
				Add preferences to customize your notification settings
			</p>
		</div>
	{:else}
		<div class="space-y-4">
			{#each $userNotificationPreferencesStore.preferences as pref (pref.id)}
				<div class="bg-white p-6 rounded-lg shadow">
					<div class="flex justify-between items-start">
						<div class="flex-1">
							<div class="flex items-center gap-3 mb-2">
								<h3 class="text-lg font-semibold text-gray-900">
									{getChannelName(pref.channel_id)}
								</h3>
								<span
									class="px-2 py-1 rounded text-xs font-medium {pref.is_enabled
										? 'bg-green-100 text-green-800'
										: 'bg-gray-100 text-gray-800'}"
								>
									{pref.is_enabled ? 'Enabled' : 'Disabled'}
								</span>
								<span class="px-2 py-1 bg-blue-100 text-blue-800 rounded text-xs font-medium">
									{getChannelType(pref.channel_id)}
								</span>
							</div>

							<div class="space-y-2 text-sm text-gray-600">
								{#if pref.dnd_enabled && pref.dnd_start_time && pref.dnd_end_time}
									<p>
										Do Not Disturb: {pref.dnd_start_time.substring(0, 5)} - {pref.dnd_end_time.substring(
											0,
											5
										)}
									</p>
								{/if}
								{#if pref.min_priority}
									<p>Minimum Priority: {pref.min_priority}</p>
								{/if}
							</div>
						</div>

						<div class="flex gap-2">
							<Button
								variant={pref.is_enabled ? 'secondary' : 'primary'}
								size="sm"
								on:click={() => handleTogglePreference(pref.id, pref.is_enabled)}
							>
								{pref.is_enabled ? 'Disable' : 'Enable'}
							</Button>
							<Button variant="danger" size="sm" on:click={() => handleDeletePreference(pref.id)}>
								Delete
							</Button>
						</div>
					</div>
				</div>
			{/each}
		</div>
	{/if}
</div>
