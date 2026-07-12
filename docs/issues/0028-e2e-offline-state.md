# 0028 — E2E: Offline and disconnected state

## What to build

Test app behavior when the network is unavailable or the backend is unreachable. Verify graceful degradation and user-friendly error messages.

## Acceptance criteria

- [ ] Floor plan shows cached/empty state when API unreachable
- [ ] Login shows "cannot connect" error when backend is down
- [ ] Order creation fails gracefully with retry suggestion
- [ ] KDS dashboard shows connection lost indicator
- [ ] App does not crash or show blank white screen on disconnect
- [ ] Reconnection after restore works without page reload

## Blocked by

- 0001 — Project scaffold + auth
