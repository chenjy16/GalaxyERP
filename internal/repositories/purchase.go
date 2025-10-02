package repositories

import (
	"context"
	"errors"
	"github.com/galaxyerp/galaxyErp/internal/models"
	"gorm.io/gorm"
)

// SupplierRepository 供应商仓储接口
type SupplierRepository interface {
	Create(ctx context.Context, supplier *models.Supplier) error
	GetByID(ctx context.Context, id uint) (*models.Supplier, error)
	Update(ctx context.Context, supplier *models.Supplier) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, offset, limit int) ([]*models.Supplier, int64, error)
	Search(ctx context.Context, query string, offset, limit int) ([]*models.Supplier, int64, error)
}

// SupplierRepositoryImpl 供应商仓储实现
type SupplierRepositoryImpl struct {
	db *gorm.DB
}

// NewSupplierRepository 创建供应商仓储实例
func NewSupplierRepository(db *gorm.DB) SupplierRepository {
	return &SupplierRepositoryImpl{
		db: db,
	}
}

// Create 创建供应商
func (r *SupplierRepositoryImpl) Create(ctx context.Context, supplier *models.Supplier) error {
	return r.db.WithContext(ctx).Create(supplier).Error
}

// GetByID 根据ID获取供应商
func (r *SupplierRepositoryImpl) GetByID(ctx context.Context, id uint) (*models.Supplier, error) {
	var supplier models.Supplier
	err := r.db.WithContext(ctx).First(&supplier, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &supplier, nil
}

// Update 更新供应商
func (r *SupplierRepositoryImpl) Update(ctx context.Context, supplier *models.Supplier) error {
	return r.db.WithContext(ctx).Save(supplier).Error
}

// Delete 删除供应商
func (r *SupplierRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Supplier{}, id).Error
}

// List 获取供应商列表
func (r *SupplierRepositoryImpl) List(ctx context.Context, offset, limit int) ([]*models.Supplier, int64, error) {
	var suppliers []*models.Supplier
	var total int64

	// 获取总数
	if err := r.db.WithContext(ctx).Model(&models.Supplier{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&suppliers).Error
	if err != nil {
		return nil, 0, err
	}

	return suppliers, total, nil
}

// Search 搜索供应商
func (r *SupplierRepositoryImpl) Search(ctx context.Context, query string, offset, limit int) ([]*models.Supplier, int64, error) {
	var suppliers []*models.Supplier
	var total int64

	searchQuery := "%" + query + "%"

	// 获取总数
	if err := r.db.WithContext(ctx).Model(&models.Supplier{}).
		Where("name LIKE ? OR email LIKE ? OR phone LIKE ?",
			searchQuery, searchQuery, searchQuery).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := r.db.WithContext(ctx).
		Where("name LIKE ? OR email LIKE ? OR phone LIKE ?",
			searchQuery, searchQuery, searchQuery).
		Offset(offset).Limit(limit).Find(&suppliers).Error
	if err != nil {
		return nil, 0, err
	}

	return suppliers, total, nil
}

// PurchaseRequestRepository 采购申请仓储接口
type PurchaseRequestRepository interface {
	Create(ctx context.Context, request *models.PurchaseRequest) error
	GetByID(ctx context.Context, id uint) (*models.PurchaseRequest, error)
	Update(ctx context.Context, request *models.PurchaseRequest) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, offset, limit int) ([]*models.PurchaseRequest, int64, error)
	GetByDepartmentID(ctx context.Context, departmentID uint, offset, limit int) ([]*models.PurchaseRequest, int64, error)
	GetByStatus(ctx context.Context, status string, offset, limit int) ([]*models.PurchaseRequest, int64, error)
}

// PurchaseRequestRepositoryImpl 采购申请仓储实现
type PurchaseRequestRepositoryImpl struct {
	db *gorm.DB
}

// NewPurchaseRequestRepository 创建采购申请仓储实例
func NewPurchaseRequestRepository(db *gorm.DB) PurchaseRequestRepository {
	return &PurchaseRequestRepositoryImpl{
		db: db,
	}
}

// Create 创建采购申请
func (r *PurchaseRequestRepositoryImpl) Create(ctx context.Context, request *models.PurchaseRequest) error {
	return r.db.WithContext(ctx).Create(request).Error
}

// GetByID 根据ID获取采购申请
func (r *PurchaseRequestRepositoryImpl) GetByID(ctx context.Context, id uint) (*models.PurchaseRequest, error) {
	var request models.PurchaseRequest
	err := r.db.WithContext(ctx).Preload("Items").First(&request, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &request, nil
}

// Update 更新采购申请
func (r *PurchaseRequestRepositoryImpl) Update(ctx context.Context, request *models.PurchaseRequest) error {
	return r.db.WithContext(ctx).Save(request).Error
}

// Delete 删除采购申请
func (r *PurchaseRequestRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.PurchaseRequest{}, id).Error
}

// List 获取采购申请列表
func (r *PurchaseRequestRepositoryImpl) List(ctx context.Context, offset, limit int) ([]*models.PurchaseRequest, int64, error) {
	var requests []*models.PurchaseRequest
	var total int64

	// 获取总数
	if err := r.db.WithContext(ctx).Model(&models.PurchaseRequest{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&requests).Error
	if err != nil {
		return nil, 0, err
	}

	return requests, total, nil
}

// GetByDepartmentID 根据部门ID获取采购申请
func (r *PurchaseRequestRepositoryImpl) GetByDepartmentID(ctx context.Context, departmentID uint, offset, limit int) ([]*models.PurchaseRequest, int64, error) {
	var requests []*models.PurchaseRequest
	var total int64

	// 获取总数
	if err := r.db.WithContext(ctx).Model(&models.PurchaseRequest{}).
		Where("department_id = ?", departmentID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := r.db.WithContext(ctx).Where("department_id = ?", departmentID).
		Offset(offset).Limit(limit).Find(&requests).Error
	if err != nil {
		return nil, 0, err
	}

	return requests, total, nil
}

// GetByStatus 根据状态获取采购申请
func (r *PurchaseRequestRepositoryImpl) GetByStatus(ctx context.Context, status string, offset, limit int) ([]*models.PurchaseRequest, int64, error) {
	var requests []*models.PurchaseRequest
	var total int64

	// 获取总数
	if err := r.db.WithContext(ctx).Model(&models.PurchaseRequest{}).
		Where("status = ?", status).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := r.db.WithContext(ctx).Where("status = ?", status).
		Offset(offset).Limit(limit).Find(&requests).Error
	if err != nil {
		return nil, 0, err
	}

	return requests, total, nil
}

// PurchaseOrderRepository 采购订单仓储接口
type PurchaseOrderRepository interface {
	Create(ctx context.Context, order *models.PurchaseOrder) error
	GetByID(ctx context.Context, id uint) (*models.PurchaseOrder, error)
	Update(ctx context.Context, order *models.PurchaseOrder) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, offset, limit int) ([]*models.PurchaseOrder, int64, error)
	GetBySupplierID(ctx context.Context, supplierID uint, offset, limit int) ([]*models.PurchaseOrder, int64, error)
}

// PurchaseOrderRepositoryImpl 采购订单仓储实现
type PurchaseOrderRepositoryImpl struct {
	db *gorm.DB
}

// NewPurchaseOrderRepository 创建采购订单仓储实例
func NewPurchaseOrderRepository(db *gorm.DB) PurchaseOrderRepository {
	return &PurchaseOrderRepositoryImpl{
		db: db,
	}
}

// Create 创建采购订单
func (r *PurchaseOrderRepositoryImpl) Create(ctx context.Context, order *models.PurchaseOrder) error {
	return r.db.WithContext(ctx).Create(order).Error
}

// GetByID 根据ID获取采购订单
func (r *PurchaseOrderRepositoryImpl) GetByID(ctx context.Context, id uint) (*models.PurchaseOrder, error) {
	var order models.PurchaseOrder
	err := r.db.WithContext(ctx).Preload("Supplier").Preload("Items").Preload("Items.Item").First(&order, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &order, nil
}

// Update 更新采购订单
func (r *PurchaseOrderRepositoryImpl) Update(ctx context.Context, order *models.PurchaseOrder) error {
	return r.db.WithContext(ctx).Save(order).Error
}

// Delete 删除采购订单
func (r *PurchaseOrderRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.PurchaseOrder{}, id).Error
}

// List 获取采购订单列表
func (r *PurchaseOrderRepositoryImpl) List(ctx context.Context, offset, limit int) ([]*models.PurchaseOrder, int64, error) {
	var orders []*models.PurchaseOrder
	var total int64

	// 获取总数
	if err := r.db.WithContext(ctx).Model(&models.PurchaseOrder{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := r.db.WithContext(ctx).Preload("Supplier").Offset(offset).Limit(limit).Find(&orders).Error
	if err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

// GetBySupplierID 根据供应商ID获取采购订单
func (r *PurchaseOrderRepositoryImpl) GetBySupplierID(ctx context.Context, supplierID uint, offset, limit int) ([]*models.PurchaseOrder, int64, error) {
	var orders []*models.PurchaseOrder
	var total int64

	// 获取总数
	if err := r.db.WithContext(ctx).Model(&models.PurchaseOrder{}).
		Where("supplier_id = ?", supplierID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := r.db.WithContext(ctx).Preload("Supplier").
		Where("supplier_id = ?", supplierID).
		Offset(offset).Limit(limit).Find(&orders).Error
	if err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}
