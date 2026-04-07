package studenthandlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/service"
)

func TestListByClass_SuccessAndErrorMapping(t *testing.T) {
	gin.SetMode(gin.TestMode)
	schoolID := uuid.New()
	classID := uuid.New()

	t.Run("invalid class id", func(t *testing.T) {
		h := &StudentHandler{}
		r := gin.New()
		r.GET("/students/by-class/:class_id", h.ListByClass)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/students/by-class/not-a-uuid", nil)
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
	})

	t.Run("invalid pagination", func(t *testing.T) {
		h := &StudentHandler{}
		r := gin.New()
		r.GET("/students/by-class/:class_id", h.ListByClass)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/students/by-class/"+classID.String()+"?limit=0&offset=-1", nil)
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
	})

	t.Run("success", func(t *testing.T) {
		h := &StudentHandler{studentService: &fakeStudentService{listByClassFn: func(_ context.Context, adminSchoolID *uuid.UUID, gotClassID uuid.UUID, limit, offset int) ([]model.Student, int, error) {
			if adminSchoolID == nil || *adminSchoolID != schoolID {
				t.Fatalf("expected admin school scope to be forwarded")
			}
			if gotClassID != classID || limit != 10 || offset != 2 {
				t.Fatalf("unexpected params forwarded")
			}
			return []model.Student{{StudentID: uuid.New(), SchoolID: schoolID, CurrentClassID: classID, FullName: "Student A", DOB: time.Now(), Gender: "male"}}, 1, nil
		}}}

		r := gin.New()
		r.Use(withStudentAdminScope(schoolID))
		r.GET("/students/by-class/:class_id", h.ListByClass)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/students/by-class/"+classID.String()+"?limit=10&offset=2", nil)
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
		{name: "internal", err: errors.New("boom"), wantStatus: http.StatusInternalServerError, wantError: "failed to fetch students"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &StudentHandler{studentService: &fakeStudentService{listByClassFn: func(context.Context, *uuid.UUID, uuid.UUID, int, int) ([]model.Student, int, error) {
				return nil, 0, tt.err
			}}}

			r := gin.New()
			r.Use(withStudentAdminScope(schoolID))
			r.GET("/students/by-class/:class_id", h.ListByClass)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/students/by-class/"+classID.String()+"?limit=20&offset=0", nil)
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

func TestGetProfile_SuccessAndErrorMapping(t *testing.T) {
	gin.SetMode(gin.TestMode)
	schoolID := uuid.New()
	studentID := uuid.New()

	t.Run("invalid student id", func(t *testing.T) {
		h := &StudentHandler{}
		r := gin.New()
		r.GET("/students/:student_id", h.GetProfile)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/students/not-a-uuid", nil)
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
	})

	t.Run("success", func(t *testing.T) {
		h := &StudentHandler{studentService: &fakeStudentService{getProfileFn: func(_ context.Context, adminSchoolID *uuid.UUID, gotStudentID uuid.UUID) (*model.StudentProfile, error) {
			if adminSchoolID == nil || *adminSchoolID != schoolID {
				t.Fatalf("expected admin school scope to be forwarded")
			}
			if gotStudentID != studentID {
				t.Fatalf("student_id = %s, want %s", gotStudentID, studentID)
			}
			className := "Class A"
			return &model.StudentProfile{Student: model.Student{StudentID: studentID, SchoolID: schoolID, FullName: "Student A", CurrentClassName: &className}}, nil
		}}}

		r := gin.New()
		r.Use(withStudentAdminScope(schoolID))
		r.GET("/students/:student_id", h.GetProfile)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/students/"+studentID.String(), nil)
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
		{name: "student not found", err: service.ErrFailedToGetStudent, wantStatus: http.StatusNotFound, wantError: "student not found"},
		{name: "internal", err: errors.New("boom"), wantStatus: http.StatusInternalServerError, wantError: "failed to fetch student profile"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &StudentHandler{studentService: &fakeStudentService{getProfileFn: func(context.Context, *uuid.UUID, uuid.UUID) (*model.StudentProfile, error) {
				return nil, tt.err
			}}}

			r := gin.New()
			r.Use(withStudentAdminScope(schoolID))
			r.GET("/students/:student_id", h.GetProfile)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/students/"+studentID.String(), nil)
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
