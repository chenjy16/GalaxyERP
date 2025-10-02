package routes

import (
	"github.com/galaxyerp/galaxyErp/internal/container"
	"github.com/gin-gonic/gin"
)

// RegisterAccountingRoutes 注册会计相关路由
func RegisterAccountingRoutes(router *gin.RouterGroup, container *container.Container) {
	accountingController := container.AccountingController

	// 会计科目管理
	accounts := router.Group("/accounts")
	{
		accounts.POST("/", accountingController.CreateAccount)
		accounts.GET("/", accountingController.GetAccountList)
		accounts.GET("/:id", accountingController.GetAccount)
		accounts.PUT("/:id", accountingController.UpdateAccount)
		accounts.DELETE("/:id", accountingController.DeleteAccount)
		accounts.GET("/code/:code", accountingController.GetAccountByCode)
		accounts.GET("/:id/children", accountingController.GetAccountChildren)
	}

	// 会计分录管理
	journalEntries := router.Group("/journal-entries")
	{
		journalEntries.POST("/", accountingController.CreateJournalEntry)
		journalEntries.GET("/", accountingController.GetJournalEntryList)
		journalEntries.GET("/:id", accountingController.GetJournalEntry)
		journalEntries.PUT("/:id", accountingController.UpdateJournalEntry)
		journalEntries.DELETE("/:id", accountingController.DeleteJournalEntry)
	}

	// 科目类型
	router.GET("/account-types", accountingController.GetAccountTypes)

	// 财务报表 (暂时保留占位符)
	reports := router.Group("/reports")
	{
		reports.GET("/balance-sheet", func(c *gin.Context) {
			c.JSON(501, gin.H{"message": "财务报表功能待实现"})
		})
		reports.GET("/income-statement", func(c *gin.Context) {
			c.JSON(501, gin.H{"message": "财务报表功能待实现"})
		})
	}
}
