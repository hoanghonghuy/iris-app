package studenthandlers

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
	classID := uuid.New()

	t.Run("invalid request body", func(t *testing.T) {
		h := &StudentHandler{}
		r := gin.New()
		r.POST("/students", h.Create)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/students", strings.NewReader("{bad-json"))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
	})

	t.Run("success", func(t *testing.T) {
		h := &StudentHandler{studentService: &fakeStudentService{createFn: func(_ context.Context, adminSchoolID *uuid.UUID, gotSchoolID, gotClassID uuid.UUID, fullName, dobStr, gender string) (*model.Student, error) {
			if adminSchoolID == nil || *adminSchoolID != schoolID {
				t.Fatalf("expected admin school scope to be forwarded")
			}
			if gotSchoolID != schoolID || gotClassID != classID {
				t.Fatalf("unexpected school/class ids forwarded")
			}
			if fullName != "Student A" || dobStr != "2018-05-10" || gender != "male" {
				t.Fatalf("unexpected payload forwarded")
			}
			id := uuid.New()
			return &model.Student{StudentID: id}, nil
		}}}

		r := gin.New()
		r.Use(withStudentAdminScope(schoolID))
		r.POST("/students", h.Create)

		body := `{"school_id":"` + schoolID.String() + `","current_class_id":"` + classID.String() + `","full_name":"Student A","dob":"2018-05-10","gender":"male"}`
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/students", strings.NewReader(body))
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
		{name: "invalid dob", err: service.ErrInvalidValue, wantStatus: http.StatusBadRequest, wantError: "invalid dob format (expected YYYY-MM-DD)"},
		{name: "other create error", err: errors.New("boom"), wantStatus: http.StatusBadRequest, wantError: "failed to create student"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &StudentHandler{studentService: &fakeStudentService{createFn: func(context.Context, *uuid.UUID, uuid.UUID, uuid.UUID, string, string, string) (*model.Student, error) {
				return nil, tt.err
			}}}
			r := gin.New()
			r.Use(withStudentAdminScope(schoolID))
			r.POST("/students", h.Create)

			body := `{"school_id":"` + schoolID.String() + `","current_class_id":"` + classID.String() + `","full_name":"Student A","dob":"2018-05-10","gender":"male"}`
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/students", strings.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if got := decodeStudentError(t, rec); got != tt.wantError {
				t.Fatalf("error = %q, want %q", got, tt.wantError)
			}
		})
	}
}

func TestUpdate_SuccessAndErrorMapping(t *testing.T) {
	gin.SetMode(gin.TestMode)
	schoolID := uuid.New()
	studentID := uuid.New()

	t.Run("invalid student id", func(t *testing.T) {
		h := &StudentHandler{}
		r := gin.New()
		r.PUT("/students/:student_id", h.Update)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, "/students/not-a-uuid", strings.NewReader(`{"full_name":"A","dob":"2018-01-01","gender":"male"}`))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
	})

	t.Run("invalid request body", func(t *testing.T) {
		h := &StudentHandler{}
		r := gin.New()
		r.PUT("/students/:student_id", h.Update)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, "/students/"+studentID.String(), strings.NewReader("{bad-json"))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
	})

	t.Run("success", func(t *testing.T) {
		h := &StudentHandler{studentService: &fakeStudentService{updateFn: func(_ context.Context, adminSchoolID *uuid.UUID, gotStudentID uuid.UUID, fullName, dobStr, gender string) error {
			if adminSchoolID == nil || *adminSchoolID != schoolID {
				t.Fatalf("expected admin school scope to be forwarded")
			}
			if gotStudentID != studentID {
				t.Fatalf("student_id = %s, want %s", gotStudentID, studentID)
			}
			if fullName != "Student A" || dobStr != "2018-05-10" || gender != "male" {
				t.Fatalf("unexpected payload forwarded")
			}
			return nil
		}}}
		r := gin.New()
		r.Use(withStudentAdminScope(schoolID))
		r.PUT("/students/:student_id", h.Update)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, "/students/"+studentID.String(), strings.NewReader(`{"full_name":"Student A","dob":"2018-05-10","gender":"male"}`))
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
		{name: "student not found", err: service.ErrStudentNotFound, wantStatus: http.StatusNotFound, wantError: "student not found"},
		{name: "invalid value", err: service.ErrInvalidValue, wantStatus: http.StatusBadRequest, wantError: "invalid dob format (expected YYYY-MM-DD)"},
		{name: "internal", err: errors.New("boom"), wantStatus: http.StatusInternalServerError, wantError: "failed to update student"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &StudentHandler{studentService: &fakeStudentService{updateFn: func(context.Context, *uuid.UUID, uuid.UUID, string, string, string) error {
				return tt.err
			}}}
			r := gin.New()
			r.Use(withStudentAdminScope(schoolID))
			r.PUT("/students/:student_id", h.Update)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPut, "/students/"+studentID.String(), strings.NewReader(`{"full_name":"Student A","dob":"2018-05-10","gender":"male"}`))
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if got := decodeStudentError(t, rec); got != tt.wantError {
				t.Fatalf("error = %q, want %q", got, tt.wantError)
			}
		})
	}
}

func TestDelete_SuccessAndErrorMapping(t *testing.T) {
	gin.SetMode(gin.TestMode)
	schoolID := uuid.New()
	studentID := uuid.New()

	t.Run("invalid student id", func(t *testing.T) {
		h := &StudentHandler{}
		r := gin.New()
		r.DELETE("/students/:student_id", h.Delete)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/students/not-a-uuid", nil)
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
	})

	t.Run("success", func(t *testing.T) {
		h := &StudentHandler{studentService: &fakeStudentService{deleteFn: func(_ context.Context, adminSchoolID *uuid.UUID, gotStudentID uuid.UUID) error {
			if adminSchoolID == nil || *adminSchoolID != schoolID {
				t.Fatalf("expected admin school scope to be forwarded")
			}
			if gotStudentID != studentID {
				t.Fatalf("student_id = %s, want %s", gotStudentID, studentID)
			}
			return nil
		}}}

		r := gin.New()
		r.Use(withStudentAdminScope(schoolID))
		r.DELETE("/students/:student_id", h.Delete)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/students/"+studentID.String(), nil)
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
		{name: "student not found", err: service.ErrStudentNotFound, wantStatus: http.StatusNotFound, wantError: "student not found"},
		{name: "internal", err: errors.New("boom"), wantStatus: http.StatusInternalServerError, wantError: "failed to delete student"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &StudentHandler{studentService: &fakeStudentService{deleteFn: func(context.Context, *uuid.UUID, uuid.UUID) error {
				return tt.err
			}}}
			r := gin.New()
			r.Use(withStudentAdminScope(schoolID))
			r.DELETE("/students/:student_id", h.Delete)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodDelete, "/students/"+studentID.String(), nil)
			r.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if got := decodeStudentError(t, rec); got != tt.wantError {
				t.Fatalf("error = %q, want %q", got, tt.wantError)
			}
		})
	}
}
