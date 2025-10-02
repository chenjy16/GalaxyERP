package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/galaxyerp/galaxyErp/internal/dto"
	"github.com/galaxyerp/galaxyErp/internal/services"
	"github.com/galaxyerp/galaxyErp/internal/utils"
)

// DeliveryNoteController 发货单控制器
type DeliveryNoteController struct {
	deliveryNoteService services.DeliveryNoteServiceInterface
}

// NewDeliveryNoteController 创建发货单控制器
func NewDeliveryNoteController(deliveryNoteService services.DeliveryNoteServiceInterface) *DeliveryNoteController {
	return &DeliveryNoteController{
		deliveryNoteService: deliveryNoteService,
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
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "请求参数错误")
		return
	}

	// 获取当前用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "用户未登录")
		return
	}

	deliveryNote, err := c.deliveryNoteService.Create(ctx, &req, userID.(uint))
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "创建发货单失败")
		return
	}

	utils.CreatedResponse(ctx, deliveryNote, "创建发货单成功")
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
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "无效的发货单ID")
		return
	}

	deliveryNote, err := c.deliveryNoteService.GetByID(uint(id))
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, "发货单不存在")
		return
	}

	utils.SuccessResponse(ctx, deliveryNote, "获取发货单成功")
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
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "无效的发货单ID")
		return
	}

	var req dto.DeliveryNoteUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "请求参数错误")
		return
	}

	deliveryNote, err := c.deliveryNoteService.Update(uint(id), &req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "更新发货单失败")
		return
	}

	utils.SuccessResponse(ctx, deliveryNote, "更新发货单成功")
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
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "无效的发货单ID")
		return
	}

	err = c.deliveryNoteService.Delete(uint(id))
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "删除发货单失败")
		return
	}

	utils.SuccessResponse(ctx, "删除发货单成功")
}

// List 获取发货单列表
// @Summary 获取发货单列表
// @Description 获取发货单列表，支持分页和筛选
// @Tags 发货单
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param search query string false "搜索关键词"
// @Param customer_id query int false "客户ID"
// @Param status query string false "状态"
// @Param date_from query string false "开始日期"
// @Param date_to query string false "结束日期"
// @Success 200 {object} utils.Response{data=dto.PaginatedResponse}
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/delivery-notes [get]
func (c *DeliveryNoteController) List(ctx *gin.Context) {
	var req dto.DeliveryNoteListRequest

	// 解析分页参数
	if page := ctx.Query("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil && p > 0 {
			req.Page = p
		}
	}
	if req.Page == 0 {
		req.Page = 1
	}

	if pageSize := ctx.Query("page_size"); pageSize != "" {
		if ps, err := strconv.Atoi(pageSize); err == nil && ps > 0 {
			req.PageSize = ps
		}
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}

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
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "获取发货单列表失败")
		return
	}

	utils.PaginatedResponse(ctx, deliveryNotes, req.Page, req.PageSize, total, "获取发货单列表成功")
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
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "无效的发货单ID")
		return
	}

	var req dto.DeliveryNoteStatusUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "请求参数错误")
		return
	}

	deliveryNote, err := c.deliveryNoteService.UpdateStatus(uint(id), &req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "更新状态失败")
		return
	}

	utils.SuccessResponse(ctx, deliveryNote, "更新状态成功")
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
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "请求参数错误")
		return
	}

	// 获取当前用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "用户未登录")
		return
	}

	deliveryNote, err := c.deliveryNoteService.CreateFromSalesOrder(ctx, &req, userID.(uint))
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "从销售订单创建发货单失败")
		return
	}

	utils.CreatedResponse(ctx, deliveryNote, "从销售订单创建发货单成功")
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
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "获取统计信息失败")
		return
	}

	utils.SuccessResponse(ctx, stats, "获取统计信息成功")
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
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "获取发货趋势失败")
		return
	}

	utils.SuccessResponse(ctx, trend, "获取发货趋势成功")
}