<script lang="ts">
	import { auth } from '$lib/stores/auth.svelte';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { API_BASE } from '$lib/config';

	if (auth.role !== 'manager') goto('/login');

	type Stats = {
		people_today: number;
		revenue_today_cents: number;
		orders_today: number;
		tables_occupied: number;
		top_dishes: { name: string; count: number; total_cents: number }[];
		top_wines: { name: string; count: number; total_cents: number }[];
	};

	let stats = $state<Stats | null>(null);
	let loading = $state(true);

	onMount(async () => {
		const token = localStorage.getItem('token');
		const r = await fetch(`${API_BASE}/dashboard/stats`, { headers: { Authorization: `Bearer ${token}` } });
		if (r.ok) stats = await r.json();
		loading = false;
	});

	function price(c: number) { return `€${(c / 100).toFixed(2)}`; }
</script>

<div class="space-y-6">
	<h2 class="text-2xl font-bold">Dashboard</h2>

	{#if loading}
		<div class="flex justify-center py-8"><span class="loading loading-spinner loading-lg"></span></div>
	{:else if !stats}
		<div class="flex h-32 items-center justify-center rounded-box border-2 border-dashed text-base-content/40">Could not load stats</div>
	{:else}
		<!-- KPI cards -->
		<div class="grid grid-cols-2 gap-4 md:grid-cols-4">
			<div class="card bg-base-100 shadow-xl">
				<div class="card-body items-center text-center py-4">
					<span class="text-3xl font-bold text-primary">{stats.people_today}</span>
					<span class="text-sm text-base-content/60">People Served</span>
				</div>
			</div>
			<div class="card bg-base-100 shadow-xl">
				<div class="card-body items-center text-center py-4">
					<span class="text-3xl font-bold text-success">{price(stats.revenue_today_cents)}</span>
					<span class="text-sm text-base-content/60">Revenue Today</span>
				</div>
			</div>
			<div class="card bg-base-100 shadow-xl">
				<div class="card-body items-center text-center py-4">
					<span class="text-3xl font-bold text-info">{stats.orders_today}</span>
					<span class="text-sm text-base-content/60">Orders Today</span>
				</div>
			</div>
			<div class="card bg-base-100 shadow-xl">
				<div class="card-body items-center text-center py-4">
					<span class="text-3xl font-bold text-warning">{stats.tables_occupied}</span>
					<span class="text-sm text-base-content/60">Tables Now</span>
				</div>
			</div>
		</div>

		<div class="grid grid-cols-1 gap-4 md:grid-cols-2">
			<!-- Top dishes -->
			<div class="card bg-base-100 shadow-xl">
				<div class="card-body">
					<h3 class="card-title text-base">Top Dishes Today</h3>
					{#if stats.top_dishes.length === 0}
						<p class="text-sm text-base-content/40">No data yet</p>
					{:else}
						<div class="space-y-2">
							{#each stats.top_dishes as d, i}
								<div class="flex items-center justify-between">
									<div class="flex items-center gap-2">
										<span class="text-sm font-bold text-base-content/30">#{i + 1}</span>
										<span class="text-sm">{d.name}</span>
									</div>
									<div class="flex items-center gap-3 text-sm">
										<span class="text-base-content/50">×{d.count}</span>
										<span class="font-semibold tabular-nums">{price(d.total_cents)}</span>
									</div>
								</div>
							{/each}
						</div>
					{/if}
				</div>
			</div>

			<!-- Top wines -->
			<div class="card bg-base-100 shadow-xl">
				<div class="card-body">
					<h3 class="card-title text-base">Top Wines Sold</h3>
					{#if stats.top_wines.length === 0}
						<p class="text-sm text-base-content/40">No data yet</p>
					{:else}
						<div class="space-y-2">
							{#each stats.top_wines as w, i}
								<div class="flex items-center justify-between">
									<div class="flex items-center gap-2">
										<span class="text-sm font-bold text-base-content/30">#{i + 1}</span>
										<span class="text-sm">{w.name}</span>
									</div>
									<div class="flex items-center gap-3 text-sm">
										<span class="text-base-content/50">×{w.count}</span>
										<span class="font-semibold tabular-nums">{price(w.total_cents)}</span>
									</div>
								</div>
							{/each}
						</div>
					{/if}
				</div>
			</div>
		</div>

		<!-- Quick links -->
		<div class="grid grid-cols-1 gap-4 md:grid-cols-3">
			<a href="/admin/tables" class="card bg-base-100 shadow-xl hover:shadow-2xl transition-shadow">
				<div class="card-body"><h3 class="card-title">Floor Plan</h3><p class="text-sm text-base-content/60">Manage tables layout</p></div>
			</a>
			<a href="/admin/menu" class="card bg-base-100 shadow-xl hover:shadow-2xl transition-shadow">
				<div class="card-body"><h3 class="card-title">Menu</h3><p class="text-sm text-base-content/60">Dishes, categories, allergens</p></div>
			</a>
			<a href="/admin/invoices" class="card bg-base-100 shadow-xl hover:shadow-2xl transition-shadow">
				<div class="card-body"><h3 class="card-title">Invoices</h3><p class="text-sm text-base-content/60">Send invoices and mark paid</p></div>
			</a>
		</div>
	{/if}
</div>
