package teacherhandlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/service"
)

func TestAssign_SuccessAndErrorMapping(t *testing.T) {
	gin.SetMode(gin.TestMode)
	schoolID := uuid.New()
	teacherID := uuid.New()
	classID := uuid.New()

	t.Run("invalid teacher id", func(t *testing.T) {
		h := &TeacherHandler{}
		r := gin.New()
		r.POST("/teachers/:teacher_id/classes/:class_id", h.Assign)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/teachers/not-a-uuid/classes/"+classID.String(), nil)
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
	})

	t.Run("invalid class id", func(t *testing.T) {
		h := &TeacherHandler{}
		r := gin.New()
		r.POST("/teachers/:teacher_id/classes/:class_id", h.Assign)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/teachers/"+teacherID.String()+"/classes/not-a-uuid", nil)
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
	})

	t.Run("success", func(t *testing.T) {
		h := &TeacherHandler{teacherService: &fakeTeacherService{assignFn: func(_ context.Context, adminSchoolID *uuid.UUID, gotTeacherID, gotClassID uuid.UUID) error {
			if adminSchoolID == nil || *adminSchoolID != schoolID {
				t.Fatalf("expected admin school scope to be forwarded")
			}
			if gotTeacherID != teacherID || gotClassID != classID {
				t.Fatalf("unexpected ids forwarded")
			}
			return nil
		}}}

		r := gin.New()
		r.Use(withAdminSchoolScope(schoolID))
		r.POST("/teachers/:teacher_id/classes/:class_id", h.Assign)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/teachers/"+teacherID.String()+"/classes/"+classID.String(), nil)
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
		{name: "internal", err: errors.New("boom"), wantStatus: http.StatusInternalServerError, wantError: "failed to assign teacher to class"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &TeacherHandler{teacherService: &fakeTeacherService{assignFn: func(context.Context, *uuid.UUID, uuid.UUID, uuid.UUID) error {
				return tt.err
			}}}
			r := gin.New()
			r.Use(withAdminSchoolScope(schoolID))
			r.POST("/teachers/:teacher_id/classes/:class_id", h.Assign)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/teachers/"+teacherID.String()+"/classes/"+classID.String(), nil)
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

func TestUnassign_SuccessAndErrorMapping(t *testing.T) {
	gin.SetMode(gin.TestMode)
	schoolID := uuid.New()
	teacherID := uuid.New()
	classID := uuid.New()

	t.Run("invalid teacher id", func(t *testing.T) {
		h := &TeacherHandler{}
		r := gin.New()
		r.DELETE("/teachers/:teacher_id/classes/:class_id", h.Unassign)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/teachers/not-a-uuid/classes/"+classID.String(), nil)
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
	})

	t.Run("invalid class id", func(t *testing.T) {
		h := &TeacherHandler{}
		r := gin.New()
		r.DELETE("/teachers/:teacher_id/classes/:class_id", h.Unassign)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/teachers/"+teacherID.String()+"/classes/not-a-uuid", nil)
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
	})

	t.Run("success", func(t *testing.T) {
		h := &TeacherHandler{teacherService: &fakeTeacherService{unassignFn: func(_ context.Context, adminSchoolID *uuid.UUID, gotTeacherID, gotClassID uuid.UUID) error {
			if adminSchoolID == nil || *adminSchoolID != schoolID {
				t.Fatalf("expected admin school scope to be forwarded")
			}
			if gotTeacherID != teacherID || gotClassID != classID {
				t.Fatalf("unexpected ids forwarded")
			}
			return nil
		}}}

		r := gin.New()
		r.Use(withAdminSchoolScope(schoolID))
		r.DELETE("/teachers/:teacher_id/classes/:class_id", h.Unassign)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/teachers/"+teacherID.String()+"/classes/"+classID.String(), nil)
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
		{name: "internal", err: errors.New("boom"), wantStatus: http.StatusInternalServerError, wantError: "failed to unassign teacher from class"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &TeacherHandler{teacherService: &fakeTeacherService{unassignFn: func(context.Context, *uuid.UUID, uuid.UUID, uuid.UUID) error {
				return tt.err
			}}}
			r := gin.New()
			r.Use(withAdminSchoolScope(schoolID))
			r.DELETE("/teachers/:teacher_id/classes/:class_id", h.Unassign)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodDelete, "/teachers/"+teacherID.String()+"/classes/"+classID.String(), nil)
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
