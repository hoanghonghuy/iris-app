package model

import "github.com/google/uuid"

type Class struct {
	ID         uuid.UUID `json:"class_id"`
	SchoolID   uuid.UUID `json:"school_id"`
	Name       string    `json:"name"`
	SchoolYear string    `json:"school_year"`
}
