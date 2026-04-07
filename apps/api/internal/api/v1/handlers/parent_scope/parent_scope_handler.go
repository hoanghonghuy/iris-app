package parentscope

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/service"
)

type parentScopeAppointmentService interface {
	ListAvailableSlotsForParent(ctx context.Context, parentUserID, studentID uuid.UUID, from, to *time.Time, limit, offset int) ([]model.AppointmentSlot, int, error)
	CreateAppointment(ctx context.Context, parentUserID, studentID, slotID uuid.UUID, note string) (model.Appointment, error)
	ListParentAppointments(ctx context.Context, parentUserID uuid.UUID, status string, from, to *time.Time, limit, offset int) ([]model.Appointment, int, error)
	CancelAppointmentByParent(ctx context.Context, parentUserID, appointmentID uuid.UUID, cancelReason string) (model.Appointment, error)
}

type ParentScopeHandler struct {
	parentScopeService *service.ParentScopeService
	appointmentService parentScopeAppointmentService
}

func NewParentScopeHandler(parentScopeService *service.ParentScopeService, appointmentService parentScopeAppointmentService) *ParentScopeHandler {
	return &ParentScopeHandler{
		parentScopeService: parentScopeService,
		appointmentService: appointmentService,
	}
}
