package integration

import (
	"fmt"
	"testing"
)

type dish struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	PriceCents    int    `json:"price_cents"`
	CategoryID    int64  `json:"category_id"`
	EatingTimeMin int    `json:"eating_time_min"`
}

type allergen struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Icon string `json:"icon"`
}

func TestDishesCRUD(t *testing.T) {
	token := login(t, "admin")

	// Get first category for dish creation
	catResp := apiGET(t, token, "/categories")
	cats := decode[[]cat](t, catResp)
	if len(cats) == 0 {
		t.Fatal("no categories exist")
	}
	catID := cats[0].ID
	t.Logf("Using category %s (id=%d)", cats[0].Name, catID)

	// CREATE dish
	created := apiPOST(t, token, "/dishes", map[string]interface{}{
		"name":           "IT-Dish",
		"description":    "Integration test dish",
		"price_cents":    1500,
		"category_id":    catID,
		"eating_time_min": 10,
	})
	assertStatus(t, created, 200)
	d := decode[dish](t, created)
	if d.Name != "IT-Dish" || d.PriceCents != 1500 || d.CategoryID != catID {
		t.Errorf("unexpected created dish: %+v", d)
	}
	t.Logf("CREATE ok: id=%d name=%s", d.ID, d.Name)

	// Allergens — get list and assign one
	allResp := apiGET(t, token, "/allergens")
	allergens := decode[[]allergen](t, allResp)
	if len(allergens) > 0 {
		setResp := apiPUT(t, token, "/dishes/"+fmt.Sprintf("%d", d.ID)+"/allergens", map[string]interface{}{
			"allergen_ids": []int64{allergens[0].ID},
		})
		assertStatus(t, setResp, 200)
		t.Logf("ALLERGEN assign ok: %s", allergens[0].Name)

		// READ dish allergens
		daResp := apiGET(t, token, "/dishes/"+fmt.Sprintf("%d", d.ID)+"/allergens")
		assertStatus(t, daResp, 200)
		da := decode[[]allergen](t, daResp)
		if len(da) != 1 || da[0].ID != allergens[0].ID {
			t.Errorf("unexpected dish allergens: %+v", da)
		}
		t.Logf("ALLERGEN read ok: %d assigned", len(da))
	}

	// UPDATE dish
	updated := apiPATCH(t, token, "/dishes/"+fmt.Sprintf("%d", d.ID), map[string]interface{}{
		"name":        "IT-Dish-renamed",
		"price_cents": 2000,
	})
	assertStatus(t, updated, 200)
	u := decode[dish](t, updated)
	if u.Name != "IT-Dish-renamed" || u.PriceCents != 2000 {
		t.Errorf("unexpected updated dish: %+v", u)
	}
	t.Logf("UPDATE ok: %s €%.2f", u.Name, float64(u.PriceCents)/100)

	// DELETE dish
	delResp := apiDELETE(t, token, "/dishes/"+fmt.Sprintf("%d", d.ID))
	assertStatus(t, delResp, 204)
	t.Logf("DELETE ok")
}
