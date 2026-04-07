package auditloghandlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/auth"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/middleware"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/service"
)

type fakeAuditLogQueryService struct {
	parseTimeRangeFn func(string, string) (*time.Time, *time.Time, error)
	listFn           func(context.Context, model.AuditLogFilter) ([]model.AuditLog, int, error)
}

func (f *fakeAuditLogQueryService) ParseTimeRange(fromRaw, toRaw string) (*time.Time, *time.Time, error) {
	if f.parseTimeRangeFn == nil {
		return nil, nil, errors.New("unexpected ParseTimeRange call")
	}
	return f.parseTimeRangeFn(fromRaw, toRaw)
}

func (f *fakeAuditLogQueryService) List(ctx context.Context, filter model.AuditLogFilter) ([]model.AuditLog, int, error) {
	if f.listFn == nil {
		return nil, 0, errors.New("unexpected List call")
	}
	return f.listFn(ctx, filter)
}

func decodeAuditLogError(t *testing.T, rec *httptest.ResponseRecorder) string {
	t.Helper()
	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	return body["error"]
}

func TestList_RejectsSchoolAdminWithoutSchoolScope(t *testing.T) {
	gin.SetMode(gin.TestMode)

	h := &AuditLogHandler{auditLogService: &fakeAuditLogQueryService{}}
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set(middleware.CtxClaims, &auth.Claims{
			UserID: "00000000-0000-0000-0000-000000000001",
			Roles:  []string{"SCHOOL_ADMIN"},
		})
		c.Next()
	})
	r.GET("/audit-logs", h.List)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/audit-logs", nil)
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d", http.StatusForbidden, rec.Code)
	}
}

func TestList_RejectsInvalidActorUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	h := &AuditLogHandler{auditLogService: &fakeAuditLogQueryService{
		parseTimeRangeFn: func(string, string) (*time.Time, *time.Time, error) {
			return nil, nil, nil
		},
	}}
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set(middleware.CtxClaims, &auth.Claims{
			UserID: "00000000-0000-0000-0000-000000000001",
			Roles:  []string{"SUPER_ADMIN"},
		})
		c.Next()
	})
	r.GET("/audit-logs", h.List)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/audit-logs?actor_user_id=not-a-uuid", nil)
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestList_RejectsInvalidFromQuery(t *testing.T) {
	gin.SetMode(gin.TestMode)

	h := &AuditLogHandler{auditLogService: &fakeAuditLogQueryService{
		parseTimeRangeFn: func(string, string) (*time.Time, *time.Time, error) {
			return nil, nil, service.ErrInvalidDate
		},
	}}
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set(middleware.CtxClaims, &auth.Claims{
			UserID: "00000000-0000-0000-0000-000000000001",
			Roles:  []string{"SUPER_ADMIN"},
		})
		c.Next()
	})
	r.GET("/audit-logs", h.List)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/audit-logs?from=invalid-rfc3339", nil)
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestList_RejectsInvalidToQuery(t *testing.T) {
	gin.SetMode(gin.TestMode)

	h := &AuditLogHandler{auditLogService: &fakeAuditLogQueryService{
		parseTimeRangeFn: func(string, string) (*time.Time, *time.Time, error) {
			return nil, nil, service.ErrInvalidDate
		},
	}}
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set(middleware.CtxClaims, &auth.Claims{
			UserID: "00000000-0000-0000-0000-000000000001",
			Roles:  []string{"SUPER_ADMIN"},
		})
		c.Next()
	})
	r.GET("/audit-logs", h.List)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/audit-logs?to=invalid-rfc3339", nil)
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestList_SuccessAndErrorMapping(t *testing.T) {
	gin.SetMode(gin.TestMode)
	actorID := uuid.New()
	schoolID := uuid.New()

	t.Run("success", func(t *testing.T) {
		h := &AuditLogHandler{auditLogService: &fakeAuditLogQueryService{
			parseTimeRangeFn: func(string, string) (*time.Time, *time.Time, error) {
				return nil, nil, nil
			},
			listFn: func(_ context.Context, filter model.AuditLogFilter) ([]model.AuditLog, int, error) {
				if filter.Limit != 100 || filter.Offset != 0 {
					t.Fatalf("expected normalized limit/offset, got %d/%d", filter.Limit, filter.Offset)
				}
				if filter.SchoolID == nil || *filter.SchoolID != schoolID {
					t.Fatalf("expected school scope to be forwarded")
				}
				return []model.AuditLog{{AuditLogID: uuid.New(), ActorUserID: actorID, Action: "x", EntityType: "user"}}, 1, nil
			},
		}}

		r := gin.New()
		r.Use(func(c *gin.Context) {
			c.Set(middleware.CtxClaims, &auth.Claims{UserID: actorID.String(), Roles: []string{"SCHOOL_ADMIN"}})
			c.Set(middleware.CtxAdminSchoolID, schoolID.String())
			c.Next()
		})
		r.GET("/audit-logs", h.List)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/audit-logs?limit=999&offset=-2", nil)
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
		}
	})

	tests := []struct {
		name       string
		serviceErr error
		wantStatus int
	}{
		{name: "bad request", serviceErr: service.ErrInvalidValue, wantStatus: http.StatusBadRequest},
		{name: "internal", serviceErr: errors.New("boom"), wantStatus: http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &AuditLogHandler{auditLogService: &fakeAuditLogQueryService{
				parseTimeRangeFn: func(string, string) (*time.Time, *time.Time, error) {
					return nil, nil, nil
				},
				listFn: func(context.Context, model.AuditLogFilter) ([]model.AuditLog, int, error) {
					return nil, 0, tt.serviceErr
				},
			}}

			r := gin.New()
			r.Use(func(c *gin.Context) {
				c.Set(middleware.CtxClaims, &auth.Claims{UserID: actorID.String(), Roles: []string{"SUPER_ADMIN"}})
				c.Next()
			})
			r.GET("/audit-logs", h.List)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/audit-logs", nil)
			r.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
		})
	}
}

func TestList_PaginationFallbackAndNormalizationForwarding(t *testing.T) {
	gin.SetMode(gin.TestMode)
	actorID := uuid.New()

	tests := []struct {
		name       string
		query      string
		wantLimit  int
		wantOffset int
	}{
		{name: "fallback invalid query", query: "?limit=abc&offset=xyz", wantLimit: 20, wantOffset: 0},
		{name: "default when non-positive", query: "?limit=0&offset=-9", wantLimit: 20, wantOffset: 0},
		{name: "clamp max limit", query: "?limit=999&offset=7", wantLimit: 100, wantOffset: 7},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &AuditLogHandler{auditLogService: &fakeAuditLogQueryService{
				parseTimeRangeFn: func(string, string) (*time.Time, *time.Time, error) {
					return nil, nil, nil
				},
				listFn: func(_ context.Context, filter model.AuditLogFilter) ([]model.AuditLog, int, error) {
					if filter.Limit != tt.wantLimit || filter.Offset != tt.wantOffset {
						t.Fatalf("forwarded limit/offset = %d/%d, want %d/%d", filter.Limit, filter.Offset, tt.wantLimit, tt.wantOffset)
					}
					return []model.AuditLog{}, 0, nil
				},
			}}

			r := gin.New()
			r.Use(func(c *gin.Context) {
				c.Set(middleware.CtxClaims, &auth.Claims{UserID: actorID.String(), Roles: []string{"SUPER_ADMIN"}})
				c.Next()
			})
			r.GET("/audit-logs", h.List)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/audit-logs"+tt.query, nil)
			r.ServeHTTP(rec, req)

			if rec.Code != http.StatusOK {
				t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
			}
		})
	}
}

func TestList_RequiresClaims(t *testing.T) {
	gin.SetMode(gin.TestMode)

	h := &AuditLogHandler{auditLogService: &fakeAuditLogQueryService{}}
	r := gin.New()
	r.GET("/audit-logs", h.List)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/audit-logs", nil)
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
	if got := decodeAuditLogError(t, rec); got != "unauthorized" {
		t.Fatalf("error = %q, want %q", got, "unauthorized")
	}
}

func TestList_TimezoneOffsetRangeForwarding(t *testing.T) {
	gin.SetMode(gin.TestMode)
	actorID := uuid.New()
	fromRaw := "2026-04-07T12:00:00+07:00"
	toRaw := "2026-04-07T06:00:00+01:00"
	expectedFrom, _ := time.Parse(time.RFC3339, fromRaw)
	expectedTo, _ := time.Parse(time.RFC3339, toRaw)

	h := &AuditLogHandler{auditLogService: &fakeAuditLogQueryService{
		parseTimeRangeFn: func(gotFromRaw, gotToRaw string) (*time.Time, *time.Time, error) {
			if gotFromRaw != fromRaw || gotToRaw != toRaw {
				t.Fatalf("unexpected raw range: from=%q to=%q", gotFromRaw, gotToRaw)
			}
			return &expectedFrom, &expectedTo, nil
		},
		listFn: func(_ context.Context, filter model.AuditLogFilter) ([]model.AuditLog, int, error) {
			if filter.From == nil || filter.To == nil {
				t.Fatalf("expected non-nil from/to filter")
			}
			if !filter.From.Equal(expectedFrom) || !filter.To.Equal(expectedTo) {
				t.Fatalf("timezone forwarding mismatch: from=%v to=%v", *filter.From, *filter.To)
			}
			return []model.AuditLog{}, 0, nil
		},
	}}

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set(middleware.CtxClaims, &auth.Claims{UserID: actorID.String(), Roles: []string{"SUPER_ADMIN"}})
		c.Next()
	})
	r.GET("/audit-logs", h.List)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/audit-logs?from="+url.QueryEscape(fromRaw)+"&to="+url.QueryEscape(toRaw), nil)
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
}
