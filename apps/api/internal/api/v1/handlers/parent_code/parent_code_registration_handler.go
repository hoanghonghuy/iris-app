package parentcodehandlers

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/auth"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/response"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/service"
)

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

	student, err := h.parentCodeService.GetStudentInfo(ctx, codeInfo.StudentID)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "failed to get student info")
		return
	}

	response.OK(c, gin.H{
		"student": gin.H{
			"student_id": codeInfo.StudentID,
			"full_name":  student.FullName,
			"dob":        student.DOB.Format("2006-01-02"),
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
