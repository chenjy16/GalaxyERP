package services

import (
	"context"
	"fmt"
	"time"

	"github.com/galaxyerp/galaxyErp/internal/common"
	"github.com/galaxyerp/galaxyErp/internal/dto"
	"github.com/galaxyerp/galaxyErp/internal/models"
	"github.com/galaxyerp/galaxyErp/internal/repositories"
)

// SupplierService 供应商服务接口
type SupplierService interface {
	CreateSupplier(ctx context.Context, req *dto.SupplierCreateRequest) (*dto.SupplierResponse, error)
	GetSupplier(ctx context.Context, id uint) (*dto.SupplierResponse, error)
	UpdateSupplier(ctx context.Context, id uint, req *dto.SupplierUpdateRequest) (*dto.SupplierResponse, error)
	DeleteSupplier(ctx context.Context, id uint) error
	ListSuppliers(ctx context.Context, filter *dto.SupplierFilter) (*dto.PaginatedResponse[dto.SupplierResponse], error)
}

// SupplierServiceImpl 供应商服务实现
type SupplierServiceImpl struct {
	supplierRepo repositories.SupplierRepository
}

// NewSupplierService 创建供应商服务实例
func NewSupplierService(supplierRepo repositories.SupplierRepository) SupplierService {
	return &SupplierServiceImpl{
		supplierRepo: supplierRepo,
	}
}

// CreateSupplier 创建供应商
func (s *SupplierServiceImpl) CreateSupplier(ctx context.Context, req *dto.SupplierCreateRequest) (*dto.SupplierResponse, error) {
	supplier := &models.Supplier{
		Code:          req.Code,
		Name:          req.Name,
		Email:         req.Email,
		Phone:         req.Phone,
		Address:       req.Address,
		ContactPerson: req.ContactName,
		CreditLimit:   req.CreditLimit,
	}

	if err := s.supplierRepo.Create(ctx, supplier); err != nil {
		return nil, fmt.Errorf("创建供应商失败: %w", err)
	}

	return s.convertToSupplierResponse(supplier), nil
}

// GetSupplier 获取供应商详情
func (s *SupplierServiceImpl) GetSupplier(ctx context.Context, id uint) (*dto.SupplierResponse, error) {
	supplier, err := s.supplierRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("获取供应商失败: %w", err)
	}

	return s.convertToSupplierResponse(supplier), nil
}

// UpdateSupplier 更新供应商
func (s *SupplierServiceImpl) UpdateSupplier(ctx context.Context, id uint, req *dto.SupplierUpdateRequest) (*dto.SupplierResponse, error) {
	supplier, err := s.supplierRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("获取供应商失败: %w", err)
	}

	// 更新字段
	if req.Name != nil && *req.Name != "" {
		supplier.Name = *req.Name
	}
	if req.Email != nil && *req.Email != "" {
		supplier.Email = *req.Email
	}
	if req.Phone != nil && *req.Phone != "" {
		supplier.Phone = *req.Phone
	}
	if req.Address != nil && *req.Address != "" {
		supplier.Address = *req.Address
	}
	if req.ContactName != nil && *req.ContactName != "" {
		supplier.ContactPerson = *req.ContactName
	}
	if req.CreditLimit != nil && *req.CreditLimit > 0 {
		supplier.CreditLimit = *req.CreditLimit
	}

	if err := s.supplierRepo.Update(ctx, supplier); err != nil {
		return nil, fmt.Errorf("更新供应商失败: %w", err)
	}

	return s.convertToSupplierResponse(supplier), nil
}

// DeleteSupplier 删除供应商
func (s *SupplierServiceImpl) DeleteSupplier(ctx context.Context, id uint) error {
	if err := s.supplierRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("删除供应商失败: %w", err)
	}
	return nil
}

// ListSuppliers 获取供应商列表
func (s *SupplierServiceImpl) ListSuppliers(ctx context.Context, filter *dto.SupplierFilter) (*dto.PaginatedResponse[dto.SupplierResponse], error) {
	page := 1
	limit := 10
	if filter != nil {
		if filter.Page > 0 {
			page = filter.Page
		}
		if filter.PageSize > 0 {
			limit = filter.PageSize
		}
	}

	paginationReq := &dto.PaginationRequest{
		Page:     page,
		PageSize: limit,
	}

	suppliers, total, err := s.supplierRepo.List(ctx, &common.QueryOptions{
		Pagination: paginationReq,
	})
	if err != nil {
		return nil, fmt.Errorf("获取供应商列表失败: %w", err)
	}

	var responses []dto.SupplierResponse
	for _, supplier := range suppliers {
		responses = append(responses, *s.convertToSupplierResponse(supplier))
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	return &dto.PaginatedResponse[dto.SupplierResponse]{
		Data:       responses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

// convertToSupplierResponse 转换为供应商响应
func (s *SupplierServiceImpl) convertToSupplierResponse(supplier *models.Supplier) *dto.SupplierResponse {
	return &dto.SupplierResponse{
		BaseModel: dto.BaseModel{
			ID:        supplier.ID,
			CreatedAt: supplier.CreatedAt,
			UpdatedAt: supplier.UpdatedAt,
		},
		Code:        supplier.Code,
		Name:        supplier.Name,
		Email:       supplier.Email,
		Phone:       supplier.Phone,
		Address:     supplier.Address,
		ContactName: supplier.ContactPerson,
		CreditLimit: supplier.CreditLimit,
		IsActive:    true, // 默认为活跃状态
	}
}

// PurchaseRequestService 采购申请服务接口
type PurchaseRequestService interface {
	CreatePurchaseRequest(ctx context.Context, req *dto.PurchaseRequestCreateRequest) (*dto.PurchaseRequestResponse, error)
	GetPurchaseRequest(ctx context.Context, id uint) (*dto.PurchaseRequestResponse, error)
	UpdatePurchaseRequest(ctx context.Context, id uint, req *dto.PurchaseRequestUpdateRequest) (*dto.PurchaseRequestResponse, error)
	DeletePurchaseRequest(ctx context.Context, id uint) error
	ListPurchaseRequests(ctx context.Context, filter *dto.PurchaseSearchRequest) (*dto.PaginatedResponse[dto.PurchaseRequestResponse], error)
	SubmitPurchaseRequest(ctx context.Context, id uint) error
	ApprovePurchaseRequest(ctx context.Context, id uint, userID uint) error
	RejectPurchaseRequest(ctx context.Context, id uint, userID uint, reason string) error
}

// PurchaseRequestServiceImpl 采购申请服务实现
type PurchaseRequestServiceImpl struct {
	purchaseRequestRepo repositories.PurchaseRequestRepository
}

// NewPurchaseRequestService 创建采购申请服务实例
func NewPurchaseRequestService(purchaseRequestRepo repositories.PurchaseRequestRepository) PurchaseRequestService {
	return &PurchaseRequestServiceImpl{
		purchaseRequestRepo: purchaseRequestRepo,
	}
}

// CreatePurchaseRequest 创建采购申请
func (s *PurchaseRequestServiceImpl) CreatePurchaseRequest(ctx context.Context, req *dto.PurchaseRequestCreateRequest) (*dto.PurchaseRequestResponse, error) {
	// 生成申请编号
	requestNumber := fmt.Sprintf("PR%s%06d", time.Now().Format("20060102"), time.Now().Unix()%1000000)

	// 创建采购申请明细
	var items []models.PurchaseRequestItem
	for _, itemReq := range req.Items {
		item := models.PurchaseRequestItem{
			ItemID:        itemReq.ItemID,
			Quantity:      itemReq.Quantity,
			EstimatedCost: itemReq.UnitPrice,
			Notes:         itemReq.Notes,
		}
		items = append(items, item)
	}

	purchaseRequest := &models.PurchaseRequest{
		RequestNumber: requestNumber,
		Title:         req.Title,
		Description:   req.Description,
		Priority:      req.Priority,
		RequestDate:   time.Now(),
		RequiredBy:    req.RequiredDate,
		Status:        "draft",
		Items:         items, // GORM 会自动创建关联的明细项
	}

	// 创建采购申请（包含明细）
	if err := s.purchaseRequestRepo.Create(ctx, purchaseRequest); err != nil {
		return nil, fmt.Errorf("创建采购申请失败: %w", err)
	}

	// 重新获取完整的采购申请数据（包含明细和关联数据）
	createdRequest, err := s.purchaseRequestRepo.GetByID(ctx, purchaseRequest.ID)
	if err != nil {
		return nil, fmt.Errorf("获取创建的采购申请失败: %w", err)
	}

	return s.convertToPurchaseRequestResponse(createdRequest), nil
}

// GetPurchaseRequest 获取采购申请详情
func (s *PurchaseRequestServiceImpl) GetPurchaseRequest(ctx context.Context, id uint) (*dto.PurchaseRequestResponse, error) {
	purchaseRequest, err := s.purchaseRequestRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("获取采购申请失败: %w", err)
	}

	return s.convertToPurchaseRequestResponse(purchaseRequest), nil
}

// UpdatePurchaseRequest 更新采购申请
func (s *PurchaseRequestServiceImpl) UpdatePurchaseRequest(ctx context.Context, id uint, req *dto.PurchaseRequestUpdateRequest) (*dto.PurchaseRequestResponse, error) {
	purchaseRequest, err := s.purchaseRequestRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("获取采购申请失败: %w", err)
	}

	// 只有草稿状态的申请才能修改
	if purchaseRequest.Status != "draft" {
		return nil, fmt.Errorf("只有草稿状态的申请才能修改")
	}

	// 更新字段
	if req.RequiredDate != nil {
		purchaseRequest.RequiredBy = *req.RequiredDate
	}
	if req.Title != "" {
		// 注意：模型中没有Title字段，这里可能需要映射到其他字段或忽略
	}

	if err := s.purchaseRequestRepo.Update(ctx, purchaseRequest); err != nil {
		return nil, fmt.Errorf("更新采购申请失败: %w", err)
	}

	return s.convertToPurchaseRequestResponse(purchaseRequest), nil
}

// DeletePurchaseRequest 删除采购申请
func (s *PurchaseRequestServiceImpl) DeletePurchaseRequest(ctx context.Context, id uint) error {
	purchaseRequest, err := s.purchaseRequestRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("获取采购申请失败: %w", err)
	}

	// 只有草稿状态的申请才能删除
	if purchaseRequest.Status != "draft" {
		return fmt.Errorf("只有草稿状态的申请才能删除")
	}

	if err := s.purchaseRequestRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("删除采购申请失败: %w", err)
	}
	return nil
}

// ListPurchaseRequests 获取采购申请列表
func (s *PurchaseRequestServiceImpl) ListPurchaseRequests(ctx context.Context, filter *dto.PurchaseSearchRequest) (*dto.PaginatedResponse[dto.PurchaseRequestResponse], error) {
	page := 1
	limit := 10
	if filter != nil {
		if filter.Page > 0 {
			page = filter.Page
		}
		if filter.PageSize > 0 {
			limit = filter.PageSize
		}
	}

	paginationReq := &dto.PaginationRequest{
		Page:     page,
		PageSize: limit,
	}

	purchaseRequests, total, err := s.purchaseRequestRepo.List(ctx, &common.QueryOptions{
		Pagination: paginationReq,
	})
	if err != nil {
		return nil, fmt.Errorf("获取采购申请列表失败: %w", err)
	}

	var responses []dto.PurchaseRequestResponse
	for _, purchaseRequest := range purchaseRequests {
		responses = append(responses, *s.convertToPurchaseRequestResponse(purchaseRequest))
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	return &dto.PaginatedResponse[dto.PurchaseRequestResponse]{
		Data:       responses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

// SubmitPurchaseRequest 提交采购申请
func (s *PurchaseRequestServiceImpl) SubmitPurchaseRequest(ctx context.Context, id uint) error {
	purchaseRequest, err := s.purchaseRequestRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("获取采购申请失败: %w", err)
	}

	if purchaseRequest.Status != "draft" {
		return fmt.Errorf("只有草稿状态的申请才能提交")
	}

	purchaseRequest.Status = "submitted"
	if err := s.purchaseRequestRepo.Update(ctx, purchaseRequest); err != nil {
		return fmt.Errorf("提交采购申请失败: %w", err)
	}

	return nil
}

// ApprovePurchaseRequest 审批采购申请
func (s *PurchaseRequestServiceImpl) ApprovePurchaseRequest(ctx context.Context, id uint, userID uint) error {
	purchaseRequest, err := s.purchaseRequestRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("获取采购申请失败: %w", err)
	}

	if purchaseRequest.Status != "submitted" {
		return fmt.Errorf("只能审批已提交的采购申请")
	}

	// now := time.Now() // 暂时注释，因为模型中没有ApprovedAt字段
	purchaseRequest.Status = "approved"
	purchaseRequest.ApprovedBy = &userID
	// purchaseRequest.ApprovedAt = &now // 字段不存在，需要添加到模型中

	if err := s.purchaseRequestRepo.Update(ctx, purchaseRequest); err != nil {
		return fmt.Errorf("更新采购申请失败: %w", err)
	}

	return nil
}

// RejectPurchaseRequest 拒绝采购申请
func (s *PurchaseRequestServiceImpl) RejectPurchaseRequest(ctx context.Context, id uint, userID uint, reason string) error {
	purchaseRequest, err := s.purchaseRequestRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("获取采购申请失败: %w", err)
	}

	if purchaseRequest.Status != "submitted" {
		return fmt.Errorf("只能拒绝已提交的采购申请")
	}

	// now := time.Now() // 暂时注释，因为模型中没有相关字段
	purchaseRequest.Status = "rejected"
	// purchaseRequest.RejectedBy = &userID // 字段不存在，需要添加到模型中
	// purchaseRequest.RejectedAt = &now // 字段不存在，需要添加到模型中
	// purchaseRequest.RejectReason = reason // 字段不存在，需要添加到模型中

	if err := s.purchaseRequestRepo.Update(ctx, purchaseRequest); err != nil {
		return fmt.Errorf("更新采购申请失败: %w", err)
	}

	return nil
}

// convertToPurchaseRequestResponse 转换为响应格式
func (s *PurchaseRequestServiceImpl) convertToPurchaseRequestResponse(purchaseRequest *models.PurchaseRequest) *dto.PurchaseRequestResponse {
	// 转换采购申请项目
	var items []dto.PurchaseRequestItemResponse
	var totalAmount float64
	for _, item := range purchaseRequest.Items {
		amount := item.Quantity * item.EstimatedCost
		totalAmount += amount

		itemResponse := dto.PurchaseRequestItemResponse{
			ID:        item.ID,
			Quantity:  item.Quantity,
			UnitPrice: item.EstimatedCost,
			Amount:    amount,
			Notes:     item.Notes,
			Item: dto.ItemResponse{
				ID:          item.ItemID,
				Name:        item.Description, // 使用描述作为名称
				Description: item.Description,
			},
		}
		items = append(items, itemResponse)
	}

	return &dto.PurchaseRequestResponse{
		ID:           purchaseRequest.ID,
		Number:       purchaseRequest.RequestNumber,
		Title:        purchaseRequest.Title,
		Description:  purchaseRequest.Description,
		Priority:     purchaseRequest.Priority,
		Status:       purchaseRequest.Status,
		Department:   purchaseRequest.Department,
		RequiredDate: purchaseRequest.RequiredBy,
		TotalAmount:  totalAmount,
		Items:        items,
		CreatedBy: dto.UserResponse{
			ID:        purchaseRequest.CreatedBy,
			FirstName: "用户", // 临时占位符，实际应该从用户表获取
			LastName:  "",
		},
		ApprovedBy: nil, // 需要单独处理
		CreatedAt:  purchaseRequest.CreatedAt,
		UpdatedAt:  purchaseRequest.UpdatedAt,
	}
}

// PurchaseOrderService 采购订单服务接口
type PurchaseOrderService interface {
	CreatePurchaseOrder(ctx context.Context, req *dto.PurchaseOrderCreateRequest) (*dto.PurchaseOrderResponse, error)
	GetPurchaseOrder(ctx context.Context, id uint) (*dto.PurchaseOrderResponse, error)
	UpdatePurchaseOrder(ctx context.Context, id uint, req *dto.PurchaseOrderUpdateRequest) (*dto.PurchaseOrderResponse, error)
	DeletePurchaseOrder(ctx context.Context, id uint) error
	ListPurchaseOrders(ctx context.Context, filter *dto.PurchaseOrderFilter) (*dto.PaginatedResponse[dto.PurchaseOrderResponse], error)
	ConfirmPurchaseOrder(ctx context.Context, id uint) error
	CancelPurchaseOrder(ctx context.Context, id uint) error
}

// PurchaseOrderServiceImpl 采购订单服务实现
type PurchaseOrderServiceImpl struct {
	purchaseOrderRepo repositories.PurchaseOrderRepository
}

// NewPurchaseOrderService 创建采购订单服务实例
func NewPurchaseOrderService(purchaseOrderRepo repositories.PurchaseOrderRepository) PurchaseOrderService {
	return &PurchaseOrderServiceImpl{
		purchaseOrderRepo: purchaseOrderRepo,
	}
}

// CreatePurchaseOrder 创建采购订单
func (s *PurchaseOrderServiceImpl) CreatePurchaseOrder(ctx context.Context, req *dto.PurchaseOrderCreateRequest) (*dto.PurchaseOrderResponse, error) {
	orderNumber := fmt.Sprintf("PO%s%06d", time.Now().Format("20060102"), time.Now().Unix()%1000000)

	var deliveryDate time.Time
	if req.DeliveryDate != nil {
		deliveryDate = *req.DeliveryDate
	} else {
		deliveryDate = req.ExpectedDate
	}

	// 计算总金额
	var totalAmount float64
	var items []models.PurchaseOrderItem

	for _, itemReq := range req.Items {
		amount := itemReq.Quantity * itemReq.UnitPrice
		taxAmount := amount * itemReq.TaxRate / 100
		totalAmount += amount + taxAmount

		item := models.PurchaseOrderItem{
			ItemID:      itemReq.ItemID,
			Description: itemReq.Notes,
			Quantity:    itemReq.Quantity,
			Rate:        itemReq.UnitPrice,
			Amount:      amount,
			TaxRate:     itemReq.TaxRate,
			TaxAmount:   taxAmount,
			TotalAmount: amount + taxAmount,
		}
		items = append(items, item)
	}

	purchaseOrder := &models.PurchaseOrder{
		OrderNumber:  orderNumber,
		SupplierID:   req.SupplierID,
		OrderDate:    req.OrderDate,
		DeliveryDate: deliveryDate,
		Status:       "draft",
		Terms:        req.PaymentTerms,
		TotalAmount:  totalAmount,
		GrandTotal:   totalAmount,
		Items:        items,
	}

	if err := s.purchaseOrderRepo.Create(ctx, purchaseOrder); err != nil {
		return nil, fmt.Errorf("创建采购订单失败: %w", err)
	}

	return s.convertToPurchaseOrderResponse(purchaseOrder), nil
}

// GetPurchaseOrder 获取采购订单详情
func (s *PurchaseOrderServiceImpl) GetPurchaseOrder(ctx context.Context, id uint) (*dto.PurchaseOrderResponse, error) {
	purchaseOrder, err := s.purchaseOrderRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("获取采购订单失败: %w", err)
	}

	if purchaseOrder == nil {
		return nil, fmt.Errorf("采购订单不存在")
	}

	return s.convertToPurchaseOrderResponse(purchaseOrder), nil
}

// UpdatePurchaseOrder 更新采购订单
func (s *PurchaseOrderServiceImpl) UpdatePurchaseOrder(ctx context.Context, id uint, req *dto.PurchaseOrderUpdateRequest) (*dto.PurchaseOrderResponse, error) {
	purchaseOrder, err := s.purchaseOrderRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("获取采购订单失败: %w", err)
	}

	// 更新字段
	if req.SupplierID != nil {
		purchaseOrder.SupplierID = *req.SupplierID
	}
	// if req.ExpectedDate != nil {
	// 	purchaseOrder.ExpectedDate = *req.ExpectedDate // 字段不存在
	// }
	if req.DeliveryDate != nil {
		purchaseOrder.DeliveryDate = *req.DeliveryDate
	}
	if req.Status != nil && *req.Status != "" {
		purchaseOrder.Status = *req.Status
	}
	if req.PaymentTerms != nil && *req.PaymentTerms != "" {
		purchaseOrder.Terms = *req.PaymentTerms
	}

	if err := s.purchaseOrderRepo.Update(ctx, purchaseOrder); err != nil {
		return nil, fmt.Errorf("更新采购订单失败: %w", err)
	}

	return s.convertToPurchaseOrderResponse(purchaseOrder), nil
}

// DeletePurchaseOrder 删除采购订单
func (s *PurchaseOrderServiceImpl) DeletePurchaseOrder(ctx context.Context, id uint) error {
	purchaseOrder, err := s.purchaseOrderRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("获取采购订单失败: %w", err)
	}

	// 只有草稿状态的订单才能删除
	if purchaseOrder.Status != "draft" {
		return fmt.Errorf("只有草稿状态的订单才能删除")
	}

	if err := s.purchaseOrderRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("删除采购订单失败: %w", err)
	}
	return nil
}

// ListPurchaseOrders 获取采购订单列表
func (s *PurchaseOrderServiceImpl) ListPurchaseOrders(ctx context.Context, filter *dto.PurchaseOrderFilter) (*dto.PaginatedResponse[dto.PurchaseOrderResponse], error) {
	page := 1
	limit := 10
	if filter != nil {
		if filter.Page > 0 {
			page = filter.Page
		}
		if filter.PageSize > 0 {
			limit = filter.PageSize
		}
	}

	paginationReq := &dto.PaginationRequest{
		Page:     page,
		PageSize: limit,
	}

	purchaseOrders, total, err := s.purchaseOrderRepo.List(ctx, &common.QueryOptions{
		Pagination: paginationReq,
	})
	if err != nil {
		return nil, fmt.Errorf("获取采购订单列表失败: %w", err)
	}

	var responses []dto.PurchaseOrderResponse
	for _, purchaseOrder := range purchaseOrders {
		responses = append(responses, *s.convertToPurchaseOrderResponse(purchaseOrder))
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	return &dto.PaginatedResponse[dto.PurchaseOrderResponse]{
		Data:       responses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

// ConfirmPurchaseOrder 确认采购订单
func (s *PurchaseOrderServiceImpl) ConfirmPurchaseOrder(ctx context.Context, id uint) error {
	purchaseOrder, err := s.purchaseOrderRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("获取采购订单失败: %w", err)
	}

	// 允许从draft或sent状态确认订单
	if purchaseOrder.Status != "draft" && purchaseOrder.Status != "sent" {
		return fmt.Errorf("只有草稿或已发送的订单才能确认")
	}

	purchaseOrder.Status = "confirmed"
	if err := s.purchaseOrderRepo.Update(ctx, purchaseOrder); err != nil {
		return fmt.Errorf("确认采购订单失败: %w", err)
	}

	return nil
}

// CancelPurchaseOrder 取消采购订单
func (s *PurchaseOrderServiceImpl) CancelPurchaseOrder(ctx context.Context, id uint) error {
	purchaseOrder, err := s.purchaseOrderRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("获取采购订单失败: %w", err)
	}

	if purchaseOrder.Status == "completed" || purchaseOrder.Status == "cancelled" {
		return fmt.Errorf("已完成或已取消的订单不能再次取消")
	}

	purchaseOrder.Status = "cancelled"
	if err := s.purchaseOrderRepo.Update(ctx, purchaseOrder); err != nil {
		return fmt.Errorf("取消采购订单失败: %w", err)
	}

	return nil
}

// convertToPurchaseOrderResponse 转换为响应格式
func (s *PurchaseOrderServiceImpl) convertToPurchaseOrderResponse(purchaseOrder *models.PurchaseOrder) *dto.PurchaseOrderResponse {
	// 转换采购订单项目
	items := make([]dto.PurchaseOrderItemResponse, len(purchaseOrder.Items))
	for i, item := range purchaseOrder.Items {
		items[i] = dto.PurchaseOrderItemResponse{
			ID:          item.ID,
			Quantity:    item.Quantity,
			UnitPrice:   item.Rate,
			TaxRate:     item.TaxRate,
			TaxAmount:   item.TaxAmount,
			Amount:      item.Amount,
			ReceivedQty: item.ReceivedQty,
			Notes:       item.Description,
			Item:        s.convertToItemResponse(&item.Item),
		}
	}

	// 转换供应商信息
	var supplier *dto.SupplierResponse
	if purchaseOrder.Supplier.ID != 0 {
		supplier = &dto.SupplierResponse{
			BaseModel: dto.BaseModel{
				ID:        purchaseOrder.Supplier.ID,
				CreatedAt: purchaseOrder.Supplier.CreatedAt,
				UpdatedAt: purchaseOrder.Supplier.UpdatedAt,
			},
			Name:        purchaseOrder.Supplier.Name,
			Code:        purchaseOrder.Supplier.Code,
			ContactName: purchaseOrder.Supplier.ContactPerson,
			Email:       purchaseOrder.Supplier.Email,
			Phone:       purchaseOrder.Supplier.Phone,
			Address:     purchaseOrder.Supplier.Address,
			CreditLimit: purchaseOrder.Supplier.CreditLimit,
			IsActive:    purchaseOrder.Supplier.IsActive,
		}
	}

	return &dto.PurchaseOrderResponse{
		BaseModel: dto.BaseModel{
			ID:        purchaseOrder.ID,
			CreatedAt: purchaseOrder.CreatedAt,
			UpdatedAt: purchaseOrder.UpdatedAt,
		},
		OrderNumber:   purchaseOrder.OrderNumber,
		SupplierID:    purchaseOrder.SupplierID,
		OrderDate:     purchaseOrder.OrderDate,
		ExpectedDate:  purchaseOrder.DeliveryDate, // 使用DeliveryDate作为ExpectedDate
		DeliveryDate:  &purchaseOrder.DeliveryDate,
		Status:        purchaseOrder.Status,
		Currency:      "CNY",
		ExchangeRate:  1.0,
		PaymentTerms:  purchaseOrder.Terms, // 使用Terms作为PaymentTerms
		Terms:         purchaseOrder.Terms,
		Notes:         purchaseOrder.Notes,
		SubTotal:      purchaseOrder.TotalAmount,
		TotalDiscount: purchaseOrder.DiscountAmount,
		TotalTax:      purchaseOrder.TaxAmount,
		TotalAmount:   purchaseOrder.GrandTotal,
		Supplier:      supplier,
		Items:         items,
	}
}

// convertToItemResponse 转换为物料响应格式
func (s *PurchaseOrderServiceImpl) convertToItemResponse(item *models.Item) dto.ItemResponse {
	if item == nil {
		return dto.ItemResponse{}
	}

	// 需要根据实际的Item模型字段进行映射
	return dto.ItemResponse{
		ID:          item.ID,
		Code:        item.Code,
		Name:        item.Name,
		Description: item.Description,
		Type:        "raw_material", // 默认类型，可以根据实际情况调整
		MinStock:    0,
		MaxStock:    0,
		UnitCost:    item.Cost,
		SalePrice:   item.Price,
		IsActive:    true,
		// Category和Unit需要根据关联关系填充
		Category:  dto.CategoryResponse{},
		Unit:      dto.UnitResponse{},
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}
}
