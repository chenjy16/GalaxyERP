package models

import (
	"time"
)

// Customer 客户模型
type Customer struct {
	BaseModel
	Name          string  `json:"name" gorm:"not null"`
	Code          string  `json:"code" gorm:"uniqueIndex;not null"`
	Email         string  `json:"email,omitempty"`
	Phone         string  `json:"phone,omitempty"`
	Address       string  `json:"address,omitempty"`
	City          string  `json:"city,omitempty"`
	State         string  `json:"state,omitempty"`
	PostalCode    string  `json:"postal_code,omitempty"`
	Country       string  `json:"country,omitempty"`
	ContactPerson string  `json:"contact_person,omitempty"`
	CreditLimit   float64 `json:"credit_limit" gorm:"default:0"`
	CustomerGroup string  `json:"customer_group,omitempty"`
	Territory     string  `json:"territory,omitempty"`
	IsActive      bool    `json:"is_active" gorm:"default:true"`

	// 关联关系
	Quotations  []Quotation  `json:"quotations,omitempty" gorm:"foreignKey:CustomerID"`
	SalesOrders []SalesOrder `json:"salesOrders,omitempty" gorm:"foreignKey:CustomerID"`
}

// Quotation 报价单模型
type Quotation struct {
	BaseModel
	QuotationNumber string    `json:"quotation_number" gorm:"uniqueIndex;not null"`
	CustomerID      uint      `json:"customer_id" gorm:"not null"`
	TemplateID      *uint     `json:"template_id,omitempty"`
	Date            time.Time `json:"date" gorm:"not null"`
	ValidTill       time.Time `json:"valid_till" gorm:"not null"`
	Status          string    `json:"status" gorm:"default:'Draft'"`
	Subject         string    `json:"subject,omitempty"`
	TotalAmount     float64   `json:"total_amount" gorm:"default:0"`
	DiscountAmount  float64   `json:"discount_amount" gorm:"default:0"`
	TaxAmount       float64   `json:"tax_amount" gorm:"default:0"`
	GrandTotal      float64   `json:"grand_total" gorm:"default:0"`
	Terms           string    `json:"terms,omitempty"`
	Notes           string    `json:"notes,omitempty"`
	CreatedBy       uint      `json:"created_by,omitempty"`

	// 版本管理字段
	CurrentVersion int  `json:"current_version" gorm:"default:1"`
	TotalVersions  int  `json:"total_versions" gorm:"default:1"`
	HasVersions    bool `json:"has_versions" gorm:"default:false"`

	// 关联
	Customer Customer           `json:"customer,omitempty" gorm:"foreignKey:CustomerID"`
	Template *QuotationTemplate `json:"template,omitempty" gorm:"foreignKey:TemplateID"`
	Items    []QuotationItem    `json:"items,omitempty" gorm:"foreignKey:QuotationID"`
	Versions []QuotationVersion `json:"versions,omitempty" gorm:"foreignKey:QuotationID"`
}

// QuotationItem 报价单明细模型
type QuotationItem struct {
	BaseModel
	QuotationID    uint    `json:"quotation_id" gorm:"not null"`
	ItemID         uint    `json:"item_id" gorm:"not null"`
	Description    string  `json:"description,omitempty"`
	Quantity       float64 `json:"quantity" gorm:"default:1"`
	Rate           float64 `json:"rate" gorm:"default:0"`
	Amount         float64 `json:"amount" gorm:"default:0"`
	DiscountRate   float64 `json:"discount_rate" gorm:"default:0"`
	DiscountAmount float64 `json:"discount_amount" gorm:"default:0"`
	TaxRate        float64 `json:"tax_rate" gorm:"default:0"`
	TaxAmount      float64 `json:"tax_amount" gorm:"default:0"`
	TotalAmount    float64 `json:"total_amount" gorm:"default:0"`

	// 关联
	Quotation Quotation `json:"quotation,omitempty" gorm:"foreignKey:QuotationID"`
	Item      Item      `json:"item,omitempty" gorm:"foreignKey:ItemID"`
}

// SalesOrder 销售订单模型
type SalesOrder struct {
	BaseModel
	OrderNumber    string    `json:"order_number" gorm:"uniqueIndex;not null"`
	CustomerID     uint      `json:"customer_id" gorm:"not null"`
	Date           time.Time `json:"date" gorm:"not null"`
	DeliveryDate   time.Time `json:"delivery_date" gorm:"not null"`
	Status         string    `json:"status" gorm:"default:'Draft'"`
	QuotationID    *uint     `json:"quotation_id,omitempty"`
	TotalAmount    float64   `json:"total_amount" gorm:"default:0"`
	DiscountAmount float64   `json:"discount_amount" gorm:"default:0"`
	TaxAmount      float64   `json:"tax_amount" gorm:"default:0"`
	GrandTotal     float64   `json:"grand_total" gorm:"default:0"`
	Terms          string    `json:"terms,omitempty"`
	Notes          string    `json:"notes,omitempty"`
	CreatedBy      uint      `json:"created_by,omitempty"`

	// 关联
	Customer      Customer         `json:"customer,omitempty" gorm:"foreignKey:CustomerID"`
	CreatedByUser *User            `json:"created_by_user,omitempty" gorm:"foreignKey:CreatedBy"`
	Quotation     *Quotation       `json:"quotation,omitempty" gorm:"foreignKey:QuotationID"`
	Items         []SalesOrderItem `json:"items,omitempty" gorm:"foreignKey:SalesOrderID"`
}

// SalesOrderItem 销售订单明细模型
type SalesOrderItem struct {
	BaseModel
	SalesOrderID   uint    `json:"sales_order_id" gorm:"not null"`
	ItemID         uint    `json:"item_id" gorm:"not null"`
	Description    string  `json:"description,omitempty"`
	Quantity       float64 `json:"quantity" gorm:"default:1"`
	DeliveredQty   float64 `json:"delivered_qty" gorm:"default:0"`
	Rate           float64 `json:"rate" gorm:"default:0"`
	Amount         float64 `json:"amount" gorm:"default:0"`
	DiscountRate   float64 `json:"discount_rate" gorm:"default:0"`
	DiscountAmount float64 `json:"discount_amount" gorm:"default:0"`
	TaxRate        float64 `json:"tax_rate" gorm:"default:0"`
	TaxAmount      float64 `json:"tax_amount" gorm:"default:0"`
	TotalAmount    float64 `json:"total_amount" gorm:"default:0"`
	WarehouseID    *uint   `json:"warehouse_id,omitempty"`

	// 关联
	SalesOrder SalesOrder `json:"sales_order,omitempty" gorm:"foreignKey:SalesOrderID"`
	Item       Item       `json:"item,omitempty" gorm:"foreignKey:ItemID"`
}

// 注意：Delivery、DeliveryItem、Invoice、InvoiceItem模型已移除，因为数据库中没有对应的表

// DeliveryNote 送货单模型
type DeliveryNote struct {
	BaseModel
	DeliveryNumber string    `json:"delivery_number" gorm:"uniqueIndex;not null"`
	CustomerID     uint      `json:"customer_id" gorm:"not null"`
	SalesOrderID   *uint     `json:"sales_order_id,omitempty"`
	Date           time.Time `json:"date" gorm:"not null"`
	Status         string    `json:"status" gorm:"default:'Draft'"`
	TotalQuantity  float64   `json:"total_quantity" gorm:"default:0"`
	Transporter    string    `json:"transporter,omitempty"`
	DriverName     string    `json:"driver_name,omitempty"`
	VehicleNumber  string    `json:"vehicle_number,omitempty"`
	Destination    string    `json:"destination,omitempty"`
	Notes          string    `json:"notes,omitempty"`
	CreatedBy      uint      `json:"created_by,omitempty"`

	// 关联
	Customer   Customer           `json:"customer,omitempty" gorm:"foreignKey:CustomerID"`
	SalesOrder *SalesOrder        `json:"sales_order,omitempty" gorm:"foreignKey:SalesOrderID"`
	Items      []DeliveryNoteItem `json:"items,omitempty" gorm:"foreignKey:DeliveryNoteID"`
}

// DeliveryNoteItem 送货单明细模型
type DeliveryNoteItem struct {
	BaseModel
	DeliveryNoteID   uint    `json:"delivery_note_id" gorm:"not null"`
	SalesOrderItemID *uint   `json:"sales_order_item_id,omitempty"`
	ItemID           uint    `json:"item_id" gorm:"not null"`
	Description      string  `json:"description,omitempty"`
	Quantity         float64 `json:"quantity" gorm:"not null"`
	BatchNo          string  `json:"batch_no,omitempty"`
	SerialNo         string  `json:"serial_no,omitempty"`
	WarehouseID      *uint   `json:"warehouse_id,omitempty"`

	// 关联
	DeliveryNote   DeliveryNote    `json:"delivery_note,omitempty" gorm:"foreignKey:DeliveryNoteID"`
	SalesOrderItem *SalesOrderItem `json:"sales_order_item,omitempty" gorm:"foreignKey:SalesOrderItemID"`
	Item           Item            `json:"item,omitempty" gorm:"foreignKey:ItemID"`
	Warehouse      *Warehouse      `json:"warehouse,omitempty" gorm:"foreignKey:WarehouseID"`
}

// QuotationTemplate 报价单模板模型
type QuotationTemplate struct {
	BaseModel
	Name         string  `json:"name" gorm:"not null"`
	Description  string  `json:"description,omitempty"`
	IsDefault    bool    `json:"is_default" gorm:"default:false"`
	IsActive     bool    `json:"is_active" gorm:"default:true"`
	ValidityDays int     `json:"validity_days" gorm:"default:30"`
	Terms        string  `json:"terms,omitempty"`
	Notes        string  `json:"notes,omitempty"`
	DiscountRate float64 `json:"discount_rate" gorm:"default:0"`
	TaxRate      float64 `json:"tax_rate" gorm:"default:0"`
	CreatedBy    uint    `json:"created_by,omitempty"`

	// 关联
	Items []QuotationTemplateItem `json:"items,omitempty" gorm:"foreignKey:TemplateID"`
}

// QuotationTemplateItem 报价单模板项目模型
type QuotationTemplateItem struct {
	BaseModel
	TemplateID   uint    `json:"template_id" gorm:"not null"`
	ItemID       uint    `json:"item_id" gorm:"not null"`
	Description  string  `json:"description,omitempty"`
	Quantity     float64 `json:"quantity" gorm:"default:1"`
	Rate         float64 `json:"rate" gorm:"default:0"`
	DiscountRate float64 `json:"discount_rate" gorm:"default:0"`
	TaxRate      float64 `json:"tax_rate" gorm:"default:0"`
	SortOrder    int     `json:"sort_order" gorm:"default:0"`

	// 关联
	Template QuotationTemplate `json:"template,omitempty" gorm:"foreignKey:TemplateID"`
	Item     Item              `json:"item,omitempty" gorm:"foreignKey:ItemID"`
}

// QuotationVersion 报价单版本模型
type QuotationVersion struct {
	BaseModel
	QuotationID   uint   `json:"quotation_id" gorm:"not null"`
	VersionNumber int    `json:"version_number" gorm:"not null"`
	VersionName   string `json:"version_name,omitempty"`
	VersionData   string `json:"version_data" gorm:"type:text"` // JSON格式存储版本数据
	ChangeReason  string `json:"change_reason,omitempty"`
	IsActive      bool   `json:"is_active" gorm:"default:false"`
	CreatedBy     uint   `json:"created_by,omitempty"`

	// 关联
	Quotation Quotation `json:"quotation,omitempty" gorm:"foreignKey:QuotationID"`
}

// QuotationVersionComparison 报价单版本比较结果
type QuotationVersionComparison struct {
	FieldName   string      `json:"field_name"`
	OldValue    interface{} `json:"old_value"`
	NewValue    interface{} `json:"new_value"`
	ChangeType  string      `json:"change_type"` // added, modified, deleted
	Description string      `json:"description,omitempty"`
}

// QuotationVersionHistory 报价单版本历史记录
type QuotationVersionHistory struct {
	QuotationID    uint                         `json:"quotation_id"`
	Versions       []QuotationVersion           `json:"versions"`
	Comparisons    []QuotationVersionComparison `json:"comparisons,omitempty"`
	TotalVersions  int                          `json:"total_versions"`
	CurrentVersion int                          `json:"current_version"`
}
