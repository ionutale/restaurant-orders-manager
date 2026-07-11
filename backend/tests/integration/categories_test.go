package integration

import (
	"fmt"
	"testing"
)

type cat struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	DisplayOrder int    `json:"display_order"`
	Icon         string `json:"icon"`
}

func TestCategoriesCRUD(t *testing.T) {
	token := login(t, "admin")

	// CREATE
	created := apiPOST(t, token, "/categories", map[string]string{
		"name": "IT-Cat",
		"icon": "🧪",
	})
	assertStatus(t, created, 200)
	c := decode[cat](t, created)
	if c.Name != "IT-Cat" || c.Icon != "🧪" {
		t.Errorf("unexpected created category: %+v", c)
	}
	t.Logf("CREATE ok: id=%d name=%s order=%d", c.ID, c.Name, c.DisplayOrder)

	// READ (list)
	listResp := apiGET(t, token, "/categories")
	assertStatus(t, listResp, 200)
	list := decode[[]cat](t, listResp)
	found := false
	for _, cat := range list {
		if cat.ID == c.ID {
			found = true
			break
		}
	}
	if !found {
		t.Error("created category not found in listing")
	}
	t.Logf("READ ok: %d categories listed", len(list))

	// UPDATE
	updated := apiPATCH(t, token, "/categories/"+fmt.Sprintf("%d", c.ID), map[string]string{
		"name": "IT-Cat-renamed",
		"icon": "🔬",
	})
	assertStatus(t, updated, 200)
	u := decode[cat](t, updated)
	if u.Name != "IT-Cat-renamed" || u.Icon != "🔬" {
		t.Errorf("unexpected updated category: %+v", u)
	}
	t.Logf("UPDATE ok: renamed to %s", u.Name)

	// DELETE
	delResp := apiDELETE(t, token, "/categories/"+fmt.Sprintf("%d", c.ID))
	assertStatus(t, delResp, 204)
	t.Logf("DELETE ok")
}
