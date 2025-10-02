package models

import (
	"time"
)

// Item 物料模型 - 根据数据库结构调整
type Item struct {
	BaseModel
	Code         string  `json:"code" gorm:"uniqueIndex;size:100;not null"`
	Name         string  `json:"name" gorm:"size:255;not null"`
	Description  string  `json:"description,omitempty" gorm:"type:text"`
	Category     string  `json:"category,omitempty" gorm:"size:100"`
	Unit         string  `json:"unit,omitempty" gorm:"size:50"`
	Cost         float64 `json:"cost" gorm:"default:0"`
	Price        float64 `json:"price" gorm:"default:0"`
	ReorderLevel int     `json:"reorder_level" gorm:"default:0"`
	IsActive     bool    `json:"is_active" gorm:"default:true"`

	// 关联
	Stocks    []Stock    `json:"stocks,omitempty" gorm:"foreignKey:ItemID"`
	Movements []Movement `json:"movements,omitempty" gorm:"foreignKey:ItemID"`
}

// 注意：Category和Unit模型已移除，因为数据库中使用字符串字段而非关联表

// Warehouse 仓库模型 - 根据数据库结构调整
type Warehouse struct {
	BaseModel
	Code        string `json:"code" gorm:"uniqueIndex;size:100;not null"`
	Name        string `json:"name" gorm:"size:255;not null"`
	Description string `json:"description,omitempty" gorm:"type:text"`
	Address     string `json:"address,omitempty" gorm:"type:text"`
	IsActive    bool   `json:"is_active" gorm:"default:true"`

	// 关联
	Stocks    []Stock    `json:"stocks,omitempty" gorm:"foreignKey:WarehouseID"`
	Movements []Movement `json:"movements,omitempty" gorm:"foreignKey:WarehouseID"`
}

// 注意：Location模型已移除，因为数据库中没有对应的表

// Stock 库存模型 - 根据数据库结构调整
type Stock struct {
	BaseModel
	ItemID      uint    `json:"item_id" gorm:"index;not null"`
	WarehouseID uint    `json:"warehouse_id" gorm:"index;not null"`
	Quantity    float64 `json:"quantity" gorm:"default:0"`

	// 关联
	Item      Item      `json:"item,omitempty" gorm:"foreignKey:ItemID"`
	Warehouse Warehouse `json:"warehouse,omitempty" gorm:"foreignKey:WarehouseID"`
}

// Movement 库存移动模型 - 根据数据库结构调整
type Movement struct {
	BaseModel
	ItemID        *uint      `json:"item_id,omitempty"`
	WarehouseID   *uint      `json:"warehouse_id,omitempty"`
	Quantity      *float64   `json:"quantity,omitempty"`
	MovementType  string     `json:"movement_type,omitempty"`
	Reference     string     `json:"reference,omitempty"`
	Notes         string     `json:"notes,omitempty"`
	UnitCost      float64    `json:"unit_cost" gorm:"default:0"`
	TotalCost     float64    `json:"total_cost" gorm:"default:0"`
	ReferenceType string     `json:"reference_type,omitempty"`
	ReferenceID   *uint      `json:"reference_id,omitempty"`
	BatchNo       string     `json:"batch_no,omitempty"`
	SerialNo      string     `json:"serial_no,omitempty"`
	ExpiryDate    *time.Time `json:"expiry_date,omitempty"`
	CreatedBy     *uint      `json:"created_by,omitempty"`

	// 关联
	Item      *Item      `json:"item,omitempty" gorm:"foreignKey:ItemID"`
	Warehouse *Warehouse `json:"warehouse,omitempty" gorm:"foreignKey:WarehouseID"`
}

// StockAdjustment 库存调整模型
type StockAdjustment struct {
	AuditableModel
	AdjustmentNumber string    `json:"adjustment_number" gorm:"uniqueIndex;size:100;not null"`
	WarehouseID      uint      `json:"warehouse_id" gorm:"index;not null"`
	AdjustmentDate   time.Time `json:"adjustment_date" gorm:"index;not null"`
	Reason           string    `json:"reason" gorm:"size:255;not null"`
	Notes            string    `json:"notes,omitempty" gorm:"type:text"`
	Status           string    `json:"status" gorm:"size:50;default:'draft';index"` // draft, confirmed, cancelled

	// 关联
	Warehouse Warehouse             `json:"warehouse,omitempty" gorm:"foreignKey:WarehouseID"`
	Items     []StockAdjustmentItem `json:"items,omitempty" gorm:"foreignKey:AdjustmentID"`
}

// StockAdjustmentItem 库存调整明细模型
type StockAdjustmentItem struct {
	BaseModel
	AdjustmentID  uint    `json:"adjustment_id" gorm:"index;not null"`
	ItemID        uint    `json:"item_id" gorm:"index;not null"`
	LocationID    *uint   `json:"location_id,omitempty" gorm:"index"`
	SystemQty     float64 `json:"system_qty" gorm:"not null"`
	ActualQty     float64 `json:"actual_qty" gorm:"not null"`
	DifferenceQty float64 `json:"difference_qty" gorm:"not null"`
	UnitCost      float64 `json:"unit_cost" gorm:"default:0"`
	TotalCost     float64 `json:"total_cost" gorm:"default:0"`
	Notes         string  `json:"notes,omitempty" gorm:"type:text"`

	// 关联
	Adjustment StockAdjustment `json:"adjustment,omitempty" gorm:"foreignKey:AdjustmentID"`
	Item       Item            `json:"item,omitempty" gorm:"foreignKey:ItemID"`
}

// StockMovement 库存移动类型别名，用于兼容repository
type StockMovement = Movement

// StockTransfer 库存调拨模型 - 根据数据库结构调整
type StockTransfer struct {
	BaseModel
	ItemID          *uint    `json:"item_id,omitempty"`
	FromWarehouseID *uint    `json:"from_warehouse_id,omitempty"`
	ToWarehouseID   *uint    `json:"to_warehouse_id,omitempty"`
	Quantity        *float64 `json:"quantity,omitempty"`
	Notes           string   `json:"notes,omitempty"`

	// 关联
	Item          *Item      `json:"item,omitempty" gorm:"foreignKey:ItemID"`
	FromWarehouse *Warehouse `json:"from_warehouse,omitempty" gorm:"foreignKey:FromWarehouseID"`
	ToWarehouse   *Warehouse `json:"to_warehouse,omitempty" gorm:"foreignKey:ToWarehouseID"`
}
