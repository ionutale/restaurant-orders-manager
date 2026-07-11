import { test, expect } from '@playwright/test';

test.describe.configure({ mode: 'serial' });

test('waiter: floor plan loads with canvas and list toggle', async ({ page }) => {
	await page.goto('/login');
	await page.fill('input[type="email"]', 'waiter@restaurant.com');
	await page.fill('input[type="password"]', 'waiter');
	await page.click('button[type="submit"]');
	await page.waitForURL(/\/waiter/);

	await expect(page.locator('h2:has-text("Floor Plan")')).toBeVisible();
	await expect(page.locator('.join')).toBeVisible();

	// Click List toggle
	await page.locator('.join button:has-text("List")').click();
	await page.waitForTimeout(500);
	// Click back to Canvas
	await page.locator('.join button:has-text("Canvas")').click();
});

test('waiter: browse menu with dish details', async ({ page }) => {
	await page.goto('/login');
	await page.fill('input[type="email"]', 'waiter@restaurant.com');
	await page.fill('input[type="password"]', 'waiter');
	await page.click('button[type="submit"]');
	await page.waitForURL(/\/waiter/);
	await page.goto('/waiter/menu');
	await expect(page.locator('h2:has-text("Menu")')).toBeVisible();
	await page.locator('.join button:has-text("Mains")').click();
	await page.waitForTimeout(300);
	const dish = page.locator('button.card').first();
	if (await dish.isVisible()) {
		await dish.click();
		await page.waitForTimeout(300);
		await expect(page.locator('.modal-box')).toBeVisible();
		await page.locator('.modal-box button:has-text("Close")').click();
	}
});

test('waiter: orders page shows New Order button', async ({ page }) => {
	await page.goto('/login');
	await page.fill('input[type="email"]', 'waiter@restaurant.com');
	await page.fill('input[type="password"]', 'waiter');
	await page.click('button[type="submit"]');
	await page.waitForURL(/\/waiter/);
	await page.goto('/waiter/orders');
	await expect(page.locator('h2:has-text("Orders")')).toBeVisible();

	await page.locator('button:has-text("New Order")').click();
	await page.waitForTimeout(300);
	await expect(page.locator('.modal-open')).toBeVisible();
	await page.locator('.modal-open button:has-text("Cancel")').click();
});

test('waiter: seat a free table', async ({ page }) => {
	await page.goto('/login');
	await page.fill('input[type="email"]', 'waiter@restaurant.com');
	await page.fill('input[type="password"]', 'waiter');
	await page.click('button[type="submit"]');
	await page.waitForURL(/\/waiter/);
	await page.locator('.join button:has-text("List")').click();
	await page.waitForTimeout(500);
	const freeRow = page.locator('tr').filter({ has: page.locator('span.badge-success') }).first();
	if (await freeRow.isVisible()) {
		await freeRow.click();
		await page.waitForTimeout(300);
		await expect(page.locator('.modal-open')).toBeVisible();
		await page.locator('.modal-open input[type="number"]').fill('3');
		await page.locator('.modal-open button:has-text("Seat")').click();
		await page.waitForTimeout(500);
	}
});

test('waiter: occupied table shows info', async ({ page }) => {
	await page.goto('/login');
	await page.fill('input[type="email"]', 'waiter@restaurant.com');
	await page.fill('input[type="password"]', 'waiter');
	await page.click('button[type="submit"]');
	await page.waitForURL(/\/waiter/);
	await page.locator('.join button:has-text("List")').click();
	await page.waitForTimeout(500);
	const occRow = page.locator('tr').filter({ has: page.locator('span.badge-error') }).first();
	if (await occRow.isVisible()) {
		await occRow.click();
		await page.waitForTimeout(300);
		const close = page.locator('.modal-open button:has-text("Close")');
		if (await close.isVisible()) await close.click();
	}
});
