package integration

import (
	"fmt"
	"testing"
	"time"
)

type suggestion struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	PriceCents  int    `json:"price_cents"`
	ChefName    string `json:"chef_name,omitempty"`
}

func TestChefSuggestionsCRUD(t *testing.T) {
	chefToken := login(t, "chef")
	adminToken := login(t, "admin")

	expiresAt := time.Now().Add(8 * time.Hour).Format(time.RFC3339)

	// CREATE
	created := apiPOST(t, chefToken, "/chef-suggestions", map[string]interface{}{
		"name": "IT-Special", "description": "test", "price_cents": 2500, "expires_at": expiresAt,
	})
	assertStatus(t, created, 200)
	s := decode[suggestion](t, created)
	if s.Name != "IT-Special" || s.PriceCents != 2500 {
		t.Errorf("unexpected created: %+v", s)
	}
	t.Logf("CREATE ok: id=%d", s.ID)

	// READ active
	listResp := apiGET(t, chefToken, "/chef-suggestions")
	assertStatus(t, listResp, 200)
	t.Logf("READ active ok")

	// READ all
	allResp := apiGET(t, chefToken, "/chef-suggestions?all=true")
	assertStatus(t, allResp, 200)
	t.Logf("READ all ok")

	// RENEW
	renewed := apiPOST(t, chefToken, "/chef-suggestions/"+fmt.Sprintf("%d", s.ID)+"/renew", nil)
	assertStatus(t, renewed, 200)
	r := decode[suggestion](t, renewed)
	t.Logf("RENEW ok: new id=%d", r.ID)

	// DELETE (as admin)
	delResp := apiDELETE(t, adminToken, "/chef-suggestions/"+fmt.Sprintf("%d", r.ID))
	assertStatus(t, delResp, 204)
	t.Logf("DELETE ok")
}
