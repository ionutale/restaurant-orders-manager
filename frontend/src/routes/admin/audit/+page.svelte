<script lang="ts">
	import { auth } from '$lib/stores/auth.svelte';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { API_BASE } from '$lib/config';

	if (auth.role !== 'manager') goto('/login');

	type Event = { id: number; actor_name: string; action: string; entity: string; entity_id: number | null; created_at: string };
	let events = $state<Event[]>([]);
	let loading = $state(true);

	onMount(async () => {
		const token = localStorage.getItem('token');
		const r = await fetch(`${API_BASE}/audit-events`, { headers: { Authorization: `Bearer ${token}` } });
		if (r.ok) events = await r.json();
		loading = false;
	});

	function actionBadge(action: string) {
		const m: Record<string, string> = {};
		return 'badge-ghost';
	}
</script>

<div class="space-y-4">
	<h2 class="text-2xl font-bold">Audit Log</h2>

	{#if loading}
		<div class="flex justify-center py-8"><span class="loading loading-spinner loading-lg"></span></div>
	{:else if events.length === 0}
		<div class="flex h-32 items-center justify-center rounded-box border-2 border-dashed text-base-content/40">No events yet</div>
	{:else}
		<div class="overflow-x-auto">
			<table class="table table-zebra table-sm">
				<thead><tr><th>Time</th><th>Actor</th><th>Action</th><th>Entity</th></tr></thead>
				<tbody>
					{#each events as e (e.id)}
						<tr>
							<td class="text-xs text-base-content/50">{new Date(e.created_at).toLocaleString()}</td>
							<td>{e.actor_name}</td>
							<td><span class="badge badge-sm">{e.action}</span></td>
							<td class="text-sm">{e.entity}{#if e.entity_id} #{e.entity_id}{/if}</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	{/if}
</div>
