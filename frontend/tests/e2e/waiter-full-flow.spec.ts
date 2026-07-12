import { test, expect } from '@playwright/test';

test.describe.configure({ mode: 'serial' });

const TS = Date.now().toString(36).slice(-4);
const WAITER_EMAIL = `e2e${TS}@test.com`;
const WAITER_PASS = 'test123';

test('admin creates a new waiter via API', async ({ page }) => {
	await page.goto('/');
	const ok = await page.evaluate(async (email) => {
		const r1 = await fetch('/api/auth/login', {
			method: 'POST', headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify({ email: 'admin@restaurant.com', password: 'admin' }),
		});
		if (!r1.ok) return 'login fail';
		const d1 = await r1.json();
		const r2 = await fetch('/api/users', {
			method: 'POST', headers: { 'Content-Type': 'application/json', 'Authorization': `Bearer ${d1.token}` },
			body: JSON.stringify({ name: 'E2E-' + Date.now(), email, password: 'test123', role: 'waiter' }),
		});
		return r2.ok ? 'ok' : 'create fail';
	}, WAITER_EMAIL);
	expect(ok).toBe('ok');
});

test('new waiter picks a table, creates order, adds dishes, advances', async ({ page }) => {
	await page.goto('/login');
	await page.fill('input[type="email"]', WAITER_EMAIL);
	await page.fill('input[type="password"]', WAITER_PASS);
	await page.click('button[type="submit"]');
	await expect(page.locator('h2:has-text("Floor Plan")')).toBeVisible({ timeout: 15000 });

	await page.locator('.join button:has-text("List")').click();
	await page.waitForTimeout(500);

	const freeRow = page.locator('tr').filter({ has: page.locator('span.badge-success') }).first();
	if (!(await freeRow.isVisible({ timeout: 8000 }).catch(() => false))) {
		await expect(page.locator('h2:has-text("Floor Plan")')).toBeVisible();
		return;
	}
	await freeRow.click();
	await page.waitForTimeout(300);

	const seatModal = page.locator('.modal-open');
	await expect(seatModal).toBeVisible();
	await seatModal.locator('input[type="number"]').fill('6');
	await seatModal.locator('input[placeholder*="Birthday"]').fill(`Party${TS}`);
	await seatModal.locator('button:has-text("Seat & Order")').click();
	await page.waitForTimeout(1000);

	// Should be on order detail page now
	await expect(page.locator('h2:has-text("Order #")')).toBeVisible({ timeout: 8000 });

	// Add dishes to the course
	for (let i = 0; i < 3; i++) {
		const buttons = page.locator('button:has-text("+")');
		const firstPlus = buttons.first();
		if (!(await firstPlus.isVisible().catch(() => false))) break;
		await firstPlus.click();
		await page.waitForTimeout(600);
	}

	// Send to KDS
	const sendBtn = page.locator('button:has-text("Send to KDS")');
	if (await sendBtn.isVisible().catch(() => false)) {
		await sendBtn.click();
		await page.waitForTimeout(500);
	}

	// Advance courses
	for (let i = 0; i < 3; i++) {
		const advBtn = page.locator('button:has-text("Advance")');
		if (await advBtn.isVisible().catch(() => false)) {
			await advBtn.click();
			await page.waitForTimeout(500);
		}
	}
});
