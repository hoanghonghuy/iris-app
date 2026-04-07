package schooladminhandlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/response"
)

// List lấy danh sách tất cả school admins (SUPER_ADMIN only).
func (h *SchoolAdminHandler) List(c *gin.Context) {
	var params PaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid pagination params")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	admins, total, err := h.schoolAdminService.List(ctx, params.Limit, params.Offset)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "failed to fetch school admins")
		return
	}

	response.OKPaginated(c, admins, response.Pagination{
		Total:   total,
		Limit:   params.Limit,
		Offset:  params.Offset,
		HasMore: params.Offset+len(admins) < total,
	})
}
