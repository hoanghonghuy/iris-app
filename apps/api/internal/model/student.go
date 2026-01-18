package model

import (
	"time"

	"github.com/google/uuid"
)

type Student struct {
	ID             uuid.UUID `json:"student_id"`
	SchoolID       uuid.UUID `json:"school_id"`
	CurrentClassID uuid.UUID `json:"current_class_id"`
	FullName       string    `json:"full_name"`
	DOB            time.Time `json:"dob"`
	Gender         string    `json:"gender"`
}
