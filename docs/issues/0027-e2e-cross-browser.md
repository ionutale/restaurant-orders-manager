# 0027 — E2E: Cross-browser testing (Firefox + WebKit)

## What to build

Run the full E2E suite against Firefox and WebKit (Safari) in addition to Chromium. The Playwright config currently only tests Chromium.

## Acceptance criteria

- [ ] Playwright config includes firefox and webkit projects
- [ ] All 43 tests pass on Firefox
- [ ] All 43 tests pass on WebKit
- [ ] Browser-specific failures are fixed or conditionally handled

## Blocked by

- 0025 — E2E: Mobile and tablet viewport testing
