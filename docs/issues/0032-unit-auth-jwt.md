# 0032 — Unit tests: Go auth middleware and JWT

## What to build

Unit test the auth middleware and JWT token generation/validation functions in `internal/auth/`. Test token creation, parsing, expiry, malformed tokens, and role checks.

## Acceptance criteria

- [ ] `JWT.Generate()` creates valid tokens for all roles
- [ ] `JWT.Validate()` parses valid tokens and returns correct claims
- [ ] `JWT.Validate()` rejects expired tokens
- [ ] `JWT.Validate()` rejects malformed tokens
- [ ] `JWT.Validate()` rejects tokens signed with wrong secret
- [ ] `RequireRole()` middleware passes correct roles
- [ ] `RequireRole()` rejects incorrect roles
- [ ] `RequireRole()` rejects unauthenticated requests
- [ ] Tests use `httptest.NewRecorder` and `httptest.NewRequest`

## Blocked by

- 0001 — Project scaffold + auth
