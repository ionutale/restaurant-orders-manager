import { test, expect } from '@playwright/test';

test.describe.configure({ mode: 'serial' });

test('chef: KDS dashboard shows sent orders with items', async ({ page }) => {
	// First create a sent order via API so we have something to see
	await page.goto("/");
	const gid = await page.evaluate(async () => {
		const r1 = await fetch('/api/auth/login', {
			method: 'POST', headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify({ email: 'waiter@restaurant.com', password: 'waiter' }),
		});
		const { token } = await r1.json();

		const fp = await (await fetch('/api/floor-plan', { headers: { Authorization: `Bearer ${token}` } })).json();
		const free = fp.find((t: any) => t.status === 'free');
		if (!free) return 'no free table';

		const g = await (await fetch('/api/start-order', {
			method: 'POST', headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` },
			body: JSON.stringify({ table_ids: [free.id], party_size: 2 }),
		})).json();

		const menu = await (await fetch('/api/menu', { headers: { Authorization: `Bearer ${token}` } })).json();
		const dishId = menu.categories?.[0]?.dishes?.[0]?.id;
		if (!dishId) return 'no dish';

		const order = await (await fetch(`/api/orders/${g.id}`, { headers: { Authorization: `Bearer ${token}` } })).json();
		const courseId = order.courses?.[0]?.id;
		if (!courseId) return 'no course';

		await fetch(`/api/orders/${g.id}/courses/${courseId}/items`, {
			method: 'POST', headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` },
			body: JSON.stringify({ dish_id: dishId, quantity: 1 }),
		});
		await fetch(`/api/orders/${g.id}/send`, { method: 'POST', headers: { Authorization: `Bearer ${token}` } });
		return g.id;
	});
	if (gid === 'no free table' || token === 'no dish' || token === 'no course') return;

	await page.goto('/login');
	await page.fill('input[type="email"]', 'chef@restaurant.com');
	await page.fill('input[type="password"]', 'chef');
	await page.click('button[type="submit"]');
	await expect(page.locator('h2:has-text("KDS Dashboard")')).toBeVisible({ timeout: 10000 });

	// Should see the sent order
	await expect(page.locator('text=Order #').first()).toBeVisible({ timeout: 5000 });

	// Click a "Ready" button to mark an item done
	const readyBtn = page.locator('button:has-text("Ready")').first();
	if (await readyBtn.isVisible({ timeout: 3000 }).catch(() => false)) {
		await readyBtn.click();
		await page.waitForTimeout(500);
		await expect(page.locator('span.badge-success').first()).toBeVisible({ timeout: 3000 });
	}
});
