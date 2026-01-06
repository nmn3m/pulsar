<script lang="ts">
	export let type: 'button' | 'submit' | 'reset' = 'button';
	export let variant: 'primary' | 'secondary' | 'danger' | 'ghost' = 'primary';
	export let disabled = false;
	export let fullWidth = false;

	const variantClasses = {
		primary: 'bg-primary-500 hover:bg-primary-400 text-space-900 font-semibold shadow-neon-cyan hover:shadow-lg hover:shadow-primary-500/50 border border-primary-400/50',
		secondary: 'bg-space-700 hover:bg-space-600 text-gray-200 border border-space-500 hover:border-primary-500/50',
		danger: 'bg-accent-600 hover:bg-accent-500 text-white shadow-neon-pink hover:shadow-lg hover:shadow-accent-500/50 border border-accent-400/50',
		ghost: 'bg-transparent hover:bg-space-700/50 text-gray-300 hover:text-primary-400 border border-transparent hover:border-primary-500/30'
	};
</script>

<button
	{type}
	{disabled}
	class="px-4 py-2 rounded-lg font-medium transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed {variantClasses[variant]} {fullWidth ? 'w-full' : ''}"
	on:click
>
	<slot />
</button>
