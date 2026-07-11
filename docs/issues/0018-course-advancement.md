# 0018 — Course advancement

## What to build

Waiter advances the current course to the next one on their device. The KDS view flips: the completed course is hidden, and the next course becomes active. If no more courses remain, the order status changes to "completed".

## Acceptance criteria

- [ ] Go endpoint: POST /orders/:id/advance-course (sets current course completed, activates next, or marks order completed if none remain)
- [ ] Waiter UI: "Advance to next course" button on the order detail view
- [ ] KDS UI: auto-updates (poll/SSE) — previous course disappears, next course appears
- [ ] Last course advancement marks order as "completed"
- [ ] Mutations recorded in audit_events

## Blocked by

- 0017 — Chef marks ready
