<script lang="ts">
	import { onMount } from 'svelte';
	import { api } from '$lib/api/client';
	import type {
		DNDSettings,
		DNDSchedule,
		DNDTimeSlot,
		DNDOverride,
		DayOfWeek,
		UpdateDNDSettingsRequest
	} from '$lib/types/dnd';
	import { DAYS_OF_WEEK, DAY_LABELS, COMMON_TIMEZONES } from '$lib/types/dnd';

	let settings: DNDSettings | null = null;
	let loading = true;
	let saving = false;
	let error: string | null = null;
	let successMessage: string | null = null;

	// Form state
	let enabled = false;
	let allowP1Override = true;
	let timezone = 'UTC';
	let weeklySlots: DNDTimeSlot[] = [];

	// Override form state
	let showAddOverride = false;
	let overrideStart = '';
	let overrideEnd = '';
	let overrideReason = '';

	onMount(async () => {
		await loadSettings();
	});

	async function loadSettings() {
		loading = true;
		error = null;
		try {
			settings = await api.getDNDSettings();
			if (settings) {
				enabled = settings.enabled;
				allowP1Override = settings.allow_p1_override;

				// Parse schedule
				let schedule: DNDSchedule = { weekly: [], timezone: 'UTC' };
				if (settings.schedule && typeof settings.schedule === 'object') {
					schedule = settings.schedule as unknown as DNDSchedule;
				}
				timezone = schedule.timezone || 'UTC';
				weeklySlots = schedule.weekly || [];
			}
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to load DND settings';
		} finally {
			loading = false;
		}
	}

	async function saveSettings() {
		saving = true;
		error = null;
		successMessage = null;

		try {
			const schedule: DNDSchedule = {
				weekly: weeklySlots,
				timezone
			};

			const updateRequest: UpdateDNDSettingsRequest = {
				enabled,
				allow_p1_override: allowP1Override,
				schedule
			};

			settings = await api.updateDNDSettings(updateRequest);
			successMessage = 'Settings saved successfully';
			setTimeout(() => (successMessage = null), 3000);
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to save settings';
		} finally {
			saving = false;
		}
	}

	function addTimeSlot(day: DayOfWeek) {
		weeklySlots = [
			...weeklySlots,
			{
				day,
				start: '22:00',
				end: '08:00'
			}
		];
	}

	function removeTimeSlot(index: number) {
		weeklySlots = weeklySlots.filter((_, i) => i !== index);
	}

	function updateTimeSlot(index: number, field: 'start' | 'end', value: string) {
		weeklySlots = weeklySlots.map((slot, i) => {
			if (i === index) {
				return { ...slot, [field]: value };
			}
			return slot;
		});
	}

	async function addOverride() {
		if (!overrideStart || !overrideEnd) {
			error = 'Please select both start and end dates';
			return;
		}

		saving = true;
		error = null;

		try {
			settings = await api.addDNDOverride({
				start: new Date(overrideStart).toISOString(),
				end: new Date(overrideEnd).toISOString(),
				reason: overrideReason || undefined
			});
			showAddOverride = false;
			overrideStart = '';
			overrideEnd = '';
			overrideReason = '';
			successMessage = 'Override added successfully';
			setTimeout(() => (successMessage = null), 3000);
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to add override';
		} finally {
			saving = false;
		}
	}

	async function removeOverride(index: number) {
		if (!confirm('Are you sure you want to remove this override?')) return;

		saving = true;
		error = null;

		try {
			settings = await api.removeDNDOverride(index);
			successMessage = 'Override removed successfully';
			setTimeout(() => (successMessage = null), 3000);
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to remove override';
		} finally {
			saving = false;
		}
	}

	function formatDate(dateStr: string): string {
		return new Date(dateStr).toLocaleString();
	}

	function getSlotsForDay(day: DayOfWeek): { slot: DNDTimeSlot; index: number }[] {
		return weeklySlots
			.map((slot, index) => ({ slot, index }))
			.filter(({ slot }) => slot.day === day);
	}

	function getOverrides(): DNDOverride[] {
		if (!settings?.overrides) return [];
		if (Array.isArray(settings.overrides)) return settings.overrides;
		return [];
	}
</script>

<svelte:head>
	<title>Do Not Disturb Settings - Pulsar</title>
</svelte:head>

<div class="container mx-auto max-w-4xl px-4 py-8">
	<div class="mb-8">
		<h1 class="text-2xl font-bold text-gray-900">Do Not Disturb Settings</h1>
		<p class="mt-2 text-gray-600">
			Configure quiet hours to prevent notifications during specific times.
		</p>
	</div>

	{#if error}
		<div class="mb-6 rounded-lg border border-red-200 bg-red-50 p-4">
			<p class="text-sm text-red-600">{error}</p>
		</div>
	{/if}

	{#if successMessage}
		<div class="mb-6 rounded-lg border border-green-200 bg-green-50 p-4">
			<p class="text-sm text-green-600">{successMessage}</p>
		</div>
	{/if}

	{#if loading}
		<div class="flex items-center justify-center py-12">
			<div class="h-8 w-8 animate-spin rounded-full border-b-2 border-blue-600"></div>
		</div>
	{:else}
		<div class="space-y-8">
			<!-- Enable/Disable DND -->
			<div class="rounded-lg border border-gray-200 bg-white p-6 shadow-sm">
				<div class="flex items-center justify-between">
					<div>
						<h2 class="text-lg font-medium text-gray-900">Enable Do Not Disturb</h2>
						<p class="mt-1 text-sm text-gray-500">
							When enabled, notifications will be suppressed during quiet hours.
						</p>
					</div>
					<label class="relative inline-flex cursor-pointer items-center">
						<input type="checkbox" bind:checked={enabled} class="peer sr-only" />
						<div
							class="peer h-6 w-11 rounded-full bg-gray-200 after:absolute after:left-[2px] after:top-[2px] after:h-5 after:w-5 after:rounded-full after:border after:border-gray-300 after:bg-white after:transition-all after:content-[''] peer-checked:bg-blue-600 peer-checked:after:translate-x-full peer-checked:after:border-white peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300"
						></div>
					</label>
				</div>

				<div class="mt-6 flex items-center justify-between border-t border-gray-200 pt-6">
					<div>
						<h3 class="font-medium text-gray-900">Allow P1 Override</h3>
						<p class="mt-1 text-sm text-gray-500">
							P1 (Critical) alerts will still notify you even during DND.
						</p>
					</div>
					<label class="relative inline-flex cursor-pointer items-center">
						<input type="checkbox" bind:checked={allowP1Override} class="peer sr-only" />
						<div
							class="peer h-6 w-11 rounded-full bg-gray-200 after:absolute after:left-[2px] after:top-[2px] after:h-5 after:w-5 after:rounded-full after:border after:border-gray-300 after:bg-white after:transition-all after:content-[''] peer-checked:bg-blue-600 peer-checked:after:translate-x-full peer-checked:after:border-white peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300"
						></div>
					</label>
				</div>
			</div>

			<!-- Timezone -->
			<div class="rounded-lg border border-gray-200 bg-white p-6 shadow-sm">
				<h2 class="text-lg font-medium text-gray-900">Timezone</h2>
				<p class="mt-1 text-sm text-gray-500">
					Select your timezone for the DND schedule.
				</p>
				<div class="mt-4">
					<select
						bind:value={timezone}
						class="block w-full rounded-lg border border-gray-300 bg-white px-3 py-2 text-gray-900 focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
					>
						{#each COMMON_TIMEZONES as tz}
							<option value={tz}>{tz}</option>
						{/each}
					</select>
				</div>
			</div>

			<!-- Weekly Schedule -->
			<div class="rounded-lg border border-gray-200 bg-white p-6 shadow-sm">
				<h2 class="text-lg font-medium text-gray-900">Weekly Schedule</h2>
				<p class="mt-1 text-sm text-gray-500">
					Set recurring quiet hours for each day of the week.
				</p>

				<div class="mt-6 space-y-4">
					{#each DAYS_OF_WEEK as day}
						<div class="border-b border-gray-200 pb-4 last:border-0">
							<div class="flex items-center justify-between">
								<span class="font-medium text-gray-900">{DAY_LABELS[day]}</span>
								<button
									type="button"
									on:click={() => addTimeSlot(day)}
									class="text-sm text-blue-600 hover:text-blue-800"
								>
									+ Add time slot
								</button>
							</div>

							{#each getSlotsForDay(day) as { slot, index }}
								<div class="mt-3 flex items-center gap-4">
									<div class="flex items-center gap-2">
										<label class="text-sm text-gray-500">From</label>
										<input
											type="time"
											value={slot.start}
											on:change={(e) => updateTimeSlot(index, 'start', e.currentTarget.value)}
											class="rounded-md border border-gray-300 bg-white px-2 py-1 text-sm text-gray-900"
										/>
									</div>
									<div class="flex items-center gap-2">
										<label class="text-sm text-gray-500">To</label>
										<input
											type="time"
											value={slot.end}
											on:change={(e) => updateTimeSlot(index, 'end', e.currentTarget.value)}
											class="rounded-md border border-gray-300 bg-white px-2 py-1 text-sm text-gray-900"
										/>
									</div>
									<button
										type="button"
										on:click={() => removeTimeSlot(index)}
										class="text-red-600 hover:text-red-800"
									>
										<svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
										</svg>
									</button>
								</div>
							{/each}

							{#if getSlotsForDay(day).length === 0}
								<p class="mt-2 text-sm text-gray-400">No quiet hours set</p>
							{/if}
						</div>
					{/each}
				</div>
			</div>

			<!-- Temporary Overrides -->
			<div class="rounded-lg border border-gray-200 bg-white p-6 shadow-sm">
				<div class="flex items-center justify-between">
					<div>
						<h2 class="text-lg font-medium text-gray-900">Temporary Overrides</h2>
						<p class="mt-1 text-sm text-gray-500">
							Add one-time DND periods (e.g., vacation, personal time).
						</p>
					</div>
					<button
						type="button"
						on:click={() => (showAddOverride = true)}
						class="rounded-lg bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700"
					>
						Add Override
					</button>
				</div>

				{#if showAddOverride}
					<div class="mt-4 rounded-lg border border-gray-200 bg-gray-50 p-4">
						<div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
							<div>
								<label class="block text-sm font-medium text-gray-700">Start Date/Time</label>
								<input
									type="datetime-local"
									bind:value={overrideStart}
									class="mt-1 block w-full rounded-md border border-gray-300 bg-white px-3 py-2 text-sm text-gray-900"
								/>
							</div>
							<div>
								<label class="block text-sm font-medium text-gray-700">End Date/Time</label>
								<input
									type="datetime-local"
									bind:value={overrideEnd}
									class="mt-1 block w-full rounded-md border border-gray-300 bg-white px-3 py-2 text-sm text-gray-900"
								/>
							</div>
						</div>
						<div class="mt-4">
							<label class="block text-sm font-medium text-gray-700">Reason (optional)</label>
							<input
								type="text"
								bind:value={overrideReason}
								placeholder="e.g., Vacation, Personal time"
								class="mt-1 block w-full rounded-md border border-gray-300 bg-white px-3 py-2 text-sm text-gray-900"
							/>
						</div>
						<div class="mt-4 flex justify-end gap-2">
							<button
								type="button"
								on:click={() => (showAddOverride = false)}
								class="rounded-lg border border-gray-300 bg-white px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50"
							>
								Cancel
							</button>
							<button
								type="button"
								on:click={addOverride}
								disabled={saving}
								class="rounded-lg bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700 disabled:opacity-50"
							>
								{saving ? 'Adding...' : 'Add Override'}
							</button>
						</div>
					</div>
				{/if}

				{#if getOverrides().length > 0}
					<div class="mt-4 space-y-3">
						{#each getOverrides() as override, index}
							<div class="flex items-center justify-between rounded-lg border border-gray-200 bg-gray-50 p-3">
								<div>
									<p class="font-medium text-gray-900">
										{formatDate(override.start)} - {formatDate(override.end)}
									</p>
									{#if override.reason}
										<p class="text-sm text-gray-500">{override.reason}</p>
									{/if}
								</div>
								<button
									type="button"
									on:click={() => removeOverride(index)}
									class="text-red-600 hover:text-red-800"
								>
									<svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
									</svg>
								</button>
							</div>
						{/each}
					</div>
				{:else}
					<p class="mt-4 text-sm text-gray-400">No temporary overrides set</p>
				{/if}
			</div>

			<!-- Save Button -->
			<div class="flex justify-end">
				<button
					type="button"
					on:click={saveSettings}
					disabled={saving}
					class="rounded-lg bg-blue-600 px-6 py-2 text-sm font-medium text-white hover:bg-blue-700 disabled:opacity-50"
				>
					{saving ? 'Saving...' : 'Save Settings'}
				</button>
			</div>
		</div>
	{/if}
</div>
