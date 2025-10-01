package dto

import (
	"time"
)

// CustomerCreateRequest 客户创建请求
type CustomerCreateRequest struct {
	Name         string  `json:"name" binding:"required,max=100"`
	Code         string  `json:"code" binding:"required,max=50"`
	Type         string  `json:"type" binding:"required,oneof=individual corporate"`
	ContactName  string  `json:"contact_name,omitempty" binding:"max=50"`
	Email        string  `json:"email,omitempty" binding:"omitempty,email"`
	Phone        string  `json:"phone,omitempty" binding:"max=20"`
	Address      string  `json:"address,omitempty"`
	TaxNumber    string  `json:"tax_number,omitempty" binding:"max=50"`
	PaymentTerms string  `json:"payment_terms,omitempty"`
	CreditLimit  float64 `json:"credit_limit,omitempty" binding:"min=0"`
}

// CustomerUpdateRequest 客户更新请求
type CustomerUpdateRequest struct {
	Name         string   `json:"name,omitempty" binding:"omitempty,max=100"`
	Type         string   `json:"type,omitempty" binding:"omitempty,oneof=individual corporate"`
	ContactName  string   `json:"contact_name,omitempty" binding:"omitempty,max=50"`
	Email        string   `json:"email,omitempty" binding:"omitempty,email"`
	Phone        string   `json:"phone,omitempty" binding:"omitempty,max=20"`
	Address      string   `json:"address,omitempty"`
	TaxNumber    string   `json:"tax_number,omitempty" binding:"omitempty,max=50"`
	PaymentTerms string   `json:"payment_terms,omitempty"`
	CreditLimit  *float64 `json:"credit_limit,omitempty" binding:"omitempty,min=0"`
	IsActive     *bool    `json:"is_active,omitempty"`
}

// CustomerResponse 客户响应
type CustomerResponse struct {
	ID           uint      `json:"id"`
	Name         string    `json:"name"`
	Code         string    `json:"code"`
	Type         string    `json:"type"`
	ContactName  string    `json:"contact_name,omitempty"`
	Email        string    `json:"email,omitempty"`
	Phone        string    `json:"phone,omitempty"`
	Address      string    `json:"address,omitempty"`
	TaxNumber    string    `json:"tax_number,omitempty"`
	PaymentTerms string    `json:"payment_terms,omitempty"`
	CreditLimit  float64   `json:"credit_limit"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// QuotationCreateRequest 报价单创建请求
type QuotationCreateRequest struct {
	CustomerID   uint                    `json:"customer_id" binding:"required"`
	Title        string                  `json:"title" binding:"required,max=200"`
	Description  string                  `json:"description,omitempty"`
	ValidUntil   time.Time               `json:"valid_until" binding:"required"`
	PaymentTerms string                  `json:"payment_terms,omitempty"`
	Notes        string                  `json:"notes,omitempty"`
	Items        []QuotationItemRequest  `json:"items" binding:"required,min=1"`
}

// QuotationItemRequest 报价单项目请求
type QuotationItemRequest struct {
	ItemID      uint    `json:"item_id" binding:"required"`
	Quantity    float64 `json:"quantity" binding:"required,gt=0"`
	UnitPrice   float64 `json:"unit_price" binding:"required,gt=0"`
	Discount    float64 `json:"discount,omitempty" binding:"min=0,max=100"`
	TaxRate     float64 `json:"tax_rate,omitempty" binding:"min=0,max=100"`
	Notes       string  `json:"notes,omitempty"`
}

// QuotationUpdateRequest 报价单更新请求
type QuotationUpdateRequest struct {
	Title        string                  `json:"title,omitempty" binding:"omitempty,max=200"`
	Description  string                  `json:"description,omitempty"`
	ValidUntil   *time.Time              `json:"valid_until,omitempty"`
	PaymentTerms string                  `json:"payment_terms,omitempty"`
	Notes        string                  `json:"notes,omitempty"`
	Items        []QuotationItemRequest  `json:"items,omitempty"`
}

// QuotationResponse 报价单响应
type QuotationResponse struct {
	ID           uint                     `json:"id"`
	Number       string                   `json:"number"`
	Title        string                   `json:"title"`
	Description  string                   `json:"description,omitempty"`
	Status       string                   `json:"status"`
	ValidUntil   time.Time                `json:"valid_until"`
	PaymentTerms string                   `json:"payment_terms,omitempty"`
	Notes        string                   `json:"notes,omitempty"`
	SubTotal     float64                  `json:"sub_total"`
	DiscountAmount float64                `json:"discount_amount"`
	TaxAmount    float64                  `json:"tax_amount"`
	TotalAmount  float64                  `json:"total_amount"`
	Customer     CustomerResponse         `json:"customer"`
	Items        []QuotationItemResponse  `json:"items"`
	CreatedBy    UserResponse             `json:"created_by"`
	CreatedAt    time.Time                `json:"created_at"`
	UpdatedAt    time.Time                `json:"updated_at"`
}

// QuotationItemResponse 报价单项目响应
type QuotationItemResponse struct {
	ID            uint         `json:"id"`
	Quantity      float64      `json:"quantity"`
	UnitPrice     float64      `json:"unit_price"`
	Discount      float64      `json:"discount"`
	DiscountAmount float64     `json:"discount_amount"`
	TaxRate       float64      `json:"tax_rate"`
	TaxAmount     float64      `json:"tax_amount"`
	Amount        float64      `json:"amount"`
	Notes         string       `json:"notes,omitempty"`
	Item          ItemResponse `json:"item"`
}

// SalesOrderCreateRequest 销售订单创建请求
type SalesOrderCreateRequest struct {
	CustomerID      uint                     `json:"customer_id" binding:"required"`
	QuotationID     *uint                    `json:"quotation_id,omitempty"`
	OrderDate       time.Time                `json:"order_date" binding:"required"`
	DeliveryDate    *time.Time               `json:"delivery_date,omitempty"`
	ExpectedDate    time.Time                `json:"expected_date" binding:"required"`
	Status          string                   `json:"status,omitempty"`
	PaymentTerms    string                   `json:"payment_terms,omitempty"`
	ShippingAddress string                   `json:"shipping_address,omitempty"`
	Notes           string                   `json:"notes,omitempty"`
	Items           []SalesOrderItemRequest  `json:"items" binding:"required,min=1"`
}

// SalesOrderItemRequest 销售订单项目请求
type SalesOrderItemRequest struct {
	ItemID      uint    `json:"item_id" binding:"required"`
	Quantity    float64 `json:"quantity" binding:"required,gt=0"`
	UnitPrice   float64 `json:"unit_price" binding:"required,gt=0"`
	Discount    float64 `json:"discount,omitempty" binding:"min=0,max=100"`
	TaxRate     float64 `json:"tax_rate,omitempty" binding:"min=0,max=100"`
	Notes       string  `json:"notes,omitempty"`
}

// SalesOrderUpdateRequest 销售订单更新请求
type SalesOrderUpdateRequest struct {
	CustomerID      *uint                    `json:"customer_id,omitempty"`
	DeliveryDate    *time.Time               `json:"delivery_date,omitempty"`
	ExpectedDate    *time.Time               `json:"expected_date,omitempty"`
	Status          *string                  `json:"status,omitempty"`
	PaymentTerms    string                   `json:"payment_terms,omitempty"`
	ShippingAddress string                   `json:"shipping_address,omitempty"`
	Notes           *string                  `json:"notes,omitempty"`
	Items           []SalesOrderItemRequest  `json:"items,omitempty"`
}

// SalesOrderResponse 销售订单响应
type SalesOrderResponse struct {
	ID              uint                      `json:"id"`
	Number          string                    `json:"number"`
	Status          string                    `json:"status"`
	ExpectedDate    time.Time                 `json:"expected_date"`
	PaymentTerms    string                    `json:"payment_terms,omitempty"`
	ShippingAddress string                    `json:"shipping_address,omitempty"`
	Notes           string                    `json:"notes,omitempty"`
	SubTotal        float64                   `json:"sub_total"`
	DiscountAmount  float64                   `json:"discount_amount"`
	TaxAmount       float64                   `json:"tax_amount"`
	TotalAmount     float64                   `json:"total_amount"`
	Customer        CustomerResponse          `json:"customer"`
	Quotation       *QuotationResponse        `json:"quotation,omitempty"`
	Items           []SalesOrderItemResponse  `json:"items"`
	CreatedBy       UserResponse              `json:"created_by"`
	ApprovedBy      *UserResponse             `json:"approved_by,omitempty"`
	CreatedAt       time.Time                 `json:"created_at"`
	UpdatedAt       time.Time                 `json:"updated_at"`
}

// SalesOrderItemResponse 销售订单项目响应
type SalesOrderItemResponse struct {
	ID             uint         `json:"id"`
	SalesOrderID   uint         `json:"sales_order_id"`
	ItemID         uint         `json:"item_id"`
	Quantity       float64      `json:"quantity"`
	UnitPrice      float64      `json:"unit_price"`
	Discount       float64      `json:"discount"`
	DiscountAmount float64      `json:"discount_amount"`
	TaxRate        float64      `json:"tax_rate"`
	TaxAmount      float64      `json:"tax_amount"`
	LineTotal      float64      `json:"line_total"`
	DeliveredQty   float64      `json:"delivered_qty"`
	Description    string       `json:"description,omitempty"`
	Item           ItemResponse `json:"item"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
}

// DeliveryCreateRequest 发货创建请求
type DeliveryCreateRequest struct {
	OrderID      uint                   `json:"order_id" binding:"required"`
	WarehouseID  uint                   `json:"warehouse_id" binding:"required"`
	DeliveryDate time.Time              `json:"delivery_date" binding:"required"`
	TrackingNumber string               `json:"tracking_number,omitempty"`
	Carrier      string                 `json:"carrier,omitempty"`
	Notes        string                 `json:"notes,omitempty"`
	Items        []DeliveryItemRequest  `json:"items" binding:"required,min=1"`
}

// DeliveryItemRequest 发货项目请求
type DeliveryItemRequest struct {
	OrderItemID  uint    `json:"order_item_id" binding:"required"`
	LocationID   uint    `json:"location_id" binding:"required"`
	DeliveredQty float64 `json:"delivered_qty" binding:"required,gt=0"`
	Notes        string  `json:"notes,omitempty"`
}

// DeliveryResponse 发货响应
type DeliveryResponse struct {
	ID             uint                    `json:"id"`
	Number         string                  `json:"number"`
	DeliveryDate   time.Time               `json:"delivery_date"`
	TrackingNumber string                  `json:"tracking_number,omitempty"`
	Carrier        string                  `json:"carrier,omitempty"`
	Notes          string                  `json:"notes,omitempty"`
	Status         string                  `json:"status"`
	Order          SalesOrderResponse      `json:"order"`
	Warehouse      WarehouseResponse       `json:"warehouse"`
	Items          []DeliveryItemResponse  `json:"items"`
	CreatedBy      UserResponse            `json:"created_by"`
	CreatedAt      time.Time               `json:"created_at"`
	UpdatedAt      time.Time               `json:"updated_at"`
}

// DeliveryItemResponse 发货项目响应
type DeliveryItemResponse struct {
	ID           uint                    `json:"id"`
	DeliveredQty float64                 `json:"delivered_qty"`
	Notes        string                  `json:"notes,omitempty"`
	OrderItem    SalesOrderItemResponse  `json:"order_item"`
	Location     LocationResponse        `json:"location"`
}

// InvoiceCreateRequest 发票创建请求
type InvoiceCreateRequest struct {
	OrderID     uint                  `json:"order_id" binding:"required"`
	InvoiceDate time.Time             `json:"invoice_date" binding:"required"`
	DueDate     time.Time             `json:"due_date" binding:"required"`
	Notes       string                `json:"notes,omitempty"`
	Items       []InvoiceItemRequest  `json:"items" binding:"required,min=1"`
}

// InvoiceItemRequest 发票项目请求
type InvoiceItemRequest struct {
	OrderItemID uint    `json:"order_item_id" binding:"required"`
	Quantity    float64 `json:"quantity" binding:"required,gt=0"`
	UnitPrice   float64 `json:"unit_price" binding:"required,gt=0"`
	Notes       string  `json:"notes,omitempty"`
}

// InvoiceResponse 发票响应
type InvoiceResponse struct {
	ID             uint                   `json:"id"`
	Number         string                 `json:"number"`
	InvoiceDate    time.Time              `json:"invoice_date"`
	DueDate        time.Time              `json:"due_date"`
	Notes          string                 `json:"notes,omitempty"`
	Status         string                 `json:"status"`
	SubTotal       float64                `json:"sub_total"`
	TaxAmount      float64                `json:"tax_amount"`
	TotalAmount    float64                `json:"total_amount"`
	PaidAmount     float64                `json:"paid_amount"`
	BalanceAmount  float64                `json:"balance_amount"`
	Order          SalesOrderResponse     `json:"order"`
	Items          []InvoiceItemResponse  `json:"items"`
	CreatedBy      UserResponse           `json:"created_by"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
}

// InvoiceItemResponse 发票项目响应
type InvoiceItemResponse struct {
	ID        uint                    `json:"id"`
	Quantity  float64                 `json:"quantity"`
	UnitPrice float64                 `json:"unit_price"`
	Amount    float64                 `json:"amount"`
	Notes     string                  `json:"notes,omitempty"`
	OrderItem SalesOrderItemResponse  `json:"order_item"`
}

// SalesOrderFilter 销售订单过滤器
type SalesOrderFilter struct {
	SearchRequest
	CustomerID   *uint      `json:"customer_id,omitempty" form:"customer_id"`
	Status       string     `json:"status,omitempty" form:"status"`
	StartDate    *time.Time `json:"start_date,omitempty" form:"start_date"`
	EndDate      *time.Time `json:"end_date,omitempty" form:"end_date"`
	MinAmount    *float64   `json:"min_amount,omitempty" form:"min_amount"`
	MaxAmount    *float64   `json:"max_amount,omitempty" form:"max_amount"`
}

// SalesSearchRequest 销售搜索请求
type SalesSearchRequest struct {
	SearchRequest
	CustomerID   *uint      `json:"customer_id,omitempty" form:"customer_id"`
	StartDate    *time.Time `json:"start_date,omitempty" form:"start_date"`
	EndDate      *time.Time `json:"end_date,omitempty" form:"end_date"`
	MinAmount    *float64   `json:"min_amount,omitempty" form:"min_amount"`
	MaxAmount    *float64   `json:"max_amount,omitempty" form:"max_amount"`
}

// SalesApprovalRequest 销售审批请求
type SalesApprovalRequest struct {
	Action   string `json:"action" binding:"required,oneof=approve reject"`
	Comments string `json:"comments,omitempty"`
}

// SalesStatisticsResponse 销售统计响应
type SalesStatisticsResponse struct {
	TotalOrders      int64   `json:"total_orders"`
	TotalAmount      float64 `json:"total_amount"`
	PendingOrders    int64   `json:"pending_orders"`
	ApprovedOrders   int64   `json:"approved_orders"`
	CompletedOrders  int64   `json:"completed_orders"`
	TopCustomers     []CustomerStatistics `json:"top_customers"`
	TopProducts      []ProductStatistics  `json:"top_products"`
	MonthlyTrend     []MonthlySales       `json:"monthly_trend"`
}

// CustomerStatistics 客户统计
type CustomerStatistics struct {
	CustomerID   uint    `json:"customer_id"`
	CustomerName string  `json:"customer_name"`
	OrderCount   int64   `json:"order_count"`
	TotalAmount  float64 `json:"total_amount"`
}

// ProductStatistics 产品统计
type ProductStatistics struct {
	ItemID      uint    `json:"item_id"`
	ItemName    string  `json:"item_name"`
	SoldQty     float64 `json:"sold_qty"`
	TotalAmount float64 `json:"total_amount"`
}

// MonthlySales 月度销售
type MonthlySales struct {
	Month       string  `json:"month"`
	OrderCount  int64   `json:"order_count"`
	TotalAmount float64 `json:"total_amount"`
}

// SalesReportRequest 销售报告请求
type SalesReportRequest struct {
	DateRangeRequest
	CustomerID *uint  `json:"customer_id,omitempty" form:"customer_id"`
	Status     string `json:"status,omitempty" form:"status"`
	Format     string `json:"format" form:"format" binding:"required,oneof=excel pdf csv"`
}

// 类型别名，用于兼容service中的命名
type CreateSalesOrderRequest = SalesOrderCreateRequest
type UpdateSalesOrderRequest = SalesOrderUpdateRequest
type CreateSalesOrderItemRequest = SalesOrderItemRequest
type UpdateSalesOrderItemRequest = SalesOrderItemRequest