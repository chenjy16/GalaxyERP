package repositories

import (
	"context"
	"github.com/galaxyerp/galaxyErp/internal/models"
	"gorm.io/gorm"
)

// SupplierRepository 供应商仓储接口
type SupplierRepository interface {
	BaseRepository[models.Supplier]
	Search(ctx context.Context, query string, offset, limit int) ([]*models.Supplier, int64, error)
}

// SupplierRepositoryImpl 供应商仓储实现
type SupplierRepositoryImpl struct {
	BaseRepository[models.Supplier]
	db *gorm.DB
}

// NewSupplierRepository 创建供应商仓储实例
func NewSupplierRepository(db *gorm.DB) SupplierRepository {
	return &SupplierRepositoryImpl{
		BaseRepository: NewBaseRepository[models.Supplier](db),
		db:             db,
	}
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
	BaseRepository[models.PurchaseRequest]
	GetByDepartmentID(ctx context.Context, departmentID uint, offset, limit int) ([]*models.PurchaseRequest, int64, error)
	GetByStatus(ctx context.Context, status string, offset, limit int) ([]*models.PurchaseRequest, int64, error)
}

// PurchaseRequestRepositoryImpl 采购申请仓储实现
type PurchaseRequestRepositoryImpl struct {
	BaseRepository[models.PurchaseRequest]
	db *gorm.DB
}

// NewPurchaseRequestRepository 创建采购申请仓储实例
func NewPurchaseRequestRepository(db *gorm.DB) PurchaseRequestRepository {
	return &PurchaseRequestRepositoryImpl{
		BaseRepository: NewBaseRepository[models.PurchaseRequest](db),
		db:             db,
	}
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
	BaseRepository[models.PurchaseOrder]
	GetBySupplierID(ctx context.Context, supplierID uint, offset, limit int) ([]*models.PurchaseOrder, int64, error)
}

// PurchaseOrderRepositoryImpl 采购订单仓储实现
type PurchaseOrderRepositoryImpl struct {
	BaseRepository[models.PurchaseOrder]
	db *gorm.DB
}

// NewPurchaseOrderRepository 创建采购订单仓储实例
func NewPurchaseOrderRepository(db *gorm.DB) PurchaseOrderRepository {
	return &PurchaseOrderRepositoryImpl{
		BaseRepository: NewBaseRepository[models.PurchaseOrder](db),
		db:             db,
	}
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
