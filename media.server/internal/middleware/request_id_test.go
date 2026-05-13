package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TestRequestIDUsesProvidedHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(RequestID())
	router.GET("/", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Request-ID", "req-123")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if got := rec.Header().Get("X-Request-ID"); got != "req-123" {
		t.Fatalf("unexpected response request id: got %q, want %q", got, "req-123")
	}
}

func TestRequestIDGeneratesUUIDWhenHeaderMissing(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(RequestID())
	router.GET("/", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	generatedID := rec.Header().Get("X-Request-ID")
	if generatedID == "" {
		t.Fatalf("expected generated X-Request-ID header to be non-empty")
	}

	if _, err := uuid.Parse(generatedID); err != nil {
		t.Fatalf("expected generated X-Request-ID to be a valid UUID, got %q", generatedID)
	}
}
