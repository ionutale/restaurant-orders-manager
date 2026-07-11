# 0021 — Availability prediction

## What to build

The waiter floor plan shows predicted free time for each occupied table, calculated from the longest running ready timer. When dessert/coffee course is served (identified by course name or items), the prediction assumes the table will finish soon and shows "almost done" status.

## Acceptance criteria

- [ ] Go endpoint: GET /floor-plan (extended with prediction fields: estimated_free_at, status_hint)
- [ ] Prediction logic: the latest expected_end_at among running timers for the group = estimated free time
- [ ] When dessert/coffee course is active, adjust prediction to a shorter window
- [ ] Floor plan UI shows: "Free ~20:35" or "Almost done" beneath occupied tables
- [ ] Live updates as timers progress and courses advance

## Blocked by

- 0019 — Ready timers
- 0018 — Course advancement
