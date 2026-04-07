package parenthandlers

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

func TestAssignStudent_SuccessAndErrorMapping(t *testing.T) {
	gin.SetMode(gin.TestMode)
	schoolID := uuid.New()
	parentID := uuid.New()
	studentID := uuid.New()

	t.Run("invalid parent id", func(t *testing.T) {
		h := &ParentHandler{}
		r := gin.New()
		r.POST("/parents/:parent_id/students/:student_id", h.AssignStudent)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/parents/not-a-uuid/students/"+studentID.String(), nil)
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
	})

	t.Run("invalid student id", func(t *testing.T) {
		h := &ParentHandler{}
		r := gin.New()
		r.POST("/parents/:parent_id/students/:student_id", h.AssignStudent)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/parents/"+parentID.String()+"/students/not-a-uuid", nil)
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
	})

	t.Run("invalid json falls back to empty relationship and still calls service", func(t *testing.T) {
		called := false
		h := &ParentHandler{parentService: &fakeParentService{assignStudentFn: func(_ context.Context, adminSchoolID *uuid.UUID, gotParentID, gotStudentID uuid.UUID, relationship string) error {
			called = true
			if adminSchoolID == nil || *adminSchoolID != schoolID {
				t.Fatalf("expected admin school scope to be forwarded")
			}
			if gotParentID != parentID || gotStudentID != studentID {
				t.Fatalf("unexpected ids forwarded")
			}
			if relationship != "" {
				t.Fatalf("relationship = %q, want empty", relationship)
			}
			return nil
		}}}

		r := gin.New()
		r.Use(withParentAdminScope(schoolID))
		r.POST("/parents/:parent_id/students/:student_id", h.AssignStudent)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/parents/"+parentID.String()+"/students/"+studentID.String(), strings.NewReader("{bad-json"))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
		}
		if !called {
			t.Fatalf("expected AssignStudent to be called")
		}
	})

	t.Run("success", func(t *testing.T) {
		h := &ParentHandler{parentService: &fakeParentService{assignStudentFn: func(_ context.Context, adminSchoolID *uuid.UUID, gotParentID, gotStudentID uuid.UUID, relationship string) error {
			if adminSchoolID == nil || *adminSchoolID != schoolID {
				t.Fatalf("expected admin school scope to be forwarded")
			}
			if gotParentID != parentID || gotStudentID != studentID || relationship != "father" {
				t.Fatalf("unexpected assign payload")
			}
			return nil
		}}}

		r := gin.New()
		r.Use(withParentAdminScope(schoolID))
		r.POST("/parents/:parent_id/students/:student_id", h.AssignStudent)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/parents/"+parentID.String()+"/students/"+studentID.String(), strings.NewReader(`{"relationship":"father"}`))
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
		{name: "internal", err: errors.New("boom"), wantStatus: http.StatusInternalServerError, wantError: "failed to assign parent to student"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &ParentHandler{parentService: &fakeParentService{assignStudentFn: func(context.Context, *uuid.UUID, uuid.UUID, uuid.UUID, string) error {
				return tt.err
			}}}
			r := gin.New()
			r.Use(withParentAdminScope(schoolID))
			r.POST("/parents/:parent_id/students/:student_id", h.AssignStudent)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/parents/"+parentID.String()+"/students/"+studentID.String(), strings.NewReader(`{"relationship":"mother"}`))
			req.Header.Set("Content-Type", "application/json")
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

func TestUnassignStudent_SuccessAndErrorMapping(t *testing.T) {
	gin.SetMode(gin.TestMode)
	schoolID := uuid.New()
	parentID := uuid.New()
	studentID := uuid.New()

	t.Run("invalid parent id", func(t *testing.T) {
		h := &ParentHandler{}
		r := gin.New()
		r.DELETE("/parents/:parent_id/students/:student_id", h.UnassignStudent)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/parents/not-a-uuid/students/"+studentID.String(), nil)
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
	})

	t.Run("invalid student id", func(t *testing.T) {
		h := &ParentHandler{}
		r := gin.New()
		r.DELETE("/parents/:parent_id/students/:student_id", h.UnassignStudent)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/parents/"+parentID.String()+"/students/not-a-uuid", nil)
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
	})

	t.Run("success", func(t *testing.T) {
		h := &ParentHandler{parentService: &fakeParentService{unassignFn: func(_ context.Context, adminSchoolID *uuid.UUID, gotParentID, gotStudentID uuid.UUID) error {
			if adminSchoolID == nil || *adminSchoolID != schoolID {
				t.Fatalf("expected admin school scope to be forwarded")
			}
			if gotParentID != parentID || gotStudentID != studentID {
				t.Fatalf("unexpected ids forwarded")
			}
			return nil
		}}}

		r := gin.New()
		r.Use(withParentAdminScope(schoolID))
		r.DELETE("/parents/:parent_id/students/:student_id", h.UnassignStudent)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/parents/"+parentID.String()+"/students/"+studentID.String(), nil)
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
		{name: "internal", err: errors.New("boom"), wantStatus: http.StatusInternalServerError, wantError: "failed to unassign parent from student"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &ParentHandler{parentService: &fakeParentService{unassignFn: func(context.Context, *uuid.UUID, uuid.UUID, uuid.UUID) error {
				return tt.err
			}}}
			r := gin.New()
			r.Use(withParentAdminScope(schoolID))
			r.DELETE("/parents/:parent_id/students/:student_id", h.UnassignStudent)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodDelete, "/parents/"+parentID.String()+"/students/"+studentID.String(), nil)
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
