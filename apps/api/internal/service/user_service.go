package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/auth"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/repo"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo       *repo.UserRepo
	resetTokenRepo *repo.ResetTokenRepo
	jwtAuth        *auth.Authenticator
	emailSender    EmailSender
	frontendURL    string
}

func NewUserService(
	userRepo *repo.UserRepo,
	resetTokenRepo *repo.ResetTokenRepo,
	jwtAuth *auth.Authenticator,
	emailSender EmailSender,
	frontendURL string,
) *UserService {
	return &UserService{
		userRepo:       userRepo,
		resetTokenRepo: resetTokenRepo,
		jwtAuth:        jwtAuth,
		emailSender:    emailSender,
		frontendURL:    frontendURL,
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
	// TODO: chặn SUPER_ADMIN ở flow thường; chỉ cho phép qua quy trình promote có kiểm soát + audit.
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
	err = s.userRepo.Update(ctx, user.UserID, user.Email, string(passwordHash))
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
	err = s.userRepo.ActivateWithPassword(ctx, user.UserID, string(passwordHash))
	if err != nil {
		return ErrFailedToActivateUser
	}

	return nil
}

// AssignRole gán role cho user
func (s *UserService) AssignRole(ctx context.Context, userID uuid.UUID, roleName string) error {
	// Validate role name
	// TODO: Tách endpoint/flow riêng cho SUPER_ADMIN (2-person approval), không dùng AssignRole chung.
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
func (s *UserService) List(ctx context.Context, adminSchoolID *uuid.UUID, roleFilter string, limit, offset int) ([]model.UserInfo, int, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	validRoles := map[string]bool{"SUPER_ADMIN": true, "SCHOOL_ADMIN": true, "TEACHER": true, "PARENT": true}
	if roleFilter != "" && !validRoles[roleFilter] {
		return nil, 0, fmt.Errorf("%w: %s", ErrInvalidRoleName, roleFilter)
	}

	return s.userRepo.List(ctx, adminSchoolID, roleFilter, limit, offset)
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
		return ErrPasswordCannotBeEmpty
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

const resetTokenExpiry = 15 * time.Minute

// RequestPasswordReset tạo reset token và gửi email cho user.
// Luôn trả nil để tránh lộ thông tin email có tồn tại hay không.
func (s *UserService) RequestPasswordReset(ctx context.Context, email string) error {
	if email == "" {
		return ErrEmailCannotBeEmpty
	}

	// Tìm user (bỏ qua nếu không tìm thấy)
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil // không leak thông tin
	}
	_ = user // user.UserID được dùng bên dưới

	// Tạo crypto-random token (32 bytes = 256 bits)
	rawToken := make([]byte, 32)
	if _, err := rand.Read(rawToken); err != nil {
		return nil
	}
	plainToken := hex.EncodeToString(rawToken) // 64 hex chars

	// SHA-256 hash trước khi lưu
	hash := sha256.Sum256([]byte(plainToken))
	tokenHash := hex.EncodeToString(hash[:])

	// Lưu vào DB
	expiresAt := time.Now().Add(resetTokenExpiry)
	if err := s.resetTokenRepo.Create(ctx, user.UserID, tokenHash, expiresAt); err != nil {
		return nil
	}

	resetURL, err := url.Parse(s.frontendURL)
	if err != nil {
		log.Printf("[Fatal Error] Unable to parse frontendURL: %v", err)
		return nil
	}

	resetURL.Path = "/reset-password"

	htmlBody := fmt.Sprintf(`
<h2>Đặt lại mật khẩu — Iris</h2>
<p>Bạn đã yêu cầu đặt lại mật khẩu. Sử dụng mã bên dưới (hết hạn sau 15 phút):</p>
<p><strong>%s</strong></p>

<p><a href="%s" style="display: inline-block; background: #007bff; color: white; padding: 10px 20px; border-radius: 6px; text-decoration: none;">Đặt lại mật khẩu</a></p>

<p style="color: #666; font-size: 12px;">Sao chép mã trên và nhập vào form đặt lại mật khẩu cùng với email của bạn</p>
<p>Nếu bạn không yêu cầu, hãy bỏ qua email này để bảo vệ tài khoản.</p>
`, plainToken, resetURL.String())

	// _ = s.emailSender.Send(email, "Đặt lại mật khẩu — Iris", htmlBody)
	// gửi email ngầm bằng goroutine để không block API
	// Mục đích: tránh để API endpoint phải chờ kết nối SMTP (3s),
	// giúp phản hồi HTTP về FE lập tức (giảm rủi ro nghẽn Connection Pool khi có nhiều người reset password).
	go func(targetEmail, subject, body string) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[Recover] Email sending panicked: %v", r)
			}
		}()
		if err := s.emailSender.Send(targetEmail, subject, body); err != nil {
			log.Printf("[EmailService] Failed to send reset password email to %s: %v", targetEmail, err)
		}
	}(email, "Đặt lại mật khẩu — Iris", htmlBody)

	return nil
}

// ResetPasswordWithToken xác thực token và đặt mật khẩu mới
func (s *UserService) ResetPasswordWithToken(ctx context.Context, email, plainToken, newPassword string) error {
	if email == "" {
		return ErrResetTokenInvalid
	}
	if plainToken == "" {
		return ErrResetTokenInvalid
	}
	if newPassword == "" {
		return ErrPasswordCannotBeEmpty
	}

	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return ErrResetTokenInvalid
	}

	// Hash token đầu vào để so khớp với DB
	hash := sha256.Sum256([]byte(plainToken))
	tokenHash := hex.EncodeToString(hash[:])

	// Tìm token
	rt, err := s.resetTokenRepo.FindByTokenHash(ctx, tokenHash)
	if err != nil {
		return ErrResetTokenInvalid
	}

	// Token phải thuộc đúng email gửi kèm.
	if rt.UserID != user.UserID {
		return ErrResetTokenInvalid
	}

	// Kiểm tra hết hạn
	if rt.ExpiresAt.Before(time.Now()) {
		return ErrResetTokenInvalid
	}

	// Hash mật khẩu mới
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return ErrFailedToHashPassword
	}

	// Đánh dấu token đã sử dụng TRƯỚC khi đặt mật khẩu (fail-fast chống TOCTOU race)
	// Nếu 2 request đồng thời: request thứ 2 sẽ thất bại ở bước này
	// repo.ErrNoRowsUpdated cho biết token đã bị mark bởi request khác
	if err := s.resetTokenRepo.MarkUsed(ctx, rt.ID); err != nil {
		if errors.Is(err, repo.ErrNoRowsUpdated) {
			return ErrResetTokenInvalid // token đã bị dùng bởi request khác
		}
		return ErrResetTokenInvalid // lỗi DB khác, treat như invalid
	}

	// Cập nhật mật khẩu
	if err := s.userRepo.UpdatePassword(ctx, rt.UserID, user.Email, string(passwordHash)); err != nil {
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
