package auth

import (
	"context"
	"errors"
	"testing"
)

func TestStringClaim(t *testing.T) {
	claims := map[string]interface{}{
		"email": "user@example.com",
		"n":     123,
	}

	if got := stringClaim(claims, "email"); got != "user@example.com" {
		t.Fatalf("stringClaim(email) = %q, want %q", got, "user@example.com")
	}
	if got := stringClaim(claims, "n"); got != "" {
		t.Fatalf("stringClaim(non-string) = %q, want empty", got)
	}
	if got := stringClaim(claims, "missing"); got != "" {
		t.Fatalf("stringClaim(missing) = %q, want empty", got)
	}
}

func TestBoolClaim(t *testing.T) {
	claims := map[string]interface{}{
		"email_verified": true,
		"s":              "true",
	}

	if got := boolClaim(claims, "email_verified"); !got {
		t.Fatalf("boolClaim(email_verified) = false, want true")
	}
	if got := boolClaim(claims, "s"); got {
		t.Fatalf("boolClaim(non-bool) = true, want false")
	}
	if got := boolClaim(claims, "missing"); got {
		t.Fatalf("boolClaim(missing) = true, want false")
	}
}

func TestGoogleIDTokenVerifierVerifyRejectsMissingInputs(t *testing.T) {
	v := &googleIDTokenVerifier{clientID: ""}
	_, err := v.Verify(context.Background(), "some-token")
	if !errors.Is(err, ErrGoogleTokenInvalid) {
		t.Fatalf("error = %v, want %v", err, ErrGoogleTokenInvalid)
	}

	v = &googleIDTokenVerifier{clientID: "client-id"}
	_, err = v.Verify(context.Background(), "")
	if !errors.Is(err, ErrGoogleTokenInvalid) {
		t.Fatalf("error = %v, want %v", err, ErrGoogleTokenInvalid)
	}
}
