package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/auth"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/repo"
	"github.com/jackc/pgx/v5"
)

var (
	ErrGoogleLoginDisabled         = errors.New("google login disabled")
	ErrGoogleLinkPasswordRequired  = errors.New("password is required to link google account")
	ErrGoogleAccountNotProvisioned = errors.New("google account is not provisioned")
	ErrGoogleDomainNotAllowed      = errors.New("google hosted domain is not allowed")
)

type AuthService struct {
	userRepo        *repo.UserRepo
	schoolAdminRepo *repo.SchoolAdminRepo
	jwtAuth         *auth.Authenticator
	googleVerifier  auth.GoogleTokenVerifier
	googleEnabled   bool
	googleHD        string
}

func NewAuthService(
	userRepo *repo.UserRepo,
	schoolAdminRepo *repo.SchoolAdminRepo,
	jwtAuth *auth.Authenticator,
	googleVerifier auth.GoogleTokenVerifier,
	googleEnabled bool,
	googleHD string,
) *AuthService {
	return &AuthService{
		userRepo:        userRepo,
		schoolAdminRepo: schoolAdminRepo,
		jwtAuth:         jwtAuth,
		googleVerifier:  googleVerifier,
		googleEnabled:   googleEnabled,
		googleHD:        googleHD,
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

	// 4. Build response + JWT theo flow hiện tại
	return s.buildLoginResponse(ctx, user)
}

// LoginWithGoogleToken xử lý đăng nhập bằng Google ID token.
// Chính sách phase 1:
// - Chỉ user đã tồn tại local
// - Nếu user chưa liên kết Google, phải xác nhận mật khẩu local trước khi link
func (s *AuthService) LoginWithGoogleToken(ctx context.Context, googleIDToken, password string) (*LoginResponse, error) {
	if !s.googleEnabled || s.googleVerifier == nil {
		return nil, ErrGoogleLoginDisabled
	}

	identity, err := s.googleVerifier.Verify(ctx, googleIDToken)
	if err != nil {
		return nil, auth.ErrInvalidCredentials
	}
	if s.googleHD != "" && identity.HostedDomain != s.googleHD {
		return nil, ErrGoogleDomainNotAllowed
	}

	// 1) Ưu tiên lookup theo sub nếu đã link trước đó
	linkedUser, err := s.userRepo.FindByGoogleSub(ctx, identity.Sub)
	if err == nil {
		if linkedUser.Status == "locked" {
			return nil, auth.ErrUserLocked
		}
		return s.buildLoginResponse(ctx, linkedUser)
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	// 2) Chưa link theo sub: buộc phải có local account cùng email
	user, err := s.userRepo.FindByEmail(ctx, identity.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrGoogleAccountNotProvisioned
		}
		return nil, err
	}
	if user.Status == "locked" {
		return nil, auth.ErrUserLocked
	}

	// 3) Đã có account local nhưng chưa link: yêu cầu password để xác nhận chủ tài khoản
	if user.GoogleSub == "" {
		if password == "" {
			return nil, ErrGoogleLinkPasswordRequired
		}
		if !auth.VerifyPassword(user.PasswordHash, password) {
			return nil, auth.ErrInvalidCredentials
		}
		if err := s.userRepo.LinkGoogleSub(ctx, user.UserID, identity.Sub); err != nil {
			return nil, err
		}
	}

	return s.buildLoginResponse(ctx, user)
}

func (s *AuthService) buildLoginResponse(ctx context.Context, user *model.User) (*LoginResponse, error) {
	var schoolID string
	// 1. Lấy roles của user
	roles, err := s.userRepo.RolesOfUser(ctx, user.UserID)
	if err != nil {
		return nil, err
	}

	// 2. Nếu user có role SCHOOL_ADMIN → lấy school_id từ bảng school_admins
	for _, r := range roles {
		if r == "SCHOOL_ADMIN" {
			admin, err := s.schoolAdminRepo.GetByUserID(ctx, user.UserID)
			if err == nil {
				schoolID = admin.SchoolID.String()
			}
			break
		}
	}

	// 3. Tạo JWT token (schoolID rỗng cho SUPER_ADMIN/TEACHER/PARENT)
	token, err := s.jwtAuth.SignToken(user.UserID.String(), user.Email, roles, schoolID)
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
