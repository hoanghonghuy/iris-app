package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/response"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/service"
)

type AuditLogHandler struct {
	auditLogService *service.AuditLogService
}

func NewAuditLogHandler(auditLogService *service.AuditLogService) *AuditLogHandler {
	return &AuditLogHandler{auditLogService: auditLogService}
}

func (h *AuditLogHandler) List(c *gin.Context) {
	from, to, err := h.auditLogService.ParseTimeRange(c.Query("from"), c.Query("to"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	var actorUserID *uuid.UUID
	if actorRaw := c.Query("actor_user_id"); actorRaw != "" {
		id, err := uuid.Parse(actorRaw)
		if err != nil {
			response.Fail(c, http.StatusBadRequest, "invalid actor_user_id")
			return
		}
		actorUserID = &id
	}

	limit := 20
	if raw := c.Query("limit"); raw != "" {
		if n, err := strconv.Atoi(raw); err == nil {
			limit = n
		}
	}
	offset := 0
	if raw := c.Query("offset"); raw != "" {
		if n, err := strconv.Atoi(raw); err == nil {
			offset = n
		}
	}

	items, total, err := h.auditLogService.List(c.Request.Context(), model.AuditLogFilter{
		Action:      c.Query("action"),
		EntityType:  c.Query("entity_type"),
		ActorUserID: actorUserID,
		From:        from,
		To:          to,
		Search:      c.Query("q"),
		Limit:       limit,
		Offset:      offset,
	})
	if err != nil {
		if errors.Is(err, service.ErrInvalidValue) {
			response.Fail(c, http.StatusBadRequest, err.Error())
			return
		}
		response.Fail(c, http.StatusInternalServerError, "failed to fetch audit logs")
		return
	}

	response.OKPaginated(c, items, response.Pagination{
		Total:   total,
		Limit:   limit,
		Offset:  offset,
		HasMore: offset+len(items) < total,
	})
}
