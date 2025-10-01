-- ============================================================================
-- GalaxyERP 生产管理模块 - PostgreSQL 建表脚本
-- 生成时间: 2025-01-01
-- 说明: 基于 /internal/models/production.go 的结构，包含BOM、工艺路线、生产订单管理
-- ============================================================================

-- ============================================================================
-- 物料清单(BOM)管理
-- ============================================================================

-- BOM表
CREATE TABLE IF NOT EXISTS boms (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL,
  created_by BIGINT NULL,
  updated_by BIGINT NULL,
  is_active BOOLEAN DEFAULT TRUE,
  code VARCHAR(50) NOT NULL,
  name VARCHAR(255) NOT NULL,
  description TEXT NULL,
  item_id BIGINT NOT NULL,
  version VARCHAR(20) DEFAULT '1.0',
  quantity DECIMAL(15,4) DEFAULT 1,
  unit_id BIGINT NOT NULL,
  bom_type VARCHAR(20) DEFAULT 'PRODUCTION',
  effective_date DATE NOT NULL,
  expiry_date DATE NULL,
  is_default BOOLEAN DEFAULT FALSE,
  status VARCHAR(20) DEFAULT 'DRAFT',
  approved_by BIGINT NULL,
  approved_at TIMESTAMP NULL,
  CONSTRAINT uq_boms_code UNIQUE (code),
  CONSTRAINT fk_boms_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE,
  CONSTRAINT fk_boms_unit FOREIGN KEY (unit_id) REFERENCES units(id) ON DELETE CASCADE,
  CONSTRAINT fk_boms_approved_by FOREIGN KEY (approved_by) REFERENCES employees(id) ON DELETE SET NULL,
  CONSTRAINT chk_boms_bom_type CHECK (bom_type IN ('PRODUCTION', 'ENGINEERING', 'SALES', 'COSTING'))
);
CREATE INDEX IF NOT EXISTS idx_boms_deleted_at ON boms (deleted_at);
CREATE INDEX IF NOT EXISTS idx_boms_is_active ON boms (is_active);
CREATE INDEX IF NOT EXISTS idx_boms_code ON boms (code);
CREATE INDEX IF NOT EXISTS idx_boms_item_id ON boms (item_id);
CREATE INDEX IF NOT EXISTS idx_boms_bom_type ON boms (bom_type);
CREATE INDEX IF NOT EXISTS idx_boms_status ON boms (status);
CREATE INDEX IF NOT EXISTS idx_boms_effective_date ON boms (effective_date);
CREATE INDEX IF NOT EXISTS idx_boms_is_default ON boms (is_default);
CREATE INDEX IF NOT EXISTS idx_boms_name ON boms (name);

-- BOM明细表
CREATE TABLE IF NOT EXISTS bom_items (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL,
  bom_id BIGINT NOT NULL,
  item_id BIGINT NOT NULL,
  sequence INTEGER DEFAULT 1,
  quantity DECIMAL(15,4) NOT NULL,
  unit_id BIGINT NOT NULL,
  scrap_rate DECIMAL(5,2) DEFAULT 0,
  yield_rate DECIMAL(5,2) DEFAULT 100,
  is_phantom BOOLEAN DEFAULT FALSE,
  is_optional BOOLEAN DEFAULT FALSE,
  effective_date DATE NOT NULL,
  expiry_date DATE NULL,
  notes TEXT NULL,
  CONSTRAINT fk_bom_items_bom FOREIGN KEY (bom_id) REFERENCES boms(id) ON DELETE CASCADE,
  CONSTRAINT fk_bom_items_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE,
  CONSTRAINT fk_bom_items_unit FOREIGN KEY (unit_id) REFERENCES units(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_bom_items_deleted_at ON bom_items (deleted_at);
CREATE INDEX IF NOT EXISTS idx_bom_items_bom_id ON bom_items (bom_id);
CREATE INDEX IF NOT EXISTS idx_bom_items_item_id ON bom_items (item_id);
CREATE INDEX IF NOT EXISTS idx_bom_items_sequence ON bom_items (sequence);
CREATE INDEX IF NOT EXISTS idx_bom_items_effective_date ON bom_items (effective_date);
CREATE INDEX IF NOT EXISTS idx_bom_items_is_phantom ON bom_items (is_phantom);

-- ============================================================================
-- 工艺路线管理
-- ============================================================================

-- 工作中心表
CREATE TABLE IF NOT EXISTS work_centers (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL,
  created_by BIGINT NULL,
  updated_by BIGINT NULL,
  is_active BOOLEAN DEFAULT TRUE,
  code VARCHAR(50) NOT NULL,
  name VARCHAR(255) NOT NULL,
  description TEXT NULL,
  work_center_type VARCHAR(20) DEFAULT 'MACHINE',
  department_id BIGINT NULL,
  capacity_per_hour DECIMAL(10,2) DEFAULT 1,
  efficiency_rate DECIMAL(5,2) DEFAULT 100,
  utilization_rate DECIMAL(5,2) DEFAULT 80,
  setup_time_minutes INTEGER DEFAULT 0,
  teardown_time_minutes INTEGER DEFAULT 0,
  cost_per_hour DECIMAL(10,2) DEFAULT 0,
  is_bottleneck BOOLEAN DEFAULT FALSE,
  status VARCHAR(20) DEFAULT 'ACTIVE',
  CONSTRAINT uq_work_centers_code UNIQUE (code),
  CONSTRAINT fk_work_centers_department FOREIGN KEY (department_id) REFERENCES departments(id) ON DELETE SET NULL,
  CONSTRAINT chk_work_centers_type CHECK (work_center_type IN ('MACHINE', 'LABOR', 'ASSEMBLY', 'INSPECTION', 'OTHER'))
);
CREATE INDEX IF NOT EXISTS idx_work_centers_deleted_at ON work_centers (deleted_at);
CREATE INDEX IF NOT EXISTS idx_work_centers_is_active ON work_centers (is_active);
CREATE INDEX IF NOT EXISTS idx_work_centers_code ON work_centers (code);
CREATE INDEX IF NOT EXISTS idx_work_centers_work_center_type ON work_centers (work_center_type);
CREATE INDEX IF NOT EXISTS idx_work_centers_department_id ON work_centers (department_id);
CREATE INDEX IF NOT EXISTS idx_work_centers_status ON work_centers (status);
CREATE INDEX IF NOT EXISTS idx_work_centers_is_bottleneck ON work_centers (is_bottleneck);
CREATE INDEX IF NOT EXISTS idx_work_centers_name ON work_centers (name);

-- 工艺路线表
CREATE TABLE IF NOT EXISTS routings (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL,
  created_by BIGINT NULL,
  updated_by BIGINT NULL,
  is_active BOOLEAN DEFAULT TRUE,
  code VARCHAR(50) NOT NULL,
  name VARCHAR(255) NOT NULL,
  description TEXT NULL,
  item_id BIGINT NOT NULL,
  version VARCHAR(20) DEFAULT '1.0',
  routing_type VARCHAR(20) DEFAULT 'PRODUCTION',
  effective_date DATE NOT NULL,
  expiry_date DATE NULL,
  is_default BOOLEAN DEFAULT FALSE,
  total_setup_time INTEGER DEFAULT 0,
  total_run_time INTEGER DEFAULT 0,
  status VARCHAR(20) DEFAULT 'DRAFT',
  approved_by BIGINT NULL,
  approved_at TIMESTAMP NULL,
  CONSTRAINT uq_routings_code UNIQUE (code),
  CONSTRAINT fk_routings_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE,
  CONSTRAINT fk_routings_approved_by FOREIGN KEY (approved_by) REFERENCES employees(id) ON DELETE SET NULL,
  CONSTRAINT chk_routings_routing_type CHECK (routing_type IN ('PRODUCTION', 'ENGINEERING', 'ALTERNATE'))
);
CREATE INDEX IF NOT EXISTS idx_routings_deleted_at ON routings (deleted_at);
CREATE INDEX IF NOT EXISTS idx_routings_is_active ON routings (is_active);
CREATE INDEX IF NOT EXISTS idx_routings_code ON routings (code);
CREATE INDEX IF NOT EXISTS idx_routings_item_id ON routings (item_id);
CREATE INDEX IF NOT EXISTS idx_routings_routing_type ON routings (routing_type);
CREATE INDEX IF NOT EXISTS idx_routings_status ON routings (status);
CREATE INDEX IF NOT EXISTS idx_routings_effective_date ON routings (effective_date);
CREATE INDEX IF NOT EXISTS idx_routings_is_default ON routings (is_default);
CREATE INDEX IF NOT EXISTS idx_routings_name ON routings (name);

-- 工艺路线操作表
CREATE TABLE IF NOT EXISTS routing_operations (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL,
  routing_id BIGINT NOT NULL,
  operation_number INTEGER NOT NULL,
  operation_name VARCHAR(255) NOT NULL,
  description TEXT NULL,
  work_center_id BIGINT NOT NULL,
  setup_time_minutes INTEGER DEFAULT 0,
  run_time_minutes DECIMAL(10,2) DEFAULT 0,
  teardown_time_minutes INTEGER DEFAULT 0,
  move_time_minutes INTEGER DEFAULT 0,
  queue_time_minutes INTEGER DEFAULT 0,
  overlap_percentage DECIMAL(5,2) DEFAULT 0,
  is_subcontract BOOLEAN DEFAULT FALSE,
  cost_per_unit DECIMAL(10,2) DEFAULT 0,
  notes TEXT NULL,
  CONSTRAINT fk_routing_operations_routing FOREIGN KEY (routing_id) REFERENCES routings(id) ON DELETE CASCADE,
  CONSTRAINT fk_routing_operations_work_center FOREIGN KEY (work_center_id) REFERENCES work_centers(id) ON DELETE CASCADE,
  CONSTRAINT uq_routing_operations_routing_number UNIQUE (routing_id, operation_number)
);
CREATE INDEX IF NOT EXISTS idx_routing_operations_deleted_at ON routing_operations (deleted_at);
CREATE INDEX IF NOT EXISTS idx_routing_operations_routing_id ON routing_operations (routing_id);
CREATE INDEX IF NOT EXISTS idx_routing_operations_work_center_id ON routing_operations (work_center_id);
CREATE INDEX IF NOT EXISTS idx_routing_operations_operation_number ON routing_operations (operation_number);
CREATE INDEX IF NOT EXISTS idx_routing_operations_is_subcontract ON routing_operations (is_subcontract);

-- ============================================================================
-- 生产计划管理
-- ============================================================================

-- 生产计划表
CREATE TABLE IF NOT EXISTS production_plans (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL,
  created_by BIGINT NULL,
  updated_by BIGINT NULL,
  is_active BOOLEAN DEFAULT TRUE,
  code VARCHAR(50) NOT NULL,
  name VARCHAR(255) NOT NULL,
  description TEXT NULL,
  plan_type VARCHAR(20) DEFAULT 'MASTER',
  plan_period VARCHAR(20) DEFAULT 'MONTHLY',
  start_date DATE NOT NULL,
  end_date DATE NOT NULL,
  planner_id BIGINT NULL,
  status VARCHAR(20) DEFAULT 'DRAFT',
  approved_by BIGINT NULL,
  approved_at TIMESTAMP NULL,
  CONSTRAINT uq_production_plans_code UNIQUE (code),
  CONSTRAINT fk_production_plans_planner FOREIGN KEY (planner_id) REFERENCES employees(id) ON DELETE SET NULL,
  CONSTRAINT fk_production_plans_approved_by FOREIGN KEY (approved_by) REFERENCES employees(id) ON DELETE SET NULL,
  CONSTRAINT chk_production_plans_plan_type CHECK (plan_type IN ('MASTER', 'DETAILED', 'CAPACITY')),
  CONSTRAINT chk_production_plans_plan_period CHECK (plan_period IN ('DAILY', 'WEEKLY', 'MONTHLY', 'QUARTERLY'))
);
CREATE INDEX IF NOT EXISTS idx_production_plans_deleted_at ON production_plans (deleted_at);
CREATE INDEX IF NOT EXISTS idx_production_plans_is_active ON production_plans (is_active);
CREATE INDEX IF NOT EXISTS idx_production_plans_code ON production_plans (code);
CREATE INDEX IF NOT EXISTS idx_production_plans_plan_type ON production_plans (plan_type);
CREATE INDEX IF NOT EXISTS idx_production_plans_plan_period ON production_plans (plan_period);
CREATE INDEX IF NOT EXISTS idx_production_plans_start_date ON production_plans (start_date);
CREATE INDEX IF NOT EXISTS idx_production_plans_end_date ON production_plans (end_date);
CREATE INDEX IF NOT EXISTS idx_production_plans_status ON production_plans (status);
CREATE INDEX IF NOT EXISTS idx_production_plans_planner_id ON production_plans (planner_id);

-- 生产计划明细表
CREATE TABLE IF NOT EXISTS production_plan_items (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL,
  production_plan_id BIGINT NOT NULL,
  item_id BIGINT NOT NULL,
  planned_quantity DECIMAL(15,4) NOT NULL,
  unit_id BIGINT NOT NULL,
  planned_start_date DATE NOT NULL,
  planned_end_date DATE NOT NULL,
  priority INTEGER DEFAULT 5,
  sales_order_id BIGINT NULL,
  notes TEXT NULL,
  CONSTRAINT fk_production_plan_items_production_plan FOREIGN KEY (production_plan_id) REFERENCES production_plans(id) ON DELETE CASCADE,
  CONSTRAINT fk_production_plan_items_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE,
  CONSTRAINT fk_production_plan_items_unit FOREIGN KEY (unit_id) REFERENCES units(id) ON DELETE CASCADE,
  CONSTRAINT fk_production_plan_items_sales_order FOREIGN KEY (sales_order_id) REFERENCES sales_orders(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_production_plan_items_deleted_at ON production_plan_items (deleted_at);
CREATE INDEX IF NOT EXISTS idx_production_plan_items_production_plan_id ON production_plan_items (production_plan_id);
CREATE INDEX IF NOT EXISTS idx_production_plan_items_item_id ON production_plan_items (item_id);
CREATE INDEX IF NOT EXISTS idx_production_plan_items_planned_start_date ON production_plan_items (planned_start_date);
CREATE INDEX IF NOT EXISTS idx_production_plan_items_planned_end_date ON production_plan_items (planned_end_date);
CREATE INDEX IF NOT EXISTS idx_production_plan_items_priority ON production_plan_items (priority);
CREATE INDEX IF NOT EXISTS idx_production_plan_items_sales_order_id ON production_plan_items (sales_order_id);

-- ============================================================================
-- 生产订单管理
-- ============================================================================

-- 生产订单表
CREATE TABLE IF NOT EXISTS production_orders (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL,
  created_by BIGINT NULL,
  updated_by BIGINT NULL,
  is_active BOOLEAN DEFAULT TRUE,
  code VARCHAR(50) NOT NULL,
  item_id BIGINT NOT NULL,
  bom_id BIGINT NULL,
  routing_id BIGINT NULL,
  planned_quantity DECIMAL(15,4) NOT NULL,
  produced_quantity DECIMAL(15,4) DEFAULT 0,
  scrapped_quantity DECIMAL(15,4) DEFAULT 0,
  unit_id BIGINT NOT NULL,
  planned_start_date DATE NOT NULL,
  planned_end_date DATE NOT NULL,
  actual_start_date DATE NULL,
  actual_end_date DATE NULL,
  priority INTEGER DEFAULT 5,
  sales_order_id BIGINT NULL,
  production_plan_item_id BIGINT NULL,
  supervisor_id BIGINT NULL,
  warehouse_id BIGINT NOT NULL,
  location_id BIGINT NULL,
  notes TEXT NULL,
  status VARCHAR(20) DEFAULT 'PLANNED',
  CONSTRAINT uq_production_orders_code UNIQUE (code),
  CONSTRAINT fk_production_orders_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE,
  CONSTRAINT fk_production_orders_bom FOREIGN KEY (bom_id) REFERENCES boms(id) ON DELETE SET NULL,
  CONSTRAINT fk_production_orders_routing FOREIGN KEY (routing_id) REFERENCES routings(id) ON DELETE SET NULL,
  CONSTRAINT fk_production_orders_unit FOREIGN KEY (unit_id) REFERENCES units(id) ON DELETE CASCADE,
  CONSTRAINT fk_production_orders_sales_order FOREIGN KEY (sales_order_id) REFERENCES sales_orders(id) ON DELETE SET NULL,
  CONSTRAINT fk_production_orders_production_plan_item FOREIGN KEY (production_plan_item_id) REFERENCES production_plan_items(id) ON DELETE SET NULL,
  CONSTRAINT fk_production_orders_supervisor FOREIGN KEY (supervisor_id) REFERENCES employees(id) ON DELETE SET NULL,
  CONSTRAINT fk_production_orders_warehouse FOREIGN KEY (warehouse_id) REFERENCES warehouses(id) ON DELETE CASCADE,
  CONSTRAINT fk_production_orders_location FOREIGN KEY (location_id) REFERENCES locations(id) ON DELETE SET NULL,
  CONSTRAINT chk_production_orders_status CHECK (status IN ('PLANNED', 'RELEASED', 'IN_PROGRESS', 'COMPLETED', 'CANCELLED', 'ON_HOLD'))
);
CREATE INDEX IF NOT EXISTS idx_production_orders_deleted_at ON production_orders (deleted_at);
CREATE INDEX IF NOT EXISTS idx_production_orders_is_active ON production_orders (is_active);
CREATE INDEX IF NOT EXISTS idx_production_orders_code ON production_orders (code);
CREATE INDEX IF NOT EXISTS idx_production_orders_item_id ON production_orders (item_id);
CREATE INDEX IF NOT EXISTS idx_production_orders_bom_id ON production_orders (bom_id);
CREATE INDEX IF NOT EXISTS idx_production_orders_routing_id ON production_orders (routing_id);
CREATE INDEX IF NOT EXISTS idx_production_orders_planned_start_date ON production_orders (planned_start_date);
CREATE INDEX IF NOT EXISTS idx_production_orders_planned_end_date ON production_orders (planned_end_date);
CREATE INDEX IF NOT EXISTS idx_production_orders_status ON production_orders (status);
CREATE INDEX IF NOT EXISTS idx_production_orders_priority ON production_orders (priority);
CREATE INDEX IF NOT EXISTS idx_production_orders_sales_order_id ON production_orders (sales_order_id);
CREATE INDEX IF NOT EXISTS idx_production_orders_supervisor_id ON production_orders (supervisor_id);
CREATE INDEX IF NOT EXISTS idx_production_orders_warehouse_id ON production_orders (warehouse_id);

-- 生产订单物料需求表
CREATE TABLE IF NOT EXISTS production_order_materials (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL,
  production_order_id BIGINT NOT NULL,
  item_id BIGINT NOT NULL,
  required_quantity DECIMAL(15,4) NOT NULL,
  issued_quantity DECIMAL(15,4) DEFAULT 0,
  consumed_quantity DECIMAL(15,4) DEFAULT 0,
  returned_quantity DECIMAL(15,4) DEFAULT 0,
  unit_id BIGINT NOT NULL,
  unit_cost DECIMAL(15,2) DEFAULT 0,
  warehouse_id BIGINT NULL,
  location_id BIGINT NULL,
  notes TEXT NULL,
  CONSTRAINT fk_production_order_materials_production_order FOREIGN KEY (production_order_id) REFERENCES production_orders(id) ON DELETE CASCADE,
  CONSTRAINT fk_production_order_materials_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE,
  CONSTRAINT fk_production_order_materials_unit FOREIGN KEY (unit_id) REFERENCES units(id) ON DELETE CASCADE,
  CONSTRAINT fk_production_order_materials_warehouse FOREIGN KEY (warehouse_id) REFERENCES warehouses(id) ON DELETE SET NULL,
  CONSTRAINT fk_production_order_materials_location FOREIGN KEY (location_id) REFERENCES locations(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_production_order_materials_deleted_at ON production_order_materials (deleted_at);
CREATE INDEX IF NOT EXISTS idx_production_order_materials_production_order_id ON production_order_materials (production_order_id);
CREATE INDEX IF NOT EXISTS idx_production_order_materials_item_id ON production_order_materials (item_id);
CREATE INDEX IF NOT EXISTS idx_production_order_materials_warehouse_id ON production_order_materials (warehouse_id);
CREATE INDEX IF NOT EXISTS idx_production_order_materials_location_id ON production_order_materials (location_id);

-- 生产订单操作表
CREATE TABLE IF NOT EXISTS production_order_operations (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL,
  production_order_id BIGINT NOT NULL,
  routing_operation_id BIGINT NOT NULL,
  operation_number INTEGER NOT NULL,
  operation_name VARCHAR(255) NOT NULL,
  work_center_id BIGINT NOT NULL,
  planned_setup_time INTEGER DEFAULT 0,
  planned_run_time DECIMAL(10,2) DEFAULT 0,
  actual_setup_time INTEGER DEFAULT 0,
  actual_run_time DECIMAL(10,2) DEFAULT 0,
  planned_start_date TIMESTAMP NULL,
  planned_end_date TIMESTAMP NULL,
  actual_start_date TIMESTAMP NULL,
  actual_end_date TIMESTAMP NULL,
  quantity_completed DECIMAL(15,4) DEFAULT 0,
  quantity_scrapped DECIMAL(15,4) DEFAULT 0,
  operator_id BIGINT NULL,
  status VARCHAR(20) DEFAULT 'PLANNED',
  notes TEXT NULL,
  CONSTRAINT fk_production_order_operations_production_order FOREIGN KEY (production_order_id) REFERENCES production_orders(id) ON DELETE CASCADE,
  CONSTRAINT fk_production_order_operations_routing_operation FOREIGN KEY (routing_operation_id) REFERENCES routing_operations(id) ON DELETE CASCADE,
  CONSTRAINT fk_production_order_operations_work_center FOREIGN KEY (work_center_id) REFERENCES work_centers(id) ON DELETE CASCADE,
  CONSTRAINT fk_production_order_operations_operator FOREIGN KEY (operator_id) REFERENCES employees(id) ON DELETE SET NULL,
  CONSTRAINT chk_production_order_operations_status CHECK (status IN ('PLANNED', 'QUEUED', 'IN_PROGRESS', 'COMPLETED', 'CANCELLED'))
);
CREATE INDEX IF NOT EXISTS idx_production_order_operations_deleted_at ON production_order_operations (deleted_at);
CREATE INDEX IF NOT EXISTS idx_production_order_operations_production_order_id ON production_order_operations (production_order_id);
CREATE INDEX IF NOT EXISTS idx_production_order_operations_routing_operation_id ON production_order_operations (routing_operation_id);
CREATE INDEX IF NOT EXISTS idx_production_order_operations_work_center_id ON production_order_operations (work_center_id);
CREATE INDEX IF NOT EXISTS idx_production_order_operations_operation_number ON production_order_operations (operation_number);
CREATE INDEX IF NOT EXISTS idx_production_order_operations_status ON production_order_operations (status);
CREATE INDEX IF NOT EXISTS idx_production_order_operations_operator_id ON production_order_operations (operator_id);
CREATE INDEX IF NOT EXISTS idx_production_order_operations_planned_start_date ON production_order_operations (planned_start_date);
CREATE INDEX IF NOT EXISTS idx_production_order_operations_actual_start_date ON production_order_operations (actual_start_date);

-- ============================================================================
-- 生产报工管理
-- ============================================================================

-- 生产报工表
CREATE TABLE IF NOT EXISTS work_orders (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL,
  created_by BIGINT NULL,
  updated_by BIGINT NULL,
  is_active BOOLEAN DEFAULT TRUE,
  code VARCHAR(50) NOT NULL,
  production_order_id BIGINT NOT NULL,
  production_order_operation_id BIGINT NOT NULL,
  work_center_id BIGINT NOT NULL,
  operator_id BIGINT NOT NULL,
  start_time TIMESTAMP NOT NULL,
  end_time TIMESTAMP NULL,
  setup_time_minutes INTEGER DEFAULT 0,
  run_time_minutes DECIMAL(10,2) DEFAULT 0,
  quantity_completed DECIMAL(15,4) DEFAULT 0,
  quantity_scrapped DECIMAL(15,4) DEFAULT 0,
  scrap_reason TEXT NULL,
  efficiency_rate DECIMAL(5,2) DEFAULT 100,
  notes TEXT NULL,
  status VARCHAR(20) DEFAULT 'IN_PROGRESS',
  CONSTRAINT uq_work_orders_code UNIQUE (code),
  CONSTRAINT fk_work_orders_production_order FOREIGN KEY (production_order_id) REFERENCES production_orders(id) ON DELETE CASCADE,
  CONSTRAINT fk_work_orders_production_order_operation FOREIGN KEY (production_order_operation_id) REFERENCES production_order_operations(id) ON DELETE CASCADE,
  CONSTRAINT fk_work_orders_work_center FOREIGN KEY (work_center_id) REFERENCES work_centers(id) ON DELETE CASCADE,
  CONSTRAINT fk_work_orders_operator FOREIGN KEY (operator_id) REFERENCES employees(id) ON DELETE CASCADE,
  CONSTRAINT chk_work_orders_status CHECK (status IN ('IN_PROGRESS', 'COMPLETED', 'CANCELLED'))
);
CREATE INDEX IF NOT EXISTS idx_work_orders_deleted_at ON work_orders (deleted_at);
CREATE INDEX IF NOT EXISTS idx_work_orders_is_active ON work_orders (is_active);
CREATE INDEX IF NOT EXISTS idx_work_orders_code ON work_orders (code);
CREATE INDEX IF NOT EXISTS idx_work_orders_production_order_id ON work_orders (production_order_id);
CREATE INDEX IF NOT EXISTS idx_work_orders_production_order_operation_id ON work_orders (production_order_operation_id);
CREATE INDEX IF NOT EXISTS idx_work_orders_work_center_id ON work_orders (work_center_id);
CREATE INDEX IF NOT EXISTS idx_work_orders_operator_id ON work_orders (operator_id);
CREATE INDEX IF NOT EXISTS idx_work_orders_start_time ON work_orders (start_time);
CREATE INDEX IF NOT EXISTS idx_work_orders_status ON work_orders (status);

-- ============================================================================
-- 质量检验管理
-- ============================================================================

-- 质量检验表
CREATE TABLE IF NOT EXISTS quality_inspections (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL,
  created_by BIGINT NULL,
  updated_by BIGINT NULL,
  is_active BOOLEAN DEFAULT TRUE,
  code VARCHAR(50) NOT NULL,
  inspection_type VARCHAR(20) NOT NULL,
  reference_type VARCHAR(50) NOT NULL,
  reference_id BIGINT NOT NULL,
  item_id BIGINT NOT NULL,
  batch_number VARCHAR(100) NULL,
  inspection_date DATE NOT NULL,
  inspector_id BIGINT NOT NULL,
  sample_size DECIMAL(15,4) DEFAULT 0,
  inspected_quantity DECIMAL(15,4) NOT NULL,
  passed_quantity DECIMAL(15,4) DEFAULT 0,
  failed_quantity DECIMAL(15,4) DEFAULT 0,
  unit_id BIGINT NOT NULL,
  inspection_result VARCHAR(20) DEFAULT 'PENDING',
  notes TEXT NULL,
  status VARCHAR(20) DEFAULT 'PENDING',
  CONSTRAINT uq_quality_inspections_code UNIQUE (code),
  CONSTRAINT fk_quality_inspections_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE,
  CONSTRAINT fk_quality_inspections_inspector FOREIGN KEY (inspector_id) REFERENCES employees(id) ON DELETE CASCADE,
  CONSTRAINT fk_quality_inspections_unit FOREIGN KEY (unit_id) REFERENCES units(id) ON DELETE CASCADE,
  CONSTRAINT chk_quality_inspections_inspection_type CHECK (inspection_type IN ('INCOMING', 'IN_PROCESS', 'FINAL', 'OUTGOING')),
  CONSTRAINT chk_quality_inspections_inspection_result CHECK (inspection_result IN ('PENDING', 'PASSED', 'FAILED', 'CONDITIONAL')),
  CONSTRAINT chk_quality_inspections_status CHECK (status IN ('PENDING', 'IN_PROGRESS', 'COMPLETED', 'CANCELLED'))
);
CREATE INDEX IF NOT EXISTS idx_quality_inspections_deleted_at ON quality_inspections (deleted_at);
CREATE INDEX IF NOT EXISTS idx_quality_inspections_is_active ON quality_inspections (is_active);
CREATE INDEX IF NOT EXISTS idx_quality_inspections_code ON quality_inspections (code);
CREATE INDEX IF NOT EXISTS idx_quality_inspections_inspection_type ON quality_inspections (inspection_type);
CREATE INDEX IF NOT EXISTS idx_quality_inspections_reference_type ON quality_inspections (reference_type);
CREATE INDEX IF NOT EXISTS idx_quality_inspections_reference_id ON quality_inspections (reference_id);
CREATE INDEX IF NOT EXISTS idx_quality_inspections_item_id ON quality_inspections (item_id);
CREATE INDEX IF NOT EXISTS idx_quality_inspections_inspection_date ON quality_inspections (inspection_date);
CREATE INDEX IF NOT EXISTS idx_quality_inspections_inspector_id ON quality_inspections (inspector_id);
CREATE INDEX IF NOT EXISTS idx_quality_inspections_inspection_result ON quality_inspections (inspection_result);
CREATE INDEX IF NOT EXISTS idx_quality_inspections_status ON quality_inspections (status);

-- ============================================================================
-- 初始数据插入
-- ============================================================================

-- 插入默认工作中心
INSERT INTO work_centers (code, name, description, work_center_type, capacity_per_hour, efficiency_rate, utilization_rate, cost_per_hour, status, is_active) VALUES
('WC_ASSEMBLY', '装配线', '产品装配工作中心', 'ASSEMBLY', 10.00, 95.00, 85.00, 50.00, 'ACTIVE', TRUE),
('WC_MACHINE', '机加工', '机械加工工作中心', 'MACHINE', 5.00, 90.00, 80.00, 80.00, 'ACTIVE', TRUE),
('WC_INSPECTION', '质检', '质量检验工作中心', 'INSPECTION', 20.00, 98.00, 90.00, 30.00, 'ACTIVE', TRUE)
ON CONFLICT (code) DO NOTHING;

-- ============================================================================
-- 结束
-- ============================================================================