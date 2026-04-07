package classhandlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/service"
)

func TestListBySchool_SuccessAndErrorMapping(t *testing.T) {
	gin.SetMode(gin.TestMode)
	schoolID := uuid.New()

	t.Run("invalid school id", func(t *testing.T) {
		h := &ClassHandler{}
		r := gin.New()
		r.GET("/classes/by-school/:school_id", h.ListBySchool)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/classes/by-school/not-a-uuid", nil)
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
	})

	t.Run("invalid pagination", func(t *testing.T) {
		h := &ClassHandler{}
		r := gin.New()
		r.GET("/classes/by-school/:school_id", h.ListBySchool)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/classes/by-school/"+schoolID.String()+"?limit=0&offset=-1", nil)
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
	})

	t.Run("success", func(t *testing.T) {
		h := &ClassHandler{classService: &fakeClassService{listBySchoolFn: func(_ context.Context, adminSchoolID *uuid.UUID, gotSchoolID uuid.UUID, limit, offset int) ([]model.Class, int, error) {
			if adminSchoolID == nil || *adminSchoolID != schoolID {
				t.Fatalf("expected admin school scope to be forwarded")
			}
			if gotSchoolID != schoolID || limit != 20 || offset != 0 {
				t.Fatalf("unexpected params forwarded")
			}
			return []model.Class{{ClassID: uuid.New(), SchoolID: schoolID, Name: "A1", SchoolYear: "2026"}}, 1, nil
		}}}
		r := gin.New()
		r.Use(withClassAdminScope(schoolID))
		r.GET("/classes/by-school/:school_id", h.ListBySchool)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/classes/by-school/"+schoolID.String()+"?limit=20&offset=0", nil)
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
		{name: "internal", err: errors.New("boom"), wantStatus: http.StatusInternalServerError, wantError: "failed to fetch classes"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &ClassHandler{classService: &fakeClassService{listBySchoolFn: func(context.Context, *uuid.UUID, uuid.UUID, int, int) ([]model.Class, int, error) {
				return nil, 0, tt.err
			}}}
			r := gin.New()
			r.Use(withClassAdminScope(schoolID))
			r.GET("/classes/by-school/:school_id", h.ListBySchool)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/classes/by-school/"+schoolID.String()+"?limit=20&offset=0", nil)
			r.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if got := decodeClassError(t, rec); got != tt.wantError {
				t.Fatalf("error = %q, want %q", got, tt.wantError)
			}
		})
	}
}
