package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/repo"
)

type AppointmentService struct {
	appointmentRepo *repo.AppointmentRepo
}

const (
	defaultSlotBufferMinutes     = 10
	defaultSlotMaxBookingsPerDay = 12
	parentCancelCutoff           = 2 * time.Hour
	maxSlotBufferMinutes         = 180
	maxSlotBookingsPerDayHardCap = 100
)

func NewAppointmentService(appointmentRepo *repo.AppointmentRepo) *AppointmentService {
	return &AppointmentService{appointmentRepo: appointmentRepo}
}

func (s *AppointmentService) CreateSlot(
	ctx context.Context,
	teacherUserID,
	classID uuid.UUID,
	startTime,
	endTime time.Time,
	note string,
	bufferMinutes,
	maxBookingsPerDay int,
) (model.AppointmentSlot, error) {
	if teacherUserID == uuid.Nil {
		return model.AppointmentSlot{}, ErrInvalidUserID
	}
	if classID == uuid.Nil {
		return model.AppointmentSlot{}, ErrInvalidClassID
	}
	if !endTime.After(startTime) {
		return model.AppointmentSlot{}, fmt.Errorf("%w: end_time must be greater than start_time", ErrInvalidValue)
	}
	if startTime.Before(time.Now().Add(-5 * time.Minute)) {
		return model.AppointmentSlot{}, fmt.Errorf("%w: start_time cannot be in the past", ErrInvalidValue)
	}

	bufferMinutes = normalizeSlotBufferMinutes(bufferMinutes)
	maxBookingsPerDay = normalizeSlotMaxBookingsPerDay(maxBookingsPerDay)

	dayStart := time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 0, 0, 0, 0, startTime.Location())
	dayEnd := dayStart.Add(24*time.Hour - time.Nanosecond)

	dailyCount, err := s.appointmentRepo.CountTeacherActiveSlotsForDay(ctx, teacherUserID, dayStart, dayEnd)
	if err != nil {
		return model.AppointmentSlot{}, fmt.Errorf("failed to check daily slot count: %w", err)
	}
	if dailyCount >= maxBookingsPerDay {
		return model.AppointmentSlot{}, fmt.Errorf("%w: maximum %d slots/day reached", ErrInvalidValue, maxBookingsPerDay)
	}

	conflict, err := s.appointmentRepo.FindTeacherSlotConflict(ctx, teacherUserID, startTime, endTime, bufferMinutes)
	if err != nil {
		return model.AppointmentSlot{}, fmt.Errorf("failed to validate slot overlap: %w", err)
	}
	if conflict != nil {
		return model.AppointmentSlot{}, fmt.Errorf(
			"%w: slot conflicts with existing slot %s - %s (buffer %d minutes)",
			ErrInvalidValue,
			conflict.StartTime.Format(time.RFC3339),
			conflict.EndTime.Format(time.RFC3339),
			bufferMinutes,
		)
	}

	slot, err := s.appointmentRepo.CreateSlot(ctx, teacherUserID, classID, startTime, endTime, strings.TrimSpace(note))
	if err != nil {
		if err == repo.ErrNoRowsUpdated {
			return model.AppointmentSlot{}, ErrForbidden
		}
		return model.AppointmentSlot{}, fmt.Errorf("failed to create appointment slot: %w", err)
	}
	return slot, nil
}

func (s *AppointmentService) ListTeacherAppointments(ctx context.Context, teacherUserID uuid.UUID, status string, from, to *time.Time, limit, offset int) ([]model.Appointment, int, error) {
	if teacherUserID == uuid.Nil {
		return nil, 0, ErrInvalidUserID
	}
	if status != "" && !isValidAppointmentStatus(status) {
		return nil, 0, fmt.Errorf("%w: invalid appointment status", ErrInvalidValue)
	}
	limit = normalizeListLimit(limit)
	if offset < 0 {
		offset = 0
	}
	return s.appointmentRepo.ListTeacherAppointments(ctx, teacherUserID, status, from, to, limit, offset)
}

func (s *AppointmentService) ListParentAppointments(ctx context.Context, parentUserID uuid.UUID, status string, from, to *time.Time, limit, offset int) ([]model.Appointment, int, error) {
	if parentUserID == uuid.Nil {
		return nil, 0, ErrInvalidUserID
	}
	if status != "" && !isValidAppointmentStatus(status) {
		return nil, 0, fmt.Errorf("%w: invalid appointment status", ErrInvalidValue)
	}
	limit = normalizeListLimit(limit)
	if offset < 0 {
		offset = 0
	}
	return s.appointmentRepo.ListParentAppointments(ctx, parentUserID, status, from, to, limit, offset)
}

func (s *AppointmentService) ListAvailableSlotsForParent(ctx context.Context, parentUserID, studentID uuid.UUID, from, to *time.Time, limit, offset int) ([]model.AppointmentSlot, int, error) {
	if parentUserID == uuid.Nil {
		return nil, 0, ErrInvalidUserID
	}
	if studentID == uuid.Nil {
		return nil, 0, ErrInvalidValue
	}
	limit = normalizeListLimit(limit)
	if offset < 0 {
		offset = 0
	}
	return s.appointmentRepo.ListAvailableSlotsForParent(ctx, parentUserID, studentID, from, to, limit, offset)
}

func (s *AppointmentService) CreateAppointment(ctx context.Context, parentUserID, studentID, slotID uuid.UUID, note string) (model.Appointment, error) {
	if parentUserID == uuid.Nil {
		return model.Appointment{}, ErrInvalidUserID
	}
	if studentID == uuid.Nil || slotID == uuid.Nil {
		return model.Appointment{}, fmt.Errorf("%w: student_id and slot_id are required", ErrInvalidValue)
	}

	a, err := s.appointmentRepo.CreateAppointment(ctx, parentUserID, studentID, slotID, strings.TrimSpace(note))
	if err != nil {
		if err == repo.ErrAppointmentSlotUnavailable {
			return model.Appointment{}, ErrAppointmentSlotUnavailable
		}
		if err == repo.ErrNoRowsUpdated {
			return model.Appointment{}, ErrForbidden
		}
		return model.Appointment{}, fmt.Errorf("failed to create appointment: %w", err)
	}
	return a, nil
}

func (s *AppointmentService) UpdateAppointmentStatusByTeacher(ctx context.Context, teacherUserID, appointmentID uuid.UUID, status, cancelReason string) (model.Appointment, error) {
	if teacherUserID == uuid.Nil {
		return model.Appointment{}, ErrInvalidUserID
	}
	if appointmentID == uuid.Nil {
		return model.Appointment{}, fmt.Errorf("%w: appointment_id is required", ErrInvalidValue)
	}
	if !isValidAppointmentStatus(status) {
		return model.Appointment{}, fmt.Errorf("%w: invalid appointment status", ErrInvalidValue)
	}
	if status == "cancelled" && strings.TrimSpace(cancelReason) == "" {
		return model.Appointment{}, fmt.Errorf("%w: cancel_reason is required when cancelled", ErrInvalidValue)
	}
	a, err := s.appointmentRepo.UpdateAppointmentStatusByTeacher(ctx, teacherUserID, appointmentID, status, strings.TrimSpace(cancelReason))
	if err != nil {
		if err == repo.ErrNoRowsUpdated {
			return model.Appointment{}, ErrForbidden
		}
		return model.Appointment{}, fmt.Errorf("failed to update appointment status: %w", err)
	}
	return a, nil
}

func (s *AppointmentService) CancelAppointmentByParent(ctx context.Context, parentUserID, appointmentID uuid.UUID, cancelReason string) (model.Appointment, error) {
	if parentUserID == uuid.Nil {
		return model.Appointment{}, ErrInvalidUserID
	}
	if appointmentID == uuid.Nil {
		return model.Appointment{}, fmt.Errorf("%w: appointment_id is required", ErrInvalidValue)
	}
	if strings.TrimSpace(cancelReason) == "" {
		cancelReason = "parent_cancelled"
	}
	a, err := s.appointmentRepo.CancelAppointmentByParent(ctx, parentUserID, appointmentID, cancelReason, parentCancelCutoff)
	if err != nil {
		if err == repo.ErrAppointmentCancellationWindowPassed {
			return model.Appointment{}, ErrAppointmentCancellationWindowPassed
		}
		if err == repo.ErrNoRowsUpdated {
			return model.Appointment{}, ErrForbidden
		}
		return model.Appointment{}, fmt.Errorf("failed to cancel appointment: %w", err)
	}
	return a, nil
}

func (s *AppointmentService) CountParentUpcomingAppointments(ctx context.Context, parentUserID uuid.UUID) (int, error) {
	if parentUserID == uuid.Nil {
		return 0, ErrInvalidUserID
	}
	return s.appointmentRepo.CountParentUpcomingAppointments(ctx, parentUserID)
}

func normalizeListLimit(limit int) int {
	if limit <= 0 {
		return 20
	}
	if limit > 100 {
		return 100
	}
	return limit
}

func isValidAppointmentStatus(status string) bool {
	switch status {
	case "pending", "confirmed", "cancelled", "completed", "no_show":
		return true
	default:
		return false
	}
}

func normalizeSlotBufferMinutes(bufferMinutes int) int {
	if bufferMinutes < 0 {
		return 0
	}
	if bufferMinutes == 0 {
		return defaultSlotBufferMinutes
	}
	if bufferMinutes > maxSlotBufferMinutes {
		return maxSlotBufferMinutes
	}
	return bufferMinutes
}

func normalizeSlotMaxBookingsPerDay(maxBookingsPerDay int) int {
	if maxBookingsPerDay <= 0 {
		return defaultSlotMaxBookingsPerDay
	}
	if maxBookingsPerDay > maxSlotBookingsPerDayHardCap {
		return maxSlotBookingsPerDayHardCap
	}
	return maxBookingsPerDay
}
