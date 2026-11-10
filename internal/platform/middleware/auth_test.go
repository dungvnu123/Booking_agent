package middleware

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http/httptest"
	"testing"
)

type fakeVerifier struct {
	principal Principal
	err       error
}

func (v fakeVerifier) Verify(context.Context, string) (Principal, error) { return v.principal, v.err }
func TestAuthAndRolesFailClosed(t *testing.T) {
	gin.SetMode("test")
	for _, tc := range []struct {
		name, header string
		verifier     TokenVerifier
		status       int
	}{
		{"no verifier", "Bearer x", nil, 401},
		{"missing header", "", fakeVerifier{}, 401},
		{"bad token", "Bearer x", fakeVerifier{err: errors.New("revoked")}, 401},
		{"admin has no implicit doctor access", "Bearer x", fakeVerifier{principal: Principal{"u", []string{"ADMIN"}}}, 403},
		{"doctor", "Bearer x", fakeVerifier{principal: Principal{"u", []string{"DOCTOR"}}}, 200},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := gin.New()
			r.GET("/protected", Authenticate(tc.verifier), RequireRoles("DOCTOR"), func(c *gin.Context) { c.Status(200) })
			req := httptest.NewRequest("GET", "/protected", nil)
			req.Header.Set("Authorization", tc.header)
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)
			if rec.Code != tc.status {
				t.Fatalf("got %d want %d", rec.Code, tc.status)
			}
		})
	}
}
func TestCORS(t *testing.T) {
	gin.SetMode("test")
	for _, tc := range []struct {
		origin string
		status int
	}{{"http://localhost:5173", 204}, {"https://other.example", 403}} {
		r := gin.New()
		r.Use(CORS([]string{"http://localhost:5173"}))
		req := httptest.NewRequest("OPTIONS", "/api/patients", nil)
		req.Header.Set("Origin", tc.origin)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		if rec.Code != tc.status {
			t.Fatalf("%s: %d", tc.origin, rec.Code)
		}
	}
}
