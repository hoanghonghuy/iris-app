package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/auth"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/response"
)

// lưu trữ thông tin claims của JWT vào context của Gin
const CtxClaims = "claims"

// lưu trữ school_id của admin (rỗng => SUPER_ADMIN, có giá trị => SCHOOL_ADMIN)
const CtxAdminSchoolID = "admin_school_id"

// AuthJWT là một middleware dùng để xác thực JWT trong các request đến server.
//   - Hàm nhận vào một chuỗi secret dùng để giải mã và xác thực token.
//   - Middleware sẽ kiểm tra header "Authorization" của request, đảm bảo có định dạng "Bearer <token>".
//   - Nếu không có hoặc token không hợp lệ, middleware trả về lỗi 401 Unauthorized.
//   - Nếu token hợp lệ, middleware sẽ giải mã và lưu thông tin claims vào context của Gin với key CtxClaims.
//   - Các middleware hoặc handler phía sau có thể lấy thông tin claims này để kiểm tra phân quyền hoặc lấy thông tin người dùng.
func AuthJWT(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if !strings.HasPrefix(h, "Bearer ") {
			// c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
			response.Fail(c, http.StatusUnauthorized, "missing bearer token")
			c.Abort()
			return
		}
		tokenStr := strings.TrimPrefix(h, "Bearer ")
		claims, err := auth.Parse(secret, tokenStr)
		if err != nil {
			// c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			response.Fail(c, http.StatusUnauthorized, "invalid token")
			c.Abort()
			return
		}
		c.Set(CtxClaims, claims)
		c.Next()
	}
}

func RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		v, ok := c.Get(CtxClaims)
		if !ok {
			response.Fail(c, http.StatusUnauthorized, "unauthorized")
			c.Abort()
			return
		}
		claims := v.(*auth.Claims)

		for _, r := range claims.Roles {
			if r == role {
				c.Next()
				return
			}
		}
		response.Fail(c, http.StatusForbidden, "access denied")
		c.Abort()
	}
}

// RequireAnyRole kiểm tra user có ít nhất 1 trong các roles được chỉ định.
// Dùng cho admin routes mà cả SUPER_ADMIN và SCHOOL_ADMIN đều truy cập được.
func RequireAnyRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		v, ok := c.Get(CtxClaims)
		if !ok {
			response.Fail(c, http.StatusUnauthorized, "unauthorized")
			c.Abort()
			return
		}
		claims := v.(*auth.Claims)

		for _, userRole := range claims.Roles {
			for _, allowed := range roles {
				if userRole == allowed {
					c.Next()
					return
				}
			}
		}
		response.Fail(c, http.StatusForbidden, "access denied")
		c.Abort()
	}
}

// InjectAdminScope đọc SchoolID từ JWT claims và lưu vào context.
//   - SUPER_ADMIN: SchoolID rỗng → context value = "" (không giới hạn trường)
//   - SCHOOL_ADMIN: SchoolID có giá trị → context value = school_id (chỉ truy cập trường mình)
//
// Middleware này cũng validate rằng SCHOOL_ADMIN bắt buộc phải có school_id trong token.
// Nếu không có → trả 403 (trường hợp data không nhất quán).
func InjectAdminScope() gin.HandlerFunc {
	return func(c *gin.Context) {
		v, ok := c.Get(CtxClaims)
		if !ok {
			response.Fail(c, http.StatusUnauthorized, "unauthorized")
			c.Abort()
			return
		}
		claims := v.(*auth.Claims)

		// Kiểm tra SCHOOL_ADMIN phải có school_id
		for _, r := range claims.Roles {
			if r == "SCHOOL_ADMIN" && claims.SchoolID == "" {
				response.Fail(c, http.StatusForbidden, "school admin account not linked to any school")
				c.Abort()
				return
			}
		}

		// Lưu school_id vào context (rỗng nếu SUPER_ADMIN)
		c.Set(CtxAdminSchoolID, claims.SchoolID)
		c.Next()
	}
}
