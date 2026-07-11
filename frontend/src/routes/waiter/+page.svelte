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
		group_id?: number;
		group_name?: string;
		party_size?: number;
	};

	type Group = {
		id: number;
		name: string | null;
		party_size: number;
		status: string;
		table_ids: number[];
	};

	let tables = $state<Table[]>([]);
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
		const res = await fetch(`${API_BASE}/floor-plan`, { headers: { Authorization: `Bearer ${token()}` } });
		if (res.ok) tables = await res.json();
		loading = false;
	}

	onMount(load);

	async function seatGuests() {
		seatingError = '';
		const res = await fetch(`${API_BASE}/table-groups`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token()}` },
			body: JSON.stringify({ table_ids: seatTargetIds, party_size: seatPartySize, name: seatName || undefined }),
		});
		if (!res.ok) { seatingError = 'Failed to seat'; return; }
		seating = false; selectedTable = null; seatName = ''; seatPartySize = 2; seatTargetIds = [];
		await load();
	}

	async function loadGroup(id: number) {
		const res = await fetch(`${API_BASE}/table-groups/${id}`, { headers: { Authorization: `Bearer ${token()}` } });
		if (res.ok) selectedGroup = await res.json();
	}

	async function closeGroup(id: number) {
		await fetch(`${API_BASE}/table-groups/${id}/close`, { method: 'POST', headers: { Authorization: `Bearer ${token()}` } });
		selectedTable = null; selectedGroup = null;
		await load();
	}

	async function mergeTables() {
		if (!selectedGroup) return;
		await fetch(`${API_BASE}/table-groups/${selectedGroup.id}/tables`, {
			method: 'PATCH',
			headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token()}` },
			body: JSON.stringify({ add_table_ids: mergeTargetIds, remove_table_ids: [] }),
		});
		merging = false; mergeTargetIds = [];
		await loadGroup(selectedGroup.id);
		await load();
	}

	async function splitTable(tableId: number) {
		if (!selectedGroup) return;
		await fetch(`${API_BASE}/table-groups/${selectedGroup.id}/tables`, {
			method: 'PATCH',
			headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token()}` },
			body: JSON.stringify({ add_table_ids: [], remove_table_ids: [tableId] }),
		});
		await loadGroup(selectedGroup.id);
		await load();
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

{#if seating}
	<div class="modal modal-open">
		<div class="modal-box">
			<h3 class="font-bold text-lg">Seat {selectedTable?.name}</h3>
			{#if seatingError}<div class="alert alert-error mt-2">{seatingError}</div>{/if}
			<div class="py-4 space-y-3">
				<p class="text-sm text-base-content/60">{selectedTable?.capacity} seats</p>
				<label class="form-control"><div class="label"><span class="label-text">Party size</span></div><input type="number" bind:value={seatPartySize} class="input input-bordered" min="1" max={selectedTable?.capacity} /></label>
				<label class="form-control"><div class="label"><span class="label-text">Name (optional)</span></div><input type="text" bind:value={seatName} placeholder="e.g. Birthday party" class="input input-bordered" /></label>
			</div>
			<div class="modal-action">
				<button class="btn" onclick={() => { seating = false; selectedTable = null; }}>Cancel</button>
				<button class="btn btn-primary" onclick={seatGuests}>Seat Guests</button>
			</div>
		</div>
	</div>
{/if}

{#if selectedGroup}
	<div class="modal modal-open">
		<div class="modal-box">
			<h3 class="font-bold text-lg">{selectedGroup.name || 'Table Group'}</h3>
			<div class="py-4 space-y-3">
				<p class="text-sm">Party of {selectedGroup.party_size} | Tables: {selectedGroup.table_ids?.map(id => tableName(id)).join(', ') || '—'}</p>
				<div class="flex gap-2">
					<button class="btn btn-sm" onclick={() => { merging = true; mergeTargetIds = []; }}>Merge Tables</button>
					<button class="btn btn-sm btn-error" onclick={() => closeGroup(selectedGroup!.id)}>Close Group</button>
				</div>
				{#if selectedGroup.table_ids && selectedGroup.table_ids.length > 1}
					<div class="mt-2">
						<p class="text-sm font-semibold mb-1">Split table from group:</p>
						<div class="flex flex-wrap gap-2">
							{#each selectedGroup.table_ids as tid}
								<button class="btn btn-ghost btn-xs" onclick={() => splitTable(tid)}>Remove {tableName(tid)}</button>
							{/each}
						</div>
					</div>
				{/if}
				{#if merging}
					<div class="mt-2">
						<p class="text-sm font-semibold mb-1">Add free tables:</p>
						<div class="flex flex-wrap gap-2">
							{#each freeTables() as t}
								<button
									class="btn btn-outline btn-xs"
									class:btn-active={mergeTargetIds.includes(t.id)}
									onclick={() => {
										if (mergeTargetIds.includes(t.id)) {
											mergeTargetIds = mergeTargetIds.filter(id => id !== t.id);
										} else {
											mergeTargetIds = [...mergeTargetIds, t.id];
										}
									}}
								>{t.name}</button>
							{/each}
						</div>
						<button class="btn btn-primary btn-sm mt-2" onclick={mergeTables} disabled={mergeTargetIds.length === 0}>Merge</button>
					</div>
				{/if}
			</div>
			<div class="modal-action">
				<button class="btn" onclick={() => { selectedGroup = null; merging = false; }}>Close</button>
			</div>
		</div>
	</div>
{/if}
