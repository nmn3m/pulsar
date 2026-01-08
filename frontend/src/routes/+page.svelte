<script lang="ts">
  import { goto } from '$app/navigation';
  import { authStore } from '$lib/stores/auth';
  import { onMount } from 'svelte';

  onMount(() => {
    const unsubscribe = authStore.subscribe((state) => {
      if (!state.isLoading) {
        if (state.isAuthenticated) {
          goto('/dashboard');
        } else {
          goto('/login');
        }
      }
    });

    return unsubscribe;
  });
</script>

<div class="min-h-screen flex items-center justify-center">
  <div class="text-center">
    <div
      class="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"
    ></div>
    <p class="mt-4 text-gray-600">Loading...</p>
  </div>
</div>
