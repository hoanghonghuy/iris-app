package studenthandlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/middleware"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
)

type fakeStudentService struct {
	createFn      func(context.Context, *uuid.UUID, uuid.UUID, uuid.UUID, string, string, string) (*model.Student, error)
	listByClassFn func(context.Context, *uuid.UUID, uuid.UUID, int, int) ([]model.Student, int, error)
	getProfileFn  func(context.Context, *uuid.UUID, uuid.UUID) (*model.StudentProfile, error)
	updateFn      func(context.Context, *uuid.UUID, uuid.UUID, string, string, string) error
	deleteFn      func(context.Context, *uuid.UUID, uuid.UUID) error
}

func (f *fakeStudentService) Create(ctx context.Context, adminSchoolID *uuid.UUID, schoolID, classID uuid.UUID, fullName string, dobStr string, gender string) (*model.Student, error) {
	if f.createFn == nil {
		return nil, errors.New("unexpected Create call")
	}
	return f.createFn(ctx, adminSchoolID, schoolID, classID, fullName, dobStr, gender)
}

func (f *fakeStudentService) ListByClass(ctx context.Context, adminSchoolID *uuid.UUID, classID uuid.UUID, limit, offset int) ([]model.Student, int, error) {
	if f.listByClassFn == nil {
		return nil, 0, errors.New("unexpected ListByClass call")
	}
	return f.listByClassFn(ctx, adminSchoolID, classID, limit, offset)
}

func (f *fakeStudentService) GetProfile(ctx context.Context, adminSchoolID *uuid.UUID, studentID uuid.UUID) (*model.StudentProfile, error) {
	if f.getProfileFn == nil {
		return nil, errors.New("unexpected GetProfile call")
	}
	return f.getProfileFn(ctx, adminSchoolID, studentID)
}

func (f *fakeStudentService) Update(ctx context.Context, adminSchoolID *uuid.UUID, studentID uuid.UUID, fullName, dobStr, gender string) error {
	if f.updateFn == nil {
		return errors.New("unexpected Update call")
	}
	return f.updateFn(ctx, adminSchoolID, studentID, fullName, dobStr, gender)
}

func (f *fakeStudentService) Delete(ctx context.Context, adminSchoolID *uuid.UUID, studentID uuid.UUID) error {
	if f.deleteFn == nil {
		return errors.New("unexpected Delete call")
	}
	return f.deleteFn(ctx, adminSchoolID, studentID)
}

func withStudentAdminScope(schoolID uuid.UUID) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(middleware.CtxAdminSchoolID, schoolID.String())
		c.Next()
	}
}

func decodeStudentError(t *testing.T, rec *httptest.ResponseRecorder) string {
	t.Helper()
	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	return body["error"]
}
