package integration

import (
	"fmt"
	"testing"
)

type tbl struct {
	ID       int64   `json:"id"`
	Name     string  `json:"name"`
	Capacity int     `json:"capacity"`
	X        float64 `json:"x"`
	Y        float64 `json:"y"`
}

func TestTablesCRUD(t *testing.T) {
	token := login(t, "admin")

	// CREATE
	created := apiPOST(t, token, "/tables", map[string]interface{}{
		"name":     "IT-T1",
		"capacity": 4,
		"x":        100,
		"y":        200,
	})
	assertStatus(t, created, 200)
	c := decode[tbl](t, created)
	if c.Name != "IT-T1" || c.Capacity != 4 {
		t.Errorf("unexpected created table: %+v", c)
	}
	if c.X != 100 || c.Y != 200 {
		t.Errorf("unexpected position: %.0f,%.0f", c.X, c.Y)
	}
	t.Logf("CREATE ok: id=%d name=%s", c.ID, c.Name)

	// READ (list)
	listResp := apiGET(t, token, "/tables")
	assertStatus(t, listResp, 200)
	list := decode[[]tbl](t, listResp)
	found := false
	for _, tb := range list {
		if tb.Name == "IT-T1" {
			found = true
			break
		}
	}
	if !found {
		t.Error("table IT-T1 not found in listing")
	}
	t.Logf("READ ok: %d tables listed", len(list))

	// UPDATE
	updated := apiPATCH(t, token, "/tables/"+fmt.Sprintf("%d", c.ID), map[string]interface{}{
		"name":     "IT-T1-renamed",
		"capacity": 6,
	})
	assertStatus(t, updated, 200)
	u := decode[tbl](t, updated)
	if u.Name != "IT-T1-renamed" || u.Capacity != 6 {
		t.Errorf("unexpected updated table: %+v", u)
	}
	t.Logf("UPDATE ok: renamed to %s capacity %d", u.Name, u.Capacity)

	// DELETE
	delResp := apiDELETE(t, token, "/tables/"+fmt.Sprintf("%d", c.ID))
	assertStatus(t, delResp, 204)
	t.Logf("DELETE ok")
}
