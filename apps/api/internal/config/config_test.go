package config

import "testing"

func TestLoadParsesAuthResetRateLimitFromEnv(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/iris?sslmode=disable")
	t.Setenv("JWT_SECRET", "test-secret")
	t.Setenv("AUTH_RESET_PASSWORD_RATE_LIMIT", "7")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.AuthResetRateLimit != 7 {
		t.Fatalf("AuthResetRateLimit = %d, want %d", cfg.AuthResetRateLimit, 7)
	}
}

func TestLoadUsesDefaultAuthResetRateLimitWhenEnvInvalid(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/iris?sslmode=disable")
	t.Setenv("JWT_SECRET", "test-secret")
	t.Setenv("AUTH_RESET_PASSWORD_RATE_LIMIT", "-1")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.AuthResetRateLimit != 5 {
		t.Fatalf("AuthResetRateLimit = %d, want default %d", cfg.AuthResetRateLimit, 5)
	}
}
