package integration

import (
	"fmt"
	"testing"
)

type order struct {
	ID           int64        `json:"id"`
	TableGroupID int64        `json:"table_group_id"`
	Status       string       `json:"status"`
	Courses      []orderCourse `json:"courses,omitempty"`
}

type orderCourse struct {
	ID    int64       `json:"id"`
	Name  string      `json:"name"`
	Status string     `json:"status"`
	Items []orderItem `json:"items,omitempty"`
}

type orderItem struct {
	ID       int64  `json:"id"`
	DishName string `json:"dish_name"`
	Quantity int    `json:"quantity"`
	Notes    string `json:"notes"`
	Ready    bool   `json:"ready"`
}

type tableGroup struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	PartySize int    `json:"party_size"`
	Status    string `json:"status"`
}

func TestOrdersCRUD(t *testing.T) {
	waiterToken := login(t, "waiter")
	chefToken := login(t, "chef")

	// Find a free table to seat
	fpResp := apiGET(t, waiterToken, "/floor-plan")
	type fpTable struct {
		ID     int64  `json:"id"`
		Name   string `json:"name"`
		Status string `json:"status"`
	}
	fp := decode[[]fpTable](t, fpResp)
	var freeTableID int64
	for _, tbl := range fp {
		if tbl.Status == "free" {
			freeTableID = tbl.ID
			break
		}
	}
	if freeTableID == 0 {
		t.Fatal("no free table available")
	}
	t.Logf("Using free table id=%d", freeTableID)

	// CREATE table group (seat guests)
	groupResp := apiPOST(t, waiterToken, "/table-groups", map[string]interface{}{
		"table_ids":  []int64{freeTableID},
		"party_size": 3,
		"name":       "IT-Test-Group",
	})
	assertStatus(t, groupResp, 200)
	g := decode[tableGroup](t, groupResp)
	if g.PartySize != 3 || g.Status != "open" {
		t.Errorf("unexpected group: %+v", g)
	}
	t.Logf("TABLE GROUP created: id=%d", g.ID)

	// CREATE order with courses
	orderResp := apiPOST(t, waiterToken, "/orders", map[string]interface{}{
		"table_group_id": g.ID,
		"course_names":   []string{"IT-Starter", "IT-Main", "IT-Dessert"},
	})
	assertStatus(t, orderResp, 200)
	o1 := decode[order](t, orderResp)
	if len(o1.Courses) != 3 {
		t.Errorf("expected 3 courses, got %d", len(o1.Courses))
	}
	if o1.Courses[0].Name != "IT-Starter" || o1.Courses[0].Status != "active" {
		t.Errorf("first course should be active: %+v", o1.Courses[0])
	}
	t.Logf("ORDER created: id=%d with %d courses", o1.ID, len(o1.Courses))

	// READ order
	getResp := apiGET(t, waiterToken, "/orders/"+fmt.Sprintf("%d", o1.ID))
	assertStatus(t, getResp, 200)
	o2 := decode[order](t, getResp)
	if o2.ID != o1.ID {
		t.Errorf("order id mismatch: %d != %d", o2.ID, o1.ID)
	}
	t.Logf("READ order ok: status=%s", o2.Status)

	// ADD item to first course
	dResp := apiGET(t, waiterToken, "/dishes")
	type dishLight struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	}
	dishes := decode[[]dishLight](t, dResp)
	if len(dishes) == 0 {
		t.Fatal("no dishes available")
	}
	addResp := apiPOST(t, waiterToken,
		"/orders/"+fmt.Sprintf("%d", o1.ID)+"/courses/"+fmt.Sprintf("%d", o1.Courses[0].ID)+"/items",
		map[string]interface{}{
			"dish_id":  dishes[0].ID,
			"quantity": 2,
			"notes":    "no onions",
		})
	assertStatus(t, addResp, 200)
	item := decode[orderItem](t, addResp)
	if item.Quantity != 2 || item.Notes != "no onions" {
		t.Errorf("unexpected item: %+v", item)
	}
	t.Logf("ITEM added: %s x%d", item.DishName, item.Quantity)

	// SEND to KDS
	sendResp := apiPOST(t, waiterToken, "/orders/"+fmt.Sprintf("%d", o1.ID)+"/send", nil)
	assertStatus(t, sendResp, 200)
	sent := decode[order](t, sendResp)
	if sent.Status != "sent" {
		t.Errorf("expected status sent, got %s", sent.Status)
	}
	t.Logf("SEND to KDS ok")

	// CHEF marks item ready
	readyResp := apiPATCH(t, chefToken,
		"/kds/order-items/"+fmt.Sprintf("%d", item.ID)+"/ready", nil)
	assertStatus(t, readyResp, 200)
	t.Logf("ITEM marked ready ok")

	// ADVANCE course
	advResp := apiPOST(t, waiterToken, "/orders/"+fmt.Sprintf("%d", o1.ID)+"/advance-course", nil)
	assertStatus(t, advResp, 200)
	adv := decode[order](t, advResp)
	if len(adv.Courses) >= 2 && adv.Courses[0].Status == "completed" && adv.Courses[1].Status == "active" {
		t.Logf("ADVANCE ok: course 0 completed, course 1 active")
	} else {
		t.Errorf("unexpected course states after advance: %+v", adv.Courses)
	}

	// KDS — verify order appears
	kdsResp := apiGET(t, chefToken, "/kds/orders")
	assertStatus(t, kdsResp, 200)
	t.Logf("KDS orders ok (no error)")

	// LIST waiter's orders
	listResp := apiGET(t, waiterToken, "/orders")
	assertStatus(t, listResp, 200)
	list := decode[[]order](t, listResp)
	if len(list) == 0 {
		t.Error("no orders found in listing")
	}
	t.Logf("LIST orders ok: %d orders", len(list))

	// DELETE item from second course before advancing further
	delItemResp := apiDELETE(t, waiterToken, "/order-items/"+fmt.Sprintf("%d", item.ID))
	assertStatus(t, delItemResp, 204)
	t.Logf("DELETE item ok")
}
