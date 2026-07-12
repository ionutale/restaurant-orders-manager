# 0035 — Integration tests: Input validation and error responses

## What to build

Add integration tests for every API endpoint that send invalid, empty, or malformed request bodies. Verify the API returns correct HTTP status codes (400, 404, 409, 500) and meaningful error messages.

## Acceptance criteria

- [ ] POST /tables with empty name → 400
- [ ] POST /tables with negative capacity → 400
- [ ] PATCH /tables/:id with non-existent ID → 404
- [ ] DELETE /tables/:id with non-existent ID → 404
- [ ] POST /categories with empty name → 400
- [ ] POST /categories with duplicate name → 409 or 400
- [ ] POST /dishes without category_id → 400
- [ ] PATCH /dishes/:id with invalid ID → 404
- [ ] POST /orders without table_group_id → 400
- [ ] POST /orders with non-existent table_group_id → 404
- [ ] POST /orders/:id/send on already-sent order → 400
- [ ] POST /auth/login with wrong password → 401
- [ ] POST /auth/register with duplicate email → 409
- [ ] POST /auth/register with invalid role → 400
- [ ] PATCH /users/:id with non-existent ID → 404
- [ ] DELETE /users/:id when user has orders → 500 or 409
- [ ] POST /upload with no file → 400
- [ ] POST /upload with invalid file type → 400

## Blocked by

- 0001 — Project scaffold + auth
