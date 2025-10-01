package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/galaxyerp/galaxyErp/internal/dto"
	"github.com/galaxyerp/galaxyErp/internal/models"
	"github.com/galaxyerp/galaxyErp/internal/repositories"
)

// ItemService 物料服务接口
type ItemService interface {
	CreateItem(ctx context.Context, req *dto.ItemCreateRequest) (*dto.ItemResponse, error)
	GetItem(ctx context.Context, id uint) (*dto.ItemResponse, error)
	UpdateItem(ctx context.Context, id uint, req *dto.ItemUpdateRequest) (*dto.ItemResponse, error)
	DeleteItem(ctx context.Context, id uint) error
	GetItems(ctx context.Context, req *dto.PaginationRequest) (*dto.BaseResponse, error)
	SearchItems(ctx context.Context, req *dto.ItemSearchRequest) (*dto.BaseResponse, error)
}

// ItemServiceImpl 物料服务实现
type ItemServiceImpl struct {
	itemRepo repositories.ItemRepository
}

// NewItemService 创建物料服务实例
func NewItemService(itemRepo repositories.ItemRepository) ItemService {
	return &ItemServiceImpl{
		itemRepo: itemRepo,
	}
}

// CreateItem 创建物料
func (s *ItemServiceImpl) CreateItem(ctx context.Context, req *dto.ItemCreateRequest) (*dto.ItemResponse, error) {
	// 创建物料
	item := &models.Item{
		Code:         req.Code,
		Name:         req.Name,
		Description:  req.Description,
		// TODO: 需要根据CategoryID和UnitID查询对应的名称
		// Category:     req.Category,
		// Unit:         req.Unit,
		Cost:         req.UnitCost,
		Price:        req.SalePrice,
		ReorderLevel: int(req.MinStock),
		IsActive:     true,
	}

	if err := s.itemRepo.Create(ctx, item); err != nil {
		return nil, fmt.Errorf("创建物料失败: %w", err)
	}

	return s.toItemResponse(item), nil
}

// GetItem 获取物料
func (s *ItemServiceImpl) GetItem(ctx context.Context, id uint) (*dto.ItemResponse, error) {
	item, err := s.itemRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("获取物料失败: %w", err)
	}
	if item == nil {
		return nil, errors.New("物料不存在")
	}

	return s.toItemResponse(item), nil
}

// UpdateItem 更新物料
func (s *ItemServiceImpl) UpdateItem(ctx context.Context, id uint, req *dto.ItemUpdateRequest) (*dto.ItemResponse, error) {
	item, err := s.itemRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("获取物料失败: %w", err)
	}
	if item == nil {
		return nil, errors.New("物料不存在")
	}

	// 更新物料信息
	if req.Name != "" {
		item.Name = req.Name
	}
	if req.Description != "" {
		item.Description = req.Description
	}
	// TODO: 需要根据CategoryID和UnitID查询对应的名称
	// if req.CategoryID != nil {
	//     item.Category = getCategoryName(*req.CategoryID)
	// }
	// if req.UnitID != nil {
	//     item.Unit = getUnitName(*req.UnitID)
	// }
	if req.UnitCost != nil {
		item.Cost = *req.UnitCost
	}
	if req.SalePrice != nil {
		item.Price = *req.SalePrice
	}
	if req.IsActive != nil {
		item.IsActive = *req.IsActive
	}
	if req.IsActive != nil {
		item.IsActive = *req.IsActive
	}
	item.UpdatedAt = time.Now()

	if err := s.itemRepo.Update(ctx, item); err != nil {
		return nil, fmt.Errorf("更新物料失败: %w", err)
	}

	return s.toItemResponse(item), nil
}

// DeleteItem 删除物料
func (s *ItemServiceImpl) DeleteItem(ctx context.Context, id uint) error {
	item, err := s.itemRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("获取物料失败: %w", err)
	}
	if item == nil {
		return errors.New("物料不存在")
	}

	if err := s.itemRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("删除物料失败: %w", err)
	}

	return nil
}

// GetItems 获取物料列表
func (s *ItemServiceImpl) GetItems(ctx context.Context, req *dto.PaginationRequest) (*dto.BaseResponse, error) {
	offset := req.GetOffset()
	limit := req.GetLimit()

	items, _, err := s.itemRepo.List(ctx, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("获取物料列表失败: %w", err)
	}

	// 转换为响应格式
	itemResponses := make([]dto.ItemResponse, len(items))
	for i, item := range items {
		itemResponses[i] = *s.toItemResponse(item)
	}

	return &dto.BaseResponse{
		Success: true,
		Message: "获取物料列表成功",
		Data:    itemResponses,
	}, nil
}

// SearchItems 搜索物料
func (s *ItemServiceImpl) SearchItems(ctx context.Context, req *dto.ItemSearchRequest) (*dto.BaseResponse, error) {
	offset := req.GetOffset()
	limit := req.GetLimit()

	items, _, err := s.itemRepo.Search(ctx, req.Keyword, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("搜索物料失败: %w", err)
	}

	// 转换为响应格式
	itemResponses := make([]dto.ItemResponse, len(items))
	for i, item := range items {
		itemResponses[i] = *s.toItemResponse(item)
	}

	return &dto.BaseResponse{
		Success: true,
		Message: "搜索物料成功",
		Data:    itemResponses,
	}, nil
}

// toItemResponse 转换为物料响应格式
func (s *ItemServiceImpl) toItemResponse(item *models.Item) *dto.ItemResponse {
	return &dto.ItemResponse{
		ID:          item.ID,
		Name:        item.Name,
		Code:        item.Code,
		Description: item.Description,
		// TODO: 需要根据字符串字段映射到对应的ID和对象
		// MinStock:    float64(item.ReorderLevel),
		// MaxStock:    0, // 当前模型中没有MaxStock字段
		UnitCost:    item.Cost,
		SalePrice:   item.Price,
		// Barcode:     "", // 当前模型中没有Barcode字段
		IsActive:    item.IsActive,
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	}
}

// StockService 库存服务接口
type StockService interface {
	CreateStock(ctx context.Context, req *dto.MovementCreateRequest) (*dto.StockResponse, error)
	GetStock(ctx context.Context, id uint) (*dto.StockResponse, error)
	UpdateStock(ctx context.Context, id uint, quantity float64) (*dto.StockResponse, error)
	DeleteStock(ctx context.Context, id uint) error
	GetStocks(ctx context.Context, req *dto.PaginationRequest) (*dto.BaseResponse, error)
	GetByItemID(ctx context.Context, itemID uint) (*dto.BaseResponse, error)
	AdjustStock(ctx context.Context, req *dto.StockAdjustmentCreateRequest) (*dto.StockAdjustmentResponse, error)
}

// StockServiceImpl 库存服务实现
type StockServiceImpl struct {
	stockRepo repositories.StockRepository
}

// NewStockService 创建库存服务实例
func NewStockService(stockRepo repositories.StockRepository) StockService {
	return &StockServiceImpl{
		stockRepo: stockRepo,
	}
}

// CreateStock 创建库存
func (s *StockServiceImpl) CreateStock(ctx context.Context, req *dto.MovementCreateRequest) (*dto.StockResponse, error) {
	// 这里应该创建库存移动记录，而不是直接创建库存
	// 实际实现需要根据业务逻辑调整
	return nil, errors.New("方法未实现")
}

// GetStock 获取库存
func (s *StockServiceImpl) GetStock(ctx context.Context, id uint) (*dto.StockResponse, error) {
	stock, err := s.stockRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("获取库存失败: %w", err)
	}
	if stock == nil {
		return nil, errors.New("库存不存在")
	}

	return s.toStockResponse(stock), nil
}

// UpdateStock 更新库存
func (s *StockServiceImpl) UpdateStock(ctx context.Context, id uint, quantity float64) (*dto.StockResponse, error) {
	stock, err := s.stockRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("获取库存失败: %w", err)
	}
	if stock == nil {
		return nil, errors.New("库存不存在")
	}

	// 更新库存数量
	stock.Quantity = quantity
	// TODO: 需要实现预留数量逻辑

	if err := s.stockRepo.Update(ctx, stock); err != nil {
		return nil, fmt.Errorf("更新库存失败: %w", err)
	}

	return s.toStockResponse(stock), nil
}

// DeleteStock 删除库存
func (s *StockServiceImpl) DeleteStock(ctx context.Context, id uint) error {
	stock, err := s.stockRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("获取库存失败: %w", err)
	}
	if stock == nil {
		return errors.New("库存不存在")
	}

	if err := s.stockRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("删除库存失败: %w", err)
	}

	return nil
}

// GetStocks 获取库存列表
func (s *StockServiceImpl) GetStocks(ctx context.Context, req *dto.PaginationRequest) (*dto.BaseResponse, error) {
	offset := req.GetOffset()
	limit := req.GetLimit()

	stocks, _, err := s.stockRepo.List(ctx, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("获取库存列表失败: %w", err)
	}

	// 转换为响应格式
	stockResponses := make([]dto.StockResponse, len(stocks))
	for i, stock := range stocks {
		stockResponses[i] = *s.toStockResponse(stock)
	}

	return &dto.BaseResponse{
		Success: true,
		Message: "获取库存列表成功",
		Data:    stockResponses,
	}, nil
}

// GetByItemID 根据物料ID获取库存
func (s *StockServiceImpl) GetByItemID(ctx context.Context, itemID uint) (*dto.BaseResponse, error) {
	stocks, err := s.stockRepo.GetByItemID(ctx, itemID)
	if err != nil {
		return nil, fmt.Errorf("获取物料库存失败: %w", err)
	}

	var stockResponses []*dto.StockResponse
	for _, stock := range stocks {
		stockResponses = append(stockResponses, s.toStockResponse(stock))
	}

	return &dto.BaseResponse{
		Success: true,
		Message: "获取物料库存成功",
		Data:    stockResponses,
	}, nil
}

// AdjustStock 调整库存
func (s *StockServiceImpl) AdjustStock(ctx context.Context, req *dto.StockAdjustmentCreateRequest) (*dto.StockAdjustmentResponse, error) {
	// 这里应该实现库存调整逻辑
	// 实际实现需要根据业务逻辑调整
	return nil, errors.New("方法未实现")
}

// toStockResponse 转换为库存响应格式
func (s *StockServiceImpl) toStockResponse(stock *models.Stock) *dto.StockResponse {
	return &dto.StockResponse{
		ID:        stock.ID,
		Quantity:  stock.Quantity,
		// TODO: 需要实现预留数量和可用数量逻辑
		// ReservedQty:  0,
		// AvailableQty: stock.Quantity,
		UpdatedAt: stock.UpdatedAt,
	}
}

// WarehouseService 仓库服务接口
type WarehouseService interface {
	CreateWarehouse(ctx context.Context, req *dto.WarehouseCreateRequest) (*dto.WarehouseResponse, error)
	GetWarehouse(ctx context.Context, id uint) (*dto.WarehouseResponse, error)
	UpdateWarehouse(ctx context.Context, id uint, req *dto.WarehouseUpdateRequest) (*dto.WarehouseResponse, error)
	DeleteWarehouse(ctx context.Context, id uint) error
	GetWarehouses(ctx context.Context, req *dto.PaginationRequest) (*dto.BaseResponse, error)
}

// WarehouseServiceImpl 仓库服务实现
type WarehouseServiceImpl struct {
	warehouseRepo repositories.WarehouseRepository
}

// NewWarehouseService 创建仓库服务实例
func NewWarehouseService(warehouseRepo repositories.WarehouseRepository) WarehouseService {
	return &WarehouseServiceImpl{
		warehouseRepo: warehouseRepo,
	}
}

// CreateWarehouse 创建仓库
func (s *WarehouseServiceImpl) CreateWarehouse(ctx context.Context, req *dto.WarehouseCreateRequest) (*dto.WarehouseResponse, error) {
	// 创建仓库
	warehouse := &models.Warehouse{
		Code:        req.Code,
		Name:        req.Name,
		Address:     req.Address,
		Description: req.Description,
		IsActive:    true,
	}

	if err := s.warehouseRepo.Create(ctx, warehouse); err != nil {
		return nil, fmt.Errorf("创建仓库失败: %w", err)
	}

	return s.toWarehouseResponse(warehouse), nil
}

// GetWarehouse 获取仓库
func (s *WarehouseServiceImpl) GetWarehouse(ctx context.Context, id uint) (*dto.WarehouseResponse, error) {
	warehouse, err := s.warehouseRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("获取仓库失败: %w", err)
	}
	if warehouse == nil {
		return nil, errors.New("仓库不存在")
	}

	return s.toWarehouseResponse(warehouse), nil
}

// UpdateWarehouse 更新仓库
func (s *WarehouseServiceImpl) UpdateWarehouse(ctx context.Context, id uint, req *dto.WarehouseUpdateRequest) (*dto.WarehouseResponse, error) {
	warehouse, err := s.warehouseRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("获取仓库失败: %w", err)
	}
	if warehouse == nil {
		return nil, errors.New("仓库不存在")
	}

	// 更新字段
	if req.Name != "" {
		warehouse.Name = req.Name
	}
	if req.Address != "" {
		warehouse.Address = req.Address
	}
	if req.Description != "" {
		warehouse.Description = req.Description
	}
	if req.IsActive != nil {
		warehouse.IsActive = *req.IsActive
	}

	if err := s.warehouseRepo.Update(ctx, warehouse); err != nil {
		return nil, fmt.Errorf("更新仓库失败: %w", err)
	}

	return s.toWarehouseResponse(warehouse), nil
}

// DeleteWarehouse 删除仓库
func (s *WarehouseServiceImpl) DeleteWarehouse(ctx context.Context, id uint) error {
	warehouse, err := s.warehouseRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("获取仓库失败: %w", err)
	}
	if warehouse == nil {
		return errors.New("仓库不存在")
	}

	if err := s.warehouseRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("删除仓库失败: %w", err)
	}

	return nil
}

// GetWarehouses 获取仓库列表
func (s *WarehouseServiceImpl) GetWarehouses(ctx context.Context, req *dto.PaginationRequest) (*dto.BaseResponse, error) {
	warehouses, _, err := s.warehouseRepo.List(ctx, req.GetOffset(), req.GetLimit())
	if err != nil {
		return nil, fmt.Errorf("获取仓库列表失败: %w", err)
	}

	var warehouseResponses []*dto.WarehouseResponse
	for _, warehouse := range warehouses {
		warehouseResponses = append(warehouseResponses, s.toWarehouseResponse(warehouse))
	}

	return &dto.BaseResponse{
		Success: true,
		Message: "获取仓库列表成功",
		Data:    warehouseResponses,
	}, nil
}

// toWarehouseResponse 转换为仓库响应
func (s *WarehouseServiceImpl) toWarehouseResponse(warehouse *models.Warehouse) *dto.WarehouseResponse {
	return &dto.WarehouseResponse{
		ID:          warehouse.ID,
		Name:        warehouse.Name,
		Code:        warehouse.Code,
		Address:     warehouse.Address,
		Description: warehouse.Description,
		IsActive:    warehouse.IsActive,
		CreatedAt:   warehouse.CreatedAt,
		UpdatedAt:   warehouse.UpdatedAt,
	}
}

// MovementService 库存移动服务接口
type MovementService interface {
	CreateMovement(ctx context.Context, req dto.MovementCreateRequest) (*dto.MovementResponse, error)
	GetMovementByID(ctx context.Context, id uint) (*dto.MovementResponse, error)
	GetMovementsByItemID(ctx context.Context, itemID uint, page, limit int) (*dto.PaginatedResponse[dto.MovementResponse], error)
	GetMovementsByWarehouseID(ctx context.Context, warehouseID uint, page, limit int) (*dto.PaginatedResponse[dto.MovementResponse], error)
	GetMovementsByType(ctx context.Context, movementType string, page, limit int) (*dto.PaginatedResponse[dto.MovementResponse], error)
}

// MovementServiceImpl 库存移动服务实现
type MovementServiceImpl struct {
	movementRepo repositories.MovementRepository
	stockRepo    repositories.StockRepository
	itemRepo     repositories.ItemRepository
	warehouseRepo repositories.WarehouseRepository
}

// NewMovementService 创建库存移动服务
func NewMovementService(
	movementRepo repositories.MovementRepository,
	stockRepo repositories.StockRepository,
	itemRepo repositories.ItemRepository,
	warehouseRepo repositories.WarehouseRepository,
) MovementService {
	return &MovementServiceImpl{
		movementRepo:  movementRepo,
		stockRepo:     stockRepo,
		itemRepo:      itemRepo,
		warehouseRepo: warehouseRepo,
	}
}

// CreateMovement 创建库存移动
func (s *MovementServiceImpl) CreateMovement(ctx context.Context, req dto.MovementCreateRequest) (*dto.MovementResponse, error) {
	// 验证物料是否存在
	item, err := s.itemRepo.GetByID(ctx, req.ItemID)
	if err != nil {
		return nil, fmt.Errorf("物料不存在: %w", err)
	}

	// 验证仓库是否存在
	warehouse, err := s.warehouseRepo.GetByID(ctx, req.WarehouseID)
	if err != nil {
		return nil, fmt.Errorf("仓库不存在: %w", err)
	}

	// 创建库存移动记录
	movement := &models.Movement{
		ItemID:       &req.ItemID,
		WarehouseID:  &req.WarehouseID,
		Quantity:     &req.Quantity,
		MovementType: req.Type,
		Reference:    req.Reference,
		Notes:        req.Notes,
	}

	err = s.movementRepo.Create(ctx, movement)
	if err != nil {
		return nil, fmt.Errorf("创建库存移动失败: %w", err)
	}

	// 更新库存数量
	if err := s.updateStockQuantity(ctx, req.ItemID, req.WarehouseID, req.Quantity, req.Type); err != nil {
		return nil, fmt.Errorf("更新库存失败: %w", err)
	}

	return s.convertToMovementResponse(movement, item, warehouse), nil
}

// GetMovementByID 根据ID获取库存移动
func (s *MovementServiceImpl) GetMovementByID(ctx context.Context, id uint) (*dto.MovementResponse, error) {
	movement, err := s.movementRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("获取库存移动失败: %w", err)
	}

	// 获取关联的物料和仓库信息
	var item *models.Item
	var warehouse *models.Warehouse

	if movement.ItemID != nil {
		item, _ = s.itemRepo.GetByID(ctx, *movement.ItemID)
	}

	if movement.WarehouseID != nil {
		warehouse, _ = s.warehouseRepo.GetByID(ctx, *movement.WarehouseID)
	}

	return s.convertToMovementResponse(movement, item, warehouse), nil
}

// GetMovementsByItemID 根据物料ID获取库存移动列表
func (s *MovementServiceImpl) GetMovementsByItemID(ctx context.Context, itemID uint, page, limit int) (*dto.PaginatedResponse[dto.MovementResponse], error) {
	offset := (page - 1) * limit
	movements, total, err := s.movementRepo.GetByItemID(ctx, itemID, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("获取物料库存移动失败: %w", err)
	}

	responses := make([]dto.MovementResponse, 0, len(movements))
	for _, movement := range movements {
		var item *models.Item
		var warehouse *models.Warehouse

		if movement.ItemID != nil {
			item, _ = s.itemRepo.GetByID(ctx, *movement.ItemID)
		}

		if movement.WarehouseID != nil {
			warehouse, _ = s.warehouseRepo.GetByID(ctx, *movement.WarehouseID)
		}

		responses = append(responses, *s.convertToMovementResponse(movement, item, warehouse))
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	return &dto.PaginatedResponse[dto.MovementResponse]{
		Data:       responses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

// GetMovementsByWarehouseID 根据仓库ID获取库存移动列表
func (s *MovementServiceImpl) GetMovementsByWarehouseID(ctx context.Context, warehouseID uint, page, limit int) (*dto.PaginatedResponse[dto.MovementResponse], error) {
	offset := (page - 1) * limit
	movements, total, err := s.movementRepo.GetByWarehouseID(ctx, warehouseID, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("获取仓库库存移动失败: %w", err)
	}

	responses := make([]dto.MovementResponse, 0, len(movements))
	for _, movement := range movements {
		var item *models.Item
		var warehouse *models.Warehouse

		if movement.ItemID != nil {
			item, _ = s.itemRepo.GetByID(ctx, *movement.ItemID)
		}

		if movement.WarehouseID != nil {
			warehouse, _ = s.warehouseRepo.GetByID(ctx, *movement.WarehouseID)
		}

		responses = append(responses, *s.convertToMovementResponse(movement, item, warehouse))
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	return &dto.PaginatedResponse[dto.MovementResponse]{
		Data:       responses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

// GetMovementsByType 根据移动类型获取库存移动列表
func (s *MovementServiceImpl) GetMovementsByType(ctx context.Context, movementType string, page, limit int) (*dto.PaginatedResponse[dto.MovementResponse], error) {
	offset := (page - 1) * limit
	movements, total, err := s.movementRepo.GetByType(ctx, movementType, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("获取类型库存移动失败: %w", err)
	}

	responses := make([]dto.MovementResponse, 0, len(movements))
	for _, movement := range movements {
		var item *models.Item
		var warehouse *models.Warehouse

		if movement.ItemID != nil {
			item, _ = s.itemRepo.GetByID(ctx, *movement.ItemID)
		}

		if movement.WarehouseID != nil {
			warehouse, _ = s.warehouseRepo.GetByID(ctx, *movement.WarehouseID)
		}

		responses = append(responses, *s.convertToMovementResponse(movement, item, warehouse))
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	return &dto.PaginatedResponse[dto.MovementResponse]{
		Data:       responses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

// updateStockQuantity 更新库存数量
func (s *MovementServiceImpl) updateStockQuantity(ctx context.Context, itemID, warehouseID uint, quantity float64, movementType string) error {
	// 查找现有库存记录
	stocks, err := s.stockRepo.GetByItemID(ctx, itemID)
	if err != nil {
		return err
	}

	var targetStock *models.Stock
	for _, stock := range stocks {
		if stock.WarehouseID == warehouseID {
			targetStock = stock
			break
		}
	}

	// 如果没有找到库存记录，创建新的
	if targetStock == nil {
		targetStock = &models.Stock{
			ItemID:      itemID,
			WarehouseID: warehouseID,
			Quantity:    0,
		}
		err = s.stockRepo.Create(ctx, targetStock)
		if err != nil {
			return fmt.Errorf("创建库存记录失败: %w", err)
		}
	}

	// 根据移动类型调整库存数量
	switch movementType {
	case "in":
		targetStock.Quantity += quantity
	case "out":
		if targetStock.Quantity < quantity {
			return fmt.Errorf("库存不足，当前库存: %.2f，需要: %.2f", targetStock.Quantity, quantity)
		}
		targetStock.Quantity -= quantity
	case "adjustment":
		targetStock.Quantity = quantity
	}

	return s.stockRepo.Update(ctx, targetStock)
}

// convertToMovementResponse 转换为移动响应
func (s *MovementServiceImpl) convertToMovementResponse(movement *models.Movement, item *models.Item, warehouse *models.Warehouse) *dto.MovementResponse {
	response := &dto.MovementResponse{
		ID:        movement.ID,
		Type:      movement.MovementType,
		Reference: movement.Reference,
		Notes:     movement.Notes,
		CreatedAt: movement.CreatedAt,
	}

	if movement.Quantity != nil {
		response.Quantity = *movement.Quantity
	}

	if item != nil {
		response.Item = dto.ItemResponse{
			ID:          item.ID,
			Code:        item.Code,
			Name:        item.Name,
			Description: item.Description,
			Type:        "",  // 模型中没有Type字段，使用空字符串
			MinStock:    0,   // 模型中没有MinStock字段
			MaxStock:    0,   // 模型中没有MaxStock字段
			UnitCost:    item.Cost,
			SalePrice:   item.Price,
			IsActive:    item.IsActive,
			Category: dto.CategoryResponse{
				Name: item.Category, // 模型中Category是字符串
			},
			Unit: dto.UnitResponse{
				Symbol: item.Unit, // 模型中Unit是字符串
			},
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
		}
	}

	if warehouse != nil {
		response.Warehouse = dto.WarehouseResponse{
			ID:          warehouse.ID,
			Name:        warehouse.Name,
			Code:        warehouse.Code,
			Address:     warehouse.Address,
			Description: warehouse.Description,
			IsActive:    warehouse.IsActive,
			CreatedAt:   warehouse.CreatedAt,
			UpdatedAt:   warehouse.UpdatedAt,
		}
	}

	return response
}