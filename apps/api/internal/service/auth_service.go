package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/auth"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/repo"
)

type AuthService struct {
	userRepo        *repo.UserRepo
	schoolAdminRepo *repo.SchoolAdminRepo
	jwtAuth         *auth.Authenticator
}

func NewAuthService(userRepo *repo.UserRepo, schoolAdminRepo *repo.SchoolAdminRepo, jwtAuth *auth.Authenticator) *AuthService {
	return &AuthService{
		userRepo:        userRepo,
		schoolAdminRepo: schoolAdminRepo,
		jwtAuth:         jwtAuth,
	}
}

// LoginResponse chứa thông tin trả về sau khi đăng nhập thành công
type LoginResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// Login xử lý logic đăng nhập
func (s *AuthService) Login(ctx context.Context, email, password string) (*LoginResponse, error) {
	// 1. Tìm user theo email
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err // pgx.ErrNoRows sẽ được handler xử lý
	}

	// 2. Kiểm tra trạng thái tài khoản
	if user.Status == "locked" {
		return nil, auth.ErrUserLocked
	}

	// 3. Verify password
	if !auth.VerifyPassword(user.PasswordHash, password) {
		return nil, auth.ErrInvalidCredentials
	}

	// 4. Lấy roles của user
	roles, err := s.userRepo.RolesOfUser(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	// 5. Nếu user có role SCHOOL_ADMIN → lấy school_id từ bảng school_admins
	//    để nhúng vào JWT claims (không cần query DB mỗi request)
	var schoolID string
	for _, r := range roles {
		if r == "SCHOOL_ADMIN" {
			admin, err := s.schoolAdminRepo.GetByUserID(ctx, user.ID)
			if err == nil {
				schoolID = admin.SchoolID.String()
			}
			break
		}
	}

	// 6. Tạo JWT token (schoolID rỗng cho SUPER_ADMIN/TEACHER/PARENT)
	token, err := s.jwtAuth.SignToken(user.ID.String(), user.Email, roles, schoolID)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   s.jwtAuth.TTLSeconds,
	}, nil
}

// GetUserInfo trả về thông tin user theo ID, hiện tại chưa dùng.
func (s *AuthService) GetUserInfo(ctx context.Context, userID uuid.UUID) (*model.UserInfo, error) {
	return s.userRepo.FindByID(ctx, userID)
}
