package classhandlers

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

// Create tạo mới lớp học
func (h *ClassHandler) Create(c *gin.Context) {
	adminSchoolID := shared.ExtractAdminSchoolID(c)

	var req CreateClassRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid request body")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	class, err := h.classService.Create(ctx, adminSchoolID, req.SchoolID, req.Name, req.SchoolYear)
	if err != nil {
		if errors.Is(err, service.ErrSchoolAccessDenied) {
			response.Fail(c, http.StatusForbidden, "access denied")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "failed to create class")
		return
	}

	response.Created(c, gin.H{"class_id": class.ClassID})
}

// Update cập nhật lớp học
func (h *ClassHandler) Update(c *gin.Context) {
	adminSchoolID := shared.ExtractAdminSchoolID(c)

	classID, err := uuid.Parse(c.Param("class_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid class_id format")
		return
	}

	var req UpdateClassRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid request body")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	if err := h.classService.Update(ctx, adminSchoolID, classID, req.Name, req.SchoolYear); err != nil {
		if errors.Is(err, service.ErrSchoolAccessDenied) {
			response.Fail(c, http.StatusForbidden, "access denied")
			return
		}
		if errors.Is(err, service.ErrClassNotFound) {
			response.Fail(c, http.StatusNotFound, "class not found")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "failed to update class")
		return
	}

	response.OK(c, gin.H{"message": "class updated successfully", "class_id": classID.String()})
}

// Delete xóa lớp học
func (h *ClassHandler) Delete(c *gin.Context) {
	adminSchoolID := shared.ExtractAdminSchoolID(c)

	classID, err := uuid.Parse(c.Param("class_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid class_id format")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	if err := h.classService.Delete(ctx, adminSchoolID, classID); err != nil {
		if errors.Is(err, service.ErrSchoolAccessDenied) {
			response.Fail(c, http.StatusForbidden, "access denied")
			return
		}
		if errors.Is(err, service.ErrClassNotFound) {
			response.Fail(c, http.StatusNotFound, "class not found")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "failed to delete class")
		return
	}

	response.OK(c, gin.H{"message": "class deleted successfully", "class_id": classID.String()})
}
