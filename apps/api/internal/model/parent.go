package model

import "github.com/google/uuid"

type Parent struct {
	ParentID uuid.UUID `json:"parent_id"`
	UserID   uuid.UUID `json:"user_id"`
	Email    string    `json:"email"`
	FullName string    `json:"full_name"`
	Phone    string    `json:"phone"`
	SchoolID uuid.UUID `json:"school_id"`
}
