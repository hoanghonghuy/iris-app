package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
)

type AuditLogRepo struct {
	pool *pgxpool.Pool
}

func NewAuditLogRepo(pool *pgxpool.Pool) *AuditLogRepo {
	return &AuditLogRepo{pool: pool}
}

func (r *AuditLogRepo) Create(ctx context.Context, in model.AuditLogCreate) error {
	detailsJSON, err := json.Marshal(in.Details)
	if err != nil {
		return err
	}

	const q = `
		INSERT INTO audit_logs (actor_user_id, actor_role, school_id, action, entity_type, entity_id, details)
		VALUES ($1, $2, $3, $4, $5, $6, $7::jsonb);
	`
	_, err = r.pool.Exec(ctx, q, in.ActorUserID, in.ActorRole, in.SchoolID, in.Action, in.EntityType, in.EntityID, detailsJSON)
	return err
}

func (r *AuditLogRepo) List(ctx context.Context, filter model.AuditLogFilter) ([]model.AuditLog, int, error) {
	q := `
		SELECT audit_log_id, actor_user_id, actor_role, school_id, action, entity_type, entity_id, details, created_at,
		       COUNT(*) OVER() AS total_count
		FROM audit_logs
		WHERE 1=1
	`

	args := make([]any, 0, 8)
	argPos := 1
	if filter.Action != "" {
		q += fmt.Sprintf(" AND action = $%d", argPos)
		args = append(args, filter.Action)
		argPos++
	}
	if filter.EntityType != "" {
		q += fmt.Sprintf(" AND entity_type = $%d", argPos)
		args = append(args, filter.EntityType)
		argPos++
	}
	if filter.ActorUserID != nil {
		q += fmt.Sprintf(" AND actor_user_id = $%d", argPos)
		args = append(args, *filter.ActorUserID)
		argPos++
	}
	if filter.SchoolID != nil {
		q += fmt.Sprintf(" AND school_id = $%d", argPos)
		args = append(args, *filter.SchoolID)
		argPos++
	}
	if filter.From != nil {
		q += fmt.Sprintf(" AND created_at >= $%d", argPos)
		args = append(args, *filter.From)
		argPos++
	}
	if filter.To != nil {
		q += fmt.Sprintf(" AND created_at <= $%d", argPos)
		args = append(args, *filter.To)
		argPos++
	}
	if filter.Search != "" {
		q += fmt.Sprintf(" AND (action ILIKE $%d OR entity_type ILIKE $%d OR CAST(details AS text) ILIKE $%d)", argPos, argPos, argPos)
		args = append(args, "%"+filter.Search+"%")
		argPos++
	}

	q += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argPos, argPos+1)
	args = append(args, filter.Limit, filter.Offset)

	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	items := make([]model.AuditLog, 0)
	total := 0
	for rows.Next() {
		var item model.AuditLog
		var schoolID *uuid.UUID
		var entityID *uuid.UUID
		var createdAt time.Time
		var detailsRaw []byte
		if err := rows.Scan(
			&item.AuditLogID,
			&item.ActorUserID,
			&item.ActorRole,
			&schoolID,
			&item.Action,
			&item.EntityType,
			&entityID,
			&detailsRaw,
			&createdAt,
			&total,
		); err != nil {
			return nil, 0, err
		}
		item.SchoolID = schoolID
		item.EntityID = entityID
		item.CreatedAt = createdAt
		if len(detailsRaw) > 0 {
			var details any
			if err := json.Unmarshal(detailsRaw, &details); err == nil {
				item.Details = details
			} else {
				item.Details = string(detailsRaw)
			}
		}
		items = append(items, item)
	}

	return items, total, rows.Err()
}
