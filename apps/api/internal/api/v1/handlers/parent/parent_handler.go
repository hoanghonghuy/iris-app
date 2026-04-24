package parenthandlers

import (
	"github.com/google/uuid"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/service"
)

type ParentHandler struct {
	parentService *service.ParentService
}

func NewParentHandler(parentService *service.ParentService) *ParentHandler {
	return &ParentHandler{
		parentService: parentService,
	}
}

type AssignStudentRequest struct {
	Relationship string `json:"relationship"`
}

// UpdateParentRequest input để admin cập nhật thông tin phụ huynh.
type UpdateParentRequest struct {
	FullName string    `json:"full_name"`
	Phone    string    `json:"phone"`
	SchoolID uuid.UUID `json:"school_id"`
}
