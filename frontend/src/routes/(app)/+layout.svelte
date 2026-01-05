<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { goto } from '$app/navigation';
	import { authStore } from '$lib/stores/auth';
	import { wsStore } from '$lib/stores/websocket';
	import Button from '$lib/components/ui/Button.svelte';

	$: if (!$authStore.isLoading && !$authStore.isAuthenticated) {
		wsStore.disconnect();
		goto('/login');
	}

	onMount(() => {
		if ($authStore.isAuthenticated) {
			wsStore.connect();
		}
	});

	onDestroy(() => {
		wsStore.disconnect();
	});

	async function handleLogout() {
		wsStore.disconnect();
		await authStore.logout();
		goto('/login');
	}

	function getStatusColor(status: string): string {
		switch (status) {
			case 'connected':
				return 'bg-green-500';
			case 'connecting':
				return 'bg-yellow-500';
			case 'disconnected':
				return 'bg-gray-400';
			case 'error':
				return 'bg-red-500';
			default:
				return 'bg-gray-400';
		}
	}

	function getStatusText(status: string): string {
		switch (status) {
			case 'connected':
				return 'Connected';
			case 'connecting':
				return 'Connecting...';
			case 'disconnected':
				return 'Disconnected';
			case 'error':
				return 'Connection Error';
			default:
				return 'Unknown';
		}
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
						<!-- WebSocket Status Indicator -->
						<div class="flex items-center space-x-2" title={getStatusText($wsStore.status)}>
							<span class="relative flex h-3 w-3">
								{#if $wsStore.status === 'connecting'}
									<span
										class="animate-ping absolute inline-flex h-full w-full rounded-full {getStatusColor(
											$wsStore.status
										)} opacity-75"
									/>
								{/if}
								<span
									class="relative inline-flex rounded-full h-3 w-3 {getStatusColor(
										$wsStore.status
									)}"
								/>
							</span>
							<span class="text-xs text-gray-600 hidden sm:inline">
								{getStatusText($wsStore.status)}
							</span>
						</div>

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
