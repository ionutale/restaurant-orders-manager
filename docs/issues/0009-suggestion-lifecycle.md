# 0009 — Suggestion lifecycle

## What to build

Chef suggestions auto-expire after their shift. The admin can view expired/archived suggestions, and the system cleans up or archives them. Suggestions that expired are hidden from the waiter menu. Optionally, the chef can extend or renew an expired suggestion.

## Acceptance criteria

- [ ] Expired suggestions are automatically excluded from the waiter menu query
- [ ] Admin/chef can view a history of past suggestions
- [ ] Chef can renew an expired suggestion (creates a new one with updated shift date)
- [ ] Optional: admin can permanently delete archived suggestions
- [ ] Mutations recorded in audit_events

## Blocked by

- 0008 — Suggestions in waiter menu
