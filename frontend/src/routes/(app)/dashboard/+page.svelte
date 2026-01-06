<script lang="ts">
	import { authStore } from '$lib/stores/auth';
</script>

<svelte:head>
	<title>Dashboard - Pulsar</title>
</svelte:head>

<div class="space-y-6">
	<div>
		<h2 class="text-3xl font-bold text-gray-900 dark:text-gray-100">Dashboard</h2>
		<p class="mt-2 text-gray-500 dark:text-gray-400">Welcome to Pulsar incident management</p>
	</div>

	<div class="grid grid-cols-1 md:grid-cols-3 gap-6">
		<!-- Active Alerts Card -->
		<div class="bg-white dark:bg-space-800/50 backdrop-blur-sm p-6 rounded-xl border border-accent-200 dark:border-accent-500/20 hover:border-accent-400 dark:hover:border-accent-500/40 transition-all duration-300 hover:shadow-lg dark:hover:shadow-accent-500/10 shadow-sm">
			<div class="flex items-center justify-between">
				<h3 class="text-sm font-medium text-gray-500 dark:text-gray-400">Active Alerts</h3>
				<span class="text-accent-500 dark:text-accent-400">
					<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
					</svg>
				</span>
			</div>
			<p class="mt-3 text-4xl font-bold text-gray-900 dark:text-gray-100">0</p>
			<p class="mt-2 text-sm text-gray-400 dark:text-gray-500">No active alerts</p>
		</div>

		<!-- Open Incidents Card -->
		<div class="bg-white dark:bg-space-800/50 backdrop-blur-sm p-6 rounded-xl border border-accent-200 dark:border-accent-500/20 hover:border-accent-400 dark:hover:border-accent-500/40 transition-all duration-300 hover:shadow-lg dark:hover:shadow-accent-500/10 shadow-sm">
			<div class="flex items-center justify-between">
				<h3 class="text-sm font-medium text-gray-500 dark:text-gray-400">Open Incidents</h3>
				<span class="text-accent-600 dark:text-accent-500">
					<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
					</svg>
				</span>
			</div>
			<p class="mt-3 text-4xl font-bold text-gray-900 dark:text-gray-100">0</p>
			<p class="mt-2 text-sm text-gray-400 dark:text-gray-500">No open incidents</p>
		</div>

		<!-- On-Call Status Card -->
		<div class="bg-white dark:bg-space-800/50 backdrop-blur-sm p-6 rounded-xl border border-primary-200 dark:border-primary-500/20 hover:border-primary-400 dark:hover:border-primary-500/40 transition-all duration-300 hover:shadow-lg dark:hover:shadow-primary-500/10 shadow-sm">
			<div class="flex items-center justify-between">
				<h3 class="text-sm font-medium text-gray-500 dark:text-gray-400">On-Call Status</h3>
				<span class="text-primary-500 dark:text-primary-400">
					<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
					</svg>
				</span>
			</div>
			<p class="mt-3 text-4xl font-bold text-gray-900 dark:text-gray-100">Off</p>
			<p class="mt-2 text-sm text-gray-400 dark:text-gray-500">You are not currently on-call</p>
		</div>
	</div>

	<!-- Welcome Card -->
	<div class="bg-gradient-to-r from-primary-50 to-accent-50 dark:from-primary-900/30 dark:to-accent-900/20 border border-primary-200 dark:border-primary-500/30 p-6 rounded-xl backdrop-blur-sm">
		<h3 class="text-lg font-semibold text-primary-600 dark:text-primary-300 dark:text-glow-cyan">Welcome to Pulsar!</h3>
		<p class="mt-2 text-gray-600 dark:text-gray-300">
			You've successfully set up your account. Here's what you can do next:
		</p>
		<ul class="mt-4 space-y-2 text-gray-500 dark:text-gray-400">
			<li class="flex items-center gap-2">
				<span class="text-primary-500 dark:text-primary-400">&#8226;</span> Set up your first on-call schedule
			</li>
			<li class="flex items-center gap-2">
				<span class="text-primary-500 dark:text-primary-400">&#8226;</span> Configure integrations for alert sources
			</li>
			<li class="flex items-center gap-2">
				<span class="text-primary-500 dark:text-primary-400">&#8226;</span> Create escalation policies
			</li>
			<li class="flex items-center gap-2">
				<span class="text-primary-500 dark:text-primary-400">&#8226;</span> Invite team members
			</li>
		</ul>
	</div>

	<!-- User Info Card -->
	<div class="bg-white dark:bg-space-800/50 backdrop-blur-sm p-6 rounded-xl border border-gray-200 dark:border-space-600 shadow-sm">
		<h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">Your Information</h3>
		<dl class="mt-4 space-y-3">
			<div class="flex justify-between items-center py-2 border-b border-gray-100 dark:border-space-700">
				<dt class="text-sm font-medium text-gray-400 dark:text-gray-500">Email:</dt>
				<dd class="text-sm text-primary-600 dark:text-primary-300">{$authStore.user?.email}</dd>
			</div>
			<div class="flex justify-between items-center py-2 border-b border-gray-100 dark:border-space-700">
				<dt class="text-sm font-medium text-gray-400 dark:text-gray-500">Username:</dt>
				<dd class="text-sm text-gray-700 dark:text-gray-300">{$authStore.user?.username}</dd>
			</div>
			<div class="flex justify-between items-center py-2 border-b border-gray-100 dark:border-space-700">
				<dt class="text-sm font-medium text-gray-400 dark:text-gray-500">Organization:</dt>
				<dd class="text-sm text-gray-700 dark:text-gray-300">{$authStore.organization?.name}</dd>
			</div>
			<div class="flex justify-between items-center py-2">
				<dt class="text-sm font-medium text-gray-400 dark:text-gray-500">Plan:</dt>
				<dd class="text-sm text-primary-600 dark:text-primary-400 capitalize">{$authStore.organization?.plan}</dd>
			</div>
		</dl>
	</div>
</div>
