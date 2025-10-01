-- ============================================================================
-- GalaxyERP 销售管理模块 - PostgreSQL 建表脚本
-- 生成时间: 2025-01-01
-- 说明: 基于 /internal/models/sales.go 的结构，包含客户、销售订单、报价管理
-- ============================================================================

-- ============================================================================
-- 客户管理
-- ============================================================================

-- 客户表
CREATE TABLE IF NOT EXISTS customers (
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
  type VARCHAR(20) DEFAULT 'INDIVIDUAL',
  industry VARCHAR(100) NULL,
  tax_number VARCHAR(100) NULL,
  registration_number VARCHAR(100) NULL,
  website VARCHAR(255) NULL,
  phone VARCHAR(50) NULL,
  email VARCHAR(255) NULL,
  fax VARCHAR(50) NULL,
  address TEXT NULL,
  billing_address TEXT NULL,
  shipping_address TEXT NULL,
  credit_limit DECIMAL(15,2) DEFAULT 0,
  payment_terms INTEGER DEFAULT 30,
  discount_rate DECIMAL(5,2) DEFAULT 0,
  sales_rep_id BIGINT NULL,
  status VARCHAR(20) DEFAULT 'ACTIVE',
  CONSTRAINT uq_customers_code UNIQUE (code),
  CONSTRAINT fk_customers_sales_rep FOREIGN KEY (sales_rep_id) REFERENCES employees(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_customers_deleted_at ON customers (deleted_at);
CREATE INDEX IF NOT EXISTS idx_customers_is_active ON customers (is_active);
CREATE INDEX IF NOT EXISTS idx_customers_code ON customers (code);
CREATE INDEX IF NOT EXISTS idx_customers_type ON customers (type);
CREATE INDEX IF NOT EXISTS idx_customers_status ON customers (status);
CREATE INDEX IF NOT EXISTS idx_customers_sales_rep_id ON customers (sales_rep_id);
CREATE INDEX IF NOT EXISTS idx_customers_name ON customers (name);

-- 客户联系人表
CREATE TABLE IF NOT EXISTS customer_contacts (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL,
  created_by BIGINT NULL,
  updated_by BIGINT NULL,
  customer_id BIGINT NOT NULL,
  name VARCHAR(255) NOT NULL,
  title VARCHAR(100) NULL,
  department VARCHAR(100) NULL,
  phone VARCHAR(50) NULL,
  mobile VARCHAR(50) NULL,
  email VARCHAR(255) NULL,
  is_primary BOOLEAN DEFAULT FALSE,
  notes TEXT NULL,
  CONSTRAINT fk_customer_contacts_customer FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_customer_contacts_deleted_at ON customer_contacts (deleted_at);
CREATE INDEX IF NOT EXISTS idx_customer_contacts_customer_id ON customer_contacts (customer_id);
CREATE INDEX IF NOT EXISTS idx_customer_contacts_is_primary ON customer_contacts (is_primary);
CREATE INDEX IF NOT EXISTS idx_customer_contacts_name ON customer_contacts (name);

-- ============================================================================
-- 销售机会管理
-- ============================================================================

-- 销售机会表
CREATE TABLE IF NOT EXISTS opportunities (
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
  customer_id BIGINT NOT NULL,
  sales_rep_id BIGINT NOT NULL,
  stage VARCHAR(50) DEFAULT 'PROSPECTING',
  probability DECIMAL(5,2) DEFAULT 0,
  estimated_value DECIMAL(15,2) DEFAULT 0,
  expected_close_date DATE NULL,
  actual_close_date DATE NULL,
  source VARCHAR(100) NULL,
  competitor VARCHAR(255) NULL,
  next_action TEXT NULL,
  status VARCHAR(20) DEFAULT 'OPEN',
  CONSTRAINT uq_opportunities_code UNIQUE (code),
  CONSTRAINT fk_opportunities_customer FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE,
  CONSTRAINT fk_opportunities_sales_rep FOREIGN KEY (sales_rep_id) REFERENCES employees(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_opportunities_deleted_at ON opportunities (deleted_at);
CREATE INDEX IF NOT EXISTS idx_opportunities_is_active ON opportunities (is_active);
CREATE INDEX IF NOT EXISTS idx_opportunities_code ON opportunities (code);
CREATE INDEX IF NOT EXISTS idx_opportunities_customer_id ON opportunities (customer_id);
CREATE INDEX IF NOT EXISTS idx_opportunities_sales_rep_id ON opportunities (sales_rep_id);
CREATE INDEX IF NOT EXISTS idx_opportunities_stage ON opportunities (stage);
CREATE INDEX IF NOT EXISTS idx_opportunities_status ON opportunities (status);
CREATE INDEX IF NOT EXISTS idx_opportunities_expected_close_date ON opportunities (expected_close_date);

-- ============================================================================
-- 报价管理
-- ============================================================================

-- 报价单表
CREATE TABLE IF NOT EXISTS quotations (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL,
  created_by BIGINT NULL,
  updated_by BIGINT NULL,
  is_active BOOLEAN DEFAULT TRUE,
  code VARCHAR(50) NOT NULL,
  customer_id BIGINT NOT NULL,
  opportunity_id BIGINT NULL,
  sales_rep_id BIGINT NOT NULL,
  quote_date DATE NOT NULL,
  valid_until DATE NOT NULL,
  subtotal DECIMAL(15,2) DEFAULT 0,
  tax_amount DECIMAL(15,2) DEFAULT 0,
  discount_amount DECIMAL(15,2) DEFAULT 0,
  total_amount DECIMAL(15,2) DEFAULT 0,
  currency VARCHAR(10) DEFAULT 'CNY',
  payment_terms INTEGER DEFAULT 30,
  delivery_terms TEXT NULL,
  notes TEXT NULL,
  status VARCHAR(20) DEFAULT 'DRAFT',
  approved_by BIGINT NULL,
  approved_at TIMESTAMP NULL,
  CONSTRAINT uq_quotations_code UNIQUE (code),
  CONSTRAINT fk_quotations_customer FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE,
  CONSTRAINT fk_quotations_opportunity FOREIGN KEY (opportunity_id) REFERENCES opportunities(id) ON DELETE SET NULL,
  CONSTRAINT fk_quotations_sales_rep FOREIGN KEY (sales_rep_id) REFERENCES employees(id) ON DELETE CASCADE,
  CONSTRAINT fk_quotations_approved_by FOREIGN KEY (approved_by) REFERENCES employees(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_quotations_deleted_at ON quotations (deleted_at);
CREATE INDEX IF NOT EXISTS idx_quotations_is_active ON quotations (is_active);
CREATE INDEX IF NOT EXISTS idx_quotations_code ON quotations (code);
CREATE INDEX IF NOT EXISTS idx_quotations_customer_id ON quotations (customer_id);
CREATE INDEX IF NOT EXISTS idx_quotations_opportunity_id ON quotations (opportunity_id);
CREATE INDEX IF NOT EXISTS idx_quotations_sales_rep_id ON quotations (sales_rep_id);
CREATE INDEX IF NOT EXISTS idx_quotations_status ON quotations (status);
CREATE INDEX IF NOT EXISTS idx_quotations_quote_date ON quotations (quote_date);
CREATE INDEX IF NOT EXISTS idx_quotations_valid_until ON quotations (valid_until);

-- 报价单明细表
CREATE TABLE IF NOT EXISTS quotation_items (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL,
  quotation_id BIGINT NOT NULL,
  item_id BIGINT NOT NULL,
  quantity DECIMAL(15,4) NOT NULL,
  unit_price DECIMAL(15,2) NOT NULL,
  discount_rate DECIMAL(5,2) DEFAULT 0,
  discount_amount DECIMAL(15,2) DEFAULT 0,
  line_total DECIMAL(15,2) NOT NULL,
  notes TEXT NULL,
  CONSTRAINT fk_quotation_items_quotation FOREIGN KEY (quotation_id) REFERENCES quotations(id) ON DELETE CASCADE,
  CONSTRAINT fk_quotation_items_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_quotation_items_deleted_at ON quotation_items (deleted_at);
CREATE INDEX IF NOT EXISTS idx_quotation_items_quotation_id ON quotation_items (quotation_id);
CREATE INDEX IF NOT EXISTS idx_quotation_items_item_id ON quotation_items (item_id);

-- ============================================================================
-- 销售订单管理
-- ============================================================================

-- 销售订单表
CREATE TABLE IF NOT EXISTS sales_orders (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL,
  created_by BIGINT NULL,
  updated_by BIGINT NULL,
  is_active BOOLEAN DEFAULT TRUE,
  code VARCHAR(50) NOT NULL,
  customer_id BIGINT NOT NULL,
  quotation_id BIGINT NULL,
  sales_rep_id BIGINT NOT NULL,
  order_date DATE NOT NULL,
  required_date DATE NULL,
  promised_date DATE NULL,
  subtotal DECIMAL(15,2) DEFAULT 0,
  tax_amount DECIMAL(15,2) DEFAULT 0,
  discount_amount DECIMAL(15,2) DEFAULT 0,
  shipping_amount DECIMAL(15,2) DEFAULT 0,
  total_amount DECIMAL(15,2) DEFAULT 0,
  currency VARCHAR(10) DEFAULT 'CNY',
  payment_terms INTEGER DEFAULT 30,
  shipping_address TEXT NULL,
  billing_address TEXT NULL,
  notes TEXT NULL,
  status VARCHAR(20) DEFAULT 'PENDING',
  approved_by BIGINT NULL,
  approved_at TIMESTAMP NULL,
  CONSTRAINT uq_sales_orders_code UNIQUE (code),
  CONSTRAINT fk_sales_orders_customer FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE,
  CONSTRAINT fk_sales_orders_quotation FOREIGN KEY (quotation_id) REFERENCES quotations(id) ON DELETE SET NULL,
  CONSTRAINT fk_sales_orders_sales_rep FOREIGN KEY (sales_rep_id) REFERENCES employees(id) ON DELETE CASCADE,
  CONSTRAINT fk_sales_orders_approved_by FOREIGN KEY (approved_by) REFERENCES employees(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_sales_orders_deleted_at ON sales_orders (deleted_at);
CREATE INDEX IF NOT EXISTS idx_sales_orders_is_active ON sales_orders (is_active);
CREATE INDEX IF NOT EXISTS idx_sales_orders_code ON sales_orders (code);
CREATE INDEX IF NOT EXISTS idx_sales_orders_customer_id ON sales_orders (customer_id);
CREATE INDEX IF NOT EXISTS idx_sales_orders_quotation_id ON sales_orders (quotation_id);
CREATE INDEX IF NOT EXISTS idx_sales_orders_sales_rep_id ON sales_orders (sales_rep_id);
CREATE INDEX IF NOT EXISTS idx_sales_orders_status ON sales_orders (status);
CREATE INDEX IF NOT EXISTS idx_sales_orders_order_date ON sales_orders (order_date);
CREATE INDEX IF NOT EXISTS idx_sales_orders_required_date ON sales_orders (required_date);

-- 销售订单明细表
CREATE TABLE IF NOT EXISTS sales_order_items (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL,
  sales_order_id BIGINT NOT NULL,
  item_id BIGINT NOT NULL,
  quantity DECIMAL(15,4) NOT NULL,
  unit_price DECIMAL(15,2) NOT NULL,
  discount_rate DECIMAL(5,2) DEFAULT 0,
  discount_amount DECIMAL(15,2) DEFAULT 0,
  line_total DECIMAL(15,2) NOT NULL,
  delivered_quantity DECIMAL(15,4) DEFAULT 0,
  remaining_quantity DECIMAL(15,4) NOT NULL,
  notes TEXT NULL,
  CONSTRAINT fk_sales_order_items_sales_order FOREIGN KEY (sales_order_id) REFERENCES sales_orders(id) ON DELETE CASCADE,
  CONSTRAINT fk_sales_order_items_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_sales_order_items_deleted_at ON sales_order_items (deleted_at);
CREATE INDEX IF NOT EXISTS idx_sales_order_items_sales_order_id ON sales_order_items (sales_order_id);
CREATE INDEX IF NOT EXISTS idx_sales_order_items_item_id ON sales_order_items (item_id);

-- ============================================================================
-- 发货管理
-- ============================================================================

-- 发货单表
CREATE TABLE IF NOT EXISTS shipments (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL,
  created_by BIGINT NULL,
  updated_by BIGINT NULL,
  is_active BOOLEAN DEFAULT TRUE,
  code VARCHAR(50) NOT NULL,
  sales_order_id BIGINT NOT NULL,
  customer_id BIGINT NOT NULL,
  shipment_date DATE NOT NULL,
  expected_delivery_date DATE NULL,
  actual_delivery_date DATE NULL,
  shipping_address TEXT NOT NULL,
  carrier VARCHAR(255) NULL,
  tracking_number VARCHAR(255) NULL,
  shipping_cost DECIMAL(15,2) DEFAULT 0,
  notes TEXT NULL,
  status VARCHAR(20) DEFAULT 'PREPARING',
  CONSTRAINT uq_shipments_code UNIQUE (code),
  CONSTRAINT fk_shipments_sales_order FOREIGN KEY (sales_order_id) REFERENCES sales_orders(id) ON DELETE CASCADE,
  CONSTRAINT fk_shipments_customer FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_shipments_deleted_at ON shipments (deleted_at);
CREATE INDEX IF NOT EXISTS idx_shipments_is_active ON shipments (is_active);
CREATE INDEX IF NOT EXISTS idx_shipments_code ON shipments (code);
CREATE INDEX IF NOT EXISTS idx_shipments_sales_order_id ON shipments (sales_order_id);
CREATE INDEX IF NOT EXISTS idx_shipments_customer_id ON shipments (customer_id);
CREATE INDEX IF NOT EXISTS idx_shipments_status ON shipments (status);
CREATE INDEX IF NOT EXISTS idx_shipments_shipment_date ON shipments (shipment_date);
CREATE INDEX IF NOT EXISTS idx_shipments_tracking_number ON shipments (tracking_number);

-- 发货单明细表
CREATE TABLE IF NOT EXISTS shipment_items (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL,
  shipment_id BIGINT NOT NULL,
  sales_order_item_id BIGINT NOT NULL,
  item_id BIGINT NOT NULL,
  quantity DECIMAL(15,4) NOT NULL,
  notes TEXT NULL,
  CONSTRAINT fk_shipment_items_shipment FOREIGN KEY (shipment_id) REFERENCES shipments(id) ON DELETE CASCADE,
  CONSTRAINT fk_shipment_items_sales_order_item FOREIGN KEY (sales_order_item_id) REFERENCES sales_order_items(id) ON DELETE CASCADE,
  CONSTRAINT fk_shipment_items_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_shipment_items_deleted_at ON shipment_items (deleted_at);
CREATE INDEX IF NOT EXISTS idx_shipment_items_shipment_id ON shipment_items (shipment_id);
CREATE INDEX IF NOT EXISTS idx_shipment_items_sales_order_item_id ON shipment_items (sales_order_item_id);
CREATE INDEX IF NOT EXISTS idx_shipment_items_item_id ON shipment_items (item_id);

-- ============================================================================
-- 销售发票管理
-- ============================================================================

-- 销售发票表
CREATE TABLE IF NOT EXISTS sales_invoices (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL,
  created_by BIGINT NULL,
  updated_by BIGINT NULL,
  is_active BOOLEAN DEFAULT TRUE,
  code VARCHAR(50) NOT NULL,
  customer_id BIGINT NOT NULL,
  sales_order_id BIGINT NULL,
  invoice_date DATE NOT NULL,
  due_date DATE NOT NULL,
  subtotal DECIMAL(15,2) DEFAULT 0,
  tax_amount DECIMAL(15,2) DEFAULT 0,
  discount_amount DECIMAL(15,2) DEFAULT 0,
  total_amount DECIMAL(15,2) DEFAULT 0,
  paid_amount DECIMAL(15,2) DEFAULT 0,
  balance_amount DECIMAL(15,2) DEFAULT 0,
  currency VARCHAR(10) DEFAULT 'CNY',
  payment_terms INTEGER DEFAULT 30,
  notes TEXT NULL,
  status VARCHAR(20) DEFAULT 'DRAFT',
  CONSTRAINT uq_sales_invoices_code UNIQUE (code),
  CONSTRAINT fk_sales_invoices_customer FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE,
  CONSTRAINT fk_sales_invoices_sales_order FOREIGN KEY (sales_order_id) REFERENCES sales_orders(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_sales_invoices_deleted_at ON sales_invoices (deleted_at);
CREATE INDEX IF NOT EXISTS idx_sales_invoices_is_active ON sales_invoices (is_active);
CREATE INDEX IF NOT EXISTS idx_sales_invoices_code ON sales_invoices (code);
CREATE INDEX IF NOT EXISTS idx_sales_invoices_customer_id ON sales_invoices (customer_id);
CREATE INDEX IF NOT EXISTS idx_sales_invoices_sales_order_id ON sales_invoices (sales_order_id);
CREATE INDEX IF NOT EXISTS idx_sales_invoices_status ON sales_invoices (status);
CREATE INDEX IF NOT EXISTS idx_sales_invoices_invoice_date ON sales_invoices (invoice_date);
CREATE INDEX IF NOT EXISTS idx_sales_invoices_due_date ON sales_invoices (due_date);

-- 销售发票明细表
CREATE TABLE IF NOT EXISTS sales_invoice_items (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL,
  sales_invoice_id BIGINT NOT NULL,
  item_id BIGINT NOT NULL,
  quantity DECIMAL(15,4) NOT NULL,
  unit_price DECIMAL(15,2) NOT NULL,
  discount_rate DECIMAL(5,2) DEFAULT 0,
  discount_amount DECIMAL(15,2) DEFAULT 0,
  line_total DECIMAL(15,2) NOT NULL,
  notes TEXT NULL,
  CONSTRAINT fk_sales_invoice_items_sales_invoice FOREIGN KEY (sales_invoice_id) REFERENCES sales_invoices(id) ON DELETE CASCADE,
  CONSTRAINT fk_sales_invoice_items_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_sales_invoice_items_deleted_at ON sales_invoice_items (deleted_at);
CREATE INDEX IF NOT EXISTS idx_sales_invoice_items_sales_invoice_id ON sales_invoice_items (sales_invoice_id);
CREATE INDEX IF NOT EXISTS idx_sales_invoice_items_item_id ON sales_invoice_items (item_id);

-- ============================================================================
-- 初始数据插入
-- ============================================================================

-- 插入默认客户类型数据
INSERT INTO customers (code, name, description, type, phone, email, address, credit_limit, payment_terms, status, is_active) VALUES
('CUST_DEFAULT', '默认客户', '系统默认客户', 'INDIVIDUAL', '400-000-0000', 'default@example.com', '默认地址', 100000.00, 30, 'ACTIVE', TRUE)
ON CONFLICT (code) DO NOTHING;

-- ============================================================================
-- 结束
-- ============================================================================