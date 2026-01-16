package httpapi

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/auth"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/middleware"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/repo"
)

type Handlers struct {
	UserRepo  *repo.UserRepo
	JWTSecret string
	TTL       time.Duration
}

type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// NewRouter tạo và cấu hình một router HTTP sử dụng Gin.
// Hàm này nhận vào một UserRepo để thao tác với dữ liệu người dùng,
// một chuỗi bí mật JWT để ký token, và thời gian sống (TTL) của token tính bằng phút.
// Các middleware được sử dụng để kiểm tra JWT và quyền truy cập.
func NewRouter(repos *repo.Repositories, jwtSecret string, ttlMinutes int) *gin.Engine {
	r := gin.Default()
	h := &Handlers{
		UserRepo:  repos.UserRepo,
		JWTSecret: jwtSecret,
		TTL:       time.Duration(ttlMinutes) * time.Minute,
	}

	adminSchoolHandler := &AdminSchoolHandler{
		Schools: repos.SchoolRepo,
	}

	v1 := r.Group("/api/v1")
	{
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"ok": true,
			})
		})
		v1.POST("/auth/login", h.Login)

		// route yêu cầu xác thực (JWT)
		protected := v1.Group("/")
		protected.Use(middleware.AuthJWT(jwtSecret))
		protected.GET("/me", h.Me)

		// route admin (yêu cầu role "admin")
		admin := v1.Group("/admin")
		admin.Use(middleware.RequireRole("ADMIN"))
		admin.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"pong": "admin",
			})
		})

		schools := admin.Group("/schools")
		schools.POST("", adminSchoolHandler.Create)
		schools.GET("", adminSchoolHandler.List)
	}

	return r
}

// Login là handler xử lý yêu cầu đăng nhập của người dùng.
//   - Nhận vào JSON gồm email và password.
//   - Kiểm tra tính hợp lệ của dữ liệu đầu vào.
//   - Tìm người dùng theo email trong cơ sở dữ liệu.
//   - Kiểm tra trạng thái tài khoản (ví dụ: bị khóa).
//   - Xác thực mật khẩu bằng cách so sánh với hash lưu trong DB.
//   - Lấy danh sách vai trò (roles) của người dùng.
//   - Nếu hợp lệ, tạo JWT token chứa thông tin người dùng và roles.
//   - Trả về access_token, loại token và thời gian hết hạn.
//
// Các trường hợp lỗi được xử lý gồm:
//   - Thiếu email/password hoặc dữ liệu không hợp lệ: trả về 400.
//   - Không tìm thấy người dùng hoặc mật khẩu sai: trả về 401.
//   - Tài khoản bị khóa: trả về 403.
//   - Lỗi hệ thống hoặc DB: trả về 500.
func (h *Handlers) Login(c *gin.Context) {
	var req LoginReq
	if err := c.ShouldBindJSON(&req); err != nil || req.Email == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email/password required"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	u, err := h.UserRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		c.JSON(500, gin.H{"error": "db error"})
		return
	}
	if u.Status == "locked" {
		c.JSON(http.StatusForbidden, gin.H{"error": "user locked"})
		return
	}
	if !auth.VerifyPassword(u.PasswordHash, req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	roles, err := h.UserRepo.RolesOfUser(ctx, u.ID)
	if err != nil {
		c.JSON(500, gin.H{"error": "db error"})
		return
	}

	token, err := auth.Sign(h.JWTSecret, h.TTL, u.ID.String(), u.Email, roles)
	if err != nil {
		c.JSON(500, gin.H{"error": "token error"})
		return
	}

	c.JSON(200, gin.H{
		"access_token": token,
		"token_type":   "Bearer",
		"expires_in":   int(h.TTL.Seconds()),
	})
}

func (h *Handlers) Me(c *gin.Context) {
	claimsAny, _ := c.Get(middleware.CtxClaims)
	claims := claimsAny.(*auth.Claims)
	c.JSON(200, gin.H{"user_id": claims.UserID, "email": claims.Email, "roles": claims.Roles})
}
