package parentcodehandlers

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
	"github.com/hoanghonghuy/iris-app/apps/api/internal/service"
)

type fakeParentCodeService struct {
	generateCodeForStudentFn   func(context.Context, *uuid.UUID, uuid.UUID) (string, error)
	revokeCodeFn               func(context.Context, *uuid.UUID, uuid.UUID) error
	registerParentWithGoogleFn func(context.Context, string, string) (*service.LoginResponse, error)
	registerParentFn           func(context.Context, string, string, string) (*service.LoginResponse, error)
	verifyCodeFn               func(context.Context, string) (*model.StudentParentCode, error)
	getStudentInfoFn           func(context.Context, uuid.UUID) (*model.Student, error)
}

func (f *fakeParentCodeService) GenerateCodeForStudent(ctx context.Context, adminSchoolID *uuid.UUID, studentID uuid.UUID) (string, error) {
	if f.generateCodeForStudentFn == nil {
		return "", errors.New("unexpected GenerateCodeForStudent call")
	}
	return f.generateCodeForStudentFn(ctx, adminSchoolID, studentID)
}

func (f *fakeParentCodeService) RevokeCode(ctx context.Context, adminSchoolID *uuid.UUID, studentID uuid.UUID) error {
	if f.revokeCodeFn == nil {
		return errors.New("unexpected RevokeCode call")
	}
	return f.revokeCodeFn(ctx, adminSchoolID, studentID)
}

func (f *fakeParentCodeService) RegisterParentWithGoogle(ctx context.Context, idToken string, parentCode string) (*service.LoginResponse, error) {
	if f.registerParentWithGoogleFn == nil {
		return nil, errors.New("unexpected RegisterParentWithGoogle call")
	}
	return f.registerParentWithGoogleFn(ctx, idToken, parentCode)
}

func (f *fakeParentCodeService) RegisterParent(ctx context.Context, email, password, parentCode string) (*service.LoginResponse, error) {
	if f.registerParentFn == nil {
		return nil, errors.New("unexpected RegisterParent call")
	}
	return f.registerParentFn(ctx, email, password, parentCode)
}

func (f *fakeParentCodeService) VerifyCode(ctx context.Context, code string) (*model.StudentParentCode, error) {
	if f.verifyCodeFn == nil {
		return nil, errors.New("unexpected VerifyCode call")
	}
	return f.verifyCodeFn(ctx, code)
}

func (f *fakeParentCodeService) GetStudentInfo(ctx context.Context, studentID uuid.UUID) (*model.Student, error) {
	if f.getStudentInfoFn == nil {
		return nil, errors.New("unexpected GetStudentInfo call")
	}
	return f.getStudentInfoFn(ctx, studentID)
}

func withParentCodeAdminScope(schoolID uuid.UUID) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(middleware.CtxAdminSchoolID, schoolID.String())
		c.Next()
	}
}

func decodeParentCodeError(t *testing.T, rec *httptest.ResponseRecorder) string {
	t.Helper()
	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	return body["error"]
}
