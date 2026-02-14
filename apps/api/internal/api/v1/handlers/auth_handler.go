package handlers

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/auth"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/middleware"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/response"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
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
	response.OK(c, gin.H{
		"user_id": claims.UserID,
		"email":   claims.Email,
		"roles":   claims.Roles,
	})
}
