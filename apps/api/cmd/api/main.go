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
	"github.com/hoanghonghuy/iris-app/apps/api/internal/repo"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/service"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/ws"

	v1handlers "github.com/hoanghonghuy/iris-app/apps/api/internal/api/v1/handlers"
	v1chathandlers "github.com/hoanghonghuy/iris-app/apps/api/internal/api/v1/handlers/chat"
	v1parentcodehandlers "github.com/hoanghonghuy/iris-app/apps/api/internal/api/v1/handlers/parent_code"
	v1parentscopehandlers "github.com/hoanghonghuy/iris-app/apps/api/internal/api/v1/handlers/parent_scope"
	v1studenthandlers "github.com/hoanghonghuy/iris-app/apps/api/internal/api/v1/handlers/student"
	v1teacherhandlers "github.com/hoanghonghuy/iris-app/apps/api/internal/api/v1/handlers/teacher"
	v1teacherscopehandlers "github.com/hoanghonghuy/iris-app/apps/api/internal/api/v1/handlers/teacher_scope"
	v1userhandlers "github.com/hoanghonghuy/iris-app/apps/api/internal/api/v1/handlers/user"
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
	// TODO: tách ra hàm helper initRepositories(pool) *repo.Repositories
	repos := &repo.Repositories{
		UserRepo:            repo.NewUserRepo(pool),
		SchoolRepo:          repo.NewSchoolRepo(pool),
		ClassRepo:           repo.NewClassRepo(pool),
		StudentRepo:         repo.NewStudentRepo(pool),
		StudentParentRepo:   repo.NewStudentParentRepo(pool),
		ParentRepo:          repo.NewParentRepo(pool),
		ParentCodeRepo:      repo.NewParentCodeRepo(pool),
		TeacherRepo:         repo.NewTeacherRepo(pool),
		TeacherClassRepo:    repo.NewTeacherClassRepo(pool),
		TeacherScopeRepo:    repo.NewTeacherScopeRepo(pool),
		HealthLogRepo:       repo.NewHealthLogRepo(pool),
		ParentScopeRepo:     repo.NewParentScopeRepo(pool),
		PostInteractionRepo: repo.NewPostInteractionRepo(pool),
		AppointmentRepo:     repo.NewAppointmentRepo(pool),
		AuditLogRepo:        repo.NewAuditLogRepo(pool),
		SchoolAdminRepo:     repo.NewSchoolAdminRepo(pool),
		ResetTokenRepo:      repo.NewResetTokenRepo(pool),
		ChatRepo:            repo.NewChatRepo(pool),
	}

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
		emailSender = service.NewSMTPEmailSender(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPass, frontendURL)
		log.Println("Email: SMTP mode (", cfg.SMTPHost, ")")
	} else {
		emailSender = service.NewLogEmailSender(frontendURL)
		log.Println("Email: LOG mode (no SMTP_HOST set, emails printed to console)")
	}

	// Services
	// TODO: tách ra hàm helper initServices(repos, jwtAuth) *Services
	var (
		authService         = service.NewAuthService(repos.UserRepo, repos.SchoolAdminRepo, jwtAuth, googleVerifier, cfg.GoogleLoginEnabled, cfg.GoogleHostedDomain)
		schoolService       = service.NewSchoolService(repos.SchoolRepo)
		classService        = service.NewClassService(repos.ClassRepo)
		studentService      = service.NewStudentService(repos.StudentRepo, repos.ClassRepo)
		userService         = service.NewUserService(repos.UserRepo, repos.ResetTokenRepo, jwtAuth, emailSender, frontendURL)
		teacherService      = service.NewTeacherService(repos.TeacherRepo, repos.TeacherClassRepo, repos.ClassRepo)
		teacherScopeService = service.NewTeacherScopeService(repos.TeacherScopeRepo, repos.HealthLogRepo, repos.TeacherRepo, repos.PostInteractionRepo)
		appointmentService  = service.NewAppointmentService(repos.AppointmentRepo)
		parentService       = service.NewParentService(repos.ParentRepo, repos.StudentParentRepo, repos.StudentRepo)
		parentScopeService  = service.NewParentScopeService(repos.ParentScopeRepo, repos.PostInteractionRepo, repos.AppointmentRepo)
		auditLogService     = service.NewAuditLogService(repos.AuditLogRepo)
		parentCodeService   = service.NewParentCodeService(repos.ParentCodeRepo, repos.UserRepo, repos.ParentRepo, repos.StudentParentRepo, repos.StudentRepo, jwtAuth, googleVerifier, cfg.GoogleLoginEnabled, cfg.GoogleHostedDomain)
		schoolAdminService  = service.NewSchoolAdminService(repos.SchoolAdminRepo, repos.UserRepo)
		analyticsService    = service.NewAnalyticsService(repos)
		chatService         = service.NewChatService(repos.ChatRepo)
	)

	// Handlers
	// TODO: tách ra hàm helper initHandlers(services) *Handlers
	var (
		authHandler         = v1handlers.NewAuthHandler(authService, userService)
		schoolHandler       = v1handlers.NewSchoolHandler(schoolService)
		classHandler        = v1handlers.NewClassHandler(classService)
		studentHandler      = v1studenthandlers.NewStudentHandler(studentService)
		userHandler         = v1userhandlers.NewUserHandler(userService)
		teacherHandler      = v1teacherhandlers.NewTeacherHandler(teacherService)
		teacherScopeHandler = v1teacherscopehandlers.NewTeacherScopeHandler(teacherScopeService, appointmentService)
		parentHandler       = v1handlers.NewParentHandler(parentService)
		parentScopeHandler  = v1parentscopehandlers.NewParentScopeHandler(parentScopeService, appointmentService)
		auditLogHandler     = v1handlers.NewAuditLogHandler(auditLogService)
		parentCodeHandler   = v1parentcodehandlers.NewParentCodeHandler(parentCodeService)
		schoolAdminHandler  = v1handlers.NewSchoolAdminHandler(schoolAdminService)
		analyticsHandler    = v1handlers.NewAnalyticsHandler(analyticsService)
	)

	// WebSocket Hub (chạy goroutine background)
	hub := ws.NewHub()
	go hub.Run()

	chatHandler := v1chathandlers.NewChatHandler(chatService, hub, cfg.JWTSecret, cfg.AllowedOrigins)
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
		authHandler,
		schoolHandler,
		classHandler,
		studentHandler,
		userHandler,
		teacherScopeHandler,
		teacherHandler,
		parentHandler,
		parentScopeHandler,
		parentCodeHandler,
		schoolAdminHandler,
		analyticsHandler,
		auditLogHandler,
		auditLogService,
		chatHandler,
	)

	// Start server
	log.Println("listening on :" + cfg.Port)
	log.Fatal(r.Run(":" + cfg.Port))
}
