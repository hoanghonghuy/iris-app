package schoolhandlers

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/api/v1/handlers/shared"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/response"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/service"
)

// List lấy danh sách trường học.
func (h *SchoolHandler) List(c *gin.Context) {
	adminSchoolID := shared.ExtractAdminSchoolID(c)

	var params shared.PaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid pagination params")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	schools, total, err := h.schoolService.List(ctx, adminSchoolID, params.Limit, params.Offset)
	if err != nil {
		if errors.Is(err, service.ErrSchoolAccessDenied) {
			response.Fail(c, http.StatusForbidden, "access denied")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "failed to fetch schools")
		return
	}

	response.OKPaginated(c, schools, response.Pagination{
		Total:   total,
		Limit:   params.Limit,
		Offset:  params.Offset,
		HasMore: params.Offset+len(schools) < total,
	})
}
