package parentcodehandlers

import "github.com/hoanghonghuy/iris-app/apps/api/internal/service"

type ParentCodeHandler struct {
	parentCodeService *service.ParentCodeService
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
