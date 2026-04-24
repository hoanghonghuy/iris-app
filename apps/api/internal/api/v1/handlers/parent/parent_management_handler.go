package parenthandlers

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

// Update cập nhật thông tin phụ huynh
func (h *ParentHandler) Update(c *gin.Context) {
	// lấy id trường học từ token
	adminSchoolID := shared.ExtractAdminSchoolID(c)

	parentID, err := uuid.Parse(c.Param("parent_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid parent_id format")
		return
	}

	var req UpdateParentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid request body")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	err = h.parentService.Update(ctx, adminSchoolID, parentID, req.FullName, req.Phone, req.SchoolID)
	if err != nil {
		if errors.Is(err, service.ErrSchoolAccessDenied) {
			response.Fail(c, http.StatusForbidden, "access denied")
			return
		}
		if errors.Is(err, service.ErrParentNotFound) {
			response.Fail(c, http.StatusNotFound, "parent not found")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "failed to update parent")
		return
	}

	response.OK(c, gin.H{
		"message":   "parent updated successfully",
		"parent_id": parentID.String(),
	})
}
