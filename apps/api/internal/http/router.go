package httpapi

import (
	"github.com/gin-gonic/gin"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/middleware"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/response"

	v1handlers "github.com/hoanghonghuy/iris-app/apps/api/internal/api/v1/handlers"
)

// NewRouter tạo và cấu hình router HTTP sử dụng Gin
func NewRouter(
	jwtSecret string,
	ttlMinutes int,
	authHandler *v1handlers.AuthHandler,
	schoolHandler *v1handlers.SchoolHandler,
	classHandler *v1handlers.ClassHandler,
	studentHandler *v1handlers.StudentHandler,
	userHandler *v1handlers.UserHandler,
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
		v1.POST("/users/activate", userHandler.ActivateUser)

		// Protected routes (require valid JWT)
		protected := v1.Group("/")
		protected.Use(middleware.AuthJWT(jwtSecret))
		// /me endpoint trả về thông tin user hiện tại từ JWT claims (không cần query DB)
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
		classes.POST("/school", classHandler.Create)
		classes.GET("/school/:school_id", classHandler.ListBySchool)

		// Student routes
		students := admin.Group("/students")
		students.POST("/student", studentHandler.Create)
		students.GET("/student/:class_id", studentHandler.ListByClass)

		// User routes (ADMIN only - quản lý users)
		// Pattern:
		//   - GET /api/v1/admin/users - list all users
		//   - GET /api/v1/admin/users/:userid - get user by ID
		//   - POST /api/v1/admin/users/:userid/lock - lock user account
		//   - POST /api/v1/admin/users/:userid/unlock - unlock user account
		//   - POST /api/v1/admin/users/:userid/roles - assign role to user
		users := admin.Group("/users")
		users.POST("", userHandler.CreateUser)
		users.GET("", userHandler.List)
		users.GET("/:userid", userHandler.GetByID)
		users.POST("/:userid/lock", userHandler.Lock)
		users.POST("/:userid/unlock", userHandler.Unlock)
		users.POST("/:userid/roles", userHandler.AssignRole)
	}

	return r
}
