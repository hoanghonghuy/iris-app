package model

import (
	"time"

	"github.com/google/uuid"
)

type AuditLog struct {
	AuditLogID  uuid.UUID  `json:"audit_log_id"`
	ActorUserID uuid.UUID  `json:"actor_user_id"`
	ActorRole   string     `json:"actor_role,omitempty"`
	SchoolID    *uuid.UUID `json:"school_id,omitempty"`
	Action      string     `json:"action"`
	EntityType  string     `json:"entity_type"`
	EntityID    *uuid.UUID `json:"entity_id,omitempty"`
	Details     any        `json:"details"`
	CreatedAt   time.Time  `json:"created_at"`
}

type AuditLogFilter struct {
	Action      string
	EntityType  string
	ActorUserID *uuid.UUID
	SchoolID    *uuid.UUID
	From        *time.Time
	To          *time.Time
	Search      string
	Limit       int
	Offset      int
}

type AuditLogCreate struct {
	ActorUserID uuid.UUID
	ActorRole   string
	SchoolID    *uuid.UUID
	Action      string
	EntityType  string
	EntityID    *uuid.UUID
	Details     any
}

type ParentAnalytics struct {
	TotalChildren        int `json:"total_children"`
	UpcomingAppointments int `json:"upcoming_appointments"`
	RecentPosts7d        int `json:"recent_posts_7d"`
	RecentHealthAlerts7d int `json:"recent_health_alerts_7d"`
}
