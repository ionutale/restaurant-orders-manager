<script lang="ts">
	import { auth } from '$lib/stores/auth.svelte';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { API_BASE } from '$lib/config';

	if (auth.role !== 'manager') goto('/login');

	type Order = { id: number; status: string; created_at: string; courses?: Course[] };
	type Course = { name: string; items: Item[] };
	type Item = { dish_name: string; quantity: number; price_cents: number };

	let orders = $state<Order[]>([]);
	let detail = $state<Order | null>(null);
	let loading = $state(true);
	let detailLoading = $state(false);
	let selectedId = $state<number | null>(null);
	let customerEmail = $state('');
	let invoiceResult = $state<string | null>(null);
	let payResult = $state<string | null>(null);

	const token = () => localStorage.getItem('token') ?? '';

	onMount(async () => {
		const r = await fetch(`${API_BASE}/orders`, { headers: { Authorization: `Bearer ${token()}` } });
		if (r.ok) orders = (await r.json()).filter((o: Order) => o.status === 'completed' || o.status === 'paid');
		loading = false;
	});

	async function selectOrder(o: Order) {
		selectedId = o.id;
		detailLoading = true;
		customerEmail = ''; invoiceResult = null; payResult = null;
		const r = await fetch(`${API_BASE}/orders/${o.id}`, { headers: { Authorization: `Bearer ${token()}` } });
		if (r.ok) detail = await r.json();
		detailLoading = false;
	}

	function courseTotal(items: Item[]) {
		return items.reduce((s, i) => s + i.price_cents * i.quantity, 0);
	}

	function orderTotal() {
		if (!detail) return 0;
		return detail.courses?.reduce((s, c) => s + courseTotal(c.items), 0) ?? 0;
	}

	function price(c: number) { return `€${(c / 100).toFixed(2)}`; }

	async function sendInvoice() {
		if (!selectedId || !customerEmail) return;
		const r = await fetch(`${API_BASE}/orders/${selectedId}/send-invoice`, { method: 'POST', headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token()}` }, body: JSON.stringify({ email: customerEmail }) });
		invoiceResult = r.ok ? 'Invoice sent!' : 'Failed';
	}

	async function markPaid() {
		if (!selectedId) return;
		const r = await fetch(`${API_BASE}/orders/${selectedId}/pay`, { method: 'POST', headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token()}` }, body: JSON.stringify({ payment_method: 'cash' }) });
		if (r.ok) { payResult = 'Marked as paid'; if (detail) detail.status = 'paid'; }
		else payResult = 'Failed';
	}
</script>

<div class="space-y-4">
	<h2 class="text-2xl font-bold">Invoicing</h2>

	{#if loading}
		<div class="flex justify-center py-8"><span class="loading loading-spinner loading-lg"></span></div>
	{:else if orders.length === 0}
		<div class="flex h-32 items-center justify-center rounded-box border-2 border-dashed text-base-content/40">No completed orders</div>
	{:else}
		<div class="overflow-x-auto">
			<table class="table table-zebra">
				<thead><tr><th>Order</th><th>Status</th><th>Date</th><th>Items</th><th>Total</th></tr></thead>
				<tbody>
					{#each orders as o (o.id)}
						<tr class="cursor-pointer" class:bg-base-200={selectedId === o.id} onclick={() => selectOrder(o)}>
							<td class="font-medium">#{o.id}</td>
							<td><span class="badge" class:badge-success={o.status === 'paid'}>{o.status}</span></td>
							<td class="text-sm text-base-content/50">{new Date(o.created_at).toLocaleDateString()}</td>
							<td class="text-sm">{o.courses?.reduce((s, c) => s + c.items.length, 0) ?? '—'}</td>
							<td class="font-semibold">{price(orderTotal())}</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>

		{#if detailLoading}
			<div class="flex justify-center py-4"><span class="loading loading-spinner loading-md"></span></div>
		{:else if detail}
			<div class="card bg-base-100 shadow-xl">
				<div class="card-body">
					<h3 class="card-title text-lg">Order #{detail.id}</h3>
					<p class="text-sm text-base-content/50">{new Date(detail.created_at).toLocaleDateString()}</p>

					<!-- Order preview with prices -->
					<div class="divider text-sm font-semibold">Items</div>
					<div class="space-y-4">
						{#each detail.courses ?? [] as course}
							<div>
								<h4 class="font-semibold text-sm mb-1">{course.name}</h4>
								<div class="space-y-1">
									{#each course.items as item}
										<div class="flex justify-between text-sm">
											<div class="flex items-center gap-2">
												<span class="text-base-content/50">×{item.quantity}</span>
												<span>{item.dish_name}</span>
											</div>
											<span class="tabular-nums">{price(item.price_cents * item.quantity)}</span>
										</div>
									{/each}
								</div>
								<div class="flex justify-between text-sm font-medium border-t border-base-300 pt-1 mt-1">
									<span>{course.name} subtotal</span>
									<span class="tabular-nums">{price(courseTotal(course.items))}</span>
								</div>
							</div>
						{/each}
					</div>

					<div class="divider"></div>
					<div class="flex justify-between text-xl font-bold">
						<span>Total</span>
						<span class="tabular-nums">{price(orderTotal())}</span>
					</div>

					<div class="divider"></div>

					{#if detail.status === 'completed'}
						<div class="flex gap-2 items-end">
							<label class="form-control flex-1">
								<div class="label"><span class="label-text">Customer Email</span></div>
								<input type="email" bind:value={customerEmail} class="input input-bordered" />
							</label>
							<button class="btn btn-primary" onclick={sendInvoice} disabled={!customerEmail}>Send Invoice</button>
						</div>
						{#if invoiceResult}<p class="text-sm text-success">{invoiceResult}</p>{/if}
						<div class="border-t pt-3">
							<button class="btn btn-success" onclick={markPaid}>Mark as Paid</button>
							{#if payResult}<p class="text-sm text-success mt-1">{payResult}</p>{/if}
						</div>
					{:else}
						<div class="badge badge-success badge-lg">Paid</div>
						<button class="btn btn-primary btn-sm mt-2" onclick={sendInvoice} disabled={!customerEmail}>Resend Invoice</button>
					{/if}
				</div>
			</div>
		{/if}
	{/if}
</div>
