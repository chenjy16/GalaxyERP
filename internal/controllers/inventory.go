package controllers

import (
	"strconv"

	"github.com/galaxyerp/galaxyErp/internal/dto"
	"github.com/galaxyerp/galaxyErp/internal/services"
	"github.com/gin-gonic/gin"
)

// InventoryController 库存控制器
type InventoryController struct {
	itemService      services.ItemService
	stockService     services.StockService
	warehouseService services.WarehouseService
	movementService  services.MovementService
	utils            *ControllerUtils
}

// NewInventoryController 创建库存控制器实例
func NewInventoryController(
	itemService services.ItemService,
	stockService services.StockService,
	warehouseService services.WarehouseService,
	movementService services.MovementService,
) *InventoryController {
	return &InventoryController{
		itemService:      itemService,
		stockService:     stockService,
		warehouseService: warehouseService,
		movementService:  movementService,
		utils:            NewControllerUtils(),
	}
}

// CreateItem 创建物料
// @Summary 创建物料
// @Description 创建新物料
// @Tags 物料管理
// @Accept json
// @Produce json
// @Param item body dto.ItemCreateRequest true "物料信息"
// @Success 201 {object} dto.ItemResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/items [post]
func (c *InventoryController) CreateItem(ctx *gin.Context) {
	var req dto.ItemCreateRequest
	if !c.utils.BindAndValidateJSON(ctx, &req) {
		return
	}

	item, err := c.itemService.CreateItem(ctx.Request.Context(), &req)
	if err != nil {
		c.utils.RespondInternalError(ctx, "创建物料失败")
		return
	}

	c.utils.RespondCreated(ctx, item)
}

// GetItem 获取物料
// @Summary 获取物料
// @Description 根据ID获取物料信息
// @Tags 物料管理
// @Accept json
// @Produce json
// @Param id path int true "物料ID"
// @Success 200 {object} dto.ItemResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/items/{id} [get]
func (c *InventoryController) GetItem(ctx *gin.Context) {
	id, ok := c.utils.ParseIDParam(ctx, "id")
	if !ok {
		return
	}

	item, err := c.itemService.GetItem(ctx.Request.Context(), id)
	if err != nil {
		c.utils.RespondInternalError(ctx, "获取物料失败")
		return
	}

	c.utils.RespondOK(ctx, item)
}

// UpdateItem 更新物料
// @Summary 更新物料
// @Description 更新物料信息
// @Tags 物料管理
// @Accept json
// @Produce json
// @Param id path int true "物料ID"
// @Param item body dto.ItemUpdateRequest true "物料信息"
// @Success 200 {object} dto.BaseResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/items/{id} [put]
func (c *InventoryController) UpdateItem(ctx *gin.Context) {
	id, ok := c.utils.ParseIDParam(ctx, "id")
	if !ok {
		return
	}

	var req dto.ItemUpdateRequest
	if !c.utils.BindAndValidateJSON(ctx, &req) {
		return
	}

	_, err := c.itemService.UpdateItem(ctx.Request.Context(), id, &req)
	if err != nil {
		c.utils.RespondInternalError(ctx, "更新物料失败")
		return
	}

	c.utils.RespondSuccess(ctx, "更新物料成功")
}

// DeleteItem 删除物料
// @Summary 删除物料
// @Description 删除物料
// @Tags 物料管理
// @Accept json
// @Produce json
// @Param id path int true "物料ID"
// @Success 200 {object} dto.BaseResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/items/{id} [delete]
func (c *InventoryController) DeleteItem(ctx *gin.Context) {
	id, ok := c.utils.ParseIDParam(ctx, "id")
	if !ok {
		return
	}

	err := c.itemService.DeleteItem(ctx.Request.Context(), id)
	if err != nil {
		c.utils.RespondInternalError(ctx, "删除物料失败")
		return
	}

	c.utils.RespondSuccess(ctx, "删除物料成功")
}

// ListItems 获取物料列表
// @Summary 获取物料列表
// @Description 获取物料列表
// @Tags 物料管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} dto.PaginatedResponse[dto.ItemResponse]
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/items [get]
func (c *InventoryController) ListItems(ctx *gin.Context) {
	req := c.utils.ParsePaginationParams(ctx)

	response, err := c.itemService.GetItems(ctx.Request.Context(), req)
	if err != nil {
		c.utils.RespondInternalError(ctx, "获取物料列表失败")
		return
	}

	// 转换为统一的分页响应格式
	pagination := c.utils.CreatePagination(response.Page, response.Limit, response.Total)
	c.utils.RespondPaginated(ctx, response.Data, pagination, "获取物料列表成功")
}

// SearchItems 搜索物料
// @Summary 搜索物料
// @Description 搜索物料
// @Tags 物料管理
// @Accept json
// @Produce json
// @Param keyword query string false "搜索关键词"
// @Param category query string false "物料分类"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} dto.PaginatedResponse[dto.ItemResponse]
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/items/search [get]
func (c *InventoryController) SearchItems(ctx *gin.Context) {
	keyword := ctx.Query("keyword")
	req := c.utils.ParsePaginationParams(ctx)

	searchReq := &dto.ItemSearchRequest{
		SearchRequest: dto.SearchRequest{
			PaginationRequest: *req,
			Keyword:          keyword,
		},
	}

	response, err := c.itemService.SearchItems(ctx.Request.Context(), searchReq)
	if err != nil {
		c.utils.RespondInternalError(ctx, "搜索物料失败")
		return
	}

	// 使用统一的分页响应格式
	pagination := c.utils.CreatePagination(response.Page, response.Limit, response.Total)
	c.utils.RespondPaginated(ctx, response.Data, pagination, "搜索物料成功")
}

// GetStockByItemID 根据物料ID获取库存
// @Summary 根据物料ID获取库存
// @Description 根据物料ID获取库存信息
// @Tags 库存管理
// @Accept json
// @Produce json
// @Param item_id path int true "物料ID"
// @Success 200 {object} dto.BaseResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/stock/item/{item_id} [get]
func (c *InventoryController) GetStockByItemID(ctx *gin.Context) {
	itemID, ok := c.utils.ParseIDParam(ctx, "item_id")
	if !ok {
		return
	}

	stock, err := c.stockService.GetByItemID(ctx.Request.Context(), uint(itemID))
	if err != nil {
		c.utils.RespondInternalError(ctx, "获取库存失败")
		return
	}

	c.utils.RespondOK(ctx, stock)
}

// ListStocks 获取库存列表
func (c *InventoryController) ListStocks(ctx *gin.Context) {
	req := c.utils.ParsePaginationParams(ctx)

	response, err := c.stockService.List(ctx.Request.Context(), req)
	if err != nil {
		c.utils.RespondInternalError(ctx, "获取库存列表失败")
		return
	}

	// 转换为统一的分页响应格式
	pagination := c.utils.CreatePagination(response.Page, response.Limit, response.Total)
	c.utils.RespondPaginated(ctx, response.Data, pagination, "获取库存列表成功")
}

// CreateStock 创建库存
func (c *InventoryController) CreateStock(ctx *gin.Context) {
	var req dto.MovementCreateRequest
	if !c.utils.BindAndValidateJSON(ctx, &req) {
		return
	}

	stock, err := c.stockService.CreateStock(ctx.Request.Context(), &req)
	if err != nil {
		c.utils.RespondInternalError(ctx, "创建库存失败")
		return
	}

	c.utils.RespondCreated(ctx, stock)
}

// GetStock 获取单个库存
func (c *InventoryController) GetStock(ctx *gin.Context) {
	id, ok := c.utils.ParseIDParam(ctx, "id")
	if !ok {
		return
	}

	stock, err := c.stockService.GetByID(ctx.Request.Context(), uint(id))
	if err != nil {
		c.utils.RespondInternalError(ctx, "获取库存失败")
		return
	}

	c.utils.RespondOK(ctx, stock)
}

// UpdateStock 更新库存
func (c *InventoryController) UpdateStock(ctx *gin.Context) {
	id, ok := c.utils.ParseIDParam(ctx, "id")
	if !ok {
		return
	}

	var req struct {
		Quantity float64 `json:"quantity" validate:"required,min=0"`
	}
	if !c.utils.BindAndValidateJSON(ctx, &req) {
		return
	}

	stock, err := c.stockService.UpdateStock(ctx.Request.Context(), uint(id), req.Quantity)
	if err != nil {
		c.utils.RespondInternalError(ctx, "更新库存失败")
		return
	}

	c.utils.RespondOK(ctx, stock)
}

// DeleteStock 删除库存
func (c *InventoryController) DeleteStock(ctx *gin.Context) {
	id, ok := c.utils.ParseIDParam(ctx, "id")
	if !ok {
		return
	}

	_, err := c.stockService.Delete(ctx.Request.Context(), uint(id))
	if err != nil {
		c.utils.RespondInternalError(ctx, "删除库存失败")
		return
	}

	c.utils.RespondSuccess(ctx, "删除库存成功")
}

// 库存移动相关方法 - 占位符实现
// ListStockMovements 获取库存移动列表
// @Summary 获取库存移动列表
// @Description 获取库存移动列表，支持分页和筛选
// @Tags 库存移动
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param limit query int false "每页数量" default(10)
// @Param item_id query int false "物料ID"
// @Param warehouse_id query int false "仓库ID"
// @Param type query string false "移动类型"
// @Success 200 {object} dto.PaginatedResponse[dto.MovementResponse]
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/stock-movements [get]
func (c *InventoryController) ListStockMovements(ctx *gin.Context) {
	req := c.utils.ParsePaginationParams(ctx)
	page := req.Page
	limit := req.PageSize

	itemIDStr := ctx.Query("item_id")
	warehouseIDStr := ctx.Query("warehouse_id")
	movementType := ctx.Query("type")

	var result *dto.PaginatedResponse[dto.MovementResponse]
	var err error

	if itemIDStr != "" {
		itemID, parseErr := strconv.ParseUint(itemIDStr, 10, 32)
		if parseErr != nil {
			c.utils.RespondBadRequest(ctx, "物料ID格式错误")
			return
		}
		result, err = c.movementService.GetMovementsByItemID(ctx, uint(itemID), page, limit)
	} else if warehouseIDStr != "" {
		warehouseID, parseErr := strconv.ParseUint(warehouseIDStr, 10, 32)
		if parseErr != nil {
			c.utils.RespondBadRequest(ctx, "仓库ID格式错误")
			return
		}
		result, err = c.movementService.GetMovementsByWarehouseID(ctx, uint(warehouseID), page, limit)
	} else if movementType != "" {
		result, err = c.movementService.GetMovementsByType(ctx, movementType, page, limit)
	} else {
		// 如果没有筛选条件，返回所有库存移动记录
		result, err = c.movementService.ListMovements(ctx, page, limit)
	}

	if err != nil {
		c.utils.RespondInternalError(ctx, "获取库存移动列表失败")
		return
	}

	// 转换为统一的分页响应格式
	pagination := c.utils.CreatePagination(result.Page, result.Limit, result.Total)
	c.utils.RespondPaginated(ctx, result.Data, pagination, "获取库存移动列表成功")
}

// CreateStockMovement 创建库存移动
// @Summary 创建库存移动
// @Description 创建新的库存移动记录
// @Tags 库存移动
// @Accept json
// @Produce json
// @Param movement body dto.MovementCreateRequest true "库存移动信息"
// @Success 201 {object} dto.MovementResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/stock-movements [post]
func (c *InventoryController) CreateStockMovement(ctx *gin.Context) {
	var req dto.MovementCreateRequest
	if !c.utils.BindAndValidateJSON(ctx, &req) {
		return
	}

	movement, err := c.movementService.CreateMovement(ctx, req)
	if err != nil {
		c.utils.RespondInternalError(ctx, "创建库存移动失败: "+err.Error())
		return
	}

	c.utils.RespondCreated(ctx, movement)
}

// StockIn 入库操作
// @Summary 入库操作
// @Description 执行入库操作
// @Tags 库存移动
// @Accept json
// @Produce json
// @Param movement body dto.MovementCreateRequest true "入库信息"
// @Success 201 {object} dto.MovementResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/stock-movements/in [post]
func (c *InventoryController) StockIn(ctx *gin.Context) {
	var req dto.MovementCreateRequest
	if !c.utils.BindAndValidateJSON(ctx, &req) {
		return
	}

	// 强制设置为入库类型
	req.Type = "in"

	movement, err := c.movementService.CreateMovement(ctx, req)
	if err != nil {
		c.utils.RespondInternalError(ctx, "入库操作失败: "+err.Error())
		return
	}

	c.utils.RespondCreated(ctx, movement)
}

// StockOut 出库操作
// @Summary 出库操作
// @Description 执行出库操作
// @Tags 库存移动
// @Accept json
// @Produce json
// @Param movement body dto.MovementCreateRequest true "出库信息"
// @Success 201 {object} dto.MovementResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/stock-movements/out [post]
func (c *InventoryController) StockOut(ctx *gin.Context) {
	var req dto.MovementCreateRequest
	if !c.utils.BindAndValidateJSON(ctx, &req) {
		return
	}

	// 强制设置为出库类型
	req.Type = "out"

	movement, err := c.movementService.CreateMovement(ctx, req)
	if err != nil {
		c.utils.RespondInternalError(ctx, "出库操作失败: "+err.Error())
		return
	}

	c.utils.RespondCreated(ctx, movement)
}

// StockAdjustment 库存调整
// @Summary 库存调整
// @Description 执行库存调整操作
// @Tags 库存移动
// @Accept json
// @Produce json
// @Param movement body dto.MovementCreateRequest true "库存调整信息"
// @Success 201 {object} dto.MovementResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/stock-movements/adjustment [post]
func (c *InventoryController) StockAdjustment(ctx *gin.Context) {
	var req dto.MovementCreateRequest
	if !c.utils.BindAndValidateJSON(ctx, &req) {
		return
	}

	// 强制设置为调整类型
	req.Type = "adjustment"

	movement, err := c.movementService.CreateMovement(ctx, req)
	if err != nil {
		c.utils.RespondInternalError(ctx, "库存调整失败: "+err.Error())
		return
	}

	c.utils.RespondCreated(ctx, movement)
}

// StockTransfer 库存调拨
// @Summary 库存调拨
// @Description 执行库存调拨操作
// @Tags 库存移动
// @Accept json
// @Produce json
// @Param movement body dto.MovementCreateRequest true "库存调拨信息"
// @Success 201 {object} dto.MovementResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/stock-movements/transfer [post]
func (c *InventoryController) StockTransfer(ctx *gin.Context) {
	var req dto.MovementCreateRequest
	if !c.utils.BindAndValidateJSON(ctx, &req) {
		return
	}

	// 强制设置为调拨类型
	req.Type = "transfer"

	movement, err := c.movementService.CreateMovement(ctx, req)
	if err != nil {
		c.utils.RespondInternalError(ctx, "库存调拨失败: "+err.Error())
		return
	}

	c.utils.RespondCreated(ctx, movement)
}

// 仓库管理相关方法
func (c *InventoryController) ListWarehouses(ctx *gin.Context) {
	req := c.utils.ParsePaginationParams(ctx)

	response, err := c.warehouseService.GetWarehouses(ctx.Request.Context(), req)
	if err != nil {
		c.utils.RespondInternalError(ctx, "获取仓库列表失败")
		return
	}

	// 转换为统一的分页响应格式
	pagination := c.utils.CreatePagination(response.Page, response.Limit, response.Total)
	c.utils.RespondPaginated(ctx, response.Data, pagination, "获取仓库列表成功")
}

func (c *InventoryController) CreateWarehouse(ctx *gin.Context) {
	var req dto.WarehouseCreateRequest
	if !c.utils.BindAndValidateJSON(ctx, &req) {
		return
	}

	warehouse, err := c.warehouseService.CreateWarehouse(ctx.Request.Context(), &req)
	if err != nil {
		c.utils.RespondInternalError(ctx, "创建仓库失败")
		return
	}

	c.utils.RespondCreated(ctx, warehouse)
}

func (c *InventoryController) GetWarehouse(ctx *gin.Context) {
	id, ok := c.utils.ParseIDParam(ctx, "id")
	if !ok {
		return
	}

	warehouse, err := c.warehouseService.GetWarehouse(ctx.Request.Context(), uint(id))
	if err != nil {
		c.utils.RespondInternalError(ctx, "获取仓库失败")
		return
	}

	c.utils.RespondOK(ctx, warehouse)
}

func (c *InventoryController) UpdateWarehouse(ctx *gin.Context) {
	id, ok := c.utils.ParseIDParam(ctx, "id")
	if !ok {
		return
	}

	var req dto.WarehouseUpdateRequest
	if !c.utils.BindAndValidateJSON(ctx, &req) {
		return
	}

	warehouse, err := c.warehouseService.UpdateWarehouse(ctx.Request.Context(), uint(id), &req)
	if err != nil {
		c.utils.RespondInternalError(ctx, "更新仓库失败")
		return
	}

	c.utils.RespondOK(ctx, warehouse)
}

func (c *InventoryController) DeleteWarehouse(ctx *gin.Context) {
	id, ok := c.utils.ParseIDParam(ctx, "id")
	if !ok {
		return
	}

	err := c.warehouseService.DeleteWarehouse(ctx.Request.Context(), uint(id))
	if err != nil {
		c.utils.RespondInternalError(ctx, "删除仓库失败")
		return
	}

	c.utils.RespondSuccess(ctx, "删除仓库成功")
}

// 库存报告和统计相关方法
func (c *InventoryController) GetInventoryStats(ctx *gin.Context) {
	c.utils.RespondNotImplemented(ctx, "功能暂未实现")
}

func (c *InventoryController) GetInventoryReport(ctx *gin.Context) {
	c.utils.RespondNotImplemented(ctx, "功能暂未实现")
}

func (c *InventoryController) GetABCAnalysis(ctx *gin.Context) {
	c.utils.RespondNotImplemented(ctx, "功能暂未实现")
}

func (c *InventoryController) ExportInventoryReport(ctx *gin.Context) {
	c.utils.RespondNotImplemented(ctx, "功能暂未实现")
}
