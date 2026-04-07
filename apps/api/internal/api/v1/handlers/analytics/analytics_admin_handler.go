package analyticshandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/api/v1/handlers/shared"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/response"
)

// AdminDashboardStats lấy thống kê cho Admin.
func (h *AnalyticsHandler) AdminDashboardStats(c *gin.Context) {
	schoolID := shared.ExtractAdminSchoolID(c) // Trả về nil nếu là SUPER_ADMIN

	stats, err := h.analyticsService.GetAdminAnalytics(c.Request.Context(), schoolID)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "Lỗi khi thống kê dữ liệu: "+err.Error())
		return
	}

	response.OK(c, stats)
}
