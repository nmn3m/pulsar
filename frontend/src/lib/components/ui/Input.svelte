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
    <label for={id} class="block text-sm font-medium text-gray-600 mb-1">
      {label}
      {#if required}<span class="text-error">*</span>{/if}
    </label>
  {/if}
  <input
    {id}
    {type}
    {placeholder}
    {required}
    {value}
    on:input={handleInput}
    class="w-full px-3 py-2 bg-white border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-400 focus:border-primary-400 text-gray-900 placeholder-gray-400 caret-gray-900 transition-all duration-200 {error
      ? 'border-error focus:ring-error'
      : ''}"
    on:input
    on:blur
  />
  {#if error}
    <p class="mt-1 text-sm text-error">{error}</p>
  {/if}
</div>
