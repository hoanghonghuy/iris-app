package teacherhandlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/service"
)

func TestUpdate_SuccessAndErrorMapping(t *testing.T) {
	gin.SetMode(gin.TestMode)
	schoolID := uuid.New()
	teacherID := uuid.New()

	t.Run("invalid teacher id", func(t *testing.T) {
		h := &TeacherHandler{}
		r := gin.New()
		r.PUT("/teachers/:teacher_id", h.Update)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, "/teachers/not-a-uuid", strings.NewReader(`{"full_name":"A","phone":"1","school_id":"`+schoolID.String()+`"}`))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
	})

	t.Run("invalid request body", func(t *testing.T) {
		h := &TeacherHandler{}
		r := gin.New()
		r.PUT("/teachers/:teacher_id", h.Update)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, "/teachers/"+teacherID.String(), strings.NewReader("{bad-json"))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
	})

	t.Run("success", func(t *testing.T) {
		h := &TeacherHandler{teacherService: &fakeTeacherService{updateFn: func(_ context.Context, adminSchoolID *uuid.UUID, gotTeacherID uuid.UUID, fullName, phone string, gotSchoolID uuid.UUID) error {
			if adminSchoolID == nil || *adminSchoolID != schoolID {
				t.Fatalf("expected admin school scope to be forwarded")
			}
			if gotTeacherID != teacherID || gotSchoolID != schoolID {
				t.Fatalf("unexpected ids forwarded")
			}
			if fullName != "Teacher A" || phone != "0909" {
				t.Fatalf("unexpected payload forwarded")
			}
			return nil
		}}}
		r := gin.New()
		r.Use(withAdminSchoolScope(schoolID))
		r.PUT("/teachers/:teacher_id", h.Update)

		body := `{"full_name":"Teacher A","phone":"0909","school_id":"` + schoolID.String() + `"}`
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, "/teachers/"+teacherID.String(), strings.NewReader(body))
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
		{name: "internal", err: errors.New("boom"), wantStatus: http.StatusInternalServerError, wantError: "failed to update teacher"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &TeacherHandler{teacherService: &fakeTeacherService{updateFn: func(context.Context, *uuid.UUID, uuid.UUID, string, string, uuid.UUID) error {
				return tt.err
			}}}
			r := gin.New()
			r.Use(withAdminSchoolScope(schoolID))
			r.PUT("/teachers/:teacher_id", h.Update)

			body := `{"full_name":"Teacher A","phone":"0909","school_id":"` + schoolID.String() + `"}`
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPut, "/teachers/"+teacherID.String(), strings.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if got := decodeTeacherError(t, rec); got != tt.wantError {
				t.Fatalf("error = %q, want %q", got, tt.wantError)
			}
		})
	}
}

func TestDelete_SuccessAndErrorMapping(t *testing.T) {
	gin.SetMode(gin.TestMode)
	schoolID := uuid.New()
	teacherID := uuid.New()

	t.Run("invalid teacher id", func(t *testing.T) {
		h := &TeacherHandler{}
		r := gin.New()
		r.DELETE("/teachers/:teacher_id", h.Delete)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/teachers/not-a-uuid", nil)
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
	})

	t.Run("success", func(t *testing.T) {
		h := &TeacherHandler{teacherService: &fakeTeacherService{deleteFn: func(_ context.Context, adminSchoolID *uuid.UUID, gotTeacherID uuid.UUID) error {
			if adminSchoolID == nil || *adminSchoolID != schoolID {
				t.Fatalf("expected admin school scope to be forwarded")
			}
			if gotTeacherID != teacherID {
				t.Fatalf("teacher_id = %s, want %s", gotTeacherID, teacherID)
			}
			return nil
		}}}
		r := gin.New()
		r.Use(withAdminSchoolScope(schoolID))
		r.DELETE("/teachers/:teacher_id", h.Delete)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/teachers/"+teacherID.String(), nil)
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
		{name: "teacher not found", err: service.ErrTeacherNotFound, wantStatus: http.StatusNotFound, wantError: "teacher not found"},
		{name: "internal", err: errors.New("boom"), wantStatus: http.StatusInternalServerError, wantError: "failed to delete teacher"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &TeacherHandler{teacherService: &fakeTeacherService{deleteFn: func(context.Context, *uuid.UUID, uuid.UUID) error {
				return tt.err
			}}}
			r := gin.New()
			r.Use(withAdminSchoolScope(schoolID))
			r.DELETE("/teachers/:teacher_id", h.Delete)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodDelete, "/teachers/"+teacherID.String(), nil)
			r.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if got := decodeTeacherError(t, rec); got != tt.wantError {
				t.Fatalf("error = %q, want %q", got, tt.wantError)
			}
		})
	}
}
