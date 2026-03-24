package config

import (
	"log"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	DatabaseURL    string
	DBMaxConns     int32
	JWTSecret      string
	JWTTTLMinutes  int
	Port           string
	AllowedOrigins []string // CORS + WS origin allowlist

	// SMTP (empty = dev mode, uses LogEmailSender)
	SMTPHost    string
	SMTPPort    string
	SMTPUser    string
	SMTPPass    string
	FrontendURL string
}

// Load reads configuration from environment variables.
func Load() Config {
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

	return Config{
		DatabaseURL:    must("DATABASE_URL"),
		DBMaxConns:     maxConns,
		JWTSecret:      must("JWT_SECRET"),
		JWTTTLMinutes:  ttl,
		Port:           port,
		AllowedOrigins: allowedOrigins,
		SMTPHost:       os.Getenv("SMTP_HOST"),
		SMTPPort:       os.Getenv("SMTP_PORT"),
		SMTPUser:       os.Getenv("SMTP_USER"),
		SMTPPass:       os.Getenv("SMTP_PASS"),
		FrontendURL:    os.Getenv("FRONTEND_URL"),
	}
}

func must(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("[config] missing required env variable: %s", k)
	}
	return v
}
