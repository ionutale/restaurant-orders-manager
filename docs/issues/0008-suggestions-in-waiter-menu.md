# 0008 — Suggestions in waiter menu

## What to build

When the waiter browses the menu to build an order, active chef suggestions appear in a special "Chef's Suggestions" section alongside the regular menu categories. The waiter can add them to an order just like regular dishes.

## Acceptance criteria

- [ ] Go endpoint: GET /menu returns both regular dishes (grouped by category) and active chef suggestions in a separate "suggestions" section
- [ ] Waiter UI menu browser shows "Chef's Suggestions" as a section with the same add-to-order interaction as regular dishes
- [ ] Suggestion items in an order are distinguishable from regular dishes

## Blocked by

- 0007 — Chef creates suggestions
- 0001 — Project scaffold + auth
