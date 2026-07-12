package integration

import (
	"net/http"
	"testing"
)

type apiReq struct {
	method string
	path   string
	body   interface{}
}

func TestAPIWithoutAuth(t *testing.T) {
	tt := []apiReq{
		{"GET", "/tables", nil},
		{"POST", "/tables", map[string]string{"name": "Test"}},
		{"GET", "/categories", nil},
		{"GET", "/dishes", nil},
		{"GET", "/orders", nil},
		{"GET", "/kds/orders", nil},
		{"GET", "/audit-events", nil},
		{"GET", "/floor-plan", nil},
		{"GET", "/predictions", nil},
	}

	for _, ep := range tt {
		var resp *http.Response
		switch ep.method {
		case "GET":
			resp = apiGET(t, "", ep.path)
		case "POST":
			resp = apiPOST(t, "", ep.path, ep.body)
		case "PATCH":
			resp = apiPATCH(t, "", ep.path, ep.body)
		case "DELETE":
			resp = apiDELETE(t, "", ep.path)
		}
		if resp.StatusCode != 401 {
			t.Errorf("%s %s: expected 401, got %d", ep.method, ep.path, resp.StatusCode)
		}
	}
}

func TestAuthenticatedAccessWorks(t *testing.T) {
	token := login(t, "waiter")

	// All these should be accessible to any authenticated user
	paths := []string{
		"/tables", "/categories", "/dishes", "/menu",
		"/floor-plan", "/orders", "/chef-suggestions",
		"/predictions", "/notifications", "/kds/orders",
	}
	for _, p := range paths {
		resp := apiGET(t, token, p)
		if resp.StatusCode == 401 {
			t.Errorf("waiter should access %s, got 401", p)
		}
	}
}
