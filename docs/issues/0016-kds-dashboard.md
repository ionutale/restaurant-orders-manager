# 0016 — KDS dashboard

## What to build

Chef logs into the KDS view and sees a list of all sent orders. Each order shows the table group name, waiter, and the current active course with its items. Only the active course per order is visible. Orders are sorted by recency.

## Acceptance criteria

- [ ] Go endpoint: GET /kds/orders (returns all sent orders with only the active course items visible)
- [ ] KDS UI: full-screen tablet-friendly view, auto-refreshing (poll or SSE)
- [ ] Each order card shows: table name, waiter name, course name, item list (qty, name, notes highlighted)
- [ ] Orders sorted by send time, newest first
- [ ] Allergen icons shown per dish
- [ ] Chef authentication scoped to KDS role

## Blocked by

- 0015 — Send to KDS
