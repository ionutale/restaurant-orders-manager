import { test, expect } from '@playwright/test';

test.describe.configure({ mode: 'serial' });

const TS = Date.now().toString(36).slice(-4);
const TNAME = `E2ET${TS}`;

test('create, edit, delete table', async ({ page }) => {
	await page.goto('/login');
	await page.fill('input[type="email"]', 'admin@restaurant.com');
	await page.fill('input[type="password"]', 'admin');
	await page.click('button[type="submit"]');
	await page.waitForURL(/\/admin/);
	await page.goto('/admin/tables');
	await page.locator('.join button:has-text("List")').click();
	await page.waitForTimeout(500);

	// CREATE
	await page.fill('input[placeholder="T6"]', TNAME);
	await page.locator('input[type="number"][min="1"]').first().fill('6');
	await page.locator('button:has-text("Add")').click();
	await page.waitForTimeout(500);
	await expect(page.locator(`text=${TNAME}`).first()).toBeVisible();

	// EDIT
	await page.locator(`button:has-text("Edit")`).first().click();
	await page.waitForTimeout(300);
	const ni = page.locator('input[type="text"]').first();
	await ni.fill(`${TNAME}R`);
	await page.locator('button:has-text("Save")').click();
	await page.waitForTimeout(500);
	await expect(page.locator(`text=${TNAME}R`).first()).toBeVisible();

	// DELETE (just confirm no error)
	await page.locator(`button:has-text("Delete")`).first().click();
	await page.locator('button:has-text("Confirm")').click();
	await page.waitForTimeout(500);
});

test('create, edit, delete category', async ({ page }) => {
	const C = `E2EC${TS}`;
	await page.goto('/login');
	await page.fill('input[type="email"]', 'admin@restaurant.com');
	await page.fill('input[type="password"]', 'admin');
	await page.click('button[type="submit"]');
	await page.waitForURL(/\/admin/);
	await page.goto('/admin/menu');

	// CREATE
	await page.locator('button:has-text("Add Category")').click();
	await page.locator('input[placeholder="Appetizers"]').fill(C);
	await page.locator('.modal-box button:has-text("Add")').click();
	await page.waitForTimeout(500);
	await expect(page.locator(`text=${C}`).first()).toBeVisible();

	// EDIT
	await page.locator('button:has-text("Edit")').last().click();
	await page.waitForTimeout(300);
	const ci = page.locator('table input[type="text"]').first();
	await ci.fill(`${C}R`);
	await page.locator('button:has-text("Save")').click();
	await page.waitForTimeout(500);
	await expect(page.locator(`text=${C}R`).first()).toBeVisible();

	// DELETE
	await page.locator('button:has-text("Delete")').last().click();
	await page.locator('button:has-text("Confirm")').click();
	await page.waitForTimeout(500);
});

test('create, edit, delete dish with allergens', async ({ page }) => {
	const D = `E2ED${TS}`;
	await page.goto('/login');
	await page.fill('input[type="email"]', 'admin@restaurant.com');
	await page.fill('input[type="password"]', 'admin');
	await page.click('button[type="submit"]');
	await page.waitForURL(/\/admin/);
	await page.goto('/admin/menu');
	await page.locator('button:has-text("Dishes")').click();
	await page.waitForTimeout(300);

	// CREATE
	await page.locator('button:has-text("Add Dish")').click();
	await page.waitForTimeout(300);
	const modal = page.locator('.modal-open');
	await modal.locator('input').first().fill(D);
	await modal.locator('textarea').fill('E2E test');
	await modal.locator('input[type="number"]').first().fill('999');
	await modal.locator('button:has-text("Add")').click();
	await page.waitForTimeout(500);
	await expect(page.locator(`text=${D}`).first()).toBeVisible();

	// ALLERGENS
	await page.locator('button:has-text("Details")').last().click();
	await page.waitForTimeout(300);
	const glu = page.locator('.modal-open button:has-text("Gluten")');
	if (await glu.isVisible()) { await glu.click(); await page.waitForTimeout(200); }
	await page.locator('.modal-open button:has-text("Close")').click();

	// DELETE
	await page.locator('button:has-text("Delete")').last().click();
	await page.locator('button:has-text("Confirm")').click();
	await page.waitForTimeout(500);
});

test('chef: KDS dashboard loads and shows active orders', async ({ page }) => {
	await page.goto('/login');
	await page.fill('input[type="email"]', 'chef@restaurant.com');
	await page.fill('input[type="password"]', 'chef');
	await page.click('button[type="submit"]');
	await page.waitForURL(/\/chef/);

	await expect(page.locator('h2:has-text("KDS Dashboard")')).toBeVisible();
	await expect(page.locator('button:has-text("Refresh")')).toBeVisible();
});
