package models

import (
	"time"
)

// Product 产品模型
type Product struct {
	BaseModel
	Code        string  `json:"code" gorm:"uniqueIndex;size:100;not null"`
	Name        string  `json:"name" gorm:"size:255;not null"`
	Description string  `json:"description" gorm:"type:text"`
	Category    string  `json:"category" gorm:"size:100;not null;index"`
	Unit        string  `json:"unit" gorm:"size:50;not null"`
	Price       float64 `json:"price" gorm:"default:0"`
	Cost        float64 `json:"cost" gorm:"default:0"`
	Status      string  `json:"status" gorm:"size:50;default:'active';index"`
}

// BOM 物料清单模型
type BOM struct {
	AuditableModel
	ProductID     uint       `json:"product_id" gorm:"index;not null"`
	Version       string     `json:"version" gorm:"size:50;not null"`
	EffectiveDate time.Time  `json:"effective_date" gorm:"index;not null"`
	ExpiryDate    *time.Time `json:"expiry_date,omitempty" gorm:"index"`
	Quantity      float64    `json:"quantity" gorm:"not null"` // 产出数量
	IsActive      bool       `json:"is_active" gorm:"default:true;index"`
	Notes         string     `json:"notes,omitempty" gorm:"type:text"`

	// 关联
	Product Product   `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	Items   []BOMItem `json:"items,omitempty" gorm:"foreignKey:BOMID"`
}

// BOMItem 物料清单明细模型
type BOMItem struct {
	BaseModel
	BOMID     uint    `json:"bom_id" gorm:"index;not null"`
	ItemID    uint    `json:"item_id" gorm:"index;not null"`
	Quantity  float64 `json:"quantity" gorm:"not null"`
	UnitCost  float64 `json:"unit_cost,omitempty" gorm:"default:0"`
	TotalCost float64 `json:"total_cost,omitempty" gorm:"default:0"`
	ScrapRate float64 `json:"scrap_rate,omitempty" gorm:"default:0"` // 损耗率
	Notes     string  `json:"notes,omitempty" gorm:"type:text"`

	// 关联
	BOM  BOM  `json:"bom,omitempty" gorm:"foreignKey:BOMID"`
	Item Item `json:"item,omitempty" gorm:"foreignKey:ItemID"`
}

// WorkOrder 生产工单模型
type WorkOrder struct {
	AuditableModel
	WorkOrderNumber string     `json:"work_order_number" gorm:"uniqueIndex;size:100;not null"`
	ProductID       uint       `json:"product_id" gorm:"index;not null"`
	BOMID           uint       `json:"bom_id" gorm:"index;not null"`
	PlannedQty      float64    `json:"planned_qty" gorm:"not null"`
	ProducedQty     float64    `json:"produced_qty" gorm:"default:0"`
	ScrapQty        float64    `json:"scrap_qty" gorm:"default:0"`
	StartDate       time.Time  `json:"start_date" gorm:"index;not null"`
	EndDate         time.Time  `json:"end_date" gorm:"index;not null"`
	ActualStartDate *time.Time `json:"actual_start_date,omitempty"`
	ActualEndDate   *time.Time `json:"actual_end_date,omitempty"`
	Priority        string     `json:"priority" gorm:"size:20;default:'normal';index"` // low, normal, high, urgent
	Status          string     `json:"status" gorm:"size:50;default:'planned';index"`  // planned, released, in_progress, completed, cancelled
	Notes           string     `json:"notes,omitempty" gorm:"type:text"`

	// 关联
	Product    Product              `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	BOM        BOM                  `json:"bom,omitempty" gorm:"foreignKey:BOMID"`
	Operations []WorkOrderOperation `json:"operations,omitempty" gorm:"foreignKey:WorkOrderID"`
	Materials  []WorkOrderMaterial  `json:"materials,omitempty" gorm:"foreignKey:WorkOrderID"`
}

// WorkOrderOperation 生产工单工序模型
type WorkOrderOperation struct {
	BaseModel
	WorkOrderID      uint       `json:"work_order_id" gorm:"index;not null"`
	OperationID      uint       `json:"operation_id" gorm:"index;not null"`
	Sequence         int        `json:"sequence" gorm:"not null"`
	PlannedStartDate time.Time  `json:"planned_start_date" gorm:"index;not null"`
	PlannedEndDate   time.Time  `json:"planned_end_date" gorm:"index;not null"`
	ActualStartDate  *time.Time `json:"actual_start_date,omitempty"`
	ActualEndDate    *time.Time `json:"actual_end_date,omitempty"`
	PlannedHours     float64    `json:"planned_hours" gorm:"not null"`
	ActualHours      float64    `json:"actual_hours" gorm:"default:0"`
	Status           string     `json:"status" gorm:"size:50;default:'pending';index"` // pending, in_progress, completed, skipped
	Notes            string     `json:"notes,omitempty" gorm:"type:text"`

	// 关联
	WorkOrder WorkOrder `json:"work_order,omitempty" gorm:"foreignKey:WorkOrderID"`
	Operation Operation `json:"operation,omitempty" gorm:"foreignKey:OperationID"`
}

// Operation 工序模型
type Operation struct {
	CodeModel
	WorkCenterID uint    `json:"work_center_id" gorm:"index;not null"`
	SetupTime    float64 `json:"setup_time" gorm:"default:0"`    // 准备时间（小时）
	RunTime      float64 `json:"run_time" gorm:"default:0"`      // 运行时间（小时/单位）
	StandardCost float64 `json:"standard_cost" gorm:"default:0"` // 标准成本（元/小时）

	// 关联
	WorkCenter          WorkCenter           `json:"work_center,omitempty" gorm:"foreignKey:WorkCenterID"`
	WorkOrderOperations []WorkOrderOperation `json:"work_order_operations,omitempty" gorm:"foreignKey:OperationID"`
}

// WorkCenter 工作中心模型
type WorkCenter struct {
	CodeModel
	WorkCenterType string  `json:"work_center_type" gorm:"size:50;not null;index"` // machine, manual, assembly
	Capacity       float64 `json:"capacity" gorm:"default:8"`                      // 产能（小时/天）
	Efficiency     float64 `json:"efficiency" gorm:"default:100"`                  // 效率（%）
	CostPerHour    float64 `json:"cost_per_hour" gorm:"default:0"`                 // 每小时成本

	// 关联
	Operations []Operation `json:"operations,omitempty" gorm:"foreignKey:WorkCenterID"`
}

// WorkOrderMaterial 生产工单物料模型
type WorkOrderMaterial struct {
	BaseModel
	WorkOrderID uint    `json:"work_order_id" gorm:"index;not null"`
	ItemID      uint    `json:"item_id" gorm:"index;not null"`
	RequiredQty float64 `json:"required_qty" gorm:"not null"`
	IssuedQty   float64 `json:"issued_qty" gorm:"default:0"`
	ConsumedQty float64 `json:"consumed_qty" gorm:"default:0"`
	ReturnedQty float64 `json:"returned_qty" gorm:"default:0"`
	UnitCost    float64 `json:"unit_cost" gorm:"default:0"`
	TotalCost   float64 `json:"total_cost" gorm:"default:0"`
	Status      string  `json:"status" gorm:"size:50;default:'pending';index"` // pending, issued, consumed, returned

	// 关联
	WorkOrder WorkOrder `json:"work_order,omitempty" gorm:"foreignKey:WorkOrderID"`
	Item      Item      `json:"item,omitempty" gorm:"foreignKey:ItemID"`
}

// QualityCheck 质量检验模型
type QualityCheck struct {
	AuditableModel
	CheckNumber string    `json:"check_number" gorm:"uniqueIndex;size:100;not null"`
	CheckDate   time.Time `json:"check_date" gorm:"index;not null"`
	CheckType   string    `json:"check_type" gorm:"size:50;not null;index"` // incoming, in_process, final, outgoing
	ItemID      uint      `json:"item_id" gorm:"index;not null"`
	BatchNumber string    `json:"batch_number,omitempty" gorm:"size:100;index"`
	CheckedQty  float64   `json:"checked_qty" gorm:"not null"`
	PassedQty   float64   `json:"passed_qty" gorm:"default:0"`
	FailedQty   float64   `json:"failed_qty" gorm:"default:0"`
	ReworkQty   float64   `json:"rework_qty" gorm:"default:0"`
	InspectorID uint      `json:"inspector_id" gorm:"index;not null"`
	Status      string    `json:"status" gorm:"size:50;default:'pending';index"` // pending, passed, failed, rework
	Notes       string    `json:"notes,omitempty" gorm:"type:text"`

	// 关联
	Item      Item                 `json:"item,omitempty" gorm:"foreignKey:ItemID"`
	Inspector Employee             `json:"inspector,omitempty" gorm:"foreignKey:InspectorID"`
	Details   []QualityCheckDetail `json:"details,omitempty" gorm:"foreignKey:CheckID"`
}

// QualityCheckDetail 质量检验明细模型
type QualityCheckDetail struct {
	BaseModel
	CheckID     uint   `json:"check_id" gorm:"index;not null"`
	CheckPoint  string `json:"check_point" gorm:"size:255;not null"`
	Standard    string `json:"standard" gorm:"type:text"`
	ActualValue string `json:"actual_value" gorm:"type:text"`
	Result      string `json:"result" gorm:"size:50;not null;index"` // pass, fail, na
	Notes       string `json:"notes,omitempty" gorm:"type:text"`

	// 关联
	Check QualityCheck `json:"check,omitempty" gorm:"foreignKey:CheckID"`
}

// ProductionPlan 生产计划模型 - 根据数据库结构调整
type ProductionPlan struct {
	BaseModel
	Name        string     `json:"name" gorm:"size:255;not null"`
	Description string     `json:"description,omitempty"`
	PlanType    string     `json:"plan_type" gorm:"size:50;not null"`
	StartDate   time.Time  `json:"start_date" gorm:"index;not null"`
	EndDate     time.Time  `json:"end_date" gorm:"index;not null"`
	Status      string     `json:"status" gorm:"size:50;default:'draft'"`
	CreatedBy   uint       `json:"created_by" gorm:"index;not null"`
	ApprovedBy  *uint      `json:"approved_by,omitempty"`
	ApprovedAt  *time.Time `json:"approved_at,omitempty"`

	// 关联
	MaterialRequirements []MaterialRequirement `json:"material_requirements,omitempty" gorm:"foreignKey:PlanID"`
}

// MaterialRequirement 物料需求模型 - 根据数据库结构调整
type MaterialRequirement struct {
	BaseModel
	PlanID            uint      `json:"plan_id" gorm:"index;not null"`
	ItemID            uint      `json:"item_id" gorm:"index;not null"`
	RequiredQuantity  float64   `json:"required_quantity" gorm:"not null"`
	AvailableQuantity float64   `json:"available_quantity" gorm:"default:0"`
	NetRequirement    float64   `json:"net_requirement" gorm:"not null"`
	DueDate           time.Time `json:"due_date" gorm:"index;not null"`
	Priority          int       `json:"priority" gorm:"default:0"`

	// 关联
	Plan *ProductionPlan `json:"plan,omitempty" gorm:"foreignKey:PlanID"`
	Item *Item           `json:"item,omitempty" gorm:"foreignKey:ItemID"`
}

// ProcessRoute 工艺路线模型 - 根据数据库结构调整
type ProcessRoute struct {
	BaseModel
	Name        string `json:"name" gorm:"size:255;not null"`
	Code        string `json:"code" gorm:"uniqueIndex;size:100;not null"`
	Description string `json:"description,omitempty"`
	ItemID      uint   `json:"item_id" gorm:"index;not null"`
	IsActive    bool   `json:"is_active" gorm:"default:true"`

	// 关联
	Item       *Item              `json:"item,omitempty" gorm:"foreignKey:ItemID"`
	Operations []ProcessOperation `json:"operations,omitempty" gorm:"foreignKey:RouteID"`
}

// ProcessOperation 工艺工序模型 - 根据数据库结构调整
type ProcessOperation struct {
	BaseModel
	RouteID         uint    `json:"route_id" gorm:"index;not null"`
	OperationNumber int     `json:"operation_number" gorm:"not null"`
	Name            string  `json:"name" gorm:"size:255;not null"`
	Description     string  `json:"description,omitempty"`
	WorkCenterID    *uint   `json:"work_center_id,omitempty"`
	StandardHours   float64 `json:"standard_hours" gorm:"default:0"`
	SetupHours      float64 `json:"setup_hours" gorm:"default:0"`
	WaitHours       float64 `json:"wait_hours" gorm:"default:0"`
	MoveHours       float64 `json:"move_hours" gorm:"default:0"`
	Sequence        int     `json:"sequence" gorm:"not null"`

	// 关联
	Route      *ProcessRoute `json:"route,omitempty" gorm:"foreignKey:RouteID"`
	WorkCenter *WorkCenter   `json:"work_center,omitempty" gorm:"foreignKey:WorkCenterID"`
}

// ProductionOrder 生产订单模型 - 根据数据库结构调整
type ProductionOrder struct {
	BaseModel
	OrderNumber      string     `json:"order_number" gorm:"uniqueIndex;size:100;not null"`
	Name             string     `json:"name" gorm:"size:255;not null"`
	Description      string     `json:"description,omitempty"`
	ItemID           uint       `json:"item_id" gorm:"index;not null"`
	Quantity         float64    `json:"quantity" gorm:"not null"`
	ProducedQuantity float64    `json:"produced_quantity" gorm:"default:0"`
	Unit             string     `json:"unit,omitempty"`
	StartDate        time.Time  `json:"start_date" gorm:"index;not null"`
	EndDate          time.Time  `json:"end_date" gorm:"index;not null"`
	ActualStartDate  *time.Time `json:"actual_start_date,omitempty"`
	ActualEndDate    *time.Time `json:"actual_end_date,omitempty"`
	Status           string     `json:"status" gorm:"size:50;default:'draft';index"`
	Priority         int        `json:"priority" gorm:"default:0"`
	RouteID          *uint      `json:"route_id,omitempty"`
	CreatedBy        uint       `json:"created_by" gorm:"index;not null"`
	ApprovedBy       *uint      `json:"approved_by,omitempty"`
	ApprovedAt       *time.Time `json:"approved_at,omitempty"`

	// 关联
	Item     *Item                 `json:"item,omitempty" gorm:"foreignKey:ItemID"`
	Route    *ProcessRoute         `json:"route,omitempty" gorm:"foreignKey:RouteID"`
	Items    []ProductionOrderItem `json:"items,omitempty" gorm:"foreignKey:OrderID"`
	Progress []ProductionProgress  `json:"progress,omitempty" gorm:"foreignKey:OrderID"`
}

// ProductionOrderItem 生产订单明细模型 - 根据数据库结构调整
type ProductionOrderItem struct {
	BaseModel
	OrderID          uint    `json:"order_id" gorm:"index;not null"`
	ItemID           uint    `json:"item_id" gorm:"index;not null"`
	RequiredQuantity float64 `json:"required_quantity" gorm:"not null"`
	IssuedQuantity   float64 `json:"issued_quantity" gorm:"default:0"`
	ConsumedQuantity float64 `json:"consumed_quantity" gorm:"default:0"`
	Unit             string  `json:"unit,omitempty"`

	// 关联
	Order *ProductionOrder `json:"order,omitempty" gorm:"foreignKey:OrderID"`
	Item  *Item            `json:"item,omitempty" gorm:"foreignKey:ItemID"`
}

// ProductionProgress 生产进度模型
// ProductionProgress 生产进度模型 - 根据数据库结构调整
type ProductionProgress struct {
	BaseModel
	OrderID          uint       `json:"order_id" gorm:"index;not null"`
	OperationID      *uint      `json:"operation_id,omitempty"`
	WorkCenterID     *uint      `json:"work_center_id,omitempty"`
	StartDate        time.Time  `json:"start_date" gorm:"index;not null"`
	EndDate          time.Time  `json:"end_date" gorm:"index;not null"`
	ActualStartDate  *time.Time `json:"actual_start_date,omitempty"`
	ActualEndDate    *time.Time `json:"actual_end_date,omitempty"`
	Quantity         float64    `json:"quantity" gorm:"not null"`
	ProducedQuantity float64    `json:"produced_quantity" gorm:"default:0"`
	Status           string     `json:"status" gorm:"size:50;default:'pending'"`
	Notes            string     `json:"notes,omitempty"`

	// 关联
	Order      *ProductionOrder  `json:"order,omitempty" gorm:"foreignKey:OrderID"`
	Operation  *ProcessOperation `json:"operation,omitempty" gorm:"foreignKey:OperationID"`
	WorkCenter *WorkCenter       `json:"work_center,omitempty" gorm:"foreignKey:WorkCenterID"`
}

// QualityInspection 质量检验模型
type QualityInspection struct {
	AuditableModel
	InspectionNumber  string    `json:"inspection_number" gorm:"uniqueIndex;size:100;not null"`
	ProductionOrderID *uint     `json:"production_order_id,omitempty" gorm:"index"`
	ItemID            uint      `json:"item_id" gorm:"index;not null"`
	BatchNumber       string    `json:"batch_number" gorm:"size:100;index"`
	InspectionDate    time.Time `json:"inspection_date" gorm:"index;not null"`
	InspectionType    string    `json:"inspection_type" gorm:"size:50;not null;index"` // incoming, in_process, final, random
	SampleSize        int       `json:"sample_size" gorm:"not null"`
	PassedQuantity    int       `json:"passed_quantity" gorm:"default:0"`
	FailedQuantity    int       `json:"failed_quantity" gorm:"default:0"`
	Result            string    `json:"result" gorm:"size:50;default:'pending';index"` // pending, passed, failed, conditional
	InspectorID       uint      `json:"inspector_id" gorm:"index;not null"`
	Notes             string    `json:"notes,omitempty" gorm:"type:text"`

	// 关联
	ProductionOrder *ProductionOrder `json:"production_order,omitempty" gorm:"foreignKey:ProductionOrderID"`
	Item            *Item            `json:"item,omitempty" gorm:"foreignKey:ItemID"`
}

// NonConformingProduct 不合格品模型
type NonConformingProduct struct {
	AuditableModel
	NCRNumber         string  `json:"ncr_number" gorm:"uniqueIndex;size:100;not null"`
	InspectionID      *uint   `json:"inspection_id,omitempty" gorm:"index"`
	ProductionOrderID *uint   `json:"production_order_id,omitempty" gorm:"index"`
	ItemID            uint    `json:"item_id" gorm:"index;not null"`
	Quantity          float64 `json:"quantity" gorm:"not null"`
	DefectType        string  `json:"defect_type" gorm:"size:100;not null;index"`
	DefectDescription string  `json:"defect_description" gorm:"type:text;not null"`
	Disposition       string  `json:"disposition" gorm:"size:50;default:'pending';index"` // pending, rework, scrap, return, use_as_is
	ResponsibleParty  string  `json:"responsible_party" gorm:"size:100"`
	CorrectiveAction  string  `json:"corrective_action" gorm:"type:text"`
	Status            string  `json:"status" gorm:"size:50;default:'open';index"` // open, in_progress, closed

	// 关联
	Inspection      *QualityInspection `json:"inspection,omitempty" gorm:"foreignKey:InspectionID"`
	ProductionOrder *ProductionOrder   `json:"production_order,omitempty" gorm:"foreignKey:ProductionOrderID"`
	Item            *Item              `json:"item,omitempty" gorm:"foreignKey:ItemID"`
}

// Equipment 设备模型
type Equipment struct {
	CodeModel
	EquipmentNumber string    `json:"equipment_number" gorm:"uniqueIndex;size:100;not null"`
	EquipmentType   string    `json:"equipment_type" gorm:"size:100;not null;index"`
	Manufacturer    string    `json:"manufacturer" gorm:"size:100"`
	Model           string    `json:"model" gorm:"size:100"`
	SerialNumber    string    `json:"serial_number" gorm:"size:100"`
	PurchaseDate    time.Time `json:"purchase_date" gorm:"index"`
	InstallDate     time.Time `json:"install_date" gorm:"index"`
	WarrantyExpiry  time.Time `json:"warranty_expiry" gorm:"index"`
	Location        string    `json:"location" gorm:"size:255"`
	WorkCenterID    *uint     `json:"work_center_id,omitempty" gorm:"index"`
	Capacity        float64   `json:"capacity" gorm:"default:0"`
	Status          string    `json:"status" gorm:"size:50;default:'active';index"` // active, maintenance, breakdown, retired

	// 关联
	WorkCenter  *WorkCenter            `json:"work_center,omitempty" gorm:"foreignKey:WorkCenterID"`
	Maintenance []EquipmentMaintenance `json:"maintenance,omitempty" gorm:"foreignKey:EquipmentID"`
}

// EquipmentMaintenance 设备维护模型
type EquipmentMaintenance struct {
	AuditableModel
	MaintenanceNumber string     `json:"maintenance_number" gorm:"uniqueIndex;size:100;not null"`
	EquipmentID       uint       `json:"equipment_id" gorm:"index;not null"`
	MaintenanceType   string     `json:"maintenance_type" gorm:"size:50;not null;index"` // preventive, corrective, emergency
	ScheduledDate     time.Time  `json:"scheduled_date" gorm:"index;not null"`
	ActualDate        *time.Time `json:"actual_date,omitempty" gorm:"index"`
	Duration          float64    `json:"duration" gorm:"default:0"` // 维护时长（小时）
	Cost              float64    `json:"cost" gorm:"default:0"`
	TechnicianID      uint       `json:"technician_id" gorm:"index;not null"`
	Description       string     `json:"description" gorm:"type:text;not null"`
	PartsUsed         string     `json:"parts_used" gorm:"type:text"`
	Status            string     `json:"status" gorm:"size:50;default:'scheduled';index"` // scheduled, in_progress, completed, cancelled

	// 关联
	Equipment *Equipment `json:"equipment,omitempty" gorm:"foreignKey:EquipmentID"`
}

// EquipmentFailure 设备故障模型
type EquipmentFailure struct {
	AuditableModel
	FailureNumber string     `json:"failure_number" gorm:"uniqueIndex;size:100;not null"`
	EquipmentID   uint       `json:"equipment_id" gorm:"index;not null"`
	FailureDate   time.Time  `json:"failure_date" gorm:"index;not null"`
	FailureType   string     `json:"failure_type" gorm:"size:100;not null;index"` // mechanical, electrical, software, other
	Severity      string     `json:"severity" gorm:"size:50;not null;index"`      // low, medium, high, critical
	Description   string     `json:"description" gorm:"type:text;not null"`
	CauseAnalysis string     `json:"cause_analysis" gorm:"type:text"`
	RepairAction  string     `json:"repair_action" gorm:"type:text"`
	RepairCost    float64    `json:"repair_cost" gorm:"default:0"`
	DowntimeHours float64    `json:"downtime_hours" gorm:"default:0"`
	TechnicianID  *uint      `json:"technician_id,omitempty" gorm:"index"`
	ReportedBy    uint       `json:"reported_by" gorm:"index;not null"`
	ResolvedBy    *uint      `json:"resolved_by,omitempty" gorm:"index"`
	ResolvedAt    *time.Time `json:"resolved_at,omitempty" gorm:"index"`
	Status        string     `json:"status" gorm:"size:50;default:'reported';index"` // reported, investigating, repairing, resolved, closed

	// 关联
	Equipment *Equipment `json:"equipment,omitempty" gorm:"foreignKey:EquipmentID"`
}
