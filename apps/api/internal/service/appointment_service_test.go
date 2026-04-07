package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/repo"
)

type fakeAppointmentRepo struct {
	createSlotFn                       func(context.Context, uuid.UUID, uuid.UUID, time.Time, time.Time, string) (model.AppointmentSlot, error)
	listTeacherAppointmentsFn          func(context.Context, uuid.UUID, string, *time.Time, *time.Time, int, int) ([]model.Appointment, int, error)
	listParentAppointmentsFn           func(context.Context, uuid.UUID, string, *time.Time, *time.Time, int, int) ([]model.Appointment, int, error)
	listAvailableSlotsForParentFn      func(context.Context, uuid.UUID, uuid.UUID, *time.Time, *time.Time, int, int) ([]model.AppointmentSlot, int, error)
	createAppointmentFn                func(context.Context, uuid.UUID, uuid.UUID, uuid.UUID, string) (model.Appointment, error)
	updateAppointmentStatusByTeacherFn func(context.Context, uuid.UUID, uuid.UUID, string, string) (model.Appointment, error)
	cancelAppointmentByParentFn        func(context.Context, uuid.UUID, uuid.UUID, string) (model.Appointment, error)
	countParentUpcomingAppointmentsFn  func(context.Context, uuid.UUID) (int, error)
}

func (f *fakeAppointmentRepo) CreateSlot(ctx context.Context, teacherUserID, classID uuid.UUID, startTime, endTime time.Time, note string) (model.AppointmentSlot, error) {
	if f.createSlotFn == nil {
		return model.AppointmentSlot{}, errors.New("unexpected CreateSlot call")
	}
	return f.createSlotFn(ctx, teacherUserID, classID, startTime, endTime, note)
}

func (f *fakeAppointmentRepo) ListTeacherAppointments(ctx context.Context, teacherUserID uuid.UUID, status string, from, to *time.Time, limit, offset int) ([]model.Appointment, int, error) {
	if f.listTeacherAppointmentsFn == nil {
		return nil, 0, errors.New("unexpected ListTeacherAppointments call")
	}
	return f.listTeacherAppointmentsFn(ctx, teacherUserID, status, from, to, limit, offset)
}

func (f *fakeAppointmentRepo) ListParentAppointments(ctx context.Context, parentUserID uuid.UUID, status string, from, to *time.Time, limit, offset int) ([]model.Appointment, int, error) {
	if f.listParentAppointmentsFn == nil {
		return nil, 0, errors.New("unexpected ListParentAppointments call")
	}
	return f.listParentAppointmentsFn(ctx, parentUserID, status, from, to, limit, offset)
}

func (f *fakeAppointmentRepo) ListAvailableSlotsForParent(ctx context.Context, parentUserID, studentID uuid.UUID, from, to *time.Time, limit, offset int) ([]model.AppointmentSlot, int, error) {
	if f.listAvailableSlotsForParentFn == nil {
		return nil, 0, errors.New("unexpected ListAvailableSlotsForParent call")
	}
	return f.listAvailableSlotsForParentFn(ctx, parentUserID, studentID, from, to, limit, offset)
}

func (f *fakeAppointmentRepo) CreateAppointment(ctx context.Context, parentUserID, studentID, slotID uuid.UUID, note string) (model.Appointment, error) {
	if f.createAppointmentFn == nil {
		return model.Appointment{}, errors.New("unexpected CreateAppointment call")
	}
	return f.createAppointmentFn(ctx, parentUserID, studentID, slotID, note)
}

func (f *fakeAppointmentRepo) UpdateAppointmentStatusByTeacher(ctx context.Context, teacherUserID, appointmentID uuid.UUID, status, cancelReason string) (model.Appointment, error) {
	if f.updateAppointmentStatusByTeacherFn == nil {
		return model.Appointment{}, errors.New("unexpected UpdateAppointmentStatusByTeacher call")
	}
	return f.updateAppointmentStatusByTeacherFn(ctx, teacherUserID, appointmentID, status, cancelReason)
}

func (f *fakeAppointmentRepo) CancelAppointmentByParent(ctx context.Context, parentUserID, appointmentID uuid.UUID, cancelReason string) (model.Appointment, error) {
	if f.cancelAppointmentByParentFn == nil {
		return model.Appointment{}, errors.New("unexpected CancelAppointmentByParent call")
	}
	return f.cancelAppointmentByParentFn(ctx, parentUserID, appointmentID, cancelReason)
}

func (f *fakeAppointmentRepo) CountParentUpcomingAppointments(ctx context.Context, parentUserID uuid.UUID) (int, error) {
	if f.countParentUpcomingAppointmentsFn == nil {
		return 0, errors.New("unexpected CountParentUpcomingAppointments call")
	}
	return f.countParentUpcomingAppointmentsFn(ctx, parentUserID)
}

func TestCreateSlotValidation(t *testing.T) {
	svc := &AppointmentService{}
	now := time.Now()

	tests := []struct {
		name          string
		teacherUserID uuid.UUID
		classID       uuid.UUID
		startTime     time.Time
		endTime       time.Time
		wantErr       error
	}{
		{
			name:          "invalid teacher user id",
			teacherUserID: uuid.Nil,
			classID:       uuid.New(),
			startTime:     now.Add(10 * time.Minute),
			endTime:       now.Add(40 * time.Minute),
			wantErr:       ErrInvalidUserID,
		},
		{
			name:          "invalid class id",
			teacherUserID: uuid.New(),
			classID:       uuid.Nil,
			startTime:     now.Add(10 * time.Minute),
			endTime:       now.Add(40 * time.Minute),
			wantErr:       ErrInvalidClassID,
		},
		{
			name:          "end time must be after start time",
			teacherUserID: uuid.New(),
			classID:       uuid.New(),
			startTime:     now.Add(40 * time.Minute),
			endTime:       now.Add(10 * time.Minute),
			wantErr:       ErrInvalidValue,
		},
		{
			name:          "start time in the past",
			teacherUserID: uuid.New(),
			classID:       uuid.New(),
			startTime:     now.Add(-10 * time.Minute),
			endTime:       now.Add(20 * time.Minute),
			wantErr:       ErrInvalidValue,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.CreateSlot(context.Background(), tt.teacherUserID, tt.classID, tt.startTime, tt.endTime, "")
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestListTeacherAppointmentsValidation(t *testing.T) {
	svc := &AppointmentService{}

	_, _, err := svc.ListTeacherAppointments(context.Background(), uuid.Nil, "", nil, nil, 20, 0)
	if !errors.Is(err, ErrInvalidUserID) {
		t.Fatalf("error = %v, want %v", err, ErrInvalidUserID)
	}

	_, _, err = svc.ListTeacherAppointments(context.Background(), uuid.New(), "wrong-status", nil, nil, 20, 0)
	if !errors.Is(err, ErrInvalidValue) {
		t.Fatalf("error = %v, want %v", err, ErrInvalidValue)
	}

	from := time.Date(2026, time.April, 7, 12, 0, 0, 0, time.UTC)
	to := from.Add(-time.Hour)
	_, _, err = svc.ListTeacherAppointments(context.Background(), uuid.New(), "", &from, &to, 20, 0)
	if !errors.Is(err, ErrInvalidValue) {
		t.Fatalf("error = %v, want %v", err, ErrInvalidValue)
	}
}

func TestListParentAppointmentsValidation(t *testing.T) {
	svc := &AppointmentService{}

	_, _, err := svc.ListParentAppointments(context.Background(), uuid.Nil, "", nil, nil, 20, 0)
	if !errors.Is(err, ErrInvalidUserID) {
		t.Fatalf("error = %v, want %v", err, ErrInvalidUserID)
	}

	_, _, err = svc.ListParentAppointments(context.Background(), uuid.New(), "wrong-status", nil, nil, 20, 0)
	if !errors.Is(err, ErrInvalidValue) {
		t.Fatalf("error = %v, want %v", err, ErrInvalidValue)
	}

	from := time.Date(2026, time.April, 7, 12, 0, 0, 0, time.UTC)
	to := from.Add(-time.Hour)
	_, _, err = svc.ListParentAppointments(context.Background(), uuid.New(), "", &from, &to, 20, 0)
	if !errors.Is(err, ErrInvalidValue) {
		t.Fatalf("error = %v, want %v", err, ErrInvalidValue)
	}
}

func TestListAvailableSlotsForParentValidation(t *testing.T) {
	svc := &AppointmentService{}

	_, _, err := svc.ListAvailableSlotsForParent(context.Background(), uuid.Nil, uuid.New(), nil, nil, 20, 0)
	if !errors.Is(err, ErrInvalidUserID) {
		t.Fatalf("error = %v, want %v", err, ErrInvalidUserID)
	}

	_, _, err = svc.ListAvailableSlotsForParent(context.Background(), uuid.New(), uuid.Nil, nil, nil, 20, 0)
	if !errors.Is(err, ErrInvalidValue) {
		t.Fatalf("error = %v, want %v", err, ErrInvalidValue)
	}

	from := time.Date(2026, time.April, 7, 12, 0, 0, 0, time.UTC)
	to := from.Add(-time.Hour)
	_, _, err = svc.ListAvailableSlotsForParent(context.Background(), uuid.New(), uuid.New(), &from, &to, 20, 0)
	if !errors.Is(err, ErrInvalidValue) {
		t.Fatalf("error = %v, want %v", err, ErrInvalidValue)
	}
}

func TestCreateAppointmentValidation(t *testing.T) {
	svc := &AppointmentService{}

	_, err := svc.CreateAppointment(context.Background(), uuid.Nil, uuid.New(), uuid.New(), "")
	if !errors.Is(err, ErrInvalidUserID) {
		t.Fatalf("error = %v, want %v", err, ErrInvalidUserID)
	}

	_, err = svc.CreateAppointment(context.Background(), uuid.New(), uuid.Nil, uuid.New(), "")
	if !errors.Is(err, ErrInvalidValue) {
		t.Fatalf("error = %v, want %v", err, ErrInvalidValue)
	}

	_, err = svc.CreateAppointment(context.Background(), uuid.New(), uuid.New(), uuid.Nil, "")
	if !errors.Is(err, ErrInvalidValue) {
		t.Fatalf("error = %v, want %v", err, ErrInvalidValue)
	}
}

func TestUpdateAppointmentStatusByTeacherValidation(t *testing.T) {
	svc := &AppointmentService{}

	_, err := svc.UpdateAppointmentStatusByTeacher(context.Background(), uuid.Nil, uuid.New(), "confirmed", "")
	if !errors.Is(err, ErrInvalidUserID) {
		t.Fatalf("error = %v, want %v", err, ErrInvalidUserID)
	}

	_, err = svc.UpdateAppointmentStatusByTeacher(context.Background(), uuid.New(), uuid.Nil, "confirmed", "")
	if !errors.Is(err, ErrInvalidValue) {
		t.Fatalf("error = %v, want %v", err, ErrInvalidValue)
	}

	_, err = svc.UpdateAppointmentStatusByTeacher(context.Background(), uuid.New(), uuid.New(), "bad-status", "")
	if !errors.Is(err, ErrInvalidValue) {
		t.Fatalf("error = %v, want %v", err, ErrInvalidValue)
	}

	_, err = svc.UpdateAppointmentStatusByTeacher(context.Background(), uuid.New(), uuid.New(), "cancelled", "   ")
	if !errors.Is(err, ErrInvalidValue) {
		t.Fatalf("error = %v, want %v", err, ErrInvalidValue)
	}
}

func TestCancelAppointmentByParentValidation(t *testing.T) {
	svc := &AppointmentService{}

	_, err := svc.CancelAppointmentByParent(context.Background(), uuid.Nil, uuid.New(), "")
	if !errors.Is(err, ErrInvalidUserID) {
		t.Fatalf("error = %v, want %v", err, ErrInvalidUserID)
	}

	_, err = svc.CancelAppointmentByParent(context.Background(), uuid.New(), uuid.Nil, "")
	if !errors.Is(err, ErrInvalidValue) {
		t.Fatalf("error = %v, want %v", err, ErrInvalidValue)
	}
}

func TestCountParentUpcomingAppointmentsValidation(t *testing.T) {
	svc := &AppointmentService{}

	_, err := svc.CountParentUpcomingAppointments(context.Background(), uuid.Nil)
	if !errors.Is(err, ErrInvalidUserID) {
		t.Fatalf("error = %v, want %v", err, ErrInvalidUserID)
	}
}

func TestNormalizeListLimit(t *testing.T) {
	tests := []struct {
		name  string
		limit int
		want  int
	}{
		{name: "default for non-positive", limit: 0, want: 20},
		{name: "clamp to max", limit: 101, want: 100},
		{name: "keep valid value", limit: 50, want: 50},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := normalizeListLimit(tt.limit)
			if got != tt.want {
				t.Fatalf("normalizeListLimit(%d) = %d, want %d", tt.limit, got, tt.want)
			}
		})
	}
}

func TestIsValidAppointmentStatus(t *testing.T) {
	valid := []string{"pending", "confirmed", "cancelled", "completed", "no_show"}
	for _, status := range valid {
		if !isValidAppointmentStatus(status) {
			t.Fatalf("status %q must be valid", status)
		}
	}

	invalid := []string{"", "PENDING", "done", "scheduled"}
	for _, status := range invalid {
		if isValidAppointmentStatus(status) {
			t.Fatalf("status %q must be invalid", status)
		}
	}
}

func TestValidateTimeRange_TimezoneOffsets(t *testing.T) {
	fromRaw := "2026-04-07T12:00:00+07:00"
	toRaw := "2026-04-07T06:00:00+01:00"
	from, err := time.Parse(time.RFC3339, fromRaw)
	if err != nil {
		t.Fatalf("time.Parse(from) error = %v", err)
	}
	to, err := time.Parse(time.RFC3339, toRaw)
	if err != nil {
		t.Fatalf("time.Parse(to) error = %v", err)
	}

	if err := validateTimeRange(&from, &to); err != nil {
		t.Fatalf("validateTimeRange(equal instant with different offsets) error = %v", err)
	}

	invalidToRaw := "2026-04-07T05:59:59+01:00"
	invalidTo, err := time.Parse(time.RFC3339, invalidToRaw)
	if err != nil {
		t.Fatalf("time.Parse(invalidTo) error = %v", err)
	}
	if err := validateTimeRange(&from, &invalidTo); !errors.Is(err, ErrInvalidValue) {
		t.Fatalf("error = %v, want %v", err, ErrInvalidValue)
	}
}

func TestCreateAppointment_ErrorMappingAndTrimmedNote(t *testing.T) {
	parentUserID := uuid.New()
	studentID := uuid.New()
	slotID := uuid.New()

	t.Run("maps appointment slot unavailable", func(t *testing.T) {
		svc := &AppointmentService{appointmentRepo: &fakeAppointmentRepo{
			createAppointmentFn: func(context.Context, uuid.UUID, uuid.UUID, uuid.UUID, string) (model.Appointment, error) {
				return model.Appointment{}, repo.ErrAppointmentSlotUnavailable
			},
		}}

		_, err := svc.CreateAppointment(context.Background(), parentUserID, studentID, slotID, "note")
		if !errors.Is(err, ErrAppointmentSlotUnavailable) {
			t.Fatalf("error = %v, want %v", err, ErrAppointmentSlotUnavailable)
		}
	})

	t.Run("maps no rows updated to forbidden", func(t *testing.T) {
		svc := &AppointmentService{appointmentRepo: &fakeAppointmentRepo{
			createAppointmentFn: func(context.Context, uuid.UUID, uuid.UUID, uuid.UUID, string) (model.Appointment, error) {
				return model.Appointment{}, repo.ErrNoRowsUpdated
			},
		}}

		_, err := svc.CreateAppointment(context.Background(), parentUserID, studentID, slotID, "note")
		if !errors.Is(err, ErrForbidden) {
			t.Fatalf("error = %v, want %v", err, ErrForbidden)
		}
	})

	t.Run("wraps unexpected repository error", func(t *testing.T) {
		sentinel := errors.New("db down")
		svc := &AppointmentService{appointmentRepo: &fakeAppointmentRepo{
			createAppointmentFn: func(context.Context, uuid.UUID, uuid.UUID, uuid.UUID, string) (model.Appointment, error) {
				return model.Appointment{}, sentinel
			},
		}}

		_, err := svc.CreateAppointment(context.Background(), parentUserID, studentID, slotID, "note")
		if !errors.Is(err, sentinel) {
			t.Fatalf("error = %v, want wrapped sentinel", err)
		}
	})

	t.Run("trims note before forwarding", func(t *testing.T) {
		gotNote := ""
		svc := &AppointmentService{appointmentRepo: &fakeAppointmentRepo{
			createAppointmentFn: func(_ context.Context, _ uuid.UUID, _ uuid.UUID, _ uuid.UUID, note string) (model.Appointment, error) {
				gotNote = note
				return model.Appointment{AppointmentID: uuid.New()}, nil
			},
		}}

		_, err := svc.CreateAppointment(context.Background(), parentUserID, studentID, slotID, "  note with spaces  ")
		if err != nil {
			t.Fatalf("CreateAppointment() error = %v", err)
		}
		if gotNote != "note with spaces" {
			t.Fatalf("forwarded note = %q, want %q", gotNote, "note with spaces")
		}
	})
}

func TestCreateSlot_ErrorMappingAndTrimmedNote(t *testing.T) {
	teacherUserID := uuid.New()
	classID := uuid.New()
	start := time.Now().Add(30 * time.Minute)
	end := start.Add(30 * time.Minute)

	t.Run("maps no rows updated to forbidden", func(t *testing.T) {
		svc := &AppointmentService{appointmentRepo: &fakeAppointmentRepo{
			createSlotFn: func(context.Context, uuid.UUID, uuid.UUID, time.Time, time.Time, string) (model.AppointmentSlot, error) {
				return model.AppointmentSlot{}, repo.ErrNoRowsUpdated
			},
		}}

		_, err := svc.CreateSlot(context.Background(), teacherUserID, classID, start, end, "x")
		if !errors.Is(err, ErrForbidden) {
			t.Fatalf("error = %v, want %v", err, ErrForbidden)
		}
	})

	t.Run("trims note before forwarding", func(t *testing.T) {
		gotNote := ""
		svc := &AppointmentService{appointmentRepo: &fakeAppointmentRepo{
			createSlotFn: func(_ context.Context, _ uuid.UUID, _ uuid.UUID, _ time.Time, _ time.Time, note string) (model.AppointmentSlot, error) {
				gotNote = note
				return model.AppointmentSlot{SlotID: uuid.New()}, nil
			},
		}}

		_, err := svc.CreateSlot(context.Background(), teacherUserID, classID, start, end, "  slot note  ")
		if err != nil {
			t.Fatalf("CreateSlot() error = %v", err)
		}
		if gotNote != "slot note" {
			t.Fatalf("forwarded note = %q, want %q", gotNote, "slot note")
		}
	})
}

func TestUpdateAppointmentStatusByTeacher_ErrorMappingAndTrimmedReason(t *testing.T) {
	teacherUserID := uuid.New()
	appointmentID := uuid.New()

	t.Run("maps no rows updated to forbidden", func(t *testing.T) {
		svc := &AppointmentService{appointmentRepo: &fakeAppointmentRepo{
			updateAppointmentStatusByTeacherFn: func(context.Context, uuid.UUID, uuid.UUID, string, string) (model.Appointment, error) {
				return model.Appointment{}, repo.ErrNoRowsUpdated
			},
		}}

		_, err := svc.UpdateAppointmentStatusByTeacher(context.Background(), teacherUserID, appointmentID, "confirmed", "")
		if !errors.Is(err, ErrForbidden) {
			t.Fatalf("error = %v, want %v", err, ErrForbidden)
		}
	})

	t.Run("trims cancel reason before forwarding", func(t *testing.T) {
		gotReason := ""
		svc := &AppointmentService{appointmentRepo: &fakeAppointmentRepo{
			updateAppointmentStatusByTeacherFn: func(_ context.Context, _ uuid.UUID, _ uuid.UUID, _ string, cancelReason string) (model.Appointment, error) {
				gotReason = cancelReason
				return model.Appointment{AppointmentID: appointmentID}, nil
			},
		}}

		_, err := svc.UpdateAppointmentStatusByTeacher(context.Background(), teacherUserID, appointmentID, "cancelled", "  teacher reason  ")
		if err != nil {
			t.Fatalf("UpdateAppointmentStatusByTeacher() error = %v", err)
		}
		if gotReason != "teacher reason" {
			t.Fatalf("forwarded cancel_reason = %q, want %q", gotReason, "teacher reason")
		}
	})
}

func TestCancelAppointmentByParent_DefaultReasonAndErrorMapping(t *testing.T) {
	parentUserID := uuid.New()
	appointmentID := uuid.New()

	t.Run("maps no rows updated to forbidden", func(t *testing.T) {
		svc := &AppointmentService{appointmentRepo: &fakeAppointmentRepo{
			cancelAppointmentByParentFn: func(context.Context, uuid.UUID, uuid.UUID, string) (model.Appointment, error) {
				return model.Appointment{}, repo.ErrNoRowsUpdated
			},
		}}

		_, err := svc.CancelAppointmentByParent(context.Background(), parentUserID, appointmentID, "any")
		if !errors.Is(err, ErrForbidden) {
			t.Fatalf("error = %v, want %v", err, ErrForbidden)
		}
	})

	t.Run("uses default reason when blank", func(t *testing.T) {
		gotReason := ""
		svc := &AppointmentService{appointmentRepo: &fakeAppointmentRepo{
			cancelAppointmentByParentFn: func(_ context.Context, _ uuid.UUID, _ uuid.UUID, cancelReason string) (model.Appointment, error) {
				gotReason = cancelReason
				return model.Appointment{AppointmentID: appointmentID}, nil
			},
		}}

		_, err := svc.CancelAppointmentByParent(context.Background(), parentUserID, appointmentID, "   ")
		if err != nil {
			t.Fatalf("CancelAppointmentByParent() error = %v", err)
		}
		if gotReason != "parent_cancelled" {
			t.Fatalf("forwarded cancel_reason = %q, want %q", gotReason, "parent_cancelled")
		}
	})
}

func TestListTeacherAppointments_NormalizesLimitAndOffsetBeforeRepoCall(t *testing.T) {
	teacherUserID := uuid.New()
	capturedLimit := -1
	capturedOffset := -1

	svc := &AppointmentService{appointmentRepo: &fakeAppointmentRepo{
		listTeacherAppointmentsFn: func(_ context.Context, _ uuid.UUID, _ string, _ *time.Time, _ *time.Time, limit, offset int) ([]model.Appointment, int, error) {
			capturedLimit = limit
			capturedOffset = offset
			return []model.Appointment{}, 0, nil
		},
	}}

	_, _, err := svc.ListTeacherAppointments(context.Background(), teacherUserID, "", nil, nil, 1000, -9)
	if err != nil {
		t.Fatalf("ListTeacherAppointments() error = %v", err)
	}
	if capturedLimit != 100 {
		t.Fatalf("forwarded limit = %d, want %d", capturedLimit, 100)
	}
	if capturedOffset != 0 {
		t.Fatalf("forwarded offset = %d, want %d", capturedOffset, 0)
	}
}

func TestListParentAppointments_NormalizesLimitAndOffsetBeforeRepoCall(t *testing.T) {
	parentUserID := uuid.New()
	capturedLimit := -1
	capturedOffset := -1

	svc := &AppointmentService{appointmentRepo: &fakeAppointmentRepo{
		listParentAppointmentsFn: func(_ context.Context, _ uuid.UUID, _ string, _ *time.Time, _ *time.Time, limit, offset int) ([]model.Appointment, int, error) {
			capturedLimit = limit
			capturedOffset = offset
			return []model.Appointment{}, 0, nil
		},
	}}

	_, _, err := svc.ListParentAppointments(context.Background(), parentUserID, "", nil, nil, 1000, -9)
	if err != nil {
		t.Fatalf("ListParentAppointments() error = %v", err)
	}
	if capturedLimit != 100 {
		t.Fatalf("forwarded limit = %d, want %d", capturedLimit, 100)
	}
	if capturedOffset != 0 {
		t.Fatalf("forwarded offset = %d, want %d", capturedOffset, 0)
	}
}

func TestListAvailableSlotsForParent_NormalizesLimitAndOffsetBeforeRepoCall(t *testing.T) {
	parentUserID := uuid.New()
	studentID := uuid.New()
	capturedLimit := -1
	capturedOffset := -1

	svc := &AppointmentService{appointmentRepo: &fakeAppointmentRepo{
		listAvailableSlotsForParentFn: func(_ context.Context, gotParentUserID, gotStudentID uuid.UUID, _ *time.Time, _ *time.Time, limit, offset int) ([]model.AppointmentSlot, int, error) {
			if gotParentUserID != parentUserID || gotStudentID != studentID {
				t.Fatalf("unexpected ids forwarded")
			}
			capturedLimit = limit
			capturedOffset = offset
			return []model.AppointmentSlot{}, 0, nil
		},
	}}

	_, _, err := svc.ListAvailableSlotsForParent(context.Background(), parentUserID, studentID, nil, nil, 1000, -9)
	if err != nil {
		t.Fatalf("ListAvailableSlotsForParent() error = %v", err)
	}
	if capturedLimit != 100 {
		t.Fatalf("forwarded limit = %d, want %d", capturedLimit, 100)
	}
	if capturedOffset != 0 {
		t.Fatalf("forwarded offset = %d, want %d", capturedOffset, 0)
	}
}
