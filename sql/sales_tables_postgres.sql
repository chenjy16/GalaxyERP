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
    credit_limit DECIMAL(15,2) DEFAULT 0,
    customer_group VARCHAR(255),
    territory VARCHAR(255),
    is_active BOOLEAN DEFAULT true,
    CONSTRAINT fk_customers_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
    CONSTRAINT fk_customers_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_customers_deleted_at ON customers (deleted_at);
CREATE INDEX IF NOT EXISTS idx_customers_code ON customers (code);
CREATE INDEX IF NOT EXISTS idx_customers_name ON customers (name);
CREATE INDEX IF NOT EXISTS idx_customers_is_active ON customers (is_active);
CREATE INDEX IF NOT EXISTS idx_customers_customer_group ON customers (customer_group);
CREATE INDEX IF NOT EXISTS idx_customers_territory ON customers (territory);



-- ============================================================================
-- 报价管理
-- ============================================================================

-- 报价单表
CREATE TABLE IF NOT EXISTS quotations (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE NULL,
  created_by INTEGER NULL,
  updated_by INTEGER NULL,
  quotation_number VARCHAR(255) UNIQUE NOT NULL,
  customer_id INTEGER NOT NULL,
  template_id INTEGER,
  date TIMESTAMP WITH TIME ZONE NOT NULL,
  valid_till TIMESTAMP WITH TIME ZONE NOT NULL,
  status VARCHAR(255) DEFAULT 'Draft',
  subject VARCHAR(255),
  total_amount DECIMAL(15,2) DEFAULT 0,
  discount_amount DECIMAL(15,2) DEFAULT 0,
  tax_amount DECIMAL(15,2) DEFAULT 0,
  grand_total DECIMAL(15,2) DEFAULT 0,
  terms TEXT,
  notes TEXT,
  CONSTRAINT fk_quotations_customer FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE,
  CONSTRAINT fk_quotations_template FOREIGN KEY (template_id) REFERENCES quotation_templates(id) ON DELETE SET NULL,
  CONSTRAINT fk_quotations_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_quotations_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_quotations_deleted_at ON quotations (deleted_at);
CREATE INDEX IF NOT EXISTS idx_quotations_quotation_number ON quotations (quotation_number);
CREATE INDEX IF NOT EXISTS idx_quotations_customer_id ON quotations (customer_id);
CREATE INDEX IF NOT EXISTS idx_quotations_template_id ON quotations (template_id);
CREATE INDEX IF NOT EXISTS idx_quotations_status ON quotations (status);
CREATE INDEX IF NOT EXISTS idx_quotations_date ON quotations (date);
CREATE INDEX IF NOT EXISTS idx_quotations_valid_till ON quotations (valid_till);
CREATE INDEX IF NOT EXISTS idx_quotations_created_by ON quotations (created_by);

-- 报价单明细表
CREATE TABLE IF NOT EXISTS quotation_items (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE NULL,
  created_by INTEGER NULL,
  updated_by INTEGER NULL,
  quotation_id INTEGER NOT NULL,
  item_id INTEGER NOT NULL,
  description VARCHAR(255),
  quantity DECIMAL(15,4) DEFAULT 1,
  rate DECIMAL(15,2) DEFAULT 0,
  amount DECIMAL(15,2) DEFAULT 0,
  discount_rate DECIMAL(5,2) DEFAULT 0,
  discount_amount DECIMAL(15,2) DEFAULT 0,
  tax_rate DECIMAL(5,2) DEFAULT 0,
  tax_amount DECIMAL(15,2) DEFAULT 0,
  total_amount DECIMAL(15,2) DEFAULT 0,
  CONSTRAINT fk_quotation_items_quotation FOREIGN KEY (quotation_id) REFERENCES quotations(id) ON DELETE CASCADE,
  CONSTRAINT fk_quotation_items_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE,
  CONSTRAINT fk_quotation_items_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_quotation_items_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_quotation_items_deleted_at ON quotation_items (deleted_at);
CREATE INDEX IF NOT EXISTS idx_quotation_items_quotation_id ON quotation_items (quotation_id);
CREATE INDEX IF NOT EXISTS idx_quotation_items_item_id ON quotation_items (item_id);

-- ============================================================================
-- 报价单模板管理
-- ============================================================================

-- 报价单模板表
CREATE TABLE IF NOT EXISTS quotation_templates (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE NULL,
  created_by INTEGER NULL,
  updated_by INTEGER NULL,
  name VARCHAR(255) NOT NULL,
  description TEXT,
  is_default BOOLEAN DEFAULT false,
  is_active BOOLEAN DEFAULT true,
  valid_days INTEGER DEFAULT 30,
  terms TEXT,
  notes TEXT,
  discount_rate DECIMAL(5,2) DEFAULT 0,
  tax_rate DECIMAL(5,2) DEFAULT 0,
  CONSTRAINT fk_quotation_templates_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_quotation_templates_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_quotation_templates_deleted_at ON quotation_templates (deleted_at);
CREATE INDEX IF NOT EXISTS idx_quotation_templates_name ON quotation_templates (name);
CREATE INDEX IF NOT EXISTS idx_quotation_templates_is_default ON quotation_templates (is_default);
CREATE INDEX IF NOT EXISTS idx_quotation_templates_is_active ON quotation_templates (is_active);
CREATE INDEX IF NOT EXISTS idx_quotation_templates_created_by ON quotation_templates (created_by);

-- 报价单模板明细表
CREATE TABLE IF NOT EXISTS quotation_template_items (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE NULL,
  created_by INTEGER NULL,
  updated_by INTEGER NULL,
  template_id INTEGER NOT NULL,
  item_id INTEGER NOT NULL,
  description VARCHAR(255),
  quantity DECIMAL(15,4) DEFAULT 1,
  rate DECIMAL(15,2) DEFAULT 0,
  discount_rate DECIMAL(5,2) DEFAULT 0,
  tax_rate DECIMAL(5,2) DEFAULT 0,
  sort_order INTEGER DEFAULT 0,
  CONSTRAINT fk_quotation_template_items_template FOREIGN KEY (template_id) REFERENCES quotation_templates(id) ON DELETE CASCADE,
  CONSTRAINT fk_quotation_template_items_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE,
  CONSTRAINT fk_quotation_template_items_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_quotation_template_items_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_quotation_template_items_deleted_at ON quotation_template_items (deleted_at);
CREATE INDEX IF NOT EXISTS idx_quotation_template_items_template_id ON quotation_template_items (template_id);
CREATE INDEX IF NOT EXISTS idx_quotation_template_items_item_id ON quotation_template_items (item_id);
CREATE INDEX IF NOT EXISTS idx_quotation_template_items_sort_order ON quotation_template_items (sort_order);

-- ============================================================================
-- 报价单版本管理
-- ============================================================================

-- 报价单版本表
CREATE TABLE IF NOT EXISTS quotation_versions (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE NULL,
  created_by INTEGER NULL,
  updated_by INTEGER NULL,
  quotation_id INTEGER NOT NULL,
  version_number INTEGER NOT NULL,
  version_name VARCHAR(255),
  version_data JSONB NOT NULL,
  change_reason TEXT,
  is_active BOOLEAN DEFAULT false,
  CONSTRAINT fk_quotation_versions_quotation FOREIGN KEY (quotation_id) REFERENCES quotations(id) ON DELETE CASCADE,
  CONSTRAINT fk_quotation_versions_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_quotation_versions_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT uk_quotation_versions_quotation_version UNIQUE (quotation_id, version_number)
);
CREATE INDEX IF NOT EXISTS idx_quotation_versions_deleted_at ON quotation_versions (deleted_at);
CREATE INDEX IF NOT EXISTS idx_quotation_versions_quotation_id ON quotation_versions (quotation_id);
CREATE INDEX IF NOT EXISTS idx_quotation_versions_version_number ON quotation_versions (version_number);
CREATE INDEX IF NOT EXISTS idx_quotation_versions_is_active ON quotation_versions (is_active);
CREATE INDEX IF NOT EXISTS idx_quotation_versions_created_by ON quotation_versions (created_by);
CREATE INDEX IF NOT EXISTS idx_quotation_versions_created_at ON quotation_versions (created_at);

-- ============================================================================
-- 销售订单管理
-- ============================================================================

-- 销售订单表
CREATE TABLE IF NOT EXISTS sales_orders (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE NULL,
  created_by INTEGER NULL,
  updated_by INTEGER NULL,
  order_number VARCHAR(255) UNIQUE NOT NULL,
  customer_id INTEGER NOT NULL,
  date TIMESTAMP WITH TIME ZONE NOT NULL,
  delivery_date TIMESTAMP WITH TIME ZONE NOT NULL,
  status VARCHAR(255) DEFAULT 'Draft',
  quotation_id INTEGER,
  total_amount DECIMAL(15,2) DEFAULT 0,
  discount_amount DECIMAL(15,2) DEFAULT 0,
  tax_amount DECIMAL(15,2) DEFAULT 0,
  grand_total DECIMAL(15,2) DEFAULT 0,
  terms TEXT,
  notes TEXT,
  CONSTRAINT fk_sales_orders_customer FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE,
  CONSTRAINT fk_sales_orders_quotation FOREIGN KEY (quotation_id) REFERENCES quotations(id) ON DELETE SET NULL,
  CONSTRAINT fk_sales_orders_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_sales_orders_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_sales_orders_deleted_at ON sales_orders (deleted_at);
CREATE INDEX IF NOT EXISTS idx_sales_orders_order_number ON sales_orders (order_number);
CREATE INDEX IF NOT EXISTS idx_sales_orders_customer_id ON sales_orders (customer_id);
CREATE INDEX IF NOT EXISTS idx_sales_orders_quotation_id ON sales_orders (quotation_id);
CREATE INDEX IF NOT EXISTS idx_sales_orders_status ON sales_orders (status);
CREATE INDEX IF NOT EXISTS idx_sales_orders_date ON sales_orders (date);
CREATE INDEX IF NOT EXISTS idx_sales_orders_delivery_date ON sales_orders (delivery_date);
CREATE INDEX IF NOT EXISTS idx_sales_orders_created_by ON sales_orders (created_by);

-- 销售订单明细表
CREATE TABLE IF NOT EXISTS sales_order_items (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE NULL,
  created_by INTEGER NULL,
  updated_by INTEGER NULL,
  sales_order_id INTEGER NOT NULL,
  item_id INTEGER NOT NULL,
  description VARCHAR(255),
  quantity DECIMAL(15,4) DEFAULT 1,
  delivered_qty DECIMAL(15,4) DEFAULT 0,
  rate DECIMAL(15,2) DEFAULT 0,
  amount DECIMAL(15,2) DEFAULT 0,
  discount_rate DECIMAL(5,2) DEFAULT 0,
  discount_amount DECIMAL(15,2) DEFAULT 0,
  tax_rate DECIMAL(5,2) DEFAULT 0,
  tax_amount DECIMAL(15,2) DEFAULT 0,
  total_amount DECIMAL(15,2) DEFAULT 0,
  warehouse_id INTEGER,
  CONSTRAINT fk_sales_order_items_sales_order FOREIGN KEY (sales_order_id) REFERENCES sales_orders(id) ON DELETE CASCADE,
  CONSTRAINT fk_sales_order_items_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE,
  CONSTRAINT fk_sales_order_items_warehouse FOREIGN KEY (warehouse_id) REFERENCES warehouses(id) ON DELETE SET NULL,
  CONSTRAINT fk_sales_order_items_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_sales_order_items_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_sales_order_items_deleted_at ON sales_order_items (deleted_at);
CREATE INDEX IF NOT EXISTS idx_sales_order_items_sales_order_id ON sales_order_items (sales_order_id);
CREATE INDEX IF NOT EXISTS idx_sales_order_items_item_id ON sales_order_items (item_id);
CREATE INDEX IF NOT EXISTS idx_sales_order_items_warehouse_id ON sales_order_items (warehouse_id);

-- ============================================================================
-- 送货管理
-- ============================================================================

-- 送货单表
CREATE TABLE IF NOT EXISTS delivery_notes (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE NULL,
  created_by INTEGER NULL,
  updated_by INTEGER NULL,
  delivery_number VARCHAR(255) UNIQUE NOT NULL,
  customer_id INTEGER NOT NULL,
  sales_order_id INTEGER,
  date TIMESTAMP WITH TIME ZONE NOT NULL,
  status VARCHAR(255) DEFAULT 'Draft',
  total_quantity DECIMAL(15,4) DEFAULT 0,
  transporter VARCHAR(255),
  driver_name VARCHAR(255),
  vehicle_number VARCHAR(255),
  destination VARCHAR(255),
  notes TEXT,
  CONSTRAINT fk_delivery_notes_customer FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE,
  CONSTRAINT fk_delivery_notes_sales_order FOREIGN KEY (sales_order_id) REFERENCES sales_orders(id) ON DELETE SET NULL,
  CONSTRAINT fk_delivery_notes_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_delivery_notes_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_delivery_notes_deleted_at ON delivery_notes (deleted_at);
CREATE INDEX IF NOT EXISTS idx_delivery_notes_delivery_number ON delivery_notes (delivery_number);
CREATE INDEX IF NOT EXISTS idx_delivery_notes_customer_id ON delivery_notes (customer_id);
CREATE INDEX IF NOT EXISTS idx_delivery_notes_sales_order_id ON delivery_notes (sales_order_id);
CREATE INDEX IF NOT EXISTS idx_delivery_notes_date ON delivery_notes (date);
CREATE INDEX IF NOT EXISTS idx_delivery_notes_status ON delivery_notes (status);
CREATE INDEX IF NOT EXISTS idx_delivery_notes_created_by ON delivery_notes (created_by);

-- 送货单明细表
CREATE TABLE IF NOT EXISTS delivery_note_items (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE NULL,
  created_by INTEGER NULL,
  updated_by INTEGER NULL,
  delivery_note_id INTEGER NOT NULL,
  sales_order_item_id INTEGER,
  item_id INTEGER NOT NULL,
  description VARCHAR(255),
  quantity DECIMAL(15,4) NOT NULL,
  batch_no VARCHAR(255),
  serial_no VARCHAR(255),
  warehouse_id INTEGER,
  CONSTRAINT fk_delivery_note_items_delivery_note FOREIGN KEY (delivery_note_id) REFERENCES delivery_notes(id) ON DELETE CASCADE,
  CONSTRAINT fk_delivery_note_items_sales_order_item FOREIGN KEY (sales_order_item_id) REFERENCES sales_order_items(id) ON DELETE SET NULL,
  CONSTRAINT fk_delivery_note_items_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE,
  CONSTRAINT fk_delivery_note_items_warehouse FOREIGN KEY (warehouse_id) REFERENCES warehouses(id) ON DELETE SET NULL,
  CONSTRAINT fk_delivery_note_items_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_delivery_note_items_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_delivery_note_items_deleted_at ON delivery_note_items (deleted_at);
CREATE INDEX IF NOT EXISTS idx_delivery_note_items_delivery_note_id ON delivery_note_items (delivery_note_id);
CREATE INDEX IF NOT EXISTS idx_delivery_note_items_sales_order_item_id ON delivery_note_items (sales_order_item_id);
CREATE INDEX IF NOT EXISTS idx_delivery_note_items_item_id ON delivery_note_items (item_id);
CREATE INDEX IF NOT EXISTS idx_delivery_note_items_warehouse_id ON delivery_note_items (warehouse_id);

-- ============================================================================
-- 销售发票管理
-- ============================================================================

-- 销售发票表
CREATE TABLE IF NOT EXISTS sales_invoices (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE NULL,
  created_by INTEGER NULL,
  updated_by INTEGER NULL,
  invoice_number VARCHAR(255) UNIQUE NOT NULL,
  customer_id INTEGER NOT NULL,
  sales_order_id INTEGER,
  delivery_note_id INTEGER,
  invoice_date TIMESTAMP WITH TIME ZONE NOT NULL,
  due_date TIMESTAMP WITH TIME ZONE NOT NULL,
  posting_date TIMESTAMP WITH TIME ZONE NOT NULL,
  doc_status VARCHAR(255) DEFAULT 'Draft',
  payment_status VARCHAR(255) DEFAULT 'Unpaid',
  currency VARCHAR(10) DEFAULT 'CNY',
  exchange_rate DECIMAL(15,6) DEFAULT 1,
  sub_total DECIMAL(15,2) DEFAULT 0,
  discount_amount DECIMAL(15,2) DEFAULT 0,
  tax_amount DECIMAL(15,2) DEFAULT 0,
  shipping_amount DECIMAL(15,2) DEFAULT 0,
  grand_total DECIMAL(15,2) DEFAULT 0,
  outstanding_amount DECIMAL(15,2) DEFAULT 0,
  paid_amount DECIMAL(15,2) DEFAULT 0,
  billing_address TEXT,
  shipping_address TEXT,
  payment_terms VARCHAR(255),
  payment_terms_days INTEGER,
  sales_person_id INTEGER,
  territory VARCHAR(255),
  customer_po_number VARCHAR(255),
  project VARCHAR(255),
  cost_center VARCHAR(255),
  terms TEXT,
  notes TEXT,
  CONSTRAINT fk_sales_invoices_customer FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE,
  CONSTRAINT fk_sales_invoices_sales_order FOREIGN KEY (sales_order_id) REFERENCES sales_orders(id) ON DELETE SET NULL,
  CONSTRAINT fk_sales_invoices_delivery_note FOREIGN KEY (delivery_note_id) REFERENCES delivery_notes(id) ON DELETE SET NULL,
  CONSTRAINT fk_sales_invoices_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_sales_invoices_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_sales_invoices_deleted_at ON sales_invoices (deleted_at);
CREATE INDEX IF NOT EXISTS idx_sales_invoices_invoice_number ON sales_invoices (invoice_number);
CREATE INDEX IF NOT EXISTS idx_sales_invoices_customer_id ON sales_invoices (customer_id);
CREATE INDEX IF NOT EXISTS idx_sales_invoices_sales_order_id ON sales_invoices (sales_order_id);
CREATE INDEX IF NOT EXISTS idx_sales_invoices_delivery_note_id ON sales_invoices (delivery_note_id);
CREATE INDEX IF NOT EXISTS idx_sales_invoices_invoice_date ON sales_invoices (invoice_date);
CREATE INDEX IF NOT EXISTS idx_sales_invoices_due_date ON sales_invoices (due_date);
CREATE INDEX IF NOT EXISTS idx_sales_invoices_doc_status ON sales_invoices (doc_status);
CREATE INDEX IF NOT EXISTS idx_sales_invoices_payment_status ON sales_invoices (payment_status);
CREATE INDEX IF NOT EXISTS idx_sales_invoices_created_by ON sales_invoices (created_by);

-- 销售发票明细表
CREATE TABLE IF NOT EXISTS sales_invoice_items (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE NULL,
  created_by INTEGER NULL,
  updated_by INTEGER NULL,
  sales_invoice_id INTEGER NOT NULL,
  sales_order_item_id INTEGER,
  delivery_note_item_id INTEGER,
  item_id INTEGER NOT NULL,
  item_code VARCHAR(255) NOT NULL,
  item_name VARCHAR(255) NOT NULL,
  description VARCHAR(255),
  quantity DECIMAL(15,4) NOT NULL,
  uom VARCHAR(50) NOT NULL,
  conversion_factor DECIMAL(15,4) DEFAULT 1,
  stock_uom VARCHAR(50),
  rate DECIMAL(15,2) NOT NULL,
  price_list_rate DECIMAL(15,2) DEFAULT 0,
  amount DECIMAL(15,2) NOT NULL,
  discount_percentage DECIMAL(5,2) DEFAULT 0,
  discount_amount DECIMAL(15,2) DEFAULT 0,
  tax_category VARCHAR(255),
  tax_rate DECIMAL(5,2) DEFAULT 0,
  tax_amount DECIMAL(15,2) DEFAULT 0,
  net_rate DECIMAL(15,2) NOT NULL,
  net_amount DECIMAL(15,2) NOT NULL,
  warehouse_id INTEGER,
  batch_no VARCHAR(255),
  serial_no VARCHAR(255),
  project VARCHAR(255),
  cost_center VARCHAR(255),
  CONSTRAINT fk_sales_invoice_items_sales_invoice FOREIGN KEY (sales_invoice_id) REFERENCES sales_invoices(id) ON DELETE CASCADE,
  CONSTRAINT fk_sales_invoice_items_sales_order_item FOREIGN KEY (sales_order_item_id) REFERENCES sales_order_items(id) ON DELETE SET NULL,
  CONSTRAINT fk_sales_invoice_items_delivery_note_item FOREIGN KEY (delivery_note_item_id) REFERENCES delivery_note_items(id) ON DELETE SET NULL,
  CONSTRAINT fk_sales_invoice_items_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE,
  CONSTRAINT fk_sales_invoice_items_warehouse FOREIGN KEY (warehouse_id) REFERENCES warehouses(id) ON DELETE SET NULL,
  CONSTRAINT fk_sales_invoice_items_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_sales_invoice_items_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_sales_invoice_items_deleted_at ON sales_invoice_items (deleted_at);
CREATE INDEX IF NOT EXISTS idx_sales_invoice_items_sales_invoice_id ON sales_invoice_items (sales_invoice_id);
CREATE INDEX IF NOT EXISTS idx_sales_invoice_items_sales_order_item_id ON sales_invoice_items (sales_order_item_id);
CREATE INDEX IF NOT EXISTS idx_sales_invoice_items_delivery_note_item_id ON sales_invoice_items (delivery_note_item_id);
CREATE INDEX IF NOT EXISTS idx_sales_invoice_items_item_id ON sales_invoice_items (item_id);
CREATE INDEX IF NOT EXISTS idx_sales_invoice_items_warehouse_id ON sales_invoice_items (warehouse_id);

-- 发票付款记录表
CREATE TABLE IF NOT EXISTS invoice_payments (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE NULL,
  created_by INTEGER NULL,
  updated_by INTEGER NULL,
  sales_invoice_id INTEGER NOT NULL,
  payment_entry_id INTEGER,
  payment_date TIMESTAMP WITH TIME ZONE NOT NULL,
  payment_method VARCHAR(255) NOT NULL,
  amount DECIMAL(15,2) NOT NULL,
  currency VARCHAR(10) DEFAULT 'CNY',
  exchange_rate DECIMAL(15,6) DEFAULT 1,
  reference_number VARCHAR(255),
  bank_account_id INTEGER,
  notes TEXT,
  status VARCHAR(255) DEFAULT 'Pending',
  CONSTRAINT fk_invoice_payments_sales_invoice FOREIGN KEY (sales_invoice_id) REFERENCES sales_invoices(id) ON DELETE CASCADE,
  CONSTRAINT fk_invoice_payments_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_invoice_payments_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_invoice_payments_deleted_at ON invoice_payments (deleted_at);
CREATE INDEX IF NOT EXISTS idx_invoice_payments_sales_invoice_id ON invoice_payments (sales_invoice_id);
CREATE INDEX IF NOT EXISTS idx_invoice_payments_payment_entry_id ON invoice_payments (payment_entry_id);
CREATE INDEX IF NOT EXISTS idx_invoice_payments_payment_date ON invoice_payments (payment_date);
CREATE INDEX IF NOT EXISTS idx_invoice_payments_status ON invoice_payments (status);
CREATE INDEX IF NOT EXISTS idx_invoice_payments_created_by ON invoice_payments (created_by);




-- 定价规则表
CREATE TABLE IF NOT EXISTS pricing_rules (
  id SERIAL PRIMARY KEY,
  rule_name VARCHAR(255) NOT NULL,
  rule_type VARCHAR(255) NOT NULL, -- 'discount', 'markup', 'fixed_price'
  applicable_for VARCHAR(255) NOT NULL, -- 'item', 'item_group', 'customer', 'customer_group'
  item_id INTEGER,
  item_group VARCHAR(255),
  customer_id INTEGER,
  customer_group VARCHAR(255),
  territory VARCHAR(255),
  sales_person_id INTEGER,
  min_quantity DECIMAL(15,4) DEFAULT 0,
  max_quantity DECIMAL(15,4),
  min_amount DECIMAL(15,2) DEFAULT 0,
  max_amount DECIMAL(15,2),
  valid_from TIMESTAMP WITH TIME ZONE,
  valid_to TIMESTAMP WITH TIME ZONE,
  rate DECIMAL(15,2),
  discount_percentage DECIMAL(5,2),
  margin_type VARCHAR(255),
  margin_rate_or_amount DECIMAL(15,2),
  currency VARCHAR(10) DEFAULT 'CNY',
  price_list VARCHAR(255),
  warehouse VARCHAR(255),
  is_active BOOLEAN DEFAULT TRUE,
  priority INTEGER DEFAULT 1,
  company VARCHAR(255),
  created_by INTEGER,
  updated_by INTEGER,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_pricing_rules_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE,
  CONSTRAINT fk_pricing_rules_customer FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE
);
CREATE INDEX idx_pricing_rules_rule_type ON pricing_rules (rule_type);
CREATE INDEX idx_pricing_rules_applicable_for ON pricing_rules (applicable_for);
CREATE INDEX idx_pricing_rules_customer_group ON pricing_rules (customer_group);
CREATE INDEX idx_pricing_rules_item_group ON pricing_rules (item_group);
CREATE INDEX idx_pricing_rules_item_id ON pricing_rules (item_id);
CREATE INDEX idx_pricing_rules_customer_id ON pricing_rules (customer_id);
CREATE INDEX idx_pricing_rules_territory ON pricing_rules (territory);
CREATE INDEX idx_pricing_rules_sales_person_id ON pricing_rules (sales_person_id);
CREATE INDEX idx_pricing_rules_valid_from ON pricing_rules (valid_from);
CREATE INDEX idx_pricing_rules_valid_to ON pricing_rules (valid_to);
CREATE INDEX idx_pricing_rules_is_active ON pricing_rules (is_active);
CREATE INDEX idx_pricing_rules_priority ON pricing_rules (priority);

-- ============================================================================
-- 初始数据插入
-- ============================================================================

-- 插入默认客户数据
INSERT INTO customers (name, code, customer_group, territory, phone, email, address, is_active) VALUES
('默认客户', 'DEFAULT001', 'Default', 'Default', '400-000-0000', 'default@example.com', '默认地址', TRUE)
ON CONFLICT DO NOTHING;

-- ============================================================================
-- 结束
-- ============================================================================