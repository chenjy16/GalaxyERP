package controllers

import (
	"github.com/galaxyerp/galaxyErp/internal/dto"
	"github.com/galaxyerp/galaxyErp/internal/services"
	"github.com/galaxyerp/galaxyErp/internal/utils"
	"github.com/gin-gonic/gin"
)

// SalesController 销售控制器
type SalesController struct {
	customerService     services.CustomerService
	salesOrderService   services.SalesOrderService
	quotationService    services.QuotationService
	templateService     services.QuotationTemplateService
	salesInvoiceService services.SalesInvoiceService
	versionService      services.QuotationVersionService
	utils               *ControllerUtils
}

// NewSalesController 创建销售控制器实例
func NewSalesController(
	customerService services.CustomerService,
	salesOrderService services.SalesOrderService,
	quotationService services.QuotationService,
	templateService services.QuotationTemplateService,
	salesInvoiceService services.SalesInvoiceService,
	versionService services.QuotationVersionService,
) *SalesController {
	return &SalesController{
		customerService:     customerService,
		salesOrderService:   salesOrderService,
		quotationService:    quotationService,
		templateService:     templateService,
		salesInvoiceService: salesInvoiceService,
		versionService:      versionService,
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
	if !c.utils.BindAndValidateJSON(ctx, &req) {
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
	if !c.utils.BindAndValidateJSON(ctx, &req) {
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

	// 转换为统一的分页响应格式
	pagination := c.utils.CreatePagination(response.Page, response.Limit, response.Total)
	c.utils.RespondPaginated(ctx, response.Data, pagination, "获取客户列表成功")
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

	// 转换为统一的分页响应格式
	pagination2 := c.utils.CreatePagination(response.Page, response.Limit, response.Total)
	c.utils.RespondPaginated(ctx, response.Data, pagination2, "搜索客户成功")
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
	var req dto.CreateSalesOrderRequest

	if !c.utils.BindAndValidateJSON(ctx, &req) {
		return
	}

	// 从上下文获取用户ID
	userID := utils.GetUserIDFromContext(ctx)
	// 添加调试日志
	utils.Debug("userID from context", utils.Uint("user_id", userID))
	if userID == 0 {
		c.utils.RespondUnauthorized(ctx, "用户未认证")
		return
	}

	order, err := c.salesOrderService.CreateSalesOrder(ctx.Request.Context(), &req, userID)
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
	if !c.utils.BindAndValidateJSON(ctx, &req) {
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

	// 转换为统一的分页响应格式
	pagination := c.utils.CreatePagination(response.Page, response.Limit, response.Total)
	c.utils.RespondPaginated(ctx, response.Data, pagination, "获取销售订单列表成功")
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
	if !c.utils.BindAndValidateJSON(ctx, &req) {
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
	if !c.utils.BindAndValidateJSON(ctx, &req) {
		return
	}

	quotation, err := c.quotationService.CreateQuotation(ctx.Request.Context(), &req)
	if err != nil {
		c.utils.RespondInternalError(ctx, "创建报价单失败")
		return
	}

	c.utils.RespondCreated(ctx, quotation)
}

// ==================== 报价单版本管理 ====================

// CreateQuotationVersion 创建报价单版本
// @Summary 创建报价单版本
// @Description 为报价单创建新版本
// @Tags 报价单版本
// @Accept json
// @Produce json
// @Param request body dto.QuotationVersionCreateRequest true "版本创建请求"
// @Success 201 {object} dto.QuotationVersionResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/quotation-versions [post]
func (c *SalesController) CreateQuotationVersion(ctx *gin.Context) {
	var req dto.QuotationVersionCreateRequest
	if !c.utils.BindAndValidateJSON(ctx, &req) {
		return
	}

	version, err := c.versionService.CreateVersion(ctx.Request.Context(), &req)
	if err != nil {
		c.utils.RespondInternalError(ctx, "创建版本失败")
		return
	}

	c.utils.RespondCreated(ctx, version)
}

// GetQuotationVersion 获取报价单版本
// @Summary 获取报价单版本
// @Description 根据ID获取报价单版本信息
// @Tags 报价单版本
// @Accept json
// @Produce json
// @Param id path int true "版本ID"
// @Success 200 {object} dto.QuotationVersionResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/quotation-versions/{id} [get]
func (c *SalesController) GetQuotationVersion(ctx *gin.Context) {
	id, ok := c.utils.ParseIDParam(ctx, "id")
	if !ok {
		return
	}

	version, err := c.versionService.GetVersion(ctx.Request.Context(), id)
	if err != nil {
		c.utils.RespondInternalError(ctx, "获取版本失败")
		return
	}

	c.utils.RespondOK(ctx, version)
}

// SetActiveQuotationVersion 设置活跃版本
// @Summary 设置活跃版本
// @Description 将指定版本设置为活跃版本
// @Tags 报价单版本
// @Accept json
// @Produce json
// @Param quotation_id path int true "报价单ID"
// @Param version_number path int true "版本号"
// @Success 200 {object} dto.BaseResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/quotations/{quotation_id}/versions/{version_number}/set-active [put]
func (c *SalesController) SetActiveQuotationVersion(ctx *gin.Context) {
	quotationID, ok := c.utils.ParseIDParam(ctx, "quotation_id")
	if !ok {
		return
	}

	versionNumber, ok := c.utils.ParseIDParam(ctx, "version_number")
	if !ok {
		return
	}

	err := c.versionService.SetActiveVersion(ctx.Request.Context(), quotationID, int(versionNumber))
	if err != nil {
		c.utils.RespondInternalError(ctx, "设置活跃版本失败")
		return
	}

	c.utils.RespondSuccess(ctx, "设置活跃版本成功")
}

// CompareQuotationVersions 比较报价单版本
// @Summary 比较报价单版本
// @Description 比较两个报价单版本的差异
// @Tags 报价单版本
// @Accept json
// @Produce json
// @Param request body dto.QuotationVersionCompareRequest true "版本比较请求"
// @Success 200 {object} dto.QuotationVersionComparisonResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/quotation-versions/compare [post]
func (c *SalesController) CompareQuotationVersions(ctx *gin.Context) {
	var req dto.QuotationVersionCompareRequest
	if !c.utils.BindAndValidateJSON(ctx, &req) {
		return
	}

	comparison, err := c.versionService.CompareVersions(ctx.Request.Context(), &req)
	if err != nil {
		c.utils.RespondInternalError(ctx, "比较版本失败")
		return
	}

	c.utils.RespondOK(ctx, comparison)
}

// GetQuotationVersionHistory 获取报价单版本历史
// @Summary 获取报价单版本历史
// @Description 获取指定报价单的版本历史记录
// @Tags 报价单版本
// @Accept json
// @Produce json
// @Param quotation_id path int true "报价单ID"
// @Success 200 {object} dto.QuotationVersionHistoryResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/quotations/{quotation_id}/version-history [get]
func (c *SalesController) GetQuotationVersionHistory(ctx *gin.Context) {
	quotationID, ok := c.utils.ParseIDParam(ctx, "quotation_id")
	if !ok {
		return
	}

	history, err := c.versionService.GetVersionHistory(ctx.Request.Context(), quotationID)
	if err != nil {
		c.utils.RespondInternalError(ctx, "获取版本历史失败")
		return
	}

	c.utils.RespondOK(ctx, history)
}

// RollbackQuotationVersion 回滚报价单版本
// @Summary 回滚报价单版本
// @Description 将报价单回滚到指定版本
// @Tags 报价单版本
// @Accept json
// @Produce json
// @Param request body dto.QuotationVersionRollbackRequest true "版本回滚请求"
// @Success 200 {object} dto.QuotationVersionResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/quotation-versions/rollback [post]
func (c *SalesController) RollbackQuotationVersion(ctx *gin.Context) {
	var req dto.QuotationVersionRollbackRequest
	if !c.utils.BindAndValidateJSON(ctx, &req) {
		return
	}

	version, err := c.versionService.RollbackToVersion(ctx.Request.Context(), &req)
	if err != nil {
		c.utils.RespondInternalError(ctx, "回滚版本失败")
		return
	}

	c.utils.RespondOK(ctx, version)
}

// DeleteQuotationVersion 删除报价单版本
// @Summary 删除报价单版本
// @Description 删除指定的报价单版本
// @Tags 报价单版本
// @Accept json
// @Produce json
// @Param id path int true "版本ID"
// @Success 200 {object} dto.BaseResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/quotation-versions/{id} [delete]
func (c *SalesController) DeleteQuotationVersion(ctx *gin.Context) {
	id, ok := c.utils.ParseIDParam(ctx, "id")
	if !ok {
		return
	}

	err := c.versionService.DeleteVersion(ctx.Request.Context(), id)
	if err != nil {
		c.utils.RespondInternalError(ctx, "删除版本失败")
		return
	}

	c.utils.RespondSuccess(ctx, "删除版本成功")
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
	if !c.utils.BindAndValidateJSON(ctx, &req) {
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

	// 转换为统一的分页响应格式
	pagination := c.utils.CreatePagination(response.Page, response.Limit, response.Total)
	c.utils.RespondPaginated(ctx, response.Data, pagination, "获取报价单列表成功")
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

	// 转换为统一的分页响应格式
	pagination2 := c.utils.CreatePagination(response.Page, response.Limit, response.Total)
	c.utils.RespondPaginated(ctx, response.Data, pagination2, "搜索报价单成功")
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
	if !c.utils.BindAndValidateJSON(ctx, &req) {
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

// AddInvoicePayment 为销售发票添加付款记录
// @Summary 为销售发票添加付款记录
// @Description 为销售发票添加付款记录
// @Tags 销售发票管理
// @Accept json
// @Produce json
// @Param id path int true "销售发票ID"
// @Param payment body dto.InvoicePaymentCreateRequest true "付款信息"
// @Success 200 {object} dto.SalesInvoiceResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/sales-invoices/{id}/payments [post]
func (c *SalesController) AddInvoicePayment(ctx *gin.Context) {
	id, ok := c.utils.ParseIDParam(ctx, "id")
	if !ok {
		return
	}

	var req dto.InvoicePaymentCreateRequest
	if !c.utils.BindAndValidateJSON(ctx, &req) {
		return
	}

	invoice, err := c.salesInvoiceService.AddPayment(ctx, id, &req)
	if err != nil {
		c.utils.RespondInternalError(ctx, "添加付款记录失败")
		return
	}

	c.utils.RespondOK(ctx, invoice)
}

// GetInvoicePayments 获取销售发票的付款记录
// @Summary 获取销售发票的付款记录
// @Description 获取销售发票的所有付款记录
// @Tags 销售发票管理
// @Accept json
// @Produce json
// @Param id path int true "销售发票ID"
// @Success 200 {array} dto.InvoicePaymentResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/sales-invoices/{id}/payments [get]
func (c *SalesController) GetInvoicePayments(ctx *gin.Context) {
	id, ok := c.utils.ParseIDParam(ctx, "id")
	if !ok {
		return
	}

	payments, err := c.salesInvoiceService.GetPayments(id)
	if err != nil {
		c.utils.RespondInternalError(ctx, "获取付款记录失败")
		return
	}

	c.utils.RespondOK(ctx, payments)
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
	if !c.utils.BindAndValidateJSON(ctx, &req) {
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

	// 转换为统一的分页响应格式
	pagination2 := c.utils.CreatePagination(response.Page, response.Limit, response.Total)
	c.utils.RespondPaginated(ctx, response.Data, pagination2, "获取销售发票列表成功")
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

// ==================== 报价单模板管理 ====================

// CreateQuotationTemplate 创建报价单模板
// @Summary 创建报价单模板
// @Description 创建新的报价单模板
// @Tags 报价单模板
// @Accept json
// @Produce json
// @Param template body dto.QuotationTemplateCreateRequest true "模板信息"
// @Success 201 {object} dto.QuotationTemplateResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/quotation-templates [post]
func (c *SalesController) CreateQuotationTemplate(ctx *gin.Context) {
	var req dto.QuotationTemplateCreateRequest
	if !c.utils.BindAndValidateJSON(ctx, &req) {
		return
	}

	template, err := c.templateService.CreateTemplate(ctx.Request.Context(), &req)
	if err != nil {
		c.utils.RespondInternalError(ctx, "创建模板失败")
		return
	}

	c.utils.RespondCreated(ctx, template)
}

// GetQuotationTemplate 获取报价单模板
// @Summary 获取报价单模板
// @Description 根据ID获取报价单模板信息
// @Tags 报价单模板
// @Accept json
// @Produce json
// @Param id path int true "模板ID"
// @Success 200 {object} dto.QuotationTemplateResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/quotation-templates/{id} [get]
func (c *SalesController) GetQuotationTemplate(ctx *gin.Context) {
	id, ok := c.utils.ParseIDParam(ctx, "id")
	if !ok {
		return
	}

	template, err := c.templateService.GetTemplate(ctx.Request.Context(), id)
	if err != nil {
		c.utils.RespondInternalError(ctx, "获取模板失败")
		return
	}

	c.utils.RespondOK(ctx, template)
}

// UpdateQuotationTemplate 更新报价单模板
// @Summary 更新报价单模板
// @Description 更新报价单模板信息
// @Tags 报价单模板
// @Accept json
// @Produce json
// @Param id path int true "模板ID"
// @Param template body dto.QuotationTemplateUpdateRequest true "模板信息"
// @Success 200 {object} dto.BaseResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/quotation-templates/{id} [put]
func (c *SalesController) UpdateQuotationTemplate(ctx *gin.Context) {
	id, ok := c.utils.ParseIDParam(ctx, "id")
	if !ok {
		return
	}

	var req dto.QuotationTemplateUpdateRequest
	if !c.utils.BindAndValidateJSON(ctx, &req) {
		return
	}

	if err := c.templateService.UpdateTemplate(ctx.Request.Context(), id, &req); err != nil {
		c.utils.RespondInternalError(ctx, "更新模板失败")
		return
	}

	c.utils.RespondSuccess(ctx, "更新模板成功")
}

// DeleteQuotationTemplate 删除报价单模板
// @Summary 删除报价单模板
// @Description 删除报价单模板
// @Tags 报价单模板
// @Accept json
// @Produce json
// @Param id path int true "模板ID"
// @Success 200 {object} dto.BaseResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/quotation-templates/{id} [delete]
func (c *SalesController) DeleteQuotationTemplate(ctx *gin.Context) {
	id, ok := c.utils.ParseIDParam(ctx, "id")
	if !ok {
		return
	}

	if err := c.templateService.DeleteTemplate(ctx.Request.Context(), id); err != nil {
		c.utils.RespondInternalError(ctx, "删除模板失败")
		return
	}

	c.utils.RespondSuccess(ctx, "删除模板成功")
}

// ListQuotationTemplates 获取报价单模板列表
// @Summary 获取报价单模板列表
// @Description 获取报价单模板列表
// @Tags 报价单模板
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} dto.PaginatedResponse[dto.QuotationTemplateResponse]
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/quotation-templates [get]
func (c *SalesController) ListQuotationTemplates(ctx *gin.Context) {
	req := c.utils.ParsePaginationParams(ctx)

	response, err := c.templateService.ListTemplates(ctx.Request.Context(), req)
	if err != nil {
		c.utils.RespondInternalError(ctx, "获取模板列表失败")
		return
	}

	// 转换为统一的分页响应格式
	pagination := c.utils.CreatePagination(response.Page, response.Limit, response.Total)
	c.utils.RespondPaginated(ctx, response.Data, pagination, "获取模板列表成功")
}

// GetActiveQuotationTemplates 获取活跃的报价单模板
// @Summary 获取活跃的报价单模板
// @Description 获取所有活跃的报价单模板
// @Tags 报价单模板
// @Accept json
// @Produce json
// @Success 200 {array} dto.QuotationTemplateResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/quotation-templates/active [get]
func (c *SalesController) GetActiveQuotationTemplates(ctx *gin.Context) {
	templates, err := c.templateService.GetActiveTemplates(ctx.Request.Context())
	if err != nil {
		c.utils.RespondInternalError(ctx, "获取活跃模板失败")
		return
	}

	c.utils.RespondOK(ctx, templates)
}

// GetDefaultQuotationTemplate 获取默认报价单模板
// @Summary 获取默认报价单模板
// @Description 获取默认的报价单模板
// @Tags 报价单模板
// @Accept json
// @Produce json
// @Success 200 {object} dto.QuotationTemplateResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/quotation-templates/default [get]
func (c *SalesController) GetDefaultQuotationTemplate(ctx *gin.Context) {
	template, err := c.templateService.GetDefaultTemplate(ctx.Request.Context())
	if err != nil {
		c.utils.RespondInternalError(ctx, "获取默认模板失败")
		return
	}

	if template == nil {
		c.utils.RespondNotFound(ctx, "未找到默认模板")
		return
	}

	c.utils.RespondOK(ctx, template)
}

// SetDefaultQuotationTemplate 设置默认报价单模板
// @Summary 设置默认报价单模板
// @Description 设置指定模板为默认模板
// @Tags 报价单模板
// @Accept json
// @Produce json
// @Param id path int true "模板ID"
// @Success 200 {object} dto.BaseResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/quotation-templates/{id}/set-default [post]
func (c *SalesController) SetDefaultQuotationTemplate(ctx *gin.Context) {
	id, ok := c.utils.ParseIDParam(ctx, "id")
	if !ok {
		return
	}

	if err := c.templateService.SetAsDefault(ctx.Request.Context(), id); err != nil {
		c.utils.RespondInternalError(ctx, "设置默认模板失败")
		return
	}

	c.utils.RespondSuccess(ctx, "设置默认模板成功")
}

// CreateQuotationFromTemplate 从模板创建报价单
// @Summary 从模板创建报价单
// @Description 使用指定模板创建新的报价单
// @Tags 报价单模板
// @Accept json
// @Produce json
// @Param request body dto.CreateQuotationFromTemplateRequest true "创建请求"
// @Success 201 {object} dto.QuotationResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/quotation-templates/create-quotation [post]
func (c *SalesController) CreateQuotationFromTemplate(ctx *gin.Context) {
	var req dto.CreateQuotationFromTemplateRequest
	if !c.utils.BindAndValidateJSON(ctx, &req) {
		return
	}

	quotation, err := c.templateService.CreateQuotationFromTemplate(ctx.Request.Context(), req.TemplateID, req.CustomerID)
	if err != nil {
		c.utils.RespondInternalError(ctx, "从模板创建报价单失败")
		return
	}

	c.utils.RespondCreated(ctx, quotation)
}
