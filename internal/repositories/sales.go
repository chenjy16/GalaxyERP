package repositories

import (
	"context"
	"errors"
	"fmt"
	"github.com/galaxyerp/galaxyErp/internal/dto"
	"github.com/galaxyerp/galaxyErp/internal/models"
	"gorm.io/gorm"
	"time"
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
	CreateItem(ctx context.Context, item *models.SalesOrderItem) error
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
	err := r.db.WithContext(ctx).
		Preload("Customer").
		Preload("CreatedByUser").
		Preload("Items").
		Preload("Items.Item").
		First(&order, id).Error
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

// CreateItem 创建销售订单项目
func (r *SalesOrderRepositoryImpl) CreateItem(ctx context.Context, item *models.SalesOrderItem) error {
	return r.db.WithContext(ctx).Create(item).Error
}

// ==================== 销售发票仓储 ====================

// SalesInvoiceRepository 销售发票仓储接口
type SalesInvoiceRepository interface {
	Create(ctx context.Context, invoice *models.SalesInvoice) error
	GetByID(ctx context.Context, id uint) (*models.SalesInvoice, error)
	Update(ctx context.Context, invoice *models.SalesInvoice) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, filter *dto.SalesInvoiceListRequest) ([]*models.SalesInvoice, int64, error)
	GetByCustomerID(ctx context.Context, customerID uint, offset, limit int) ([]*models.SalesInvoice, int64, error)
	GetByInvoiceNumber(ctx context.Context, invoiceNumber string) (*models.SalesInvoice, error)
	Search(ctx context.Context, query string, offset, limit int) ([]*models.SalesInvoice, int64, error)
	UpdateStatus(ctx context.Context, id uint, status string) error
	GetNextInvoiceNumber(ctx context.Context) (string, error)
}

// SalesInvoiceRepositoryImpl 销售发票仓储实现
type SalesInvoiceRepositoryImpl struct {
	db *gorm.DB
}

// NewSalesInvoiceRepository 创建销售发票仓储实例
func NewSalesInvoiceRepository(db *gorm.DB) SalesInvoiceRepository {
	return &SalesInvoiceRepositoryImpl{db: db}
}

// Create 创建销售发票
func (r *SalesInvoiceRepositoryImpl) Create(ctx context.Context, invoice *models.SalesInvoice) error {
	return r.db.WithContext(ctx).Create(invoice).Error
}

// GetByID 根据ID获取销售发票
func (r *SalesInvoiceRepositoryImpl) GetByID(ctx context.Context, id uint) (*models.SalesInvoice, error) {
	var invoice models.SalesInvoice
	err := r.db.WithContext(ctx).
		Preload("Customer").
		Preload("SalesOrder").
		Preload("DeliveryNote").
		Preload("Items").
		Preload("Items.Item").
		Preload("Payments").
		Preload("StatusLogs").
		First(&invoice, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &invoice, nil
}

// Update 更新销售发票
func (r *SalesInvoiceRepositoryImpl) Update(ctx context.Context, invoice *models.SalesInvoice) error {
	return r.db.WithContext(ctx).Save(invoice).Error
}

// Delete 删除销售发票
func (r *SalesInvoiceRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.SalesInvoice{}, id).Error
}

// List 获取销售发票列表
func (r *SalesInvoiceRepositoryImpl) List(ctx context.Context, filter *dto.SalesInvoiceListRequest) ([]*models.SalesInvoice, int64, error) {
	var invoices []*models.SalesInvoice
	var total int64

	query := r.db.WithContext(ctx).Model(&models.SalesInvoice{})

	// 应用过滤条件
	if filter.CustomerID != nil {
		query = query.Where("customer_id = ?", *filter.CustomerID)
	}
	if filter.DocStatus != "" {
		query = query.Where("doc_status = ?", filter.DocStatus)
	}
	if filter.PaymentStatus != "" {
		query = query.Where("payment_status = ?", filter.PaymentStatus)
	}
	if filter.Currency != "" {
		query = query.Where("currency = ?", filter.Currency)
	}
	if filter.SalesPersonID != nil {
		query = query.Where("sales_person_id = ?", *filter.SalesPersonID)
	}
	if filter.Territory != "" {
		query = query.Where("territory = ?", filter.Territory)
	}
	if filter.DateFrom != "" {
		query = query.Where("invoice_date >= ?", filter.DateFrom)
	}
	if filter.DateTo != "" {
		query = query.Where("invoice_date <= ?", filter.DateTo)
	}
	if filter.Search != "" {
		searchQuery := "%" + filter.Search + "%"
		query = query.Where("invoice_number LIKE ? OR customer_po_number LIKE ?", searchQuery, searchQuery)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 应用分页和排序
	offset := (filter.Page - 1) * filter.PageSize
	query = query.Offset(offset).Limit(filter.PageSize)

	if filter.SortBy != "" {
		order := filter.SortBy
		if filter.SortDesc {
			order += " DESC"
		}
		query = query.Order(order)
	} else {
		query = query.Order("created_at DESC")
	}

	// 获取数据
	err := r.db.WithContext(ctx).Preload("Customer").Find(&invoices).Error
	if err != nil {
		return nil, 0, err
	}

	return invoices, total, nil
}

// GetByCustomerID 根据客户ID获取销售发票列表
func (r *SalesInvoiceRepositoryImpl) GetByCustomerID(ctx context.Context, customerID uint, offset, limit int) ([]*models.SalesInvoice, int64, error) {
	var invoices []*models.SalesInvoice
	var total int64

	// 获取总数
	if err := r.db.WithContext(ctx).Model(&models.SalesInvoice{}).
		Where("customer_id = ?", customerID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := r.db.WithContext(ctx).Preload("Customer").
		Where("customer_id = ?", customerID).
		Offset(offset).Limit(limit).Find(&invoices).Error
	if err != nil {
		return nil, 0, err
	}

	return invoices, total, nil
}

// GetByInvoiceNumber 根据发票号获取销售发票
func (r *SalesInvoiceRepositoryImpl) GetByInvoiceNumber(ctx context.Context, invoiceNumber string) (*models.SalesInvoice, error) {
	var invoice models.SalesInvoice
	err := r.db.WithContext(ctx).Where("invoice_number = ?", invoiceNumber).First(&invoice).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &invoice, nil
}

// Search 搜索销售发票
func (r *SalesInvoiceRepositoryImpl) Search(ctx context.Context, query string, offset, limit int) ([]*models.SalesInvoice, int64, error) {
	var invoices []*models.SalesInvoice
	var total int64

	searchQuery := "%" + query + "%"

	// 获取总数
	if err := r.db.WithContext(ctx).Model(&models.SalesInvoice{}).
		Joins("LEFT JOIN customers ON sales_invoices.customer_id = customers.id").
		Where("sales_invoices.invoice_number LIKE ? OR customers.name LIKE ? OR sales_invoices.customer_po_number LIKE ?",
			searchQuery, searchQuery, searchQuery).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := r.db.WithContext(ctx).Preload("Customer").
		Joins("LEFT JOIN customers ON sales_invoices.customer_id = customers.id").
		Where("sales_invoices.invoice_number LIKE ? OR customers.name LIKE ? OR sales_invoices.customer_po_number LIKE ?",
			searchQuery, searchQuery, searchQuery).
		Offset(offset).Limit(limit).Find(&invoices).Error
	if err != nil {
		return nil, 0, err
	}

	return invoices, total, nil
}

// UpdateStatus 更新销售发票状态
func (r *SalesInvoiceRepositoryImpl) UpdateStatus(ctx context.Context, id uint, status string) error {
	return r.db.WithContext(ctx).Model(&models.SalesInvoice{}).
		Where("id = ?", id).Update("doc_status", status).Error
}

// GetNextInvoiceNumber 获取下一个发票号
func (r *SalesInvoiceRepositoryImpl) GetNextInvoiceNumber(ctx context.Context) (string, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.SalesInvoice{}).Count(&count).Error
	if err != nil {
		return "", err
	}

	// 生成格式: INV-YYYYMM-000001
	now := time.Now()
	prefix := fmt.Sprintf("INV-%04d%02d", now.Year(), now.Month())

	// 查询当月已有的发票数量
	var monthlyCount int64
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Nanosecond)

	err = r.db.WithContext(ctx).Model(&models.SalesInvoice{}).
		Where("created_at >= ? AND created_at <= ?", startOfMonth, endOfMonth).
		Count(&monthlyCount).Error
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s-%06d", prefix, monthlyCount+1), nil
}

// QuotationTemplateRepository 报价单模板仓储接口
type QuotationTemplateRepository interface {
	Create(ctx context.Context, template *models.QuotationTemplate) error
	GetByID(ctx context.Context, id uint) (*models.QuotationTemplate, error)
	Update(ctx context.Context, template *models.QuotationTemplate) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, offset, limit int) ([]*models.QuotationTemplate, int64, error)
	GetActiveTemplates(ctx context.Context) ([]*models.QuotationTemplate, error)
	GetDefaultTemplate(ctx context.Context) (*models.QuotationTemplate, error)
	SetAsDefault(ctx context.Context, id uint) error
}

// QuotationTemplateRepositoryImpl 报价单模板仓储实现
type QuotationTemplateRepositoryImpl struct {
	db *gorm.DB
}

// NewQuotationTemplateRepository 创建报价单模板仓储实例
func NewQuotationTemplateRepository(db *gorm.DB) QuotationTemplateRepository {
	return &QuotationTemplateRepositoryImpl{
		db: db,
	}
}

// Create 创建模板
func (r *QuotationTemplateRepositoryImpl) Create(ctx context.Context, template *models.QuotationTemplate) error {
	return r.db.WithContext(ctx).Create(template).Error
}

// GetByID 根据ID获取模板
func (r *QuotationTemplateRepositoryImpl) GetByID(ctx context.Context, id uint) (*models.QuotationTemplate, error) {
	var template models.QuotationTemplate
	err := r.db.WithContext(ctx).
		Preload("Items").
		Preload("Items.Item").
		First(&template, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &template, nil
}

// Update 更新模板
func (r *QuotationTemplateRepositoryImpl) Update(ctx context.Context, template *models.QuotationTemplate) error {
	return r.db.WithContext(ctx).Save(template).Error
}

// Delete 删除模板
func (r *QuotationTemplateRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.QuotationTemplate{}, id).Error
}

// List 获取模板列表
func (r *QuotationTemplateRepositoryImpl) List(ctx context.Context, offset, limit int) ([]*models.QuotationTemplate, int64, error) {
	var templates []*models.QuotationTemplate
	var total int64

	// 获取总数
	if err := r.db.WithContext(ctx).Model(&models.QuotationTemplate{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取数据
	err := r.db.WithContext(ctx).
		Preload("Items").
		Preload("Items.Item").
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&templates).Error

	return templates, total, err
}

// GetActiveTemplates 获取活跃的模板
func (r *QuotationTemplateRepositoryImpl) GetActiveTemplates(ctx context.Context) ([]*models.QuotationTemplate, error) {
	var templates []*models.QuotationTemplate
	err := r.db.WithContext(ctx).
		Where("is_active = ?", true).
		Preload("Items").
		Preload("Items.Item").
		Order("name ASC").
		Find(&templates).Error
	return templates, err
}

// GetDefaultTemplate 获取默认模板
func (r *QuotationTemplateRepositoryImpl) GetDefaultTemplate(ctx context.Context) (*models.QuotationTemplate, error) {
	var template models.QuotationTemplate
	err := r.db.WithContext(ctx).
		Where("is_default = ? AND is_active = ?", true, true).
		Preload("Items").
		Preload("Items.Item").
		First(&template).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &template, nil
}

// SetAsDefault 设置为默认模板
func (r *QuotationTemplateRepositoryImpl) SetAsDefault(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 先取消所有默认模板
		if err := tx.Model(&models.QuotationTemplate{}).
			Where("is_default = ?", true).
			Update("is_default", false).Error; err != nil {
			return err
		}

		// 设置新的默认模板
		return tx.Model(&models.QuotationTemplate{}).
			Where("id = ?", id).
			Update("is_default", true).Error
	})
}


// QuotationVersionRepository 报价单版本仓储接口
type QuotationVersionRepository interface {
	Create(ctx context.Context, version *models.QuotationVersion) error
	GetByID(ctx context.Context, id uint) (*models.QuotationVersion, error)
	GetByQuotationID(ctx context.Context, quotationID uint) ([]*models.QuotationVersion, error)
	GetActiveVersion(ctx context.Context, quotationID uint) (*models.QuotationVersion, error)
	SetActiveVersion(ctx context.Context, quotationID, versionID uint) error
	GetVersionHistory(ctx context.Context, quotationID uint) ([]*models.QuotationVersion, error)
	GetNextVersionNumber(ctx context.Context, quotationID uint) (int, error)
	Delete(ctx context.Context, id uint) error
}

// QuotationVersionRepositoryImpl 报价单版本仓储实现
type QuotationVersionRepositoryImpl struct {
	db *gorm.DB
}

// NewQuotationVersionRepository 创建报价单版本仓储实例
func NewQuotationVersionRepository(db *gorm.DB) QuotationVersionRepository {
	return &QuotationVersionRepositoryImpl{
		db: db,
	}
}

// Create 创建报价单版本
func (r *QuotationVersionRepositoryImpl) Create(ctx context.Context, version *models.QuotationVersion) error {
	return r.db.WithContext(ctx).Create(version).Error
}

// GetByID 根据ID获取报价单版本
func (r *QuotationVersionRepositoryImpl) GetByID(ctx context.Context, id uint) (*models.QuotationVersion, error) {
	var version models.QuotationVersion
	err := r.db.WithContext(ctx).
		Preload("Quotation").
		First(&version, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &version, nil
}

// GetByQuotationID 根据报价单ID获取所有版本
func (r *QuotationVersionRepositoryImpl) GetByQuotationID(ctx context.Context, quotationID uint) ([]*models.QuotationVersion, error) {
	var versions []*models.QuotationVersion
	err := r.db.WithContext(ctx).
		Where("quotation_id = ?", quotationID).
		Order("version_number DESC").
		Find(&versions).Error
	return versions, err
}

// GetActiveVersion 获取报价单的活跃版本
func (r *QuotationVersionRepositoryImpl) GetActiveVersion(ctx context.Context, quotationID uint) (*models.QuotationVersion, error) {
	var version models.QuotationVersion
	err := r.db.WithContext(ctx).
		Where("quotation_id = ? AND is_active = ?", quotationID, true).
		First(&version).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &version, nil
}

// SetActiveVersion 设置活跃版本
func (r *QuotationVersionRepositoryImpl) SetActiveVersion(ctx context.Context, quotationID, versionID uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 先将所有版本设为非活跃
		if err := tx.Model(&models.QuotationVersion{}).
			Where("quotation_id = ?", quotationID).
			Update("is_active", false).Error; err != nil {
			return err
		}
		
		// 设置指定版本为活跃
		return tx.Model(&models.QuotationVersion{}).
			Where("id = ? AND quotation_id = ?", versionID, quotationID).
			Update("is_active", true).Error
	})
}

// GetVersionHistory 获取版本历史
func (r *QuotationVersionRepositoryImpl) GetVersionHistory(ctx context.Context, quotationID uint) ([]*models.QuotationVersion, error) {
	var versions []*models.QuotationVersion
	err := r.db.WithContext(ctx).
		Where("quotation_id = ?", quotationID).
		Order("version_number ASC").
		Find(&versions).Error
	return versions, err
}

// GetNextVersionNumber 获取下一个版本号
func (r *QuotationVersionRepositoryImpl) GetNextVersionNumber(ctx context.Context, quotationID uint) (int, error) {
	var maxVersion int
	err := r.db.WithContext(ctx).Model(&models.QuotationVersion{}).
		Where("quotation_id = ?", quotationID).
		Select("COALESCE(MAX(version_number), 0)").
		Scan(&maxVersion).Error
	if err != nil {
		return 0, err
	}
	return maxVersion + 1, nil
}

// Delete 删除报价单版本
func (r *QuotationVersionRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.QuotationVersion{}, id).Error
}
