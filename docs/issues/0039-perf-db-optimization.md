# 0039 — Performance: Database query optimization

## What to build

Profile and optimize slow database queries. Measure query execution time for the most frequently called endpoints, add missing indexes, and optimize N+1 query patterns.

## Acceptance criteria

- [ ] Profile GET /menu — query plan shows index scans, not sequential scans
- [ ] Profile GET /floor-plan — join with table_groups should use index
- [ ] Profile GET /audit-events with 100K+ audit rows
- [ ] Profile GET /orders (list all) with 1000+ orders per waiter
- [ ] Profile GET /dashboard/stats — aggregation queries should be fast
- [ ] Add indexes on: orders(waiter_id, created_at), order_items(course_id), audit_events(created_at)
- [ ] Add composite index on: orders(waiter_id, status, created_at) for waiter's order list
- [ ] Verify all queries use index-only scans where possible
- [ ] Document EXPLAIN ANALYZE output before and after optimization

## Blocked by

- 0038 — Performance: API load testing with k6
