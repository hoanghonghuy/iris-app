package schooladminhandlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/response"
)

// Create tạo mới school admin (SUPER_ADMIN only).
func (h *SchoolAdminHandler) Create(c *gin.Context) {
	var req CreateSchoolAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid request body")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	admin, err := h.schoolAdminService.Create(ctx, req.UserID, req.SchoolID, req.FullName, req.Phone)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "failed to create school admin")
		return
	}

	c.Header("Location", c.Request.URL.Path+"/"+admin.AdminID.String())
	response.Created(c, admin)
}

// Delete xóa school admin theo admin_id (SUPER_ADMIN only).
// Chỉ xóa record trong school_admins, không xóa user.
func (h *SchoolAdminHandler) Delete(c *gin.Context) {
	adminID, err := uuid.Parse(c.Param("admin_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid admin_id format")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	if err := h.schoolAdminService.Delete(ctx, adminID); err != nil {
		response.Fail(c, http.StatusInternalServerError, "failed to delete school admin")
		return
	}

	response.OK(c, gin.H{
		"message":  "school admin deleted successfully",
		"admin_id": adminID.String(),
	})
}
