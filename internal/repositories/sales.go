package repositories

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"github.com/galaxyerp/galaxyErp/internal/models"
)

// CustomerRepository 客户仓储接口
type CustomerRepository interface {
	Create(ctx context.Context, customer *models.Customer) error
	GetByID(ctx context.Context, id uint) (*models.Customer, error)
	Update(ctx context.Context, customer *models.Customer) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, offset, limit int) ([]*models.Customer, int64, error)
	Search(ctx context.Context, query string, offset, limit int) ([]*models.Customer, int64, error)
}

// CustomerRepositoryImpl 客户仓储实现
type CustomerRepositoryImpl struct {
	db *gorm.DB
}

// NewCustomerRepository 创建客户仓储实例
func NewCustomerRepository(db *gorm.DB) CustomerRepository {
	return &CustomerRepositoryImpl{
		db: db,
	}
}

// Create 创建客户
func (r *CustomerRepositoryImpl) Create(ctx context.Context, customer *models.Customer) error {
	return r.db.WithContext(ctx).Create(customer).Error
}

// GetByID 根据ID获取客户
func (r *CustomerRepositoryImpl) GetByID(ctx context.Context, id uint) (*models.Customer, error) {
	var customer models.Customer
	err := r.db.WithContext(ctx).First(&customer, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &customer, nil
}

// Update 更新客户
func (r *CustomerRepositoryImpl) Update(ctx context.Context, customer *models.Customer) error {
	return r.db.WithContext(ctx).Save(customer).Error
}

// Delete 删除客户
func (r *CustomerRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Customer{}, id).Error
}

// List 获取客户列表
func (r *CustomerRepositoryImpl) List(ctx context.Context, offset, limit int) ([]*models.Customer, int64, error) {
	var customers []*models.Customer
	var total int64
	
	// 获取总数
	if err := r.db.WithContext(ctx).Model(&models.Customer{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// 获取分页数据
	err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&customers).Error
	if err != nil {
		return nil, 0, err
	}
	
	return customers, total, nil
}

// Search 搜索客户
func (r *CustomerRepositoryImpl) Search(ctx context.Context, query string, offset, limit int) ([]*models.Customer, int64, error) {
	var customers []*models.Customer
	var total int64
	
	searchQuery := "%" + query + "%"
	
	// 获取总数
	if err := r.db.WithContext(ctx).Model(&models.Customer{}).
		Where("name LIKE ? OR email LIKE ? OR phone LIKE ?", 
			searchQuery, searchQuery, searchQuery).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// 获取分页数据
	err := r.db.WithContext(ctx).
		Where("name LIKE ? OR email LIKE ? OR phone LIKE ?", 
			searchQuery, searchQuery, searchQuery).
		Offset(offset).Limit(limit).Find(&customers).Error
	if err != nil {
		return nil, 0, err
	}
	
	return customers, total, nil
}

// SalesOrderRepository 销售订单仓储接口
type SalesOrderRepository interface {
	Create(ctx context.Context, order *models.SalesOrder) error
	GetByID(ctx context.Context, id uint) (*models.SalesOrder, error)
	Update(ctx context.Context, order *models.SalesOrder) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, offset, limit int) ([]*models.SalesOrder, int64, error)
	GetByCustomerID(ctx context.Context, customerID uint, offset, limit int) ([]*models.SalesOrder, int64, error)
}

// SalesOrderRepositoryImpl 销售订单仓储实现
type SalesOrderRepositoryImpl struct {
	db *gorm.DB
}

// NewSalesOrderRepository 创建销售订单仓储实例
func NewSalesOrderRepository(db *gorm.DB) SalesOrderRepository {
	return &SalesOrderRepositoryImpl{
		db: db,
	}
}

// Create 创建销售订单
func (r *SalesOrderRepositoryImpl) Create(ctx context.Context, order *models.SalesOrder) error {
	return r.db.WithContext(ctx).Create(order).Error
}

// GetByID 根据ID获取销售订单
func (r *SalesOrderRepositoryImpl) GetByID(ctx context.Context, id uint) (*models.SalesOrder, error) {
	var order models.SalesOrder
	err := r.db.WithContext(ctx).Preload("Customer").Preload("Items").Preload("Items.Product").First(&order, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &order, nil
}

// Update 更新销售订单
func (r *SalesOrderRepositoryImpl) Update(ctx context.Context, order *models.SalesOrder) error {
	return r.db.WithContext(ctx).Save(order).Error
}

// Delete 删除销售订单
func (r *SalesOrderRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.SalesOrder{}, id).Error
}

// List 获取销售订单列表
func (r *SalesOrderRepositoryImpl) List(ctx context.Context, offset, limit int) ([]*models.SalesOrder, int64, error) {
	var orders []*models.SalesOrder
	var total int64
	
	// 获取总数
	if err := r.db.WithContext(ctx).Model(&models.SalesOrder{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// 获取分页数据
	err := r.db.WithContext(ctx).Preload("Customer").Offset(offset).Limit(limit).Find(&orders).Error
	if err != nil {
		return nil, 0, err
	}
	
	return orders, total, nil
}

// QuotationRepository 报价仓储接口
type QuotationRepository interface {
	Create(ctx context.Context, quotation *models.Quotation) error
	GetByID(ctx context.Context, id uint) (*models.Quotation, error)
	Update(ctx context.Context, quotation *models.Quotation) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, offset, limit int) ([]*models.Quotation, int64, error)
	GetByCustomerID(ctx context.Context, customerID uint, offset, limit int) ([]*models.Quotation, int64, error)
	Search(ctx context.Context, query string, offset, limit int) ([]*models.Quotation, int64, error)
}

// QuotationRepositoryImpl 报价仓储实现
type QuotationRepositoryImpl struct {
	db *gorm.DB
}

// NewQuotationRepository 创建报价仓储实例
func NewQuotationRepository(db *gorm.DB) QuotationRepository {
	return &QuotationRepositoryImpl{
		db: db,
	}
}

// Create 创建报价
func (r *QuotationRepositoryImpl) Create(ctx context.Context, quotation *models.Quotation) error {
	return r.db.WithContext(ctx).Create(quotation).Error
}

// GetByID 根据ID获取报价
func (r *QuotationRepositoryImpl) GetByID(ctx context.Context, id uint) (*models.Quotation, error) {
	var quotation models.Quotation
	err := r.db.WithContext(ctx).Preload("Items").Preload("Customer").First(&quotation, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &quotation, nil
}

// Update 更新报价
func (r *QuotationRepositoryImpl) Update(ctx context.Context, quotation *models.Quotation) error {
	return r.db.WithContext(ctx).Save(quotation).Error
}

// Delete 删除报价
func (r *QuotationRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Quotation{}, id).Error
}

// List 获取报价列表
func (r *QuotationRepositoryImpl) List(ctx context.Context, offset, limit int) ([]*models.Quotation, int64, error) {
	var quotations []*models.Quotation
	var total int64
	
	// 获取总数
	if err := r.db.WithContext(ctx).Model(&models.Quotation{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// 获取分页数据
	err := r.db.WithContext(ctx).Preload("Customer").Offset(offset).Limit(limit).Find(&quotations).Error
	if err != nil {
		return nil, 0, err
	}
	
	return quotations, total, nil
}

// GetByCustomerID 根据客户ID获取报价列表
func (r *QuotationRepositoryImpl) GetByCustomerID(ctx context.Context, customerID uint, offset, limit int) ([]*models.Quotation, int64, error) {
	var quotations []*models.Quotation
	var total int64
	
	// 获取总数
	if err := r.db.WithContext(ctx).Model(&models.Quotation{}).Where("customer_id = ?", customerID).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// 获取分页数据
	err := r.db.WithContext(ctx).Preload("Customer").Where("customer_id = ?", customerID).Offset(offset).Limit(limit).Find(&quotations).Error
	if err != nil {
		return nil, 0, err
	}
	
	return quotations, total, nil
}

// Search 搜索报价
func (r *QuotationRepositoryImpl) Search(ctx context.Context, query string, offset, limit int) ([]*models.Quotation, int64, error) {
	var quotations []*models.Quotation
	var total int64
	
	searchQuery := "%" + query + "%"
	
	// 获取总数
	if err := r.db.WithContext(ctx).Model(&models.Quotation{}).Joins("LEFT JOIN customers ON quotations.customer_id = customers.id").Where("quotations.quotation_number LIKE ? OR customers.name LIKE ?", searchQuery, searchQuery).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// 获取分页数据
	err := r.db.WithContext(ctx).Preload("Customer").Joins("LEFT JOIN customers ON quotations.customer_id = customers.id").Where("quotations.quotation_number LIKE ? OR customers.name LIKE ?", searchQuery, searchQuery).Offset(offset).Limit(limit).Find(&quotations).Error
	if err != nil {
		return nil, 0, err
	}
	
	return quotations, total, nil
}

// GetByCustomerID 根据客户ID获取销售订单
func (r *SalesOrderRepositoryImpl) GetByCustomerID(ctx context.Context, customerID uint, offset, limit int) ([]*models.SalesOrder, int64, error) {
	var orders []*models.SalesOrder
	var total int64
	
	// 获取总数
	if err := r.db.WithContext(ctx).Model(&models.SalesOrder{}).
		Where("customer_id = ?", customerID).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// 获取分页数据
	err := r.db.WithContext(ctx).Preload("Customer").
		Where("customer_id = ?", customerID).
		Offset(offset).Limit(limit).Find(&orders).Error
	if err != nil {
		return nil, 0, err
	}
	
	return orders, total, nil
}