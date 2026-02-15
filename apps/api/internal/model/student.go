package model

import (
	"time"

	"github.com/google/uuid"
)

type Student struct {
	StudentID      uuid.UUID `json:"student_id"`
	SchoolID       uuid.UUID `json:"school_id"`
	CurrentClassID uuid.UUID `json:"current_class_id"`
	FullName       string    `json:"full_name"`
	DOB            time.Time `json:"dob"`
	Gender         string    `json:"gender"`
}

// StudentParentCode thông tin parent code của student
type StudentParentCode struct {
	CodeID     uuid.UUID
	StudentID  uuid.UUID
	Code       string
	UsageCount int
	MaxUsage   int
	ExpiresAt  time.Time
}
