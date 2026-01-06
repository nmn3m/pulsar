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
				return 'bg-neon-green shadow-neon-green';
			case 'connecting':
				return 'bg-neon-yellow animate-pulse';
			case 'disconnected':
				return 'bg-gray-500';
			case 'error':
				return 'bg-neon-red shadow-neon-pink';
			default:
				return 'bg-gray-500';
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
			<div class="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-primary-500"></div>
			<p class="mt-4 text-gray-400">Loading...</p>
		</div>
	</div>
{:else if $authStore.isAuthenticated}
	<div class="min-h-screen">
		<!-- Header -->
		<header class="bg-space-800/80 backdrop-blur-md border-b border-primary-500/20 shadow-lg shadow-primary-500/5">
			<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
				<div class="flex justify-between items-center h-16">
					<div class="flex items-center">
						<h1 class="text-2xl font-bold text-primary-400 text-glow-cyan">Pulsar</h1>
						<nav class="ml-10 flex space-x-2">
							<a href="/dashboard" class="text-gray-300 hover:text-primary-400 hover:bg-space-700/50 px-4 py-2 rounded-lg text-sm font-medium transition-all duration-200">
								Dashboard
							</a>
							<a href="/alerts" class="text-gray-300 hover:text-accent-400 hover:bg-space-700/50 px-4 py-2 rounded-lg text-sm font-medium transition-all duration-200">
								Alerts
							</a>
							<a href="/incidents" class="text-gray-300 hover:text-accent-400 hover:bg-space-700/50 px-4 py-2 rounded-lg text-sm font-medium transition-all duration-200">
								Incidents
							</a>
							<a href="/schedules" class="text-gray-300 hover:text-primary-400 hover:bg-space-700/50 px-4 py-2 rounded-lg text-sm font-medium transition-all duration-200">
								Schedules
							</a>
							<a href="/webhooks/endpoints" class="text-gray-300 hover:text-primary-400 hover:bg-space-700/50 px-4 py-2 rounded-lg text-sm font-medium transition-all duration-200">
								Webhooks
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
							<span class="text-xs text-gray-400 hidden sm:inline">
								{getStatusText($wsStore.status)}
							</span>
						</div>

						<span class="text-sm text-primary-300">
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
