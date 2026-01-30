package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CorsHandler 跨域中间件
func CorsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		// 设置允许的源（生产环境建议配置具体域名）
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
		}

		// 允许的请求方法
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")

		// 允许的请求头
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Requested-With")

		// 允许携带凭证（cookies）
		c.Header("Access-Control-Allow-Credentials", "true")

		// 预检请求缓存时间（秒）
		c.Header("Access-Control-Max-Age", "86400")

		// 处理 OPTIONS 预检请求
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
