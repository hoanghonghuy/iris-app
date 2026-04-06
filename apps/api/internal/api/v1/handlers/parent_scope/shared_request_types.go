package parentscope

// PaginationParams input chung cho phân trang trong parent scope endpoints.
type PaginationParams struct {
	Limit  int `form:"limit" binding:"omitempty,min=1,max=100"`
	Offset int `form:"offset" binding:"omitempty,min=0"`
}

// CreatePostCommentRequest input chung cho tạo bình luận bài đăng trong parent scope.
type CreatePostCommentRequest struct {
	Content string `json:"content" binding:"required"`
}
