package schoolhandlers

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

type fakeSchoolService struct {
	createFn func(context.Context, string, string) (*model.School, error)
	listFn   func(context.Context, *uuid.UUID, int, int) ([]model.School, int, error)
	updateFn func(context.Context, uuid.UUID, string, string) error
	deleteFn func(context.Context, uuid.UUID) error
}

func (f *fakeSchoolService) Create(ctx context.Context, name, address string) (*model.School, error) {
	if f.createFn == nil {
		return nil, errors.New("unexpected Create call")
	}
	return f.createFn(ctx, name, address)
}

func (f *fakeSchoolService) List(ctx context.Context, adminSchoolID *uuid.UUID, limit, offset int) ([]model.School, int, error) {
	if f.listFn == nil {
		return nil, 0, errors.New("unexpected List call")
	}
	return f.listFn(ctx, adminSchoolID, limit, offset)
}

func (f *fakeSchoolService) Update(ctx context.Context, schoolID uuid.UUID, name, address string) error {
	if f.updateFn == nil {
		return errors.New("unexpected Update call")
	}
	return f.updateFn(ctx, schoolID, name, address)
}

func (f *fakeSchoolService) Delete(ctx context.Context, schoolID uuid.UUID) error {
	if f.deleteFn == nil {
		return errors.New("unexpected Delete call")
	}
	return f.deleteFn(ctx, schoolID)
}

func withSchoolAdminScope(schoolID uuid.UUID) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(middleware.CtxAdminSchoolID, schoolID.String())
		c.Next()
	}
}

func decodeSchoolError(t *testing.T, rec *httptest.ResponseRecorder) string {
	t.Helper()
	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	return body["error"]
}
