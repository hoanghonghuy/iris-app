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

func (h *TeacherHandler) List(c *gin.Context) {
	adminSchoolID := shared.ExtractAdminSchoolID(c)

	var params shared.PaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid pagination params")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	teachers, total, err := h.teacherService.List(ctx, adminSchoolID, params.Limit, params.Offset)
	if err != nil {
		if errors.Is(err, service.ErrSchoolAccessDenied) {
			response.Fail(c, http.StatusForbidden, "access denied")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "failed to fetch teachers")
		return
	}

	response.OKPaginated(c, teachers, response.Pagination{
		Total:   total,
		Limit:   params.Limit,
		Offset:  params.Offset,
		HasMore: params.Offset+len(teachers) < total,
	})
}

func (h *TeacherHandler) ListTeachersOfClass(c *gin.Context) {
	adminSchoolID := shared.ExtractAdminSchoolID(c)

	classID, err := uuid.Parse(c.Param("class_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid class_id format")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	teachers, err := h.teacherService.ListTeachersOfClass(ctx, adminSchoolID, classID)
	if err != nil {
		if errors.Is(err, service.ErrSchoolAccessDenied) {
			response.Fail(c, http.StatusForbidden, "access denied")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "failed to fetch teachers of class")
		return
	}
	response.OK(c, teachers)
}

func (h *TeacherHandler) GetByTeacherID(c *gin.Context) {
	adminSchoolID := shared.ExtractAdminSchoolID(c)

	teacherID, err := uuid.Parse(c.Param("teacher_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid teacher_id format")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	teacher, err := h.teacherService.GetByTeacherID(ctx, adminSchoolID, teacherID)
	if err != nil {
		if errors.Is(err, service.ErrSchoolAccessDenied) {
			response.Fail(c, http.StatusForbidden, "access denied")
			return
		}
		if errors.Is(err, service.ErrTeacherNotFound) {
			response.Fail(c, http.StatusNotFound, "teacher not found")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "failed to fetch teacher")
		return
	}
	response.OK(c, teacher)
}
