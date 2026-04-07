package teacherscope

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/api/v1/handlers/shared"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/response"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/service"
)

// CreatePost tạo bài đăng mới cho lớp hoặc học sinh
func (h *TeacherScopeHandler) CreatePost(c *gin.Context) {
	var req CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid request body")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	// Get userID from JWT claims
	userID, ok := shared.RequireCurrentUserID(c)
	if !ok {
		return
	}

	var postID uuid.UUID

	switch req.ScopeType {
	case "class":
		classID, err := uuid.Parse(req.ClassID)
		if err != nil {
			response.Fail(c, http.StatusBadRequest, "invalid class_id")
			return
		}
		postID, err = h.teacherScopeService.CreateClassPost(ctx, userID, classID, req.Type, req.Content)
		if err != nil {
			if errors.Is(err, service.ErrInvalidUserID) || errors.Is(err, service.ErrInvalidClassID) || errors.Is(err, service.ErrInvalidValue) {
				response.Fail(c, http.StatusBadRequest, err.Error())
				return
			}
			if errors.Is(err, service.ErrForbidden) {
				response.Fail(c, http.StatusForbidden, "forbidden: you can only create posts for your assigned classes")
				return
			}
			response.Fail(c, http.StatusInternalServerError, "failed to create post")
			return
		}

	case "student":
		studentID, err := uuid.Parse(req.StudentID)
		if err != nil {
			response.Fail(c, http.StatusBadRequest, "invalid student_id")
			return
		}
		postID, err = h.teacherScopeService.CreateStudentPost(ctx, userID, studentID, req.Type, req.Content)
		if err != nil {
			if errors.Is(err, service.ErrInvalidUserID) || errors.Is(err, service.ErrInvalidValue) {
				response.Fail(c, http.StatusBadRequest, err.Error())
				return
			}
			if errors.Is(err, service.ErrForbidden) {
				response.Fail(c, http.StatusForbidden, "forbidden: you can only create posts for students in your assigned classes")
				return
			}
			response.Fail(c, http.StatusInternalServerError, "failed to create post")
			return
		}

	default:
		response.Fail(c, http.StatusBadRequest, "invalid scope_type (class|student)")
		return
	}

	response.Created(c, gin.H{
		"message":    "post created successfully",
		"post_id":    postID.String(),
		"scope_type": req.ScopeType,
		"type":       req.Type,
	})
}

// UpdatePost cập nhật nội dung bài đăng của chính giáo viên.
func (h *TeacherScopeHandler) UpdatePost(c *gin.Context) {
	postID, err := uuid.Parse(c.Param("post_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid post_id")
		return
	}

	var req UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid request body")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	userID, ok := shared.RequireCurrentUserID(c)
	if !ok {
		return
	}

	err = h.teacherScopeService.UpdatePost(ctx, userID, postID, req.Content)
	if err != nil {
		if errors.Is(err, service.ErrInvalidUserID) || errors.Is(err, service.ErrInvalidValue) {
			response.Fail(c, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, service.ErrForbidden) {
			response.Fail(c, http.StatusForbidden, "forbidden: you can only edit your own posts")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "failed to update post")
		return
	}

	response.OK(c, gin.H{
		"message": "post updated successfully",
		"post_id": postID.String(),
	})
}

// DeletePost xóa bài đăng của chính giáo viên.
func (h *TeacherScopeHandler) DeletePost(c *gin.Context) {
	postID, err := uuid.Parse(c.Param("post_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid post_id")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	userID, ok := shared.RequireCurrentUserID(c)
	if !ok {
		return
	}

	err = h.teacherScopeService.DeletePost(ctx, userID, postID)
	if err != nil {
		if errors.Is(err, service.ErrInvalidUserID) || errors.Is(err, service.ErrInvalidValue) {
			response.Fail(c, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, service.ErrForbidden) {
			response.Fail(c, http.StatusForbidden, "forbidden: you can only delete your own posts")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "failed to delete post")
		return
	}

	response.OK(c, gin.H{
		"message": "post deleted successfully",
		"post_id": postID.String(),
	})
}

// ListClassPosts liệt kê bài đăng của một lớp học
func (h *TeacherScopeHandler) ListClassPosts(c *gin.Context) {
	classID, err := uuid.Parse(c.Param("class_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid class_id")
		return
	}

	var req shared.PaginationParams
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid query parameters")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	// Get userID from JWT claims
	userID, ok := shared.RequireCurrentUserID(c)
	if !ok {
		return
	}

	posts, total, err := h.teacherScopeService.ListClassPosts(ctx, userID, classID, req.Limit, req.Offset)
	if err != nil {
		if errors.Is(err, service.ErrInvalidUserID) || errors.Is(err, service.ErrInvalidClassID) {
			response.Fail(c, http.StatusBadRequest, "invalid user ID or class ID")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "failed to list class posts")
		return
	}

	response.OKPaginated(c, posts, response.Pagination{
		Total:   total,
		Limit:   req.Limit,
		Offset:  req.Offset,
		HasMore: req.Offset+len(posts) < total,
	})
}

// ListStudentPosts liệt kê bài đăng của một học sinh
func (h *TeacherScopeHandler) ListStudentPosts(c *gin.Context) {
	studentID, err := uuid.Parse(c.Param("student_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid student_id")
		return
	}

	var req shared.PaginationParams
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid query parameters")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	// Get userID from JWT claims
	userID, ok := shared.RequireCurrentUserID(c)
	if !ok {
		return
	}

	posts, total, err := h.teacherScopeService.ListStudentPosts(ctx, userID, studentID, req.Limit, req.Offset)
	if err != nil {
		if errors.Is(err, service.ErrInvalidUserID) {
			response.Fail(c, http.StatusBadRequest, "invalid user ID or student ID")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "failed to list student posts")
		return
	}

	response.OKPaginated(c, posts, response.Pagination{
		Total:   total,
		Limit:   req.Limit,
		Offset:  req.Offset,
		HasMore: req.Offset+len(posts) < total,
	})
}

// TogglePostLike bật/tắt like cho bài đăng.
func (h *TeacherScopeHandler) TogglePostLike(c *gin.Context) {
	postID, err := uuid.Parse(c.Param("post_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid post_id")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	userID, ok := shared.RequireCurrentUserID(c)
	if !ok {
		return
	}

	liked, likeCount, err := h.teacherScopeService.TogglePostLike(ctx, userID, postID)
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

// ListPostComments liệt kê bình luận của bài đăng.
func (h *TeacherScopeHandler) ListPostComments(c *gin.Context) {
	postID, err := uuid.Parse(c.Param("post_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid post_id")
		return
	}

	var req shared.PaginationParams
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid query parameters")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	userID, ok := shared.RequireCurrentUserID(c)
	if !ok {
		return
	}

	comments, total, err := h.teacherScopeService.ListPostComments(ctx, userID, postID, req.Limit, req.Offset)
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

// CreatePostComment thêm bình luận cho bài đăng.
func (h *TeacherScopeHandler) CreatePostComment(c *gin.Context) {
	postID, err := uuid.Parse(c.Param("post_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid post_id")
		return
	}

	var req shared.CreatePostCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid request body")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	userID, ok := shared.RequireCurrentUserID(c)
	if !ok {
		return
	}

	comment, commentCount, err := h.teacherScopeService.AddPostComment(ctx, userID, postID, req.Content)
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

// SharePost ghi nhận chia sẻ bài đăng.
func (h *TeacherScopeHandler) SharePost(c *gin.Context) {
	postID, err := uuid.Parse(c.Param("post_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid post_id")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	userID, ok := shared.RequireCurrentUserID(c)
	if !ok {
		return
	}

	shareCount, err := h.teacherScopeService.SharePost(ctx, userID, postID)
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
