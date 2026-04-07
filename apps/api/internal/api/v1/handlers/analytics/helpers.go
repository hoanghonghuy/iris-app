package analyticshandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/auth"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/middleware"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/response"
)

func requireCurrentClaims(c *gin.Context) (*auth.Claims, bool) {
	claimsAny, exists := c.Get(middleware.CtxClaims)
	if !exists {
		response.Fail(c, http.StatusUnauthorized, "unauthorized")
		return nil, false
	}

	claims, ok := claimsAny.(*auth.Claims)
	if !ok {
		response.Fail(c, http.StatusUnauthorized, "unauthorized")
		return nil, false
	}

	return claims, true
}

func requireCurrentUserID(c *gin.Context) (uuid.UUID, bool) {
	claims, ok := requireCurrentClaims(c)
	if !ok {
		return uuid.Nil, false
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid user ID")
		return uuid.Nil, false
	}

	return userID, true
}

func extractAdminSchoolID(c *gin.Context) *uuid.UUID {
	value, exists := c.Get(middleware.CtxAdminSchoolID)
	if !exists {
		return nil
	}

	schoolIDStr, ok := value.(string)
	if !ok || schoolIDStr == "" {
		return nil
	}

	schoolID, err := uuid.Parse(schoolIDStr)
	if err != nil {
		return nil
	}
	return &schoolID
}
