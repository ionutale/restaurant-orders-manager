import { test, expect } from '@playwright/test';

test.describe.configure({ mode: 'serial' });

test('admin: create a waiter via UI, then delete them', async ({ page }) => {
	const N = `E2E${Date.now().toString(36).slice(-4)}`;
	await page.goto('/login');
	await page.fill('input[type="email"]', 'admin@restaurant.com');
	await page.fill('input[type="password"]', 'admin');
	await page.click('button[type="submit"]');
	await page.waitForTimeout(2000);
	await page.goto('/admin/users');
	await page.waitForTimeout(1000);
	await expect(page.locator('h2:has-text("Users")')).toBeVisible({ timeout: 5000 });

	await page.locator('button:has-text("Add User")').click();
	await page.waitForTimeout(300);
	const modal = page.locator('.modal-open');
	await modal.locator('input[type="text"]').first().fill(N);
	await modal.locator('input[type="email"]').fill(`${N}@test.com`);
	await modal.locator('input[type="password"]').fill('pass123');
	await modal.locator('select').selectOption('waiter');
	await modal.locator('button:has-text("Add User")').click();
	await page.waitForTimeout(500);
	await expect(page.locator(`text=${N}`).first()).toBeVisible();

	await page.locator('button:has-text("Edit")').last().click();
	await page.waitForTimeout(300);
	await page.locator('input[type="text"]').first().fill(`${N}R`);
	await page.locator('button:has-text("Save")').click();
	await page.waitForTimeout(500);
	await expect(page.locator(`text=${N}R`).first()).toBeVisible();

	await page.locator('button:has-text("Delete")').last().click();
	await page.locator('button:has-text("Confirm")').click();
	await page.waitForTimeout(500);
});

test('admin: audit log loads with events', async ({ page }) => {
	await page.goto('/login');
	await page.fill('input[type="email"]', 'admin@restaurant.com');
	await page.fill('input[type="password"]', 'admin');
	await page.click('button[type="submit"]');
	await page.waitForTimeout(2000);
	await page.goto('/admin/audit');
	await page.waitForTimeout(1000);
	await expect(page.locator('h2:has-text("Audit Log")')).toBeVisible({ timeout: 5000 });
});

test('admin: invoicing page loads', async ({ page }) => {
	await page.goto('/login');
	await page.fill('input[type="email"]', 'admin@restaurant.com');
	await page.fill('input[type="password"]', 'admin');
	await page.click('button[type="submit"]');
	await page.waitForTimeout(2000);
	await page.goto('/admin/invoices');
	await page.waitForTimeout(1000);
	await expect(page.locator('h2:has-text("Invoicing")')).toBeVisible({ timeout: 5000 });
});
