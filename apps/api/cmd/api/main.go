package main

import (
	"context"
	"log"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/config"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/db"
	httpapi "github.com/hoanghonghuy/iris-app/apps/api/internal/http"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/repo"
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
		UserRepo: repo.NewUserRepo(pool),
		SchoolRepo: repo.NewSchoolRepo(pool),
	}
	r := httpapi.NewRouter(repos, cfg.JWTSecret, cfg.JWTTTLMinutes)

	log.Println("listening on :" + cfg.Port)
	log.Fatal(r.Run(":" + cfg.Port))
}
