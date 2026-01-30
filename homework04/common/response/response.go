package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Success 成功响应（带数据）
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// SuccessWithMessage 成功响应（带消息）
func SuccessWithMessage(c *gin.Context, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: message,
	})
}

// SuccessWithData 成功响应（带消息和数据）
func SuccessWithData(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: message,
		Data:    data,
	})
}

// Fail 失败响应（使用 AppError）
func Fail(c *gin.Context, err *AppError) {
	c.JSON(err.HttpStatus, Response{
		Code:    err.Code,
		Message: err.Message,
	})
}

// FailWithMessage 失败响应（自定义消息）
func FailWithMessage(c *gin.Context, err *AppError, message string) {
	c.JSON(err.HttpStatus, Response{
		Code:    err.Code,
		Message: message,
	})
}
