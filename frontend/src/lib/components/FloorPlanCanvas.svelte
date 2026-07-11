<script lang="ts">
	type Table = {
		id: number;
		name: string;
		capacity: number;
		x: number;
		y: number;
		label: string | null;
		status?: 'free' | 'occupied';
		group_name?: string;
	};

	let {
		tables,
		readonly = false,
		onUpdate,
		onAdd,
		onDelete,
		onTableClick,
	}: {
		tables: Table[];
		readonly?: boolean;
		onUpdate?: (id: number, data: Partial<Table>) => Promise<void>;
		onAdd?: () => void;
		onDelete?: (id: number) => void;
		onTableClick?: (t: Table) => void;
	} = $props();

	const GRID_SIZE = 20;

	let container: HTMLDivElement;
	let draggingId = $state<number | null>(null);
	let dragOffsetX = $state(0);
	let dragOffsetY = $state(0);
	let currentX = $state(0);
	let currentY = $state(0);

	function handlePointerDown(e: PointerEvent, t: Table) {
		if (readonly || !onUpdate) return;
		e.preventDefault();
		draggingId = t.id;
		dragOffsetX = e.clientX - t.x;
		dragOffsetY = e.clientY - t.y;
		currentX = t.x;
		currentY = t.y;
		container.setPointerCapture(e.pointerId);
	}

	function handlePointerMove(e: PointerEvent) {
		if (draggingId === null) return;
		currentX = Math.max(0, e.clientX - dragOffsetX);
		currentY = Math.max(0, e.clientY - dragOffsetY);
	}

	async function handlePointerUp(e: PointerEvent) {
		if (draggingId === null || !onUpdate) return;
		const id = draggingId;
		const newX = Math.max(0, e.clientX - dragOffsetX);
		const newY = Math.max(0, e.clientY - dragOffsetY);
		draggingId = null;
		await onUpdate(id, { x: newX, y: newY });
	}

	function rectStyle(t: Table): string {
		const x = draggingId === t.id ? currentX : t.x;
		const y = draggingId === t.id ? currentY : t.y;
		return `left: ${x}px; top: ${y}px; width: 120px; height: 80px;`;
	}

	function statusClass(t: Table): string {
		if (t.status === 'occupied') return 'table-occupied';
		if (t.status === 'free') return 'table-free';
		return 'table-default';
	}
</script>

<div class="flex flex-col gap-4">
	{#if !readonly}
		<div class="flex items-center gap-2">
			<button class="btn btn-primary btn-sm" onclick={onAdd}>
				<svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" /></svg>
				Add Table
			</button>
			<span class="text-sm text-base-content/50">Drag tables to reposition</span>
		</div>
	{/if}

	<!-- svelte-ignore a11y_no_static_element_interactions -->
	<div
		bind:this={container}
		class="relative overflow-auto rounded-box border border-base-300 bg-base-100"
		class:cursor-grab={!readonly}
		style="min-height: 500px; min-width: 600px; background-image: radial-gradient(circle, oklch(var(--bc) / 0.08) 1px, transparent 1px); background-size: {GRID_SIZE}px {GRID_SIZE}px;"
	>
		{#each tables as t (t.id)}
			<!-- svelte-ignore a11y_no_static_element_interactions -->
			<div
				class="table-rect {statusClass(t)}"
				class:cursor-grab={!readonly && draggingId !== t.id}
				class:cursor-grabbing={draggingId === t.id}
				class:cursor-pointer={readonly}
				class:opacity-60={draggingId === t.id}
				class:z-10={draggingId === t.id}
				style={rectStyle(t)}
				onpointerdown={(e) => handlePointerDown(e, t)}
				onpointermove={handlePointerMove}
				onpointerup={handlePointerUp}
				onpointercancel={() => (draggingId = null)}
				onclick={() => onTableClick?.(t)}
				role="button"
				tabindex="0"
			>
				<div class="flex h-full flex-col items-center justify-center p-1 text-center">
					<span class="text-sm font-bold leading-tight">{t.name}</span>
					<span class="text-xs text-base-content/60">{t.status === 'occupied' ? (t.group_name || 'Occupied') : `${t.capacity} seats`}</span>
					{#if t.label && t.status !== 'occupied'}
						<span class="truncate text-xs text-base-content/40">{t.label}</span>
					{/if}
				</div>
				{#if !readonly}
					<button
						class="btn btn-ghost btn-xs absolute right-0.5 top-0.5 text-error opacity-0 hover:opacity-100"
						onclick={(e) => { e.stopPropagation(); onDelete?.(t.id); }}
					>✕</button>
				{/if}
			</div>
		{/each}

		{#if tables.length === 0}
			<div class="flex h-full min-h-[300px] items-center justify-center text-base-content/30">
				No tables
			</div>
		{/if}
	</div>
</div>

<style>
	.table-rect {
		position: absolute;
		border-radius: 8px;
		border: 2px solid oklch(var(--p));
		background: oklch(var(--p) / 0.1);
		display: flex;
		align-items: center;
		justify-content: center;
		touch-action: none;
		user-select: none;
		transition: box-shadow 0.15s;
	}
	.table-rect:hover {
		box-shadow: 0 2px 8px oklch(var(--p) / 0.3);
	}
	.table-free {
		border-color: oklch(var(--su));
		background: oklch(var(--su) / 0.1);
	}
	.table-free:hover {
		box-shadow: 0 2px 8px oklch(var(--su) / 0.3);
	}
	.table-occupied {
		border-color: oklch(var(--er));
		background: oklch(var(--er) / 0.1);
	}
	.table-occupied:hover {
		box-shadow: 0 2px 8px oklch(var(--er) / 0.3);
	}
</style>
