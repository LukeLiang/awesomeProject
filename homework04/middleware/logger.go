package middleware

import (
	"time"

	"awesomeProject/homework04/global"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RequestLogger 请求日志中间件
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		method := c.Request.Method
		clientIP := c.ClientIP()

		// 处理请求
		c.Next()

		// 计算耗时
		latency := time.Since(start)
		statusCode := c.Writer.Status()

		// 构建日志字段
		fields := []zap.Field{
			zap.String("method", method),
			zap.String("path", path),
			zap.String("query", query),
			zap.Int("status", statusCode),
			zap.Duration("latency", latency),
			zap.String("client_ip", clientIP),
		}

		// 根据状态码选择日志级别
		if statusCode >= 500 {
			global.Logger.Error("Server error", fields...)
		} else if statusCode >= 400 {
			global.Logger.Warn("Client error", fields...)
		} else {
			global.Logger.Info("Request completed", fields...)
		}

		// 慢请求警告（超过 3 秒）
		if latency > 3*time.Second {
			global.Logger.Warn("Slow request detected",
				zap.String("path", path),
				zap.Duration("latency", latency),
			)
		}
	}
}
