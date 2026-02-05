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
