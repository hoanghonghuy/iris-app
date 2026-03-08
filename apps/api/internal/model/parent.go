package model

import "github.com/google/uuid"

// ChildInfo chứa thông tin cơ bản của học sinh (con) hiển thị cùng phụ huynh
type ChildInfo struct {
	StudentID uuid.UUID `json:"student_id"`
	FullName  string    `json:"full_name"`
}

type Parent struct {
	ParentID uuid.UUID    `json:"parent_id"`
	UserID   uuid.UUID    `json:"user_id"`
	Email    string       `json:"email"`
	FullName string       `json:"full_name"`
	Phone    string       `json:"phone"`
	SchoolID uuid.UUID    `json:"school_id"`
	Children []ChildInfo `json:"children,omitempty"`
}
