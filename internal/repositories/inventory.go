package repositories

import (
	"context"
	"errors"
	"github.com/galaxyerp/galaxyErp/internal/models"
	"gorm.io/gorm"
)

// ItemRepository 物料仓储接口
type ItemRepository interface {
	BaseRepository[models.Item]
	GetBySKU(ctx context.Context, sku string) (*models.Item, error)
	ListItems(ctx context.Context, offset, limit int) ([]*models.Item, int64, error)
	Search(ctx context.Context, query string, offset, limit int) ([]*models.Item, int64, error)
}

// ItemRepositoryImpl 物料仓储实现
type ItemRepositoryImpl struct {
	BaseRepository[models.Item]
	db *gorm.DB
}

// NewItemRepository 创建物料仓储实例
func NewItemRepository(db *gorm.DB) ItemRepository {
	return &ItemRepositoryImpl{
		BaseRepository: NewBaseRepository[models.Item](db),
		db:             db,
	}
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

// ListItems 获取物料列表
func (r *ItemRepositoryImpl) ListItems(ctx context.Context, offset, limit int) ([]*models.Item, int64, error) {
	var items []*models.Item
	var total int64

	// 获取总数
	if err := r.db.WithContext(ctx).Model(&models.Item{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&items).Error
	if err != nil {
		return nil, 0, err
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
	BaseRepository[models.Stock]
	GetByItemID(ctx context.Context, itemID uint) ([]*models.Stock, error)
}

// StockRepositoryImpl 库存仓储实现
type StockRepositoryImpl struct {
	BaseRepository[models.Stock]
	db *gorm.DB
}

// NewStockRepository 创建库存仓储实例
func NewStockRepository(db *gorm.DB) StockRepository {
	return &StockRepositoryImpl{
		BaseRepository: NewBaseRepository[models.Stock](db),
		db:             db,
	}
}

// GetByItemID 根据物料ID获取库存
func (r *StockRepositoryImpl) GetByItemID(ctx context.Context, itemID uint) ([]*models.Stock, error) {
	var stocks []*models.Stock
	err := r.db.WithContext(ctx).Preload("Item").Preload("Warehouse").Where("item_id = ?", itemID).Find(&stocks).Error
	return stocks, err
}

// WarehouseRepository 仓库仓储接口
type WarehouseRepository interface {
	BaseRepository[models.Warehouse]
}

// WarehouseRepositoryImpl 仓库仓储实现
type WarehouseRepositoryImpl struct {
	BaseRepository[models.Warehouse]
	db *gorm.DB
}

// NewWarehouseRepository 创建仓库仓储实例
func NewWarehouseRepository(db *gorm.DB) WarehouseRepository {
	return &WarehouseRepositoryImpl{
		BaseRepository: NewBaseRepository[models.Warehouse](db),
		db:             db,
	}
}

// MovementRepository 库存移动仓储接口
type MovementRepository interface {
	BaseRepository[models.Movement]
	GetByItemID(ctx context.Context, itemID uint, offset, limit int) ([]*models.Movement, int64, error)
	GetByWarehouseID(ctx context.Context, warehouseID uint, offset, limit int) ([]*models.Movement, int64, error)
	GetByType(ctx context.Context, movementType string, offset, limit int) ([]*models.Movement, int64, error)
}

// MovementRepositoryImpl 库存移动仓储实现
type MovementRepositoryImpl struct {
	BaseRepository[models.Movement]
	db *gorm.DB
}

// NewMovementRepository 创建库存移动仓储实例
func NewMovementRepository(db *gorm.DB) MovementRepository {
	return &MovementRepositoryImpl{
		BaseRepository: NewBaseRepository[models.Movement](db),
		db:             db,
	}
}

// GetByItemID 根据物料ID获取库存移动记录
func (r *MovementRepositoryImpl) GetByItemID(ctx context.Context, itemID uint, offset, limit int) ([]*models.Movement, int64, error) {
	var movements []*models.Movement
	var total int64

	// 获取总数
	if err := r.db.WithContext(ctx).Model(&models.Movement{}).Where("item_id = ?", itemID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := r.db.WithContext(ctx).
		Preload("Item").
		Preload("Warehouse").
		Preload("User").
		Where("item_id = ?", itemID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&movements).Error

	return movements, total, err
}

// GetByWarehouseID 根据仓库ID获取库存移动记录
func (r *MovementRepositoryImpl) GetByWarehouseID(ctx context.Context, warehouseID uint, offset, limit int) ([]*models.Movement, int64, error) {
	var movements []*models.Movement
	var total int64

	// 获取总数
	if err := r.db.WithContext(ctx).Model(&models.Movement{}).Where("warehouse_id = ?", warehouseID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := r.db.WithContext(ctx).
		Preload("Item").
		Preload("Warehouse").
		Preload("User").
		Where("warehouse_id = ?", warehouseID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&movements).Error

	return movements, total, err
}

// GetByType 根据移动类型获取库存移动记录
func (r *MovementRepositoryImpl) GetByType(ctx context.Context, movementType string, offset, limit int) ([]*models.Movement, int64, error) {
	var movements []*models.Movement
	var total int64

	// 获取总数
	if err := r.db.WithContext(ctx).Model(&models.Movement{}).Where("movement_type = ?", movementType).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := r.db.WithContext(ctx).
		Preload("Item").
		Preload("Warehouse").
		Preload("User").
		Where("movement_type = ?", movementType).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&movements).Error

	return movements, total, err
}
