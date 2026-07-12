# 0037 — Integration tests: Edge cases and data integrity

## What to build

Add integration tests for edge cases that could break data integrity: concurrent requests, partial updates, cascading deletes, and boundary values.

## Acceptance criteria

- [ ] Delete a category that has dishes → should fail (foreign key)
- [ ] Delete a table that is in an active group → should fail
- [ ] Delete a user who has active orders → should fail
- [ ] Create order for a closed table group → should fail
- [ ] Add item to a sent order → should fail (order locked)
- [ ] Advance course when no more courses → order becomes completed
- [ ] Advance course on a pending order → should fail (not sent yet)
- [ ] Mark item ready twice → second call is idempotent
- [ ] Seat a table that is already occupied → should fail
- [ ] Partial update (PATCH with single field) preserves other fields
- [ ] Create order with 0 courses (API default) → creates one default course
- [ ] String truncation: very long dish name (5000 chars) → handled gracefully

## Blocked by

- 0035 — Integration tests: Input validation and error responses
