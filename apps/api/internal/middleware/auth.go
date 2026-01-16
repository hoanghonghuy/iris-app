package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/auth"
)
// lưu trữ thông tin claims của JWT vào context của Gin
const CtxClaims = "claims"

// AuthJWT là một middleware dùng để xác thực JWT trong các request đến server.
// 	- Hàm nhận vào một chuỗi secret dùng để giải mã và xác thực token.
// 	- Middleware sẽ kiểm tra header "Authorization" của request, đảm bảo có định dạng "Bearer <token>".
// 	- Nếu không có hoặc token không hợp lệ, middleware trả về lỗi 401 Unauthorized.
// 	- Nếu token hợp lệ, middleware sẽ giải mã và lưu thông tin claims vào context của Gin với key CtxClaims.
// 	- Các middleware hoặc handler phía sau có thể lấy thông tin claims này để kiểm tra phân quyền hoặc lấy thông tin người dùng.
func AuthJWT(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if !strings.HasPrefix(h, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
			return
		}
		tokenStr := strings.TrimPrefix(h, "Bearer ")
		claims, err := auth.Parse(secret, tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
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
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing claims"})
			return
		}
		claims := v.(*auth.Claims)
		
		for _, r := range claims.Roles {
			if r == role {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
	}
}
