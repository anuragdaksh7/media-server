package errors

import (
	stderrors "errors"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestNewAndError(t *testing.T) {
	err := New(http.StatusBadRequest, "BAD_REQUEST", "bad request")

	if err.Status != http.StatusBadRequest {
		t.Fatalf("unexpected status: got %d, want %d", err.Status, http.StatusBadRequest)
	}

	if err.Code != "BAD_REQUEST" || err.Message != "bad request" {
		t.Fatalf("unexpected error payload: %+v", err)
	}

	if err.Error() != "bad request" {
		t.Fatalf("unexpected Error() value: got %q", err.Error())
	}
}

func TestNewWithDetails(t *testing.T) {
	details := gin.H{"field": "email"}
	err := NewWithDetails(
		http.StatusUnprocessableEntity,
		"VALIDATION_ERROR",
		"validation failed",
		details,
	)

	if err.Details == nil {
		t.Fatalf("expected details to be present")
	}
}

func TestRespondWithAPIError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/", func(c *gin.Context) {
		Respond(c, NewWithDetails(
			http.StatusForbidden,
			"FORBIDDEN",
			"forbidden",
			gin.H{"resource": "file"},
		))
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("unexpected status: got %d, want %d", rec.Code, http.StatusForbidden)
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

	if errorData["code"] != "FORBIDDEN" {
		t.Fatalf("unexpected code: got %v", errorData["code"])
	}
}

func TestRespondWithGenericError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/", func(c *gin.Context) {
		Respond(c, stderrors.New("boom"))
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("unexpected status: got %d, want %d", rec.Code, http.StatusInternalServerError)
	}

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to parse response body: %v", err)
	}

	errorData, ok := body["error"].(map[string]any)
	if !ok {
		t.Fatalf("expected error object in response body")
	}

	if errorData["code"] != "INTERNAL_SERVER_ERROR" {
		t.Fatalf("unexpected code: got %v", errorData["code"])
	}
}
