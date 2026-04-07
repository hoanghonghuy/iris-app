package authhandlers

import (
	"context"

	"github.com/google/uuid"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/service"
)

type authService interface {
	Login(ctx context.Context, email, password string) (*service.LoginResponse, error)
	LoginWithGoogleToken(ctx context.Context, googleIDToken, password string) (*service.LoginResponse, error)
}

type userService interface {
	RequestPasswordReset(ctx context.Context, email string) error
	ResetPasswordWithToken(ctx context.Context, email, plainToken, newPassword string) error
	FindByID(ctx context.Context, adminSchoolID *uuid.UUID, userID uuid.UUID) (*model.UserInfo, error)
}

type AuthHandler struct {
	authService authService
	userService userService
}

func NewAuthHandler(authService authService, userService userService) *AuthHandler {
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

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}
