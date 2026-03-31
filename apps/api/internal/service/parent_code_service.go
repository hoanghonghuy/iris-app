package service

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/google/uuid"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/auth"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/repo"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type ParentCodeService struct {
	parentCodeRepo    *repo.ParentCodeRepo
	userRepo          *repo.UserRepo
	parentRepo        *repo.ParentRepo
	studentParentRepo *repo.StudentParentRepo
	studentRepo       *repo.StudentRepo
	jwtAuth           *auth.Authenticator
	googleVerifier    auth.GoogleTokenVerifier
	googleEnabled     bool
	googleHD          string
}

func NewParentCodeService(parentCodeRepo *repo.ParentCodeRepo, userRepo *repo.UserRepo, parentRepo *repo.ParentRepo,
	studentParentRepo *repo.StudentParentRepo, studentRepo *repo.StudentRepo, jwtAuth *auth.Authenticator, googleVerifier auth.GoogleTokenVerifier, googleEnabled bool, googleHD string) *ParentCodeService {
	return &ParentCodeService{
		parentCodeRepo:    parentCodeRepo,
		userRepo:          userRepo,
		parentRepo:        parentRepo,
		studentParentRepo: studentParentRepo,
		studentRepo:       studentRepo,
		jwtAuth:           jwtAuth,
		googleVerifier:    googleVerifier,
		googleEnabled:     googleEnabled,
		googleHD:          googleHD,
	}
}

// GenerateCodeForStudent tạo parent code cho student (admin only)
func (s *ParentCodeService) GenerateCodeForStudent(ctx context.Context, adminSchoolID *uuid.UUID, studentID uuid.UUID) (string, error) {
	// SCHOOL_ADMIN: validate student thuộc cùng school với admin
	// adminSchoolID == nil => SUPER_ADMIN: không cần validate
	if adminSchoolID != nil {
		student, err := s.studentRepo.GetByStudentID(ctx, studentID)
		if err != nil {
			return "", ErrFailedToGetStudent
		}
		if student.SchoolID != *adminSchoolID {
			return "", ErrSchoolAccessDenied
		}
	}

	// Generate random code
	code := generateRandomCode(8) // 8 ký tự

	maxUsage := 4
	expiresAt := time.Now().AddDate(0, 0, 7) // 7 ngày

	// xoá mã cũ nếu có
	_ = s.parentCodeRepo.DeleteByStudentID(ctx, studentID)

	// Lưu vào DB
	err := s.parentCodeRepo.Create(ctx, studentID, code, maxUsage, expiresAt)
	if err != nil {
		return "", err
	}

	return code, nil
}

// RevokeCode thu hồi toàn bộ parent code đang active của student
func (s *ParentCodeService) RevokeCode(ctx context.Context, adminSchoolID *uuid.UUID, studentID uuid.UUID) error {
	if adminSchoolID != nil {
		student, err := s.studentRepo.GetByStudentID(ctx, studentID)
		if err != nil {
			return ErrFailedToGetStudent
		}
		if student.SchoolID != *adminSchoolID {
			return ErrSchoolAccessDenied
		}
	}
	return s.parentCodeRepo.DeleteByStudentID(ctx, studentID)
}

// VerifyCode xác minh parent code hợp lệ và chưa vượt giới hạn (read-only, dùng để preview thông tin).
// Lưu ý: do không atomic, không dùng để guard việc đăng ký — dùng IncrementUsageIfNotMaxed thay thế.
func (s *ParentCodeService) VerifyCode(ctx context.Context, code string) (*model.StudentParentCode, error) {
	codeInfo, err := s.parentCodeRepo.FindByCode(ctx, code)
	if err != nil {
		return nil, ErrInvalidParentCode
	}

	// Check usage chưa vượt quá max
	if codeInfo.UsageCount >= codeInfo.MaxUsage {
		return nil, ErrParentCodeMaxUsageReached
	}

	// Check code chưa expired
	if codeInfo.ExpiresAt.Before(time.Now()) {
		return nil, ErrParentCodeExpired
	}

	return codeInfo, nil
}

// RegisterParent đăng ký parent mới sử dụng parent code
func (s *ParentCodeService) RegisterParent(ctx context.Context, email, password, parentCode string) (*LoginResponse, error) {
	// Đọc thông tin code để lấy studentID (read-only, chỉ dùng để preview)
	codeInfo, err := s.parentCodeRepo.FindByCode(ctx, parentCode)
	if err != nil {
		return nil, ErrInvalidParentCode
	}
	if codeInfo.ExpiresAt.Before(time.Now()) {
		return nil, ErrParentCodeExpired
	}

	// Check email chưa tồn tại
	_, err = s.userRepo.FindByEmail(ctx, email)
	if err == nil {
		// Nếu không có lỗi → email đã tồn tại
		return nil, ErrEmailAlreadyExists
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	// Hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, ErrFailedToHashPassword
	}

	// Tạo user với status='active' (parent không cần approval)
	userID, err := s.userRepo.CreateActive(ctx, email, string(passwordHash))
	if err != nil {
		return nil, ErrFailedToCreateUser
	}

	// Assign role PARENT
	err = s.userRepo.AssignRole(ctx, userID, "PARENT")
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrFailedToAssignRole, "PARENT")
	}

	// Tạo parent record (link user <-> parent)
	// Lấy student từ codeInfo
	studentID := codeInfo.StudentID

	// Cần schoolID để tạo parent, lấy từ student record
	schoolID, err := s.studentRepo.GetSchoolIDByStudentID(ctx, studentID)
	if err != nil {
		return nil, ErrFailedToGetStudent
	}

	// Tạo parent trong bảng parents
	parentName := "" // TODO: Parent có thể cập nhật sau, hoặc dùng email làm tên
	parentID, err := s.parentRepo.Create(ctx, userID, schoolID, parentName, "")
	if err != nil {
		return nil, ErrFailedToCreateParent
	}

	// Link parent <-> student qua bảng student_parents
	err = s.studentParentRepo.Assign(ctx, studentID, parentID, "parent")
	if err != nil {
		return nil, ErrFailedToLinkParentToStudent
	}

	// kiểm tra + tăng usage_count trong 1 câu SQL (an toàn khi nhiều người dùng dùng chung mã)
	// Nếu code đã đạt max_usage giữa lúc kiểm tra và tăng → repo trả ErrNoRowsUpdated
	// Service map sang business error ErrParentCodeMaxUsageReached
	if err = s.parentCodeRepo.IncrementUsageIfNotMaxed(ctx, parentCode); err != nil {
		if errors.Is(err, repo.ErrNoRowsUpdated) {
			return nil, ErrParentCodeMaxUsageReached
		}
		return nil, err
	}

	// Generate JWT token cho parent (auto login sau khi register)
	roles := []string{"PARENT"}

	token, err := s.jwtAuth.SignToken(userID.String(), email, roles, "")
	if err != nil {
		// Nếu generate token thất bại, vẫn trả về success
		// Parent có thể login sau với email/password
		token = "" // hoặc có thể return error
	}

	return &LoginResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   s.jwtAuth.TTLSeconds,
	}, nil
}

// RegisterParentWithGoogle đăng ký parent mới bằng Google
func (s *ParentCodeService) RegisterParentWithGoogle(ctx context.Context, idToken string, parentCode string) (*LoginResponse, error) {
	if !s.googleEnabled || s.googleVerifier == nil {
		return nil, ErrGoogleLoginDisabled
	}

	// 1. Verify Google Token
	claims, err := s.googleVerifier.Verify(ctx, idToken)
	if err != nil {
		return nil, auth.ErrInvalidCredentials
	}
	if s.googleHD != "" && claims.HostedDomain != s.googleHD {
		return nil, ErrGoogleDomainNotAllowed
	}

	email := claims.Email
	if !claims.EmailVerified {
		return nil, auth.ErrInvalidCredentials
	}

	// 2. Load Parent Code
	codeInfo, err := s.parentCodeRepo.FindByCode(ctx, parentCode)
	if err != nil {
		return nil, ErrInvalidParentCode
	}
	if codeInfo.ExpiresAt.Before(time.Now()) {
		return nil, ErrParentCodeExpired
	}

	// 3. Check email chưa tồn tại
	_, err = s.userRepo.FindByEmail(ctx, email)
	if err == nil {
		return nil, ErrEmailAlreadyExists
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	// 4. Lấy thông tin student (để lấy school_id)
	schoolID, err := s.studentRepo.GetSchoolIDByStudentID(ctx, codeInfo.StudentID)
	if err != nil {
		return nil, ErrFailedToGetStudent
	}

	// 5. Hash random password fallback
	randomPwd := generateRandomCode(16)
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(randomPwd), bcrypt.DefaultCost)
	if err != nil {
		return nil, ErrFailedToHashPassword
	}

	fullName := claims.Name
	if fullName == "" {
		fullName = "Google Parent"
	}

	// 6. Execute transaction
	userID := uuid.New()
	roles := []string{"PARENT"}
	token, err := s.jwtAuth.SignToken(userID.String(), email, roles, "")
	if err != nil {
		return nil, err
	}

	txParams := repo.RegisterParentTxParams{
		UserID:       userID,
		Email:        email,
		PasswordHash: string(passwordHash),
		FullName:     fullName,
		Phone:        "", // Optional
		SchoolID:     schoolID,
		StudentID:    codeInfo.StudentID,
		Code:         parentCode,
		GoogleSub:    claims.Sub,
	}

	_, err = s.parentCodeRepo.RegisterParentTx(ctx, txParams)
	if err != nil {
		if errors.Is(err, repo.ErrNoRowsUpdated) {
			return nil, ErrParentCodeMaxUsageReached
		}
		return nil, fmt.Errorf("failed to register parent: %w", err)
	}

	return &LoginResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   s.jwtAuth.TTLSeconds,
	}, nil
}

// GetStudentInfo lấy thông tin student cho VerifyCode endpoint
func (s *ParentCodeService) GetStudentInfo(ctx context.Context, studentID uuid.UUID) (*model.Student, error) {
	return s.studentRepo.GetByStudentID(ctx, studentID)
}

// generateRandomCode generate random alphanumeric code có độ dài n
func generateRandomCode(length int) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		b[i] = charset[n.Int64()]
	}
	return string(b)
}
