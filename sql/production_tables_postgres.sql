-- ============================================================================
-- GalaxyERP 生产管理模块 - PostgreSQL 建表脚本
-- 生成时间: 2025-01-01
-- 说明: 基于 /internal/models/production.go 的结构，包含BOM、工艺路线、生产订单管理
-- ============================================================================

-- ============================================================================
-- 生产管理模块数据库表结构
-- ============================================================================

-- 产品表
CREATE TABLE IF NOT EXISTS products (
  id SERIAL PRIMARY KEY,
  code VARCHAR(255) NOT NULL UNIQUE,
  name VARCHAR(255) NOT NULL,
  description TEXT,
  category VARCHAR(255) NOT NULL,
  unit VARCHAR(255) NOT NULL,
  price DECIMAL(15,2) DEFAULT 0,
  cost DECIMAL(15,2) DEFAULT 0,
  status VARCHAR(255) DEFAULT 'active',
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_products_code ON products (code);
CREATE INDEX IF NOT EXISTS idx_products_category ON products (category);
CREATE INDEX IF NOT EXISTS idx_products_status ON products (status);

-- 物料清单表
CREATE TABLE IF NOT EXISTS boms (
  id SERIAL PRIMARY KEY,
  product_id INTEGER NOT NULL,
  version VARCHAR(255) NOT NULL,
  effective_date TIMESTAMP WITH TIME ZONE NOT NULL,
  expiry_date TIMESTAMP WITH TIME ZONE,
  quantity DECIMAL(15,4) NOT NULL,
  is_active BOOLEAN DEFAULT TRUE,
  notes TEXT,
  created_by INTEGER NOT NULL,
  updated_by INTEGER,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_boms_product FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_boms_product_id ON boms (product_id);
CREATE INDEX IF NOT EXISTS idx_boms_effective_date ON boms (effective_date);
CREATE INDEX IF NOT EXISTS idx_boms_expiry_date ON boms (expiry_date);
CREATE INDEX IF NOT EXISTS idx_boms_is_active ON boms (is_active);

-- BOM明细表
CREATE TABLE IF NOT EXISTS bom_items (
  id SERIAL PRIMARY KEY,
  bom_id INTEGER NOT NULL,
  item_id INTEGER NOT NULL,
  quantity DECIMAL(15,4) NOT NULL,
  unit_cost DECIMAL(15,2) DEFAULT 0,
  total_cost DECIMAL(15,2) DEFAULT 0,
  scrap_rate DECIMAL(15,4) DEFAULT 0,
  notes TEXT,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_bom_items_bom FOREIGN KEY (bom_id) REFERENCES boms(id) ON DELETE CASCADE,
  CONSTRAINT fk_bom_items_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_bom_items_bom_id ON bom_items (bom_id);
CREATE INDEX IF NOT EXISTS idx_bom_items_item_id ON bom_items (item_id);

-- ============================================================================
-- 工艺路线管理
-- ============================================================================

-- 工作中心表
CREATE TABLE IF NOT EXISTS work_centers (
  id SERIAL PRIMARY KEY,
  code VARCHAR(255) NOT NULL UNIQUE,
  name VARCHAR(255) NOT NULL,
  description TEXT,
  work_center_type VARCHAR(255) NOT NULL,
  capacity DECIMAL(15,2) DEFAULT 8,
  efficiency DECIMAL(15,2) DEFAULT 100,
  cost_per_hour DECIMAL(15,2) DEFAULT 0,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_work_centers_code ON work_centers (code);
CREATE INDEX IF NOT EXISTS idx_work_centers_work_center_type ON work_centers (work_center_type);

-- 操作表
CREATE TABLE IF NOT EXISTS operations (
  id SERIAL PRIMARY KEY,
  code VARCHAR(255) NOT NULL UNIQUE,
  name VARCHAR(255) NOT NULL,
  description TEXT,
  work_center_id INTEGER NOT NULL,
  setup_time DECIMAL(15,2) DEFAULT 0,
  run_time DECIMAL(15,2) DEFAULT 0,
  standard_cost DECIMAL(15,2) DEFAULT 0,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_operations_work_center FOREIGN KEY (work_center_id) REFERENCES work_centers(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_operations_code ON operations (code);
CREATE INDEX IF NOT EXISTS idx_operations_work_center_id ON operations (work_center_id);

-- 生产计划表
CREATE TABLE IF NOT EXISTS production_plans (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  description TEXT,
  plan_type VARCHAR(255) NOT NULL,
  start_date TIMESTAMP WITH TIME ZONE NOT NULL,
  end_date TIMESTAMP WITH TIME ZONE NOT NULL,
  status VARCHAR(255) DEFAULT 'draft',
  created_by INTEGER NOT NULL,
  approved_by INTEGER,
  approved_at TIMESTAMP WITH TIME ZONE,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_production_plans_status ON production_plans (status);
CREATE INDEX IF NOT EXISTS idx_production_plans_start_date ON production_plans (start_date);
CREATE INDEX IF NOT EXISTS idx_production_plans_end_date ON production_plans (end_date);
CREATE INDEX IF NOT EXISTS idx_production_plans_created_by ON production_plans (created_by);

-- 物料需求表
CREATE TABLE IF NOT EXISTS material_requirements (
  id SERIAL PRIMARY KEY,
  plan_id INTEGER NOT NULL,
  item_id INTEGER NOT NULL,
  required_quantity DECIMAL(15,4) NOT NULL,
  available_quantity DECIMAL(15,4) DEFAULT 0,
  net_requirement DECIMAL(15,4) NOT NULL,
  due_date TIMESTAMP WITH TIME ZONE NOT NULL,
  priority INTEGER DEFAULT 0,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_material_requirements_plan FOREIGN KEY (plan_id) REFERENCES production_plans(id) ON DELETE CASCADE,
  CONSTRAINT fk_material_requirements_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_material_requirements_plan_id ON material_requirements (plan_id);
CREATE INDEX IF NOT EXISTS idx_material_requirements_item_id ON material_requirements (item_id);
CREATE INDEX IF NOT EXISTS idx_material_requirements_due_date ON material_requirements (due_date);

-- 工艺路线表
CREATE TABLE IF NOT EXISTS process_routes (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  code VARCHAR(255) NOT NULL UNIQUE,
  description TEXT,
  item_id INTEGER NOT NULL,
  is_active BOOLEAN DEFAULT TRUE,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_process_routes_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_process_routes_code ON process_routes (code);
CREATE INDEX IF NOT EXISTS idx_process_routes_item_id ON process_routes (item_id);
CREATE INDEX IF NOT EXISTS idx_process_routes_is_active ON process_routes (is_active);

-- 工艺操作表
CREATE TABLE IF NOT EXISTS process_operations (
  id SERIAL PRIMARY KEY,
  route_id INTEGER NOT NULL,
  operation_number INTEGER NOT NULL,
  name VARCHAR(255) NOT NULL,
  description TEXT,
  work_center_id INTEGER,
  standard_hours DECIMAL(15,2) DEFAULT 0,
  setup_hours DECIMAL(15,2) DEFAULT 0,
  wait_hours DECIMAL(15,2) DEFAULT 0,
  move_hours DECIMAL(15,2) DEFAULT 0,
  sequence INTEGER NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_process_operations_route FOREIGN KEY (route_id) REFERENCES process_routes(id) ON DELETE CASCADE,
  CONSTRAINT fk_process_operations_work_center FOREIGN KEY (work_center_id) REFERENCES work_centers(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_process_operations_route_id ON process_operations (route_id);
CREATE INDEX IF NOT EXISTS idx_process_operations_work_center_id ON process_operations (work_center_id);
CREATE INDEX IF NOT EXISTS idx_process_operations_sequence ON process_operations (sequence);

-- 生产订单表
CREATE TABLE IF NOT EXISTS production_orders (
  id SERIAL PRIMARY KEY,
  order_number VARCHAR(255) NOT NULL UNIQUE,
  name VARCHAR(255) NOT NULL,
  description TEXT,
  item_id INTEGER NOT NULL,
  quantity DECIMAL(15,4) NOT NULL,
  produced_quantity DECIMAL(15,4) DEFAULT 0,
  unit VARCHAR(255),
  start_date TIMESTAMP WITH TIME ZONE NOT NULL,
  end_date TIMESTAMP WITH TIME ZONE NOT NULL,
  actual_start_date TIMESTAMP WITH TIME ZONE,
  actual_end_date TIMESTAMP WITH TIME ZONE,
  status VARCHAR(255) DEFAULT 'draft',
  priority INTEGER DEFAULT 0,
  route_id INTEGER,
  created_by INTEGER NOT NULL,
  approved_by INTEGER,
  approved_at TIMESTAMP WITH TIME ZONE,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_production_orders_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE,
  CONSTRAINT fk_production_orders_route FOREIGN KEY (route_id) REFERENCES process_routes(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_production_orders_order_number ON production_orders (order_number);
CREATE INDEX IF NOT EXISTS idx_production_orders_item_id ON production_orders (item_id);
CREATE INDEX IF NOT EXISTS idx_production_orders_route_id ON production_orders (route_id);
CREATE INDEX IF NOT EXISTS idx_production_orders_status ON production_orders (status);
CREATE INDEX IF NOT EXISTS idx_production_orders_start_date ON production_orders (start_date);
CREATE INDEX IF NOT EXISTS idx_production_orders_end_date ON production_orders (end_date);
CREATE INDEX IF NOT EXISTS idx_production_orders_created_by ON production_orders (created_by);

-- 生产订单明细表
CREATE TABLE IF NOT EXISTS production_order_items (
  id SERIAL PRIMARY KEY,
  order_id INTEGER NOT NULL,
  item_id INTEGER NOT NULL,
  required_quantity DECIMAL(15,4) NOT NULL,
  issued_quantity DECIMAL(15,4) DEFAULT 0,
  consumed_quantity DECIMAL(15,4) DEFAULT 0,
  unit VARCHAR(255),
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_production_order_items_order FOREIGN KEY (order_id) REFERENCES production_orders(id) ON DELETE CASCADE,
  CONSTRAINT fk_production_order_items_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_production_order_items_order_id ON production_order_items (order_id);
CREATE INDEX IF NOT EXISTS idx_production_order_items_item_id ON production_order_items (item_id);

-- 工作订单表
CREATE TABLE IF NOT EXISTS work_orders (
  id SERIAL PRIMARY KEY,
  work_order_number VARCHAR(255) NOT NULL UNIQUE,
  product_id INTEGER NOT NULL,
  bom_id INTEGER NOT NULL,
  planned_qty DECIMAL(15,4) NOT NULL,
  produced_qty DECIMAL(15,4) DEFAULT 0,
  scrap_qty DECIMAL(15,4) DEFAULT 0,
  start_date TIMESTAMP WITH TIME ZONE NOT NULL,
  end_date TIMESTAMP WITH TIME ZONE NOT NULL,
  actual_start_date TIMESTAMP WITH TIME ZONE,
  actual_end_date TIMESTAMP WITH TIME ZONE,
  priority VARCHAR(255) DEFAULT 'normal',
  status VARCHAR(255) DEFAULT 'planned',
  notes TEXT,
  created_by INTEGER NOT NULL,
  updated_by INTEGER,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_work_orders_product FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
  CONSTRAINT fk_work_orders_bom FOREIGN KEY (bom_id) REFERENCES boms(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_work_orders_work_order_number ON work_orders (work_order_number);
CREATE INDEX IF NOT EXISTS idx_work_orders_product_id ON work_orders (product_id);
CREATE INDEX IF NOT EXISTS idx_work_orders_bom_id ON work_orders (bom_id);
CREATE INDEX IF NOT EXISTS idx_work_orders_start_date ON work_orders (start_date);
CREATE INDEX IF NOT EXISTS idx_work_orders_end_date ON work_orders (end_date);
CREATE INDEX IF NOT EXISTS idx_work_orders_priority ON work_orders (priority);
CREATE INDEX IF NOT EXISTS idx_work_orders_status ON work_orders (status);

-- 工作订单操作表
CREATE TABLE IF NOT EXISTS work_order_operations (
  id SERIAL PRIMARY KEY,
  work_order_id INTEGER NOT NULL,
  operation_number INTEGER NOT NULL,
  name VARCHAR(255) NOT NULL,
  description TEXT,
  work_center_id INTEGER NOT NULL,
  planned_hours DECIMAL(15,2) DEFAULT 0,
  actual_hours DECIMAL(15,2) DEFAULT 0,
  setup_hours DECIMAL(15,2) DEFAULT 0,
  wait_hours DECIMAL(15,2) DEFAULT 0,
  move_hours DECIMAL(15,2) DEFAULT 0,
  status VARCHAR(255) DEFAULT 'pending',
  start_time TIMESTAMP WITH TIME ZONE,
  end_time TIMESTAMP WITH TIME ZONE,
  notes TEXT,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_work_order_operations_work_order FOREIGN KEY (work_order_id) REFERENCES work_orders(id) ON DELETE CASCADE,
  CONSTRAINT fk_work_order_operations_work_center FOREIGN KEY (work_center_id) REFERENCES work_centers(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_work_order_operations_work_order_id ON work_order_operations (work_order_id);
CREATE INDEX IF NOT EXISTS idx_work_order_operations_work_center_id ON work_order_operations (work_center_id);
CREATE INDEX IF NOT EXISTS idx_work_order_operations_operation_number ON work_order_operations (operation_number);
CREATE INDEX IF NOT EXISTS idx_work_order_operations_status ON work_order_operations (status);

-- 工作订单物料表
CREATE TABLE IF NOT EXISTS work_order_materials (
  id SERIAL PRIMARY KEY,
  work_order_id INTEGER NOT NULL,
  item_id INTEGER NOT NULL,
  required_quantity DECIMAL(15,4) NOT NULL,
  issued_quantity DECIMAL(15,4) DEFAULT 0,
  consumed_quantity DECIMAL(15,4) DEFAULT 0,
  unit VARCHAR(255),
  unit_cost DECIMAL(15,2) DEFAULT 0,
  total_cost DECIMAL(15,2) DEFAULT 0,
  notes TEXT,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_work_order_materials_work_order FOREIGN KEY (work_order_id) REFERENCES work_orders(id) ON DELETE CASCADE,
  CONSTRAINT fk_work_order_materials_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_work_order_materials_work_order_id ON work_order_materials (work_order_id);
CREATE INDEX IF NOT EXISTS idx_work_order_materials_item_id ON work_order_materials (item_id);

-- 生产进度表
CREATE TABLE IF NOT EXISTS production_progress (
  id SERIAL PRIMARY KEY,
  work_order_id INTEGER NOT NULL,
  operation_id INTEGER NOT NULL,
  quantity_completed DECIMAL(15,4) NOT NULL,
  quantity_scrapped DECIMAL(15,4) DEFAULT 0,
  progress_date TIMESTAMP WITH TIME ZONE NOT NULL,
  shift VARCHAR(255),
  operator_id INTEGER,
  notes TEXT,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_production_progress_work_order FOREIGN KEY (work_order_id) REFERENCES work_orders(id) ON DELETE CASCADE,
  CONSTRAINT fk_production_progress_operation FOREIGN KEY (operation_id) REFERENCES work_order_operations(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_production_progress_work_order_id ON production_progress (work_order_id);
CREATE INDEX IF NOT EXISTS idx_production_progress_operation_id ON production_progress (operation_id);
CREATE INDEX IF NOT EXISTS idx_production_progress_progress_date ON production_progress (progress_date);
CREATE INDEX IF NOT EXISTS idx_production_progress_operator_id ON production_progress (operator_id);

-- 设备表
CREATE TABLE IF NOT EXISTS equipment (
  id SERIAL PRIMARY KEY,
  code VARCHAR(255) NOT NULL UNIQUE,
  name VARCHAR(255) NOT NULL,
  type VARCHAR(255),
  model VARCHAR(255),
  manufacturer VARCHAR(255),
  purchase_date TIMESTAMP WITH TIME ZONE,
  purchase_cost DECIMAL(15,2),
  work_center_id INTEGER,
  status VARCHAR(255) DEFAULT 'active',
  location VARCHAR(255),
  specifications TEXT,
  notes TEXT,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_equipment_work_center FOREIGN KEY (work_center_id) REFERENCES work_centers(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_equipment_code ON equipment (code);
CREATE INDEX IF NOT EXISTS idx_equipment_name ON equipment (name);
CREATE INDEX IF NOT EXISTS idx_equipment_type ON equipment (type);
CREATE INDEX IF NOT EXISTS idx_equipment_work_center_id ON equipment (work_center_id);
CREATE INDEX IF NOT EXISTS idx_equipment_status ON equipment (status);

-- 设备维护表
CREATE TABLE IF NOT EXISTS equipment_maintenance (
  id SERIAL PRIMARY KEY,
  equipment_id INTEGER NOT NULL,
  maintenance_type VARCHAR(255) NOT NULL,
  scheduled_date TIMESTAMP WITH TIME ZONE NOT NULL,
  completed_date TIMESTAMP WITH TIME ZONE,
  description TEXT,
  cost DECIMAL(15,2) DEFAULT 0,
  performed_by VARCHAR(255),
  status VARCHAR(255) DEFAULT 'scheduled',
  notes TEXT,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_equipment_maintenance_equipment FOREIGN KEY (equipment_id) REFERENCES equipment(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_equipment_maintenance_equipment_id ON equipment_maintenance (equipment_id);
CREATE INDEX IF NOT EXISTS idx_equipment_maintenance_type ON equipment_maintenance (maintenance_type);
CREATE INDEX IF NOT EXISTS idx_equipment_maintenance_scheduled_date ON equipment_maintenance (scheduled_date);
CREATE INDEX IF NOT EXISTS idx_equipment_maintenance_status ON equipment_maintenance (status);

-- 质量检查表
CREATE TABLE IF NOT EXISTS quality_checks (
  id SERIAL PRIMARY KEY,
  work_order_id INTEGER NOT NULL,
  operation_id INTEGER NOT NULL,
  check_type VARCHAR(255) NOT NULL,
  check_date TIMESTAMP WITH TIME ZONE NOT NULL,
  inspector_id INTEGER NOT NULL,
  sample_size INTEGER DEFAULT 0,
  defect_count INTEGER DEFAULT 0,
  pass_count INTEGER DEFAULT 0,
  fail_count INTEGER DEFAULT 0,
  status VARCHAR(255) DEFAULT 'pending',
  notes TEXT,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_quality_checks_work_order FOREIGN KEY (work_order_id) REFERENCES work_orders(id) ON DELETE CASCADE,
  CONSTRAINT fk_quality_checks_operation FOREIGN KEY (operation_id) REFERENCES work_order_operations(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_quality_checks_work_order_id ON quality_checks (work_order_id);
CREATE INDEX IF NOT EXISTS idx_quality_checks_operation_id ON quality_checks (operation_id);
CREATE INDEX IF NOT EXISTS idx_quality_checks_check_type ON quality_checks (check_type);
CREATE INDEX IF NOT EXISTS idx_quality_checks_check_date ON quality_checks (check_date);
CREATE INDEX IF NOT EXISTS idx_quality_checks_inspector_id ON quality_checks (inspector_id);
CREATE INDEX IF NOT EXISTS idx_quality_checks_status ON quality_checks (status);

-- 生产成本表
CREATE TABLE IF NOT EXISTS production_costs (
  id SERIAL PRIMARY KEY,
  work_order_id INTEGER NOT NULL,
  cost_type VARCHAR(255) NOT NULL,
  cost_category VARCHAR(255) NOT NULL,
  amount DECIMAL(15,2) NOT NULL,
  currency VARCHAR(10) DEFAULT 'CNY',
  cost_date TIMESTAMP WITH TIME ZONE NOT NULL,
  description TEXT,
  reference_id INTEGER,
  reference_type VARCHAR(255),
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_production_costs_work_order FOREIGN KEY (work_order_id) REFERENCES work_orders(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_production_costs_work_order_id ON production_costs (work_order_id);
CREATE INDEX IF NOT EXISTS idx_production_costs_cost_type ON production_costs (cost_type);
CREATE INDEX IF NOT EXISTS idx_production_costs_cost_category ON production_costs (cost_category);
CREATE INDEX IF NOT EXISTS idx_production_costs_cost_date ON production_costs (cost_date);
CREATE INDEX IF NOT EXISTS idx_production_costs_reference ON production_costs (reference_id, reference_type);

-- 生产报告表
CREATE TABLE IF NOT EXISTS production_reports (
  id SERIAL PRIMARY KEY,
  report_number VARCHAR(255) NOT NULL UNIQUE,
  report_type VARCHAR(255) NOT NULL,
  title VARCHAR(255) NOT NULL,
  description TEXT,
  report_date TIMESTAMP WITH TIME ZONE NOT NULL,
  period_start TIMESTAMP WITH TIME ZONE NOT NULL,
  period_end TIMESTAMP WITH TIME ZONE NOT NULL,
  status VARCHAR(255) DEFAULT 'draft',
  data JSONB,
  generated_by INTEGER NOT NULL,
  approved_by INTEGER,
  approved_at TIMESTAMP WITH TIME ZONE,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_production_reports_report_number ON production_reports (report_number);
CREATE INDEX IF NOT EXISTS idx_production_reports_report_type ON production_reports (report_type);
CREATE INDEX IF NOT EXISTS idx_production_reports_report_date ON production_reports (report_date);
CREATE INDEX IF NOT EXISTS idx_production_reports_period ON production_reports (period_start, period_end);
CREATE INDEX IF NOT EXISTS idx_production_reports_status ON production_reports (status);
CREATE INDEX IF NOT EXISTS idx_production_reports_generated_by ON production_reports (generated_by);

-- ============================================================================
-- 初始数据插入
-- ============================================================================

-- 插入示例产品数据
INSERT INTO products (id, code, name, description, category, unit, price, cost, status) VALUES
(1, 'PROD001', '产品A', '高质量产品A', '电子产品', 'pcs', 299.99, 150.00, 'active'),
(2, 'PROD002', '产品B', '经济型产品B', '电子产品', 'pcs', 199.99, 100.00, 'active'),
(3, 'MAT001', '原材料A', '基础原材料A', '原材料', 'kg', 50.00, 30.00, 'active'),
(4, 'MAT002', '原材料B', '基础原材料B', '原材料', 'kg', 75.00, 45.00, 'active'),
(5, 'MAT003', '原材料C', '基础原材料C', '原材料', 'kg', 25.00, 15.00, 'active');

-- 插入示例BOM数据
INSERT INTO boms (id, product_id, version, effective_date, quantity, notes, created_by) VALUES
(1, 1, '1.0', CURRENT_TIMESTAMP, 1.00, '产品A物料清单', 1),
(2, 2, '1.0', CURRENT_TIMESTAMP, 1.00, '产品B物料清单', 1);

-- 插入示例BOM明细数据
INSERT INTO bom_items (id, bom_id, item_id, quantity, unit_cost, total_cost, scrap_rate) VALUES
(1, 1, 3, 2.00, 30.00, 60.00, 0.05),
(2, 1, 4, 1.50, 45.00, 67.50, 0.03),
(3, 2, 3, 1.00, 30.00, 30.00, 0.05),
(4, 2, 5, 3.00, 15.00, 45.00, 0.02);

-- 插入示例工作中心数据
INSERT INTO work_centers (id, code, name, description, work_center_type, capacity, efficiency, cost_per_hour) VALUES
(1, 'WC001', '装配线1', '主要装配线', 'assembly', 10.00, 95.00, 50.00),
(2, 'WC002', '包装线1', '产品包装线', 'packaging', 20.00, 98.00, 30.00),
(3, 'WC003', '质检站1', '质量检验工作站', 'inspection', 15.00, 90.00, 40.00);

-- 插入示例操作数据
INSERT INTO operations (id, code, name, description, work_center_id, setup_time, run_time, standard_cost) VALUES
(1, 'OP001', '装配操作', '产品装配操作', 1, 30.00, 60.00, 50.00),
(2, 'OP002', '包装操作', '产品包装操作', 2, 15.00, 30.00, 30.00),
(3, 'OP003', '质检操作', '质量检验操作', 3, 10.00, 20.00, 40.00);

-- 插入示例设备数据
INSERT INTO equipment (id, code, name, type, model, manufacturer, status, location, work_center_id) VALUES
(1, 'EQ001', '装配机器人', '机器人', 'AR-2000', 'RoboCorp', 'active', '车间A', 1),
(2, 'EQ002', '包装机', '包装设备', 'PK-500', 'PackTech', 'active', '车间B', 2),
(3, 'EQ003', '检测仪器', '检测设备', 'QC-100', 'QualityTech', 'active', '质检室', 3);

-- ============================================================================
-- 结束
-- ============================================================================