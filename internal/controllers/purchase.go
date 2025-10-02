package controllers

import (
	"net/http"
	"strconv"

	"github.com/galaxyerp/galaxyErp/internal/dto"
	"github.com/galaxyerp/galaxyErp/internal/services"
	"github.com/gin-gonic/gin"
)

// PurchaseController 采购管理控制器
type PurchaseController struct {
	utils                  *ControllerUtils
	supplierService        services.SupplierService
	purchaseRequestService services.PurchaseRequestService
	purchaseOrderService   services.PurchaseOrderService
}

// NewPurchaseController 创建采购管理控制器实例
func NewPurchaseController(
	utils *ControllerUtils,
	supplierService services.SupplierService,
	purchaseRequestService services.PurchaseRequestService,
	purchaseOrderService services.PurchaseOrderService,
) *PurchaseController {
	return &PurchaseController{
		utils:                  utils,
		supplierService:        supplierService,
		purchaseRequestService: purchaseRequestService,
		purchaseOrderService:   purchaseOrderService,
	}
}

// ===== 供应商管理 =====

// CreateSupplier 创建供应商
func (c *PurchaseController) CreateSupplier(ctx *gin.Context) {
	var req dto.SupplierCreateRequest
	if !c.utils.BindJSON(ctx, &req) {
		return
	}

	supplier, err := c.supplierService.CreateSupplier(ctx.Request.Context(), &req)
	if err != nil {
		c.utils.RespondInternalError(ctx, "创建供应商失败")
		return
	}

	c.utils.RespondCreated(ctx, supplier)
}

// GetSupplier 获取供应商详情
func (c *PurchaseController) GetSupplier(ctx *gin.Context) {
	id, ok := c.utils.ParseIDParam(ctx, "id")
	if !ok {
		return
	}

	supplier, err := c.supplierService.GetSupplier(ctx.Request.Context(), id)
	if err != nil {
		c.utils.RespondNotFound(ctx, "供应商不存在")
		return
	}

	c.utils.RespondOK(ctx, supplier)
}

// UpdateSupplier 更新供应商
func (c *PurchaseController) UpdateSupplier(ctx *gin.Context) {
	id, ok := c.utils.ParseIDParam(ctx, "id")
	if !ok {
		return
	}

	var req dto.SupplierUpdateRequest
	if !c.utils.BindJSON(ctx, &req) {
		return
	}

	supplier, err := c.supplierService.UpdateSupplier(ctx.Request.Context(), id, &req)
	if err != nil {
		c.utils.RespondInternalError(ctx, "更新供应商失败")
		return
	}

	c.utils.RespondOK(ctx, supplier)
}

// DeleteSupplier 删除供应商
func (c *PurchaseController) DeleteSupplier(ctx *gin.Context) {
	id, ok := c.utils.ParseIDParam(ctx, "id")
	if !ok {
		return
	}

	if err := c.supplierService.DeleteSupplier(ctx.Request.Context(), id); err != nil {
		c.utils.RespondInternalError(ctx, "删除供应商失败")
		return
	}

	c.utils.RespondSuccess(ctx, "供应商删除成功")
}

// ListSuppliers 获取供应商列表
func (c *PurchaseController) ListSuppliers(ctx *gin.Context) {
	var filter dto.SupplierFilter
	if err := ctx.ShouldBindQuery(&filter); err != nil {
		c.utils.RespondBadRequest(ctx, "查询参数无效")
		return
	}

	suppliers, err := c.supplierService.ListSuppliers(ctx.Request.Context(), &filter)
	if err != nil {
		c.utils.RespondInternalError(ctx, "获取供应商列表失败")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    suppliers,
	})
}

// ===== 采购申请管理 =====

// CreatePurchaseRequest 创建采购申请
func (c *PurchaseController) CreatePurchaseRequest(ctx *gin.Context) {
	var req dto.PurchaseRequestCreateRequest
	if !c.utils.BindJSON(ctx, &req) {
		return
	}

	purchaseRequest, err := c.purchaseRequestService.CreatePurchaseRequest(ctx.Request.Context(), &req)
	if err != nil {
		c.utils.RespondInternalError(ctx, "创建采购申请失败")
		return
	}

	c.utils.RespondCreated(ctx, purchaseRequest)
}

// GetPurchaseRequest 获取采购申请详情
func (c *PurchaseController) GetPurchaseRequest(ctx *gin.Context) {
	id, ok := c.utils.ParseIDParam(ctx, "id")
	if !ok {
		return
	}

	purchaseRequest, err := c.purchaseRequestService.GetPurchaseRequest(ctx.Request.Context(), id)
	if err != nil {
		c.utils.RespondNotFound(ctx, "采购申请不存在")
		return
	}

	c.utils.RespondOK(ctx, purchaseRequest)
}

// UpdatePurchaseRequest 更新采购申请
func (c *PurchaseController) UpdatePurchaseRequest(ctx *gin.Context) {
	id, ok := c.utils.ParseIDParam(ctx, "id")
	if !ok {
		return
	}

	var req dto.PurchaseRequestUpdateRequest
	if !c.utils.BindJSON(ctx, &req) {
		return
	}

	purchaseRequest, err := c.purchaseRequestService.UpdatePurchaseRequest(ctx.Request.Context(), id, &req)
	if err != nil {
		c.utils.RespondInternalError(ctx, "更新采购申请失败")
		return
	}

	c.utils.RespondOK(ctx, purchaseRequest)
}

// DeletePurchaseRequest 删除采购申请
func (c *PurchaseController) DeletePurchaseRequest(ctx *gin.Context) {
	id, ok := c.utils.ParseIDParam(ctx, "id")
	if !ok {
		return
	}

	if err := c.purchaseRequestService.DeletePurchaseRequest(ctx.Request.Context(), id); err != nil {
		c.utils.RespondInternalError(ctx, "删除采购申请失败")
		return
	}

	c.utils.RespondSuccess(ctx, "采购申请删除成功")
}

// ListPurchaseRequests 获取采购申请列表
func (c *PurchaseController) ListPurchaseRequests(ctx *gin.Context) {
	var filter dto.PurchaseSearchRequest
	if err := ctx.ShouldBindQuery(&filter); err != nil {
		c.utils.RespondBadRequest(ctx, "查询参数无效")
		return
	}

	purchaseRequests, err := c.purchaseRequestService.ListPurchaseRequests(ctx.Request.Context(), &filter)
	if err != nil {
		c.utils.RespondInternalError(ctx, "获取采购申请列表失败")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    purchaseRequests,
	})
}

// SubmitPurchaseRequest 提交采购申请
func (c *PurchaseController) SubmitPurchaseRequest(ctx *gin.Context) {
	id, ok := c.utils.ParseIDParam(ctx, "id")
	if !ok {
		return
	}

	if err := c.purchaseRequestService.SubmitPurchaseRequest(ctx.Request.Context(), id); err != nil {
		c.utils.RespondInternalError(ctx, "提交采购申请失败")
		return
	}

	c.utils.RespondSuccess(ctx, "采购申请提交成功")
}

// ApprovePurchaseRequest 审批采购申请
func (c *PurchaseController) ApprovePurchaseRequest(ctx *gin.Context) {
	id, ok := c.utils.ParseIDParam(ctx, "id")
	if !ok {
		return
	}

	// 从上下文获取用户ID（需要中间件支持）
	userID := uint(1) // 临时硬编码，实际应从JWT token获取

	if err := c.purchaseRequestService.ApprovePurchaseRequest(ctx.Request.Context(), id, userID); err != nil {
		c.utils.RespondInternalError(ctx, "审批采购申请失败")
		return
	}

	c.utils.RespondSuccess(ctx, "采购申请审批成功")
}

// RejectPurchaseRequest 拒绝采购申请
func (c *PurchaseController) RejectPurchaseRequest(ctx *gin.Context) {
	id, ok := c.utils.ParseIDParam(ctx, "id")
	if !ok {
		return
	}

	var req struct {
		Reason string `json:"reason" binding:"required"`
	}
	if !c.utils.BindJSON(ctx, &req) {
		return
	}

	// 从上下文获取用户ID（需要中间件支持）
	userID := uint(1) // 临时硬编码，实际应从JWT token获取

	if err := c.purchaseRequestService.RejectPurchaseRequest(ctx.Request.Context(), id, userID, req.Reason); err != nil {
		c.utils.RespondInternalError(ctx, "拒绝采购申请失败")
		return
	}

	c.utils.RespondSuccess(ctx, "采购申请拒绝成功")
}

// ===== 采购订单管理 =====

// CreatePurchaseOrder 创建采购订单
func (c *PurchaseController) CreatePurchaseOrder(ctx *gin.Context) {
	var req dto.PurchaseOrderCreateRequest
	if !c.utils.BindJSON(ctx, &req) {
		return
	}

	purchaseOrder, err := c.purchaseOrderService.CreatePurchaseOrder(ctx.Request.Context(), &req)
	if err != nil {
		c.utils.RespondInternalError(ctx, "创建采购订单失败")
		return
	}

	c.utils.RespondCreated(ctx, purchaseOrder)
}

// GetPurchaseOrder 获取采购订单详情
func (c *PurchaseController) GetPurchaseOrder(ctx *gin.Context) {
	id, ok := c.utils.ParseIDParam(ctx, "id")
	if !ok {
		return
	}

	purchaseOrder, err := c.purchaseOrderService.GetPurchaseOrder(ctx.Request.Context(), id)
	if err != nil {
		c.utils.RespondNotFound(ctx, "采购订单不存在")
		return
	}

	c.utils.RespondOK(ctx, purchaseOrder)
}

// UpdatePurchaseOrder 更新采购订单
func (c *PurchaseController) UpdatePurchaseOrder(ctx *gin.Context) {
	id, ok := c.utils.ParseIDParam(ctx, "id")
	if !ok {
		return
	}

	var req dto.PurchaseOrderUpdateRequest
	if !c.utils.BindJSON(ctx, &req) {
		return
	}

	purchaseOrder, err := c.purchaseOrderService.UpdatePurchaseOrder(ctx.Request.Context(), id, &req)
	if err != nil {
		c.utils.RespondInternalError(ctx, "更新采购订单失败")
		return
	}

	c.utils.RespondOK(ctx, purchaseOrder)
}

// DeletePurchaseOrder 删除采购订单
func (c *PurchaseController) DeletePurchaseOrder(ctx *gin.Context) {
	id, ok := c.utils.ParseIDParam(ctx, "id")
	if !ok {
		return
	}

	if err := c.purchaseOrderService.DeletePurchaseOrder(ctx.Request.Context(), id); err != nil {
		c.utils.RespondInternalError(ctx, "删除采购订单失败")
		return
	}

	c.utils.RespondSuccess(ctx, "采购订单删除成功")
}

// ListPurchaseOrders 获取采购订单列表
func (c *PurchaseController) ListPurchaseOrders(ctx *gin.Context) {
	var filter dto.PurchaseOrderFilter
	if err := ctx.ShouldBindQuery(&filter); err != nil {
		c.utils.RespondBadRequest(ctx, "无效的查询参数")
		return
	}

	result, err := c.purchaseOrderService.ListPurchaseOrders(ctx, &filter)
	if err != nil {
		c.utils.RespondInternalError(ctx, "获取采购订单列表失败")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// ConfirmPurchaseOrder 确认采购订单
func (c *PurchaseController) ConfirmPurchaseOrder(ctx *gin.Context) {
	id, ok := c.utils.ParseIDParam(ctx, "id")
	if !ok {
		return
	}

	if err := c.purchaseOrderService.ConfirmPurchaseOrder(ctx.Request.Context(), id); err != nil {
		c.utils.RespondInternalError(ctx, "确认采购订单失败")
		return
	}

	c.utils.RespondSuccess(ctx, "采购订单确认成功")
}

// CancelPurchaseOrder 取消采购订单
func (c *PurchaseController) CancelPurchaseOrder(ctx *gin.Context) {
	id, ok := c.utils.ParseIDParam(ctx, "id")
	if !ok {
		return
	}

	if err := c.purchaseOrderService.CancelPurchaseOrder(ctx.Request.Context(), id); err != nil {
		c.utils.RespondInternalError(ctx, "取消采购订单失败")
		return
	}

	c.utils.RespondSuccess(ctx, "采购订单取消成功")
}

// ===== 统计和报表 =====

// GetPurchaseStats 获取采购统计信息
func (c *PurchaseController) GetPurchaseStats(ctx *gin.Context) {
	// 获取查询参数
	startDate := ctx.Query("start_date")
	endDate := ctx.Query("end_date")

	// 这里应该调用相应的服务方法获取统计数据
	// 暂时返回模拟数据
	stats := map[string]interface{}{
		"total_suppliers":         50,
		"active_suppliers":        45,
		"total_purchase_requests": 120,
		"pending_requests":        15,
		"approved_requests":       95,
		"rejected_requests":       10,
		"total_purchase_orders":   80,
		"pending_orders":          8,
		"confirmed_orders":        65,
		"completed_orders":        7,
		"total_amount":            1250000.00,
		"period": map[string]string{
			"start_date": startDate,
			"end_date":   endDate,
		},
	}

	c.utils.RespondOK(ctx, stats)
}

// GetSupplierPerformance 获取供应商绩效
func (c *PurchaseController) GetSupplierPerformance(ctx *gin.Context) {
	supplierIDStr := ctx.Param("supplier_id")
	supplierID, err := strconv.ParseUint(supplierIDStr, 10, 32)
	if err != nil {
		c.utils.RespondBadRequest(ctx, "无效的供应商ID")
		return
	}

	// 这里应该调用相应的服务方法获取供应商绩效数据
	// 暂时返回模拟数据
	performance := map[string]interface{}{
		"supplier_id":       uint(supplierID),
		"total_orders":      25,
		"completed_orders":  23,
		"on_time_delivery":  0.92,
		"quality_rating":    4.5,
		"total_amount":      350000.00,
		"average_lead_time": 7.5,
		"return_rate":       0.02,
	}

	c.utils.RespondOK(ctx, performance)
}
