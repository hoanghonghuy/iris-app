package httpapi

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/middleware"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/response"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/service"

	v1handlers "github.com/hoanghonghuy/iris-app/apps/api/internal/api/v1/handlers"
	v1chathandlers "github.com/hoanghonghuy/iris-app/apps/api/internal/api/v1/handlers/chat"
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
)

// NewRouter tạo và cấu hình router HTTP sử dụng Gin
func NewRouter(
	jwtSecret string,
	ttlMinutes int,
	allowedOrigins []string,
	authLoginRateLimit gin.HandlerFunc,
	authForgotPasswordRateLimit gin.HandlerFunc,
	authResetPasswordRateLimit gin.HandlerFunc,
	authHandler *v1handlers.AuthHandler,
	schoolHandler *v1schoolhandlers.SchoolHandler,
	classHandler *v1classhandlers.ClassHandler,
	studentHandler *v1studenthandlers.StudentHandler,
	userHandler *v1userhandlers.UserHandler,
	teacherScopeHandler *v1teacherscopehandlers.TeacherScopeHandler,
	teacherHandler *v1teacherhandlers.TeacherHandler,
	parentHandler *v1parenthandlers.ParentHandler,
	parentScopeHandler *v1parentscopehandlers.ParentScopeHandler,
	parentCodeHandler *v1parentcodehandlers.ParentCodeHandler,
	schoolAdminHandler *v1schooladminhandlers.SchoolAdminHandler,
	analyticsHandler *v1handlers.AnalyticsHandler,
	auditLogHandler *v1handlers.AuditLogHandler,
	auditLogService *service.AuditLogService,
	chatHandler *v1chathandlers.ChatHandler,
) *gin.Engine {
	r := gin.Default()
	// Trust no proxy headers by default (safer for direct internet deployment).
	_ = r.SetTrustedProxies(nil)

	// Build origin set for O(1) lookup
	originSet := make(map[string]struct{}, len(allowedOrigins))
	for _, o := range allowedOrigins {
		originSet[o] = struct{}{}
	}

	// Security headers middleware cho API responses.
	r.Use(func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Header("Content-Security-Policy", "default-src 'none'; frame-ancestors 'none'; base-uri 'none'")
		c.Next()
	})

	// CORS middleware — chỉ cho phép origin trong allowlist
	r.Use(func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin != "" {
			if _, ok := originSet[origin]; ok {
				c.Header("Access-Control-Allow-Origin", origin)
				c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
				c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
				c.Header("Access-Control-Allow-Credentials", "true")
				c.Header("Access-Control-Max-Age", "86400")
			} else if c.Request.Method == "OPTIONS" {
				// Preflight từ origin không hợp lệ → reject ngay
				c.AbortWithStatus(http.StatusForbidden)
				return
			}
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Setup routes
	r.GET("/", func(c *gin.Context) {
		response.OK(c, gin.H{"ok": true, "service": "iris-api"})
	})

	v1 := r.Group("/api/v1")
	{
		// Public routes
		v1.GET("/health", func(c *gin.Context) {
			response.OK(c, gin.H{"ok": true})
		})
		v1.POST("/auth/login", authLoginRateLimit, authHandler.Login)
		v1.POST("/auth/login/google", authLoginRateLimit, authHandler.LoginWithGoogle)
		v1.POST("/auth/forgot-password", authForgotPasswordRateLimit, authHandler.ForgotPassword)
		v1.POST("/auth/reset-password", authResetPasswordRateLimit, authHandler.ResetPassword)
		v1.POST("/users/activate-token", userHandler.ActivateUserWithToken)
		v1.GET("/register/parent/verify", parentCodeHandler.VerifyCode)
		v1.POST("/register/parent", parentCodeHandler.RegisterParent)
		v1.POST("/register/parent/google", parentCodeHandler.RegisterParentWithGoogle)

		// WebSocket endpoint (auth qua Sec-WebSocket-Protocol header)
		v1.GET("/chat/ws", chatHandler.HandleWS)

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
				teacherScope.DELETE("/attendance", teacherScopeHandler.CancelAttendance)

				// Giáo viên tạo nhật ký sức khỏe cho học sinh trong lớp của mình
				teacherScope.POST("/health", teacherScopeHandler.CreateHealth)

				// giáo viên xem nhật ký sức khỏe của một học sinh cụ thể trong lớp của mình
				teacherScope.GET("/students/:student_id/health", teacherScopeHandler.ListHealth)

				// giáo viên xem lịch sử điểm danh của một học sinh trong lớp của mình
				teacherScope.GET("/students/:student_id/attendance", teacherScopeHandler.ListAttendance)
				teacherScope.GET("/students/:student_id/attendance-changes", teacherScopeHandler.ListAttendanceChangeLogs)
				teacherScope.GET("/classes/:class_id/attendance-changes", teacherScopeHandler.ListClassAttendanceChangeLogs)

				// giáo viên cập nhật hồ sơ cá nhân của mình
				teacherScope.PUT("/profile", teacherScopeHandler.UpdateMyProfile)

				// bài đăng (posts)
				teacherScope.POST("/posts", teacherScopeHandler.CreatePost)
				teacherScope.PUT("/posts/:post_id", teacherScopeHandler.UpdatePost)
				teacherScope.DELETE("/posts/:post_id", teacherScopeHandler.DeletePost)
				teacherScope.POST("/posts/:post_id/like", teacherScopeHandler.TogglePostLike)
				teacherScope.GET("/posts/:post_id/comments", teacherScopeHandler.ListPostComments)
				teacherScope.POST("/posts/:post_id/comments", teacherScopeHandler.CreatePostComment)
				teacherScope.POST("/posts/:post_id/share", teacherScopeHandler.SharePost)
				teacherScope.GET("/classes/:class_id/posts", teacherScopeHandler.ListClassPosts)
				teacherScope.GET("/students/:student_id/posts", teacherScopeHandler.ListStudentPosts)

				// teacher analytics
				teacherScope.GET("/analytics", analyticsHandler.TeacherDashboardStats)
				teacherScope.POST("/appointments/slots", teacherScopeHandler.CreateAppointmentSlot)
				teacherScope.GET("/appointments", teacherScopeHandler.ListMyAppointments)
				teacherScope.PATCH("/appointments/:appointment_id/status", teacherScopeHandler.UpdateAppointmentStatus)
			}

			// parent routes (phụ huynh xem thông tin con mình)
			parentScope := protected.Group("/parent")
			parentScope.Use(middleware.RequireRole("PARENT"))
			{
				// phụ huynh xem danh sách con của mình
				parentScope.GET("/children", parentScopeHandler.MyChildren)
				parentScope.GET("/analytics", parentScopeHandler.GetMyAnalytics)
				parentScope.GET("/appointments/slots", parentScopeHandler.ListAvailableAppointmentSlots)
				parentScope.POST("/appointments", parentScopeHandler.CreateAppointment)
				parentScope.GET("/appointments", parentScopeHandler.ListMyAppointments)
				parentScope.PATCH("/appointments/:appointment_id/cancel", parentScopeHandler.CancelAppointment)

				// phụ huynh xem feed tổng hợp của tất cả con (aggregated feed)
				parentScope.GET("/feed", parentScopeHandler.GetMyFeed)
				parentScope.POST("/posts/:post_id/like", parentScopeHandler.TogglePostLike)
				parentScope.GET("/posts/:post_id/comments", parentScopeHandler.ListPostComments)
				parentScope.POST("/posts/:post_id/comments", parentScopeHandler.CreatePostComment)
				parentScope.POST("/posts/:post_id/share", parentScopeHandler.SharePost)

				// phụ huynh xem bài đăng của lớp con mình
				parentScope.GET("/children/:student_id/class-posts", parentScopeHandler.ListMyChildClassPosts)

				// phụ huynh xem bài đăng riêng của con mình (student scope)
				parentScope.GET("/children/:student_id/student-posts", parentScopeHandler.ListMyChildStudentPosts)

				// phụ huynh xem tất cả bài đăng liên quan đến con mình
				parentScope.GET("/children/:student_id/posts", parentScopeHandler.ListAllMyChildPosts)
			}

			// chat routes (tất cả authenticated users đều có thể dùng)
			chat := protected.Group("/chat")
			{
				// tìm kiếm user để chat
				chat.GET("/users/search", chatHandler.SearchUsers)

				// tạo cuộc hội thoại direct (1-1)
				chat.POST("/conversations/direct", chatHandler.CreateDirectConversation)

				// lấy danh sách cuộc hội thoại của user hiện tại
				chat.GET("/conversations", chatHandler.ListConversations)

				// lấy danh sách tin nhắn của cuộc hội thoại
				chat.GET("/conversations/:conversation_id/messages", chatHandler.ListMessages)
			}

			// admin routes (SUPER_ADMIN + SCHOOL_ADMIN đều truy cập được)
			// InjectAdminScope đọc school_id từ JWT → lưu vào context
			// SUPER_ADMIN: school_id rỗng (truy cập tất cả trường)
			// SCHOOL_ADMIN: school_id có giá trị (chỉ truy cập trường mình)
			admin := protected.Group("/admin")
			admin.Use(middleware.RequireAnyRole("SUPER_ADMIN", "SCHOOL_ADMIN"))
			admin.Use(middleware.InjectAdminScope())
			admin.Use(middleware.AdminAuditLogger(auditLogService))
			{
				// Admin ping/health check
				admin.GET("/ping", func(c *gin.Context) {
					response.OK(c, gin.H{"pong": "admin"})
				})

				// Admin analytics
				admin.GET("/analytics", analyticsHandler.AdminDashboardStats)
				admin.GET("/audit-logs", auditLogHandler.List)

				// School routes (GET: cả 2 roles, POST: chỉ SUPER_ADMIN — đăng ký ở superOnly bên dưới)
				admin.GET("/schools", schoolHandler.List)

				// Class routes
				classes := admin.Group("/classes")
				classes.POST("", classHandler.Create)
				classes.GET("/by-school/:school_id", classHandler.ListBySchool)
				classes.PUT("/:class_id", classHandler.Update)
				classes.DELETE("/:class_id", classHandler.Delete)

				// Student routes
				students := admin.Group("/students")
				students.POST("", studentHandler.Create)
				students.GET("/by-class/:class_id", studentHandler.ListByClass)
				students.GET("/:student_id", studentHandler.GetProfile)
				students.PUT("/:student_id", studentHandler.Update)
				students.DELETE("/:student_id", studentHandler.Delete)

				// User routes (quản lý users)
				// AssignRole chỉ SUPER_ADMIN → đăng ký ở superOnly bên dưới
				users := admin.Group("/users")
				users.POST("", userHandler.CreateUser)
				users.GET("", userHandler.List)
				users.GET("/:user_id", userHandler.GetByID)
				users.POST("/:user_id/lock", userHandler.Lock)
				users.POST("/:user_id/unlock", userHandler.Unlock)

				// teacher routes (quản lý giáo viên)
				teachers := admin.Group("/teachers")
				{
					// lấy danh sách tất cả giáo viên
					teachers.GET("", teacherHandler.List)

					// lấy thông tin giáo viên theo ID
					teachers.GET("/:teacher_id", teacherHandler.GetByTeacherID)

					// cập nhật thông tin giáo viên
					teachers.PUT("/:teacher_id", teacherHandler.Update)

					// lấy danh sách giáo viên của một lớp
					teachers.GET("/class/:class_id", teacherHandler.ListTeachersOfClass)

					// gán giáo viên vào lớp
					teachers.POST("/:teacher_id/classes/:class_id", teacherHandler.Assign)

					// hủy gán giáo viên khỏi lớp
					teachers.DELETE("/:teacher_id/classes/:class_id", teacherHandler.Unassign)

					// xóa giáo viên
					teachers.DELETE("/:teacher_id", teacherHandler.Delete)
				}

				// parent routes (quản lý phụ huynh)
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

				// parent code routes (tạo parent codes)
				parentCodes := admin.Group("/students")
				{
					// tạo parent code cho student
					parentCodes.POST("/:student_id/generate-parent-code", parentCodeHandler.GenerateCodeForStudent)
					// thu hoi parent code hien tai
					parentCodes.DELETE("/:student_id/parent-code", parentCodeHandler.RevokeParentCode)
				}

				// Các routes chỉ SUPER_ADMIN mới truy cập được
				superOnly := admin.Group("/")
				superOnly.Use(middleware.RequireRole("SUPER_ADMIN"))
				{
					// tạo trường mới (chỉ SUPER_ADMIN)
					superOnly.POST("/schools", schoolHandler.Create)
					superOnly.PUT("/schools/:school_id", schoolHandler.Update)
					superOnly.DELETE("/schools/:school_id", schoolHandler.Delete)

					// gán role cho user (chỉ SUPER_ADMIN — tránh SCHOOL_ADMIN tự nâng quyền)
					superOnly.POST("/users/:user_id/roles", userHandler.AssignRole)

					// quản lý school admins (chỉ SUPER_ADMIN)
					schoolAdmins := superOnly.Group("/school-admins")
					{
						schoolAdmins.POST("", schoolAdminHandler.Create)
						schoolAdmins.GET("", schoolAdminHandler.List)
						schoolAdmins.DELETE("/:admin_id", schoolAdminHandler.Delete)
					}
				}
			}
		}
	}

	return r
}
