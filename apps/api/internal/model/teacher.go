package model

import "github.com/google/uuid"

// ClassInfo chứa thông tin cơ bản của lớp hiển thị cùng giáo viên
type ClassInfo struct {
	ClassID uuid.UUID `json:"class_id"`
	Name    string    `json:"name"`
}

type Teacher struct {
	TeacherID uuid.UUID    `json:"teacher_id"`
	UserID    uuid.UUID    `json:"user_id"`
	Email     string       `json:"email"`
	FullName  string       `json:"full_name"`
	Phone     string       `json:"phone"`
	SchoolID  uuid.UUID    `json:"school_id"`
	Classes   []ClassInfo `json:"classes,omitempty"`
}
