package integration

import (
	"fmt"
	"testing"
)

func TestForeignKeyConstraints(t *testing.T) {
	adminToken := login(t, "admin")

	// Delete a category that has dishes → should fail (500)
	// Category 1 (Appetizers) has dishes
	resp := apiDELETE(t, adminToken, "/categories/1")
	if resp.StatusCode == 204 {
		t.Error("deleting category with dishes should fail")
	}
	t.Logf("DELETE category with dishes → %d (expected error)", resp.StatusCode)

	// Delete a user who created orders → should fail or succeed gracefully
	// User 1 is admin, has many events
	resp = apiDELETE(t, adminToken, "/users/1")
	if resp.StatusCode == 204 {
		t.Error("deleting admin with audit events should fail")
	}
	t.Logf("DELETE admin user → %d (expected error)", resp.StatusCode)
}

func TestIdempotentOperations(t *testing.T) {
	token := login(t, "waiter")

	// Mark item ready twice — second should be idempotent
	// First create a sent order with an item
	fpResp := apiGET(t, token, "/floor-plan")
	type fpTbl struct{ ID int64 }
	fp := decode[[]fpTbl](t, fpResp)
	var freeID int64
	for _, tbl := range fp {
		freeID = tbl.ID
		break
	}
	if freeID == 0 {
		t.Skip("no tables")
	}

	g := decode[struct{ ID int64 }](t, apiPOST(t, token, "/start-order", map[string]interface{}{
		"table_ids": []int64{freeID}, "party_size": 2,
	}))

	menu := decode[struct {
		Categories []struct {
			Dishes []struct{ ID int64 } `json:"dishes"`
		} `json:"categories"`
	}](t, apiGET(t, token, "/menu"))

	if len(menu.Categories) == 0 || len(menu.Categories[0].Dishes) == 0 {
		t.Skip("no dishes")
	}
	dishID := menu.Categories[0].Dishes[0].ID

	o := decode[struct {
		ID      int64 `json:"id"`
		Courses []struct {
			ID int64 `json:"id"`
		} `json:"courses"`
	}](t, apiGET(t, token, "/orders/"+fmt.Sprintf("%d", g.ID)))

	if len(o.Courses) == 0 {
		t.Skip("no courses")
	}

	item := decode[struct {
		ID int64 `json:"id"`
	}](t, apiPOST(t, token,
		"/orders/"+fmt.Sprintf("%d", o.ID)+"/courses/"+fmt.Sprintf("%d", o.Courses[0].ID)+"/items",
		map[string]interface{}{"dish_id": dishID, "quantity": 1}))

	apiPOST(t, token, "/orders/"+fmt.Sprintf("%d", o.ID)+"/send", nil)

	chefToken := login(t, "chef")
	resp1 := apiPATCH(t, chefToken, "/kds/order-items/"+fmt.Sprintf("%d", item.ID)+"/ready", nil)
	assertStatus(t, resp1, 200)

	resp2 := apiPATCH(t, chefToken, "/kds/order-items/"+fmt.Sprintf("%d", item.ID)+"/ready", nil)
	assertStatus(t, resp2, 200)
	t.Log("Mark ready twice is idempotent")

	// Close the group to free the table for subsequent tests
	closeGroupForOrder(t, token, o.ID)
}

func TestPartialUpdate(t *testing.T) {
	token := login(t, "admin")

	// PATCH with single field should preserve others
	created := apiPOST(t, token, "/tables", map[string]interface{}{
		"name": "PartialTest", "capacity": 4, "x": 100, "y": 200,
	})
	c := decode[struct {
		ID       int64   `json:"id"`
		Name     string  `json:"name"`
		Capacity int     `json:"capacity"`
		X        float64 `json:"x"`
		Y        float64 `json:"y"`
	}](t, created)

	// Update only name
	updated := apiPATCH(t, token, "/tables/"+fmt.Sprintf("%d", c.ID), map[string]string{
		"name": "PartialTestRenamed",
	})
	u := decode[struct {
		Name     string  `json:"name"`
		Capacity int     `json:"capacity"`
		X        float64 `json:"x"`
		Y        float64 `json:"y"`
	}](t, updated)

	if u.Name != "PartialTestRenamed" {
		t.Errorf("name not updated: %s", u.Name)
	}
	if u.Capacity != 4 {
		t.Errorf("capacity changed: %d", u.Capacity)
	}
	if u.X != 100 || u.Y != 200 {
		t.Errorf("position changed: %.0f,%.0f", u.X, u.Y)
	}
	t.Log("Partial update preserves other fields")

	// Cleanup
	apiDELETE(t, token, "/tables/"+fmt.Sprintf("%d", c.ID))
}

func TestDefaultOrderCourse(t *testing.T) {
	token := login(t, "waiter")

	// Create order without course_names → should get one default course
	fpResp := apiGET(t, token, "/floor-plan")
	type fpTbl struct{ ID int64 }
	fp := decode[[]fpTbl](t, fpResp)
	var freeID int64
	for _, tbl := range fp {
		freeID = tbl.ID
		break
	}
	if freeID == 0 {
		t.Skip("no free tables")
	}

	g := decode[struct{ ID int64 }](t, apiPOST(t, token, "/start-order", map[string]interface{}{
		"table_ids": []int64{freeID}, "party_size": 2,
	}))

	o := decode[struct {
		ID      int64 `json:"id"`
		Courses []struct {
			ID   int64  `json:"id"`
			Name string `json:"name"`
		} `json:"courses"`
	}](t, apiGET(t, token, "/orders/"+fmt.Sprintf("%d", g.ID)))

	if len(o.Courses) != 1 {
		t.Errorf("expected 1 default course, got %d", len(o.Courses))
	}
	if o.Courses[0].Name != "Course 1" {
		t.Errorf("expected 'Course 1', got '%s'", o.Courses[0].Name)
	}
	t.Logf("Default order: 1 course '%s'", o.Courses[0].Name)

	// Cleanup: close the group to free the table
	// The group ID is in the order's table_group_id
	// We need to close it via the table-groups endpoint
	type orderWithGroup struct {
		ID           int64 `json:"id"`
		TableGroupID int64 `json:"table_group_id"`
	}
	o2 := decode[orderWithGroup](t, apiGET(t, token, "/orders/"+fmt.Sprintf("%d", g.ID)))
	if o2.TableGroupID != 0 {
		apiPOST(t, token, "/table-groups/"+fmt.Sprintf("%d", o2.TableGroupID)+"/close", nil)
		t.Logf("Cleaned up group %d", o2.TableGroupID)
	}
}
