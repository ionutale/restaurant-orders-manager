# 0050 — Black-box testing: Equivalence partitioning and boundary analysis

## What to build

Apply systematic black-box testing techniques to the API: equivalence partitioning (group inputs into valid/invalid classes), boundary value analysis (test at edges of ranges), and state transition testing (test all order/group lifecycle states).

## Acceptance criteria

- [ ] Equivalence partitioning for price: 0, 1, 999999, negative, non-integer
- [ ] Equivalence partitioning for eating_time_min: 0, 1, 999, negative
- [ ] Equivalence partitioning for capacity: 0, 1, 99, negative
- [ ] Equivalence partitioning for party_size: 0, 1, 99, negative
- [ ] Boundary analysis for page size / limit parameters
- [ ] State transition testing for order lifecycle: pending → sent → completed → paid
- [ ] State transition testing for order lifecycle invalid transitions: pending → completed (skip sent)
- [ ] State transition testing for table group: open → in_progress → closed
- [ ] State transition testing for course: pending → active → completed
- [ ] State transition testing for course invalid: completed → active (cannot go back)
- [ ] Decision table for auth: all combinations of {no token, invalid token, expired token, waiter, chef, manager} × {admin endpoint, waiter endpoint, chef endpoint, public endpoint}
- [ ] Error guessing: empty arrays, null values, missing required fields, extra unknown fields

## Blocked by

- 0035 — Integration tests: Input validation and error responses
- 0036 — Integration tests: Authorization and role-based access
- 0037 — Integration tests: Edge cases and data integrity
