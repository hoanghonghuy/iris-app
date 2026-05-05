package authhandlers

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/auth"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/response"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/service"
)

// Refresh đổi refresh token lấy cặp token mới.
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid request body")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	resp, err := h.authService.Refresh(ctx, req.RefreshToken)
	if err != nil {
		if errors.Is(err, service.ErrRefreshTokenInvalid) {
			response.Fail(c, http.StatusUnauthorized, "invalid refresh token")
			return
		}
		if errors.Is(err, auth.ErrUserLocked) {
			response.Fail(c, http.StatusForbidden, "user account locked")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "server error")
		return
	}

	response.OK(c, resp)
}
