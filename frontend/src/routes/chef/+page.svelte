<script lang="ts">
	import { auth } from '$lib/stores/auth.svelte';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { API_BASE } from '$lib/config';

	if (auth.role !== 'chef') goto('/login');

	type Suggestion = {
		id: number;
		name: string;
		description: string;
		price_cents: number;
		shift_date: string;
		expires_at: string;
		chef_name: string;
	};

	let activeSuggestions = $state<Suggestion[]>([]);
	let expiredSuggestions = $state<Suggestion[]>([]);
	let loading = $state(true);
	let showForm = $state(false);
	let showExpired = $state(false);
	let newName = $state('');
	let newDesc = $state('');
	let newPrice = $state(0);
	let deleteId = $state<number | null>(null);
	let renewing = $state<number | null>(null);
	let error = $state('');

	const token = () => localStorage.getItem('token') ?? '';

	async function load() {
		const [activeRes, allRes] = await Promise.all([
			fetch(`${API_BASE}/chef-suggestions`, { headers: { Authorization: `Bearer ${token()}` } }),
			fetch(`${API_BASE}/chef-suggestions?all=true`, { headers: { Authorization: `Bearer ${token()}` } }),
		]);
		if (activeRes.ok) activeSuggestions = await activeRes.json();
		if (allRes.ok) {
			const all = await allRes.json();
			const activeIds = new Set(activeSuggestions.map(s => s.id));
			expiredSuggestions = all.filter((s: Suggestion) => !activeIds.has(s.id));
		}
		loading = false;
	}

	onMount(load);

	async function create() {
		if (!newName) return;
		error = '';
		const expiresAt = new Date(); expiresAt.setHours(expiresAt.getHours() + 8);
		const res = await fetch(`${API_BASE}/chef-suggestions`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token()}` },
			body: JSON.stringify({ name: newName, description: newDesc, price_cents: Math.round(newPrice * 100), expires_at: expiresAt.toISOString() }),
		});
		if (!res.ok) { error = 'Failed to create'; return; }
		newName = ''; newDesc = ''; newPrice = 0; showForm = false;
		await load();
	}

	async function remove(id: number) {
		await fetch(`${API_BASE}/chef-suggestions/${id}`, { method: 'DELETE', headers: { Authorization: `Bearer ${token()}` } });
		deleteId = null;
		await load();
	}

	async function renew(id: number) {
		renewing = id;
		await fetch(`${API_BASE}/chef-suggestions/${id}/renew`, { method: 'POST', headers: { Authorization: `Bearer ${token()}` } });
		renewing = null;
		await load();
	}

	function price(cents: number) { return `€${(cents / 100).toFixed(2)}`; }
</script>

<div class="space-y-4">
	<div class="flex items-center justify-between">
		<h2 class="text-2xl font-bold">KDS — Chef Suggestions</h2>
		<button class="btn btn-primary btn-sm" onclick={() => (showForm = true)}>New Suggestion</button>
	</div>

	{#if error}<div class="alert alert-error">{error}</div>{/if}

	{#if loading}
		<div class="flex justify-center py-8"><span class="loading loading-spinner loading-lg"></span></div>
	{:else}
		<h3 class="font-semibold">Active Suggestions</h3>
		{#if activeSuggestions.length === 0}
			<div class="flex h-24 items-center justify-center rounded-box border-2 border-dashed text-base-content/40 mb-4">No active suggestions</div>
		{:else}
			<div class="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3 mb-6">
				{#each activeSuggestions as s (s.id)}
					<div class="card bg-base-100 shadow-xl">
						<div class="card-body">
							<div class="flex items-start justify-between">
								<h4 class="card-title text-lg">{s.name}</h4>
								{#if deleteId === s.id}
									<div class="flex gap-1">
										<button class="btn btn-ghost btn-xs text-error" onclick={() => remove(s.id)}>Confirm</button>
										<button class="btn btn-ghost btn-xs" onclick={() => (deleteId = null)}>Cancel</button>
									</div>
								{:else}
									<button class="btn btn-ghost btn-xs text-error" onclick={() => (deleteId = s.id)}>✕</button>
								{/if}
							</div>
							{#if s.description}<p class="text-sm text-base-content/70">{s.description}</p>{/if}
							<div class="mt-2 flex items-center justify-between text-sm">
								<span class="badge badge-primary">{price(s.price_cents)}</span>
								<span class="text-base-content/50">by {s.chef_name}</span>
							</div>
						</div>
					</div>
				{/each}
			</div>
		{/if}

		<button class="btn btn-ghost btn-sm mb-4" onclick={() => (showExpired = !showExpired)}>
			{showExpired ? 'Hide' : 'Show'} Expired ({expiredSuggestions.length})
		</button>

		{#if showExpired}
			<h3 class="font-semibold">Expired Suggestions</h3>
			{#if expiredSuggestions.length === 0}
				<div class="flex h-24 items-center justify-center rounded-box border-2 border-dashed text-base-content/40">No expired suggestions</div>
			{:else}
				<div class="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3">
					{#each expiredSuggestions as s (s.id)}
						<div class="card bg-base-200 shadow-sm opacity-70">
							<div class="card-body">
								<div class="flex items-start justify-between">
									<h4 class="card-title text-base">{s.name}</h4>
									<div class="flex gap-1">
										<button class="btn btn-ghost btn-xs" onclick={() => renew(s.id)} disabled={renewing === s.id}>
											{renewing === s.id ? '...' : 'Renew'}
										</button>
										<button class="btn btn-ghost btn-xs text-error" onclick={() => { deleteId = s.id; remove(s.id); }}>✕</button>
									</div>
								</div>
								{#if s.description}<p class="text-xs text-base-content/60">{s.description}</p>{/if}
								<div class="mt-1 flex items-center justify-between text-xs">
									<span class="badge badge-ghost badge-sm">{price(s.price_cents)}</span>
									<span class="text-base-content/40">by {s.chef_name}</span>
								</div>
							</div>
						</div>
					{/each}
				</div>
			{/if}
		{/if}
	{/if}
</div>

{#if showForm}
	<div class="modal modal-open">
		<div class="modal-box">
			<h3 class="font-bold text-lg">New Chef Suggestion</h3>
			<div class="py-4 space-y-3">
				<label class="form-control"><div class="label"><span class="label-text">Name</span></div><input type="text" bind:value={newName} class="input input-bordered" /></label>
				<label class="form-control"><div class="label"><span class="label-text">Description</span></div><textarea bind:value={newDesc} class="textarea textarea-bordered" rows="2"></textarea></label>
				<label class="form-control"><div class="label"><span class="label-text">Price (EUR)</span></div><input type="number" bind:value={newPrice} class="input input-bordered" min="0" step="0.01" /></label>
			</div>
			<div class="modal-action">
				<button class="btn" onclick={() => (showForm = false)}>Cancel</button>
				<button class="btn btn-primary" onclick={create}>Create</button>
			</div>
		</div>
	</div>
{/if}
