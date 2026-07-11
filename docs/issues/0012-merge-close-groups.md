# 0012 — Merge + close groups

## What to build

Waiter can add free tables to an existing table group (merge) and split tables out of a group. Closing a group frees all constituent tables simultaneously. A multi-table group shows on the floor plan as a combined highlighted area.

## Acceptance criteria

- [ ] Go endpoint: PATCH /table-groups/:id/tables (add/remove tables)
- [ ] Go endpoint: PATCH /table-groups/:id/close (sets status=closed, frees tables)
- [ ] Waiter UI: from group detail, "Add table" picker showing only free adjacent tables, "Remove" button per table in group
- [ ] Waiter UI: "Close group" button with confirmation
- [ ] Floor plan updates in real-time when group changes
- [ ] Mutations recorded in audit_events

## Blocked by

- 0011 — Seat a table
