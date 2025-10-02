package routes

import (
	"github.com/galaxyerp/galaxyErp/internal/container"
	"github.com/gin-gonic/gin"
)

// RegisterPurchaseRoutes 注册采购相关路由
func RegisterPurchaseRoutes(router *gin.RouterGroup, container *container.Container) {
	purchaseController := container.PurchaseController

	// 供应商管理
	suppliers := router.Group("/suppliers")
	{
		suppliers.POST("/", purchaseController.CreateSupplier)
		suppliers.GET("/", purchaseController.ListSuppliers)
		suppliers.GET("/:id", purchaseController.GetSupplier)
		suppliers.PUT("/:id", purchaseController.UpdateSupplier)
		suppliers.DELETE("/:id", purchaseController.DeleteSupplier)
	}

	// 采购订单管理
	orders := router.Group("/purchase-orders")
	{
		orders.POST("/", purchaseController.CreatePurchaseOrder)
		orders.GET("/", purchaseController.ListPurchaseOrders)
		orders.GET("/:id", purchaseController.GetPurchaseOrder)
		orders.PUT("/:id", purchaseController.UpdatePurchaseOrder)
		orders.DELETE("/:id", purchaseController.DeletePurchaseOrder)
		orders.POST("/:id/confirm", purchaseController.ConfirmPurchaseOrder)
		orders.POST("/:id/cancel", purchaseController.CancelPurchaseOrder)
	}

	// 采购申请管理
	requests := router.Group("/purchase-requests")
	{
		requests.POST("/", purchaseController.CreatePurchaseRequest)
		requests.GET("/", purchaseController.ListPurchaseRequests)
		requests.GET("/:id", purchaseController.GetPurchaseRequest)
		requests.PUT("/:id", purchaseController.UpdatePurchaseRequest)
		requests.DELETE("/:id", purchaseController.DeletePurchaseRequest)
		requests.POST("/:id/submit", purchaseController.SubmitPurchaseRequest)
		requests.POST("/:id/approve", purchaseController.ApprovePurchaseRequest)
		requests.POST("/:id/reject", purchaseController.RejectPurchaseRequest)
	}

	// 采购统计
	router.GET("/purchase/stats", purchaseController.GetPurchaseStats)
}
