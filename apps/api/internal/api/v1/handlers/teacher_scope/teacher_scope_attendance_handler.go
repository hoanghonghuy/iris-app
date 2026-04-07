package teacherscope

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/api/v1/handlers/shared"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/response"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/service"
)

// MarkAttendance đánh dấu hoặc cập nhật điểm danh cho học sinh
// Teacher chỉ có thể điểm danh cho sinh viên trong các lớp được phân công.
func (h *TeacherScopeHandler) MarkAttendance(c *gin.Context) {
	var req MarkAttendanceRequest
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

	userID, ok := shared.RequireCurrentUserID(c)
	if !ok {
		return
	}

	var checkIn *time.Time
	if req.CheckInAt != nil && *req.CheckInAt != "" {
		t, err := time.Parse(time.RFC3339, *req.CheckInAt)
		if err != nil {
			response.Fail(c, http.StatusBadRequest, "invalid check_in_at (RFC3339)")
			return
		}
		checkIn = &t
	}

	var checkOut *time.Time
	if req.CheckOutAt != nil && *req.CheckOutAt != "" {
		t, err := time.Parse(time.RFC3339, *req.CheckOutAt)
		if err != nil {
			response.Fail(c, http.StatusBadRequest, "invalid check_out_at (RFC3339)")
			return
		}
		checkOut = &t
	}

	err = h.teacherScopeService.UpsertAttendance(ctx, userID, studentID, req.Date, req.Status, checkIn, checkOut, req.Note)
	if err != nil {
		if errors.Is(err, service.ErrInvalidUserID) || errors.Is(err, service.ErrInvalidDate) || errors.Is(err, service.ErrInvalidStatus) {
			response.Fail(c, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, service.ErrForbidden) {
			response.Fail(c, http.StatusForbidden, "forbidden: you can only mark attendance for students in your assigned classes")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "failed to mark attendance")
		return
	}

	response.OK(c, gin.H{
		"message":      "attendance marked successfully",
		"student_id":   studentID.String(),
		"date":         req.Date,
		"status":       req.Status,
		"check_in_at":  req.CheckInAt,
		"check_out_at": req.CheckOutAt,
		"note":         req.Note,
	})
}

// CancelAttendance hủy điểm danh đã lưu của một học sinh trong ngày.
func (h *TeacherScopeHandler) CancelAttendance(c *gin.Context) {
	studentIDRaw := c.Query("student_id")
	dateRaw := c.Query("date")

	if studentIDRaw == "" || dateRaw == "" {
		response.Fail(c, http.StatusBadRequest, "student_id and date are required")
		return
	}

	studentID, err := uuid.Parse(studentIDRaw)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid student_id")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	userID, ok := shared.RequireCurrentUserID(c)
	if !ok {
		return
	}

	err = h.teacherScopeService.CancelAttendanceForDate(ctx, userID, studentID, dateRaw)
	if err != nil {
		if errors.Is(err, service.ErrInvalidUserID) || errors.Is(err, service.ErrInvalidDate) {
			response.Fail(c, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, service.ErrForbidden) {
			response.Fail(c, http.StatusForbidden, "forbidden: you can only cancel attendance for students in your assigned classes")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "failed to cancel attendance")
		return
	}

	response.OK(c, gin.H{
		"message":    "attendance canceled successfully",
		"student_id": studentID.String(),
		"date":       dateRaw,
	})
}

func (h *TeacherScopeHandler) ListAttendance(c *gin.Context) {
	studentID, err := uuid.Parse(c.Param("student_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid student_id")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	userID, ok := shared.RequireCurrentUserID(c)
	if !ok {
		return
	}

	// Default 30 ngày
	to := time.Now()
	from := to.AddDate(0, 0, -30)

	if v := c.Query("from"); v != "" {
		if t, e := time.Parse("2006-01-02", v); e == nil {
			from = t
		}
	}
	if v := c.Query("to"); v != "" {
		if t, e := time.Parse("2006-01-02", v); e == nil {
			to = t
		}
	}

	attendanceRecords, err := h.teacherScopeService.ListAttendanceByStudent(ctx, userID, studentID, from, to)
	if err != nil {
		if errors.Is(err, service.ErrInvalidUserID) {
			response.Fail(c, http.StatusBadRequest, "invalid user ID")
			return
		}
		if errors.Is(err, service.ErrForbidden) {
			response.Fail(c, http.StatusForbidden, "forbidden: you can only view attendance for students in your assigned classes")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "failed to fetch attendance records")
		return
	}

	response.OK(c, attendanceRecords)
}

func (h *TeacherScopeHandler) ListAttendanceChangeLogs(c *gin.Context) {
	studentID, err := uuid.Parse(c.Param("student_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid student_id")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	userID, ok := shared.RequireCurrentUserID(c)
	if !ok {
		return
	}

	to := time.Now()
	from := to.AddDate(0, 0, -30)

	if v := c.Query("from"); v != "" {
		if t, e := time.Parse("2006-01-02", v); e == nil {
			from = t
		}
	}
	if v := c.Query("to"); v != "" {
		if t, e := time.Parse("2006-01-02", v); e == nil {
			to = t
		}
	}

	logs, err := h.teacherScopeService.ListAttendanceChangeLogsByStudent(ctx, userID, studentID, from, to)
	if err != nil {
		if errors.Is(err, service.ErrInvalidUserID) {
			response.Fail(c, http.StatusBadRequest, "invalid user ID")
			return
		}
		if errors.Is(err, service.ErrForbidden) {
			response.Fail(c, http.StatusForbidden, "forbidden: you can only view attendance change logs for students in your assigned classes")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "failed to fetch attendance change logs")
		return
	}

	response.OK(c, logs)
}

func (h *TeacherScopeHandler) ListClassAttendanceChangeLogs(c *gin.Context) {
	classID, err := uuid.Parse(c.Param("class_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid class_id")
		return
	}

	var req shared.PaginationParams
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid query parameters")
		return
	}

	status := c.Query("status")
	if status == "" {
		status = ""
	}

	studentIDRaw := c.Query("student_id")
	var studentID *uuid.UUID
	if studentIDRaw != "" {
		parsedStudentID, err := uuid.Parse(studentIDRaw)
		if err != nil {
			response.Fail(c, http.StatusBadRequest, "invalid student_id")
			return
		}
		studentID = &parsedStudentID
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	userID, ok := shared.RequireCurrentUserID(c)
	if !ok {
		return
	}

	to := time.Now()
	from := to.AddDate(0, 0, -30)

	if v := c.Query("from"); v != "" {
		if t, e := time.Parse("2006-01-02", v); e == nil {
			from = t
		}
	}
	if v := c.Query("to"); v != "" {
		if t, e := time.Parse("2006-01-02", v); e == nil {
			to = t
		}
	}

	var statusPtr *string
	if status != "" {
		statusPtr = &status
	}

	logs, total, err := h.teacherScopeService.ListAttendanceChangeLogsByClass(ctx, userID, classID, studentID, statusPtr, from, to, req.Limit, req.Offset)
	if err != nil {
		if errors.Is(err, service.ErrInvalidUserID) || errors.Is(err, service.ErrInvalidClassID) || errors.Is(err, service.ErrInvalidStatus) {
			response.Fail(c, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, service.ErrForbidden) {
			response.Fail(c, http.StatusForbidden, "forbidden")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "failed to fetch class attendance change logs")
		return
	}

	response.OKPaginated(c, logs, response.Pagination{
		Total:   total,
		Limit:   req.Limit,
		Offset:  req.Offset,
		HasMore: req.Offset+len(logs) < total,
	})
}
