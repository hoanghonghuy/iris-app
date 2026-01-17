package httpapi

import (
	"github.com/gin-gonic/gin"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/middleware"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/response"

	v1handlers "github.com/hoanghonghuy/iris-app/apps/api/internal/api/v1/handlers"
)

// NewRouter tạo và cấu hình router HTTP sử dụng Gin
func NewRouter(
	repos interface{}, // giữ cho tương thích với main.go, không dùng
	jwtSecret string,
	ttlMinutes int,
	authHandler *v1handlers.AuthHandler,
	schoolHandler *v1handlers.SchoolHandler,
	classHandler *v1handlers.ClassHandler,
) *gin.Engine {
	r := gin.Default()

	// Setup routes
	v1 := r.Group("/api/v1")
	{
		// Public routes
		v1.GET("/health", func(c *gin.Context) {
			response.OK(c, gin.H{"ok": true})
		})
		v1.POST("/auth/login", authHandler.Login)

		// Protected routes (require valid JWT)
		protected := v1.Group("/")
		protected.Use(middleware.AuthJWT(jwtSecret))
		protected.GET("/me", authHandler.Me)

		// Admin routes (require ADMIN role)
		admin := v1.Group("/admin")
		admin.Use(middleware.AuthJWT(jwtSecret))
		admin.Use(middleware.RequireRole("ADMIN"))

		// Admin pong/health
		admin.GET("/ping", func(c *gin.Context) {
			response.OK(c, gin.H{"pong": "admin"})
		})

		// School routes
		schools := admin.Group("/schools")
		schools.POST("", schoolHandler.Create)
		schools.GET("", schoolHandler.List)

		// Class routes
		classes := admin.Group("/classes")
		classes.POST("/school/:school_id", classHandler.Create)
		classes.GET("/school/:school_id", classHandler.ListBySchool)
	}

	return r
}
