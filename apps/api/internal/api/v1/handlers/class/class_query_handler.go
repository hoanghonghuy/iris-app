package classhandlers

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

// ListBySchool lấy danh sách lớp theo trường
func (h *ClassHandler) ListBySchool(c *gin.Context) {
	adminSchoolID := extractAdminSchoolID(c)

	schoolID, err := uuid.Parse(c.Param("school_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid school_id format")
		return
	}

	var params PaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid pagination params")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	classes, total, err := h.classService.ListBySchool(ctx, adminSchoolID, schoolID, params.Limit, params.Offset)
	if err != nil {
		if errors.Is(err, service.ErrSchoolAccessDenied) {
			response.Fail(c, http.StatusForbidden, "access denied")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "failed to fetch classes")
		return
	}

	response.OKPaginated(c, classes, response.Pagination{
		Total:   total,
		Limit:   params.Limit,
		Offset:  params.Offset,
		HasMore: params.Offset+len(classes) < total,
	})
}
