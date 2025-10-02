package dto

import (
	"time"
)

// SupplierCreateRequest 供应商创建请求
type SupplierCreateRequest struct {
	Name         string  `json:"name" binding:"required,max=100"`
	Code         string  `json:"code" binding:"required,max=50"`
	ContactName  string  `json:"contact_name,omitempty" binding:"max=50"`
	Email        string  `json:"email,omitempty" binding:"omitempty,email"`
	Phone        string  `json:"phone,omitempty" binding:"max=20"`
	Address      string  `json:"address,omitempty"`
	TaxNumber    string  `json:"tax_number,omitempty" binding:"max=50"`
	PaymentTerms string  `json:"payment_terms,omitempty"`
	CreditLimit  float64 `json:"credit_limit,omitempty" binding:"min=0"`
}

// SupplierUpdateRequest 供应商更新请求
type SupplierUpdateRequest struct {
	Code         *string  `json:"code,omitempty" binding:"omitempty,max=50"`
	Name         *string  `json:"name,omitempty" binding:"omitempty,max=100"`
	ContactName  *string  `json:"contact_name,omitempty" binding:"omitempty,max=50"`
	Email        *string  `json:"email,omitempty" binding:"omitempty,email"`
	Phone        *string  `json:"phone,omitempty" binding:"omitempty,max=20"`
	Address      *string  `json:"address,omitempty"`
	TaxNumber    *string  `json:"tax_number,omitempty" binding:"omitempty,max=50"`
	PaymentTerms *string  `json:"payment_terms,omitempty"`
	CreditLimit  *float64 `json:"credit_limit,omitempty" binding:"omitempty,min=0"`
	IsActive     *bool    `json:"is_active,omitempty"`
}

// SupplierResponse 供应商响应
type SupplierResponse struct {
	BaseModel
	Name         string  `json:"name"`
	Code         string  `json:"code"`
	ContactName  string  `json:"contact_name,omitempty"`
	Email        string  `json:"email,omitempty"`
	Phone        string  `json:"phone,omitempty"`
	Address      string  `json:"address,omitempty"`
	TaxNumber    string  `json:"tax_number,omitempty"`
	PaymentTerms string  `json:"payment_terms,omitempty"`
	CreditLimit  float64 `json:"credit_limit"`
	IsActive     bool    `json:"is_active"`
}

// PurchaseRequestCreateRequest 采购申请创建请求
type PurchaseRequestCreateRequest struct {
	Title        string                       `json:"title" binding:"required,max=200"`
	Description  string                       `json:"description,omitempty"`
	Priority     string                       `json:"priority" binding:"required,oneof=low medium high urgent"`
	RequiredDate time.Time                    `json:"required_date" binding:"required"`
	Items        []PurchaseRequestItemRequest `json:"items" binding:"required,min=1"`
}

// PurchaseRequestItemRequest 采购申请项目请求
type PurchaseRequestItemRequest struct {
	ItemID    uint    `json:"item_id" binding:"required"`
	Quantity  float64 `json:"quantity" binding:"required,gt=0"`
	UnitPrice float64 `json:"unit_price,omitempty" binding:"min=0"`
	Notes     string  `json:"notes,omitempty"`
}

// PurchaseRequestUpdateRequest 采购申请更新请求
type PurchaseRequestUpdateRequest struct {
	Title        string                       `json:"title,omitempty" binding:"omitempty,max=200"`
	Description  string                       `json:"description,omitempty"`
	Priority     string                       `json:"priority,omitempty" binding:"omitempty,oneof=low medium high urgent"`
	RequiredDate *time.Time                   `json:"required_date,omitempty"`
	Items        []PurchaseRequestItemRequest `json:"items,omitempty"`
}

// PurchaseRequestResponse 采购申请响应
type PurchaseRequestResponse struct {
	ID           uint                          `json:"id"`
	Number       string                        `json:"number"`
	Title        string                        `json:"title"`
	Description  string                        `json:"description,omitempty"`
	Priority     string                        `json:"priority"`
	Status       string                        `json:"status"`
	RequiredDate time.Time                     `json:"required_date"`
	TotalAmount  float64                       `json:"total_amount"`
	Items        []PurchaseRequestItemResponse `json:"items"`
	CreatedBy    UserResponse                  `json:"created_by"`
	ApprovedBy   *UserResponse                 `json:"approved_by,omitempty"`
	CreatedAt    time.Time                     `json:"created_at"`
	UpdatedAt    time.Time                     `json:"updated_at"`
}

// PurchaseRequestItemResponse 采购申请项目响应
type PurchaseRequestItemResponse struct {
	ID        uint         `json:"id"`
	Quantity  float64      `json:"quantity"`
	UnitPrice float64      `json:"unit_price"`
	Amount    float64      `json:"amount"`
	Notes     string       `json:"notes,omitempty"`
	Item      ItemResponse `json:"item"`
}

// PurchaseOrderCreateRequest 采购订单创建请求
type PurchaseOrderCreateRequest struct {
	SupplierID   uint                       `json:"supplier_id" binding:"required"`
	RequestID    *uint                      `json:"request_id,omitempty"`
	OrderDate    time.Time                  `json:"order_date" binding:"required"`
	ExpectedDate time.Time                  `json:"expected_date" binding:"required"`
	DeliveryDate *time.Time                 `json:"delivery_date,omitempty"`
	Status       string                     `json:"status,omitempty" binding:"omitempty,oneof=draft sent confirmed partial completed cancelled"`
	PaymentTerms string                     `json:"payment_terms,omitempty"`
	Terms        string                     `json:"terms,omitempty"`
	Notes        string                     `json:"notes,omitempty"`
	Items        []PurchaseOrderItemRequest `json:"items" binding:"required,min=1"`
}

// PurchaseOrderItemRequest 采购订单项目请求
type PurchaseOrderItemRequest struct {
	ItemID    uint    `json:"item_id" binding:"required"`
	Quantity  float64 `json:"quantity" binding:"required,gt=0"`
	UnitPrice float64 `json:"unit_price" binding:"required,gt=0"`
	TaxRate   float64 `json:"tax_rate,omitempty" binding:"min=0,max=100"`
	Notes     string  `json:"notes,omitempty"`
}

// PurchaseOrderUpdateRequest 采购订单更新请求
type PurchaseOrderUpdateRequest struct {
	SupplierID   *uint                      `json:"supplier_id,omitempty"`
	OrderDate    *time.Time                 `json:"order_date,omitempty"`
	ExpectedDate *time.Time                 `json:"expected_date,omitempty"`
	DeliveryDate *time.Time                 `json:"delivery_date,omitempty"`
	Status       *string                    `json:"status,omitempty"`
	PaymentTerms *string                    `json:"payment_terms,omitempty"`
	Terms        *string                    `json:"terms,omitempty"`
	Notes        *string                    `json:"notes,omitempty"`
	Items        []PurchaseOrderItemRequest `json:"items,omitempty"`
}

// PurchaseOrderResponse 采购订单响应
type PurchaseOrderResponse struct {
	BaseModel
	OrderNumber       string                      `json:"order_number"`
	SupplierID        uint                        `json:"supplier_id"`
	OrderDate         time.Time                   `json:"order_date"`
	ExpectedDate      time.Time                   `json:"expected_date"`
	DeliveryDate      *time.Time                  `json:"delivery_date,omitempty"`
	Status            string                      `json:"status"`
	Currency          string                      `json:"currency"`
	ExchangeRate      float64                     `json:"exchange_rate"`
	PaymentTerms      string                      `json:"payment_terms,omitempty"`
	DeliveryAddress   string                      `json:"delivery_address,omitempty"`
	BillingAddress    string                      `json:"billing_address,omitempty"`
	Terms             string                      `json:"terms,omitempty"`
	Notes             string                      `json:"notes,omitempty"`
	PurchaseRequestID *uint                       `json:"purchase_request_id,omitempty"`
	SubTotal          float64                     `json:"sub_total"`
	TotalDiscount     float64                     `json:"total_discount"`
	TotalTax          float64                     `json:"total_tax"`
	TotalAmount       float64                     `json:"total_amount"`
	Supplier          *SupplierResponse           `json:"supplier,omitempty"`
	Request           *PurchaseRequestResponse    `json:"request,omitempty"`
	Items             []PurchaseOrderItemResponse `json:"items"`
	CreatedBy         *UserResponse               `json:"created_by,omitempty"`
	ApprovedBy        *UserResponse               `json:"approved_by,omitempty"`
}

// PurchaseOrderItemResponse 采购订单项目响应
type PurchaseOrderItemResponse struct {
	ID          uint         `json:"id"`
	Quantity    float64      `json:"quantity"`
	UnitPrice   float64      `json:"unit_price"`
	TaxRate     float64      `json:"tax_rate"`
	TaxAmount   float64      `json:"tax_amount"`
	Amount      float64      `json:"amount"`
	ReceivedQty float64      `json:"received_qty"`
	Notes       string       `json:"notes,omitempty"`
	Item        ItemResponse `json:"item"`
}

// PurchaseReceiptCreateRequest 采购收货创建请求
type PurchaseReceiptCreateRequest struct {
	OrderID      uint                         `json:"order_id" binding:"required"`
	WarehouseID  uint                         `json:"warehouse_id" binding:"required"`
	ReceivedDate time.Time                    `json:"received_date" binding:"required"`
	Notes        string                       `json:"notes,omitempty"`
	Items        []PurchaseReceiptItemRequest `json:"items" binding:"required,min=1"`
}

// PurchaseReceiptItemRequest 采购收货项目请求
type PurchaseReceiptItemRequest struct {
	OrderItemID   uint    `json:"order_item_id" binding:"required"`
	LocationID    uint    `json:"location_id" binding:"required"`
	ReceivedQty   float64 `json:"received_qty" binding:"required,gt=0"`
	QualityStatus string  `json:"quality_status" binding:"required,oneof=passed failed pending"`
	Notes         string  `json:"notes,omitempty"`
}

// PurchaseReceiptResponse 采购收货响应
type PurchaseReceiptResponse struct {
	ID           uint                          `json:"id"`
	Number       string                        `json:"number"`
	ReceivedDate time.Time                     `json:"received_date"`
	Notes        string                        `json:"notes,omitempty"`
	Status       string                        `json:"status"`
	Order        PurchaseOrderResponse         `json:"order"`
	Warehouse    WarehouseResponse             `json:"warehouse"`
	Items        []PurchaseReceiptItemResponse `json:"items"`
	CreatedBy    UserResponse                  `json:"created_by"`
	CreatedAt    time.Time                     `json:"created_at"`
	UpdatedAt    time.Time                     `json:"updated_at"`
}

// PurchaseReceiptItemResponse 采购收货项目响应
type PurchaseReceiptItemResponse struct {
	ID            uint                      `json:"id"`
	ReceivedQty   float64                   `json:"received_qty"`
	QualityStatus string                    `json:"quality_status"`
	Notes         string                    `json:"notes,omitempty"`
	OrderItem     PurchaseOrderItemResponse `json:"order_item"`
	Location      LocationResponse          `json:"location"`
}

// SupplierFilter 供应商过滤器
type SupplierFilter struct {
	SearchRequest
	IsActive *bool `json:"is_active,omitempty" form:"is_active"`
}

// PurchaseOrderFilter 采购订单过滤器
type PurchaseOrderFilter struct {
	SearchRequest
	SupplierID *uint      `json:"supplier_id,omitempty" form:"supplier_id"`
	Status     string     `json:"status,omitempty" form:"status"`
	StartDate  *time.Time `json:"start_date,omitempty" form:"start_date"`
	EndDate    *time.Time `json:"end_date,omitempty" form:"end_date"`
	MinAmount  *float64   `json:"min_amount,omitempty" form:"min_amount"`
	MaxAmount  *float64   `json:"max_amount,omitempty" form:"max_amount"`
}

// PurchaseSearchRequest 采购搜索请求
type PurchaseSearchRequest struct {
	SearchRequest
	SupplierID *uint      `json:"supplier_id,omitempty" form:"supplier_id"`
	StartDate  *time.Time `json:"start_date,omitempty" form:"start_date"`
	EndDate    *time.Time `json:"end_date,omitempty" form:"end_date"`
	MinAmount  *float64   `json:"min_amount,omitempty" form:"min_amount"`
	MaxAmount  *float64   `json:"max_amount,omitempty" form:"max_amount"`
	Priority   string     `json:"priority,omitempty" form:"priority"`
}

// PurchaseApprovalRequest 采购审批请求
type PurchaseApprovalRequest struct {
	Action   string `json:"action" binding:"required,oneof=approve reject"`
	Comments string `json:"comments,omitempty"`
}

// PurchaseStatisticsResponse 采购统计响应
type PurchaseStatisticsResponse struct {
	TotalOrders     int64                `json:"total_orders"`
	TotalAmount     float64              `json:"total_amount"`
	PendingOrders   int64                `json:"pending_orders"`
	ApprovedOrders  int64                `json:"approved_orders"`
	CompletedOrders int64                `json:"completed_orders"`
	TopSuppliers    []SupplierStatistics `json:"top_suppliers"`
	MonthlyTrend    []MonthlyPurchase    `json:"monthly_trend"`
}

// SupplierStatistics 供应商统计
type SupplierStatistics struct {
	SupplierID   uint    `json:"supplier_id"`
	SupplierName string  `json:"supplier_name"`
	OrderCount   int64   `json:"order_count"`
	TotalAmount  float64 `json:"total_amount"`
}

// MonthlyPurchase 月度采购
type MonthlyPurchase struct {
	Month       string  `json:"month"`
	OrderCount  int64   `json:"order_count"`
	TotalAmount float64 `json:"total_amount"`
}

// PurchaseReportRequest 采购报告请求
type PurchaseReportRequest struct {
	DateRangeRequest
	SupplierID *uint  `json:"supplier_id,omitempty" form:"supplier_id"`
	Status     string `json:"status,omitempty" form:"status"`
	Format     string `json:"format" form:"format" binding:"required,oneof=excel pdf csv"`
}
