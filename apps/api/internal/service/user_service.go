package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/auth"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/repo"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	UserRepo *repo.UserRepo
	JWTAuth  *auth.Authenticator
}

// CreateUserInput chứa input để tạo user mới
type CreateUserInput struct {
	Email    string   `json:"email"`
	Password string   `json:"password"`
	Roles    []string `json:"roles"`
}

// CreateUserResponse chứa response sau khi tạo user thành công
type CreateUserResponse struct {
	UserID  uuid.UUID `json:"user_id"`
	Email   string    `json:"email"`
	Status  string    `json:"status"`
	Roles   []string  `json:"roles"`
	Message string    `json:"message"`
}

// Create tạo mới user với password được hash và assign roles
func (s *UserService) Create(ctx context.Context, input *CreateUserInput) (*CreateUserResponse, error) {
	// Validate input
	if input.Email == "" {
		return nil, errors.New("email cannot be empty")
	}
	if input.Password == "" {
		return nil, errors.New("password cannot be empty")
	}

	// Hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Create user
	userID, err := s.UserRepo.Create(ctx, input.Email, string(passwordHash))
	if err != nil {
		return nil, errors.New("failed to create user")
	}

	// Assign roles nếu có
	if len(input.Roles) > 0 {
		for _, role := range input.Roles {
			if err := s.UserRepo.AssignRole(ctx, userID, role); err != nil {
				return nil, errors.New("failed to assign role: " + role)
			}
		}
	}

	return &CreateUserResponse{
		UserID:  userID,
		Email:   input.Email,
		Status:  "active",
		Roles:   input.Roles,
		Message: "user created successfully",
	}, nil
}

// AssignRole gán role cho user
func (s *UserService) AssignRole(ctx context.Context, userID uuid.UUID, roleName string) error {
	// Validate role name
	validRoles := map[string]bool{"ADMIN": true, "TEACHER": true, "PARENT": true}
	if !validRoles[roleName] {
		return errors.New("invalid role name: " + roleName)
	}

	return s.UserRepo.AssignRole(ctx, userID, roleName)
}

// FindByEmail tìm user theo email
func (s *UserService) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	return s.UserRepo.FindByEmail(ctx, email)
}

// RolesOfUser lấy danh sách roles của user
func (s *UserService) RolesOfUser(ctx context.Context, userID uuid.UUID) ([]string, error) {
	return s.UserRepo.RolesOfUser(ctx, userID)
}

// FindByID lấy thông tin user theo ID
func (s *UserService) FindByID(ctx context.Context, userID uuid.UUID) (*model.UserInfo, error) {
	return s.UserRepo.FindByID(ctx, userID)
}

// Update cập nhật thông tin user (email và password) - hỗ trợ partial update
func (s *UserService) Update(ctx context.Context, userID uuid.UUID, email, password string) error {
	// Validate input
	if email == "" && password == "" {
		return errors.New("email or password must be provided")
	}

	// Lấy user hiện tại
	currentUser, err := s.UserRepo.FindByID(ctx, userID)
	if err != nil {
		return errors.New("failed to get current user")
	}

	// Xử lý email - nếu không có email mới, giữ nguyên email hiện tại
	if email == "" {
		email = currentUser.Email
	}

	// Xử lý password
	var passwordHash string
	if password != "" {
		// Hash password mới
		passwordHashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return errors.New("failed to hash password")
		}
		passwordHash = string(passwordHashBytes)
	} else {
		// Lấy current password hash từ user hiện tại
		user, err := s.UserRepo.FindByEmail(ctx, currentUser.Email)
		if err != nil {
			return errors.New("failed to get current password hash")
		}
		passwordHash = user.PasswordHash
	}

	return s.UserRepo.Update(ctx, userID, email, passwordHash)
}

// Delete xóa user (hard delete)
func (s *UserService) Delete(ctx context.Context, userID uuid.UUID) error {
	return s.UserRepo.Delete(ctx, userID)
}

// Lock khóa tài khoản user
func (s *UserService) Lock(ctx context.Context, userID uuid.UUID) error {
	return s.UserRepo.Lock(ctx, userID)
}

// Unlock mở khóa tài khoản user
func (s *UserService) Unlock(ctx context.Context, userID uuid.UUID) error {
	return s.UserRepo.Unlock(ctx, userID)
}
