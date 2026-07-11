<script lang="ts">
	import { auth } from '$lib/stores/auth.svelte';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { API_BASE } from '$lib/config';

	if (auth.role !== 'manager') goto('/login');

	type User = { id: number; name: string; email: string; role: string };
	let users = $state<User[]>([]);
	let loading = $state(true);
	let showAdd = $state(false);
	let newName = $state('');
	let newEmail = $state('');
	let newPass = $state('');
	let newRole = $state('waiter');
	let editId = $state<number | null>(null);
	let editName = $state('');
	let editEmail = $state('');
	let editRole = $state('');
	let deleteId = $state<number | null>(null);

	const token = () => localStorage.getItem('token') ?? '';

	async function load() {
		const r = await fetch(`${API_BASE}/users`, { headers: { Authorization: `Bearer ${token()}` } });
		if (r.ok) users = await r.json();
		loading = false;
	}

	onMount(load);

	async function create() {
		if (!newName || !newEmail || !newPass) return;
		const r = await fetch(`${API_BASE}/users`, {
			method: 'POST', headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token()}` },
			body: JSON.stringify({ name: newName, email: newEmail, password: newPass, role: newRole }),
		});
		if (r.ok) { showAdd = false; newName = ''; newEmail = ''; newPass = ''; newRole = 'waiter'; await load(); }
	}

	async function save(id: number) {
		await fetch(`${API_BASE}/users/${id}`, {
			method: 'PATCH', headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token()}` },
			body: JSON.stringify({ name: editName, email: editEmail, role: editRole }),
		});
		editId = null; await load();
	}

	async function remove(id: number) {
		await fetch(`${API_BASE}/users/${id}`, { method: 'DELETE', headers: { Authorization: `Bearer ${token()}` } });
		deleteId = null; await load();
	}

	function startEdit(u: User) { editId = u.id; editName = u.name; editEmail = u.email; editRole = u.role; }
</script>

<div class="space-y-4">
	<div class="flex items-center justify-between">
		<h2 class="text-2xl font-bold">Users</h2>
		<button class="btn btn-primary btn-sm" onclick={() => (showAdd = true)}>Add User</button>
	</div>

	{#if loading}
		<div class="flex justify-center py-8"><span class="loading loading-spinner loading-lg"></span></div>
	{:else}
		<div class="overflow-x-auto">
			<table class="table table-zebra">
				<thead><tr><th>Name</th><th>Email</th><th>Role</th><th>Actions</th></tr></thead>
				<tbody>
					{#each users as u (u.id)}
						<tr>
							{#if editId === u.id}
								<td><input type="text" bind:value={editName} class="input input-bordered input-sm w-full" /></td>
								<td><input type="email" bind:value={editEmail} class="input input-bordered input-sm w-full" /></td>
								<td><select bind:value={editRole} class="select select-bordered select-sm">
									<option value="waiter">Waiter</option>
									<option value="chef">Chef</option>
									<option value="manager">Manager</option>
								</select></td>
								<td class="flex gap-1">
									<button class="btn btn-ghost btn-xs text-success" onclick={() => save(u.id)}>Save</button>
									<button class="btn btn-ghost btn-xs" onclick={() => (editId = null)}>Cancel</button>
								</td>
							{:else if deleteId === u.id}
								<td colspan="2" class="text-warning">Delete {u.name}?</td>
								<td colspan="2" class="flex gap-1">
									<button class="btn btn-ghost btn-xs text-error" onclick={() => remove(u.id)}>Confirm</button>
									<button class="btn btn-ghost btn-xs" onclick={() => (deleteId = null)}>Cancel</button>
								</td>
							{:else}
								<td class="font-medium">{u.name}</td>
								<td>{u.email}</td>
								<td><span class="badge badge-sm">{u.role}</span></td>
								<td class="flex gap-1">
									<button class="btn btn-ghost btn-xs" onclick={() => startEdit(u)}>Edit</button>
									<button class="btn btn-ghost btn-xs text-error" onclick={() => (deleteId = u.id)}>Delete</button>
								</td>
							{/if}
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	{/if}
</div>

{#if showAdd}
	<div class="modal modal-open">
		<div class="modal-box">
			<h3 class="font-bold text-lg mb-4">Add User</h3>
			<div class="space-y-3">
				<label class="form-control w-full">
					<div class="label"><span class="label-text">Name</span></div>
					<input type="text" bind:value={newName} placeholder="Full name" class="input input-bordered w-full" />
				</label>
				<label class="form-control w-full">
					<div class="label"><span class="label-text">Email</span></div>
					<input type="email" bind:value={newEmail} placeholder="email@example.com" class="input input-bordered w-full" />
				</label>
				<label class="form-control w-full">
					<div class="label"><span class="label-text">Password</span></div>
					<input type="password" bind:value={newPass} placeholder="••••••••" class="input input-bordered w-full" />
				</label>
				<label class="form-control w-full">
					<div class="label"><span class="label-text">Role</span></div>
					<select bind:value={newRole} class="select select-bordered w-full">
						<option value="waiter">Waiter</option>
						<option value="chef">Chef</option>
						<option value="manager">Manager</option>
					</select>
				</label>
			</div>
			<div class="modal-action">
				<button class="btn" onclick={() => (showAdd = false)}>Cancel</button>
				<button class="btn btn-primary" onclick={create}>Add User</button>
			</div>
		</div>
	</div>
{/if}
