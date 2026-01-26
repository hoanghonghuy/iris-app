package handlers

import (
	"context"
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
	TeacherScopeService *service.TeacherScopeService
}

type MarkAttendanceRequest struct {
	StudentID  string     `json:"student_id" binding:"required"`
	Date       string     `json:"date" binding:"required"`   // YYYY-MM-DD
	Status     string     `json:"status" binding:"required"` // present/absent/late/excused
	CheckInAt  *time.Time `json:"check_in_at,omitempty"`     // optional
	CheckOutAt *time.Time `json:"check_out_at,omitempty"`    // optional
	Note       string     `json:"note"`
}

// MyClasses returns list of classes that the teacher is assigned to teach
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

	classes, err := h.TeacherScopeService.ListMyClasses(ctx, userID)
	if err != nil {
		if err == service.ErrInvalidUserID {
			response.Fail(c, http.StatusBadRequest, "invalid user ID")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "failed to fetch classes")
		return
	}

	response.OK(c, classes)
}

// MyStudentsInClass returns list of students in a class if the teacher is assigned to that class
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

	students, err := h.TeacherScopeService.ListMyStudentsInClass(ctx, userID, classID)
	if err != nil {
		if err == service.ErrInvalidUserID || err == service.ErrInvalidClassID {
			response.Fail(c, http.StatusBadRequest, err.Error())
			return
		}
		if err == service.ErrForbidden {
			response.Fail(c, http.StatusForbidden, "forbidden: you can only view students in your assigned classes")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "failed to fetch students")
		return
	}

	response.OK(c, students)
}

// MarkAttendance marks or updates attendance for a student
// Teacher can only mark attendance for students in their assigned classes
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

	err = h.TeacherScopeService.UpsertAttendance(ctx, userID, studentID, req.Date, req.Status, req.CheckInAt, req.CheckOutAt, req.Note)
	if err != nil {
		if err == service.ErrInvalidUserID || err == service.ErrInvalidDate || err == service.ErrInvalidStatus {
			response.Fail(c, http.StatusBadRequest, err.Error())
			return
		}
		if err == service.ErrForbidden {
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

// UpdateMyProfile updates teacher's own profile (teacher only - can only update phone)
func (h *TeacherScopeHandler) UpdateMyProfile(c *gin.Context) {
	var req service.UpdateMyProfileRequest
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

	err = h.TeacherScopeService.UpdateMyProfile(ctx, userID, req)
	if err != nil {
		if err == service.ErrInvalidUserID {
			response.Fail(c, http.StatusBadRequest, "invalid user ID")
			return
		}
		if err == service.ErrTeacherNotFound {
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
