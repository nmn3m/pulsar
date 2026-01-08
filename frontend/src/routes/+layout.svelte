<script lang="ts">
  import '../app.css';
  import { onMount } from 'svelte';
  import { authStore } from '$lib/stores/auth';
  import { themeStore } from '$lib/stores/theme';

  onMount(() => {
    authStore.init();
    themeStore.initialize();
  });
</script>

<slot />
