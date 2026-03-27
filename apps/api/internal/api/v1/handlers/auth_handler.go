package handlers

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/auth"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/middleware"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/response"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
	userService *service.UserService
}

func NewAuthHandler(authService *service.AuthService, userService *service.UserService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		userService: userService,
	}
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type GoogleLoginRequest struct {
	IDToken  string `json:"id_token" binding:"required"`
	Password string `json:"password"`
}

// Login xử lý đăng nhập
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
			response.Fail(c, http.StatusNotFound, "google login is disabled")
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

// Me trả về thông tin user đã đăng nhập
func (h *AuthHandler) Me(c *gin.Context) {
	// Lấy claims từ context (được middleware AuthJWT set)
	claimsAny, exists := c.Get(middleware.CtxClaims)
	if !exists {
		response.Fail(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	claims := claimsAny.(*auth.Claims)

	// Trả về thông tin từ JWT claims (đã validate và xác thực)
	result := gin.H{
		"user_id": claims.UserID,
		"email":   claims.Email,
		"roles":   claims.Roles,
	}

	// Lấy thêm full_name từ Database do JWT không chứa
	if uid, err := uuid.Parse(claims.UserID); err == nil {
		if userInfo, err := h.userService.FindByID(c.Request.Context(), nil, uid); err == nil {
			result["full_name"] = userInfo.FullName
		}
	}

	// Nếu user là SCHOOL_ADMIN → trả thêm school_id
	if claims.SchoolID != "" {
		result["school_id"] = claims.SchoolID
	}

	response.OK(c, result)
}

// ForgotPassword xử lý yêu cầu reset password
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid request body")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	// Luôn trả success để tránh lộ thông tin email có tồn tại hay không
	_ = h.userService.RequestPasswordReset(ctx, req.Email)

	response.OK(c, gin.H{"message": "Nếu email tồn tại, bạn sẽ nhận được link đặt lại mật khẩu."})
}

// ResetPassword xử lý đặt lại mật khẩu bằng token
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Token    string `json:"token" binding:"required"`
		Password string `json:"password" binding:"required,min=6"`
	}
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
