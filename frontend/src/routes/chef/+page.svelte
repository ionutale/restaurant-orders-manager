<script lang="ts">
	import { auth } from '$lib/stores/auth.svelte';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { API_BASE } from '$lib/config';

	if (auth.role !== 'chef') goto('/login');

	type Item = { id: number; dish_name: string; quantity: number; notes: string; ready: boolean };
	type Course = { id: number; name: string; items: Item[] };
	type Order = { id: number; courses: Course[] };

	let orders = $state<Order[]>([]);
	let loading = $state(true);

	const token = () => localStorage.getItem('token') ?? '';

	async function load() {
		const res = await fetch(`${API_BASE}/kds/orders`, { headers: { Authorization: `Bearer ${token()}` } });
		if (res.ok) orders = await res.json();
		loading = false;
	}

	onMount(load);
	setInterval(load, 10000);

	async function markReady(itemId: number) {
		await fetch(`${API_BASE}/kds/order-items/${itemId}/ready`, { method: 'PATCH', headers: { Authorization: `Bearer ${token()}` } });
		await load();
	}

	function readyCount(items: Item[]) { return items.filter(i => i.ready).length; }
</script>

<div class="space-y-4">
	<div class="flex items-center justify-between">
		<h2 class="text-2xl font-bold">KDS Dashboard</h2>
		<button class="btn btn-sm" onclick={load}>Refresh</button>
	</div>

	{#if loading}
		<div class="flex justify-center py-8"><span class="loading loading-spinner loading-lg"></span></div>
	{:else if orders.length === 0}
		<div class="flex h-48 items-center justify-center rounded-box border-2 border-dashed text-base-content/40">No active orders</div>
	{:else}
		<div class="grid grid-cols-1 gap-4 md:grid-cols-2 xl:grid-cols-3">
			{#each orders as o (o.id)}
				<div class="card bg-base-100 shadow-xl">
					<div class="card-body">
						<div class="flex items-center justify-between">
							<h3 class="card-title">Order #{o.id}</h3>
						</div>
						{#each o.courses as c (c.id)}
							<div class="mt-2">
								<h4 class="font-semibold text-sm badge badge-warning">{c.name}</h4>
								<div class="mt-2 space-y-2">
									{#each c.items as item (item.id)}
										<div class="flex items-center justify-between rounded-box bg-base-200 p-2" class:opacity-50={item.ready}>
											<div class="flex-1">
												<div class="font-medium text-sm">×{item.quantity} {item.dish_name}</div>
												{#if item.notes}<div class="text-xs text-warning">"{item.notes}"</div>{/if}
											</div>
											{#if !item.ready}
												<button class="btn btn-primary btn-xs" onclick={() => markReady(item.id)}>Ready</button>
											{:else}
												<span class="badge badge-success badge-xs">✓</span>
											{/if}
										</div>
									{/each}
								</div>
							</div>
						{/each}
					</div>
				</div>
			{/each}
		</div>
	{/if}
</div>
