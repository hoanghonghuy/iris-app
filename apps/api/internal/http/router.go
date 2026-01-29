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
	teacherScopeHandler *v1handlers.TeacherScopeHandler,
	teacherHandler *v1handlers.TeacherHandler,
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
		{
			// /me endpoint trả về thông tin user hiện tại từ JWT claims (không cần query DB)
			protected.GET("/me", authHandler.Me)

			// user cập nhật mật khẩu của chính mình
			protected.PUT("/me/password", userHandler.UpdateMyPassword)

			// user xóa tài khoản của chính mình
			protected.DELETE("/me", userHandler.Delete)

			// teacher routes (chỉ có thể truy cập các lớp và học sinh được phân công)
			teacherScope := protected.Group("/teacher")
			teacherScope.Use(middleware.RequireRole("TEACHER"))
			{
				// giáo viên xem danh sách lớp của mình
				teacherScope.GET("/classes", teacherScopeHandler.MyClasses)

				// giáo viên xem danh sách học sinh trong một lớp cụ thể
				teacherScope.GET("/classes/:class_id/students", teacherScopeHandler.MyStudentsInClass)

				// giáo viên điểm danh cho học sinh trong lớp của mình
				teacherScope.POST("/attendance", teacherScopeHandler.MarkAttendance)

				// Giáo viên tạo nhật ký sức khỏe cho học sinh trong lớp của mình
				teacherScope.POST("/health", teacherScopeHandler.CreateHealth)

				// giáo viên xem nhật ký sức khỏe của một học sinh cụ thể trong lớp của mình
				teacherScope.GET("/students/:student_id/health", teacherScopeHandler.ListHealth)

				// giáo viên cập nhật hồ sơ cá nhân của mình
				teacherScope.PUT("/profile", teacherScopeHandler.UpdateMyProfile)
			}

			// admin routes (require ADMIN role)
			admin := protected.Group("/admin")
			admin.Use(middleware.RequireRole("ADMIN"))
			{
				// Admin ping/health check
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

				// teacher routes (ADMIN only - quản lý giáo viên)
				teachers := admin.Group("/teachers")
				{
					// lấy danh sách tất cả giáo viên
					teachers.GET("", teacherHandler.List)

					// lấy thông tin giáo viên theo ID
					teachers.GET("/:teacher_id", teacherHandler.GetByTeacherID)

					// cập nhật thông tin giáo viên
					teachers.PUT("/:teacher_id", teacherHandler.Update)

					// lấy danh sách giáo viên của một lớp
					teachers.GET("/class/:class_id", teacherHandler.ListTeacherOfClass)

					// gán giáo viên vào lớp
					teachers.POST("/:teacher_id/classes/:class_id", teacherHandler.Assign)

					// hủy gán giáo viên khỏi lớp
					teachers.DELETE("/:teacher_id/classes/:class_id", teacherHandler.Unassign)
				}
			}
		}
	}

	return r
}
