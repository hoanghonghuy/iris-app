package parentscope

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

// UpdateMyProfile cập nhật hồ sơ cá nhân của phụ huynh (parent only - chỉ có thể cập nhật số điện thoại)
func (h *ParentScopeHandler) UpdateMyProfile(c *gin.Context) {
	var req UpdateMyProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid request body")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	userID, ok := shared.RequireCurrentUserID(c)
	if !ok {
		return
	}

	err := h.parentScopeService.UpdateMyProfile(ctx, userID, req.Phone)
	if err != nil {
		if errors.Is(err, service.ErrInvalidUserID) {
			response.Fail(c, http.StatusBadRequest, "invalid user ID")
			return
		}
		if errors.Is(err, service.ErrParentNotFound) {
			response.Fail(c, http.StatusNotFound, "parent profile not found")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "failed to update profile")
		return
	}

	response.OK(c, gin.H{
		"message": "profile updated successfully",
		"phone":   req.Phone,
	})
}

// UpdateMyProfileRequest input để phụ huynh cập nhật thông tin cá nhân (chỉ phone)
type UpdateMyProfileRequest struct {
	Phone string `json:"phone"`
}
