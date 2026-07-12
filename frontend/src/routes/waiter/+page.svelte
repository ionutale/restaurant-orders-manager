<script lang="ts">
	import { auth } from '$lib/stores/auth.svelte';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { API_BASE } from '$lib/config';
	import FloorPlanCanvas from '$lib/components/FloorPlanCanvas.svelte';
	import type { FloorPlanTable as Table } from '$lib/types';

	if (auth.role !== 'waiter') goto('/login');

	type Group = { id: number; name: string | null; party_size: number; status: string; table_ids: number[] };
	type Notif = { dish_name: string; quantity: number; table_name: string; order_id: number };
	type Pred = { table_id: number; table_name: string; estimated_free?: string };

	let tables = $state<Table[]>([]);
	let notifications = $state<Notif[]>([]);
	let predictions = $state<Pred[]>([]);
	let loading = $state(true);
	let viewMode = $state<'canvas' | 'list'>('canvas');
	let selectedTable = $state<Table | null>(null);
	let selectedGroup = $state<Group | null>(null);
	let seating = $state(false);
	let merging = $state(false);
	let seatName = $state('');
	let seatPartySize = $state(2);
	let seatTargetIds = $state<number[]>([]);
	let mergeTargetIds = $state<number[]>([]);
	let seatingError = $state('');

	const token = () => localStorage.getItem('token') ?? '';

	async function load() {
		const [fRes, nRes, pRes] = await Promise.all([
			fetch(`${API_BASE}/floor-plan`, { headers: { Authorization: `Bearer ${token()}` } }),
			fetch(`${API_BASE}/notifications`, { headers: { Authorization: `Bearer ${token()}` } }),
			fetch(`${API_BASE}/predictions`, { headers: { Authorization: `Bearer ${token()}` } }),
		]);
		if (fRes.ok) tables = await fRes.json();
		if (nRes.ok) notifications = await nRes.json();
		if (pRes.ok) predictions = await pRes.json();
		loading = false;
	}

	onMount(load);
	setInterval(load, 15000);

	async function seatAndOrder() {
		seatingError = '';
		const res = await fetch(`${API_BASE}/start-order`, {
			method: 'POST', headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token()}` },
			body: JSON.stringify({ table_ids: seatTargetIds, party_size: seatPartySize, name: seatName || undefined }),
		});
		if (!res.ok) { seatingError = 'Failed'; return; }
		const order = await res.json();
		goto(`/waiter/orders/${order.id}`);
	}

	async function loadGroup(id: number) {
		const r = await fetch(`${API_BASE}/table-groups/${id}`, { headers: { Authorization: `Bearer ${token()}` } });
		if (r.ok) selectedGroup = await r.json();
	}
	async function closeGroup(id: number) {
		await fetch(`${API_BASE}/table-groups/${id}/close`, { method: 'POST', headers: { Authorization: `Bearer ${token()}` } });
		selectedTable = null; selectedGroup = null; await load();
	}
	async function mergeTables() {
		if (!selectedGroup) return;
		await fetch(`${API_BASE}/table-groups/${selectedGroup.id}/tables`, {
			method: 'PATCH', headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token()}` },
			body: JSON.stringify({ add_table_ids: mergeTargetIds, remove_table_ids: [] }),
		});
		merging = false; mergeTargetIds = []; await loadGroup(selectedGroup.id); await load();
	}
	async function splitTable(tid: number) {
		if (!selectedGroup) return;
		await fetch(`${API_BASE}/table-groups/${selectedGroup.id}/tables`, {
			method: 'PATCH', headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token()}` },
			body: JSON.stringify({ add_table_ids: [], remove_table_ids: [tid] }),
		});
		await loadGroup(selectedGroup.id); await load();
	}

	function handleTableClick(t: Table) {
		selectedTable = t;
		if (t.status === 'free') {
			seating = true; seatingError = ''; seatName = '';
			seatPartySize = Math.min(2, t.capacity);
			seatTargetIds = [t.id];
		} else {
			seating = false;
			if (t.group_id) loadGroup(t.group_id);
		}
	}

	function freeTables() { return tables.filter(t => t.status === 'free'); }
	function occupiedTables() { return tables.filter(t => t.status === 'occupied'); }
	function tableName(id: number) { return tables.find(t => t.id === id)?.name ?? `T${id}`; }
	function estimatedFree(tableId: number) { return predictions.find(p => p.table_id === tableId)?.estimated_free; }
</script>

<div class="space-y-4">
	{#if notifications.length > 0}
		<div class="alert alert-success">
			<div class="font-semibold">Ready! ({notifications.length})</div>
			{#each notifications.slice(0, 3) as n}
				<div class="text-sm">×{n.quantity} {n.dish_name} — {n.table_name} <a href="/waiter/orders/{n.order_id}" class="link">View</a></div>
			{/each}
		</div>
	{/if}

	<div class="flex flex-wrap items-center justify-between gap-2">
		<h2 class="text-2xl font-bold">Floor Plan</h2>
		<div class="flex items-center gap-4">
			<div class="flex items-center gap-2 text-sm">
				<span class="flex items-center gap-1"><span class="badge badge-success badge-xs"></span> {freeTables().length}</span>
				<span class="flex items-center gap-1"><span class="badge badge-error badge-xs"></span> {occupiedTables().length}</span>
			</div>
			<div class="join">
				<button class="join-item btn btn-sm" class:btn-active={viewMode === 'canvas'} onclick={() => (viewMode = 'canvas')}>Canvas</button>
				<button class="join-item btn btn-sm" class:btn-active={viewMode === 'list'} onclick={() => (viewMode = 'list')}>List</button>
			</div>
		</div>
	</div>

	{#if loading}
		<div class="flex justify-center py-8"><span class="loading loading-spinner loading-lg"></span></div>
	{:else if viewMode === 'canvas'}
		<FloorPlanCanvas tables={tables} readonly={true} onTableClick={handleTableClick} />
	{:else}
		<div class="overflow-x-auto">
			<table class="table table-zebra">
				<thead><tr><th>Table</th><th>Capacity</th><th>Status</th><th>Group</th><th>Free</th></tr></thead>
				<tbody>
					{#each tables as t (t.id)}
						<tr class="cursor-pointer" onclick={() => handleTableClick(t)}>
							<td class="font-medium">{t.name}</td>
							<td>{t.capacity}</td>
							<td><span class="badge" class:badge-success={t.status === 'free'} class:badge-error={t.status === 'occupied'}>{t.status}</span></td>
							<td>{t.group_name ?? '—'}</td>
							<td class="text-sm text-base-content/50">{estimatedFree(t.id) ? new Date(estimatedFree(t.id)!).toLocaleTimeString() : ''}</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	{/if}
</div>

{#if seating && selectedTable}
	<div class="modal modal-open">
		<div class="modal-box">
			<h3 class="font-bold text-lg">Seat {selectedTable.name} &amp; Order</h3>
			{#if seatingError}<div class="alert alert-error mt-2">{seatingError}</div>{/if}
			<div class="py-4 space-y-3">
				<p class="text-sm text-base-content/60">{selectedTable.capacity} seats — 6 people</p>
				<label class="form-control">
					<div class="label"><span class="label-text">Party size</span></div>
					<input type="number" bind:value={seatPartySize} class="input input-bordered" min="1" max={selectedTable.capacity} />
				</label>
				<label class="form-control">
					<div class="label"><span class="label-text">Name (optional)</span></div>
					<input type="text" bind:value={seatName} placeholder="e.g. Birthday party" class="input input-bordered" />
				</label>
			</div>
			<div class="modal-action gap-2">
				<button class="btn flex-1" onclick={() => { seating = false; selectedTable = null; }}>Cancel</button>
				<button class="btn btn-primary flex-1" onclick={seatAndOrder}>Seat &amp; Order</button>
			</div>
		</div>
	</div>
{/if}

{#if selectedGroup}
	<div class="modal modal-open">
		<div class="modal-box">
			<h3 class="font-bold text-lg">{selectedGroup.name || 'Group'}</h3>
			<div class="py-4 space-y-3">
				<p class="text-sm">Party of {selectedGroup.party_size} | Tables: {selectedGroup.table_ids?.map(id => tableName(id)).join(', ') || '—'}</p>
				{#if estimatedFree(selectedGroup.table_ids?.[0] ?? 0)}
					<p class="text-sm text-base-content/50">Est. free: {new Date(estimatedFree(selectedGroup.table_ids![0])!).toLocaleTimeString()}</p>
				{/if}
				<div class="flex gap-2">
					<button class="btn btn-sm" onclick={() => { merging = true; mergeTargetIds = []; }}>Merge</button>
					<button class="btn btn-sm btn-error" onclick={() => closeGroup(selectedGroup!.id)}>Close</button>
				</div>
				{#if merging}
					<div class="mt-2"><p class="text-sm font-semibold mb-1">Add tables:</p>
						<div class="flex flex-wrap gap-2">
							{#each freeTables() as t}
								<button class="btn btn-outline btn-xs" class:btn-active={mergeTargetIds.includes(t.id)}
									onclick={() => { mergeTargetIds = mergeTargetIds.includes(t.id) ? mergeTargetIds.filter(id => id !== t.id) : [...mergeTargetIds, t.id]; }}>{t.name}</button>
							{/each}
						</div>
						<button class="btn btn-primary btn-sm mt-2" onclick={mergeTables} disabled={mergeTargetIds.length === 0}>Merge</button>
					</div>
				{/if}
				{#if selectedGroup.table_ids && selectedGroup.table_ids.length > 1}
					<div class="mt-2"><p class="text-sm font-semibold mb-1">Split:</p>
						{#each selectedGroup.table_ids as tid}
							<button class="btn btn-ghost btn-xs" onclick={() => splitTable(tid)}>Remove {tableName(tid)}</button>
						{/each}
					</div>
				{/if}
			</div>
			<div class="modal-action"><button class="btn" onclick={() => { selectedGroup = null; merging = false; }}>Close</button></div>
		</div>
	</div>
{/if}
