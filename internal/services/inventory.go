package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/galaxyerp/galaxyErp/internal/common"
	"github.com/galaxyerp/galaxyErp/internal/dto"
	"github.com/galaxyerp/galaxyErp/internal/models"
	"github.com/galaxyerp/galaxyErp/internal/repositories"
)

// ItemService 物料服务接口
type ItemService interface {
	CRUDService[models.Item, dto.ItemCreateRequest, dto.ItemUpdateRequest, dto.ItemResponse]
	CreateItem(ctx context.Context, req *dto.ItemCreateRequest) (*dto.ItemResponse, error)
	GetItem(ctx context.Context, id uint) (*dto.ItemResponse, error)
	UpdateItem(ctx context.Context, id uint, req *dto.ItemUpdateRequest) (*dto.ItemResponse, error)
	DeleteItem(ctx context.Context, id uint) error
	GetItems(ctx context.Context, req *dto.PaginationRequest) (*dto.PaginatedResponse[dto.ItemResponse], error)
	SearchItems(ctx context.Context, req *dto.ItemSearchRequest) (*dto.PaginatedResponse[dto.ItemResponse], error)
}

// ItemServiceImpl 物料服务实现
type ItemServiceImpl struct {
	*BaseService
	itemRepo repositories.ItemRepository
}

// NewItemService 创建物料服务实例
func NewItemService(itemRepo repositories.ItemRepository) ItemService {
	config := &BaseServiceConfig{
		EnableAudit:      true,
		EnableValidation: true,
		EnableCache:      true,
		CacheExpiry:      time.Hour,
		EnableMetrics:    true,
	}

	return &ItemServiceImpl{
		BaseService: NewBaseService(config),
		itemRepo:    itemRepo,
	}
}

// Create 实现 CRUDService 接口的 Create 方法
func (s *ItemServiceImpl) Create(ctx context.Context, req *dto.ItemCreateRequest) (*dto.CreateResponse, error) {
	// 验证请求
	if err := s.ValidateRequest(ctx, "item", req); err != nil {
		return nil, err
	}

	// 创建物料
	item := &models.Item{
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
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

	// 清除相关缓存
	s.DeleteFromCache(ctx, "items_list")

	return s.CreateCreateResponse(item.ID, s.toItemResponse(item), "物料创建成功"), nil
}

// CreateItem 创建物料
func (s *ItemServiceImpl) CreateItem(ctx context.Context, req *dto.ItemCreateRequest) (*dto.ItemResponse, error) {
	createResp, err := s.Create(ctx, req)
	if err != nil {
		return nil, err
	}
	return createResp.Data.(*dto.ItemResponse), nil
}

// GetByID 实现 CRUDService 接口的 GetByID 方法
func (s *ItemServiceImpl) GetByID(ctx context.Context, id uint) (*dto.ItemResponse, error) {
	// 尝试从缓存获取
	cacheKey := fmt.Sprintf("item_%d", id)
	if cached, found := s.GetFromCache(ctx, cacheKey); found {
		return cached.(*dto.ItemResponse), nil
	}

	item, err := s.itemRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("获取物料失败: %w", err)
	}
	if item == nil {
		return nil, errors.New("物料不存在")
	}

	response := s.toItemResponse(item)
	
	// 缓存结果
	s.SetToCache(ctx, cacheKey, response)

	return response, nil
}

// GetItem 获取物料
func (s *ItemServiceImpl) GetItem(ctx context.Context, id uint) (*dto.ItemResponse, error) {
	return s.GetByID(ctx, id)
}

// Update 实现 CRUDService 接口的 Update 方法
func (s *ItemServiceImpl) Update(ctx context.Context, id uint, req *dto.ItemUpdateRequest) (*dto.UpdateResponse, error) {
	// 验证请求
	if err := s.ValidateRequest(ctx, "item", req); err != nil {
		return nil, err
	}

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
	item.UpdatedAt = time.Now()

	if err := s.itemRepo.Update(ctx, item); err != nil {
		return nil, fmt.Errorf("更新物料失败: %w", err)
	}

	// 清除相关缓存
	cacheKey := fmt.Sprintf("item_%d", id)
	s.DeleteFromCache(ctx, cacheKey)
	s.DeleteFromCache(ctx, "items_list")

	return s.CreateUpdateResponse(s.toItemResponse(item), "物料更新成功"), nil
}

// UpdateItem 更新物料
func (s *ItemServiceImpl) UpdateItem(ctx context.Context, id uint, req *dto.ItemUpdateRequest) (*dto.ItemResponse, error) {
	updateResp, err := s.Update(ctx, id, req)
	if err != nil {
		return nil, err
	}
	return updateResp.Data.(*dto.ItemResponse), nil
}

// Delete 实现 CRUDService 接口的 Delete 方法
func (s *ItemServiceImpl) Delete(ctx context.Context, id uint) (*dto.DeleteResponse, error) {
	item, err := s.itemRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("获取物料失败: %w", err)
	}
	if item == nil {
		return nil, errors.New("物料不存在")
	}

	if err := s.itemRepo.Delete(ctx, id); err != nil {
		return nil, fmt.Errorf("删除物料失败: %w", err)
	}

	// 清除相关缓存
	cacheKey := fmt.Sprintf("item_%d", id)
	s.DeleteFromCache(ctx, cacheKey)
	s.DeleteFromCache(ctx, "items_list")

	return s.CreateDeleteResponse("物料删除成功"), nil
}

// DeleteItem 删除物料
func (s *ItemServiceImpl) DeleteItem(ctx context.Context, id uint) error {
	_, err := s.Delete(ctx, id)
	return err
}

// List 实现 CRUDService 接口的 List 方法
func (s *ItemServiceImpl) List(ctx context.Context, req *dto.PaginationRequest) (*dto.PaginatedResponse[dto.ItemResponse], error) {
	offset := req.GetOffset()
	limit := req.GetLimit()

	items, total, err := s.itemRepo.ListItems(ctx, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("获取物料列表失败: %w", err)
	}

	// 转换为响应格式
	itemResponses := make([]dto.ItemResponse, len(items))
	for i, item := range items {
		itemResponses[i] = *s.toItemResponse(item)
	}

	// 计算总页数
	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return &dto.PaginatedResponse[dto.ItemResponse]{
		Data:       itemResponses,
		Total:      total,
		Page:       req.Page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

// Search 实现 CRUDService 接口的 Search 方法
func (s *ItemServiceImpl) Search(ctx context.Context, req *dto.SearchRequest) (*dto.PaginatedResponse[dto.ItemResponse], error) {
	offset := req.GetOffset()
	limit := req.GetLimit()

	items, total, err := s.itemRepo.Search(ctx, req.Keyword, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("搜索物料失败: %w", err)
	}

	// 转换为响应格式
	itemResponses := make([]dto.ItemResponse, len(items))
	for i, item := range items {
		itemResponses[i] = *s.toItemResponse(item)
	}

	// 计算总页数
	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return &dto.PaginatedResponse[dto.ItemResponse]{
		Data:       itemResponses,
		Total:      total,
		Page:       req.Page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

// GetItems 获取物料列表
func (s *ItemServiceImpl) GetItems(ctx context.Context, req *dto.PaginationRequest) (*dto.PaginatedResponse[dto.ItemResponse], error) {
	return s.List(ctx, req)
}

// SearchItems 搜索物料
func (s *ItemServiceImpl) SearchItems(ctx context.Context, req *dto.ItemSearchRequest) (*dto.PaginatedResponse[dto.ItemResponse], error) {
	// 转换为通用搜索请求
	searchReq := &dto.SearchRequest{
		PaginationRequest: req.PaginationRequest,
		Keyword:           req.Keyword,
	}
	return s.Search(ctx, searchReq)
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
		UnitCost:  item.Cost,
		SalePrice: item.Price,
		// Barcode:     "", // 当前模型中没有Barcode字段
		IsActive:  item.IsActive,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}
}

// StockService 库存服务接口
type StockService interface {
	CRUDService[models.Stock, dto.StockCreateRequest, dto.StockUpdateRequest, dto.StockResponse]
	
	// 原有的特定方法
	CreateStock(ctx context.Context, req *dto.MovementCreateRequest) (*dto.StockResponse, error)
	UpdateStock(ctx context.Context, id uint, quantity float64) (*dto.StockResponse, error)
	GetByItemID(ctx context.Context, itemID uint) ([]*dto.StockResponse, error)
	AdjustStock(ctx context.Context, req *dto.StockAdjustmentCreateRequest) (*dto.StockAdjustmentResponse, error)
}

// StockServiceImpl 库存服务实现
type StockServiceImpl struct {
	*BaseService
	stockRepo repositories.StockRepository
}

// NewStockService 创建库存服务
func NewStockService(stockRepo repositories.StockRepository) StockService {
	config := &BaseServiceConfig{
		EnableAudit:      true,
		EnableValidation: true,
		EnableCache:      true,
		EnableMetrics:    true,
		CacheExpiry:      time.Hour,
	}
	
	return &StockServiceImpl{
		BaseService: NewBaseService(config),
		stockRepo:   stockRepo,
	}
}

// Create 实现 CRUDService 接口的 Create 方法
func (s *StockServiceImpl) Create(ctx context.Context, req *dto.StockCreateRequest) (*dto.CreateResponse, error) {
	// 验证请求
	if err := s.ValidateRequest(ctx, "stock", req); err != nil {
		return nil, err
	}

	// 创建库存记录
	stock := &models.Stock{
		ItemID:      req.ItemID,
		WarehouseID: req.WarehouseID,
		Quantity:    req.Quantity,
	}

	if err := s.stockRepo.Create(ctx, stock); err != nil {
		return nil, err
	}

	// 清除相关缓存
	s.DeleteFromCache(ctx, fmt.Sprintf("stock:list"))
	s.DeleteFromCache(ctx, fmt.Sprintf("stock:item:%d", req.ItemID))

	return s.CreateCreateResponse(stock.ID, s.toStockResponse(stock), "库存创建成功"), nil
}

// GetByID 实现 CRUDService 接口的 GetByID 方法
func (s *StockServiceImpl) GetByID(ctx context.Context, id uint) (*dto.StockResponse, error) {
	// 尝试从缓存获取
	if cached, found := s.GetFromCache(ctx, fmt.Sprintf("stock:%d", id)); found {
		if stock, ok := cached.(*dto.StockResponse); ok {
			return stock, nil
		}
	}

	stock, err := s.stockRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	response := s.toStockResponse(stock)
	
	// 缓存结果
	s.SetToCache(ctx, fmt.Sprintf("stock:%d", id), response)
	
	return response, nil
}

// Update 实现 CRUDService 接口的 Update 方法
func (s *StockServiceImpl) Update(ctx context.Context, id uint, req *dto.StockUpdateRequest) (*dto.UpdateResponse, error) {
	// 验证请求
	if err := s.ValidateRequest(ctx, "stock", req); err != nil {
		return nil, err
	}

	// 获取现有库存
	stock, err := s.stockRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 更新字段
	if req.Quantity != nil {
		stock.Quantity = *req.Quantity
	}

	if err := s.stockRepo.Update(ctx, stock); err != nil {
		return nil, err
	}

	// 清除缓存
	s.DeleteFromCache(ctx, fmt.Sprintf("stock:%d", id))
	s.DeleteFromCache(ctx, fmt.Sprintf("stock:list"))
	s.DeleteFromCache(ctx, fmt.Sprintf("stock:item:%d", stock.ItemID))

	return s.CreateUpdateResponse(s.toStockResponse(stock), "库存更新成功"), nil
}

// Delete 实现 CRUDService 接口的 Delete 方法
func (s *StockServiceImpl) Delete(ctx context.Context, id uint) (*dto.DeleteResponse, error) {
	// 获取库存信息用于清除缓存
	stock, err := s.stockRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := s.stockRepo.Delete(ctx, id); err != nil {
		return nil, err
	}

	// 清除缓存
	s.DeleteFromCache(ctx, fmt.Sprintf("stock:%d", id))
	s.DeleteFromCache(ctx, fmt.Sprintf("stock:list"))
	s.DeleteFromCache(ctx, fmt.Sprintf("stock:item:%d", stock.ItemID))

	return s.CreateDeleteResponse("库存删除成功"), nil
}

// List 实现 CRUDService 接口的 List 方法
func (s *StockServiceImpl) List(ctx context.Context, req *dto.PaginationRequest) (*dto.PaginatedResponse[dto.StockResponse], error) {
	// 尝试从缓存获取
	cacheKey := fmt.Sprintf("stock:list:%d:%d", req.Page, req.PageSize)
	if cached, found := s.GetFromCache(ctx, cacheKey); found {
		if response, ok := cached.(*dto.PaginatedResponse[dto.StockResponse]); ok {
			return response, nil
		}
	}

	options := &common.QueryOptions{
		Pagination: req,
	}
	stocks, total, err := s.stockRepo.List(ctx, options)
	if err != nil {
		return nil, err
	}

	stockResponses := make([]dto.StockResponse, len(stocks))
	for i, stock := range stocks {
		stockResponses[i] = *s.toStockResponse(stock)
	}

	response := &dto.PaginatedResponse[dto.StockResponse]{
		Data:       stockResponses,
		Total:      total,
		Page:       req.Page,
		Limit:      req.PageSize,
		TotalPages: int((total + int64(req.PageSize) - 1) / int64(req.PageSize)),
	}

	// 缓存结果
	s.SetToCache(ctx, cacheKey, response)

	return response, nil
}

// Search 实现 CRUDService 接口的 Search 方法
func (s *StockServiceImpl) Search(ctx context.Context, req *dto.SearchRequest) (*dto.PaginatedResponse[dto.StockResponse], error) {
	// 这里可以实现更复杂的搜索逻辑
	// 暂时使用 List 方法的逻辑
	paginationReq := &dto.PaginationRequest{
		Page:     req.Page,
		PageSize: req.PageSize,
	}
	
	return s.List(ctx, paginationReq)
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
func (s *StockServiceImpl) GetStocks(ctx context.Context, req *dto.PaginationRequest) (*dto.PaginatedResponse[dto.StockResponse], error) {
	stocks, total, err := s.stockRepo.List(ctx, &common.QueryOptions{
		Pagination: req,
	})
	if err != nil {
		return nil, fmt.Errorf("获取库存列表失败: %w", err)
	}

	// 转换为响应格式
	stockResponses := make([]dto.StockResponse, len(stocks))
	for i, stock := range stocks {
		stockResponses[i] = *s.toStockResponse(stock)
	}

	return &dto.PaginatedResponse[dto.StockResponse]{
		Data:       stockResponses,
		Total:      total,
		Page:       req.Page,
		Limit:      req.PageSize,
		TotalPages: int((total + int64(req.PageSize) - 1) / int64(req.PageSize)),
	}, nil
}

// GetByItemID 根据物料ID获取库存
func (s *StockServiceImpl) GetByItemID(ctx context.Context, itemID uint) ([]*dto.StockResponse, error) {
	stocks, err := s.stockRepo.GetByItemID(ctx, itemID)
	if err != nil {
		return nil, fmt.Errorf("获取物料库存失败: %w", err)
	}

	var stockResponses []*dto.StockResponse
	for _, stock := range stocks {
		stockResponses = append(stockResponses, s.toStockResponse(stock))
	}

	return stockResponses, nil
}

// AdjustStock 调整库存
func (s *StockServiceImpl) AdjustStock(ctx context.Context, req *dto.StockAdjustmentCreateRequest) (*dto.StockAdjustmentResponse, error) {
	// 这里应该实现库存调整逻辑
	// 实际实现需要根据业务逻辑调整
	return nil, errors.New("方法未实现")
}

// toStockResponse 转换为库存响应格式
func (s *StockServiceImpl) toStockResponse(stock *models.Stock) *dto.StockResponse {
	response := &dto.StockResponse{
		ID:           stock.ID,
		Quantity:     stock.Quantity,
		ReservedQty:  0,              // TODO: 需要实现预留数量逻辑
		AvailableQty: stock.Quantity, // 暂时设为总数量
		UpdatedAt:    stock.UpdatedAt,
	}

	// 填充Item信息
	if stock.Item.ID != 0 {
		response.Item = dto.ItemResponse{
			ID:          stock.Item.ID,
			Code:        stock.Item.Code,
			Name:        stock.Item.Name,
			Description: stock.Item.Description,
			Type:        stock.Item.Category, // 暂时使用Category字段作为Type
			MinStock:    0,                   // TODO: 从Item模型获取
			MaxStock:    0,                   // TODO: 从Item模型获取
			UnitCost:    stock.Item.Cost,
			SalePrice:   stock.Item.Price,
			IsActive:    stock.Item.IsActive,
			CreatedAt:   stock.Item.CreatedAt,
			UpdatedAt:   stock.Item.UpdatedAt,
		}
	}

	// 填充Warehouse信息
	if stock.Warehouse.ID != 0 {
		response.Warehouse = dto.WarehouseResponse{
			ID:          stock.Warehouse.ID,
			Name:        stock.Warehouse.Name,
			Code:        stock.Warehouse.Code,
			Address:     stock.Warehouse.Address,
			Description: stock.Warehouse.Description,
			IsActive:    stock.Warehouse.IsActive,
			CreatedAt:   stock.Warehouse.CreatedAt,
			UpdatedAt:   stock.Warehouse.UpdatedAt,
		}
	}

	// Location字段暂时留空，因为当前数据库结构中没有location_id
	response.Location = dto.LocationResponse{}

	return response
}

// WarehouseService 仓库服务接口
type WarehouseService interface {
	CRUDService[models.Warehouse, dto.WarehouseCreateRequest, dto.WarehouseUpdateRequest, dto.WarehouseResponse]
	CreateWarehouse(ctx context.Context, req *dto.WarehouseCreateRequest) (*dto.WarehouseResponse, error)
	GetWarehouse(ctx context.Context, id uint) (*dto.WarehouseResponse, error)
	UpdateWarehouse(ctx context.Context, id uint, req *dto.WarehouseUpdateRequest) (*dto.WarehouseResponse, error)
	DeleteWarehouse(ctx context.Context, id uint) error
	GetWarehouses(ctx context.Context, req *dto.PaginationRequest) (*dto.PaginatedResponse[dto.WarehouseResponse], error)
}

// WarehouseServiceImpl 仓库服务实现
type WarehouseServiceImpl struct {
	*BaseService
	warehouseRepo repositories.WarehouseRepository
}

// NewWarehouseService 创建仓库服务实例
func NewWarehouseService(warehouseRepo repositories.WarehouseRepository) WarehouseService {
	config := &BaseServiceConfig{
		EnableAudit:      true,
		EnableValidation: true,
		EnableCache:      true,
		CacheExpiry:      time.Hour,
		EnableMetrics:    true,
	}

	return &WarehouseServiceImpl{
		BaseService:   NewBaseService(config),
		warehouseRepo: warehouseRepo,
	}
}

// Create 实现 CRUDService 接口的 Create 方法
func (s *WarehouseServiceImpl) Create(ctx context.Context, req *dto.WarehouseCreateRequest) (*dto.CreateResponse, error) {
	// 验证请求
	if err := s.ValidateRequest(ctx, "warehouse", req); err != nil {
		return nil, err
	}

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

	// 清除相关缓存
	s.DeleteFromCache(ctx, "warehouses_list")

	return s.CreateCreateResponse(warehouse.ID, s.toWarehouseResponse(warehouse), "仓库创建成功"), nil
}

// GetByID 实现 CRUDService 接口的 GetByID 方法
func (s *WarehouseServiceImpl) GetByID(ctx context.Context, id uint) (*dto.WarehouseResponse, error) {
	// 尝试从缓存获取
	cacheKey := fmt.Sprintf("warehouse_%d", id)
	if cached, found := s.GetFromCache(ctx, cacheKey); found {
		return cached.(*dto.WarehouseResponse), nil
	}

	warehouse, err := s.warehouseRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("获取仓库失败: %w", err)
	}
	if warehouse == nil {
		return nil, errors.New("仓库不存在")
	}

	response := s.toWarehouseResponse(warehouse)
	
	// 缓存结果
	s.SetToCache(ctx, cacheKey, response)

	return response, nil
}

// Update 实现 CRUDService 接口的 Update 方法
func (s *WarehouseServiceImpl) Update(ctx context.Context, id uint, req *dto.WarehouseUpdateRequest) (*dto.UpdateResponse, error) {
	// 验证请求
	if err := s.ValidateRequest(ctx, "warehouse", req); err != nil {
		return nil, err
	}

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
	warehouse.UpdatedAt = time.Now()

	if err := s.warehouseRepo.Update(ctx, warehouse); err != nil {
		return nil, fmt.Errorf("更新仓库失败: %w", err)
	}

	// 清除相关缓存
	cacheKey := fmt.Sprintf("warehouse_%d", id)
	s.DeleteFromCache(ctx, cacheKey)
	s.DeleteFromCache(ctx, "warehouses_list")

	return s.CreateUpdateResponse(s.toWarehouseResponse(warehouse), "仓库更新成功"), nil
}

// Delete 实现 CRUDService 接口的 Delete 方法
func (s *WarehouseServiceImpl) Delete(ctx context.Context, id uint) (*dto.DeleteResponse, error) {
	warehouse, err := s.warehouseRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("获取仓库失败: %w", err)
	}
	if warehouse == nil {
		return nil, errors.New("仓库不存在")
	}

	if err := s.warehouseRepo.Delete(ctx, id); err != nil {
		return nil, fmt.Errorf("删除仓库失败: %w", err)
	}

	// 清除相关缓存
	cacheKey := fmt.Sprintf("warehouse_%d", id)
	s.DeleteFromCache(ctx, cacheKey)
	s.DeleteFromCache(ctx, "warehouses_list")

	return s.CreateDeleteResponse("仓库删除成功"), nil
}

// List 实现 CRUDService 接口的 List 方法
func (s *WarehouseServiceImpl) List(ctx context.Context, req *dto.PaginationRequest) (*dto.PaginatedResponse[dto.WarehouseResponse], error) {
	warehouses, total, err := s.warehouseRepo.List(ctx, &common.QueryOptions{
		Pagination: req,
	})
	if err != nil {
		return nil, fmt.Errorf("获取仓库列表失败: %w", err)
	}

	// 转换为响应格式
	warehouseResponses := make([]dto.WarehouseResponse, len(warehouses))
	for i, warehouse := range warehouses {
		warehouseResponses[i] = *s.toWarehouseResponse(warehouse)
	}

	// 计算总页数
	limit := req.GetLimit()
	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return &dto.PaginatedResponse[dto.WarehouseResponse]{
		Data:       warehouseResponses,
		Total:      total,
		Page:       req.Page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

// Search 实现 CRUDService 接口的 Search 方法
func (s *WarehouseServiceImpl) Search(ctx context.Context, req *dto.SearchRequest) (*dto.PaginatedResponse[dto.WarehouseResponse], error) {
	// 目前使用 List 方法，后续可以在 repository 层添加搜索功能
	return s.List(ctx, &req.PaginationRequest)
}

// CreateWarehouse 创建仓库
func (s *WarehouseServiceImpl) CreateWarehouse(ctx context.Context, req *dto.WarehouseCreateRequest) (*dto.WarehouseResponse, error) {
	createResp, err := s.Create(ctx, req)
	if err != nil {
		return nil, err
	}
	return createResp.Data.(*dto.WarehouseResponse), nil
}

// GetWarehouse 获取仓库
func (s *WarehouseServiceImpl) GetWarehouse(ctx context.Context, id uint) (*dto.WarehouseResponse, error) {
	return s.GetByID(ctx, id)
}

// UpdateWarehouse 更新仓库
func (s *WarehouseServiceImpl) UpdateWarehouse(ctx context.Context, id uint, req *dto.WarehouseUpdateRequest) (*dto.WarehouseResponse, error) {
	updateResp, err := s.Update(ctx, id, req)
	if err != nil {
		return nil, err
	}
	return updateResp.Data.(*dto.WarehouseResponse), nil
}

// DeleteWarehouse 删除仓库
func (s *WarehouseServiceImpl) DeleteWarehouse(ctx context.Context, id uint) error {
	_, err := s.Delete(ctx, id)
	return err
}

// GetWarehouses 获取仓库列表
func (s *WarehouseServiceImpl) GetWarehouses(ctx context.Context, req *dto.PaginationRequest) (*dto.PaginatedResponse[dto.WarehouseResponse], error) {
	return s.List(ctx, req)
}

// toWarehouseResponse 转换为仓库响应
func (s *WarehouseServiceImpl) toWarehouseResponse(warehouse *models.Warehouse) *dto.WarehouseResponse {
	response := &dto.WarehouseResponse{
		ID:          warehouse.ID,
		Name:        warehouse.Name,
		Code:        warehouse.Code,
		Address:     warehouse.Address,
		Description: warehouse.Description,
		IsActive:    warehouse.IsActive,
		CreatedAt:   warehouse.CreatedAt,
		UpdatedAt:   warehouse.UpdatedAt,
	}

	// 填充Manager字段
	if warehouse.Manager != nil {
		response.Manager = &dto.UserResponse{
			ID:        warehouse.Manager.ID,
			Username:  warehouse.Manager.Username,
			Email:     warehouse.Manager.Email,
			FirstName: warehouse.Manager.FirstName,
			LastName:  warehouse.Manager.LastName,
			Phone:     warehouse.Manager.Phone,
			IsActive:  warehouse.Manager.IsActive,
			CreatedAt: warehouse.Manager.CreatedAt,
			UpdatedAt: warehouse.Manager.UpdatedAt,
		}
	}

	return response
}

// MovementService 库存移动服务接口
type MovementService interface {
	CRUDService[models.Movement, dto.MovementCreateRequest, dto.MovementCreateRequest, dto.MovementResponse]
	CreateMovement(ctx context.Context, req dto.MovementCreateRequest) (*dto.MovementResponse, error)
	GetMovementByID(ctx context.Context, id uint) (*dto.MovementResponse, error)
	ListMovements(ctx context.Context, page, limit int) (*dto.PaginatedResponse[dto.MovementResponse], error)
	GetMovementsByItemID(ctx context.Context, itemID uint, page, limit int) (*dto.PaginatedResponse[dto.MovementResponse], error)
	GetMovementsByWarehouseID(ctx context.Context, warehouseID uint, page, limit int) (*dto.PaginatedResponse[dto.MovementResponse], error)
	GetMovementsByType(ctx context.Context, movementType string, page, limit int) (*dto.PaginatedResponse[dto.MovementResponse], error)
}

// MovementServiceImpl 库存移动服务实现
type MovementServiceImpl struct {
	*BaseService
	movementRepo  repositories.MovementRepository
	stockRepo     repositories.StockRepository
	itemRepo      repositories.ItemRepository
	warehouseRepo repositories.WarehouseRepository
}

// NewMovementService 创建库存移动服务
func NewMovementService(
	movementRepo repositories.MovementRepository,
	stockRepo repositories.StockRepository,
	itemRepo repositories.ItemRepository,
	warehouseRepo repositories.WarehouseRepository,
) MovementService {
	config := &BaseServiceConfig{
		EnableAudit:      true,
		EnableValidation: true,
		EnableCache:      true,
		CacheExpiry:      time.Hour,
		EnableMetrics:    true,
	}
	
	return &MovementServiceImpl{
		BaseService:   NewBaseService(config),
		movementRepo:  movementRepo,
		stockRepo:     stockRepo,
		itemRepo:      itemRepo,
		warehouseRepo: warehouseRepo,
	}
}

// Create 实现 CRUDService 接口的 Create 方法
func (s *MovementServiceImpl) Create(ctx context.Context, req *dto.MovementCreateRequest) (*dto.CreateResponse, error) {
	resp, err := s.CreateMovement(ctx, *req)
	if err != nil {
		return nil, err
	}
	return s.CreateCreateResponse(resp.ID, resp, "库存移动创建成功"), nil
}

// GetByID 实现 CRUDService 接口的 GetByID 方法
func (s *MovementServiceImpl) GetByID(ctx context.Context, id uint) (*dto.MovementResponse, error) {
	return s.GetMovementByID(ctx, id)
}

// Update 实现 CRUDService 接口的 Update 方法
func (s *MovementServiceImpl) Update(ctx context.Context, id uint, req *dto.MovementCreateRequest) (*dto.UpdateResponse, error) {
	// Movement 通常不允许更新，返回错误
	return nil, errors.New("movement records cannot be updated")
}

// Delete 实现 CRUDService 接口的 Delete 方法
func (s *MovementServiceImpl) Delete(ctx context.Context, id uint) (*dto.DeleteResponse, error) {
	// Movement 通常不允许删除，返回错误
	return nil, errors.New("movement records cannot be deleted")
}

// List 实现 CRUDService 接口的 List 方法
func (s *MovementServiceImpl) List(ctx context.Context, req *dto.PaginationRequest) (*dto.PaginatedResponse[dto.MovementResponse], error) {
	page := 1
	limit := 10
	if req.Page > 0 {
		page = req.Page
	}
	if req.PageSize > 0 {
		limit = req.PageSize
	}
	return s.ListMovements(ctx, page, limit)
}

// Search 实现 CRUDService 接口的 Search 方法
func (s *MovementServiceImpl) Search(ctx context.Context, req *dto.SearchRequest) (*dto.PaginatedResponse[dto.MovementResponse], error) {
	// 使用 List 方法实现搜索，后续可以在 repository 层添加搜索功能
	return s.List(ctx, &req.PaginationRequest)
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

// ListMovements 获取库存移动列表
func (s *MovementServiceImpl) ListMovements(ctx context.Context, page, limit int) (*dto.PaginatedResponse[dto.MovementResponse], error) {
	req := &dto.PaginationRequest{
		Page:     page,
		PageSize: limit,
	}
	movements, total, err := s.movementRepo.List(ctx, &common.QueryOptions{
		Pagination: req,
	})
	if err != nil {
		return nil, fmt.Errorf("获取库存移动列表失败: %w", err)
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
			Type:        "", // 模型中没有Type字段，使用空字符串
			MinStock:    0,  // 模型中没有MinStock字段
			MaxStock:    0,  // 模型中没有MaxStock字段
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
