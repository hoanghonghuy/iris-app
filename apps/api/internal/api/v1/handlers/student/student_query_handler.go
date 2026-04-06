package studenthandlers

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

func (h *StudentHandler) ListByClass(c *gin.Context) {
	adminSchoolID := extractAdminSchoolID(c)

	classID, err := uuid.Parse(c.Param("class_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid class_id format")
		return
	}

	var params PaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid pagination params")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	students, total, err := h.studentService.ListByClass(ctx, adminSchoolID, classID, params.Limit, params.Offset)
	if err != nil {
		if errors.Is(err, service.ErrSchoolAccessDenied) {
			response.Fail(c, http.StatusForbidden, "access denied")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "failed to fetch students")
		return
	}

	response.OKPaginated(c, students, response.Pagination{
		Total:   total,
		Limit:   params.Limit,
		Offset:  params.Offset,
		HasMore: params.Offset+len(students) < total,
	})
}

// GetProfile lấy thông tin chi tiết của một học sinh
func (h *StudentHandler) GetProfile(c *gin.Context) {
	adminSchoolID := extractAdminSchoolID(c)

	studentID, err := uuid.Parse(c.Param("student_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid student_id format")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	profile, err := h.studentService.GetProfile(ctx, adminSchoolID, studentID)
	if err != nil {
		if errors.Is(err, service.ErrSchoolAccessDenied) {
			response.Fail(c, http.StatusForbidden, "access denied")
			return
		}
		if errors.Is(err, service.ErrFailedToGetStudent) {
			response.Fail(c, http.StatusNotFound, "student not found")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "failed to fetch student profile")
		return
	}

	response.OK(c, profile)
}
