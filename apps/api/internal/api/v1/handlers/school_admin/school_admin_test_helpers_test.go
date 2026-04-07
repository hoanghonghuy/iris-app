package schooladminhandlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
)

type fakeSchoolAdminService struct {
	createFn func(context.Context, uuid.UUID, uuid.UUID, string, string) (*model.SchoolAdmin, error)
	listFn   func(context.Context, int, int) ([]model.SchoolAdmin, int, error)
	deleteFn func(context.Context, uuid.UUID) error
}

func (f *fakeSchoolAdminService) Create(ctx context.Context, userID, schoolID uuid.UUID, fullName, phone string) (*model.SchoolAdmin, error) {
	if f.createFn == nil {
		return nil, errors.New("unexpected Create call")
	}
	return f.createFn(ctx, userID, schoolID, fullName, phone)
}

func (f *fakeSchoolAdminService) List(ctx context.Context, limit, offset int) ([]model.SchoolAdmin, int, error) {
	if f.listFn == nil {
		return nil, 0, errors.New("unexpected List call")
	}
	return f.listFn(ctx, limit, offset)
}

func (f *fakeSchoolAdminService) Delete(ctx context.Context, adminID uuid.UUID) error {
	if f.deleteFn == nil {
		return errors.New("unexpected Delete call")
	}
	return f.deleteFn(ctx, adminID)
}

func decodeSchoolAdminError(t *testing.T, rec *httptest.ResponseRecorder) string {
	t.Helper()
	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	return body["error"]
}
