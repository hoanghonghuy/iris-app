// Package response cung cấp các helper function để tạo HTTP response nhất quán
package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

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
