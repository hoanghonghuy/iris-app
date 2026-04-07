package auditloghandlers

import (
	"context"
	"time"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
)

type auditLogQueryService interface {
	ParseTimeRange(fromRaw, toRaw string) (*time.Time, *time.Time, error)
	List(ctx context.Context, filter model.AuditLogFilter) ([]model.AuditLog, int, error)
}

type AuditLogHandler struct {
	auditLogService auditLogQueryService
}

func NewAuditLogHandler(auditLogService auditLogQueryService) *AuditLogHandler {
	return &AuditLogHandler{auditLogService: auditLogService}
}
