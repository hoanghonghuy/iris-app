package handlers

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/auth"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/middleware"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/response"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/service"
)

type TeacherScopeHandler struct {
	teacherScopeService *service.TeacherScopeService
}

func NewTeacherScopeHandler(teacherScopeService *service.TeacherScopeService) *TeacherScopeHandler {
	return &TeacherScopeHandler{
		teacherScopeService: teacherScopeService,
	}
}

type MarkAttendanceRequest struct {
	StudentID  string  `json:"student_id" binding:"required"`
	Date       string  `json:"date" binding:"required"`   // YYYY-MM-DD
	Status     string  `json:"status" binding:"required"` // present/absent/late/excused
	CheckInAt  *string `json:"check_in_at,omitempty"`     // RFC3339 or empty
	CheckOutAt *string `json:"check_out_at,omitempty"`    // RFC3339 or empty
	Note       string  `json:"note"`
}

type CreateHealthRequest struct {
	StudentID   string   `json:"student_id" binding:"required"`
	RecordedAt  *string  `json:"recorded_at"` // RFC3339 optional
	Temperature *float64 `json:"temperature"`
	Symptoms    string   `json:"symptoms"`
	Severity    *string  `json:"severity"` // normal|watch|urgent optional
	Note        string   `json:"note"`
}

type CreatePostRequest struct {
	ScopeType string `json:"scope_type" binding:"required"` // class|student
	ClassID   string `json:"class_id"`                      // required if scope_type=class
	StudentID string `json:"student_id"`                    // required if scope_type=student
	Type      string `json:"type" binding:"required"`       // announcement|activity|daily_note|health_note
	Content   string `json:"content" binding:"required"`
}

// PaginationParams input chung cho phân trang (dùng cho tất cả list endpoints)
type PaginationParams struct {
	Limit  int `form:"limit"`
	Offset int `form:"offset"`
}

// UpdateMyProfileRequest input để giáo viên cập nhật thông tin cá nhân (chỉ phone)
type UpdateMyProfileRequest struct {
	Phone string `json:"phone"`
}

// MyClasses trả về danh sách các lớp mà giáo viên được phân công giảng dạy.
func (h *TeacherScopeHandler) MyClasses(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	// Get userID from JWT claims
	claimsAny, exists := c.Get(middleware.CtxClaims)
	if !exists {
		response.Fail(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	claims := claimsAny.(*auth.Claims)

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid user ID")
		return
	}

	classes, err := h.teacherScopeService.ListMyClasses(ctx, userID)
	if err != nil {
		if errors.Is(err, service.ErrInvalidUserID) {
			response.Fail(c, http.StatusBadRequest, "invalid user ID")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "failed to fetch classes")
		return
	}

	response.OK(c, classes)
}

// MyStudentsInClass trả về danh sách học sinh trong một lớp nếu giáo viên đó được phân công dạy lớp đó.
func (h *TeacherScopeHandler) MyStudentsInClass(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	classID, err := uuid.Parse(c.Param("class_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid class_id")
		return
	}

	// Get userID from JWT claims
	claimsAny, exists := c.Get(middleware.CtxClaims)
	if !exists {
		response.Fail(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	claims := claimsAny.(*auth.Claims)

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid user ID")
		return
	}

	students, err := h.teacherScopeService.ListMyStudentsInClass(ctx, userID, classID)
	if err != nil {
		if errors.Is(err, service.ErrInvalidUserID) || errors.Is(err, service.ErrInvalidClassID) {
			response.Fail(c, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, service.ErrForbidden) {
			response.Fail(c, http.StatusForbidden, "forbidden: you can only view students in your assigned classes")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "failed to fetch students")
		return
	}

	response.OK(c, students)
}

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

	// Get userID from JWT claims
	claimsAny, exists := c.Get(middleware.CtxClaims)
	if !exists {
		response.Fail(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	claims := claimsAny.(*auth.Claims)

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid user ID")
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

// UpdateMyProfile cập nhật hồ sơ cá nhân của giáo viên (teacher only - chỉ có thể cập nhật số điện thoại)
func (h *TeacherScopeHandler) UpdateMyProfile(c *gin.Context) {
	var req UpdateMyProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid request body")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	// Get userID from JWT claims
	claimsAny, exists := c.Get(middleware.CtxClaims)
	if !exists {
		response.Fail(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	claims := claimsAny.(*auth.Claims)

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid user ID")
		return
	}

	err = h.teacherScopeService.UpdateMyProfile(ctx, userID, req.Phone)
	if err != nil {
		if errors.Is(err, service.ErrInvalidUserID) {
			response.Fail(c, http.StatusBadRequest, "invalid user ID")
			return
		}
		if errors.Is(err, service.ErrTeacherNotFound) {
			response.Fail(c, http.StatusNotFound, "teacher profile not found")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "failed to update profile")
		return
	}

	response.OK(c, gin.H{
		"message": "profile updated successfully",
		"phone":   req.Phone,
	})
}

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

	// Get userID from JWT claims
	claimsAny, exists := c.Get(middleware.CtxClaims)
	if !exists {
		response.Fail(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	claims := claimsAny.(*auth.Claims)

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid user ID")
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

	// Get userID from JWT claims
	claimsAny, exists := c.Get(middleware.CtxClaims)
	if !exists {
		response.Fail(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	claims := claimsAny.(*auth.Claims)

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid user ID")
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

func (h *TeacherScopeHandler) ListAttendance(c *gin.Context) {
	studentID, err := uuid.Parse(c.Param("student_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid student_id")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	claimsAny, exists := c.Get(middleware.CtxClaims)
	if !exists {
		response.Fail(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	claims := claimsAny.(*auth.Claims)

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid user ID")
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

// CreatePost tạo bài đăng mới cho lớp hoặc học sinh
func (h *TeacherScopeHandler) CreatePost(c *gin.Context) {
	var req CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid request body")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	// Get userID from JWT claims
	claimsAny, exists := c.Get(middleware.CtxClaims)
	if !exists {
		response.Fail(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	claims := claimsAny.(*auth.Claims)

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid user ID")
		return
	}

	var postID uuid.UUID

	switch req.ScopeType {
	case "class":
		classID, err := uuid.Parse(req.ClassID)
		if err != nil {
			response.Fail(c, http.StatusBadRequest, "invalid class_id")
			return
		}
		postID, err = h.teacherScopeService.CreateClassPost(ctx, userID, classID, req.Type, req.Content)
		if err != nil {
			if errors.Is(err, service.ErrInvalidUserID) || errors.Is(err, service.ErrInvalidClassID) || errors.Is(err, service.ErrInvalidValue) {
				response.Fail(c, http.StatusBadRequest, err.Error())
				return
			}
			if errors.Is(err, service.ErrForbidden) {
				response.Fail(c, http.StatusForbidden, "forbidden: you can only create posts for your assigned classes")
				return
			}
			response.Fail(c, http.StatusInternalServerError, "failed to create post")
			return
		}

	case "student":
		studentID, err := uuid.Parse(req.StudentID)
		if err != nil {
			response.Fail(c, http.StatusBadRequest, "invalid student_id")
			return
		}
		postID, err = h.teacherScopeService.CreateStudentPost(ctx, userID, studentID, req.Type, req.Content)
		if err != nil {
			if errors.Is(err, service.ErrInvalidUserID) || errors.Is(err, service.ErrInvalidValue) {
				response.Fail(c, http.StatusBadRequest, err.Error())
				return
			}
			if errors.Is(err, service.ErrForbidden) {
				response.Fail(c, http.StatusForbidden, "forbidden: you can only create posts for students in your assigned classes")
				return
			}
			response.Fail(c, http.StatusInternalServerError, "failed to create post")
			return
		}

	default:
		response.Fail(c, http.StatusBadRequest, "invalid scope_type (class|student)")
		return
	}

	response.Created(c, gin.H{
		"message":    "post created successfully",
		"post_id":    postID.String(),
		"scope_type": req.ScopeType,
		"type":       req.Type,
	})
}

// ListClassPosts liệt kê bài đăng của một lớp học
func (h *TeacherScopeHandler) ListClassPosts(c *gin.Context) {
	classID, err := uuid.Parse(c.Param("class_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid class_id")
		return
	}

	var req PaginationParams
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid query parameters")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	// Get userID from JWT claims
	claimsAny, exists := c.Get(middleware.CtxClaims)
	if !exists {
		response.Fail(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	claims := claimsAny.(*auth.Claims)

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid user ID")
		return
	}

	posts, total, err := h.teacherScopeService.ListClassPosts(ctx, userID, classID, req.Limit, req.Offset)
	if err != nil {
		if errors.Is(err, service.ErrInvalidUserID) || errors.Is(err, service.ErrInvalidClassID) {
			response.Fail(c, http.StatusBadRequest, "invalid user ID or class ID")
			return
		}
		if errors.Is(err, service.ErrForbidden) {
			response.Fail(c, http.StatusForbidden, "not assigned to this class")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "failed to list class posts")
		return
	}

	response.OKPaginated(c, posts, response.Pagination{
		Total:   total,
		Limit:   req.Limit,
		Offset:  req.Offset,
		HasMore: req.Offset+len(posts) < total,
	})
}

// ListStudentPosts liệt kê bài đăng của một học sinh
func (h *TeacherScopeHandler) ListStudentPosts(c *gin.Context) {
	studentID, err := uuid.Parse(c.Param("student_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid student_id")
		return
	}

	var req PaginationParams
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid query parameters")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	// Get userID from JWT claims
	claimsAny, exists := c.Get(middleware.CtxClaims)
	if !exists {
		response.Fail(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	claims := claimsAny.(*auth.Claims)

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid user ID")
		return
	}

	posts, total, err := h.teacherScopeService.ListStudentPosts(ctx, userID, studentID, req.Limit, req.Offset)
	if err != nil {
		if errors.Is(err, service.ErrInvalidUserID) {
			response.Fail(c, http.StatusBadRequest, "invalid user ID or student ID")
			return
		}
		if errors.Is(err, service.ErrForbidden) {
			response.Fail(c, http.StatusForbidden, "not assigned to this student's class")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "failed to list student posts")
		return
	}

	response.OKPaginated(c, posts, response.Pagination{
		Total:   total,
		Limit:   req.Limit,
		Offset:  req.Offset,
		HasMore: req.Offset+len(posts) < total,
	})
}
