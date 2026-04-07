package teacherhandlers

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

// Update updates a teacher's information (admin only - can update all fields)
func (h *TeacherHandler) Update(c *gin.Context) {
	adminSchoolID := shared.ExtractAdminSchoolID(c)

	teacherID, err := uuid.Parse(c.Param("teacher_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid teacher_id format")
		return
	}

	var req UpdateTeacherRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid request body")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	err = h.teacherService.Update(ctx, adminSchoolID, teacherID, req.FullName, req.Phone, req.SchoolID)
	if err != nil {
		if errors.Is(err, service.ErrSchoolAccessDenied) {
			response.Fail(c, http.StatusForbidden, "access denied")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "failed to update teacher")
		return
	}

	response.OK(c, gin.H{
		"message":    "teacher updated successfully",
		"teacher_id": teacherID.String(),
	})
}

// Delete xóa giáo viên
func (h *TeacherHandler) Delete(c *gin.Context) {
	adminSchoolID := shared.ExtractAdminSchoolID(c)

	teacherID, err := uuid.Parse(c.Param("teacher_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid teacher_id format")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	if err := h.teacherService.Delete(ctx, adminSchoolID, teacherID); err != nil {
		if errors.Is(err, service.ErrSchoolAccessDenied) {
			response.Fail(c, http.StatusForbidden, "access denied")
			return
		}
		if errors.Is(err, service.ErrTeacherNotFound) {
			response.Fail(c, http.StatusNotFound, "teacher not found")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "failed to delete teacher")
		return
	}

	response.OK(c, gin.H{"message": "teacher deleted successfully", "teacher_id": teacherID.String()})
}
