<script lang="ts">
	export let type: 'button' | 'submit' | 'reset' = 'button';
	export let variant: 'primary' | 'secondary' | 'danger' = 'primary';
	export let disabled = false;
	export let fullWidth = false;

	const variantClasses = {
		primary: 'bg-primary-600 hover:bg-primary-700 text-white',
		secondary: 'bg-gray-200 hover:bg-gray-300 text-gray-900',
		danger: 'bg-red-600 hover:bg-red-700 text-white'
	};
</script>

<button
	{type}
	{disabled}
	class="px-4 py-2 rounded-lg font-medium transition-colors disabled:opacity-50 disabled:cursor-not-allowed {variantClasses[variant]} {fullWidth ? 'w-full' : ''}"
	on:click
>
	<slot />
</button>
