import { defineConfig } from '@playwright/test';

export default defineConfig({
	testDir: './tests/e2e',
	timeout: 60000,
	fullyParallel: false,
	workers: 1,
	retries: 1,
	use: {
		baseURL: 'http://localhost:3000',
		headless: true,
	},
	webServer: {
		command: 'echo "Using existing Docker setup"',
		url: 'http://localhost:3000',
		reuseExistingServer: true,
	},
});
