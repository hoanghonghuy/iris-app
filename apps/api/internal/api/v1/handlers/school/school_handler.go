package schoolhandlers

import (
	"context"

	"github.com/google/uuid"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/service"
)

type schoolService interface {
	Create(ctx context.Context, name, address string) (*model.School, error)
	List(ctx context.Context, adminSchoolID *uuid.UUID, limit, offset int) ([]model.School, int, error)
	Update(ctx context.Context, schoolID uuid.UUID, name, address string) error
	Delete(ctx context.Context, schoolID uuid.UUID) error
}

type SchoolHandler struct {
	schoolService schoolService
}

func NewSchoolHandler(schoolService *service.SchoolService) *SchoolHandler {
	return &SchoolHandler{
		schoolService: schoolService,
	}
}

type CreateSchoolRequest struct {
	Name    string `json:"name" binding:"required,min=2"`
	Address string `json:"address"`
}

// UpdateSchoolRequest input để cập nhật trường học.
type UpdateSchoolRequest struct {
	Name    string `json:"name" binding:"required,min=2"`
	Address string `json:"address"`
}
