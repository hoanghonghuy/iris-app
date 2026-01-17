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
