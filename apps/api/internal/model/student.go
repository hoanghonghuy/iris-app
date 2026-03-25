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
	
	// Active parent code information (added from join)
	ActiveParentCode *string    `json:"active_parent_code,omitempty"`
	CodeExpiresAt    *time.Time `json:"code_expires_at,omitempty"`
	CodeUsageCount   *int       `json:"code_usage_count,omitempty"`
	CodeMaxUsage     *int       `json:"code_max_usage,omitempty"`
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
