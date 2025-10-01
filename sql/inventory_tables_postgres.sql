-- ============================================================================
-- GalaxyERP 库存管理模块 - PostgreSQL 建表脚本
-- 生成时间: 2025-01-01
-- 说明: 基于 /internal/models/inventory.go 的结构，包含物料、仓库、库存管理
-- ============================================================================

-- ============================================================================
-- 物料分类管理
-- ============================================================================

-- 物料分类表
CREATE TABLE IF NOT EXISTS item_categories (
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
  parent_id BIGINT NULL,
  CONSTRAINT uq_item_categories_code UNIQUE (code),
  CONSTRAINT fk_item_categories_parent FOREIGN KEY (parent_id) REFERENCES item_categories(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_item_categories_deleted_at ON item_categories (deleted_at);
CREATE INDEX IF NOT EXISTS idx_item_categories_is_active ON item_categories (is_active);
CREATE INDEX IF NOT EXISTS idx_item_categories_code ON item_categories (code);
CREATE INDEX IF NOT EXISTS idx_item_categories_parent_id ON item_categories (parent_id);

-- ============================================================================
-- 计量单位管理
-- ============================================================================

-- 计量单位表
CREATE TABLE IF NOT EXISTS units (
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
  symbol VARCHAR(20) NULL,
  type VARCHAR(50) DEFAULT 'QUANTITY',
  CONSTRAINT uq_units_code UNIQUE (code)
);
CREATE INDEX IF NOT EXISTS idx_units_deleted_at ON units (deleted_at);
CREATE INDEX IF NOT EXISTS idx_units_is_active ON units (is_active);
CREATE INDEX IF NOT EXISTS idx_units_code ON units (code);
CREATE INDEX IF NOT EXISTS idx_units_type ON units (type);

-- ============================================================================
-- 物料管理
-- ============================================================================

-- 物料表
CREATE TABLE IF NOT EXISTS items (
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
  category_id BIGINT NULL,
  unit_id BIGINT NOT NULL,
  type VARCHAR(20) DEFAULT 'PRODUCT',
  specification TEXT NULL,
  brand VARCHAR(255) NULL,
  model VARCHAR(255) NULL,
  barcode VARCHAR(255) NULL,
  qr_code VARCHAR(255) NULL,
  weight DECIMAL(15,4) NULL,
  volume DECIMAL(15,4) NULL,
  length DECIMAL(15,4) NULL,
  width DECIMAL(15,4) NULL,
  height DECIMAL(15,4) NULL,
  color VARCHAR(100) NULL,
  material VARCHAR(255) NULL,
  origin_country VARCHAR(100) NULL,
  hs_code VARCHAR(50) NULL,
  standard_cost DECIMAL(15,2) DEFAULT 0,
  list_price DECIMAL(15,2) DEFAULT 0,
  min_stock_level DECIMAL(15,4) DEFAULT 0,
  max_stock_level DECIMAL(15,4) DEFAULT 0,
  reorder_point DECIMAL(15,4) DEFAULT 0,
  reorder_quantity DECIMAL(15,4) DEFAULT 0,
  lead_time_days INTEGER DEFAULT 0,
  shelf_life_days INTEGER NULL,
  is_serialized BOOLEAN DEFAULT FALSE,
  is_batch_tracked BOOLEAN DEFAULT FALSE,
  status VARCHAR(20) DEFAULT 'ACTIVE',
  CONSTRAINT uq_items_code UNIQUE (code),
  CONSTRAINT uq_items_barcode UNIQUE (barcode),
  CONSTRAINT fk_items_category FOREIGN KEY (category_id) REFERENCES item_categories(id) ON DELETE SET NULL,
  CONSTRAINT fk_items_unit FOREIGN KEY (unit_id) REFERENCES units(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_items_deleted_at ON items (deleted_at);
CREATE INDEX IF NOT EXISTS idx_items_is_active ON items (is_active);
CREATE INDEX IF NOT EXISTS idx_items_code ON items (code);
CREATE INDEX IF NOT EXISTS idx_items_category_id ON items (category_id);
CREATE INDEX IF NOT EXISTS idx_items_unit_id ON items (unit_id);
CREATE INDEX IF NOT EXISTS idx_items_type ON items (type);
CREATE INDEX IF NOT EXISTS idx_items_status ON items (status);
CREATE INDEX IF NOT EXISTS idx_items_barcode ON items (barcode);
CREATE INDEX IF NOT EXISTS idx_items_name ON items (name);

-- ============================================================================
-- 仓库管理
-- ============================================================================

-- 仓库表
CREATE TABLE IF NOT EXISTS warehouses (
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
  type VARCHAR(20) DEFAULT 'STANDARD',
  address TEXT NULL,
  contact_person VARCHAR(255) NULL,
  contact_phone VARCHAR(50) NULL,
  contact_email VARCHAR(255) NULL,
  manager_id BIGINT NULL,
  is_default BOOLEAN DEFAULT FALSE,
  CONSTRAINT uq_warehouses_code UNIQUE (code),
  CONSTRAINT fk_warehouses_manager FOREIGN KEY (manager_id) REFERENCES employees(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_warehouses_deleted_at ON warehouses (deleted_at);
CREATE INDEX IF NOT EXISTS idx_warehouses_is_active ON warehouses (is_active);
CREATE INDEX IF NOT EXISTS idx_warehouses_code ON warehouses (code);
CREATE INDEX IF NOT EXISTS idx_warehouses_type ON warehouses (type);
CREATE INDEX IF NOT EXISTS idx_warehouses_manager_id ON warehouses (manager_id);
CREATE INDEX IF NOT EXISTS idx_warehouses_is_default ON warehouses (is_default);

-- 库位表
CREATE TABLE IF NOT EXISTS locations (
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
  warehouse_id BIGINT NOT NULL,
  parent_id BIGINT NULL,
  type VARCHAR(20) DEFAULT 'STORAGE',
  zone VARCHAR(100) NULL,
  aisle VARCHAR(100) NULL,
  rack VARCHAR(100) NULL,
  shelf VARCHAR(100) NULL,
  bin VARCHAR(100) NULL,
  capacity DECIMAL(15,4) NULL,
  max_weight DECIMAL(15,4) NULL,
  CONSTRAINT uq_locations_code UNIQUE (code),
  CONSTRAINT fk_locations_warehouse FOREIGN KEY (warehouse_id) REFERENCES warehouses(id) ON DELETE CASCADE,
  CONSTRAINT fk_locations_parent FOREIGN KEY (parent_id) REFERENCES locations(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_locations_deleted_at ON locations (deleted_at);
CREATE INDEX IF NOT EXISTS idx_locations_is_active ON locations (is_active);
CREATE INDEX IF NOT EXISTS idx_locations_code ON locations (code);
CREATE INDEX IF NOT EXISTS idx_locations_warehouse_id ON locations (warehouse_id);
CREATE INDEX IF NOT EXISTS idx_locations_parent_id ON locations (parent_id);
CREATE INDEX IF NOT EXISTS idx_locations_type ON locations (type);

-- ============================================================================
-- 库存管理
-- ============================================================================

-- 库存表
CREATE TABLE IF NOT EXISTS inventories (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL,
  item_id BIGINT NOT NULL,
  warehouse_id BIGINT NOT NULL,
  location_id BIGINT NULL,
  quantity_on_hand DECIMAL(15,4) DEFAULT 0,
  quantity_reserved DECIMAL(15,4) DEFAULT 0,
  quantity_available DECIMAL(15,4) DEFAULT 0,
  quantity_on_order DECIMAL(15,4) DEFAULT 0,
  average_cost DECIMAL(15,2) DEFAULT 0,
  last_cost DECIMAL(15,2) DEFAULT 0,
  last_count_date DATE NULL,
  last_movement_date TIMESTAMP NULL,
  CONSTRAINT fk_inventories_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE,
  CONSTRAINT fk_inventories_warehouse FOREIGN KEY (warehouse_id) REFERENCES warehouses(id) ON DELETE CASCADE,
  CONSTRAINT fk_inventories_location FOREIGN KEY (location_id) REFERENCES locations(id) ON DELETE SET NULL,
  CONSTRAINT uq_inventories_item_warehouse_location UNIQUE (item_id, warehouse_id, location_id)
);
CREATE INDEX IF NOT EXISTS idx_inventories_deleted_at ON inventories (deleted_at);
CREATE INDEX IF NOT EXISTS idx_inventories_item_id ON inventories (item_id);
CREATE INDEX IF NOT EXISTS idx_inventories_warehouse_id ON inventories (warehouse_id);
CREATE INDEX IF NOT EXISTS idx_inventories_location_id ON inventories (location_id);
CREATE INDEX IF NOT EXISTS idx_inventories_quantity_on_hand ON inventories (quantity_on_hand);
CREATE INDEX IF NOT EXISTS idx_inventories_last_movement_date ON inventories (last_movement_date);

-- ============================================================================
-- 批次管理
-- ============================================================================

-- 批次表
CREATE TABLE IF NOT EXISTS batches (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL,
  created_by BIGINT NULL,
  updated_by BIGINT NULL,
  item_id BIGINT NOT NULL,
  batch_number VARCHAR(100) NOT NULL,
  production_date DATE NULL,
  expiry_date DATE NULL,
  supplier_batch VARCHAR(100) NULL,
  notes TEXT NULL,
  status VARCHAR(20) DEFAULT 'ACTIVE',
  CONSTRAINT fk_batches_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE,
  CONSTRAINT uq_batches_item_batch_number UNIQUE (item_id, batch_number)
);
CREATE INDEX IF NOT EXISTS idx_batches_deleted_at ON batches (deleted_at);
CREATE INDEX IF NOT EXISTS idx_batches_item_id ON batches (item_id);
CREATE INDEX IF NOT EXISTS idx_batches_batch_number ON batches (batch_number);
CREATE INDEX IF NOT EXISTS idx_batches_expiry_date ON batches (expiry_date);
CREATE INDEX IF NOT EXISTS idx_batches_status ON batches (status);

-- 序列号表
CREATE TABLE IF NOT EXISTS serial_numbers (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL,
  created_by BIGINT NULL,
  updated_by BIGINT NULL,
  item_id BIGINT NOT NULL,
  serial_number VARCHAR(100) NOT NULL,
  batch_id BIGINT NULL,
  warehouse_id BIGINT NULL,
  location_id BIGINT NULL,
  status VARCHAR(20) DEFAULT 'AVAILABLE',
  notes TEXT NULL,
  CONSTRAINT fk_serial_numbers_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE,
  CONSTRAINT fk_serial_numbers_batch FOREIGN KEY (batch_id) REFERENCES batches(id) ON DELETE SET NULL,
  CONSTRAINT fk_serial_numbers_warehouse FOREIGN KEY (warehouse_id) REFERENCES warehouses(id) ON DELETE SET NULL,
  CONSTRAINT fk_serial_numbers_location FOREIGN KEY (location_id) REFERENCES locations(id) ON DELETE SET NULL,
  CONSTRAINT uq_serial_numbers_item_serial UNIQUE (item_id, serial_number)
);
CREATE INDEX IF NOT EXISTS idx_serial_numbers_deleted_at ON serial_numbers (deleted_at);
CREATE INDEX IF NOT EXISTS idx_serial_numbers_item_id ON serial_numbers (item_id);
CREATE INDEX IF NOT EXISTS idx_serial_numbers_serial_number ON serial_numbers (serial_number);
CREATE INDEX IF NOT EXISTS idx_serial_numbers_batch_id ON serial_numbers (batch_id);
CREATE INDEX IF NOT EXISTS idx_serial_numbers_warehouse_id ON serial_numbers (warehouse_id);
CREATE INDEX IF NOT EXISTS idx_serial_numbers_status ON serial_numbers (status);

-- ============================================================================
-- 库存移动管理
-- ============================================================================

-- 库存移动表
CREATE TABLE IF NOT EXISTS inventory_movements (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL,
  created_by BIGINT NULL,
  updated_by BIGINT NULL,
  item_id BIGINT NOT NULL,
  warehouse_id BIGINT NOT NULL,
  location_id BIGINT NULL,
  movement_type VARCHAR(50) NOT NULL,
  reference_type VARCHAR(50) NULL,
  reference_id BIGINT NULL,
  reference_number VARCHAR(100) NULL,
  quantity DECIMAL(15,4) NOT NULL,
  unit_cost DECIMAL(15,2) DEFAULT 0,
  total_cost DECIMAL(15,2) DEFAULT 0,
  batch_id BIGINT NULL,
  serial_number VARCHAR(100) NULL,
  movement_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  notes TEXT NULL,
  CONSTRAINT fk_inventory_movements_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE,
  CONSTRAINT fk_inventory_movements_warehouse FOREIGN KEY (warehouse_id) REFERENCES warehouses(id) ON DELETE CASCADE,
  CONSTRAINT fk_inventory_movements_location FOREIGN KEY (location_id) REFERENCES locations(id) ON DELETE SET NULL,
  CONSTRAINT fk_inventory_movements_batch FOREIGN KEY (batch_id) REFERENCES batches(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_inventory_movements_deleted_at ON inventory_movements (deleted_at);
CREATE INDEX IF NOT EXISTS idx_inventory_movements_item_id ON inventory_movements (item_id);
CREATE INDEX IF NOT EXISTS idx_inventory_movements_warehouse_id ON inventory_movements (warehouse_id);
CREATE INDEX IF NOT EXISTS idx_inventory_movements_location_id ON inventory_movements (location_id);
CREATE INDEX IF NOT EXISTS idx_inventory_movements_movement_type ON inventory_movements (movement_type);
CREATE INDEX IF NOT EXISTS idx_inventory_movements_reference_type ON inventory_movements (reference_type);
CREATE INDEX IF NOT EXISTS idx_inventory_movements_reference_id ON inventory_movements (reference_id);
CREATE INDEX IF NOT EXISTS idx_inventory_movements_movement_date ON inventory_movements (movement_date);
CREATE INDEX IF NOT EXISTS idx_inventory_movements_batch_id ON inventory_movements (batch_id);

-- ============================================================================
-- 库存调整管理
-- ============================================================================

-- 库存调整表
CREATE TABLE IF NOT EXISTS inventory_adjustments (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL,
  created_by BIGINT NULL,
  updated_by BIGINT NULL,
  is_active BOOLEAN DEFAULT TRUE,
  code VARCHAR(50) NOT NULL,
  warehouse_id BIGINT NOT NULL,
  adjustment_date DATE NOT NULL,
  reason VARCHAR(255) NOT NULL,
  notes TEXT NULL,
  status VARCHAR(20) DEFAULT 'DRAFT',
  approved_by BIGINT NULL,
  approved_at TIMESTAMP NULL,
  CONSTRAINT uq_inventory_adjustments_code UNIQUE (code),
  CONSTRAINT fk_inventory_adjustments_warehouse FOREIGN KEY (warehouse_id) REFERENCES warehouses(id) ON DELETE CASCADE,
  CONSTRAINT fk_inventory_adjustments_approved_by FOREIGN KEY (approved_by) REFERENCES employees(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_inventory_adjustments_deleted_at ON inventory_adjustments (deleted_at);
CREATE INDEX IF NOT EXISTS idx_inventory_adjustments_is_active ON inventory_adjustments (is_active);
CREATE INDEX IF NOT EXISTS idx_inventory_adjustments_code ON inventory_adjustments (code);
CREATE INDEX IF NOT EXISTS idx_inventory_adjustments_warehouse_id ON inventory_adjustments (warehouse_id);
CREATE INDEX IF NOT EXISTS idx_inventory_adjustments_status ON inventory_adjustments (status);
CREATE INDEX IF NOT EXISTS idx_inventory_adjustments_adjustment_date ON inventory_adjustments (adjustment_date);

-- 库存调整明细表
CREATE TABLE IF NOT EXISTS inventory_adjustment_items (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL,
  adjustment_id BIGINT NOT NULL,
  item_id BIGINT NOT NULL,
  location_id BIGINT NULL,
  batch_id BIGINT NULL,
  serial_number VARCHAR(100) NULL,
  current_quantity DECIMAL(15,4) NOT NULL,
  adjusted_quantity DECIMAL(15,4) NOT NULL,
  difference_quantity DECIMAL(15,4) NOT NULL,
  unit_cost DECIMAL(15,2) DEFAULT 0,
  total_cost DECIMAL(15,2) DEFAULT 0,
  reason VARCHAR(255) NULL,
  notes TEXT NULL,
  CONSTRAINT fk_inventory_adjustment_items_adjustment FOREIGN KEY (adjustment_id) REFERENCES inventory_adjustments(id) ON DELETE CASCADE,
  CONSTRAINT fk_inventory_adjustment_items_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE,
  CONSTRAINT fk_inventory_adjustment_items_location FOREIGN KEY (location_id) REFERENCES locations(id) ON DELETE SET NULL,
  CONSTRAINT fk_inventory_adjustment_items_batch FOREIGN KEY (batch_id) REFERENCES batches(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_inventory_adjustment_items_deleted_at ON inventory_adjustment_items (deleted_at);
CREATE INDEX IF NOT EXISTS idx_inventory_adjustment_items_adjustment_id ON inventory_adjustment_items (adjustment_id);
CREATE INDEX IF NOT EXISTS idx_inventory_adjustment_items_item_id ON inventory_adjustment_items (item_id);
CREATE INDEX IF NOT EXISTS idx_inventory_adjustment_items_location_id ON inventory_adjustment_items (location_id);
CREATE INDEX IF NOT EXISTS idx_inventory_adjustment_items_batch_id ON inventory_adjustment_items (batch_id);

-- ============================================================================
-- 库存盘点管理
-- ============================================================================

-- 库存盘点表
CREATE TABLE IF NOT EXISTS stock_counts (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL,
  created_by BIGINT NULL,
  updated_by BIGINT NULL,
  is_active BOOLEAN DEFAULT TRUE,
  code VARCHAR(50) NOT NULL,
  warehouse_id BIGINT NOT NULL,
  count_date DATE NOT NULL,
  count_type VARCHAR(20) DEFAULT 'FULL',
  status VARCHAR(20) DEFAULT 'PLANNING',
  notes TEXT NULL,
  approved_by BIGINT NULL,
  approved_at TIMESTAMP NULL,
  CONSTRAINT uq_stock_counts_code UNIQUE (code),
  CONSTRAINT fk_stock_counts_warehouse FOREIGN KEY (warehouse_id) REFERENCES warehouses(id) ON DELETE CASCADE,
  CONSTRAINT fk_stock_counts_approved_by FOREIGN KEY (approved_by) REFERENCES employees(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_stock_counts_deleted_at ON stock_counts (deleted_at);
CREATE INDEX IF NOT EXISTS idx_stock_counts_is_active ON stock_counts (is_active);
CREATE INDEX IF NOT EXISTS idx_stock_counts_code ON stock_counts (code);
CREATE INDEX IF NOT EXISTS idx_stock_counts_warehouse_id ON stock_counts (warehouse_id);
CREATE INDEX IF NOT EXISTS idx_stock_counts_status ON stock_counts (status);
CREATE INDEX IF NOT EXISTS idx_stock_counts_count_date ON stock_counts (count_date);

-- 库存盘点明细表
CREATE TABLE IF NOT EXISTS stock_count_items (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL,
  stock_count_id BIGINT NOT NULL,
  item_id BIGINT NOT NULL,
  location_id BIGINT NULL,
  batch_id BIGINT NULL,
  serial_number VARCHAR(100) NULL,
  system_quantity DECIMAL(15,4) NOT NULL,
  counted_quantity DECIMAL(15,4) NULL,
  difference_quantity DECIMAL(15,4) DEFAULT 0,
  unit_cost DECIMAL(15,2) DEFAULT 0,
  variance_cost DECIMAL(15,2) DEFAULT 0,
  counter_id BIGINT NULL,
  count_time TIMESTAMP NULL,
  notes TEXT NULL,
  CONSTRAINT fk_stock_count_items_stock_count FOREIGN KEY (stock_count_id) REFERENCES stock_counts(id) ON DELETE CASCADE,
  CONSTRAINT fk_stock_count_items_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE,
  CONSTRAINT fk_stock_count_items_location FOREIGN KEY (location_id) REFERENCES locations(id) ON DELETE SET NULL,
  CONSTRAINT fk_stock_count_items_batch FOREIGN KEY (batch_id) REFERENCES batches(id) ON DELETE SET NULL,
  CONSTRAINT fk_stock_count_items_counter FOREIGN KEY (counter_id) REFERENCES employees(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_stock_count_items_deleted_at ON stock_count_items (deleted_at);
CREATE INDEX IF NOT EXISTS idx_stock_count_items_stock_count_id ON stock_count_items (stock_count_id);
CREATE INDEX IF NOT EXISTS idx_stock_count_items_item_id ON stock_count_items (item_id);
CREATE INDEX IF NOT EXISTS idx_stock_count_items_location_id ON stock_count_items (location_id);
CREATE INDEX IF NOT EXISTS idx_stock_count_items_batch_id ON stock_count_items (batch_id);
CREATE INDEX IF NOT EXISTS idx_stock_count_items_counter_id ON stock_count_items (counter_id);

-- ============================================================================
-- 初始数据插入
-- ============================================================================

-- 插入默认物料分类
INSERT INTO item_categories (code, name, description, is_active) VALUES
('RAW_MATERIAL', '原材料', '生产用原材料', TRUE),
('FINISHED_GOODS', '成品', '最终产品', TRUE),
('SEMI_FINISHED', '半成品', '半成品物料', TRUE),
('CONSUMABLES', '消耗品', '消耗性物料', TRUE),
('SPARE_PARTS', '备件', '设备备件', TRUE)
ON CONFLICT (code) DO NOTHING;

-- 插入默认计量单位
INSERT INTO units (code, name, description, symbol, type, is_active) VALUES
('PCS', '件', '计件单位', 'pcs', 'QUANTITY', TRUE),
('KG', '千克', '重量单位', 'kg', 'WEIGHT', TRUE),
('M', '米', '长度单位', 'm', 'LENGTH', TRUE),
('L', '升', '体积单位', 'L', 'VOLUME', TRUE),
('SET', '套', '成套单位', 'set', 'QUANTITY', TRUE),
('BOX', '箱', '包装单位', 'box', 'PACKAGE', TRUE)
ON CONFLICT (code) DO NOTHING;

-- 插入默认仓库
INSERT INTO warehouses (code, name, description, type, address, is_default, is_active) VALUES
('WH_MAIN', '主仓库', '主要仓库', 'STANDARD', '主仓库地址', TRUE, TRUE),
('WH_RAW', '原料仓', '原材料仓库', 'RAW_MATERIAL', '原料仓地址', FALSE, TRUE),
('WH_FINISHED', '成品仓', '成品仓库', 'FINISHED_GOODS', '成品仓地址', FALSE, TRUE)
ON CONFLICT (code) DO NOTHING;

-- 插入默认库位
INSERT INTO locations (code, name, description, warehouse_id, type, zone, is_active)
SELECT 'LOC_A01', 'A区01位', 'A区第1个库位', w.id, 'STORAGE', 'A', TRUE
FROM warehouses w WHERE w.code = 'WH_MAIN'
ON CONFLICT (code) DO NOTHING;

-- 插入示例物料
INSERT INTO items (code, name, description, category_id, unit_id, type, standard_cost, list_price, min_stock_level, max_stock_level, reorder_point, reorder_quantity, status, is_active)
SELECT 'ITEM_DEMO', '示例物料', '演示用物料', c.id, u.id, 'PRODUCT', 100.00, 150.00, 10.0, 1000.0, 50.0, 100.0, 'ACTIVE', TRUE
FROM item_categories c, units u 
WHERE c.code = 'FINISHED_GOODS' AND u.code = 'PCS'
ON CONFLICT (code) DO NOTHING;

-- ============================================================================
-- 结束
-- ============================================================================