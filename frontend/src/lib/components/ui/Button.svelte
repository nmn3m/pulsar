<script lang="ts">
  export let type: 'button' | 'submit' | 'reset' = 'button';
  export let variant: 'primary' | 'secondary' | 'danger' | 'ghost' = 'primary';
  export let size: 'sm' | 'md' | 'lg' = 'md';
  export let disabled = false;
  export let fullWidth = false;

  const variantClasses = {
    primary:
      'bg-primary-600 hover:bg-primary-500 text-white font-semibold hover:shadow-lg hover:shadow-primary-400/50 border border-primary-500',
    secondary:
      'bg-gray-100 hover:bg-gray-200 text-gray-700 border border-gray-300 hover:border-primary-400/50',
    danger:
      'bg-error hover:bg-error-dark text-white hover:shadow-lg hover:shadow-error/50 border border-error',
    ghost:
      'bg-transparent hover:bg-gray-100 text-gray-600 hover:text-primary-600 border border-transparent hover:border-primary-400/30',
  };

  const sizeClasses = {
    sm: 'px-3 py-1.5 text-sm',
    md: 'px-4 py-2 text-base',
    lg: 'px-6 py-3 text-lg',
  };
</script>

<button
  {type}
  {disabled}
  class="rounded-lg font-medium transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed {variantClasses[
    variant
  ]} {sizeClasses[size]} {fullWidth ? 'w-full' : ''}"
  on:click
>
  <slot />
</button>
