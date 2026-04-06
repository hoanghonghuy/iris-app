package teacherscope

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/response"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/service"
)

// CreateHealth tạo nhật ký sức khỏe mới cho học sinh
func (h *TeacherScopeHandler) CreateHealth(c *gin.Context) {
	var req CreateHealthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid request body")
		return
	}

	studentID, err := uuid.Parse(req.StudentID)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid student_id")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	userID, ok := requireCurrentUserID(c)
	if !ok {
		return
	}

	var recordedAt *time.Time
	if req.RecordedAt != nil && *req.RecordedAt != "" {
		t, err := time.Parse(time.RFC3339, *req.RecordedAt)
		if err != nil {
			response.Fail(c, http.StatusBadRequest, "invalid recorded_at (RFC3339)")
			return
		}
		recordedAt = &t
	}

	id, err := h.teacherScopeService.CreateHealthLog(ctx, userID, studentID, recordedAt, req.Temperature, req.Symptoms, req.Note, req.Severity)
	if err != nil {
		if errors.Is(err, service.ErrInvalidUserID) || errors.Is(err, service.ErrInvalidValue) {
			response.Fail(c, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, service.ErrForbidden) {
			response.Fail(c, http.StatusForbidden, "forbidden: you can only create health logs for students in your assigned classes")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "failed to create health log")
		return
	}

	response.Created(c, gin.H{
		"message":       "health log created successfully",
		"health_log_id": id.String(),
		"student_id":    req.StudentID,
		"recorded_at":   req.RecordedAt,
		"temperature":   req.Temperature,
		"symptoms":      req.Symptoms,
		"severity":      req.Severity,
		"note":          req.Note,
	})
}

// ListHealth liệt kê nhật ký sức khỏe của học sinh nếu giáo viên đó được phân công dạy lớp của học sinh đó.
func (h *TeacherScopeHandler) ListHealth(c *gin.Context) {
	studentID, err := uuid.Parse(c.Param("student_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid student_id")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	userID, ok := requireCurrentUserID(c)
	if !ok {
		return
	}

	// default 30 ngày
	to := time.Now()
	from := to.AddDate(0, 0, -30)

	if v := c.Query("from"); v != "" {
		if t, e := time.Parse("2006-01-02", v); e == nil {
			from = t
		}
	}
	if v := c.Query("to"); v != "" {
		if t, e := time.Parse("2006-01-02", v); e == nil {
			// end-of-day: recorded_at là TIMESTAMP, cần bao gồm cả ngày cuối
			to = t.Add(24*time.Hour - time.Nanosecond)
		}
	}

	healthLogs, err := h.teacherScopeService.ListHealthLogs(ctx, userID, studentID, from, to)
	if err != nil {
		if errors.Is(err, service.ErrInvalidUserID) {
			response.Fail(c, http.StatusBadRequest, "invalid user ID")
			return
		}
		if errors.Is(err, service.ErrForbidden) {
			response.Fail(c, http.StatusForbidden, "forbidden: you can only view health logs for students in your assigned classes")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "failed to fetch health logs")
		return
	}

	response.OK(c, healthLogs)
}
