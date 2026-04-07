package shared

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/middleware"
)

func TestExtractAdminSchoolID(t *testing.T) {
	validSchoolID := uuid.New()

	tests := []struct {
		name  string
		setup func(*gin.Context)
		want  *uuid.UUID
	}{
		{
			name:  "missing context value",
			setup: func(*gin.Context) {},
			want:  nil,
		},
		{
			name: "non-string context value",
			setup: func(c *gin.Context) {
				c.Set(middleware.CtxAdminSchoolID, 123)
			},
			want: nil,
		},
		{
			name: "invalid uuid string",
			setup: func(c *gin.Context) {
				c.Set(middleware.CtxAdminSchoolID, "not-a-uuid")
			},
			want: nil,
		},
		{
			name: "valid uuid string",
			setup: func(c *gin.Context) {
				c.Set(middleware.CtxAdminSchoolID, validSchoolID.String())
			},
			want: &validSchoolID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := gin.CreateTestContext(nil)
			tt.setup(c)

			got := ExtractAdminSchoolID(c)
			if tt.want == nil {
				if got != nil {
					t.Fatalf("expected nil school id, got %v", *got)
				}
				return
			}

			if got == nil {
				t.Fatalf("expected school id %s, got nil", tt.want.String())
			}
			if *got != *tt.want {
				t.Fatalf("school id = %s, want %s", got.String(), tt.want.String())
			}
		})
	}
}
