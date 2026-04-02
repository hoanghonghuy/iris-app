package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/response"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/service"
)

type SchoolAdminHandler struct {
	schoolAdminService *service.SchoolAdminService
}

func NewSchoolAdminHandler(schoolAdminService *service.SchoolAdminService) *SchoolAdminHandler {
	return &SchoolAdminHandler{
		schoolAdminService: schoolAdminService,
	}
}

// CreateSchoolAdminRequest input để SUPER_ADMIN tạo school admin mới
type CreateSchoolAdminRequest struct {
	UserID   uuid.UUID `json:"user_id" binding:"required"`
	SchoolID uuid.UUID `json:"school_id" binding:"required"`
	FullName string    `json:"full_name"`
	Phone    string    `json:"phone"`
}

// Create tạo mới school admin (SUPER_ADMIN only)
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

// List lấy danh sách tất cả school admins (SUPER_ADMIN only)
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

// Delete xóa school admin theo admin_id (SUPER_ADMIN only)
// Chỉ xóa record trong school_admins, không xóa user
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
