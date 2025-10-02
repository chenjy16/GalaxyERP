package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

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
}

// SalesInvoiceServiceImpl 销售发票服务实现
type SalesInvoiceServiceImpl struct {
	db         *gorm.DB
	repository repositories.SalesInvoiceRepository
}

// NewSalesInvoiceService 创建销售发票服务实例
func NewSalesInvoiceService(db *gorm.DB) SalesInvoiceService {
	return &SalesInvoiceServiceImpl{
		db:         db,
		repository: repositories.NewSalesInvoiceRepository(db),
	}
}

// CreateSalesInvoice 创建销售发票
func (s *SalesInvoiceServiceImpl) CreateSalesInvoice(ctx *gin.Context, req *dto.SalesInvoiceCreateRequest) (*dto.SalesInvoiceResponse, error) {
	userID := utils.GetUserIDFromContext(ctx)
	if userID == 0 {
		return nil, errors.New("用户未认证")
	}

	// 验证客户是否存在
	var customer models.Customer
	if err := s.db.First(&customer, req.CustomerID).Error; err != nil {
		return nil, errors.New("客户不存在")
	}

	// 验证销售订单（如果提供）
	if req.SalesOrderID != nil {
		var salesOrder models.SalesOrder
		if err := s.db.First(&salesOrder, *req.SalesOrderID).Error; err != nil {
			return nil, errors.New("销售订单不存在")
		}
		if salesOrder.CustomerID != req.CustomerID {
			return nil, errors.New("销售订单与客户不匹配")
		}
	}

	// 验证送货单（如果提供）
	if req.DeliveryNoteID != nil {
		var deliveryNote models.DeliveryNote
		if err := s.db.First(&deliveryNote, *req.DeliveryNoteID).Error; err != nil {
			return nil, errors.New("送货单不存在")
		}
		if deliveryNote.CustomerID != req.CustomerID {
			return nil, errors.New("送货单与客户不匹配")
		}
	}

	// 生成发票编号
	invoiceNumber, err := s.generateInvoiceNumber()
	if err != nil {
		return nil, fmt.Errorf("生成发票编号失败: %v", err)
	}

	// 开始事务
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

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

	if err := tx.Create(invoice).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("创建销售发票失败: %v", err)
	}

	// 创建发票明细
	var totalAmount float64
	for _, itemReq := range req.Items {
		// 验证物料
		var item models.Item
		if err := tx.First(&item, itemReq.ItemID).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("物料 %d 不存在", itemReq.ItemID)
		}

		// 计算金额
		amount := itemReq.Quantity * itemReq.Rate
		discountAmount := amount * itemReq.DiscountPercentage / 100
		netAmount := amount - discountAmount
		taxAmount := netAmount * itemReq.TaxRate / 100

		invoiceItem := &models.SalesInvoiceItem{
			SalesInvoiceID:     invoice.ID,
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

		if err := tx.Create(invoiceItem).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("创建发票明细失败: %v", err)
		}

		totalAmount += netAmount + taxAmount
	}

	// 更新发票总金额
	invoice.SubTotal = totalAmount
	invoice.GrandTotal = totalAmount
	invoice.OutstandingAmount = totalAmount

	if err := tx.Save(invoice).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("更新发票总金额失败: %v", err)
	}

	// 记录状态变更日志
	statusLog := &models.InvoiceStatusLog{
		SalesInvoiceID: invoice.ID,
		ToStatus:       "Draft",
		StatusType:     "DocStatus",
		ChangedBy:      userID,
		ChangedAt:      time.Now(),
		Reason:         "发票创建",
	}

	if err := tx.Create(statusLog).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("创建状态日志失败: %v", err)
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("提交事务失败: %v", err)
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

	var invoice models.SalesInvoice
	if err := s.db.First(&invoice, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("销售发票不存在")
		}
		return nil, fmt.Errorf("获取销售发票失败: %v", err)
	}

	// 检查发票状态，只有草稿状态才能修改
	if invoice.DocStatus != "Draft" {
		return nil, errors.New("只有草稿状态的发票才能修改")
	}

	// 开始事务
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

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
	if req.Currency != "" {
		invoice.Currency = req.Currency
	}
	if req.ExchangeRate != nil {
		invoice.ExchangeRate = *req.ExchangeRate
	}
	if req.BillingAddress != "" {
		invoice.BillingAddress = req.BillingAddress
	}
	if req.ShippingAddress != "" {
		invoice.ShippingAddress = req.ShippingAddress
	}
	if req.PaymentTerms != "" {
		invoice.PaymentTerms = req.PaymentTerms
	}
	if req.PaymentTermsDays != nil {
		invoice.PaymentTermsDays = *req.PaymentTermsDays
	}
	if req.SalesPersonID != nil {
		invoice.SalesPersonID = req.SalesPersonID
	}
	if req.Territory != "" {
		invoice.Territory = req.Territory
	}
	if req.CustomerPONumber != "" {
		invoice.CustomerPONumber = req.CustomerPONumber
	}
	if req.Project != "" {
		invoice.Project = req.Project
	}
	if req.CostCenter != "" {
		invoice.CostCenter = req.CostCenter
	}
	if req.Terms != "" {
		invoice.Terms = req.Terms
	}
	if req.Notes != "" {
		invoice.Notes = req.Notes
	}

	invoice.UpdatedBy = userID

	// 如果有明细更新，先删除原有明细
	if len(req.Items) > 0 {
		if err := tx.Where("sales_invoice_id = ?", invoice.ID).Delete(&models.SalesInvoiceItem{}).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("删除原有明细失败: %v", err)
		}

		// 重新创建明细
		var totalAmount float64
		for _, itemReq := range req.Items {
			// 验证物料
			var item models.Item
			if err := tx.First(&item, itemReq.ItemID).Error; err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("物料 %d 不存在", itemReq.ItemID)
			}

			// 计算金额
			amount := itemReq.Quantity * itemReq.Rate
			discountAmount := amount * itemReq.DiscountPercentage / 100
			netAmount := amount - discountAmount
			taxAmount := netAmount * itemReq.TaxRate / 100

			invoiceItem := &models.SalesInvoiceItem{
				SalesInvoiceID:     invoice.ID,
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

			if err := tx.Create(invoiceItem).Error; err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("创建发票明细失败: %v", err)
			}

			totalAmount += netAmount + taxAmount
		}

		// 更新发票总金额
		invoice.SubTotal = totalAmount
		invoice.GrandTotal = totalAmount
		invoice.OutstandingAmount = totalAmount - invoice.PaidAmount
	}

	if err := tx.Save(&invoice).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("更新销售发票失败: %v", err)
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("提交事务失败: %v", err)
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
	query := s.db.Model(&models.SalesInvoice{}).
		Preload("Customer").
		Preload("SalesOrder").
		Preload("DeliveryNote")

	// 应用过滤条件
	if req.CustomerID != nil {
		query = query.Where("customer_id = ?", *req.CustomerID)
	}
	if req.DocStatus != "" {
		query = query.Where("doc_status = ?", req.DocStatus)
	}
	if req.PaymentStatus != "" {
		query = query.Where("payment_status = ?", req.PaymentStatus)
	}
	if req.Currency != "" {
		query = query.Where("currency = ?", req.Currency)
	}
	if req.SalesPersonID != nil {
		query = query.Where("sales_person_id = ?", *req.SalesPersonID)
	}
	if req.Territory != "" {
		query = query.Where("territory = ?", req.Territory)
	}
	if req.DateFrom != "" {
		query = query.Where("invoice_date >= ?", req.DateFrom)
	}
	if req.DateTo != "" {
		query = query.Where("invoice_date <= ?", req.DateTo)
	}
	if req.Search != "" {
		query = query.Where("invoice_number LIKE ? OR customer_po_number LIKE ?",
			"%"+req.Search+"%", "%"+req.Search+"%")
	}

	// 获取总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("获取发票总数失败: %v", err)
	}

	// 应用分页和排序
	offset := (req.Page - 1) * req.PageSize
	query = query.Offset(offset).Limit(req.PageSize)

	if req.SortBy != "" {
		order := req.SortBy
		if req.SortDesc {
			order += " DESC"
		}
		query = query.Order(order)
	} else {
		query = query.Order("created_at DESC")
	}

	var invoices []models.SalesInvoice
	if err := query.Find(&invoices).Error; err != nil {
		return nil, fmt.Errorf("获取发票列表失败: %v", err)
	}

	// 转换为响应格式
	var items []dto.SalesInvoiceResponse
	for _, invoice := range invoices {
		items = append(items, *s.convertToSalesInvoiceResponse(&invoice))
	}

	return &dto.PaginatedResponse[dto.SalesInvoiceResponse]{
		Data:       items,
		Total:      total,
		Page:       req.Page,
		Limit:      req.PageSize,
		TotalPages: int((total + int64(req.PageSize) - 1) / int64(req.PageSize)),
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