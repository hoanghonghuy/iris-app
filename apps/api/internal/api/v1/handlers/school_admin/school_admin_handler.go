package schooladminhandlers

import (
	"github.com/google/uuid"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/service"
)

type SchoolAdminHandler struct {
	schoolAdminService *service.SchoolAdminService
}

func NewSchoolAdminHandler(schoolAdminService *service.SchoolAdminService) *SchoolAdminHandler {
	return &SchoolAdminHandler{
		schoolAdminService: schoolAdminService,
	}
}

// CreateSchoolAdminRequest input để SUPER_ADMIN tạo school admin mới.
type CreateSchoolAdminRequest struct {
	UserID   uuid.UUID `json:"user_id" binding:"required"`
	SchoolID uuid.UUID `json:"school_id" binding:"required"`
	FullName string    `json:"full_name"`
	Phone    string    `json:"phone"`
}
