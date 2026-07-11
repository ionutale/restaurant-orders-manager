package integration

import (
	"fmt"
	"testing"
)

type userResp struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

func TestUsersCRUD(t *testing.T) {
	adminToken := login(t, "admin")

	// CREATE
	created := apiPOST(t, adminToken, "/users", map[string]string{
		"name": "IT-TestWaiter", "email": "itwaiter@test.com", "password": "test123", "role": "waiter",
	})
	assertStatus(t, created, 200)
	u := decode[userResp](t, created)
	if u.Name != "IT-TestWaiter" || u.Role != "waiter" {
		t.Errorf("unexpected created user: %+v", u)
	}
	t.Logf("CREATE ok: id=%d %s (%s)", u.ID, u.Name, u.Role)

	// READ (list)
	listResp := apiGET(t, adminToken, "/users")
	assertStatus(t, listResp, 200)
	list := decode[[]userResp](t, listResp)
	found := false
	for _, us := range list {
		if us.ID == u.ID {
			found = true
			break
		}
	}
	if !found {
		t.Error("created user not found in listing")
	}
	t.Logf("LIST ok: %d users", len(list))

	// UPDATE
	updated := apiPATCH(t, adminToken, "/users/"+fmt.Sprintf("%d", u.ID), map[string]string{
		"name": "IT-TestWaiterRenamed",
	})
	assertStatus(t, updated, 200)
	u2 := decode[userResp](t, updated)
	if u2.Name != "IT-TestWaiterRenamed" {
		t.Errorf("unexpected updated name: %s", u2.Name)
	}
	t.Logf("UPDATE ok: renamed to %s", u2.Name)

	// DELETE
	delResp := apiDELETE(t, adminToken, "/users/"+fmt.Sprintf("%d", u.ID))
	assertStatus(t, delResp, 204)
	t.Logf("DELETE ok")

	// Verify deletion
	listResp2 := apiGET(t, adminToken, "/users")
	list2 := decode[[]userResp](t, listResp2)
	for _, us := range list2 {
		if us.ID == u.ID {
			t.Error("user was not deleted")
		}
	}
	t.Log("VERIFY deletion ok")
}
