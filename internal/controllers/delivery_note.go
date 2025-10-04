package controllers

import (
	"strconv"

	"github.com/galaxyerp/galaxyErp/internal/dto"
	"github.com/galaxyerp/galaxyErp/internal/services"
	"github.com/gin-gonic/gin"
)

// DeliveryNoteController 发货单控制器
type DeliveryNoteController struct {
	deliveryNoteService services.DeliveryNoteServiceInterface
	utils               *ControllerUtils
}

// NewDeliveryNoteController 创建发货单控制器
func NewDeliveryNoteController(deliveryNoteService services.DeliveryNoteServiceInterface) *DeliveryNoteController {
	return &DeliveryNoteController{
		deliveryNoteService: deliveryNoteService,
		utils:               NewControllerUtils(),
	}
}

// Create 创建发货单
// @Summary 创建发货单
// @Description 创建新的发货单
// @Tags 发货单
// @Accept json
// @Produce json
// @Param request body dto.DeliveryNoteCreateRequest true "发货单创建请求"
// @Success 201 {object} utils.Response{data=models.DeliveryNote}
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/delivery-notes [post]
func (c *DeliveryNoteController) Create(ctx *gin.Context) {
	var req dto.DeliveryNoteCreateRequest
	if !c.utils.BindAndValidateJSON(ctx, &req) {
		return
	}

	// 获取当前用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		c.utils.RespondUnauthorized(ctx, "用户未登录")
		return
	}

	deliveryNote, err := c.deliveryNoteService.Create(ctx, &req, userID.(uint))
	if err != nil {
		c.utils.RespondInternalError(ctx, "创建发货单失败")
		return
	}

	c.utils.RespondCreated(ctx, deliveryNote)
}

// GetByID 根据ID获取发货单
// @Summary 根据ID获取发货单
// @Description 根据ID获取发货单详情
// @Tags 发货单
// @Accept json
// @Produce json
// @Param id path int true "发货单ID"
// @Success 200 {object} utils.Response{data=models.DeliveryNote}
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/delivery-notes/{id} [get]
func (c *DeliveryNoteController) GetByID(ctx *gin.Context) {
	id, ok := c.utils.ParseIDParam(ctx, "id")
	if !ok {
		return
	}

	deliveryNote, err := c.deliveryNoteService.GetByID(id)
	if err != nil {
		c.utils.RespondNotFound(ctx, "发货单不存在")
		return
	}

	c.utils.RespondOK(ctx, deliveryNote)
}

// Update 更新发货单
// @Summary 更新发货单
// @Description 更新发货单信息
// @Tags 发货单
// @Accept json
// @Produce json
// @Param id path int true "发货单ID"
// @Param request body dto.DeliveryNoteUpdateRequest true "发货单更新请求"
// @Success 200 {object} utils.Response{data=models.DeliveryNote}
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/delivery-notes/{id} [put]
func (c *DeliveryNoteController) Update(ctx *gin.Context) {
	id, ok := c.utils.ParseIDParam(ctx, "id")
	if !ok {
		return
	}

	var req dto.DeliveryNoteUpdateRequest
	if !c.utils.BindAndValidateJSON(ctx, &req) {
		return
	}

	deliveryNote, err := c.deliveryNoteService.Update(id, &req)
	if err != nil {
		c.utils.RespondInternalError(ctx, "更新发货单失败")
		return
	}

	c.utils.RespondOK(ctx, deliveryNote)
}

// Delete 删除发货单
// @Summary 删除发货单
// @Description 删除发货单
// @Tags 发货单
// @Accept json
// @Produce json
// @Param id path int true "发货单ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/delivery-notes/{id} [delete]
func (c *DeliveryNoteController) Delete(ctx *gin.Context) {
	id, ok := c.utils.ParseIDParam(ctx, "id")
	if !ok {
		return
	}

	err := c.deliveryNoteService.Delete(id)
	if err != nil {
		c.utils.RespondInternalError(ctx, "删除发货单失败")
		return
	}

	c.utils.RespondOK(ctx, "删除发货单成功")
}

// List 获取发货单列表
// @Summary 获取发货单列表
// @Description 获取发货单列表，支持分页和筛选
// @Tags 发货单
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param search query string false "搜索关键字"
// @Param status query string false "状态"
// @Param date_from query string false "开始日期"
// @Param date_to query string false "结束日期"
// @Param customer_id query int false "客户ID"
// @Success 200 {object} dto.PaginatedResponse{data=[]dto.DeliveryNoteResponse}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/delivery-notes [get]
func (c *DeliveryNoteController) List(ctx *gin.Context) {
	var req dto.DeliveryNoteListRequest

	// 使用 ControllerUtils 解析分页参数
	pagination := c.utils.ParsePaginationParams(ctx)
	req.PaginationRequest = *pagination

	// 解析筛选参数
	req.Search = ctx.Query("search")
	req.Status = ctx.Query("status")
	req.DateFrom = ctx.Query("date_from")
	req.DateTo = ctx.Query("date_to")

	if customerID := ctx.Query("customer_id"); customerID != "" {
		if cid, err := strconv.ParseUint(customerID, 10, 32); err == nil {
			id := uint(cid)
			req.CustomerID = &id
		}
	}

	deliveryNotes, total, err := c.deliveryNoteService.List(&req)
	if err != nil {
		c.utils.RespondInternalError(ctx, "获取发货单列表失败")
		return
	}

	// 转换为统一的分页响应格式
	pagination2 := c.utils.CreatePagination(req.Page, req.PageSize, total)
	c.utils.RespondPaginated(ctx, deliveryNotes, pagination2, "获取发货单列表成功")
}

// UpdateStatus 更新发货单状态
// @Summary 更新发货单状态
// @Description 更新发货单状态
// @Tags 发货单
// @Accept json
// @Produce json
// @Param id path int true "发货单ID"
// @Param request body dto.DeliveryNoteStatusUpdateRequest true "状态更新请求"
// @Success 200 {object} utils.Response{data=models.DeliveryNote}
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/delivery-notes/{id}/status [patch]
func (c *DeliveryNoteController) UpdateStatus(ctx *gin.Context) {
	id, ok := c.utils.ParseIDParam(ctx, "id")
	if !ok {
		return
	}

	var req dto.DeliveryNoteStatusUpdateRequest
	if !c.utils.BindAndValidateJSON(ctx, &req) {
		return
	}

	deliveryNote, err := c.deliveryNoteService.UpdateStatus(id, &req)
	if err != nil {
		c.utils.RespondInternalError(ctx, "更新状态失败")
		return
	}

	c.utils.RespondOK(ctx, deliveryNote)
}

// CreateFromSalesOrder 从销售订单创建发货单
// @Summary 从销售订单创建发货单
// @Description 从销售订单创建发货单
// @Tags 发货单
// @Accept json
// @Produce json
// @Param request body dto.DeliveryNoteBatchCreateRequest true "批量创建请求"
// @Success 201 {object} utils.Response{data=models.DeliveryNote}
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/delivery-notes/from-sales-order [post]
func (c *DeliveryNoteController) CreateFromSalesOrder(ctx *gin.Context) {
	var req dto.DeliveryNoteBatchCreateRequest
	if !c.utils.BindAndValidateJSON(ctx, &req) {
		return
	}

	// 获取当前用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		c.utils.RespondUnauthorized(ctx, "用户未登录")
		return
	}

	deliveryNote, err := c.deliveryNoteService.CreateFromSalesOrder(ctx, &req, userID.(uint))
	if err != nil {
		c.utils.RespondInternalError(ctx, "从销售订单创建发货单失败")
		return
	}

	c.utils.RespondCreated(ctx, deliveryNote)
}

// GetStatistics 获取发货单统计信息
// @Summary 获取发货单统计信息
// @Description 获取发货单统计信息
// @Tags 发货单
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response{data=dto.DeliveryNoteStatisticsResponse}
// @Failure 500 {object} utils.Response
// @Router /api/delivery-notes/statistics [get]
func (c *DeliveryNoteController) GetStatistics(ctx *gin.Context) {
	stats, err := c.deliveryNoteService.GetStatistics(ctx)
	if err != nil {
		c.utils.RespondInternalError(ctx, "获取统计信息失败")
		return
	}

	c.utils.RespondOK(ctx, stats)
}

// GetDeliveryTrend 获取发货趋势
// @Summary 获取发货趋势
// @Description 获取发货趋势数据
// @Tags 发货单
// @Accept json
// @Produce json
// @Param days query int false "天数" default(30)
// @Success 200 {object} utils.Response{data=[]dto.DeliveryTrendData}
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/delivery-notes/trend [get]
func (c *DeliveryNoteController) GetDeliveryTrend(ctx *gin.Context) {
	days := 30
	if daysStr := ctx.Query("days"); daysStr != "" {
		if d, err := strconv.Atoi(daysStr); err == nil && d > 0 {
			days = d
		}
	}

	trend, err := c.deliveryNoteService.GetDeliveryTrend(days)
	if err != nil {
		c.utils.RespondInternalError(ctx, "获取发货趋势失败")
		return
	}

	c.utils.RespondOK(ctx, trend)
}
