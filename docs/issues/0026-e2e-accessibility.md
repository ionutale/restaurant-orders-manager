# 0026 — E2E: Accessibility (keyboard + screen reader)

## What to build

Test that all major user flows work with keyboard-only navigation and screen readers. Verify aria labels, roles, focus management, and announcements.

## Acceptance criteria

- [ ] Tab order follows visual layout
- [ ] All buttons, links, inputs reachable via keyboard
- [ ] Modals trap focus when open
- [ ] Toast/alert announcements are read by screen readers
- [ ] Floor plan canvas tables are focusable and activatable
- [ ] Menu dish items are keyboard-navigable
- [ ] Drag-and-drop has keyboard alternative
- [ ] Login error messages announced

## Blocked by

- 0001 — Project scaffold + auth
