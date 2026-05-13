package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestSuccessWritesExpectedPayload(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/", func(c *gin.Context) {
		Success(c, http.StatusCreated, gin.H{"id": "123"})
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("unexpected status: got %d, want %d", rec.Code, http.StatusCreated)
	}

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to parse response body: %v", err)
	}

	if success, ok := body["success"].(bool); !ok || !success {
		t.Fatalf("expected success=true in response body")
	}

	data, ok := body["data"].(map[string]any)
	if !ok {
		t.Fatalf("expected data object in response body")
	}

	if data["id"] != "123" {
		t.Fatalf("unexpected data.id: got %v, want %v", data["id"], "123")
	}
}
