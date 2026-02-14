// Package response cung cấp các helper function để tạo HTTP response nhất quán
package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// const (
// 	DefaultErrorType = "about:blank"
// )

// TODO: RFC 7807 standard

// type ProblemDetail struct {
// 	Type     string            `json:"type"`
// 	Title    string            `json:"title`
// 	Status   int               `json:"status"` // http status code
// 	Detail   string            `json:"detail`
// 	Instance string            `json:"instance,omitempty"`
// 	Errors   map[string]string `string,omitempty"`
// }

// OK trả về response HTTP 200 OK với data
func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, gin.H{"data": data})
}

// Created trả về response HTTP 201 Created với data
func Created(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, gin.H{"data": data})
}

// Fail trả về response lỗi với status code và message
func Fail(c *gin.Context, status int, msg string) {
	c.JSON(status, gin.H{"error": msg})
}

// Pagination chứa metadata phân trang cho response
type Pagination struct {
	Total   int  `json:"total"`
	Limit   int  `json:"limit"`
	Offset  int  `json:"offset"`
	HasMore bool `json:"has_more"`
}

// OKPaginated trả về response HTTP 200 OK với data và pagination metadata
func OKPaginated(c *gin.Context, data any, p Pagination) {
	c.JSON(http.StatusOK, gin.H{"data": data, "pagination": p})
}

// TODO: thêm hàm CreatedWithLocation(c *gin.Context, data any, location string)
// để hỗ trợ tiêu chuẩn RFC 7231 cho phản hồi 201 Created kèm theo header Location.
// Hàm này sẽ hữu ích cho tất cả các endpoint POST tạo tài nguyên mới.
// Tham khảo: https://datatracker.ietf.org/doc/html/rfc7231#section-7.1.2
//
// func CreatedWithLocation(c *gin.Context, data any, location string) {
//     c.Header("Location", location)
//     c.JSON(http.StatusCreated, gin.H{"data": data})
// }
//
// hiện tại, tự thêm header Location:
// c.Header("Location", "/api/v1/users/" + userID.String())
// response.Created(c, data)
