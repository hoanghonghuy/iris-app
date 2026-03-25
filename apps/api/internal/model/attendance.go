package model

import (
	"time"

	"github.com/google/uuid"
)

type AttendanceRecord struct {
	AttendanceID uuid.UUID  `json:"attendance_id"`
	StudentID    uuid.UUID  `json:"student_id"`
	Date         time.Time  `json:"date"`
	Status       string     `json:"status"` // present, absent, late, excused
	CheckInAt    *time.Time `json:"check_in_at,omitempty"`
	CheckOutAt   *time.Time `json:"check_out_at,omitempty"`
	Note         string     `json:"note,omitempty"`
	RecordedBy   uuid.UUID  `json:"recorded_by"` // user_id of teacher who recorded
}

type AttendanceChangeLog struct {
	ChangeID      uuid.UUID  `json:"change_id"`
	AttendanceID  uuid.UUID  `json:"attendance_id"`
	StudentID     uuid.UUID  `json:"student_id"`
	Date          time.Time  `json:"date"`
	ChangeType    string     `json:"change_type"` // create|update
	OldStatus     *string    `json:"old_status,omitempty"`
	NewStatus     string     `json:"new_status"`
	OldNote       *string    `json:"old_note,omitempty"`
	NewNote       *string    `json:"new_note,omitempty"`
	ChangedBy     uuid.UUID  `json:"changed_by"`
	ChangedAt     time.Time  `json:"changed_at"`
}
