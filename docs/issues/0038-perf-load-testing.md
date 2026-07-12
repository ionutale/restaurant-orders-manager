# 0038 — Performance: API load testing with k6

## What to build

Add load testing scripts using k6 (or similar) that simulate real-world usage patterns: multiple waiters taking orders, chef marking items ready, and admin managing the menu concurrently.

## Acceptance criteria

- [ ] k6 script for waiter workflow: login → seat → order → add dishes → send to KDS
- [ ] k6 script for chef workflow: view KDS → mark items ready
- [ ] k6 script for admin: manage menu, view audit log
- [ ] Run with 10, 25, 50, 100 concurrent virtual users
- [ ] Document P50, P95, P99 response times for each endpoint
- [ ] Identify and fix any endpoint slower than 500ms at P95
- [ ] Test POST /auth/login under load (password hashing is CPU-intensive)
- [ ] Test GET /menu (largest payload with dish + allergen + suggestion data)
- [ ] Report peak memory and CPU usage during load test

## Blocked by

- 0001 — Project scaffold + auth
