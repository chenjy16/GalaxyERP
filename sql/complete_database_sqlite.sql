-- ============================================================================
-- GalaxyERP 完整数据库初始化脚本 - SQLite 版本
-- 生成时间: 2025-01-01
-- 说明: 包含所有101个表的完整建表脚本，基于Go模型结构
-- 版本: 1.0.0
-- ============================================================================

-- 暂时禁用外键约束以避免插入顺序问题
PRAGMA foreign_keys = OFF;

-- ============================================================================
-- 1. 用户管理模块 (8个表)
-- ============================================================================

-- 公司表
CREATE TABLE IF NOT EXISTS companies (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  created_by INTEGER,
  updated_by INTEGER,
  is_active INTEGER DEFAULT 1,
  code TEXT NOT NULL,
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
  CONSTRAINT uq_companies_code UNIQUE (code)
);

-- 部门表
CREATE TABLE IF NOT EXISTS departments (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  created_by INTEGER,
  updated_by INTEGER,
  is_active INTEGER DEFAULT 1,
  code TEXT NOT NULL,
  name TEXT NOT NULL,
  description TEXT,
  company_id INTEGER NOT NULL,
  parent_id INTEGER,
  manager_id INTEGER,
  level INTEGER DEFAULT 1,
  sort_order INTEGER DEFAULT 0,
  status TEXT DEFAULT 'ACTIVE',
  CONSTRAINT uq_departments_code UNIQUE (code),
  CONSTRAINT fk_departments_company FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE,
  CONSTRAINT fk_departments_parent FOREIGN KEY (parent_id) REFERENCES departments(id) ON DELETE SET NULL
);

-- 用户表
CREATE TABLE IF NOT EXISTS users (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  created_by INTEGER,
  updated_by INTEGER,
  is_active INTEGER DEFAULT 1,
  username TEXT NOT NULL,
  email TEXT NOT NULL,
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
  status TEXT DEFAULT 'ACTIVE',
  CONSTRAINT uq_users_username UNIQUE (username),
  CONSTRAINT uq_users_email UNIQUE (email),
  CONSTRAINT fk_users_company FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE SET NULL,
  CONSTRAINT fk_users_department FOREIGN KEY (department_id) REFERENCES departments(id) ON DELETE SET NULL
);

-- 角色表
CREATE TABLE IF NOT EXISTS roles (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  created_by INTEGER,
  updated_by INTEGER,
  is_active INTEGER DEFAULT 1,
  code TEXT NOT NULL,
  name TEXT NOT NULL,
  description TEXT,
  role_type TEXT DEFAULT 'CUSTOM',
  is_system INTEGER DEFAULT 0,
  sort_order INTEGER DEFAULT 0,
  status TEXT DEFAULT 'ACTIVE',
  CONSTRAINT uq_roles_code UNIQUE (code)
);

-- 权限表
CREATE TABLE IF NOT EXISTS permissions (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  created_by INTEGER,
  updated_by INTEGER,
  is_active INTEGER DEFAULT 1,
  code TEXT NOT NULL,
  name TEXT NOT NULL,
  description TEXT,
  resource TEXT NOT NULL,
  action TEXT NOT NULL,
  permission_type TEXT DEFAULT 'FUNCTION',
  parent_id INTEGER,
  is_system INTEGER DEFAULT 0,
  sort_order INTEGER DEFAULT 0,
  status TEXT DEFAULT 'ACTIVE',
  CONSTRAINT uq_permissions_code UNIQUE (code),
  CONSTRAINT fk_permissions_parent FOREIGN KEY (parent_id) REFERENCES permissions(id) ON DELETE SET NULL
);

-- 用户角色关联表
CREATE TABLE IF NOT EXISTS user_roles (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  user_id INTEGER NOT NULL,
  role_id INTEGER NOT NULL,
  assigned_by INTEGER,
  assigned_at DATETIME,
  expires_at DATETIME,
  is_active INTEGER DEFAULT 1,
  CONSTRAINT uq_user_roles_user_role UNIQUE (user_id, role_id),
  CONSTRAINT fk_user_roles_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  CONSTRAINT fk_user_roles_role FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
  CONSTRAINT fk_user_roles_assigned_by FOREIGN KEY (assigned_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 角色权限关联表
CREATE TABLE IF NOT EXISTS role_permissions (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  role_id INTEGER NOT NULL,
  permission_id INTEGER NOT NULL,
  assigned_by INTEGER,
  assigned_at DATETIME,
  is_active INTEGER DEFAULT 1,
  CONSTRAINT uq_role_permissions_role_permission UNIQUE (role_id, permission_id),
  CONSTRAINT fk_role_permissions_role FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
  CONSTRAINT fk_role_permissions_permission FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE,
  CONSTRAINT fk_role_permissions_assigned_by FOREIGN KEY (assigned_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 用户会话表
CREATE TABLE IF NOT EXISTS user_sessions (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  user_id INTEGER NOT NULL,
  session_id TEXT NOT NULL,
  ip_address TEXT,
  user_agent TEXT,
  device_info TEXT,
  location TEXT,
  login_at DATETIME NOT NULL,
  last_activity_at DATETIME,
  expires_at DATETIME,
  logout_at DATETIME,
  is_active INTEGER DEFAULT 1,
  session_data TEXT,
  CONSTRAINT uq_user_sessions_session_id UNIQUE (session_id),
  CONSTRAINT fk_user_sessions_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- ============================================================================
-- 2. 人力资源模块 (10个表)
-- ============================================================================

-- 职位表
CREATE TABLE IF NOT EXISTS positions (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  created_by INTEGER,
  updated_by INTEGER,
  is_active INTEGER DEFAULT 1,
  code TEXT NOT NULL,
  name TEXT NOT NULL,
  description TEXT,
  department_id INTEGER,
  level INTEGER DEFAULT 1,
  min_salary REAL DEFAULT 0,
  max_salary REAL DEFAULT 0,
  requirements TEXT,
  responsibilities TEXT,
  status TEXT DEFAULT 'ACTIVE',
  CONSTRAINT uq_positions_code UNIQUE (code),
  CONSTRAINT fk_positions_department FOREIGN KEY (department_id) REFERENCES departments(id) ON DELETE SET NULL
);

-- 员工表
CREATE TABLE IF NOT EXISTS employees (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  created_by INTEGER,
  updated_by INTEGER,
  is_active INTEGER DEFAULT 1,
  code TEXT NOT NULL,
  name TEXT NOT NULL,
  description TEXT,
  user_id INTEGER,
  employee_number TEXT NOT NULL,
  department_id INTEGER,
  position_id INTEGER,
  manager_id INTEGER,
  hire_date DATE,
  probation_end_date DATE,
  contract_start_date DATE,
  contract_end_date DATE,
  employment_type TEXT DEFAULT 'FULL_TIME',
  work_location TEXT,
  salary_grade_id INTEGER,
  base_salary REAL DEFAULT 0,
  status TEXT DEFAULT 'ACTIVE',
  CONSTRAINT uq_employees_code UNIQUE (code),
  CONSTRAINT uq_employees_employee_number UNIQUE (employee_number),
  CONSTRAINT fk_employees_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_employees_department FOREIGN KEY (department_id) REFERENCES departments(id) ON DELETE SET NULL,
  CONSTRAINT fk_employees_position FOREIGN KEY (position_id) REFERENCES positions(id) ON DELETE SET NULL,
  CONSTRAINT fk_employees_manager FOREIGN KEY (manager_id) REFERENCES employees(id) ON DELETE SET NULL
);

-- 薪资等级表
CREATE TABLE IF NOT EXISTS salary_grades (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  created_by INTEGER,
  updated_by INTEGER,
  is_active INTEGER DEFAULT 1,
  code TEXT NOT NULL,
  name TEXT NOT NULL,
  description TEXT,
  grade_level INTEGER NOT NULL,
  min_salary REAL NOT NULL,
  max_salary REAL NOT NULL,
  currency TEXT DEFAULT 'CNY',
  effective_date DATE,
  status TEXT DEFAULT 'ACTIVE',
  CONSTRAINT uq_salary_grades_code UNIQUE (code),
  CONSTRAINT uq_salary_grades_grade_level UNIQUE (grade_level)
);

-- 薪资记录表
CREATE TABLE IF NOT EXISTS salary_records (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  created_by INTEGER,
  updated_by INTEGER,
  is_active INTEGER DEFAULT 1,
  employee_id INTEGER NOT NULL,
  salary_month TEXT NOT NULL,
  base_salary REAL DEFAULT 0,
  allowances REAL DEFAULT 0,
  overtime_pay REAL DEFAULT 0,
  bonus REAL DEFAULT 0,
  deductions REAL DEFAULT 0,
  social_insurance REAL DEFAULT 0,
  housing_fund REAL DEFAULT 0,
  tax REAL DEFAULT 0,
  net_salary REAL DEFAULT 0,
  pay_date DATE,
  status TEXT DEFAULT 'PENDING',
  CONSTRAINT fk_salary_records_employee FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE CASCADE
);

-- 考勤表
CREATE TABLE IF NOT EXISTS attendances (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  created_by INTEGER,
  updated_by INTEGER,
  is_active INTEGER DEFAULT 1,
  employee_id INTEGER NOT NULL,
  attendance_date DATE NOT NULL,
  check_in_time DATETIME,
  check_out_time DATETIME,
  work_hours REAL DEFAULT 0,
  overtime_hours REAL DEFAULT 0,
  late_minutes INTEGER DEFAULT 0,
  early_leave_minutes INTEGER DEFAULT 0,
  attendance_type TEXT DEFAULT 'NORMAL',
  status TEXT DEFAULT 'NORMAL',
  remarks TEXT,
  CONSTRAINT fk_attendances_employee FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE CASCADE
);

-- 请假类型表
CREATE TABLE IF NOT EXISTS leave_types (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  created_by INTEGER,
  updated_by INTEGER,
  is_active INTEGER DEFAULT 1,
  code TEXT NOT NULL,
  name TEXT NOT NULL,
  description TEXT,
  max_days_per_year INTEGER DEFAULT 0,
  is_paid INTEGER DEFAULT 1,
  requires_approval INTEGER DEFAULT 1,
  approval_levels INTEGER DEFAULT 1,
  advance_notice_days INTEGER DEFAULT 1,
  status TEXT DEFAULT 'ACTIVE',
  CONSTRAINT uq_leave_types_code UNIQUE (code)
);

-- 请假申请表
CREATE TABLE IF NOT EXISTS leave_requests (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  created_by INTEGER,
  updated_by INTEGER,
  is_active INTEGER DEFAULT 1,
  employee_id INTEGER NOT NULL,
  leave_type_id INTEGER NOT NULL,
  start_date DATE NOT NULL,
  end_date DATE NOT NULL,
  days_requested REAL NOT NULL,
  reason TEXT,
  emergency_contact TEXT,
  emergency_phone TEXT,
  approved_by INTEGER,
  approved_at DATETIME,
  status TEXT DEFAULT 'PENDING',
  CONSTRAINT fk_leave_requests_employee FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE CASCADE,
  CONSTRAINT fk_leave_requests_leave_type FOREIGN KEY (leave_type_id) REFERENCES leave_types(id) ON DELETE CASCADE,
  CONSTRAINT fk_leave_requests_approved_by FOREIGN KEY (approved_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 绩效评估表
CREATE TABLE IF NOT EXISTS performance_reviews (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  created_by INTEGER,
  updated_by INTEGER,
  is_active INTEGER DEFAULT 1,
  employee_id INTEGER NOT NULL,
  reviewer_id INTEGER NOT NULL,
  review_period_start DATE NOT NULL,
  review_period_end DATE NOT NULL,
  overall_score REAL DEFAULT 0,
  goals_achievement REAL DEFAULT 0,
  competency_score REAL DEFAULT 0,
  behavior_score REAL DEFAULT 0,
  strengths TEXT,
  areas_for_improvement TEXT,
  development_plan TEXT,
  employee_comments TEXT,
  reviewer_comments TEXT,
  status TEXT DEFAULT 'DRAFT',
  CONSTRAINT fk_performance_reviews_employee FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE CASCADE,
  CONSTRAINT fk_performance_reviews_reviewer FOREIGN KEY (reviewer_id) REFERENCES users(id) ON DELETE CASCADE
);

-- 培训课程表
CREATE TABLE IF NOT EXISTS training_courses (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  created_by INTEGER,
  updated_by INTEGER,
  is_active INTEGER DEFAULT 1,
  code TEXT NOT NULL,
  name TEXT NOT NULL,
  description TEXT,
  category TEXT,
  instructor TEXT,
  duration_hours INTEGER DEFAULT 0,
  max_participants INTEGER DEFAULT 0,
  cost REAL DEFAULT 0,
  location TEXT,
  start_date DATE,
  end_date DATE,
  registration_deadline DATE,
  status TEXT DEFAULT 'PLANNED',
  CONSTRAINT uq_training_courses_code UNIQUE (code)
);

-- 培训记录表
CREATE TABLE IF NOT EXISTS training_records (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  created_by INTEGER,
  updated_by INTEGER,
  is_active INTEGER DEFAULT 1,
  employee_id INTEGER NOT NULL,
  training_course_id INTEGER NOT NULL,
  registration_date DATE,
  attendance_status TEXT DEFAULT 'REGISTERED',
  completion_date DATE,
  score REAL DEFAULT 0,
  certificate_number TEXT,
  feedback TEXT,
  status TEXT DEFAULT 'REGISTERED',
  CONSTRAINT fk_training_records_employee FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE CASCADE,
  CONSTRAINT fk_training_records_training_course FOREIGN KEY (training_course_id) REFERENCES training_courses(id) ON DELETE CASCADE
);

-- ============================================================================
-- 3. 销售管理模块 (12个表)
-- ============================================================================

-- 客户表
CREATE TABLE IF NOT EXISTS customers (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  created_by INTEGER,
  updated_by INTEGER,
  is_active INTEGER DEFAULT 1,
  code TEXT NOT NULL,
  name TEXT NOT NULL,
  description TEXT,
  customer_type TEXT DEFAULT 'ENTERPRISE',
  industry TEXT,
  company_size TEXT,
  tax_number TEXT,
  registration_number TEXT,
  website TEXT,
  address TEXT,
  city TEXT,
  province TEXT,
  country TEXT DEFAULT 'China',
  postal_code TEXT,
  phone TEXT,
  fax TEXT,
  email TEXT,
  credit_limit REAL DEFAULT 0,
  payment_terms INTEGER DEFAULT 30,
  currency TEXT DEFAULT 'CNY',
  sales_rep_id INTEGER,
  status TEXT DEFAULT 'ACTIVE',
  CONSTRAINT uq_customers_code UNIQUE (code),
  CONSTRAINT fk_customers_sales_rep FOREIGN KEY (sales_rep_id) REFERENCES users(id) ON DELETE SET NULL
);

-- 客户联系人表
CREATE TABLE IF NOT EXISTS customer_contacts (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  created_by INTEGER,
  updated_by INTEGER,
  is_active INTEGER DEFAULT 1,
  customer_id INTEGER NOT NULL,
  name TEXT NOT NULL,
  title TEXT,
  department TEXT,
  phone TEXT,
  mobile TEXT,
  email TEXT,
  address TEXT,
  is_primary INTEGER DEFAULT 0,
  is_decision_maker INTEGER DEFAULT 0,
  notes TEXT,
  status TEXT DEFAULT 'ACTIVE',
  CONSTRAINT fk_customer_contacts_customer FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE
);

-- 销售机会表
CREATE TABLE IF NOT EXISTS opportunities (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  created_by INTEGER,
  updated_by INTEGER,
  is_active INTEGER DEFAULT 1,
  code TEXT NOT NULL,
  name TEXT NOT NULL,
  description TEXT,
  customer_id INTEGER NOT NULL,
  contact_id INTEGER,
  sales_rep_id INTEGER,
  opportunity_source TEXT,
  stage TEXT DEFAULT 'PROSPECTING',
  probability REAL DEFAULT 0,
  estimated_value REAL DEFAULT 0,
  estimated_close_date DATE,
  actual_close_date DATE,
  next_action TEXT,
  next_action_date DATE,
  competitor_info TEXT,
  status TEXT DEFAULT 'OPEN',
  CONSTRAINT uq_opportunities_code UNIQUE (code),
  CONSTRAINT fk_opportunities_customer FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE,
  CONSTRAINT fk_opportunities_contact FOREIGN KEY (contact_id) REFERENCES customer_contacts(id) ON DELETE SET NULL,
  CONSTRAINT fk_opportunities_sales_rep FOREIGN KEY (sales_rep_id) REFERENCES users(id) ON DELETE SET NULL
);

-- 报价单表
CREATE TABLE IF NOT EXISTS quotations (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  created_by INTEGER,
  updated_by INTEGER,
  is_active INTEGER DEFAULT 1,
  code TEXT NOT NULL,
  name TEXT NOT NULL,
  description TEXT,
  customer_id INTEGER NOT NULL,
  contact_id INTEGER,
  opportunity_id INTEGER,
  sales_rep_id INTEGER,
  quotation_date DATE NOT NULL,
  valid_until DATE,
  currency TEXT DEFAULT 'CNY',
  exchange_rate REAL DEFAULT 1,
  subtotal REAL DEFAULT 0,
  discount_amount REAL DEFAULT 0,
  tax_amount REAL DEFAULT 0,
  total_amount REAL DEFAULT 0,
  payment_terms INTEGER DEFAULT 30,
  delivery_terms TEXT,
  notes TEXT,
  status TEXT DEFAULT 'DRAFT',
  CONSTRAINT uq_quotations_code UNIQUE (code),
  CONSTRAINT fk_quotations_customer FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE,
  CONSTRAINT fk_quotations_contact FOREIGN KEY (contact_id) REFERENCES customer_contacts(id) ON DELETE SET NULL,
  CONSTRAINT fk_quotations_opportunity FOREIGN KEY (opportunity_id) REFERENCES opportunities(id) ON DELETE SET NULL,
  CONSTRAINT fk_quotations_sales_rep FOREIGN KEY (sales_rep_id) REFERENCES users(id) ON DELETE SET NULL
);

-- 报价单明细表
CREATE TABLE IF NOT EXISTS quotation_items (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  quotation_id INTEGER NOT NULL,
  item_id INTEGER NOT NULL,
  quantity REAL NOT NULL,
  unit_price REAL NOT NULL,
  discount_percent REAL DEFAULT 0,
  discount_amount REAL DEFAULT 0,
  line_total REAL NOT NULL,
  delivery_date DATE,
  notes TEXT,
  CONSTRAINT fk_quotation_items_quotation FOREIGN KEY (quotation_id) REFERENCES quotations(id) ON DELETE CASCADE,
  CONSTRAINT fk_quotation_items_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE
);

-- 销售订单表
CREATE TABLE IF NOT EXISTS sales_orders (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  created_by INTEGER,
  updated_by INTEGER,
  is_active INTEGER DEFAULT 1,
  code TEXT NOT NULL,
  name TEXT NOT NULL,
  description TEXT,
  customer_id INTEGER NOT NULL,
  contact_id INTEGER,
  quotation_id INTEGER,
  sales_rep_id INTEGER,
  order_date DATE NOT NULL,
  required_date DATE,
  promised_date DATE,
  currency TEXT DEFAULT 'CNY',
  exchange_rate REAL DEFAULT 1,
  subtotal REAL DEFAULT 0,
  discount_amount REAL DEFAULT 0,
  tax_amount REAL DEFAULT 0,
  shipping_amount REAL DEFAULT 0,
  total_amount REAL DEFAULT 0,
  payment_terms INTEGER DEFAULT 30,
  shipping_address TEXT,
  billing_address TEXT,
  notes TEXT,
  status TEXT DEFAULT 'PENDING',
  CONSTRAINT uq_sales_orders_code UNIQUE (code),
  CONSTRAINT fk_sales_orders_customer FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE,
  CONSTRAINT fk_sales_orders_contact FOREIGN KEY (contact_id) REFERENCES customer_contacts(id) ON DELETE SET NULL,
  CONSTRAINT fk_sales_orders_quotation FOREIGN KEY (quotation_id) REFERENCES quotations(id) ON DELETE SET NULL,
  CONSTRAINT fk_sales_orders_sales_rep FOREIGN KEY (sales_rep_id) REFERENCES users(id) ON DELETE SET NULL
);

-- 销售订单明细表
CREATE TABLE IF NOT EXISTS sales_order_items (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  sales_order_id INTEGER NOT NULL,
  item_id INTEGER NOT NULL,
  quantity REAL NOT NULL,
  unit_price REAL NOT NULL,
  discount_percent REAL DEFAULT 0,
  discount_amount REAL DEFAULT 0,
  line_total REAL NOT NULL,
  shipped_quantity REAL DEFAULT 0,
  remaining_quantity REAL DEFAULT 0,
  delivery_date DATE,
  notes TEXT,
  CONSTRAINT fk_sales_order_items_sales_order FOREIGN KEY (sales_order_id) REFERENCES sales_orders(id) ON DELETE CASCADE,
  CONSTRAINT fk_sales_order_items_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE
);

-- 发货单表
CREATE TABLE IF NOT EXISTS shipments (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  created_by INTEGER,
  updated_by INTEGER,
  is_active INTEGER DEFAULT 1,
  code TEXT NOT NULL,
  name TEXT NOT NULL,
  description TEXT,
  sales_order_id INTEGER NOT NULL,
  customer_id INTEGER NOT NULL,
  warehouse_id INTEGER,
  shipment_date DATE NOT NULL,
  tracking_number TEXT,
  carrier TEXT,
  shipping_method TEXT,
  shipping_cost REAL DEFAULT 0,
  shipping_address TEXT,
  notes TEXT,
  status TEXT DEFAULT 'PENDING',
  CONSTRAINT uq_shipments_code UNIQUE (code),
  CONSTRAINT fk_shipments_sales_order FOREIGN KEY (sales_order_id) REFERENCES sales_orders(id) ON DELETE CASCADE,
  CONSTRAINT fk_shipments_customer FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE,
  CONSTRAINT fk_shipments_warehouse FOREIGN KEY (warehouse_id) REFERENCES warehouses(id) ON DELETE SET NULL
);

-- 发货单明细表
CREATE TABLE IF NOT EXISTS shipment_items (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  shipment_id INTEGER NOT NULL,
  sales_order_item_id INTEGER NOT NULL,
  item_id INTEGER NOT NULL,
  quantity REAL NOT NULL,
  unit_price REAL NOT NULL,
  line_total REAL NOT NULL,
  batch_number TEXT,
  serial_numbers TEXT,
  notes TEXT,
  CONSTRAINT fk_shipment_items_shipment FOREIGN KEY (shipment_id) REFERENCES shipments(id) ON DELETE CASCADE,
  CONSTRAINT fk_shipment_items_sales_order_item FOREIGN KEY (sales_order_item_id) REFERENCES sales_order_items(id) ON DELETE CASCADE,
  CONSTRAINT fk_shipment_items_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE
);

-- 销售发票表
CREATE TABLE IF NOT EXISTS sales_invoices (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  created_by INTEGER,
  updated_by INTEGER,
  is_active INTEGER DEFAULT 1,
  code TEXT NOT NULL,
  name TEXT NOT NULL,
  description TEXT,
  customer_id INTEGER NOT NULL,
  sales_order_id INTEGER,
  shipment_id INTEGER,
  invoice_date DATE NOT NULL,
  due_date DATE,
  currency TEXT DEFAULT 'CNY',
  exchange_rate REAL DEFAULT 1,
  subtotal REAL DEFAULT 0,
  discount_amount REAL DEFAULT 0,
  tax_amount REAL DEFAULT 0,
  total_amount REAL DEFAULT 0,
  paid_amount REAL DEFAULT 0,
  balance_amount REAL DEFAULT 0,
  payment_terms INTEGER DEFAULT 30,
  billing_address TEXT,
  notes TEXT,
  status TEXT DEFAULT 'PENDING',
  CONSTRAINT uq_sales_invoices_code UNIQUE (code),
  CONSTRAINT fk_sales_invoices_customer FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE,
  CONSTRAINT fk_sales_invoices_sales_order FOREIGN KEY (sales_order_id) REFERENCES sales_orders(id) ON DELETE SET NULL,
  CONSTRAINT fk_sales_invoices_shipment FOREIGN KEY (shipment_id) REFERENCES shipments(id) ON DELETE SET NULL
);

-- 销售发票明细表
CREATE TABLE IF NOT EXISTS sales_invoice_items (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  sales_invoice_id INTEGER NOT NULL,
  item_id INTEGER NOT NULL,
  quantity REAL NOT NULL,
  unit_price REAL NOT NULL,
  discount_percent REAL DEFAULT 0,
  discount_amount REAL DEFAULT 0,
  tax_rate REAL DEFAULT 0,
  tax_amount REAL DEFAULT 0,
  line_total REAL NOT NULL,
  notes TEXT,
  CONSTRAINT fk_sales_invoice_items_sales_invoice FOREIGN KEY (sales_invoice_id) REFERENCES sales_invoices(id) ON DELETE CASCADE,
  CONSTRAINT fk_sales_invoice_items_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE
);

-- ============================================================================
-- 4. 库存管理模块 (13个表)
-- ============================================================================

-- 物料分类表
CREATE TABLE IF NOT EXISTS item_categories (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  created_by INTEGER,
  updated_by INTEGER,
  is_active INTEGER DEFAULT 1,
  code TEXT NOT NULL,
  name TEXT NOT NULL,
  description TEXT,
  parent_id INTEGER,
  level INTEGER DEFAULT 1,
  sort_order INTEGER DEFAULT 0,
  status TEXT DEFAULT 'ACTIVE',
  CONSTRAINT uq_item_categories_code UNIQUE (code),
  CONSTRAINT fk_item_categories_parent FOREIGN KEY (parent_id) REFERENCES item_categories(id) ON DELETE SET NULL
);

-- 计量单位表
CREATE TABLE IF NOT EXISTS units (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  created_by INTEGER,
  updated_by INTEGER,
  is_active INTEGER DEFAULT 1,
  code TEXT NOT NULL,
  name TEXT NOT NULL,
  description TEXT,
  unit_type TEXT DEFAULT 'BASE',
  base_unit_id INTEGER,
  conversion_factor REAL DEFAULT 1,
  symbol TEXT,
  status TEXT DEFAULT 'ACTIVE',
  CONSTRAINT uq_units_code UNIQUE (code),
  CONSTRAINT fk_units_base_unit FOREIGN KEY (base_unit_id) REFERENCES units(id) ON DELETE SET NULL
);

-- 物料表
CREATE TABLE IF NOT EXISTS items (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  created_by INTEGER,
  updated_by INTEGER,
  is_active INTEGER DEFAULT 1,
  code TEXT NOT NULL,
  name TEXT NOT NULL,
  description TEXT,
  item_type TEXT DEFAULT 'PRODUCT',
  category_id INTEGER,
  unit_id INTEGER,
  purchase_unit_id INTEGER,
  sales_unit_id INTEGER,
  inventory_unit_id INTEGER,
  barcode TEXT,
  sku TEXT,
  brand TEXT,
  model TEXT,
  specification TEXT,
  weight REAL DEFAULT 0,
  volume REAL DEFAULT 0,
  length REAL DEFAULT 0,
  width REAL DEFAULT 0,
  height REAL DEFAULT 0,
  color TEXT,
  material TEXT,
  origin_country TEXT,
  hs_code TEXT,
  standard_cost REAL DEFAULT 0,
  list_price REAL DEFAULT 0,
  min_stock_level REAL DEFAULT 0,
  max_stock_level REAL DEFAULT 0,
  reorder_point REAL DEFAULT 0,
  reorder_quantity REAL DEFAULT 0,
  lead_time_days INTEGER DEFAULT 0,
  shelf_life_days INTEGER DEFAULT 0,
  is_serialized INTEGER DEFAULT 0,
  is_batch_tracked INTEGER DEFAULT 0,
  is_perishable INTEGER DEFAULT 0,
  is_hazardous INTEGER DEFAULT 0,
  image_url TEXT,
  notes TEXT,
  status TEXT DEFAULT 'ACTIVE',
  CONSTRAINT uq_items_code UNIQUE (code),
  CONSTRAINT fk_items_category FOREIGN KEY (category_id) REFERENCES item_categories(id) ON DELETE SET NULL,
  CONSTRAINT fk_items_unit FOREIGN KEY (unit_id) REFERENCES units(id) ON DELETE SET NULL,
  CONSTRAINT fk_items_purchase_unit FOREIGN KEY (purchase_unit_id) REFERENCES units(id) ON DELETE SET NULL,
  CONSTRAINT fk_items_sales_unit FOREIGN KEY (sales_unit_id) REFERENCES units(id) ON DELETE SET NULL,
  CONSTRAINT fk_items_inventory_unit FOREIGN KEY (inventory_unit_id) REFERENCES units(id) ON DELETE SET NULL
);

-- 仓库表
CREATE TABLE IF NOT EXISTS warehouses (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  created_by INTEGER,
  updated_by INTEGER,
  is_active INTEGER DEFAULT 1,
  code TEXT NOT NULL,
  name TEXT NOT NULL,
  description TEXT,
  warehouse_type TEXT DEFAULT 'GENERAL',
  address TEXT,
  city TEXT,
  province TEXT,
  country TEXT DEFAULT 'China',
  postal_code TEXT,
  phone TEXT,
  email TEXT,
  manager_id INTEGER,
  capacity REAL DEFAULT 0,
  area REAL DEFAULT 0,
  is_default INTEGER DEFAULT 0,
  status TEXT DEFAULT 'ACTIVE',
  CONSTRAINT uq_warehouses_code UNIQUE (code),
  CONSTRAINT fk_warehouses_manager FOREIGN KEY (manager_id) REFERENCES users(id) ON DELETE SET NULL
);

-- 库位表
CREATE TABLE IF NOT EXISTS locations (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  created_by INTEGER,
  updated_by INTEGER,
  is_active INTEGER DEFAULT 1,
  code TEXT NOT NULL,
  name TEXT NOT NULL,
  description TEXT,
  warehouse_id INTEGER NOT NULL,
  parent_id INTEGER,
  location_type TEXT DEFAULT 'STORAGE',
  zone TEXT,
  aisle TEXT,
  rack TEXT,
  shelf TEXT,
  bin TEXT,
  level INTEGER DEFAULT 1,
  capacity REAL DEFAULT 0,
  max_weight REAL DEFAULT 0,
  is_pickable INTEGER DEFAULT 1,
  is_receivable INTEGER DEFAULT 1,
  status TEXT DEFAULT 'ACTIVE',
  CONSTRAINT uq_locations_code UNIQUE (code),
  CONSTRAINT fk_locations_warehouse FOREIGN KEY (warehouse_id) REFERENCES warehouses(id) ON DELETE CASCADE,
  CONSTRAINT fk_locations_parent FOREIGN KEY (parent_id) REFERENCES locations(id) ON DELETE SET NULL
);

-- 库存表
CREATE TABLE IF NOT EXISTS inventories (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  item_id INTEGER NOT NULL,
  warehouse_id INTEGER NOT NULL,
  location_id INTEGER,
  quantity_on_hand REAL DEFAULT 0,
  quantity_available REAL DEFAULT 0,
  quantity_reserved REAL DEFAULT 0,
  quantity_on_order REAL DEFAULT 0,
  quantity_allocated REAL DEFAULT 0,
  average_cost REAL DEFAULT 0,
  last_cost REAL DEFAULT 0,
  last_movement_date DATETIME,
  last_count_date DATETIME,
  CONSTRAINT uq_inventories_item_warehouse_location UNIQUE (item_id, warehouse_id, location_id),
  CONSTRAINT fk_inventories_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE,
  CONSTRAINT fk_inventories_warehouse FOREIGN KEY (warehouse_id) REFERENCES warehouses(id) ON DELETE CASCADE,
  CONSTRAINT fk_inventories_location FOREIGN KEY (location_id) REFERENCES locations(id) ON DELETE SET NULL
);

-- 批次表
CREATE TABLE IF NOT EXISTS batches (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  created_by INTEGER,
  updated_by INTEGER,
  is_active INTEGER DEFAULT 1,
  batch_number TEXT NOT NULL,
  item_id INTEGER NOT NULL,
  warehouse_id INTEGER NOT NULL,
  location_id INTEGER,
  quantity REAL NOT NULL,
  production_date DATE,
  expiry_date DATE,
  supplier_id INTEGER,
  purchase_order_id INTEGER,
  cost REAL DEFAULT 0,
  notes TEXT,
  status TEXT DEFAULT 'ACTIVE',
  CONSTRAINT uq_batches_batch_number_item UNIQUE (batch_number, item_id),
  CONSTRAINT fk_batches_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE,
  CONSTRAINT fk_batches_warehouse FOREIGN KEY (warehouse_id) REFERENCES warehouses(id) ON DELETE CASCADE,
  CONSTRAINT fk_batches_location FOREIGN KEY (location_id) REFERENCES locations(id) ON DELETE SET NULL,
  CONSTRAINT fk_batches_supplier FOREIGN KEY (supplier_id) REFERENCES suppliers(id) ON DELETE SET NULL
);

-- 序列号表
CREATE TABLE IF NOT EXISTS serial_numbers (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  created_by INTEGER,
  updated_by INTEGER,
  is_active INTEGER DEFAULT 1,
  serial_number TEXT NOT NULL,
  item_id INTEGER NOT NULL,
  warehouse_id INTEGER NOT NULL,
  location_id INTEGER,
  batch_id INTEGER,
  status TEXT DEFAULT 'AVAILABLE',
  purchase_order_id INTEGER,
  sales_order_id INTEGER,
  customer_id INTEGER,
  warranty_start_date DATE,
  warranty_end_date DATE,
  notes TEXT,
  CONSTRAINT uq_serial_numbers_serial_number UNIQUE (serial_number),
  CONSTRAINT fk_serial_numbers_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE,
  CONSTRAINT fk_serial_numbers_warehouse FOREIGN KEY (warehouse_id) REFERENCES warehouses(id) ON DELETE CASCADE,
  CONSTRAINT fk_serial_numbers_location FOREIGN KEY (location_id) REFERENCES locations(id) ON DELETE SET NULL,
  CONSTRAINT fk_serial_numbers_batch FOREIGN KEY (batch_id) REFERENCES batches(id) ON DELETE SET NULL,
  CONSTRAINT fk_serial_numbers_customer FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE SET NULL
);

-- 库存移动表
CREATE TABLE IF NOT EXISTS inventory_movements (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  created_by INTEGER,
  updated_by INTEGER,
  movement_date DATETIME NOT NULL,
  movement_type TEXT NOT NULL,
  reference_type TEXT,
  reference_id INTEGER,
  reference_number TEXT,
  item_id INTEGER NOT NULL,
  from_warehouse_id INTEGER,
  from_location_id INTEGER,
  to_warehouse_id INTEGER,
  to_location_id INTEGER,
  quantity REAL NOT NULL,
  unit_cost REAL DEFAULT 0,
  total_cost REAL DEFAULT 0,
  batch_number TEXT,
  serial_number TEXT,
  reason TEXT,
  notes TEXT,
  CONSTRAINT fk_inventory_movements_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE,
  CONSTRAINT fk_inventory_movements_from_warehouse FOREIGN KEY (from_warehouse_id) REFERENCES warehouses(id) ON DELETE SET NULL,
  CONSTRAINT fk_inventory_movements_from_location FOREIGN KEY (from_location_id) REFERENCES locations(id) ON DELETE SET NULL,
  CONSTRAINT fk_inventory_movements_to_warehouse FOREIGN KEY (to_warehouse_id) REFERENCES warehouses(id) ON DELETE SET NULL,
  CONSTRAINT fk_inventory_movements_to_location FOREIGN KEY (to_location_id) REFERENCES locations(id) ON DELETE SET NULL
);

-- 库存调整表
CREATE TABLE IF NOT EXISTS inventory_adjustments (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  created_by INTEGER,
  updated_by INTEGER,
  is_active INTEGER DEFAULT 1,
  code TEXT NOT NULL,
  name TEXT NOT NULL,
  description TEXT,
  adjustment_date DATE NOT NULL,
  warehouse_id INTEGER NOT NULL,
  reason TEXT,
  approved_by INTEGER,
  approved_at DATETIME,
  status TEXT DEFAULT 'DRAFT',
  CONSTRAINT uq_inventory_adjustments_code UNIQUE (code),
  CONSTRAINT fk_inventory_adjustments_warehouse FOREIGN KEY (warehouse_id) REFERENCES warehouses(id) ON DELETE CASCADE,
  CONSTRAINT fk_inventory_adjustments_approved_by FOREIGN KEY (approved_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 库存调整明细表
CREATE TABLE IF NOT EXISTS inventory_adjustment_items (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  inventory_adjustment_id INTEGER NOT NULL,
  item_id INTEGER NOT NULL,
  location_id INTEGER,
  system_quantity REAL DEFAULT 0,
  actual_quantity REAL DEFAULT 0,
  adjustment_quantity REAL DEFAULT 0,
  unit_cost REAL DEFAULT 0,
  total_cost REAL DEFAULT 0,
  batch_number TEXT,
  serial_number TEXT,
  reason TEXT,
  notes TEXT,
  CONSTRAINT fk_inventory_adjustment_items_inventory_adjustment FOREIGN KEY (inventory_adjustment_id) REFERENCES inventory_adjustments(id) ON DELETE CASCADE,
  CONSTRAINT fk_inventory_adjustment_items_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE,
  CONSTRAINT fk_inventory_adjustment_items_location FOREIGN KEY (location_id) REFERENCES locations(id) ON DELETE SET NULL
);

-- 库存盘点表
CREATE TABLE IF NOT EXISTS stock_counts (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  created_by INTEGER,
  updated_by INTEGER,
  is_active INTEGER DEFAULT 1,
  code TEXT NOT NULL,
  name TEXT NOT NULL,
  description TEXT,
  count_date DATE NOT NULL,
  warehouse_id INTEGER NOT NULL,
  count_type TEXT DEFAULT 'FULL',
  scheduled_date DATE,
  started_at DATETIME,
  completed_at DATETIME,
  approved_by INTEGER,
  approved_at DATETIME,
  status TEXT DEFAULT 'PLANNED',
  CONSTRAINT uq_stock_counts_code UNIQUE (code),
  CONSTRAINT fk_stock_counts_warehouse FOREIGN KEY (warehouse_id) REFERENCES warehouses(id) ON DELETE CASCADE,
  CONSTRAINT fk_stock_counts_approved_by FOREIGN KEY (approved_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 库存盘点明细表
CREATE TABLE IF NOT EXISTS stock_count_items (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  stock_count_id INTEGER NOT NULL,
  item_id INTEGER NOT NULL,
  location_id INTEGER,
  system_quantity REAL DEFAULT 0,
  counted_quantity REAL DEFAULT 0,
  variance_quantity REAL DEFAULT 0,
  unit_cost REAL DEFAULT 0,
  variance_value REAL DEFAULT 0,
  batch_number TEXT,
  serial_number TEXT,
  counter_id INTEGER,
  count_date DATETIME,
  notes TEXT,
  CONSTRAINT fk_stock_count_items_stock_count FOREIGN KEY (stock_count_id) REFERENCES stock_counts(id) ON DELETE CASCADE,
  CONSTRAINT fk_stock_count_items_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE,
  CONSTRAINT fk_stock_count_items_location FOREIGN KEY (location_id) REFERENCES locations(id) ON DELETE SET NULL,
  CONSTRAINT fk_stock_count_items_counter FOREIGN KEY (counter_id) REFERENCES users(id) ON DELETE SET NULL
);

-- ============================================================================
-- 5. 采购管理模块 (15个表)
-- ============================================================================

-- 供应商表
CREATE TABLE IF NOT EXISTS suppliers (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  created_by INTEGER,
  updated_by INTEGER,
  is_active INTEGER DEFAULT 1,
  code TEXT NOT NULL,
  name TEXT NOT NULL,
  description TEXT,
  supplier_type TEXT DEFAULT 'VENDOR',
  industry TEXT,
  company_size TEXT,
  tax_number TEXT,
  registration_number TEXT,
  website TEXT,
  address TEXT,
  city TEXT,
  province TEXT,
  country TEXT DEFAULT 'China',
  postal_code TEXT,
  phone TEXT,
  fax TEXT,
  email TEXT,
  credit_rating TEXT,
  payment_terms INTEGER DEFAULT 30,
  currency TEXT DEFAULT 'CNY',
  buyer_id INTEGER,
  status TEXT DEFAULT 'ACTIVE',
  CONSTRAINT uq_suppliers_code UNIQUE (code),
  CONSTRAINT fk_suppliers_buyer FOREIGN KEY (buyer_id) REFERENCES users(id) ON DELETE SET NULL
);

-- 供应商联系人表
CREATE TABLE IF NOT EXISTS supplier_contacts (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  created_by INTEGER,
  updated_by INTEGER,
  is_active INTEGER DEFAULT 1,
  supplier_id INTEGER NOT NULL,
  name TEXT NOT NULL,
  title TEXT,
  department TEXT,
  phone TEXT,
  mobile TEXT,
  email TEXT,
  address TEXT,
  is_primary INTEGER DEFAULT 0,
  is_decision_maker INTEGER DEFAULT 0,
  notes TEXT,
  status TEXT DEFAULT 'ACTIVE',
  CONSTRAINT fk_supplier_contacts_supplier FOREIGN KEY (supplier_id) REFERENCES suppliers(id) ON DELETE CASCADE
);

-- 采购申请表
CREATE TABLE IF NOT EXISTS purchase_requests (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  created_by INTEGER,
  updated_by INTEGER,
  is_active INTEGER DEFAULT 1,
  code TEXT NOT NULL,
  name TEXT NOT NULL,
  description TEXT,
  request_date DATE NOT NULL,
  required_date DATE,
  department_id INTEGER,
  requester_id INTEGER,
  priority TEXT DEFAULT 'NORMAL',
  total_amount REAL DEFAULT 0,
  approved_by INTEGER,
  approved_at DATETIME,
  notes TEXT,
  status TEXT DEFAULT 'DRAFT',
  CONSTRAINT uq_purchase_requests_code UNIQUE (code),
  CONSTRAINT fk_purchase_requests_department FOREIGN KEY (department_id) REFERENCES departments(id) ON DELETE SET NULL,
  CONSTRAINT fk_purchase_requests_requester FOREIGN KEY (requester_id) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_purchase_requests_approved_by FOREIGN KEY (approved_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 采购申请明细表
CREATE TABLE IF NOT EXISTS purchase_request_items (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  purchase_request_id INTEGER NOT NULL,
  item_id INTEGER NOT NULL,
  quantity REAL NOT NULL,
  estimated_price REAL DEFAULT 0,
  line_total REAL DEFAULT 0,
  required_date DATE,
  specification TEXT,
  notes TEXT,
  CONSTRAINT fk_purchase_request_items_purchase_request FOREIGN KEY (purchase_request_id) REFERENCES purchase_requests(id) ON DELETE CASCADE,
  CONSTRAINT fk_purchase_request_items_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE
);

-- 询价单表
CREATE TABLE IF NOT EXISTS rfqs (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  created_by INTEGER,
  updated_by INTEGER,
  is_active INTEGER DEFAULT 1,
  code TEXT NOT NULL,
  name TEXT NOT NULL,
  description TEXT,
  rfq_date DATE NOT NULL,
  response_deadline DATE,
  buyer_id INTEGER,
  currency TEXT DEFAULT 'CNY',
  payment_terms INTEGER DEFAULT 30,
  delivery_terms TEXT,
  notes TEXT,
  status TEXT DEFAULT 'DRAFT',
  CONSTRAINT uq_rfqs_code UNIQUE (code),
  CONSTRAINT fk_rfqs_buyer FOREIGN KEY (buyer_id) REFERENCES users(id) ON DELETE SET NULL
);

-- 询价单明细表
CREATE TABLE IF NOT EXISTS rfq_items (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  rfq_id INTEGER NOT NULL,
  item_id INTEGER NOT NULL,
  quantity REAL NOT NULL,
  required_date DATE,
  specification TEXT,
  notes TEXT,
  CONSTRAINT fk_rfq_items_rfq FOREIGN KEY (rfq_id) REFERENCES rfqs(id) ON DELETE CASCADE,
  CONSTRAINT fk_rfq_items_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE
);

-- 询价供应商表
CREATE TABLE IF NOT EXISTS rfq_suppliers (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  rfq_id INTEGER NOT NULL,
  supplier_id INTEGER NOT NULL,
  contact_id INTEGER,
  sent_date DATE,
  response_date DATE,
  status TEXT DEFAULT 'SENT',
  CONSTRAINT fk_rfq_suppliers_rfq FOREIGN KEY (rfq_id) REFERENCES rfqs(id) ON DELETE CASCADE,
  CONSTRAINT fk_rfq_suppliers_supplier FOREIGN KEY (supplier_id) REFERENCES suppliers(id) ON DELETE CASCADE,
  CONSTRAINT fk_rfq_suppliers_contact FOREIGN KEY (contact_id) REFERENCES supplier_contacts(id) ON DELETE SET NULL
);

-- 供应商报价表
CREATE TABLE IF NOT EXISTS supplier_quotes (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  created_by INTEGER,
  updated_by INTEGER,
  is_active INTEGER DEFAULT 1,
  code TEXT NOT NULL,
  name TEXT NOT NULL,
  description TEXT,
  rfq_id INTEGER,
  supplier_id INTEGER NOT NULL,
  contact_id INTEGER,
  quote_date DATE NOT NULL,
  valid_until DATE,
  currency TEXT DEFAULT 'CNY',
  exchange_rate REAL DEFAULT 1,
  payment_terms INTEGER DEFAULT 30,
  delivery_terms TEXT,
  total_amount REAL DEFAULT 0,
  notes TEXT,
  status TEXT DEFAULT 'RECEIVED',
  CONSTRAINT uq_supplier_quotes_code UNIQUE (code),
  CONSTRAINT fk_supplier_quotes_rfq FOREIGN KEY (rfq_id) REFERENCES rfqs(id) ON DELETE SET NULL,
  CONSTRAINT fk_supplier_quotes_supplier FOREIGN KEY (supplier_id) REFERENCES suppliers(id) ON DELETE CASCADE,
  CONSTRAINT fk_supplier_quotes_contact FOREIGN KEY (contact_id) REFERENCES supplier_contacts(id) ON DELETE SET NULL
);

-- 供应商报价明细表
CREATE TABLE IF NOT EXISTS supplier_quote_items (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  supplier_quote_id INTEGER NOT NULL,
  rfq_item_id INTEGER,
  item_id INTEGER NOT NULL,
  quantity REAL NOT NULL,
  unit_price REAL NOT NULL,
  line_total REAL NOT NULL,
  delivery_date DATE,
  notes TEXT,
  CONSTRAINT fk_supplier_quote_items_supplier_quote FOREIGN KEY (supplier_quote_id) REFERENCES supplier_quotes(id) ON DELETE CASCADE,
  CONSTRAINT fk_supplier_quote_items_rfq_item FOREIGN KEY (rfq_item_id) REFERENCES rfq_items(id) ON DELETE SET NULL,
  CONSTRAINT fk_supplier_quote_items_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE
);

-- 采购订单表
CREATE TABLE IF NOT EXISTS purchase_orders (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  created_by INTEGER,
  updated_by INTEGER,
  is_active INTEGER DEFAULT 1,
  code TEXT NOT NULL,
  name TEXT NOT NULL,
  description TEXT,
  supplier_id INTEGER NOT NULL,
  contact_id INTEGER,
  purchase_request_id INTEGER,
  supplier_quote_id INTEGER,
  buyer_id INTEGER,
  order_date DATE NOT NULL,
  required_date DATE,
  promised_date DATE,
  currency TEXT DEFAULT 'CNY',
  exchange_rate REAL DEFAULT 1,
  subtotal REAL DEFAULT 0,
  discount_amount REAL DEFAULT 0,
  tax_amount REAL DEFAULT 0,
  shipping_amount REAL DEFAULT 0,
  total_amount REAL DEFAULT 0,
  payment_terms INTEGER DEFAULT 30,
  delivery_address TEXT,
  notes TEXT,
  status TEXT DEFAULT 'PENDING',
  CONSTRAINT uq_purchase_orders_code UNIQUE (code),
  CONSTRAINT fk_purchase_orders_supplier FOREIGN KEY (supplier_id) REFERENCES suppliers(id) ON DELETE CASCADE,
  CONSTRAINT fk_purchase_orders_contact FOREIGN KEY (contact_id) REFERENCES supplier_contacts(id) ON DELETE SET NULL,
  CONSTRAINT fk_purchase_orders_purchase_request FOREIGN KEY (purchase_request_id) REFERENCES purchase_requests(id) ON DELETE SET NULL,
  CONSTRAINT fk_purchase_orders_supplier_quote FOREIGN KEY (supplier_quote_id) REFERENCES supplier_quotes(id) ON DELETE SET NULL,
  CONSTRAINT fk_purchase_orders_buyer FOREIGN KEY (buyer_id) REFERENCES users(id) ON DELETE SET NULL
);

-- 采购订单明细表
CREATE TABLE IF NOT EXISTS purchase_order_items (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  purchase_order_id INTEGER NOT NULL,
  item_id INTEGER NOT NULL,
  quantity REAL NOT NULL,
  unit_price REAL NOT NULL,
  discount_percent REAL DEFAULT 0,
  discount_amount REAL DEFAULT 0,
  line_total REAL NOT NULL,
  received_quantity REAL DEFAULT 0,
  remaining_quantity REAL DEFAULT 0,
  delivery_date DATE,
  notes TEXT,
  CONSTRAINT fk_purchase_order_items_purchase_order FOREIGN KEY (purchase_order_id) REFERENCES purchase_orders(id) ON DELETE CASCADE,
  CONSTRAINT fk_purchase_order_items_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE
);

-- 收货单表
CREATE TABLE IF NOT EXISTS receipts (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  created_by INTEGER,
  updated_by INTEGER,
  is_active INTEGER DEFAULT 1,
  code TEXT NOT NULL,
  name TEXT NOT NULL,
  description TEXT,
  purchase_order_id INTEGER NOT NULL,
  supplier_id INTEGER NOT NULL,
  warehouse_id INTEGER,
  receipt_date DATE NOT NULL,
  delivery_note_number TEXT,
  carrier TEXT,
  tracking_number TEXT,
  received_by INTEGER,
  inspected_by INTEGER,
  inspection_date DATE,
  notes TEXT,
  status TEXT DEFAULT 'PENDING',
  CONSTRAINT uq_receipts_code UNIQUE (code),
  CONSTRAINT fk_receipts_purchase_order FOREIGN KEY (purchase_order_id) REFERENCES purchase_orders(id) ON DELETE CASCADE,
  CONSTRAINT fk_receipts_supplier FOREIGN KEY (supplier_id) REFERENCES suppliers(id) ON DELETE CASCADE,
  CONSTRAINT fk_receipts_warehouse FOREIGN KEY (warehouse_id) REFERENCES warehouses(id) ON DELETE SET NULL,
  CONSTRAINT fk_receipts_received_by FOREIGN KEY (received_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_receipts_inspected_by FOREIGN KEY (inspected_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 收货单明细表
CREATE TABLE IF NOT EXISTS receipt_items (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  receipt_id INTEGER NOT NULL,
  purchase_order_item_id INTEGER NOT NULL,
  item_id INTEGER NOT NULL,
  ordered_quantity REAL NOT NULL,
  received_quantity REAL NOT NULL,
  accepted_quantity REAL NOT NULL,
  rejected_quantity REAL DEFAULT 0,
  unit_cost REAL NOT NULL,
  line_total REAL NOT NULL,
  batch_number TEXT,
  expiry_date DATE,
  location_id INTEGER,
  quality_status TEXT DEFAULT 'PASSED',
  notes TEXT,
  CONSTRAINT fk_receipt_items_receipt FOREIGN KEY (receipt_id) REFERENCES receipts(id) ON DELETE CASCADE,
  CONSTRAINT fk_receipt_items_purchase_order_item FOREIGN KEY (purchase_order_item_id) REFERENCES purchase_order_items(id) ON DELETE CASCADE,
  CONSTRAINT fk_receipt_items_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE,
  CONSTRAINT fk_receipt_items_location FOREIGN KEY (location_id) REFERENCES locations(id) ON DELETE SET NULL
);

-- 采购发票表
CREATE TABLE IF NOT EXISTS purchase_invoices (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  created_by INTEGER,
  updated_by INTEGER,
  is_active INTEGER DEFAULT 1,
  code TEXT NOT NULL,
  name TEXT NOT NULL,
  description TEXT,
  supplier_id INTEGER NOT NULL,
  purchase_order_id INTEGER,
  receipt_id INTEGER,
  invoice_number TEXT NOT NULL,
  invoice_date DATE NOT NULL,
  due_date DATE,
  currency TEXT DEFAULT 'CNY',
  exchange_rate REAL DEFAULT 1,
  subtotal REAL DEFAULT 0,
  discount_amount REAL DEFAULT 0,
  tax_amount REAL DEFAULT 0,
  total_amount REAL DEFAULT 0,
  paid_amount REAL DEFAULT 0,
  balance_amount REAL DEFAULT 0,
  payment_terms INTEGER DEFAULT 30,
  notes TEXT,
  status TEXT DEFAULT 'PENDING',
  CONSTRAINT uq_purchase_invoices_code UNIQUE (code),
  CONSTRAINT fk_purchase_invoices_supplier FOREIGN KEY (supplier_id) REFERENCES suppliers(id) ON DELETE CASCADE,
  CONSTRAINT fk_purchase_invoices_purchase_order FOREIGN KEY (purchase_order_id) REFERENCES purchase_orders(id) ON DELETE SET NULL,
  CONSTRAINT fk_purchase_invoices_receipt FOREIGN KEY (receipt_id) REFERENCES receipts(id) ON DELETE SET NULL
);

-- 采购发票明细表
CREATE TABLE IF NOT EXISTS purchase_invoice_items (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  purchase_invoice_id INTEGER NOT NULL,
  item_id INTEGER NOT NULL,
  quantity REAL NOT NULL,
  unit_price REAL NOT NULL,
  discount_percent REAL DEFAULT 0,
  discount_amount REAL DEFAULT 0,
  tax_rate REAL DEFAULT 0,
  tax_amount REAL DEFAULT 0,
  line_total REAL NOT NULL,
  notes TEXT,
  CONSTRAINT fk_purchase_invoice_items_purchase_invoice FOREIGN KEY (purchase_invoice_id) REFERENCES purchase_invoices(id) ON DELETE CASCADE,
  CONSTRAINT fk_purchase_invoice_items_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE
);

-- ============================================================================
-- 6. 财务会计模块 (10个表)
-- ============================================================================

-- 会计科目表
CREATE TABLE IF NOT EXISTS accounts (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  created_by INTEGER,
  updated_by INTEGER,
  is_active INTEGER DEFAULT 1,
  code TEXT NOT NULL,
  name TEXT NOT NULL,
  description TEXT,
  account_type TEXT NOT NULL,
  account_category TEXT,
  parent_id INTEGER,
  level INTEGER DEFAULT 1,
  is_leaf INTEGER DEFAULT 1,
  normal_balance TEXT NOT NULL,
  is_system INTEGER DEFAULT 0,
  is_control INTEGER DEFAULT 0,
  currency TEXT DEFAULT 'CNY',
  sort_order INTEGER DEFAULT 0,
  status TEXT DEFAULT 'ACTIVE',
  CONSTRAINT uq_accounts_code UNIQUE (code),
  CONSTRAINT fk_accounts_parent FOREIGN KEY (parent_id) REFERENCES accounts(id) ON DELETE SET NULL
);

-- 会计期间表
CREATE TABLE IF NOT EXISTS fiscal_periods (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  created_by INTEGER,
  updated_by INTEGER,
  is_active INTEGER DEFAULT 1,
  code TEXT NOT NULL,
  name TEXT NOT NULL,
  description TEXT,
  fiscal_year INTEGER NOT NULL,
  period_number INTEGER NOT NULL,
  start_date DATE NOT NULL,
  end_date DATE NOT NULL,
  is_closed INTEGER DEFAULT 0,
  closed_by INTEGER,
  closed_at DATETIME,
  status TEXT DEFAULT 'OPEN',
  CONSTRAINT uq_fiscal_periods_code UNIQUE (code),
  CONSTRAINT uq_fiscal_periods_year_period UNIQUE (fiscal_year, period_number),
  CONSTRAINT fk_fiscal_periods_closed_by FOREIGN KEY (closed_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 会计凭证表
CREATE TABLE IF NOT EXISTS journal_entries (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  created_by INTEGER,
  updated_by INTEGER,
  is_active INTEGER DEFAULT 1,
  code TEXT NOT NULL,
  name TEXT NOT NULL,
  description TEXT,
  entry_date DATE NOT NULL,
  fiscal_period_id INTEGER NOT NULL,
  reference_type TEXT,
  reference_id INTEGER,
  reference_number TEXT,
  total_debit REAL DEFAULT 0,
  total_credit REAL DEFAULT 0,
  currency TEXT DEFAULT 'CNY',
  exchange_rate REAL DEFAULT 1,
  posted_by INTEGER,
  posted_at DATETIME,
  reversed_by INTEGER,
  reversed_at DATETIME,
  reversal_entry_id INTEGER,
  status TEXT DEFAULT 'DRAFT',
  CONSTRAINT uq_journal_entries_code UNIQUE (code),
  CONSTRAINT fk_journal_entries_fiscal_period FOREIGN KEY (fiscal_period_id) REFERENCES fiscal_periods(id) ON DELETE CASCADE,
  CONSTRAINT fk_journal_entries_posted_by FOREIGN KEY (posted_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_journal_entries_reversed_by FOREIGN KEY (reversed_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_journal_entries_reversal_entry FOREIGN KEY (reversal_entry_id) REFERENCES journal_entries(id) ON DELETE SET NULL
);

-- 会计凭证明细表
CREATE TABLE IF NOT EXISTS journal_entry_lines (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  journal_entry_id INTEGER NOT NULL,
  line_number INTEGER NOT NULL,
  account_id INTEGER NOT NULL,
  debit_amount REAL DEFAULT 0,
  credit_amount REAL DEFAULT 0,
  description TEXT,
  reference_type TEXT,
  reference_id INTEGER,
  cost_center TEXT,
  project_id INTEGER,
  CONSTRAINT fk_journal_entry_lines_journal_entry FOREIGN KEY (journal_entry_id) REFERENCES journal_entries(id) ON DELETE CASCADE,
  CONSTRAINT fk_journal_entry_lines_account FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE,
  CONSTRAINT fk_journal_entry_lines_project FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE SET NULL
);

-- 总账余额表
CREATE TABLE IF NOT EXISTS general_ledger_balances (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  account_id INTEGER NOT NULL,
  fiscal_period_id INTEGER NOT NULL,
  beginning_balance REAL DEFAULT 0,
  debit_amount REAL DEFAULT 0,
  credit_amount REAL DEFAULT 0,
  ending_balance REAL DEFAULT 0,
  CONSTRAINT uq_general_ledger_balances_account_period UNIQUE (account_id, fiscal_period_id),
  CONSTRAINT fk_general_ledger_balances_account FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE,
  CONSTRAINT fk_general_ledger_balances_fiscal_period FOREIGN KEY (fiscal_period_id) REFERENCES fiscal_periods(id) ON DELETE CASCADE
);

-- 应收账款表
CREATE TABLE IF NOT EXISTS accounts_receivable (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  created_by INTEGER,
  updated_by INTEGER,
  is_active INTEGER DEFAULT 1,
  code TEXT NOT NULL,
  customer_id INTEGER NOT NULL,
  sales_invoice_id INTEGER,
  invoice_number TEXT,
  invoice_date DATE NOT NULL,
  due_date DATE,
  original_amount REAL NOT NULL,
  paid_amount REAL DEFAULT 0,
  balance_amount REAL NOT NULL,
  currency TEXT DEFAULT 'CNY',
  exchange_rate REAL DEFAULT 1,
  aging_days INTEGER DEFAULT 0,
  status TEXT DEFAULT 'OPEN',
  CONSTRAINT uq_accounts_receivable_code UNIQUE (code),
  CONSTRAINT fk_accounts_receivable_customer FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE,
  CONSTRAINT fk_accounts_receivable_sales_invoice FOREIGN KEY (sales_invoice_id) REFERENCES sales_invoices(id) ON DELETE SET NULL
);

-- 应付账款表
CREATE TABLE IF NOT EXISTS accounts_payable (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  created_by INTEGER,
  updated_by INTEGER,
  is_active INTEGER DEFAULT 1,
  code TEXT NOT NULL,
  supplier_id INTEGER NOT NULL,
  purchase_invoice_id INTEGER,
  invoice_number TEXT,
  invoice_date DATE NOT NULL,
  due_date DATE,
  original_amount REAL NOT NULL,
  paid_amount REAL DEFAULT 0,
  balance_amount REAL NOT NULL,
  currency TEXT DEFAULT 'CNY',
  exchange_rate REAL DEFAULT 1,
  aging_days INTEGER DEFAULT 0,
  status TEXT DEFAULT 'OPEN',
  CONSTRAINT uq_accounts_payable_code UNIQUE (code),
  CONSTRAINT fk_accounts_payable_supplier FOREIGN KEY (supplier_id) REFERENCES suppliers(id) ON DELETE CASCADE,
  CONSTRAINT fk_accounts_payable_purchase_invoice FOREIGN KEY (purchase_invoice_id) REFERENCES purchase_invoices(id) ON DELETE SET NULL
);

-- 付款单表
CREATE TABLE IF NOT EXISTS payments (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  created_by INTEGER,
  updated_by INTEGER,
  is_active INTEGER DEFAULT 1,
  code TEXT NOT NULL,
  name TEXT NOT NULL,
  description TEXT,
  supplier_id INTEGER NOT NULL,
  payment_date DATE NOT NULL,
  payment_method TEXT DEFAULT 'BANK_TRANSFER',
  bank_account_id INTEGER,
  reference_number TEXT,
  total_amount REAL NOT NULL,
  currency TEXT DEFAULT 'CNY',
  exchange_rate REAL DEFAULT 1,
  notes TEXT,
  status TEXT DEFAULT 'PENDING',
  CONSTRAINT uq_payments_code UNIQUE (code),
  CONSTRAINT fk_payments_supplier FOREIGN KEY (supplier_id) REFERENCES suppliers(id) ON DELETE CASCADE,
  CONSTRAINT fk_payments_bank_account FOREIGN KEY (bank_account_id) REFERENCES bank_accounts(id) ON DELETE SET NULL
);

-- 付款单明细表
CREATE TABLE IF NOT EXISTS payment_items (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  payment_id INTEGER NOT NULL,
  accounts_payable_id INTEGER NOT NULL,
  invoice_number TEXT,
  original_amount REAL NOT NULL,
  payment_amount REAL NOT NULL,
  discount_amount REAL DEFAULT 0,
  notes TEXT,
  CONSTRAINT fk_payment_items_payment FOREIGN KEY (payment_id) REFERENCES payments(id) ON DELETE CASCADE,
  CONSTRAINT fk_payment_items_accounts_payable FOREIGN KEY (accounts_payable_id) REFERENCES accounts_payable(id) ON DELETE CASCADE
);

-- 收款单表
CREATE TABLE IF NOT EXISTS receipts_ar (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  created_by INTEGER,
  updated_by INTEGER,
  is_active INTEGER DEFAULT 1,
  code TEXT NOT NULL,
  name TEXT NOT NULL,
  description TEXT,
  customer_id INTEGER NOT NULL,
  receipt_date DATE NOT NULL,
  payment_method TEXT DEFAULT 'BANK_TRANSFER',
  bank_account_id INTEGER,
  reference_number TEXT,
  total_amount REAL NOT NULL,
  currency TEXT DEFAULT 'CNY',
  exchange_rate REAL DEFAULT 1,
  notes TEXT,
  status TEXT DEFAULT 'PENDING',
  CONSTRAINT uq_receipts_ar_code UNIQUE (code),
  CONSTRAINT fk_receipts_ar_customer FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE,
  CONSTRAINT fk_receipts_ar_bank_account FOREIGN KEY (bank_account_id) REFERENCES bank_accounts(id) ON DELETE SET NULL
);

-- 收款单明细表
CREATE TABLE IF NOT EXISTS receipt_ar_items (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  receipt_ar_id INTEGER NOT NULL,
  accounts_receivable_id INTEGER NOT NULL,
  invoice_number TEXT,
  original_amount REAL NOT NULL,
  receipt_amount REAL NOT NULL,
  discount_amount REAL DEFAULT 0,
  notes TEXT,
  CONSTRAINT fk_receipt_ar_items_receipt_ar FOREIGN KEY (receipt_ar_id) REFERENCES receipts_ar(id) ON DELETE CASCADE,
  CONSTRAINT fk_receipt_ar_items_accounts_receivable FOREIGN KEY (accounts_receivable_id) REFERENCES accounts_receivable(id) ON DELETE CASCADE
);

-- 银行账户表
CREATE TABLE IF NOT EXISTS bank_accounts (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  created_by INTEGER,
  updated_by INTEGER,
  is_active INTEGER DEFAULT 1,
  code TEXT NOT NULL,
  name TEXT NOT NULL,
  description TEXT,
  bank_name TEXT NOT NULL,
  account_number TEXT NOT NULL,
  account_type TEXT DEFAULT 'CHECKING',
  currency TEXT DEFAULT 'CNY',
  opening_balance REAL DEFAULT 0,
  current_balance REAL DEFAULT 0,
  is_default INTEGER DEFAULT 0,
  status TEXT DEFAULT 'ACTIVE',
  CONSTRAINT uq_bank_accounts_code UNIQUE (code),
  CONSTRAINT uq_bank_accounts_account_number UNIQUE (account_number)
);

-- ============================================================================
-- 7. 生产管理模块 (12个表)
-- ============================================================================

-- 物料清单表
CREATE TABLE IF NOT EXISTS boms (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  created_by INTEGER,
  updated_by INTEGER,
  is_active INTEGER DEFAULT 1,
  code TEXT NOT NULL,
  name TEXT NOT NULL,
  description TEXT,
  item_id INTEGER NOT NULL,
  version TEXT DEFAULT '1.0',
  effective_date DATE,
  expiry_date DATE,
  bom_type TEXT DEFAULT 'PRODUCTION',
  base_quantity REAL DEFAULT 1,
  unit_id INTEGER,
  approved_by INTEGER,
  approved_at DATETIME,
  status TEXT DEFAULT 'DRAFT',
  CONSTRAINT uq_boms_code UNIQUE (code),
  CONSTRAINT fk_boms_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE,
  CONSTRAINT fk_boms_unit FOREIGN KEY (unit_id) REFERENCES units(id) ON DELETE SET NULL,
  CONSTRAINT fk_boms_approved_by FOREIGN KEY (approved_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 物料清单明细表
CREATE TABLE IF NOT EXISTS bom_items (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  bom_id INTEGER NOT NULL,
  item_id INTEGER NOT NULL,
  quantity REAL NOT NULL,
  unit_id INTEGER,
  scrap_factor REAL DEFAULT 0,
  line_number INTEGER DEFAULT 0,
  notes TEXT,
  CONSTRAINT fk_bom_items_bom FOREIGN KEY (bom_id) REFERENCES boms(id) ON DELETE CASCADE,
  CONSTRAINT fk_bom_items_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE,
  CONSTRAINT fk_bom_items_unit FOREIGN KEY (unit_id) REFERENCES units(id) ON DELETE SET NULL
);

-- 工作中心表
CREATE TABLE IF NOT EXISTS work_centers (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  created_by INTEGER,
  updated_by INTEGER,
  is_active INTEGER DEFAULT 1,
  code TEXT NOT NULL,
  name TEXT NOT NULL,
  description TEXT,
  work_center_type TEXT DEFAULT 'MACHINE',
  capacity_per_hour REAL DEFAULT 0,
  efficiency_factor REAL DEFAULT 1,
  cost_per_hour REAL DEFAULT 0,
  setup_time_minutes INTEGER DEFAULT 0,
  queue_time_minutes INTEGER DEFAULT 0,
  location TEXT,
  responsible_person TEXT,
  status TEXT DEFAULT 'ACTIVE',
  CONSTRAINT uq_work_centers_code UNIQUE (code)
);

-- 工艺路线表
CREATE TABLE IF NOT EXISTS routings (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  created_by INTEGER,
  updated_by INTEGER,
  is_active INTEGER DEFAULT 1,
  code TEXT NOT NULL,
  name TEXT NOT NULL,
  description TEXT,
  item_id INTEGER NOT NULL,
  version TEXT DEFAULT '1.0',
  effective_date DATE,
  expiry_date DATE,
  base_quantity REAL DEFAULT 1,
  unit_id INTEGER,
  approved_by INTEGER,
  approved_at DATETIME,
  status TEXT DEFAULT 'DRAFT',
  CONSTRAINT uq_routings_code UNIQUE (code),
  CONSTRAINT fk_routings_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE,
  CONSTRAINT fk_routings_unit FOREIGN KEY (unit_id) REFERENCES units(id) ON DELETE SET NULL,
  CONSTRAINT fk_routings_approved_by FOREIGN KEY (approved_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 工艺路线操作表
CREATE TABLE IF NOT EXISTS routing_operations (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  routing_id INTEGER NOT NULL,
  operation_number INTEGER NOT NULL,
  work_center_id INTEGER NOT NULL,
  operation_name TEXT NOT NULL,
  description TEXT,
  setup_time_minutes INTEGER DEFAULT 0,
  run_time_minutes INTEGER DEFAULT 0,
  queue_time_minutes INTEGER DEFAULT 0,
  move_time_minutes INTEGER DEFAULT 0,
  overlap_quantity REAL DEFAULT 0,
  notes TEXT,
  CONSTRAINT fk_routing_operations_routing FOREIGN KEY (routing_id) REFERENCES routings(id) ON DELETE CASCADE,
  CONSTRAINT fk_routing_operations_work_center FOREIGN KEY (work_center_id) REFERENCES work_centers(id) ON DELETE CASCADE
);

-- 生产计划表
CREATE TABLE IF NOT EXISTS production_plans (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  created_by INTEGER,
  updated_by INTEGER,
  is_active INTEGER DEFAULT 1,
  code TEXT NOT NULL,
  name TEXT NOT NULL,
  description TEXT,
  plan_date DATE NOT NULL,
  start_date DATE,
  end_date DATE,
  planner_id INTEGER,
  approved_by INTEGER,
  approved_at DATETIME,
  status TEXT DEFAULT 'DRAFT',
  CONSTRAINT uq_production_plans_code UNIQUE (code),
  CONSTRAINT fk_production_plans_planner FOREIGN KEY (planner_id) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_production_plans_approved_by FOREIGN KEY (approved_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 生产计划明细表
CREATE TABLE IF NOT EXISTS production_plan_items (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  production_plan_id INTEGER NOT NULL,
  item_id INTEGER NOT NULL,
  planned_quantity REAL NOT NULL,
  unit_id INTEGER,
  required_date DATE,
  priority INTEGER DEFAULT 5,
  sales_order_id INTEGER,
  notes TEXT,
  CONSTRAINT fk_production_plan_items_production_plan FOREIGN KEY (production_plan_id) REFERENCES production_plans(id) ON DELETE CASCADE,
  CONSTRAINT fk_production_plan_items_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE,
  CONSTRAINT fk_production_plan_items_unit FOREIGN KEY (unit_id) REFERENCES units(id) ON DELETE SET NULL,
  CONSTRAINT fk_production_plan_items_sales_order FOREIGN KEY (sales_order_id) REFERENCES sales_orders(id) ON DELETE SET NULL
);

-- 生产订单表
CREATE TABLE IF NOT EXISTS production_orders (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  created_by INTEGER,
  updated_by INTEGER,
  is_active INTEGER DEFAULT 1,
  code TEXT NOT NULL,
  name TEXT NOT NULL,
  description TEXT,
  item_id INTEGER NOT NULL,
  bom_id INTEGER,
  routing_id INTEGER,
  planned_quantity REAL NOT NULL,
  produced_quantity REAL DEFAULT 0,
  scrapped_quantity REAL DEFAULT 0,
  unit_id INTEGER,
  planned_start_date DATE,
  planned_end_date DATE,
  actual_start_date DATE,
  actual_end_date DATE,
  priority INTEGER DEFAULT 5,
  sales_order_id INTEGER,
  production_plan_item_id INTEGER,
  responsible_person_id INTEGER,
  status TEXT DEFAULT 'PLANNED',
  CONSTRAINT uq_production_orders_code UNIQUE (code),
  CONSTRAINT fk_production_orders_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE,
  CONSTRAINT fk_production_orders_bom FOREIGN KEY (bom_id) REFERENCES boms(id) ON DELETE SET NULL,
  CONSTRAINT fk_production_orders_routing FOREIGN KEY (routing_id) REFERENCES routings(id) ON DELETE SET NULL,
  CONSTRAINT fk_production_orders_unit FOREIGN KEY (unit_id) REFERENCES units(id) ON DELETE SET NULL,
  CONSTRAINT fk_production_orders_sales_order FOREIGN KEY (sales_order_id) REFERENCES sales_orders(id) ON DELETE SET NULL,
  CONSTRAINT fk_production_orders_production_plan_item FOREIGN KEY (production_plan_item_id) REFERENCES production_plan_items(id) ON DELETE SET NULL,
  CONSTRAINT fk_production_orders_responsible_person FOREIGN KEY (responsible_person_id) REFERENCES users(id) ON DELETE SET NULL
);

-- 生产订单物料需求表
CREATE TABLE IF NOT EXISTS production_order_materials (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  production_order_id INTEGER NOT NULL,
  item_id INTEGER NOT NULL,
  required_quantity REAL NOT NULL,
  issued_quantity REAL DEFAULT 0,
  unit_id INTEGER,
  warehouse_id INTEGER,
  location_id INTEGER,
  required_date DATE,
  status TEXT DEFAULT 'PLANNED',
  CONSTRAINT fk_production_order_materials_production_order FOREIGN KEY (production_order_id) REFERENCES production_orders(id) ON DELETE CASCADE,
  CONSTRAINT fk_production_order_materials_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE,
  CONSTRAINT fk_production_order_materials_unit FOREIGN KEY (unit_id) REFERENCES units(id) ON DELETE SET NULL,
  CONSTRAINT fk_production_order_materials_warehouse FOREIGN KEY (warehouse_id) REFERENCES warehouses(id) ON DELETE SET NULL,
  CONSTRAINT fk_production_order_materials_location FOREIGN KEY (location_id) REFERENCES locations(id) ON DELETE SET NULL
);

-- 生产订单操作表
CREATE TABLE IF NOT EXISTS production_order_operations (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  production_order_id INTEGER NOT NULL,
  operation_number INTEGER NOT NULL,
  work_center_id INTEGER NOT NULL,
  operation_name TEXT NOT NULL,
  description TEXT,
  planned_setup_time INTEGER DEFAULT 0,
  planned_run_time INTEGER DEFAULT 0,
  actual_setup_time INTEGER DEFAULT 0,
  actual_run_time INTEGER DEFAULT 0,
  planned_start_date DATETIME,
  planned_end_date DATETIME,
  actual_start_date DATETIME,
  actual_end_date DATETIME,
  status TEXT DEFAULT 'PLANNED',
  CONSTRAINT fk_production_order_operations_production_order FOREIGN KEY (production_order_id) REFERENCES production_orders(id) ON DELETE CASCADE,
  CONSTRAINT fk_production_order_operations_work_center FOREIGN KEY (work_center_id) REFERENCES work_centers(id) ON DELETE CASCADE
);

-- 生产报工表
CREATE TABLE IF NOT EXISTS work_orders (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  created_by INTEGER,
  updated_by INTEGER,
  is_active INTEGER DEFAULT 1,
  code TEXT NOT NULL,
  production_order_id INTEGER NOT NULL,
  operation_id INTEGER NOT NULL,
  work_center_id INTEGER NOT NULL,
  operator_id INTEGER,
  work_date DATE NOT NULL,
  start_time DATETIME,
  end_time DATETIME,
  setup_time_minutes INTEGER DEFAULT 0,
  run_time_minutes INTEGER DEFAULT 0,
  quantity_completed REAL DEFAULT 0,
  quantity_scrapped REAL DEFAULT 0,
  notes TEXT,
  status TEXT DEFAULT 'COMPLETED',
  CONSTRAINT uq_work_orders_code UNIQUE (code),
  CONSTRAINT fk_work_orders_production_order FOREIGN KEY (production_order_id) REFERENCES production_orders(id) ON DELETE CASCADE,
  CONSTRAINT fk_work_orders_operation FOREIGN KEY (operation_id) REFERENCES production_order_operations(id) ON DELETE CASCADE,
  CONSTRAINT fk_work_orders_work_center FOREIGN KEY (work_center_id) REFERENCES work_centers(id) ON DELETE CASCADE,
  CONSTRAINT fk_work_orders_operator FOREIGN KEY (operator_id) REFERENCES users(id) ON DELETE SET NULL
);

-- 质量检验表
CREATE TABLE IF NOT EXISTS quality_inspections (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  created_by INTEGER,
  updated_by INTEGER,
  is_active INTEGER DEFAULT 1,
  code TEXT NOT NULL,
  name TEXT NOT NULL,
  description TEXT,
  inspection_type TEXT DEFAULT 'INCOMING',
  reference_type TEXT,
  reference_id INTEGER,
  item_id INTEGER NOT NULL,
  batch_number TEXT,
  quantity_inspected REAL NOT NULL,
  quantity_passed REAL DEFAULT 0,
  quantity_failed REAL DEFAULT 0,
  unit_id INTEGER,
  inspector_id INTEGER,
  inspection_date DATE NOT NULL,
  inspection_criteria TEXT,
  inspection_results TEXT,
  defect_description TEXT,
  corrective_action TEXT,
  status TEXT DEFAULT 'PENDING',
  CONSTRAINT uq_quality_inspections_code UNIQUE (code),
  CONSTRAINT fk_quality_inspections_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE,
  CONSTRAINT fk_quality_inspections_unit FOREIGN KEY (unit_id) REFERENCES units(id) ON DELETE SET NULL,
  CONSTRAINT fk_quality_inspections_inspector FOREIGN KEY (inspector_id) REFERENCES users(id) ON DELETE SET NULL
);

-- ============================================================================
-- 8. 项目管理模块 (14个表) - 注意：这些表已存在于 project_tables_sqlite.sql 中
-- ============================================================================

-- 项目表 (已存在)
CREATE TABLE IF NOT EXISTS projects (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  created_by INTEGER,
  updated_by INTEGER,
  is_active INTEGER DEFAULT 1,
  code TEXT NOT NULL,
  name TEXT NOT NULL,
  description TEXT,
  customer_id INTEGER,
  project_manager_id INTEGER,
  start_date DATE,
  end_date DATE,
  budget REAL DEFAULT 0,
  actual_cost REAL DEFAULT 0,
  status TEXT DEFAULT 'PLANNING',
  priority TEXT DEFAULT 'MEDIUM',
  CONSTRAINT uq_projects_code UNIQUE (code),
  CONSTRAINT fk_projects_customer FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE SET NULL,
  CONSTRAINT fk_projects_project_manager FOREIGN KEY (project_manager_id) REFERENCES users(id) ON DELETE SET NULL
);

-- 任务表 (已存在)
CREATE TABLE IF NOT EXISTS tasks (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  created_by INTEGER,
  updated_by INTEGER,
  is_active INTEGER DEFAULT 1,
  code TEXT NOT NULL,
  name TEXT NOT NULL,
  description TEXT,
  project_id INTEGER NOT NULL,
  parent_task_id INTEGER,
  assigned_to INTEGER,
  start_date DATE,
  end_date DATE,
  estimated_hours REAL DEFAULT 0,
  actual_hours REAL DEFAULT 0,
  progress INTEGER DEFAULT 0,
  priority TEXT DEFAULT 'MEDIUM',
  status TEXT DEFAULT 'TODO',
  CONSTRAINT uq_tasks_code UNIQUE (code),
  CONSTRAINT fk_tasks_project FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
  CONSTRAINT fk_tasks_parent_task FOREIGN KEY (parent_task_id) REFERENCES tasks(id) ON DELETE SET NULL,
  CONSTRAINT fk_tasks_assigned_to FOREIGN KEY (assigned_to) REFERENCES users(id) ON DELETE SET NULL
);

-- 里程碑表 (已存在)
CREATE TABLE IF NOT EXISTS milestones (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  created_by INTEGER,
  updated_by INTEGER,
  is_active INTEGER DEFAULT 1,
  code TEXT NOT NULL,
  name TEXT NOT NULL,
  description TEXT,
  project_id INTEGER NOT NULL,
  due_date DATE NOT NULL,
  completion_date DATE,
  status TEXT DEFAULT 'PENDING',
  CONSTRAINT uq_milestones_code UNIQUE (code),
  CONSTRAINT fk_milestones_project FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
);

-- 时间记录表 (已存在)
CREATE TABLE IF NOT EXISTS time_entries (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  created_by INTEGER,
  updated_by INTEGER,
  is_active INTEGER DEFAULT 1,
  user_id INTEGER NOT NULL,
  project_id INTEGER,
  task_id INTEGER,
  work_date DATE NOT NULL,
  start_time DATETIME,
  end_time DATETIME,
  hours REAL NOT NULL,
  description TEXT,
  billable INTEGER DEFAULT 1,
  hourly_rate REAL DEFAULT 0,
  CONSTRAINT fk_time_entries_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  CONSTRAINT fk_time_entries_project FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE SET NULL,
  CONSTRAINT fk_time_entries_task FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE SET NULL
);

-- 项目成员表 (已存在)
CREATE TABLE IF NOT EXISTS project_members (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  project_id INTEGER NOT NULL,
  user_id INTEGER NOT NULL,
  role TEXT DEFAULT 'MEMBER',
  hourly_rate REAL DEFAULT 0,
  joined_date DATE,
  left_date DATE,
  is_active INTEGER DEFAULT 1,
  CONSTRAINT fk_project_members_project FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
  CONSTRAINT fk_project_members_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  CONSTRAINT uq_project_members_project_user UNIQUE (project_id, user_id)
);

-- 项目费用表 (已存在)
CREATE TABLE IF NOT EXISTS project_expenses (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  created_by INTEGER,
  updated_by INTEGER,
  is_active INTEGER DEFAULT 1,
  code TEXT NOT NULL,
  project_id INTEGER NOT NULL,
  expense_date DATE NOT NULL,
  category TEXT NOT NULL,
  description TEXT,
  amount REAL NOT NULL,
  currency TEXT DEFAULT 'CNY',
  receipt_number TEXT,
  approved_by INTEGER,
  approved_at DATETIME,
  status TEXT DEFAULT 'PENDING',
  CONSTRAINT uq_project_expenses_code UNIQUE (code),
  CONSTRAINT fk_project_expenses_project FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
  CONSTRAINT fk_project_expenses_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_project_expenses_approved_by FOREIGN KEY (approved_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 任务评论表 (已存在)
CREATE TABLE IF NOT EXISTS task_comments (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  created_by INTEGER,
  updated_by INTEGER,
  is_active INTEGER DEFAULT 1,
  task_id INTEGER NOT NULL,
  comment TEXT NOT NULL,
  CONSTRAINT fk_task_comments_task FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE,
  CONSTRAINT fk_task_comments_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 项目文档表 (已存在)
CREATE TABLE IF NOT EXISTS project_documents (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  created_by INTEGER,
  updated_by INTEGER,
  is_active INTEGER DEFAULT 1,
  project_id INTEGER NOT NULL,
  name TEXT NOT NULL,
  description TEXT,
  file_path TEXT NOT NULL,
  file_size INTEGER DEFAULT 0,
  file_type TEXT,
  version TEXT DEFAULT '1.0',
  CONSTRAINT fk_project_documents_project FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
  CONSTRAINT fk_project_documents_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 项目资源表 (已存在)
CREATE TABLE IF NOT EXISTS project_resources (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  created_by INTEGER,
  updated_by INTEGER,
  is_active INTEGER DEFAULT 1,
  project_id INTEGER NOT NULL,
  resource_type TEXT NOT NULL,
  resource_name TEXT NOT NULL,
  description TEXT,
  quantity REAL DEFAULT 1,
  unit TEXT,
  cost_per_unit REAL DEFAULT 0,
  total_cost REAL DEFAULT 0,
  allocated_date DATE,
  CONSTRAINT fk_project_resources_project FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
  CONSTRAINT fk_project_resources_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 项目报告表 (已存在)
CREATE TABLE IF NOT EXISTS project_reports (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  created_by INTEGER,
  updated_by INTEGER,
  is_active INTEGER DEFAULT 1,
  project_id INTEGER NOT NULL,
  report_type TEXT NOT NULL,
  report_date DATE NOT NULL,
  title TEXT NOT NULL,
  content TEXT,
  status TEXT DEFAULT 'DRAFT',
  CONSTRAINT fk_project_reports_project FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
  CONSTRAINT fk_project_reports_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL
);

-- ============================================================================
-- 9. 系统管理模块 (9个表)
-- ============================================================================

-- 系统配置表
CREATE TABLE IF NOT EXISTS system_configs (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  created_by INTEGER,
  updated_by INTEGER,
  is_active INTEGER DEFAULT 1,
  config_key TEXT NOT NULL,
  config_value TEXT,
  config_type TEXT DEFAULT 'STRING',
  description TEXT,
  category TEXT DEFAULT 'GENERAL',
  is_encrypted INTEGER DEFAULT 0,
  is_readonly INTEGER DEFAULT 0,
  CONSTRAINT uq_system_configs_config_key UNIQUE (config_key),
  CONSTRAINT fk_system_configs_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_system_configs_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 审计日志表
CREATE TABLE IF NOT EXISTS audit_logs (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  user_id INTEGER,
  action TEXT NOT NULL,
  resource_type TEXT NOT NULL,
  resource_id INTEGER,
  old_values TEXT,
  new_values TEXT,
  ip_address TEXT,
  user_agent TEXT,
  session_id TEXT,
  CONSTRAINT fk_audit_logs_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
);

-- 通知模板表
CREATE TABLE IF NOT EXISTS notification_templates (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  created_by INTEGER,
  updated_by INTEGER,
  is_active INTEGER DEFAULT 1,
  code TEXT NOT NULL,
  name TEXT NOT NULL,
  description TEXT,
  template_type TEXT DEFAULT 'EMAIL',
  subject TEXT,
  content TEXT NOT NULL,
  variables TEXT,
  is_system INTEGER DEFAULT 0,
  CONSTRAINT uq_notification_templates_code UNIQUE (code),
  CONSTRAINT fk_notification_templates_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_notification_templates_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 通知表
CREATE TABLE IF NOT EXISTS notifications (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  created_by INTEGER,
  updated_by INTEGER,
  is_active INTEGER DEFAULT 1,
  user_id INTEGER NOT NULL,
  title TEXT NOT NULL,
  content TEXT NOT NULL,
  notification_type TEXT DEFAULT 'INFO',
  reference_type TEXT,
  reference_id INTEGER,
  is_read INTEGER DEFAULT 0,
  read_at DATETIME,
  sent_at DATETIME,
  CONSTRAINT fk_notifications_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  CONSTRAINT fk_notifications_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 文件管理表
CREATE TABLE IF NOT EXISTS files (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  created_by INTEGER,
  updated_by INTEGER,
  is_active INTEGER DEFAULT 1,
  original_name TEXT NOT NULL,
  stored_name TEXT NOT NULL,
  file_path TEXT NOT NULL,
  file_size INTEGER NOT NULL,
  file_type TEXT,
  mime_type TEXT,
  md5_hash TEXT,
  reference_type TEXT,
  reference_id INTEGER,
  is_public INTEGER DEFAULT 0,
  download_count INTEGER DEFAULT 0,
  CONSTRAINT uq_files_stored_name UNIQUE (stored_name),
  CONSTRAINT fk_files_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 数据字典表
CREATE TABLE IF NOT EXISTS dictionaries (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  created_by INTEGER,
  updated_by INTEGER,
  is_active INTEGER DEFAULT 1,
  dict_type TEXT NOT NULL,
  dict_key TEXT NOT NULL,
  dict_value TEXT NOT NULL,
  dict_label TEXT NOT NULL,
  description TEXT,
  sort_order INTEGER DEFAULT 0,
  is_system INTEGER DEFAULT 0,
  CONSTRAINT uq_dictionaries_type_key UNIQUE (dict_type, dict_key),
  CONSTRAINT fk_dictionaries_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_dictionaries_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 定时任务表
CREATE TABLE IF NOT EXISTS scheduled_jobs (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  created_by INTEGER,
  updated_by INTEGER,
  is_active INTEGER DEFAULT 1,
  job_name TEXT NOT NULL,
  job_group TEXT DEFAULT 'DEFAULT',
  description TEXT,
  cron_expression TEXT NOT NULL,
  job_class TEXT NOT NULL,
  job_data TEXT,
  is_enabled INTEGER DEFAULT 1,
  last_run_time DATETIME,
  next_run_time DATETIME,
  run_count INTEGER DEFAULT 0,
  CONSTRAINT uq_scheduled_jobs_name_group UNIQUE (job_name, job_group),
  CONSTRAINT fk_scheduled_jobs_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_scheduled_jobs_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 任务执行日志表
CREATE TABLE IF NOT EXISTS job_execution_logs (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  job_id INTEGER NOT NULL,
  execution_time DATETIME NOT NULL,
  duration_ms INTEGER DEFAULT 0,
  status TEXT DEFAULT 'SUCCESS',
  result_message TEXT,
  error_message TEXT,
  CONSTRAINT fk_job_execution_logs_job FOREIGN KEY (job_id) REFERENCES scheduled_jobs(id) ON DELETE CASCADE
);

-- 系统监控表
CREATE TABLE IF NOT EXISTS system_monitors (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  monitor_time DATETIME NOT NULL,
  cpu_usage REAL DEFAULT 0,
  memory_usage REAL DEFAULT 0,
  disk_usage REAL DEFAULT 0,
  active_users INTEGER DEFAULT 0,
  active_sessions INTEGER DEFAULT 0,
  database_connections INTEGER DEFAULT 0,
  response_time_ms INTEGER DEFAULT 0,
  error_count INTEGER DEFAULT 0
);

-- ============================================================================
-- 10. 遗留模块 (2个表)
-- ============================================================================

-- 数据迁移记录表
CREATE TABLE IF NOT EXISTS data_migrations (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  created_by INTEGER,
  updated_by INTEGER,
  is_active INTEGER DEFAULT 1,
  migration_name TEXT NOT NULL,
  source_system TEXT NOT NULL,
  target_system TEXT DEFAULT 'GalaxyERP',
  migration_type TEXT DEFAULT 'DATA_IMPORT',
  start_time DATETIME,
  end_time DATETIME,
  total_records INTEGER DEFAULT 0,
  success_records INTEGER DEFAULT 0,
  failed_records INTEGER DEFAULT 0,
  status TEXT DEFAULT 'PENDING',
  error_message TEXT,
  migration_config TEXT,
  CONSTRAINT uq_data_migrations_migration_name UNIQUE (migration_name),
  CONSTRAINT fk_data_migrations_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_data_migrations_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 数据迁移详细日志表
CREATE TABLE IF NOT EXISTS migration_logs (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME,
  updated_at DATETIME,
  deleted_at DATETIME,
  migration_id INTEGER NOT NULL,
  table_name TEXT NOT NULL,
  record_id TEXT,
  operation_type TEXT DEFAULT 'INSERT',
  status TEXT DEFAULT 'SUCCESS',
  error_message TEXT,
  old_data TEXT,
  new_data TEXT,
  CONSTRAINT fk_migration_logs_migration FOREIGN KEY (migration_id) REFERENCES data_migrations(id) ON DELETE CASCADE
);

-- ============================================================================
-- 初始数据插入
-- ============================================================================

-- 1. 用户管理模块初始数据
INSERT OR IGNORE INTO companies (id, created_at, updated_at, code, name, description, address, phone, email, website, tax_number, legal_representative, registration_date, status) VALUES
(1, datetime('now'), datetime('now'), 'GALAXY001', '银河科技有限公司', '专业的ERP系统开发公司', '北京市朝阳区科技园区1号楼', '010-12345678', 'info@galaxy-tech.com', 'https://www.galaxy-tech.com', '91110000123456789X', '张三', '2020-01-01', 'ACTIVE');

INSERT OR IGNORE INTO departments (id, created_at, updated_at, code, name, description, parent_id, manager_id, status) VALUES
(1, datetime('now'), datetime('now'), 'IT001', '信息技术部', '负责公司IT系统开发和维护', NULL, NULL, 'ACTIVE'),
(2, datetime('now'), datetime('now'), 'HR001', '人力资源部', '负责人力资源管理', NULL, NULL, 'ACTIVE'),
(3, datetime('now'), datetime('now'), 'FIN001', '财务部', '负责财务管理', NULL, NULL, 'ACTIVE'),
(4, datetime('now'), datetime('now'), 'SALES001', '销售部', '负责销售业务', NULL, NULL, 'ACTIVE'),
(5, datetime('now'), datetime('now'), 'PROD001', '生产部', '负责生产管理', NULL, NULL, 'ACTIVE');

INSERT OR IGNORE INTO roles (id, created_at, updated_at, code, name, description, is_system) VALUES
(1, datetime('now'), datetime('now'), 'ADMIN', '系统管理员', '拥有系统所有权限', 1),
(2, datetime('now'), datetime('now'), 'USER', '普通用户', '基本用户权限', 1),
(3, datetime('now'), datetime('now'), 'MANAGER', '部门经理', '部门管理权限', 0),
(4, datetime('now'), datetime('now'), 'FINANCE', '财务人员', '财务模块权限', 0),
(5, datetime('now'), datetime('now'), 'SALES', '销售人员', '销售模块权限', 0);

INSERT OR IGNORE INTO permissions (id, created_at, updated_at, code, name, description, resource, action, is_system) VALUES
(1, datetime('now'), datetime('now'), 'USER_READ', '查看用户', '查看用户信息', 'USER', 'READ', 1),
(2, datetime('now'), datetime('now'), 'USER_WRITE', '编辑用户', '创建和编辑用户', 'USER', 'WRITE', 1),
(3, datetime('now'), datetime('now'), 'USER_DELETE', '删除用户', '删除用户', 'USER', 'DELETE', 1),
(4, datetime('now'), datetime('now'), 'ROLE_MANAGE', '角色管理', '管理角色和权限', 'ROLE', 'MANAGE', 1),
(5, datetime('now'), datetime('now'), 'SYSTEM_CONFIG', '系统配置', '系统配置管理', 'SYSTEM', 'CONFIG', 1);

INSERT OR IGNORE INTO users (id, created_at, updated_at, username, email, password_hash, first_name, last_name, phone, is_active, is_admin, last_login_at) VALUES
(1, datetime('now'), datetime('now'), 'admin', 'admin@galaxy-tech.com', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iKWTn5qUjRLrbWZ2xYoda0xbdJHK', '系统', '管理员', '13800138000', 1, 1, NULL);

-- 更新用户的部门信息
UPDATE users SET department_id = 1 WHERE id = 1;

INSERT OR IGNORE INTO user_roles (created_at, updated_at, user_id, role_id, assigned_at, is_active) VALUES 
(datetime('now'), datetime('now'), 1, 1, datetime('now'), 1);

INSERT OR IGNORE INTO role_permissions (role_id, permission_id) VALUES
(1, 1), (1, 2), (1, 3), (1, 4), (1, 5),
(2, 1),
(3, 1), (3, 2),
(4, 1),
(5, 1);

-- 2. 人力资源模块初始数据
INSERT OR IGNORE INTO salary_grades (id, created_at, updated_at, code, name, description, min_salary, max_salary, currency) VALUES
(1, datetime('now'), datetime('now'), 'G01', '初级', '初级员工薪资等级', 5000, 8000, 'CNY'),
(2, datetime('now'), datetime('now'), 'G02', '中级', '中级员工薪资等级', 8000, 15000, 'CNY'),
(3, datetime('now'), datetime('now'), 'G03', '高级', '高级员工薪资等级', 15000, 25000, 'CNY'),
(4, datetime('now'), datetime('now'), 'G04', '专家', '专家级员工薪资等级', 25000, 40000, 'CNY'),
(5, datetime('now'), datetime('now'), 'G05', '总监', '总监级员工薪资等级', 40000, 80000, 'CNY');

INSERT OR IGNORE INTO leave_types (id, created_at, updated_at, code, name, description, max_days_per_year, is_paid, requires_approval) VALUES
(1, datetime('now'), datetime('now'), 'ANNUAL', '年假', '带薪年假', 10, 1, 1),
(2, datetime('now'), datetime('now'), 'SICK', '病假', '病假', 30, 1, 1),
(3, datetime('now'), datetime('now'), 'PERSONAL', '事假', '个人事假', 5, 0, 1),
(4, datetime('now'), datetime('now'), 'MATERNITY', '产假', '产假', 98, 1, 1),
(5, datetime('now'), datetime('now'), 'PATERNITY', '陪产假', '陪产假', 15, 1, 1);

-- 3. 库存管理模块初始数据
INSERT OR IGNORE INTO item_categories (id, created_at, updated_at, code, name, description, parent_id) VALUES
(1, datetime('now'), datetime('now'), 'RAW', '原材料', '生产用原材料', NULL),
(2, datetime('now'), datetime('now'), 'SEMI', '半成品', '生产过程中的半成品', NULL),
(3, datetime('now'), datetime('now'), 'FINISHED', '成品', '最终产品', NULL),
(4, datetime('now'), datetime('now'), 'CONSUMABLE', '消耗品', '办公用品等消耗品', NULL),
(5, datetime('now'), datetime('now'), 'SPARE', '备件', '设备备件', NULL);

INSERT OR IGNORE INTO units (id, created_at, updated_at, code, name, description, unit_type, base_unit_id, conversion_factor) VALUES
(1, datetime('now'), datetime('now'), 'PCS', '个', '计数单位', 'COUNT', NULL, 1),
(2, datetime('now'), datetime('now'), 'KG', '千克', '重量单位', 'WEIGHT', NULL, 1),
(3, datetime('now'), datetime('now'), 'G', '克', '重量单位', 'WEIGHT', 2, 0.001),
(4, datetime('now'), datetime('now'), 'M', '米', '长度单位', 'LENGTH', NULL, 1),
(5, datetime('now'), datetime('now'), 'CM', '厘米', '长度单位', 'LENGTH', 4, 0.01),
(6, datetime('now'), datetime('now'), 'L', '升', '体积单位', 'VOLUME', NULL, 1),
(7, datetime('now'), datetime('now'), 'ML', '毫升', '体积单位', 'VOLUME', 6, 0.001);

INSERT OR IGNORE INTO warehouses (id, created_at, updated_at, code, name, description, address, manager_id, warehouse_type) VALUES
(1, datetime('now'), datetime('now'), 'WH001', '主仓库', '公司主要仓库', '北京市朝阳区仓储区1号', NULL, 'MAIN'),
(2, datetime('now'), datetime('now'), 'WH002', '原料仓', '原材料专用仓库', '北京市朝阳区仓储区2号', NULL, 'RAW_MATERIAL'),
(3, datetime('now'), datetime('now'), 'WH003', '成品仓', '成品存储仓库', '北京市朝阳区仓储区3号', NULL, 'FINISHED_GOODS');

INSERT OR IGNORE INTO locations (id, created_at, updated_at, warehouse_id, code, name, description, location_type) VALUES
(1, datetime('now'), datetime('now'), 1, 'A01-01', 'A区1排1位', '主仓库A区第1排第1位', 'SHELF'),
(2, datetime('now'), datetime('now'), 1, 'A01-02', 'A区1排2位', '主仓库A区第1排第2位', 'SHELF'),
(3, datetime('now'), datetime('now'), 2, 'B01-01', 'B区1排1位', '原料仓B区第1排第1位', 'SHELF'),
(4, datetime('now'), datetime('now'), 3, 'C01-01', 'C区1排1位', '成品仓C区第1排第1位', 'SHELF');

-- 4. 财务会计模块初始数据
INSERT OR IGNORE INTO accounts (id, created_at, updated_at, code, name, description, account_type, parent_id, is_active) VALUES
(1, datetime('now'), datetime('now'), '1001', '库存现金', '企业库存的现金', 'ASSET', NULL, 1),
(2, datetime('now'), datetime('now'), '1002', '银行存款', '企业在银行的存款', 'ASSET', NULL, 1),
(3, datetime('now'), datetime('now'), '1122', '应收账款', '企业因销售商品提供劳务而应收的款项', 'ASSET', NULL, 1),
(4, datetime('now'), datetime('now'), '1403', '原材料', '企业库存的各种原材料', 'ASSET', NULL, 1),
(5, datetime('now'), datetime('now'), '2202', '应付账款', '企业因购买材料商品接受劳务而应付的款项', 'LIABILITY', NULL, 1),
(6, datetime('now'), datetime('now'), '4001', '主营业务收入', '企业确认的销售商品提供劳务的收入', 'REVENUE', NULL, 1),
(7, datetime('now'), datetime('now'), '5401', '主营业务成本', '企业确认销售商品提供劳务的成本', 'EXPENSE', NULL, 1);

INSERT OR IGNORE INTO fiscal_periods (id, created_at, updated_at, code, name, start_date, end_date, fiscal_year, period_number, is_closed) VALUES
(1, datetime('now'), datetime('now'), '2024-01', '2024年1月', '2024-01-01', '2024-01-31', 2024, 1, 0),
(2, datetime('now'), datetime('now'), '2024-02', '2024年2月', '2024-02-01', '2024-02-29', 2024, 2, 0),
(3, datetime('now'), datetime('now'), '2024-03', '2024年3月', '2024-03-01', '2024-03-31', 2024, 3, 0);

-- 5. 生产管理模块初始数据
INSERT OR IGNORE INTO work_centers (id, created_at, updated_at, code, name, description, work_center_type, capacity_per_hour, efficiency_factor, cost_per_hour) VALUES
(1, datetime('now'), datetime('now'), 'WC001', '装配线1', '主要装配生产线', 'ASSEMBLY', 10, 0.95, 50),
(2, datetime('now'), datetime('now'), 'WC002', '加工中心1', 'CNC加工中心', 'MACHINE', 5, 0.90, 80),
(3, datetime('now'), datetime('now'), 'WC003', '包装线1', '产品包装生产线', 'PACKAGING', 20, 0.98, 30),
(4, datetime('now'), datetime('now'), 'WC004', '质检站1', '产品质量检验工作站', 'INSPECTION', 15, 1.0, 40);

-- 6. 系统管理模块初始数据
INSERT OR IGNORE INTO system_configs (id, created_at, updated_at, config_key, config_value, config_type, description, category) VALUES
(1, datetime('now'), datetime('now'), 'SYSTEM_NAME', 'Galaxy ERP', 'STRING', '系统名称', 'GENERAL'),
(2, datetime('now'), datetime('now'), 'SYSTEM_VERSION', '1.0.0', 'STRING', '系统版本', 'GENERAL'),
(3, datetime('now'), datetime('now'), 'DEFAULT_CURRENCY', 'CNY', 'STRING', '默认货币', 'FINANCIAL'),
(4, datetime('now'), datetime('now'), 'SESSION_TIMEOUT', '3600', 'INTEGER', '会话超时时间(秒)', 'SECURITY'),
(5, datetime('now'), datetime('now'), 'PASSWORD_MIN_LENGTH', '8', 'INTEGER', '密码最小长度', 'SECURITY'),
(6, datetime('now'), datetime('now'), 'BACKUP_RETENTION_DAYS', '30', 'INTEGER', '备份保留天数', 'SYSTEM'),
(7, datetime('now'), datetime('now'), 'EMAIL_SMTP_HOST', 'smtp.galaxy-tech.com', 'STRING', 'SMTP服务器地址', 'EMAIL'),
(8, datetime('now'), datetime('now'), 'EMAIL_SMTP_PORT', '587', 'INTEGER', 'SMTP服务器端口', 'EMAIL');

INSERT OR IGNORE INTO dictionaries (id, created_at, updated_at, dict_type, dict_key, dict_value, dict_label, description, sort_order, is_system) VALUES
(1, datetime('now'), datetime('now'), 'USER_STATUS', 'ACTIVE', 'ACTIVE', '激活', '用户状态-激活', 1, 1),
(2, datetime('now'), datetime('now'), 'USER_STATUS', 'INACTIVE', 'INACTIVE', '停用', '用户状态-停用', 2, 1),
(3, datetime('now'), datetime('now'), 'USER_STATUS', 'LOCKED', 'LOCKED', '锁定', '用户状态-锁定', 3, 1),
(4, datetime('now'), datetime('now'), 'GENDER', 'MALE', 'MALE', '男', '性别-男', 1, 1),
(5, datetime('now'), datetime('now'), 'GENDER', 'FEMALE', 'FEMALE', '女', '性别-女', 2, 1),
(6, datetime('now'), datetime('now'), 'PRIORITY', 'LOW', 'LOW', '低', '优先级-低', 1, 1),
(7, datetime('now'), datetime('now'), 'PRIORITY', 'MEDIUM', 'MEDIUM', '中', '优先级-中', 2, 1),
(8, datetime('now'), datetime('now'), 'PRIORITY', 'HIGH', 'HIGH', '高', '优先级-高', 3, 1),
(9, datetime('now'), datetime('now'), 'PRIORITY', 'URGENT', 'URGENT', '紧急', '优先级-紧急', 4, 1);

INSERT OR IGNORE INTO notification_templates (id, created_at, updated_at, code, name, description, template_type, subject, content, is_system) VALUES
(1, datetime('now'), datetime('now'), 'USER_WELCOME', '用户欢迎邮件', '新用户注册欢迎邮件模板', 'EMAIL', '欢迎加入Galaxy ERP系统', '亲爱的{{user_name}}，欢迎您加入Galaxy ERP系统！您的账号已经创建成功。', 1),
(2, datetime('now'), datetime('now'), 'PASSWORD_RESET', '密码重置邮件', '用户密码重置邮件模板', 'EMAIL', '密码重置通知', '您好{{user_name}}，您的密码重置请求已处理，请使用新密码登录系统。', 1),
(3, datetime('now'), datetime('now'), 'TASK_ASSIGNED', '任务分配通知', '任务分配通知模板', 'SYSTEM', '新任务分配', '您有新的任务被分配：{{task_name}}，请及时处理。', 1);

-- 7. 遗留模块初始数据
INSERT OR IGNORE INTO data_migrations (id, created_at, updated_at, migration_name, source_system, target_system, migration_type, status) VALUES
(1, datetime('now'), datetime('now'), '用户数据迁移', 'Legacy System', 'Galaxy ERP', 'DATA_IMPORT', 'COMPLETED'),
(2, datetime('now'), datetime('now'), '产品数据迁移', 'Legacy System', 'Galaxy ERP', 'DATA_IMPORT', 'COMPLETED'),
(3, datetime('now'), datetime('now'), '客户数据迁移', 'Legacy System', 'Galaxy ERP', 'DATA_IMPORT', 'PENDING');

-- 重新启用外键约束
PRAGMA foreign_keys = ON;