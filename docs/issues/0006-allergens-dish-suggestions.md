# 0006 — Allergens + dish suggestions

## What to build

Admin manages a predefined allergen list (shared pool, e.g., 14 EU allergens) and assigns them to dishes. Also manages dish-to-dish suggestions: "suggested wine" and "suggested side" cross-references between dishes.

## Acceptance criteria

- [ ] Database schema for `allergens` and `dish_allergens` (many-to-many)
- [ ] Database schema for `dish_suggestions` (from_dish_id, to_dish_id, type: wine/side)
- [ ] Go endpoints: CRUD allergens, assign/remove allergens on dishes, CRUD dish suggestions
- [ ] Admin UI: allergen list management, per-dish allergen picker, per-dish suggestion picker (filtered by category — wine list for wines, sides for sides)
- [ ] Mutations recorded in audit_events

## Blocked by

- 0005 — Dishes CRUD
