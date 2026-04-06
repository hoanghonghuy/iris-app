package teacherhandlers

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

func (h *TeacherHandler) Assign(c *gin.Context) {
	adminSchoolID := extractAdminSchoolID(c)

	teacherID, err := uuid.Parse(c.Param("teacher_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid teacher_id format")
		return
	}
	classID, err := uuid.Parse(c.Param("class_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid class_id format")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	err = h.teacherService.Assign(ctx, adminSchoolID, teacherID, classID)
	if err != nil {
		if errors.Is(err, service.ErrSchoolAccessDenied) {
			response.Fail(c, http.StatusForbidden, "access denied")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "failed to assign teacher to class")
		return
	}
	response.OK(c, gin.H{
		"message":    "teacher assigned to class successfully",
		"teacher_id": teacherID,
		"class_id":   classID,
	})
}

func (h *TeacherHandler) Unassign(c *gin.Context) {
	adminSchoolID := extractAdminSchoolID(c)

	teacherID, err := uuid.Parse(c.Param("teacher_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid teacher_id format")
		return
	}
	classID, err := uuid.Parse(c.Param("class_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid class_id format")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	err = h.teacherService.Unassign(ctx, adminSchoolID, teacherID, classID)
	if err != nil {
		if errors.Is(err, service.ErrSchoolAccessDenied) {
			response.Fail(c, http.StatusForbidden, "access denied")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "failed to unassign teacher from class")
		return
	}
	response.OK(c, gin.H{
		"message":    "teacher unassigned from class successfully",
		"teacher_id": teacherID,
		"class_id":   classID,
	})
}
