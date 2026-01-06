<script lang="ts">
	export let type: 'button' | 'submit' | 'reset' = 'button';
	export let variant: 'primary' | 'secondary' | 'danger' | 'ghost' = 'primary';
	export let size: 'sm' | 'md' | 'lg' = 'md';
	export let disabled = false;
	export let fullWidth = false;

	const variantClasses = {
		primary: 'bg-primary-600 hover:bg-primary-500 dark:bg-primary-500 dark:hover:bg-primary-400 text-white dark:text-space-900 font-semibold dark:shadow-neon-cyan hover:shadow-lg hover:shadow-primary-500/50 border border-primary-500 dark:border-primary-400/50',
		secondary: 'bg-gray-100 hover:bg-gray-200 dark:bg-space-700 dark:hover:bg-space-600 text-gray-700 dark:text-gray-200 border border-gray-300 dark:border-space-500 hover:border-primary-500/50',
		danger: 'bg-accent-600 hover:bg-accent-500 text-white dark:shadow-neon-pink hover:shadow-lg hover:shadow-accent-500/50 border border-accent-500 dark:border-accent-400/50',
		ghost: 'bg-transparent hover:bg-gray-100 dark:hover:bg-space-700/50 text-gray-600 dark:text-gray-300 hover:text-primary-600 dark:hover:text-primary-400 border border-transparent hover:border-primary-500/30'
	};

	const sizeClasses = {
		sm: 'px-3 py-1.5 text-sm',
		md: 'px-4 py-2 text-base',
		lg: 'px-6 py-3 text-lg'
	};
</script>

<button
	{type}
	{disabled}
	class="rounded-lg font-medium transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed {variantClasses[variant]} {sizeClasses[size]} {fullWidth ? 'w-full' : ''}"
	on:click
>
	<slot />
</button>
