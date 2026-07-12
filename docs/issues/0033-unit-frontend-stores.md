# 0033 — Unit tests: Svelte stores and API client

## What to build

Unit test the frontend auth store (`auth.svelte.ts`) and API client (`api.ts`). Test login/logout state transitions, token persistence, API error handling, and role-based redirects.

## Acceptance criteria

- [ ] Auth store starts with `loading=true` and `user=null`
- [ ] Successful login sets token, user, and role
- [ ] Login failure does not change state
- [ ] Logout clears token and user
- [ ] `isLoggedIn` reflects token presence
- [ ] API client throws `ApiError` with status and message on 4xx/5xx
- [ ] API client sends Authorization header when token is present
- [ ] API client does not send Authorization header when no token
- [ ] Tests use Vitest (already in devDependencies)

## Blocked by

- 0001 — Project scaffold + auth
