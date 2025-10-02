package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Register 用户注册处理函数
func Register(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"error": "功能暂未实现",
	})
}

// Login 用户登录处理函数
func Login(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"error": "功能暂未实现",
	})
}
