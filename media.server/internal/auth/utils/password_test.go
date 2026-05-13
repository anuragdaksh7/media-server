package utils

import "testing"

func TestHashPasswordAndCheckPassword(t *testing.T) {
	password := "StrongP@ssw0rd!"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}

	if hash == password {
		t.Fatalf("expected hash to differ from original password")
	}

	if err := CheckPassword(password, hash); err != nil {
		t.Fatalf("CheckPassword should succeed for correct password: %v", err)
	}
}

func TestCheckPasswordRejectsWrongPassword(t *testing.T) {
	hash, err := HashPassword("correct-password")
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}

	if err := CheckPassword("wrong-password", hash); err == nil {
		t.Fatalf("expected CheckPassword to fail for incorrect password")
	}
}
