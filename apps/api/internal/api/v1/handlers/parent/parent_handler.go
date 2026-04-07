package parenthandlers

import "github.com/hoanghonghuy/iris-app/apps/api/internal/service"

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
