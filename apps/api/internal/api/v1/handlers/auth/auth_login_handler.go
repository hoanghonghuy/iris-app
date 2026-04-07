package authhandlers

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/auth"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/response"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/service"
)

// Login xử lý đăng nhập.
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid request body")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	resp, err := h.authService.Login(ctx, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response.Fail(c, http.StatusUnauthorized, "invalid credentials")
			return
		}
		if errors.Is(err, auth.ErrInvalidCredentials) {
			response.Fail(c, http.StatusUnauthorized, "email or password incorrect")
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

// LoginWithGoogle xử lý đăng nhập bằng Google ID token.
func (h *AuthHandler) LoginWithGoogle(c *gin.Context) {
	var req GoogleLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid request body")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 8*time.Second)
	defer cancel()

	resp, err := h.authService.LoginWithGoogleToken(ctx, req.IDToken, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrGoogleLoginDisabled):
			response.Fail(c, http.StatusForbidden, "google login is disabled")
			return
		case errors.Is(err, service.ErrGoogleDomainNotAllowed):
			response.Fail(c, http.StatusForbidden, "google account domain not allowed")
			return
		case errors.Is(err, service.ErrGoogleAccountNotProvisioned):
			response.Fail(c, http.StatusUnauthorized, "google account is not provisioned")
			return
		case errors.Is(err, service.ErrGoogleLinkPasswordRequired):
			response.FailWithCode(c, http.StatusForbidden, "password confirmation required to link google account", "GOOGLE_LINK_PASSWORD_REQUIRED")
			return
		case errors.Is(err, auth.ErrInvalidCredentials):
			response.Fail(c, http.StatusUnauthorized, "invalid credentials")
			return
		case errors.Is(err, auth.ErrUserLocked):
			response.Fail(c, http.StatusForbidden, "user account locked")
			return
		default:
			response.Fail(c, http.StatusInternalServerError, "server error")
			return
		}
	}

	response.OK(c, resp)
}
