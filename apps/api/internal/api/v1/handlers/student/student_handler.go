package studenthandlers

import (
	"github.com/google/uuid"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/service"
)

type StudentHandler struct {
	studentService *service.StudentService
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
