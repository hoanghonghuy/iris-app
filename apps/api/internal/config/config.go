package config

import (
	"os"
	"strconv"
)

type Config struct {
	DatabaseURL   string
	DBMaxConns    int32
	JWTSecret     string
	JWTTTLMinutes int
	Port          string

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

	return Config{
		DatabaseURL:   must("DATABASE_URL"),
		DBMaxConns:    maxConns,
		JWTSecret:     must("JWT_SECRET"),
		JWTTTLMinutes: ttl,
		Port:          port,
		SMTPHost:      os.Getenv("SMTP_HOST"),
		SMTPPort:      os.Getenv("SMTP_PORT"),
		SMTPUser:      os.Getenv("SMTP_USER"),
		SMTPPass:      os.Getenv("SMTP_PASS"),
		FrontendURL:   os.Getenv("FRONTEND_URL"),
	}
}

func must(k string) string {
	v := os.Getenv(k)
	if v == "" {
		panic("missing env: " + k)
	}
	return v
}
