<script lang="ts">
	import { auth } from '$lib/stores/auth.svelte';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';

	if (auth.role !== 'manager') goto('/login');

	const API = 'http://localhost:8080/api';

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

	function startEdit(t: Table) {
		editingId = t.id;
		editName = t.name;
		editCapacity = t.capacity;
		editLabel = t.label ?? '';
	}

	async function load() {
		try {
			const token = localStorage.getItem('token');
			const res = await fetch(`${API}/tables`, { headers: { Authorization: `Bearer ${token}` } });
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
		const sorted = [...tables].sort((a, b) => {
			const av = a[sortKey] ?? '';
			const bv = b[sortKey] ?? '';
			const cmp = String(av).localeCompare(String(bv), undefined, { numeric: true });
			return sortDir === 'asc' ? cmp : -cmp;
		});
		return sorted;
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
		const res = await fetch(`${API}/tables`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` },
			body: JSON.stringify({ name: newName, capacity: newCapacity, label: newLabel || null }),
		});
		if (!res.ok) { error = 'Failed to create table'; return; }
		newName = '';
		newCapacity = 4;
		newLabel = '';
		await load();
	}

	async function save(id: number, name: string, capacity: number, label: string) {
		const token = localStorage.getItem('token');
		const res = await fetch(`${API}/tables/${id}`, {
			method: 'PATCH',
			headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` },
			body: JSON.stringify({ name, capacity, label: label || null }),
		});
		if (!res.ok) { error = 'Failed to update'; return; }
		editingId = null;
		await load();
	}

	async function remove(id: number) {
		const token = localStorage.getItem('token');
		const res = await fetch(`${API}/tables/${id}`, {
			method: 'DELETE',
			headers: { Authorization: `Bearer ${token}` },
		});
		if (!res.ok) { error = 'Failed to delete'; return; }
		deleteId = null;
		await load();
	}
</script>

<div class="space-y-4">
	<div class="flex items-center justify-between">
		<h2 class="text-2xl font-bold">Floor Plan — Tables</h2>
	</div>

	{#if error}
		<div class="alert alert-error">{error}</div>
	{/if}

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
						{#if loading}
							<tr><td colspan="4" class="text-center"><span class="loading loading-spinner loading-sm"></span> Loading...</td></tr>
						{:else if sorted().length === 0}
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
											<button class="btn btn-ghost btn-xs text-error" onclick={() => remove(table.id)}>Confirm</button>
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
				<button class="btn btn-primary btn-sm" onclick={create}>Add Table</button>
			</div>
		</div>
	</div>
</div>
