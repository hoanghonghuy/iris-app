package teacherhandlers

import (
	"github.com/google/uuid"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/service"
)

type TeacherHandler struct {
	teacherService *service.TeacherService
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
