package auth

import (
	"errors"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func TestSignAndParseRoundTrip(t *testing.T) {
	secret := "unit-test-secret"

	token, err := Sign(secret, time.Minute, "user-1", "user-1@example.com", []string{"PARENT", "SCHOOL_ADMIN"}, "school-1")
	if err != nil {
		t.Fatalf("Sign() error = %v", err)
	}

	claims, err := Parse(secret, token)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if claims.UserID != "user-1" {
		t.Fatalf("claims.UserID = %q, want %q", claims.UserID, "user-1")
	}
	if claims.Email != "user-1@example.com" {
		t.Fatalf("claims.Email = %q, want %q", claims.Email, "user-1@example.com")
	}
	if len(claims.Roles) != 2 {
		t.Fatalf("len(claims.Roles) = %d, want %d", len(claims.Roles), 2)
	}
	if claims.SchoolID != "school-1" {
		t.Fatalf("claims.SchoolID = %q, want %q", claims.SchoolID, "school-1")
	}
}

func TestParseRejectsInvalidTokens(t *testing.T) {
	tests := []struct {
		name     string
		secret   string
		tokenStr string
	}{
		{
			name:     "malformed token string",
			secret:   "secret",
			tokenStr: "not-a-jwt",
		},
		{
			name:     "token signed by different secret",
			secret:   "secret-b",
			tokenStr: mustSignToken(t, "secret-a", time.Minute),
		},
		{
			name:     "expired token",
			secret:   "secret-exp",
			tokenStr: mustSignToken(t, "secret-exp", -time.Second),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse(tt.secret, tt.tokenStr)
			if !errors.Is(err, ErrInvalidToken) {
				t.Fatalf("Parse() error = %v, want %v", err, ErrInvalidToken)
			}
		})
	}
}

func TestVerifyPassword(t *testing.T) {
	hash, err := bcrypt.GenerateFromPassword([]byte("correct-password"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("GenerateFromPassword() error = %v", err)
	}

	if !VerifyPassword(string(hash), "correct-password") {
		t.Fatalf("VerifyPassword() = false, want true")
	}
	if VerifyPassword(string(hash), "wrong-password") {
		t.Fatalf("VerifyPassword() = true, want false")
	}
}

func TestNewAuthenticatorAndSignToken(t *testing.T) {
	a := NewAuthenticator("auth-secret", 15)
	if a.TTLSeconds != 900 {
		t.Fatalf("TTLSeconds = %d, want %d", a.TTLSeconds, 900)
	}

	token, err := a.SignToken("user-auth", "user-auth@example.com", []string{"SCHOOL_ADMIN"}, "school-auth")
	if err != nil {
		t.Fatalf("SignToken() error = %v", err)
	}

	claims, err := Parse("auth-secret", token)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if claims.UserID != "user-auth" {
		t.Fatalf("claims.UserID = %q, want %q", claims.UserID, "user-auth")
	}
	if claims.SchoolID != "school-auth" {
		t.Fatalf("claims.SchoolID = %q, want %q", claims.SchoolID, "school-auth")
	}
}

func mustSignToken(t *testing.T, secret string, ttl time.Duration) string {
	t.Helper()
	token, err := Sign(secret, ttl, "u", "u@example.com", []string{"PARENT"}, "")
	if err != nil {
		t.Fatalf("Sign() error = %v", err)
	}
	return token
}
