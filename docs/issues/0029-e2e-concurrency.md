# 0029 — E2E: Concurrent waiter scenarios

## What to build

Test race conditions when two waiters interact with the same table group or order simultaneously. For example: both try to seat the same table, both add dishes to the same order, or one closes a group while the other is adding items.

## Acceptance criteria

- [ ] Two simultaneous attempts to seat the same table — second gets error
- [ ] Adding items to a just-closed group shows appropriate error
- [ ] Sending an already-sent order to KDS returns error
- [ ] Advancing a course on a completed order is blocked
- [ ] No data corruption occurs from concurrent requests

## Blocked by

- 0012 — Merge + close groups
- 0015 — Send to KDS
- 0018 — Course advancement
