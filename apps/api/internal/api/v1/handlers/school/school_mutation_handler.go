package schoolhandlers

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

// Create tạo mới trường học.
func (h *SchoolHandler) Create(c *gin.Context) {
	var req CreateSchoolRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid request body")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	school, err := h.schoolService.Create(ctx, req.Name, req.Address)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "failed to create school")
		return
	}

	response.Created(c, gin.H{"school_id": school.SchoolID})
}

// Update cập nhật trường học (SUPER_ADMIN only).
func (h *SchoolHandler) Update(c *gin.Context) {
	schoolID, err := uuid.Parse(c.Param("school_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid school_id format")
		return
	}

	var req UpdateSchoolRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid request body")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	if err := h.schoolService.Update(ctx, schoolID, req.Name, req.Address); err != nil {
		if errors.Is(err, service.ErrSchoolNotFound) {
			response.Fail(c, http.StatusNotFound, "school not found")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "failed to update school")
		return
	}

	response.OK(c, gin.H{"message": "school updated successfully", "school_id": schoolID.String()})
}

// Delete xóa trường học (SUPER_ADMIN only).
func (h *SchoolHandler) Delete(c *gin.Context) {
	schoolID, err := uuid.Parse(c.Param("school_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid school_id format")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	if err := h.schoolService.Delete(ctx, schoolID); err != nil {
		if errors.Is(err, service.ErrSchoolNotFound) {
			response.Fail(c, http.StatusNotFound, "school not found")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "failed to delete school")
		return
	}

	response.OK(c, gin.H{"message": "school deleted successfully", "school_id": schoolID.String()})
}
