# 0005 — Dishes CRUD

## What to build

Admin can create, edit, and delete dishes. Each dish has a name, description, price (EUR), category reference, eating time (5–15 min), and optional image. Created in a specific category. The dish list is filterable by category.

## Acceptance criteria

- [ ] Database schema for `dishes` (id, name, description, price_cents, category_id, eating_time_min, image_url, created_at)
- [ ] Go CRUD endpoints for dishes, filtered by category
- [ ] Admin UI: dish list grouped/filterable by category, add/edit form with all fields
- [ ] Mutations recorded in audit_events

## Blocked by

- 0004 — Categories CRUD
