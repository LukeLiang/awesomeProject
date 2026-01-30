package middleware

import (
	"fmt"
	"strings"

	"awesomeProject/homework04/common"
	"awesomeProject/homework04/common/response"
	"awesomeProject/homework04/global"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func VerifyToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		if path == "/login" || path == "/register" {
			c.Next()
			return
		}

		authHeader := c.GetHeader("Authorization")
		// 按空格分割，取出真正的 token 部分 (如果是 Bearer 格式)
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			global.Logger.Warn("Token format error",
				zap.String("path", path),
				zap.String("client_ip", c.ClientIP()),
				zap.String("auth_header", authHeader),
			)
			response.Fail(c, response.ErrTokenFormat)
			c.Abort()
			return
		}
		tokenString := parts[1]

		// 解析并校验 Token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// 必须验证签名算法是否是你预期的 (防止 "none" 算法攻击)
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return global.Config.Secret, nil
		})

		// 判断 Token 是否有效
		if err != nil || !token.Valid {
			global.Logger.Warn("Token invalid or expired",
				zap.String("path", path),
				zap.String("client_ip", c.ClientIP()),
				zap.Error(err),
			)
			response.Fail(c, response.ErrTokenInvalid)
			c.Abort()
			return
		}

		// 将用户信息存储到 gin.Context 中
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			currentUser := common.GlobalUser{
				ID:       uint(claims["id"].(float64)),
				Username: claims["username"].(string),
			}
			c.Set("currentUser", currentUser)

			global.Logger.Debug("Token verified",
				zap.Uint("user_id", currentUser.ID),
				zap.String("username", currentUser.Username),
			)
		}

		c.Next()
	}
}
