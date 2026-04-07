package schoolhandlers

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

func TestList_SuccessAndErrorMapping(t *testing.T) {
	gin.SetMode(gin.TestMode)
	schoolID := uuid.New()

	t.Run("invalid pagination", func(t *testing.T) {
		h := &SchoolHandler{}
		r := gin.New()
		r.GET("/schools", h.List)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/schools?limit=0&offset=-1", nil)
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
	})

	t.Run("success", func(t *testing.T) {
		h := &SchoolHandler{schoolService: &fakeSchoolService{listFn: func(_ context.Context, adminSchoolID *uuid.UUID, limit, offset int) ([]model.School, int, error) {
			if adminSchoolID == nil || *adminSchoolID != schoolID {
				t.Fatalf("expected admin school scope to be forwarded")
			}
			if limit != 20 || offset != 0 {
				t.Fatalf("unexpected pagination forwarded")
			}
			return []model.School{{SchoolID: schoolID, Name: "School A"}}, 1, nil
		}}}
		r := gin.New()
		r.Use(withSchoolAdminScope(schoolID))
		r.GET("/schools", h.List)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/schools?limit=20&offset=0", nil)
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
		{name: "internal", err: errors.New("boom"), wantStatus: http.StatusInternalServerError, wantError: "failed to fetch schools"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &SchoolHandler{schoolService: &fakeSchoolService{listFn: func(context.Context, *uuid.UUID, int, int) ([]model.School, int, error) {
				return nil, 0, tt.err
			}}}
			r := gin.New()
			r.Use(withSchoolAdminScope(schoolID))
			r.GET("/schools", h.List)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/schools?limit=20&offset=0", nil)
			r.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if got := decodeSchoolError(t, rec); got != tt.wantError {
				t.Fatalf("error = %q, want %q", got, tt.wantError)
			}
		})
	}
}
