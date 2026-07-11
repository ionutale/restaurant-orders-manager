<script lang="ts">
	import { auth } from '$lib/stores/auth.svelte';
	import { goto } from '$app/navigation';

	if (auth.isLoggedIn) {
		if (auth.role === 'manager') goto('/admin');
		else if (auth.role === 'waiter') goto('/waiter');
		else if (auth.role === 'chef') goto('/chef');
	}
</script>
