<script lang="ts">
	import { authStore } from '$lib/stores/auth';
</script>

<svelte:head>
	<title>Dashboard - Pulsar</title>
</svelte:head>

<div class="space-y-6">
	<div>
		<h2 class="text-3xl font-bold text-gray-900">Dashboard</h2>
		<p class="mt-2 text-gray-600">Welcome to Pulsar incident management</p>
	</div>

	<div class="grid grid-cols-1 md:grid-cols-3 gap-6">
		<!-- Active Alerts Card -->
		<div class="bg-white p-6 rounded-lg shadow">
			<h3 class="text-sm font-medium text-gray-500">Active Alerts</h3>
			<p class="mt-2 text-3xl font-semibold text-gray-900">0</p>
			<p class="mt-2 text-sm text-gray-600">No active alerts</p>
		</div>

		<!-- Open Incidents Card -->
		<div class="bg-white p-6 rounded-lg shadow">
			<h3 class="text-sm font-medium text-gray-500">Open Incidents</h3>
			<p class="mt-2 text-3xl font-semibold text-gray-900">0</p>
			<p class="mt-2 text-sm text-gray-600">No open incidents</p>
		</div>

		<!-- On-Call Status Card -->
		<div class="bg-white p-6 rounded-lg shadow">
			<h3 class="text-sm font-medium text-gray-500">On-Call Status</h3>
			<p class="mt-2 text-3xl font-semibold text-gray-900">Off</p>
			<p class="mt-2 text-sm text-gray-600">You are not currently on-call</p>
		</div>
	</div>

	<!-- Welcome Card -->
	<div class="bg-primary-50 border border-primary-200 p-6 rounded-lg">
		<h3 class="text-lg font-semibold text-primary-900">Welcome to Pulsar!</h3>
		<p class="mt-2 text-primary-700">
			You've successfully set up your account. Here's what you can do next:
		</p>
		<ul class="mt-4 space-y-2 text-primary-700">
			<li>• Set up your first on-call schedule</li>
			<li>• Configure integrations for alert sources</li>
			<li>• Create escalation policies</li>
			<li>• Invite team members</li>
		</ul>
	</div>

	<!-- User Info Card -->
	<div class="bg-white p-6 rounded-lg shadow">
		<h3 class="text-lg font-semibold text-gray-900">Your Information</h3>
		<dl class="mt-4 space-y-2">
			<div class="flex justify-between">
				<dt class="text-sm font-medium text-gray-500">Email:</dt>
				<dd class="text-sm text-gray-900">{$authStore.user?.email}</dd>
			</div>
			<div class="flex justify-between">
				<dt class="text-sm font-medium text-gray-500">Username:</dt>
				<dd class="text-sm text-gray-900">{$authStore.user?.username}</dd>
			</div>
			<div class="flex justify-between">
				<dt class="text-sm font-medium text-gray-500">Organization:</dt>
				<dd class="text-sm text-gray-900">{$authStore.organization?.name}</dd>
			</div>
			<div class="flex justify-between">
				<dt class="text-sm font-medium text-gray-500">Plan:</dt>
				<dd class="text-sm text-gray-900 capitalize">{$authStore.organization?.plan}</dd>
			</div>
		</dl>
	</div>
</div>
