package routes

import (
	"github.com/galaxyerp/galaxyErp/internal/container"
	"github.com/gin-gonic/gin"
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

		// 报价单版本管理
		quotations.POST("/:id/versions", container.SalesController.CreateQuotationVersion)
		quotations.GET("/:id/versions/:versionNumber", container.SalesController.GetQuotationVersion)
		quotations.PUT("/:id/versions/:versionNumber/set-active", container.SalesController.SetActiveQuotationVersion)
		quotations.POST("/:id/versions/compare", container.SalesController.CompareQuotationVersions)
		quotations.GET("/:id/versions", container.SalesController.GetQuotationVersionHistory)
		quotations.POST("/:id/versions/:versionNumber/rollback", container.SalesController.RollbackQuotationVersion)
		quotations.DELETE("/:id/versions/:versionNumber", container.SalesController.DeleteQuotationVersion)
	}

	// 报价单模板管理
	quotationTemplates := router.Group("/quotation-templates")
	{
		quotationTemplates.POST("/", container.SalesController.CreateQuotationTemplate)
		quotationTemplates.GET("/:id", container.SalesController.GetQuotationTemplate)
		quotationTemplates.PUT("/:id", container.SalesController.UpdateQuotationTemplate)
		quotationTemplates.DELETE("/:id", container.SalesController.DeleteQuotationTemplate)
		quotationTemplates.GET("/", container.SalesController.ListQuotationTemplates)
		quotationTemplates.GET("/active", container.SalesController.GetActiveQuotationTemplates)
		quotationTemplates.GET("/default", container.SalesController.GetDefaultQuotationTemplate)
		quotationTemplates.PUT("/:id/set-default", container.SalesController.SetDefaultQuotationTemplate)
		quotationTemplates.POST("/:id/create-quotation", container.SalesController.CreateQuotationFromTemplate)
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
		// 发票付款管理
		invoices.POST("/:id/payments", container.SalesController.AddInvoicePayment)
		invoices.GET("/:id/payments", container.SalesController.GetInvoicePayments)
	}

	// 发货单管理
	deliveryNotes := router.Group("/delivery-notes")
	{
		deliveryNotes.POST("/", container.DeliveryNoteController.Create)
		deliveryNotes.GET("/:id", container.DeliveryNoteController.GetByID)
		deliveryNotes.PUT("/:id", container.DeliveryNoteController.Update)
		deliveryNotes.DELETE("/:id", container.DeliveryNoteController.Delete)
		deliveryNotes.GET("/", container.DeliveryNoteController.List)
		deliveryNotes.PATCH("/:id/status", container.DeliveryNoteController.UpdateStatus)
		deliveryNotes.POST("/from-sales-order", container.DeliveryNoteController.CreateFromSalesOrder)
		deliveryNotes.GET("/statistics", container.DeliveryNoteController.GetStatistics)
		deliveryNotes.GET("/trend", container.DeliveryNoteController.GetDeliveryTrend)
	}
}
