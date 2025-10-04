package dto

import (
	"time"
)

// CustomerCreateRequest 客户创建请求
type CustomerCreateRequest struct {
	Name         string  `json:"name" validate:"required,max=100"`
	Code         string  `json:"code" validate:"required,max=50"`
	Type         string  `json:"type" validate:"required,oneof=individual corporate"`
	ContactName  string  `json:"contact_name,omitempty" validate:"max=50"`
	Email        string  `json:"email,omitempty" validate:"omitempty,email"`
	Phone        string  `json:"phone,omitempty" validate:"omitempty,chinese_mobile"`
	Address      string  `json:"address,omitempty"`
	TaxNumber    string  `json:"tax_number,omitempty" validate:"max=50"`
	PaymentTerms string  `json:"payment_terms,omitempty"`
	CreditLimit  float64 `json:"credit_limit,omitempty" validate:"min=0,currency"`
}

// CustomerUpdateRequest 客户更新请求
type CustomerUpdateRequest struct {
	Name         string   `json:"name,omitempty" validate:"omitempty,max=100"`
	Type         string   `json:"type,omitempty" validate:"omitempty,oneof=individual corporate"`
	ContactName  string   `json:"contact_name,omitempty" validate:"omitempty,max=50"`
	Email        string   `json:"email,omitempty" validate:"omitempty,email"`
	Phone        string   `json:"phone,omitempty" validate:"omitempty,chinese_mobile"`
	Address      string   `json:"address,omitempty"`
	TaxNumber    string   `json:"tax_number,omitempty" validate:"omitempty,max=50"`
	PaymentTerms string   `json:"payment_terms,omitempty"`
	CreditLimit  *float64 `json:"credit_limit,omitempty" validate:"omitempty,min=0,currency"`
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
	CustomerID   uint                   `json:"customer_id" validate:"required"`
	Title        string                 `json:"title" validate:"required,max=200"`
	Description  string                 `json:"description,omitempty"`
	ValidUntil   time.Time              `json:"valid_until" validate:"required"`
	PaymentTerms string                 `json:"payment_terms,omitempty"`
	Notes        string                 `json:"notes,omitempty"`
	Items        []QuotationItemRequest `json:"items" validate:"required,min=1"`
}

// QuotationItemRequest 报价单项目请求
type QuotationItemRequest struct {
	ItemID    uint    `json:"item_id" validate:"required"`
	Quantity  float64 `json:"quantity" validate:"required,gt=0"`
	UnitPrice float64 `json:"unit_price" validate:"required,gt=0,currency"`
	Discount  float64 `json:"discount,omitempty" validate:"min=0,max=100"`
	TaxRate   float64 `json:"tax_rate,omitempty" validate:"min=0,max=100"`
	Notes     string  `json:"notes,omitempty"`
}

// QuotationUpdateRequest 报价单更新请求
type QuotationUpdateRequest struct {
	CustomerID   *uint                  `json:"customer_id,omitempty"`
	Title        string                 `json:"title,omitempty" validate:"omitempty,max=200"`
	Description  string                 `json:"description,omitempty"`
	ValidUntil   *time.Time             `json:"valid_until,omitempty"`
	PaymentTerms string                 `json:"payment_terms,omitempty"`
	Notes        string                 `json:"notes,omitempty"`
	TotalAmount  *float64               `json:"total_amount,omitempty"`
	Items        []QuotationItemRequest `json:"items,omitempty"`
}

// QuotationResponse 报价单响应
type QuotationResponse struct {
	ID              uint                    `json:"id"`
	Number          string                  `json:"number"`
	QuotationNumber string                  `json:"quotationNumber"` // 前端期望的字段名
	Title           string                  `json:"title"`
	Subject         string                  `json:"subject"` // 前端期望的字段名
	Description     string                  `json:"description,omitempty"`
	Status          string                  `json:"status"`
	ValidUntil      time.Time               `json:"valid_until"`
	ValidTill       time.Time               `json:"validTill"` // 前端期望的字段名
	Date            time.Time               `json:"date"`      // 前端期望的字段名
	PaymentTerms    string                  `json:"payment_terms,omitempty"`
	Notes           string                  `json:"notes,omitempty"`
	SubTotal        float64                 `json:"sub_total"`
	DiscountAmount  float64                 `json:"discount_amount"`
	TaxAmount       float64                 `json:"tax_amount"`
	TotalAmount     float64                 `json:"total_amount"`
	GrandTotal      float64                 `json:"grand_total"` // 前端期望的字段名
	Customer        CustomerResponse        `json:"customer"`
	Items           []QuotationItemResponse `json:"items"`
	CreatedBy       UserResponse            `json:"created_by"`
	CreatedAt       time.Time               `json:"created_at"`
	UpdatedAt       time.Time               `json:"updated_at"`
}

// QuotationItemResponse 报价单项目响应
type QuotationItemResponse struct {
	ID             uint         `json:"id"`
	Quantity       float64      `json:"quantity"`
	UnitPrice      float64      `json:"unit_price"`
	Discount       float64      `json:"discount"`
	DiscountAmount float64      `json:"discount_amount"`
	TaxRate        float64      `json:"tax_rate"`
	TaxAmount      float64      `json:"tax_amount"`
	Amount         float64      `json:"amount"`
	Notes          string       `json:"notes,omitempty"`
	Item           ItemResponse `json:"item"`
}

// SalesOrderCreateRequest 销售订单创建请求
type SalesOrderCreateRequest struct {
	CustomerID      uint                    `json:"customer_id" validate:"required"`
	QuotationID     *uint                   `json:"quotation_id,omitempty"`
	OrderDate       time.Time               `json:"order_date" validate:"required"`
	DeliveryDate    *time.Time              `json:"delivery_date,omitempty"`
	ExpectedDate    time.Time               `json:"expected_date" validate:"required"`
	Status          string                  `json:"status,omitempty"`
	PaymentTerms    string                  `json:"payment_terms,omitempty"`
	ShippingAddress string                  `json:"shipping_address,omitempty"`
	Notes           string                  `json:"notes,omitempty"`
	Items           []SalesOrderItemRequest `json:"items" validate:"required,min=1"`
}

// SalesOrderItemRequest 销售订单项目请求
type SalesOrderItemRequest struct {
	ItemID    uint    `json:"item_id" validate:"required"`
	Quantity  float64 `json:"quantity" validate:"required,gt=0"`
	UnitPrice float64 `json:"unit_price" validate:"required,gt=0,currency"`
	Discount  float64 `json:"discount,omitempty" validate:"min=0,max=100"`
	TaxRate   float64 `json:"tax_rate,omitempty" validate:"min=0,max=100"`
	Notes     string  `json:"notes,omitempty"`
}

// SalesOrderUpdateRequest 销售订单更新请求
type SalesOrderUpdateRequest struct {
	CustomerID      *uint                   `json:"customer_id,omitempty"`
	OrderDate       *time.Time              `json:"order_date,omitempty"`
	DeliveryDate    *time.Time              `json:"delivery_date,omitempty"`
	ExpectedDate    *time.Time              `json:"expected_date,omitempty"`
	Status          *string                 `json:"status,omitempty"`
	PaymentTerms    string                  `json:"payment_terms,omitempty"`
	ShippingAddress string                  `json:"shipping_address,omitempty"`
	Notes           *string                 `json:"notes,omitempty"`
	TotalAmount     *float64                `json:"total_amount,omitempty"`
	Items           []SalesOrderItemRequest `json:"items,omitempty"`
}

// SalesOrderResponse 销售订单响应
type SalesOrderResponse struct {
	ID              uint                     `json:"id"`
	Number          string                   `json:"number"`
	OrderNumber     string                   `json:"orderNumber"` // 前端期望的字段名
	Status          string                   `json:"status"`
	OrderDate       time.Time                `json:"orderDate"`    // 前端期望的订单日期字段
	DeliveryDate    time.Time                `json:"deliveryDate"` // 前端期望的交付日期字段
	ExpectedDate    time.Time                `json:"expected_date"`
	PaymentTerms    string                   `json:"payment_terms,omitempty"`
	ShippingAddress string                   `json:"shipping_address,omitempty"`
	Notes           string                   `json:"notes,omitempty"`
	SubTotal        float64                  `json:"sub_total"`
	DiscountAmount  float64                  `json:"discount_amount"`
	TaxAmount       float64                  `json:"tax_amount"`
	TotalAmount     float64                  `json:"total_amount"`
	Customer        CustomerResponse         `json:"customer"`
	Quotation       *QuotationResponse       `json:"quotation,omitempty"`
	Items           []SalesOrderItemResponse `json:"items"`
	CreatedBy       UserResponse             `json:"created_by"`
	ApprovedBy      *UserResponse            `json:"approved_by,omitempty"`
	CreatedAt       time.Time                `json:"created_at"`
	UpdatedAt       time.Time                `json:"updated_at"`
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
	OrderID        uint                  `json:"order_id" validate:"required"`
	WarehouseID    uint                  `json:"warehouse_id" validate:"required"`
	DeliveryDate   time.Time             `json:"delivery_date" validate:"required"`
	TrackingNumber string                `json:"tracking_number,omitempty"`
	Carrier        string                `json:"carrier,omitempty"`
	Notes          string                `json:"notes,omitempty"`
	Items          []DeliveryItemRequest `json:"items" validate:"required,min=1"`
}

// DeliveryItemRequest 发货项目请求
type DeliveryItemRequest struct {
	OrderItemID  uint    `json:"order_item_id" validate:"required"`
	LocationID   uint    `json:"location_id" validate:"required"`
	DeliveredQty float64 `json:"delivered_qty" validate:"required,gt=0"`
	Notes        string  `json:"notes,omitempty"`
}

// DeliveryResponse 发货响应
type DeliveryResponse struct {
	ID             uint                   `json:"id"`
	Number         string                 `json:"number"`
	DeliveryDate   time.Time              `json:"delivery_date"`
	TrackingNumber string                 `json:"tracking_number,omitempty"`
	Carrier        string                 `json:"carrier,omitempty"`
	Notes          string                 `json:"notes,omitempty"`
	Status         string                 `json:"status"`
	Order          SalesOrderResponse     `json:"order"`
	Warehouse      WarehouseResponse      `json:"warehouse"`
	Items          []DeliveryItemResponse `json:"items"`
	CreatedBy      UserResponse           `json:"created_by"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
}

// DeliveryItemResponse 发货项目响应
type DeliveryItemResponse struct {
	ID           uint                   `json:"id"`
	DeliveredQty float64                `json:"delivered_qty"`
	Notes        string                 `json:"notes,omitempty"`
	OrderItem    SalesOrderItemResponse `json:"order_item"`
	Location     LocationResponse       `json:"location"`
}

// InvoiceCreateRequest 发票创建请求
type InvoiceCreateRequest struct {
	OrderID     uint                 `json:"order_id" validate:"required"`
	InvoiceDate time.Time            `json:"invoice_date" validate:"required"`
	DueDate     time.Time            `json:"due_date" validate:"required"`
	Notes       string               `json:"notes,omitempty"`
	Items       []InvoiceItemRequest `json:"items" validate:"required,min=1"`
}

// InvoiceItemRequest 发票项目请求
type InvoiceItemRequest struct {
	OrderItemID uint    `json:"order_item_id" validate:"required"`
	Quantity    float64 `json:"quantity" validate:"required,gt=0"`
	UnitPrice   float64 `json:"unit_price" validate:"required,gt=0"`
	Notes       string  `json:"notes,omitempty"`
}

// InvoiceResponse 发票响应
type InvoiceResponse struct {
	ID            uint                  `json:"id"`
	Number        string                `json:"number"`
	InvoiceDate   time.Time             `json:"invoice_date"`
	DueDate       time.Time             `json:"due_date"`
	Notes         string                `json:"notes,omitempty"`
	Status        string                `json:"status"`
	SubTotal      float64               `json:"sub_total"`
	TaxAmount     float64               `json:"tax_amount"`
	TotalAmount   float64               `json:"total_amount"`
	PaidAmount    float64               `json:"paid_amount"`
	BalanceAmount float64               `json:"balance_amount"`
	Order         SalesOrderResponse    `json:"order"`
	Items         []InvoiceItemResponse `json:"items"`
	CreatedBy     UserResponse          `json:"created_by"`
	CreatedAt     time.Time             `json:"created_at"`
	UpdatedAt     time.Time             `json:"updated_at"`
}

// InvoiceItemResponse 发票项目响应
type InvoiceItemResponse struct {
	ID        uint                   `json:"id"`
	Quantity  float64                `json:"quantity"`
	UnitPrice float64                `json:"unit_price"`
	Amount    float64                `json:"amount"`
	Notes     string                 `json:"notes,omitempty"`
	OrderItem SalesOrderItemResponse `json:"order_item"`
}

// SalesOrderFilter 销售订单过滤器
type SalesOrderFilter struct {
	SearchRequest
	CustomerID *uint      `json:"customer_id,omitempty" form:"customer_id"`
	Status     string     `json:"status,omitempty" form:"status"`
	StartDate  *time.Time `json:"start_date,omitempty" form:"start_date"`
	EndDate    *time.Time `json:"end_date,omitempty" form:"end_date"`
	MinAmount  *float64   `json:"min_amount,omitempty" form:"min_amount"`
	MaxAmount  *float64   `json:"max_amount,omitempty" form:"max_amount"`
}

// SalesSearchRequest 销售搜索请求
type SalesSearchRequest struct {
	SearchRequest
	CustomerID *uint      `json:"customer_id,omitempty" form:"customer_id"`
	StartDate  *time.Time `json:"start_date,omitempty" form:"start_date"`
	EndDate    *time.Time `json:"end_date,omitempty" form:"end_date"`
	MinAmount  *float64   `json:"min_amount,omitempty" form:"min_amount"`
	MaxAmount  *float64   `json:"max_amount,omitempty" form:"max_amount"`
}

// SalesApprovalRequest 销售审批请求
type SalesApprovalRequest struct {
	Action   string `json:"action" validate:"required,oneof=approve reject"`
	Comments string `json:"comments,omitempty"`
}

// SalesStatisticsResponse 销售统计响应
type SalesStatisticsResponse struct {
	TotalOrders     int64                `json:"total_orders"`
	TotalAmount     float64              `json:"total_amount"`
	PendingOrders   int64                `json:"pending_orders"`
	ApprovedOrders  int64                `json:"approved_orders"`
	CompletedOrders int64                `json:"completed_orders"`
	TopCustomers    []CustomerStatistics `json:"top_customers"`
	TopProducts     []ProductStatistics  `json:"top_products"`
	MonthlyTrend    []MonthlySales       `json:"monthly_trend"`
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
	Format     string `json:"format" form:"format" validate:"required,oneof=excel pdf csv"`
}

// 类型别名，用于兼容service中的命名
type CreateSalesOrderRequest = SalesOrderCreateRequest
type UpdateSalesOrderRequest = SalesOrderUpdateRequest
type CreateSalesOrderItemRequest = SalesOrderItemRequest
type UpdateSalesOrderItemRequest = SalesOrderItemRequest

// QuotationTemplateCreateRequest 报价单模板创建请求
type QuotationTemplateCreateRequest struct {
	Name         string                               `json:"name" validate:"required,max=100"`
	Description  string                               `json:"description,omitempty"`
	IsDefault    bool                                 `json:"is_default,omitempty"`
	ValidityDays int                                  `json:"validity_days" validate:"required,min=1"`
	Terms        string                               `json:"terms,omitempty"`
	Notes        string                               `json:"notes,omitempty"`
	DiscountRate float64                              `json:"discount_rate,omitempty" validate:"min=0,max=100"`
	TaxRate      float64                              `json:"tax_rate,omitempty" validate:"min=0,max=100"`
	CreatedBy    uint                                 `json:"created_by" validate:"required"`
	Items        []QuotationTemplateItemCreateRequest `json:"items" validate:"required,min=1"`
}

// QuotationTemplateUpdateRequest 报价单模板更新请求
type QuotationTemplateUpdateRequest struct {
	Name         string                               `json:"name,omitempty" validate:"omitempty,max=100"`
	Description  string                               `json:"description,omitempty"`
	IsDefault    bool                                 `json:"is_default,omitempty"`
	IsActive     bool                                 `json:"is_active,omitempty"`
	ValidityDays int                                  `json:"validity_days,omitempty" validate:"omitempty,min=1"`
	Terms        string                               `json:"terms,omitempty"`
	Notes        string                               `json:"notes,omitempty"`
	DiscountRate float64                              `json:"discount_rate,omitempty" validate:"min=0,max=100"`
	TaxRate      float64                              `json:"tax_rate,omitempty" validate:"min=0,max=100"`
	Items        []QuotationTemplateItemCreateRequest `json:"items,omitempty"`
}

// QuotationTemplateItemCreateRequest 报价单模板项目创建请求
type QuotationTemplateItemCreateRequest struct {
	ItemID       uint    `json:"item_id" validate:"required"`
	Description  string  `json:"description,omitempty"`
	Quantity     float64 `json:"quantity" validate:"required,gt=0"`
	Rate         float64 `json:"rate" validate:"required,gt=0"`
	DiscountRate float64 `json:"discount_rate,omitempty" validate:"min=0,max=100"`
	TaxRate      float64 `json:"tax_rate,omitempty" validate:"min=0,max=100"`
	SortOrder    int     `json:"sort_order,omitempty"`
}

// QuotationTemplateResponse 报价单模板响应
type QuotationTemplateResponse struct {
	ID           uint                            `json:"id"`
	Name         string                          `json:"name"`
	Description  string                          `json:"description,omitempty"`
	IsDefault    bool                            `json:"is_default"`
	IsActive     bool                            `json:"is_active"`
	ValidityDays int                             `json:"validity_days"`
	Terms        string                          `json:"terms,omitempty"`
	Notes        string                          `json:"notes,omitempty"`
	DiscountRate float64                         `json:"discount_rate"`
	TaxRate      float64                         `json:"tax_rate"`
	Items        []QuotationTemplateItemResponse `json:"items,omitempty"`
	CreatedAt    time.Time                       `json:"created_at"`
	UpdatedAt    time.Time                       `json:"updated_at"`
}

// QuotationTemplateItemResponse 报价单模板项目响应
type QuotationTemplateItemResponse struct {
	ID           uint    `json:"id"`
	ItemID       uint    `json:"item_id"`
	Description  string  `json:"description,omitempty"`
	Quantity     float64 `json:"quantity"`
	Rate         float64 `json:"rate"`
	DiscountRate float64 `json:"discount_rate"`
	TaxRate      float64 `json:"tax_rate"`
	SortOrder    int     `json:"sort_order"`
}

// CreateQuotationFromTemplateRequest 从模板创建报价单请求
type CreateQuotationFromTemplateRequest struct {
	TemplateID uint `json:"template_id" validate:"required"`
	CustomerID uint `json:"customer_id" validate:"required"`
}

// QuotationVersionCreateRequest 创建报价单版本请求
type QuotationVersionCreateRequest struct {
	QuotationID  uint   `json:"quotation_id" validate:"required"`
	VersionName  string `json:"version_name,omitempty" validate:"max=100"`
	ChangeReason string `json:"change_reason,omitempty" validate:"max=500"`
}

// QuotationVersionResponse 报价单版本响应
type QuotationVersionResponse struct {
	ID            uint      `json:"id"`
	QuotationID   uint      `json:"quotation_id"`
	VersionNumber int       `json:"version_number"`
	VersionName   string    `json:"version_name,omitempty"`
	ChangeReason  string    `json:"change_reason,omitempty"`
	IsActive      bool      `json:"is_active"`
	CreatedBy     uint      `json:"created_by"`
	CreatedAt     time.Time `json:"created_at"`
}

// QuotationVersionComparisonResponse 报价单版本比较响应
type QuotationVersionComparisonResponse struct {
	FieldName   string      `json:"field_name"`
	OldValue    interface{} `json:"old_value"`
	NewValue    interface{} `json:"new_value"`
	ChangeType  string      `json:"change_type"` // added, modified, deleted
	Description string      `json:"description,omitempty"`
}

// QuotationVersionHistoryResponse 报价单版本历史响应
type QuotationVersionHistoryResponse struct {
	QuotationID    uint                                 `json:"quotation_id"`
	Versions       []QuotationVersionResponse           `json:"versions"`
	Comparisons    []QuotationVersionComparisonResponse `json:"comparisons,omitempty"`
	TotalVersions  int                                  `json:"total_versions"`
	CurrentVersion int                                  `json:"current_version"`
}

// QuotationVersionCompareRequest 报价单版本比较请求
type QuotationVersionCompareRequest struct {
	QuotationID   uint `json:"quotation_id" validate:"required"`
	FromVersionID uint `json:"from_version_id" validate:"required"`
	ToVersionID   uint `json:"to_version_id" validate:"required"`
}

// QuotationVersionRollbackRequest 报价单版本回滚请求
type QuotationVersionRollbackRequest struct {
	QuotationID uint   `json:"quotation_id" validate:"required"`
	VersionID   uint   `json:"version_id" validate:"required"`
	Reason      string `json:"reason,omitempty" validate:"max=500"`
}
