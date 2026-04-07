package parentcodehandlers

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

func TestGenerateCodeForStudent(t *testing.T) {
	gin.SetMode(gin.TestMode)
	schoolID := uuid.New()
	studentID := uuid.New()

	t.Run("invalid student id", func(t *testing.T) {
		h := &ParentCodeHandler{}
		r := gin.New()
		r.POST("/students/:student_id/parent-code", h.GenerateCodeForStudent)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/students/not-a-uuid/parent-code", nil)
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
	})

	t.Run("success", func(t *testing.T) {
		h := &ParentCodeHandler{parentCodeService: &fakeParentCodeService{generateCodeForStudentFn: func(_ context.Context, adminSchoolID *uuid.UUID, gotStudentID uuid.UUID) (string, error) {
			if adminSchoolID == nil || *adminSchoolID != schoolID || gotStudentID != studentID {
				t.Fatalf("unexpected params")
			}
			return "ABCD1234", nil
		}}}
		r := gin.New()
		r.Use(withParentCodeAdminScope(schoolID))
		r.POST("/students/:student_id/parent-code", h.GenerateCodeForStudent)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/students/"+studentID.String()+"/parent-code", nil)
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
		}
	})

	t.Run("service mappings", func(t *testing.T) {
		cases := []struct {
			name       string
			err        error
			wantStatus int
			wantError  string
		}{
			{name: "access denied", err: service.ErrSchoolAccessDenied, wantStatus: http.StatusForbidden, wantError: "access denied"},
			{name: "internal", err: errors.New("boom"), wantStatus: http.StatusInternalServerError, wantError: "failed to generate parent code"},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				h := &ParentCodeHandler{parentCodeService: &fakeParentCodeService{generateCodeForStudentFn: func(context.Context, *uuid.UUID, uuid.UUID) (string, error) {
					return "", tc.err
				}}}
				r := gin.New()
				r.Use(withParentCodeAdminScope(schoolID))
				r.POST("/students/:student_id/parent-code", h.GenerateCodeForStudent)

				rec := httptest.NewRecorder()
				req := httptest.NewRequest(http.MethodPost, "/students/"+studentID.String()+"/parent-code", nil)
				r.ServeHTTP(rec, req)

				if rec.Code != tc.wantStatus {
					t.Fatalf("status = %d, want %d", rec.Code, tc.wantStatus)
				}
				if got := decodeParentCodeError(t, rec); got != tc.wantError {
					t.Fatalf("error = %q, want %q", got, tc.wantError)
				}
			})
		}
	})
}

func TestRevokeParentCode(t *testing.T) {
	gin.SetMode(gin.TestMode)
	schoolID := uuid.New()
	studentID := uuid.New()

	t.Run("invalid student id", func(t *testing.T) {
		h := &ParentCodeHandler{}
		r := gin.New()
		r.DELETE("/students/:student_id/parent-code", h.RevokeParentCode)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/students/not-a-uuid/parent-code", nil)
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
	})

	t.Run("success", func(t *testing.T) {
		h := &ParentCodeHandler{parentCodeService: &fakeParentCodeService{revokeCodeFn: func(_ context.Context, adminSchoolID *uuid.UUID, gotStudentID uuid.UUID) error {
			if adminSchoolID == nil || *adminSchoolID != schoolID || gotStudentID != studentID {
				t.Fatalf("unexpected params")
			}
			return nil
		}}}
		r := gin.New()
		r.Use(withParentCodeAdminScope(schoolID))
		r.DELETE("/students/:student_id/parent-code", h.RevokeParentCode)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/students/"+studentID.String()+"/parent-code", nil)
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
		}
	})

	t.Run("service mappings", func(t *testing.T) {
		cases := []struct {
			name       string
			err        error
			wantStatus int
			wantError  string
		}{
			{name: "access denied", err: service.ErrSchoolAccessDenied, wantStatus: http.StatusForbidden, wantError: "access denied"},
			{name: "internal", err: errors.New("boom"), wantStatus: http.StatusInternalServerError, wantError: "failed to revoke parent code"},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				h := &ParentCodeHandler{parentCodeService: &fakeParentCodeService{revokeCodeFn: func(context.Context, *uuid.UUID, uuid.UUID) error {
					return tc.err
				}}}
				r := gin.New()
				r.Use(withParentCodeAdminScope(schoolID))
				r.DELETE("/students/:student_id/parent-code", h.RevokeParentCode)

				rec := httptest.NewRecorder()
				req := httptest.NewRequest(http.MethodDelete, "/students/"+studentID.String()+"/parent-code", nil)
				r.ServeHTTP(rec, req)

				if rec.Code != tc.wantStatus {
					t.Fatalf("status = %d, want %d", rec.Code, tc.wantStatus)
				}
				if got := decodeParentCodeError(t, rec); got != tc.wantError {
					t.Fatalf("error = %q, want %q", got, tc.wantError)
				}
			})
		}
	})
}
