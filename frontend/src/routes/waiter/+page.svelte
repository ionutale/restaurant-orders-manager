<script lang="ts">
	import { auth } from '$lib/stores/auth.svelte';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { API_BASE } from '$lib/config';
	import FloorPlanCanvas from '$lib/components/FloorPlanCanvas.svelte';

	if (auth.role !== 'waiter') goto('/login');

	type Table = {
		id: number;
		name: string;
		capacity: number;
		x: number;
		y: number;
		label: string | null;
		status: 'free' | 'occupied';
		group_name?: string;
		party_size?: number;
	};

	let tables = $state<Table[]>([]);
	let loading = $state(true);
	let viewMode = $state<'canvas' | 'list'>('canvas');
	let selectedTable = $state<Table | null>(null);
	let seating = $state(false);
	let seatName = $state('');
	let seatPartySize = $state(2);
	let seatingError = $state('');

	const token = () => localStorage.getItem('token') ?? '';

	async function load() {
		const res = await fetch(`${API_BASE}/floor-plan`, { headers: { Authorization: `Bearer ${token()}` } });
		if (res.ok) tables = await res.json();
		loading = false;
	}

	onMount(load);

	async function seatGuests() {
		if (!selectedTable) return;
		seatingError = '';
		const res = await fetch(`${API_BASE}/table-groups`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token()}` },
			body: JSON.stringify({
				table_ids: [selectedTable.id],
				party_size: seatPartySize,
				name: seatName || undefined,
			}),
		});
		if (!res.ok) { seatingError = 'Failed to seat'; return; }
		selectedTable = null;
		seating = false;
		seatName = '';
		seatPartySize = 2;
		await load();
	}

	function freeTables() { return tables.filter(t => t.status === 'free'); }
	function occupiedTables() { return tables.filter(t => t.status === 'occupied'); }

	function handleTableClick(t: Table) {
		selectedTable = t;
		if (t.status === 'free') {
			seating = true;
			seatName = '';
			seatPartySize = Math.min(2, t.capacity);
		}
	}
</script>

<div class="space-y-4">
	<div class="flex flex-wrap items-center justify-between gap-2">
		<h2 class="text-2xl font-bold">Floor Plan</h2>
		<div class="flex items-center gap-4">
			<div class="flex items-center gap-2 text-sm">
				<span class="flex items-center gap-1"><span class="badge badge-success badge-xs"></span> {freeTables().length} free</span>
				<span class="flex items-center gap-1"><span class="badge badge-error badge-xs"></span> {occupiedTables().length} occupied</span>
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
				<thead>
					<tr><th>Table</th><th>Capacity</th><th>Label</th><th>Status</th><th>Group</th></tr>
				</thead>
				<tbody>
					{#each tables as t (t.id)}
						<tr class="cursor-pointer" onclick={() => handleTableClick(t)}>
							<td class="font-medium">{t.name}</td>
							<td>{t.capacity}</td>
							<td>{t.label ?? '—'}</td>
							<td><span class="badge" class:badge-success={t.status === 'free'} class:badge-error={t.status === 'occupied'}>{t.status}</span></td>
							<td>{t.group_name ?? '—'}</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	{/if}
</div>

{#if selectedTable && seating}
	<div class="modal modal-open">
		<div class="modal-box">
			<h3 class="font-bold text-lg">Seat {selectedTable.name}</h3>
			{#if seatingError}<div class="alert alert-error mt-2">{seatingError}</div>{/if}
			<div class="py-4 space-y-3">
				<p class="text-sm text-base-content/60">{selectedTable.capacity} seats available</p>
				<label class="form-control">
					<div class="label"><span class="label-text">Party size</span></div>
					<input type="number" bind:value={seatPartySize} class="input input-bordered" min="1" max={selectedTable.capacity} />
				</label>
				<label class="form-control">
					<div class="label"><span class="label-text">Name (optional)</span></div>
					<input type="text" bind:value={seatName} placeholder="e.g. Birthday party" class="input input-bordered" />
				</label>
			</div>
			<div class="modal-action">
				<button class="btn" onclick={() => { selectedTable = null; seating = false; }}>Cancel</button>
				<button class="btn btn-primary" onclick={seatGuests}>Seat Guests</button>
			</div>
		</div>
	</div>
{:else if selectedTable}
	<div class="modal modal-open">
		<div class="modal-box">
			<h3 class="font-bold text-lg">{selectedTable.name}</h3>
			<div class="py-4 space-y-3">
				<div class="flex items-center gap-2">
					<span class="badge" class:badge-success={selectedTable.status === 'free'} class:badge-error={selectedTable.status === 'occupied'}>{selectedTable.status}</span>
					<span class="text-base-content/60">{selectedTable.capacity} seats</span>
				</div>
				{#if selectedTable.status === 'occupied'}
					{#if selectedTable.group_name}<p class="font-semibold">{selectedTable.group_name}</p>{/if}
					{#if selectedTable.party_size}<p class="text-sm text-base-content/60">Party of {selectedTable.party_size}</p>{/if}
				{/if}
				{#if selectedTable.label}<p class="text-sm text-base-content/50">{selectedTable.label}</p>{/if}
			</div>
			<div class="modal-action">
				<button class="btn" onclick={() => (selectedTable = null)}>Close</button>
			</div>
		</div>
	</div>
{/if}
