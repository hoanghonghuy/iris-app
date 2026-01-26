package model

import (
	"time"

	"github.com/google/uuid"
)

type HealthLog struct {
	HealthLogID uuid.UUID `json:"health_log_id"`
	StudentID   uuid.UUID `json:"student_id"`
	RecordedAt  time.Time `json:"recorded_at"`
	Temperature *float64  `json:"temperature,omitempty"`
	Symptoms    string    `json:"symptoms,omitempty"`
	Severity    *string   `json:"severity,omitempty"` // normal|watch|urgent
	Note        string    `json:"note,omitempty"`
	RecordedBy  uuid.UUID `json:"recorded_by"`
}
