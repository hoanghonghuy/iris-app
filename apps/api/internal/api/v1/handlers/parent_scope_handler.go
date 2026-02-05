package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/auth"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/middleware"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/response"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/service"
)

type ParentScopeHandler struct {
	parentScopeService *service.ParentScopeService
}

func NewParentScopeHandler(parentScopeService *service.ParentScopeService) *ParentScopeHandler {
	return &ParentScopeHandler{
		parentScopeService: parentScopeService,
	}
}

type ListChildPostsRequest struct {
	Limit  int `form:"limit"`
	Offset int `form:"offset"`
}

type GetMyFeedRequest struct {
	Limit  int `form:"limit"`
	Offset int `form:"offset"`
}

// MyChildren trả về danh sách các học sinh (con) của phụ huynh
func (h *ParentScopeHandler) MyChildren(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	// Get userID from JWT claims
	claimsAny, exists := c.Get(middleware.CtxClaims)
	if !exists {
		response.Fail(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	claims := claimsAny.(*auth.Claims)

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid user ID")
		return
	}

	students, err := h.parentScopeService.ListMyChildren(ctx, userID)
	if err != nil {
		if err == service.ErrInvalidUserID {
			response.Fail(c, http.StatusBadRequest, "invalid user ID")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "failed to fetch children")
		return
	}

	response.OK(c, students)
}

// ListMyChildClassPosts liệt kê bài đăng của lớp con mình đang học
func (h *ParentScopeHandler) ListMyChildClassPosts(c *gin.Context) {
	studentID, err := uuid.Parse(c.Param("student_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid student_id")
		return
	}

	var req ListChildPostsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid query parameters")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	// Get userID from JWT claims
	claimsAny, exists := c.Get(middleware.CtxClaims)
	if !exists {
		response.Fail(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	claims := claimsAny.(*auth.Claims)

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid user ID")
		return
	}

	posts, err := h.parentScopeService.ListMyChildClassPosts(ctx, userID, studentID, req.Limit, req.Offset)
	if err != nil {
		if err == service.ErrInvalidUserID {
			response.Fail(c, http.StatusBadRequest, err.Error())
			return
		}
		if err == service.ErrForbidden {
			response.Fail(c, http.StatusForbidden, "forbidden: you can only view posts for your own children")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "failed to fetch class posts")
		return
	}

	response.OK(c, posts)
}

// ListMyChildStudentPosts liệt kê bài đăng riêng của con mình (student scope)
func (h *ParentScopeHandler) ListMyChildStudentPosts(c *gin.Context) {
	studentID, err := uuid.Parse(c.Param("student_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid student_id")
		return
	}

	var req ListChildPostsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid query parameters")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	// Get userID from JWT claims
	claimsAny, exists := c.Get(middleware.CtxClaims)
	if !exists {
		response.Fail(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	claims := claimsAny.(*auth.Claims)

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid user ID")
		return
	}

	posts, err := h.parentScopeService.ListMyChildStudentPosts(ctx, userID, studentID, req.Limit, req.Offset)
	if err != nil {
		if err == service.ErrInvalidUserID {
			response.Fail(c, http.StatusBadRequest, err.Error())
			return
		}
		if err == service.ErrForbidden {
			response.Fail(c, http.StatusForbidden, "forbidden: you can only view posts for your own children")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "failed to fetch student posts")
		return
	}

	response.OK(c, posts)
}

// ListAllMyChildPosts liệt kê tất cả bài đăng liên quan đến con mình (cả class và student scope)
func (h *ParentScopeHandler) ListAllMyChildPosts(c *gin.Context) {
	studentID, err := uuid.Parse(c.Param("student_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid student_id")
		return
	}

	var req ListChildPostsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid query parameters")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	// Get userID from JWT claims
	claimsAny, exists := c.Get(middleware.CtxClaims)
	if !exists {
		response.Fail(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	claims := claimsAny.(*auth.Claims)

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid user ID")
		return
	}

	posts, err := h.parentScopeService.ListAllMyChildPosts(ctx, userID, studentID, req.Limit, req.Offset)
	if err != nil {
		if err == service.ErrInvalidUserID {
			response.Fail(c, http.StatusBadRequest, err.Error())
			return
		}
		if err == service.ErrForbidden {
			response.Fail(c, http.StatusForbidden, "forbidden: you can only view posts for your own children")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "failed to fetch posts")
		return
	}

	response.OK(c, posts)
}

func (h *ParentScopeHandler) GetMyFeed(c *gin.Context) {
	var req GetMyFeedRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid query parameters")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	// Get userID from JWT claims
	claimsAny, exists := c.Get(middleware.CtxClaims)
	if !exists {
		response.Fail(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	claims := claimsAny.(*auth.Claims)

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid user ID")
		return
	}

	posts, err := h.parentScopeService.GetMyFeed(ctx, userID, req.Limit, req.Offset)
	if err != nil {
		if err == service.ErrInvalidUserID {
			response.Fail(c, http.StatusBadRequest, err.Error())
			return
		}
		response.Fail(c, http.StatusInternalServerError, "failed to fetch feed")
		return
	}

	response.OK(c, posts)
}
