<script lang="ts">
	import { auth } from '$lib/stores/auth.svelte';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { API_BASE } from '$lib/config';

	if (auth.role !== 'manager') goto('/login');

	type Order = { id: number; table_group_id: number; status: string; created_at: string };

	let orders = $state<Order[]>([]);
	let loading = $state(true);
	let selectedOrder = $state<Order | null>(null);
	let customerEmail = $state('');
	let invoiceResult = $state<string | null>(null);
	let payResult = $state<string | null>(null);

	const token = () => localStorage.getItem('token') ?? '';

	onMount(async () => {
		const r = await fetch(`${API_BASE}/orders`, { headers: { Authorization: `Bearer ${token()}` } });
		if (r.ok) orders = (await r.json()).filter((o: Order) => o.status === 'completed' || o.status === 'paid');
		loading = false;
	});

	async function sendInvoice() {
		if (!selectedOrder || !customerEmail) return;
		const r = await fetch(`${API_BASE}/orders/${selectedOrder.id}/send-invoice`, { method: 'POST', headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token()}` }, body: JSON.stringify({ email: customerEmail }) });
		if (r.ok) invoiceResult = 'Invoice sent!';
		else invoiceResult = 'Failed to send';
	}

	async function markPaid() {
		if (!selectedOrder) return;
		const r = await fetch(`${API_BASE}/orders/${selectedOrder.id}/pay`, { method: 'POST', headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token()}` }, body: JSON.stringify({ payment_method: 'cash' }) });
		if (r.ok) { payResult = 'Marked as paid!'; selectedOrder.status = 'paid'; }
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
				<thead><tr><th>Order</th><th>Status</th><th>Date</th><th></th></tr></thead>
				<tbody>
					{#each orders as o (o.id)}
						<tr class="cursor-pointer" class:bg-base-200={selectedOrder?.id === o.id} onclick={() => { selectedOrder = o; customerEmail = ''; invoiceResult = null; payResult = null; }}>
							<td class="font-medium">#{o.id}</td>
							<td><span class="badge" class:badge-success={o.status === 'paid'}>{o.status}</span></td>
							<td class="text-sm text-base-content/50">{new Date(o.created_at).toLocaleDateString()}</td>
							<td><button class="btn btn-ghost btn-xs" onclick={() => goto(`/waiter/orders/${o.id}`)}>View</button></td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>

		{#if selectedOrder}
			<div class="card bg-base-100 shadow-xl">
				<div class="card-body space-y-3">
					<h3 class="card-title">Order #{selectedOrder.id}</h3>

					{#if selectedOrder.status === 'completed'}
						<div>
							<h4 class="font-semibold text-sm">Send Invoice</h4>
							<div class="flex gap-2 items-end mt-1">
								<label class="form-control flex-1"><div class="label"><span class="label-text">Customer Email</span></div><input type="email" bind:value={customerEmail} class="input input-bordered input-sm" /></label>
								<button class="btn btn-primary btn-sm" onclick={sendInvoice} disabled={!customerEmail}>Send Invoice</button>
							</div>
							{#if invoiceResult}<p class="text-sm text-success mt-1">{invoiceResult}</p>{/if}
						</div>

						<div class="border-t pt-3">
							<h4 class="font-semibold text-sm">Payment</h4>
							<p class="text-sm text-base-content/60 mb-2">Customer pays manager, then mark as paid</p>
							<button class="btn btn-success btn-sm" onclick={markPaid}>Mark as Paid</button>
							{#if payResult}<p class="text-sm text-success mt-1">{payResult}</p>{/if}
						</div>
					{:else}
						<div class="badge badge-success">Already paid</div>
						<button class="btn btn-primary btn-sm" onclick={sendInvoice}>Resend Invoice</button>
					{/if}
				</div>
			</div>
		{/if}
	{/if}
</div>
