# 0010 — Waiter floor plan view

## What to build

Waiter sees the floor plan with all tables, color-coded by status: free (green), occupied (red). Clicking a table shows its current details (group name, party size, waiter, status). Free tables show capacity. Uses the same canvas component from the admin floor plan but read-only for positioning.

## Acceptance criteria

- [ ] Go endpoint: GET /floor-plan returns tables with current status (free/occupied/grouped)
- [ ] Waiter UI: canvas view of floor plan with status colors and tooltip on click
- [ ] Toggle between canvas and list view
- [ ] Free tables show seat capacity, occupied tables show party name
- [ ] Responsive — mobile and tablet friendly

## Blocked by

- 0003 — Drag-drop floor plan
