package handlers

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/response"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/service"
)

type ParentScopeHandler struct {
	parentScopeService *service.ParentScopeService
	appointmentService *service.AppointmentService
}

func NewParentScopeHandler(parentScopeService *service.ParentScopeService, appointmentService *service.AppointmentService) *ParentScopeHandler {
	return &ParentScopeHandler{
		parentScopeService: parentScopeService,
		appointmentService: appointmentService,
	}
}

// MyChildren trả về danh sách các học sinh (con) của phụ huynh
func (h *ParentScopeHandler) MyChildren(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	userID, ok := requireCurrentUserID(c)
	if !ok {
		return
	}

	students, err := h.parentScopeService.ListMyChildren(ctx, userID)
	if err != nil {
		if errors.Is(err, service.ErrInvalidUserID) {
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

	var req PaginationParams
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid query parameters")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	userID, ok := requireCurrentUserID(c)
	if !ok {
		return
	}

	posts, total, err := h.parentScopeService.ListMyChildClassPosts(ctx, userID, studentID, req.Limit, req.Offset)
	if err != nil {
		if errors.Is(err, service.ErrInvalidUserID) {
			response.Fail(c, http.StatusBadRequest, "invalid user ID")
			return
		}
		if errors.Is(err, service.ErrForbidden) {
			response.Fail(c, http.StatusForbidden, "forbidden")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "failed to fetch class posts")
		return
	}

	response.OKPaginated(c, posts, response.Pagination{
		Total:   total,
		Limit:   req.Limit,
		Offset:  req.Offset,
		HasMore: req.Offset+len(posts) < total,
	})
}

// ListMyChildStudentPosts liệt kê bài đăng riêng của con mình (student scope)
func (h *ParentScopeHandler) ListMyChildStudentPosts(c *gin.Context) {
	studentID, err := uuid.Parse(c.Param("student_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid student_id")
		return
	}

	var req PaginationParams
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid query parameters")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	userID, ok := requireCurrentUserID(c)
	if !ok {
		return
	}

	posts, total, err := h.parentScopeService.ListMyChildStudentPosts(ctx, userID, studentID, req.Limit, req.Offset)
	if err != nil {
		if errors.Is(err, service.ErrInvalidUserID) {
			response.Fail(c, http.StatusBadRequest, "invalid user ID")
			return
		}
		if errors.Is(err, service.ErrForbidden) {
			response.Fail(c, http.StatusForbidden, "forbidden")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "failed to fetch student posts")
		return
	}

	response.OKPaginated(c, posts, response.Pagination{
		Total:   total,
		Limit:   req.Limit,
		Offset:  req.Offset,
		HasMore: req.Offset+len(posts) < total,
	})
}

// ListAllMyChildPosts liệt kê tất cả bài đăng liên quan đến con mình (cả class và student scope)
func (h *ParentScopeHandler) ListAllMyChildPosts(c *gin.Context) {
	studentID, err := uuid.Parse(c.Param("student_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid student_id")
		return
	}

	var req PaginationParams
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid query parameters")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	userID, ok := requireCurrentUserID(c)
	if !ok {
		return
	}

	posts, total, err := h.parentScopeService.ListAllMyChildPosts(ctx, userID, studentID, req.Limit, req.Offset)
	if err != nil {
		if errors.Is(err, service.ErrInvalidUserID) {
			response.Fail(c, http.StatusBadRequest, "invalid user ID")
			return
		}
		if errors.Is(err, service.ErrForbidden) {
			response.Fail(c, http.StatusForbidden, "forbidden")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "failed to fetch posts")
		return
	}

	response.OKPaginated(c, posts, response.Pagination{
		Total:   total,
		Limit:   req.Limit,
		Offset:  req.Offset,
		HasMore: req.Offset+len(posts) < total,
	})
}

func (h *ParentScopeHandler) GetMyFeed(c *gin.Context) {
	var req PaginationParams
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid query parameters")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	userID, ok := requireCurrentUserID(c)
	if !ok {
		return
	}

	posts, total, err := h.parentScopeService.GetMyFeed(ctx, userID, req.Limit, req.Offset)
	if err != nil {
		if errors.Is(err, service.ErrInvalidUserID) {
			response.Fail(c, http.StatusBadRequest, "invalid user ID")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "failed to fetch feed")
		return
	}

	response.OKPaginated(c, posts, response.Pagination{
		Total:   total,
		Limit:   req.Limit,
		Offset:  req.Offset,
		HasMore: req.Offset+len(posts) < total,
	})
}

// TogglePostLike bật/tắt like cho bài đăng trong feed phụ huynh.
func (h *ParentScopeHandler) TogglePostLike(c *gin.Context) {
	postID, err := uuid.Parse(c.Param("post_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid post_id")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	userID, ok := requireCurrentUserID(c)
	if !ok {
		return
	}

	liked, likeCount, err := h.parentScopeService.TogglePostLike(ctx, userID, postID)
	if err != nil {
		if errors.Is(err, service.ErrInvalidUserID) || errors.Is(err, service.ErrInvalidValue) {
			response.Fail(c, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, service.ErrForbidden) {
			response.Fail(c, http.StatusForbidden, "forbidden")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "failed to toggle like")
		return
	}

	response.OK(c, gin.H{
		"post_id":     postID.String(),
		"liked_by_me": liked,
		"like_count":  likeCount,
	})
}

// ListPostComments liệt kê bình luận của một bài đăng trong feed phụ huynh.
func (h *ParentScopeHandler) ListPostComments(c *gin.Context) {
	postID, err := uuid.Parse(c.Param("post_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid post_id")
		return
	}

	var req PaginationParams
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid query parameters")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	userID, ok := requireCurrentUserID(c)
	if !ok {
		return
	}

	comments, total, err := h.parentScopeService.ListPostComments(ctx, userID, postID, req.Limit, req.Offset)
	if err != nil {
		if errors.Is(err, service.ErrInvalidUserID) || errors.Is(err, service.ErrInvalidValue) {
			response.Fail(c, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, service.ErrForbidden) {
			response.Fail(c, http.StatusForbidden, "forbidden")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "failed to list comments")
		return
	}

	response.OKPaginated(c, comments, response.Pagination{
		Total:   total,
		Limit:   req.Limit,
		Offset:  req.Offset,
		HasMore: req.Offset+len(comments) < total,
	})
}

// CreatePostComment thêm bình luận vào bài đăng trong feed phụ huynh.
func (h *ParentScopeHandler) CreatePostComment(c *gin.Context) {
	postID, err := uuid.Parse(c.Param("post_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid post_id")
		return
	}

	var req CreatePostCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid request body")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	userID, ok := requireCurrentUserID(c)
	if !ok {
		return
	}

	comment, commentCount, err := h.parentScopeService.AddPostComment(ctx, userID, postID, req.Content)
	if err != nil {
		if errors.Is(err, service.ErrInvalidUserID) || errors.Is(err, service.ErrInvalidValue) {
			response.Fail(c, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, service.ErrForbidden) {
			response.Fail(c, http.StatusForbidden, "forbidden")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "failed to create comment")
		return
	}

	response.Created(c, gin.H{
		"comment":       comment,
		"comment_count": commentCount,
		"post_id":       postID.String(),
	})
}

// SharePost ghi nhận chia sẻ bài đăng trong feed phụ huynh.
func (h *ParentScopeHandler) SharePost(c *gin.Context) {
	postID, err := uuid.Parse(c.Param("post_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid post_id")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	userID, ok := requireCurrentUserID(c)
	if !ok {
		return
	}

	shareCount, err := h.parentScopeService.SharePost(ctx, userID, postID)
	if err != nil {
		if errors.Is(err, service.ErrInvalidUserID) || errors.Is(err, service.ErrInvalidValue) {
			response.Fail(c, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, service.ErrForbidden) {
			response.Fail(c, http.StatusForbidden, "forbidden")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "failed to share post")
		return
	}

	response.Created(c, gin.H{
		"post_id":     postID.String(),
		"share_count": shareCount,
	})
}
