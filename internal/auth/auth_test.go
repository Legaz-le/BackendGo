package auth

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Setenv("JWT_SECRET", "testsecret")
	os.Exit(m.Run())
}

func TestPasswordHash(t *testing.T) {
	value, err := HashPassword("mypassword123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	passwordCheck := CheckPasswordHash("mypassword123", value)
	wrongCheck := CheckPasswordHash("wrongpassword123", value)

	if passwordCheck != nil {
		t.Fatalf("expected match, got error: %v", passwordCheck)
	}

	if wrongCheck == nil {
		t.Fatal("expected error for wrong password, got nil")
	}
}

func TestHashRefreshToken(t *testing.T) {
	value := HashRefreshToken("rawToken")
	secondValue := HashRefreshToken("rawToken")

	if value != secondValue {
		t.Fatal("expected same hash for same input")
	}

	checkValue := HashRefreshToken("differentHash")
	anotherValue := HashRefreshToken("notSameHash")

	if checkValue == anotherValue {
		t.Fatal("expected different hash for different input")
	}
}

func TestGenerateRefreshToken(t *testing.T) {
	rawToken, tokenHash, err := GenerateRefreshToken()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if rawToken == "" || tokenHash == "" {
		t.Fatal("expected non-empty  tokens")
	}

	if rawToken == tokenHash {
		t.Fatal("expected to be different tokens")
	}
}

func TestGenerateAndValidateAcessToken(t *testing.T) {
	value, err := GenerateAccessToken(1, "employer")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	check, err := ValidateAccessToken(value)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if check.UserID != 1 {
		t.Fatalf("expected  UserID 1, got  %d", check.UserID)
	}

	if check.Role != "employer" {
		t.Fatalf("expected role employer, got %s", check.Role)
	}

	_, err = ValidateAccessToken(value + "bad")
	if err == nil {
		t.Fatal("expected error for invalid token, got nil")
	}
}
