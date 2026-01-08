<script lang="ts">
  import { goto } from '$app/navigation';
  import { onMount } from 'svelte';
  import Button from '$lib/components/ui/Button.svelte';
  import { authStore } from '$lib/stores/auth';

  let otp = ['', '', '', '', '', ''];
  let error = '';
  let success = '';
  let loading = false;
  let resendLoading = false;
  let resendCooldown = 0;
  let email = '';

  let inputs: HTMLInputElement[] = [];

  onMount(() => {
    // Get email from auth store
    const unsubscribe = authStore.subscribe((state) => {
      if (state.pendingVerificationEmail) {
        email = state.pendingVerificationEmail;
      } else if (!state.isLoading && !state.pendingVerificationEmail) {
        // No pending verification, redirect to login
        goto('/login');
      }
    });

    // Focus first input
    if (inputs[0]) {
      inputs[0].focus();
    }

    return unsubscribe;
  });

  function handleInput(index: number, event: Event) {
    const input = event.target as HTMLInputElement;
    const value = input.value;

    // Only allow digits
    if (value && !/^\d$/.test(value)) {
      otp[index] = '';
      return;
    }

    otp[index] = value;

    // Move to next input if value entered
    if (value && index < 5) {
      inputs[index + 1]?.focus();
    }

    // Auto-submit when all digits entered
    if (otp.every((digit) => digit !== '')) {
      handleVerify();
    }
  }

  function handleKeyDown(index: number, event: KeyboardEvent) {
    // Handle backspace
    if (event.key === 'Backspace' && !otp[index] && index > 0) {
      inputs[index - 1]?.focus();
    }
  }

  function handlePaste(event: ClipboardEvent) {
    event.preventDefault();
    const pastedData = event.clipboardData?.getData('text') || '';
    const digits = pastedData.replace(/\D/g, '').slice(0, 6).split('');

    digits.forEach((digit, index) => {
      if (index < 6) {
        otp[index] = digit;
      }
    });

    // Focus last filled input or next empty one
    const lastIndex = Math.min(digits.length, 5);
    inputs[lastIndex]?.focus();

    // Auto-submit if all digits pasted
    if (digits.length === 6) {
      handleVerify();
    }
  }

  async function handleVerify() {
    if (otp.some((digit) => !digit)) {
      error = 'Please enter all 6 digits';
      return;
    }

    error = '';
    success = '';
    loading = true;

    try {
      await authStore.verifyEmail({
        email,
        otp: otp.join(''),
      });
      success = 'Email verified successfully! Redirecting...';
      setTimeout(() => {
        goto('/login');
      }, 1500);
    } catch (err) {
      error = err instanceof Error ? err.message : 'Verification failed';
      // Clear OTP on error
      otp = ['', '', '', '', '', ''];
      inputs[0]?.focus();
    } finally {
      loading = false;
    }
  }

  async function handleResend() {
    if (resendCooldown > 0) return;

    error = '';
    success = '';
    resendLoading = true;

    try {
      await authStore.resendOTP({ email });
      success = 'Verification code sent!';
      resendCooldown = 60;

      // Start cooldown timer
      const interval = setInterval(() => {
        resendCooldown -= 1;
        if (resendCooldown <= 0) {
          clearInterval(interval);
        }
      }, 1000);
    } catch (err) {
      error = err instanceof Error ? err.message : 'Failed to resend code';
    } finally {
      resendLoading = false;
    }
  }

  function maskEmail(email: string): string {
    if (!email) return '';
    const [localPart, domain] = email.split('@');
    if (localPart.length <= 2) {
      return `${localPart[0]}*@${domain}`;
    }
    return `${localPart[0]}${'*'.repeat(localPart.length - 2)}${localPart[localPart.length - 1]}@${domain}`;
  }
</script>

<svelte:head>
  <title>Verify Email - Pulsar</title>
</svelte:head>

<div class="min-h-screen flex items-center justify-center px-4 py-12">
  <div class="max-w-md w-full space-y-8">
    <div class="text-center">
      <h1 class="text-4xl font-bold text-primary-600 dark:text-primary-400 dark:text-glow-cyan">
        Pulsar
      </h1>
      <p class="mt-2 text-gray-500 dark:text-gray-400">Verify your email</p>
    </div>

    <div
      class="bg-white dark:bg-space-800/50 backdrop-blur-sm p-8 rounded-xl border border-gray-200 dark:border-space-600 shadow-lg"
    >
      <div class="text-center mb-8">
        <div
          class="mx-auto w-16 h-16 bg-primary-100 dark:bg-primary-900/30 rounded-full flex items-center justify-center mb-4"
        >
          <svg
            class="w-8 h-8 text-primary-600 dark:text-primary-400"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"
            />
          </svg>
        </div>
        <p class="text-gray-600 dark:text-gray-300">We've sent a 6-digit verification code to</p>
        <p class="font-medium text-gray-900 dark:text-white mt-1">
          {maskEmail(email)}
        </p>
      </div>

      <form on:submit|preventDefault={handleVerify} class="space-y-6">
        <div>
          <label
            class="block text-sm font-medium text-gray-700 dark:text-gray-300 text-center mb-4"
          >
            Enter verification code
          </label>
          <div class="flex justify-center gap-2" on:paste={handlePaste}>
            {#each otp as digit, index}
              <input
                bind:this={inputs[index]}
                type="text"
                inputmode="numeric"
                maxlength="1"
                value={digit}
                on:input={(e) => handleInput(index, e)}
                on:keydown={(e) => handleKeyDown(index, e)}
                class="w-12 h-14 text-center text-2xl font-bold border-2 rounded-lg
									   bg-white dark:bg-space-700
									   border-gray-300 dark:border-space-500
									   focus:border-primary-500 dark:focus:border-primary-400
									   focus:ring-2 focus:ring-primary-500/20
									   text-gray-900 dark:text-white
									   transition-colors"
                disabled={loading}
              />
            {/each}
          </div>
        </div>

        {#if error}
          <div
            class="bg-red-50 dark:bg-accent-900/30 border border-red-200 dark:border-accent-500/50 text-red-600 dark:text-accent-300 px-4 py-3 rounded-lg text-center"
          >
            {error}
          </div>
        {/if}

        {#if success}
          <div
            class="bg-green-50 dark:bg-green-900/30 border border-green-200 dark:border-green-500/50 text-green-600 dark:text-green-300 px-4 py-3 rounded-lg text-center"
          >
            {success}
          </div>
        {/if}

        <Button type="submit" variant="primary" fullWidth disabled={loading || otp.some((d) => !d)}>
          {loading ? 'Verifying...' : 'Verify Email'}
        </Button>
      </form>

      <div class="mt-6 text-center">
        <p class="text-sm text-gray-500 dark:text-gray-400">
          Didn't receive the code?
          <button
            on:click={handleResend}
            class="text-primary-600 dark:text-primary-400 hover:text-primary-500 dark:hover:text-primary-300 font-medium disabled:opacity-50 disabled:cursor-not-allowed"
            disabled={resendLoading || resendCooldown > 0}
          >
            {#if resendLoading}
              Sending...
            {:else if resendCooldown > 0}
              Resend in {resendCooldown}s
            {:else}
              Resend code
            {/if}
          </button>
        </p>
      </div>

      <div class="mt-4 text-center">
        <a
          href="/login"
          class="text-sm text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-300"
        >
          Back to login
        </a>
      </div>
    </div>
  </div>
</div>
