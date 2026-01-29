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
		TeacherRepo:       repo.NewTeacherRepo(pool),
		TeacherClassRepo:  repo.NewTeacherClassRepo(pool),
		TeacherScopeRepo:  repo.NewTeacherScopeRepo(pool),
		ParentScopeRepo:   repo.NewParentScopeRepo(pool),
	}

	// Authenticator
	jwtAuth := auth.NewAuthenticator(cfg.JWTSecret, cfg.JWTTTLMinutes)

	// Services
	// TODO: tách ra hàm helper initServices(repos, jwtAuth) *Services
	var (
		authService         = service.NewAuthService(repos.UserRepo, jwtAuth)
		schoolService       = service.NewSchoolService(repos.SchoolRepo)
		classService        = service.NewClassService(repos.ClassRepo)
		studentService      = service.NewStudentService(repos.StudentRepo)
		userService         = service.NewUserService(repos.UserRepo, jwtAuth)
		teacherService      = service.NewTeacherService(repos.TeacherRepo, repos.TeacherClassRepo)
		teacherScopeService = service.NewTeacherScopeService(repos.TeacherScopeRepo, repos.TeacherRepo)
		parentScopeService  = service.NewParentScopeService(repos.ParentScopeRepo)
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
		parentScopeHandler  = v1handlers.NewParentScopeHandler(parentScopeService)
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
		parentScopeHandler,
	)

	// Start server
	log.Println("listening on :" + cfg.Port)
	log.Fatal(r.Run(":" + cfg.Port))
}
