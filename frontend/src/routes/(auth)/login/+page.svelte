<script lang="ts">
  import { goto } from '$app/navigation';
  import Button from '$lib/components/ui/Button.svelte';
  import Input from '$lib/components/ui/Input.svelte';
  import { authStore } from '$lib/stores/auth';

  let email = '';
  let password = '';
  let error = '';
  let loading = false;

  async function handleLogin() {
    error = '';
    loading = true;

    try {
      const response = await authStore.login({ email, password });
      if (response.requires_email_verification) {
        goto('/verify-email');
      } else {
        goto('/dashboard');
      }
    } catch (err) {
      error = err instanceof Error ? err.message : 'Login failed';
    } finally {
      loading = false;
    }
  }
</script>

<svelte:head>
  <title>Login - Pulsar</title>
</svelte:head>

<div class="min-h-screen flex items-center justify-center px-4 py-12">
  <div class="max-w-md w-full space-y-8">
    <div class="text-center">
      <h1 class="text-4xl font-bold text-primary-600">
        Pulsar
      </h1>
      <p class="mt-2 text-gray-500">Sign in to your account</p>
    </div>

    <div
      class="bg-white backdrop-blur-sm p-8 rounded-xl border border-gray-200 shadow-lg"
    >
      <form on:submit|preventDefault={handleLogin} class="space-y-6">
        <Input
          id="email"
          type="email"
          label="Email address"
          bind:value={email}
          placeholder="you@example.com"
          required
        />

        <Input
          id="password"
          type="password"
          label="Password"
          bind:value={password}
          placeholder="••••••••"
          required
        />

        {#if error}
          <div
            class="bg-red-50 border border-red-200 text-red-600 px-4 py-3 rounded-lg"
          >
            {error}
          </div>
        {/if}

        <Button type="submit" variant="primary" fullWidth disabled={loading}>
          {loading ? 'Signing in...' : 'Sign in'}
        </Button>
      </form>

      <div class="mt-6 text-center">
        <p class="text-sm text-gray-500">
          Don't have an account?
          <a
            href="/register"
            class="text-primary-600 hover:text-primary-500 font-medium"
          >
            Sign up
          </a>
        </p>
      </div>
    </div>
  </div>
</div>
