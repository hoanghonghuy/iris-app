package model

import "github.com/google/uuid"

type Teacher struct {
	TeacherID uuid.UUID `json:"teacher_id"`
	UserID    uuid.UUID `json:"user_id"`
	Email     string    `json:"email"`
	FullName  string    `json:"full_name"`
	Phone     string    `json:"phone"`
	SchoolID  uuid.UUID `json:"school_id"`
}
