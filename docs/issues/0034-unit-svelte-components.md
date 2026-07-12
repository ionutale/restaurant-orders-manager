# 0034 — Unit tests: Svelte component rendering

## What to build

Add component-level unit tests for key Svelte components: FloorPlanCanvas, login form, order detail page, and dish menu browser. Test conditional rendering, prop handling, and user interactions.

## Acceptance criteria

- [ ] FloorPlanCanvas renders tables at correct positions
- [ ] FloorPlanCanvas shows free/occupied colors correctly
- [ ] FloorPlanCanvas readonly mode disables drag
- [ ] Login form shows error on API failure
- [ ] Login form disables button while loading
- [ ] Order detail page shows loading state then content
- [ ] Add Dish button hidden after order is sent
- [ ] Course pills show correct item counts
- [ ] Components use Vitest + @testing-library/svelte or svelte-testing-library
- [ ] Components render without crashing (smoke test)

## Blocked by

- 0033 — Unit tests: Svelte stores and API client
