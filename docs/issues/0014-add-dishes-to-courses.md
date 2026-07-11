# 0014 — Add dishes to courses

## What to build

Waiter browses the menu (categories + chef suggestions) and adds items to specific course slots. Each dish has a quantity and optional note. Suggested wines/sides are shown when a dish is selected. Allergen icons are visible in the menu browser.

## Acceptance criteria

- [ ] Database schema for `order_items` (id, course_id, dish_id, is_chef_suggestion, quantity, notes, added_at)
- [ ] Go endpoint: POST /orders/:id/courses/:course_id/items (add item with qty, notes)
- [ ] Go endpoint: DELETE /order-items/:id
- [ ] Go endpoint: GET /menu returns categories + dishes + allergens + suggestions
- [ ] Waiter UI: course tab view with "+ Add Dish" → menu browser with category filters, allergen icons, suggestion pills, dish details
- [ ] Selecting a dish shows suggested wines/sides with one-tap add
- [ ] Each dish item has a notes field (e.g., "no onions, well done")
- [ ] Order summary shows all courses with their items

## Blocked by

- 0006 — Allergens + dish suggestions
- 0008 — Suggestions in waiter menu
- 0013 — Create order with courses
