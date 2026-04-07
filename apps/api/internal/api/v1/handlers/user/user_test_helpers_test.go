package userhandlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/auth"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/middleware"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
)

type fakeUserService struct {
	createUserWithoutPasswordFn func(context.Context, *uuid.UUID, string, []string) (*model.UserInfo, error)
	findByIDFn                  func(context.Context, *uuid.UUID, uuid.UUID) (*model.UserInfo, error)
	listFn                      func(context.Context, *uuid.UUID, string, int, int) ([]model.UserInfo, int, error)
	lockFn                      func(context.Context, *uuid.UUID, uuid.UUID) error
	unlockFn                    func(context.Context, *uuid.UUID, uuid.UUID) error
	assignRoleFn                func(context.Context, uuid.UUID, string) error
	activateUserWithTokenFn     func(context.Context, string, string) error
	updateMyPasswordFn          func(context.Context, uuid.UUID, string) error
	deleteFn                    func(context.Context, uuid.UUID) error
}

func (f *fakeUserService) CreateUserWithoutPassword(ctx context.Context, adminSchoolID *uuid.UUID, email string, roles []string) (*model.UserInfo, error) {
	if f.createUserWithoutPasswordFn == nil {
		return nil, errors.New("unexpected CreateUserWithoutPassword call")
	}
	return f.createUserWithoutPasswordFn(ctx, adminSchoolID, email, roles)
}

func (f *fakeUserService) FindByID(ctx context.Context, adminSchoolID *uuid.UUID, userID uuid.UUID) (*model.UserInfo, error) {
	if f.findByIDFn == nil {
		return nil, errors.New("unexpected FindByID call")
	}
	return f.findByIDFn(ctx, adminSchoolID, userID)
}

func (f *fakeUserService) List(ctx context.Context, adminSchoolID *uuid.UUID, roleFilter string, limit, offset int) ([]model.UserInfo, int, error) {
	if f.listFn == nil {
		return nil, 0, errors.New("unexpected List call")
	}
	return f.listFn(ctx, adminSchoolID, roleFilter, limit, offset)
}

func (f *fakeUserService) Lock(ctx context.Context, adminSchoolID *uuid.UUID, userID uuid.UUID) error {
	if f.lockFn == nil {
		return errors.New("unexpected Lock call")
	}
	return f.lockFn(ctx, adminSchoolID, userID)
}

func (f *fakeUserService) Unlock(ctx context.Context, adminSchoolID *uuid.UUID, userID uuid.UUID) error {
	if f.unlockFn == nil {
		return errors.New("unexpected Unlock call")
	}
	return f.unlockFn(ctx, adminSchoolID, userID)
}

func (f *fakeUserService) AssignRole(ctx context.Context, userID uuid.UUID, roleName string) error {
	if f.assignRoleFn == nil {
		return errors.New("unexpected AssignRole call")
	}
	return f.assignRoleFn(ctx, userID, roleName)
}

func (f *fakeUserService) ActivateUserWithToken(ctx context.Context, token, password string) error {
	if f.activateUserWithTokenFn == nil {
		return errors.New("unexpected ActivateUserWithToken call")
	}
	return f.activateUserWithTokenFn(ctx, token, password)
}

func (f *fakeUserService) UpdateMyPassword(ctx context.Context, userID uuid.UUID, password string) error {
	if f.updateMyPasswordFn == nil {
		return errors.New("unexpected UpdateMyPassword call")
	}
	return f.updateMyPasswordFn(ctx, userID, password)
}

func (f *fakeUserService) Delete(ctx context.Context, userID uuid.UUID) error {
	if f.deleteFn == nil {
		return errors.New("unexpected Delete call")
	}
	return f.deleteFn(ctx, userID)
}

func withUserAdminScope(schoolID uuid.UUID) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(middleware.CtxAdminSchoolID, schoolID.String())
		c.Next()
	}
}

func withUserClaims(userID string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(middleware.CtxClaims, &auth.Claims{UserID: userID})
		c.Next()
	}
}

func decodeUserError(t *testing.T, rec *httptest.ResponseRecorder) string {
	t.Helper()
	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	return body["error"]
}
