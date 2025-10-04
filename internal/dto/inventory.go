package dto

import (
	"time"
)

// ItemCreateRequest 物料创建请求
type ItemCreateRequest struct {
	Code        string  `json:"code" validate:"required,max=50"`
	Name        string  `json:"name" validate:"required,max=100"`
	Description string  `json:"description,omitempty"`
	CategoryID  uint    `json:"category_id" validate:"required"`
	UnitID      uint    `json:"unit_id" validate:"required"`
	Type        string  `json:"type" validate:"required,oneof=raw_material finished_goods semi_finished consumable"`
	MinStock    float64 `json:"min_stock,omitempty" validate:"min=0"`
	MaxStock    float64 `json:"max_stock,omitempty" validate:"min=0"`
	UnitCost    float64 `json:"unit_cost,omitempty" validate:"min=0"`
	SalePrice   float64 `json:"sale_price,omitempty" validate:"min=0"`
	Barcode     string  `json:"barcode,omitempty"`
	ImageURL    string  `json:"image_url,omitempty"`
}

// ItemUpdateRequest 物料更新请求
type ItemUpdateRequest struct {
	Name        string   `json:"name,omitempty" validate:"omitempty,max=100"`
	Description string   `json:"description,omitempty"`
	CategoryID  *uint    `json:"category_id,omitempty"`
	UnitID      *uint    `json:"unit_id,omitempty"`
	MinStock    *float64 `json:"min_stock,omitempty" validate:"omitempty,min=0"`
	MaxStock    *float64 `json:"max_stock,omitempty" validate:"omitempty,min=0"`
	UnitCost    *float64 `json:"unit_cost,omitempty" validate:"omitempty,min=0"`
	SalePrice   *float64 `json:"sale_price,omitempty" validate:"omitempty,min=0"`
	Barcode     string   `json:"barcode,omitempty"`
	ImageURL    string   `json:"image_url,omitempty"`
	IsActive    *bool    `json:"is_active,omitempty"`
}

// ItemResponse 物料响应
type ItemResponse struct {
	ID          uint             `json:"id"`
	Code        string           `json:"code"`
	Name        string           `json:"name"`
	Description string           `json:"description,omitempty"`
	Type        string           `json:"type"`
	MinStock    float64          `json:"min_stock"`
	MaxStock    float64          `json:"max_stock"`
	UnitCost    float64          `json:"unit_cost"`
	SalePrice   float64          `json:"sale_price"`
	Barcode     string           `json:"barcode,omitempty"`
	ImageURL    string           `json:"image_url,omitempty"`
	IsActive    bool             `json:"is_active"`
	Category    CategoryResponse `json:"category"`
	Unit        UnitResponse     `json:"unit"`
	Stock       []StockResponse  `json:"stock,omitempty"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

// ItemListResponse 物料列表响应
type ItemListResponse struct {
	ID         uint    `json:"id"`
	Code       string  `json:"code"`
	Name       string  `json:"name"`
	Type       string  `json:"type"`
	Category   string  `json:"category"`
	Unit       string  `json:"unit"`
	UnitCost   float64 `json:"unit_cost"`
	SalePrice  float64 `json:"sale_price"`
	TotalStock float64 `json:"total_stock"`
	IsActive   bool    `json:"is_active"`
}



// CategoryCreateRequest 分类创建请求
type CategoryCreateRequest struct {
	Name        string `json:"name" validate:"required,max=100"`
	Code        string `json:"code" validate:"required,max=50"`
	Description string `json:"description,omitempty"`
	ParentID    *uint  `json:"parent_id,omitempty"`
	IsActive    bool   `json:"is_active"`
}

// CategoryUpdateRequest 分类更新请求
type CategoryUpdateRequest struct {
	Name        string `json:"name,omitempty" validate:"omitempty,max=100"`
	Description string `json:"description,omitempty"`
	ParentID    *uint  `json:"parent_id,omitempty"`
	IsActive    *bool  `json:"is_active,omitempty"`
}

// CategoryResponse 分类响应
type CategoryResponse struct {
	ID          uint               `json:"id"`
	Name        string             `json:"name"`
	Code        string             `json:"code"`
	Description string             `json:"description,omitempty"`
	IsActive    bool               `json:"is_active"`
	Parent      *CategoryResponse  `json:"parent,omitempty"`
	Children    []CategoryResponse `json:"children,omitempty"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
}

// UnitCreateRequest 单位创建请求
type UnitCreateRequest struct {
	Name         string `json:"name" validate:"required,max=50"`
	Symbol       string `json:"symbol" validate:"required,max=10"`
	Description  string `json:"description,omitempty"`
	BaseUnitID   *uint  `json:"base_unit_id,omitempty"`
	ConversionRate float64 `json:"conversion_rate,omitempty" validate:"omitempty,min=0"`
	IsActive     bool   `json:"is_active"`
}

// UnitUpdateRequest 单位更新请求
type UnitUpdateRequest struct {
	Name         string   `json:"name,omitempty" validate:"omitempty,max=50"`
	Symbol       string   `json:"symbol,omitempty" validate:"omitempty,max=10"`
	Description  string   `json:"description,omitempty"`
	BaseUnitID   *uint    `json:"base_unit_id,omitempty"`
	ConversionRate *float64 `json:"conversion_rate,omitempty" validate:"omitempty,min=0"`
	IsActive     *bool    `json:"is_active,omitempty"`
}

// UnitResponse 单位响应
type UnitResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Symbol      string    `json:"symbol"`
	Description string    `json:"description,omitempty"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// WarehouseCreateRequest 仓库创建请求
type WarehouseCreateRequest struct {
	Name        string `json:"name" validate:"required,max=100"`
	Code        string `json:"code" validate:"required,max=50"`
	Description string `json:"description,omitempty"`
	Address     string `json:"address,omitempty"`
	ManagerID   *uint  `json:"manager_id,omitempty"`
	IsActive    bool   `json:"is_active"`
}

// WarehouseUpdateRequest 仓库更新请求
type WarehouseUpdateRequest struct {
	Name        string `json:"name,omitempty" validate:"omitempty,max=100"`
	Description string `json:"description,omitempty"`
	Address     string `json:"address,omitempty"`
	ManagerID   *uint  `json:"manager_id,omitempty"`
	IsActive    *bool  `json:"is_active,omitempty"`
}

// WarehouseResponse 仓库响应
type WarehouseResponse struct {
	ID          uint               `json:"id"`
	Name        string             `json:"name"`
	Code        string             `json:"code"`
	Address     string             `json:"address,omitempty"`
	Description string             `json:"description,omitempty"`
	IsActive    bool               `json:"is_active"`
	Manager     *UserResponse      `json:"manager,omitempty"`
	Locations   []LocationResponse `json:"locations,omitempty"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
}

// LocationCreateRequest 库位创建请求
type LocationCreateRequest struct {
	Name        string `json:"name" validate:"required,max=100"`
	Code        string `json:"code" validate:"required,max=50"`
	WarehouseID uint   `json:"warehouse_id" validate:"required"`
	Type        string `json:"type" validate:"required,oneof=storage picking shipping receiving"`
	Description string `json:"description,omitempty"`
}

// LocationUpdateRequest 库位更新请求
type LocationUpdateRequest struct {
	Name        string `json:"name,omitempty" validate:"omitempty,max=100"`
	Type        string `json:"type,omitempty" validate:"omitempty,oneof=storage picking shipping receiving"`
	Description string `json:"description,omitempty"`
	IsActive    *bool  `json:"is_active,omitempty"`
}

// LocationResponse 库位响应
type LocationResponse struct {
	ID          uint              `json:"id"`
	Name        string            `json:"name"`
	Code        string            `json:"code"`
	Type        string            `json:"type"`
	Description string            `json:"description,omitempty"`
	IsActive    bool              `json:"is_active"`
	Warehouse   WarehouseResponse `json:"warehouse"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

// StockCreateRequest 库存创建请求
type StockCreateRequest struct {
	ItemID      uint    `json:"item_id" validate:"required"`
	WarehouseID uint    `json:"warehouse_id" validate:"required"`
	Quantity    float64 `json:"quantity" validate:"required,min=0"`
}

// StockUpdateRequest 库存更新请求
type StockUpdateRequest struct {
	Quantity *float64 `json:"quantity,omitempty" validate:"omitempty,min=0"`
}

// StockResponse 库存响应
type StockResponse struct {
	ID           uint              `json:"id"`
	Quantity     float64           `json:"quantity"`
	ReservedQty  float64           `json:"reserved_qty"`
	AvailableQty float64           `json:"available_qty"`
	Item         ItemResponse      `json:"item"`
	Warehouse    WarehouseResponse `json:"warehouse"`
	Location     LocationResponse  `json:"location"`
	UpdatedAt    time.Time         `json:"updated_at"`
}

// MovementCreateRequest 库存移动创建请求
type MovementCreateRequest struct {
	ItemID      uint    `json:"item_id" validate:"required"`
	WarehouseID uint    `json:"warehouse_id" validate:"required"`
	LocationID  uint    `json:"location_id" validate:"required"`
	Type        string  `json:"type" validate:"required,oneof=in out transfer adjustment"`
	Quantity    float64 `json:"quantity" validate:"required,gt=0"`
	Reference   string  `json:"reference,omitempty"`
	Notes       string  `json:"notes,omitempty"`
}

// MovementResponse 库存移动响应
type MovementResponse struct {
	ID        uint              `json:"id"`
	Type      string            `json:"type"`
	Quantity  float64           `json:"quantity"`
	Reference string            `json:"reference,omitempty"`
	Notes     string            `json:"notes,omitempty"`
	Item      ItemResponse      `json:"item"`
	Warehouse WarehouseResponse `json:"warehouse"`
	Location  LocationResponse  `json:"location"`
	CreatedBy UserResponse      `json:"created_by"`
	CreatedAt time.Time         `json:"created_at"`
}

// StockAdjustmentCreateRequest 库存调整创建请求
type StockAdjustmentCreateRequest struct {
	WarehouseID uint                         `json:"warehouse_id" validate:"required"`
	Reason      string                       `json:"reason" validate:"required"`
	Notes       string                       `json:"notes,omitempty"`
	Items       []StockAdjustmentItemRequest `json:"items" validate:"required,min=1"`
}

// StockAdjustmentItemRequest 库存调整项目请求
type StockAdjustmentItemRequest struct {
	ItemID     uint    `json:"item_id" validate:"required"`
	LocationID uint    `json:"location_id" validate:"required"`
	SystemQty  float64 `json:"system_qty" validate:"required,min=0"`
	ActualQty  float64 `json:"actual_qty" validate:"required,min=0"`
	Reason     string  `json:"reason,omitempty"`
}

// StockAdjustmentResponse 库存调整响应
type StockAdjustmentResponse struct {
	ID        uint                          `json:"id"`
	Number    string                        `json:"number"`
	Reason    string                        `json:"reason"`
	Notes     string                        `json:"notes,omitempty"`
	Status    string                        `json:"status"`
	Warehouse WarehouseResponse             `json:"warehouse"`
	Items     []StockAdjustmentItemResponse `json:"items"`
	CreatedBy UserResponse                  `json:"created_by"`
	CreatedAt time.Time                     `json:"created_at"`
	UpdatedAt time.Time                     `json:"updated_at"`
}

// StockAdjustmentItemResponse 库存调整项目响应
type StockAdjustmentItemResponse struct {
	ID            uint             `json:"id"`
	SystemQty     float64          `json:"system_qty"`
	ActualQty     float64          `json:"actual_qty"`
	DifferenceQty float64          `json:"difference_qty"`
	Reason        string           `json:"reason,omitempty"`
	Item          ItemResponse     `json:"item"`
	Location      LocationResponse `json:"location"`
}

// ItemSearchRequest 物料搜索请求
type ItemSearchRequest struct {
	SearchRequest
	CategoryID *uint    `json:"category_id,omitempty" form:"category_id"`
	Type       string   `json:"type,omitempty" form:"type"`
	IsActive   *bool    `json:"is_active,omitempty" form:"is_active"`
	MinPrice   *float64 `json:"min_price,omitempty" form:"min_price"`
	MaxPrice   *float64 `json:"max_price,omitempty" form:"max_price"`
}

// StockSearchRequest 库存搜索请求
type StockSearchRequest struct {
	SearchRequest
	WarehouseID *uint `json:"warehouse_id,omitempty" form:"warehouse_id"`
	LocationID  *uint `json:"location_id,omitempty" form:"location_id"`
	ItemID      *uint `json:"item_id,omitempty" form:"item_id"`
	LowStock    *bool `json:"low_stock,omitempty" form:"low_stock"`
}

// StockReportRequest 库存报告请求
type StockReportRequest struct {
	WarehouseID *uint     `json:"warehouse_id,omitempty" form:"warehouse_id"`
	CategoryID  *uint     `json:"category_id,omitempty" form:"category_id"`
	AsOfDate    time.Time `json:"as_of_date,omitempty" form:"as_of_date"`
	Format      string    `json:"format" form:"format" validate:"required,oneof=excel pdf csv"`
}
