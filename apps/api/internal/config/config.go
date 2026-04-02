package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	DatabaseURL   string
	DBMaxConns    int32
	JWTSecret     string
	JWTTTLMinutes int
	// Cấu hình rate limit dạng fixed-window cho login/forgot-password xác thực.
	AuthLoginRateLimit              int
	AuthForgotRateLimit             int
	AuthRateLimitWindowSeconds      int
	AuthRateLimitCleanupEvery       int
	AuthRateLimitStaleTTLMultiplier int
	GoogleLoginEnabled              bool
	GoogleClientID                  string
	GoogleHostedDomain              string
	Port                            string
	AllowedOrigins                  []string // CORS + WS origin allowlist

	// SMTP (empty = dev mode, dùng LogEmailSender)
	SMTPHost    string
	SMTPPort    string
	SMTPUser    string
	SMTPPass    string
	FrontendURL string
}

// Load đọc cấu hình từ các biến môi trường.
func Load() (Config, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	ttl := 60
	if v := os.Getenv("JWT_TTL_MINUTES"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			ttl = n
		}
	}

	maxConns := int32(50)
	if v := os.Getenv("DB_MAX_CONNS"); v != "" {
		if n, err := strconv.ParseInt(v, 10, 32); err == nil && n > 0 {
			maxConns = int32(n)
		}
	}

	// Parse auth rate limit config (fall back to safe defaults on invalid env values).
	authLoginRateLimit := parsePositiveIntEnv("AUTH_LOGIN_RATE_LIMIT", 10)
	authForgotRateLimit := parsePositiveIntEnv("AUTH_FORGOT_PASSWORD_RATE_LIMIT", 5)
	authRateLimitWindowSeconds := parsePositiveIntEnv("AUTH_RATE_LIMIT_WINDOW_SECONDS", 60)
	authRateLimitCleanupEvery := parsePositiveIntEnv("AUTH_RATE_LIMIT_CLEANUP_EVERY", 256)
	authRateLimitStaleTTLMultiplier := parsePositiveIntEnv("AUTH_RATE_LIMIT_STALE_TTL_MULTIPLIER", 5)

	// Parse ALLOWED_ORIGINS (comma-separated, e.g. "https://app.example.com,http://localhost:3000")
	allowedOrigins := []string{"http://localhost:3000"} // default: dev
	if raw := os.Getenv("ALLOWED_ORIGINS"); raw != "" {
		allowedOrigins = strings.Split(raw, ",")
		for i, o := range allowedOrigins {
			allowedOrigins[i] = strings.TrimSpace(o)
		}
	}

	databaseURL, err := must("DATABASE_URL")
	if err != nil {
		return Config{}, err
	}
	jwtSecret, err := must("JWT_SECRET")
	if err != nil {
		return Config{}, err
	}

	return Config{
		DatabaseURL:                     databaseURL,
		DBMaxConns:                      maxConns,
		JWTSecret:                       jwtSecret,
		JWTTTLMinutes:                   ttl,
		AuthLoginRateLimit:              authLoginRateLimit,
		AuthForgotRateLimit:             authForgotRateLimit,
		AuthRateLimitWindowSeconds:      authRateLimitWindowSeconds,
		AuthRateLimitCleanupEvery:       authRateLimitCleanupEvery,
		AuthRateLimitStaleTTLMultiplier: authRateLimitStaleTTLMultiplier,
		GoogleLoginEnabled:              parseBoolEnv("GOOGLE_LOGIN_ENABLED", false),
		GoogleClientID:                  os.Getenv("GOOGLE_CLIENT_ID"),
		GoogleHostedDomain:              os.Getenv("GOOGLE_HOSTED_DOMAIN"),
		Port:                            port,
		AllowedOrigins:                  allowedOrigins,
		SMTPHost:                        os.Getenv("SMTP_HOST"),
		SMTPPort:                        os.Getenv("SMTP_PORT"),
		SMTPUser:                        os.Getenv("SMTP_USER"),
		SMTPPass:                        os.Getenv("SMTP_PASS"),
		FrontendURL:                     os.Getenv("FRONTEND_URL"),
	}, nil
}

func must(k string) (string, error) {
	v := os.Getenv(k)
	if v == "" {
		return "", fmt.Errorf("[config] missing required env variable: %s", k)
	}
	return v, nil
}

func parseBoolEnv(key string, defaultVal bool) bool {
	raw := strings.TrimSpace(strings.ToLower(os.Getenv(key)))
	switch raw {
	case "1", "true", "yes", "y", "on":
		return true
	case "0", "false", "no", "n", "off":
		return false
	case "":
		return defaultVal
	default:
		return defaultVal
	}
}

func parsePositiveIntEnv(key string, defaultVal int) int {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return defaultVal
	}
	n, err := strconv.Atoi(raw)
	if err != nil || n <= 0 {
		return defaultVal
	}
	return n
}
