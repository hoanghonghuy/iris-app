package middleware

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/auth"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/service"
)

func AdminAuditLogger(auditLogService *service.AuditLogService) gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		if method != http.MethodPost && method != http.MethodPut && method != http.MethodPatch && method != http.MethodDelete {
			c.Next()
			return
		}

		var bodyPreview string
		if c.Request.Body != nil {
			rawBody, _ := io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(rawBody))
			if len(rawBody) > 2048 {
				bodyPreview = string(rawBody[:2048])
			} else {
				bodyPreview = string(rawBody)
			}
		}

		path := c.FullPath()
		entityType := "admin"
		if path != "" {
			parts := strings.Split(strings.Trim(path, "/"), "/")
			if len(parts) >= 3 {
				entityType = parts[2]
			}
		}

		c.Next()

		status := c.Writer.Status()
		if status >= 500 {
			return
		}

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

		entityID := extractEntityIDFromPath(c.Param("school_id"), c.Param("class_id"), c.Param("student_id"), c.Param("teacher_id"), c.Param("parent_id"), c.Param("user_id"), c.Param("admin_id"))

		actorRole := ""
		if len(claims.Roles) > 0 {
			actorRole = claims.Roles[0]
		}

		ctx, cancel := contextWithShortTimeout()
		defer cancel()

		_ = auditLogService.Create(ctx, model.AuditLogCreate{
			ActorUserID: actorID,
			ActorRole:   actorRole,
			Action:      method + " " + path,
			EntityType:  entityType,
			EntityID:    entityID,
			Details: map[string]any{
				"status":       status,
				"request_path": c.Request.URL.Path,
				"query":        c.Request.URL.RawQuery,
				"body":         bodyPreview,
				"at":           time.Now().UTC().Format(time.RFC3339),
			},
		})
	}
}

func extractEntityIDFromPath(values ...string) *uuid.UUID {
	for _, v := range values {
		if v == "" {
			continue
		}
		id, err := uuid.Parse(v)
		if err == nil {
			return &id
		}
	}
	return nil
}

func contextWithShortTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 1200*time.Millisecond)
}
