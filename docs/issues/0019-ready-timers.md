# 0019 — Ready timers

## What to build

When the chef marks a dish as ready, a timer starts for the duration of that dish's eating time. The system tracks which timers are running per table group. Timers are used to calculate when the table group will be free.

## Acceptance criteria

- [ ] Go service: when item marked ready, create a timer record with expected_end = now + eating_time
- [ ] Database schema for `ready_timers` (id, order_item_id, started_at, eating_time_min, expected_end_at)
- [ ] Go endpoint: GET /table-groups/:id/timers (return running and completed timers for the group)
- [ ] Timer resolution: 1-minute granularity is sufficient
- [ ] When a table group closes, all its timers are invalidated

## Blocked by

- 0017 — Chef marks ready
