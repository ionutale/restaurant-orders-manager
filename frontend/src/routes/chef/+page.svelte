<script lang="ts">
	import { auth } from '$lib/stores/auth.svelte';
	import { goto } from '$app/navigation';

	if (auth.role !== 'chef') goto('/login');
</script>

<div class="space-y-4">
	<h2 class="text-2xl font-bold">KDS Dashboard</h2>
	<p class="text-base-content/70">Orders will appear here when the waiter sends them to the kitchen.</p>
	<div class="flex h-64 items-center justify-center rounded-box border-2 border-dashed border-base-300">
		<p class="text-base-content/50">No active orders</p>
	</div>
</div>
