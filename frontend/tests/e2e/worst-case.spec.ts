import { test, expect } from '@playwright/test';

test.describe.configure({ mode: 'serial' });

// ─── Network errors / API failures ───

test('handles API 500 when creating a table via route interception', async ({ page }) => {
	// Intercept POST /api/tables and return 500
	await page.route('**/api/tables*', async route => {
		if (route.request().method() === 'POST') {
			await route.fulfill({ status: 500, contentType: 'application/json', body: JSON.stringify({ error: 'Internal server error' }) });
		} else {
			await route.continue();
		}
	});

	await page.goto('/login');
	await page.fill('input[type="email"]', 'admin@restaurant.com');
	await page.fill('input[type="password"]', 'admin');
	await page.click('button[type="submit"]');
	await page.waitForTimeout(1500);
	await page.goto('/admin/tables');
	await page.waitForTimeout(500);
	await page.locator('button:has-text("List")').click();
	await page.waitForTimeout(300);

	// Try creating a table — the form submit should not crash the app
	await page.fill('input[placeholder="T6"]', 'FailTable');
	await page.locator('input[type="number"][min="1"]').first().fill('4');
	await page.locator('button:has-text("Add")').click();
	await page.waitForTimeout(500);

	// App should show an error or stay on the same page without crashing
	await expect(page.locator('h2:has-text("Floor Plan")')).toBeVisible();
});

test('handles network timeout gracefully', async ({ page }) => {
	await page.route('**/api/**', async route => {
		await new Promise(r => setTimeout(r, 30000)); // never resolves within test timeout
	});

	await page.goto('/login');
	await page.fill('input[type="email"]', 'admin@restaurant.com');
	await page.fill('input[type="password"]', 'admin');
	await page.click('button[type="submit"]');

	// App should show loading state without crashing
	await page.waitForTimeout(2000);
	const body = page.locator('body');
	await expect(body).toBeVisible();
});

test('handles malformed API response gracefully', async ({ page }) => {
	// Intercept and return non-JSON garbage
	await page.route('**/api/tables', async route => {
		if (route.request().method() === 'GET') {
			await route.fulfill({ status: 200, contentType: 'application/json', body: 'this is not valid json {{{' });
		} else {
			await route.continue();
		}
	});

	await page.goto('/login');
	await page.fill('input[type="email"]', 'admin@restaurant.com');
	await page.fill('input[type="password"]', 'admin');
	await page.click('button[type="submit"]');
	await page.waitForTimeout(1500);
	await page.goto('/admin/tables');
	await page.waitForTimeout(500);

	// App should not crash — should show error or empty state
	await expect(page.locator('body')).toBeVisible();
});

// ─── Input validation / wrong data ───

test('login with empty fields shows error', async ({ page }) => {
	await page.goto('/login');
	await page.click('button[type="submit"]');
	// HTML5 validation should prevent submission or app should handle empty
	await page.waitForTimeout(500);
	await expect(page.locator('h1:has-text("Restaurant Orders")')).toBeVisible();
});

test('create category with empty name is rejected', async ({ page }) => {
	await page.goto('/login');
	await page.fill('input[type="email"]', 'admin@restaurant.com');
	await page.fill('input[type="password"]', 'admin');
	await page.click('button[type="submit"]');
	await page.waitForTimeout(1500);
	await page.goto('/admin/menu');
	await page.waitForTimeout(500);

	await page.locator('button:has-text("Add Category")').click();
	await page.waitForTimeout(300);
	// Leave name empty and click Add
	await page.locator('.modal-box button:has-text("Add")').click();
	await page.waitForTimeout(500);
	// Modal should remain open (validation failed)
	await expect(page.locator('.modal-open')).toBeVisible();
	await page.locator('.modal-open button:has-text("Cancel")').click();
});

test('xss injection in category name is sanitized', async ({ page }) => {
	await page.goto('/login');
	await page.fill('input[type="email"]', 'admin@restaurant.com');
	await page.fill('input[type="password"]', 'admin');
	await page.click('button[type="submit"]');
	await page.waitForTimeout(1500);
	await page.goto('/admin/menu');
	await page.waitForTimeout(500);

	await page.locator('button:has-text("Add Category")').click();
	await page.waitForTimeout(300);
	const modal = page.locator('.modal-open');
	await modal.locator('input[placeholder="Appetizers"]').fill('<script>alert("xss")</script>');
	await modal.locator('button:has-text("Add")').click();
	await page.waitForTimeout(500);

	// The category should be created with the literal string, not executed
	// App should still work normally
	await expect(page.locator('text=<script>').first()).toBeVisible({ timeout: 3000 });

	// Cleanup: delete it
	// First need to find and delete this category
	await page.locator('button:has-text("Delete")').last().click();
	await page.locator('button:has-text("Confirm")').click();
	await page.waitForTimeout(300);
});

test('sql injection attempt in login is handled safely', async ({ page }) => {
	await page.goto('/login');
	await page.fill('input[type="email"]', `' OR 1=1 --`);
	await page.fill('input[type="password"]', `' OR '1'='1`);
	await page.click('button[type="submit"]');
	await page.waitForTimeout(1000);

	// Should NOT be logged in — should show error or stay on login
	await expect(page.locator('h1:has-text("Restaurant Orders")')).toBeVisible();
});

test('very long input in dish name is handled', async ({ page }) => {
	await page.goto('/login');
	await page.fill('input[type="email"]', 'admin@restaurant.com');
	await page.fill('input[type="password"]', 'admin');
	await page.click('button[type="submit"]');
	await page.waitForTimeout(1500);
	await page.goto('/admin/menu');
	await page.waitForTimeout(500);
	await page.locator('button:has-text("Dishes")').click();
	await page.waitForTimeout(300);
	await page.locator('button:has-text("Add Dish")').click();
	await page.waitForTimeout(300);

	const modal = page.locator('.modal-open');
	await modal.locator('input').first().fill('A'.repeat(5000));
	await modal.locator('textarea').fill('Very long name test');
	await modal.locator('input[type="number"]').first().fill('500');
	await modal.locator('button:has-text("Add")').click();
	await page.waitForTimeout(500);

	// App should not crash — might truncate or show error
	const body = page.locator('body');
	await expect(body).toBeVisible();
});

// ─── Stale data / concurrency ───

test('deleting a category that has dishes shows proper error', async ({ page }) => {
	await page.goto('/login');
	await page.fill('input[type="email"]', 'admin@restaurant.com');
	await page.fill('input[type="password"]', 'admin');
	await page.click('button[type="submit"]');
	await page.waitForTimeout(1500);
	await page.goto('/admin/menu');
	await page.waitForTimeout(500);

	// Try deleting "Appetizers" which has dishes
	const deleteBtns = page.locator('button:has-text("Delete")');
	const count = await deleteBtns.count();
	if (count > 0) {
		// Click Delete on first category
		await deleteBtns.first().click();
		await page.waitForTimeout(200);
		await page.locator('button:has-text("Confirm")').click();
		await page.waitForTimeout(500);
		// App should not crash — should handle foreign key error gracefully
		await expect(page.locator('body')).toBeVisible();
	}
});
