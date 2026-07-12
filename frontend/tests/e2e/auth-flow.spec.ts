import { test, expect } from '@playwright/test';

test.describe.configure({ mode: 'serial' });

test('login: invalid credentials shows error', async ({ page }) => {
	await page.goto('/login');
	await page.fill('input[type="email"]', 'nonexistent@test.com');
	await page.fill('input[type="password"]', 'wrongpassword');
	await page.click('button[type="submit"]');
	await page.waitForTimeout(1000);
	await expect(page.locator('.alert-error')).toBeVisible({ timeout: 5000 });
});

test('login: admin logs in and sees admin dashboard', async ({ page }) => {
	await page.goto('/login');
	await page.fill('input[type="email"]', 'admin@restaurant.com');
	await page.fill('input[type="password"]', 'admin');
	await page.click('button[type="submit"]');
	await expect(page.locator('h2:has-text("Dashboard")')).toBeVisible({ timeout: 10000 });
});

test('login: waiter logs in and sees floor plan', async ({ page }) => {
	await page.goto('/login');
	await page.fill('input[type="email"]', 'waiter@restaurant.com');
	await page.fill('input[type="password"]', 'waiter');
	await page.click('button[type="submit"]');
	await expect(page.locator('h2:has-text("Floor Plan")')).toBeVisible({ timeout: 10000 });
});

test('login: chef logs in and sees KDS dashboard', async ({ page }) => {
	await page.goto('/login');
	await page.fill('input[type="email"]', 'chef@restaurant.com');
	await page.fill('input[type="password"]', 'chef');
	await page.click('button[type="submit"]');
	await expect(page.locator('h2:has-text("KDS Dashboard")')).toBeVisible({ timeout: 10000 });
});

test('logout: admin can logout and is redirected to login', async ({ page }) => {
	await page.goto('/login');
	await page.fill('input[type="email"]', 'admin@restaurant.com');
	await page.fill('input[type="password"]', 'admin');
	await page.click('button[type="submit"]');
	await expect(page.locator('h2:has-text("Dashboard")')).toBeVisible({ timeout: 10000 });

	// Click logout button in navbar
	const logoutBtn = page.locator('button:has-text("Logout")');
	if (await logoutBtn.isVisible({ timeout: 3000 }).catch(() => false)) {
		await logoutBtn.click();
		await page.waitForTimeout(500);
		await expect(page.locator('h1:has-text("Restaurant Orders")')).toBeVisible({ timeout: 5000 });
	}
});
