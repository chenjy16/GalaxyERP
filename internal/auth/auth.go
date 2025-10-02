package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Register 用户注册处理函数
func Register(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"message": "用户注册功能待实现，请使用 /api/v1/users/register",
	})
}

// Login 用户登录处理函数
func Login(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"message": "用户登录功能待实现，请使用 /api/v1/users/login",
	})
}
