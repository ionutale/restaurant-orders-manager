import { test, expect } from '@playwright/test';

test('admin CRUD: create, edit, and delete tables', async ({ page }) => {
	await page.goto('/login');
	await page.fill('input[type="email"]', 'admin@restaurant.com');
	await page.fill('input[type="password"]', 'admin');
	await page.click('button[type="submit"]');
	await expect(page).toHaveURL(/\/admin/);

	await page.goto('/admin/tables');
	await page.click('text=List');
	await expect(page.locator('text=Add Table')).toBeVisible();

	// CREATE
	await page.fill('input[placeholder="T6"]', 'ZA');
	await page.fill('input[type="number"][min="1"]', '6');
	await page.click('button:has-text("Add")');
	await expect(page.locator('text=ZA').first()).toBeVisible();

	// READ — verify table appears
	const tableRows = page.locator('table.table-zebra tbody tr');
	await expect(tableRows.first()).toBeVisible();

	// UPDATE (inline edit)
	const editBtn = page.locator('button:has-text("Edit")').first();
	await editBtn.click();
	await page.waitForTimeout(300);
	const nameInput = page.locator('input[type="text"]').first();
	await nameInput.fill('ZB');
	await page.locator('button:has-text("Save")').click();
	await page.waitForTimeout(500);
	await expect(page.locator('text=ZB').first()).toBeVisible();

	// DELETE
	await page.locator('button:has-text("Delete")').first().click();
	await page.locator('button:has-text("Confirm")').click();
	await page.waitForTimeout(500);
	await expect(page.locator('text=ZB')).toHaveCount(0);
});
