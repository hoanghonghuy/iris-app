package teacherscope

import "github.com/hoanghonghuy/iris-app/apps/api/internal/service"

type TeacherScopeHandler struct {
	teacherScopeService *service.TeacherScopeService
	appointmentService  *service.AppointmentService
}

func NewTeacherScopeHandler(teacherScopeService *service.TeacherScopeService, appointmentService *service.AppointmentService) *TeacherScopeHandler {
	return &TeacherScopeHandler{
		teacherScopeService: teacherScopeService,
		appointmentService:  appointmentService,
	}
}

type MarkAttendanceRequest struct {
	StudentID  string  `json:"student_id" binding:"required"`
	Date       string  `json:"date" binding:"required"`   // YYYY-MM-DD
	Status     string  `json:"status" binding:"required"` // present/absent/late/excused
	CheckInAt  *string `json:"check_in_at,omitempty"`     // RFC3339 or empty
	CheckOutAt *string `json:"check_out_at,omitempty"`    // RFC3339 or empty
	Note       string  `json:"note"`
}

type CreateHealthRequest struct {
	StudentID   string   `json:"student_id" binding:"required"`
	RecordedAt  *string  `json:"recorded_at"` // RFC3339 optional
	Temperature *float64 `json:"temperature"`
	Symptoms    string   `json:"symptoms"`
	Severity    *string  `json:"severity"` // normal|watch|urgent optional
	Note        string   `json:"note"`
}

type CreatePostRequest struct {
	ScopeType string `json:"scope_type" binding:"required"` // class|student
	ClassID   string `json:"class_id"`                      // required if scope_type=class
	StudentID string `json:"student_id"`                    // required if scope_type=student
	Type      string `json:"type" binding:"required"`       // announcement|activity|daily_note|health_note
	Content   string `json:"content" binding:"required"`
}

type UpdatePostRequest struct {
	Content string `json:"content" binding:"required"`
}

// UpdateMyProfileRequest input để giáo viên cập nhật thông tin cá nhân (chỉ phone)
type UpdateMyProfileRequest struct {
	Phone string `json:"phone"`
}
