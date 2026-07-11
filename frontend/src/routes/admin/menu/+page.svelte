<script lang="ts">
	import { auth } from '$lib/stores/auth.svelte';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';

	if (auth.role !== 'manager') goto('/login');

	import { API_BASE } from '$lib/config';

	type Category = {
		id: number;
		name: string;
		display_order: number;
		icon: string;
	};

	let categories = $state<Category[]>([]);
	let loading = $state(true);
	let editingId = $state<number | null>(null);
	let editName = $state('');
	let editIcon = $state('');
	let deleteId = $state<number | null>(null);
	let adding = $state(false);
	let newName = $state('');
	let newIcon = $state('');

	async function load() {
		const token = localStorage.getItem('token');
		const res = await fetch(`${API_BASE}/categories`, { headers: { Authorization: `Bearer ${token}` } });
		if (res.ok) categories = await res.json();
		loading = false;
	}

	onMount(load);

	async function create() {
		if (!newName) return;
		const token = localStorage.getItem('token');
		const res = await fetch(`${API_BASE}/categories`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` },
			body: JSON.stringify({ name: newName, icon: newIcon || undefined }),
		});
		if (res.ok) {
			newName = '';
			newIcon = '';
			adding = false;
			await load();
		}
	}

	async function save(id: number) {
		const token = localStorage.getItem('token');
		await fetch(`${API_BASE}/categories/${id}`, {
			method: 'PATCH',
			headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` },
			body: JSON.stringify({ name: editName, icon: editIcon || undefined }),
		});
		editingId = null;
		await load();
	}

	async function remove(id: number) {
		const token = localStorage.getItem('token');
		await fetch(`${API_BASE}/categories/${id}`, {
			method: 'DELETE',
			headers: { Authorization: `Bearer ${token}` },
		});
		deleteId = null;
		await load();
	}

	async function move(id: number, delta: number) {
		const token = localStorage.getItem('token');
		const res = await fetch(`${API_BASE}/categories/reorder`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` },
			body: JSON.stringify({ id, delta }),
		});
		if (res.ok) await load();
	}

	function startEdit(c: Category) {
		editingId = c.id;
		editName = c.name;
		editIcon = c.icon;
	}
</script>

<div class="space-y-4">
	<div class="flex items-center justify-between">
		<h2 class="text-2xl font-bold">Menu — Categories</h2>
		<button class="btn btn-primary btn-sm" onclick={() => (adding = true)}>Add Category</button>
	</div>

	{#if loading}
		<div class="flex justify-center py-8"><span class="loading loading-spinner loading-lg"></span></div>
	{:else if categories.length === 0}
		<div class="flex h-32 items-center justify-center rounded-box border-2 border-dashed text-base-content/40">No categories yet</div>
	{:else}
		<div class="overflow-x-auto">
			<table class="table table-zebra">
				<thead>
					<tr><th>Order</th><th>Name</th><th>Icon</th><th>Actions</th></tr>
				</thead>
				<tbody>
					{#each categories as c (c.id)}
						<tr>
							{#if editingId === c.id}
								<td class="text-center text-base-content/50">{c.display_order}</td>
								<td><input type="text" bind:value={editName} class="input input-bordered input-sm w-full" /></td>
								<td><input type="text" bind:value={editIcon} class="input input-bordered input-sm w-24" placeholder="🍕" /></td>
								<td class="flex gap-1">
									<button class="btn btn-ghost btn-xs text-success" onclick={() => save(c.id)}>Save</button>
									<button class="btn btn-ghost btn-xs" onclick={() => (editingId = null)}>Cancel</button>
								</td>
							{:else if deleteId === c.id}
								<td colspan="2" class="text-warning">Delete "{c.name}"?</td>
								<td colspan="2" class="flex gap-1">
									<button class="btn btn-ghost btn-xs text-error" onclick={() => remove(c.id)}>Confirm</button>
									<button class="btn btn-ghost btn-xs" onclick={() => (deleteId = null)}>Cancel</button>
								</td>
							{:else}
								<td class="text-center text-base-content/50">{c.display_order}</td>
								<td class="font-medium">{c.name}</td>
								<td>{c.icon || '—'}</td>
								<td class="flex gap-1">
									<button class="btn btn-ghost btn-xs" onclick={() => move(c.id, -1)} disabled={c.display_order <= 1}>↑</button>
									<button class="btn btn-ghost btn-xs" onclick={() => move(c.id, 1)} disabled={c.display_order >= categories.length}>↓</button>
									<button class="btn btn-ghost btn-xs" onclick={() => startEdit(c)}>Edit</button>
									<button class="btn btn-ghost btn-xs text-error" onclick={() => (deleteId = c.id)}>Delete</button>
								</td>
							{/if}
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	{/if}
</div>

{#if adding}
	<div class="modal modal-open">
		<div class="modal-box">
			<h3 class="font-bold text-lg">Add Category</h3>
			<div class="py-4 space-y-3">
				<label class="form-control">
					<div class="label"><span class="label-text">Name</span></div>
					<input type="text" bind:value={newName} placeholder="Appetizers" class="input input-bordered" />
				</label>
				<label class="form-control">
					<div class="label"><span class="label-text">Icon (optional)</span></div>
					<input type="text" bind:value={newIcon} placeholder="🍕" class="input input-bordered" />
				</label>
			</div>
			<div class="modal-action">
				<button class="btn" onclick={() => (adding = false)}>Cancel</button>
				<button class="btn btn-primary" onclick={create}>Add</button>
			</div>
		</div>
		<div class="modal-backdrop" onclick={() => (adding = false)}></div>
	</div>
{/if}
