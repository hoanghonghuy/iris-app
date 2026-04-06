package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/response"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/service"
)

type AnalyticsHandler struct {
	analyticsService *service.AnalyticsService
}

func NewAnalyticsHandler(analyticsService *service.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{
		analyticsService: analyticsService,
	}
}

// AdminDashboardStats lấy thống kê cho Admin.
func (h *AnalyticsHandler) AdminDashboardStats(c *gin.Context) {
	schoolID := extractAdminSchoolID(c) // Trả về nil nếu là SUPER_ADMIN

	stats, err := h.analyticsService.GetAdminAnalytics(c.Request.Context(), schoolID)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "Lỗi khi thống kê dữ liệu: "+err.Error())
		return
	}

	response.OK(c, stats)
}

// TeacherDashboardStats lấy thống kê cho Giáo viên.
func (h *AnalyticsHandler) TeacherDashboardStats(c *gin.Context) {
	userID, ok := requireCurrentUserID(c)
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
