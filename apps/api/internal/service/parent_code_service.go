package service

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/google/uuid"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/auth"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/repo"
	"golang.org/x/crypto/bcrypt"
)

type ParentCodeService struct {
	parentCodeRepo    *repo.ParentCodeRepo
	userRepo          *repo.UserRepo
	parentRepo        *repo.ParentRepo
	studentParentRepo *repo.StudentParentRepo
	studentRepo       *repo.StudentRepo
	jwtAuth           *auth.Authenticator
}

func NewParentCodeService(parentCodeRepo *repo.ParentCodeRepo, userRepo *repo.UserRepo, parentRepo *repo.ParentRepo,
	studentParentRepo *repo.StudentParentRepo, studentRepo *repo.StudentRepo, jwtAuth *auth.Authenticator) *ParentCodeService {
	return &ParentCodeService{
		parentCodeRepo:    parentCodeRepo,
		userRepo:          userRepo,
		parentRepo:        parentRepo,
		studentParentRepo: studentParentRepo,
		studentRepo:       studentRepo,
		jwtAuth:           jwtAuth,
	}
}

// GenerateCodeForStudent tạo parent code cho student (admin only)
func (s *ParentCodeService) GenerateCodeForStudent(ctx context.Context, studentID uuid.UUID) (string, error) {
	// Check student tồn tại không (optional - có thể validate ở handler)

	// Generate random code
	code := generateRandomCode(8) // 8 ký tự

	maxUsage := 4

	// Lưu vào DB
	err := s.parentCodeRepo.Create(ctx, studentID, code, maxUsage)
	if err != nil {
		return "", err
	}

	return code, nil
}

// VerifyCode xác minh parent code hợp lệ và chưa vượt giới hạn
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
	// Verify parent code
	codeInfo, err := s.VerifyCode(ctx, parentCode)
	if err != nil {
		return nil, err
	}

	// Check email chưa tồn tại
	_, err = s.userRepo.FindByEmail(ctx, email)
	if err == nil {
		// Nếu không có lỗi → email đã tồn tại
		return nil, ErrEmailAlreadyExists
	}
	// Nếu lỗi != ErrNotFound → có lỗi khác → báo lỗi

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

	// Increment usage count
	err = s.parentCodeRepo.IncrementUsage(ctx, parentCode)
	if err != nil {
		// Nếu increment không thành công, không sao, chỉ log
		// (parent vẫn đăng ký thành công)
		// TODO: log error
	}

	// Generate JWT token cho parent (auto login sau khi register)
	roles := []string{"PARENT"}

	token, err := s.jwtAuth.SignToken(userID.String(), email, roles)
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
