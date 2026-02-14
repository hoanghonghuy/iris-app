package main

import (
	"context"
	"log"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/auth"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/config"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/db"
	httpapi "github.com/hoanghonghuy/iris-app/apps/api/internal/http"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/repo"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/service"

	v1handlers "github.com/hoanghonghuy/iris-app/apps/api/internal/api/v1/handlers"
	"github.com/joho/godotenv"
)

func main() {
	// load config
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Warning: No .env file found or error parsing it:", err)
	}
	cfg := config.Load()

	// Database connection
	pool, err := db.NewPool(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	// Repositories
	// TODO: tách ra hàm helper initRepositories(pool) *repo.Repositories
	repos := &repo.Repositories{
		UserRepo:          repo.NewUserRepo(pool),
		SchoolRepo:        repo.NewSchoolRepo(pool),
		ClassRepo:         repo.NewClassRepo(pool),
		StudentRepo:       repo.NewStudentRepo(pool),
		StudentParentRepo: repo.NewStudentParentRepo(pool),
		ParentRepo:        repo.NewParentRepo(pool),
		ParentCodeRepo:    repo.NewParentCodeRepo(pool),
		TeacherRepo:       repo.NewTeacherRepo(pool),
		TeacherClassRepo:  repo.NewTeacherClassRepo(pool),
		TeacherScopeRepo:  repo.NewTeacherScopeRepo(pool),
		ParentScopeRepo:   repo.NewParentScopeRepo(pool),
		SchoolAdminRepo:   repo.NewSchoolAdminRepo(pool),
	}

	// Authenticator
	jwtAuth := auth.NewAuthenticator(cfg.JWTSecret, cfg.JWTTTLMinutes)

	// Services
	// TODO: tách ra hàm helper initServices(repos, jwtAuth) *Services
	var (
		authService         = service.NewAuthService(repos.UserRepo, repos.SchoolAdminRepo, jwtAuth)
		schoolService       = service.NewSchoolService(repos.SchoolRepo)
		classService        = service.NewClassService(repos.ClassRepo)
		studentService      = service.NewStudentService(repos.StudentRepo)
		userService         = service.NewUserService(repos.UserRepo, jwtAuth)
		teacherService      = service.NewTeacherService(repos.TeacherRepo, repos.TeacherClassRepo)
		teacherScopeService = service.NewTeacherScopeService(repos.TeacherScopeRepo, repos.TeacherRepo)
		parentService       = service.NewParentService(repos.ParentRepo, repos.StudentParentRepo)
		parentScopeService  = service.NewParentScopeService(repos.ParentScopeRepo)
		parentCodeService   = service.NewParentCodeService(repos.ParentCodeRepo, repos.UserRepo, repos.ParentRepo, repos.StudentParentRepo, repos.StudentRepo, jwtAuth)
		schoolAdminService  = service.NewSchoolAdminService(repos.SchoolAdminRepo, repos.UserRepo)
	)

	// Handlers
	// TODO: tách ra hàm helper initHandlers(services) *Handlers
	var (
		authHandler         = v1handlers.NewAuthHandler(authService)
		schoolHandler       = v1handlers.NewSchoolHandler(schoolService)
		classHandler        = v1handlers.NewClassHandler(classService)
		studentHandler      = v1handlers.NewStudentHandler(studentService)
		userHandler         = v1handlers.NewUserHandler(userService)
		teacherHandler      = v1handlers.NewTeacherHandler(teacherService)
		teacherScopeHandler = v1handlers.NewTeacherScopeHandler(teacherScopeService)
		parentHandler       = v1handlers.NewParentHandler(parentService)
		parentScopeHandler  = v1handlers.NewParentScopeHandler(parentScopeService)
		parentCodeHandler   = v1handlers.NewParentCodeHandler(parentCodeService)
		schoolAdminHandler  = v1handlers.NewSchoolAdminHandler(schoolAdminService)
	)

	// Router
	r := httpapi.NewRouter(
		cfg.JWTSecret,
		cfg.JWTTTLMinutes,
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
	)

	// Start server
	log.Println("listening on :" + cfg.Port)
	log.Fatal(r.Run(":" + cfg.Port))
}
