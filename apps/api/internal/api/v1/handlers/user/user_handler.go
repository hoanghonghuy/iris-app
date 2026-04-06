package userhandlers

import "github.com/hoanghonghuy/iris-app/apps/api/internal/service"

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
