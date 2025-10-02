package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/galaxyerp/galaxyErp/internal/container"
)

// RegisterSalesRoutes 注册销售相关路由
func RegisterSalesRoutes(router *gin.RouterGroup, container *container.Container) {
	// 客户管理
	customers := router.Group("/customers")
	{
		customers.POST("/", container.SalesController.CreateCustomer)
		customers.GET("/:id", container.SalesController.GetCustomer)
		customers.PUT("/:id", container.SalesController.UpdateCustomer)
		customers.DELETE("/:id", container.SalesController.DeleteCustomer)
		customers.GET("/", container.SalesController.ListCustomers)
		customers.POST("/search", container.SalesController.SearchCustomers)
	}

	// 销售订单管理
	orders := router.Group("/sales-orders")
	{
		orders.POST("/", container.SalesController.CreateSalesOrder)
		orders.GET("/:id", container.SalesController.GetSalesOrder)
		orders.PUT("/:id", container.SalesController.UpdateSalesOrder)
		orders.DELETE("/:id", container.SalesController.DeleteSalesOrder)
		orders.GET("/", container.SalesController.ListSalesOrders)
		orders.PUT("/:id/status", container.SalesController.UpdateOrderStatus)
	}

	// 报价单管理
	quotations := router.Group("/quotations")
	{
		quotations.POST("/", container.SalesController.CreateQuotation)
		quotations.GET("/:id", container.SalesController.GetQuotation)
		quotations.PUT("/:id", container.SalesController.UpdateQuotation)
		quotations.DELETE("/:id", container.SalesController.DeleteQuotation)
		quotations.GET("/", container.SalesController.ListQuotations)
		quotations.GET("/search", container.SalesController.SearchQuotations)
	}

	// 销售发票管理
	invoices := router.Group("/sales-invoices")
	{
		invoices.POST("/", container.SalesController.CreateSalesInvoice)
		invoices.GET("/:id", container.SalesController.GetSalesInvoice)
		invoices.PUT("/:id", container.SalesController.UpdateSalesInvoice)
		invoices.DELETE("/:id", container.SalesController.DeleteSalesInvoice)
		invoices.GET("/", container.SalesController.ListSalesInvoices)
		invoices.PUT("/:id/submit", container.SalesController.SubmitSalesInvoice)
		invoices.PUT("/:id/cancel", container.SalesController.CancelSalesInvoice)
	}
}