package parentcodehandlers

import (
	"context"

	"github.com/google/uuid"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/service"
)

type parentCodeService interface {
	GenerateCodeForStudent(context.Context, *uuid.UUID, uuid.UUID) (string, error)
	RevokeCode(context.Context, *uuid.UUID, uuid.UUID) error
	RegisterParentWithGoogle(context.Context, string, string) (*service.LoginResponse, error)
	RegisterParent(context.Context, string, string, string) (*service.LoginResponse, error)
	VerifyCode(context.Context, string) (*model.StudentParentCode, error)
	GetStudentInfo(context.Context, uuid.UUID) (*model.Student, error)
}

type ParentCodeHandler struct {
	parentCodeService parentCodeService
}

func NewParentCodeHandler(parentCodeService *service.ParentCodeService) *ParentCodeHandler {
	return &ParentCodeHandler{
		parentCodeService: parentCodeService,
	}
}

// RegisterParentRequest request để parent tự đăng ký.
type RegisterParentRequest struct {
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required,min=6"`
	ParentCode string `json:"parent_code" binding:"required"`
}
