<script lang="ts">
	import './layout.css';
	import { auth } from '$lib/stores/auth.svelte';
	import { page } from '$app/state';
	import { goto } from '$app/navigation';

	let { children } = $props();
</script>

<svelte:head>
	<title>Restaurant Orders</title>
</svelte:head>

{#if auth.loading}
	<div class="flex h-screen items-center justify-center">
		<span class="loading loading-spinner loading-lg"></span>
	</div>
{:else if !auth.isLoggedIn && $page.url.pathname !== '/login'}
	{goto('/login')}
{:else}
	{#if auth.isLoggedIn}
		<div class="drawer">
			<input id="drawer" type="checkbox" class="drawer-toggle" />
			<div class="drawer-content flex flex-col">
				<nav class="navbar bg-base-300 px-4 shadow-sm">
					<div class="flex-none lg:hidden">
						<label for="drawer" aria-label="open sidebar" class="btn btn-square btn-ghost">
							<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" class="inline-block h-6 w-6 stroke-current">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16"></path>
							</svg>
						</label>
					</div>
					<div class="flex-1">
						<a href="/" class="btn btn-ghost text-xl">Restaurant Orders</a>
					</div>
					<div class="flex-none">
						<span class="mr-4 text-sm">{auth.user?.name}</span>
						<button class="btn btn-ghost btn-sm" onclick={() => { auth.logout(); goto('/login'); }}>Logout</button>
					</div>
				</nav>
				<main class="flex-1 p-4">
					{@render children()}
				</main>
			</div>
			<div class="drawer-side z-50">
				<label for="drawer" aria-label="close sidebar" class="drawer-overlay"></label>
				<ul class="menu bg-base-200 text-base-content min-h-full w-64 p-4">
					<li class="menu-title"><span>{auth.user?.name}</span></li>
					{#if auth.role === 'manager'}
						<li><a href="/admin">Dashboard</a></li>
						<li><a href="/admin/tables">Floor Plan</a></li>
						<li><a href="/admin/menu">Menu</a></li>
					{:else if auth.role === 'waiter'}
						<li><a href="/waiter">Floor Plan</a></li>
						<li><a href="/waiter/orders">Orders</a></li>
					{:else if auth.role === 'chef'}
						<li><a href="/chef">KDS Dashboard</a></li>
					{/if}
				</ul>
			</div>
		</div>
	{:else}
		{@render children()}
	{/if}
{/if}
