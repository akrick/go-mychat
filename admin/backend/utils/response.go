package utils

import (
	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"msg"`
	Data    interface{} `json:"data,omitempty"`
}

// PageResponse 分页响应结构
type PageResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"msg"`
	Data    interface{} `json:"data"`
	Total   int64       `json:"total"`
	Page    int         `json:"page"`
	PageSize int         `json:"page_size"`
}

// Success 成功响应
func Success(c *gin.Context, msg string, data interface{}) {
	c.JSON(200, Response{
		Code:    200,
		Message: msg,
		Data:    data,
	})
}

// Error 错误响应
func Error(c *gin.Context, code int, msg string) {
	c.JSON(code, Response{
		Code:    code,
		Message: msg,
	})
}

// PageSuccess 分页成功响应
func PageSuccess(c *gin.Context, msg string, data interface{}, total int64, page, pageSize int) {
	c.JSON(200, PageResponse{
		Code:    200,
		Message: msg,
		Data:    data,
		Total:   total,
		Page:    page,
		PageSize: pageSize,
	})
}

// BadRequest 400错误
func BadRequest(c *gin.Context, msg string) {
	Error(c, 400, msg)
}

// Unauthorized 401错误
func Unauthorized(c *gin.Context, msg string) {
	Error(c, 401, msg)
}

// Forbidden 403错误
func Forbidden(c *gin.Context, msg string) {
	Error(c, 403, msg)
}

// NotFound 404错误
func NotFound(c *gin.Context, msg string) {
	Error(c, 404, msg)
}

// InternalServerError 500错误
func InternalServerError(c *gin.Context, msg string) {
	Error(c, 500, msg)
}
