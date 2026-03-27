package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	DatabaseURL               string
	DBMaxConns                int32
	JWTSecret                 string
	JWTTTLMinutes             int
	Port                      string
	AllowedOrigins            []string // CORS + WS origin allowlist
	WSAllowQueryTokenFallback bool

	// SMTP (empty = dev mode, uses LogEmailSender)
	SMTPHost    string
	SMTPPort    string
	SMTPUser    string
	SMTPPass    string
	FrontendURL string
}

// Load reads configuration from environment variables.
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
		DatabaseURL:               databaseURL,
		DBMaxConns:                maxConns,
		JWTSecret:                 jwtSecret,
		JWTTTLMinutes:             ttl,
		Port:                      port,
		AllowedOrigins:            allowedOrigins,
		WSAllowQueryTokenFallback: parseBoolEnv("WS_ALLOW_QUERY_TOKEN_FALLBACK", false),
		SMTPHost:                  os.Getenv("SMTP_HOST"),
		SMTPPort:                  os.Getenv("SMTP_PORT"),
		SMTPUser:                  os.Getenv("SMTP_USER"),
		SMTPPass:                  os.Getenv("SMTP_PASS"),
		FrontendURL:               os.Getenv("FRONTEND_URL"),
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
