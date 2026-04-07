package parenthandlers

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/api/v1/handlers/shared"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/response"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/service"
)

// List lấy danh sách phụ huynh (admin only)
func (h *ParentHandler) List(c *gin.Context) {
	adminSchoolID := shared.ExtractAdminSchoolID(c)

	var params shared.PaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid pagination params")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	parents, total, err := h.parentService.List(ctx, adminSchoolID, params.Limit, params.Offset)
	if err != nil {
		if errors.Is(err, service.ErrSchoolAccessDenied) {
			response.Fail(c, http.StatusForbidden, "access denied")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "failed to fetch parents")
		return
	}

	response.OKPaginated(c, parents, response.Pagination{
		Total:   total,
		Limit:   params.Limit,
		Offset:  params.Offset,
		HasMore: params.Offset+len(parents) < total,
	})
}

// GetByID lấy thông tin phụ huynh theo parent_id (admin only)
func (h *ParentHandler) GetByID(c *gin.Context) {
	adminSchoolID := shared.ExtractAdminSchoolID(c)

	parentID, err := uuid.Parse(c.Param("parent_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid parent_id format")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	parent, err := h.parentService.GetByParentID(ctx, adminSchoolID, parentID)
	if err != nil {
		if errors.Is(err, service.ErrSchoolAccessDenied) {
			response.Fail(c, http.StatusForbidden, "access denied")
			return
		}
		if errors.Is(err, pgx.ErrNoRows) {
			response.Fail(c, http.StatusNotFound, "parent not found")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "failed to fetch parent")
		return
	}

	response.OK(c, parent)
}
