# 0020 — Waiter notification

## What to build

When the chef marks items as ready, the responsible waiter receives a notification on their device (e.g., "Order ready — Table T3: 4 items ready"). This can be implemented via polling, SSE, or WebSocket.

## Acceptance criteria

- [ ] Go notification endpoint or SSE stream: GET /notifications/stream (per-waiter, events for "items_ready")
- [ ] Waiter UI: notification toast/badge when items are ready, clickable to navigate to the order
- [ ] Notification includes: table name, order ID, count of ready items
- [ ] Notifications are dismissed when the waiter navigates to the order
- [ ] Integrates with the ready timers — when all items in a course are ready, send one notification

## Blocked by

- 0017 — Chef marks ready
