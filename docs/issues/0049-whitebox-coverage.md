# 0049 — White-box testing: Code coverage and internal logic

## What to build

Add code coverage analysis and write white-box unit tests for internal functions that are not directly tested through the API. This includes helper functions, formatters, validators, and business logic.

## Acceptance criteria

- [ ] `go test -coverprofile` reaches > 60% coverage on `internal/` package
- [ ] Test price formatting: cents to EUR string conversion
- [ ] Test time parsing: various date formats, timezone handling
- [ ] Test order status transitions: pending → sent → completed → paid
- [ ] Test course advancement logic: which course becomes active next
- [ ] Test notification query logic: which items are "ready to notify"
- [ ] Test prediction calculation: estimated free time from item timers
- [ ] Test invoice totals: sum of items, per-course subtotals
- [ ] Test audit event payload serialization
- [ ] Test all error paths: database errors, JSON parsing errors, empty results
- [ ] Add `make cover` target that runs coverage and opens HTML report

## Blocked by

- 0031 — Unit tests: Go backend handlers
- 0048 — Mutation testing: Go backend code
