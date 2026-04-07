package teacherhandlers

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

type fakeTeacherService struct {
	listFn                func(context.Context, *uuid.UUID, int, int) ([]model.Teacher, int, error)
	listTeachersOfClassFn func(context.Context, *uuid.UUID, uuid.UUID) ([]model.Teacher, error)
	getByTeacherIDFn      func(context.Context, *uuid.UUID, uuid.UUID) (*model.Teacher, error)
	assignFn              func(context.Context, *uuid.UUID, uuid.UUID, uuid.UUID) error
	unassignFn            func(context.Context, *uuid.UUID, uuid.UUID, uuid.UUID) error
	updateFn              func(context.Context, *uuid.UUID, uuid.UUID, string, string, uuid.UUID) error
	deleteFn              func(context.Context, *uuid.UUID, uuid.UUID) error
}

func (f *fakeTeacherService) List(ctx context.Context, adminSchoolID *uuid.UUID, limit, offset int) ([]model.Teacher, int, error) {
	if f.listFn == nil {
		return nil, 0, errors.New("unexpected List call")
	}
	return f.listFn(ctx, adminSchoolID, limit, offset)
}

func (f *fakeTeacherService) ListTeachersOfClass(ctx context.Context, adminSchoolID *uuid.UUID, classID uuid.UUID) ([]model.Teacher, error) {
	if f.listTeachersOfClassFn == nil {
		return nil, errors.New("unexpected ListTeachersOfClass call")
	}
	return f.listTeachersOfClassFn(ctx, adminSchoolID, classID)
}

func (f *fakeTeacherService) GetByTeacherID(ctx context.Context, adminSchoolID *uuid.UUID, teacherID uuid.UUID) (*model.Teacher, error) {
	if f.getByTeacherIDFn == nil {
		return nil, errors.New("unexpected GetByTeacherID call")
	}
	return f.getByTeacherIDFn(ctx, adminSchoolID, teacherID)
}

func (f *fakeTeacherService) Assign(ctx context.Context, adminSchoolID *uuid.UUID, teacherID, classID uuid.UUID) error {
	if f.assignFn == nil {
		return errors.New("unexpected Assign call")
	}
	return f.assignFn(ctx, adminSchoolID, teacherID, classID)
}

func (f *fakeTeacherService) Unassign(ctx context.Context, adminSchoolID *uuid.UUID, teacherID, classID uuid.UUID) error {
	if f.unassignFn == nil {
		return errors.New("unexpected Unassign call")
	}
	return f.unassignFn(ctx, adminSchoolID, teacherID, classID)
}

func (f *fakeTeacherService) Update(ctx context.Context, adminSchoolID *uuid.UUID, teacherID uuid.UUID, fullName, phone string, schoolID uuid.UUID) error {
	if f.updateFn == nil {
		return errors.New("unexpected Update call")
	}
	return f.updateFn(ctx, adminSchoolID, teacherID, fullName, phone, schoolID)
}

func (f *fakeTeacherService) Delete(ctx context.Context, adminSchoolID *uuid.UUID, teacherID uuid.UUID) error {
	if f.deleteFn == nil {
		return errors.New("unexpected Delete call")
	}
	return f.deleteFn(ctx, adminSchoolID, teacherID)
}

func withAdminSchoolScope(schoolID uuid.UUID) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(middleware.CtxAdminSchoolID, schoolID.String())
		c.Next()
	}
}

func decodeTeacherError(t *testing.T, rec *httptest.ResponseRecorder) string {
	t.Helper()
	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	return body["error"]
}
