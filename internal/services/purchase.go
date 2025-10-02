package services

import (
	"context"
	"fmt"
	"time"

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

	offset := (page - 1) * limit

	suppliers, total, err := s.supplierRepo.List(ctx, offset, limit)
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

	purchaseRequest := &models.PurchaseRequest{
		RequestNumber: requestNumber,
		RequestDate:   time.Now(),
		RequiredBy:    req.RequiredDate,
		Status:        "draft",
	}

	if err := s.purchaseRequestRepo.Create(ctx, purchaseRequest); err != nil {
		return nil, fmt.Errorf("创建采购申请失败: %w", err)
	}

	return s.convertToPurchaseRequestResponse(purchaseRequest), nil
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

	offset := (page - 1) * limit

	purchaseRequests, total, err := s.purchaseRequestRepo.List(ctx, offset, limit)
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
	return &dto.PurchaseRequestResponse{
		ID:           purchaseRequest.ID,
		Number:       purchaseRequest.RequestNumber,
		Title:        "",       // 模型中没有Title字段
		Description:  "",       // 模型中没有Description字段
		Priority:     "medium", // 默认优先级，模型中没有Priority字段
		Status:       purchaseRequest.Status,
		RequiredDate: purchaseRequest.RequiredBy,
		TotalAmount:  0,                                   // 需要计算，模型中没有TotalAmount字段
		Items:        []dto.PurchaseRequestItemResponse{}, // 需要单独处理
		CreatedBy:    dto.UserResponse{},                  // 需要单独处理
		ApprovedBy:   nil,                                 // 需要单独处理
		CreatedAt:    purchaseRequest.CreatedAt,
		UpdatedAt:    purchaseRequest.UpdatedAt,
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

	purchaseOrder := &models.PurchaseOrder{
		OrderNumber:  orderNumber,
		SupplierID:   req.SupplierID,
		OrderDate:    req.OrderDate,
		DeliveryDate: deliveryDate,
		Status:       "draft",
		Terms:        req.PaymentTerms,
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

	offset := (page - 1) * limit

	purchaseOrders, total, err := s.purchaseOrderRepo.List(ctx, offset, limit)
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

	if purchaseOrder.Status != "sent" {
		return fmt.Errorf("只有已发送的订单才能确认")
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
		PaymentTerms:  purchaseOrder.Terms,
		Terms:         purchaseOrder.Terms,
		Notes:         purchaseOrder.Notes,
		SubTotal:      purchaseOrder.TotalAmount,
		TotalDiscount: purchaseOrder.DiscountAmount,
		TotalTax:      purchaseOrder.TaxAmount,
		TotalAmount:   purchaseOrder.GrandTotal,
		Items:         []dto.PurchaseOrderItemResponse{},
	}
}
