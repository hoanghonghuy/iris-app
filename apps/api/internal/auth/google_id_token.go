package auth

import (
	"context"
	"errors"

	"google.golang.org/api/idtoken"
)

var ErrGoogleTokenInvalid = errors.New("google token invalid")

type GoogleIdentity struct {
	Sub           string
	Email         string
	EmailVerified bool
	HostedDomain  string
	Name          string
}

type GoogleTokenVerifier interface {
	Verify(ctx context.Context, rawIDToken string) (*GoogleIdentity, error)
}

type googleIDTokenVerifier struct {
	validator *idtoken.Validator
	clientID  string
}

func NewGoogleIDTokenVerifier(clientID string) (GoogleTokenVerifier, error) {
	v, err := idtoken.NewValidator(context.Background())
	if err != nil {
		return nil, err
	}
	return &googleIDTokenVerifier{validator: v, clientID: clientID}, nil
}

func (g *googleIDTokenVerifier) Verify(ctx context.Context, rawIDToken string) (*GoogleIdentity, error) {
	if g.clientID == "" || rawIDToken == "" {
		return nil, ErrGoogleTokenInvalid
	}
	payload, err := g.validator.Validate(ctx, rawIDToken, g.clientID)
	if err != nil {
		return nil, ErrGoogleTokenInvalid
	}

	identity := &GoogleIdentity{
		Sub:           payload.Subject,
		Email:         stringClaim(payload.Claims, "email"),
		EmailVerified: boolClaim(payload.Claims, "email_verified"),
		HostedDomain:  stringClaim(payload.Claims, "hd"),
		Name:          stringClaim(payload.Claims, "name"),
	}

	if identity.Sub == "" || identity.Email == "" || !identity.EmailVerified {
		return nil, ErrGoogleTokenInvalid
	}

	return identity, nil
}

func stringClaim(claims map[string]interface{}, key string) string {
	v, ok := claims[key]
	if !ok || v == nil {
		return ""
	}
	s, ok := v.(string)
	if !ok {
		return ""
	}
	return s
}

func boolClaim(claims map[string]interface{}, key string) bool {
	v, ok := claims[key]
	if !ok || v == nil {
		return false
	}
	b, ok := v.(bool)
	if !ok {
		return false
	}
	return b
}
