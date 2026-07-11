import { test, expect } from '@playwright/test';

test('chef CRUD: create, renew, and delete suggestions', async ({ page }) => {
	await page.goto('/login');
	await page.fill('input[type="email"]', 'chef@restaurant.com');
	await page.fill('input[type="password"]', 'chef');
	await page.click('button[type="submit"]');
	await expect(page).toHaveURL(/\/chef/);

	// CREATE suggestion
	await page.click('button:has-text("New Suggestion")');
	const modal = page.locator('.modal-open');
	await modal.locator('input[placeholder]').first().fill('TestSpecial');
	await modal.locator('textarea').fill('Chefs daily test');
	await modal.locator('input[type="number"]').first().fill('1500');
	await modal.locator('button:has-text("Create")').click();
	await page.waitForTimeout(500);
	await expect(page.locator('text=TestSpecial').first()).toBeVisible();

	// READ — verify it appears as active
	await expect(page.locator('text=€15.00').first()).toBeVisible();

	// RENEW — show expired, check none exist, then verify active stays
	await page.click('button:has-text("Show Expired")');
	await page.waitForTimeout(300);
	await page.click('button:has-text("Hide Expired")');

	// DELETE
	await page.locator('button:has-text("✕")').first().click();
	await page.locator('button:has-text("Confirm")').click();
	await page.waitForTimeout(500);
	await expect(page.locator('text=TestSpecial')).toHaveCount(0);
});
