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
  type VARCHAR(20) DEFAULT 'VENDOR',
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
  lead_time_days INTEGER DEFAULT 7,
  quality_rating DECIMAL(3,2) DEFAULT 0,
  delivery_rating DECIMAL(3,2) DEFAULT 0,
  service_rating DECIMAL(3,2) DEFAULT 0,
  overall_rating DECIMAL(3,2) DEFAULT 0,
  buyer_id BIGINT NULL,
  status VARCHAR(20) DEFAULT 'ACTIVE',
  CONSTRAINT uq_suppliers_code UNIQUE (code),
  CONSTRAINT fk_suppliers_buyer FOREIGN KEY (buyer_id) REFERENCES employees(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_suppliers_deleted_at ON suppliers (deleted_at);
CREATE INDEX IF NOT EXISTS idx_suppliers_is_active ON suppliers (is_active);
CREATE INDEX IF NOT EXISTS idx_suppliers_code ON suppliers (code);
CREATE INDEX IF NOT EXISTS idx_suppliers_type ON suppliers (type);
CREATE INDEX IF NOT EXISTS idx_suppliers_status ON suppliers (status);
CREATE INDEX IF NOT EXISTS idx_suppliers_buyer_id ON suppliers (buyer_id);
CREATE INDEX IF NOT EXISTS idx_suppliers_name ON suppliers (name);

-- 供应商联系人表
CREATE TABLE IF NOT EXISTS supplier_contacts (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL,
  created_by BIGINT NULL,
  updated_by BIGINT NULL,
  supplier_id BIGINT NOT NULL,
  name VARCHAR(255) NOT NULL,
  title VARCHAR(100) NULL,
  department VARCHAR(100) NULL,
  phone VARCHAR(50) NULL,
  mobile VARCHAR(50) NULL,
  email VARCHAR(255) NULL,
  is_primary BOOLEAN DEFAULT FALSE,
  notes TEXT NULL,
  CONSTRAINT fk_supplier_contacts_supplier FOREIGN KEY (supplier_id) REFERENCES suppliers(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_supplier_contacts_deleted_at ON supplier_contacts (deleted_at);
CREATE INDEX IF NOT EXISTS idx_supplier_contacts_supplier_id ON supplier_contacts (supplier_id);
CREATE INDEX IF NOT EXISTS idx_supplier_contacts_is_primary ON supplier_contacts (is_primary);
CREATE INDEX IF NOT EXISTS idx_supplier_contacts_name ON supplier_contacts (name);

-- ============================================================================
-- 采购申请管理
-- ============================================================================

-- 采购申请表
CREATE TABLE IF NOT EXISTS purchase_requests (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL,
  created_by BIGINT NULL,
  updated_by BIGINT NULL,
  is_active BOOLEAN DEFAULT TRUE,
  code VARCHAR(50) NOT NULL,
  requester_id BIGINT NOT NULL,
  department_id BIGINT NOT NULL,
  request_date DATE NOT NULL,
  required_date DATE NOT NULL,
  priority VARCHAR(20) DEFAULT 'NORMAL',
  purpose TEXT NULL,
  justification TEXT NULL,
  total_amount DECIMAL(15,2) DEFAULT 0,
  currency VARCHAR(10) DEFAULT 'CNY',
  status VARCHAR(20) DEFAULT 'DRAFT',
  approved_by BIGINT NULL,
  approved_at TIMESTAMP NULL,
  rejection_reason TEXT NULL,
  CONSTRAINT uq_purchase_requests_code UNIQUE (code),
  CONSTRAINT fk_purchase_requests_requester FOREIGN KEY (requester_id) REFERENCES employees(id) ON DELETE CASCADE,
  CONSTRAINT fk_purchase_requests_department FOREIGN KEY (department_id) REFERENCES departments(id) ON DELETE CASCADE,
  CONSTRAINT fk_purchase_requests_approved_by FOREIGN KEY (approved_by) REFERENCES employees(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_purchase_requests_deleted_at ON purchase_requests (deleted_at);
CREATE INDEX IF NOT EXISTS idx_purchase_requests_is_active ON purchase_requests (is_active);
CREATE INDEX IF NOT EXISTS idx_purchase_requests_code ON purchase_requests (code);
CREATE INDEX IF NOT EXISTS idx_purchase_requests_requester_id ON purchase_requests (requester_id);
CREATE INDEX IF NOT EXISTS idx_purchase_requests_department_id ON purchase_requests (department_id);
CREATE INDEX IF NOT EXISTS idx_purchase_requests_status ON purchase_requests (status);
CREATE INDEX IF NOT EXISTS idx_purchase_requests_request_date ON purchase_requests (request_date);
CREATE INDEX IF NOT EXISTS idx_purchase_requests_required_date ON purchase_requests (required_date);
CREATE INDEX IF NOT EXISTS idx_purchase_requests_priority ON purchase_requests (priority);

-- 采购申请明细表
CREATE TABLE IF NOT EXISTS purchase_request_items (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL,
  purchase_request_id BIGINT NOT NULL,
  item_id BIGINT NOT NULL,
  quantity DECIMAL(15,4) NOT NULL,
  estimated_unit_price DECIMAL(15,2) DEFAULT 0,
  estimated_total DECIMAL(15,2) DEFAULT 0,
  specification TEXT NULL,
  notes TEXT NULL,
  CONSTRAINT fk_purchase_request_items_request FOREIGN KEY (purchase_request_id) REFERENCES purchase_requests(id) ON DELETE CASCADE,
  CONSTRAINT fk_purchase_request_items_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_purchase_request_items_deleted_at ON purchase_request_items (deleted_at);
CREATE INDEX IF NOT EXISTS idx_purchase_request_items_request_id ON purchase_request_items (purchase_request_id);
CREATE INDEX IF NOT EXISTS idx_purchase_request_items_item_id ON purchase_request_items (item_id);

-- ============================================================================
-- 询价管理
-- ============================================================================

-- 询价单表
CREATE TABLE IF NOT EXISTS rfqs (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL,
  created_by BIGINT NULL,
  updated_by BIGINT NULL,
  is_active BOOLEAN DEFAULT TRUE,
  code VARCHAR(50) NOT NULL,
  purchase_request_id BIGINT NULL,
  buyer_id BIGINT NOT NULL,
  rfq_date DATE NOT NULL,
  response_deadline DATE NOT NULL,
  delivery_terms TEXT NULL,
  payment_terms TEXT NULL,
  notes TEXT NULL,
  status VARCHAR(20) DEFAULT 'DRAFT',
  CONSTRAINT uq_rfqs_code UNIQUE (code),
  CONSTRAINT fk_rfqs_purchase_request FOREIGN KEY (purchase_request_id) REFERENCES purchase_requests(id) ON DELETE SET NULL,
  CONSTRAINT fk_rfqs_buyer FOREIGN KEY (buyer_id) REFERENCES employees(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_rfqs_deleted_at ON rfqs (deleted_at);
CREATE INDEX IF NOT EXISTS idx_rfqs_is_active ON rfqs (is_active);
CREATE INDEX IF NOT EXISTS idx_rfqs_code ON rfqs (code);
CREATE INDEX IF NOT EXISTS idx_rfqs_purchase_request_id ON rfqs (purchase_request_id);
CREATE INDEX IF NOT EXISTS idx_rfqs_buyer_id ON rfqs (buyer_id);
CREATE INDEX IF NOT EXISTS idx_rfqs_status ON rfqs (status);
CREATE INDEX IF NOT EXISTS idx_rfqs_rfq_date ON rfqs (rfq_date);
CREATE INDEX IF NOT EXISTS idx_rfqs_response_deadline ON rfqs (response_deadline);

-- 询价单明细表
CREATE TABLE IF NOT EXISTS rfq_items (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL,
  rfq_id BIGINT NOT NULL,
  item_id BIGINT NOT NULL,
  quantity DECIMAL(15,4) NOT NULL,
  specification TEXT NULL,
  notes TEXT NULL,
  CONSTRAINT fk_rfq_items_rfq FOREIGN KEY (rfq_id) REFERENCES rfqs(id) ON DELETE CASCADE,
  CONSTRAINT fk_rfq_items_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_rfq_items_deleted_at ON rfq_items (deleted_at);
CREATE INDEX IF NOT EXISTS idx_rfq_items_rfq_id ON rfq_items (rfq_id);
CREATE INDEX IF NOT EXISTS idx_rfq_items_item_id ON rfq_items (item_id);

-- 供应商询价表
CREATE TABLE IF NOT EXISTS rfq_suppliers (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL,
  rfq_id BIGINT NOT NULL,
  supplier_id BIGINT NOT NULL,
  sent_date DATE NULL,
  response_date DATE NULL,
  status VARCHAR(20) DEFAULT 'SENT',
  notes TEXT NULL,
  CONSTRAINT fk_rfq_suppliers_rfq FOREIGN KEY (rfq_id) REFERENCES rfqs(id) ON DELETE CASCADE,
  CONSTRAINT fk_rfq_suppliers_supplier FOREIGN KEY (supplier_id) REFERENCES suppliers(id) ON DELETE CASCADE,
  CONSTRAINT uq_rfq_suppliers_rfq_supplier UNIQUE (rfq_id, supplier_id)
);
CREATE INDEX IF NOT EXISTS idx_rfq_suppliers_deleted_at ON rfq_suppliers (deleted_at);
CREATE INDEX IF NOT EXISTS idx_rfq_suppliers_rfq_id ON rfq_suppliers (rfq_id);
CREATE INDEX IF NOT EXISTS idx_rfq_suppliers_supplier_id ON rfq_suppliers (supplier_id);
CREATE INDEX IF NOT EXISTS idx_rfq_suppliers_status ON rfq_suppliers (status);

-- ============================================================================
-- 报价管理
-- ============================================================================

-- 供应商报价表
CREATE TABLE IF NOT EXISTS supplier_quotes (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL,
  created_by BIGINT NULL,
  updated_by BIGINT NULL,
  is_active BOOLEAN DEFAULT TRUE,
  code VARCHAR(50) NOT NULL,
  rfq_id BIGINT NOT NULL,
  supplier_id BIGINT NOT NULL,
  quote_date DATE NOT NULL,
  valid_until DATE NOT NULL,
  subtotal DECIMAL(15,2) DEFAULT 0,
  tax_amount DECIMAL(15,2) DEFAULT 0,
  total_amount DECIMAL(15,2) DEFAULT 0,
  currency VARCHAR(10) DEFAULT 'CNY',
  payment_terms TEXT NULL,
  delivery_terms TEXT NULL,
  lead_time_days INTEGER DEFAULT 0,
  notes TEXT NULL,
  status VARCHAR(20) DEFAULT 'RECEIVED',
  CONSTRAINT uq_supplier_quotes_code UNIQUE (code),
  CONSTRAINT fk_supplier_quotes_rfq FOREIGN KEY (rfq_id) REFERENCES rfqs(id) ON DELETE CASCADE,
  CONSTRAINT fk_supplier_quotes_supplier FOREIGN KEY (supplier_id) REFERENCES suppliers(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_supplier_quotes_deleted_at ON supplier_quotes (deleted_at);
CREATE INDEX IF NOT EXISTS idx_supplier_quotes_is_active ON supplier_quotes (is_active);
CREATE INDEX IF NOT EXISTS idx_supplier_quotes_code ON supplier_quotes (code);
CREATE INDEX IF NOT EXISTS idx_supplier_quotes_rfq_id ON supplier_quotes (rfq_id);
CREATE INDEX IF NOT EXISTS idx_supplier_quotes_supplier_id ON supplier_quotes (supplier_id);
CREATE INDEX IF NOT EXISTS idx_supplier_quotes_status ON supplier_quotes (status);
CREATE INDEX IF NOT EXISTS idx_supplier_quotes_quote_date ON supplier_quotes (quote_date);
CREATE INDEX IF NOT EXISTS idx_supplier_quotes_valid_until ON supplier_quotes (valid_until);

-- 供应商报价明细表
CREATE TABLE IF NOT EXISTS supplier_quote_items (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL,
  supplier_quote_id BIGINT NOT NULL,
  rfq_item_id BIGINT NOT NULL,
  item_id BIGINT NOT NULL,
  quantity DECIMAL(15,4) NOT NULL,
  unit_price DECIMAL(15,2) NOT NULL,
  line_total DECIMAL(15,2) NOT NULL,
  lead_time_days INTEGER DEFAULT 0,
  notes TEXT NULL,
  CONSTRAINT fk_supplier_quote_items_quote FOREIGN KEY (supplier_quote_id) REFERENCES supplier_quotes(id) ON DELETE CASCADE,
  CONSTRAINT fk_supplier_quote_items_rfq_item FOREIGN KEY (rfq_item_id) REFERENCES rfq_items(id) ON DELETE CASCADE,
  CONSTRAINT fk_supplier_quote_items_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_supplier_quote_items_deleted_at ON supplier_quote_items (deleted_at);
CREATE INDEX IF NOT EXISTS idx_supplier_quote_items_quote_id ON supplier_quote_items (supplier_quote_id);
CREATE INDEX IF NOT EXISTS idx_supplier_quote_items_rfq_item_id ON supplier_quote_items (rfq_item_id);
CREATE INDEX IF NOT EXISTS idx_supplier_quote_items_item_id ON supplier_quote_items (item_id);

-- ============================================================================
-- 采购订单管理
-- ============================================================================

-- 采购订单表
CREATE TABLE IF NOT EXISTS purchase_orders (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL,
  created_by BIGINT NULL,
  updated_by BIGINT NULL,
  is_active BOOLEAN DEFAULT TRUE,
  code VARCHAR(50) NOT NULL,
  supplier_id BIGINT NOT NULL,
  supplier_quote_id BIGINT NULL,
  buyer_id BIGINT NOT NULL,
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
  delivery_address TEXT NULL,
  billing_address TEXT NULL,
  notes TEXT NULL,
  status VARCHAR(20) DEFAULT 'PENDING',
  approved_by BIGINT NULL,
  approved_at TIMESTAMP NULL,
  CONSTRAINT uq_purchase_orders_code UNIQUE (code),
  CONSTRAINT fk_purchase_orders_supplier FOREIGN KEY (supplier_id) REFERENCES suppliers(id) ON DELETE CASCADE,
  CONSTRAINT fk_purchase_orders_supplier_quote FOREIGN KEY (supplier_quote_id) REFERENCES supplier_quotes(id) ON DELETE SET NULL,
  CONSTRAINT fk_purchase_orders_buyer FOREIGN KEY (buyer_id) REFERENCES employees(id) ON DELETE CASCADE,
  CONSTRAINT fk_purchase_orders_approved_by FOREIGN KEY (approved_by) REFERENCES employees(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_purchase_orders_deleted_at ON purchase_orders (deleted_at);
CREATE INDEX IF NOT EXISTS idx_purchase_orders_is_active ON purchase_orders (is_active);
CREATE INDEX IF NOT EXISTS idx_purchase_orders_code ON purchase_orders (code);
CREATE INDEX IF NOT EXISTS idx_purchase_orders_supplier_id ON purchase_orders (supplier_id);
CREATE INDEX IF NOT EXISTS idx_purchase_orders_supplier_quote_id ON purchase_orders (supplier_quote_id);
CREATE INDEX IF NOT EXISTS idx_purchase_orders_buyer_id ON purchase_orders (buyer_id);
CREATE INDEX IF NOT EXISTS idx_purchase_orders_status ON purchase_orders (status);
CREATE INDEX IF NOT EXISTS idx_purchase_orders_order_date ON purchase_orders (order_date);
CREATE INDEX IF NOT EXISTS idx_purchase_orders_required_date ON purchase_orders (required_date);

-- 采购订单明细表
CREATE TABLE IF NOT EXISTS purchase_order_items (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL,
  purchase_order_id BIGINT NOT NULL,
  item_id BIGINT NOT NULL,
  quantity DECIMAL(15,4) NOT NULL,
  unit_price DECIMAL(15,2) NOT NULL,
  discount_rate DECIMAL(5,2) DEFAULT 0,
  discount_amount DECIMAL(15,2) DEFAULT 0,
  line_total DECIMAL(15,2) NOT NULL,
  received_quantity DECIMAL(15,4) DEFAULT 0,
  remaining_quantity DECIMAL(15,4) NOT NULL,
  notes TEXT NULL,
  CONSTRAINT fk_purchase_order_items_purchase_order FOREIGN KEY (purchase_order_id) REFERENCES purchase_orders(id) ON DELETE CASCADE,
  CONSTRAINT fk_purchase_order_items_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_purchase_order_items_deleted_at ON purchase_order_items (deleted_at);
CREATE INDEX IF NOT EXISTS idx_purchase_order_items_purchase_order_id ON purchase_order_items (purchase_order_id);
CREATE INDEX IF NOT EXISTS idx_purchase_order_items_item_id ON purchase_order_items (item_id);

-- ============================================================================
-- 收货管理
-- ============================================================================

-- 收货单表
CREATE TABLE IF NOT EXISTS receipts (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL,
  created_by BIGINT NULL,
  updated_by BIGINT NULL,
  is_active BOOLEAN DEFAULT TRUE,
  code VARCHAR(50) NOT NULL,
  purchase_order_id BIGINT NOT NULL,
  supplier_id BIGINT NOT NULL,
  warehouse_id BIGINT NOT NULL,
  receipt_date DATE NOT NULL,
  delivery_note_number VARCHAR(100) NULL,
  carrier VARCHAR(255) NULL,
  tracking_number VARCHAR(255) NULL,
  received_by BIGINT NULL,
  notes TEXT NULL,
  status VARCHAR(20) DEFAULT 'PENDING',
  CONSTRAINT uq_receipts_code UNIQUE (code),
  CONSTRAINT fk_receipts_purchase_order FOREIGN KEY (purchase_order_id) REFERENCES purchase_orders(id) ON DELETE CASCADE,
  CONSTRAINT fk_receipts_supplier FOREIGN KEY (supplier_id) REFERENCES suppliers(id) ON DELETE CASCADE,
  CONSTRAINT fk_receipts_warehouse FOREIGN KEY (warehouse_id) REFERENCES warehouses(id) ON DELETE CASCADE,
  CONSTRAINT fk_receipts_received_by FOREIGN KEY (received_by) REFERENCES employees(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_receipts_deleted_at ON receipts (deleted_at);
CREATE INDEX IF NOT EXISTS idx_receipts_is_active ON receipts (is_active);
CREATE INDEX IF NOT EXISTS idx_receipts_code ON receipts (code);
CREATE INDEX IF NOT EXISTS idx_receipts_purchase_order_id ON receipts (purchase_order_id);
CREATE INDEX IF NOT EXISTS idx_receipts_supplier_id ON receipts (supplier_id);
CREATE INDEX IF NOT EXISTS idx_receipts_warehouse_id ON receipts (warehouse_id);
CREATE INDEX IF NOT EXISTS idx_receipts_status ON receipts (status);
CREATE INDEX IF NOT EXISTS idx_receipts_receipt_date ON receipts (receipt_date);

-- 收货单明细表
CREATE TABLE IF NOT EXISTS receipt_items (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL,
  receipt_id BIGINT NOT NULL,
  purchase_order_item_id BIGINT NOT NULL,
  item_id BIGINT NOT NULL,
  ordered_quantity DECIMAL(15,4) NOT NULL,
  received_quantity DECIMAL(15,4) NOT NULL,
  rejected_quantity DECIMAL(15,4) DEFAULT 0,
  unit_cost DECIMAL(15,2) NOT NULL,
  total_cost DECIMAL(15,2) NOT NULL,
  batch_number VARCHAR(100) NULL,
  expiry_date DATE NULL,
  location_id BIGINT NULL,
  quality_status VARCHAR(20) DEFAULT 'PENDING',
  notes TEXT NULL,
  CONSTRAINT fk_receipt_items_receipt FOREIGN KEY (receipt_id) REFERENCES receipts(id) ON DELETE CASCADE,
  CONSTRAINT fk_receipt_items_purchase_order_item FOREIGN KEY (purchase_order_item_id) REFERENCES purchase_order_items(id) ON DELETE CASCADE,
  CONSTRAINT fk_receipt_items_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE,
  CONSTRAINT fk_receipt_items_location FOREIGN KEY (location_id) REFERENCES locations(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_receipt_items_deleted_at ON receipt_items (deleted_at);
CREATE INDEX IF NOT EXISTS idx_receipt_items_receipt_id ON receipt_items (receipt_id);
CREATE INDEX IF NOT EXISTS idx_receipt_items_purchase_order_item_id ON receipt_items (purchase_order_item_id);
CREATE INDEX IF NOT EXISTS idx_receipt_items_item_id ON receipt_items (item_id);
CREATE INDEX IF NOT EXISTS idx_receipt_items_location_id ON receipt_items (location_id);
CREATE INDEX IF NOT EXISTS idx_receipt_items_quality_status ON receipt_items (quality_status);

-- ============================================================================
-- 采购发票管理
-- ============================================================================

-- 采购发票表
CREATE TABLE IF NOT EXISTS purchase_invoices (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL,
  created_by BIGINT NULL,
  updated_by BIGINT NULL,
  is_active BOOLEAN DEFAULT TRUE,
  code VARCHAR(50) NOT NULL,
  supplier_id BIGINT NOT NULL,
  purchase_order_id BIGINT NULL,
  supplier_invoice_number VARCHAR(100) NOT NULL,
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
  status VARCHAR(20) DEFAULT 'RECEIVED',
  CONSTRAINT uq_purchase_invoices_code UNIQUE (code),
  CONSTRAINT uq_purchase_invoices_supplier_number UNIQUE (supplier_id, supplier_invoice_number),
  CONSTRAINT fk_purchase_invoices_supplier FOREIGN KEY (supplier_id) REFERENCES suppliers(id) ON DELETE CASCADE,
  CONSTRAINT fk_purchase_invoices_purchase_order FOREIGN KEY (purchase_order_id) REFERENCES purchase_orders(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_purchase_invoices_deleted_at ON purchase_invoices (deleted_at);
CREATE INDEX IF NOT EXISTS idx_purchase_invoices_is_active ON purchase_invoices (is_active);
CREATE INDEX IF NOT EXISTS idx_purchase_invoices_code ON purchase_invoices (code);
CREATE INDEX IF NOT EXISTS idx_purchase_invoices_supplier_id ON purchase_invoices (supplier_id);
CREATE INDEX IF NOT EXISTS idx_purchase_invoices_purchase_order_id ON purchase_invoices (purchase_order_id);
CREATE INDEX IF NOT EXISTS idx_purchase_invoices_status ON purchase_invoices (status);
CREATE INDEX IF NOT EXISTS idx_purchase_invoices_invoice_date ON purchase_invoices (invoice_date);
CREATE INDEX IF NOT EXISTS idx_purchase_invoices_due_date ON purchase_invoices (due_date);

-- 采购发票明细表
CREATE TABLE IF NOT EXISTS purchase_invoice_items (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL,
  purchase_invoice_id BIGINT NOT NULL,
  item_id BIGINT NOT NULL,
  quantity DECIMAL(15,4) NOT NULL,
  unit_price DECIMAL(15,2) NOT NULL,
  discount_rate DECIMAL(5,2) DEFAULT 0,
  discount_amount DECIMAL(15,2) DEFAULT 0,
  line_total DECIMAL(15,2) NOT NULL,
  notes TEXT NULL,
  CONSTRAINT fk_purchase_invoice_items_purchase_invoice FOREIGN KEY (purchase_invoice_id) REFERENCES purchase_invoices(id) ON DELETE CASCADE,
  CONSTRAINT fk_purchase_invoice_items_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_purchase_invoice_items_deleted_at ON purchase_invoice_items (deleted_at);
CREATE INDEX IF NOT EXISTS idx_purchase_invoice_items_purchase_invoice_id ON purchase_invoice_items (purchase_invoice_id);
CREATE INDEX IF NOT EXISTS idx_purchase_invoice_items_item_id ON purchase_invoice_items (item_id);

-- ============================================================================
-- 初始数据插入
-- ============================================================================

-- 插入默认供应商
INSERT INTO suppliers (code, name, description, type, phone, email, address, credit_limit, payment_terms, lead_time_days, status, is_active) VALUES
('SUPP_DEFAULT', '默认供应商', '系统默认供应商', 'VENDOR', '400-000-0000', 'supplier@example.com', '默认供应商地址', 500000.00, 30, 7, 'ACTIVE', TRUE)
ON CONFLICT (code) DO NOTHING;

-- ============================================================================
-- 结束
-- ============================================================================