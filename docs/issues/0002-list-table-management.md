# 0002 — List-based table management

## What to build

Admin can view, create, edit, and delete physical tables via a list/table view. Each table has a name/number, capacity, coordinates (x, y for floor plan later), and an optional label. No drag-drop canvas yet — this is the data entry foundation.

## Acceptance criteria

- [ ] Database schema for `tables` (id, name, capacity, x, y, label, created_at, updated_at)
- [ ] Go CRUD endpoints: GET /tables, POST /tables, PATCH /tables/:id, DELETE /tables/:id
- [ ] Admin UI page with sortable table list, inline edit, delete confirmation
- [ ] All mutations recorded in audit_events

## Blocked by

- 0001 — Project scaffold + auth
