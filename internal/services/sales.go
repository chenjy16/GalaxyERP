package services

import (
	"context"
	"errors"
	"fmt"
	"time"
	"github.com/galaxyerp/galaxyErp/internal/models"
	"github.com/galaxyerp/galaxyErp/internal/repositories"
	"github.com/galaxyerp/galaxyErp/internal/dto"
)

// CustomerService 客户服务接口
type CustomerService interface {
	CreateCustomer(ctx context.Context, req *dto.CustomerCreateRequest) (*dto.CustomerResponse, error)
	GetCustomer(ctx context.Context, id uint) (*dto.CustomerResponse, error)
	UpdateCustomer(ctx context.Context, id uint, req *dto.CustomerUpdateRequest) error
	DeleteCustomer(ctx context.Context, id uint) error
	ListCustomers(ctx context.Context, req *dto.PaginationRequest) (*dto.PaginatedResponse[dto.CustomerResponse], error)
	SearchCustomers(ctx context.Context, req *dto.SearchRequest) (*dto.PaginatedResponse[dto.CustomerResponse], error)
}

// CustomerServiceImpl 客户服务实现
type CustomerServiceImpl struct {
	customerRepo repositories.CustomerRepository
}

// NewCustomerService 创建客户服务实例
func NewCustomerService(customerRepo repositories.CustomerRepository) CustomerService {
	return &CustomerServiceImpl{
		customerRepo: customerRepo,
	}
}

// CreateCustomer 创建客户
func (s *CustomerServiceImpl) CreateCustomer(ctx context.Context, req *dto.CustomerCreateRequest) (*dto.CustomerResponse, error) {
	// 创建客户
	customer := &models.Customer{
		Email:         req.Email,
		Phone:         req.Phone,
		Address:       req.Address,
		ContactPerson: req.ContactName,
		CreditLimit:   req.CreditLimit,
		// PaymentTerms:  req.PaymentTerms, // 字段不存在
		// Status:        "active", // 字段不存在，使用IsActive
	}
	
	// 设置CodeModel字段
	customer.Name = req.Name
	customer.Code = req.Code
	customer.IsActive = true

	if err := s.customerRepo.Create(ctx, customer); err != nil {
		return nil, fmt.Errorf("创建客户失败: %w", err)
	}

	return s.toCustomerResponse(customer), nil
}

// GetCustomer 获取客户
func (s *CustomerServiceImpl) GetCustomer(ctx context.Context, id uint) (*dto.CustomerResponse, error) {
	customer, err := s.customerRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("获取客户失败: %w", err)
	}
	if customer == nil {
		return nil, errors.New("客户不存在")
	}

	return s.toCustomerResponse(customer), nil
}

// UpdateCustomer 更新客户
func (s *CustomerServiceImpl) UpdateCustomer(ctx context.Context, id uint, req *dto.CustomerUpdateRequest) error {
	customer, err := s.customerRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("获取客户失败: %w", err)
	}

	// 更新字段
	if req.Name != "" {
		customer.Name = req.Name
	}
	// if req.Code != "" {
	// 	customer.Code = req.Code // 字段不存在于UpdateRequest
	// }
	if req.Email != "" {
		customer.Email = req.Email
	}
	if req.Phone != "" {
		customer.Phone = req.Phone
	}
	if req.Address != "" {
		customer.Address = req.Address
	}
	if req.ContactName != "" {
		customer.ContactPerson = req.ContactName
	}
	if req.CreditLimit != nil && *req.CreditLimit != 0 {
		customer.CreditLimit = *req.CreditLimit
	}

	if err := s.customerRepo.Update(ctx, customer); err != nil {
		return fmt.Errorf("更新客户失败: %w", err)
	}

	return nil
}

// DeleteCustomer 删除客户
func (s *CustomerServiceImpl) DeleteCustomer(ctx context.Context, id uint) error {
	customer, err := s.customerRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("获取客户失败: %w", err)
	}
	if customer == nil {
		return errors.New("客户不存在")
	}

	return s.customerRepo.Delete(ctx, id)
}

// ListCustomers 获取客户列表
func (s *CustomerServiceImpl) ListCustomers(ctx context.Context, req *dto.PaginationRequest) (*dto.PaginatedResponse[dto.CustomerResponse], error) {
	offset := req.GetOffset()
	limit := req.GetLimit()

	customers, total, err := s.customerRepo.List(ctx, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("获取客户列表失败: %w", err)
	}

	// 转换为响应格式
	customerResponses := make([]dto.CustomerResponse, len(customers))
	for i, customer := range customers {
		customerResponses[i] = *s.toCustomerResponse(customer)
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return &dto.PaginatedResponse[dto.CustomerResponse]{
		Data:       customerResponses,
		Total:      total,
		Page:       req.Page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

// SearchCustomers 搜索客户
func (s *CustomerServiceImpl) SearchCustomers(ctx context.Context, req *dto.SearchRequest) (*dto.PaginatedResponse[dto.CustomerResponse], error) {
	offset := req.GetOffset()
	limit := req.GetLimit()

	customers, total, err := s.customerRepo.Search(ctx, req.Keyword, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("搜索客户失败: %w", err)
	}

	// 转换为响应格式
	customerResponses := make([]dto.CustomerResponse, len(customers))
	for i, customer := range customers {
		customerResponses[i] = *s.toCustomerResponse(customer)
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return &dto.PaginatedResponse[dto.CustomerResponse]{
		Data:       customerResponses,
		Total:      total,
		Page:       req.Page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

// toCustomerResponse 转换为客户响应格式
func (s *CustomerServiceImpl) toCustomerResponse(customer *models.Customer) *dto.CustomerResponse {
	return &dto.CustomerResponse{
		ID:           customer.ID,
		Name:         customer.Name,
		Code:         customer.Code,
		Email:        customer.Email,
		Phone:        customer.Phone,
		Address:      customer.Address,
		ContactName:  customer.ContactPerson,
		CreditLimit:  customer.CreditLimit,
		// PaymentTerms: customer.PaymentTerms, // 字段不存在
		IsActive:     customer.IsActive,
		CreatedAt:    customer.CreatedAt,
		UpdatedAt:    customer.UpdatedAt,
	}
}

// SalesOrderService 销售订单服务接口
type SalesOrderService interface {
	CreateSalesOrder(ctx context.Context, req *dto.SalesOrderCreateRequest) (*dto.SalesOrderResponse, error)
	GetSalesOrder(ctx context.Context, id uint) (*dto.SalesOrderResponse, error)
	UpdateSalesOrder(ctx context.Context, id uint, req *dto.SalesOrderUpdateRequest) error
	DeleteSalesOrder(ctx context.Context, id uint) error
	ListSalesOrders(ctx context.Context, req *dto.PaginationRequest) (*dto.PaginatedResponse[dto.SalesOrderResponse], error)
	GetSalesOrdersByCustomer(ctx context.Context, customerID uint, req *dto.PaginationRequest) (*dto.PaginatedResponse[dto.SalesOrderResponse], error)
	UpdateOrderStatus(ctx context.Context, id uint, status string) error
}

// SalesOrderServiceImpl 销售订单服务实现
type SalesOrderServiceImpl struct {
	salesOrderRepo repositories.SalesOrderRepository
	customerRepo   repositories.CustomerRepository
}

// NewSalesOrderService 创建销售订单服务实例
func NewSalesOrderService(salesOrderRepo repositories.SalesOrderRepository, customerRepo repositories.CustomerRepository) SalesOrderService {
	return &SalesOrderServiceImpl{
		salesOrderRepo: salesOrderRepo,
		customerRepo:   customerRepo,
	}
}

// CreateSalesOrder 创建销售订单
func (s *SalesOrderServiceImpl) CreateSalesOrder(ctx context.Context, req *dto.SalesOrderCreateRequest) (*dto.SalesOrderResponse, error) {
	// 生成订单编号
	orderNumber := fmt.Sprintf("SO%s%06d", time.Now().Format("20060102"), time.Now().Unix()%1000000)

	var deliveryDate time.Time
	if req.DeliveryDate != nil {
		deliveryDate = *req.DeliveryDate
	} else {
		deliveryDate = time.Now().AddDate(0, 0, 7) // 默认7天后交货
	}

	salesOrder := &models.SalesOrder{
		OrderNumber:  orderNumber,
		CustomerID:   req.CustomerID,
		Date:         time.Now(), // 使用Date字段
		// OrderDate:    req.OrderDate, // 字段不存在
		// RequiredDate: req.RequiredDate, // 字段不存在
		DeliveryDate: deliveryDate,
		Status:       "draft",
		// PaymentTerms:  req.PaymentTerms, // 字段不存在，使用Terms
		// DeliveryTerms: req.DeliveryTerms, // 字段不存在
		Terms:        req.PaymentTerms,
		Notes:        req.Notes,
	}

	if err := s.salesOrderRepo.Create(ctx, salesOrder); err != nil {
		return nil, fmt.Errorf("创建销售订单失败: %w", err)
	}

	return s.toSalesOrderResponse(salesOrder), nil
}

// GetSalesOrder 获取销售订单
func (s *SalesOrderServiceImpl) GetSalesOrder(ctx context.Context, id uint) (*dto.SalesOrderResponse, error) {
	salesOrder, err := s.salesOrderRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("获取销售订单失败: %w", err)
	}
	if salesOrder == nil {
		return nil, errors.New("销售订单不存在")
	}

	return s.toSalesOrderResponse(salesOrder), nil
}

// UpdateSalesOrder 更新销售订单
func (s *SalesOrderServiceImpl) UpdateSalesOrder(ctx context.Context, id uint, req *dto.SalesOrderUpdateRequest) error {
	salesOrder, err := s.salesOrderRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("获取销售订单失败: %w", err)
	}
	if salesOrder == nil {
		return errors.New("销售订单不存在")
	}

	// 更新字段
	if req.CustomerID != nil {
		salesOrder.CustomerID = *req.CustomerID
	}
	if req.DeliveryDate != nil {
		salesOrder.DeliveryDate = *req.DeliveryDate
	}
	if req.Status != nil && *req.Status != "" {
		salesOrder.Status = *req.Status
	}
	if req.Notes != nil && *req.Notes != "" {
		salesOrder.Notes = *req.Notes
	}

	if err := s.salesOrderRepo.Update(ctx, salesOrder); err != nil {
		return fmt.Errorf("更新销售订单失败: %w", err)
	}

	return nil
}

// DeleteSalesOrder 删除销售订单
func (s *SalesOrderServiceImpl) DeleteSalesOrder(ctx context.Context, id uint) error {
	salesOrder, err := s.salesOrderRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("获取销售订单失败: %w", err)
	}
	if salesOrder == nil {
		return errors.New("销售订单不存在")
	}

	return s.salesOrderRepo.Delete(ctx, id)
}

// ListSalesOrders 获取销售订单列表
func (s *SalesOrderServiceImpl) ListSalesOrders(ctx context.Context, req *dto.PaginationRequest) (*dto.PaginatedResponse[dto.SalesOrderResponse], error) {
	offset := req.GetOffset()
	limit := req.GetLimit()

	salesOrders, total, err := s.salesOrderRepo.List(ctx, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("获取销售订单列表失败: %w", err)
	}

	// 转换为响应格式
	salesOrderResponses := make([]dto.SalesOrderResponse, len(salesOrders))
	for i, salesOrder := range salesOrders {
		salesOrderResponses[i] = *s.toSalesOrderResponse(salesOrder)
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return &dto.PaginatedResponse[dto.SalesOrderResponse]{
		Data:       salesOrderResponses,
		Total:      total,
		Page:       req.Page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

// GetSalesOrdersByCustomer 根据客户获取销售订单
func (s *SalesOrderServiceImpl) GetSalesOrdersByCustomer(ctx context.Context, customerID uint, req *dto.PaginationRequest) (*dto.PaginatedResponse[dto.SalesOrderResponse], error) {
	offset := req.GetOffset()
	limit := req.GetLimit()

	salesOrders, total, err := s.salesOrderRepo.GetByCustomerID(ctx, customerID, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("获取客户销售订单失败: %w", err)
	}

	// 转换为响应格式
	salesOrderResponses := make([]dto.SalesOrderResponse, len(salesOrders))
	for i, salesOrder := range salesOrders {
		salesOrderResponses[i] = *s.toSalesOrderResponse(salesOrder)
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return &dto.PaginatedResponse[dto.SalesOrderResponse]{
		Data:       salesOrderResponses,
		Total:      total,
		Page:       req.Page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

// UpdateOrderStatus 更新订单状态
func (s *SalesOrderServiceImpl) UpdateOrderStatus(ctx context.Context, id uint, status string) error {
	salesOrder, err := s.salesOrderRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("获取销售订单失败: %w", err)
	}
	if salesOrder == nil {
		return errors.New("销售订单不存在")
	}

	salesOrder.Status = status
	return s.salesOrderRepo.Update(ctx, salesOrder)
}

// toSalesOrderResponse 转换为销售订单响应格式
func (s *SalesOrderServiceImpl) toSalesOrderResponse(salesOrder *models.SalesOrder) *dto.SalesOrderResponse {
	return &dto.SalesOrderResponse{
		ID:              salesOrder.ID,
		Number:          salesOrder.OrderNumber,
		Status:          salesOrder.Status,
		ExpectedDate:    salesOrder.DeliveryDate, // 使用DeliveryDate作为ExpectedDate
		PaymentTerms:    salesOrder.Terms, // 使用Terms字段
		ShippingAddress: "", // 模型中没有ShippingAddress字段
		Notes:           salesOrder.Notes,
		SubTotal:        salesOrder.TotalAmount,
		DiscountAmount:  salesOrder.DiscountAmount,
		TaxAmount:       salesOrder.TaxAmount,
		TotalAmount:     salesOrder.GrandTotal,
		CreatedAt:       salesOrder.CreatedAt,
		UpdatedAt:       salesOrder.UpdatedAt,
	}
}

// QuotationService 报价服务接口
type QuotationService interface {
	CreateQuotation(ctx context.Context, req *dto.QuotationCreateRequest) (*dto.QuotationResponse, error)
	GetQuotation(ctx context.Context, id uint) (*dto.QuotationResponse, error)
	UpdateQuotation(ctx context.Context, id uint, req *dto.QuotationUpdateRequest) error
	DeleteQuotation(ctx context.Context, id uint) error
	ListQuotations(ctx context.Context, req *dto.PaginationRequest) (*dto.PaginatedResponse[dto.QuotationResponse], error)
	GetQuotationsByCustomer(ctx context.Context, customerID uint, req *dto.PaginationRequest) (*dto.PaginatedResponse[dto.QuotationResponse], error)
	SearchQuotations(ctx context.Context, req *dto.SearchRequest) (*dto.PaginatedResponse[dto.QuotationResponse], error)
	ConvertToSalesOrder(ctx context.Context, quotationID uint) (*dto.SalesOrderResponse, error)
}

// QuotationServiceImpl 报价服务实现
type QuotationServiceImpl struct {
	quotationRepo repositories.QuotationRepository
	customerRepo  repositories.CustomerRepository
}

// NewQuotationService 创建报价服务实例
func NewQuotationService(quotationRepo repositories.QuotationRepository, customerRepo repositories.CustomerRepository) QuotationService {
	return &QuotationServiceImpl{
		quotationRepo: quotationRepo,
		customerRepo:  customerRepo,
	}
}

// CreateQuotation 创建报价单
func (s *QuotationServiceImpl) CreateQuotation(ctx context.Context, req *dto.QuotationCreateRequest) (*dto.QuotationResponse, error) {
	// 生成报价单编号
	quotationNumber := fmt.Sprintf("QT%s%06d", time.Now().Format("20060102"), time.Now().Unix()%1000000)

	quotation := &models.Quotation{
		QuotationNumber: quotationNumber,
		CustomerID:      req.CustomerID,
		Date:            time.Now(), // 使用Date字段
		ValidTill:       req.ValidUntil, // 使用ValidUntil字段
		Status:          "draft",
		Subject:         req.Title, // 使用Title作为Subject
		Terms:           req.PaymentTerms, // 使用PaymentTerms作为Terms
		Notes:           req.Notes,
	}

	if err := s.quotationRepo.Create(ctx, quotation); err != nil {
		return nil, fmt.Errorf("创建报价单失败: %w", err)
	}

	return s.toQuotationResponse(quotation), nil
}

// GetQuotation 获取报价
func (s *QuotationServiceImpl) GetQuotation(ctx context.Context, id uint) (*dto.QuotationResponse, error) {
	quotation, err := s.quotationRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("获取报价失败: %w", err)
	}
	if quotation == nil {
		return nil, errors.New("报价不存在")
	}

	return s.toQuotationResponse(quotation), nil
}

// UpdateQuotation 更新报价
func (s *QuotationServiceImpl) UpdateQuotation(ctx context.Context, id uint, req *dto.QuotationUpdateRequest) error {
	quotation, err := s.quotationRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("获取报价单失败: %w", err)
	}

	// 更新字段
	if req.ValidUntil != nil {
		quotation.ValidTill = *req.ValidUntil // 使用ValidTill字段
	}
	if req.Title != "" {
		quotation.Subject = req.Title // 使用Title作为Subject
	}
	if req.PaymentTerms != "" {
		quotation.Terms = req.PaymentTerms // 使用PaymentTerms作为Terms
	}
	if req.Notes != "" {
		quotation.Notes = req.Notes
	}

	if err := s.quotationRepo.Update(ctx, quotation); err != nil {
		return fmt.Errorf("更新报价单失败: %w", err)
	}

	return nil
}

// DeleteQuotation 删除报价
func (s *QuotationServiceImpl) DeleteQuotation(ctx context.Context, id uint) error {
	quotation, err := s.quotationRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("获取报价失败: %w", err)
	}
	if quotation == nil {
		return errors.New("报价不存在")
	}

	return s.quotationRepo.Delete(ctx, id)
}

// ListQuotations 获取报价列表
func (s *QuotationServiceImpl) ListQuotations(ctx context.Context, req *dto.PaginationRequest) (*dto.PaginatedResponse[dto.QuotationResponse], error) {
	offset := req.GetOffset()
	limit := req.GetLimit()

	quotations, total, err := s.quotationRepo.List(ctx, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("获取报价列表失败: %w", err)
	}

	// 转换为响应格式
	quotationResponses := make([]dto.QuotationResponse, len(quotations))
	for i, quotation := range quotations {
		quotationResponses[i] = *s.toQuotationResponse(quotation)
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return &dto.PaginatedResponse[dto.QuotationResponse]{
		Data:       quotationResponses,
		Total:      total,
		Page:       req.Page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

// GetQuotationsByCustomer 根据客户获取报价
func (s *QuotationServiceImpl) GetQuotationsByCustomer(ctx context.Context, customerID uint, req *dto.PaginationRequest) (*dto.PaginatedResponse[dto.QuotationResponse], error) {
	offset := req.GetOffset()
	limit := req.GetLimit()

	quotations, total, err := s.quotationRepo.GetByCustomerID(ctx, customerID, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("获取客户报价失败: %w", err)
	}

	// 转换为响应格式
	quotationResponses := make([]dto.QuotationResponse, len(quotations))
	for i, quotation := range quotations {
		quotationResponses[i] = *s.toQuotationResponse(quotation)
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return &dto.PaginatedResponse[dto.QuotationResponse]{
		Data:       quotationResponses,
		Total:      total,
		Page:       req.Page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

// SearchQuotations 搜索报价
func (s *QuotationServiceImpl) SearchQuotations(ctx context.Context, req *dto.SearchRequest) (*dto.PaginatedResponse[dto.QuotationResponse], error) {
	offset := req.GetOffset()
	limit := req.GetLimit()

	quotations, total, err := s.quotationRepo.Search(ctx, req.Keyword, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("搜索报价失败: %w", err)
	}

	// 转换为响应格式
	quotationResponses := make([]dto.QuotationResponse, len(quotations))
	for i, quotation := range quotations {
		quotationResponses[i] = *s.toQuotationResponse(quotation)
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return &dto.PaginatedResponse[dto.QuotationResponse]{
		Data:       quotationResponses,
		Total:      total,
		Page:       req.Page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

// ConvertToSalesOrder 将报价转换为销售订单
func (s *QuotationServiceImpl) ConvertToSalesOrder(ctx context.Context, quotationID uint) (*dto.SalesOrderResponse, error) {
	// 这个方法需要与SalesOrderService协作实现
	// 暂时返回错误，提示需要实现
	return nil, errors.New("报价转销售订单功能暂未实现")
}

// toQuotationResponse 转换为报价响应格式
func (s *QuotationServiceImpl) toQuotationResponse(quotation *models.Quotation) *dto.QuotationResponse {
	return &dto.QuotationResponse{
		ID:              quotation.ID,
		Number:          quotation.QuotationNumber,
		Title:           quotation.Subject, // 使用Subject作为Title
		Description:     "", // 模型中没有Description字段
		Status:          quotation.Status,
		ValidUntil:      quotation.ValidTill, // 使用ValidTill字段
		PaymentTerms:    quotation.Terms, // 使用Terms字段
		Notes:           quotation.Notes,
		SubTotal:        quotation.TotalAmount,
		DiscountAmount:  quotation.DiscountAmount,
		TaxAmount:       quotation.TaxAmount,
		TotalAmount:     quotation.GrandTotal,
		CreatedAt:       quotation.CreatedAt,
		UpdatedAt:       quotation.UpdatedAt,
	}
}

func (s *QuotationServiceImpl) toQuotationItemResponse(item models.QuotationItem) dto.QuotationItemResponse {
	return dto.QuotationItemResponse{
		ID:        item.ID,
		Quantity:  item.Quantity,
		UnitPrice: item.Rate, // 使用Rate字段作为UnitPrice
		Amount:    item.Amount,
		// Notes字段在QuotationItem模型中不存在，使用Description
	}
}