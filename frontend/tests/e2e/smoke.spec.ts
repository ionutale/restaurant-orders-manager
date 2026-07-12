import { test, expect } from '@playwright/test';

// Smoke test suite — runs in < 30 seconds, verifies critical paths
test.describe.configure({ mode: 'parallel' });

test('health endpoint returns ok', async ({ page }) => {
	await page.goto('/');
	const status = await page.evaluate(async () => {
		const r = await fetch('/api/health');
		if (!r.ok) return r.status;
		const d = await r.json();
		return d.status === 'ok' ? 200 : 500;
	});
	expect(status).toBe(200);
});

test('admin can log in', async ({ page }) => {
	await page.goto('/login');
	await page.fill('input[type="email"]', 'admin@restaurant.com');
	await page.fill('input[type="password"]', 'admin');
	await page.click('button[type="submit"]');
	await page.waitForTimeout(2000);
	await expect(page.locator('h2').first()).toBeVisible({ timeout: 5000 });
});

test('waiter can log in and see floor plan', async ({ page }) => {
	await page.goto('/login');
	await page.fill('input[type="email"]', 'waiter@restaurant.com');
	await page.fill('input[type="password"]', 'waiter');
	await page.click('button[type="submit"]');
	await page.waitForTimeout(3000);
	await expect(page.locator('h2:has-text("Floor Plan")')).toBeVisible({ timeout: 10000 });
});

test('chef can log in and see KDS', async ({ page }) => {
	await page.goto('/login');
	await page.fill('input[type="email"]', 'chef@restaurant.com');
	await page.fill('input[type="password"]', 'chef');
	await page.click('button[type="submit"]');
	await expect(page.locator('h2:has-text("KDS Dashboard")')).toBeVisible({ timeout: 8000 });
});

test('menu API returns categories', async ({ page }) => {
	await page.goto('/');
	const cats = await page.evaluate(async () => {
		const r1 = await fetch('/api/auth/login', {
			method: 'POST', headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify({ email: 'waiter@restaurant.com', password: 'waiter' }),
		});
		const { token } = await r1.json();
		const r2 = await fetch('/api/menu', { headers: { Authorization: `Bearer ${token}` } });
		if (!r2.ok) return [];
		const d = await r2.json();
		return d.categories || [];
	});
	expect(cats.length).toBeGreaterThanOrEqual(1);
});
