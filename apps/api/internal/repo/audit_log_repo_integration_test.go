package repo

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
)

type auditLogRepoFixture struct {
	schoolID         uuid.UUID
	otherSchoolID    uuid.UUID
	actorUserID      uuid.UUID
	otherActorUserID uuid.UUID
}

func mustNewAuditLogTestPool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	return mustNewAppointmentTestPool(t)
}

func seedAuditLogFixture(t *testing.T, pool *pgxpool.Pool) auditLogRepoFixture {
	t.Helper()
	ctx := context.Background()

	fx := auditLogRepoFixture{
		schoolID:         uuid.New(),
		otherSchoolID:    uuid.New(),
		actorUserID:      uuid.New(),
		otherActorUserID: uuid.New(),
	}

	actorEmail := fmt.Sprintf("audit-actor-%s@example.com", uuid.NewString())
	otherActorEmail := fmt.Sprintf("audit-actor-other-%s@example.com", uuid.NewString())

	if _, err := pool.Exec(ctx, `INSERT INTO schools (school_id, name) VALUES ($1, $2), ($3, $4)`, fx.schoolID, "Audit School A", fx.otherSchoolID, "Audit School B"); err != nil {
		t.Fatalf("insert schools error = %v", err)
	}
	if _, err := pool.Exec(ctx, `INSERT INTO users (user_id, email, password_hash) VALUES ($1, $2, $3), ($4, $5, $6)`,
		fx.actorUserID, actorEmail, "hash",
		fx.otherActorUserID, otherActorEmail, "hash",
	); err != nil {
		t.Fatalf("insert users error = %v", err)
	}

	t.Cleanup(func() {
		cleanupCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, _ = pool.Exec(cleanupCtx, `DELETE FROM users WHERE user_id = ANY($1::uuid[])`, []uuid.UUID{fx.actorUserID, fx.otherActorUserID})
		_, _ = pool.Exec(cleanupCtx, `DELETE FROM schools WHERE school_id = ANY($1::uuid[])`, []uuid.UUID{fx.schoolID, fx.otherSchoolID})
	})

	return fx
}

func insertAuditLogRecordAt(
	t *testing.T,
	pool *pgxpool.Pool,
	actorUserID uuid.UUID,
	schoolID *uuid.UUID,
	action, entityType string,
	detailsJSON string,
	createdAt time.Time,
) uuid.UUID {
	t.Helper()
	ctx := context.Background()

	auditLogID := uuid.New()
	if _, err := pool.Exec(ctx, `
		INSERT INTO audit_logs (audit_log_id, actor_user_id, actor_role, school_id, action, entity_type, details, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7::jsonb, $8)
	`, auditLogID, actorUserID, "SCHOOL_ADMIN", schoolID, action, entityType, detailsJSON, createdAt); err != nil {
		t.Fatalf("insert audit log record error = %v", err)
	}
	return auditLogID
}

func TestAuditLogRepoCreateAndList_FilterPagination(t *testing.T) {
	pool := mustNewAuditLogTestPool(t)
	fx := seedAuditLogFixture(t, pool)
	repo := NewAuditLogRepo(pool)
	ctx := context.Background()

	base := time.Now().UTC().Add(-2 * time.Hour).Truncate(time.Second)
	log1ID := insertAuditLogRecordAt(t, pool, fx.actorUserID, &fx.schoolID, "user.created", "user", `{"message":"created profile"}`, base.Add(1*time.Minute))
	log2ID := insertAuditLogRecordAt(t, pool, fx.otherActorUserID, &fx.schoolID, "user.updated", "user", `{"message":"updated profile"}`, base.Add(2*time.Minute))
	log3ID := insertAuditLogRecordAt(t, pool, fx.actorUserID, &fx.otherSchoolID, "post.deleted", "post", `{"message":"deleted post"}`, base.Add(3*time.Minute))

	items, total, err := repo.List(ctx, model.AuditLogFilter{Limit: 2, Offset: 0})
	if err != nil {
		t.Fatalf("List(all) error = %v", err)
	}
	if total != 3 || len(items) != 2 {
		t.Fatalf("all pagination mismatch: total=%d len=%d", total, len(items))
	}
	if items[0].AuditLogID != log3ID || items[1].AuditLogID != log2ID {
		t.Fatalf("all list order mismatch")
	}
	if detailsMap, ok := items[0].Details.(map[string]any); !ok || detailsMap["message"] != "deleted post" {
		t.Fatalf("expected decoded details map for newest record")
	}

	pagedItems, pagedTotal, err := repo.List(ctx, model.AuditLogFilter{Limit: 1, Offset: 1})
	if err != nil {
		t.Fatalf("List(offset) error = %v", err)
	}
	if pagedTotal != 3 || len(pagedItems) != 1 || pagedItems[0].AuditLogID != log2ID {
		t.Fatalf("offset pagination mismatch: total=%d len=%d", pagedTotal, len(pagedItems))
	}

	actionItems, actionTotal, err := repo.List(ctx, model.AuditLogFilter{Action: "user.updated", Limit: 10, Offset: 0})
	if err != nil {
		t.Fatalf("List(action) error = %v", err)
	}
	if actionTotal != 1 || len(actionItems) != 1 || actionItems[0].AuditLogID != log2ID {
		t.Fatalf("action filter mismatch: total=%d len=%d", actionTotal, len(actionItems))
	}

	entityItems, entityTotal, err := repo.List(ctx, model.AuditLogFilter{EntityType: "user", Limit: 10, Offset: 0})
	if err != nil {
		t.Fatalf("List(entity_type) error = %v", err)
	}
	if entityTotal != 2 || len(entityItems) != 2 {
		t.Fatalf("entity_type filter mismatch: total=%d len=%d", entityTotal, len(entityItems))
	}

	actorItems, actorTotal, err := repo.List(ctx, model.AuditLogFilter{ActorUserID: &fx.actorUserID, Limit: 10, Offset: 0})
	if err != nil {
		t.Fatalf("List(actor_user_id) error = %v", err)
	}
	if actorTotal != 2 || len(actorItems) != 2 {
		t.Fatalf("actor filter mismatch: total=%d len=%d", actorTotal, len(actorItems))
	}
	if actorItems[0].AuditLogID != log3ID || actorItems[1].AuditLogID != log1ID {
		t.Fatalf("actor filter order mismatch")
	}

	schoolItems, schoolTotal, err := repo.List(ctx, model.AuditLogFilter{SchoolID: &fx.schoolID, Limit: 10, Offset: 0})
	if err != nil {
		t.Fatalf("List(school_id) error = %v", err)
	}
	if schoolTotal != 2 || len(schoolItems) != 2 {
		t.Fatalf("school filter mismatch: total=%d len=%d", schoolTotal, len(schoolItems))
	}

	searchItems, searchTotal, err := repo.List(ctx, model.AuditLogFilter{Search: "profile", Limit: 10, Offset: 0})
	if err != nil {
		t.Fatalf("List(search) error = %v", err)
	}
	if searchTotal != 2 || len(searchItems) != 2 {
		t.Fatalf("search filter mismatch: total=%d len=%d", searchTotal, len(searchItems))
	}

	from := base.Add(2 * time.Minute)
	to := base.Add(3 * time.Minute)
	timeItems, timeTotal, err := repo.List(ctx, model.AuditLogFilter{From: &from, To: &to, Limit: 10, Offset: 0})
	if err != nil {
		t.Fatalf("List(time range) error = %v", err)
	}
	if timeTotal != 2 || len(timeItems) != 2 {
		t.Fatalf("time range filter mismatch: total=%d len=%d", timeTotal, len(timeItems))
	}
	if timeItems[0].AuditLogID != log3ID || timeItems[1].AuditLogID != log2ID {
		t.Fatalf("time range order mismatch")
	}
}

func TestAuditLogRepoCreate_PersistsJSONDetails(t *testing.T) {
	pool := mustNewAuditLogTestPool(t)
	fx := seedAuditLogFixture(t, pool)
	repo := NewAuditLogRepo(pool)
	ctx := context.Background()

	err := repo.Create(ctx, model.AuditLogCreate{
		ActorUserID: fx.actorUserID,
		ActorRole:   "SUPER_ADMIN",
		SchoolID:    &fx.schoolID,
		Action:      "school.created",
		EntityType:  "school",
		Details: map[string]any{
			"source": "integration-test",
			"ok":     true,
		},
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	items, total, err := repo.List(ctx, model.AuditLogFilter{Action: "school.created", Limit: 10, Offset: 0})
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if total != 1 || len(items) != 1 {
		t.Fatalf("expected one created audit log, got total=%d len=%d", total, len(items))
	}
	if items[0].SchoolID == nil || *items[0].SchoolID != fx.schoolID {
		t.Fatalf("expected school_id = %s", fx.schoolID)
	}
	detailsMap, ok := items[0].Details.(map[string]any)
	if !ok {
		t.Fatalf("expected details map, got %T", items[0].Details)
	}
	if detailsMap["source"] != "integration-test" {
		t.Fatalf("unexpected details.source = %v", detailsMap["source"])
	}
}
