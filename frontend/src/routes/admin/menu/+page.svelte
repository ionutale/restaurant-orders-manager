<script lang="ts">
	import { auth } from '$lib/stores/auth.svelte';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { API_BASE } from '$lib/config';

	if (auth.role !== 'manager') goto('/login');

	let tab = $state<'categories' | 'dishes'>('categories');
	let categories = $state<{ id: number; name: string; display_order: number; icon: string }[]>([]);
	let dishes = $state<any[]>([]);
	let allAllergens = $state<{ id: number; name: string; icon: string }[]>([]);
	let loading = $state(true);
	let selCategory = $state<string>('');
	let editingId = $state<number | null>(null);
	let editName = $state('');
	let editIcon = $state('');
	let deleteId = $state<number | null>(null);
	let adding = $state(false);
	let newName = $state('');
	let newIcon = $state('');

	let dishEditingId = $state<number | null>(null);
	let dishEditName = $state('');
	let dishEditDesc = $state('');
	let dishEditPrice = $state(0);
	let dishEditCat = $state(0);
	let dishEditTime = $state(10);
	let dishEditImage = $state('');
	let dishDeleteId = $state<number | null>(null);
	let dishAdding = $state(false);
	let dName = $state('');
	let dDesc = $state('');
	let dPrice = $state(0);
	let dCat = $state(0);
	let dTime = $state(10);
	let dImage = $state('');

	let detailDishId = $state<number | null>(null);
	let detailDishName = $state('');
	let detailDishAllergens = $state<number[]>([]);
	let detailWineSuggestions = $state<any[]>([]);
	let detailSideSuggestions = $state<any[]>([]);
	let addWineId = $state(0);
	let addSideId = $state(0);

	const token = () => localStorage.getItem('token') ?? '';

	async function loadCats() {
		const res = await fetch(`${API_BASE}/categories`, { headers: { Authorization: `Bearer ${token()}` } });
		if (res.ok) categories = await res.json();
	}

	async function loadDishes() {
		const url = selCategory ? `${API_BASE}/dishes?category_id=${selCategory}` : `${API_BASE}/dishes`;
		const res = await fetch(url, { headers: { Authorization: `Bearer ${token()}` } });
		if (res.ok) dishes = await res.json();
	}

	async function loadAllergens() {
		const res = await fetch(`${API_BASE}/allergens`, { headers: { Authorization: `Bearer ${token()}` } });
		if (res.ok) allAllergens = await res.json();
	}

	async function loadAll() {
		loading = true;
		await Promise.all([loadCats(), loadDishes(), loadAllergens()]);
		loading = false;
	}

	onMount(loadAll);

	async function openDetail(d: any) {
		detailDishId = d.id;
		detailDishName = d.name;

		const [aRes, sRes] = await Promise.all([
			fetch(`${API_BASE}/dishes/${d.id}/allergens`, { headers: { Authorization: `Bearer ${token()}` } }),
			fetch(`${API_BASE}/dishes/${d.id}/suggestions`, { headers: { Authorization: `Bearer ${token()}` } }),
		]);
		if (aRes.ok) detailDishAllergens = (await aRes.json()).map((a: any) => a.id);
		if (sRes.ok) {
			const all = await sRes.json();
			detailWineSuggestions = all.filter((s: any) => s.suggestion_type === 'wine');
			detailSideSuggestions = all.filter((s: any) => s.suggestion_type === 'side');
		}
		addWineId = 0;
		addSideId = 0;
	}

	async function saveAllergens() {
		if (!detailDishId) return;
		await fetch(`${API_BASE}/dishes/${detailDishId}/allergens`, {
			method: 'PUT', headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token()}` },
			body: JSON.stringify({ allergen_ids: detailDishAllergens }),
		});
	}

	async function addSuggestion(type: string) {
		if (!detailDishId) return;
		const toId = type === 'wine' ? addWineId : addSideId;
		if (!toId) return;
		await fetch(`${API_BASE}/dishes/${detailDishId}/suggestions`, {
			method: 'POST', headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token()}` },
			body: JSON.stringify({ to_dish_id: toId, suggestion_type: type }),
		});
		const sRes = await fetch(`${API_BASE}/dishes/${detailDishId}/suggestions`, { headers: { Authorization: `Bearer ${token()}` } });
		if (sRes.ok) {
			const all = await sRes.json();
			detailWineSuggestions = all.filter((s: any) => s.suggestion_type === 'wine');
			detailSideSuggestions = all.filter((s: any) => s.suggestion_type === 'side');
		}
		addWineId = 0; addSideId = 0;
	}

	async function removeSuggestion(id: number) {
		await fetch(`${API_BASE}/dish-suggestions/${id}`, { method: 'DELETE', headers: { Authorization: `Bearer ${token()}` } });
		const sRes = await fetch(`${API_BASE}/dishes/${detailDishId}/suggestions`, { headers: { Authorization: `Bearer ${token()}` } });
		if (sRes.ok) {
			const all = await sRes.json();
			detailWineSuggestions = all.filter((s: any) => s.suggestion_type === 'wine');
			detailSideSuggestions = all.filter((s: any) => s.suggestion_type === 'side');
		}
	}

	async function create() {
		if (!newName) return;
		await fetch(`${API_BASE}/categories`, {
			method: 'POST', headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token()}` },
			body: JSON.stringify({ name: newName, icon: newIcon || undefined }),
		});
		newName = ''; newIcon = ''; adding = false; await loadCats();
	}

	async function saveCat(id: number) {
		await fetch(`${API_BASE}/categories/${id}`, {
			method: 'PATCH', headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token()}` },
			body: JSON.stringify({ name: editName, icon: editIcon || undefined }),
		});
		editingId = null; await loadCats();
	}

	async function removeCat(id: number) {
		await fetch(`${API_BASE}/categories/${id}`, { method: 'DELETE', headers: { Authorization: `Bearer ${token()}` } });
		deleteId = null; await loadAll();
	}

	async function moveCat(id: number, delta: number) {
		await fetch(`${API_BASE}/categories/reorder`, {
			method: 'POST', headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token()}` },
			body: JSON.stringify({ id, delta }),
		});
		await loadCats();
	}

	function startEditCat(c: any) { editingId = c.id; editName = c.name; editIcon = c.icon; }

	async function createDish() {
		if (!dName || !dCat) return;
		let url = dImage;
		const fileInput = document.getElementById('create-image-upload') as HTMLInputElement;
		if (fileInput?.files?.[0]) url = await uploadImage(fileInput.files[0]) || dImage;
		await fetch(`${API_BASE}/dishes`, {
			method: 'POST', headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token()}` },
			body: JSON.stringify({ name: dName, description: dDesc, price_cents: Math.round(dPrice * 100), category_id: dCat, eating_time_min: dTime, image_url: url }),
		});
		dName = ''; dDesc = ''; dPrice = 0; dCat = 0; dTime = 10; dImage = ''; dishAdding = false; await loadDishes();
	}

	async function uploadImage(file: File): Promise<string> {
		const fd = new FormData();
		fd.append('file', file);
		const r = await fetch(`${API_BASE}/upload`, { method: 'POST', headers: { Authorization: `Bearer ${token()}` }, body: fd });
		if (r.ok) { const d = await r.json(); return d.url; }
		return '';
	}

	function startEditDish(d: any) {
		dishEditingId = d.id; dishEditName = d.name; dishEditDesc = d.description;
		dishEditPrice = d.price_cents; dishEditCat = d.category_id; dishEditTime = d.eating_time_min;
		dishEditImage = d.image_url || '';
	}

	async function saveDish(id: number) {
		let url = dishEditImage;
		const fileInput = document.getElementById('edit-image-upload') as HTMLInputElement;
		if (fileInput?.files?.[0]) url = await uploadImage(fileInput.files[0]) || dishEditImage;
		await fetch(`${API_BASE}/dishes/${id}`, {
			method: 'PATCH', headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token()}` },
			body: JSON.stringify({ name: dishEditName, description: dishEditDesc, price_cents: dishEditPrice, category_id: dishEditCat, eating_time_min: dishEditTime, image_url: url }),
		});
		dishEditingId = null; await loadDishes();
	}

	async function removeDish(id: number) {
		await fetch(`${API_BASE}/dishes/${id}`, { method: 'DELETE', headers: { Authorization: `Bearer ${token()}` } });
		dishDeleteId = null; await loadDishes();
	}

	function catName(id: number) { return categories.find(c => c.id === id)?.name ?? '—'; }

	function wines() { return dishes.filter(d => catName(d.category_id).toLowerCase() === 'wines').filter(d => d.id !== detailDishId); }
	function sides() { return dishes.filter(d => catName(d.category_id).toLowerCase() !== 'wines').filter(d => d.id !== detailDishId); }
</script>

<div class="space-y-4">
	<div class="flex items-center justify-between">
		<h2 class="text-2xl font-bold">Menu</h2>
		<div class="tabs tabs-box">
			<button class="tab" class:tab-active={tab === 'categories'} onclick={() => (tab = 'categories')}>Categories</button>
			<button class="tab" class:tab-active={tab === 'dishes'} onclick={() => (tab = 'dishes')}>Dishes</button>
		</div>
	</div>

	{#if loading}
		<div class="flex justify-center py-8"><span class="loading loading-spinner loading-lg"></span></div>
	{:else if tab === 'categories'}
		<div class="flex justify-end"><button class="btn btn-primary btn-sm" onclick={() => (adding = true)}>Add Category</button></div>
		<div class="overflow-x-auto">
			<table class="table table-zebra">
				<thead><tr><th>#</th><th>Name</th><th>Icon</th><th>Actions</th></tr></thead>
				<tbody>
					{#each categories as c (c.id)}
						<tr>
							{#if editingId === c.id}
								<td class="text-center text-base-content/50">{c.display_order}</td>
								<td><input type="text" bind:value={editName} class="input input-bordered input-sm w-full" /></td>
								<td><input type="text" bind:value={editIcon} class="input input-bordered input-sm w-16" /></td>
								<td class="flex gap-1">
									<button class="btn btn-ghost btn-xs text-success" onclick={() => saveCat(c.id)}>Save</button>
									<button class="btn btn-ghost btn-xs" onclick={() => (editingId = null)}>Cancel</button>
								</td>
							{:else if deleteId === c.id}
								<td colspan="2" class="text-warning">Delete "{c.name}"?</td>
								<td colspan="2" class="flex gap-1">
									<button class="btn btn-ghost btn-xs text-error" onclick={() => removeCat(c.id)}>Confirm</button>
									<button class="btn btn-ghost btn-xs" onclick={() => (deleteId = null)}>Cancel</button>
								</td>
							{:else}
								<td class="text-center text-base-content/50">{c.display_order}</td>
								<td class="font-medium">{c.name}</td>
								<td>{c.icon || '—'}</td>
								<td class="flex gap-1">
									<button class="btn btn-ghost btn-xs" onclick={() => moveCat(c.id, -1)} disabled={c.display_order <= 1}>↑</button>
									<button class="btn btn-ghost btn-xs" onclick={() => moveCat(c.id, 1)} disabled={c.display_order >= categories.length}>↓</button>
									<button class="btn btn-ghost btn-xs" onclick={() => startEditCat(c)}>Edit</button>
									<button class="btn btn-ghost btn-xs text-error" onclick={() => (deleteId = c.id)}>Delete</button>
								</td>
							{/if}
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	{:else}
		<div class="flex flex-wrap items-center justify-between gap-2">
			<div class="join">
				<button class="join-item btn btn-sm" class:btn-active={!selCategory} onclick={() => { selCategory = ''; loadDishes(); }}>All</button>
				{#each categories as c}
					<button class="join-item btn btn-sm" class:btn-active={selCategory === String(c.id)} onclick={() => { selCategory = String(c.id); loadDishes(); }}>{c.name}</button>
				{/each}
			</div>
			<button class="btn btn-primary btn-sm" onclick={() => { dCat = categories[0]?.id ?? 0; dishAdding = true; }}>Add Dish</button>
		</div>

		<div class="overflow-x-auto">
			<table class="table table-zebra">
				<thead><tr><th>Name</th><th>Category</th><th>Price</th><th>Time</th><th>Actions</th></tr></thead>
				<tbody>
					{#if dishes.length === 0}
						<tr><td colspan="5" class="text-center text-base-content/50">No dishes</td></tr>
					{:else}
						{#each dishes as d (d.id)}
							<tr>
								{#if dishEditingId === d.id}
									<td><input type="text" bind:value={dishEditName} class="input input-bordered input-sm w-full" /></td>
									<td><select bind:value={dishEditCat} class="select select-bordered select-sm">{#each categories as c}<option value={c.id}>{c.name}</option>{/each}</select></td>
									<td><input type="number" bind:value={dishEditPrice} class="input input-bordered input-sm w-24" step="1" /></td>
									<td><input type="number" bind:value={dishEditTime} class="input input-bordered input-sm w-16" min="1" /></td>
									<td class="flex gap-1">
										<button class="btn btn-ghost btn-xs text-success" onclick={() => saveDish(d.id)}>Save</button>
										<button class="btn btn-ghost btn-xs" onclick={() => (dishEditingId = null)}>Cancel</button>
									</td>
								{:else if dishDeleteId === d.id}
									<td colspan="3" class="text-warning">Delete "{d.name}"?</td>
									<td colspan="2" class="flex gap-1">
										<button class="btn btn-ghost btn-xs text-error" onclick={() => removeDish(d.id)}>Confirm</button>
										<button class="btn btn-ghost btn-xs" onclick={() => (dishDeleteId = null)}>Cancel</button>
									</td>
								{:else}
									<td class="font-medium">{d.name}</td>
									<td>{d.category_name}</td>
									<td>€{(d.price_cents / 100).toFixed(2)}</td>
									<td>{d.eating_time_min} min</td>
									<td class="flex gap-1">
										<button class="btn btn-ghost btn-xs" onclick={() => startEditDish(d)}>Edit</button>
										<button class="btn btn-ghost btn-xs" onclick={() => openDetail(d)}>Details</button>
										<button class="btn btn-ghost btn-xs text-error" onclick={() => (dishDeleteId = d.id)}>Delete</button>
									</td>
								{/if}
							</tr>
						{/each}
					{/if}
				</tbody>
			</table>
		</div>
	{/if}
</div>

{#if adding}
	<div class="modal modal-open">
		<div class="modal-box"><h3 class="font-bold text-lg">Add Category</h3>
			<div class="py-4 space-y-3">
				<label class="form-control"><div class="label"><span class="label-text">Name</span></div><input type="text" bind:value={newName} placeholder="Appetizers" class="input input-bordered" /></label>
				<label class="form-control"><div class="label"><span class="label-text">Icon</span></div><input type="text" bind:value={newIcon} placeholder="🍕" class="input input-bordered" /></label>
			</div>
			<div class="modal-action"><button class="btn" onclick={() => (adding = false)}>Cancel</button><button class="btn btn-primary" onclick={create}>Add</button></div>
		</div>
	</div>
{/if}

{#if dishAdding}
	<div class="modal modal-open">
		<div class="modal-box"><h3 class="font-bold text-lg">Add Dish</h3>
			<div class="py-4 space-y-3">
				<label class="form-control"><div class="label"><span class="label-text">Name</span></div><input type="text" bind:value={dName} class="input input-bordered" /></label>
				<label class="form-control"><div class="label"><span class="label-text">Description</span></div><textarea bind:value={dDesc} class="textarea textarea-bordered" rows="2"></textarea></label>
				<div class="flex gap-3">
					<label class="form-control flex-1"><div class="label"><span class="label-text">Price (EUR)</span></div><input type="number" bind:value={dPrice} class="input input-bordered" min="0" step="0.01" /></label>
					<label class="form-control w-24"><div class="label"><span class="label-text">Time (min)</span></div><input type="number" bind:value={dTime} class="input input-bordered" min="1" /></label>
				</div>
				<label class="form-control"><div class="label"><span class="label-text">Category</span></div>
					<select bind:value={dCat} class="select select-bordered"><option value={0} disabled>Select</option>{#each categories as c}<option value={c.id}>{c.name}</option>{/each}</select>
				</label>
				<label class="form-control"><div class="label"><span class="label-text">Image</span></div>
					<input type="file" id="create-image-upload" accept="image/jpeg,image/png,image/webp" class="file-input file-input-bordered" />
				</label>
			</div>
			<div class="modal-action"><button class="btn" onclick={() => (dishAdding = false)}>Cancel</button><button class="btn btn-primary" onclick={createDish}>Add</button></div>
		</div>
	</div>
{/if}

{#if detailDishId !== null}
	<div class="modal modal-open">
		<div class="modal-box max-w-2xl">
			<h3 class="font-bold text-lg">{detailDishName} — Details</h3>

			<div class="py-4 space-y-6">
				<div>
					<h4 class="font-semibold mb-2">Allergens</h4>
					<div class="flex flex-wrap gap-2">
						{#each allAllergens as a}
							<button
								class="btn btn-sm"
								class:btn-primary={detailDishAllergens.includes(a.id)}
								class:btn-outline={!detailDishAllergens.includes(a.id)}
								onclick={() => {
									if (detailDishAllergens.includes(a.id)) {
										detailDishAllergens = detailDishAllergens.filter(id => id !== a.id);
									} else {
										detailDishAllergens = [...detailDishAllergens, a.id];
									}
									saveAllergens();
								}}
							>{a.icon} {a.name}</button>
						{/each}
					</div>
				</div>

				<div>
					<h4 class="font-semibold mb-2">Suggested Wines</h4>
					<div class="flex flex-wrap gap-2 mb-2">
						{#each detailWineSuggestions as s}
							<div class="badge badge-lg gap-1">
								{s.to_dish_name}
								<button class="btn btn-ghost btn-xs text-error" onclick={() => removeSuggestion(s.id)}>✕</button>
							</div>
						{:else}
							<span class="text-sm text-base-content/40">None</span>
						{/each}
					</div>
					<div class="flex gap-2">
						<select bind:value={addWineId} class="select select-bordered select-sm flex-1">
							<option value={0} disabled>Add a wine...</option>
							{#each wines() as w}
								<option value={w.id}>€{(w.price_cents / 100).toFixed(2)} — {w.name}</option>
							{/each}
						</select>
						<button class="btn btn-primary btn-sm" onclick={() => addSuggestion('wine')} disabled={!addWineId}>Add</button>
					</div>
				</div>

				<div>
					<h4 class="font-semibold mb-2">Suggested Sides</h4>
					<div class="flex flex-wrap gap-2 mb-2">
						{#each detailSideSuggestions as s}
							<div class="badge badge-lg gap-1">
								{s.to_dish_name}
								<button class="btn btn-ghost btn-xs text-error" onclick={() => removeSuggestion(s.id)}>✕</button>
							</div>
						{:else}
							<span class="text-sm text-base-content/40">None</span>
						{/each}
					</div>
					<div class="flex gap-2">
						<select bind:value={addSideId} class="select select-bordered select-sm flex-1">
							<option value={0} disabled>Add a side...</option>
							{#each sides() as s}
								<option value={s.id}>€{(s.price_cents / 100).toFixed(2)} — {s.name}</option>
							{/each}
						</select>
						<button class="btn btn-primary btn-sm" onclick={() => addSuggestion('side')} disabled={!addSideId}>Add</button>
					</div>
				</div>
			</div>

			<div class="modal-action"><button class="btn" onclick={() => (detailDishId = null)}>Close</button></div>
		</div>
	</div>
{/if}
