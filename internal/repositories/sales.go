package repositories

import (
	"context"
	"fmt"
	"github.com/galaxyerp/galaxyErp/internal/common"
	"github.com/galaxyerp/galaxyErp/internal/dto"
	"github.com/galaxyerp/galaxyErp/internal/models"
	"gorm.io/gorm"
)

// CustomerRepository 客户仓储接口
type CustomerRepository interface {
	BaseRepository[models.Customer]
	GetByCode(ctx context.Context, code string) (*models.Customer, error)
	GetByEmail(ctx context.Context, email string) (*models.Customer, error)
	Search(ctx context.Context, keyword string, options *common.QueryOptions) ([]*models.Customer, error)
	GetActiveCustomers(ctx context.Context) ([]*models.Customer, error)
	UpdateStatus(ctx context.Context, id uint, isActive bool) error
}

// CustomerRepositoryImpl 客户仓储实现
type CustomerRepositoryImpl struct {
	BaseRepository[models.Customer]
	db *gorm.DB
}

// NewCustomerRepository 创建客户仓储
func NewCustomerRepository(db *gorm.DB) CustomerRepository {
	return &CustomerRepositoryImpl{
		BaseRepository: NewBaseRepository[models.Customer](db),
		db:             db,
	}
}

// GetByCode 根据编码获取客户
func (r *CustomerRepositoryImpl) GetByCode(ctx context.Context, code string) (*models.Customer, error) {
	var customer models.Customer
	err := r.db.WithContext(ctx).Where("code = ?", code).First(&customer).Error
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

// GetByEmail 根据邮箱获取客户
func (r *CustomerRepositoryImpl) GetByEmail(ctx context.Context, email string) (*models.Customer, error) {
	var customer models.Customer
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&customer).Error
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

// Search 搜索客户
func (r *CustomerRepositoryImpl) Search(ctx context.Context, keyword string, options *common.QueryOptions) ([]*models.Customer, error) {
	var customers []*models.Customer
	query := r.db.WithContext(ctx).Model(&models.Customer{})
	
	if keyword != "" {
		query = query.Where("name LIKE ? OR code LIKE ? OR email LIKE ?", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}
	
	if options != nil {
		query = r.buildQuery(options)
	}
	
	err := query.Find(&customers).Error
	return customers, err
}

// GetActiveCustomers 获取活跃客户
func (r *CustomerRepositoryImpl) GetActiveCustomers(ctx context.Context) ([]*models.Customer, error) {
	var customers []*models.Customer
	err := r.db.WithContext(ctx).Where("is_active = ?", true).Find(&customers).Error
	return customers, err
}

// UpdateStatus 更新客户状态
func (r *CustomerRepositoryImpl) UpdateStatus(ctx context.Context, id uint, isActive bool) error {
	return r.db.WithContext(ctx).Model(&models.Customer{}).Where("id = ?", id).Update("is_active", isActive).Error
}

// buildQuery 构建查询
func (r *CustomerRepositoryImpl) buildQuery(options *common.QueryOptions) *gorm.DB {
	query := r.db.Model(&models.Customer{})

	if options == nil {
		return query
	}

	// 应用过滤条件
	for _, filter := range options.Filters {
		query = r.applyFilter(query, filter)
	}

	// 应用排序
	for _, sort := range options.Sorts {
		order := "ASC"
		if sort.Order == common.SortOrderDesc {
			order = "DESC"
		}
		query = query.Order(sort.Field + " " + order)
	}

	// 应用分页
	if options.Pagination != nil {
		query = query.Offset(options.Pagination.GetOffset()).Limit(options.Pagination.GetLimit())
	}

	// 应用关联查询
	for _, include := range options.Includes {
		query = query.Preload(include)
	}

	return query
}

// applyFilter 应用过滤条件
func (r *CustomerRepositoryImpl) applyFilter(query *gorm.DB, filter common.FilterCondition) *gorm.DB {
	switch filter.Operator {
	case common.FilterOperatorEq:
		return query.Where(filter.Field+" = ?", filter.Value)
	case common.FilterOperatorNe:
		return query.Where(filter.Field+" != ?", filter.Value)
	case common.FilterOperatorLike:
		return query.Where(filter.Field+" LIKE ?", "%"+filter.Value.(string)+"%")
	case common.FilterOperatorIn:
		return query.Where(filter.Field+" IN ?", filter.Values)
	case common.FilterOperatorNotIn:
		return query.Where(filter.Field+" NOT IN ?", filter.Values)
	case common.FilterOperatorGt:
		return query.Where(filter.Field+" > ?", filter.Value)
	case common.FilterOperatorGte:
		return query.Where(filter.Field+" >= ?", filter.Value)
	case common.FilterOperatorLt:
		return query.Where(filter.Field+" < ?", filter.Value)
	case common.FilterOperatorLte:
		return query.Where(filter.Field+" <= ?", filter.Value)
	case common.FilterOperatorIsNull:
		return query.Where(filter.Field + " IS NULL")
	case common.FilterOperatorNotNull:
		return query.Where(filter.Field + " IS NOT NULL")
	default:
		return query
	}
}

// GetNextInvoiceNumber 获取下一个发票编号
func (r *SalesInvoiceRepositoryImpl) GetNextInvoiceNumber(ctx context.Context) (string, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.SalesInvoice{}).Count(&count).Error
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("INV-%06d", count+1), nil
}

// AddPaymentWithTransaction 在事务中添加付款记录
func (r *SalesInvoiceRepositoryImpl) AddPaymentWithTransaction(ctx context.Context, invoiceID uint, payment *models.InvoicePayment) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 创建付款记录
		payment.SalesInvoiceID = invoiceID
		if err := tx.Create(payment).Error; err != nil {
			return err
		}

		// 更新发票的已付金额和付款状态
		var invoice models.SalesInvoice
		if err := tx.First(&invoice, invoiceID).Error; err != nil {
			return err
		}

		invoice.PaidAmount += payment.Amount
		invoice.OutstandingAmount = invoice.GrandTotal - invoice.PaidAmount

		// 更新付款状态
		if invoice.OutstandingAmount <= 0 {
			invoice.PaymentStatus = "Paid"
		} else if invoice.PaidAmount > 0 {
			invoice.PaymentStatus = "Partially Paid"
		}

		return tx.Save(&invoice).Error
	})
}

// GetPayments 获取发票的付款记录
func (r *SalesInvoiceRepositoryImpl) GetPayments(ctx context.Context, invoiceID uint) ([]*models.InvoicePayment, error) {
	var payments []*models.InvoicePayment
	err := r.db.WithContext(ctx).Where("sales_invoice_id = ?", invoiceID).Find(&payments).Error
	return payments, err
}

// SalesOrderRepository 销售订单仓储接口
type SalesOrderRepository interface {
	BaseRepository[models.SalesOrder]
	GetByOrderNumber(ctx context.Context, orderNumber string) (*models.SalesOrder, error)
	GetByCustomerID(ctx context.Context, customerID uint) ([]*models.SalesOrder, error)
	GetByStatus(ctx context.Context, status string) ([]*models.SalesOrder, error)
	UpdateStatus(ctx context.Context, id uint, status string) error
	GetStatistics(ctx context.Context, startDate, endDate string) (*dto.SalesStatisticsResponse, error)
	GetSalesTrend(ctx context.Context, period string) ([]*dto.MonthlySales, error)
	Search(ctx context.Context, keyword string, options *common.QueryOptions) ([]*models.SalesOrder, error)
}

// SalesOrderRepositoryImpl 销售订单仓储实现
type SalesOrderRepositoryImpl struct {
	BaseRepository[models.SalesOrder]
	db *gorm.DB
}

// NewSalesOrderRepository 创建销售订单仓储
func NewSalesOrderRepository(db *gorm.DB) SalesOrderRepository {
	return &SalesOrderRepositoryImpl{
		BaseRepository: NewBaseRepository[models.SalesOrder](db),
		db:             db,
	}
}

// GetByOrderNumber 根据订单号获取销售订单
func (r *SalesOrderRepositoryImpl) GetByOrderNumber(ctx context.Context, orderNumber string) (*models.SalesOrder, error) {
	var order models.SalesOrder
	err := r.db.WithContext(ctx).Where("order_number = ?", orderNumber).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// GetByCustomerID 根据客户ID获取销售订单
func (r *SalesOrderRepositoryImpl) GetByCustomerID(ctx context.Context, customerID uint) ([]*models.SalesOrder, error) {
	var orders []*models.SalesOrder
	err := r.db.WithContext(ctx).Where("customer_id = ?", customerID).Find(&orders).Error
	return orders, err
}

// GetByStatus 根据状态获取销售订单
func (r *SalesOrderRepositoryImpl) GetByStatus(ctx context.Context, status string) ([]*models.SalesOrder, error) {
	var orders []*models.SalesOrder
	err := r.db.WithContext(ctx).Where("status = ?", status).Find(&orders).Error
	return orders, err
}

// UpdateStatus 更新销售订单状态
func (r *SalesOrderRepositoryImpl) UpdateStatus(ctx context.Context, id uint, status string) error {
	return r.db.WithContext(ctx).Model(&models.SalesOrder{}).Where("id = ?", id).Update("status", status).Error
}

// GetStatistics 获取销售统计
func (r *SalesOrderRepositoryImpl) GetStatistics(ctx context.Context, startDate, endDate string) (*dto.SalesStatisticsResponse, error) {
	var stats dto.SalesStatisticsResponse
	
	// 获取基本统计信息
	err := r.db.WithContext(ctx).Model(&models.SalesOrder{}).
		Select("COUNT(*) as total_orders, SUM(total_amount) as total_amount").
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Scan(&stats).Error
	
	if err != nil {
		return nil, err
	}
	
	// 获取各状态订单数量
	r.db.WithContext(ctx).Model(&models.SalesOrder{}).
		Where("status = ? AND created_at BETWEEN ? AND ?", "pending", startDate, endDate).
		Count(&stats.PendingOrders)
	
	r.db.WithContext(ctx).Model(&models.SalesOrder{}).
		Where("status = ? AND created_at BETWEEN ? AND ?", "approved", startDate, endDate).
		Count(&stats.ApprovedOrders)
	
	r.db.WithContext(ctx).Model(&models.SalesOrder{}).
		Where("status = ? AND created_at BETWEEN ? AND ?", "completed", startDate, endDate).
		Count(&stats.CompletedOrders)
	
	return &stats, nil
}

// GetSalesTrend 获取销售趋势
func (r *SalesOrderRepositoryImpl) GetSalesTrend(ctx context.Context, period string) ([]*dto.MonthlySales, error) {
	var trends []*dto.MonthlySales
	
	err := r.db.WithContext(ctx).Model(&models.SalesOrder{}).
		Select("DATE_FORMAT(created_at, '%Y-%m') as month, COUNT(*) as order_count, SUM(total_amount) as total_amount").
		Group("DATE_FORMAT(created_at, '%Y-%m')").
		Order("month").
		Scan(&trends).Error
	
	return trends, err
}

// Search 搜索销售订单
func (r *SalesOrderRepositoryImpl) Search(ctx context.Context, keyword string, options *common.QueryOptions) ([]*models.SalesOrder, error) {
	var orders []*models.SalesOrder
	query := r.db.WithContext(ctx).Model(&models.SalesOrder{})
	
	if keyword != "" {
		query = query.Where("order_number LIKE ? OR notes LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	
	if options != nil {
		query = r.buildQueryForSalesOrder(options)
	}
	
	err := query.Find(&orders).Error
	return orders, err
}

// buildQueryForSalesOrder 构建销售订单查询
func (r *SalesOrderRepositoryImpl) buildQueryForSalesOrder(options *common.QueryOptions) *gorm.DB {
	query := r.db.Model(&models.SalesOrder{})

	if options == nil {
		return query
	}

	// 应用过滤条件
	for _, filter := range options.Filters {
		query = r.applySalesOrderFilter(query, filter)
	}

	// 应用排序
	for _, sort := range options.Sorts {
		order := "ASC"
		if sort.Order == common.SortOrderDesc {
			order = "DESC"
		}
		query = query.Order(sort.Field + " " + order)
	}

	// 应用分页
	if options.Pagination != nil {
		query = query.Offset(options.Pagination.GetOffset()).Limit(options.Pagination.GetLimit())
	}

	// 应用关联查询
	for _, include := range options.Includes {
		query = query.Preload(include)
	}

	return query
}

// applySalesOrderFilter 应用销售订单过滤条件
func (r *SalesOrderRepositoryImpl) applySalesOrderFilter(query *gorm.DB, filter common.FilterCondition) *gorm.DB {
	switch filter.Operator {
	case common.FilterOperatorEq:
		return query.Where(filter.Field+" = ?", filter.Value)
	case common.FilterOperatorNe:
		return query.Where(filter.Field+" != ?", filter.Value)
	case common.FilterOperatorLike:
		return query.Where(filter.Field+" LIKE ?", "%"+filter.Value.(string)+"%")
	case common.FilterOperatorIn:
		return query.Where(filter.Field+" IN ?", filter.Values)
	case common.FilterOperatorNotIn:
		return query.Where(filter.Field+" NOT IN ?", filter.Values)
	case common.FilterOperatorGt:
		return query.Where(filter.Field+" > ?", filter.Value)
	case common.FilterOperatorGte:
		return query.Where(filter.Field+" >= ?", filter.Value)
	case common.FilterOperatorLt:
		return query.Where(filter.Field+" < ?", filter.Value)
	case common.FilterOperatorLte:
		return query.Where(filter.Field+" <= ?", filter.Value)
	case common.FilterOperatorIsNull:
		return query.Where(filter.Field + " IS NULL")
	case common.FilterOperatorNotNull:
		return query.Where(filter.Field + " IS NOT NULL")
	default:
		return query
	}
}

// QuotationRepository 报价单仓储接口
type QuotationRepository interface {
	BaseRepository[models.Quotation]
	GetByQuotationNumber(ctx context.Context, quotationNumber string) (*models.Quotation, error)
	GetByCustomerID(ctx context.Context, customerID uint) ([]*models.Quotation, error)
	GetByStatus(ctx context.Context, status string) ([]*models.Quotation, error)
	UpdateStatus(ctx context.Context, id uint, status string) error
	Search(ctx context.Context, keyword string, options *common.QueryOptions) ([]*models.Quotation, error)
}

// QuotationRepositoryImpl 报价单仓储实现
type QuotationRepositoryImpl struct {
	BaseRepository[models.Quotation]
	db *gorm.DB
}

// NewQuotationRepository 创建报价单仓储
func NewQuotationRepository(db *gorm.DB) QuotationRepository {
	return &QuotationRepositoryImpl{
		BaseRepository: NewBaseRepository[models.Quotation](db),
		db:             db,
	}
}

// GetByQuotationNumber 根据报价单号获取报价单
func (r *QuotationRepositoryImpl) GetByQuotationNumber(ctx context.Context, quotationNumber string) (*models.Quotation, error) {
	var quotation models.Quotation
	err := r.db.WithContext(ctx).Where("quotation_number = ?", quotationNumber).First(&quotation).Error
	if err != nil {
		return nil, err
	}
	return &quotation, nil
}

// GetByCustomerID 根据客户ID获取报价单
func (r *QuotationRepositoryImpl) GetByCustomerID(ctx context.Context, customerID uint) ([]*models.Quotation, error) {
	var quotations []*models.Quotation
	err := r.db.WithContext(ctx).Where("customer_id = ?", customerID).Find(&quotations).Error
	return quotations, err
}

// GetByStatus 根据状态获取报价单
func (r *QuotationRepositoryImpl) GetByStatus(ctx context.Context, status string) ([]*models.Quotation, error) {
	var quotations []*models.Quotation
	err := r.db.WithContext(ctx).Where("status = ?", status).Find(&quotations).Error
	return quotations, err
}

// UpdateStatus 更新报价单状态
func (r *QuotationRepositoryImpl) UpdateStatus(ctx context.Context, id uint, status string) error {
	return r.db.WithContext(ctx).Model(&models.Quotation{}).Where("id = ?", id).Update("status", status).Error
}

// Search 搜索报价单
func (r *QuotationRepositoryImpl) Search(ctx context.Context, keyword string, options *common.QueryOptions) ([]*models.Quotation, error) {
	var quotations []*models.Quotation
	query := r.db.WithContext(ctx).Model(&models.Quotation{})
	
	if keyword != "" {
		query = query.Where("quotation_number LIKE ? OR subject LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	
	if options != nil {
		query = r.buildQueryForQuotation(options)
	}
	
	err := query.Find(&quotations).Error
	return quotations, err
}

// buildQueryForQuotation 构建报价单查询
func (r *QuotationRepositoryImpl) buildQueryForQuotation(options *common.QueryOptions) *gorm.DB {
	query := r.db.Model(&models.Quotation{})

	if options == nil {
		return query
	}

	// 应用过滤条件
	for _, filter := range options.Filters {
		query = r.applyQuotationFilter(query, filter)
	}

	// 应用排序
	for _, sort := range options.Sorts {
		order := "ASC"
		if sort.Order == common.SortOrderDesc {
			order = "DESC"
		}
		query = query.Order(sort.Field + " " + order)
	}

	// 应用分页
	if options.Pagination != nil {
		query = query.Offset(options.Pagination.GetOffset()).Limit(options.Pagination.GetLimit())
	}

	// 应用关联查询
	for _, include := range options.Includes {
		query = query.Preload(include)
	}

	return query
}

// applyQuotationFilter 应用报价单过滤条件
func (r *QuotationRepositoryImpl) applyQuotationFilter(query *gorm.DB, filter common.FilterCondition) *gorm.DB {
	switch filter.Operator {
	case common.FilterOperatorEq:
		return query.Where(filter.Field+" = ?", filter.Value)
	case common.FilterOperatorNe:
		return query.Where(filter.Field+" != ?", filter.Value)
	case common.FilterOperatorLike:
		return query.Where(filter.Field+" LIKE ?", "%"+filter.Value.(string)+"%")
	case common.FilterOperatorIn:
		return query.Where(filter.Field+" IN ?", filter.Values)
	case common.FilterOperatorNotIn:
		return query.Where(filter.Field+" NOT IN ?", filter.Values)
	case common.FilterOperatorGt:
		return query.Where(filter.Field+" > ?", filter.Value)
	case common.FilterOperatorGte:
		return query.Where(filter.Field+" >= ?", filter.Value)
	case common.FilterOperatorLt:
		return query.Where(filter.Field+" < ?", filter.Value)
	case common.FilterOperatorLte:
		return query.Where(filter.Field+" <= ?", filter.Value)
	case common.FilterOperatorIsNull:
		return query.Where(filter.Field + " IS NULL")
	case common.FilterOperatorNotNull:
		return query.Where(filter.Field + " IS NOT NULL")
	default:
		return query
	}
}

// SalesInvoiceRepository 销售发票仓储接口
type SalesInvoiceRepository interface {
	BaseRepository[models.SalesInvoice]
	GetByInvoiceNumber(ctx context.Context, invoiceNumber string) (*models.SalesInvoice, error)
	GetByCustomerID(ctx context.Context, customerID uint) ([]*models.SalesInvoice, error)
	GetByStatus(ctx context.Context, status string) ([]*models.SalesInvoice, error)
	UpdateStatus(ctx context.Context, id uint, status string) error
	Search(ctx context.Context, keyword string, options *common.QueryOptions) ([]*models.SalesInvoice, error)
	GetNextInvoiceNumber(ctx context.Context) (string, error)
	AddPaymentWithTransaction(ctx context.Context, invoiceID uint, payment *models.InvoicePayment) error
	GetPayments(ctx context.Context, invoiceID uint) ([]*models.InvoicePayment, error)
}

// SalesInvoiceRepositoryImpl 销售发票仓储实现
type SalesInvoiceRepositoryImpl struct {
	BaseRepository[models.SalesInvoice]
	db *gorm.DB
}

// NewSalesInvoiceRepository 创建销售发票仓储
func NewSalesInvoiceRepository(db *gorm.DB) SalesInvoiceRepository {
	return &SalesInvoiceRepositoryImpl{
		BaseRepository: NewBaseRepository[models.SalesInvoice](db),
		db:             db,
	}
}

// GetByInvoiceNumber 根据发票号获取销售发票
func (r *SalesInvoiceRepositoryImpl) GetByInvoiceNumber(ctx context.Context, invoiceNumber string) (*models.SalesInvoice, error) {
	var invoice models.SalesInvoice
	err := r.db.WithContext(ctx).Where("invoice_number = ?", invoiceNumber).First(&invoice).Error
	if err != nil {
		return nil, err
	}
	return &invoice, nil
}

// GetByCustomerID 根据客户ID获取销售发票
func (r *SalesInvoiceRepositoryImpl) GetByCustomerID(ctx context.Context, customerID uint) ([]*models.SalesInvoice, error) {
	var invoices []*models.SalesInvoice
	err := r.db.WithContext(ctx).Where("customer_id = ?", customerID).Find(&invoices).Error
	return invoices, err
}

// GetByStatus 根据状态获取销售发票
func (r *SalesInvoiceRepositoryImpl) GetByStatus(ctx context.Context, status string) ([]*models.SalesInvoice, error) {
	var invoices []*models.SalesInvoice
	err := r.db.WithContext(ctx).Where("doc_status = ?", status).Find(&invoices).Error
	return invoices, err
}

// UpdateStatus 更新销售发票状态
func (r *SalesInvoiceRepositoryImpl) UpdateStatus(ctx context.Context, id uint, status string) error {
	return r.db.WithContext(ctx).Model(&models.SalesInvoice{}).Where("id = ?", id).Update("doc_status", status).Error
}

// Search 搜索销售发票
func (r *SalesInvoiceRepositoryImpl) Search(ctx context.Context, keyword string, options *common.QueryOptions) ([]*models.SalesInvoice, error) {
	var invoices []*models.SalesInvoice
	query := r.db.WithContext(ctx).Model(&models.SalesInvoice{})
	
	if keyword != "" {
		query = query.Where("invoice_number LIKE ?", "%"+keyword+"%")
	}
	
	if options != nil {
		query = r.buildQueryForSalesInvoice(options)
	}
	
	err := query.Find(&invoices).Error
	return invoices, err
}

// buildQueryForSalesInvoice 构建销售发票查询
func (r *SalesInvoiceRepositoryImpl) buildQueryForSalesInvoice(options *common.QueryOptions) *gorm.DB {
	query := r.db.Model(&models.SalesInvoice{})

	if options == nil {
		return query
	}

	// 应用过滤条件
	for _, filter := range options.Filters {
		query = r.applySalesInvoiceFilter(query, filter)
	}

	// 应用排序
	for _, sort := range options.Sorts {
		order := "ASC"
		if sort.Order == common.SortOrderDesc {
			order = "DESC"
		}
		query = query.Order(sort.Field + " " + order)
	}

	// 应用分页
	if options.Pagination != nil {
		query = query.Offset(options.Pagination.GetOffset()).Limit(options.Pagination.GetLimit())
	}

	// 应用关联查询
	for _, include := range options.Includes {
		query = query.Preload(include)
	}

	return query
}

// applySalesInvoiceFilter 应用销售发票过滤条件
func (r *SalesInvoiceRepositoryImpl) applySalesInvoiceFilter(query *gorm.DB, filter common.FilterCondition) *gorm.DB {
	switch filter.Operator {
	case common.FilterOperatorEq:
		return query.Where(filter.Field+" = ?", filter.Value)
	case common.FilterOperatorNe:
		return query.Where(filter.Field+" != ?", filter.Value)
	case common.FilterOperatorLike:
		return query.Where(filter.Field+" LIKE ?", "%"+filter.Value.(string)+"%")
	case common.FilterOperatorIn:
		return query.Where(filter.Field+" IN ?", filter.Values)
	case common.FilterOperatorNotIn:
		return query.Where(filter.Field+" NOT IN ?", filter.Values)
	case common.FilterOperatorGt:
		return query.Where(filter.Field+" > ?", filter.Value)
	case common.FilterOperatorGte:
		return query.Where(filter.Field+" >= ?", filter.Value)
	case common.FilterOperatorLt:
		return query.Where(filter.Field+" < ?", filter.Value)
	case common.FilterOperatorLte:
		return query.Where(filter.Field+" <= ?", filter.Value)
	case common.FilterOperatorIsNull:
		return query.Where(filter.Field + " IS NULL")
	case common.FilterOperatorNotNull:
		return query.Where(filter.Field + " IS NOT NULL")
	default:
		return query
	}
}

// QuotationTemplateRepository 报价单模板仓储接口
type QuotationTemplateRepository interface {
	BaseRepository[models.QuotationTemplate]
	GetByName(ctx context.Context, name string) (*models.QuotationTemplate, error)
	GetActiveTemplates(ctx context.Context) ([]*models.QuotationTemplate, error)
	GetDefaultTemplate(ctx context.Context) (*models.QuotationTemplate, error)
	SetAsDefault(ctx context.Context, id uint) error
	Search(ctx context.Context, keyword string, options *common.QueryOptions) ([]*models.QuotationTemplate, error)
}

// QuotationTemplateRepositoryImpl 报价单模板仓储实现
type QuotationTemplateRepositoryImpl struct {
	BaseRepository[models.QuotationTemplate]
	db *gorm.DB
}

// NewQuotationTemplateRepository 创建报价单模板仓储
func NewQuotationTemplateRepository(db *gorm.DB) QuotationTemplateRepository {
	return &QuotationTemplateRepositoryImpl{
		BaseRepository: NewBaseRepository[models.QuotationTemplate](db),
		db:             db,
	}
}

// GetByName 根据名称获取报价单模板
func (r *QuotationTemplateRepositoryImpl) GetByName(ctx context.Context, name string) (*models.QuotationTemplate, error) {
	var template models.QuotationTemplate
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&template).Error
	if err != nil {
		return nil, err
	}
	return &template, nil
}

// GetActiveTemplates 获取活跃的报价单模板
func (r *QuotationTemplateRepositoryImpl) GetActiveTemplates(ctx context.Context) ([]*models.QuotationTemplate, error) {
	var templates []*models.QuotationTemplate
	err := r.db.WithContext(ctx).Where("is_active = ?", true).Find(&templates).Error
	return templates, err
}

// GetDefaultTemplate 获取默认报价单模板
func (r *QuotationTemplateRepositoryImpl) GetDefaultTemplate(ctx context.Context) (*models.QuotationTemplate, error) {
	var template models.QuotationTemplate
	err := r.db.WithContext(ctx).Where("is_default = ? AND is_active = ?", true, true).First(&template).Error
	if err != nil {
		return nil, err
	}
	return &template, nil
}

// SetAsDefault 设置为默认模板
func (r *QuotationTemplateRepositoryImpl) SetAsDefault(ctx context.Context, id uint) error {
	// 开启事务
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 先将所有模板的默认状态设为false
		if err := tx.Model(&models.QuotationTemplate{}).Where("is_default = ?", true).Update("is_default", false).Error; err != nil {
			return err
		}
		
		// 将指定模板设为默认
		return tx.Model(&models.QuotationTemplate{}).Where("id = ?", id).Update("is_default", true).Error
	})
}

// Search 搜索报价单模板
func (r *QuotationTemplateRepositoryImpl) Search(ctx context.Context, keyword string, options *common.QueryOptions) ([]*models.QuotationTemplate, error) {
	var templates []*models.QuotationTemplate
	query := r.db.WithContext(ctx).Model(&models.QuotationTemplate{})
	
	if keyword != "" {
		query = query.Where("name LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	
	if options != nil {
		query = r.buildQueryForQuotationTemplate(options)
	}
	
	err := query.Find(&templates).Error
	return templates, err
}

// buildQueryForQuotationTemplate 构建报价单模板查询
func (r *QuotationTemplateRepositoryImpl) buildQueryForQuotationTemplate(options *common.QueryOptions) *gorm.DB {
	query := r.db.Model(&models.QuotationTemplate{})

	if options == nil {
		return query
	}

	// 应用过滤条件
	for _, filter := range options.Filters {
		query = r.applyQuotationTemplateFilter(query, filter)
	}

	// 应用排序
	for _, sort := range options.Sorts {
		order := "ASC"
		if sort.Order == common.SortOrderDesc {
			order = "DESC"
		}
		query = query.Order(sort.Field + " " + order)
	}

	// 应用分页
	if options.Pagination != nil {
		query = query.Offset(options.Pagination.GetOffset()).Limit(options.Pagination.GetLimit())
	}

	// 应用关联查询
	for _, include := range options.Includes {
		query = query.Preload(include)
	}

	return query
}

// applyQuotationTemplateFilter 应用报价单模板过滤条件
func (r *QuotationTemplateRepositoryImpl) applyQuotationTemplateFilter(query *gorm.DB, filter common.FilterCondition) *gorm.DB {
	switch filter.Operator {
	case common.FilterOperatorEq:
		return query.Where(filter.Field+" = ?", filter.Value)
	case common.FilterOperatorNe:
		return query.Where(filter.Field+" != ?", filter.Value)
	case common.FilterOperatorLike:
		return query.Where(filter.Field+" LIKE ?", "%"+filter.Value.(string)+"%")
	case common.FilterOperatorIn:
		return query.Where(filter.Field+" IN ?", filter.Values)
	case common.FilterOperatorNotIn:
		return query.Where(filter.Field+" NOT IN ?", filter.Values)
	case common.FilterOperatorGt:
		return query.Where(filter.Field+" > ?", filter.Value)
	case common.FilterOperatorGte:
		return query.Where(filter.Field+" >= ?", filter.Value)
	case common.FilterOperatorLt:
		return query.Where(filter.Field+" < ?", filter.Value)
	case common.FilterOperatorLte:
		return query.Where(filter.Field+" <= ?", filter.Value)
	case common.FilterOperatorIsNull:
		return query.Where(filter.Field + " IS NULL")
	case common.FilterOperatorNotNull:
		return query.Where(filter.Field + " IS NOT NULL")
	default:
		return query
	}
}

// QuotationVersionRepository 报价单版本仓储接口
type QuotationVersionRepository interface {
	BaseRepository[models.QuotationVersion]
	GetByQuotationID(ctx context.Context, quotationID uint) ([]*models.QuotationVersion, error)
	GetActiveVersion(ctx context.Context, quotationID uint) (*models.QuotationVersion, error)
	GetVersionByNumber(ctx context.Context, quotationID uint, versionNumber int) (*models.QuotationVersion, error)
	GetNextVersionNumber(ctx context.Context, quotationID uint) (int, error)
	SetActiveVersion(ctx context.Context, quotationID uint, versionNumber uint) error
	Search(ctx context.Context, keyword string, options *common.QueryOptions) ([]*models.QuotationVersion, error)
}

// QuotationVersionRepositoryImpl 报价单版本仓储实现
type QuotationVersionRepositoryImpl struct {
	BaseRepository[models.QuotationVersion]
	db *gorm.DB
}

// NewQuotationVersionRepository 创建报价单版本仓储
func NewQuotationVersionRepository(db *gorm.DB) QuotationVersionRepository {
	return &QuotationVersionRepositoryImpl{
		BaseRepository: NewBaseRepository[models.QuotationVersion](db),
		db:             db,
	}
}

// GetByQuotationID 根据报价单ID获取版本列表
func (r *QuotationVersionRepositoryImpl) GetByQuotationID(ctx context.Context, quotationID uint) ([]*models.QuotationVersion, error) {
	var versions []*models.QuotationVersion
	err := r.db.WithContext(ctx).Where("quotation_id = ?", quotationID).Order("version_number DESC").Find(&versions).Error
	return versions, err
}

// GetActiveVersion 获取活跃版本
func (r *QuotationVersionRepositoryImpl) GetActiveVersion(ctx context.Context, quotationID uint) (*models.QuotationVersion, error) {
	var version models.QuotationVersion
	err := r.db.WithContext(ctx).Where("quotation_id = ? AND is_active = ?", quotationID, true).First(&version).Error
	if err != nil {
		return nil, err
	}
	return &version, nil
}

// GetVersionByNumber 根据版本号获取版本
func (r *QuotationVersionRepositoryImpl) GetVersionByNumber(ctx context.Context, quotationID uint, versionNumber int) (*models.QuotationVersion, error) {
	var version models.QuotationVersion
	err := r.db.WithContext(ctx).Where("quotation_id = ? AND version_number = ?", quotationID, versionNumber).First(&version).Error
	if err != nil {
		return nil, err
	}
	return &version, nil
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

// SetActiveVersion 设置活跃版本
func (r *QuotationVersionRepositoryImpl) SetActiveVersion(ctx context.Context, quotationID uint, versionNumber uint) error {
	// 开启事务
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 先将该报价单的所有版本设为非活跃
		if err := tx.Model(&models.QuotationVersion{}).
			Where("quotation_id = ?", quotationID).
			Update("is_active", false).Error; err != nil {
			return err
		}
		
		// 将指定版本设为活跃
		return tx.Model(&models.QuotationVersion{}).
			Where("quotation_id = ? AND version_number = ?", quotationID, versionNumber).
			Update("is_active", true).Error
	})
}

// Search 搜索报价单版本
func (r *QuotationVersionRepositoryImpl) Search(ctx context.Context, keyword string, options *common.QueryOptions) ([]*models.QuotationVersion, error) {
	var versions []*models.QuotationVersion
	query := r.db.WithContext(ctx).Model(&models.QuotationVersion{})
	
	if keyword != "" {
		query = query.Where("version_name LIKE ? OR change_reason LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	
	if options != nil {
		query = r.buildQueryForQuotationVersion(options)
	}
	
	err := query.Find(&versions).Error
	return versions, err
}

// buildQueryForQuotationVersion 构建报价单版本查询
func (r *QuotationVersionRepositoryImpl) buildQueryForQuotationVersion(options *common.QueryOptions) *gorm.DB {
	query := r.db.Model(&models.QuotationVersion{})

	if options == nil {
		return query
	}

	// 应用过滤条件
	for _, filter := range options.Filters {
		query = r.applyQuotationVersionFilter(query, filter)
	}

	// 应用排序
	for _, sort := range options.Sorts {
		order := "ASC"
		if sort.Order == common.SortOrderDesc {
			order = "DESC"
		}
		query = query.Order(sort.Field + " " + order)
	}

	// 应用分页
	if options.Pagination != nil {
		query = query.Offset(options.Pagination.GetOffset()).Limit(options.Pagination.GetLimit())
	}

	// 应用关联查询
	for _, include := range options.Includes {
		query = query.Preload(include)
	}

	return query
}

// applyQuotationVersionFilter 应用报价单版本过滤条件
func (r *QuotationVersionRepositoryImpl) applyQuotationVersionFilter(query *gorm.DB, filter common.FilterCondition) *gorm.DB {
	switch filter.Operator {
	case common.FilterOperatorEq:
		return query.Where(filter.Field+" = ?", filter.Value)
	case common.FilterOperatorNe:
		return query.Where(filter.Field+" != ?", filter.Value)
	case common.FilterOperatorLike:
		return query.Where(filter.Field+" LIKE ?", "%"+filter.Value.(string)+"%")
	case common.FilterOperatorIn:
		return query.Where(filter.Field+" IN ?", filter.Values)
	case common.FilterOperatorNotIn:
		return query.Where(filter.Field+" NOT IN ?", filter.Values)
	case common.FilterOperatorGt:
		return query.Where(filter.Field+" > ?", filter.Value)
	case common.FilterOperatorGte:
		return query.Where(filter.Field+" >= ?", filter.Value)
	case common.FilterOperatorLt:
		return query.Where(filter.Field+" < ?", filter.Value)
	case common.FilterOperatorLte:
		return query.Where(filter.Field+" <= ?", filter.Value)
	case common.FilterOperatorIsNull:
		return query.Where(filter.Field + " IS NULL")
	case common.FilterOperatorNotNull:
		return query.Where(filter.Field + " IS NOT NULL")
	default:
		return query
	}
}
