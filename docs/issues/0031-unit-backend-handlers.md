# 0031 — Unit tests: Go backend handlers

## What to build

Add Go unit tests for all handler functions in `internal/handler/`. Test input validation, error responses, edge cases, and business logic decoupled from the database using mocked interfaces.

## Acceptance criteria

- [ ] Table handler: invalid name, negative capacity, missing fields
- [ ] Category handler: empty name, duplicate name, reorder out of bounds
- [ ] Dish handler: missing category_id, price_cents overflow, invalid category
- [ ] Order handler: missing table_group_id, empty course names, send already-sent order
- [ ] Auth handler: invalid credentials, duplicate email, missing fields
- [ ] Chef suggestion handler: expired dates, missing chef_id
- [ ] All tests run with `go test ./internal/handler/` (no database needed)

## Blocked by

- 0001 — Project scaffold + auth
