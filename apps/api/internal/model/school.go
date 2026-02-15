package model

import "github.com/google/uuid"

type School struct {
	SchoolID uuid.UUID `json:"school_id"`
	Name     string    `json:"name"`
	Address  string    `json:"address"`
}
