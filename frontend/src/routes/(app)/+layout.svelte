<script lang="ts">
	import { goto } from '$app/navigation';
	import { authStore } from '$lib/stores/auth';
	import Button from '$lib/components/ui/Button.svelte';

	$: if (!$authStore.isLoading && !$authStore.isAuthenticated) {
		goto('/login');
	}

	async function handleLogout() {
		await authStore.logout();
		goto('/login');
	}
</script>

{#if $authStore.isLoading}
	<div class="min-h-screen flex items-center justify-center">
		<div class="text-center">
			<div class="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
			<p class="mt-4 text-gray-600">Loading...</p>
		</div>
	</div>
{:else if $authStore.isAuthenticated}
	<div class="min-h-screen bg-gray-50">
		<!-- Header -->
		<header class="bg-white shadow">
			<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
				<div class="flex justify-between items-center h-16">
					<div class="flex items-center">
						<h1 class="text-2xl font-bold text-primary-600">Pulsar</h1>
						<nav class="ml-10 flex space-x-8">
							<a href="/dashboard" class="text-gray-700 hover:text-gray-900 px-3 py-2 text-sm font-medium">
								Dashboard
							</a>
							<a href="/alerts" class="text-gray-700 hover:text-gray-900 px-3 py-2 text-sm font-medium">
								Alerts
							</a>
							<a href="/incidents" class="text-gray-700 hover:text-gray-900 px-3 py-2 text-sm font-medium">
								Incidents
							</a>
							<a href="/schedules" class="text-gray-700 hover:text-gray-900 px-3 py-2 text-sm font-medium">
								Schedules
							</a>
						</nav>
					</div>

					<div class="flex items-center space-x-4">
						<span class="text-sm text-gray-700">
							{$authStore.user?.email}
						</span>
						<Button variant="secondary" on:click={handleLogout}>
							Logout
						</Button>
					</div>
				</div>
			</div>
		</header>

		<!-- Main content -->
		<main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
			<slot />
		</main>
	</div>
{/if}
