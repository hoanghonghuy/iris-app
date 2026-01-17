package model

import "github.com/google/uuid"

type School struct {
	ID      uuid.UUID `json:"school_id"`
	Name    string    `json:"name"`
	Address string    `json:"address"`
}
