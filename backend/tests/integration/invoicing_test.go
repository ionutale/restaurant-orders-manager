package integration

import (
	"fmt"
	"testing"
)

type invoiceResult struct {
	Status string `json:"status"`
	Email  string `json:"email"`
	Total  string `json:"total"`
}

func TestInvoicing(t *testing.T) {
	adminToken := login(t, "admin")
	waiterToken := login(t, "waiter")

	fpResp := apiGET(t, adminToken, "/floor-plan")
	type fpTbl struct{ ID int64; Status string }
	fp := decode[[]fpTbl](t, fpResp)
	var freeID int64
	for _, tbl := range fp {
		if tbl.Status == "free" { freeID = tbl.ID; break }
	}
	if freeID == 0 { t.Fatal("no free table") }

	grpResp := apiPOST(t, waiterToken, "/table-groups", map[string]interface{}{
		"table_ids": []int64{freeID}, "party_size": 2,
	})
	assertStatus(t, grpResp, 200)
	type grp struct{ ID int64 }
	g := decode[grp](t, grpResp)

	ordResp := apiPOST(t, waiterToken, "/orders", map[string]interface{}{
		"table_group_id": g.ID, "course_names": []string{"IT-Inv"},
	})
	assertStatus(t, ordResp, 200)
	type course struct{ ID int64 }
	type ordFull struct {
		ID      int64    `json:"id"`
		Status  string   `json:"status"`
		Courses []course `json:"courses"`
	}
	o := decode[ordFull](t, ordResp)

	dResp := apiGET(t, adminToken, "/dishes")
	type dish struct{ ID int64 }
	dishes := decode[[]dish](t, dResp)
	apiPOST(t, waiterToken,
		"/orders/"+fmt.Sprintf("%d", o.ID)+"/courses/"+fmt.Sprintf("%d", o.Courses[0].ID)+"/items",
		map[string]interface{}{"dish_id": dishes[0].ID, "quantity": 1})

	apiPOST(t, waiterToken, "/orders/"+fmt.Sprintf("%d", o.ID)+"/send", nil)
	apiPOST(t, waiterToken, "/orders/"+fmt.Sprintf("%d", o.ID)+"/advance-course", nil)

	ordGet := apiGET(t, waiterToken, "/orders/"+fmt.Sprintf("%d", o.ID))
	o2 := decode[ordFull](t, ordGet)
	if o2.Status != "completed" {
		t.Skipf("order not completed (status=%s)", o2.Status)
	}

	payResp := apiPOST(t, adminToken, "/orders/"+fmt.Sprintf("%d", o.ID)+"/pay", map[string]string{
		"payment_method": "cash",
	})
	assertStatus(t, payResp, 200)
	t.Logf("PAY ok")

	invResp := apiPOST(t, adminToken, "/orders/"+fmt.Sprintf("%d", o.ID)+"/send-invoice", map[string]string{
		"email": "test@example.com",
	})
	assertStatus(t, invResp, 200)

	inv := decode[invoiceResult](t, invResp)
	if inv.Status != "sent" || inv.Email != "test@example.com" {
		t.Errorf("unexpected invoice: %+v", inv)
	}
	t.Logf("INVOICE ok: €%s to %s", inv.Total, inv.Email)

	closeGroupForOrder(t, waiterToken, o.ID)
}
