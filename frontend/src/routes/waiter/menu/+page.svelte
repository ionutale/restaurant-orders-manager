<script lang="ts">
	import { auth } from '$lib/stores/auth.svelte';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { API_BASE } from '$lib/config';

	if (auth.role !== 'waiter') goto('/login');

	type MenuItem = {
		id: number;
		name: string;
		description: string;
		price_cents: number;
		category_id: number;
		eating_time_min: number;
	};

	type MenuCategory = {
		category: { id: number; name: string; icon: string };
		dishes: MenuItem[];
	};

	type Suggestion = {
		id: number;
		name: string;
		description: string;
		price_cents: number;
		chef_name: string;
	};

	type Allergen = { id: number; name: string; icon: string };

	let categories = $state<MenuCategory[]>([]);
	let suggestions = $state<Suggestion[]>([]);
	let dishAllergens = $state<Record<number, Allergen[]>>({});
	let loading = $state(true);
	let selectedCategory = $state<string | null>(null);
	let selectedDish = $state<MenuItem | Suggestion | null>(null);

	const token = () => localStorage.getItem('token') ?? '';

	onMount(async () => {
		const res = await fetch(`${API_BASE}/menu`, { headers: { Authorization: `Bearer ${token()}` } });
		if (res.ok) {
			const data = await res.json();
			categories = data.categories;
			suggestions = data.suggestions;
			dishAllergens = data.dish_allergens || {};
		}
		loading = false;
	});

	function price(cents: number) {
		return `€${(cents / 100).toFixed(2)}`;
	}

	function allergensFor(dishId: number) {
		return dishAllergens[dishId] || [];
	}
</script>

<div class="space-y-4">
	<div class="flex items-center justify-between">
		<h2 class="text-2xl font-bold">Menu</h2>
	</div>

	{#if loading}
		<div class="flex justify-center py-8"><span class="loading loading-spinner loading-lg"></span></div>
	{:else}
		<div class="join flex-wrap">
			<button class="join-item btn btn-sm" class:btn-active={selectedCategory === null} onclick={() => (selectedCategory = null)}>All</button>
			<button class="join-item btn btn-sm" class:btn-active={selectedCategory === '_suggestions'} onclick={() => (selectedCategory = '_suggestions')}>
				Chef's Suggestions {#if suggestions.length > 0}<span class="badge badge-sm">{suggestions.length}</span>{/if}
			</button>
			{#each categories as cat}
				<button class="join-item btn btn-sm" class:btn-active={selectedCategory === cat.category.name} onclick={() => (selectedCategory = cat.category.name)}>
					{cat.category.icon} {cat.category.name}
				</button>
			{/each}
		</div>

		{#if selectedCategory === '_suggestions'}
			<div class="space-y-2">
				<h3 class="text-lg font-bold">Chef's Suggestions</h3>
				{#if suggestions.length === 0}
					<p class="text-base-content/50">No chef suggestions for today</p>
				{:else}
					<div class="grid grid-cols-1 gap-3 md:grid-cols-2">
						{#each suggestions as s (s.id)}
							<button class="card card-compact bg-base-100 shadow-xl text-left hover:shadow-2xl transition-shadow" onclick={() => (selectedDish = s)}>
								<div class="card-body">
									<div class="flex items-start justify-between">
										<h4 class="card-title text-base">{s.name}</h4>
										<span class="badge badge-primary">{price(s.price_cents)}</span>
									</div>
									{#if s.description}
										<p class="text-sm text-base-content/60">{s.description}</p>
									{/if}
									<p class="text-xs text-base-content/40">by {s.chef_name}</p>
								</div>
							</button>
						{/each}
					</div>
				{/if}
			</div>
		{:else}
			{#each categories as cat}
				{#if selectedCategory === null || selectedCategory === cat.category.name}
					<div class="space-y-2">
						<h3 class="text-lg font-bold">{cat.category.icon} {cat.category.name}</h3>
						{#if cat.dishes.length === 0}
							<p class="text-sm text-base-content/40">No dishes in this category</p>
						{:else}
							<div class="grid grid-cols-1 gap-2 md:grid-cols-2">
								{#each cat.dishes as d (d.id)}
									<button class="card card-compact bg-base-100 shadow-sm text-left hover:shadow-md transition-shadow" onclick={() => (selectedDish = d)}>
										<div class="card-body">
											<div class="flex items-start justify-between">
												<h4 class="card-title text-sm">{d.name}</h4>
												<span class="whitespace-nowrap text-sm font-semibold">{price(d.price_cents)}</span>
											</div>
											{#if d.description}
												<p class="text-xs text-base-content/60">{d.description}</p>
											{/if}
											<div class="flex flex-wrap gap-1">
												{#each allergensFor(d.id) as a}
													<span class="tooltip" data-tip={a.name}>{a.icon}</span>
												{/each}
												<span class="text-xs text-base-content/40 ml-auto">{d.eating_time_min} min</span>
											</div>
										</div>
									</button>
								{/each}
							</div>
						{/if}
					</div>
				{/if}
			{/each}
		{/if}
	{/if}
</div>

{#if selectedDish}
	<div class="modal modal-open">
		<div class="modal-box">
			<h3 class="font-bold text-lg">{selectedDish.name}</h3>
			<div class="py-4 space-y-3">
				{#if 'description' in selectedDish && selectedDish.description}
					<p class="text-sm text-base-content/70">{selectedDish.description}</p>
				{/if}
				<p class="text-2xl font-bold text-primary">{price(selectedDish.price_cents)}</p>
				{#if 'eating_time_min' in selectedDish}
					<p class="text-sm text-base-content/50">Prep time: {selectedDish.eating_time_min} min</p>
				{/if}
				{#if 'chef_name' in selectedDish && selectedDish.chef_name}
					<p class="text-sm text-base-content/50">by {selectedDish.chef_name}</p>
				{/if}
				{#if 'id' in selectedDish && allergensFor(selectedDish.id).length > 0}
					<div>
						<p class="text-sm font-semibold mb-1">Allergens:</p>
						<div class="flex flex-wrap gap-2">
							{#each allergensFor(selectedDish.id) as a}
								<span class="badge badge-ghost">{a.icon} {a.name}</span>
							{/each}
						</div>
					</div>
				{/if}
			</div>
			<div class="modal-action">
				<button class="btn" onclick={() => (selectedDish = null)}>Close</button>
			</div>
		</div>
	</div>
{/if}
