package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/auth"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
)

type fakeAuditLogCreator struct {
	calls int
	last  model.AuditLogCreate
}

func (f *fakeAuditLogCreator) Create(_ context.Context, in model.AuditLogCreate) error {
	f.calls++
	f.last = in
	return nil
}

func TestAdminAuditLogger_StoresOnlyWhitelistedMetadata(t *testing.T) {
	gin.SetMode(gin.TestMode)

	fake := &fakeAuditLogCreator{}
	actorID := uuid.New()
	schoolID := uuid.New()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set(CtxClaims, &auth.Claims{
			UserID: actorID.String(),
			Roles:  []string{"SCHOOL_ADMIN"},
		})
		c.Set(CtxAdminSchoolID, schoolID.String())
		c.Next()
	})
	r.Use(AdminAuditLogger(fake))
	r.POST("/api/v1/admin/users", func(c *gin.Context) {
		c.Status(http.StatusCreated)
	})

	body := `{"password":"super-secret","token":"top-secret"}`
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/users?foo=bar", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}
	if fake.calls != 1 {
		t.Fatalf("expected one audit log call, got %d", fake.calls)
	}
	if fake.last.SchoolID == nil || *fake.last.SchoolID != schoolID {
		t.Fatalf("expected school_id %s to be propagated", schoolID.String())
	}

	details, ok := fake.last.Details.(map[string]any)
	if !ok {
		t.Fatalf("expected details as map, got %T", fake.last.Details)
	}
	if _, exists := details["body"]; exists {
		t.Fatalf("expected raw request body not to be logged")
	}
	if details["school_id"] != schoolID.String() {
		t.Fatalf("expected details.school_id to equal %s, got %v", schoolID.String(), details["school_id"])
	}
}

func TestAdminAuditLogger_SkipsReadOnlyMethods(t *testing.T) {
	gin.SetMode(gin.TestMode)

	fake := &fakeAuditLogCreator{}
	r := gin.New()
	r.Use(AdminAuditLogger(fake))
	r.GET("/api/v1/admin/ping", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/ping", nil)
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	if fake.calls != 0 {
		t.Fatalf("expected no audit log call for GET request, got %d", fake.calls)
	}
}
