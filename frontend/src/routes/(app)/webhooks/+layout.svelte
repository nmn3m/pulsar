<script lang="ts">
  import { page } from '$app/stores';

  $: currentPath = $page.url.pathname;
</script>

<div class="space-y-6">
  <!-- Sub-navigation tabs -->
  <div class="border-b border-gray-200">
    <nav class="-mb-px flex space-x-8">
      <a
        href="/webhooks/endpoints"
        class="border-b-2 py-4 px-1 text-sm font-medium {currentPath === '/webhooks/endpoints'
          ? 'border-primary-500 text-primary-600'
          : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'}"
      >
        Outgoing Endpoints
      </a>
      <a
        href="/webhooks/incoming"
        class="border-b-2 py-4 px-1 text-sm font-medium {currentPath === '/webhooks/incoming'
          ? 'border-primary-500 text-primary-600'
          : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'}"
      >
        Incoming Webhooks
      </a>
      <a
        href="/webhooks/deliveries"
        class="border-b-2 py-4 px-1 text-sm font-medium {currentPath === '/webhooks/deliveries'
          ? 'border-primary-500 text-primary-600'
          : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'}"
      >
        Delivery Logs
      </a>
    </nav>
  </div>

  <!-- Page content -->
  <slot />
</div>
