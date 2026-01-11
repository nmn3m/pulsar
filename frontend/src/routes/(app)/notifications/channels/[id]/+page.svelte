<script lang="ts">
  import { onMount } from 'svelte';
  import { page } from '$app/stores';
  import { goto } from '$app/navigation';
  import { api } from '$lib/api/client';
  import type { NotificationChannel, ChannelType } from '$lib/types/notification';
  import Button from '$lib/components/ui/Button.svelte';
  import Input from '$lib/components/ui/Input.svelte';

  let channelId = $page.params.id!;
  let channel: NotificationChannel | null = null;
  let isLoading = true;
  let error = '';
  let saveError = '';
  let saving = false;
  let successMessage = '';

  // Form fields
  let name = '';
  let isEnabled = true;

  // Provider-specific config fields
  let emailConfig = {
    smtp_host: '',
    smtp_port: '587',
    smtp_username: '',
    smtp_password: '',
    from_address: '',
    from_name: '',
    use_tls: true,
  };

  let slackConfig = {
    webhook_url: '',
    channel: '',
    username: 'Pulsar',
    icon_emoji: ':bell:',
  };

  let teamsConfig = {
    webhook_url: '',
    theme_color: '0078D4',
  };

  let webhookConfig = {
    url: '',
    method: 'POST',
    headers: {} as Record<string, string>,
    timeout: '30',
  };

  onMount(async () => {
    await loadChannel();
  });

  async function loadChannel() {
    try {
      isLoading = true;
      error = '';
      channel = await api.getNotificationChannel(channelId);

      // Populate form fields
      name = channel.name;
      isEnabled = channel.is_enabled;

      // Parse config based on channel type
      const config = (channel.config || {}) as Record<string, any>;
      switch (channel.channel_type) {
        case 'email':
          emailConfig = {
            smtp_host: config.smtp_host || '',
            smtp_port: String(config.smtp_port || 587),
            smtp_username: config.smtp_username || '',
            smtp_password: config.smtp_password || '',
            from_address: config.from_address || '',
            from_name: config.from_name || '',
            use_tls: config.use_tls !== false,
          };
          break;
        case 'slack':
          slackConfig = {
            webhook_url: config.webhook_url || '',
            channel: config.channel || '',
            username: config.username || 'Pulsar',
            icon_emoji: config.icon_emoji || ':bell:',
          };
          break;
        case 'teams':
          teamsConfig = {
            webhook_url: config.webhook_url || '',
            theme_color: config.theme_color || '0078D4',
          };
          break;
        case 'webhook':
          webhookConfig = {
            url: config.url || '',
            method: config.method || 'POST',
            headers: config.headers || {},
            timeout: String(config.timeout || 30),
          };
          break;
      }
    } catch (err) {
      error = err instanceof Error ? err.message : 'Failed to load notification channel';
    } finally {
      isLoading = false;
    }
  }

  async function handleSave() {
    if (!channel) return;

    saveError = '';
    saving = true;

    try {
      let config: Record<string, unknown> = {};

      switch (channel.channel_type) {
        case 'email':
          config = { ...emailConfig, smtp_port: parseInt(emailConfig.smtp_port, 10) };
          break;
        case 'slack':
          config = slackConfig;
          break;
        case 'teams':
          config = teamsConfig;
          break;
        case 'webhook':
          config = { ...webhookConfig, timeout: parseInt(webhookConfig.timeout, 10) };
          break;
      }

      await api.updateNotificationChannel(channelId, {
        name,
        is_enabled: isEnabled,
        config,
      });

      // Reload to get updated data
      await loadChannel();

      // Show success briefly
      successMessage = 'Channel saved successfully!';
      setTimeout(() => (successMessage = ''), 3000);
    } catch (err) {
      saveError = err instanceof Error ? err.message : 'Failed to save notification channel';
    } finally {
      saving = false;
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
  <title>{channel?.name || 'Configure Channel'} - Pulsar</title>
</svelte:head>

<div class="space-y-6">
  <!-- Header -->
  <div class="flex items-center gap-3">
    <button
      on:click={() => goto('/notifications/channels')}
      class="text-gray-600 hover:text-gray-900"
      aria-label="Back to notification channels"
    >
      <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
      </svg>
    </button>
    <div>
      <h2 class="text-3xl font-bold text-gray-900">
        {#if channel}
          <span class="mr-2">{getChannelTypeIcon(channel.channel_type)}</span>
          {channel.name}
        {:else}
          Loading...
        {/if}
      </h2>
      {#if channel}
        <p class="mt-1 text-gray-500">{getChannelTypeDisplay(channel.channel_type)}</p>
      {/if}
    </div>
  </div>

  {#if isLoading}
    <div class="text-center py-12">
      <div
        class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600"
      ></div>
      <p class="mt-2 text-gray-600">Loading channel configuration...</p>
    </div>
  {:else if error}
    <div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg">
      {error}
    </div>
  {:else if channel}
    <div class="bg-white p-6 rounded-lg shadow">
      <form on:submit|preventDefault={handleSave} class="space-y-6">
        <!-- Basic Settings -->
        <div class="space-y-4">
          <h3 class="text-lg font-semibold text-gray-900">Basic Settings</h3>

          <Input id="name" label="Channel Name" bind:value={name} required />

          <div class="flex items-center gap-2">
            <input
              id="is-enabled"
              type="checkbox"
              bind:checked={isEnabled}
              class="rounded bg-white border-gray-300 text-primary-600 focus:ring-primary-500"
            />
            <label for="is-enabled" class="text-sm font-medium text-gray-700">
              Enable this channel
            </label>
          </div>
        </div>

        <!-- Email Configuration -->
        {#if channel.channel_type === 'email'}
          <div class="space-y-4 p-4 bg-gray-50 rounded-lg border border-gray-200">
            <h4 class="text-md font-semibold text-gray-900">Email Configuration</h4>
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
              placeholder="Leave empty to keep existing"
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
                class="rounded bg-white border-gray-300 text-primary-600 focus:ring-primary-500"
              />
              <label for="use-tls" class="text-sm font-medium text-gray-700">Use TLS</label>
            </div>
          </div>
        {/if}

        <!-- Slack Configuration -->
        {#if channel.channel_type === 'slack'}
          <div class="space-y-4 p-4 bg-gray-50 rounded-lg border border-gray-200">
            <h4 class="text-md font-semibold text-gray-900">Slack Configuration</h4>
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
        {#if channel.channel_type === 'teams'}
          <div class="space-y-4 p-4 bg-gray-50 rounded-lg border border-gray-200">
            <h4 class="text-md font-semibold text-gray-900">Microsoft Teams Configuration</h4>
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
        {#if channel.channel_type === 'webhook'}
          <div class="space-y-4 p-4 bg-gray-50 rounded-lg border border-gray-200">
            <h4 class="text-md font-semibold text-gray-900">Webhook Configuration</h4>
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
                class="w-full px-3 py-2 bg-white border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500 text-gray-900"
              >
                <option value="POST">POST</option>
                <option value="PUT">PUT</option>
                <option value="PATCH">PATCH</option>
              </select>
            </div>
            <div>
              <label for="timeout" class="block text-sm font-medium text-gray-700 mb-1">
                Timeout (seconds)
              </label>
              <input
                id="timeout"
                type="number"
                bind:value={webhookConfig.timeout}
                min="1"
                max="300"
                class="w-full px-3 py-2 bg-white border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500 text-gray-900"
              />
            </div>
          </div>
        {/if}

        <!-- Error/Success Messages -->
        {#if saveError}
          <div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg">
            {saveError}
          </div>
        {/if}

        {#if successMessage}
          <div class="bg-green-50 border border-green-200 text-green-700 px-4 py-3 rounded-lg">
            {successMessage}
          </div>
        {/if}

        <!-- Action Buttons -->
        <div class="flex gap-3 pt-4 border-t border-gray-200">
          <Button type="submit" variant="primary" disabled={saving}>
            {saving ? 'Saving...' : 'Save Changes'}
          </Button>
          <Button
            type="button"
            variant="secondary"
            on:click={() => goto('/notifications/channels')}
          >
            Cancel
          </Button>
        </div>
      </form>
    </div>

    <!-- Channel Info -->
    <div class="bg-gray-50 p-4 rounded-lg border border-gray-200">
      <h4 class="text-sm font-semibold text-gray-700 mb-2">Channel Information</h4>
      <div class="text-sm text-gray-600 space-y-1">
        <p><strong>ID:</strong> {channel.id}</p>
        <p><strong>Created:</strong> {new Date(channel.created_at).toLocaleString()}</p>
        <p><strong>Last Updated:</strong> {new Date(channel.updated_at).toLocaleString()}</p>
      </div>
    </div>
  {/if}
</div>
