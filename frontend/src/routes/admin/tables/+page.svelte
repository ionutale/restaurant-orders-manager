<script lang="ts">
	import { auth } from '$lib/stores/auth.svelte';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import FloorPlanCanvas from '$lib/components/FloorPlanCanvas.svelte';

	if (auth.role !== 'manager') goto('/login');

	import { API_BASE } from '$lib/config';

	type Table = {
		id: number;
		name: string;
		capacity: number;
		x: number;
		y: number;
		label: string | null;
	};

	let tables = $state<Table[]>([]);
	let loading = $state(true);
	let error = $state('');
	let viewMode = $state<'canvas' | 'list'>('canvas');
	let editingId = $state<number | null>(null);
	let editName = $state('');
	let editCapacity = $state(4);
	let editLabel = $state('');
	let deleteId = $state<number | null>(null);
	let sortKey = $state<keyof Table>('name');
	let sortDir = $state<'asc' | 'desc'>('asc');

	let newName = $state('');
	let newCapacity = $state(4);
	let newLabel = $state('');

	let showAddDialog = $state(false);
	let addName = $state('');
	let addCapacity = $state(4);
	let addLabel = $state('');

	function startEdit(t: Table) {
		editingId = t.id;
		editName = t.name;
		editCapacity = t.capacity;
		editLabel = t.label ?? '';
	}

	async function load() {
		try {
			const token = localStorage.getItem('token');
			const res = await fetch(`${API_BASE}/tables`, { headers: { Authorization: `Bearer ${token}` } });
			if (!res.ok) throw new Error('Failed to load');
			tables = await res.json();
		} catch (e: any) {
			error = e.message;
		} finally {
			loading = false;
		}
	}

	onMount(load);

	function sorted() {
		return [...tables].sort((a, b) => {
			const av = a[sortKey] ?? '';
			const bv = b[sortKey] ?? '';
			const cmp = String(av).localeCompare(String(bv), undefined, { numeric: true });
			return sortDir === 'asc' ? cmp : -cmp;
		});
	}

	function toggleSort(key: keyof Table) {
		if (sortKey === key) {
			sortDir = sortDir === 'asc' ? 'desc' : 'asc';
		} else {
			sortKey = key;
			sortDir = 'asc';
		}
	}

	async function create() {
		if (!newName) return;
		const token = localStorage.getItem('token');
		await fetch(`${API_BASE}/tables`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` },
			body: JSON.stringify({ name: newName, capacity: newCapacity, label: newLabel || null }),
		});
		newName = '';
		newCapacity = 4;
		newLabel = '';
		await load();
	}

	async function save(id: number, name: string, capacity: number, label: string) {
		const token = localStorage.getItem('token');
		await fetch(`${API_BASE}/tables/${id}`, {
			method: 'PATCH',
			headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` },
			body: JSON.stringify({ name, capacity, label: label || null }),
		});
		editingId = null;
		await load();
	}

	async function remove(id: number) {
		const token = localStorage.getItem('token');
		await fetch(`${API_BASE}/tables/${id}`, {
			method: 'DELETE',
			headers: { Authorization: `Bearer ${token}` },
		});
		deleteId = null;
		await load();
	}

	async function updatePosition(id: number, data: Partial<Table>) {
		const token = localStorage.getItem('token');
		await fetch(`${API_BASE}/tables/${id}`, {
			method: 'PATCH',
			headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` },
			body: JSON.stringify(data),
		});
		await load();
	}

	function openAddDialog() {
		addName = '';
		addCapacity = 4;
		addLabel = '';
		showAddDialog = true;
	}

	async function addFromDialog() {
		if (!addName) return;
		const token = localStorage.getItem('token');
		const res = await fetch(`${API_BASE}/tables`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` },
			body: JSON.stringify({ name: addName, capacity: addCapacity, label: addLabel || null }),
		});
		if (res.ok) {
			showAddDialog = false;
			await load();
		}
	}
</script>

<div class="space-y-4">
	<div class="flex flex-wrap items-center justify-between gap-2">
		<h2 class="text-2xl font-bold">Floor Plan</h2>
		<div class="join">
			<button
				class="join-item btn btn-sm"
				class:btn-active={viewMode === 'canvas'}
				onclick={() => (viewMode = 'canvas')}
			>Canvas</button>
			<button
				class="join-item btn btn-sm"
				class:btn-active={viewMode === 'list'}
				onclick={() => (viewMode = 'list')}
			>List</button>
		</div>
	</div>

	{#if error}
		<div class="alert alert-error">{error}</div>
	{/if}

	{#if loading}
		<div class="flex justify-center py-8"><span class="loading loading-spinner loading-lg"></span></div>
	{:else if viewMode === 'canvas'}
		<FloorPlanCanvas {tables} onUpdate={updatePosition} onAdd={openAddDialog} onDelete={(id) => (deleteId = id)} />
	{:else}
		<div class="card bg-base-100 shadow-xl">
			<div class="card-body p-0">
				<div class="overflow-x-auto">
					<table class="table table-zebra">
						<thead>
							<tr>
								<th class="cursor-pointer" onclick={() => toggleSort('name')}>
									Name {sortKey === 'name' ? (sortDir === 'asc' ? '↑' : '↓') : ''}
								</th>
								<th class="cursor-pointer" onclick={() => toggleSort('capacity')}>
									Capacity {sortKey === 'capacity' ? (sortDir === 'asc' ? '↑' : '↓') : ''}
								</th>
								<th>Label</th>
								<th class="w-32">Actions</th>
							</tr>
						</thead>
						<tbody>
							{#if sorted().length === 0}
								<tr><td colspan="4" class="text-center text-base-content/50">No tables yet</td></tr>
							{:else}
								{#each sorted() as table (table.id)}
									<tr>
										{#if editingId === table.id}
											<td><input type="text" bind:value={editName} class="input input-bordered input-sm w-full" /></td>
											<td><input type="number" bind:value={editCapacity} class="input input-bordered input-sm w-20" min="1" /></td>
											<td><input type="text" bind:value={editLabel} class="input input-bordered input-sm w-full" placeholder="Optional" /></td>
											<td class="flex gap-1">
												<button class="btn btn-ghost btn-xs text-success" onclick={() => save(table.id, editName, editCapacity, editLabel)}>Save</button>
												<button class="btn btn-ghost btn-xs" onclick={() => (editingId = null)}>Cancel</button>
											</td>
										{:else if deleteId === table.id}
											<td colspan="3" class="text-warning">Delete {table.name}?</td>
											<td class="flex gap-1">
												<button class="btn btn-ghost btn-xs text-error" onclick={() => { remove(table.id); }}>Confirm</button>
												<button class="btn btn-ghost btn-xs" onclick={() => (deleteId = null)}>Cancel</button>
											</td>
										{:else}
											<td class="font-medium">{table.name}</td>
											<td>{table.capacity}</td>
											<td>{table.label ?? '—'}</td>
											<td class="flex gap-1">
												<button class="btn btn-ghost btn-xs" onclick={() => startEdit(table)}>Edit</button>
												<button class="btn btn-ghost btn-xs text-error" onclick={() => (deleteId = table.id)}>Delete</button>
											</td>
										{/if}
									</tr>
								{/each}
							{/if}
						</tbody>
					</table>
				</div>
			</div>
		</div>

		<div class="card bg-base-100 shadow-xl">
			<div class="card-body">
				<h3 class="card-title">Add Table</h3>
				<div class="flex flex-wrap items-end gap-3">
					<label class="form-control w-32">
						<div class="label"><span class="label-text">Name</span></div>
						<input type="text" bind:value={newName} placeholder="T6" class="input input-bordered input-sm" />
					</label>
					<label class="form-control w-24">
						<div class="label"><span class="label-text">Capacity</span></div>
						<input type="number" bind:value={newCapacity} class="input input-bordered input-sm" min="1" />
					</label>
					<label class="form-control w-40">
						<div class="label"><span class="label-text">Label (optional)</span></div>
						<input type="text" bind:value={newLabel} placeholder="Near window" class="input input-bordered input-sm" />
					</label>
					<button class="btn btn-primary btn-sm" onclick={create}>Add</button>
				</div>
			</div>
		</div>
	{/if}

	{#if viewMode === 'canvas'}
		<div class="card bg-base-100 shadow-xl">
			<div class="card-body">
				<h3 class="card-title">All Tables</h3>
				<div class="overflow-x-auto">
					<table class="table table-sm">
						<thead>
							<tr><th>Name</th><th>Capacity</th><th>Label</th><th>Actions</th></tr>
						</thead>
						<tbody>
							{#each tables as t (t.id)}
								<tr>
									<td class="font-medium">{t.name}</td>
									<td>{t.capacity}</td>
									<td>{t.label ?? '—'}</td>
									<td>
										<button class="btn btn-ghost btn-xs text-error" onclick={() => (deleteId = t.id)}>Delete</button>
									</td>
								</tr>
							{/each}
						</tbody>
					</table>
				</div>
			</div>
		</div>
	{/if}

	{#if showAddDialog}
		<div class="modal modal-open">
			<div class="modal-box">
				<h3 class="font-bold text-lg">Add Table</h3>
				<div class="py-4 space-y-3">
					<label class="form-control">
						<div class="label"><span class="label-text">Name</span></div>
						<input type="text" bind:value={addName} placeholder="T6" class="input input-bordered" />
					</label>
					<label class="form-control">
						<div class="label"><span class="label-text">Capacity</span></div>
						<input type="number" bind:value={addCapacity} class="input input-bordered" min="1" />
					</label>
					<label class="form-control">
						<div class="label"><span class="label-text">Label (optional)</span></div>
						<input type="text" bind:value={addLabel} placeholder="Near window" class="input input-bordered" />
					</label>
				</div>
				<div class="modal-action">
					<button class="btn" onclick={() => (showAddDialog = false)}>Cancel</button>
					<button class="btn btn-primary" onclick={addFromDialog}>Add</button>
				</div>
			</div>
			<div class="modal-backdrop" onclick={() => (showAddDialog = false)}></div>
		</div>
	{/if}

	{#if deleteId !== null && viewMode === 'canvas'}
		<div class="modal modal-open">
			<div class="modal-box">
				<h3 class="font-bold text-lg">Delete table?</h3>
				<p>Are you sure you want to delete this table?</p>
				<div class="modal-action">
					<button class="btn" onclick={() => (deleteId = null)}>Cancel</button>
					<button class="btn btn-error" onclick={() => { const id = deleteId; deleteId = null; remove(id!); }}>Delete</button>
				</div>
			</div>
			<div class="modal-backdrop" onclick={() => (deleteId = null)}></div>
		</div>
	{/if}
</div>
