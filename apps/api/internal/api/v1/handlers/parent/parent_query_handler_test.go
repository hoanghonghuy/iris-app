package parenthandlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/service"
)

func TestList_SuccessAndErrorMapping(t *testing.T) {
	gin.SetMode(gin.TestMode)
	schoolID := uuid.New()

	t.Run("invalid pagination", func(t *testing.T) {
		h := &ParentHandler{}
		r := gin.New()
		r.GET("/parents", h.List)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/parents?limit=0&offset=-1", nil)
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
	})

	t.Run("success", func(t *testing.T) {
		h := &ParentHandler{parentService: &fakeParentService{listFn: func(_ context.Context, adminSchoolID *uuid.UUID, limit, offset int) ([]model.Parent, int, error) {
			if adminSchoolID == nil || *adminSchoolID != schoolID {
				t.Fatalf("expected admin school scope to be forwarded")
			}
			if limit != 20 || offset != 0 {
				t.Fatalf("unexpected pagination forwarded")
			}
			return []model.Parent{{ParentID: uuid.New(), FullName: "Parent A", SchoolID: schoolID}}, 1, nil
		}}}

		r := gin.New()
		r.Use(withParentAdminScope(schoolID))
		r.GET("/parents", h.List)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/parents?limit=20&offset=0", nil)
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
		}
	})

	tests := []struct {
		name       string
		err        error
		wantStatus int
		wantError  string
	}{
		{name: "access denied", err: service.ErrSchoolAccessDenied, wantStatus: http.StatusForbidden, wantError: "access denied"},
		{name: "internal", err: errors.New("boom"), wantStatus: http.StatusInternalServerError, wantError: "failed to fetch parents"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &ParentHandler{parentService: &fakeParentService{listFn: func(context.Context, *uuid.UUID, int, int) ([]model.Parent, int, error) {
				return nil, 0, tt.err
			}}}
			r := gin.New()
			r.Use(withParentAdminScope(schoolID))
			r.GET("/parents", h.List)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/parents?limit=20&offset=0", nil)
			r.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if got := decodeParentError(t, rec); got != tt.wantError {
				t.Fatalf("error = %q, want %q", got, tt.wantError)
			}
		})
	}
}

func TestGetByID_SuccessAndErrorMapping(t *testing.T) {
	gin.SetMode(gin.TestMode)
	schoolID := uuid.New()
	parentID := uuid.New()

	t.Run("invalid parent id", func(t *testing.T) {
		h := &ParentHandler{}
		r := gin.New()
		r.GET("/parents/:parent_id", h.GetByID)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/parents/not-a-uuid", nil)
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
		if got := decodeParentError(t, rec); got != "invalid parent_id format" {
			t.Fatalf("error = %q, want %q", got, "invalid parent_id format")
		}
	})

	t.Run("success", func(t *testing.T) {
		h := &ParentHandler{parentService: &fakeParentService{getByParentIDFn: func(_ context.Context, adminSchoolID *uuid.UUID, gotParentID uuid.UUID) (*model.Parent, error) {
			if adminSchoolID == nil || *adminSchoolID != schoolID {
				t.Fatalf("expected admin school scope to be forwarded")
			}
			if gotParentID != parentID {
				t.Fatalf("parent_id = %s, want %s", gotParentID, parentID)
			}
			return &model.Parent{ParentID: parentID, SchoolID: schoolID, FullName: "Parent A"}, nil
		}}}

		r := gin.New()
		r.Use(withParentAdminScope(schoolID))
		r.GET("/parents/:parent_id", h.GetByID)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/parents/"+parentID.String(), nil)
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
		}
	})

	tests := []struct {
		name       string
		err        error
		wantStatus int
		wantError  string
	}{
		{name: "access denied", err: service.ErrSchoolAccessDenied, wantStatus: http.StatusForbidden, wantError: "access denied"},
		{name: "not found", err: pgx.ErrNoRows, wantStatus: http.StatusNotFound, wantError: "parent not found"},
		{name: "internal", err: errors.New("boom"), wantStatus: http.StatusInternalServerError, wantError: "failed to fetch parent"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &ParentHandler{parentService: &fakeParentService{getByParentIDFn: func(context.Context, *uuid.UUID, uuid.UUID) (*model.Parent, error) {
				return nil, tt.err
			}}}
			r := gin.New()
			r.Use(withParentAdminScope(schoolID))
			r.GET("/parents/:parent_id", h.GetByID)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/parents/"+parentID.String(), nil)
			r.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if got := decodeParentError(t, rec); got != tt.wantError {
				t.Fatalf("error = %q, want %q", got, tt.wantError)
			}
		})
	}
}
