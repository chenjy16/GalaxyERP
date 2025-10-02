-- ============================================================================
-- GalaxyERP 库存管理模块 - PostgreSQL 建表脚本
-- 生成时间: 2025-01-01
-- 说明: 基于 /internal/models/inventory.go 的结构，包含物料、仓库、库存管理
-- ============================================================================

-- ============================================================================
-- 物料管理
-- ============================================================================

-- 物料表
CREATE TABLE IF NOT EXISTS items (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE NULL,
  created_by INTEGER NULL,
  updated_by INTEGER NULL,
  code VARCHAR(100) NOT NULL,
  name VARCHAR(255) NOT NULL,
  description TEXT NULL,
  category VARCHAR(100) NULL,
  unit VARCHAR(50) NULL,
  cost DECIMAL(15,2) DEFAULT 0,
  price DECIMAL(15,2) DEFAULT 0,
  reorder_level INTEGER DEFAULT 0,
  is_active BOOLEAN DEFAULT TRUE,
  CONSTRAINT uq_items_code UNIQUE (code),
  CONSTRAINT fk_items_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_items_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_items_deleted_at ON items (deleted_at);
CREATE INDEX IF NOT EXISTS idx_items_code ON items (code);
CREATE INDEX IF NOT EXISTS idx_items_category ON items (category);
CREATE INDEX IF NOT EXISTS idx_items_is_active ON items (is_active);

-- ============================================================================
-- 仓库管理
-- ============================================================================

-- 仓库表
CREATE TABLE IF NOT EXISTS warehouses (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE NULL,
  created_by INTEGER NULL,
  updated_by INTEGER NULL,
  code VARCHAR(100) NOT NULL,
  name VARCHAR(255) NOT NULL,
  description TEXT NULL,
  address TEXT NULL,
  is_active BOOLEAN DEFAULT TRUE,
  CONSTRAINT uq_warehouses_code UNIQUE (code),
  CONSTRAINT fk_warehouses_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_warehouses_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_warehouses_deleted_at ON warehouses (deleted_at);
CREATE INDEX IF NOT EXISTS idx_warehouses_code ON warehouses (code);
CREATE INDEX IF NOT EXISTS idx_warehouses_is_active ON warehouses (is_active);



-- ============================================================================
-- 库存管理
-- ============================================================================

-- 库存表
CREATE TABLE IF NOT EXISTS stocks (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE NULL,
  created_by INTEGER NULL,
  updated_by INTEGER NULL,
  item_id INTEGER NOT NULL,
  warehouse_id INTEGER NOT NULL,
  quantity DECIMAL(15,4) DEFAULT 0,
  CONSTRAINT fk_stocks_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE,
  CONSTRAINT fk_stocks_warehouse FOREIGN KEY (warehouse_id) REFERENCES warehouses(id) ON DELETE CASCADE,
  CONSTRAINT fk_stocks_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_stocks_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_stocks_deleted_at ON stocks (deleted_at);
CREATE INDEX IF NOT EXISTS idx_stocks_item_id ON stocks (item_id);
CREATE INDEX IF NOT EXISTS idx_stocks_warehouse_id ON stocks (warehouse_id);

-- ============================================================================
-- 库存移动管理
-- ============================================================================

-- 库存移动表
CREATE TABLE IF NOT EXISTS movements (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE NULL,
  created_by INTEGER NULL,
  updated_by INTEGER NULL,
  item_id INTEGER NULL,
  warehouse_id INTEGER NULL,
  quantity DECIMAL(15,4) NULL,
  movement_type VARCHAR(50) NULL,
  reference VARCHAR(255) NULL,
  notes TEXT NULL,
  unit_cost DECIMAL(15,2) DEFAULT 0,
  total_cost DECIMAL(15,2) DEFAULT 0,
  reference_type VARCHAR(50) NULL,
  reference_id INTEGER NULL,
  batch_no VARCHAR(100) NULL,
  serial_no VARCHAR(100) NULL,
  expiry_date TIMESTAMP WITH TIME ZONE NULL,
  CONSTRAINT fk_movements_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE,
  CONSTRAINT fk_movements_warehouse FOREIGN KEY (warehouse_id) REFERENCES warehouses(id) ON DELETE CASCADE,
  CONSTRAINT fk_movements_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_movements_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_movements_deleted_at ON movements (deleted_at);
CREATE INDEX IF NOT EXISTS idx_movements_item_id ON movements (item_id);
CREATE INDEX IF NOT EXISTS idx_movements_warehouse_id ON movements (warehouse_id);
CREATE INDEX IF NOT EXISTS idx_movements_movement_type ON movements (movement_type);



-- ============================================================================
-- 库存转移管理
-- ============================================================================

-- 库存转移表
CREATE TABLE IF NOT EXISTS stock_transfers (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE NULL,
  created_by INTEGER NULL,
  updated_by INTEGER NULL,
  item_id INTEGER NULL,
  from_warehouse_id INTEGER NULL,
  to_warehouse_id INTEGER NULL,
  quantity DECIMAL(15,4) NULL,
  notes TEXT NULL,
  CONSTRAINT fk_stock_transfers_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE,
  CONSTRAINT fk_stock_transfers_from_warehouse FOREIGN KEY (from_warehouse_id) REFERENCES warehouses(id) ON DELETE CASCADE,
  CONSTRAINT fk_stock_transfers_to_warehouse FOREIGN KEY (to_warehouse_id) REFERENCES warehouses(id) ON DELETE CASCADE,
  CONSTRAINT fk_stock_transfers_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_stock_transfers_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_stock_transfers_deleted_at ON stock_transfers (deleted_at);
CREATE INDEX IF NOT EXISTS idx_stock_transfers_item_id ON stock_transfers (item_id);
CREATE INDEX IF NOT EXISTS idx_stock_transfers_from_warehouse_id ON stock_transfers (from_warehouse_id);
CREATE INDEX IF NOT EXISTS idx_stock_transfers_to_warehouse_id ON stock_transfers (to_warehouse_id);

-- ============================================================================
-- 库存调整管理
-- ============================================================================

-- 库存调整表
CREATE TABLE IF NOT EXISTS stock_adjustments (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE NULL,
  created_by INTEGER NULL,
  updated_by INTEGER NULL,
  approved_by INTEGER NULL,
  approved_at TIMESTAMP WITH TIME ZONE NULL,
  adjustment_number VARCHAR(100) NOT NULL,
  warehouse_id INTEGER NOT NULL,
  adjustment_date TIMESTAMP WITH TIME ZONE NOT NULL,
  reason VARCHAR(255) NOT NULL,
  notes TEXT NULL,
  status VARCHAR(50) DEFAULT 'draft',
  CONSTRAINT fk_stock_adjustments_warehouse FOREIGN KEY (warehouse_id) REFERENCES warehouses(id) ON DELETE CASCADE,
  CONSTRAINT fk_stock_adjustments_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_stock_adjustments_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_stock_adjustments_approved_by FOREIGN KEY (approved_by) REFERENCES users(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_stock_adjustments_deleted_at ON stock_adjustments (deleted_at);
CREATE INDEX IF NOT EXISTS idx_stock_adjustments_warehouse_id ON stock_adjustments (warehouse_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_stock_adjustments_adjustment_number ON stock_adjustments (adjustment_number);
CREATE INDEX IF NOT EXISTS idx_stock_adjustments_adjustment_date ON stock_adjustments (adjustment_date);
CREATE INDEX IF NOT EXISTS idx_stock_adjustments_status ON stock_adjustments (status);

-- 库存调整明细表
CREATE TABLE IF NOT EXISTS stock_adjustment_items (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE NULL,
  created_by INTEGER NULL,
  updated_by INTEGER NULL,
  adjustment_id INTEGER NOT NULL,
  item_id INTEGER NOT NULL,
  location_id INTEGER NULL,
  system_qty DECIMAL(15,4) NOT NULL,
  actual_qty DECIMAL(15,4) NOT NULL,
  difference_qty DECIMAL(15,4) NOT NULL,
  unit_cost DECIMAL(15,2) DEFAULT 0,
  total_cost DECIMAL(15,2) DEFAULT 0,
  notes TEXT NULL,
  CONSTRAINT fk_stock_adjustment_items_adjustment FOREIGN KEY (adjustment_id) REFERENCES stock_adjustments(id) ON DELETE CASCADE,
  CONSTRAINT fk_stock_adjustment_items_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE,
  CONSTRAINT fk_stock_adjustment_items_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_stock_adjustment_items_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_stock_adjustment_items_deleted_at ON stock_adjustment_items (deleted_at);
CREATE INDEX IF NOT EXISTS idx_stock_adjustment_items_adjustment_id ON stock_adjustment_items (adjustment_id);
CREATE INDEX IF NOT EXISTS idx_stock_adjustment_items_item_id ON stock_adjustment_items (item_id);
CREATE INDEX IF NOT EXISTS idx_stock_adjustment_items_location_id ON stock_adjustment_items (location_id);



-- ============================================================================
-- 初始数据插入
-- ============================================================================

-- 插入默认仓库
INSERT INTO warehouses (code, name, description, address, is_active, created_at, updated_at) VALUES
('WH_MAIN', '主仓库', '主要仓库', '主仓库地址', TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('WH_RAW', '原料仓', '原材料仓库', '原料仓地址', TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('WH_FINISHED', '成品仓', '成品仓库', '成品仓地址', TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT (code) DO NOTHING;

-- 插入示例物料
INSERT INTO items (code, name, description, category, unit, cost, price, reorder_level, is_active, created_at, updated_at) VALUES
('ITEM_001', '示例物料A', '演示用物料A', 'finished_goods', 'pcs', 100.00, 150.00, 10, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('ITEM_002', '示例物料B', '演示用物料B', 'raw_material', 'kg', 50.00, 80.00, 20, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT (code) DO NOTHING;

-- ============================================================================
-- 结束
-- ============================================================================