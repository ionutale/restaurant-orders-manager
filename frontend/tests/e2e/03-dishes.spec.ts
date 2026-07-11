import { test, expect } from '@playwright/test';

test('admin CRUD: create, update, and delete dishes with allergens', async ({ page }) => {
	await page.goto('/login');
	await page.fill('input[type="email"]', 'admin@restaurant.com');
	await page.fill('input[type="password"]', 'admin');
	await page.click('button[type="submit"]');
	await expect(page).toHaveURL(/\/admin/);

	await page.goto('/admin/menu');
	await page.click('text=Dishes');
	await expect(page.locator('text=Add Dish')).toBeVisible();

	// CREATE
	await page.click('button:has-text("Add Dish")');
	const modal = page.locator('.modal-open');
	await modal.locator('input[placeholder]').first().fill('TestDish');
	await modal.locator('textarea').fill('A test dish description');
	await modal.locator('input[type="number"]').first().fill('1250');
	await modal.locator('button:has-text("Add")').click();
	await page.waitForTimeout(500);
	await expect(page.locator('text=TestDish').first()).toBeVisible();

	// UPDATE — click Edit on the test dish
	const editBtn = page.locator('button:has-text("Edit")').last();
	await editBtn.click();
	await page.waitForTimeout(300);
	const nameInput = page.locator('input[type="text"]').first();
	await nameInput.fill('TestDishUpdated');
	await page.locator('button:has-text("Save")').click();
	await page.waitForTimeout(500);
	await expect(page.locator('text=TestDishUpdated').first()).toBeVisible();

	// ALLERGENS — click Details and toggle allergen
	await page.locator('button:has-text("Details")').last().click();
	await page.waitForTimeout(300);
	const allergenBtn = page.locator('.modal-open button:has-text("Gluten")');
	await allergenBtn.click();
	await page.waitForTimeout(300);
	await page.locator('.modal-open button:has-text("Close")').click();

	// DELETE
	await page.locator('button:has-text("Delete")').last().click();
	await page.locator('button:has-text("Confirm")').click();
	await page.waitForTimeout(500);
	await expect(page.locator('text=TestDishUpdated')).toHaveCount(0);
});
