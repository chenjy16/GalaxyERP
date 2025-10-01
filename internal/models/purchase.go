package models

import (
	"time"
)

// Supplier 供应商模型
type Supplier struct {
	BaseModel
	Name           string  `json:"name" gorm:"not null"`
	Code           string  `json:"code" gorm:"uniqueIndex;not null"`
	Email          string  `json:"email,omitempty"`
	Phone          string  `json:"phone,omitempty"`
	Address        string  `json:"address,omitempty"`
	City           string  `json:"city,omitempty"`
	State          string  `json:"state,omitempty"`
	PostalCode     string  `json:"postal_code,omitempty"`
	Country        string  `json:"country,omitempty"`
	ContactPerson  string  `json:"contact_person,omitempty"`
	CreditLimit    float64 `json:"credit_limit" gorm:"default:0"`
	SupplierGroup  string  `json:"supplier_group,omitempty"`
	Territory      string  `json:"territory,omitempty"`
	QualityRating  float64 `json:"quality_rating" gorm:"default:0"`
	DeliveryRating float64 `json:"delivery_rating" gorm:"default:0"`
	IsActive       bool    `json:"is_active" gorm:"default:true"`

	// 关联关系
	PurchaseOrders []PurchaseOrder `json:"purchase_orders,omitempty" gorm:"foreignKey:SupplierID"`
}

// PurchaseRequest 采购申请模型
type PurchaseRequest struct {
	BaseModel
	RequestNumber string    `json:"request_number" gorm:"uniqueIndex;not null"`
	RequestDate   time.Time `json:"request_date" gorm:"not null"`
	RequiredBy    time.Time `json:"required_by" gorm:"not null"`
	Department    string    `json:"department,omitempty"`
	Status        string    `json:"status" gorm:"default:'Draft'"`
	Notes         string    `json:"notes,omitempty"`
	CreatedBy     uint      `json:"created_by,omitempty"`
	ApprovedBy    *uint     `json:"approved_by,omitempty"`
	
	// 关联
	Items []PurchaseRequestItem `json:"items,omitempty" gorm:"foreignKey:PurchaseRequestID"`
}

// PurchaseRequestItem 采购申请明细模型
type PurchaseRequestItem struct {
	BaseModel
	PurchaseRequestID uint    `json:"purchase_request_id" gorm:"not null"`
	ItemID            uint    `json:"item_id" gorm:"not null"`
	Description       string  `json:"description,omitempty"`
	Quantity          float64 `json:"quantity" gorm:"default:1"`
	UOM               string  `json:"uom,omitempty"`
	EstimatedCost     float64 `json:"estimated_cost" gorm:"default:0"`
	Notes             string  `json:"notes,omitempty"`
	
	// 关联
	PurchaseRequest PurchaseRequest `json:"purchase_request,omitempty" gorm:"foreignKey:PurchaseRequestID"`
	Item            Item            `json:"item,omitempty" gorm:"foreignKey:ItemID"`
}

// PurchaseOrder 采购订单模型
type PurchaseOrder struct {
	BaseModel
	OrderNumber        string    `json:"order_number" gorm:"uniqueIndex;not null"`
	SupplierID         uint      `json:"supplier_id" gorm:"not null"`
	OrderDate          time.Time `json:"order_date" gorm:"not null"`
	DeliveryDate       time.Time `json:"delivery_date" gorm:"not null"`
	Status             string    `json:"status" gorm:"default:'Draft'"`
	PurchaseRequestID  *uint     `json:"purchase_request_id,omitempty"`
	TotalAmount        float64   `json:"total_amount" gorm:"default:0"`
	DiscountAmount     float64   `json:"discount_amount" gorm:"default:0"`
	TaxAmount          float64   `json:"tax_amount" gorm:"default:0"`
	GrandTotal         float64   `json:"grand_total" gorm:"default:0"`
	Terms              string    `json:"terms,omitempty"`
	Notes              string    `json:"notes,omitempty"`
	CreatedBy          uint      `json:"created_by,omitempty"`
	
	// 关联
	Supplier        Supplier             `json:"supplier,omitempty" gorm:"foreignKey:SupplierID"`
	PurchaseRequest *PurchaseRequest     `json:"purchase_request,omitempty" gorm:"foreignKey:PurchaseRequestID"`
	Items           []PurchaseOrderItem  `json:"items,omitempty" gorm:"foreignKey:PurchaseOrderID"`
	Receipts        []PurchaseReceipt    `json:"receipts,omitempty" gorm:"foreignKey:PurchaseOrderID"`
}

// PurchaseOrderItem 采购订单明细模型
type PurchaseOrderItem struct {
	BaseModel
	PurchaseOrderID uint    `json:"purchase_order_id" gorm:"not null"`
	ItemID          uint    `json:"item_id" gorm:"not null"`
	Description     string  `json:"description,omitempty"`
	Quantity        float64 `json:"quantity" gorm:"default:1"`
	ReceivedQty     float64 `json:"received_qty" gorm:"default:0"`
	Rate            float64 `json:"rate" gorm:"default:0"`
	Amount          float64 `json:"amount" gorm:"default:0"`
	DiscountRate    float64 `json:"discount_rate" gorm:"default:0"`
	DiscountAmount  float64 `json:"discount_amount" gorm:"default:0"`
	TaxRate         float64 `json:"tax_rate" gorm:"default:0"`
	TaxAmount       float64 `json:"tax_amount" gorm:"default:0"`
	TotalAmount     float64 `json:"total_amount" gorm:"default:0"`
	WarehouseID     *uint   `json:"warehouse_id,omitempty"`
	
	// 关联
	PurchaseOrder PurchaseOrder `json:"purchase_order,omitempty" gorm:"foreignKey:PurchaseOrderID"`
	Item          Item          `json:"item,omitempty" gorm:"foreignKey:ItemID"`
	Warehouse     *Warehouse    `json:"warehouse,omitempty" gorm:"foreignKey:WarehouseID"`
}

// PurchaseReceipt 采购收货模型
type PurchaseReceipt struct {
	BaseModel
	ReceiptNumber    string    `json:"receipt_number" gorm:"uniqueIndex;not null"`
	SupplierID       uint      `json:"supplier_id" gorm:"not null"`
	PurchaseOrderID  *uint     `json:"purchase_order_id,omitempty"`
	Date             time.Time `json:"date" gorm:"not null"`
	Status           string    `json:"status" gorm:"default:'Draft'"`
	TotalQuantity    float64   `json:"total_quantity" gorm:"default:0"`
	Transporter      string    `json:"transporter,omitempty"`
	DriverName       string    `json:"driver_name,omitempty"`
	VehicleNumber    string    `json:"vehicle_number,omitempty"`
	Destination      string    `json:"destination,omitempty"`
	Notes            string    `json:"notes,omitempty"`
	CreatedBy        uint      `json:"created_by,omitempty"`
	
	// 关联
	Supplier      Supplier               `json:"supplier,omitempty" gorm:"foreignKey:SupplierID"`
	PurchaseOrder *PurchaseOrder         `json:"purchase_order,omitempty" gorm:"foreignKey:PurchaseOrderID"`
	Items         []PurchaseReceiptItem  `json:"items,omitempty" gorm:"foreignKey:PurchaseReceiptID"`
}

// PurchaseReceiptItem 采购收货明细模型
type PurchaseReceiptItem struct {
	BaseModel
	PurchaseReceiptID    uint    `json:"purchase_receipt_id" gorm:"not null"`
	PurchaseOrderItemID  *uint   `json:"purchase_order_item_id,omitempty"`
	ItemID               uint    `json:"item_id" gorm:"not null"`
	Description          string  `json:"description,omitempty"`
	Quantity             float64 `json:"quantity" gorm:"default:1"`
	AcceptedQty          float64 `json:"accepted_qty" gorm:"default:0"`
	RejectedQty          float64 `json:"rejected_qty" gorm:"default:0"`
	BatchNo              string  `json:"batch_no,omitempty"`
	SerialNo             string  `json:"serial_no,omitempty"`
	WarehouseID          *uint   `json:"warehouse_id,omitempty"`
	QualityStatus        string  `json:"quality_status" gorm:"default:'Pending'"`
	Notes                string  `json:"notes,omitempty"`
	
	// 关联
	PurchaseReceipt     PurchaseReceipt    `json:"purchase_receipt,omitempty" gorm:"foreignKey:PurchaseReceiptID"`
	PurchaseOrderItem   *PurchaseOrderItem `json:"purchase_order_item,omitempty" gorm:"foreignKey:PurchaseOrderItemID"`
	Item                Item               `json:"item,omitempty" gorm:"foreignKey:ItemID"`
	Warehouse           *Warehouse         `json:"warehouse,omitempty" gorm:"foreignKey:WarehouseID"`
}