package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/repo"
)

type AuditLogService struct {
	auditRepo *repo.AuditLogRepo
}

func NewAuditLogService(auditRepo *repo.AuditLogRepo) *AuditLogService {
	return &AuditLogService{auditRepo: auditRepo}
}

func (s *AuditLogService) Create(ctx context.Context, in model.AuditLogCreate) error {
	if in.ActorUserID == uuid.Nil {
		return ErrInvalidUserID
	}
	if in.Action == "" || in.EntityType == "" {
		return fmt.Errorf("%w: action and entity_type are required", ErrInvalidValue)
	}
	return s.auditRepo.Create(ctx, in)
}

func (s *AuditLogService) List(ctx context.Context, filter model.AuditLogFilter) ([]model.AuditLog, int, error) {
	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	if filter.Limit > 100 {
		filter.Limit = 100
	}
	if filter.Offset < 0 {
		filter.Offset = 0
	}
	return s.auditRepo.List(ctx, filter)
}

func (s *AuditLogService) ParseTimeRange(fromRaw, toRaw string) (*time.Time, *time.Time, error) {
	var fromPtr *time.Time
	var toPtr *time.Time
	if fromRaw != "" {
		from, err := time.Parse(time.RFC3339, fromRaw)
		if err != nil {
			return nil, nil, fmt.Errorf("%w: invalid from time", ErrInvalidDate)
		}
		fromPtr = &from
	}
	if toRaw != "" {
		to, err := time.Parse(time.RFC3339, toRaw)
		if err != nil {
			return nil, nil, fmt.Errorf("%w: invalid to time", ErrInvalidDate)
		}
		toPtr = &to
	}
	return fromPtr, toPtr, nil
}
