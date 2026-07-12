package integration

import (
	"testing"
)

// TestCleanupAllGroups runs first (alphabetically) and closes all open table groups
// to ensure a clean state for subsequent tests.
func TestCleanupAllGroups(t *testing.T) {
	token := login(t, "admin")
	closeAllGroups(t, token)
	t.Log("All open groups closed")
}
