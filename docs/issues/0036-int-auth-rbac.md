# 0036 — Integration tests: Authorization and role-based access

## What to build

Add integration tests that verify API endpoints correctly enforce role-based access control. Each endpoint should be tested with tokens for all roles (waiter, chef, manager) and without authentication.

## Acceptance criteria

- [ ] GET /tables without token → 401
- [ ] GET /tables with waiter token → 200
- [ ] GET /tables with chef token → 200
- [ ] POST /tables with waiter token → 200 (allowed for all authenticated)
- [ ] POST /categories with waiter token → 200 (allowed for all authenticated)
- [ ] POST /users with waiter token → 403 (admin only)
- [ ] GET /audit-events with waiter token → 200 (or 403 if manager-only)
- [ ] GET /audit-events without token → 401
- [ ] POST /orders with chef token → 200 (allowed for waiters/chefs/managers)
- [ ] PATCH /kds/order-items/:id/ready with waiter token → 403 (chef only)
- [ ] GET /kds/orders with waiter token → 403 (chef only)
- [ ] All admin-only endpoints reject waiter/chef tokens with 403
- [ ] All waiter endpoints reject unauthenticated requests with 401

## Blocked by

- 0001 — Project scaffold + auth
