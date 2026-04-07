package schoolhandlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/service"
)

func TestCreate_SuccessAndErrorMapping(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("invalid body", func(t *testing.T) {
		h := &SchoolHandler{}
		r := gin.New()
		r.POST("/schools", h.Create)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/schools", strings.NewReader("{bad-json"))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
	})

	t.Run("success", func(t *testing.T) {
		schoolID := uuid.New()
		h := &SchoolHandler{schoolService: &fakeSchoolService{createFn: func(_ context.Context, name, address string) (*model.School, error) {
			if name != "School A" || address != "Address A" {
				t.Fatalf("unexpected create payload")
			}
			return &model.School{SchoolID: schoolID}, nil
		}}}
		r := gin.New()
		r.POST("/schools", h.Create)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/schools", strings.NewReader(`{"name":"School A","address":"Address A"}`))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusCreated {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusCreated)
		}
	})

	t.Run("internal", func(t *testing.T) {
		h := &SchoolHandler{schoolService: &fakeSchoolService{createFn: func(context.Context, string, string) (*model.School, error) {
			return nil, errors.New("boom")
		}}}
		r := gin.New()
		r.POST("/schools", h.Create)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/schools", strings.NewReader(`{"name":"School A"}`))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusInternalServerError {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
		}
		if got := decodeSchoolError(t, rec); got != "failed to create school" {
			t.Fatalf("error = %q, want %q", got, "failed to create school")
		}
	})
}

func TestUpdate_SuccessAndErrorMapping(t *testing.T) {
	gin.SetMode(gin.TestMode)
	schoolID := uuid.New()

	t.Run("invalid school id", func(t *testing.T) {
		h := &SchoolHandler{}
		r := gin.New()
		r.PUT("/schools/:school_id", h.Update)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, "/schools/not-a-uuid", strings.NewReader(`{"name":"School A"}`))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
	})

	t.Run("invalid body", func(t *testing.T) {
		h := &SchoolHandler{}
		r := gin.New()
		r.PUT("/schools/:school_id", h.Update)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, "/schools/"+schoolID.String(), strings.NewReader("{bad-json"))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
	})

	t.Run("success", func(t *testing.T) {
		h := &SchoolHandler{schoolService: &fakeSchoolService{updateFn: func(_ context.Context, gotSchoolID uuid.UUID, name, address string) error {
			if gotSchoolID != schoolID || name != "School B" || address != "Address B" {
				t.Fatalf("unexpected update payload")
			}
			return nil
		}}}
		r := gin.New()
		r.PUT("/schools/:school_id", h.Update)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, "/schools/"+schoolID.String(), strings.NewReader(`{"name":"School B","address":"Address B"}`))
		req.Header.Set("Content-Type", "application/json")
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
		{name: "not found", err: service.ErrSchoolNotFound, wantStatus: http.StatusNotFound, wantError: "school not found"},
		{name: "internal", err: errors.New("boom"), wantStatus: http.StatusInternalServerError, wantError: "failed to update school"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &SchoolHandler{schoolService: &fakeSchoolService{updateFn: func(context.Context, uuid.UUID, string, string) error {
				return tt.err
			}}}
			r := gin.New()
			r.PUT("/schools/:school_id", h.Update)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPut, "/schools/"+schoolID.String(), strings.NewReader(`{"name":"School B","address":"Address B"}`))
			req.Header.Set("Content-Type", "application/json")
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

func TestDelete_SuccessAndErrorMapping(t *testing.T) {
	gin.SetMode(gin.TestMode)
	schoolID := uuid.New()

	t.Run("invalid school id", func(t *testing.T) {
		h := &SchoolHandler{}
		r := gin.New()
		r.DELETE("/schools/:school_id", h.Delete)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/schools/not-a-uuid", nil)
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
	})

	t.Run("success", func(t *testing.T) {
		h := &SchoolHandler{schoolService: &fakeSchoolService{deleteFn: func(_ context.Context, gotSchoolID uuid.UUID) error {
			if gotSchoolID != schoolID {
				t.Fatalf("unexpected school id forwarded")
			}
			return nil
		}}}
		r := gin.New()
		r.DELETE("/schools/:school_id", h.Delete)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/schools/"+schoolID.String(), nil)
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
		{name: "not found", err: service.ErrSchoolNotFound, wantStatus: http.StatusNotFound, wantError: "school not found"},
		{name: "internal", err: errors.New("boom"), wantStatus: http.StatusInternalServerError, wantError: "failed to delete school"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &SchoolHandler{schoolService: &fakeSchoolService{deleteFn: func(context.Context, uuid.UUID) error {
				return tt.err
			}}}
			r := gin.New()
			r.DELETE("/schools/:school_id", h.Delete)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodDelete, "/schools/"+schoolID.String(), nil)
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
