package models

import (
	"time"
)

// SalesInvoice 销售发票模型
type SalesInvoice struct {
	AuditableModel
	InvoiceNumber    string    `json:"invoice_number" gorm:"uniqueIndex;not null"`
	CustomerID       uint      `json:"customer_id" gorm:"not null"`
	SalesOrderID     *uint     `json:"sales_order_id,omitempty"`
	DeliveryNoteID   *uint     `json:"delivery_note_id,omitempty"`
	InvoiceDate      time.Time `json:"invoice_date" gorm:"not null"`
	DueDate          time.Time `json:"due_date" gorm:"not null"`
	PostingDate      time.Time `json:"posting_date" gorm:"not null"`
	
	// 状态管理 - 双重状态系统
	DocStatus        string    `json:"doc_status" gorm:"default:'Draft'"` // Draft, Submitted, Cancelled
	PaymentStatus    string    `json:"payment_status" gorm:"default:'Unpaid'"` // Unpaid, Partially Paid, Paid, Overdue
	
	// 币种和汇率
	Currency         string    `json:"currency" gorm:"default:'CNY'"`
	ExchangeRate     float64   `json:"exchange_rate" gorm:"default:1"`
	
	// 金额字段
	SubTotal         float64   `json:"sub_total" gorm:"default:0"`
	DiscountAmount   float64   `json:"discount_amount" gorm:"default:0"`
	TaxAmount        float64   `json:"tax_amount" gorm:"default:0"`
	ShippingAmount   float64   `json:"shipping_amount" gorm:"default:0"`
	GrandTotal       float64   `json:"grand_total" gorm:"default:0"`
	OutstandingAmount float64  `json:"outstanding_amount" gorm:"default:0"`
	PaidAmount       float64   `json:"paid_amount" gorm:"default:0"`
	
	// 地址信息
	BillingAddress   string    `json:"billing_address,omitempty"`
	ShippingAddress  string    `json:"shipping_address,omitempty"`
	
	// 付款条款
	PaymentTerms     string    `json:"payment_terms,omitempty"`
	PaymentTermsDays int       `json:"payment_terms_days" gorm:"default:30"`
	
	// 销售人员和团队
	SalesPersonID    *uint     `json:"sales_person_id,omitempty"`
	SalesTeamID      *uint     `json:"sales_team_id,omitempty"`
	Territory        string    `json:"territory,omitempty"`
	
	// 其他信息
	CustomerPONumber string    `json:"customer_po_number,omitempty"`
	Project          string    `json:"project,omitempty"`
	CostCenter       string    `json:"cost_center,omitempty"`
	Terms            string    `json:"terms,omitempty"`
	Notes            string    `json:"notes,omitempty"`
	
	// 审计字段
	CreatedBy        uint      `json:"created_by,omitempty"`
	SubmittedBy      *uint     `json:"submitted_by,omitempty"`
	SubmittedAt      *time.Time `json:"submitted_at,omitempty"`
	
	// 关联关系
	Customer         Customer            `json:"customer,omitempty" gorm:"foreignKey:CustomerID"`
	SalesOrder       *SalesOrder         `json:"sales_order,omitempty" gorm:"foreignKey:SalesOrderID"`
	DeliveryNote     *DeliveryNote       `json:"delivery_note,omitempty" gorm:"foreignKey:DeliveryNoteID"`
	SalesPerson      *User               `json:"sales_person,omitempty" gorm:"foreignKey:SalesPersonID"`
	Items            []SalesInvoiceItem  `json:"items,omitempty" gorm:"foreignKey:SalesInvoiceID"`
	Payments         []InvoicePayment    `json:"payments,omitempty" gorm:"foreignKey:SalesInvoiceID"`
	StatusLogs       []InvoiceStatusLog  `json:"status_logs,omitempty" gorm:"foreignKey:SalesInvoiceID"`
}

// SalesInvoiceItem 销售发票明细模型
type SalesInvoiceItem struct {
	AuditableModel
	SalesInvoiceID     uint    `json:"sales_invoice_id" gorm:"not null"`
	SalesOrderItemID   *uint   `json:"sales_order_item_id,omitempty"`
	DeliveryNoteItemID *uint   `json:"delivery_note_item_id,omitempty"`
	ItemID             uint    `json:"item_id" gorm:"not null"`
	ItemCode           string  `json:"item_code" gorm:"not null"`
	ItemName           string  `json:"item_name" gorm:"not null"`
	Description        string  `json:"description,omitempty"`
	
	// 数量和单位
	Quantity           float64 `json:"quantity" gorm:"not null"`
	UOM                string  `json:"uom" gorm:"not null"`
	ConversionFactor   float64 `json:"conversion_factor" gorm:"default:1"`
	StockUOM           string  `json:"stock_uom,omitempty"`
	
	// 价格和金额
	Rate               float64 `json:"rate" gorm:"not null"`
	PriceListRate      float64 `json:"price_list_rate" gorm:"default:0"`
	Amount             float64 `json:"amount" gorm:"not null"`
	
	// 折扣
	DiscountPercentage float64 `json:"discount_percentage" gorm:"default:0"`
	DiscountAmount     float64 `json:"discount_amount" gorm:"default:0"`
	
	// 税费
	TaxCategory        string  `json:"tax_category,omitempty"`
	TaxRate            float64 `json:"tax_rate" gorm:"default:0"`
	TaxAmount          float64 `json:"tax_amount" gorm:"default:0"`
	
	// 总计
	NetRate            float64 `json:"net_rate" gorm:"not null"`
	NetAmount          float64 `json:"net_amount" gorm:"not null"`
	
	// 仓库和批次信息
	WarehouseID        *uint   `json:"warehouse_id,omitempty"`
	BatchNo            string  `json:"batch_no,omitempty"`
	SerialNo           string  `json:"serial_no,omitempty"`
	
	// 会计维度
	CostCenter         string  `json:"cost_center,omitempty"`
	Project            string  `json:"project,omitempty"`
	
	// 关联关系
	SalesInvoice       SalesInvoice      `json:"sales_invoice,omitempty" gorm:"foreignKey:SalesInvoiceID"`
	SalesOrderItem     *SalesOrderItem   `json:"sales_order_item,omitempty" gorm:"foreignKey:SalesOrderItemID"`
	DeliveryNoteItem   *DeliveryNoteItem `json:"delivery_note_item,omitempty" gorm:"foreignKey:DeliveryNoteItemID"`
	Item               Item              `json:"item,omitempty" gorm:"foreignKey:ItemID"`
	Warehouse          *Warehouse        `json:"warehouse,omitempty" gorm:"foreignKey:WarehouseID"`
}

// InvoicePayment 发票付款记录模型
type InvoicePayment struct {
	AuditableModel
	SalesInvoiceID   uint      `json:"sales_invoice_id" gorm:"not null"`
	PaymentEntryID   *uint     `json:"payment_entry_id,omitempty"`
	PaymentDate      time.Time `json:"payment_date" gorm:"not null"`
	PaymentMethod    string    `json:"payment_method" gorm:"not null"`
	Amount           float64   `json:"amount" gorm:"not null"`
	Currency         string    `json:"currency" gorm:"default:'CNY'"`
	ExchangeRate     float64   `json:"exchange_rate" gorm:"default:1"`
	ReferenceNumber  string    `json:"reference_number,omitempty"`
	BankAccountID    *uint     `json:"bank_account_id,omitempty"`
	Notes            string    `json:"notes,omitempty"`
	Status           string    `json:"status" gorm:"default:'Pending'"` // Pending, Cleared, Cancelled
	
	// 关联关系
	SalesInvoice     SalesInvoice  `json:"sales_invoice,omitempty" gorm:"foreignKey:SalesInvoiceID"`
	PaymentEntry     *PaymentEntry `json:"payment_entry,omitempty" gorm:"foreignKey:PaymentEntryID"`
	BankAccount      *BankAccount  `json:"bank_account,omitempty" gorm:"foreignKey:BankAccountID"`
}

// InvoiceStatusLog 发票状态变更日志模型
type InvoiceStatusLog struct {
	AuditableModel
	SalesInvoiceID   uint      `json:"sales_invoice_id" gorm:"not null"`
	FromStatus       string    `json:"from_status,omitempty"`
	ToStatus         string    `json:"to_status" gorm:"not null"`
	StatusType       string    `json:"status_type" gorm:"not null"` // doc_status, payment_status
	ChangedBy        uint      `json:"changed_by" gorm:"not null"`
	ChangedAt        time.Time `json:"changed_at" gorm:"not null"`
	Reason           string    `json:"reason,omitempty"`
	Notes            string    `json:"notes,omitempty"`
	
	// 关联关系
	SalesInvoice     SalesInvoice `json:"sales_invoice,omitempty" gorm:"foreignKey:SalesInvoiceID"`
	User             User         `json:"user,omitempty" gorm:"foreignKey:ChangedBy"`
}

// PricingRule 定价规则模型
type PricingRule struct {
	AuditableModel
	RuleName         string    `json:"rule_name" gorm:"not null"`
	RuleType         string    `json:"rule_type" gorm:"not null"` // Discount, Price
	ApplicableFor    string    `json:"applicable_for" gorm:"not null"` // Item, Item Group, Customer, Customer Group
	Priority         int       `json:"priority" gorm:"default:1"`
	
	// 适用条件
	ItemID           *uint     `json:"item_id,omitempty"`
	ItemGroup        string    `json:"item_group,omitempty"`
	CustomerID       *uint     `json:"customer_id,omitempty"`
	CustomerGroup    string    `json:"customer_group,omitempty"`
	Territory        string    `json:"territory,omitempty"`
	
	// 数量条件
	MinQty           float64   `json:"min_qty" gorm:"default:0"`
	MaxQty           float64   `json:"max_qty" gorm:"default:0"`
	
	// 金额条件
	MinAmount        float64   `json:"min_amount" gorm:"default:0"`
	MaxAmount        float64   `json:"max_amount" gorm:"default:0"`
	
	// 价格/折扣设置
	Rate             float64   `json:"rate" gorm:"default:0"`
	DiscountPercentage float64 `json:"discount_percentage" gorm:"default:0"`
	DiscountAmount   float64   `json:"discount_amount" gorm:"default:0"`
	
	// 有效期
	ValidFrom        time.Time `json:"valid_from" gorm:"not null"`
	ValidUpto        time.Time `json:"valid_upto" gorm:"not null"`
	
	// 状态
	IsActive         bool      `json:"is_active" gorm:"default:true"`
	
	// 关联关系
	Item             *Item     `json:"item,omitempty" gorm:"foreignKey:ItemID"`
	Customer         *Customer `json:"customer,omitempty" gorm:"foreignKey:CustomerID"`
}

// 注意：PaymentEntry 模型已在 accounting.go 中定义，这里直接使用