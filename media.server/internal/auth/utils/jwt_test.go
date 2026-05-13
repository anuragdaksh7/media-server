package utils

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestGenerateAndParseJWT(t *testing.T) {
	originalSecret := JWTSecret
	JWTSecret = []byte("test-secret")
	t.Cleanup(func() { JWTSecret = originalSecret })

	userID := uuid.New()
	role := "admin"

	token, err := GenerateJWT(userID, role)
	if err != nil {
		t.Fatalf("GenerateJWT returned error: %v", err)
	}

	claims, err := ParseJWT(token)
	if err != nil {
		t.Fatalf("ParseJWT returned error: %v", err)
	}

	if claims.UserID != userID.String() {
		t.Fatalf("unexpected user_id: got %s, want %s", claims.UserID, userID.String())
	}

	if claims.Role != role {
		t.Fatalf("unexpected role: got %s, want %s", claims.Role, role)
	}

	if claims.ExpiresAt == nil || claims.ExpiresAt.Time.Before(time.Now()) {
		t.Fatalf("expected token to have a future expiration timestamp")
	}
}

func TestParseJWTRejectsMalformedToken(t *testing.T) {
	if _, err := ParseJWT("not-a-jwt-token"); err == nil {
		t.Fatalf("expected ParseJWT to fail for malformed token")
	}
}

func TestParseJWTRejectsTokenWithWrongSecret(t *testing.T) {
	originalSecret := JWTSecret
	t.Cleanup(func() { JWTSecret = originalSecret })

	userID := uuid.New()

	JWTSecret = []byte("secret-a")
	token, err := GenerateJWT(userID, "user")
	if err != nil {
		t.Fatalf("GenerateJWT returned error: %v", err)
	}

	JWTSecret = []byte("secret-b")
	if _, err := ParseJWT(token); err == nil {
		t.Fatalf("expected ParseJWT to fail with mismatched signing secret")
	}
}
