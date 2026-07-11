# 0017 — Chef marks item ready

## What to build

Chef can mark individual items in the active course as "ready" on the KDS. A visual indicator shows completion progress per order (e.g., "3/5 items ready"). Orders with all items marked ready are highlighted to signal the waiter.

## Acceptance criteria

- [ ] Go endpoint: PATCH /kds/order-items/:id/ready (marks item as ready, records timestamp)
- [ ] KDS UI: each item has a "Mark Ready" button, which changes to "Ready ✓" when done
- [ ] Order card shows progress bar: items ready / total items in active course
- [ ] When all items in the active course are ready, the order card is visually highlighted
- [ ] Mutations recorded in audit_events

## Blocked by

- 0016 — KDS dashboard
