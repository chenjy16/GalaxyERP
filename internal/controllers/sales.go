package controllers

import (
	"github.com/galaxyerp/galaxyErp/internal/dto"
	"github.com/galaxyerp/galaxyErp/internal/services"
	"github.com/gin-gonic/gin"
)

// SalesController 销售控制器
type SalesController struct {
	customerService     services.CustomerService
	salesOrderService   services.SalesOrderService
	quotationService    services.QuotationService
	salesInvoiceService services.SalesInvoiceService
	utils               *ControllerUtils
}

// NewSalesController 创建销售控制器实例
func NewSalesController(
	customerService services.CustomerService,
	salesOrderService services.SalesOrderService,
	quotationService services.QuotationService,
	salesInvoiceService services.SalesInvoiceService,
) *SalesController {
	return &SalesController{
		customerService:     customerService,
		salesOrderService:   salesOrderService,
		quotationService:    quotationService,
		salesInvoiceService: salesInvoiceService,
		utils:               NewControllerUtils(),
	}
}

// CreateCustomer 创建客户
// @Summary 创建客户
// @Description 创建新客户
// @Tags 客户管理
// @Accept json
// @Produce json
// @Param customer body dto.CustomerCreateRequest true "客户信息"
// @Success 201 {object} dto.CustomerResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/customers [post]
func (c *SalesController) CreateCustomer(ctx *gin.Context) {
	var req dto.CustomerCreateRequest
	if !c.utils.BindJSON(ctx, &req) {
		return
	}

	customer, err := c.customerService.CreateCustomer(ctx.Request.Context(), &req)
	if err != nil {
		c.utils.RespondInternalError(ctx, "创建客户失败")
		return
	}

	c.utils.RespondCreated(ctx, customer)
}

// GetCustomer 获取客户
// @Summary 获取客户
// @Description 根据ID获取客户信息
// @Tags 客户管理
// @Accept json
// @Produce json
// @Param id path int true "客户ID"
// @Success 200 {object} dto.CustomerResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/customers/{id} [get]
func (c *SalesController) GetCustomer(ctx *gin.Context) {
	id, ok := c.utils.ParseIDParam(ctx, "id")
	if !ok {
		return
	}

	customer, err := c.customerService.GetCustomer(ctx.Request.Context(), id)
	if err != nil {
		c.utils.RespondInternalError(ctx, "获取客户失败")
		return
	}

	c.utils.RespondOK(ctx, customer)
}

// UpdateCustomer 更新客户
// @Summary 更新客户
// @Description 更新客户信息
// @Tags 客户管理
// @Accept json
// @Produce json
// @Param id path int true "客户ID"
// @Param customer body dto.CustomerUpdateRequest true "客户信息"
// @Success 200 {object} dto.BaseResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/customers/{id} [put]
func (c *SalesController) UpdateCustomer(ctx *gin.Context) {
	id, ok := c.utils.ParseIDParam(ctx, "id")
	if !ok {
		return
	}

	var req dto.CustomerUpdateRequest
	if !c.utils.BindJSON(ctx, &req) {
		return
	}

	err := c.customerService.UpdateCustomer(ctx.Request.Context(), id, &req)
	if err != nil {
		c.utils.RespondInternalError(ctx, "更新客户失败")
		return
	}

	c.utils.RespondSuccess(ctx, "更新客户成功")
}

// DeleteCustomer 删除客户
// @Summary 删除客户
// @Description 删除客户
// @Tags 客户管理
// @Accept json
// @Produce json
// @Param id path int true "客户ID"
// @Success 200 {object} dto.BaseResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/customers/{id} [delete]
func (c *SalesController) DeleteCustomer(ctx *gin.Context) {
	id, ok := c.utils.ParseIDParam(ctx, "id")
	if !ok {
		return
	}

	err := c.customerService.DeleteCustomer(ctx.Request.Context(), id)
	if err != nil {
		c.utils.RespondInternalError(ctx, "删除客户失败")
		return
	}

	c.utils.RespondSuccess(ctx, "删除客户成功")
}

// ListCustomers 获取客户列表
// @Summary 获取客户列表
// @Description 获取客户列表
// @Tags 客户管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} dto.PaginatedResponse[dto.CustomerResponse]
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/customers [get]
func (c *SalesController) ListCustomers(ctx *gin.Context) {
	req := c.utils.ParsePaginationParams(ctx)

	response, err := c.customerService.ListCustomers(ctx.Request.Context(), req)
	if err != nil {
		c.utils.RespondInternalError(ctx, "获取客户列表失败")
		return
	}

	c.utils.RespondOK(ctx, response)
}

// SearchCustomers 搜索客户
// @Summary 搜索客户
// @Description 搜索客户
// @Tags 客户管理
// @Accept json
// @Produce json
// @Param keyword query string false "搜索关键词"
// @Param status query string false "客户状态"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} dto.PaginatedResponse[dto.CustomerResponse]
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/customers/search [get]
func (c *SalesController) SearchCustomers(ctx *gin.Context) {
	keyword := ctx.Query("keyword")
	status := ctx.Query("status")
	pagination := c.utils.ParsePaginationParams(ctx)

	req := &dto.SearchRequest{
		PaginationRequest: *pagination,
		Keyword:           keyword,
		Status:            status,
	}

	response, err := c.customerService.SearchCustomers(ctx.Request.Context(), req)
	if err != nil {
		c.utils.RespondInternalError(ctx, "搜索客户失败")
		return
	}

	c.utils.RespondOK(ctx, response)
}

// CreateSalesOrder 创建销售订单
// @Summary 创建销售订单
// @Description 创建新销售订单
// @Tags 销售订单
// @Accept json
// @Produce json
// @Param order body dto.SalesOrderCreateRequest true "订单信息"
// @Success 201 {object} dto.SalesOrderResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/sales/orders [post]
func (c *SalesController) CreateSalesOrder(ctx *gin.Context) {
	var req dto.SalesOrderCreateRequest
	if !c.utils.BindJSON(ctx, &req) {
		return
	}

	order, err := c.salesOrderService.CreateSalesOrder(ctx.Request.Context(), &req)
	if err != nil {
		c.utils.RespondInternalError(ctx, "创建销售订单失败")
		return
	}

	c.utils.RespondCreated(ctx, order)
}

// GetSalesOrder 获取销售订单
// @Summary 获取销售订单
// @Description 根据ID获取销售订单信息
// @Tags 销售订单
// @Accept json
// @Produce json
// @Param id path int true "订单ID"
// @Success 200 {object} dto.SalesOrderResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/sales/orders/{id} [get]
func (c *SalesController) GetSalesOrder(ctx *gin.Context) {
	id, ok := c.utils.ParseIDParam(ctx, "id")
	if !ok {
		return
	}

	order, err := c.salesOrderService.GetSalesOrder(ctx.Request.Context(), id)
	if err != nil {
		c.utils.RespondInternalError(ctx, "获取销售订单失败")
		return
	}

	c.utils.RespondOK(ctx, order)
}

// UpdateSalesOrder 更新销售订单
// @Summary 更新销售订单
// @Description 更新销售订单信息
// @Tags 销售订单
// @Accept json
// @Produce json
// @Param id path int true "订单ID"
// @Param order body dto.SalesOrderUpdateRequest true "订单信息"
// @Success 200 {object} dto.BaseResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/sales/orders/{id} [put]
func (c *SalesController) UpdateSalesOrder(ctx *gin.Context) {
	id, ok := c.utils.ParseIDParam(ctx, "id")
	if !ok {
		return
	}

	var req dto.SalesOrderUpdateRequest
	if !c.utils.BindJSON(ctx, &req) {
		return
	}

	err := c.salesOrderService.UpdateSalesOrder(ctx.Request.Context(), id, &req)
	if err != nil {
		c.utils.RespondInternalError(ctx, "更新销售订单失败")
		return
	}

	c.utils.RespondSuccess(ctx, "更新销售订单成功")
}

// DeleteSalesOrder 删除销售订单
// @Summary 删除销售订单
// @Description 删除销售订单
// @Tags 销售订单
// @Accept json
// @Produce json
// @Param id path int true "订单ID"
// @Success 200 {object} dto.BaseResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/sales/orders/{id} [delete]
func (c *SalesController) DeleteSalesOrder(ctx *gin.Context) {
	id, ok := c.utils.ParseIDParam(ctx, "id")
	if !ok {
		return
	}

	err := c.salesOrderService.DeleteSalesOrder(ctx.Request.Context(), id)
	if err != nil {
		c.utils.RespondInternalError(ctx, "删除销售订单失败")
		return
	}

	c.utils.RespondSuccess(ctx, "删除销售订单成功")
}

// ListSalesOrders 获取销售订单列表
// @Summary 获取销售订单列表
// @Description 获取销售订单列表
// @Tags 销售订单
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} dto.PaginatedResponse[dto.SalesOrderResponse]
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/sales/orders [get]
func (c *SalesController) ListSalesOrders(ctx *gin.Context) {
	req := c.utils.ParsePaginationParams(ctx)

	response, err := c.salesOrderService.ListSalesOrders(ctx.Request.Context(), req)
	if err != nil {
		c.utils.RespondInternalError(ctx, "获取销售订单列表失败")
		return
	}

	c.utils.RespondOK(ctx, response)
}

// UpdateOrderStatus 更新订单状态
// @Summary 更新订单状态
// @Description 更新销售订单状态
// @Tags 销售订单
// @Accept json
// @Produce json
// @Param id path int true "订单ID"
// @Param status body map[string]string true "状态信息"
// @Success 200 {object} dto.BaseResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/sales/orders/{id}/status [put]
func (c *SalesController) UpdateOrderStatus(ctx *gin.Context) {
	id, ok := c.utils.ParseIDParam(ctx, "id")
	if !ok {
		return
	}

	var req map[string]string
	if !c.utils.BindJSON(ctx, &req) {
		return
	}

	status, ok := req["status"]
	if !ok {
		c.utils.RespondBadRequest(ctx, "缺少状态参数")
		return
	}

	err := c.salesOrderService.UpdateOrderStatus(ctx.Request.Context(), id, status)
	if err != nil {
		c.utils.RespondInternalError(ctx, "更新订单状态失败")
		return
	}

	c.utils.RespondSuccess(ctx, "更新订单状态成功")
}

// CreateQuotation 创建报价单
// @Summary 创建报价单
// @Description 创建新的报价单
// @Tags 报价管理
// @Accept json
// @Produce json
// @Param quotation body dto.QuotationCreateRequest true "报价单信息"
// @Success 201 {object} dto.QuotationResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/quotations [post]
func (c *SalesController) CreateQuotation(ctx *gin.Context) {
	var req dto.QuotationCreateRequest
	if !c.utils.BindJSON(ctx, &req) {
		return
	}

	quotation, err := c.quotationService.CreateQuotation(ctx.Request.Context(), &req)
	if err != nil {
		c.utils.RespondInternalError(ctx, "创建报价单失败")
		return
	}

	c.utils.RespondCreated(ctx, quotation)
}

// GetQuotation 获取报价单
// @Summary 获取报价单
// @Description 根据ID获取报价单信息
// @Tags 报价管理
// @Accept json
// @Produce json
// @Param id path int true "报价单ID"
// @Success 200 {object} dto.QuotationResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/quotations/{id} [get]
func (c *SalesController) GetQuotation(ctx *gin.Context) {
	id, ok := c.utils.ParseIDParam(ctx, "id")
	if !ok {
		return
	}

	quotation, err := c.quotationService.GetQuotation(ctx.Request.Context(), id)
	if err != nil {
		c.utils.RespondInternalError(ctx, "获取报价单失败")
		return
	}

	c.utils.RespondOK(ctx, quotation)
}

// UpdateQuotation 更新报价单
// @Summary 更新报价单
// @Description 更新报价单信息
// @Tags 报价管理
// @Accept json
// @Produce json
// @Param id path int true "报价单ID"
// @Param quotation body dto.QuotationUpdateRequest true "报价单信息"
// @Success 200 {object} dto.BaseResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/quotations/{id} [put]
func (c *SalesController) UpdateQuotation(ctx *gin.Context) {
	id, ok := c.utils.ParseIDParam(ctx, "id")
	if !ok {
		return
	}

	var req dto.QuotationUpdateRequest
	if !c.utils.BindJSON(ctx, &req) {
		return
	}

	err := c.quotationService.UpdateQuotation(ctx.Request.Context(), id, &req)
	if err != nil {
		c.utils.RespondInternalError(ctx, "更新报价单失败")
		return
	}

	c.utils.RespondSuccess(ctx, "更新报价单成功")
}

// DeleteQuotation 删除报价单
// @Summary 删除报价单
// @Description 删除报价单
// @Tags 报价管理
// @Accept json
// @Produce json
// @Param id path int true "报价单ID"
// @Success 200 {object} dto.BaseResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/quotations/{id} [delete]
func (c *SalesController) DeleteQuotation(ctx *gin.Context) {
	id, ok := c.utils.ParseIDParam(ctx, "id")
	if !ok {
		return
	}

	err := c.quotationService.DeleteQuotation(ctx.Request.Context(), id)
	if err != nil {
		c.utils.RespondInternalError(ctx, "删除报价单失败")
		return
	}

	c.utils.RespondSuccess(ctx, "删除报价单成功")
}

// ListQuotations 获取报价单列表
// @Summary 获取报价单列表
// @Description 获取报价单列表
// @Tags 报价管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} dto.PaginatedResponse[dto.QuotationResponse]
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/quotations [get]
func (c *SalesController) ListQuotations(ctx *gin.Context) {
	req := c.utils.ParsePaginationParams(ctx)

	response, err := c.quotationService.ListQuotations(ctx.Request.Context(), req)
	if err != nil {
		c.utils.RespondInternalError(ctx, "获取报价单列表失败")
		return
	}

	c.utils.RespondOK(ctx, response)
}

// SearchQuotations 搜索报价单
// @Summary 搜索报价单
// @Description 根据关键词搜索报价单
// @Tags 报价管理
// @Accept json
// @Produce json
// @Param keyword query string false "搜索关键词"
// @Param status query string false "报价单状态"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} dto.PaginatedResponse[dto.QuotationResponse]
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/quotations/search [get]
func (c *SalesController) SearchQuotations(ctx *gin.Context) {
	keyword := ctx.Query("keyword")
	status := ctx.Query("status")
	pagination := c.utils.ParsePaginationParams(ctx)

	req := &dto.SearchRequest{
		PaginationRequest: *pagination,
		Keyword:           keyword,
		Status:            status,
	}

	response, err := c.quotationService.SearchQuotations(ctx.Request.Context(), req)
	if err != nil {
		c.utils.RespondInternalError(ctx, "搜索报价单失败")
		return
	}

	c.utils.RespondOK(ctx, response)
}

// ==================== 销售发票管理 ====================

// CreateSalesInvoice 创建销售发票
// @Summary 创建销售发票
// @Description 创建新的销售发票
// @Tags 销售发票管理
// @Accept json
// @Produce json
// @Param invoice body dto.SalesInvoiceCreateRequest true "销售发票信息"
// @Success 201 {object} dto.SalesInvoiceResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/sales-invoices [post]
func (c *SalesController) CreateSalesInvoice(ctx *gin.Context) {
	var req dto.SalesInvoiceCreateRequest
	if !c.utils.BindJSON(ctx, &req) {
		return
	}

	invoice, err := c.salesInvoiceService.CreateSalesInvoice(ctx, &req)
	if err != nil {
		c.utils.RespondInternalError(ctx, "创建销售发票失败")
		return
	}

	c.utils.RespondCreated(ctx, invoice)
}

// GetSalesInvoice 获取销售发票详情
// @Summary 获取销售发票详情
// @Description 根据ID获取销售发票详情
// @Tags 销售发票管理
// @Accept json
// @Produce json
// @Param id path int true "销售发票ID"
// @Success 200 {object} dto.SalesInvoiceResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/sales-invoices/{id} [get]
func (c *SalesController) GetSalesInvoice(ctx *gin.Context) {
	id, ok := c.utils.ParseIDParam(ctx, "id")
	if !ok {
		return
	}

	invoice, err := c.salesInvoiceService.GetSalesInvoice(id)
	if err != nil {
		c.utils.RespondInternalError(ctx, "获取销售发票失败")
		return
	}

	c.utils.RespondOK(ctx, invoice)
}

// UpdateSalesInvoice 更新销售发票
// @Summary 更新销售发票
// @Description 更新销售发票信息
// @Tags 销售发票管理
// @Accept json
// @Produce json
// @Param id path int true "销售发票ID"
// @Param invoice body dto.SalesInvoiceUpdateRequest true "销售发票信息"
// @Success 200 {object} dto.SalesInvoiceResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/sales-invoices/{id} [put]
func (c *SalesController) UpdateSalesInvoice(ctx *gin.Context) {
	id, ok := c.utils.ParseIDParam(ctx, "id")
	if !ok {
		return
	}

	var req dto.SalesInvoiceUpdateRequest
	if !c.utils.BindJSON(ctx, &req) {
		return
	}

	invoice, err := c.salesInvoiceService.UpdateSalesInvoice(ctx, id, &req)
	if err != nil {
		c.utils.RespondInternalError(ctx, "更新销售发票失败")
		return
	}

	c.utils.RespondOK(ctx, invoice)
}

// DeleteSalesInvoice 删除销售发票
// @Summary 删除销售发票
// @Description 删除销售发票
// @Tags 销售发票管理
// @Accept json
// @Produce json
// @Param id path int true "销售发票ID"
// @Success 200 {object} dto.BaseResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/sales-invoices/{id} [delete]
func (c *SalesController) DeleteSalesInvoice(ctx *gin.Context) {
	id, ok := c.utils.ParseIDParam(ctx, "id")
	if !ok {
		return
	}

	err := c.salesInvoiceService.DeleteSalesInvoice(ctx, id)
	if err != nil {
		c.utils.RespondInternalError(ctx, "删除销售发票失败")
		return
	}

	c.utils.RespondSuccess(ctx, "删除销售发票成功")
}

// ListSalesInvoices 获取销售发票列表
// @Summary 获取销售发票列表
// @Description 获取销售发票列表，支持分页
// @Tags 销售发票管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} dto.PaginatedResponse[dto.SalesInvoiceResponse]
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/sales-invoices [get]
func (c *SalesController) ListSalesInvoices(ctx *gin.Context) {
	pagination := c.utils.ParsePaginationParams(ctx)

	req := &dto.SalesInvoiceListRequest{
		PaginationRequest: *pagination,
	}

	response, err := c.salesInvoiceService.ListSalesInvoices(req)
	if err != nil {
		c.utils.RespondInternalError(ctx, "获取销售发票列表失败")
		return
	}

	c.utils.RespondOK(ctx, response)
}

// SubmitSalesInvoice 提交销售发票
// @Summary 提交销售发票
// @Description 提交销售发票进行审批
// @Tags 销售发票管理
// @Accept json
// @Produce json
// @Param id path int true "销售发票ID"
// @Success 200 {object} dto.SalesInvoiceResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/sales-invoices/{id}/submit [put]
func (c *SalesController) SubmitSalesInvoice(ctx *gin.Context) {
	id, ok := c.utils.ParseIDParam(ctx, "id")
	if !ok {
		return
	}

	invoice, err := c.salesInvoiceService.SubmitSalesInvoice(ctx, id)
	if err != nil {
		c.utils.RespondInternalError(ctx, "提交销售发票失败")
		return
	}

	c.utils.RespondOK(ctx, invoice)
}

// CancelSalesInvoice 取消销售发票
// @Summary 取消销售发票
// @Description 取消销售发票
// @Tags 销售发票管理
// @Accept json
// @Produce json
// @Param id path int true "销售发票ID"
// @Success 200 {object} dto.SalesInvoiceResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/sales-invoices/{id}/cancel [put]
func (c *SalesController) CancelSalesInvoice(ctx *gin.Context) {
	id, ok := c.utils.ParseIDParam(ctx, "id")
	if !ok {
		return
	}

	invoice, err := c.salesInvoiceService.CancelSalesInvoice(ctx, id)
	if err != nil {
		c.utils.RespondInternalError(ctx, "取消销售发票失败")
		return
	}

	c.utils.RespondOK(ctx, invoice)
}
