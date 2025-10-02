package dto

import (
	"time"
)

// SalesInvoiceCreateRequest 销售发票创建请求
type SalesInvoiceCreateRequest struct {
	CustomerID       uint                      `json:"customer_id" binding:"required"`
	SalesOrderID     *uint                     `json:"sales_order_id,omitempty"`
	DeliveryNoteID   *uint                     `json:"delivery_note_id,omitempty"`
	InvoiceDate      time.Time                 `json:"invoice_date" binding:"required"`
	DueDate          time.Time                 `json:"due_date" binding:"required"`
	PostingDate      time.Time                 `json:"posting_date" binding:"required"`
	Currency         string                    `json:"currency" binding:"required"`
	ExchangeRate     float64                   `json:"exchange_rate" binding:"min=0"`
	BillingAddress   string                    `json:"billing_address,omitempty"`
	ShippingAddress  string                    `json:"shipping_address,omitempty"`
	PaymentTerms     string                    `json:"payment_terms,omitempty"`
	PaymentTermsDays int                       `json:"payment_terms_days" binding:"min=0"`
	SalesPersonID    *uint                     `json:"sales_person_id,omitempty"`
	Territory        string                    `json:"territory,omitempty"`
	CustomerPONumber string                    `json:"customer_po_number,omitempty"`
	Project          string                    `json:"project,omitempty"`
	CostCenter       string                    `json:"cost_center,omitempty"`
	Terms            string                    `json:"terms,omitempty"`
	Notes            string                    `json:"notes,omitempty"`
	Items            []SalesInvoiceItemRequest `json:"items" binding:"required,min=1"`
}

// SalesInvoiceUpdateRequest 销售发票更新请求
type SalesInvoiceUpdateRequest struct {
	InvoiceDate      *time.Time                `json:"invoice_date,omitempty"`
	DueDate          *time.Time                `json:"due_date,omitempty"`
	PostingDate      *time.Time                `json:"posting_date,omitempty"`
	Currency         string                    `json:"currency,omitempty"`
	ExchangeRate     *float64                  `json:"exchange_rate,omitempty"`
	BillingAddress   string                    `json:"billing_address,omitempty"`
	ShippingAddress  string                    `json:"shipping_address,omitempty"`
	PaymentTerms     string                    `json:"payment_terms,omitempty"`
	PaymentTermsDays *int                      `json:"payment_terms_days,omitempty"`
	SalesPersonID    *uint                     `json:"sales_person_id,omitempty"`
	Territory        string                    `json:"territory,omitempty"`
	CustomerPONumber string                    `json:"customer_po_number,omitempty"`
	Project          string                    `json:"project,omitempty"`
	CostCenter       string                    `json:"cost_center,omitempty"`
	Terms            string                    `json:"terms,omitempty"`
	Notes            string                    `json:"notes,omitempty"`
	Items            []SalesInvoiceItemRequest `json:"items,omitempty"`
}

// SalesInvoiceItemRequest 销售发票明细请求
type SalesInvoiceItemRequest struct {
	SalesOrderItemID   *uint   `json:"sales_order_item_id,omitempty"`
	DeliveryNoteItemID *uint   `json:"delivery_note_item_id,omitempty"`
	ItemID             uint    `json:"item_id" binding:"required"`
	Description        string  `json:"description,omitempty"`
	Quantity           float64 `json:"quantity" binding:"required,gt=0"`
	UOM                string  `json:"uom" binding:"required"`
	Rate               float64 `json:"rate" binding:"required,gte=0"`
	DiscountPercentage float64 `json:"discount_percentage" binding:"gte=0,lte=100"`
	TaxCategory        string  `json:"tax_category,omitempty"`
	TaxRate            float64 `json:"tax_rate" binding:"gte=0,lte=100"`
	WarehouseID        *uint   `json:"warehouse_id,omitempty"`
	BatchNo            string  `json:"batch_no,omitempty"`
	SerialNo           string  `json:"serial_no,omitempty"`
	CostCenter         string  `json:"cost_center,omitempty"`
	Project            string  `json:"project,omitempty"`
}

// SalesInvoiceResponse 销售发票响应
type SalesInvoiceResponse struct {
	ID                uint                       `json:"id"`
	InvoiceNumber     string                     `json:"invoice_number"`
	CustomerID        uint                       `json:"customer_id"`
	SalesOrderID      *uint                      `json:"sales_order_id,omitempty"`
	DeliveryNoteID    *uint                      `json:"delivery_note_id,omitempty"`
	InvoiceDate       time.Time                  `json:"invoice_date"`
	DueDate           time.Time                  `json:"due_date"`
	PostingDate       time.Time                  `json:"posting_date"`
	DocStatus         string                     `json:"doc_status"`
	PaymentStatus     string                     `json:"payment_status"`
	Currency          string                     `json:"currency"`
	ExchangeRate      float64                    `json:"exchange_rate"`
	SubTotal          float64                    `json:"sub_total"`
	DiscountAmount    float64                    `json:"discount_amount"`
	TaxAmount         float64                    `json:"tax_amount"`
	ShippingAmount    float64                    `json:"shipping_amount"`
	GrandTotal        float64                    `json:"grand_total"`
	OutstandingAmount float64                    `json:"outstanding_amount"`
	PaidAmount        float64                    `json:"paid_amount"`
	BillingAddress    string                     `json:"billing_address,omitempty"`
	ShippingAddress   string                     `json:"shipping_address,omitempty"`
	PaymentTerms      string                     `json:"payment_terms,omitempty"`
	PaymentTermsDays  int                        `json:"payment_terms_days"`
	SalesPersonID     *uint                      `json:"sales_person_id,omitempty"`
	Territory         string                     `json:"territory,omitempty"`
	CustomerPONumber  string                     `json:"customer_po_number,omitempty"`
	Project           string                     `json:"project,omitempty"`
	CostCenter        string                     `json:"cost_center,omitempty"`
	Terms             string                     `json:"terms,omitempty"`
	Notes             string                     `json:"notes,omitempty"`
	CreatedBy         uint                       `json:"created_by"`
	SubmittedBy       *uint                      `json:"submitted_by,omitempty"`
	SubmittedAt       *time.Time                 `json:"submitted_at,omitempty"`
	Customer          CustomerResponse           `json:"customer"`
	SalesOrder        *SalesOrderResponse        `json:"sales_order,omitempty"`
	DeliveryNote      *DeliveryNoteResponse      `json:"delivery_note,omitempty"`
	SalesPerson       *UserResponse              `json:"sales_person,omitempty"`
	Items             []SalesInvoiceItemResponse `json:"items"`
	Payments          []InvoicePaymentResponse   `json:"payments,omitempty"`
	StatusLogs        []InvoiceStatusLogResponse `json:"status_logs,omitempty"`
	CreatedAt         time.Time                  `json:"created_at"`
	UpdatedAt         time.Time                  `json:"updated_at"`
}

// SalesInvoiceItemResponse 销售发票明细响应
type SalesInvoiceItemResponse struct {
	ID                 uint               `json:"id"`
	SalesInvoiceID     uint               `json:"sales_invoice_id"`
	SalesOrderItemID   *uint              `json:"sales_order_item_id,omitempty"`
	DeliveryNoteItemID *uint              `json:"delivery_note_item_id,omitempty"`
	ItemID             uint               `json:"item_id"`
	ItemCode           string             `json:"item_code"`
	ItemName           string             `json:"item_name"`
	Description        string             `json:"description,omitempty"`
	Quantity           float64            `json:"quantity"`
	UOM                string             `json:"uom"`
	ConversionFactor   float64            `json:"conversion_factor"`
	StockUOM           string             `json:"stock_uom,omitempty"`
	Rate               float64            `json:"rate"`
	PriceListRate      float64            `json:"price_list_rate"`
	Amount             float64            `json:"amount"`
	DiscountPercentage float64            `json:"discount_percentage"`
	DiscountAmount     float64            `json:"discount_amount"`
	TaxCategory        string             `json:"tax_category,omitempty"`
	TaxRate            float64            `json:"tax_rate"`
	TaxAmount          float64            `json:"tax_amount"`
	NetRate            float64            `json:"net_rate"`
	NetAmount          float64            `json:"net_amount"`
	WarehouseID        *uint              `json:"warehouse_id,omitempty"`
	BatchNo            string             `json:"batch_no,omitempty"`
	SerialNo           string             `json:"serial_no,omitempty"`
	CostCenter         string             `json:"cost_center,omitempty"`
	Project            string             `json:"project,omitempty"`
	Item               ItemResponse       `json:"item"`
	Warehouse          *WarehouseResponse `json:"warehouse,omitempty"`
}

// InvoicePaymentResponse 发票付款响应
type InvoicePaymentResponse struct {
	ID              uint                 `json:"id"`
	SalesInvoiceID  uint                 `json:"sales_invoice_id"`
	PaymentEntryID  *uint                `json:"payment_entry_id,omitempty"`
	PaymentDate     time.Time            `json:"payment_date"`
	PaymentMethod   string               `json:"payment_method"`
	Amount          float64              `json:"amount"`
	Currency        string               `json:"currency"`
	ExchangeRate    float64              `json:"exchange_rate"`
	ReferenceNumber string               `json:"reference_number,omitempty"`
	BankAccountID   *uint                `json:"bank_account_id,omitempty"`
	Notes           string               `json:"notes,omitempty"`
	Status          string               `json:"status"`
	BankAccount     *BankAccountResponse `json:"bank_account,omitempty"`
	CreatedAt       time.Time            `json:"created_at"`
}

// InvoiceStatusLogResponse 发票状态日志响应
type InvoiceStatusLogResponse struct {
	ID             uint         `json:"id"`
	SalesInvoiceID uint         `json:"sales_invoice_id"`
	FromStatus     string       `json:"from_status,omitempty"`
	ToStatus       string       `json:"to_status"`
	StatusType     string       `json:"status_type"`
	ChangedBy      uint         `json:"changed_by"`
	ChangedAt      time.Time    `json:"changed_at"`
	Reason         string       `json:"reason,omitempty"`
	Notes          string       `json:"notes,omitempty"`
	User           UserResponse `json:"user"`
}

// DeliveryNoteResponse 送货单响应（简化版）
type DeliveryNoteResponse struct {
	ID             uint      `json:"id"`
	DeliveryNumber string    `json:"delivery_number"`
	Date           time.Time `json:"date"`
	Status         string    `json:"status"`
}

// BankAccountResponse 银行账户响应（简化版）
type BankAccountResponse struct {
	ID            uint   `json:"id"`
	AccountName   string `json:"account_name"`
	BankName      string `json:"bank_name"`
	AccountNumber string `json:"account_number"`
	Currency      string `json:"currency"`
}

// SalesInvoiceListRequest 销售发票列表请求
type SalesInvoiceListRequest struct {
	PaginationRequest
	CustomerID    *uint  `json:"customer_id,omitempty" form:"customer_id"`
	DocStatus     string `json:"doc_status,omitempty" form:"doc_status"`
	PaymentStatus string `json:"payment_status,omitempty" form:"payment_status"`
	Currency      string `json:"currency,omitempty" form:"currency"`
	SalesPersonID *uint  `json:"sales_person_id,omitempty" form:"sales_person_id"`
	Territory     string `json:"territory,omitempty" form:"territory"`
	DateFrom      string `json:"date_from,omitempty" form:"date_from"`
	DateTo        string `json:"date_to,omitempty" form:"date_to"`
	Search        string `json:"search,omitempty" form:"search"`
}

// SalesInvoiceStatusUpdateRequest 销售发票状态更新请求
type SalesInvoiceStatusUpdateRequest struct {
	Status string `json:"status" binding:"required"`
	Reason string `json:"reason,omitempty"`
	Notes  string `json:"notes,omitempty"`
}

// InvoicePaymentCreateRequest 发票付款创建请求
type InvoicePaymentCreateRequest struct {
	PaymentDate     time.Time `json:"payment_date" binding:"required"`
	PaymentMethod   string    `json:"payment_method" binding:"required"`
	Amount          float64   `json:"amount" binding:"required,gt=0"`
	Currency        string    `json:"currency" binding:"required"`
	ExchangeRate    float64   `json:"exchange_rate" binding:"min=0"`
	ReferenceNumber string    `json:"reference_number,omitempty"`
	BankAccountID   *uint     `json:"bank_account_id,omitempty"`
	Notes           string    `json:"notes,omitempty"`
}

// PricingRuleCreateRequest 定价规则创建请求
type PricingRuleCreateRequest struct {
	RuleName           string    `json:"rule_name" binding:"required"`
	RuleType           string    `json:"rule_type" binding:"required,oneof=Discount Price"`
	ApplicableFor      string    `json:"applicable_for" binding:"required,oneof=Item ItemGroup Customer CustomerGroup"`
	Priority           int       `json:"priority" binding:"min=1"`
	ItemID             *uint     `json:"item_id,omitempty"`
	ItemGroup          string    `json:"item_group,omitempty"`
	CustomerID         *uint     `json:"customer_id,omitempty"`
	CustomerGroup      string    `json:"customer_group,omitempty"`
	Territory          string    `json:"territory,omitempty"`
	MinQty             float64   `json:"min_qty" binding:"gte=0"`
	MaxQty             float64   `json:"max_qty" binding:"gte=0"`
	MinAmount          float64   `json:"min_amount" binding:"gte=0"`
	MaxAmount          float64   `json:"max_amount" binding:"gte=0"`
	Rate               float64   `json:"rate" binding:"gte=0"`
	DiscountPercentage float64   `json:"discount_percentage" binding:"gte=0,lte=100"`
	DiscountAmount     float64   `json:"discount_amount" binding:"gte=0"`
	ValidFrom          time.Time `json:"valid_from" binding:"required"`
	ValidUpto          time.Time `json:"valid_upto" binding:"required"`
}

// PricingRuleResponse 定价规则响应
type PricingRuleResponse struct {
	ID                 uint              `json:"id"`
	RuleName           string            `json:"rule_name"`
	RuleType           string            `json:"rule_type"`
	ApplicableFor      string            `json:"applicable_for"`
	Priority           int               `json:"priority"`
	ItemID             *uint             `json:"item_id,omitempty"`
	ItemGroup          string            `json:"item_group,omitempty"`
	CustomerID         *uint             `json:"customer_id,omitempty"`
	CustomerGroup      string            `json:"customer_group,omitempty"`
	Territory          string            `json:"territory,omitempty"`
	MinQty             float64           `json:"min_qty"`
	MaxQty             float64           `json:"max_qty"`
	MinAmount          float64           `json:"min_amount"`
	MaxAmount          float64           `json:"max_amount"`
	Rate               float64           `json:"rate"`
	DiscountPercentage float64           `json:"discount_percentage"`
	DiscountAmount     float64           `json:"discount_amount"`
	ValidFrom          time.Time         `json:"valid_from"`
	ValidUpto          time.Time         `json:"valid_upto"`
	IsActive           bool              `json:"is_active"`
	Item               *ItemResponse     `json:"item,omitempty"`
	Customer           *CustomerResponse `json:"customer,omitempty"`
	CreatedAt          time.Time         `json:"created_at"`
	UpdatedAt          time.Time         `json:"updated_at"`
}

// SalesInvoiceStatistics 销售发票统计
type SalesInvoiceStatistics struct {
	TotalInvoices     int64   `json:"total_invoices"`
	DraftInvoices     int64   `json:"draft_invoices"`
	SubmittedInvoices int64   `json:"submitted_invoices"`
	PaidInvoices      int64   `json:"paid_invoices"`
	UnpaidInvoices    int64   `json:"unpaid_invoices"`
	OverdueInvoices   int64   `json:"overdue_invoices"`
	TotalAmount       float64 `json:"total_amount"`
	PaidAmount        float64 `json:"paid_amount"`
	OutstandingAmount float64 `json:"outstanding_amount"`
}

// SalesInvoiceReportRequest 销售发票报告请求
type SalesInvoiceReportRequest struct {
	DateRangeRequest
	CustomerID    *uint  `json:"customer_id,omitempty" form:"customer_id"`
	DocStatus     string `json:"doc_status,omitempty" form:"doc_status"`
	PaymentStatus string `json:"payment_status,omitempty" form:"payment_status"`
	SalesPersonID *uint  `json:"sales_person_id,omitempty" form:"sales_person_id"`
	Territory     string `json:"territory,omitempty" form:"territory"`
	Format        string `json:"format" form:"format" binding:"required,oneof=excel pdf csv"`
}
