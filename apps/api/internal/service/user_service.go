package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/auth"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/repo"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo *repo.UserRepo
	jwtAuth  *auth.Authenticator
}

func NewUserService(userRepo *repo.UserRepo, jwtAuth *auth.Authenticator) *UserService {
	return &UserService{
		userRepo: userRepo,
		jwtAuth:  jwtAuth,
	}
}

// CreateUserWithoutPassword tạo user mới không cần password (admin only).
// adminSchoolID != nil → SCHOOL_ADMIN: chỉ được tạo user với role TEACHER hoặc PARENT
func (s *UserService) CreateUserWithoutPassword(ctx context.Context, adminSchoolID *uuid.UUID, email string, roles []string) (*model.UserInfo, error) {
	// Validate input
	if email == "" {
		return nil, ErrEmailCannotBeEmpty
	}
	if len(roles) == 0 {
		return nil, ErrRolesCannotBeEmpty
	}

	// Validate tên role hợp lệ trước khi tạo user (tránh silent failure khi role không tồn tại trong DB)
	validRoles := map[string]bool{"SUPER_ADMIN": true, "SCHOOL_ADMIN": true, "TEACHER": true, "PARENT": true}
	for _, role := range roles {
		if !validRoles[role] {
			return nil, fmt.Errorf("%w: %s", ErrInvalidRoleName, role)
		}
	}

	// SCHOOL_ADMIN: chỉ được tạo user với role TEACHER hoặc PARENT (không thể tự nâng quyền)
	if adminSchoolID != nil {
		allowedForSchoolAdmin := map[string]bool{"TEACHER": true, "PARENT": true}
		for _, role := range roles {
			if !allowedForSchoolAdmin[role] {
				return nil, fmt.Errorf("%w: %s", ErrCannotAssignRole, role)
			}
		}
	}

	// Generate temporary password hash (user sẽ thay đổi khi activate)
	tempPassword := uuid.New().String()
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(tempPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, ErrFailedToGenerateTempPassword
	}

	// Check if user có role TEACHER → tạo với pending status
	hasTeacherRole := false
	for _, role := range roles {
		if role == "TEACHER" {
			hasTeacherRole = true
			break
		}
	}

	// Create user với status phù hợp
	var userID uuid.UUID
	if hasTeacherRole {
		userID, err = s.userRepo.CreatePending(ctx, email, string(passwordHash))
	} else {
		// Admin và Parent tạo trực tiếp với status active
		userID, err = s.userRepo.CreateActive(ctx, email, string(passwordHash))
	}
	if err != nil {
		return nil, ErrFailedToCreateUser
	}

	// Assign roles
	for _, role := range roles {
		if err := s.userRepo.AssignRole(ctx, userID, role); err != nil {
			return nil, fmt.Errorf("%w: %s", ErrFailedToAssignRole, role)
		}
	}

	return s.FindByID(ctx, nil, userID)
}

// ActivateUser kích hoạt tài khoản user (set password)
func (s *UserService) ActivateUser(ctx context.Context, email, password string) error {
	// Validate input
	if email == "" {
		return ErrEmailCannotBeEmpty
	}
	if password == "" {
		return ErrPasswordCannotBeEmpty
	}

	// Check user exists
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return ErrUserNotFound
	}

	// Hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return ErrFailedToHashPassword
	}

	// Update password and status
	err = s.userRepo.Update(ctx, user.ID, user.Email, string(passwordHash))
	if err != nil {
		return ErrFailedToActivateUser
	}

	return nil
}

// ActivateUserWithToken kích hoạt tài khoản bằng token (cho teacher activation flow)
func (s *UserService) ActivateUserWithToken(ctx context.Context, token, password string) error {
	// Validate input
	if token == "" {
		return ErrActivationTokenRequired
	}
	if password == "" {
		return ErrPasswordCannotBeEmpty
	}

	// Find user by activation token
	user, err := s.userRepo.FindByActivationToken(ctx, token)
	if err != nil {
		return ErrInvalidActivationToken
	}

	// Check token không hết hạn
	if user.TokenExpiresAt.Before(time.Now()) {
		return ErrActivationTokenExpired
	}

	// Hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return ErrFailedToHashPassword
	}

	// Activate user + update password
	err = s.userRepo.ActivateWithPassword(ctx, user.ID, string(passwordHash))
	if err != nil {
		return ErrFailedToActivateUser
	}

	return nil
}

// AssignRole gán role cho user
func (s *UserService) AssignRole(ctx context.Context, userID uuid.UUID, roleName string) error {
	// Validate role name
	validRoles := map[string]bool{
		"SUPER_ADMIN":  true,
		"SCHOOL_ADMIN": true,
		"TEACHER":      true,
		"PARENT":       true,
	}
	if !validRoles[roleName] {
		return fmt.Errorf("%w: %s", ErrInvalidRoleName, roleName)
	}

	return s.userRepo.AssignRole(ctx, userID, roleName)
}

// FindByEmail tìm user theo email
func (s *UserService) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	return s.userRepo.FindByEmail(ctx, email)
}

// RolesOfUser lấy danh sách roles của user
func (s *UserService) RolesOfUser(ctx context.Context, userID uuid.UUID) ([]string, error) {
	return s.userRepo.RolesOfUser(ctx, userID)
}

// FindByID lấy thông tin user theo ID.
func (s *UserService) FindByID(ctx context.Context, adminSchoolID *uuid.UUID, userID uuid.UUID) (*model.UserInfo, error) {
	info, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// SCHOOL_ADMIN: validate user thuộc cùng school với admin (qua teachers/parents)
	// adminSchoolID == nil => SUPER_ADMIN: không cần validate
	if adminSchoolID != nil {
		belongs, err := s.userRepo.IsUserInSchool(ctx, userID, *adminSchoolID)
		if err != nil {
			return nil, err
		}
		if !belongs {
			return nil, ErrSchoolAccessDenied
		}
	}

	return info, nil
}

// List lấy danh sách users.
func (s *UserService) List(ctx context.Context, adminSchoolID *uuid.UUID, limit, offset int) ([]model.UserInfo, int, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}
	return s.userRepo.List(ctx, adminSchoolID, limit, offset)
}

// UpdateEmail cập nhật email của user (admin only)
func (s *UserService) UpdateEmail(ctx context.Context, userID uuid.UUID, email string) error {
	// Validate input
	if email == "" {
		return ErrEmailCannotBeEmpty
	}

	// check user exists
	_, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return ErrUserNotFound
	}

	return s.userRepo.UpdateEmail(ctx, userID, email)
}

// UpdateMyPassword cập nhật mật khẩu của người dùng (user)
func (s *UserService) UpdateMyPassword(ctx context.Context, userID uuid.UUID, password string) error {
	// Validate userID
	if userID == uuid.Nil {
		return ErrInvalidUserID
	}

	// validate password
	if password == "" {
		return ErrInvalidPassword
	}

	// check if user exists
	userInfo, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return ErrUserNotFound
	}

	// Hash password
	passwordHashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return ErrFailedToHashPassword
	}

	err = s.userRepo.UpdatePassword(ctx, userID, userInfo.Email, string(passwordHashBytes))
	if err != nil {
		return ErrFailedToUpdatePassword
	}

	return nil
}

// Delete xóa user (hard delete)
func (s *UserService) Delete(ctx context.Context, userID uuid.UUID) error {
	return s.userRepo.Delete(ctx, userID)
}

// Lock khóa tài khoản user.
func (s *UserService) Lock(ctx context.Context, adminSchoolID *uuid.UUID, userID uuid.UUID) error {
	// SCHOOL_ADMIN: validate user thuộc cùng school với admin
	// adminSchoolID == nil => SUPER_ADMIN: không cần validate
	if adminSchoolID != nil {
		belongs, err := s.userRepo.IsUserInSchool(ctx, userID, *adminSchoolID)
		if err != nil {
			return err
		}
		if !belongs {
			return ErrSchoolAccessDenied
		}
	}

	return s.userRepo.Lock(ctx, userID)
}

// Unlock mở khóa tài khoản user.
func (s *UserService) Unlock(ctx context.Context, adminSchoolID *uuid.UUID, userID uuid.UUID) error {
	// SCHOOL_ADMIN: validate user thuộc cùng school với admin
	// adminSchoolID == nil => SUPER_ADMIN: không cần validate
	if adminSchoolID != nil {
		belongs, err := s.userRepo.IsUserInSchool(ctx, userID, *adminSchoolID)
		if err != nil {
			return err
		}
		if !belongs {
			return ErrSchoolAccessDenied
		}
	}

	return s.userRepo.Unlock(ctx, userID)
}
