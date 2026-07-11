<script lang="ts">
	import { auth } from '$lib/stores/auth.svelte';
	import { goto } from '$app/navigation';

	let email = $state('');
	let password = $state('');
	let error = $state('');
	let loading = $state(false);

	async function handleLogin(e: Event) {
		e.preventDefault();
		error = '';
		loading = true;
		try {
			await auth.login(email, password);
			if (auth.role === 'manager') goto('/admin');
			else if (auth.role === 'waiter') goto('/waiter');
			else if (auth.role === 'chef') goto('/chef');
		} catch (err: any) {
			error = err.message || 'Login failed';
		} finally {
			loading = false;
		}
	}
</script>

<div class="hero min-h-screen bg-base-200">
	<div class="hero-content w-full max-w-sm flex-col">
		<div class="text-center">
			<h1 class="text-4xl font-bold">Restaurant Orders</h1>
			<p class="py-4 text-base-content/70">Sign in to continue</p>
		</div>
		<form class="card w-full bg-base-100 shadow-xl" onsubmit={handleLogin}>
			<div class="card-body gap-4">
				{#if error}
					<div class="alert alert-error">{error}</div>
				{/if}
				<label class="form-control w-full">
					<div class="label"><span class="label-text">Email</span></div>
					<input type="email" bind:value={email} placeholder="email@example.com" class="input input-bordered w-full" required />
				</label>
				<label class="form-control w-full">
					<div class="label"><span class="label-text">Password</span></div>
					<input type="password" bind:value={password} placeholder="••••••••" class="input input-bordered w-full" required />
				</label>
				<button type="submit" class="btn btn-primary mt-2" disabled={loading}>
					{#if loading}<span class="loading loading-spinner"></span>{/if}
					Sign In
				</button>
			</div>
		</form>
	</div>
</div>
