package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/auth"
)

func TestAuthJWT(t *testing.T) {
	gin.SetMode(gin.TestMode)

	secret := "jwt-test-secret"
	validToken, err := auth.Sign(secret, time.Minute, "user-1", "user-1@example.com", []string{"PARENT"}, "")
	if err != nil {
		t.Fatalf("Sign() error = %v", err)
	}

	tests := []struct {
		name       string
		authHeader string
		wantStatus int
		wantError  string
	}{
		{name: "missing bearer token", authHeader: "", wantStatus: http.StatusUnauthorized, wantError: "missing bearer token"},
		{name: "invalid header format", authHeader: "Token abc", wantStatus: http.StatusUnauthorized, wantError: "missing bearer token"},
		{name: "invalid token", authHeader: "Bearer invalid-token", wantStatus: http.StatusUnauthorized, wantError: "invalid token"},
		{name: "valid token", authHeader: "Bearer " + validToken, wantStatus: http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			r.Use(AuthJWT(secret))
			r.GET("/private", func(c *gin.Context) {
				v, _ := c.Get(CtxClaims)
				claims := v.(*auth.Claims)
				c.JSON(http.StatusOK, gin.H{"user_id": claims.UserID})
			})

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/private", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			r.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tt.wantStatus)
			}

			if tt.wantError != "" {
				errorMsg := decodeErrorMessage(t, rec)
				if errorMsg != tt.wantError {
					t.Fatalf("error = %q, want %q", errorMsg, tt.wantError)
				}
				return
			}

			var body map[string]string
			if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
				t.Fatalf("json.Unmarshal() error = %v", err)
			}
			if body["user_id"] != "user-1" {
				t.Fatalf("user_id = %q, want %q", body["user_id"], "user-1")
			}
		})
	}
}

func TestRequireRole(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		claims     *auth.Claims
		wantStatus int
		wantError  string
	}{
		{name: "missing claims", claims: nil, wantStatus: http.StatusUnauthorized, wantError: "unauthorized"},
		{name: "role mismatch", claims: &auth.Claims{Roles: []string{"PARENT"}}, wantStatus: http.StatusForbidden, wantError: "access denied"},
		{name: "role matched", claims: &auth.Claims{Roles: []string{"SUPER_ADMIN"}}, wantStatus: http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			r.Use(setClaims(tt.claims))
			r.Use(RequireRole("SUPER_ADMIN"))
			r.GET("/admin", func(c *gin.Context) { c.Status(http.StatusOK) })

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/admin", nil)
			r.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if tt.wantError != "" {
				errorMsg := decodeErrorMessage(t, rec)
				if errorMsg != tt.wantError {
					t.Fatalf("error = %q, want %q", errorMsg, tt.wantError)
				}
			}
		})
	}
}

func TestRequireAnyRole(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		claims     *auth.Claims
		wantStatus int
		wantError  string
	}{
		{name: "missing claims", claims: nil, wantStatus: http.StatusUnauthorized, wantError: "unauthorized"},
		{name: "none of roles matched", claims: &auth.Claims{Roles: []string{"PARENT"}}, wantStatus: http.StatusForbidden, wantError: "access denied"},
		{name: "one role matched", claims: &auth.Claims{Roles: []string{"TEACHER", "SCHOOL_ADMIN"}}, wantStatus: http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			r.Use(setClaims(tt.claims))
			r.Use(RequireAnyRole("SUPER_ADMIN", "SCHOOL_ADMIN"))
			r.GET("/admin", func(c *gin.Context) { c.Status(http.StatusOK) })

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/admin", nil)
			r.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if tt.wantError != "" {
				errorMsg := decodeErrorMessage(t, rec)
				if errorMsg != tt.wantError {
					t.Fatalf("error = %q, want %q", errorMsg, tt.wantError)
				}
			}
		})
	}
}

func TestInjectAdminScope(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name          string
		claims        *auth.Claims
		wantStatus    int
		wantError     string
		wantCtxSchool string
	}{
		{name: "missing claims", claims: nil, wantStatus: http.StatusUnauthorized, wantError: "unauthorized"},
		{
			name:       "school admin without school id",
			claims:     &auth.Claims{Roles: []string{"SCHOOL_ADMIN"}},
			wantStatus: http.StatusForbidden,
			wantError:  "school admin account not linked to any school",
		},
		{
			name:          "school admin with school id",
			claims:        &auth.Claims{Roles: []string{"SCHOOL_ADMIN"}, SchoolID: "school-1"},
			wantStatus:    http.StatusOK,
			wantCtxSchool: "school-1",
		},
		{
			name:          "super admin with empty school id",
			claims:        &auth.Claims{Roles: []string{"SUPER_ADMIN"}},
			wantStatus:    http.StatusOK,
			wantCtxSchool: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			r.Use(setClaims(tt.claims))
			r.Use(InjectAdminScope())
			r.GET("/admin", func(c *gin.Context) {
				v, _ := c.Get(CtxAdminSchoolID)
				schoolID, _ := v.(string)
				c.JSON(http.StatusOK, gin.H{"school_id": schoolID})
			})

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/admin", nil)
			r.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tt.wantStatus)
			}

			if tt.wantError != "" {
				errorMsg := decodeErrorMessage(t, rec)
				if errorMsg != tt.wantError {
					t.Fatalf("error = %q, want %q", errorMsg, tt.wantError)
				}
				return
			}

			var body map[string]string
			if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
				t.Fatalf("json.Unmarshal() error = %v", err)
			}
			if body["school_id"] != tt.wantCtxSchool {
				t.Fatalf("school_id = %q, want %q", body["school_id"], tt.wantCtxSchool)
			}
		})
	}
}

func setClaims(claims *auth.Claims) gin.HandlerFunc {
	return func(c *gin.Context) {
		if claims != nil {
			c.Set(CtxClaims, claims)
		}
		c.Next()
	}
}

func decodeErrorMessage(t *testing.T, rec *httptest.ResponseRecorder) string {
	t.Helper()
	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	return body["error"]
}
