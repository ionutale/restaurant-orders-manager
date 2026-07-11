<script lang="ts">
	import { auth } from '$lib/stores/auth.svelte';
	import { goto } from '$app/navigation';

	if (auth.role !== 'manager') goto('/login');
</script>

<div class="space-y-4">
	<h2 class="text-2xl font-bold">Admin Dashboard</h2>
	<div class="grid grid-cols-1 gap-4 md:grid-cols-3">
		<div class="card bg-base-100 shadow-xl">
			<div class="card-body">
				<h3 class="card-title">Floor Plan</h3>
				<p>Manage restaurant tables layout</p>
				<div class="card-actions"><a href="/admin/tables" class="btn btn-primary">Manage Tables</a></div>
			</div>
		</div>
		<div class="card bg-base-100 shadow-xl">
			<div class="card-body">
				<h3 class="card-title">Menu</h3>
				<p>Manage dishes, categories, allergens</p>
				<div class="card-actions"><a href="/admin/menu" class="btn btn-primary">Manage Menu</a></div>
			</div>
		</div>
		<div class="card bg-base-100 shadow-xl">
			<div class="card-body">
				<h3 class="card-title">Reports</h3>
				<p>View audit logs and activity</p>
				<div class="card-actions"><a href="/admin/audit" class="btn btn-primary">View Logs</a></div>
			</div>
		</div>
	</div>
</div>
