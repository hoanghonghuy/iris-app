package authhandlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/api/v1/handlers/shared"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/response"
)

// Me trả về thông tin user đã đăng nhập.
func (h *AuthHandler) Me(c *gin.Context) {
	claims, ok := shared.RequireCurrentClaims(c)
	if !ok {
		return
	}

	// Trả về thông tin từ JWT claims (đã validate và xác thực).
	result := gin.H{
		"user_id": claims.UserID,
		"email":   claims.Email,
		"roles":   claims.Roles,
	}

	// Lấy thêm full_name từ Database do JWT không chứa.
	if uid, err := uuid.Parse(claims.UserID); err == nil {
		if userInfo, err := h.userService.FindByID(c.Request.Context(), nil, uid); err == nil {
			result["full_name"] = userInfo.FullName
		}
	}

	// Nếu user là SCHOOL_ADMIN → trả thêm school_id.
	if claims.SchoolID != "" {
		result["school_id"] = claims.SchoolID
	}

	response.OK(c, result)
}
