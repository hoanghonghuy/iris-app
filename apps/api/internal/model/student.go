package model

import (
	"time"

	"github.com/google/uuid"
)

type Student struct {
	StudentID        uuid.UUID `json:"student_id"`
	SchoolID         uuid.UUID `json:"school_id"`
	CurrentClassID   uuid.UUID `json:"current_class_id"`
	CurrentClassName *string   `json:"current_class_name,omitempty"`
	FullName         string    `json:"full_name"`
	DOB              time.Time `json:"dob"`
	Gender           string    `json:"gender"`

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

// StudentParentInfo thông tin cha mẹ của học sinh (rút gọn)
type StudentParentInfo struct {
	ParentID uuid.UUID `json:"parent_id"`
	FullName string    `json:"full_name"`
	Phone    string    `json:"phone"`
	Email    string    `json:"email"`
}

// StudentProfile Dữ liệu chi tiết về học sinh bao gồm phụ huynh
type StudentProfile struct {
	Student
	Parents []StudentParentInfo `json:"parents"`
}
