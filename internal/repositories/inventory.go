package repositories

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"github.com/galaxyerp/galaxyErp/internal/models"
)

// ItemRepository 物料仓储接口
type ItemRepository interface {
	Create(ctx context.Context, item *models.Item) error
	GetByID(ctx context.Context, id uint) (*models.Item, error)
	GetBySKU(ctx context.Context, sku string) (*models.Item, error)
	Update(ctx context.Context, item *models.Item) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, offset, limit int) ([]*models.Item, int64, error)
	Search(ctx context.Context, query string, offset, limit int) ([]*models.Item, int64, error)
}

// ItemRepositoryImpl 物料仓储实现
type ItemRepositoryImpl struct {
	db *gorm.DB
}

// NewItemRepository 创建物料仓储实例
func NewItemRepository(db *gorm.DB) ItemRepository {
	return &ItemRepositoryImpl{
		db: db,
	}
}

// Create 创建物料
func (r *ItemRepositoryImpl) Create(ctx context.Context, item *models.Item) error {
	return r.db.WithContext(ctx).Create(item).Error
}

// GetByID 根据ID获取物料
func (r *ItemRepositoryImpl) GetByID(ctx context.Context, id uint) (*models.Item, error) {
	var item models.Item
	err := r.db.WithContext(ctx).First(&item, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

// GetBySKU 根据SKU获取物料
func (r *ItemRepositoryImpl) GetBySKU(ctx context.Context, sku string) (*models.Item, error) {
	var item models.Item
	err := r.db.WithContext(ctx).Where("sku = ?", sku).First(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

// Update 更新物料
func (r *ItemRepositoryImpl) Update(ctx context.Context, item *models.Item) error {
	return r.db.WithContext(ctx).Save(item).Error
}

// Delete 删除物料
func (r *ItemRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Item{}, id).Error
}

// List 获取物料列表
func (r *ItemRepositoryImpl) List(ctx context.Context, offset, limit int) ([]*models.Item, int64, error) {
	var legacyItems []*models.LegacyItem
	var total int64
	
	// 获取总数
	if err := r.db.WithContext(ctx).Model(&models.LegacyItem{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// 获取分页数据
	err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&legacyItems).Error
	if err != nil {
		return nil, 0, err
	}
	
	// 转换为新的模型格式
	items := make([]*models.Item, len(legacyItems))
	for i, legacyItem := range legacyItems {
		var deletedAt gorm.DeletedAt
		if legacyItem.DeletedAt != nil {
			deletedAt = gorm.DeletedAt{Time: *legacyItem.DeletedAt, Valid: true}
		}
		
		items[i] = &models.Item{
			BaseModel: models.BaseModel{
				ID:        legacyItem.ID,
				CreatedAt: legacyItem.CreatedAt,
				UpdatedAt: legacyItem.UpdatedAt,
				DeletedAt: deletedAt,
			},
			Code:         legacyItem.Code,
			Name:         legacyItem.Name,
			Description:  legacyItem.Description,
			Category:     legacyItem.Category,
			Unit:         legacyItem.Unit,
			Cost:         legacyItem.Cost,
			Price:        legacyItem.Price,
			ReorderLevel: legacyItem.ReorderLevel,
			IsActive:     legacyItem.IsActive,
		}
	}
	
	return items, total, nil
}

// Search 搜索物料
func (r *ItemRepositoryImpl) Search(ctx context.Context, query string, offset, limit int) ([]*models.Item, int64, error) {
	var items []*models.Item
	var total int64
	
	searchQuery := "%" + query + "%"
	
	// 获取总数
	if err := r.db.WithContext(ctx).Model(&models.Item{}).Where("name LIKE ? OR code LIKE ? OR sku LIKE ?", searchQuery, searchQuery, searchQuery).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// 获取分页数据
	err := r.db.WithContext(ctx).Preload("Category").Preload("Unit").Where("name LIKE ? OR code LIKE ? OR sku LIKE ?", searchQuery, searchQuery, searchQuery).Offset(offset).Limit(limit).Find(&items).Error
	if err != nil {
		return nil, 0, err
	}
	
	return items, total, nil
}

// StockRepository 库存仓储接口
type StockRepository interface {
	Create(ctx context.Context, stock *models.Stock) error
	GetByID(ctx context.Context, id uint) (*models.Stock, error)
	Update(ctx context.Context, stock *models.Stock) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, offset, limit int) ([]*models.Stock, int64, error)
	GetByItemID(ctx context.Context, itemID uint) ([]*models.Stock, error)
}

// StockRepositoryImpl 库存仓储实现
type StockRepositoryImpl struct {
	db *gorm.DB
}

// NewStockRepository 创建库存仓储实例
func NewStockRepository(db *gorm.DB) StockRepository {
	return &StockRepositoryImpl{
		db: db,
	}
}

// Create 创建库存
func (r *StockRepositoryImpl) Create(ctx context.Context, stock *models.Stock) error {
	return r.db.WithContext(ctx).Create(stock).Error
}

// GetByID 根据ID获取库存
func (r *StockRepositoryImpl) GetByID(ctx context.Context, id uint) (*models.Stock, error) {
	var stock models.Stock
	err := r.db.WithContext(ctx).Preload("Item").Preload("Warehouse").First(&stock, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &stock, nil
}

// Update 更新库存
func (r *StockRepositoryImpl) Update(ctx context.Context, stock *models.Stock) error {
	return r.db.WithContext(ctx).Save(stock).Error
}

// Delete 删除库存
func (r *StockRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Stock{}, id).Error
}

// List 获取库存列表
func (r *StockRepositoryImpl) List(ctx context.Context, offset, limit int) ([]*models.Stock, int64, error) {
	var stocks []*models.Stock
	var total int64
	
	// 获取总数
	if err := r.db.WithContext(ctx).Model(&models.Stock{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// 获取分页数据
	err := r.db.WithContext(ctx).Preload("Item").Preload("Warehouse").Offset(offset).Limit(limit).Find(&stocks).Error
	if err != nil {
		return nil, 0, err
	}
	
	return stocks, total, nil
}

// GetByItemID 根据物料ID获取库存
func (r *StockRepositoryImpl) GetByItemID(ctx context.Context, itemID uint) ([]*models.Stock, error) {
	var stocks []*models.Stock
	err := r.db.WithContext(ctx).Preload("Item").Preload("Warehouse").Where("item_id = ?", itemID).Find(&stocks).Error
	if err != nil {
		return nil, err
	}
	return stocks, nil
}

// WarehouseRepository 仓库仓储接口
type WarehouseRepository interface {
	Create(ctx context.Context, warehouse *models.Warehouse) error
	GetByID(ctx context.Context, id uint) (*models.Warehouse, error)
	Update(ctx context.Context, warehouse *models.Warehouse) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, offset, limit int) ([]*models.Warehouse, int64, error)
}

// WarehouseRepositoryImpl 仓库仓储实现
type WarehouseRepositoryImpl struct {
	db *gorm.DB
}

// NewWarehouseRepository 创建仓库仓储实例
func NewWarehouseRepository(db *gorm.DB) WarehouseRepository {
	return &WarehouseRepositoryImpl{
		db: db,
	}
}

// Create 创建仓库
func (r *WarehouseRepositoryImpl) Create(ctx context.Context, warehouse *models.Warehouse) error {
	return r.db.WithContext(ctx).Create(warehouse).Error
}

// GetByID 根据ID获取仓库
func (r *WarehouseRepositoryImpl) GetByID(ctx context.Context, id uint) (*models.Warehouse, error) {
	var warehouse models.Warehouse
	err := r.db.WithContext(ctx).Preload("Locations").First(&warehouse, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &warehouse, nil
}

// Update 更新仓库
func (r *WarehouseRepositoryImpl) Update(ctx context.Context, warehouse *models.Warehouse) error {
	return r.db.WithContext(ctx).Save(warehouse).Error
}

// Delete 删除仓库
func (r *WarehouseRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Warehouse{}, id).Error
}

// List 获取仓库列表
func (r *WarehouseRepositoryImpl) List(ctx context.Context, offset, limit int) ([]*models.Warehouse, int64, error) {
	var warehouses []*models.Warehouse
	var total int64
	
	// 获取总数
	if err := r.db.WithContext(ctx).Model(&models.Warehouse{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// 获取分页数据
	err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&warehouses).Error
	if err != nil {
		return nil, 0, err
	}
	
	return warehouses, total, nil
}

// MovementRepository 库存移动仓储接口
type MovementRepository interface {
	Create(ctx context.Context, movement *models.Movement) error
	GetByID(ctx context.Context, id uint) (*models.Movement, error)
	Update(ctx context.Context, movement *models.Movement) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, offset, limit int) ([]*models.Movement, int64, error)
	GetByItemID(ctx context.Context, itemID uint, offset, limit int) ([]*models.Movement, int64, error)
	GetByWarehouseID(ctx context.Context, warehouseID uint, offset, limit int) ([]*models.Movement, int64, error)
	GetByType(ctx context.Context, movementType string, offset, limit int) ([]*models.Movement, int64, error)
}

// MovementRepositoryImpl 库存移动仓储实现
type MovementRepositoryImpl struct {
	db *gorm.DB
}

// NewMovementRepository 创建库存移动仓储实例
func NewMovementRepository(db *gorm.DB) MovementRepository {
	return &MovementRepositoryImpl{
		db: db,
	}
}

// Create 创建库存移动记录
func (r *MovementRepositoryImpl) Create(ctx context.Context, movement *models.Movement) error {
	return r.db.WithContext(ctx).Create(movement).Error
}

// GetByID 根据ID获取库存移动记录
func (r *MovementRepositoryImpl) GetByID(ctx context.Context, id uint) (*models.Movement, error) {
	var movement models.Movement
	err := r.db.WithContext(ctx).
		Preload("Item").
		Preload("Warehouse").
		First(&movement, id).Error
	if err != nil {
		return nil, err
	}
	return &movement, nil
}

// Update 更新库存移动记录
func (r *MovementRepositoryImpl) Update(ctx context.Context, movement *models.Movement) error {
	return r.db.WithContext(ctx).Save(movement).Error
}

// Delete 删除库存移动记录
func (r *MovementRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Movement{}, id).Error
}

// List 获取库存移动记录列表
func (r *MovementRepositoryImpl) List(ctx context.Context, offset, limit int) ([]*models.Movement, int64, error) {
	var movements []*models.Movement
	var total int64
	
	// 获取总数
	err := r.db.WithContext(ctx).Model(&models.Movement{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	
	// 获取数据
	err = r.db.WithContext(ctx).
		Preload("Item").
		Preload("Warehouse").
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&movements).Error
	if err != nil {
		return nil, 0, err
	}
	
	return movements, total, nil
}

// GetByItemID 根据物料ID获取库存移动记录
func (r *MovementRepositoryImpl) GetByItemID(ctx context.Context, itemID uint, offset, limit int) ([]*models.Movement, int64, error) {
	var movements []*models.Movement
	var total int64
	
	query := r.db.WithContext(ctx).Model(&models.Movement{}).Where("item_id = ?", itemID)
	
	// 获取总数
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	
	// 获取数据
	err = query.
		Preload("Item").
		Preload("Warehouse").
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&movements).Error
	if err != nil {
		return nil, 0, err
	}
	
	return movements, total, nil
}

// GetByWarehouseID 根据仓库ID获取库存移动记录
func (r *MovementRepositoryImpl) GetByWarehouseID(ctx context.Context, warehouseID uint, offset, limit int) ([]*models.Movement, int64, error) {
	var movements []*models.Movement
	var total int64
	
	query := r.db.WithContext(ctx).Model(&models.Movement{}).Where("warehouse_id = ?", warehouseID)
	
	// 获取总数
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	
	// 获取数据
	err = query.
		Preload("Item").
		Preload("Warehouse").
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&movements).Error
	if err != nil {
		return nil, 0, err
	}
	
	return movements, total, nil
}

// GetByType 根据移动类型获取库存移动记录
func (r *MovementRepositoryImpl) GetByType(ctx context.Context, movementType string, offset, limit int) ([]*models.Movement, int64, error) {
	var movements []*models.Movement
	var total int64
	
	query := r.db.WithContext(ctx).Model(&models.Movement{}).Where("movement_type = ?", movementType)
	
	// 获取总数
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	
	// 获取数据
	err = query.
		Preload("Item").
		Preload("Warehouse").
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&movements).Error
	if err != nil {
		return nil, 0, err
	}
	
	return movements, total, nil
}