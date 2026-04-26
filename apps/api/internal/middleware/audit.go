package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/auth"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
)

// AuditLogger logs all write operations (POST/PUT/PATCH/DELETE) for compliance and security
// Follows best practices:
// - Only logs metadata (no request/response body)
// - Async logging (doesn't block response)
// - Only logs successful operations (status < 500)
// - Skips health checks and websocket endpoints
func AuditLogger(auditLogService auditLogCreator) gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		// Only log write operations
		if method != http.MethodPost && method != http.MethodPut &&
			method != http.MethodPatch && method != http.MethodDelete {
			c.Next()
			return
		}

		path := c.Request.URL.Path

		// Skip health checks, websocket, and public endpoints
		if strings.Contains(path, "/health") ||
			strings.Contains(path, "/ws") ||
			strings.Contains(path, "/auth/login") ||
			strings.Contains(path, "/auth/forgot-password") ||
			strings.Contains(path, "/auth/reset-password") {
			c.Next()
			return
		}

		// Capture start time for duration tracking
		startTime := time.Now()

		// Process request first
		c.Next()

		// Only log successful operations (skip server errors)
		status := c.Writer.Status()
		if status >= 500 {
			return
		}

		// Extract actor information from JWT claims
		claimsAny, ok := c.Get(CtxClaims)
		if !ok {
			return
		}
		claims, ok := claimsAny.(*auth.Claims)
		if !ok {
			return
		}
		actorID, err := uuid.Parse(claims.UserID)
		if err != nil {
			return
		}

		// Determine actor role
		actorRole := ""
		if len(claims.Roles) > 0 {
			actorRole = claims.Roles[0]
		}

		// Extract entity type from path
		entityType := extractEntityType(path)

		// Extract entity ID from path parameters
		entityID := extractEntityIDFromPath(
			c.Param("school_id"),
			c.Param("class_id"),
			c.Param("student_id"),
			c.Param("teacher_id"),
			c.Param("parent_id"),
			c.Param("user_id"),
			c.Param("admin_id"),
			c.Param("post_id"),
			c.Param("comment_id"),
			c.Param("conversation_id"),
			c.Param("appointment_id"),
		)

		// Extract school_id if available (for admin context)
		var schoolID *uuid.UUID
		if schoolIDAny, exists := c.Get(CtxAdminSchoolID); exists {
			if schoolIDStr, ok := schoolIDAny.(string); ok && schoolIDStr != "" {
				if parsedSchoolID, parseErr := uuid.Parse(schoolIDStr); parseErr == nil {
					schoolID = &parsedSchoolID
				}
			}
		}

		// Build action string (e.g., "POST /teacher/posts", "DELETE /admin/schools/:school_id")
		action := method + " " + c.FullPath()
		if action == method+" " { // fallback if FullPath is empty
			action = method + " " + path
		}

		// Build details map (metadata only, no sensitive data)
		details := map[string]any{
			"status":      status,
			"method":      method,
			"path":        path,
			"duration_ms": time.Since(startTime).Milliseconds(),
			"user_agent":  c.GetHeader("User-Agent"),
			"ip":          c.ClientIP(),
			"timestamp":   time.Now().UTC().Format(time.RFC3339),
		}

		// Add query params if present (but redact sensitive params)
		if c.Request.URL.RawQuery != "" {
			details["query"] = redactSensitiveQuery(c.Request.URL.RawQuery)
		}

		if schoolID != nil {
			details["school_id"] = schoolID.String()
		}

		// Async logging - doesn't block response
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			_ = auditLogService.Create(ctx, model.AuditLogCreate{
				ActorUserID: actorID,
				ActorRole:   actorRole,
				SchoolID:    schoolID,
				Action:      action,
				EntityType:  entityType,
				EntityID:    entityID,
				Details:     details,
			})
		}()
	}
}

// extractEntityType determines entity type from URL path
func extractEntityType(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")

	// Skip /api/v1 prefix
	if len(parts) >= 2 && parts[0] == "api" && parts[1] == "v1" {
		parts = parts[2:]
	}

	// Extract entity type based on path structure
	// Examples:
	// /teacher/posts -> "post"
	// /admin/schools -> "school"
	// /parent/children/:id/posts -> "post"
	for i := len(parts) - 1; i >= 0; i-- {
		part := parts[i]
		// Skip path parameters (starts with :)
		if strings.HasPrefix(part, ":") {
			continue
		}
		// Skip common action words
		if part == "by-school" || part == "by-class" || part == "search" {
			continue
		}
		// Return the entity type (singular form)
		return strings.TrimSuffix(part, "s")
	}

	return "unknown"
}

// redactSensitiveQuery removes sensitive parameters from query string
func redactSensitiveQuery(query string) string {
	sensitiveParams := []string{"password", "token", "secret", "key", "api_key"}
	for _, param := range sensitiveParams {
		if strings.Contains(strings.ToLower(query), param) {
			return "[REDACTED]"
		}
	}
	return query
}
