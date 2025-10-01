package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/galaxyerp/galaxyErp/internal/container"
)

// RegisterProductionRoutes 注册生产相关路由
func RegisterProductionRoutes(router *gin.RouterGroup, container *container.Container) {
	// 产品管理
	products := router.Group("/products")
	{
		products.POST("/", container.ProductionController.CreateProduct)
		products.GET("/:id", container.ProductionController.GetProduct)
		products.PUT("/:id", container.ProductionController.UpdateProduct)
		products.DELETE("/:id", container.ProductionController.DeleteProduct)
		products.GET("/", container.ProductionController.ListProducts)
		products.POST("/search", container.ProductionController.SearchProducts)
	}
}