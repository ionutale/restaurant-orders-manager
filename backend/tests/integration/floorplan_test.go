package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
)

type loginResponse struct {
	Token string `json:"token"`
}

type tableReq struct {
	Name     string `json:"name"`
	Capacity int    `json:"capacity"`
}

type tableResp struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Capacity int    `json:"capacity"`
}

func TestCreateFloorPlan(t *testing.T) {
	baseURL := os.Getenv("API_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080/api"
	}

	// 1. Login as admin
	loginBody, _ := json.Marshal(map[string]string{
		"email":    "admin@restaurant.com",
		"password": "admin",
	})
	resp, err := http.Post(baseURL+"/auth/login", "application/json", bytes.NewReader(loginBody))
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("login status: %d", resp.StatusCode)
	}

	var lr loginResponse
	json.NewDecoder(resp.Body).Decode(&lr)
	resp.Body.Close()

	if lr.Token == "" {
		t.Fatal("no token returned")
	}

	// 2. Create 20 tables totaling 60 seats
	tables := []tableReq{
		{Name: "A1", Capacity: 4}, {Name: "A2", Capacity: 4},
		{Name: "A3", Capacity: 4}, {Name: "A4", Capacity: 4},
		{Name: "A5", Capacity: 4}, {Name: "A6", Capacity: 4},
		{Name: "A7", Capacity: 2}, {Name: "A8", Capacity: 2},
		{Name: "A9", Capacity: 2}, {Name: "A10", Capacity: 2},
		{Name: "B1", Capacity: 4}, {Name: "B2", Capacity: 4},
		{Name: "B3", Capacity: 4}, {Name: "B4", Capacity: 4},
		{Name: "B5", Capacity: 2}, {Name: "B6", Capacity: 2},
		{Name: "B7", Capacity: 2}, {Name: "B8", Capacity: 2},
		{Name: "B9", Capacity: 2}, {Name: "B10", Capacity: 2},
	}
	// 10 × 4 = 40, 10 × 2 = 20 → 20 tables, 60 seats

	var created []tableResp
	totalSeats := 0

	for _, tbl := range tables {
		body, _ := json.Marshal(tbl)
		req, _ := http.NewRequest("POST", baseURL+"/tables", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+lr.Token)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("create table %s failed: %v", tbl.Name, err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("create table %s status: %d", tbl.Name, resp.StatusCode)
		}

		var tr tableResp
		json.NewDecoder(resp.Body).Decode(&tr)
		resp.Body.Close()

		if tr.Name != tbl.Name {
			t.Errorf("expected name %s, got %s", tbl.Name, tr.Name)
		}
		created = append(created, tr)
		totalSeats += tr.Capacity
	}

	if len(created) != 20 {
		t.Errorf("expected 20 tables, got %d", len(created))
	}
	if totalSeats != 60 {
		t.Errorf("expected 60 total seats, got %d", totalSeats)
	}

	// 3. Verify via GET /tables
	req, _ := http.NewRequest("GET", baseURL+"/tables", nil)
	req.Header.Set("Authorization", "Bearer "+lr.Token)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("list tables failed: %v", err)
	}

	var listed []tableResp
	json.NewDecoder(resp.Body).Decode(&listed)
	resp.Body.Close()

	seatCount := 0
	found := map[string]bool{}
	for _, t := range listed {
		seatCount += t.Capacity
		found[t.Name] = true
	}

	for _, tbl := range tables {
		if !found[tbl.Name] {
			t.Errorf("table %s not found in listing", tbl.Name)
		}
	}

	if seatCount < 60 {
		t.Errorf("expected at least 60 seats in listing, got %d", seatCount)
	}

	t.Logf("SUCCESS: Created %d tables with %d total seats", len(created), totalSeats)
	_ = fmt.Sprintf  // avoid unused import warning
}
