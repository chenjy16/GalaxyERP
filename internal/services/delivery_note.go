package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/galaxyerp/galaxyErp/internal/common"
	"github.com/galaxyerp/galaxyErp/internal/dto"
	"github.com/galaxyerp/galaxyErp/internal/models"
	"github.com/galaxyerp/galaxyErp/internal/repositories"
)

// DeliveryNoteServiceInterface 发货单服务接口
type DeliveryNoteServiceInterface interface {
	Create(ctx *gin.Context, req *dto.DeliveryNoteCreateRequest, userID uint) (*models.DeliveryNote, error)
	GetByID(id uint) (*models.DeliveryNote, error)
	Update(id uint, req *dto.DeliveryNoteUpdateRequest) (*models.DeliveryNote, error)
	Delete(id uint) error
	List(req *dto.DeliveryNoteListRequest) ([]*models.DeliveryNote, int64, error)
	UpdateStatus(id uint, req *dto.DeliveryNoteStatusUpdateRequest) (*models.DeliveryNote, error)
	CreateFromSalesOrder(ctx *gin.Context, req *dto.DeliveryNoteBatchCreateRequest, userID uint) (*models.DeliveryNote, error)
	GetStatistics(ctx *gin.Context) (*dto.DeliveryNoteStatisticsResponse, error)
	GetDeliveryTrend(days int) ([]dto.DeliveryTrendData, error)
}

type DeliveryNoteService struct {
	deliveryNoteRepo repositories.DeliveryNoteRepository
	salesOrderRepo   repositories.SalesOrderRepository
	customerRepo     repositories.CustomerRepository
}

func NewDeliveryNoteService(
	deliveryNoteRepo repositories.DeliveryNoteRepository,
	salesOrderRepo repositories.SalesOrderRepository,
	customerRepo repositories.CustomerRepository,
) *DeliveryNoteService {
	return &DeliveryNoteService{
		deliveryNoteRepo: deliveryNoteRepo,
		salesOrderRepo:   salesOrderRepo,
		customerRepo:     customerRepo,
	}
}

// Create 创建发货单
func (s *DeliveryNoteService) Create(ctx *gin.Context, req *dto.DeliveryNoteCreateRequest, userID uint) (*models.DeliveryNote, error) {
	// 验证客户是否存在
	_, err := s.customerRepo.GetByID(ctx, req.CustomerID)
	if err != nil {
		return nil, fmt.Errorf("客户不存在: %w", err)
	}

	// 如果指定了销售订单，验证销售订单是否存在且属于该客户
	var salesOrder *models.SalesOrder
	if req.SalesOrderID != nil {
		salesOrder, err = s.salesOrderRepo.GetByID(ctx, *req.SalesOrderID)
		if err != nil {
			return nil, fmt.Errorf("销售订单不存在: %w", err)
		}
		if salesOrder.CustomerID != req.CustomerID {
			return nil, errors.New("销售订单与客户不匹配")
		}
	}

	// 生成发货单号
	deliveryNumber, err := s.deliveryNoteRepo.GenerateDeliveryNumber()
	if err != nil {
		return nil, fmt.Errorf("生成发货单号失败: %w", err)
	}

	// 计算总数量
	var totalQuantity float64
	for _, item := range req.Items {
		totalQuantity += item.Quantity
	}

	// 创建发货单
	deliveryNote := &models.DeliveryNote{
		DeliveryNumber: deliveryNumber,
		CustomerID:     req.CustomerID,
		SalesOrderID:   req.SalesOrderID,
		Date:           req.Date,
		Status:         "Draft", // 默认状态为草稿
		TotalQuantity:  totalQuantity,
		Transporter:    req.Transporter,
		DriverName:     req.DriverName,
		VehicleNumber:  req.VehicleNumber,
		Destination:    req.Destination,
		Notes:          req.Notes,
		CreatedBy:      userID,
	}

	// 创建发货单明细
	for _, itemReq := range req.Items {
		item := models.DeliveryNoteItem{
			SalesOrderItemID: itemReq.SalesOrderItemID,
			ItemID:           itemReq.ItemID,
			Description:      itemReq.Description,
			Quantity:         itemReq.Quantity,
			BatchNo:          itemReq.BatchNo,
			SerialNo:         itemReq.SerialNo,
			WarehouseID:      itemReq.WarehouseID,
		}
		deliveryNote.Items = append(deliveryNote.Items, item)
	}

	// 保存到数据库
	if err := s.deliveryNoteRepo.Create(ctx, deliveryNote); err != nil {
		return nil, fmt.Errorf("创建发货单失败: %w", err)
	}

	// 重新加载完整数据
	return s.deliveryNoteRepo.GetByID(ctx, deliveryNote.ID)
}

// GetByID 根据ID获取发货单
func (s *DeliveryNoteService) GetByID(id uint) (*models.DeliveryNote, error) {
	return s.deliveryNoteRepo.GetByID(context.Background(), id)
}

// Update 更新发货单
func (s *DeliveryNoteService) Update(id uint, req *dto.DeliveryNoteUpdateRequest) (*models.DeliveryNote, error) {
	// 获取现有发货单
	deliveryNote, err := s.deliveryNoteRepo.GetByID(context.Background(), id)
	if err != nil {
		return nil, fmt.Errorf("发货单不存在: %w", err)
	}

	// 检查状态是否允许修改
	if deliveryNote.Status == "Delivered" || deliveryNote.Status == "Cancelled" {
		return nil, errors.New("已发货或已取消的发货单不能修改")
	}

	// 更新字段
	if req.Date != nil {
		deliveryNote.Date = *req.Date
	}
	if req.Transporter != "" {
		deliveryNote.Transporter = req.Transporter
	}
	if req.DriverName != "" {
		deliveryNote.DriverName = req.DriverName
	}
	if req.VehicleNumber != "" {
		deliveryNote.VehicleNumber = req.VehicleNumber
	}
	if req.Destination != "" {
		deliveryNote.Destination = req.Destination
	}
	if req.Notes != "" {
		deliveryNote.Notes = req.Notes
	}

	// 如果有明细更新，重新计算总数量
	if len(req.Items) > 0 {
		// 清空现有明细
		deliveryNote.Items = nil
		
		// 创建新明细
		var totalQuantity float64
		for _, itemReq := range req.Items {
			item := models.DeliveryNoteItem{
				SalesOrderItemID: itemReq.SalesOrderItemID,
				ItemID:           itemReq.ItemID,
				Description:      itemReq.Description,
				Quantity:         itemReq.Quantity,
				BatchNo:          itemReq.BatchNo,
				SerialNo:         itemReq.SerialNo,
				WarehouseID:      itemReq.WarehouseID,
			}
			deliveryNote.Items = append(deliveryNote.Items, item)
			totalQuantity += itemReq.Quantity
		}

		deliveryNote.TotalQuantity = totalQuantity
	}

	// 保存更新
	if err := s.deliveryNoteRepo.Update(context.Background(), deliveryNote); err != nil {
		return nil, fmt.Errorf("更新发货单失败: %w", err)
	}

	// 重新加载完整数据
	return s.deliveryNoteRepo.GetByID(context.Background(), id)
}

// Delete 删除发货单
func (s *DeliveryNoteService) Delete(id uint) error {
	// 获取发货单
	deliveryNote, err := s.deliveryNoteRepo.GetByID(context.Background(), id)
	if err != nil {
		return fmt.Errorf("发货单不存在: %w", err)
	}

	// 检查状态是否允许删除
	if deliveryNote.Status == "Delivered" {
		return errors.New("已发货的发货单不能删除")
	}

	return s.deliveryNoteRepo.Delete(context.Background(), id)
}

// List 获取发货单列表
func (s *DeliveryNoteService) List(req *dto.DeliveryNoteListRequest) ([]*models.DeliveryNote, int64, error) {
	options := &common.QueryOptions{
		Pagination: &dto.PaginationRequest{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
	}
	return s.deliveryNoteRepo.List(context.Background(), options)
}

// UpdateStatus 更新发货单状态
func (s *DeliveryNoteService) UpdateStatus(id uint, req *dto.DeliveryNoteStatusUpdateRequest) (*models.DeliveryNote, error) {
	// 获取发货单
	deliveryNote, err := s.deliveryNoteRepo.GetByID(context.Background(), id)
	if err != nil {
		return nil, fmt.Errorf("发货单不存在: %w", err)
	}

	// 验证状态转换
	if err := s.validateStatusTransition(deliveryNote.Status, req.Status); err != nil {
		return nil, err
	}

	// 更新状态
	if err := s.deliveryNoteRepo.UpdateStatus(context.Background(), id, req.Status); err != nil {
		return nil, fmt.Errorf("更新状态失败: %w", err)
	}

	// 重新加载数据
	return s.deliveryNoteRepo.GetByID(context.Background(), id)
}

// CreateFromSalesOrder 从销售订单创建发货单
func (s *DeliveryNoteService) CreateFromSalesOrder(ctx *gin.Context, req *dto.DeliveryNoteBatchCreateRequest, userID uint) (*models.DeliveryNote, error) {
	// 获取销售订单
	salesOrder, err := s.salesOrderRepo.GetByID(ctx, req.SalesOrderID)
	if err != nil {
		return nil, fmt.Errorf("销售订单不存在: %w", err)
	}

	// 检查销售订单状态
	if salesOrder.Status != "Confirmed" {
		return nil, errors.New("只能从已确认的销售订单创建发货单")
	}

	// 生成发货单号
	deliveryNumber, err := s.deliveryNoteRepo.GenerateDeliveryNumber()
	if err != nil {
		return nil, fmt.Errorf("生成发货单号失败: %w", err)
	}

	// 计算总数量
	var totalQuantity float64
	for _, item := range req.Items {
		totalQuantity += item.Quantity
	}

	// 创建发货单
	deliveryNote := &models.DeliveryNote{
		DeliveryNumber: deliveryNumber,
		CustomerID:     salesOrder.CustomerID,
		SalesOrderID:   &req.SalesOrderID,
		Date:           req.Date,
		Status:         "Draft",
		TotalQuantity:  totalQuantity,
		Transporter:    req.Transporter,
		DriverName:     req.DriverName,
		VehicleNumber:  req.VehicleNumber,
		Destination:    req.Destination,
		Notes:          req.Notes,
		CreatedBy:      userID,
	}

	// 创建发货单明细
	for _, itemReq := range req.Items {
		// 从销售订单明细中查找对应的明细
		var salesOrderItem *models.SalesOrderItem
		for _, item := range salesOrder.Items {
			if item.ID == itemReq.SalesOrderItemID {
				salesOrderItem = &item
				break
			}
		}
		
		if salesOrderItem == nil {
			return nil, fmt.Errorf("销售订单明细不存在: %w", err)
		}

		// 检查发货数量是否超过订单数量
		if itemReq.Quantity > salesOrderItem.Quantity {
			return nil, fmt.Errorf("发货数量不能超过订单数量")
		}

		item := models.DeliveryNoteItem{
			SalesOrderItemID: &itemReq.SalesOrderItemID,
			ItemID:           salesOrderItem.ItemID,
			Description:      salesOrderItem.Description,
			Quantity:         itemReq.Quantity,
			BatchNo:          itemReq.BatchNo,
			SerialNo:         itemReq.SerialNo,
			WarehouseID:      itemReq.WarehouseID,
		}
		deliveryNote.Items = append(deliveryNote.Items, item)
	}

	// 保存到数据库
	if err := s.deliveryNoteRepo.Create(ctx, deliveryNote); err != nil {
		return nil, fmt.Errorf("创建发货单失败: %w", err)
	}

	// 重新加载完整数据
	return s.deliveryNoteRepo.GetByID(ctx, deliveryNote.ID)
}

// GetStatistics 获取发货单统计信息
func (s *DeliveryNoteService) GetStatistics(ctx *gin.Context) (*dto.DeliveryNoteStatisticsResponse, error) {
	return s.deliveryNoteRepo.GetStatistics(ctx)
}

// GetDeliveryTrend 获取发货趋势数据
func (s *DeliveryNoteService) GetDeliveryTrend(days int) ([]dto.DeliveryTrendData, error) {
	return s.deliveryNoteRepo.GetDeliveryTrend(days)
}

// validateStatusTransition 验证状态转换
func (s *DeliveryNoteService) validateStatusTransition(currentStatus, newStatus string) error {
	validTransitions := map[string][]string{
		"Draft":     {"Submitted", "Cancelled"},
		"Submitted": {"Delivered", "Cancelled"},
		"Delivered": {}, // 已发货状态不能转换
		"Cancelled": {}, // 已取消状态不能转换
	}

	allowedStatuses, exists := validTransitions[currentStatus]
	if !exists {
		return fmt.Errorf("无效的当前状态: %s", currentStatus)
	}

	for _, status := range allowedStatuses {
		if status == newStatus {
			return nil
		}
	}

	return fmt.Errorf("不能从状态 %s 转换到 %s", currentStatus, newStatus)
}
