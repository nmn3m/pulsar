<script lang="ts">
	import { goto } from '$app/navigation';
	import Button from '$lib/components/ui/Button.svelte';
	import Input from '$lib/components/ui/Input.svelte';
	import { authStore } from '$lib/stores/auth';

	let email = '';
	let username = '';
	let password = '';
	let fullName = '';
	let organizationName = '';
	let error = '';
	let loading = false;

	async function handleRegister() {
		error = '';
		loading = true;

		try {
			await authStore.register({
				email,
				username,
				password,
				full_name: fullName,
				organization_name: organizationName
			});
			goto('/dashboard');
		} catch (err) {
			error = err instanceof Error ? err.message : 'Registration failed';
		} finally {
			loading = false;
		}
	}
</script>

<svelte:head>
	<title>Register - Pulsar</title>
</svelte:head>

<div class="min-h-screen flex items-center justify-center bg-gray-50 px-4 py-12">
	<div class="max-w-md w-full space-y-8">
		<div class="text-center">
			<h1 class="text-4xl font-bold text-gray-900">Pulsar</h1>
			<p class="mt-2 text-gray-600">Create your account</p>
		</div>

		<div class="bg-white p-8 rounded-lg shadow-md">
			<form on:submit|preventDefault={handleRegister} class="space-y-6">
				<Input
					id="email"
					type="email"
					label="Email address"
					bind:value={email}
					placeholder="you@example.com"
					required
				/>

				<Input
					id="username"
					type="text"
					label="Username"
					bind:value={username}
					placeholder="johndoe"
					required
				/>

				<Input
					id="fullName"
					type="text"
					label="Full name"
					bind:value={fullName}
					placeholder="John Doe"
				/>

				<Input
					id="organizationName"
					type="text"
					label="Organization name"
					bind:value={organizationName}
					placeholder="Acme Inc."
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
					<div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg">
						{error}
					</div>
				{/if}

				<Button type="submit" variant="primary" fullWidth disabled={loading}>
					{loading ? 'Creating account...' : 'Create account'}
				</Button>
			</form>

			<div class="mt-6 text-center">
				<p class="text-sm text-gray-600">
					Already have an account?
					<a href="/login" class="text-primary-600 hover:text-primary-700 font-medium">
						Sign in
					</a>
				</p>
			</div>
		</div>
	</div>
</div>
