# 0047 — Smoke testing: Critical-path health check suite

## What to build

Create a minimal smoke test suite that verifies the most critical user paths in under 30 seconds. This runs on every deployment to quickly detect if the app is broken, before the full regression suite runs.

## Acceptance criteria

- [ ] Smoke suite runs in < 30 seconds
- [ ] Test 1: Login as admin, check dashboard returns 200
- [ ] Test 2: List tables via API, verify at least 1 table exists
- [ ] Test 3: Login as waiter, check floor plan loads
- [ ] Test 4: Login as chef, check KDS loads
- [ ] Test 5: GET /api/health returns {"status":"ok"}
- [ ] Test 6: GET /api/menu returns categories and dishes
- [ ] Smoke suite runs before E2E regression suite in CI
- [ ] Smoke failure immediately stops deployment pipeline
- [ ] Smoke suite can be run locally with `pnpm test:smoke` and `go test ./tests/smoke/`

## Blocked by

- 0046 — Regression testing: CI pipeline and automated test suite
