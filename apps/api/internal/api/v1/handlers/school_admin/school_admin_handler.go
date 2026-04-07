package schooladminhandlers

import (
	"context"

	"github.com/google/uuid"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/service"
)

type schoolAdminService interface {
	Create(ctx context.Context, userID, schoolID uuid.UUID, fullName, phone string) (*model.SchoolAdmin, error)
	List(ctx context.Context, limit, offset int) ([]model.SchoolAdmin, int, error)
	Delete(ctx context.Context, adminID uuid.UUID) error
}

type SchoolAdminHandler struct {
	schoolAdminService schoolAdminService
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
