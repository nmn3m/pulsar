<script lang="ts">
	import type { Alert } from '$lib/types/alert';
	import dayjs from 'dayjs';
	import relativeTime from 'dayjs/plugin/relativeTime';
	import Button from '../ui/Button.svelte';

	dayjs.extend(relativeTime);

	export let alert: Alert;
	export let onAcknowledge: (id: string) => void = () => {};
	export let onClose: (id: string) => void = () => {};
	export let onClick: (id: string) => void = () => {};

	const priorityColors = {
		P1: 'bg-red-100 text-red-800 border-red-200',
		P2: 'bg-orange-100 text-orange-800 border-orange-200',
		P3: 'bg-yellow-100 text-yellow-800 border-yellow-200',
		P4: 'bg-blue-100 text-blue-800 border-blue-200',
		P5: 'bg-gray-100 text-gray-800 border-gray-200'
	};

	const statusColors = {
		open: 'bg-red-100 text-red-800',
		acknowledged: 'bg-yellow-100 text-yellow-800',
		closed: 'bg-green-100 text-green-800',
		snoozed: 'bg-purple-100 text-purple-800'
	};
</script>

<div
	class="bg-white p-4 rounded-lg shadow hover:shadow-md transition-shadow cursor-pointer border-l-4 {priorityColors[alert.priority]}"
	on:click={() => onClick(alert.id)}
	on:keydown={(e) => e.key === 'Enter' && onClick(alert.id)}
	role="button"
	tabindex="0"
>
	<div class="flex items-start justify-between">
		<div class="flex-1">
			<div class="flex items-center gap-2 mb-2">
				<span class="px-2 py-1 text-xs font-semibold rounded {priorityColors[alert.priority]}">
					{alert.priority}
				</span>
				<span class="px-2 py-1 text-xs font-semibold rounded {statusColors[alert.status]}">
					{alert.status}
				</span>
				<span class="text-xs text-gray-500">{alert.source}</span>
			</div>

			<h3 class="text-lg font-semibold text-gray-900 mb-1">{alert.message}</h3>

			{#if alert.description}
				<p class="text-sm text-gray-600 mb-2">{alert.description}</p>
			{/if}

			{#if alert.tags && alert.tags.length > 0}
				<div class="flex gap-1 flex-wrap mb-2">
					{#each alert.tags as tag}
						<span class="px-2 py-0.5 text-xs bg-gray-100 text-gray-700 rounded">{tag}</span>
					{/each}
				</div>
			{/if}

			<div class="text-xs text-gray-500">
				Created {dayjs(alert.created_at).fromNow()}
			</div>
		</div>

		<div class="flex gap-2 ml-4" on:click|stopPropagation on:keydown|stopPropagation role="none">
			{#if alert.status === 'open'}
				<Button variant="primary" on:click={() => onAcknowledge(alert.id)}>
					Acknowledge
				</Button>
			{/if}
			{#if alert.status !== 'closed'}
				<Button variant="secondary" on:click={() => onClose(alert.id)}>
					Close
				</Button>
			{/if}
		</div>
	</div>
</div>
