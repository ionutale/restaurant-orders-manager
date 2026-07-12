# 0025 — E2E: Mobile and tablet viewport testing

## What to build

All existing E2E tests should run against mobile (375px) and tablet (768px) viewport sizes, not just desktop. The app is mobile-first but we only test at desktop width.

## Acceptance criteria

- [ ] Playwright config includes mobile and tablet projects
- [ ] All 43 existing tests pass at 375px viewport
- [ ] All 43 existing tests pass at 768px viewport
- [ ] Layout-specific bugs found at small viewports are fixed
- [ ] Touch interactions work (tap, long-press, swipe)

## Blocked by

- 0001 — Project scaffold + auth
