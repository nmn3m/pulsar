<script lang="ts">
  import { onMount } from 'svelte';
  import { authStore } from '$lib/stores/auth';
  import { api } from '$lib/api/client';
  import Button from '$lib/components/ui/Button.svelte';
  import Input from '$lib/components/ui/Input.svelte';

  let fullName = '';
  let phone = '';
  let timezone = '';
  let saving = false;
  let saveError = '';
  let successMessage = '';

  onMount(() => {
    if ($authStore.user) {
      fullName = $authStore.user.full_name || '';
      phone = $authStore.user.phone || '';
      timezone = $authStore.user.timezone || '';
    }
  });

  async function handleSave() {
    saveError = '';
    successMessage = '';
    saving = true;

    try {
      const updated = await api.updateProfile({
        full_name: fullName || undefined,
        phone: phone || undefined,
        timezone: timezone || undefined,
      });

      // Update the auth store with fresh data
      authStore.setUser(updated);

      successMessage = 'Profile updated successfully!';
      setTimeout(() => (successMessage = ''), 3000);
    } catch (err) {
      saveError = err instanceof Error ? err.message : 'Failed to update profile';
    } finally {
      saving = false;
    }
  }
</script>

<svelte:head>
  <title>Profile Settings - Pulsar</title>
</svelte:head>

<div class="space-y-6">
  <div>
    <h2 class="text-3xl font-bold text-gray-900">Profile Settings</h2>
    <p class="mt-2 text-gray-500">Manage your personal information and contact details</p>
  </div>

  <div class="bg-white p-6 rounded-xl border border-gray-200 shadow-sm">
    <form on:submit|preventDefault={handleSave} class="space-y-6">
      <div class="space-y-4">
        <div>
          <label for="email" class="block text-sm font-medium text-gray-600 mb-1">Email</label>
          <input
            id="email"
            type="email"
            value={$authStore.user?.email || ''}
            disabled
            class="w-full px-3 py-2 bg-gray-100 border border-gray-300 rounded-lg text-gray-500 cursor-not-allowed"
          />
          <p class="mt-1 text-xs text-gray-400">Email cannot be changed</p>
        </div>

        <div>
          <label for="username" class="block text-sm font-medium text-gray-600 mb-1">Username</label
          >
          <input
            id="username"
            type="text"
            value={$authStore.user?.username || ''}
            disabled
            class="w-full px-3 py-2 bg-gray-100 border border-gray-300 rounded-lg text-gray-500 cursor-not-allowed"
          />
        </div>

        <Input
          id="full-name"
          label="Full Name"
          bind:value={fullName}
          placeholder="Your full name"
        />

        <div>
          <Input
            id="phone"
            label="Phone Number (E.164 format)"
            bind:value={phone}
            placeholder="+1234567890"
          />
          <p class="mt-1 text-xs text-gray-400">
            Required for SMS notifications. Use international format, e.g. +20xxxxxxxxxx for Egypt,
            +1xxxxxxxxxx for US.
          </p>
        </div>

        <Input id="timezone" label="Timezone" bind:value={timezone} placeholder="Africa/Cairo" />
      </div>

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

      <div class="pt-4 border-t border-gray-200">
        <Button type="submit" variant="primary" disabled={saving}>
          {saving ? 'Saving...' : 'Save Changes'}
        </Button>
      </div>
    </form>
  </div>
</div>
