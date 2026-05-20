package main

import (
	"github.com/hoanghonghuy/iris-app/apps/api/internal/auth"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/config"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/repo"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/service"

	v1analyticshandlers "github.com/hoanghonghuy/iris-app/apps/api/internal/api/v1/handlers/analytics"
	v1auditloghandlers "github.com/hoanghonghuy/iris-app/apps/api/internal/api/v1/handlers/audit_log"
	v1authhandlers "github.com/hoanghonghuy/iris-app/apps/api/internal/api/v1/handlers/auth"
	v1classhandlers "github.com/hoanghonghuy/iris-app/apps/api/internal/api/v1/handlers/class"
	v1parenthandlers "github.com/hoanghonghuy/iris-app/apps/api/internal/api/v1/handlers/parent"
	v1parentcodehandlers "github.com/hoanghonghuy/iris-app/apps/api/internal/api/v1/handlers/parent_code"
	v1parentscopehandlers "github.com/hoanghonghuy/iris-app/apps/api/internal/api/v1/handlers/parent_scope"
	v1schoolhandlers "github.com/hoanghonghuy/iris-app/apps/api/internal/api/v1/handlers/school"
	v1schooladminhandlers "github.com/hoanghonghuy/iris-app/apps/api/internal/api/v1/handlers/school_admin"
	v1studenthandlers "github.com/hoanghonghuy/iris-app/apps/api/internal/api/v1/handlers/student"
	v1teacherhandlers "github.com/hoanghonghuy/iris-app/apps/api/internal/api/v1/handlers/teacher"
	v1teacherscopehandlers "github.com/hoanghonghuy/iris-app/apps/api/internal/api/v1/handlers/teacher_scope"
	v1userhandlers "github.com/hoanghonghuy/iris-app/apps/api/internal/api/v1/handlers/user"
	"github.com/jackc/pgx/v5/pgxpool"
)

// appServices = tập hợp tất cả các dịch vụ trong ứng dụng.
// Mỗi service chứa logic nghiệp vụ cho một domain (auth, school, student, v.v.).
type appServices struct {
	authService         *service.AuthService
	schoolService       *service.SchoolService
	classService        *service.ClassService
	studentService      *service.StudentService
	userService         *service.UserService
	teacherService      *service.TeacherService
	teacherScopeService *service.TeacherScopeService
	appointmentService  *service.AppointmentService
	parentService       *service.ParentService
	parentScopeService  *service.ParentScopeService
	auditLogService     *service.AuditLogService
	parentCodeService   *service.ParentCodeService
	schoolAdminService  *service.SchoolAdminService
	analyticsService    *service.AnalyticsService
	chatService         *service.ChatService
}

// appHandlers = tập hợp tất cả các handler HTTP trong ứng dụng.
// Mỗi handler nhận request từ route, gọi service, và trả về response.
type appHandlers struct {
	authHandler         *v1authhandlers.AuthHandler
	schoolHandler       *v1schoolhandlers.SchoolHandler
	classHandler        *v1classhandlers.ClassHandler
	studentHandler      *v1studenthandlers.StudentHandler
	userHandler         *v1userhandlers.UserHandler
	teacherHandler      *v1teacherhandlers.TeacherHandler
	teacherScopeHandler *v1teacherscopehandlers.TeacherScopeHandler
	parentHandler       *v1parenthandlers.ParentHandler
	parentScopeHandler  *v1parentscopehandlers.ParentScopeHandler
	auditLogHandler     *v1auditloghandlers.AuditLogHandler
	parentCodeHandler   *v1parentcodehandlers.ParentCodeHandler
	schoolAdminHandler  *v1schooladminhandlers.SchoolAdminHandler
	analyticsHandler    *v1analyticshandlers.AnalyticsHandler
}

// initRepositories khởi tạo tất cả các repository (data access layer).
// Repository này kết nối trực tiếp với database PostgreSQL thông qua connection pool.
// Nhận vào pgxpool.Pool là connection pool được tạo sẵn từ main function.
func initRepositories(pool *pgxpool.Pool) *repo.Repositories {
	return &repo.Repositories{
		UserRepo:                repo.NewUserRepo(pool),
		SchoolRepo:              repo.NewSchoolRepo(pool),
		ClassRepo:               repo.NewClassRepo(pool),
		StudentRepo:             repo.NewStudentRepo(pool),
		StudentParentRepo:       repo.NewStudentParentRepo(pool),
		ParentRepo:              repo.NewParentRepo(pool),
		ParentCodeRepo:          repo.NewParentCodeRepo(pool),
		TeacherRepo:             repo.NewTeacherRepo(pool),
		TeacherClassRepo:        repo.NewTeacherClassRepo(pool),
		TeacherScopeRepo:        repo.NewTeacherScopeRepo(pool),
		HealthLogRepo:           repo.NewHealthLogRepo(pool),
		ParentScopeRepo:         repo.NewParentScopeRepo(pool),
		PostInteractionRepo:     repo.NewPostInteractionRepo(pool),
		AppointmentRepo:         repo.NewAppointmentRepo(pool),
		AuditLogRepo:            repo.NewAuditLogRepo(pool),
		SchoolAdminRepo:         repo.NewSchoolAdminRepo(pool),
		ResetTokenRepo:          repo.NewResetTokenRepo(pool),
		RefreshTokenRepo:        repo.NewRefreshTokenRepo(pool),
		ChatRepo:                repo.NewChatRepo(pool),
		AnalyticsTimeseriesRepo: repo.NewAnalyticsTimeseriesRepo(pool),
	}
}

// initServices khởi tạo tất cả các service (business logic layer).
// Service này chứa logic nghiệp vụ và được inject với repositories để truy cập dữ liệu.
// Nhận vào các dependencies cần thiết: repositories, config, authenticator, xác thực Google, email sender, frontend URL.
func initServices(
	repos *repo.Repositories,
	cfg config.Config,
	jwtAuth *auth.Authenticator,
	googleVerifier auth.GoogleTokenVerifier,
	emailSender service.EmailSender,
	frontendURL string,
) *appServices {
	return &appServices{
		authService: service.NewAuthService(repos.UserRepo, repos.SchoolAdminRepo, repos.RefreshTokenRepo, jwtAuth, service.AuthServiceOptions{
			GoogleVerifier:  googleVerifier,
			GoogleEnabled:   cfg.GoogleLoginEnabled,
			GoogleHD:        cfg.GoogleHostedDomain,
			RefreshTTLHours: cfg.JWTRefreshTTLHours,
		}),
		schoolService:       service.NewSchoolService(repos.SchoolRepo),
		classService:        service.NewClassService(repos.ClassRepo),
		studentService:      service.NewStudentService(repos.StudentRepo, repos.ClassRepo),
		userService:         service.NewUserService(repos.UserRepo, repos.ResetTokenRepo, jwtAuth, emailSender, frontendURL),
		teacherService:      service.NewTeacherService(repos.TeacherRepo, repos.TeacherClassRepo, repos.ClassRepo),
		teacherScopeService: service.NewTeacherScopeService(repos.TeacherScopeRepo, repos.HealthLogRepo, repos.TeacherRepo, repos.PostInteractionRepo),
		appointmentService:  service.NewAppointmentService(repos.AppointmentRepo),
		parentService:       service.NewParentService(repos.ParentRepo, repos.StudentParentRepo, repos.StudentRepo),
		parentScopeService:  service.NewParentScopeService(repos.ParentScopeRepo, repos.PostInteractionRepo, repos.AppointmentRepo),
		auditLogService:     service.NewAuditLogService(repos.AuditLogRepo),
		parentCodeService:   service.NewParentCodeService(repos.ParentCodeRepo, repos.UserRepo, repos.ParentRepo, repos.StudentParentRepo, repos.StudentRepo, jwtAuth, googleVerifier, cfg.GoogleLoginEnabled, cfg.GoogleHostedDomain),
		schoolAdminService:  service.NewSchoolAdminService(repos.SchoolAdminRepo, repos.UserRepo),
		analyticsService:    service.NewAnalyticsService(repos),
		chatService:         service.NewChatService(repos.ChatRepo),
	}
}

// initHandlers khởi tạo tất cả các handler HTTP (presentation layer).
// Handler này nhận request từ router, gọi service, xử lý logic, rồi trả về response cho client.
// Nhận vào các service được khởi tạo sẵn từ initServices function.
func initHandlers(services *appServices) *appHandlers {
	return &appHandlers{
		authHandler:         v1authhandlers.NewAuthHandler(services.authService, services.userService),
		schoolHandler:       v1schoolhandlers.NewSchoolHandler(services.schoolService),
		classHandler:        v1classhandlers.NewClassHandler(services.classService),
		studentHandler:      v1studenthandlers.NewStudentHandler(services.studentService),
		userHandler:         v1userhandlers.NewUserHandler(services.userService),
		teacherHandler:      v1teacherhandlers.NewTeacherHandler(services.teacherService),
		teacherScopeHandler: v1teacherscopehandlers.NewTeacherScopeHandler(services.teacherScopeService, services.appointmentService),
		parentHandler:       v1parenthandlers.NewParentHandler(services.parentService),
		parentScopeHandler:  v1parentscopehandlers.NewParentScopeHandler(services.parentScopeService, services.appointmentService),
		auditLogHandler:     v1auditloghandlers.NewAuditLogHandler(services.auditLogService),
		parentCodeHandler:   v1parentcodehandlers.NewParentCodeHandler(services.parentCodeService),
		schoolAdminHandler:  v1schooladminhandlers.NewSchoolAdminHandler(services.schoolAdminService),
		analyticsHandler:    v1analyticshandlers.NewAnalyticsHandler(services.analyticsService),
	}
}
