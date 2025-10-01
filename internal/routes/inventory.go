package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/galaxyerp/galaxyErp/internal/container"
)

// RegisterInventoryRoutes 注册库存相关路由
func RegisterInventoryRoutes(router *gin.RouterGroup, container *container.Container) {
	// 物料管理
	items := router.Group("/items")
	{
		items.POST("/", container.InventoryController.CreateItem)
		items.GET("/:id", container.InventoryController.GetItem)
		items.PUT("/:id", container.InventoryController.UpdateItem)
		items.DELETE("/:id", container.InventoryController.DeleteItem)
		items.GET("/", container.InventoryController.ListItems)
		items.POST("/search", container.InventoryController.SearchItems)
	}

	// 库存管理
	stocks := router.Group("/stocks")
	{
		stocks.GET("/", container.InventoryController.ListStocks)
		stocks.POST("/", container.InventoryController.CreateStock)
		stocks.GET("/:id", container.InventoryController.GetStock)
		stocks.PUT("/:id", container.InventoryController.UpdateStock)
		stocks.DELETE("/:id", container.InventoryController.DeleteStock)
	}

	// 库存移动
	stockMovements := router.Group("/stock-movements")
	{
		stockMovements.GET("/", container.InventoryController.ListStockMovements)
		stockMovements.POST("/", container.InventoryController.CreateStockMovement)
		stockMovements.POST("/in", container.InventoryController.StockIn)
		stockMovements.POST("/out", container.InventoryController.StockOut)
		stockMovements.POST("/adjustment", container.InventoryController.StockAdjustment)
		stockMovements.POST("/transfer", container.InventoryController.StockTransfer)
	}

	// 仓库管理
	warehouses := router.Group("/warehouses")
	{
		warehouses.GET("/", container.InventoryController.ListWarehouses)
		warehouses.POST("/", container.InventoryController.CreateWarehouse)
		warehouses.GET("/:id", container.InventoryController.GetWarehouse)
		warehouses.PUT("/:id", container.InventoryController.UpdateWarehouse)
		warehouses.DELETE("/:id", container.InventoryController.DeleteWarehouse)
	}

	// 库存查询
	stock := router.Group("/stock")
	{
		stock.GET("/item/:item_id", container.InventoryController.GetStockByItemID)
	}

	// 库存报告和统计
	reports := router.Group("/inventory-reports")
	{
		reports.GET("/stats", container.InventoryController.GetInventoryStats)
		reports.GET("/report", container.InventoryController.GetInventoryReport)
		reports.GET("/abc-analysis", container.InventoryController.GetABCAnalysis)
		reports.GET("/export", container.InventoryController.ExportInventoryReport)
	}
}