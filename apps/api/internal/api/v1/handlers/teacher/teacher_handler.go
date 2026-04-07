package teacherhandlers

import (
	"context"

	"github.com/google/uuid"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/service"
)

type teacherService interface {
	List(ctx context.Context, adminSchoolID *uuid.UUID, limit, offset int) ([]model.Teacher, int, error)
	ListTeachersOfClass(ctx context.Context, adminSchoolID *uuid.UUID, classID uuid.UUID) ([]model.Teacher, error)
	GetByTeacherID(ctx context.Context, adminSchoolID *uuid.UUID, teacherID uuid.UUID) (*model.Teacher, error)
	Assign(ctx context.Context, adminSchoolID *uuid.UUID, teacherID, classID uuid.UUID) error
	Unassign(ctx context.Context, adminSchoolID *uuid.UUID, teacherID, classID uuid.UUID) error
	Update(ctx context.Context, adminSchoolID *uuid.UUID, teacherID uuid.UUID, fullName, phone string, schoolID uuid.UUID) error
	Delete(ctx context.Context, adminSchoolID *uuid.UUID, teacherID uuid.UUID) error
}

type TeacherHandler struct {
	teacherService teacherService
}

// UpdateTeacherRequest input để admin cập nhật thông tin giáo viên.
type UpdateTeacherRequest struct {
	FullName string    `json:"full_name"`
	Phone    string    `json:"phone"`
	SchoolID uuid.UUID `json:"school_id"`
}

func NewTeacherHandler(teacherService *service.TeacherService) *TeacherHandler {
	return &TeacherHandler{
		teacherService: teacherService,
	}
}
