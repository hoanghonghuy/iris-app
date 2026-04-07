package authhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/auth"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/middleware"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/response"
)

// requireCurrentClaims lấy JWT claims từ context và tự trả lỗi HTTP nếu claims không hợp lệ.
func requireCurrentClaims(c *gin.Context) (*auth.Claims, bool) {
	claimsAny, exists := c.Get(middleware.CtxClaims)
	if !exists {
		response.Fail(c, http.StatusUnauthorized, "unauthorized")
		return nil, false
	}

	claims, ok := claimsAny.(*auth.Claims)
	if !ok {
		response.Fail(c, http.StatusUnauthorized, "unauthorized")
		return nil, false
	}

	return claims, true
}
