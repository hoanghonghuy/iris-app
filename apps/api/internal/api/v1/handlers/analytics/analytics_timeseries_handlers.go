package analyticshandlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/api/v1/handlers/shared"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/response"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/service"
)

// AdminAnalyticsTimeseries GET /admin/analytics/timeseries (additive).
func (h *AnalyticsHandler) AdminAnalyticsTimeseries(c *gin.Context) {
	effectiveSchoolID := shared.ExtractAdminSchoolID(c)
	if effectiveSchoolID == nil && c.Query("school_id") != "" {
		id, err := uuid.Parse(c.Query("school_id"))
		if err != nil {
			response.Fail(c, http.StatusBadRequest, "Tham số school_id không hợp lệ")
			return
		}
		effectiveSchoolID = &id
	}

	stats, err := h.analyticsService.GetAdminAnalyticsTimeseries(
		c.Request.Context(),
		effectiveSchoolID,
		c.Query("range"),
		c.Query("interval"),
	)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrAnalyticsInvalidRange),
			errors.Is(err, service.ErrAnalyticsInvalidInterval):
			response.Fail(c, http.StatusBadRequest, err.Error())
			return
		default:
			response.Fail(c, http.StatusInternalServerError, "Lỗi khi thống kê timeseries: "+err.Error())
			return
		}
	}

	response.OK(c, stats)
}

// TeacherAnalyticsTimeseries GET /teacher/analytics/timeseries.
func (h *AnalyticsHandler) TeacherAnalyticsTimeseries(c *gin.Context) {
	userID, ok := shared.RequireCurrentUserID(c)
	if !ok {
		return
	}

	stats, err := h.analyticsService.GetTeacherAnalyticsTimeseries(
		c.Request.Context(),
		userID,
		c.Query("range"),
		c.Query("interval"),
	)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrAnalyticsInvalidRange),
			errors.Is(err, service.ErrAnalyticsInvalidInterval):
			response.Fail(c, http.StatusBadRequest, err.Error())
			return
		default:
			response.Fail(c, http.StatusInternalServerError, "Lỗi khi thống kê timeseries giáo viên: "+err.Error())
			return
		}
	}

	response.OK(c, stats)
}

// ParentAnalyticsTimeseries GET /parent/analytics/timeseries (bắt buộc student_id).
func (h *AnalyticsHandler) ParentAnalyticsTimeseries(c *gin.Context) {
	parentUserID, ok := shared.RequireCurrentUserID(c)
	if !ok {
		return
	}

	sid := c.Query("student_id")
	if sid == "" {
		response.Fail(c, http.StatusBadRequest, "Thiếu tham số student_id")
		return
	}
	studentID, err := uuid.Parse(sid)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "Tham số student_id không hợp lệ")
		return
	}

	stats, err := h.analyticsService.GetParentAnalyticsTimeseries(
		c.Request.Context(),
		parentUserID,
		studentID,
		c.Query("range"),
		c.Query("interval"),
	)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrForbidden):
			response.Fail(c, http.StatusForbidden, err.Error())
			return
		case errors.Is(err, service.ErrInvalidUserID):
			response.Fail(c, http.StatusBadRequest, err.Error())
			return
		case errors.Is(err, service.ErrAnalyticsInvalidRange),
			errors.Is(err, service.ErrAnalyticsInvalidInterval):
			response.Fail(c, http.StatusBadRequest, err.Error())
			return
		default:
			response.Fail(c, http.StatusInternalServerError, "Lỗi khi thống kê timeseries phụ huynh: "+err.Error())
			return
		}
	}

	response.OK(c, stats)
}
