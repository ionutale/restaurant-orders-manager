import { test, expect } from '@playwright/test';

test('admin CRUD: create, rename, reorder, delete categories', async ({ page }) => {
	await page.goto('/login');
	await page.fill('input[type="email"]', 'admin@restaurant.com');
	await page.fill('input[type="password"]', 'admin');
	await page.click('button[type="submit"]');
	await expect(page).toHaveURL(/\/admin/);

	await page.goto('/admin/menu');
	await expect(page.locator('text=Categories')).toBeVisible();

	// CREATE
	await page.click('button:has-text("Add Category")');
	await page.fill('input[placeholder="Appetizers"]', 'TestCat');
	await page.locator('.modal-box button:has-text("Add")').click();
	await page.waitForTimeout(500);
	await expect(page.locator('text=TestCat').first()).toBeVisible();

	// READ — verify existing categories
	await expect(page.locator('text=Appetizers').first()).toBeVisible();

	// UPDATE (inline rename)
	const editBtns = page.locator('button:has-text("Edit")');
	const lastEdit = editBtns.last();
	await lastEdit.click();
	await page.waitForTimeout(300);
	const nameInput = page.locator('table input[type="text"]').first();
	await nameInput.fill('TestCatRenamed');
	await page.locator('button:has-text("Save")').click();
	await page.waitForTimeout(500);
	await expect(page.locator('text=TestCatRenamed').first()).toBeVisible();

	// REORDER — click up arrow on the last category
	await page.locator('button:has-text("↑")').first().click();
	await page.waitForTimeout(300);

	// DELETE
	await page.locator('button:has-text("Delete")').last().click();
	await page.locator('button:has-text("Confirm")').click();
	await page.waitForTimeout(500);
	await expect(page.locator('text=TestCatRenamed')).toHaveCount(0);
});
