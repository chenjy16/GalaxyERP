package repositories

import (
	"context"

	"github.com/galaxyerp/galaxyErp/internal/models"
	"gorm.io/gorm"
)

// ProductRepository 产品仓储接口
type ProductRepository interface {
	Create(ctx context.Context, product *models.Product) error
	GetByID(ctx context.Context, id uint) (*models.Product, error)
	GetByCode(ctx context.Context, code string) (*models.Product, error)
	Update(ctx context.Context, product *models.Product) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, page, pageSize int) ([]*models.Product, int64, error)
	Search(ctx context.Context, query string, page, pageSize int) ([]*models.Product, int64, error)
}

// ProductRepositoryImpl 产品仓储实现
type ProductRepositoryImpl struct {
	db *gorm.DB
}

// NewProductRepository 创建产品仓储实例
func NewProductRepository(db *gorm.DB) ProductRepository {
	return &ProductRepositoryImpl{
		db: db,
	}
}

// Create 创建产品
func (r *ProductRepositoryImpl) Create(ctx context.Context, product *models.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

// GetByID 根据ID获取产品
func (r *ProductRepositoryImpl) GetByID(ctx context.Context, id uint) (*models.Product, error) {
	var product models.Product
	err := r.db.WithContext(ctx).First(&product, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &product, nil
}

// GetByCode 根据编码获取产品
func (r *ProductRepositoryImpl) GetByCode(ctx context.Context, code string) (*models.Product, error) {
	var product models.Product
	err := r.db.WithContext(ctx).Where("code = ?", code).First(&product).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &product, nil
}

// Update 更新产品
func (r *ProductRepositoryImpl) Update(ctx context.Context, product *models.Product) error {
	return r.db.WithContext(ctx).Save(product).Error
}

// Delete 删除产品
func (r *ProductRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Product{}, id).Error
}

// List 获取产品列表
func (r *ProductRepositoryImpl) List(ctx context.Context, page, pageSize int) ([]*models.Product, int64, error) {
	var products []*models.Product
	var total int64

	// 计算总数
	if err := r.db.WithContext(ctx).Model(&models.Product{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := r.db.WithContext(ctx).
		Offset(offset).
		Limit(pageSize).
		Find(&products).Error

	return products, total, err
}

// Search 搜索产品
func (r *ProductRepositoryImpl) Search(ctx context.Context, query string, page, pageSize int) ([]*models.Product, int64, error) {
	var products []*models.Product
	var total int64

	db := r.db.WithContext(ctx).Model(&models.Product{})
	if query != "" {
		searchPattern := "%" + query + "%"
		db = db.Where("code LIKE ? OR name LIKE ? OR description LIKE ?", searchPattern, searchPattern, searchPattern)
	}

	// 计算总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := db.Offset(offset).
		Limit(pageSize).
		Find(&products).Error

	return products, total, err
}