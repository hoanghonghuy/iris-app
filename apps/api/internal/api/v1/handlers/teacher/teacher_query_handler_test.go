package teacherhandlers

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

func TestList_InvalidPaginationParams(t *testing.T) {
	gin.SetMode(gin.TestMode)

	h := &TeacherHandler{}
	r := gin.New()
	r.GET("/teachers", h.List)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/teachers?limit=0&offset=-1", nil)
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestList_SuccessAndErrorMapping(t *testing.T) {
	gin.SetMode(gin.TestMode)
	schoolID := uuid.New()

	t.Run("success", func(t *testing.T) {
		h := &TeacherHandler{teacherService: &fakeTeacherService{listFn: func(_ context.Context, adminSchoolID *uuid.UUID, limit, offset int) ([]model.Teacher, int, error) {
			if adminSchoolID == nil || *adminSchoolID != schoolID {
				t.Fatalf("expected admin school scope to be forwarded")
			}
			if limit != 25 || offset != 3 {
				t.Fatalf("forwarded limit/offset = %d/%d, want 25/3", limit, offset)
			}
			return []model.Teacher{{TeacherID: uuid.New(), FullName: "Teacher A"}}, 1, nil
		}}}

		r := gin.New()
		r.Use(withAdminSchoolScope(schoolID))
		r.GET("/teachers", h.List)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/teachers?limit=25&offset=3", nil)
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
		{name: "internal", err: errors.New("boom"), wantStatus: http.StatusInternalServerError, wantError: "failed to fetch teachers"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &TeacherHandler{teacherService: &fakeTeacherService{listFn: func(context.Context, *uuid.UUID, int, int) ([]model.Teacher, int, error) {
				return nil, 0, tt.err
			}}}

			r := gin.New()
			r.Use(withAdminSchoolScope(schoolID))
			r.GET("/teachers", h.List)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/teachers?limit=20&offset=0", nil)
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

func TestListTeachersOfClass_SuccessAndErrorMapping(t *testing.T) {
	gin.SetMode(gin.TestMode)
	schoolID := uuid.New()
	classID := uuid.New()

	t.Run("invalid class id", func(t *testing.T) {
		h := &TeacherHandler{}
		r := gin.New()
		r.GET("/teachers/class/:class_id", h.ListTeachersOfClass)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/teachers/class/not-a-uuid", nil)
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
	})

	t.Run("success", func(t *testing.T) {
		h := &TeacherHandler{teacherService: &fakeTeacherService{listTeachersOfClassFn: func(_ context.Context, adminSchoolID *uuid.UUID, gotClassID uuid.UUID) ([]model.Teacher, error) {
			if adminSchoolID == nil || *adminSchoolID != schoolID {
				t.Fatalf("expected admin school scope to be forwarded")
			}
			if gotClassID != classID {
				t.Fatalf("class_id = %s, want %s", gotClassID, classID)
			}
			return []model.Teacher{{TeacherID: uuid.New(), FullName: "Teacher A"}}, nil
		}}}

		r := gin.New()
		r.Use(withAdminSchoolScope(schoolID))
		r.GET("/teachers/class/:class_id", h.ListTeachersOfClass)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/teachers/class/"+classID.String(), nil)
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
		{name: "internal", err: errors.New("boom"), wantStatus: http.StatusInternalServerError, wantError: "failed to fetch teachers of class"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &TeacherHandler{teacherService: &fakeTeacherService{listTeachersOfClassFn: func(context.Context, *uuid.UUID, uuid.UUID) ([]model.Teacher, error) {
				return nil, tt.err
			}}}
			r := gin.New()
			r.Use(withAdminSchoolScope(schoolID))
			r.GET("/teachers/class/:class_id", h.ListTeachersOfClass)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/teachers/class/"+classID.String(), nil)
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

func TestGetByTeacherID_SuccessAndErrorMapping(t *testing.T) {
	gin.SetMode(gin.TestMode)
	schoolID := uuid.New()
	teacherID := uuid.New()

	t.Run("invalid teacher id", func(t *testing.T) {
		h := &TeacherHandler{}
		r := gin.New()
		r.GET("/teachers/:teacher_id", h.GetByTeacherID)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/teachers/not-a-uuid", nil)
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
	})

	t.Run("success", func(t *testing.T) {
		h := &TeacherHandler{teacherService: &fakeTeacherService{getByTeacherIDFn: func(_ context.Context, adminSchoolID *uuid.UUID, gotTeacherID uuid.UUID) (*model.Teacher, error) {
			if adminSchoolID == nil || *adminSchoolID != schoolID {
				t.Fatalf("expected admin school scope to be forwarded")
			}
			if gotTeacherID != teacherID {
				t.Fatalf("teacher_id = %s, want %s", gotTeacherID, teacherID)
			}
			return &model.Teacher{TeacherID: teacherID, FullName: "Teacher A", SchoolID: schoolID}, nil
		}}}

		r := gin.New()
		r.Use(withAdminSchoolScope(schoolID))
		r.GET("/teachers/:teacher_id", h.GetByTeacherID)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/teachers/"+teacherID.String(), nil)
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
		{name: "not found", err: service.ErrTeacherNotFound, wantStatus: http.StatusNotFound, wantError: "teacher not found"},
		{name: "internal", err: errors.New("boom"), wantStatus: http.StatusInternalServerError, wantError: "failed to fetch teacher"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &TeacherHandler{teacherService: &fakeTeacherService{getByTeacherIDFn: func(context.Context, *uuid.UUID, uuid.UUID) (*model.Teacher, error) {
				return nil, tt.err
			}}}
			r := gin.New()
			r.Use(withAdminSchoolScope(schoolID))
			r.GET("/teachers/:teacher_id", h.GetByTeacherID)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/teachers/"+teacherID.String(), nil)
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
