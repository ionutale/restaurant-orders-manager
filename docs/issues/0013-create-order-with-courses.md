# 0013 — Create order with courses

## What to build

Waiter can create an order for an occupied table group. The order has multiple course slots (e.g., appetizer, main, dessert) defined by the waiter. Each course has a name and display order. The order also tracks the table group and responsible waiter.

## Acceptance criteria

- [ ] Database schema for `orders` (id, table_group_id, waiter_id, status, created_at)
- [ ] Database schema for `courses` (id, order_id, name, display_order, status: pending/active/completed)
- [ ] Go endpoint: POST /orders (creates order with initial courses)
- [ ] Go endpoint: GET /orders/:id (returns order with courses)
- [ ] Waiter UI: select table group → "New Order" → define course names (defaults: Appetizer, Main, Dessert) → create
- [ ] Order shows in the waiter's active orders list
- [ ] Mutations recorded in audit_events

## Blocked by

- 0012 — Merge + close groups
