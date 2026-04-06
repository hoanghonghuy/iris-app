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

// ActivateUserWithToken kích hoạt tài khoản bằng token (public)
func (h *UserHandler) ActivateUserWithToken(c *gin.Context) {
	var req ActivateUserWithTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid request body")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	if err := h.userService.ActivateUserWithToken(ctx, req.Token, req.Password); err != nil {
		switch {
		case errors.Is(err, service.ErrActivationTokenRequired):
			response.Fail(c, http.StatusBadRequest, "activation token is required")
			return
		case errors.Is(err, service.ErrInvalidActivationToken):
			response.Fail(c, http.StatusBadRequest, "invalid activation token")
			return
		case errors.Is(err, service.ErrActivationTokenExpired):
			response.Fail(c, http.StatusBadRequest, "activation token has expired")
			return
		case errors.Is(err, service.ErrPasswordCannotBeEmpty):
			response.Fail(c, http.StatusBadRequest, "password cannot be empty")
			return
		case errors.Is(err, service.ErrFailedToHashPassword):
			response.Fail(c, http.StatusInternalServerError, "failed to hash password")
			return
		case errors.Is(err, service.ErrFailedToActivateUser):
			response.Fail(c, http.StatusInternalServerError, "failed to activate user")
			return
		default:
			response.Fail(c, http.StatusInternalServerError, "failed to activate user")
			return
		}
	}

	response.OK(c, gin.H{
		"message": "account activated successfully",
		"status":  "active",
	})
}
