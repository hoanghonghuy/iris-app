package classhandlers

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
	schoolID := uuid.New()

	t.Run("invalid body", func(t *testing.T) {
		h := &ClassHandler{}
		r := gin.New()
		r.POST("/classes", h.Create)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/classes", strings.NewReader("{bad-json"))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
	})

	t.Run("success", func(t *testing.T) {
		classID := uuid.New()
		h := &ClassHandler{classService: &fakeClassService{createFn: func(_ context.Context, adminSchoolID *uuid.UUID, gotSchoolID uuid.UUID, name, schoolYear string) (*model.Class, error) {
			if adminSchoolID == nil || *adminSchoolID != schoolID {
				t.Fatalf("expected admin school scope to be forwarded")
			}
			if gotSchoolID != schoolID || name != "A1" || schoolYear != "2026" {
				t.Fatalf("unexpected create payload")
			}
			return &model.Class{ClassID: classID}, nil
		}}}
		r := gin.New()
		r.Use(withClassAdminScope(schoolID))
		r.POST("/classes", h.Create)

		body := `{"school_id":"` + schoolID.String() + `","name":"A1","school_year":"2026"}`
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/classes", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusCreated {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusCreated)
		}
	})

	tests := []struct {
		name       string
		err        error
		wantStatus int
		wantError  string
	}{
		{name: "access denied", err: service.ErrSchoolAccessDenied, wantStatus: http.StatusForbidden, wantError: "access denied"},
		{name: "internal", err: errors.New("boom"), wantStatus: http.StatusInternalServerError, wantError: "failed to create class"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &ClassHandler{classService: &fakeClassService{createFn: func(context.Context, *uuid.UUID, uuid.UUID, string, string) (*model.Class, error) {
				return nil, tt.err
			}}}
			r := gin.New()
			r.Use(withClassAdminScope(schoolID))
			r.POST("/classes", h.Create)

			body := `{"school_id":"` + schoolID.String() + `","name":"A1","school_year":"2026"}`
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/classes", strings.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
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

func TestUpdate_SuccessAndErrorMapping(t *testing.T) {
	gin.SetMode(gin.TestMode)
	schoolID := uuid.New()
	classID := uuid.New()

	t.Run("invalid class id", func(t *testing.T) {
		h := &ClassHandler{}
		r := gin.New()
		r.PUT("/classes/:class_id", h.Update)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, "/classes/not-a-uuid", strings.NewReader(`{"name":"A2","school_year":"2026"}`))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
	})

	t.Run("invalid body", func(t *testing.T) {
		h := &ClassHandler{}
		r := gin.New()
		r.PUT("/classes/:class_id", h.Update)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, "/classes/"+classID.String(), strings.NewReader("{bad-json"))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
	})

	t.Run("success", func(t *testing.T) {
		h := &ClassHandler{classService: &fakeClassService{updateFn: func(_ context.Context, adminSchoolID *uuid.UUID, gotClassID uuid.UUID, name, schoolYear string) error {
			if adminSchoolID == nil || *adminSchoolID != schoolID {
				t.Fatalf("expected admin school scope to be forwarded")
			}
			if gotClassID != classID || name != "A2" || schoolYear != "2026" {
				t.Fatalf("unexpected update payload")
			}
			return nil
		}}}
		r := gin.New()
		r.Use(withClassAdminScope(schoolID))
		r.PUT("/classes/:class_id", h.Update)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, "/classes/"+classID.String(), strings.NewReader(`{"name":"A2","school_year":"2026"}`))
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
		{name: "access denied", err: service.ErrSchoolAccessDenied, wantStatus: http.StatusForbidden, wantError: "access denied"},
		{name: "not found", err: service.ErrClassNotFound, wantStatus: http.StatusNotFound, wantError: "class not found"},
		{name: "internal", err: errors.New("boom"), wantStatus: http.StatusInternalServerError, wantError: "failed to update class"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &ClassHandler{classService: &fakeClassService{updateFn: func(context.Context, *uuid.UUID, uuid.UUID, string, string) error {
				return tt.err
			}}}
			r := gin.New()
			r.Use(withClassAdminScope(schoolID))
			r.PUT("/classes/:class_id", h.Update)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPut, "/classes/"+classID.String(), strings.NewReader(`{"name":"A2","school_year":"2026"}`))
			req.Header.Set("Content-Type", "application/json")
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

func TestDelete_SuccessAndErrorMapping(t *testing.T) {
	gin.SetMode(gin.TestMode)
	schoolID := uuid.New()
	classID := uuid.New()

	t.Run("invalid class id", func(t *testing.T) {
		h := &ClassHandler{}
		r := gin.New()
		r.DELETE("/classes/:class_id", h.Delete)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/classes/not-a-uuid", nil)
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
	})

	t.Run("success", func(t *testing.T) {
		h := &ClassHandler{classService: &fakeClassService{deleteFn: func(_ context.Context, adminSchoolID *uuid.UUID, gotClassID uuid.UUID) error {
			if adminSchoolID == nil || *adminSchoolID != schoolID {
				t.Fatalf("expected admin school scope to be forwarded")
			}
			if gotClassID != classID {
				t.Fatalf("unexpected class id forwarded")
			}
			return nil
		}}}
		r := gin.New()
		r.Use(withClassAdminScope(schoolID))
		r.DELETE("/classes/:class_id", h.Delete)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/classes/"+classID.String(), nil)
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
		{name: "not found", err: service.ErrClassNotFound, wantStatus: http.StatusNotFound, wantError: "class not found"},
		{name: "internal", err: errors.New("boom"), wantStatus: http.StatusInternalServerError, wantError: "failed to delete class"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &ClassHandler{classService: &fakeClassService{deleteFn: func(context.Context, *uuid.UUID, uuid.UUID) error {
				return tt.err
			}}}
			r := gin.New()
			r.Use(withClassAdminScope(schoolID))
			r.DELETE("/classes/:class_id", h.Delete)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodDelete, "/classes/"+classID.String(), nil)
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
