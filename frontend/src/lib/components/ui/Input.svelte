<script lang="ts">
	export let type: 'text' | 'email' | 'password' = 'text';
	export let placeholder = '';
	export let value = '';
	export let label = '';
	export let error = '';
	export let required = false;
	export let id = '';
</script>

<div class="w-full">
	{#if label}
		<label for={id} class="block text-sm font-medium text-gray-700 mb-1">
			{label}
			{#if required}<span class="text-red-500">*</span>{/if}
		</label>
	{/if}
	<input
		{id}
		{type}
		{placeholder}
		{required}
		bind:value
		class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent {error ? 'border-red-500' : ''}"
		on:input
		on:blur
	/>
	{#if error}
		<p class="mt-1 text-sm text-red-600">{error}</p>
	{/if}
</div>
