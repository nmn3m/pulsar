<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { notificationChannelsStore } from '$lib/stores/notifications';
	import type { ChannelType } from '$lib/types/notification';
	import Button from '$lib/components/ui/Button.svelte';
	import Input from '$lib/components/ui/Input.svelte';

	let showCreateForm = false;
	let name = '';
	let channelType: ChannelType = 'email';
	let isEnabled = true;

	// Provider-specific config fields
	let emailConfig = {
		smtp_host: '',
		smtp_port: 587,
		smtp_username: '',
		smtp_password: '',
		from_address: '',
		from_name: '',
		use_tls: true
	};

	let slackConfig = {
		webhook_url: '',
		channel: '',
		username: 'Pulsar',
		icon_emoji: ':bell:'
	};

	let teamsConfig = {
		webhook_url: '',
		theme_color: '0078D4'
	};

	let webhookConfig = {
		url: '',
		method: 'POST',
		headers: {} as Record<string, string>,
		timeout: 30
	};

	let createError = '';
	let creating = false;

	onMount(() => {
		notificationChannelsStore.load();
	});

	function resetForm() {
		name = '';
		channelType = 'email';
		isEnabled = true;
		emailConfig = {
			smtp_host: '',
			smtp_port: 587,
			smtp_username: '',
			smtp_password: '',
			from_address: '',
			from_name: '',
			use_tls: true
		};
		slackConfig = {
			webhook_url: '',
			channel: '',
			username: 'Pulsar',
			icon_emoji: ':bell:'
		};
		teamsConfig = {
			webhook_url: '',
			theme_color: '0078D4'
		};
		webhookConfig = {
			url: '',
			method: 'POST',
			headers: {},
			timeout: 30
		};
		createError = '';
	}

	async function handleCreateChannel() {
		createError = '';
		creating = true;

		try {
			let config: Record<string, unknown> = {};

			switch (channelType) {
				case 'email':
					config = emailConfig;
					break;
				case 'slack':
					config = slackConfig;
					break;
				case 'teams':
					config = teamsConfig;
					break;
				case 'webhook':
					config = webhookConfig;
					break;
			}

			await notificationChannelsStore.create({
				name,
				channel_type: channelType,
				is_enabled: isEnabled,
				config
			});

			resetForm();
			showCreateForm = false;
		} catch (err) {
			createError = err instanceof Error ? err.message : 'Failed to create notification channel';
		} finally {
			creating = false;
		}
	}

	async function handleDeleteChannel(id: string, channelName: string) {
		if (!confirm(`Are you sure you want to delete notification channel "${channelName}"?`)) {
			return;
		}

		try {
			await notificationChannelsStore.delete(id);
		} catch (err) {
			alert(err instanceof Error ? err.message : 'Failed to delete notification channel');
		}
	}

	function getChannelTypeDisplay(type: ChannelType): string {
		switch (type) {
			case 'email':
				return 'Email (SMTP)';
			case 'slack':
				return 'Slack';
			case 'teams':
				return 'Microsoft Teams';
			case 'webhook':
				return 'Webhook';
			default:
				return type;
		}
	}

	function getChannelTypeIcon(type: ChannelType): string {
		switch (type) {
			case 'email':
				return '📧';
			case 'slack':
				return '💬';
			case 'teams':
				return '👥';
			case 'webhook':
				return '🔗';
			default:
				return '📡';
		}
	}
</script>

<svelte:head>
	<title>Notification Channels - Pulsar</title>
</svelte:head>

<div class="space-y-6">
	<div class="flex justify-between items-center">
		<div>
			<h2 class="text-3xl font-bold text-gray-900">Notification Channels</h2>
			<p class="mt-2 text-gray-600">
				Configure notification delivery methods for your organization
			</p>
		</div>
		<Button
			variant="primary"
			on:click={() => {
				showCreateForm = !showCreateForm;
				if (!showCreateForm) resetForm();
			}}
		>
			{showCreateForm ? 'Cancel' : 'Add Channel'}
		</Button>
	</div>

	{#if showCreateForm}
		<div class="bg-white p-6 rounded-lg shadow">
			<h3 class="text-lg font-semibold mb-4">Create Notification Channel</h3>
			<form on:submit|preventDefault={handleCreateChannel} class="space-y-4">
				<Input id="name" label="Channel Name" bind:value={name} required />

				<div>
					<label for="channel-type" class="block text-sm font-medium text-gray-700 mb-1">
						Channel Type
					</label>
					<select
						id="channel-type"
						bind:value={channelType}
						class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500"
					>
						<option value="email">Email (SMTP)</option>
						<option value="slack">Slack</option>
						<option value="teams">Microsoft Teams</option>
						<option value="webhook">Webhook</option>
					</select>
				</div>

				<div class="flex items-center gap-2">
					<input
						id="is-enabled"
						type="checkbox"
						bind:checked={isEnabled}
						class="rounded border-gray-300 text-primary-600 focus:ring-primary-500"
					/>
					<label for="is-enabled" class="text-sm font-medium text-gray-700">
						Enable this channel
					</label>
				</div>

				<!-- Email Configuration -->
				{#if channelType === 'email'}
					<div class="space-y-3 p-4 bg-gray-50 rounded-lg">
						<h4 class="text-sm font-semibold text-gray-900">Email Configuration</h4>
						<Input
							id="smtp-host"
							label="SMTP Host"
							bind:value={emailConfig.smtp_host}
							placeholder="smtp.gmail.com"
							required
						/>
						<Input
							id="smtp-port"
							label="SMTP Port"
							type="number"
							bind:value={emailConfig.smtp_port}
							required
						/>
						<Input
							id="smtp-username"
							label="SMTP Username"
							bind:value={emailConfig.smtp_username}
							required
						/>
						<Input
							id="smtp-password"
							label="SMTP Password"
							type="password"
							bind:value={emailConfig.smtp_password}
							required
						/>
						<Input
							id="from-address"
							label="From Email Address"
							bind:value={emailConfig.from_address}
							placeholder="alerts@example.com"
							required
						/>
						<Input
							id="from-name"
							label="From Name (optional)"
							bind:value={emailConfig.from_name}
							placeholder="Pulsar Alerts"
						/>
						<div class="flex items-center gap-2">
							<input
								id="use-tls"
								type="checkbox"
								bind:checked={emailConfig.use_tls}
								class="rounded border-gray-300 text-primary-600 focus:ring-primary-500"
							/>
							<label for="use-tls" class="text-sm font-medium text-gray-700">Use TLS</label>
						</div>
					</div>
				{/if}

				<!-- Slack Configuration -->
				{#if channelType === 'slack'}
					<div class="space-y-3 p-4 bg-gray-50 rounded-lg">
						<h4 class="text-sm font-semibold text-gray-900">Slack Configuration</h4>
						<Input
							id="webhook-url"
							label="Webhook URL"
							bind:value={slackConfig.webhook_url}
							placeholder="https://hooks.slack.com/services/..."
							required
						/>
						<Input
							id="channel"
							label="Default Channel (optional)"
							bind:value={slackConfig.channel}
							placeholder="#alerts"
						/>
						<Input
							id="username"
							label="Bot Username (optional)"
							bind:value={slackConfig.username}
							placeholder="Pulsar"
						/>
						<Input
							id="icon-emoji"
							label="Icon Emoji (optional)"
							bind:value={slackConfig.icon_emoji}
							placeholder=":bell:"
						/>
					</div>
				{/if}

				<!-- Teams Configuration -->
				{#if channelType === 'teams'}
					<div class="space-y-3 p-4 bg-gray-50 rounded-lg">
						<h4 class="text-sm font-semibold text-gray-900">Microsoft Teams Configuration</h4>
						<Input
							id="teams-webhook-url"
							label="Webhook URL"
							bind:value={teamsConfig.webhook_url}
							placeholder="https://outlook.office.com/webhook/..."
							required
						/>
						<Input
							id="theme-color"
							label="Theme Color (hex without #)"
							bind:value={teamsConfig.theme_color}
							placeholder="0078D4"
						/>
					</div>
				{/if}

				<!-- Webhook Configuration -->
				{#if channelType === 'webhook'}
					<div class="space-y-3 p-4 bg-gray-50 rounded-lg">
						<h4 class="text-sm font-semibold text-gray-900">Webhook Configuration</h4>
						<Input
							id="webhook-url-custom"
							label="Webhook URL"
							bind:value={webhookConfig.url}
							placeholder="https://api.example.com/notifications"
							required
						/>
						<div>
							<label for="method" class="block text-sm font-medium text-gray-700 mb-1">
								HTTP Method
							</label>
							<select
								id="method"
								bind:value={webhookConfig.method}
								class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500"
							>
								<option value="POST">POST</option>
								<option value="PUT">PUT</option>
								<option value="PATCH">PATCH</option>
							</select>
						</div>
						<Input
							id="timeout"
							label="Timeout (seconds)"
							type="number"
							bind:value={webhookConfig.timeout}
							min="1"
							max="300"
						/>
					</div>
				{/if}

				{#if createError}
					<div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg">
						{createError}
					</div>
				{/if}

				<div class="flex gap-2">
					<Button type="submit" variant="primary" disabled={creating}>
						{creating ? 'Creating...' : 'Create Channel'}
					</Button>
					<Button
						type="button"
						variant="secondary"
						on:click={() => {
							showCreateForm = false;
							resetForm();
						}}
					>
						Cancel
					</Button>
				</div>
			</form>
		</div>
	{/if}

	{#if $notificationChannelsStore.isLoading}
		<div class="text-center py-12">
			<div
				class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600"
			></div>
			<p class="mt-2 text-gray-600">Loading notification channels...</p>
		</div>
	{:else if $notificationChannelsStore.error}
		<div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg">
			{$notificationChannelsStore.error}
		</div>
	{:else if $notificationChannelsStore.channels.length === 0}
		<div class="text-center py-12 bg-white rounded-lg shadow">
			<p class="text-gray-600">No notification channels configured</p>
			<p class="text-sm text-gray-500 mt-2">Create your first channel to get started</p>
		</div>
	{:else}
		<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
			{#each $notificationChannelsStore.channels as channel (channel.id)}
				<div class="bg-white p-6 rounded-lg shadow hover:shadow-md transition-shadow">
					<div class="flex justify-between items-start mb-4">
						<div class="flex-1">
							<div class="flex items-center gap-2 mb-2">
								<span class="text-2xl">{getChannelTypeIcon(channel.channel_type)}</span>
								<h3 class="text-lg font-semibold text-gray-900">{channel.name}</h3>
							</div>
							<p class="text-sm text-gray-600">{getChannelTypeDisplay(channel.channel_type)}</p>
						</div>
					</div>

					<div class="space-y-2 text-sm">
						<div class="flex items-center gap-2">
							<span
								class="px-2 py-1 rounded text-xs font-medium {channel.is_enabled
									? 'bg-green-100 text-green-800'
									: 'bg-gray-100 text-gray-800'}"
							>
								{channel.is_enabled ? 'Enabled' : 'Disabled'}
							</span>
						</div>
						<p class="text-xs text-gray-500">
							Created {new Date(channel.created_at).toLocaleDateString()}
						</p>
					</div>

					<div class="flex gap-2 mt-4">
						<Button
							variant="primary"
							size="sm"
							on:click={() => goto(`/notifications/channels/${channel.id}`)}
						>
							Configure
						</Button>
						<Button
							variant="danger"
							size="sm"
							on:click={() => handleDeleteChannel(channel.id, channel.name)}
						>
							Delete
						</Button>
					</div>
				</div>
			{/each}
		</div>
	{/if}
</div>
