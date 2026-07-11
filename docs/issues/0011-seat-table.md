# 0011 — Seat a table

## What to build

Waiter can seat guests at a free table: creates a table group with party size, optional custom name (e.g., "Birthday party — Ana"), and the table group status set to "open". The waiter is assigned as the responsible staff.

## Acceptance criteria

- [ ] Database schema for `table_groups` (id, name, party_size, status, waiter_id, opened_at, closed_at)
- [ ] Database schema for `table_group_tables` (group_id, table_id)
- [ ] Go endpoint: POST /table-groups (creates group + links table, waiter from JWT)
- [ ] Go endpoint: GET /table-groups/:id
- [ ] Waiter UI: click a free table → "Seat guests" dialog → party size + optional name → confirm
- [ ] Floor plan updates to show table as occupied after seating
- [ ] Mutations recorded in audit_events

## Blocked by

- 0010 — Waiter floor plan view
