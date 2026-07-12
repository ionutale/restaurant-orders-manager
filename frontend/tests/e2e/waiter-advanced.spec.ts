import { test, expect } from '@playwright/test';

test.describe.configure({ mode: 'serial' });

test('waiter: merge a free table into an occupied group', async ({ page }) => {
	await page.goto('/login');
	await page.fill('input[type="email"]', 'waiter@restaurant.com');
	await page.fill('input[type="password"]', 'waiter');
	await page.click('button[type="submit"]');
	await page.locator('.join button:has-text("List")').click();
	await page.waitForTimeout(500);

	const occRow = page.locator('tr').filter({ has: page.locator('span.badge-error') }).first();
	const freeRow = page.locator('tr').filter({ has: page.locator('span.badge-success') }).first();
	if (await occRow.isVisible({ timeout: 3000 }).catch(() => false) && await freeRow.isVisible().catch(() => false)) {
		await occRow.click();
		await page.waitForTimeout(300);
		await page.locator('button:has-text("Merge")').click();
		await page.waitForTimeout(300);
		const freeName = await freeRow.locator('td').first().textContent();
		await page.locator(`button:has-text("${freeName?.trim()}")`).click();
		await page.waitForTimeout(200);
		await page.locator('button:has-text("Merge")').last().click();
		await page.waitForTimeout(500);
	}
});

test('waiter: close an occupied table group', async ({ page }) => {
	await page.goto('/login');
	await page.fill('input[type="email"]', 'waiter@restaurant.com');
	await page.fill('input[type="password"]', 'waiter');
	await page.click('button[type="submit"]');
	await page.locator('.join button:has-text("List")').click();
	await page.waitForTimeout(500);

	const occRow = page.locator('tr').filter({ has: page.locator('span.badge-error') }).first();
	if (await occRow.isVisible({ timeout: 3000 }).catch(() => false)) {
		await occRow.click();
		await page.waitForTimeout(300);
		const closeBtn = page.locator('button:has-text("Close")');
		if (await closeBtn.isVisible().catch(() => false)) {
			await closeBtn.click();
			await page.waitForTimeout(500);
		}
	}
});

test('waiter: search dishes by name in order detail', async ({ page }) => {
	// Login as waiter, go to an existing order
	await page.goto('/login');
	await page.fill('input[type="email"]', 'waiter@restaurant.com');
	await page.fill('input[type="password"]', 'waiter');
	await page.click('button[type="submit"]');
	await page.goto('/waiter/orders');
	const orderRow = page.locator('table tbody tr').first();
	if (await orderRow.isVisible({ timeout: 3000 }).catch(() => false)) {
		await orderRow.click();
		await page.waitForTimeout(500);
		await expect(page.locator('h2:has-text("Order #")')).toBeVisible({ timeout: 5000 });

		const searchInput = page.locator('input[placeholder="Search dishes..."]');
		if (await searchInput.isVisible().catch(() => false)) {
			await searchInput.fill('Bruschetta');
			await page.waitForTimeout(300);
			await expect(page.locator('button:has-text("Bruschetta")').first()).toBeVisible({ timeout: 3000 });
			await searchInput.fill('zzzznotexist');
			await page.waitForTimeout(300);
			await expect(page.locator('text=No dishes match').first()).toBeVisible({ timeout: 3000 });
		}
	}
});
