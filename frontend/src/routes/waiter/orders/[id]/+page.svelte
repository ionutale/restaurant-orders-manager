<script lang="ts">
	import { auth } from '$lib/stores/auth.svelte';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { API_BASE } from '$lib/config';

	if (auth.role !== 'waiter') goto('/login');

	let { params } = $props();

	type Order = {
		id: number; table_group_id: number; status: string;
		courses: { id: number; name: string; status: string; items: Item[] }[];
	};
	type Item = { id: number; dish_id: number | null; dish_name: string; quantity: number; notes: string; is_chef_suggestion: boolean };
	type MenuCat = { category: { id: number; name: string; icon: string }; dishes: { id: number; name: string; description: string; price_cents: number; eating_time_min: number }[] };
	type Suggestion = { id: number; name: string; description: string; price_cents: number; chef_name: string };

	let order = $state<Order | null>(null);
	let menuCats = $state<MenuCat[]>([]);
	let menuSuggestions = $state<Suggestion[]>([]);
	let dishAllergens = $state<Record<number, { id: number; name: string; icon: string }[]>>({});
	let loading = $state(true);
	let showMenu = $state<number | null>(null);
	let sending = $state(false);
	let addQty = $state(1);
	let addNotes = $state('');

	const token = () => localStorage.getItem('token') ?? '';

	onMount(async () => {
		const [oRes, mRes] = await Promise.all([
			fetch(`${API_BASE}/orders/${params.id}`, { headers: { Authorization: `Bearer ${token()}` } }),
			fetch(`${API_BASE}/menu`, { headers: { Authorization: `Bearer ${token()}` } }),
		]);
		if (oRes.ok) order = await oRes.json();
		if (mRes.ok) { const d = await mRes.json(); menuCats = d.categories; menuSuggestions = d.suggestions; dishAllergens = d.dish_allergens || {}; }
		loading = false;
	});

	async function sendToKDS() {
		sending = true;
		const r = await fetch(`${API_BASE}/orders/${params.id}/send`, { method: 'POST', headers: { Authorization: `Bearer ${token()}` } });
		if (r.ok) order = await r.json();
		sending = false;
	}

	async function advanceCourse() {
		const r = await fetch(`${API_BASE}/orders/${params.id}/advance-course`, { method: 'POST', headers: { Authorization: `Bearer ${token()}` } });
		if (r.ok) order = await r.json();
	}

	async function addDish(courseId: number, dishId: number) {
		await fetch(`${API_BASE}/orders/${params.id}/courses/${courseId}/items`, {
			method: 'POST', headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token()}` },
			body: JSON.stringify({ dish_id: dishId, quantity: addQty, notes: addNotes }),
		});
		addQty = 1; addNotes = ''; showMenu = null;
		const r = await fetch(`${API_BASE}/orders/${params.id}`, { headers: { Authorization: `Bearer ${token()}` } });
		if (r.ok) order = await r.json();
	}

	async function addSuggestion(courseId: number, suggestionId: number) {
		await fetch(`${API_BASE}/orders/${params.id}/courses/${courseId}/items`, {
			method: 'POST', headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token()}` },
			body: JSON.stringify({ chef_suggestion_id: suggestionId, is_chef_suggestion: true, quantity: addQty, notes: addNotes }),
		});
		addQty = 1; addNotes = ''; showMenu = null;
		const r = await fetch(`${API_BASE}/orders/${params.id}`, { headers: { Authorization: `Bearer ${token()}` } });
		if (r.ok) order = await r.json();
	}

	async function removeItem(itemId: number) {
		await fetch(`${API_BASE}/order-items/${itemId}`, { method: 'DELETE', headers: { Authorization: `Bearer ${token()}` } });
		const r = await fetch(`${API_BASE}/orders/${params.id}`, { headers: { Authorization: `Bearer ${token()}` } });
		if (r.ok) order = await r.json();
	}

	function price(cents: number) { return `€${(cents / 100).toFixed(2)}`; }
	function allergensFor(dishId: number) { return dishAllergens[dishId] || []; }
</script>

<div class="space-y-4">
	<button class="btn btn-ghost btn-sm" onclick={() => goto('/waiter/orders')}>← Orders</button>

	{#if loading}
		<div class="flex justify-center py-8"><span class="loading loading-spinner loading-lg"></span></div>
	{:else if !order}
		<div class="text-center py-8 text-base-content/50">Order not found</div>
	{:else}
		<div class="flex items-center justify-between">
			<h2 class="text-2xl font-bold">Order #{order.id}</h2>
			<div class="flex items-center gap-2">
				<span class="badge badge-lg">{order.status}</span>
				{#if order.status === 'pending'}
					<button class="btn btn-primary btn-sm" onclick={sendToKDS} disabled={sending}>
						{sending ? 'Sending...' : 'Send to KDS'}
					</button>
				{/if}
				{#if order.status === 'sent'}
					<button class="btn btn-warning btn-sm" onclick={advanceCourse}>Advance Course</button>
				{/if}
			</div>
		</div>

		<div class="space-y-4">
			{#each order.courses as course (course.id)}
				<div class="card bg-base-100 shadow-xl">
					<div class="card-body">
						<div class="flex items-center justify-between">
							<h3 class="card-title text-base">
								{course.name}
								<span class="badge badge-sm">{course.status}</span>
							</h3>
							{#if order.status === 'pending'}<button class="btn btn-primary btn-sm" onclick={() => { showMenu = course.id; addQty = 1; addNotes = ''; }}>+ Add Dish</button>{/if}
						</div>

						{#if (course.items?.length ?? 0) === 0}
							<p class="text-sm text-base-content/40 py-2">No items yet</p>
						{:else}
							<div class="space-y-2">
								{#each course.items as item (item.id)}
									<div class="flex items-start justify-between gap-2 rounded-box bg-base-200 p-2">
										<div class="flex-1">
											<div class="flex items-center gap-2">
												<span class="font-medium">{item.dish_name}</span>
												<span class="text-sm text-base-content/50">×{item.quantity}</span>
												{#if item.notes}
													<span class="badge badge-ghost badge-xs">"{item.notes}"</span>
												{/if}
											</div>
										</div>
										<button class="btn btn-ghost btn-xs text-error" onclick={() => removeItem(item.id)}>✕</button>
									</div>
								{/each}
							</div>
						{/if}
					</div>
				</div>
			{/each}
		</div>
	{/if}
</div>

{#if showMenu !== null}
	<div class="modal modal-open">
		<div class="modal-box max-w-2xl">
			<h3 class="font-bold text-lg">Add Dish to {order?.courses.find(c => c.id === showMenu)?.name}</h3>
			<div class="py-4 space-y-4 max-h-[60vh] overflow-y-auto">
				{#if menuSuggestions.length > 0}
					<div>
						<h4 class="font-semibold text-sm mb-2">Chef's Suggestions</h4>
						<div class="grid grid-cols-1 gap-2">
							{#each menuSuggestions as s}
								<button class="card card-compact bg-base-200 text-left hover:bg-base-300" onclick={() => addSuggestion(showMenu!, s.id)}>
									<div class="card-body py-2">
										<div class="flex justify-between"><span class="font-medium">{s.name}</span><span class="text-sm">{price(s.price_cents)}</span></div>
										{#if s.description}<p class="text-xs text-base-content/60">{s.description}</p>{/if}
									</div>
								</button>
							{/each}
						</div>
					</div>
				{/if}

				{#each menuCats as cat}
					<div>
						<h4 class="font-semibold text-sm mb-2">{cat.category.icon} {cat.category.name}</h4>
						<div class="grid grid-cols-1 gap-2">
							{#each cat.dishes as d}
								<button class="card card-compact bg-base-200 text-left hover:bg-base-300" onclick={() => addDish(showMenu!, d.id)}>
									<div class="card-body py-2">
										<div class="flex justify-between"><span class="font-medium">{d.name}</span><span class="text-sm">{price(d.price_cents)}</span></div>
										<div class="flex items-center gap-2 text-xs text-base-content/60">
											{#if d.description}<span>{d.description}</span>{/if}
											{#each allergensFor(d.id) as a}<span title={a.name}>{a.icon}</span>{/each}
											<span>{d.eating_time_min} min</span>
										</div>
									</div>
								</button>
							{/each}
						</div>
					</div>
				{/each}
			</div>

			<div class="flex items-center gap-3 border-t pt-3">
				<label class="form-control w-20"><div class="label"><span class="label-text">Qty</span></div><input type="number" bind:value={addQty} class="input input-bordered input-sm" min="1" /></label>
				<label class="form-control flex-1"><div class="label"><span class="label-text">Notes (optional)</span></div><input type="text" bind:value={addNotes} placeholder="e.g. no onions" class="input input-bordered input-sm" /></label>
			</div>

			<div class="modal-action"><button class="btn" onclick={() => (showMenu = null)}>Close</button></div>
		</div>
	</div>
{/if}
