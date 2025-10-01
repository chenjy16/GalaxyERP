package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/galaxyerp/galaxyErp/internal/container"
)

// RegisterAccountingRoutes 注册会计相关路由
func RegisterAccountingRoutes(router *gin.RouterGroup, container *container.Container) {
	// 会计科目管理
	accounts := router.Group("/accounts")
	{
		accounts.POST("/", func(c *gin.Context) {
			c.JSON(501, gin.H{"message": "会计功能待实现"})
		})
		accounts.GET("/", func(c *gin.Context) {
			c.JSON(501, gin.H{"message": "会计功能待实现"})
		})
	}

	// 交易记录管理
	transactions := router.Group("/transactions")
	{
		transactions.POST("/", func(c *gin.Context) {
			c.JSON(501, gin.H{"message": "会计功能待实现"})
		})
		transactions.GET("/", func(c *gin.Context) {
			c.JSON(501, gin.H{"message": "会计功能待实现"})
		})
	}

	// 财务报表
	reports := router.Group("/reports")
	{
		reports.GET("/balance-sheet", func(c *gin.Context) {
			c.JSON(501, gin.H{"message": "会计功能待实现"})
		})
		reports.GET("/income-statement", func(c *gin.Context) {
			c.JSON(501, gin.H{"message": "会计功能待实现"})
		})
	}
}