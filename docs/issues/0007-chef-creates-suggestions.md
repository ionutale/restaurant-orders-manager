# 0007 — Chef creates suggestions

## What to build

Chef can create time-bounded menu suggestions (daily specials) with name, description, price, and a shift/date range. These are not linked to the permanent dish catalog — they're ad-hoc items that exist for a shift and auto-expire.

## Acceptance criteria

- [ ] Database schema for `chef_suggestions` (id, name, description, price_cents, shift_date, expires_at, chef_id, created_at)
- [ ] Go CRUD endpoints for chef suggestions (chef can create/edit/delete, admin can admin all)
- [ ] Chef UI: suggestion creation form, list of active suggestions, delete
- [ ] Suggestions auto-expire (query filters by expires_at)
- [ ] Mutations recorded in audit_events

## Blocked by

- 0001 — Project scaffold + auth
