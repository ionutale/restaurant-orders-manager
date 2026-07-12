<script lang="ts">
	import { auth } from '$lib/stores/auth.svelte';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { API_BASE } from '$lib/config';

	if (auth.role !== 'waiter') goto('/login');

	type Order = {
		id: number;
		table_group_id: number;
		status: string;
		created_at: string;
		courses: { id: number; name: string; status: string }[];
	};

	let orders = $state<Order[]>([]);
	let tableGroups = $state<{ id: number; name: string | null }[]>([]);
	let loading = $state(true);
	let showCreate = $state(false);
	let selGroupId = $state(0);

	const token = () => localStorage.getItem('token') ?? '';

	async function load() {
		const [oRes, tRes] = await Promise.all([
			fetch(`${API_BASE}/orders`, { headers: { Authorization: `Bearer ${token()}` } }),
			fetch(`${API_BASE}/table-groups`, { headers: { Authorization: `Bearer ${token()}` } }),
		]);
		if (oRes.ok) orders = await oRes.json();
		if (tRes.ok) tableGroups = await tRes.json();
		loading = false;
	}

	onMount(load);

	async function createOrder() {
		if (!selGroupId) return;
		const res = await fetch(`${API_BASE}/orders`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token()}` },
			body: JSON.stringify({ table_group_id: selGroupId }),
		});
		if (res.ok) { showCreate = false; selGroupId = 0; const o = await res.json(); goto(`/waiter/orders/${o.id}`); }
	}

	function statusBadge(s: string) {
		const m: Record<string, string> = { pending: 'badge-ghost', sent: 'badge-info', completed: 'badge-success', paid: 'badge-primary' };
		return m[s] || 'badge-ghost';
	}
</script>

<div class="space-y-4">
	<div class="flex items-center justify-between">
		<h2 class="text-2xl font-bold">Orders</h2>
		<button class="btn btn-primary btn-sm" onclick={() => (showCreate = true)}>New Order</button>
	</div>

	{#if loading}
		<div class="flex justify-center py-8"><span class="loading loading-spinner loading-lg"></span></div>
	{:else if orders.length === 0}
		<div class="flex h-32 items-center justify-center rounded-box border-2 border-dashed text-base-content/40">No orders yet</div>
	{:else}
		<div class="overflow-x-auto">
			<table class="table table-zebra">
				<thead><tr><th>Order #</th><th>Courses</th><th>Status</th><th>Created</th></tr></thead>
				<tbody>
					{#each orders as o (o.id)}
						<tr class="cursor-pointer" onclick={() => goto(`/waiter/orders/${o.id}`)}>
							<td class="font-medium">#{o.id}</td>
							<td>{o.courses?.map(c => c.name).join(' → ') || '—'}</td>
							<td><span class="badge {statusBadge(o.status)}">{o.status}</span></td>
							<td class="text-sm text-base-content/50">{new Date(o.created_at).toLocaleString()}</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	{/if}
</div>

{#if showCreate}
	<div class="modal modal-open">
		<div class="modal-box">
			<h3 class="font-bold text-lg">New Order</h3>
			<div class="py-4 space-y-3">
				<label class="form-control">
					<div class="label"><span class="label-text">Table Group</span></div>
					<select bind:value={selGroupId} class="select select-bordered">
						<option value={0} disabled>Select group...</option>
						{#each tableGroups.filter(g => g.name) as g}
							<option value={g.id}>{g.name}</option>
						{/each}
					</select>
				</label>
			</div>
			<div class="modal-action">
				<button class="btn" onclick={() => (showCreate = false)}>Cancel</button>
				<button class="btn btn-primary" onclick={createOrder}>Create Order</button>
			</div>
		</div>
	</div>
{/if}
