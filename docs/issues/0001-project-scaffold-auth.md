# 0001 — Project scaffold + auth + empty shell

## What to build

Set up the Go project structure with Postgres connection, migration tooling, user authentication (waiter/chef/manager roles), a Svelte app shell with login page, and routing skeleton for admin/waiter/chef views. Audit infrastructure is wired in from the start: all state-changing operations write to an immutable `audit_events` table.

## Acceptance criteria

- [ ] Go project with `cmd/`, `internal/` structure, `go.mod`
- [ ] Postgres connection via `pgx` or `sqlx`
- [ ] Migration tool (golang-migrate or goose) with initial schema for `users` and `audit_events`
- [ ] Auth endpoints: login, JWT/token validation, role-based middleware
- [ ] Svelte project with Tailwind + DaisyUI, login page, role-based route guards
- [ ] Empty placeholder pages for admin, waiter, chef dashboards
- [ ] Audit log records all authenticated state changes automatically

## Blocked by

None — can start immediately
