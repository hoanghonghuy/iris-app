package classhandlers

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

type fakeClassService struct {
	createFn       func(context.Context, *uuid.UUID, uuid.UUID, string, string) (*model.Class, error)
	listBySchoolFn func(context.Context, *uuid.UUID, uuid.UUID, int, int) ([]model.Class, int, error)
	updateFn       func(context.Context, *uuid.UUID, uuid.UUID, string, string) error
	deleteFn       func(context.Context, *uuid.UUID, uuid.UUID) error
}

func (f *fakeClassService) Create(ctx context.Context, adminSchoolID *uuid.UUID, schoolID uuid.UUID, name, schoolYear string) (*model.Class, error) {
	if f.createFn == nil {
		return nil, errors.New("unexpected Create call")
	}
	return f.createFn(ctx, adminSchoolID, schoolID, name, schoolYear)
}

func (f *fakeClassService) ListBySchool(ctx context.Context, adminSchoolID *uuid.UUID, schoolID uuid.UUID, limit, offset int) ([]model.Class, int, error) {
	if f.listBySchoolFn == nil {
		return nil, 0, errors.New("unexpected ListBySchool call")
	}
	return f.listBySchoolFn(ctx, adminSchoolID, schoolID, limit, offset)
}

func (f *fakeClassService) Update(ctx context.Context, adminSchoolID *uuid.UUID, classID uuid.UUID, name, schoolYear string) error {
	if f.updateFn == nil {
		return errors.New("unexpected Update call")
	}
	return f.updateFn(ctx, adminSchoolID, classID, name, schoolYear)
}

func (f *fakeClassService) Delete(ctx context.Context, adminSchoolID *uuid.UUID, classID uuid.UUID) error {
	if f.deleteFn == nil {
		return errors.New("unexpected Delete call")
	}
	return f.deleteFn(ctx, adminSchoolID, classID)
}

func withClassAdminScope(schoolID uuid.UUID) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(middleware.CtxAdminSchoolID, schoolID.String())
		c.Next()
	}
}

func decodeClassError(t *testing.T, rec *httptest.ResponseRecorder) string {
	t.Helper()
	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	return body["error"]
}
