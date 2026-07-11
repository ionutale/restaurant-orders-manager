import { test, expect } from '@playwright/test';

test('admin can login and create a floor plan with tables', async ({ page }) => {
	await page.goto('/login');

	// Login
	await page.fill('input[type="email"]', 'admin@restaurant.com');
	await page.fill('input[type="password"]', 'admin');
	await page.click('button[type="submit"]');

	// Should redirect to admin dashboard
	await expect(page).toHaveURL(/\/admin/);

	// Navigate to tables
	await page.goto('/admin/tables');
	await expect(page.locator('text=Floor Plan')).toBeVisible();

	// Switch to list view
	await page.click('text=List');
	await expect(page.locator('text=Add Table')).toBeVisible();

	// Create a few tables
	const tables = [
		{ name: 'T1', capacity: '4', label: 'Window' },
		{ name: 'T2', capacity: '2', label: '' },
		{ name: 'T3', capacity: '6', label: 'Corner' },
	];

	for (const tbl of tables) {
		await page.fill('input[placeholder="T6"]', tbl.name);
		await page.fill('input[type="number"][min="1"]', tbl.capacity);
		if (tbl.label) {
			await page.fill('input[placeholder="Near window"]', tbl.label);
		}
		await page.click('button:has-text("Add")');
	}

	// Verify tables appear in the list
	for (const tbl of tables) {
		await expect(page.locator(`text=${tbl.name}`).first()).toBeVisible();
	}

	// Switch to canvas view
	await page.click('text=Canvas');
	await expect(page.locator('text=T1')).toBeVisible();
	await expect(page.locator('text=T3')).toBeVisible();
});
