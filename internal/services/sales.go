package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/galaxyerp/galaxyErp/internal/dto"
	"github.com/galaxyerp/galaxyErp/internal/models"
	"github.com/galaxyerp/galaxyErp/internal/repositories"
	"time"
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
		ID:          customer.ID,
		Name:        customer.Name,
		Code:        customer.Code,
		Email:       customer.Email,
		Phone:       customer.Phone,
		Address:     customer.Address,
		ContactName: customer.ContactPerson,
		CreditLimit: customer.CreditLimit,
		// PaymentTerms: customer.PaymentTerms, // 字段不存在
		IsActive:  customer.IsActive,
		CreatedAt: customer.CreatedAt,
		UpdatedAt: customer.UpdatedAt,
	}
}

// SalesOrderService 销售订单服务接口
type SalesOrderService interface {
	CreateSalesOrder(ctx context.Context, req *dto.SalesOrderCreateRequest, userID uint) (*dto.SalesOrderResponse, error)
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
func (s *SalesOrderServiceImpl) CreateSalesOrder(ctx context.Context, req *dto.SalesOrderCreateRequest, userID uint) (*dto.SalesOrderResponse, error) {
	// 生成订单编号
	orderNumber := fmt.Sprintf("SO%s%06d", time.Now().Format("20060102"), time.Now().Unix()%1000000)

	var deliveryDate time.Time
	if req.DeliveryDate != nil {
		deliveryDate = *req.DeliveryDate
	} else {
		deliveryDate = req.ExpectedDate // 使用期望日期作为交货日期
	}

	// 计算订单总金额
	var totalAmount, discountAmount, taxAmount, grandTotal float64
	var orderItems []models.SalesOrderItem

	for _, itemReq := range req.Items {
		// 计算行金额
		lineAmount := itemReq.Quantity * itemReq.UnitPrice
		lineDiscountAmount := lineAmount * (itemReq.Discount / 100)
		lineNetAmount := lineAmount - lineDiscountAmount
		lineTaxAmount := lineNetAmount * (itemReq.TaxRate / 100)
		lineTotalAmount := lineNetAmount + lineTaxAmount

		orderItem := models.SalesOrderItem{
			ItemID:         itemReq.ItemID,
			Quantity:       itemReq.Quantity,
			Rate:           itemReq.UnitPrice,
			Amount:         lineAmount,
			DiscountRate:   itemReq.Discount,
			DiscountAmount: lineDiscountAmount,
			TaxRate:        itemReq.TaxRate,
			TaxAmount:      lineTaxAmount,
			TotalAmount:    lineTotalAmount,
		}

		orderItems = append(orderItems, orderItem)
		totalAmount += lineAmount
		discountAmount += lineDiscountAmount
		taxAmount += lineTaxAmount
		grandTotal += lineTotalAmount
	}

	salesOrder := &models.SalesOrder{
		OrderNumber:    orderNumber,
		CustomerID:     req.CustomerID,
		Date:           req.OrderDate,
		DeliveryDate:   deliveryDate,
		Status:         "draft",
		QuotationID:    req.QuotationID,
		TotalAmount:    totalAmount,
		DiscountAmount: discountAmount,
		TaxAmount:      taxAmount,
		GrandTotal:     grandTotal,
		Terms:          req.PaymentTerms,
		Notes:          req.Notes,
		CreatedBy:      userID,
	}

	// 创建销售订单
	if err := s.salesOrderRepo.Create(ctx, salesOrder); err != nil {
		return nil, fmt.Errorf("创建销售订单失败: %w", err)
	}

	// 创建订单项目
	for i := range orderItems {
		orderItems[i].SalesOrderID = salesOrder.ID
		if err := s.salesOrderRepo.CreateItem(ctx, &orderItems[i]); err != nil {
			return nil, fmt.Errorf("创建销售订单项目失败: %w", err)
		}
	}

	// 重新从数据库加载订单以获取完整的预加载信息
	fullSalesOrder, err := s.salesOrderRepo.GetByID(ctx, salesOrder.ID)
	if err != nil {
		return nil, fmt.Errorf("重新加载销售订单失败: %w", err)
	}

	return s.toSalesOrderResponse(fullSalesOrder), nil
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
	if req.OrderDate != nil {
		salesOrder.Date = *req.OrderDate
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
	if req.TotalAmount != nil {
		salesOrder.TotalAmount = *req.TotalAmount
		// 如果没有折扣和税费，GrandTotal等于TotalAmount
		if salesOrder.DiscountAmount == 0 && salesOrder.TaxAmount == 0 {
			salesOrder.GrandTotal = *req.TotalAmount
		}
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
	response := &dto.SalesOrderResponse{
		ID:              salesOrder.ID,
		Number:          salesOrder.OrderNumber,
		OrderNumber:     salesOrder.OrderNumber,  // 前端期望的字段名
		Status:          salesOrder.Status,
		OrderDate:       salesOrder.Date,         // 前端期望的订单日期字段
		DeliveryDate:    salesOrder.DeliveryDate, // 前端期望的交付日期字段
		ExpectedDate:    salesOrder.DeliveryDate, // 使用DeliveryDate作为ExpectedDate
		PaymentTerms:    salesOrder.Terms,        // 使用Terms字段
		ShippingAddress: "",                      // 模型中没有ShippingAddress字段
		Notes:           salesOrder.Notes,
		SubTotal:        salesOrder.TotalAmount,
		DiscountAmount:  salesOrder.DiscountAmount,
		TaxAmount:       salesOrder.TaxAmount,
		TotalAmount:     salesOrder.GrandTotal,
		CreatedAt:       salesOrder.CreatedAt,
		UpdatedAt:       salesOrder.UpdatedAt,
	}

	// 处理客户信息 - 添加空值检查
	if salesOrder.Customer.ID != 0 {
		response.Customer = dto.CustomerResponse{
			ID:          salesOrder.Customer.ID,
			Name:        salesOrder.Customer.Name,
			Code:        salesOrder.Customer.Code,
			ContactName: salesOrder.Customer.ContactPerson,
			Email:       salesOrder.Customer.Email,
			Phone:       salesOrder.Customer.Phone,
			Address:     salesOrder.Customer.Address,
			CreditLimit: salesOrder.Customer.CreditLimit,
			IsActive:    salesOrder.Customer.IsActive,
			CreatedAt:   salesOrder.Customer.CreatedAt,
			UpdatedAt:   salesOrder.Customer.UpdatedAt,
		}
	} else {
		// 如果客户信息没有预加载，创建一个空的客户响应
		response.Customer = dto.CustomerResponse{}
	}

	// 处理创建者信息 - 使用预加载的CreatedByUser
	if salesOrder.CreatedByUser != nil && salesOrder.CreatedByUser.ID != 0 {
		response.CreatedBy = dto.UserResponse{
			ID:        salesOrder.CreatedByUser.ID,
			Username:  salesOrder.CreatedByUser.Username,
			Email:     salesOrder.CreatedByUser.Email,
			FirstName: salesOrder.CreatedByUser.FirstName,
			LastName:  salesOrder.CreatedByUser.LastName,
			IsActive:  salesOrder.CreatedByUser.IsActive,
			CreatedAt: salesOrder.CreatedByUser.CreatedAt,
			UpdatedAt: salesOrder.CreatedByUser.UpdatedAt,
		}
	} else {
		// 如果创建者信息没有预加载，创建一个空的用户响应
		response.CreatedBy = dto.UserResponse{}
	}

	// 处理订单项目
	if len(salesOrder.Items) > 0 {
		response.Items = make([]dto.SalesOrderItemResponse, len(salesOrder.Items))
		for i, item := range salesOrder.Items {
			itemResponse := dto.SalesOrderItemResponse{
				ID:             item.ID,
				SalesOrderID:   item.SalesOrderID,
				ItemID:         item.ItemID,
				Quantity:       item.Quantity,
				UnitPrice:      item.Rate,           // 使用Rate字段作为UnitPrice
				Discount:       item.DiscountRate,   // 使用DiscountRate字段作为Discount
				DiscountAmount: item.DiscountAmount,
				TaxRate:        item.TaxRate,
				TaxAmount:      item.TaxAmount,
				LineTotal:      item.TotalAmount,    // 使用TotalAmount字段作为LineTotal
				DeliveredQty:   item.DeliveredQty,
				Description:    item.Description,
				CreatedAt:      item.CreatedAt,
				UpdatedAt:      item.UpdatedAt,
			}

			// 处理物料信息 - 添加空值检查
			if item.Item.ID != 0 {
				itemResponse.Item = dto.ItemResponse{
					ID:          item.Item.ID,
					Name:        item.Item.Name,
					Code:        item.Item.Code,
					Description: item.Item.Description,
					UnitCost:    item.Item.Cost,
					SalePrice:   item.Item.Price,
					IsActive:    item.Item.IsActive,
					CreatedAt:   item.Item.CreatedAt,
					UpdatedAt:   item.Item.UpdatedAt,
				}
			} else {
				// 如果物料信息没有预加载，创建一个空的物料响应
				itemResponse.Item = dto.ItemResponse{}
			}

			response.Items[i] = itemResponse
		}
	}

	return response
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
		Date:            time.Now(),     // 使用Date字段
		ValidTill:       req.ValidUntil, // 使用ValidUntil字段
		Status:          "draft",
		Subject:         req.Title,        // 使用Title作为Subject
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
	if req.CustomerID != nil {
		quotation.CustomerID = *req.CustomerID
	}
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
	if req.TotalAmount != nil {
		quotation.TotalAmount = *req.TotalAmount
		// 如果没有折扣和税费，将GrandTotal设置为TotalAmount
		if quotation.DiscountAmount == 0 && quotation.TaxAmount == 0 {
			quotation.GrandTotal = *req.TotalAmount
		}
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
	// 转换客户信息
	var customerResponse dto.CustomerResponse
	if quotation.Customer.ID != 0 {
		customerResponse = dto.CustomerResponse{
			ID:           quotation.Customer.ID,
			Name:         quotation.Customer.Name,
			Code:         quotation.Customer.Code,
			ContactName:  quotation.Customer.ContactPerson,
			Phone:        quotation.Customer.Phone,
			Email:        quotation.Customer.Email,
			Address:      quotation.Customer.Address,
			CreditLimit:  quotation.Customer.CreditLimit,
			IsActive:     quotation.Customer.IsActive,
			CreatedAt:    quotation.Customer.CreatedAt,
			UpdatedAt:    quotation.Customer.UpdatedAt,
		}
	}

	return &dto.QuotationResponse{
		ID:               quotation.ID,
		Number:           quotation.QuotationNumber,
		QuotationNumber:  quotation.QuotationNumber,  // 前端期望的字段名
		Title:            quotation.Subject,          // 使用Subject作为Title
		Subject:          quotation.Subject,          // 前端期望的字段名
		Description:      "",                         // 模型中没有Description字段
		Status:           quotation.Status,
		ValidUntil:       quotation.ValidTill,        // 使用ValidTill字段
		ValidTill:        quotation.ValidTill,        // 前端期望的字段名
		Date:             quotation.CreatedAt,        // 前端期望的字段名
		PaymentTerms:     quotation.Terms,            // 使用Terms字段
		Notes:            quotation.Notes,
		SubTotal:         quotation.TotalAmount,
		DiscountAmount:   quotation.DiscountAmount,
		TaxAmount:        quotation.TaxAmount,
		TotalAmount:      quotation.TotalAmount,      // 前端表单期望的字段名，应该是基础总金额
		GrandTotal:       quotation.GrandTotal,       // 前端表格期望的字段名，是最终总计金额
		Customer:         customerResponse,           // 包含客户信息
		CreatedAt:        quotation.CreatedAt,
		UpdatedAt:        quotation.UpdatedAt,
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

// QuotationTemplateService 报价单模板服务接口
type QuotationTemplateService interface {
	CreateTemplate(ctx context.Context, req *dto.QuotationTemplateCreateRequest) (*dto.QuotationTemplateResponse, error)
	GetTemplate(ctx context.Context, id uint) (*dto.QuotationTemplateResponse, error)
	UpdateTemplate(ctx context.Context, id uint, req *dto.QuotationTemplateUpdateRequest) error
	DeleteTemplate(ctx context.Context, id uint) error
	ListTemplates(ctx context.Context, req *dto.PaginationRequest) (*dto.PaginatedResponse[dto.QuotationTemplateResponse], error)
	GetActiveTemplates(ctx context.Context) ([]*dto.QuotationTemplateResponse, error)
	GetDefaultTemplate(ctx context.Context) (*dto.QuotationTemplateResponse, error)
	SetAsDefault(ctx context.Context, id uint) error
	CreateQuotationFromTemplate(ctx context.Context, templateID uint, customerID uint) (*dto.QuotationResponse, error)
}

// QuotationTemplateServiceImpl 报价单模板服务实现
type QuotationTemplateServiceImpl struct {
	templateRepo  repositories.QuotationTemplateRepository
	quotationRepo repositories.QuotationRepository
}

// NewQuotationTemplateService 创建报价单模板服务实例
func NewQuotationTemplateService(templateRepo repositories.QuotationTemplateRepository, quotationRepo repositories.QuotationRepository) QuotationTemplateService {
	return &QuotationTemplateServiceImpl{
		templateRepo:  templateRepo,
		quotationRepo: quotationRepo,
	}
}

// CreateTemplate 创建模板
func (s *QuotationTemplateServiceImpl) CreateTemplate(ctx context.Context, req *dto.QuotationTemplateCreateRequest) (*dto.QuotationTemplateResponse, error) {
	template := &models.QuotationTemplate{
		Name:         req.Name,
		Description:  req.Description,
		IsDefault:    req.IsDefault,
		IsActive:     true,
		ValidityDays: req.ValidityDays,
		Terms:        req.Terms,
		Notes:        req.Notes,
		DiscountRate: req.DiscountRate,
		TaxRate:      req.TaxRate,
		CreatedBy:    req.CreatedBy,
	}

	// 如果设置为默认模板，需要先取消其他默认模板
	if req.IsDefault {
		if err := s.templateRepo.SetAsDefault(ctx, 0); err != nil {
			return nil, fmt.Errorf("设置默认模板失败: %w", err)
		}
	}

	if err := s.templateRepo.Create(ctx, template); err != nil {
		return nil, fmt.Errorf("创建模板失败: %w", err)
	}

	return s.toTemplateResponse(template), nil
}

// GetTemplate 获取模板
func (s *QuotationTemplateServiceImpl) GetTemplate(ctx context.Context, id uint) (*dto.QuotationTemplateResponse, error) {
	template, err := s.templateRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("获取模板失败: %w", err)
	}
	if template == nil {
		return nil, errors.New("模板不存在")
	}

	return s.toTemplateResponse(template), nil
}

// UpdateTemplate 更新模板
func (s *QuotationTemplateServiceImpl) UpdateTemplate(ctx context.Context, id uint, req *dto.QuotationTemplateUpdateRequest) error {
	template, err := s.templateRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("获取模板失败: %w", err)
	}
	if template == nil {
		return errors.New("模板不存在")
	}

	// 更新字段
	if req.Name != "" {
		template.Name = req.Name
	}
	if req.Description != "" {
		template.Description = req.Description
	}
	if req.ValidityDays > 0 {
		template.ValidityDays = req.ValidityDays
	}
	if req.Terms != "" {
		template.Terms = req.Terms
	}
	if req.Notes != "" {
		template.Notes = req.Notes
	}
	template.DiscountRate = req.DiscountRate
	template.TaxRate = req.TaxRate
	template.IsActive = req.IsActive

	// 如果设置为默认模板，需要先取消其他默认模板
	if req.IsDefault && !template.IsDefault {
		if err := s.templateRepo.SetAsDefault(ctx, id); err != nil {
			return fmt.Errorf("设置默认模板失败: %w", err)
		}
		template.IsDefault = true
	}

	return s.templateRepo.Update(ctx, template)
}

// DeleteTemplate 删除模板
func (s *QuotationTemplateServiceImpl) DeleteTemplate(ctx context.Context, id uint) error {
	template, err := s.templateRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("获取模板失败: %w", err)
	}
	if template == nil {
		return errors.New("模板不存在")
	}

	return s.templateRepo.Delete(ctx, id)
}

// ListTemplates 获取模板列表
func (s *QuotationTemplateServiceImpl) ListTemplates(ctx context.Context, req *dto.PaginationRequest) (*dto.PaginatedResponse[dto.QuotationTemplateResponse], error) {
	offset := req.GetOffset()
	limit := req.GetLimit()

	templates, total, err := s.templateRepo.List(ctx, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("获取模板列表失败: %w", err)
	}

	// 转换为响应格式
	templateResponses := make([]dto.QuotationTemplateResponse, len(templates))
	for i, template := range templates {
		templateResponses[i] = *s.toTemplateResponse(template)
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return &dto.PaginatedResponse[dto.QuotationTemplateResponse]{
		Data:       templateResponses,
		Total:      total,
		Page:       req.Page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

// GetActiveTemplates 获取活跃的模板
func (s *QuotationTemplateServiceImpl) GetActiveTemplates(ctx context.Context) ([]*dto.QuotationTemplateResponse, error) {
	templates, err := s.templateRepo.GetActiveTemplates(ctx)
	if err != nil {
		return nil, fmt.Errorf("获取活跃模板失败: %w", err)
	}

	responses := make([]*dto.QuotationTemplateResponse, len(templates))
	for i, template := range templates {
		responses[i] = s.toTemplateResponse(template)
	}

	return responses, nil
}

// GetDefaultTemplate 获取默认模板
func (s *QuotationTemplateServiceImpl) GetDefaultTemplate(ctx context.Context) (*dto.QuotationTemplateResponse, error) {
	template, err := s.templateRepo.GetDefaultTemplate(ctx)
	if err != nil {
		return nil, fmt.Errorf("获取默认模板失败: %w", err)
	}
	if template == nil {
		return nil, nil
	}

	return s.toTemplateResponse(template), nil
}

// SetAsDefault 设置为默认模板
func (s *QuotationTemplateServiceImpl) SetAsDefault(ctx context.Context, id uint) error {
	return s.templateRepo.SetAsDefault(ctx, id)
}

// CreateQuotationFromTemplate 从模板创建报价单
func (s *QuotationTemplateServiceImpl) CreateQuotationFromTemplate(ctx context.Context, templateID uint, customerID uint) (*dto.QuotationResponse, error) {
	template, err := s.templateRepo.GetByID(ctx, templateID)
	if err != nil {
		return nil, fmt.Errorf("获取模板失败: %w", err)
	}
	if template == nil {
		return nil, errors.New("模板不存在")
	}

	// 生成报价单编号
	quotationNumber := fmt.Sprintf("QT%s%06d", time.Now().Format("20060102"), time.Now().Unix()%1000000)

	// 计算有效期
	validTill := time.Now().AddDate(0, 0, template.ValidityDays)

	quotation := &models.Quotation{
		QuotationNumber: quotationNumber,
		CustomerID:      customerID,
		TemplateID:      &templateID,
		Date:            time.Now(),
		ValidTill:       validTill,
		Status:          "draft",
		Subject:         template.Name,
		Terms:           template.Terms,
		Notes:           template.Notes,
		DiscountAmount:  0, // 将根据项目计算
		TaxAmount:       0, // 将根据项目计算
	}

	// 创建报价单项目
	var quotationItems []models.QuotationItem
	var totalAmount float64

	for _, templateItem := range template.Items {
		amount := templateItem.Quantity * templateItem.Rate
		discountAmount := amount * templateItem.DiscountRate / 100
		taxAmount := (amount - discountAmount) * templateItem.TaxRate / 100
		itemTotal := amount - discountAmount + taxAmount

		quotationItem := models.QuotationItem{
			ItemID:         templateItem.ItemID,
			Description:    templateItem.Description,
			Quantity:       templateItem.Quantity,
			Rate:           templateItem.Rate,
			Amount:         amount,
			DiscountRate:   templateItem.DiscountRate,
			DiscountAmount: discountAmount,
			TaxRate:        templateItem.TaxRate,
			TaxAmount:      taxAmount,
			TotalAmount:    itemTotal,
		}

		quotationItems = append(quotationItems, quotationItem)
		totalAmount += itemTotal
		quotation.DiscountAmount += discountAmount
		quotation.TaxAmount += taxAmount
	}

	quotation.TotalAmount = totalAmount - quotation.DiscountAmount
	quotation.GrandTotal = totalAmount
	quotation.Items = quotationItems

	if err := s.quotationRepo.Create(ctx, quotation); err != nil {
		return nil, fmt.Errorf("创建报价单失败: %w", err)
	}

	// 转换为响应格式
	return &dto.QuotationResponse{
		ID:             quotation.ID,
		Number:         quotation.QuotationNumber,
		Title:          quotation.Subject,
		Status:         quotation.Status,
		ValidUntil:     quotation.ValidTill,
		PaymentTerms:   quotation.Terms,
		Notes:          quotation.Notes,
		SubTotal:       quotation.TotalAmount,
		DiscountAmount: quotation.DiscountAmount,
		TaxAmount:      quotation.TaxAmount,
		TotalAmount:    quotation.GrandTotal,
		CreatedAt:      quotation.CreatedAt,
		UpdatedAt:      quotation.UpdatedAt,
	}, nil
}

// toTemplateResponse 转换为模板响应格式
func (s *QuotationTemplateServiceImpl) toTemplateResponse(template *models.QuotationTemplate) *dto.QuotationTemplateResponse {
	response := &dto.QuotationTemplateResponse{
		ID:           template.ID,
		Name:         template.Name,
		Description:  template.Description,
		IsDefault:    template.IsDefault,
		IsActive:     template.IsActive,
		ValidityDays: template.ValidityDays,
		Terms:        template.Terms,
		Notes:        template.Notes,
		DiscountRate: template.DiscountRate,
		TaxRate:      template.TaxRate,
		CreatedAt:    template.CreatedAt,
		UpdatedAt:    template.UpdatedAt,
	}

	// 转换项目
	if len(template.Items) > 0 {
		items := make([]dto.QuotationTemplateItemResponse, len(template.Items))
		for i, item := range template.Items {
			items[i] = dto.QuotationTemplateItemResponse{
				ID:           item.ID,
				ItemID:       item.ItemID,
				Description:  item.Description,
				Quantity:     item.Quantity,
				Rate:         item.Rate,
				DiscountRate: item.DiscountRate,
				TaxRate:      item.TaxRate,
				SortOrder:    item.SortOrder,
			}
		}
		response.Items = items
	}

	return response
}

// QuotationVersionService 报价单版本管理服务接口
type QuotationVersionService interface {
	CreateVersion(ctx context.Context, req *dto.QuotationVersionCreateRequest) (*dto.QuotationVersionResponse, error)
	GetVersionsByQuotation(ctx context.Context, quotationID uint) ([]*dto.QuotationVersionResponse, error)
	GetVersion(ctx context.Context, versionID uint) (*dto.QuotationVersionResponse, error)
	SetActiveVersion(ctx context.Context, quotationID uint, versionNumber int) error
	CompareVersions(ctx context.Context, req *dto.QuotationVersionCompareRequest) (*dto.QuotationVersionComparisonResponse, error)
	GetVersionHistory(ctx context.Context, quotationID uint) ([]*dto.QuotationVersionHistoryResponse, error)
	RollbackToVersion(ctx context.Context, req *dto.QuotationVersionRollbackRequest) (*dto.QuotationVersionResponse, error)
	DeleteVersion(ctx context.Context, versionID uint) error
}

// QuotationVersionServiceImpl 报价单版本管理服务实现
type QuotationVersionServiceImpl struct {
	versionRepo   repositories.QuotationVersionRepository
	quotationRepo repositories.QuotationRepository
}

// NewQuotationVersionService 创建报价单版本管理服务
func NewQuotationVersionService(versionRepo repositories.QuotationVersionRepository, quotationRepo repositories.QuotationRepository) QuotationVersionService {
	return &QuotationVersionServiceImpl{
		versionRepo:   versionRepo,
		quotationRepo: quotationRepo,
	}
}

// CreateVersion 创建新版本
func (s *QuotationVersionServiceImpl) CreateVersion(ctx context.Context, req *dto.QuotationVersionCreateRequest) (*dto.QuotationVersionResponse, error) {
	// 验证报价单是否存在
	quotation, err := s.quotationRepo.GetByID(ctx, req.QuotationID)
	if err != nil {
		return nil, fmt.Errorf("获取报价单失败: %w", err)
	}
	if quotation == nil {
		return nil, errors.New("报价单不存在")
	}

	// 获取下一个版本号
	nextVersion, err := s.versionRepo.GetNextVersionNumber(ctx, req.QuotationID)
	if err != nil {
		return nil, fmt.Errorf("获取下一个版本号失败: %w", err)
	}

	// 创建版本记录
	version := &models.QuotationVersion{
		QuotationID:   req.QuotationID,
		VersionNumber: nextVersion,
		VersionName:   req.VersionName,
		ChangeReason:  req.ChangeReason,
		CreatedBy:     1, // 暂时硬编码，后续从上下文获取
		IsActive:      false, // 新版本默认不激活
		VersionData:   "", // 暂时为空，后续实现版本数据序列化
	}

	err = s.versionRepo.Create(ctx, version)
	if err != nil {
		return nil, fmt.Errorf("创建版本失败: %w", err)
	}

	return s.toVersionResponse(version), nil
}

// GetVersionsByQuotation 获取报价单的所有版本
func (s *QuotationVersionServiceImpl) GetVersionsByQuotation(ctx context.Context, quotationID uint) ([]*dto.QuotationVersionResponse, error) {
	versions, err := s.versionRepo.GetByQuotationID(ctx, quotationID)
	if err != nil {
		return nil, fmt.Errorf("获取版本列表失败: %w", err)
	}

	responses := make([]*dto.QuotationVersionResponse, len(versions))
	for i, version := range versions {
		responses[i] = s.toVersionResponse(version)
	}

	return responses, nil
}

// GetVersion 获取指定版本
func (s *QuotationVersionServiceImpl) GetVersion(ctx context.Context, versionID uint) (*dto.QuotationVersionResponse, error) {
	version, err := s.versionRepo.GetByID(ctx, versionID)
	if err != nil {
		return nil, fmt.Errorf("获取版本失败: %w", err)
	}
	if version == nil {
		return nil, errors.New("版本不存在")
	}

	return s.toVersionResponse(version), nil
}

// SetActiveVersion 设置活跃版本
func (s *QuotationVersionServiceImpl) SetActiveVersion(ctx context.Context, quotationID uint, versionNumber int) error {
	err := s.versionRepo.SetActiveVersion(ctx, quotationID, uint(versionNumber))
	if err != nil {
		return fmt.Errorf("设置活跃版本失败: %w", err)
	}

	return nil
}

// CompareVersions 比较版本
func (s *QuotationVersionServiceImpl) CompareVersions(ctx context.Context, req *dto.QuotationVersionCompareRequest) (*dto.QuotationVersionComparisonResponse, error) {
	// 获取两个版本的数据
	version1, err := s.versionRepo.GetByID(ctx, req.FromVersionID)
	if err != nil {
		return nil, fmt.Errorf("获取版本1失败: %w", err)
	}
	if version1 == nil {
		return nil, errors.New("版本1不存在")
	}

	version2, err := s.versionRepo.GetByID(ctx, req.ToVersionID)
	if err != nil {
		return nil, fmt.Errorf("获取版本2失败: %w", err)
	}
	if version2 == nil {
		return nil, errors.New("版本2不存在")
	}

	// 创建比较结果
	comparison := &dto.QuotationVersionComparisonResponse{
		FieldName:   "version_name",
		OldValue:    version1.VersionName,
		NewValue:    version2.VersionName,
		ChangeType:  "modified",
		Description: "版本名称变更",
	}

	return comparison, nil
}

// GetVersionHistory 获取版本历史
func (s *QuotationVersionServiceImpl) GetVersionHistory(ctx context.Context, quotationID uint) ([]*dto.QuotationVersionHistoryResponse, error) {
	versions, err := s.versionRepo.GetByQuotationID(ctx, quotationID)
	if err != nil {
		return nil, fmt.Errorf("获取版本历史失败: %w", err)
	}

	var versionResponses []dto.QuotationVersionResponse
	var currentVersion int
	for _, version := range versions {
		versionResponses = append(versionResponses, dto.QuotationVersionResponse{
			ID:            version.ID,
			QuotationID:   version.QuotationID,
			VersionNumber: version.VersionNumber,
			VersionName:   version.VersionName,
			ChangeReason:  version.ChangeReason,
			CreatedBy:     version.CreatedBy,
			CreatedAt:     version.CreatedAt,
			IsActive:      version.IsActive,
		})
		if version.IsActive {
			currentVersion = version.VersionNumber
		}
	}

	history := []*dto.QuotationVersionHistoryResponse{
		{
			QuotationID:    quotationID,
			Versions:       versionResponses,
			TotalVersions:  len(versionResponses),
			CurrentVersion: currentVersion,
		},
	}

	return history, nil
}

// RollbackToVersion 回滚到指定版本
func (s *QuotationVersionServiceImpl) RollbackToVersion(ctx context.Context, req *dto.QuotationVersionRollbackRequest) (*dto.QuotationVersionResponse, error) {
	// 获取目标版本
	targetVersion, err := s.versionRepo.GetByID(ctx, req.VersionID)
	if err != nil {
		return nil, fmt.Errorf("获取目标版本失败: %w", err)
	}
	if targetVersion == nil {
		return nil, errors.New("目标版本不存在")
	}

	// 创建新版本（基于目标版本的数据）
	nextVersion, err := s.versionRepo.GetNextVersionNumber(ctx, targetVersion.QuotationID)
	if err != nil {
		return nil, fmt.Errorf("获取下一个版本号失败: %w", err)
	}

	newVersion := &models.QuotationVersion{
		QuotationID:   targetVersion.QuotationID,
		VersionNumber: nextVersion,
		VersionName:   fmt.Sprintf("回滚到版本 %d", targetVersion.VersionNumber),
		ChangeReason:  fmt.Sprintf("回滚到版本 %d: %s", targetVersion.VersionNumber, req.Reason),
		CreatedBy:     1, // 暂时硬编码，后续从上下文获取
		IsActive:      true,
		VersionData:   targetVersion.VersionData,
	}

	err = s.versionRepo.Create(ctx, newVersion)
	if err != nil {
		return nil, fmt.Errorf("创建回滚版本失败: %w", err)
	}

	// 设置为活跃版本
	err = s.versionRepo.SetActiveVersion(ctx, targetVersion.QuotationID, uint(nextVersion))
	if err != nil {
		return nil, fmt.Errorf("设置活跃版本失败: %w", err)
	}

	return s.toVersionResponse(newVersion), nil
}

// DeleteVersion 删除版本
func (s *QuotationVersionServiceImpl) DeleteVersion(ctx context.Context, versionID uint) error {
	// 检查版本是否存在
	version, err := s.versionRepo.GetByID(ctx, versionID)
	if err != nil {
		return fmt.Errorf("获取版本失败: %w", err)
	}
	if version == nil {
		return errors.New("版本不存在")
	}

	// 不允许删除活跃版本
	if version.IsActive {
		return errors.New("不能删除活跃版本")
	}

	err = s.versionRepo.Delete(ctx, versionID)
	if err != nil {
		return fmt.Errorf("删除版本失败: %w", err)
	}

	return nil
}

// toVersionResponse 转换为版本响应格式
func (s *QuotationVersionServiceImpl) toVersionResponse(version *models.QuotationVersion) *dto.QuotationVersionResponse {
	return &dto.QuotationVersionResponse{
		ID:            version.ID,
		QuotationID:   version.QuotationID,
		VersionNumber: version.VersionNumber,
		VersionName:   version.VersionName,
		ChangeReason:  version.ChangeReason,
		CreatedBy:     version.CreatedBy,
		CreatedAt:     version.CreatedAt,
		IsActive:      version.IsActive,
	}
}

// compareVersionData 比较版本数据
func (s *QuotationVersionServiceImpl) compareVersionData(data1, data2 string) []string {
	// 这里可以实现更复杂的数据比较逻辑
	// 暂时返回简单的比较结果
	changes := []string{}
	
	if data1 != data2 {
		changes = append(changes, "版本数据已更改")
	}
	
	if len(changes) == 0 {
		changes = append(changes, "无变更")
	}
	
	return changes
}
