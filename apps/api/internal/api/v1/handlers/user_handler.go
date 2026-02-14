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

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// CreateUserRequest input để admin tạo user mới (không cần password)
type CreateUserRequest struct {
	Email string   `json:"email" binding:"required,email"`
	Roles []string `json:"roles" binding:"required,min=1"`
}

// ActivateUserRequest input để user kích hoạt tài khoản
type ActivateUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// ActivateUserWithTokenRequest input để user kích hoạt tài khoản bằng token
type ActivateUserWithTokenRequest struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

// UpdateMyPasswordRequest input để user cập nhật mật khẩu (self-service)
type UpdateMyPasswordRequest struct {
	Password string `json:"password"`
}

// AssignRoleRequest input để gán role cho user
type AssignRoleRequest struct {
	RoleName string `json:"role_name" binding:"required"`
}

// CreateUser tạo user mới (admin only)
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req CreateUserRequest

	// Bind dữ liệu JSON từ request body vào struct req
	// Nếu có lỗi (thiếu field, sai định dạng,...), trả về lỗi 400
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid request body")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	resp, err := h.userService.CreateUserWithoutPassword(ctx, req.Email, req.Roles)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrFailedToAssignRole):
			response.Fail(c, http.StatusBadRequest, "failed to assign role")
			return
		default:
			response.Fail(c, http.StatusInternalServerError, "failed to create user")
			return
		}
	}

	c.Header("Location", c.Request.URL.Path+"/"+resp.ID.String())
	response.Created(c, gin.H{
		"user":    resp,
		"message": "user created successfully. User needs to activate account.",
	})
}

// ActivateUser kích hoạt tài khoản (public)
func (h *UserHandler) ActivateUser(c *gin.Context) {
	var req ActivateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid request body")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	if err := h.userService.ActivateUser(ctx, req.Email, req.Password); err != nil {
		switch {
		case errors.Is(err, service.ErrEmailCannotBeEmpty):
			response.Fail(c, http.StatusBadRequest, "email cannot be empty")
			return
		case errors.Is(err, service.ErrPasswordCannotBeEmpty):
			response.Fail(c, http.StatusBadRequest, "password cannot be empty")
			return
		case errors.Is(err, service.ErrUserNotFound):
			response.Fail(c, http.StatusNotFound, "user not found")
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

	response.OK(c, gin.H{"message": "account activated successfully"})
}

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

// GetByID lấy thông tin user theo ID (admin only - lấy từ URL param)
func (h *UserHandler) GetByID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	// Lấy userID từ URL param (admin xem thông tin user bất kỳ)
	userID, err := uuid.Parse(c.Param("userid"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid user ID")
		return
	}

	userInfo, err := h.userService.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response.Fail(c, http.StatusNotFound, "user not found")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "failed to get user")
		return
	}

	response.OK(c, userInfo)
}

// UpdateMyPassword cập nhật mật khẩu của người dùng (self-service only)
func (h *UserHandler) UpdateMyPassword(c *gin.Context) {
	var req UpdateMyPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid request body")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	// Lấy userID từ middleware AuthJWT
	claimsAny, exists := c.Get(middleware.CtxClaims)
	if !exists {
		response.Fail(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	claims := claimsAny.(*auth.Claims)

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid user ID")
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
		case errors.Is(err, service.ErrInvalidPassword):
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

	// Lấy userID từ middleware AuthJWT
	claimsAny, exists := c.Get(middleware.CtxClaims)
	if !exists {
		response.Fail(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	claims := claimsAny.(*auth.Claims)

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid user ID")
		return
	}

	if err := h.userService.Delete(ctx, userID); err != nil {
		response.Fail(c, http.StatusInternalServerError, "failed to delete user")
		return
	}

	response.OK(c, gin.H{"message": "user deleted successfully"})
}

// List lấy danh sách tất cả users (admin only)
func (h *UserHandler) List(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	users, err := h.userService.List(ctx)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "failed to fetch users")
		return
	}

	response.OK(c, users)
}

// Lock khóa tài khoản user
func (h *UserHandler) Lock(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	userID, err := uuid.Parse(c.Param("userid"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid user ID")
		return
	}

	if err := h.userService.Lock(ctx, userID); err != nil {
		response.Fail(c, http.StatusInternalServerError, "failed to lock user")
		return
	}

	response.OK(c, gin.H{"message": "user locked successfully"})
}

// Unlock mở khóa tài khoản user
func (h *UserHandler) Unlock(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	userID, err := uuid.Parse(c.Param("userid"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid user ID")
		return
	}

	if err := h.userService.Unlock(ctx, userID); err != nil {
		response.Fail(c, http.StatusInternalServerError, "failed to unlock user")
		return
	}

	response.OK(c, gin.H{"message": "user unlocked successfully"})
}

// AssignRole gán role cho user (admin only)
func (h *UserHandler) AssignRole(c *gin.Context) {
	var req AssignRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid request body")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	userIDString := c.Param("userid")
	userID, err := uuid.Parse(userIDString)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid user ID")
		return
	}

	if err := h.userService.AssignRole(ctx, userID, req.RoleName); err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidRoleName):
			response.Fail(c, http.StatusBadRequest, "invalid role name")
			return
		default:
			response.Fail(c, http.StatusInternalServerError, "failed to assign role")
			return
		}
	}

	response.OK(c, gin.H{
		"message":   "role assigned successfully",
		"user_id":   userIDString,
		"role_name": req.RoleName,
	})
}
