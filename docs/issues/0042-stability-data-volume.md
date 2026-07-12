# 0042 — Stability: Data volume stress test

## What to build

Test the application with production-scale data volumes: hundreds of tables, thousands of orders, tens of thousands of audit events. Verify that UI pages render and API queries complete within acceptable time.

## Acceptance criteria

- [ ] Insert 500 tables (admin floor plan scrolls/renders smoothly)
- [ ] Insert 200 categories with 5000 dishes total
- [ ] Insert 10,000 completed orders with 50,000 order items
- [ ] Insert 100,000 audit events
- [ ] Admin menu page loads in under 3 seconds
- [ ] Admin audit log page loads in under 3 seconds
- [ ] Waiter floor plan loads in under 2 seconds
- [ ] Waiter orders list loads in under 2 seconds
- [ ] Dashboard stats query completes in under 3 seconds
- [ ] Invoice preview loads in under 2 seconds
- [ ] KDS dashboard loads in under 2 seconds (only shows active orders)
- [ ] Menu search with 5000 dishes responds within 500ms

## Blocked by

- 0039 — Performance: Database query optimization
