<script lang="ts">
	import { auth } from '$lib/stores/auth.svelte';
	import { goto } from '$app/navigation';

	if (auth.role !== 'waiter') goto('/login');
</script>

<div class="space-y-4">
	<h2 class="text-2xl font-bold">Waiter Dashboard</h2>
	<div class="grid grid-cols-1 gap-4 md:grid-cols-2">
		<div class="card bg-base-100 shadow-xl">
			<div class="card-body">
				<h3 class="card-title">Floor Plan</h3>
				<p>View table availability and seat guests</p>
				<div class="card-actions"><a href="/waiter" class="btn btn-primary">View Floor Plan</a></div>
			</div>
		</div>
		<div class="card bg-base-100 shadow-xl">
			<div class="card-body">
				<h3 class="card-title">Menu</h3>
				<p>Browse dishes and today's chef suggestions</p>
				<div class="card-actions"><a href="/waiter/menu" class="btn btn-primary">Browse Menu</a></div>
			</div>
		</div>
		<div class="card bg-base-100 shadow-xl">
			<div class="card-body">
				<h3 class="card-title">Active Orders</h3>
				<p>Manage your current orders</p>
				<div class="card-actions"><a href="/waiter/orders" class="btn btn-primary">View Orders</a></div>
			</div>
		</div>
	</div>
</div>
