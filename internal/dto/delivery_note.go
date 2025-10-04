package dto

import (
	"time"
)

// DeliveryNoteCreateRequest 送货单创建请求
type DeliveryNoteCreateRequest struct {
	CustomerID     uint                      `json:"customer_id" validate:"required"`
	SalesOrderID   *uint                     `json:"sales_order_id,omitempty"`
	DeliveryDate   time.Time                 `json:"delivery_date"`
	Date           time.Time                 `json:"date" validate:"required"`
	Reference      string                    `json:"reference,omitempty"`
	Notes          string                    `json:"notes,omitempty"`
	Status         string                    `json:"status" validate:"required,oneof=Draft Submitted Delivered Cancelled"`
	WarehouseID    uint                      `json:"warehouse_id"`
	Transporter    string                    `json:"transporter,omitempty"`
	DriverName     string                    `json:"driver_name,omitempty"`
	VehicleNumber  string                    `json:"vehicle_number,omitempty"`
	Destination    string                    `json:"destination,omitempty"`
	Items          []DeliveryNoteItemRequest `json:"items" validate:"required,min=1"`
}

// DeliveryNoteUpdateRequest 送货单更新请求
type DeliveryNoteUpdateRequest struct {
	CustomerID    *uint                     `json:"customer_id,omitempty"`
	DeliveryDate  *time.Time                `json:"delivery_date,omitempty"`
	Date          *time.Time                `json:"date,omitempty"`
	Reference     string                    `json:"reference,omitempty"`
	Notes         string                    `json:"notes,omitempty"`
	Status        string                    `json:"status" validate:"required,oneof=Draft Submitted Delivered Cancelled"`
	WarehouseID   *uint                     `json:"warehouse_id,omitempty"`
	Transporter   string                    `json:"transporter,omitempty"`
	DriverName    string                    `json:"driver_name,omitempty"`
	VehicleNumber string                    `json:"vehicle_number,omitempty"`
	Destination   string                    `json:"destination,omitempty"`
	Items         []DeliveryNoteItemRequest `json:"items,omitempty"`
}

// DeliveryNoteItemRequest 送货单项目请求
type DeliveryNoteItemRequest struct {
	SalesOrderItemID *uint   `json:"sales_order_item_id,omitempty"`
	ItemID           uint    `json:"item_id" validate:"required"`
	Description      string  `json:"description,omitempty"`
	Quantity         float64 `json:"quantity" validate:"required,gt=0"`
	UOM              string  `json:"uom,omitempty"`
	Rate             float64 `json:"rate,omitempty"`
	Amount           float64 `json:"amount,omitempty"`
	BatchNo          string  `json:"batch_no,omitempty"`
	SerialNo         string  `json:"serial_no,omitempty"`
	WarehouseID      *uint   `json:"warehouse_id,omitempty"`
}

// DeliveryNoteDetailResponse 发货单详细响应
type DeliveryNoteDetailResponse struct {
	ID             uint                       `json:"id"`
	DeliveryNumber string                     `json:"delivery_number"`
	CustomerID     uint                       `json:"customer_id"`
	SalesOrderID   *uint                      `json:"sales_order_id,omitempty"`
	Date           time.Time                  `json:"date"`
	Status         string                     `json:"status"`
	TotalQuantity  float64                    `json:"total_quantity"`
	Transporter    string                     `json:"transporter,omitempty"`
	DriverName     string                     `json:"driver_name,omitempty"`
	VehicleNumber  string                     `json:"vehicle_number,omitempty"`
	Destination    string                     `json:"destination,omitempty"`
	Notes          string                     `json:"notes,omitempty"`
	CreatedBy      uint                       `json:"created_by"`
	Customer       CustomerResponse           `json:"customer"`
	SalesOrder     *SalesOrderResponse        `json:"sales_order,omitempty"`
	Items          []DeliveryNoteItemResponse `json:"items"`
	CreatedAt      time.Time                  `json:"created_at"`
	UpdatedAt      time.Time                  `json:"updated_at"`
}

// DeliveryNoteItemResponse 发货单明细响应
type DeliveryNoteItemResponse struct {
	ID               uint                    `json:"id"`
	DeliveryNoteID   uint                    `json:"delivery_note_id"`
	SalesOrderItemID *uint                   `json:"sales_order_item_id,omitempty"`
	ItemID           uint                    `json:"item_id"`
	Description      string                  `json:"description,omitempty"`
	Quantity         float64                 `json:"quantity"`
	BatchNo          string                  `json:"batch_no,omitempty"`
	SerialNo         string                  `json:"serial_no,omitempty"`
	WarehouseID      *uint                   `json:"warehouse_id,omitempty"`
	SalesOrderItem   *SalesOrderItemResponse `json:"sales_order_item,omitempty"`
	Item             ItemResponse            `json:"item"`
	Warehouse        *WarehouseResponse      `json:"warehouse,omitempty"`
	CreatedAt        time.Time               `json:"created_at"`
	UpdatedAt        time.Time               `json:"updated_at"`
}

// DeliveryNoteListRequest 发货单列表请求
type DeliveryNoteListRequest struct {
	PaginationRequest
	CustomerID   *uint  `json:"customer_id,omitempty" form:"customer_id"`
	SalesOrderID *uint  `json:"sales_order_id,omitempty" form:"sales_order_id"`
	Status       string `json:"status,omitempty" form:"status"`
	DateFrom     string `json:"date_from,omitempty" form:"date_from"`
	DateTo       string `json:"date_to,omitempty" form:"date_to"`
	Search       string `json:"search,omitempty" form:"search"`
}

// DeliveryNoteStatusUpdateRequest 发货单状态更新请求
type DeliveryNoteStatusUpdateRequest struct {
	Status string `json:"status" validate:"required,oneof=Draft Submitted Delivered Cancelled"`
	Notes  string `json:"notes,omitempty"`
}

// DeliveryNoteBatchCreateRequest 批量创建送货单请求
type DeliveryNoteBatchCreateRequest struct {
	SalesOrderID  uint                    `json:"sales_order_id" validate:"required"`
	Date          time.Time               `json:"date" validate:"required"`
	Reference     string                  `json:"reference,omitempty"`
	Notes         string                  `json:"notes,omitempty"`
	Status        string                  `json:"status,omitempty"`
	CustomerID    uint                    `json:"customer_id"`
	WarehouseID   uint                    `json:"warehouse_id"`
	Transporter   string                  `json:"transporter,omitempty"`
	DriverName    string                  `json:"driver_name,omitempty"`
	VehicleNumber string                  `json:"vehicle_number,omitempty"`
	Destination   string                  `json:"destination,omitempty"`
	Items         []DeliveryNoteBatchItem `json:"items" validate:"required,min=1"`
}

type DeliveryNoteBatchItem struct {
	SalesOrderItemID uint    `json:"sales_order_item_id" validate:"required"`
	Quantity         float64 `json:"quantity" validate:"required,gt=0"`
	BatchNo          string  `json:"batch_no,omitempty"`
	SerialNo         string  `json:"serial_no,omitempty"`
	WarehouseID      *uint   `json:"warehouse_id,omitempty"`
}

// DeliveryNoteStatistics 发货单统计
type DeliveryNoteStatistics struct {
	TotalCount      int64   `json:"total_count"`
	DraftCount      int64   `json:"draft_count"`
	SubmittedCount  int64   `json:"submitted_count"`
	DeliveredCount  int64   `json:"delivered_count"`
	CancelledCount  int64   `json:"cancelled_count"`
	TotalQuantity   float64 `json:"total_quantity"`
	MonthlyCount    int64   `json:"monthly_count"`
	MonthlyQuantity float64 `json:"monthly_quantity"`
}

// DeliveryTrendData 发货趋势数据
type DeliveryTrendData struct {
	Date     string  `json:"date"`
	Count    int64   `json:"count"`
	Quantity float64 `json:"quantity"`
}

// DeliveryNoteStatisticsResponse 发货单统计响应
type DeliveryNoteStatisticsResponse struct {
	TotalDeliveries     int64   `json:"total_deliveries"`
	PendingDeliveries   int64   `json:"pending_deliveries"`
	CompletedDeliveries int64   `json:"completed_deliveries"`
	TotalQuantity       float64 `json:"total_quantity"`
	AverageDeliveryTime float64 `json:"average_delivery_time"` // 平均发货时间（天）
}
