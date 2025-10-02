-- ============================================================================
-- GalaxyERP 完整数据库初始化脚本 - SQLite 版本 (更新版)
-- 生成时间: 2025-01-01
-- 说明: 基于最新的PostgreSQL脚本，包含deleted_at、created_by、updated_by字段
-- 版本: 2.0.0
-- ============================================================================

-- 暂时禁用外键约束以避免插入顺序问题
PRAGMA foreign_keys = OFF;

-- ============================================================================
-- 1. 用户管理模块
-- ============================================================================

-- 公司表
CREATE TABLE IF NOT EXISTS companies (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  deleted_at DATETIME NULL,
  created_by INTEGER NULL,
  updated_by INTEGER NULL,
  code TEXT NOT NULL UNIQUE,
  name TEXT NOT NULL,
  description TEXT,
  address TEXT,
  phone TEXT,
  email TEXT,
  website TEXT,
  tax_number TEXT,
  legal_representative TEXT,
  registration_date DATE,
  status TEXT DEFAULT 'ACTIVE',
  is_active INTEGER DEFAULT 1,
  FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 部门表
CREATE TABLE IF NOT EXISTS departments (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  deleted_at DATETIME NULL,
  created_by INTEGER NULL,
  updated_by INTEGER NULL,
  code TEXT NOT NULL UNIQUE,
  name TEXT NOT NULL,
  description TEXT,
  company_id INTEGER NOT NULL,
  parent_id INTEGER,
  manager_id INTEGER,
  level INTEGER DEFAULT 1,
  sort_order INTEGER DEFAULT 0,
  status TEXT DEFAULT 'ACTIVE',
  is_active INTEGER DEFAULT 1,
  FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE,
  FOREIGN KEY (parent_id) REFERENCES departments(id) ON DELETE SET NULL,
  FOREIGN KEY (manager_id) REFERENCES users(id) ON DELETE SET NULL,
  FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 用户表
CREATE TABLE IF NOT EXISTS users (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  deleted_at DATETIME NULL,
  created_by INTEGER NULL,
  updated_by INTEGER NULL,
  username TEXT NOT NULL UNIQUE,
  email TEXT NOT NULL UNIQUE,
  phone TEXT,
  password_hash TEXT NOT NULL,
  salt TEXT,
  first_name TEXT,
  last_name TEXT,
  display_name TEXT,
  avatar_url TEXT,
  company_id INTEGER,
  department_id INTEGER,
  is_admin INTEGER DEFAULT 0,
  position TEXT,
  employee_number TEXT,
  hire_date DATE,
  birth_date DATE,
  gender TEXT,
  address TEXT,
  emergency_contact TEXT,
  emergency_phone TEXT,
  last_login_at DATETIME,
  last_login_ip TEXT,
  login_count INTEGER DEFAULT 0,
  failed_login_count INTEGER DEFAULT 0,
  locked_until DATETIME,
  password_changed_at DATETIME,
  must_change_password INTEGER DEFAULT 0,
  two_factor_enabled INTEGER DEFAULT 0,
  two_factor_secret TEXT,
  preferences TEXT,
  timezone TEXT DEFAULT 'Asia/Shanghai',
  language TEXT DEFAULT 'zh-CN',
  theme TEXT DEFAULT 'light',
  status TEXT DEFAULT 'ACTIVE',
  is_active INTEGER DEFAULT 1,
  FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE SET NULL,
  FOREIGN KEY (department_id) REFERENCES departments(id) ON DELETE SET NULL,
  FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 角色表
CREATE TABLE IF NOT EXISTS roles (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  deleted_at DATETIME NULL,
  created_by INTEGER NULL,
  updated_by INTEGER NULL,
  name TEXT NOT NULL UNIQUE,
  code TEXT NOT NULL UNIQUE,
  description TEXT,
  is_system INTEGER DEFAULT 0,
  is_active INTEGER DEFAULT 1,
  FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 权限表
CREATE TABLE IF NOT EXISTS permissions (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  deleted_at DATETIME NULL,
  created_by INTEGER NULL,
  updated_by INTEGER NULL,
  name TEXT NOT NULL UNIQUE,
  code TEXT NOT NULL UNIQUE,
  description TEXT,
  resource TEXT NOT NULL,
  action TEXT NOT NULL,
  is_system INTEGER DEFAULT 0,
  FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 用户角色关联表
CREATE TABLE IF NOT EXISTS user_roles (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  deleted_at DATETIME NULL,
  created_by INTEGER NULL,
  updated_by INTEGER NULL,
  user_id INTEGER NOT NULL,
  role_id INTEGER NOT NULL,
  granted_by INTEGER,
  granted_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  expires_at DATETIME,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
  FOREIGN KEY (granted_by) REFERENCES users(id) ON DELETE SET NULL,
  FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL,
  UNIQUE(user_id, role_id)
);

-- 角色权限关联表
CREATE TABLE IF NOT EXISTS role_permissions (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  deleted_at DATETIME NULL,
  created_by INTEGER NULL,
  updated_by INTEGER NULL,
  role_id INTEGER NOT NULL,
  permission_id INTEGER NOT NULL,
  granted_by INTEGER,
  granted_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
  FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE,
  FOREIGN KEY (granted_by) REFERENCES users(id) ON DELETE SET NULL,
  FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL,
  UNIQUE(role_id, permission_id)
);

-- 用户会话表
CREATE TABLE IF NOT EXISTS user_sessions (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  deleted_at DATETIME NULL,
  user_id INTEGER NOT NULL,
  session_token TEXT NOT NULL UNIQUE,
  refresh_token TEXT,
  ip_address TEXT,
  user_agent TEXT,
  expires_at DATETIME NOT NULL,
  last_activity DATETIME DEFAULT CURRENT_TIMESTAMP,
  is_active INTEGER DEFAULT 1,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- ============================================================================
-- 2. 库存管理模块
-- ============================================================================

-- 物料表
CREATE TABLE IF NOT EXISTS items (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  deleted_at DATETIME NULL,
  created_by INTEGER NULL,
  updated_by INTEGER NULL,
  code TEXT NOT NULL UNIQUE,
  name TEXT NOT NULL,
  description TEXT,
  category TEXT,
  subcategory TEXT,
  brand TEXT,
  model TEXT,
  specification TEXT,
  unit_of_measure TEXT DEFAULT 'PCS',
  weight REAL DEFAULT 0,
  dimensions TEXT,
  color TEXT,
  material TEXT,
  origin_country TEXT,
  hs_code TEXT,
  barcode TEXT,
  qr_code TEXT,
  min_stock_level REAL DEFAULT 0,
  max_stock_level REAL DEFAULT 0,
  reorder_point REAL DEFAULT 0,
  reorder_quantity REAL DEFAULT 0,
  lead_time_days INTEGER DEFAULT 0,
  shelf_life_days INTEGER,
  storage_conditions TEXT,
  is_serialized INTEGER DEFAULT 0,
  is_batch_tracked INTEGER DEFAULT 0,
  is_active INTEGER DEFAULT 1,
  item_type TEXT DEFAULT 'STOCK',
  valuation_method TEXT DEFAULT 'FIFO',
  standard_cost REAL DEFAULT 0,
  last_purchase_rate REAL DEFAULT 0,
  average_cost REAL DEFAULT 0,
  FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 仓库表
CREATE TABLE IF NOT EXISTS warehouses (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  deleted_at DATETIME NULL,
  created_by INTEGER NULL,
  updated_by INTEGER NULL,
  code TEXT NOT NULL UNIQUE,
  name TEXT NOT NULL,
  description TEXT,
  address TEXT,
  city TEXT,
  state TEXT,
  postal_code TEXT,
  country TEXT,
  phone TEXT,
  email TEXT,
  manager_id INTEGER,
  warehouse_type TEXT DEFAULT 'GENERAL',
  is_active INTEGER DEFAULT 1,
  FOREIGN KEY (manager_id) REFERENCES users(id) ON DELETE SET NULL,
  FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 库存表
CREATE TABLE IF NOT EXISTS stocks (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  deleted_at DATETIME NULL,
  created_by INTEGER NULL,
  updated_by INTEGER NULL,
  item_id INTEGER NOT NULL,
  warehouse_id INTEGER NOT NULL,
  quantity REAL DEFAULT 0,
  reserved_quantity REAL DEFAULT 0,
  available_quantity REAL DEFAULT 0,
  cost_per_unit REAL DEFAULT 0,
  total_value REAL DEFAULT 0,
  last_updated DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE,
  FOREIGN KEY (warehouse_id) REFERENCES warehouses(id) ON DELETE CASCADE,
  FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL,
  UNIQUE(item_id, warehouse_id)
);

-- 库存移动表
CREATE TABLE IF NOT EXISTS movements (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  deleted_at DATETIME NULL,
  created_by INTEGER NULL,
  updated_by INTEGER NULL,
  item_id INTEGER NOT NULL,
  warehouse_id INTEGER NOT NULL,
  movement_type TEXT NOT NULL,
  quantity REAL NOT NULL,
  unit_cost REAL DEFAULT 0,
  total_cost REAL DEFAULT 0,
  reference_type TEXT,
  reference_id INTEGER,
  reference_number TEXT,
  batch_no TEXT,
  serial_no TEXT,
  expiry_date DATE,
  notes TEXT,
  movement_date DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE,
  FOREIGN KEY (warehouse_id) REFERENCES warehouses(id) ON DELETE CASCADE,
  FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 库存调整表
CREATE TABLE IF NOT EXISTS stock_adjustments (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  deleted_at DATETIME NULL,
  created_by INTEGER NULL,
  updated_by INTEGER NULL,
  adjustment_number TEXT NOT NULL UNIQUE,
  warehouse_id INTEGER NOT NULL,
  adjustment_date DATETIME NOT NULL,
  reason TEXT,
  status TEXT DEFAULT 'Draft',
  total_amount REAL DEFAULT 0,
  approved_by INTEGER,
  approved_at DATETIME,
  notes TEXT,
  FOREIGN KEY (warehouse_id) REFERENCES warehouses(id) ON DELETE CASCADE,
  FOREIGN KEY (approved_by) REFERENCES users(id) ON DELETE SET NULL,
  FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 库存调整明细表
CREATE TABLE IF NOT EXISTS stock_adjustment_items (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  deleted_at DATETIME NULL,
  created_by INTEGER NULL,
  updated_by INTEGER NULL,
  stock_adjustment_id INTEGER NOT NULL,
  item_id INTEGER NOT NULL,
  current_quantity REAL DEFAULT 0,
  adjusted_quantity REAL NOT NULL,
  difference_quantity REAL NOT NULL,
  unit_cost REAL DEFAULT 0,
  total_cost REAL DEFAULT 0,
  reason TEXT,
  batch_no TEXT,
  serial_no TEXT,
  FOREIGN KEY (stock_adjustment_id) REFERENCES stock_adjustments(id) ON DELETE CASCADE,
  FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE,
  FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 库存转移表
CREATE TABLE IF NOT EXISTS stock_transfers (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  deleted_at DATETIME NULL,
  created_by INTEGER NULL,
  updated_by INTEGER NULL,
  transfer_number TEXT NOT NULL UNIQUE,
  from_warehouse_id INTEGER NOT NULL,
  to_warehouse_id INTEGER NOT NULL,
  transfer_date DATETIME NOT NULL,
  status TEXT DEFAULT 'Draft',
  total_quantity REAL DEFAULT 0,
  notes TEXT,
  FOREIGN KEY (from_warehouse_id) REFERENCES warehouses(id) ON DELETE CASCADE,
  FOREIGN KEY (to_warehouse_id) REFERENCES warehouses(id) ON DELETE CASCADE,
  FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- ============================================================================
-- 3. 采购管理模块
-- ============================================================================

-- 供应商表
CREATE TABLE IF NOT EXISTS suppliers (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  deleted_at DATETIME NULL,
  created_by INTEGER NULL,
  updated_by INTEGER NULL,
  name TEXT NOT NULL,
  code TEXT NOT NULL UNIQUE,
  email TEXT,
  phone TEXT,
  address TEXT,
  city TEXT,
  state TEXT,
  postal_code TEXT,
  country TEXT,
  contact_person TEXT,
  credit_limit REAL DEFAULT 0,
  supplier_group TEXT,
  territory TEXT,
  quality_rating REAL DEFAULT 0,
  delivery_rating REAL DEFAULT 0,
  is_active INTEGER DEFAULT 1,
  FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 采购申请表
CREATE TABLE IF NOT EXISTS purchase_requests (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  deleted_at DATETIME NULL,
  created_by INTEGER NULL,
  updated_by INTEGER NULL,
  request_number TEXT NOT NULL UNIQUE,
  request_date DATETIME NOT NULL,
  required_by DATETIME NOT NULL,
  department TEXT,
  status TEXT DEFAULT 'Draft',
  notes TEXT,
  approved_by INTEGER,
  FOREIGN KEY (approved_by) REFERENCES users(id) ON DELETE SET NULL,
  FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 采购申请明细表
CREATE TABLE IF NOT EXISTS purchase_request_items (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  deleted_at DATETIME NULL,
  created_by INTEGER NULL,
  updated_by INTEGER NULL,
  purchase_request_id INTEGER NOT NULL,
  item_id INTEGER NOT NULL,
  description TEXT,
  quantity REAL DEFAULT 1,
  uom TEXT,
  estimated_cost REAL DEFAULT 0,
  notes TEXT,
  FOREIGN KEY (purchase_request_id) REFERENCES purchase_requests(id) ON DELETE CASCADE,
  FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE,
  FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 采购订单表
CREATE TABLE IF NOT EXISTS purchase_orders (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  deleted_at DATETIME NULL,
  created_by INTEGER NULL,
  updated_by INTEGER NULL,
  order_number TEXT NOT NULL UNIQUE,
  supplier_id INTEGER NOT NULL,
  order_date DATETIME NOT NULL,
  delivery_date DATETIME NOT NULL,
  status TEXT DEFAULT 'Draft',
  purchase_request_id INTEGER,
  total_amount REAL DEFAULT 0,
  discount_amount REAL DEFAULT 0,
  tax_amount REAL DEFAULT 0,
  grand_total REAL DEFAULT 0,
  terms TEXT,
  notes TEXT,
  FOREIGN KEY (supplier_id) REFERENCES suppliers(id) ON DELETE CASCADE,
  FOREIGN KEY (purchase_request_id) REFERENCES purchase_requests(id) ON DELETE SET NULL,
  FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 采购订单明细表
CREATE TABLE IF NOT EXISTS purchase_order_items (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  deleted_at DATETIME NULL,
  created_by INTEGER NULL,
  updated_by INTEGER NULL,
  purchase_order_id INTEGER NOT NULL,
  item_id INTEGER NOT NULL,
  description TEXT,
  quantity REAL DEFAULT 1,
  received_qty REAL DEFAULT 0,
  rate REAL DEFAULT 0,
  amount REAL DEFAULT 0,
  discount_rate REAL DEFAULT 0,
  discount_amount REAL DEFAULT 0,
  tax_rate REAL DEFAULT 0,
  tax_amount REAL DEFAULT 0,
  total_amount REAL DEFAULT 0,
  warehouse_id INTEGER,
  FOREIGN KEY (purchase_order_id) REFERENCES purchase_orders(id) ON DELETE CASCADE,
  FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE,
  FOREIGN KEY (warehouse_id) REFERENCES warehouses(id) ON DELETE SET NULL,
  FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 收货单表
CREATE TABLE IF NOT EXISTS purchase_receipts (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  deleted_at DATETIME NULL,
  created_by INTEGER NULL,
  updated_by INTEGER NULL,
  receipt_number TEXT NOT NULL UNIQUE,
  supplier_id INTEGER NOT NULL,
  purchase_order_id INTEGER,
  date DATETIME NOT NULL,
  status TEXT DEFAULT 'Draft',
  total_quantity REAL DEFAULT 0,
  transporter TEXT,
  driver_name TEXT,
  vehicle_number TEXT,
  destination TEXT,
  notes TEXT,
  FOREIGN KEY (supplier_id) REFERENCES suppliers(id) ON DELETE CASCADE,
  FOREIGN KEY (purchase_order_id) REFERENCES purchase_orders(id) ON DELETE SET NULL,
  FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 收货单明细表
CREATE TABLE IF NOT EXISTS purchase_receipt_items (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  deleted_at DATETIME NULL,
  created_by INTEGER NULL,
  updated_by INTEGER NULL,
  purchase_receipt_id INTEGER NOT NULL,
  purchase_order_item_id INTEGER,
  item_id INTEGER NOT NULL,
  description TEXT,
  quantity REAL DEFAULT 1,
  accepted_qty REAL DEFAULT 0,
  rejected_qty REAL DEFAULT 0,
  batch_no TEXT,
  serial_no TEXT,
  warehouse_id INTEGER,
  quality_status TEXT DEFAULT 'Pending',
  notes TEXT,
  FOREIGN KEY (purchase_receipt_id) REFERENCES purchase_receipts(id) ON DELETE CASCADE,
  FOREIGN KEY (purchase_order_item_id) REFERENCES purchase_order_items(id) ON DELETE SET NULL,
  FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE,
  FOREIGN KEY (warehouse_id) REFERENCES warehouses(id) ON DELETE SET NULL,
  FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- ============================================================================
-- 4. 销售管理模块
-- ============================================================================

-- 客户表
CREATE TABLE IF NOT EXISTS customers (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  deleted_at DATETIME NULL,
  created_by INTEGER NULL,
  updated_by INTEGER NULL,
  name TEXT NOT NULL,
  code TEXT NOT NULL UNIQUE,
  email TEXT,
  phone TEXT,
  address TEXT,
  city TEXT,
  state TEXT,
  postal_code TEXT,
  country TEXT,
  contact_person TEXT,
  credit_limit REAL DEFAULT 0,
  customer_group TEXT,
  territory TEXT,
  is_active INTEGER DEFAULT 1,
  FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 报价单表
CREATE TABLE IF NOT EXISTS quotations (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  deleted_at DATETIME NULL,
  created_by INTEGER NULL,
  updated_by INTEGER NULL,
  quotation_number TEXT NOT NULL UNIQUE,
  customer_id INTEGER NOT NULL,
  date DATETIME NOT NULL,
  valid_till DATETIME NOT NULL,
  status TEXT DEFAULT 'Draft',
  subject TEXT,
  total_amount REAL DEFAULT 0,
  discount_amount REAL DEFAULT 0,
  tax_amount REAL DEFAULT 0,
  grand_total REAL DEFAULT 0,
  terms TEXT,
  notes TEXT,
  FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE,
  FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 报价单明细表
CREATE TABLE IF NOT EXISTS quotation_items (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  deleted_at DATETIME NULL,
  created_by INTEGER NULL,
  updated_by INTEGER NULL,
  quotation_id INTEGER NOT NULL,
  item_id INTEGER NOT NULL,
  description TEXT,
  quantity REAL DEFAULT 1,
  rate REAL DEFAULT 0,
  amount REAL DEFAULT 0,
  discount_rate REAL DEFAULT 0,
  discount_amount REAL DEFAULT 0,
  tax_rate REAL DEFAULT 0,
  tax_amount REAL DEFAULT 0,
  total_amount REAL DEFAULT 0,
  FOREIGN KEY (quotation_id) REFERENCES quotations(id) ON DELETE CASCADE,
  FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE,
  FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 销售订单表
CREATE TABLE IF NOT EXISTS sales_orders (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  deleted_at DATETIME NULL,
  created_by INTEGER NULL,
  updated_by INTEGER NULL,
  order_number TEXT NOT NULL UNIQUE,
  customer_id INTEGER NOT NULL,
  date DATETIME NOT NULL,
  delivery_date DATETIME NOT NULL,
  status TEXT DEFAULT 'Draft',
  quotation_id INTEGER,
  total_amount REAL DEFAULT 0,
  discount_amount REAL DEFAULT 0,
  tax_amount REAL DEFAULT 0,
  grand_total REAL DEFAULT 0,
  terms TEXT,
  notes TEXT,
  FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE,
  FOREIGN KEY (quotation_id) REFERENCES quotations(id) ON DELETE SET NULL,
  FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 销售订单明细表
CREATE TABLE IF NOT EXISTS sales_order_items (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  deleted_at DATETIME NULL,
  created_by INTEGER NULL,
  updated_by INTEGER NULL,
  sales_order_id INTEGER NOT NULL,
  item_id INTEGER NOT NULL,
  description TEXT,
  quantity REAL DEFAULT 1,
  delivered_qty REAL DEFAULT 0,
  rate REAL DEFAULT 0,
  amount REAL DEFAULT 0,
  discount_rate REAL DEFAULT 0,
  discount_amount REAL DEFAULT 0,
  tax_rate REAL DEFAULT 0,
  tax_amount REAL DEFAULT 0,
  total_amount REAL DEFAULT 0,
  warehouse_id INTEGER,
  FOREIGN KEY (sales_order_id) REFERENCES sales_orders(id) ON DELETE CASCADE,
  FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE,
  FOREIGN KEY (warehouse_id) REFERENCES warehouses(id) ON DELETE SET NULL,
  FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 送货单表
CREATE TABLE IF NOT EXISTS delivery_notes (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  deleted_at DATETIME NULL,
  created_by INTEGER NULL,
  updated_by INTEGER NULL,
  delivery_number TEXT NOT NULL UNIQUE,
  customer_id INTEGER NOT NULL,
  sales_order_id INTEGER,
  date DATETIME NOT NULL,
  status TEXT DEFAULT 'Draft',
  total_quantity REAL DEFAULT 0,
  transporter TEXT,
  driver_name TEXT,
  vehicle_number TEXT,
  destination TEXT,
  notes TEXT,
  FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE,
  FOREIGN KEY (sales_order_id) REFERENCES sales_orders(id) ON DELETE SET NULL,
  FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 送货单明细表
CREATE TABLE IF NOT EXISTS delivery_note_items (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  deleted_at DATETIME NULL,
  created_by INTEGER NULL,
  updated_by INTEGER NULL,
  delivery_note_id INTEGER NOT NULL,
  sales_order_item_id INTEGER,
  item_id INTEGER NOT NULL,
  description TEXT,
  quantity REAL NOT NULL,
  batch_no TEXT,
  serial_no TEXT,
  warehouse_id INTEGER,
  FOREIGN KEY (delivery_note_id) REFERENCES delivery_notes(id) ON DELETE CASCADE,
  FOREIGN KEY (sales_order_item_id) REFERENCES sales_order_items(id) ON DELETE SET NULL,
  FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE,
  FOREIGN KEY (warehouse_id) REFERENCES warehouses(id) ON DELETE SET NULL,
  FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- ============================================================================
-- 5. 系统管理模块
-- ============================================================================

-- 系统配置表
CREATE TABLE IF NOT EXISTS system_configs (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  deleted_at DATETIME NULL,
  created_by INTEGER NULL,
  updated_by INTEGER NULL,
  config_key TEXT NOT NULL UNIQUE,
  config_value TEXT NOT NULL,
  config_type TEXT DEFAULT 'STRING',
  category TEXT DEFAULT 'GENERAL',
  description TEXT,
  is_encrypted INTEGER DEFAULT 0,
  is_system INTEGER DEFAULT 0,
  default_value TEXT,
  validation_rule TEXT,
  display_order INTEGER DEFAULT 0,
  FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 审计日志表
CREATE TABLE IF NOT EXISTS audit_logs (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  user_id INTEGER,
  username TEXT,
  action TEXT NOT NULL,
  resource_type TEXT NOT NULL,
  resource_id INTEGER,
  resource_name TEXT,
  old_values TEXT,
  new_values TEXT,
  ip_address TEXT,
  user_agent TEXT,
  request_id TEXT,
  session_id TEXT,
  result TEXT DEFAULT 'SUCCESS',
  error_message TEXT,
  duration_ms INTEGER,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
);

-- 数据权限表
CREATE TABLE IF NOT EXISTS data_permissions (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  deleted_at DATETIME NULL,
  created_by INTEGER NULL,
  updated_by INTEGER NULL,
  user_id INTEGER NOT NULL,
  resource_type TEXT NOT NULL,
  resource_id INTEGER,
  permission_type TEXT NOT NULL,
  granted_by INTEGER,
  granted_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  expires_at DATETIME,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL,
  FOREIGN KEY (granted_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 审批步骤表
CREATE TABLE IF NOT EXISTS approval_steps (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  deleted_at DATETIME NULL,
  created_by INTEGER NULL,
  updated_by INTEGER NULL,
  workflow_type TEXT NOT NULL,
  step_order INTEGER NOT NULL,
  step_name TEXT NOT NULL,
  approver_type TEXT NOT NULL,
  approver_id INTEGER,
  is_required INTEGER DEFAULT 1,
  conditions TEXT,
  FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL,
  FOREIGN KEY (approver_id) REFERENCES users(id) ON DELETE SET NULL
);

-- ============================================================================
-- 初始数据插入
-- ============================================================================

-- 插入默认公司
INSERT OR IGNORE INTO companies (code, name, description, status) VALUES
('DEFAULT', '默认公司', '系统默认公司', 'ACTIVE');

-- 插入默认部门
INSERT OR IGNORE INTO departments (code, name, description, company_id, status) VALUES
('DEFAULT', '默认部门', '系统默认部门', 1, 'ACTIVE');

-- 插入默认角色
INSERT OR IGNORE INTO roles (name, code, description, is_system) VALUES
('系统管理员', 'ADMIN', '系统管理员角色', 1),
('普通用户', 'USER', '普通用户角色', 1);

-- 插入默认权限
INSERT OR IGNORE INTO permissions (name, code, description, resource, action, is_system) VALUES
('用户管理', 'USER_MANAGE', '用户管理权限', 'user', 'manage', 1),
('库存查看', 'INVENTORY_VIEW', '库存查看权限', 'inventory', 'view', 1),
('销售管理', 'SALES_MANAGE', '销售管理权限', 'sales', 'manage', 1),
('采购管理', 'PURCHASE_MANAGE', '采购管理权限', 'purchase', 'manage', 1);

-- 插入默认用户 (密码: admin123)
INSERT OR IGNORE INTO users (username, email, password_hash, first_name, last_name, company_id, department_id, is_admin, status) VALUES
('admin', 'admin@galaxyerp.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', '系统', '管理员', 1, 1, 1, 'ACTIVE');

-- 插入默认仓库
INSERT OR IGNORE INTO warehouses (code, name, description, warehouse_type) VALUES
('DEFAULT', '默认仓库', '系统默认仓库', 'GENERAL');

-- 插入默认供应商
INSERT OR IGNORE INTO suppliers (name, code, supplier_group, territory, phone, email, address) VALUES
('默认供应商', 'DEFAULT001', 'Default', 'Default', '400-000-0000', 'default@example.com', '默认地址');

-- 插入默认客户
INSERT OR IGNORE INTO customers (name, code, customer_group, territory, phone, email, address) VALUES
('默认客户', 'DEFAULT001', 'Default', 'Default', '400-000-0000', 'default@example.com', '默认地址');

-- 插入系统配置
INSERT OR IGNORE INTO system_configs (config_key, config_value, config_type, category, description, is_system, default_value, display_order) VALUES
('system.name', 'GalaxyERP', 'STRING', 'BASIC', '系统名称', 1, 'GalaxyERP', 1),
('system.version', '1.0.0', 'STRING', 'BASIC', '系统版本', 1, '1.0.0', 2),
('system.timezone', 'Asia/Shanghai', 'STRING', 'BASIC', '系统时区', 1, 'Asia/Shanghai', 3),
('system.language', 'zh-CN', 'STRING', 'BASIC', '系统语言', 1, 'zh-CN', 4),
('security.password.min_length', '8', 'INTEGER', 'SECURITY', '密码最小长度', 1, '8', 10),
('security.session.timeout', '3600', 'INTEGER', 'SECURITY', '会话超时时间(秒)', 1, '3600', 15);

-- 启用外键约束
PRAGMA foreign_keys = ON;

-- ============================================================================
-- 结束
-- ============================================================================