# 0046 — Regression testing: CI pipeline and automated test suite

## What to build

Set up a CI pipeline (GitHub Actions) that runs the full test suite on every push and pull request. Add visual regression testing for UI components and snapshot testing for API responses.

## Acceptance criteria

- [ ] GitHub Actions workflow runs on push and PR to main
- [ ] Workflow: build Go backend → run integration tests → build frontend → run E2E tests
- [ ] Docker Compose services start before tests (db, backend, frontend)
- [ ] Integration tests take < 30s, E2E tests take < 3min
- [ ] Workflow fails fast: lint → unit → integration → E2E (stop on first failure)
- [ ] Visual regression: Playwright screenshot comparison for key pages (login, floor plan, order detail, KDS)
- [ ] API snapshot: record and diff API responses for GET /menu, GET /floor-plan
- [ ] Test results artifacts uploaded on failure (screenshots, traces, logs)
- [ ] Notifications on failure (email, Slack, or GitHub checks)
- [ ] Document how to add new E2E tests and update snapshots

## Blocked by

- 0001 — Project scaffold + auth
