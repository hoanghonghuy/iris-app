package userhandlers

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/api/v1/handlers/shared"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/response"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/service"
)

// CreateUser tạo user mới (admin only)
func (h *UserHandler) CreateUser(c *gin.Context) {
	adminSchoolID := shared.ExtractAdminSchoolID(c)

	var req CreateUserRequest

	// Bind dữ liệu JSON từ request body vào struct req
	// Nếu có lỗi (thiếu field, sai định dạng,...), trả về lỗi 400
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid request body")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	resp, err := h.userService.CreateUserWithoutPassword(ctx, adminSchoolID, req.Email, req.Roles)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUserAlreadyHasRole):
			response.Fail(c, http.StatusConflict, "user already has a role, cannot assign another role")
			return
		case errors.Is(err, service.ErrCannotAssignRole):
			response.Fail(c, http.StatusForbidden, "insufficient permissions to assign this role")
			return
		case errors.Is(err, service.ErrCannotAssignRoleSuperAdmin):
			response.Fail(c, http.StatusForbidden, "SUPER_ADMIN role requires dedicated promote flow")
			return
		case errors.Is(err, service.ErrFailedToAssignRole):
			response.Fail(c, http.StatusBadRequest, "failed to assign role")
			return
		default:
			response.Fail(c, http.StatusInternalServerError, "failed to create user")
			return
		}
	}

	c.Header("Location", c.Request.URL.Path+"/"+resp.UserID.String())
	response.Created(c, gin.H{
		"user":    resp,
		"message": "user created successfully. User needs to activate account.",
	})
}

// GetByID lấy thông tin user theo ID (admin only - lấy từ URL param)
func (h *UserHandler) GetByID(c *gin.Context) {
	adminSchoolID := shared.ExtractAdminSchoolID(c)

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	// Lấy userID từ URL param (admin xem thông tin user bất kỳ)
	userID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid user ID")
		return
	}

	userInfo, err := h.userService.FindByID(ctx, adminSchoolID, userID)
	if err != nil {
		if errors.Is(err, service.ErrSchoolAccessDenied) {
			response.Fail(c, http.StatusForbidden, "access denied")
			return
		}
		if errors.Is(err, pgx.ErrNoRows) {
			response.Fail(c, http.StatusNotFound, "user not found")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "failed to get user")
		return
	}

	response.OK(c, userInfo)
}

// List lấy danh sách users (admin only)
func (h *UserHandler) List(c *gin.Context) {
	adminSchoolID := shared.ExtractAdminSchoolID(c)

	var params shared.PaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid pagination params")
		return
	}
	roleFilter := c.Query("role")

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	users, total, err := h.userService.List(ctx, adminSchoolID, roleFilter, params.Limit, params.Offset)
	if err != nil {
		if errors.Is(err, service.ErrSchoolAccessDenied) {
			response.Fail(c, http.StatusForbidden, "access denied")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "failed to fetch users")
		return
	}

	response.OKPaginated(c, users, response.Pagination{
		Total:   total,
		Limit:   params.Limit,
		Offset:  params.Offset,
		HasMore: params.Offset+len(users) < total,
	})
}

// Lock khóa tài khoản user
func (h *UserHandler) Lock(c *gin.Context) {
	adminSchoolID := shared.ExtractAdminSchoolID(c)

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	userID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid user ID")
		return
	}

	if err := h.userService.Lock(ctx, adminSchoolID, userID); err != nil {
		if errors.Is(err, service.ErrSchoolAccessDenied) {
			response.Fail(c, http.StatusForbidden, "access denied")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "failed to lock user")
		return
	}

	response.OK(c, gin.H{"message": "user locked successfully"})
}

// Unlock mở khóa tài khoản user
func (h *UserHandler) Unlock(c *gin.Context) {
	adminSchoolID := shared.ExtractAdminSchoolID(c)

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	userID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid user ID")
		return
	}

	if err := h.userService.Unlock(ctx, adminSchoolID, userID); err != nil {
		if errors.Is(err, service.ErrSchoolAccessDenied) {
			response.Fail(c, http.StatusForbidden, "access denied")
			return
		}
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

	userIDString := c.Param("user_id")
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
		case errors.Is(err, service.ErrCannotAssignRoleSuperAdmin):
			response.Fail(c, http.StatusForbidden, "SUPER_ADMIN role requires dedicated promote flow")
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
