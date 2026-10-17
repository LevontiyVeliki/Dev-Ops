package http_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-kit/kit/log"

	httptransport "example.com/restaurant-service/http"
	"example.com/restaurant-service/implementation"
	"example.com/restaurant-service/repository"
	"example.com/restaurant-service/transport"
)

func newServer() *httptest.Server {
	logger := log.NewNopLogger()
	svc := implementation.NewService(repository.NewMemoryRepository(), logger)
	return httptest.NewServer(httptransport.NewService(transport.MakeEndpoints(svc), logger))
}

func do(t *testing.T, method, url, body string) (int, map[string]interface{}) {
	t.Helper()
	req, _ := http.NewRequest(method, url, strings.NewReader(body))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	var out map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&out)
	return resp.StatusCode, out
}

func TestHTTPFlow(t *testing.T) {
	srv := newServer()
	defer srv.Close()

	code, out := do(t, "POST", srv.URL+"/restaurants",
		`{"name":"Шашлычная","address":"ул. Мира, 1","rating":4.5,"open":true}`)
	id, _ := out["id"].(string)
	if code != http.StatusCreated || id == "" {
		t.Fatalf("POST: %d %v", code, out)
	}

	code, out = do(t, "GET", srv.URL+"/restaurants/"+id, "")
	r, _ := out["restaurant"].(map[string]interface{})
	if code != http.StatusOK || r["name"] != "Шашлычная" || r["rating"] != 4.5 || r["open"] != true {
		t.Fatalf("GET: %d %v", code, out)
	}

	if code, out = do(t, "PATCH", srv.URL+"/restaurants/"+id+"/rating", `{"rating":4.9}`); code != http.StatusOK {
		t.Fatalf("PATCH: %d %v", code, out)
	}
	_, out = do(t, "GET", srv.URL+"/restaurants/"+id, "")
	if r := out["restaurant"].(map[string]interface{}); r["rating"].(float64) < 4.89 {
		t.Errorf("rating not updated: %v", out)
	}
}

func TestHTTPErrors(t *testing.T) {
	srv := newServer()
	defer srv.Close()
	cases := []struct {
		method, path, body string
		code               int
	}{
		{"GET", "/restaurants/unknown", "", http.StatusNotFound},
		{"PATCH", "/restaurants/unknown/rating", `{"rating":3}`, http.StatusNotFound},
		{"POST", "/restaurants", `{"name":""}`, http.StatusBadRequest},
		{"POST", "/restaurants", `{"name":"A","rating":9}`, http.StatusBadRequest},
		{"POST", "/restaurants", `not json`, http.StatusBadRequest},
		{"PATCH", "/restaurants/x/rating", `{}`, http.StatusBadRequest},
		{"DELETE", "/restaurants/x", "", http.StatusMethodNotAllowed},
	}
	for _, c := range cases {
		code, out := do(t, c.method, srv.URL+c.path, c.body)
		if code != c.code {
			t.Errorf("%s %s %s: %d %v, want %d", c.method, c.path, c.body, code, out, c.code)
		}
	}
}
