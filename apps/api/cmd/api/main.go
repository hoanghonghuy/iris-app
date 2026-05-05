package main

import (
	"context"
	"log"
	"time"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/auth"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/config"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/db"
	httpapi "github.com/hoanghonghuy/iris-app/apps/api/internal/http"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/middleware"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/service"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/ws"
	v1chathandlers "github.com/hoanghonghuy/iris-app/apps/api/internal/api/v1/handlers/chat"
	"github.com/joho/godotenv"
)

func main() {
	// load config
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Warning: No .env file found or error parsing it:", err)
	}
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	// Database connection
	pool, err := db.NewPool(context.Background(), cfg.DatabaseURL, cfg.DBMaxConns)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	// Repositories
	repos := initRepositories(pool)

	// Authenticator
	jwtAuth := auth.NewAuthenticator(cfg.JWTSecret, cfg.JWTTTLMinutes)
	var googleVerifier auth.GoogleTokenVerifier
	if cfg.GoogleLoginEnabled {
		googleVerifier, err = auth.NewGoogleIDTokenVerifier(cfg.GoogleClientID)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Email sender (dev mode: log to console; prod: set SMTP_HOST env)
	frontendURL := cfg.FrontendURL
	if frontendURL == "" {
		frontendURL = "http://localhost:3000"
	}

	var emailSender service.EmailSender
	if cfg.SMTPHost != "" {
		emailSender = service.NewSMTPEmailSender(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPass, cfg.SMTPFrom, cfg.SMTPFromName, frontendURL)
		log.Println("Email: SMTP mode (", cfg.SMTPHost, ")")
	} else {
		emailSender = service.NewLogEmailSender(frontendURL)
		log.Println("Email: LOG mode (no SMTP_HOST set, emails printed to console)")
	}

	// Services
	services := initServices(repos, cfg, jwtAuth, googleVerifier, emailSender, frontendURL)

	// Handlers
	handlers := initHandlers(services)

	// WebSocket Hub (chạy goroutine background)
	hub := ws.NewHub()
	go hub.Run()

	chatHandler := v1chathandlers.NewChatHandler(services.chatService, hub, cfg.JWTSecret, cfg.AllowedOrigins)
	// build các middleware rate limit xác thực từ config.
	authRateLimitWindow := time.Duration(cfg.AuthRateLimitWindowSeconds) * time.Second
	authLoginLimiter := middleware.NewIPFixedWindowRateLimitWithConfig(middleware.FixedWindowRateLimitConfig{
		MaxRequests:  cfg.AuthLoginRateLimit,
		Window:       authRateLimitWindow,
		CleanupEvery: cfg.AuthRateLimitCleanupEvery,
		StaleTTL:     time.Duration(cfg.AuthRateLimitStaleTTLMultiplier) * authRateLimitWindow,
	})
	authForgotLimiter := middleware.NewIPFixedWindowRateLimitWithConfig(middleware.FixedWindowRateLimitConfig{
		MaxRequests:  cfg.AuthForgotRateLimit,
		Window:       authRateLimitWindow,
		CleanupEvery: cfg.AuthRateLimitCleanupEvery,
		StaleTTL:     time.Duration(cfg.AuthRateLimitStaleTTLMultiplier) * authRateLimitWindow,
	})
	authResetLimiter := middleware.NewIPFixedWindowRateLimitWithConfig(middleware.FixedWindowRateLimitConfig{
		MaxRequests:  cfg.AuthResetRateLimit,
		Window:       authRateLimitWindow,
		CleanupEvery: cfg.AuthRateLimitCleanupEvery,
		StaleTTL:     time.Duration(cfg.AuthRateLimitStaleTTLMultiplier) * authRateLimitWindow,
	})

	// Router
	r := httpapi.NewRouter(
		cfg.JWTSecret,
		cfg.JWTTTLMinutes,
		cfg.AllowedOrigins,
		authLoginLimiter,
		authForgotLimiter,
		authResetLimiter,
		handlers.authHandler,
		handlers.schoolHandler,
		handlers.classHandler,
		handlers.studentHandler,
		handlers.userHandler,
		handlers.teacherScopeHandler,
		handlers.teacherHandler,
		handlers.parentHandler,
		handlers.parentScopeHandler,
		handlers.parentCodeHandler,
		handlers.schoolAdminHandler,
		handlers.analyticsHandler,
		handlers.auditLogHandler,
		services.auditLogService,
		chatHandler,
	)

	// Start server
	log.Println("listening on :" + cfg.Port)
	log.Fatal(r.Run(":" + cfg.Port))
}
