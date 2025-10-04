package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/galaxyerp/galaxyErp/internal/common"
	"github.com/galaxyerp/galaxyErp/internal/dto"
	"github.com/galaxyerp/galaxyErp/internal/models"
	"github.com/galaxyerp/galaxyErp/internal/repositories"
	"github.com/galaxyerp/galaxyErp/internal/utils"
)

// SalesInvoiceService 销售发票服务接口
type SalesInvoiceService interface {
	CreateSalesInvoice(ctx *gin.Context, req *dto.SalesInvoiceCreateRequest) (*dto.SalesInvoiceResponse, error)
	GetSalesInvoice(id uint) (*dto.SalesInvoiceResponse, error)
	UpdateSalesInvoice(ctx *gin.Context, id uint, req *dto.SalesInvoiceUpdateRequest) (*dto.SalesInvoiceResponse, error)
	SubmitSalesInvoice(ctx *gin.Context, id uint) (*dto.SalesInvoiceResponse, error)
	CancelSalesInvoice(ctx *gin.Context, id uint) (*dto.SalesInvoiceResponse, error)
	ListSalesInvoices(req *dto.SalesInvoiceListRequest) (*dto.PaginatedResponse[dto.SalesInvoiceResponse], error)
	DeleteSalesInvoice(ctx *gin.Context, id uint) error
	AddPayment(ctx *gin.Context, invoiceID uint, req *dto.InvoicePaymentCreateRequest) (*dto.SalesInvoiceResponse, error)
	GetPayments(invoiceID uint) ([]dto.InvoicePaymentResponse, error)
}

// SalesInvoiceServiceImpl 销售发票服务实现
type SalesInvoiceServiceImpl struct {
	repository          repositories.SalesInvoiceRepository
	customerRepository  repositories.CustomerRepository
	salesOrderRepository repositories.SalesOrderRepository
	paymentEntryService PaymentEntryService
}

// NewSalesInvoiceService 创建销售发票服务实例
func NewSalesInvoiceService(
	repository repositories.SalesInvoiceRepository, 
	customerRepository repositories.CustomerRepository,
	salesOrderRepository repositories.SalesOrderRepository,
	paymentEntryService PaymentEntryService,
) SalesInvoiceService {
	return &SalesInvoiceServiceImpl{
		repository:           repository,
		customerRepository:   customerRepository,
		salesOrderRepository: salesOrderRepository,
		paymentEntryService:  paymentEntryService,
	}
}

// CreateSalesInvoice 创建销售发票
func (s *SalesInvoiceServiceImpl) CreateSalesInvoice(ctx *gin.Context, req *dto.SalesInvoiceCreateRequest) (*dto.SalesInvoiceResponse, error) {
	userID := utils.GetUserIDFromContext(ctx)
	if userID == 0 {
		return nil, errors.New("用户未认证")
	}

	// 验证客户是否存在
	_, err := s.customerRepository.GetByID(context.Background(), req.CustomerID)
	if err != nil {
		return nil, errors.New("客户不存在")
	}

	// 验证销售订单（如果提供）
	if req.SalesOrderID != nil {
		salesOrder, err := s.salesOrderRepository.GetByID(context.Background(), *req.SalesOrderID)
		if err != nil {
			return nil, errors.New("销售订单不存在")
		}
		if salesOrder.CustomerID != req.CustomerID {
			return nil, errors.New("销售订单与客户不匹配")
		}
	}

	// 验证送货单（如果提供）
	// TODO: 需要添加 DeliveryNoteRepository 依赖来验证送货单
	if req.DeliveryNoteID != nil {
		// 暂时跳过送货单验证，需要添加相应的仓储依赖
		// 这里可以在后续添加 DeliveryNoteRepository 依赖后实现
	}

	// 生成发票编号
	invoiceNumber, err := s.generateInvoiceNumber()
	if err != nil {
		return nil, fmt.Errorf("生成发票编号失败: %v", err)
	}

	// 创建销售发票
	invoice := &models.SalesInvoice{
		InvoiceNumber:    invoiceNumber,
		CustomerID:       req.CustomerID,
		SalesOrderID:     req.SalesOrderID,
		DeliveryNoteID:   req.DeliveryNoteID,
		InvoiceDate:      req.InvoiceDate,
		DueDate:          req.DueDate,
		PostingDate:      req.PostingDate,
		DocStatus:        "Draft",
		PaymentStatus:    "Unpaid",
		Currency:         req.Currency,
		ExchangeRate:     req.ExchangeRate,
		BillingAddress:   req.BillingAddress,
		ShippingAddress:  req.ShippingAddress,
		PaymentTerms:     req.PaymentTerms,
		PaymentTermsDays: req.PaymentTermsDays,
		SalesPersonID:    req.SalesPersonID,
		Territory:        req.Territory,
		CustomerPONumber: req.CustomerPONumber,
		Project:          req.Project,
		CostCenter:       req.CostCenter,
		Terms:            req.Terms,
		Notes:            req.Notes,
		AuditableModel: models.AuditableModel{
			BaseModel: models.BaseModel{},
			CreatedBy: userID,
			UpdatedBy: userID,
		},
	}

	// 计算发票明细和总金额
	var totalAmount float64
	for _, itemReq := range req.Items {
		// 计算金额
		amount := itemReq.Quantity * itemReq.Rate
		discountAmount := amount * itemReq.DiscountPercentage / 100
		netAmount := amount - discountAmount
		taxAmount := netAmount * itemReq.TaxRate / 100

		invoiceItem := models.SalesInvoiceItem{
			SalesOrderItemID:   itemReq.SalesOrderItemID,
			DeliveryNoteItemID: itemReq.DeliveryNoteItemID,
			ItemID:             itemReq.ItemID,
			Description:        itemReq.Description,
			Quantity:           itemReq.Quantity,
			UOM:                itemReq.UOM,
			Rate:               itemReq.Rate,
			Amount:             amount,
			DiscountPercentage: itemReq.DiscountPercentage,
			DiscountAmount:     discountAmount,
			TaxCategory:        itemReq.TaxCategory,
			TaxRate:            itemReq.TaxRate,
			TaxAmount:          taxAmount,
			NetRate:            (netAmount + taxAmount) / itemReq.Quantity,
			NetAmount:          netAmount + taxAmount,
			WarehouseID:        itemReq.WarehouseID,
			BatchNo:            itemReq.BatchNo,
			SerialNo:           itemReq.SerialNo,
			CostCenter:         itemReq.CostCenter,
			Project:            itemReq.Project,
		}

		invoice.Items = append(invoice.Items, invoiceItem)
		totalAmount += netAmount + taxAmount
	}

	// 设置发票总金额
	invoice.SubTotal = totalAmount
	invoice.GrandTotal = totalAmount
	invoice.OutstandingAmount = totalAmount

	// 使用仓储层创建发票
	if err := s.repository.Create(context.Background(), invoice); err != nil {
		return nil, fmt.Errorf("创建销售发票失败: %v", err)
	}

	// 返回创建的发票
	return s.GetSalesInvoice(invoice.ID)
}

// GetSalesInvoice 获取销售发票详情
func (s *SalesInvoiceServiceImpl) GetSalesInvoice(id uint) (*dto.SalesInvoiceResponse, error) {
	invoice, err := s.repository.GetByID(nil, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("销售发票不存在")
		}
		return nil, fmt.Errorf("获取销售发票失败: %v", err)
	}

	return s.convertToSalesInvoiceResponse(invoice), nil
}

// UpdateSalesInvoice 更新销售发票
func (s *SalesInvoiceServiceImpl) UpdateSalesInvoice(ctx *gin.Context, id uint, req *dto.SalesInvoiceUpdateRequest) (*dto.SalesInvoiceResponse, error) {
	userID := utils.GetUserIDFromContext(ctx)
	if userID == 0 {
		return nil, errors.New("用户未认证")
	}

	invoice, err := s.repository.GetByID(context.Background(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("销售发票不存在")
		}
		return nil, fmt.Errorf("获取销售发票失败: %v", err)
	}

	// 检查发票状态，只有草稿状态才能修改
	if invoice.DocStatus != "Draft" {
		return nil, errors.New("只有草稿状态的发票才能修改")
	}

	// 更新发票基本信息
	if req.InvoiceDate != nil {
		invoice.InvoiceDate = *req.InvoiceDate
	}
	if req.DueDate != nil {
		invoice.DueDate = *req.DueDate
	}
	if req.PostingDate != nil {
		invoice.PostingDate = *req.PostingDate
	}
	if req.Currency != nil && *req.Currency != "" {
		invoice.Currency = *req.Currency
	}
	if req.ExchangeRate != nil {
		invoice.ExchangeRate = *req.ExchangeRate
	}
	if req.BillingAddress != nil && *req.BillingAddress != "" {
		invoice.BillingAddress = *req.BillingAddress
	}
	if req.ShippingAddress != nil && *req.ShippingAddress != "" {
		invoice.ShippingAddress = *req.ShippingAddress
	}
	if req.PaymentTerms != nil && *req.PaymentTerms != "" {
		invoice.PaymentTerms = *req.PaymentTerms
	}
	if req.PaymentTermsDays != nil {
		invoice.PaymentTermsDays = *req.PaymentTermsDays
	}
	if req.SalesPersonID != nil {
		invoice.SalesPersonID = req.SalesPersonID
	}
	if req.Territory != nil && *req.Territory != "" {
		invoice.Territory = *req.Territory
	}
	if req.CustomerPONumber != nil && *req.CustomerPONumber != "" {
		invoice.CustomerPONumber = *req.CustomerPONumber
	}
	if req.Project != nil && *req.Project != "" {
		invoice.Project = *req.Project
	}
	if req.CostCenter != nil && *req.CostCenter != "" {
		invoice.CostCenter = *req.CostCenter
	}
	if req.Terms != nil && *req.Terms != "" {
		invoice.Terms = *req.Terms
	}
	if req.Notes != nil && *req.Notes != "" {
		invoice.Notes = *req.Notes
	}

	invoice.UpdatedBy = userID

	// 如果有明细更新，重新计算总金额
	if len(req.Items) > 0 {
		// 清空原有明细
		invoice.Items = nil
		
		// 重新创建明细
		var totalAmount float64
		for _, itemReq := range req.Items {
			// 计算金额
			amount := itemReq.Quantity * itemReq.Rate
			discountAmount := amount * itemReq.DiscountPercentage / 100
			netAmount := amount - discountAmount
			taxAmount := netAmount * itemReq.TaxRate / 100

			invoiceItem := models.SalesInvoiceItem{
				SalesOrderItemID:   itemReq.SalesOrderItemID,
				DeliveryNoteItemID: itemReq.DeliveryNoteItemID,
				ItemID:             itemReq.ItemID,
				Description:        itemReq.Description,
				Quantity:           itemReq.Quantity,
				UOM:                itemReq.UOM,
				Rate:               itemReq.Rate,
				Amount:             amount,
				DiscountPercentage: itemReq.DiscountPercentage,
				DiscountAmount:     discountAmount,
				TaxCategory:        itemReq.TaxCategory,
				TaxRate:            itemReq.TaxRate,
				TaxAmount:          taxAmount,
				NetRate:            (netAmount + taxAmount) / itemReq.Quantity,
				NetAmount:          netAmount + taxAmount,
				WarehouseID:        itemReq.WarehouseID,
				BatchNo:            itemReq.BatchNo,
				SerialNo:           itemReq.SerialNo,
				CostCenter:         itemReq.CostCenter,
				Project:            itemReq.Project,
			}

			invoice.Items = append(invoice.Items, invoiceItem)
			totalAmount += netAmount + taxAmount
		}

		// 更新发票总金额
		invoice.SubTotal = totalAmount
		invoice.GrandTotal = totalAmount
		invoice.OutstandingAmount = totalAmount - invoice.PaidAmount
	}

	// 使用仓储层更新发票
	if err := s.repository.Update(context.Background(), invoice); err != nil {
		return nil, fmt.Errorf("更新销售发票失败: %v", err)
	}

	return s.GetSalesInvoice(invoice.ID)
}

// SubmitSalesInvoice 提交销售发票
func (s *SalesInvoiceServiceImpl) SubmitSalesInvoice(ctx *gin.Context, id uint) (*dto.SalesInvoiceResponse, error) {
	userID := utils.GetUserIDFromContext(ctx)
	if userID == 0 {
		return nil, errors.New("用户未认证")
	}

	return s.updateInvoiceStatus(id, "Submitted", userID, "发票提交")
}

// CancelSalesInvoice 取消销售发票
func (s *SalesInvoiceServiceImpl) CancelSalesInvoice(ctx *gin.Context, id uint) (*dto.SalesInvoiceResponse, error) {
	userID := utils.GetUserIDFromContext(ctx)
	if userID == 0 {
		return nil, errors.New("用户未认证")
	}

	return s.updateInvoiceStatus(id, "Cancelled", userID, "发票取消")
}

// updateInvoiceStatus 更新发票状态
func (s *SalesInvoiceServiceImpl) updateInvoiceStatus(id uint, status string, userID uint, reason string) (*dto.SalesInvoiceResponse, error) {
	invoice, err := s.repository.GetByID(nil, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("销售发票不存在")
		}
		return nil, fmt.Errorf("获取销售发票失败: %v", err)
	}

	// 检查状态转换是否合法
	if !s.isValidStatusTransition(invoice.DocStatus, status) {
		return nil, fmt.Errorf("不能从 %s 状态转换到 %s 状态", invoice.DocStatus, status)
	}

	if err := s.repository.UpdateStatus(nil, id, status); err != nil {
		return nil, fmt.Errorf("更新发票状态失败: %v", err)
	}

	return s.GetSalesInvoice(id)
}

// isValidStatusTransition 检查状态转换是否合法
func (s *SalesInvoiceServiceImpl) isValidStatusTransition(from, to string) bool {
	validTransitions := map[string][]string{
		"Draft":     {"Submitted", "Cancelled"},
		"Submitted": {"Cancelled"},
		"Cancelled": {},
	}

	allowedStates, exists := validTransitions[from]
	if !exists {
		return false
	}

	for _, state := range allowedStates {
		if state == to {
			return true
		}
	}
	return false
}

// generateInvoiceNumber 生成发票编号
func (s *SalesInvoiceServiceImpl) generateInvoiceNumber() (string, error) {
	return s.repository.GetNextInvoiceNumber(nil)
}

// convertToSalesInvoiceResponse 转换为销售发票响应
func (s *SalesInvoiceServiceImpl) convertToSalesInvoiceResponse(invoice *models.SalesInvoice) *dto.SalesInvoiceResponse {
	response := &dto.SalesInvoiceResponse{
		ID:                invoice.ID,
		InvoiceNumber:     invoice.InvoiceNumber,
		CustomerID:        invoice.CustomerID,
		SalesOrderID:      invoice.SalesOrderID,
		DeliveryNoteID:    invoice.DeliveryNoteID,
		InvoiceDate:       invoice.InvoiceDate,
		DueDate:           invoice.DueDate,
		PostingDate:       invoice.PostingDate,
		DocStatus:         invoice.DocStatus,
		PaymentStatus:     invoice.PaymentStatus,
		Currency:          invoice.Currency,
		ExchangeRate:      invoice.ExchangeRate,
		SubTotal:          invoice.SubTotal,
		DiscountAmount:    invoice.DiscountAmount,
		TaxAmount:         invoice.TaxAmount,
		ShippingAmount:    invoice.ShippingAmount,
		GrandTotal:        invoice.GrandTotal,
		OutstandingAmount: invoice.OutstandingAmount,
		PaidAmount:        invoice.PaidAmount,
		BillingAddress:    invoice.BillingAddress,
		ShippingAddress:   invoice.ShippingAddress,
		PaymentTerms:      invoice.PaymentTerms,
		PaymentTermsDays:  invoice.PaymentTermsDays,
		SalesPersonID:     invoice.SalesPersonID,
		Territory:         invoice.Territory,
		CustomerPONumber:  invoice.CustomerPONumber,
		Project:           invoice.Project,
		CostCenter:        invoice.CostCenter,
		Terms:             invoice.Terms,
		Notes:             invoice.Notes,
		CreatedBy:         invoice.CreatedBy,
		SubmittedBy:       invoice.SubmittedBy,
		SubmittedAt:       invoice.SubmittedAt,
		CreatedAt:         invoice.CreatedAt,
		UpdatedAt:         invoice.UpdatedAt,
	}

	// 转换客户信息
	if invoice.Customer.ID != 0 {
		response.Customer = dto.CustomerResponse{
			ID:          invoice.Customer.ID,
			Name:        invoice.Customer.Name,
			Code:        invoice.Customer.Code,
			Email:       invoice.Customer.Email,
			Phone:       invoice.Customer.Phone,
			Address:     invoice.Customer.Address,
			ContactName: invoice.Customer.ContactPerson,
			CreditLimit: invoice.Customer.CreditLimit,
			IsActive:    invoice.Customer.IsActive,
		}
	}

	// 转换明细
	for _, item := range invoice.Items {
		itemResponse := dto.SalesInvoiceItemResponse{
			ID:                 item.ID,
			SalesInvoiceID:     item.SalesInvoiceID,
			SalesOrderItemID:   item.SalesOrderItemID,
			DeliveryNoteItemID: item.DeliveryNoteItemID,
			ItemID:             item.ItemID,
			Description:        item.Description,
			Quantity:           item.Quantity,
			UOM:                item.UOM,
			Rate:               item.Rate,
			Amount:             item.Amount,
			DiscountPercentage: item.DiscountPercentage,
			DiscountAmount:     item.DiscountAmount,
			TaxCategory:        item.TaxCategory,
			TaxRate:            item.TaxRate,
			TaxAmount:          item.TaxAmount,
			NetRate:            item.NetRate,
			NetAmount:          item.NetAmount,
			WarehouseID:        item.WarehouseID,
			BatchNo:            item.BatchNo,
			SerialNo:           item.SerialNo,
			CostCenter:         item.CostCenter,
			Project:            item.Project,
		}

		if item.Item.ID != 0 {
			itemResponse.Item = dto.ItemResponse{
				ID:   item.Item.ID,
				Code: item.Item.Code,
				Name: item.Item.Name,
			}
			itemResponse.ItemCode = item.Item.Code
			itemResponse.ItemName = item.Item.Name
		}

		response.Items = append(response.Items, itemResponse)
	}

	// 转换付款记录
	for _, payment := range invoice.Payments {
		paymentResponse := dto.InvoicePaymentResponse{
			ID:              payment.ID,
			SalesInvoiceID:  payment.SalesInvoiceID,
			PaymentEntryID:  payment.PaymentEntryID,
			PaymentDate:     payment.PaymentDate,
			PaymentMethod:   payment.PaymentMethod,
			Amount:          payment.Amount,
			Currency:        payment.Currency,
			ExchangeRate:    payment.ExchangeRate,
			ReferenceNumber: payment.ReferenceNumber,
			BankAccountID:   payment.BankAccountID,
			Notes:           payment.Notes,
			Status:          payment.Status,
			CreatedAt:       payment.CreatedAt,
		}
		response.Payments = append(response.Payments, paymentResponse)
	}

	// 转换状态日志
	for _, log := range invoice.StatusLogs {
		logResponse := dto.InvoiceStatusLogResponse{
			ID:             log.ID,
			SalesInvoiceID: log.SalesInvoiceID,
			FromStatus:     log.FromStatus,
			ToStatus:       log.ToStatus,
			StatusType:     log.StatusType,
			ChangedBy:      log.ChangedBy,
			ChangedAt:      log.ChangedAt,
			Reason:         log.Reason,
			Notes:          log.Notes,
		}
		response.StatusLogs = append(response.StatusLogs, logResponse)
	}

	return response
}

// ListSalesInvoices 获取销售发票列表
func (s *SalesInvoiceServiceImpl) ListSalesInvoices(req *dto.SalesInvoiceListRequest) (*dto.PaginatedResponse[dto.SalesInvoiceResponse], error) {
	// 构建查询选项
	queryOptions := &common.QueryOptions{
		Pagination: &req.PaginationRequest,
	}

	// 使用仓储层获取发票列表
	invoices, total, err := s.repository.List(context.Background(), queryOptions)
	if err != nil {
		return nil, fmt.Errorf("查询销售发票列表失败: %v", err)
	}

	// 转换为响应格式
	var responses []dto.SalesInvoiceResponse
	for _, invoice := range invoices {
		responses = append(responses, *s.convertToSalesInvoiceResponse(invoice))
	}

	// 计算总页数
	limit := req.GetLimit()
	totalPages := int((total + int64(limit) - 1) / int64(limit))

	return &dto.PaginatedResponse[dto.SalesInvoiceResponse]{
		Data:       responses,
		Total:      total,
		Page:       req.Page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

// DeleteSalesInvoice 删除销售发票
func (s *SalesInvoiceServiceImpl) DeleteSalesInvoice(ctx *gin.Context, id uint) error {
	userID := utils.GetUserIDFromContext(ctx)
	if userID == 0 {
		return errors.New("用户未认证")
	}

	invoice, err := s.repository.GetByID(nil, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("销售发票不存在")
		}
		return fmt.Errorf("获取销售发票失败: %v", err)
	}

	// 只有草稿状态的发票才能删除
	if invoice.DocStatus != "Draft" {
		return errors.New("只有草稿状态的发票才能删除")
	}

	return s.repository.Delete(nil, id)
}

// AddPayment 为销售发票添加付款记录
func (s *SalesInvoiceServiceImpl) AddPayment(ctx *gin.Context, invoiceID uint, req *dto.InvoicePaymentCreateRequest) (*dto.SalesInvoiceResponse, error) {
	userID := utils.GetUserIDFromContext(ctx)
	if userID == 0 {
		return nil, errors.New("用户未认证")
	}

	// 验证发票是否存在
	invoice, err := s.repository.GetByID(nil, invoiceID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("销售发票不存在")
		}
		return nil, fmt.Errorf("获取销售发票失败: %v", err)
	}

	// 只有已提交的发票才能添加付款
	if invoice.DocStatus != "Submitted" {
		return nil, errors.New("只有已提交的发票才能添加付款")
	}

	// 创建付款记录
	payment := &models.InvoicePayment{
		SalesInvoiceID:  invoiceID,
		PaymentDate:     req.PaymentDate,
		Amount:          req.Amount,
		PaymentMethod:   req.PaymentMethod,
		Currency:        req.Currency,
		ExchangeRate:    req.ExchangeRate,
		ReferenceNumber: req.ReferenceNumber,
		BankAccountID:   req.BankAccountID,
		Notes:           req.Notes,
		Status:          "Pending",
	}

	// 创建 PaymentEntry 记录
	paymentEntry := &models.PaymentEntry{
		PaymentType:    "Receive", // 销售发票收款
		PartyType:      "Customer",
		PartyID:        invoice.CustomerID,
		PostingDate:    req.PaymentDate,
		PaidAmount:     0,          // 我们收到的金额
		ReceivedAmount: req.Amount, // 客户支付的金额
		Currency:       req.Currency,
		ExchangeRate:   req.ExchangeRate,
		BankAccountID:  req.BankAccountID,
		Reference:      req.ReferenceNumber,
		Remarks:        req.Notes,
		Status:         "submitted",
		IsPosted:       true,
	}

	// 设置过账时间
	now := time.Now()
	paymentEntry.PostedAt = &now

	// 使用 PaymentEntryService 创建付款记录
	if err := s.paymentEntryService.CreatePaymentEntry(context.Background(), paymentEntry); err != nil {
		return nil, fmt.Errorf("创建付款记录失败: %v", err)
	}

	// 设置 InvoicePayment 的 PaymentEntryID
	payment.PaymentEntryID = &paymentEntry.ID

	// 使用仓储层的事务方法创建付款记录并更新发票状态
	err = s.repository.AddPaymentWithTransaction(context.Background(), invoiceID, payment)

	if err != nil {
		return nil, err
	}

	// 返回更新后的发票信息
	return s.GetSalesInvoice(invoiceID)
}

// GetPayments 获取销售发票的付款记录
func (s *SalesInvoiceServiceImpl) GetPayments(invoiceID uint) ([]dto.InvoicePaymentResponse, error) {
	// 验证发票是否存在
	_, err := s.repository.GetByID(nil, invoiceID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("销售发票不存在")
		}
		return nil, fmt.Errorf("获取销售发票失败: %v", err)
	}

	// 获取付款记录
	payments, err := s.repository.GetPayments(context.Background(), invoiceID)
	if err != nil {
		return nil, fmt.Errorf("获取付款记录失败: %v", err)
	}

	// 转换为响应格式
	var responses []dto.InvoicePaymentResponse
	for _, payment := range payments {
		responses = append(responses, dto.InvoicePaymentResponse{
			ID:              payment.ID,
			SalesInvoiceID:  payment.SalesInvoiceID,
			PaymentEntryID:  payment.PaymentEntryID,
			PaymentDate:     payment.PaymentDate,
			PaymentMethod:   payment.PaymentMethod,
			Amount:          payment.Amount,
			Currency:        payment.Currency,
			ExchangeRate:    payment.ExchangeRate,
			ReferenceNumber: payment.ReferenceNumber,
			BankAccountID:   payment.BankAccountID,
			Notes:           payment.Notes,
			Status:          payment.Status,
			CreatedAt:       payment.CreatedAt,
		})
	}

	return responses, nil
}
