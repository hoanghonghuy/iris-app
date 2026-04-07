package userhandlers

import (
	"context"

	"github.com/google/uuid"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/service"
)

type userService interface {
	CreateUserWithoutPassword(ctx context.Context, adminSchoolID *uuid.UUID, email string, roles []string) (*model.UserInfo, error)
	FindByID(ctx context.Context, adminSchoolID *uuid.UUID, userID uuid.UUID) (*model.UserInfo, error)
	List(ctx context.Context, adminSchoolID *uuid.UUID, roleFilter string, limit, offset int) ([]model.UserInfo, int, error)
	Lock(ctx context.Context, adminSchoolID *uuid.UUID, userID uuid.UUID) error
	Unlock(ctx context.Context, adminSchoolID *uuid.UUID, userID uuid.UUID) error
	AssignRole(ctx context.Context, userID uuid.UUID, roleName string) error
	ActivateUserWithToken(ctx context.Context, token, password string) error
	UpdateMyPassword(ctx context.Context, userID uuid.UUID, password string) error
	Delete(ctx context.Context, userID uuid.UUID) error
}

type UserHandler struct {
	userService userService
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

// ActivateUserWithTokenRequest input để user kích hoạt tài khoản bằng token
type ActivateUserWithTokenRequest struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

// UpdateMyPasswordRequest input để user cập nhật mật khẩu (self-service)
type UpdateMyPasswordRequest struct {
	Password string `json:"password" binding:"required,min=6"`
}

// AssignRoleRequest input để gán role cho user
type AssignRoleRequest struct {
	RoleName string `json:"role_name" binding:"required"`
}
