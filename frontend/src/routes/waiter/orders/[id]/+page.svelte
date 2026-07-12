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
	type Item = { id: number; dish_name: string; quantity: number; notes: string };
	type MenuCat = { category: { id: number; name: string; icon: string }; dishes: { id: number; name: string; description: string; price_cents: number }[] };
	type Suggestion = { id: number; name: string; description: string; price_cents: number };

	let order = $state<Order | null>(null);
	let menuCats = $state<MenuCat[]>([]);
	let menuSuggestions = $state<Suggestion[]>([]);
	let dishAllergens = $state<Record<number, { name: string; icon: string }[]>>({});
	let loading = $state(true);
	let activeCourseId = $state<number | null>(null);
	let newCourseName = $state('');
	let addingId = $state<number | null>(null);
	let sending = $state(false);
	let search = $state('');
	let menuRef: HTMLDivElement;
	let dragItemId = $state<number | null>(null);

	const token = () => localStorage.getItem('token') ?? '';

	onMount(async () => {
		const [oRes, mRes] = await Promise.all([
			fetch(`${API_BASE}/orders/${params.id}`, { headers: { Authorization: `Bearer ${token()}` } }),
			fetch(`${API_BASE}/menu`, { headers: { Authorization: `Bearer ${token()}` } }),
		]);
		if (oRes.ok) { order = await oRes.json(); if (order.courses.length > 0) activeCourseId = order.courses[0].id; }
		if (mRes.ok) { const d = await mRes.json(); menuCats = d.categories; menuSuggestions = d.suggestions; dishAllergens = d.dish_allergens || {}; }
		loading = false;
	});

	async function reload() {
		const r = await fetch(`${API_BASE}/orders/${params.id}`, { headers: { Authorization: `Bearer ${token()}` } });
		if (r.ok) order = await r.json();
	}

	async function addCourse() {
		if (!newCourseName.trim()) return;
		await fetch(`${API_BASE}/orders/${params.id}/courses`, { method: 'POST', headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token()}` }, body: JSON.stringify({ name: newCourseName.trim() }) });
		newCourseName = '';
		await reload();
	}

	async function addDish(dishId: number) {
		if (activeCourseId === null) return;
		addingId = dishId;
		await fetch(`${API_BASE}/orders/${params.id}/courses/${activeCourseId}/items`, {
			method: 'POST', headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token()}` },
			body: JSON.stringify({ dish_id: dishId, quantity: 1 }),
		});
		addingId = null;
		await reload();
	}

	async function addSuggestion(sugId: number) {
		if (activeCourseId === null) return;
		addingId = -sugId;
		await fetch(`${API_BASE}/orders/${params.id}/courses/${activeCourseId}/items`, {
			method: 'POST', headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token()}` },
			body: JSON.stringify({ chef_suggestion_id: sugId, is_chef_suggestion: true, quantity: 1 }),
		});
		addingId = null;
		await reload();
	}

	async function removeItem(itemId: number) {
		await fetch(`${API_BASE}/order-items/${itemId}`, { method: 'DELETE', headers: { Authorization: `Bearer ${token()}` } });
		await reload();
	}

	async function moveItem(itemId: number, targetCourseId: number) {
		await fetch(`${API_BASE}/order-items/${itemId}/move`, {
			method: 'PATCH', headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token()}` },
			body: JSON.stringify({ course_id: targetCourseId }),
		});
		dragItemId = null;
		await reload();
	}

	function dragStart(itemId: number) { dragItemId = itemId; }
	function allowDrop(e: DragEvent) { e.preventDefault(); }
	function dropOnCourse(e: DragEvent, courseId: number) {
		e.preventDefault();
		if (dragItemId !== null) moveItem(dragItemId, courseId);
	}

	async function sendToKDS() { sending = true; await fetch(`${API_BASE}/orders/${params.id}/send`, { method: 'POST', headers: { Authorization: `Bearer ${token()}` } }); sending = false; await reload(); }
	async function advanceCourse() { await fetch(`${API_BASE}/orders/${params.id}/advance-course`, { method: 'POST', headers: { Authorization: `Bearer ${token()}` } }); await reload(); }

	function activeCourse() { return order?.courses.find(c => c.id === activeCourseId); }
	function price(c: number) { return `€${(c / 100).toFixed(2)}`; }
	function alFor(dishId: number) { return dishAllergens[dishId] || []; }
</script>

<div class="space-y-3">
	<button class="btn btn-ghost btn-xs" onclick={() => goto('/waiter/orders')}>← Orders</button>

	{#if loading}
		<div class="flex justify-center py-8"><span class="loading loading-spinner loading-lg"></span></div>
	{:else if !order}
		<div class="text-center py-8 text-base-content/50">Order not found</div>
	{:else}
		<div class="flex items-center justify-between">
			<h2 class="text-xl font-bold">Order #{order.id}</h2>
			<div class="flex items-center gap-2">
				<span class="badge">{order.status}</span>
				{#if order.status === 'pending'}
					<button class="btn btn-primary btn-xs" onclick={sendToKDS} disabled={sending}>Send to KDS</button>
				{/if}
				{#if order.status === 'sent'}
					<button class="btn btn-warning btn-xs" onclick={advanceCourse}>Advance</button>
				{/if}
			</div>
		</div>

		<!-- Course pills -->
		<div class="flex flex-wrap gap-1">
			{#each order.courses as c (c.id)}
				<button
					class="btn btn-sm"
					class:btn-primary={activeCourseId === c.id}
					class:btn-ghost={activeCourseId !== c.id}
					class:border-dashed={dragItemId !== null}
					onclick={() => (activeCourseId = c.id)}
					ondragover={allowDrop}
					ondrop={(e) => dropOnCourse(e, c.id)}
				>
					{c.name}
					<span class="badge badge-sm">{c.items?.length ?? 0}</span>
				</button>
			{/each}
			{#if order.status === 'pending'}
				<div class="flex gap-1 items-center">
					<input type="text" bind:value={newCourseName} placeholder="+ Course" class="input input-bordered input-xs w-20"
						onkeydown={(e) => { if (e.key === 'Enter') addCourse(); }} />
					<button class="btn btn-ghost btn-xs" onclick={addCourse}>+</button>
				</div>
			{/if}
		</div>

		<!-- Active course items -->
		{#if activeCourse()}
			<div class="space-y-1">
				{#if (activeCourse()?.items?.length ?? 0) === 0}
					<p class="text-sm text-base-content/40">No dishes yet — tap + below to add</p>
				{:else}
					{#each activeCourse()!.items as item (item.id)}
						<div class="flex items-center justify-between rounded-box bg-base-200 px-3 py-1.5"
							draggable={order.status === 'pending'}
							ondragstart={() => dragStart(item.id)}
							class:opacity-50={dragItemId === item.id}
						>
							<div class="flex items-center gap-2">
								{#if order.status === 'pending'}
									<span class="cursor-grab text-base-content/30">⠿</span>
								{/if}
								<span class="font-medium text-sm">×{item.quantity}</span>
								<span class="text-sm">{item.dish_name}</span>
								{#if item.notes}<span class="badge badge-ghost badge-xs">"{item.notes}"</span>{/if}
							</div>
							{#if order.status === 'pending'}
								<button class="btn btn-ghost btn-xs text-error" onclick={() => removeItem(item.id)}>✕</button>
							{/if}
						</div>
					{/each}
				{/if}
			</div>
		{/if}

		<!-- Inline menu browser (one-tap add) -->
		{#if order.status === 'pending'}
			<div class="pt-2" bind:this={menuRef}>
				<div class="flex items-center gap-2 mb-2">
					<h3 class="text-sm font-semibold cursor-pointer" onclick={() => menuRef?.scrollIntoView({ behavior: 'smooth' })}>
						Add dishes to {activeCourse()?.name ?? 'course'}
					</h3>
					<input type="text" bind:value={search} placeholder="Search dishes..."
						class="input input-bordered input-sm flex-1 max-w-xs"
						oninput={() => {}} />
				</div>

				{#if search.trim()}
					{@const q = search.toLowerCase()}
					<div class="flex flex-wrap gap-1 mb-2">
						{#each menuSuggestions.filter(s => s.name.toLowerCase().includes(q)) as s}
							<button class="btn btn-outline btn-xs gap-1" onclick={() => addSuggestion(s.id)} disabled={!activeCourseId || addingId === -s.id}>
								{#if addingId === -s.id}<span class="loading loading-spinner loading-xs"></span>{:else}+{/if}
								{s.name} <span class="text-base-content/50">{price(s.price_cents)}</span>
							</button>
						{/each}
						{#each menuCats as cat}
							{#each cat.dishes.filter(d => d.name.toLowerCase().includes(q)) as d}
								<button class="btn btn-xs btn-ghost gap-1 {addingId === d.id ? 'btn-primary' : ''}"
									onclick={() => addDish(d.id)} disabled={!activeCourseId || addingId === d.id}>
									{#if addingId === d.id}<span class="loading loading-spinner loading-xs"></span>{:else}<span class="text-lg leading-none">+</span>{/if}
									<span class="text-xs">{d.name}</span>
									<span class="text-xs text-base-content/50">{price(d.price_cents)}</span>
								</button>
							{/each}
						{/each}
					</div>
					{#if [...menuSuggestions.filter(s => s.name.toLowerCase().includes(q)), ...menuCats.flatMap(c => c.dishes.filter(d => d.name.toLowerCase().includes(q)))].length === 0}
						<p class="text-xs text-base-content/40 py-1">No dishes match "{search}"</p>
					{/if}
				{:else}
					{#if menuSuggestions.length > 0}
					<div class="flex flex-wrap gap-1 mb-2">
						{#each menuSuggestions as s}
							<button class="btn btn-outline btn-xs gap-1" onclick={() => addSuggestion(s.id)} disabled={!activeCourseId || addingId === -s.id}>
								{#if addingId === -s.id}<span class="loading loading-spinner loading-xs"></span>{:else}+{/if}
								{s.name}
								<span class="text-base-content/50">{price(s.price_cents)}</span>
							</button>
						{/each}
					</div>
					{/if}

					{#each menuCats as cat}
						<details class="mb-1" open={false}>
							<summary class="cursor-pointer text-sm font-medium py-1">{cat.category.icon} {cat.category.name}</summary>
							<div class="flex flex-wrap gap-1 pt-1">
								{#each cat.dishes as d}
									<button class="btn btn-xs gap-1" class:btn-ghost class:btn-primary={addingId === d.id}
										onclick={() => addDish(d.id)} disabled={!activeCourseId || addingId === d.id}>
										{#if addingId === d.id}<span class="loading loading-spinner loading-xs"></span>{:else}<span class="text-lg leading-none">+</span>{/if}
										<span class="text-xs">{d.name}</span>
										<span class="text-xs text-base-content/50">{price(d.price_cents)}</span>
									</button>
								{/each}
							</div>
						</details>
					{/each}
				{/if}
			</div>
		{/if}
	{/if}
</div>
