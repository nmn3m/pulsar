<script lang="ts">
  export let type: 'text' | 'email' | 'password' | 'date' | 'time' | 'number' = 'text';
  export let placeholder = '';
  export let value = '';
  export let label = '';
  export let error = '';
  export let required = false;
  export let id = '';

  function handleInput(e: Event) {
    value = (e.target as HTMLInputElement).value;
  }
</script>

<div class="w-full">
  {#if label}
    <label for={id} class="block text-sm font-medium text-gray-600 dark:text-gray-300 mb-1">
      {label}
      {#if required}<span class="text-accent-500 dark:text-accent-400">*</span>{/if}
    </label>
  {/if}
  <input
    {id}
    {type}
    {placeholder}
    {required}
    {value}
    on:input={handleInput}
    class="w-full px-3 py-2 bg-white dark:bg-space-800 border border-gray-300 dark:border-space-500 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500 text-gray-900 dark:text-white placeholder-gray-400 dark:placeholder-gray-500 caret-gray-900 dark:caret-white transition-all duration-200 {error
      ? 'border-accent-500 focus:ring-accent-500'
      : ''}"
    on:input
    on:blur
  />
  {#if error}
    <p class="mt-1 text-sm text-accent-500 dark:text-accent-400">{error}</p>
  {/if}
</div>
