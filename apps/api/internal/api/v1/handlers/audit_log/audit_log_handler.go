package auditloghandlers

import "github.com/hoanghonghuy/iris-app/apps/api/internal/service"

type AuditLogHandler struct {
	auditLogService *service.AuditLogService
}

func NewAuditLogHandler(auditLogService *service.AuditLogService) *AuditLogHandler {
	return &AuditLogHandler{auditLogService: auditLogService}
}
