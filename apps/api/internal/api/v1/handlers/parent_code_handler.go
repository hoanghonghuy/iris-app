package handlers

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/auth"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/response"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/service"
)

type ParentCodeHandler struct {
	parentCodeService *service.ParentCodeService
}

func NewParentCodeHandler(parentCodeService *service.ParentCodeService) *ParentCodeHandler {
	return &ParentCodeHandler{
		parentCodeService: parentCodeService,
	}
}

// GenerateCodeForStudentRequest request để admin tạo parent code cho student
// RegisterParentRequest request để parent tự đăng ký
type RegisterParentRequest struct {
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required,min=6"`
	ParentCode string `json:"parent_code" binding:"required"`
}

// GenerateCodeForStudent tạo parent code cho student (admin only)
func (h *ParentCodeHandler) GenerateCodeForStudent(c *gin.Context) {
	adminSchoolID := extractAdminSchoolID(c)

	studentID, err := uuid.Parse(c.Param("student_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid student ID format")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	code, err := h.parentCodeService.GenerateCodeForStudent(ctx, adminSchoolID, studentID)
	if err != nil {
		if errors.Is(err, service.ErrSchoolAccessDenied) {
			response.Fail(c, http.StatusForbidden, "access denied")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "failed to generate parent code")
		return
	}

	response.OK(c, gin.H{
		"student_id":  studentID,
		"parent_code": code,
		"message":     "share this code with parent to allow registration",
		"max_usage":   4,
		"expires_at":  time.Now().AddDate(0, 0, 7), // 7 ngày từ khi tạo
	})
}

// RegisterParentWithGoogle xử lý việc phụ huynh đăng ký tài khoản liên kết Google bằng Parent Code.
func (h *ParentCodeHandler) RegisterParentWithGoogle(c *gin.Context) {
	var req struct {
		IDToken    string `json:"id_token" binding:"required"`
		ParentCode string `json:"parent_code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid request body")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	resp, err := h.parentCodeService.RegisterParentWithGoogle(ctx, req.IDToken, req.ParentCode)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrGoogleLoginDisabled):
			response.Fail(c, http.StatusForbidden, "tính năng đăng nhập Google đang tắt")
			return
		case errors.Is(err, service.ErrGoogleDomainNotAllowed):
			response.Fail(c, http.StatusForbidden, "tên miền Google không hợp lệ")
			return
		case errors.Is(err, auth.ErrInvalidCredentials):
			response.Fail(c, http.StatusUnauthorized, "xác thực Google thất bại")
			return
		case errors.Is(err, service.ErrInvalidParentCode):
			response.Fail(c, http.StatusBadRequest, "mã phụ huynh không hợp lệ/hết hạn/hết lượt dùng")
			return
		case errors.Is(err, service.ErrParentCodeExpired):
			response.Fail(c, http.StatusBadRequest, "mã phụ huynh không hợp lệ/hết hạn/hết lượt dùng")
			return
		case errors.Is(err, service.ErrParentCodeMaxUsageReached):
			response.Fail(c, http.StatusBadRequest, "mã phụ huynh không hợp lệ/hết hạn/hết lượt dùng")
			return
		case errors.Is(err, service.ErrEmailAlreadyExists):
			response.Fail(c, http.StatusConflict, "Email này đã được đăng ký. Vui lòng quay lại trang Đăng nhập")
			return
		default:
			response.Fail(c, http.StatusInternalServerError, "server error")
			return
		}
	}

	response.OK(c, resp)
}
// RevokeParentCode thu hồi parent code cua hs
func (h *ParentCodeHandler) RevokeParentCode(c *gin.Context) {
	adminSchoolID := extractAdminSchoolID(c)

	studentID, err := uuid.Parse(c.Param("student_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid student ID format")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	if err := h.parentCodeService.RevokeCode(ctx, adminSchoolID, studentID); err != nil {
		if errors.Is(err, service.ErrSchoolAccessDenied) {
			response.Fail(c, http.StatusForbidden, "access denied")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "failed to revoke parent code")
		return
	}

	response.OK(c, gin.H{"message": "parent code revoked successfully"})
}

// RegisterParent parent tự đăng ký (public endpoint)
func (h *ParentCodeHandler) RegisterParent(c *gin.Context) {
	var req RegisterParentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid request body")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	resp, err := h.parentCodeService.RegisterParent(ctx, req.Email, req.Password, req.ParentCode)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidParentCode):
			response.Fail(c, http.StatusBadRequest, "invalid parent code")
			return
		case errors.Is(err, service.ErrParentCodeExpired):
			response.Fail(c, http.StatusBadRequest, "parent code has expired")
			return
		case errors.Is(err, service.ErrParentCodeMaxUsageReached):
			response.Fail(c, http.StatusBadRequest, "parent code has reached maximum usage")
			return
		case errors.Is(err, service.ErrEmailAlreadyExists):
			response.Fail(c, http.StatusConflict, "email already exists")
			return
		case errors.Is(err, service.ErrPasswordCannotBeEmpty):
			response.Fail(c, http.StatusBadRequest, "password cannot be empty")
			return
		case errors.Is(err, service.ErrFailedToHashPassword):
			response.Fail(c, http.StatusInternalServerError, "failed to hash password")
			return
		case errors.Is(err, service.ErrFailedToCreateUser):
			response.Fail(c, http.StatusInternalServerError, "failed to create user")
			return
		case errors.Is(err, service.ErrFailedToCreateParent):
			response.Fail(c, http.StatusInternalServerError, "failed to create parent")
			return
		case errors.Is(err, service.ErrFailedToLinkParentToStudent):
			response.Fail(c, http.StatusInternalServerError, "failed to link parent to student")
			return
		default:
			response.Fail(c, http.StatusInternalServerError, "failed to register parent")
			return
		}
	}

	// Trả về token để auto-login
	response.Created(c, resp)
}

// VerifyCode kiểm tra tính hợp lệ của parent code và trả về thông tin student
func (h *ParentCodeHandler) VerifyCode(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		response.Fail(c, http.StatusBadRequest, "parent code is required")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	// Verify code (kiểm tra expired, max usage, v.v.)
	codeInfo, err := h.parentCodeService.VerifyCode(ctx, code)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidParentCode):
			response.Fail(c, http.StatusBadRequest, "invalid parent code")
			return
		case errors.Is(err, service.ErrParentCodeExpired):
			response.Fail(c, http.StatusBadRequest, "parent code has expired")
			return
		case errors.Is(err, service.ErrParentCodeMaxUsageReached):
			response.Fail(c, http.StatusBadRequest, "parent code has reached maximum usage")
			return
		default:
			response.Fail(c, http.StatusInternalServerError, "failed to verify code")
			return
		}
	}

	// Lấy thông tin student từ DB
	student, err := h.parentCodeService.GetStudentInfo(ctx, codeInfo.StudentID)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "failed to get student info")
		return
	}

	response.OK(c, gin.H{
		"student": gin.H{
			"student_id": codeInfo.StudentID,
			"full_name":  student.FullName,
			"dob":        student.DOB.Format("2006-01-02"), // YYYY-MM-DD
			"gender":     student.Gender,
			"school_id":  student.SchoolID,
			"class_id":   student.CurrentClassID,
		},
		"code_info": gin.H{
			"code":        codeInfo.Code,
			"usage_count": codeInfo.UsageCount,
			"max_usage":   codeInfo.MaxUsage,
			"expires_at":  codeInfo.ExpiresAt,
		},
	})
}
