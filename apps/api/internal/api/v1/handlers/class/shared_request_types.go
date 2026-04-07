package classhandlers

// PaginationParams input chung cho phân trang trong class endpoints.
type PaginationParams struct {
	Limit  int `form:"limit" binding:"omitempty,min=1,max=100"`
	Offset int `form:"offset" binding:"omitempty,min=0"`
}
