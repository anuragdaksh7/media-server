package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"fileserver/internal/auth/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func testRouterWithJWTMiddleware(t *testing.T) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(JWTMiddleware())
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"user_id": c.GetString("user_id"),
			"role":    c.GetString("role"),
		})
	})

	return router
}

func TestJWTMiddlewareRejectsMissingAuthorizationHeader(t *testing.T) {
	router := testRouterWithJWTMiddleware(t)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("unexpected status: got %d, want %d", rec.Code, http.StatusUnauthorized)
	}

	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to parse response body: %v", err)
	}

	if body["error"] != "missing auth header" {
		t.Fatalf("unexpected error message: %q", body["error"])
	}
}

func TestJWTMiddlewareRejectsInvalidToken(t *testing.T) {
	router := testRouterWithJWTMiddleware(t)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("unexpected status: got %d, want %d", rec.Code, http.StatusUnauthorized)
	}

	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to parse response body: %v", err)
	}

	if body["error"] != "invalid token" {
		t.Fatalf("unexpected error message: %q", body["error"])
	}
}

func TestJWTMiddlewareSetsUserContextForValidToken(t *testing.T) {
	originalSecret := utils.JWTSecret
	utils.JWTSecret = []byte("test-secret")
	t.Cleanup(func() { utils.JWTSecret = originalSecret })

	router := testRouterWithJWTMiddleware(t)

	userID := uuid.New()
	token, err := utils.GenerateJWT(userID, "admin")
	if err != nil {
		t.Fatalf("GenerateJWT returned error: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status: got %d, want %d", rec.Code, http.StatusOK)
	}

	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to parse response body: %v", err)
	}

	if body["user_id"] != userID.String() {
		t.Fatalf("unexpected user_id: got %q, want %q", body["user_id"], userID.String())
	}

	if body["role"] != "admin" {
		t.Fatalf("unexpected role: got %q, want %q", body["role"], "admin")
	}
}
