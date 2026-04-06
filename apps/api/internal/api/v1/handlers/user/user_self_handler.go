package userhandlers

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/response"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/service"
)

// UpdateMyPassword cập nhật mật khẩu của người dùng (self-service only)
func (h *UserHandler) UpdateMyPassword(c *gin.Context) {
	var req UpdateMyPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid request body")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	userID, ok := requireCurrentUserID(c)
	if !ok {
		return
	}

	if err := h.userService.UpdateMyPassword(ctx, userID, req.Password); err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidUserID):
			response.Fail(c, http.StatusBadRequest, "invalid user ID")
			return
		case errors.Is(err, service.ErrUserNotFound):
			response.Fail(c, http.StatusNotFound, "user not found")
			return
		case errors.Is(err, service.ErrPasswordCannotBeEmpty):
			response.Fail(c, http.StatusBadRequest, "password cannot be empty")
			return
		case errors.Is(err, service.ErrFailedToHashPassword):
			response.Fail(c, http.StatusInternalServerError, "failed to hash password")
			return
		case errors.Is(err, service.ErrFailedToUpdatePassword):
			response.Fail(c, http.StatusInternalServerError, "failed to update password")
			return
		default:
			response.Fail(c, http.StatusInternalServerError, "failed to update password")
			return
		}
	}

	response.OK(c, gin.H{"message": "password updated successfully"})
}

// Delete xóa user
func (h *UserHandler) Delete(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	userID, ok := requireCurrentUserID(c)
	if !ok {
		return
	}

	if err := h.userService.Delete(ctx, userID); err != nil {
		response.Fail(c, http.StatusInternalServerError, "failed to delete user")
		return
	}

	response.OK(c, gin.H{"message": "user deleted successfully"})
}
