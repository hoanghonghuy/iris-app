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

// MyClasses trả về danh sách các lớp mà giáo viên được phân công giảng dạy.
func (h *TeacherScopeHandler) MyClasses(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	userID, ok := shared.RequireCurrentUserID(c)
	if !ok {
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

	userID, ok := shared.RequireCurrentUserID(c)
	if !ok {
		return
	}

	students, err := h.teacherScopeService.ListMyStudentsInClass(ctx, userID, classID)
	if err != nil {
		if errors.Is(err, service.ErrInvalidUserID) || errors.Is(err, service.ErrInvalidClassID) {
			response.Fail(c, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, service.ErrForbidden) {
			response.OK(c, []any{})
			return
		}
		response.Fail(c, http.StatusInternalServerError, "failed to fetch students")
		return
	}

	response.OK(c, students)
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

	userID, ok := shared.RequireCurrentUserID(c)
	if !ok {
		return
	}

	err := h.teacherScopeService.UpdateMyProfile(ctx, userID, req.Phone)
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
