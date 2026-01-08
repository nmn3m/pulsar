<script lang="ts">
  import { goto } from '$app/navigation';
  import { onMount } from 'svelte';

  onMount(() => {
    // Redirect to endpoints by default
    goto('/webhooks/endpoints');
  });
</script>

<div class="text-center py-12">
  <p class="text-gray-500">Redirecting to webhooks...</p>
</div>
