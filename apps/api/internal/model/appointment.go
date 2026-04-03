package model

import (
	"time"

	"github.com/google/uuid"
)

type AppointmentSlot struct {
	SlotID      uuid.UUID `json:"slot_id"`
	TeacherID   uuid.UUID `json:"teacher_id"`
	TeacherName string    `json:"teacher_name,omitempty"`
	ClassID     uuid.UUID `json:"class_id"`
	ClassName   string    `json:"class_name,omitempty"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Note        string    `json:"note,omitempty"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Appointment struct {
	AppointmentID uuid.UUID  `json:"appointment_id"`
	SlotID        uuid.UUID  `json:"slot_id"`
	ParentID      uuid.UUID  `json:"parent_id"`
	ParentName    string     `json:"parent_name,omitempty"`
	StudentID     uuid.UUID  `json:"student_id"`
	StudentName   string     `json:"student_name,omitempty"`
	TeacherID     uuid.UUID  `json:"teacher_id,omitempty"`
	TeacherName   string     `json:"teacher_name,omitempty"`
	ClassID       uuid.UUID  `json:"class_id,omitempty"`
	ClassName     string     `json:"class_name,omitempty"`
	Status        string     `json:"status"`
	Note          string     `json:"note,omitempty"`
	CancelReason  string     `json:"cancel_reason,omitempty"`
	StartTime     time.Time  `json:"start_time"`
	EndTime       time.Time  `json:"end_time"`
	ConfirmedAt   *time.Time `json:"confirmed_at,omitempty"`
	CompletedAt   *time.Time `json:"completed_at,omitempty"`
	CancelledAt   *time.Time `json:"cancelled_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}
