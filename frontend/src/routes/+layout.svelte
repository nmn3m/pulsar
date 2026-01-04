<script lang="ts">
	import '../app.css';
	import { onMount } from 'svelte';
	import { authStore } from '$lib/stores/auth';

	onMount(() => {
		authStore.init();
	});
</script>

<slot />
