package parenthandlers

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

type fakeParentService struct {
	listFn          func(context.Context, *uuid.UUID, int, int) ([]model.Parent, int, error)
	getByParentIDFn func(context.Context, *uuid.UUID, uuid.UUID) (*model.Parent, error)
	assignStudentFn func(context.Context, *uuid.UUID, uuid.UUID, uuid.UUID, string) error
	unassignFn      func(context.Context, *uuid.UUID, uuid.UUID, uuid.UUID) error
}

func (f *fakeParentService) List(ctx context.Context, adminSchoolID *uuid.UUID, limit, offset int) ([]model.Parent, int, error) {
	if f.listFn == nil {
		return nil, 0, errors.New("unexpected List call")
	}
	return f.listFn(ctx, adminSchoolID, limit, offset)
}

func (f *fakeParentService) GetByParentID(ctx context.Context, adminSchoolID *uuid.UUID, parentID uuid.UUID) (*model.Parent, error) {
	if f.getByParentIDFn == nil {
		return nil, errors.New("unexpected GetByParentID call")
	}
	return f.getByParentIDFn(ctx, adminSchoolID, parentID)
}

func (f *fakeParentService) AssignStudent(ctx context.Context, adminSchoolID *uuid.UUID, parentID, studentID uuid.UUID, relationship string) error {
	if f.assignStudentFn == nil {
		return errors.New("unexpected AssignStudent call")
	}
	return f.assignStudentFn(ctx, adminSchoolID, parentID, studentID, relationship)
}

func (f *fakeParentService) UnassignStudent(ctx context.Context, adminSchoolID *uuid.UUID, parentID, studentID uuid.UUID) error {
	if f.unassignFn == nil {
		return errors.New("unexpected UnassignStudent call")
	}
	return f.unassignFn(ctx, adminSchoolID, parentID, studentID)
}

func withParentAdminScope(schoolID uuid.UUID) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(middleware.CtxAdminSchoolID, schoolID.String())
		c.Next()
	}
}

func decodeParentError(t *testing.T, rec *httptest.ResponseRecorder) string {
	t.Helper()
	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	return body["error"]
}
