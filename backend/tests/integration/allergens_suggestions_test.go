package integration

import (
	"fmt"
	"testing"
)

func TestAllergensDishSuggestions(t *testing.T) {
	token := login(t, "admin")

	// CREATE allergen
	created := apiPOST(t, token, "/allergens", map[string]string{
		"name": "IT-Allergen", "icon": "🧪",
	})
	assertStatus(t, created, 200)
	type al struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
		Icon string `json:"icon"`
	}
	a := decode[al](t, created)
	if a.Name != "IT-Allergen" {
		t.Errorf("unexpected allergen: %+v", a)
	}
	t.Logf("ALLERGEN create ok: id=%d", a.ID)

	// READ all allergens
	listResp := apiGET(t, token, "/allergens")
	assertStatus(t, listResp, 200)
	all := decode[[]al](t, listResp)
	found := false
	for _, x := range all {
		if x.ID == a.ID {
			found = true
			break
		}
	}
	if !found {
		t.Error("allergen not in list")
	}
	t.Logf("ALLERGEN list ok: %d total", len(all))

	// DELETE allergen
	delResp := apiDELETE(t, token, "/allergens/"+fmt.Sprintf("%d", a.ID))
	assertStatus(t, delResp, 204)
	t.Logf("ALLERGEN delete ok")

	// --- DISH SUGGESTIONS ---
	// Get a dish and a wine to create suggestions
	dResp := apiGET(t, token, "/menu")
	type menuDish struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	}
	type menuCat struct {
		Category struct{ ID int64 } `json:"category"`
		Dishes  []menuDish         `json:"dishes"`
	}
	type menuResp struct {
		Categories []menuCat `json:"categories"`
	}

	menu := decode[menuResp](t, dResp)
	if len(menu.Categories) < 2 {
		t.Fatal("need at least 2 categories for suggestion test")
	}

	dish1 := menu.Categories[0].Dishes[0]
	dish2 := menu.Categories[1].Dishes[0]

	// CREATE dish suggestion
	sugCreated := apiPOST(t, token, "/dishes/"+fmt.Sprintf("%d", dish1.ID)+"/suggestions", map[string]interface{}{
		"to_dish_id": dish2.ID, "suggestion_type": "side",
	})
	assertStatus(t, sugCreated, 200)
	type sug struct {
		ID int64 `json:"id"`
	}
	s := decode[sug](t, sugCreated)
	t.Logf("SUGGESTION create ok: id=%d", s.ID)

	// READ dish suggestions
	sugList := apiGET(t, token, "/dishes/"+fmt.Sprintf("%d", dish1.ID)+"/suggestions")
	assertStatus(t, sugList, 200)
	sList := decode[[]sug](t, sugList)
	if len(sList) == 0 {
		t.Error("no suggestions found")
	}
	t.Logf("SUGGESTION list ok: %d", len(sList))

	// DELETE suggestion
	delSug := apiDELETE(t, token, "/dish-suggestions/"+fmt.Sprintf("%d", s.ID))
	assertStatus(t, delSug, 204)
	t.Logf("SUGGESTION delete ok")
}
