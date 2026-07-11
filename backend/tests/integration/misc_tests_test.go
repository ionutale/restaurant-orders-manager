package integration

import (
	"fmt"
	"testing"
)

func TestMenuEndpoint(t *testing.T) {
	token := login(t, "waiter")
	resp := apiGET(t, token, "/menu")
	assertStatus(t, resp, 200)
	type menuResp struct {
		Categories  []interface{} `json:"categories"`
		Suggestions []interface{} `json:"suggestions"`
		Allergens   []interface{} `json:"allergens"`
	}
	m := decode[menuResp](t, resp)
	if len(m.Categories) == 0 {
		t.Error("menu has no categories")
	}
	if len(m.Allergens) == 0 {
		t.Error("menu has no allergens")
	}
	t.Logf("MENU ok: %d categories, %d suggestions, %d allergens",
		len(m.Categories), len(m.Suggestions), len(m.Allergens))
}

func TestPredictionsEndpoint(t *testing.T) {
	token := login(t, "waiter")
	resp := apiGET(t, token, "/predictions")
	assertStatus(t, resp, 200)
	type pred struct {
		TableID int64 `json:"table_id"`
	}
	list := decode[[]pred](t, resp)
	t.Logf("PREDICTIONS ok: %d predictions", len(list))
}

func TestNotificationsEndpoint(t *testing.T) {
	token := login(t, "waiter")
	resp := apiGET(t, token, "/notifications")
	assertStatus(t, resp, 200)
	type notif struct {
		ItemID int64 `json:"item_id"`
	}
	list := decode[[]notif](t, resp)
	t.Logf("NOTIFICATIONS ok: %d notifications", len(list))
}

func TestAuditEndpoint(t *testing.T) {
	token := login(t, "admin")
	resp := apiGET(t, token, "/audit-events")
	assertStatus(t, resp, 200)
	type ev struct {
		ID int64 `json:"id"`
	}
	list := decode[[]ev](t, resp)
	if len(list) == 0 {
		t.Error("audit log is empty — expected at least some events")
	}
	t.Logf("AUDIT ok: %d events (latest: id=%d)", len(list), list[0].ID)
}

func TestDishGetEndpoint(t *testing.T) {
	token := login(t, "admin")
	dResp := apiGET(t, token, "/menu")
	type dishBrief struct {
		ID int64 `json:"id"`
	}
	type cat struct {
		Dishes []dishBrief `json:"dishes"`
	}
	type menu struct {
		Categories []cat `json:"categories"`
	}
	m := decode[menu](t, dResp)
	if len(m.Categories) == 0 || len(m.Categories[0].Dishes) == 0 {
		t.Skip("no dishes to test single GET")
	}
	dishID := m.Categories[0].Dishes[0].ID
	resp := apiGET(t, token, "/dishes/"+fmt.Sprintf("%d", dishID))
	assertStatus(t, resp, 200)
	type dish struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	}
	d := decode[dish](t, resp)
	if d.ID != dishID {
		t.Errorf("unexpected dish: %+v", d)
	}
	t.Logf("DISH GET ok: %s (id=%d)", d.Name, d.ID)
}
