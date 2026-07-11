package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
)

func apiBase() string {
	if u := os.Getenv("API_URL"); u != "" {
		return u
	}
	return "http://localhost:8080/api"
}

func login(t testT, role string) string {
	creds := map[string]string{}
	switch role {
	case "admin":
		creds = map[string]string{"email": "admin@restaurant.com", "password": "admin"}
	case "waiter":
		creds = map[string]string{"email": "waiter@restaurant.com", "password": "waiter"}
	case "chef":
		creds = map[string]string{"email": "chef@restaurant.com", "password": "chef"}
	}
	body, _ := json.Marshal(creds)
	resp, err := http.Post(apiBase()+"/auth/login", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("login as %s failed: %v", role, err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("login as %s status: %d", role, resp.StatusCode)
	}
	var lr struct{ Token string }
	json.NewDecoder(resp.Body).Decode(&lr)
	resp.Body.Close()
	return lr.Token
}

type testT interface {
	Fatalf(format string, args ...interface{})
	Logf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

func apiGET(t testT, token, path string) *http.Response {
	req, _ := http.NewRequest("GET", apiBase()+path, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("GET %s: %v", path, err)
	}
	return resp
}

func apiPOST(t testT, token, path string, body interface{}) *http.Response {
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", apiBase()+path, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("POST %s: %v", path, err)
	}
	return resp
}

func apiPATCH(t testT, token, path string, body interface{}) *http.Response {
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest("PATCH", apiBase()+path, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("PATCH %s: %v", path, err)
	}
	return resp
}

func apiPUT(t testT, token, path string, body interface{}) *http.Response {
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest("PUT", apiBase()+path, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("PUT %s: %v", path, err)
	}
	return resp
}

func apiDELETE(t testT, token, path string) *http.Response {
	req, _ := http.NewRequest("DELETE", apiBase()+path, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("DELETE %s: %v", path, err)
	}
	return resp
}

func decode[T any](t testT, resp *http.Response) T {
	var v T
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		t.Fatalf("decode: %v", err)
	}
	resp.Body.Close()
	return v
}

func assertStatus(t testT, resp *http.Response, expected int) {
	if resp.StatusCode != expected {
		body := decode[map[string]interface{}](t, resp)
		t.Errorf("expected status %d, got %d: %v", expected, resp.StatusCode, body)
	}
}
