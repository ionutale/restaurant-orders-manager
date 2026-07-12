# 0041 — Stability: Endurance testing (long-running)

## What to build

Run the application under continuous moderate load for an extended period (4-24 hours) to detect memory leaks, connection pool exhaustion, goroutine leaks, and slow resource accumulation.

## Acceptance criteria

- [ ] Run 10 concurrent users repeating waiter workflows for 4 hours
- [ ] Monitor memory usage — no linear growth over time
- [ ] Monitor goroutine count — no goroutine leaks
- [ ] Monitor Postgres connection count — connections released properly
- [ ] Monitor open file descriptors — no file handle leaks
- [ ] Audit log grows to 10K+ entries with no performance degradation
- [ ] Database vacuum/analyze runs without blocking reads
- [ ] Log rotation does not fill disk
- [ ] Post-test: verify all metrics return to baseline after load stops

## Blocked by

- 0038 — Performance: API load testing with k6
- 0039 — Performance: Database query optimization
