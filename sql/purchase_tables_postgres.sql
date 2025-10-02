-- ============================================================================
-- GalaxyERP 采购管理模块 - PostgreSQL 建表脚本
-- 生成时间: 2025-01-01
-- 说明: 基于 /internal/models/purchase.go 的结构，包含供应商、采购订单、采购发票管理
-- ============================================================================

-- ============================================================================
-- 供应商管理
-- ============================================================================

-- 供应商表
CREATE TABLE IF NOT EXISTS suppliers (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE NULL,
  created_by INTEGER NULL,
  updated_by INTEGER NULL,
  name VARCHAR(255) NOT NULL,
  code VARCHAR(255) UNIQUE NOT NULL,
  email VARCHAR(255),
  phone VARCHAR(255),
  address VARCHAR(255),
  city VARCHAR(255),
  state VARCHAR(255),
  postal_code VARCHAR(255),
  country VARCHAR(255),
  contact_person VARCHAR(255),
  credit_limit DECIMAL(15,2) DEFAULT 0.00,
  supplier_group VARCHAR(255),
  territory VARCHAR(255),
  quality_rating DECIMAL(15,2) DEFAULT 0.00,
  delivery_rating DECIMAL(15,2) DEFAULT 0.00,
  is_active BOOLEAN DEFAULT true,
  CONSTRAINT fk_suppliers_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_suppliers_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_suppliers_deleted_at ON suppliers (deleted_at);
CREATE INDEX IF NOT EXISTS idx_suppliers_code ON suppliers (code);
CREATE INDEX IF NOT EXISTS idx_suppliers_name ON suppliers (name);
CREATE INDEX IF NOT EXISTS idx_suppliers_is_active ON suppliers (is_active);
CREATE INDEX IF NOT EXISTS idx_suppliers_supplier_group ON suppliers (supplier_group);
CREATE INDEX IF NOT EXISTS idx_suppliers_territory ON suppliers (territory);

-- ============================================================================
-- 采购申请管理
-- ============================================================================

-- 采购申请表
CREATE TABLE IF NOT EXISTS purchase_requests (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE NULL,
  created_by INTEGER NULL,
  updated_by INTEGER NULL,
  request_number VARCHAR(255) UNIQUE NOT NULL,
  request_date DATE NOT NULL,
  required_date DATE,
  department VARCHAR(255),
  priority VARCHAR(255) DEFAULT 'Medium',
  status VARCHAR(255) DEFAULT 'Draft',
  total_amount DECIMAL(15,2) DEFAULT 0.00,
  notes TEXT,
  CONSTRAINT fk_purchase_requests_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_purchase_requests_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_purchase_requests_deleted_at ON purchase_requests (deleted_at);
CREATE INDEX IF NOT EXISTS idx_purchase_requests_request_number ON purchase_requests (request_number);
CREATE INDEX IF NOT EXISTS idx_purchase_requests_status ON purchase_requests (status);
CREATE INDEX IF NOT EXISTS idx_purchase_requests_request_date ON purchase_requests (request_date);
CREATE INDEX IF NOT EXISTS idx_purchase_requests_department ON purchase_requests (department);

-- 采购申请明细表
CREATE TABLE IF NOT EXISTS purchase_request_items (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE NULL,
  created_by INTEGER NULL,
  updated_by INTEGER NULL,
  purchase_request_id INTEGER NOT NULL,
  item_id INTEGER NOT NULL,
  quantity DECIMAL(15,2) NOT NULL,
  unit_price DECIMAL(15,2) DEFAULT 0.00,
  total_amount DECIMAL(15,2) DEFAULT 0.00,
  notes TEXT,
  CONSTRAINT fk_purchase_request_items_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_purchase_request_items_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_purchase_request_items_purchase_request FOREIGN KEY (purchase_request_id) REFERENCES purchase_requests(id) ON DELETE CASCADE,
  CONSTRAINT fk_purchase_request_items_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_purchase_request_items_deleted_at ON purchase_request_items (deleted_at);
CREATE INDEX IF NOT EXISTS idx_purchase_request_items_request_id ON purchase_request_items (purchase_request_id);
CREATE INDEX IF NOT EXISTS idx_purchase_request_items_item_id ON purchase_request_items (item_id);



-- ============================================================================
-- 采购订单管理
-- ============================================================================

-- 采购订单表
CREATE TABLE IF NOT EXISTS purchase_orders (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE NULL,
  created_by INTEGER NULL,
  updated_by INTEGER NULL,
  order_number VARCHAR(255) UNIQUE NOT NULL,
  supplier_id INTEGER NOT NULL,
  order_date DATE NOT NULL,
  expected_delivery_date DATE,
  status VARCHAR(255) DEFAULT 'Draft',
  total_amount DECIMAL(15,2) DEFAULT 0.00,
  tax_amount DECIMAL(15,2) DEFAULT 0.00,
  discount_amount DECIMAL(15,2) DEFAULT 0.00,
  shipping_amount DECIMAL(15,2) DEFAULT 0.00,
  grand_total DECIMAL(15,2) DEFAULT 0.00,
  payment_terms VARCHAR(255),
  delivery_address TEXT,
  notes TEXT,
  CONSTRAINT fk_purchase_orders_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_purchase_orders_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_purchase_orders_supplier FOREIGN KEY (supplier_id) REFERENCES suppliers(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_purchase_orders_deleted_at ON purchase_orders (deleted_at);
CREATE INDEX IF NOT EXISTS idx_purchase_orders_order_number ON purchase_orders (order_number);
CREATE INDEX IF NOT EXISTS idx_purchase_orders_supplier_id ON purchase_orders (supplier_id);
CREATE INDEX IF NOT EXISTS idx_purchase_orders_status ON purchase_orders (status);
CREATE INDEX IF NOT EXISTS idx_purchase_orders_order_date ON purchase_orders (order_date);

-- 采购订单明细表
CREATE TABLE IF NOT EXISTS purchase_order_items (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE NULL,
  created_by INTEGER NULL,
  updated_by INTEGER NULL,
  purchase_order_id INTEGER NOT NULL,
  item_id INTEGER NOT NULL,
  quantity DECIMAL(15,2) NOT NULL,
  unit_price DECIMAL(15,2) NOT NULL,
  total_amount DECIMAL(15,2) NOT NULL,
  received_quantity DECIMAL(15,2) DEFAULT 0.00,
  notes TEXT,
  CONSTRAINT fk_purchase_order_items_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_purchase_order_items_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_purchase_order_items_purchase_order FOREIGN KEY (purchase_order_id) REFERENCES purchase_orders(id) ON DELETE CASCADE,
  CONSTRAINT fk_purchase_order_items_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_purchase_order_items_deleted_at ON purchase_order_items (deleted_at);
CREATE INDEX IF NOT EXISTS idx_purchase_order_items_order_id ON purchase_order_items (purchase_order_id);
CREATE INDEX IF NOT EXISTS idx_purchase_order_items_item_id ON purchase_order_items (item_id);

-- ============================================================================
-- 收货管理
-- ============================================================================

-- 收货单表
CREATE TABLE IF NOT EXISTS purchase_receipts (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE NULL,
  created_by INTEGER NULL,
  updated_by INTEGER NULL,
  receipt_number VARCHAR(255) UNIQUE NOT NULL,
  purchase_order_id INTEGER NOT NULL,
  supplier_id INTEGER NOT NULL,
  receipt_date DATE NOT NULL,
  status VARCHAR(255) DEFAULT 'Draft',
  total_received_amount DECIMAL(15,2) DEFAULT 0.00,
  notes TEXT,
  CONSTRAINT fk_purchase_receipts_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_purchase_receipts_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_purchase_receipts_purchase_order FOREIGN KEY (purchase_order_id) REFERENCES purchase_orders(id) ON DELETE CASCADE,
  CONSTRAINT fk_purchase_receipts_supplier FOREIGN KEY (supplier_id) REFERENCES suppliers(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_purchase_receipts_deleted_at ON purchase_receipts (deleted_at);
CREATE INDEX IF NOT EXISTS idx_purchase_receipts_receipt_number ON purchase_receipts (receipt_number);
CREATE INDEX IF NOT EXISTS idx_purchase_receipts_purchase_order_id ON purchase_receipts (purchase_order_id);
CREATE INDEX IF NOT EXISTS idx_purchase_receipts_supplier_id ON purchase_receipts (supplier_id);
CREATE INDEX IF NOT EXISTS idx_purchase_receipts_status ON purchase_receipts (status);
CREATE INDEX IF NOT EXISTS idx_purchase_receipts_receipt_date ON purchase_receipts (receipt_date);

-- 收货单明细表
CREATE TABLE IF NOT EXISTS purchase_receipt_items (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE NULL,
  created_by INTEGER NULL,
  updated_by INTEGER NULL,
  purchase_receipt_id INTEGER NOT NULL,
  purchase_order_item_id INTEGER,
  item_id INTEGER NOT NULL,
  quantity_received DECIMAL(15,2) NOT NULL,
  unit_price DECIMAL(15,2) DEFAULT 0.00,
  total_amount DECIMAL(15,2) DEFAULT 0.00,
  quality_status VARCHAR(255) DEFAULT 'Pending',
  warehouse_id INTEGER,
  notes TEXT,
  CONSTRAINT fk_purchase_receipt_items_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_purchase_receipt_items_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_purchase_receipt_items_purchase_receipt FOREIGN KEY (purchase_receipt_id) REFERENCES purchase_receipts(id) ON DELETE CASCADE,
  CONSTRAINT fk_purchase_receipt_items_purchase_order_item FOREIGN KEY (purchase_order_item_id) REFERENCES purchase_order_items(id) ON DELETE SET NULL,
  CONSTRAINT fk_purchase_receipt_items_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE,
  CONSTRAINT fk_purchase_receipt_items_warehouse FOREIGN KEY (warehouse_id) REFERENCES warehouses(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_purchase_receipt_items_deleted_at ON purchase_receipt_items (deleted_at);
CREATE INDEX IF NOT EXISTS idx_purchase_receipt_items_receipt_id ON purchase_receipt_items (purchase_receipt_id);
CREATE INDEX IF NOT EXISTS idx_purchase_receipt_items_order_item_id ON purchase_receipt_items (purchase_order_item_id);
CREATE INDEX IF NOT EXISTS idx_purchase_receipt_items_item_id ON purchase_receipt_items (item_id);
CREATE INDEX IF NOT EXISTS idx_purchase_receipt_items_warehouse_id ON purchase_receipt_items (warehouse_id);
CREATE INDEX IF NOT EXISTS idx_purchase_receipt_items_quality_status ON purchase_receipt_items (quality_status);



-- ============================================================================
-- 初始数据插入
-- ============================================================================

-- 插入示例供应商数据
INSERT INTO suppliers (name, code, email, phone, address, city, state, postal_code, country, contact_person, credit_limit, supplier_group, territory, is_active) VALUES
('北京科技有限公司', 'SUP001', 'zhang@bjtech.com', '010-12345678', '朝阳区科技园', '北京', '北京', '100000', '中国', '张经理', 100000.00, 'Technology', 'North', TRUE),
('上海制造集团', 'SUP002', 'li@shmanuf.com', '021-87654321', '浦东新区工业园', '上海', '上海', '200000', '中国', '李总监', 200000.00, 'Manufacturing', 'East', TRUE),
('深圳电子科技', 'SUP003', 'wang@sztech.com', '0755-11223344', '南山区高新园', '深圳', '广东', '518000', '中国', '王部长', 150000.00, 'Electronics', 'South', TRUE);

-- ============================================================================
-- 结束
-- ============================================================================