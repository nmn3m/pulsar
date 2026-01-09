<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { goto } from '$app/navigation';
  import { page } from '$app/stores';
  import { authStore } from '$lib/stores/auth';
  import { wsStore } from '$lib/stores/websocket';
  import Button from '$lib/components/ui/Button.svelte';

  let showSettingsMenu = false;

  $: if (!$authStore.isLoading && !$authStore.isAuthenticated) {
    wsStore.disconnect();
    goto('/login');
  }

  $: currentPath = $page.url.pathname;

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
        return 'bg-success shadow-glow-success';
      case 'connecting':
        return 'bg-warning animate-pulse';
      case 'disconnected':
        return 'bg-gray-500';
      case 'error':
        return 'bg-error shadow-glow-error';
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

  function isActive(path: string): boolean {
    return currentPath.startsWith(path);
  }

  function handleClickOutside(event: MouseEvent) {
    const target = event.target as HTMLElement;
    if (!target.closest('.settings-menu')) {
      showSettingsMenu = false;
    }
  }
</script>

<svelte:window on:click={handleClickOutside} />

{#if $authStore.isLoading}
  <div class="min-h-screen flex items-center justify-center">
    <div class="text-center">
      <div
        class="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-primary-500"
      ></div>
      <p class="mt-4 text-gray-500">Loading...</p>
    </div>
  </div>
{:else if $authStore.isAuthenticated}
  <div class="min-h-screen">
    <!-- Header -->
    <header
      class="bg-white/80 backdrop-blur-md border-b border-gray-200 shadow-lg relative z-40"
    >
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between items-center h-16 overflow-visible">
          <div class="flex items-center">
            <h1
              class="text-2xl font-bold text-primary-600"
            >
              Pulsar
            </h1>
            <nav class="ml-10 flex space-x-1 overflow-visible">
              <a
                href="/dashboard"
                class="px-3 py-2 rounded-lg text-sm font-medium transition-all duration-200 {isActive('/dashboard')
                  ? 'bg-primary-100 text-primary-700'
                  : 'text-gray-600 hover:text-primary-600 hover:bg-gray-100'}"
              >
                Dashboard
              </a>
              <a
                href="/alerts"
                class="px-3 py-2 rounded-lg text-sm font-medium transition-all duration-200 {isActive('/alerts')
                  ? 'bg-red-100 text-red-700'
                  : 'text-gray-600 hover:text-error hover:bg-gray-100'}"
              >
                Alerts
              </a>
              <a
                href="/incidents"
                class="px-3 py-2 rounded-lg text-sm font-medium transition-all duration-200 {isActive('/incidents')
                  ? 'bg-red-100 text-red-700'
                  : 'text-gray-600 hover:text-error hover:bg-gray-100'}"
              >
                Incidents
              </a>
              <a
                href="/teams"
                class="px-3 py-2 rounded-lg text-sm font-medium transition-all duration-200 {isActive('/teams')
                  ? 'bg-primary-100 text-primary-700'
                  : 'text-gray-600 hover:text-primary-600 hover:bg-gray-100'}"
              >
                Teams
              </a>
              <a
                href="/schedules"
                class="px-3 py-2 rounded-lg text-sm font-medium transition-all duration-200 {isActive('/schedules')
                  ? 'bg-primary-100 text-primary-700'
                  : 'text-gray-600 hover:text-primary-600 hover:bg-gray-100'}"
              >
                Schedules
              </a>
              <a
                href="/escalation-policies"
                class="px-3 py-2 rounded-lg text-sm font-medium transition-all duration-200 {isActive('/escalation-policies')
                  ? 'bg-primary-100 text-primary-700'
                  : 'text-gray-600 hover:text-primary-600 hover:bg-gray-100'}"
              >
                Escalations
              </a>

              <!-- Settings Dropdown -->
              <div class="relative settings-menu">
                <button
                  type="button"
                  on:click|stopPropagation={() => (showSettingsMenu = !showSettingsMenu)}
                  class="px-3 py-2 rounded-lg text-sm font-medium transition-all duration-200 flex items-center gap-1 {isActive('/webhooks') || isActive('/notifications') || isActive('/settings')
                    ? 'bg-primary-100 text-primary-700'
                    : 'text-gray-600 hover:text-primary-600 hover:bg-gray-100'}"
                >
                  Settings
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
                  </svg>
                </button>

                {#if showSettingsMenu}
                  <div class="absolute top-full right-0 mt-1 w-56 bg-white rounded-lg shadow-xl border border-gray-200 py-1 z-[9999]">
                    <a
                      href="/webhooks/endpoints"
                      class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 {isActive('/webhooks') ? 'bg-gray-50' : ''}"
                      on:click={() => (showSettingsMenu = false)}
                    >
                      <div class="font-medium">Webhooks</div>
                      <div class="text-xs text-gray-500">Incoming & outgoing webhooks</div>
                    </a>
                    <a
                      href="/notifications/channels"
                      class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 {isActive('/notifications') ? 'bg-gray-50' : ''}"
                      on:click={() => (showSettingsMenu = false)}
                    >
                      <div class="font-medium">Notifications</div>
                      <div class="text-xs text-gray-500">Channels & preferences</div>
                    </a>
                    <a
                      href="/settings/api-keys"
                      class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 {isActive('/settings/api-keys') ? 'bg-gray-50' : ''}"
                      on:click={() => (showSettingsMenu = false)}
                    >
                      <div class="font-medium">API Keys</div>
                      <div class="text-xs text-gray-500">Manage programmatic access</div>
                    </a>
                  </div>
                {/if}
              </div>
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
              <span class="text-xs text-gray-500 hidden sm:inline">
                {getStatusText($wsStore.status)}
              </span>
            </div>

            <span class="text-sm text-primary-600">
              {$authStore.user?.email}
            </span>
            <Button variant="secondary" on:click={handleLogout}>Logout</Button>
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
