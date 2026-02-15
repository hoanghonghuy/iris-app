package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/middleware"
)

// extractAdminSchoolID đọc school_id từ context (InjectAdminScope đã set).
//   - SUPER_ADMIN → nil (không giới hạn trường)
//   - SCHOOL_ADMIN → *uuid.UUID (chỉ trường mình)
func extractAdminSchoolID(c *gin.Context) *uuid.UUID {
	// Lấy school_id từ context. InjectAdminScope đã lưu:
	//   - "" (rỗng) cho SUPER_ADMIN
	//   - "uuid-string" cho SCHOOL_ADMIN
	// Nếu key không tồn tại → request chưa qua InjectAdminScope (lỗi config route) → trả nil.
	value, exists := c.Get(middleware.CtxAdminSchoolID)
	if !exists {
		return nil
	}

	// Type assertion: any → string.
	// Nếu rỗng "" → SUPER_ADMIN → trả nil (không giới hạn trường).
	schoolIDStr, ok := value.(string)
	if !ok || schoolIDStr == "" {
		return nil
	}

	// Parse string → uuid.UUID cho SCHOOL_ADMIN.
	schoolID, err := uuid.Parse(schoolIDStr)
	if err != nil {
		return nil
	}
	return &schoolID
}
