package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

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
	userRepo         *repo.UserRepo
	schoolAdminRepo  *repo.SchoolAdminRepo
	refreshTokenRepo *repo.RefreshTokenRepo
	jwtAuth          *auth.Authenticator
	googleVerifier   auth.GoogleTokenVerifier
	googleEnabled    bool
	googleHD         string
	refreshTTL       time.Duration
}

type AuthServiceOptions struct {
	GoogleVerifier  auth.GoogleTokenVerifier
	GoogleEnabled   bool
	GoogleHD        string
	RefreshTTLHours int
}

func NewAuthService(
	userRepo *repo.UserRepo,
	schoolAdminRepo *repo.SchoolAdminRepo,
	refreshTokenRepo *repo.RefreshTokenRepo,
	jwtAuth *auth.Authenticator,
	opts AuthServiceOptions,
) *AuthService {
	refreshTTLHours := opts.RefreshTTLHours
	if refreshTTLHours <= 0 {
		refreshTTLHours = 24 * 7
	}
	return &AuthService{
		userRepo:         userRepo,
		schoolAdminRepo:  schoolAdminRepo,
		refreshTokenRepo: refreshTokenRepo,
		jwtAuth:          jwtAuth,
		googleVerifier:   opts.GoogleVerifier,
		googleEnabled:    opts.GoogleEnabled,
		googleHD:         opts.GoogleHD,
		refreshTTL:       time.Duration(refreshTTLHours) * time.Hour,
	}
}

// LoginResponse chứa thông tin trả về sau khi đăng nhập thành công
type LoginResponse struct {
	AccessToken      string `json:"access_token"`
	RefreshToken     string `json:"refresh_token"`
	TokenType        string `json:"token_type"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
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

	// Tìm user đã link với sub này chưa
	linkedUser, err := s.userRepo.FindByGoogleSub(ctx, identity.Sub)
	if err == nil {
		if linkedUser.Status == "locked" {
			return nil, auth.ErrUserLocked
		}
		return s.buildLoginResponse(ctx, linkedUser)
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	// Nếu chưa có user nào link với sub này, tìm theo email xem có account local nào chưa link không
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

	// Nếu đã có account local nhưng chưa link: yêu cầu password để xác nhận chủ tài khoản
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

// buildLoginResponse tạo JWT token và response body sau khi xác thực thành công.
func (s *AuthService) buildLoginResponse(ctx context.Context, user *model.User) (*LoginResponse, error) {
	var schoolID string
	// Lấy roles của user
	roles, err := s.userRepo.RolesOfUser(ctx, user.UserID)
	if err != nil {
		return nil, err
	}

	// Nếu user có role SCHOOL_ADMIN → lấy school_id từ bảng school_admins
	for _, r := range roles {
		if r == "SCHOOL_ADMIN" {
			admin, err := s.schoolAdminRepo.GetByUserID(ctx, user.UserID)
			if err == nil {
				schoolID = admin.SchoolID.String()
			}
			break
		}
	}

	// Tạo JWT token (schoolID rỗng cho SUPER_ADMIN/TEACHER/PARENT)
	token, err := s.jwtAuth.SignToken(user.UserID.String(), user.Email, roles, schoolID)
	if err != nil {
		return nil, err
	}

	refreshToken, refreshTokenHash, err := generateOpaqueToken()
	if err != nil {
		return nil, err
	}
	refreshExpiresAt := time.Now().Add(s.refreshTTL)
	if err := s.refreshTokenRepo.Create(ctx, user.UserID, refreshTokenHash, refreshExpiresAt); err != nil {
		return nil, err
	}

	return &LoginResponse{
		AccessToken:      token,
		RefreshToken:     refreshToken,
		TokenType:        "Bearer",
		ExpiresIn:        s.jwtAuth.TTLSeconds,
		RefreshExpiresIn: int(s.refreshTTL.Seconds()),
	}, nil
}

// GetUserInfo trả về thông tin user theo ID.
func (s *AuthService) GetUserInfo(ctx context.Context, userID uuid.UUID) (*model.UserInfo, error) {
	return s.userRepo.FindByID(ctx, userID)
}

// Refresh đổi refresh token hợp lệ lấy cặp access/refresh token mới.
func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (*LoginResponse, error) {
	if refreshToken == "" {
		return nil, ErrRefreshTokenInvalid
	}
	refreshTokenHash := hashToken(refreshToken)
	storedToken, err := s.refreshTokenRepo.FindActiveByTokenHash(ctx, refreshTokenHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrRefreshTokenInvalid
		}
		return nil, err
	}

	userInfo, err := s.userRepo.FindByID(ctx, storedToken.UserID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrRefreshTokenInvalid
		}
		return nil, err
	}
	if userInfo.Status == "locked" {
		return nil, auth.ErrUserLocked
	}

	if err := s.refreshTokenRepo.RevokeByID(ctx, storedToken.ID); err != nil {
		if errors.Is(err, repo.ErrNoRowsUpdated) {
			return nil, ErrRefreshTokenInvalid
		}
		return nil, err
	}

	return s.buildLoginResponse(ctx, &model.User{
		UserID: userInfo.UserID,
		Email:  userInfo.Email,
		Status: userInfo.Status,
	})
}

func generateOpaqueToken() (token string, tokenHash string, err error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", "", err
	}
	token = hex.EncodeToString(buf)
	return token, hashToken(token), nil
}

func hashToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}
