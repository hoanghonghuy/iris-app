package parentcodehandlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/middleware"
)

// extractAdminSchoolID đọc school_id từ context (InjectAdminScope đã set).
//   - SUPER_ADMIN -> nil (không giới hạn trường)
//   - SCHOOL_ADMIN -> *uuid.UUID (chỉ trường mình)
func extractAdminSchoolID(c *gin.Context) *uuid.UUID {
	value, exists := c.Get(middleware.CtxAdminSchoolID)
	if !exists {
		return nil
	}

	schoolIDStr, ok := value.(string)
	if !ok || schoolIDStr == "" {
		return nil
	}

	schoolID, err := uuid.Parse(schoolIDStr)
	if err != nil {
		return nil
	}
	return &schoolID
}
