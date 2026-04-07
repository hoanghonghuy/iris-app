package authhandlers

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/response"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/service"
)

// ForgotPassword xử lý yêu cầu reset password.
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid request body")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	// Luôn trả success để tránh lộ thông tin email có tồn tại hay không.
	_ = h.userService.RequestPasswordReset(ctx, req.Email)

	response.OK(c, gin.H{"message": "Nếu email tồn tại, bạn sẽ nhận được link đặt lại mật khẩu."})
}

// ResetPassword xử lý đặt lại mật khẩu bằng token.
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid request body")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	if err := h.userService.ResetPasswordWithToken(ctx, req.Email, req.Token, req.Password); err != nil {
		if errors.Is(err, service.ErrResetTokenInvalid) {
			response.Fail(c, http.StatusBadRequest, "token không hợp lệ hoặc đã hết hạn")
			return
		}
		if errors.Is(err, service.ErrPasswordCannotBeEmpty) {
			response.Fail(c, http.StatusBadRequest, "mật khẩu không được để trống")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "server error")
		return
	}

	response.OK(c, gin.H{"message": "Đặt lại mật khẩu thành công. Vui lòng đăng nhập."})
}
