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
	parentHandler *v1handlers.ParentHandler,
	parentScopeHandler *v1handlers.ParentScopeHandler,
	parentCodeHandler *v1handlers.ParentCodeHandler,
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
		v1.POST("/users/activate-token", userHandler.ActivateUserWithToken)
		v1.POST("/register/parent", parentCodeHandler.RegisterParent)

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

				// giáo viên xem lịch sử điểm danh của một học sinh trong lớp của mình
				teacherScope.GET("/students/:student_id/attendance", teacherScopeHandler.ListAttendance)

				// giáo viên cập nhật hồ sơ cá nhân của mình
				teacherScope.PUT("/profile", teacherScopeHandler.UpdateMyProfile)

				// bài đăng (posts)
				teacherScope.POST("/posts", teacherScopeHandler.CreatePost)
				teacherScope.GET("/classes/:class_id/posts", teacherScopeHandler.ListClassPosts)
				teacherScope.GET("/students/:student_id/posts", teacherScopeHandler.ListStudentPosts)
			}

			// parent routes (phụ huynh xem thông tin con mình)
			parentScope := protected.Group("/parent")
			parentScope.Use(middleware.RequireRole("PARENT"))
			{
				// phụ huynh xem danh sách con của mình
				parentScope.GET("/children", parentScopeHandler.MyChildren)

				// phụ huynh xem feed tổng hợp của tất cả con (aggregated feed)
				parentScope.GET("/feed", parentScopeHandler.GetMyFeed)

				// phụ huynh xem bài đăng của lớp con mình
				parentScope.GET("/children/:student_id/class-posts", parentScopeHandler.ListMyChildClassPosts)

				// phụ huynh xem bài đăng riêng của con mình (student scope)
				parentScope.GET("/children/:student_id/student-posts", parentScopeHandler.ListMyChildStudentPosts)

				// phụ huynh xem tất cả bài đăng liên quan đến con mình
				parentScope.GET("/children/:student_id/posts", parentScopeHandler.ListAllMyChildPosts)
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
				students.GET("/student/:current_class_id", studentHandler.ListByClass)

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

				// parent routes (ADMIN only - quản lý phụ huynh)
				parents := admin.Group("/parents")
				{
					// lấy danh sách tất cả phụ huynh
					parents.GET("", parentHandler.List)

					// lấy thông tin phụ huynh theo ID
					parents.GET("/:parent_id", parentHandler.GetByID)

					// gán phụ huynh cho học sinh
					parents.POST("/:parent_id/students/:student_id", parentHandler.AssignStudent)

					// hủy gán phụ huynh khỏi học sinh
					parents.DELETE("/:parent_id/students/:student_id", parentHandler.UnassignStudent)
				}

				// parent code routes (ADMIN only - tạo parent codes)
				parentCodes := admin.Group("/students")
				{
					// tạo parent code cho student
					parentCodes.POST("/:student_id/generate-parent-code", parentCodeHandler.GenerateCodeForStudent)
				}
			}
		}
	}

	return r
}
