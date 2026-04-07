package analyticshandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/api/v1/handlers/shared"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/response"
)

// TeacherDashboardStats lấy thống kê cho Giáo viên.
func (h *AnalyticsHandler) TeacherDashboardStats(c *gin.Context) {
	userID, ok := shared.RequireCurrentUserID(c)
	if !ok {
		return
	}

	stats, err := h.analyticsService.GetTeacherAnalytics(c.Request.Context(), userID)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "Lỗi khi thống kê dữ liệu giáo viên: "+err.Error())
		return
	}

	response.OK(c, stats)
}
