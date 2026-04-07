package repo

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
)

type fakeAppointmentRow struct {
	err error
}

func (f *fakeAppointmentRow) Scan(dest ...any) error {
	if f.err != nil {
		return f.err
	}

	now := time.Now().UTC()
	appointmentID := uuid.New()
	slotID := uuid.New()
	parentID := uuid.New()
	studentID := uuid.New()
	confirmedAt := now.Add(5 * time.Minute)
	completedAt := now.Add(35 * time.Minute)
	cancelledAt := now.Add(65 * time.Minute)

	*(dest[0].(*uuid.UUID)) = appointmentID
	*(dest[1].(*uuid.UUID)) = slotID
	*(dest[2].(*uuid.UUID)) = parentID
	*(dest[3].(*uuid.UUID)) = studentID
	*(dest[4].(*string)) = "confirmed"
	*(dest[5].(*string)) = "parent note"
	*(dest[6].(*string)) = ""
	*(dest[7].(**time.Time)) = &confirmedAt
	*(dest[8].(**time.Time)) = &completedAt
	*(dest[9].(**time.Time)) = &cancelledAt
	*(dest[10].(*time.Time)) = now
	*(dest[11].(*time.Time)) = now.Add(10 * time.Minute)

	return nil
}

func TestScanAppointment(t *testing.T) {
	a, err := scanAppointment(&fakeAppointmentRow{})
	if err != nil {
		t.Fatalf("scanAppointment() error = %v", err)
	}

	if a.AppointmentID == uuid.Nil {
		t.Fatalf("expected non-empty appointment_id")
	}
	if a.SlotID == uuid.Nil {
		t.Fatalf("expected non-empty slot_id")
	}
	if a.ParentID == uuid.Nil || a.StudentID == uuid.Nil {
		t.Fatalf("expected non-empty parent/student IDs")
	}
	if a.Status != "confirmed" {
		t.Fatalf("status = %q, want %q", a.Status, "confirmed")
	}
	if a.Note != "parent note" {
		t.Fatalf("note = %q, want %q", a.Note, "parent note")
	}
	if a.ConfirmedAt == nil || a.CompletedAt == nil || a.CancelledAt == nil {
		t.Fatalf("expected confirmed/completed/cancelled pointers to be set")
	}
	if a.UpdatedAt.Before(a.CreatedAt) {
		t.Fatalf("expected updated_at >= created_at")
	}
}

func TestScanAppointmentPropagatesScanError(t *testing.T) {
	sentinel := errors.New("scan failed")
	_, err := scanAppointment(&fakeAppointmentRow{err: sentinel})
	if !errors.Is(err, sentinel) {
		t.Fatalf("error = %v, want %v", err, sentinel)
	}
}
