package classhandlers

import (
	"context"

	"github.com/google/uuid"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/service"
)

type classService interface {
	Create(ctx context.Context, adminSchoolID *uuid.UUID, schoolID uuid.UUID, name, schoolYear string) (*model.Class, error)
	ListBySchool(ctx context.Context, adminSchoolID *uuid.UUID, schoolID uuid.UUID, limit, offset int) ([]model.Class, int, error)
	Update(ctx context.Context, adminSchoolID *uuid.UUID, classID uuid.UUID, name, schoolYear string) error
	Delete(ctx context.Context, adminSchoolID *uuid.UUID, classID uuid.UUID) error
}

type ClassHandler struct {
	classService classService
}

func NewClassHandler(classService *service.ClassService) *ClassHandler {
	return &ClassHandler{
		classService: classService,
	}
}

type CreateClassRequest struct {
	SchoolID   uuid.UUID `json:"school_id" binding:"required"`
	Name       string    `json:"name" binding:"required,min=1,max=100"`
	SchoolYear string    `json:"school_year" binding:"required,min=4,max=20"`
}

// UpdateClassRequest input để cập nhật lớp học.
type UpdateClassRequest struct {
	Name       string `json:"name" binding:"required,min=1,max=100"`
	SchoolYear string `json:"school_year" binding:"required,min=4,max=20"`
}
