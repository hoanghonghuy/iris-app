package handlers

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

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
type GenerateCodeForStudentRequest struct {
	StudentID uuid.UUID `uri:"student_id" binding:"required"`
}

// RegisterParentRequest request để parent tự đăng ký
type RegisterParentRequest struct {
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required,min=6"`
	ParentCode string `json:"parent_code" binding:"required"`
}

// GenerateCodeForStudent tạo parent code cho student (admin only)
func (h *ParentCodeHandler) GenerateCodeForStudent(c *gin.Context) {
	var uriParams GenerateCodeForStudentRequest
	if err := c.ShouldBindUri(&uriParams); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid student ID")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	code, err := h.parentCodeService.GenerateCodeForStudent(ctx, uriParams.StudentID)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "failed to generate parent code")
		return
	}

	response.OK(c, gin.H{
		"student_id":  uriParams.StudentID,
		"parent_code": code,
		"message":     "share this code with parent to allow registration",
		"max_usage":   4,
		"expires_at":  time.Now().AddDate(1, 0, 0), // 1 year từ khi tạo
	})
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
