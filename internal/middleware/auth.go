package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware JWT认证中间件
func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取Authorization头
		authHeader := c.GetHeader("Authorization")

		// 检查是否为健康检查或公开接口
		path := c.Request.URL.Path
		if path == "/health" ||
			path == "/api/v1/users/login" ||
			path == "/api/v1/users/register" ||
			path == "/api/v1/auth/refresh" ||
			path == "/api/v1/auth/logout" {
			c.Next()
			return
		}

		// 检查Authorization头
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "未提供认证令牌",
			})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// 解析JWT token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// 验证签名方法
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "无效的认证令牌",
			})
			c.Abort()
			return
		}

		// 从token中提取用户信息
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if userIDFloat, exists := claims["user_id"]; exists {
				// JWT中的数字类型通常是float64
				if userIDFloat64, ok := userIDFloat.(float64); ok {
					c.Set("user_id", uint(userIDFloat64))
				}
			}
			if username, exists := claims["username"]; exists {
				c.Set("username", username.(string))
			}
		}

		c.Next()
	}
}
