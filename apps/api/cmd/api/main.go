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

	if err := godotenv.Load(".env"); err != nil {
		log.Println("Warning: No .env file found or error parsing it:", err)
	}
	cfg := config.Load()

	pool, err := db.NewPool(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	repos := &repo.Repositories{
		UserRepo:    repo.NewUserRepo(pool),
		SchoolRepo:  repo.NewSchoolRepo(pool),
		ClassRepo:   repo.NewClassRepo(pool),
		StudentRepo: repo.NewStudentRepo(pool),
	}

	// Khởi tạo Authenticator
	jwtAuth := auth.NewAuthenticator(cfg.JWTSecret, cfg.JWTTTLMinutes)

	// Khởi tạo Services
	authService := &service.AuthService{
		UserRepo: repos.UserRepo,
		JWTAuth:  jwtAuth,
	}

	schoolService := &service.SchoolService{
		SchoolRepo: repos.SchoolRepo,
	}

	classService := &service.ClassService{
		ClassRepo: repos.ClassRepo,
	}

	studentService := &service.StudentService{
		StudentRepo: repos.StudentRepo,
	}

	// Khởi tạo Handlers
	authHandler := &v1handlers.AuthHandler{
		AuthService: authService,
	}

	schoolHandler := &v1handlers.SchoolHandler{
		SchoolService: schoolService,
	}

	classHandler := &v1handlers.ClassHandler{
		ClassService: classService,
	}

	studentHandler := &v1handlers.StudentHandler{
		StudentService: studentService,
	}

	// Khởi tạo router với handlers đã được tạo
	r := httpapi.NewRouter(cfg.JWTSecret, cfg.JWTTTLMinutes, authHandler, schoolHandler, classHandler, studentHandler)

	log.Println("listening on :" + cfg.Port)
	log.Fatal(r.Run(":" + cfg.Port))
}
