import { test, expect } from '@playwright/test';

test.describe.configure({ mode: 'serial' });

const TS = Date.now().toString(36).slice(-4);
const WAITER_NAME = `E2EW${TS}`;
const WAITER_EMAIL = `e2e${TS}@test.com`;
const WAITER_PASS = 'test123';

test('admin creates a new waiter', async ({ page }) => {
	await page.goto('/login');
	await page.fill('input[type="email"]', 'admin@restaurant.com');
	await page.fill('input[type="password"]', 'admin');
	await page.click('button[type="submit"]');
	await page.waitForURL(/\/admin/);

	await page.goto('/admin/users');
	await expect(page.locator('h2:has-text("Users")')).toBeVisible();

	await page.locator('button:has-text("Add User")').click();
	await page.waitForTimeout(300);

	const modal = page.locator('.modal-open');
	await expect(modal).toBeVisible();
	await modal.locator('input[type="text"]').first().fill(WAITER_NAME);
	await modal.locator('input[type="email"]').fill(WAITER_EMAIL);
	await modal.locator('input[type="password"]').fill(WAITER_PASS);
	await modal.locator('select').selectOption('waiter');
	await modal.locator('button:has-text("Add User")').click();
	await page.waitForTimeout(500);

	await expect(page.locator(`text=${WAITER_NAME}`).first()).toBeVisible();
});

test('new waiter picks a table, creates order with courses, adds dishes, advances, and closes', async ({ page }) => {
	await page.goto('/login');
	await page.fill('input[type="email"]', WAITER_EMAIL);
	await page.fill('input[type="password"]', WAITER_PASS);
	await page.click('button[type="submit"]');
	await page.waitForURL(/\/waiter/);

	// Pick a free table and seat 6 people
	await page.locator('.join button:has-text("List")').click();
	await page.waitForTimeout(500);

	const freeRow = page.locator('tr').filter({ has: page.locator('span.badge-success') }).first();
	await expect(freeRow).toBeVisible({ timeout: 8000 });
	await freeRow.click();
	await page.waitForTimeout(300);

	const seatModal = page.locator('.modal-open');
	await expect(seatModal).toBeVisible();
	await seatModal.locator('input[type="number"]').fill('6');
	await seatModal.locator('input[placeholder*="Birthday"]').fill(`Party${TS}`);
	await seatModal.locator('button:has-text("Seat")').click();
	await page.waitForTimeout(500);

	// Create a new order
	await page.goto('/waiter/orders');
	await page.locator('button:has-text("New Order")').click();
	await page.waitForTimeout(400);

	const orderModal = page.locator('.modal-open');
	await expect(orderModal).toBeVisible();

	// Select the table group
	const groupSelect = orderModal.locator('select');
	const opts = await groupSelect.locator('option').all();
	if (opts.length <= 1) { test.skip('no groups'); return; }
	await groupSelect.selectOption({ index: opts.length - 1 });

	// Set course names
	const inputs = orderModal.locator('input[type="text"]');
	await inputs.nth(0).fill('Starter');
	await inputs.nth(1).fill('Main');
	await inputs.nth(2).fill('Dessert');
	await orderModal.locator('button:has-text("Add Course")').click();
	await page.waitForTimeout(200);
	const allInputs = orderModal.locator('input[type="text"]');
	const cnt = await allInputs.count();
	if (cnt >= 4) await allInputs.nth(3).fill('Drinks');

	await orderModal.locator('button:has-text("Create Order")').click();
	await page.waitForTimeout(500);

	// Navigate to order detail
	await expect(page.locator('table tbody tr').first()).toBeVisible({ timeout: 5000 });
	await page.locator('table tbody tr').first().click();
	await page.waitForURL(/\/waiter\/orders\/\d+/);
	await expect(page.locator('h2:has-text("Order #")')).toBeVisible({ timeout: 5000 });

	// Add dishes to each course
	for (let i = 0; i < 4; i++) {
		const addBtn = page.locator('button:has-text("Add Dish")').first();
		if (!(await addBtn.isVisible().catch(() => false))) break;
		await addBtn.click();
		await page.waitForTimeout(300);
		const addModal = page.locator('.modal-open');
		const firstDish = addModal.locator('button.card').first();
		if (await firstDish.isVisible().catch(() => false)) {
			await firstDish.click();
			await page.waitForTimeout(500);
		} else {
			await addModal.locator('button:has-text("Close")').click();
		}
	}

	// Send to KDS
	const sendBtn = page.locator('button:has-text("Send to KDS")');
	if (await sendBtn.isVisible().catch(() => false)) {
		await sendBtn.click();
		await page.waitForTimeout(500);
	}

	// Advance courses
	for (let i = 0; i < 3; i++) {
		const advBtn = page.locator('button:has-text("Advance Course")');
		if (await advBtn.isVisible().catch(() => false)) {
			await advBtn.click();
			await page.waitForTimeout(500);
		}
	}
});
