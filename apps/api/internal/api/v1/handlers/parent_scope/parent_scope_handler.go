package parentscope

import "github.com/hoanghonghuy/iris-app/apps/api/internal/service"

type ParentScopeHandler struct {
	parentScopeService *service.ParentScopeService
	appointmentService *service.AppointmentService
}

func NewParentScopeHandler(parentScopeService *service.ParentScopeService, appointmentService *service.AppointmentService) *ParentScopeHandler {
	return &ParentScopeHandler{
		parentScopeService: parentScopeService,
		appointmentService: appointmentService,
	}
}
