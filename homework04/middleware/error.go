package middleware

import (
	"net/http"
	"runtime/debug"

	"awesomeProject/homework04/common/response"
	"awesomeProject/homework04/global"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ErrorHandler 全局错误处理中间件（panic 恢复）
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 记录 panic 错误和堆栈信息
				global.Logger.Error("Panic recovered",
					zap.Any("error", err),
					zap.String("path", c.Request.URL.Path),
					zap.String("method", c.Request.Method),
					zap.String("stack", string(debug.Stack())),
				)

				// 返回 500 错误
				c.AbortWithStatusJSON(http.StatusInternalServerError, response.Response{
					Code:    response.ErrServerInternal.Code,
					Message: response.ErrServerInternal.Message,
				})
			}
		}()
		c.Next()
	}
}

// NotFoundHandler 404 处理器
func NotFoundHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		global.Logger.Warn("Route not found",
			zap.String("path", c.Request.URL.Path),
			zap.String("method", c.Request.Method),
			zap.String("client_ip", c.ClientIP()),
		)

		c.JSON(http.StatusNotFound, response.Response{
			Code:    response.ErrNotFound.Code,
			Message: "路由不存在: " + c.Request.URL.Path,
		})
	}
}

// MethodNotAllowedHandler 405 处理器
func MethodNotAllowedHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		global.Logger.Warn("Method not allowed",
			zap.String("path", c.Request.URL.Path),
			zap.String("method", c.Request.Method),
			zap.String("client_ip", c.ClientIP()),
		)

		c.JSON(http.StatusMethodNotAllowed, response.Response{
			Code:    response.ErrMethodNotAllow.Code,
			Message: "请求方法不允许: " + c.Request.Method,
		})
	}
}
