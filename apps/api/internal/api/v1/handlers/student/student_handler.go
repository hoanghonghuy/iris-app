package studenthandlers

import (
	"context"

	"github.com/google/uuid"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/service"
)

type studentService interface {
	Create(ctx context.Context, adminSchoolID *uuid.UUID, schoolID, classID uuid.UUID, fullName string, dobStr string, gender string) (*model.Student, error)
	ListByClass(ctx context.Context, adminSchoolID *uuid.UUID, classID uuid.UUID, limit, offset int) ([]model.Student, int, error)
	GetProfile(ctx context.Context, adminSchoolID *uuid.UUID, studentID uuid.UUID) (*model.StudentProfile, error)
	Update(ctx context.Context, adminSchoolID *uuid.UUID, studentID uuid.UUID, fullName, dobStr, gender string) error
	Delete(ctx context.Context, adminSchoolID *uuid.UUID, studentID uuid.UUID) error
}

type StudentHandler struct {
	studentService studentService
}

func NewStudentHandler(studentService *service.StudentService) *StudentHandler {
	return &StudentHandler{
		studentService: studentService,
	}
}

type CreateStudentReq struct {
	SchoolID       uuid.UUID `json:"school_id" binding:"required"`
	CurrentClassID uuid.UUID `json:"current_class_id" binding:"required"`
	FullName       string    `json:"full_name" binding:"required,min=1,max=120"`
	DOB            string    `json:"dob" binding:"required"`
	Gender         string    `json:"gender" binding:"required"`
}

// UpdateStudentReq input để cập nhật thông tin học sinh.
type UpdateStudentReq struct {
	FullName string `json:"full_name" binding:"required,min=1,max=120"`
	DOB      string `json:"dob" binding:"required"`
	Gender   string `json:"gender" binding:"required"`
}
