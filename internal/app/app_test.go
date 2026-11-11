package app

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"

	"abclinic.local/backend/api"
	"abclinic.local/backend/internal/config"
)

func testApp(t *testing.T) *App {
	t.Helper()
	t.Setenv("DATABASE_URL", "")
	t.Setenv("GIN_MODE", "test")
	t.Setenv("JWT_ACCESS_SECRET", "")
	t.Setenv("JWT_REFRESH_SECRET", "")
	cfg, err := config.FromEnv()
	if err != nil {
		t.Fatal(err)
	}
	a, err := New(context.Background(), cfg)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(a.Close)
	return a
}
func TestScaffoldRoutesAndDocumentation(t *testing.T) {
	router := testApp(t).Router()
	var spec struct {
		Paths map[string]map[string]json.RawMessage `json:"paths"`
	}
	if err := json.Unmarshal(api.Document, &spec); err != nil {
		t.Fatal(err)
	}
	count := 0
	for _, route := range router.Routes() {
		if !strings.HasPrefix(route.Path, "/api/") || route.Path == "/api/swagger-ui.html" || route.Path == "/api/v3/api-docs" {
			continue
		}
		count++
		docPath := route.Path
		requestPath := route.Path
		for _, part := range strings.Split(route.Path, "/") {
			if strings.HasPrefix(part, ":") {
				docPath = strings.ReplaceAll(docPath, part, "{"+part[1:]+"}")
				requestPath = strings.ReplaceAll(requestPath, part, "1")
			}
		}
		if _, ok := spec.Paths[docPath][strings.ToLower(route.Method)]; !ok {
			t.Errorf("missing OpenAPI route: %s %s", route.Method, docPath)
		}
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, httptest.NewRequest(route.Method, requestPath, nil))
		if route.Path == "/api/auth/verify-email" && route.Method == "GET" {
			if rec.Code != 200 || !strings.Contains(rec.Header().Get("Content-Type"), "text/html") {
				t.Fatal("verification page not served")
			}
			continue
		}
		expectedStatus, expectedCode := 501, 1501
		if strings.HasPrefix(route.Path, "/api/auth/") {
			expectedStatus, expectedCode = 400, 400
		}
		if rec.Code != expectedStatus {
			t.Errorf("%s %s: got %d, want %d", route.Method, requestPath, rec.Code, expectedStatus)
		}
		var response struct {
			StatusCode int `json:"statusCode"`
		}
		if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil || response.StatusCode != expectedCode {
			t.Fatalf("invalid scaffold response: %s", rec.Body.String())
		}
	}
	if count != 111 {
		t.Fatalf("business routes=%d; verify source inventory", count)
	}
	for _, path := range []string{"/api/patients/1/summary", "/api/patients/1/ai-sessions", "/api/v1/users"} {
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, httptest.NewRequest("GET", path, nil))
		if rec.Code != 404 {
			t.Errorf("excluded/demo route still exists: %s", path)
		}
	}
}
func TestHealthDocsAndReadinessWithoutDatabase(t *testing.T) {
	router := testApp(t).Router()
	for _, tc := range []struct {
		path   string
		status int
		marker string
	}{
		{"/healthz", 200, "ok"},
		{"/readyz", 503, "database_not_configured"},
		{"/swagger/index.html", 200, "ABClinic"},
		{"/swagger/swagger-ui-bundle.js", 200, "SwaggerUIBundle"},
		{"/swagger/swagger-ui.css", 200, "swagger-ui"},
		{"/swagger/swagger-initializer.js", 200, "openapi.json"},
		{"/swagger/openapi.json", 200, "ABClinic"},
		{"/api/v3/api-docs", 200, "ABClinic"},
	} {
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, httptest.NewRequest("GET", tc.path, nil))
		if rec.Code != tc.status || !strings.Contains(rec.Body.String(), tc.marker) {
			t.Errorf("%s: status %d", tc.path, rec.Code)
		}
		if rec.Header().Get("X-Request-ID") == "" {
			t.Error("missing request ID")
		}
	}
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest("GET", "/api/swagger-ui.html", nil))
	if rec.Code != 307 || rec.Header().Get("Location") != "/swagger/index.html" {
		t.Fatal("missing compatibility redirect")
	}
}
