# 0015 — Send to KDS

## What to build

Waiter confirms and sends the full order to the KDS. The order status changes to "sent". The full order payload (all courses + items) is stored and the first course becomes "active". The waiter can no longer modify sent items.

## Acceptance criteria

- [ ] Go endpoint: POST /orders/:id/send (validates order is complete, sets status=sent, activates first course)
- [ ] Go endpoint: GET /orders/:id (returns full order with all courses and items, ready for KDS)
- [ ] Waiter UI: review order summary → "Send to KDS" button → confirmation → success state
- [ ] After send, order items become read-only in the waiter UI
- [ ] Waiter sees a "Sent" badge on the order
- [ ] Mutations recorded in audit_events

## Blocked by

- 0014 — Add dishes to courses
