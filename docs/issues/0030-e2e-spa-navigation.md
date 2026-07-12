# 0030 — E2E: Page refresh and SPA navigation

## What to build

Test that the app preserves state across page refreshes (F5) and that browser back/forward navigation works correctly without losing data or showing stale state.

## Acceptance criteria

- [ ] Refreshing the order detail page keeps the order data
- [ ] Refreshing the floor plan shows tables with correct status
- [ ] Browser back from order detail returns to orders list
- [ ] Browser forward after back restores order detail
- [ ] Logging out and pressing back does not show cached authenticated page
- [ ] Token expiry during a session is handled (redirect to login)

## Blocked by

- 0013 — Create order with courses
