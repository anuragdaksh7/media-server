package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRecoveryHandlesPanics(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(RequestID(), Recovery())
	router.GET("/panic", func(c *gin.Context) {
		panic("boom")
	})

	req := httptest.NewRequest(http.MethodGet, "/panic", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("unexpected status: got %d, want %d", rec.Code, http.StatusInternalServerError)
	}

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to parse response body: %v", err)
	}

	success, ok := body["success"].(bool)
	if !ok || success {
		t.Fatalf("expected success=false in response body")
	}

	errorData, ok := body["error"].(map[string]any)
	if !ok {
		t.Fatalf("expected error object in response body")
	}

	if code := errorData["code"]; code != "INTERNAL_SERVER_ERROR" {
		t.Fatalf("unexpected error code: got %v", code)
	}
}
