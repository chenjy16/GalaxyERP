-- ============================================================================
-- Galaxy ERP 完整数据库初始化脚本 (PostgreSQL版本)
-- 包含所有101个表的定义和初始数据
-- 生成时间: 2024年
-- ============================================================================

-- 启用必要的扩展
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ============================================================================
-- 1. 用户管理模块 (8个表)
-- ============================================================================

-- 公司表
CREATE TABLE IF NOT EXISTS companies (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  created_by INTEGER,
  updated_by INTEGER,
  is_active BOOLEAN DEFAULT TRUE,
  code VARCHAR(50) NOT NULL,
  name VARCHAR(200) NOT NULL,
  description TEXT,
  address TEXT,
  phone VARCHAR(50),
  email VARCHAR(100),
  website VARCHAR(200),
  tax_number VARCHAR(50),
  legal_representative VARCHAR(100),
  registration_date DATE,
  status VARCHAR(20) DEFAULT 'ACTIVE',
  CONSTRAINT uq_companies_code UNIQUE (code)
);

-- 部门表
CREATE TABLE IF NOT EXISTS departments (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  created_by INTEGER,
  updated_by INTEGER,
  is_active BOOLEAN DEFAULT TRUE,
  code VARCHAR(50) NOT NULL,
  name VARCHAR(100) NOT NULL,
  description TEXT,
  parent_id INTEGER,
  manager_id INTEGER,
  status VARCHAR(20) DEFAULT 'ACTIVE',
  CONSTRAINT uq_departments_code UNIQUE (code),
  CONSTRAINT fk_departments_parent FOREIGN KEY (parent_id) REFERENCES departments(id) ON DELETE SET NULL
);

-- 用户表
CREATE TABLE IF NOT EXISTS users (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  created_by INTEGER,
  updated_by INTEGER,
  is_active BOOLEAN DEFAULT TRUE,
  username VARCHAR(50) NOT NULL,
  email VARCHAR(100) NOT NULL,
  password_hash VARCHAR(255) NOT NULL,
  first_name VARCHAR(50),
  last_name VARCHAR(50),
  phone VARCHAR(20),
  avatar VARCHAR(255),
  department_id INTEGER,
  is_admin BOOLEAN DEFAULT FALSE,
  last_login_at TIMESTAMP WITH TIME ZONE,
  password_changed_at TIMESTAMP WITH TIME ZONE,
  failed_login_attempts INTEGER DEFAULT 0,
  locked_until TIMESTAMP WITH TIME ZONE,
  CONSTRAINT uq_users_username UNIQUE (username),
  CONSTRAINT uq_users_email UNIQUE (email),
  CONSTRAINT fk_users_department FOREIGN KEY (department_id) REFERENCES departments(id) ON DELETE SET NULL
);

-- 角色表
CREATE TABLE IF NOT EXISTS roles (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  created_by INTEGER,
  updated_by INTEGER,
  is_active BOOLEAN DEFAULT TRUE,
  code VARCHAR(50) NOT NULL,
  name VARCHAR(100) NOT NULL,
  description TEXT,
  is_system BOOLEAN DEFAULT FALSE,
  CONSTRAINT uq_roles_code UNIQUE (code),
  CONSTRAINT fk_roles_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_roles_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 权限表
CREATE TABLE IF NOT EXISTS permissions (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  created_by INTEGER,
  updated_by INTEGER,
  is_active BOOLEAN DEFAULT TRUE,
  code VARCHAR(50) NOT NULL,
  name VARCHAR(100) NOT NULL,
  description TEXT,
  resource VARCHAR(50) NOT NULL,
  action VARCHAR(50) NOT NULL,
  is_system BOOLEAN DEFAULT FALSE,
  CONSTRAINT uq_permissions_code UNIQUE (code),
  CONSTRAINT fk_permissions_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_permissions_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 用户角色关联表
CREATE TABLE IF NOT EXISTS user_roles (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  user_id INTEGER NOT NULL,
  role_id INTEGER NOT NULL,
  CONSTRAINT fk_user_roles_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  CONSTRAINT fk_user_roles_role FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
  CONSTRAINT uq_user_roles_user_role UNIQUE (user_id, role_id)
);

-- 角色权限关联表
CREATE TABLE IF NOT EXISTS role_permissions (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  role_id INTEGER NOT NULL,
  permission_id INTEGER NOT NULL,
  CONSTRAINT fk_role_permissions_role FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
  CONSTRAINT fk_role_permissions_permission FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE,
  CONSTRAINT uq_role_permissions_role_permission UNIQUE (role_id, permission_id)
);

-- 用户会话表
CREATE TABLE IF NOT EXISTS user_sessions (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  user_id INTEGER NOT NULL,
  session_token VARCHAR(255) NOT NULL,
  ip_address INET,
  user_agent TEXT,
  expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
  is_active BOOLEAN DEFAULT TRUE,
  CONSTRAINT uq_user_sessions_session_token UNIQUE (session_token),
  CONSTRAINT fk_user_sessions_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- ============================================================================
-- 2. 人力资源模块 (10个表)
-- ============================================================================

-- 职位表
CREATE TABLE IF NOT EXISTS positions (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  created_by INTEGER,
  updated_by INTEGER,
  is_active BOOLEAN DEFAULT TRUE,
  code VARCHAR(50) NOT NULL,
  name VARCHAR(100) NOT NULL,
  description TEXT,
  department_id INTEGER,
  level INTEGER DEFAULT 1,
  min_salary DECIMAL(15,2) DEFAULT 0,
  max_salary DECIMAL(15,2) DEFAULT 0,
  requirements TEXT,
  responsibilities TEXT,
  status VARCHAR(20) DEFAULT 'ACTIVE',
  CONSTRAINT uq_positions_code UNIQUE (code),
  CONSTRAINT fk_positions_department FOREIGN KEY (department_id) REFERENCES departments(id) ON DELETE SET NULL,
  CONSTRAINT fk_positions_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_positions_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 员工表
CREATE TABLE IF NOT EXISTS employees (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  created_by INTEGER,
  updated_by INTEGER,
  is_active BOOLEAN DEFAULT TRUE,
  employee_number VARCHAR(50) NOT NULL,
  user_id INTEGER,
  first_name VARCHAR(50) NOT NULL,
  last_name VARCHAR(50) NOT NULL,
  email VARCHAR(100),
  phone VARCHAR(20),
  id_number VARCHAR(50),
  gender VARCHAR(10),
  birth_date DATE,
  hire_date DATE NOT NULL,
  termination_date DATE,
  department_id INTEGER,
  position_id INTEGER,
  manager_id INTEGER,
  employment_type VARCHAR(20) DEFAULT 'FULL_TIME',
  status VARCHAR(20) DEFAULT 'ACTIVE',
  CONSTRAINT uq_employees_employee_number UNIQUE (employee_number),
  CONSTRAINT uq_employees_user_id UNIQUE (user_id),
  CONSTRAINT uq_employees_id_number UNIQUE (id_number),
  CONSTRAINT fk_employees_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_employees_department FOREIGN KEY (department_id) REFERENCES departments(id) ON DELETE SET NULL,
  CONSTRAINT fk_employees_position FOREIGN KEY (position_id) REFERENCES positions(id) ON DELETE SET NULL,
  CONSTRAINT fk_employees_manager FOREIGN KEY (manager_id) REFERENCES employees(id) ON DELETE SET NULL,
  CONSTRAINT fk_employees_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_employees_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 薪资等级表
CREATE TABLE IF NOT EXISTS salary_grades (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  created_by INTEGER,
  updated_by INTEGER,
  is_active BOOLEAN DEFAULT TRUE,
  code VARCHAR(50) NOT NULL,
  name VARCHAR(100) NOT NULL,
  description TEXT,
  min_salary DECIMAL(15,2) NOT NULL,
  max_salary DECIMAL(15,2) NOT NULL,
  currency VARCHAR(3) DEFAULT 'CNY',
  CONSTRAINT uq_salary_grades_code UNIQUE (code),
  CONSTRAINT fk_salary_grades_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_salary_grades_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 薪资记录表
CREATE TABLE IF NOT EXISTS salary_records (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  created_by INTEGER,
  updated_by INTEGER,
  is_active BOOLEAN DEFAULT TRUE,
  employee_id INTEGER NOT NULL,
  salary_period VARCHAR(7) NOT NULL, -- YYYY-MM
  base_salary DECIMAL(15,2) DEFAULT 0,
  allowances DECIMAL(15,2) DEFAULT 0,
  overtime_pay DECIMAL(15,2) DEFAULT 0,
  bonus DECIMAL(15,2) DEFAULT 0,
  deductions DECIMAL(15,2) DEFAULT 0,
  tax DECIMAL(15,2) DEFAULT 0,
  social_insurance DECIMAL(15,2) DEFAULT 0,
  housing_fund DECIMAL(15,2) DEFAULT 0,
  net_salary DECIMAL(15,2) DEFAULT 0,
  currency VARCHAR(3) DEFAULT 'CNY',
  payment_date DATE,
  status VARCHAR(20) DEFAULT 'DRAFT',
  CONSTRAINT fk_salary_records_employee FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE CASCADE,
  CONSTRAINT fk_salary_records_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_salary_records_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 考勤表
CREATE TABLE IF NOT EXISTS attendances (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  created_by INTEGER,
  updated_by INTEGER,
  is_active BOOLEAN DEFAULT TRUE,
  employee_id INTEGER NOT NULL,
  attendance_date DATE NOT NULL,
  check_in_time TIMESTAMP WITH TIME ZONE,
  check_out_time TIMESTAMP WITH TIME ZONE,
  work_hours DECIMAL(4,2) DEFAULT 0,
  overtime_hours DECIMAL(4,2) DEFAULT 0,
  late_minutes INTEGER DEFAULT 0,
  early_leave_minutes INTEGER DEFAULT 0,
  status VARCHAR(20) DEFAULT 'NORMAL',
  notes TEXT,
  CONSTRAINT fk_attendances_employee FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE CASCADE,
  CONSTRAINT fk_attendances_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_attendances_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 请假类型表
CREATE TABLE IF NOT EXISTS leave_types (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  created_by INTEGER,
  updated_by INTEGER,
  is_active BOOLEAN DEFAULT TRUE,
  code VARCHAR(50) NOT NULL,
  name VARCHAR(100) NOT NULL,
  description TEXT,
  max_days_per_year INTEGER DEFAULT 0,
  is_paid BOOLEAN DEFAULT TRUE,
  requires_approval BOOLEAN DEFAULT TRUE,
  CONSTRAINT uq_leave_types_code UNIQUE (code),
  CONSTRAINT fk_leave_types_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_leave_types_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 请假申请表
CREATE TABLE IF NOT EXISTS leave_requests (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  created_by INTEGER,
  updated_by INTEGER,
  is_active BOOLEAN DEFAULT TRUE,
  employee_id INTEGER NOT NULL,
  leave_type_id INTEGER NOT NULL,
  start_date DATE NOT NULL,
  end_date DATE NOT NULL,
  days_requested DECIMAL(4,1) NOT NULL,
  reason TEXT,
  approver_id INTEGER,
  approved_at TIMESTAMP WITH TIME ZONE,
  status VARCHAR(20) DEFAULT 'PENDING',
  CONSTRAINT fk_leave_requests_employee FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE CASCADE,
  CONSTRAINT fk_leave_requests_leave_type FOREIGN KEY (leave_type_id) REFERENCES leave_types(id) ON DELETE CASCADE,
  CONSTRAINT fk_leave_requests_approver FOREIGN KEY (approver_id) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_leave_requests_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_leave_requests_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 绩效评估表
CREATE TABLE IF NOT EXISTS performance_reviews (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  created_by INTEGER,
  updated_by INTEGER,
  is_active BOOLEAN DEFAULT TRUE,
  employee_id INTEGER NOT NULL,
  reviewer_id INTEGER NOT NULL,
  review_period VARCHAR(7) NOT NULL, -- YYYY-MM
  goals TEXT,
  achievements TEXT,
  strengths TEXT,
  areas_for_improvement TEXT,
  overall_rating INTEGER CHECK (overall_rating >= 1 AND overall_rating <= 5),
  comments TEXT,
  status VARCHAR(20) DEFAULT 'DRAFT',
  CONSTRAINT fk_performance_reviews_employee FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE CASCADE,
  CONSTRAINT fk_performance_reviews_reviewer FOREIGN KEY (reviewer_id) REFERENCES users(id) ON DELETE CASCADE,
  CONSTRAINT fk_performance_reviews_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_performance_reviews_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 培训课程表
CREATE TABLE IF NOT EXISTS training_courses (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  created_by INTEGER,
  updated_by INTEGER,
  is_active BOOLEAN DEFAULT TRUE,
  code VARCHAR(50) NOT NULL,
  name VARCHAR(200) NOT NULL,
  description TEXT,
  instructor VARCHAR(100),
  duration_hours INTEGER DEFAULT 0,
  max_participants INTEGER DEFAULT 0,
  cost DECIMAL(15,2) DEFAULT 0,
  currency VARCHAR(3) DEFAULT 'CNY',
  status VARCHAR(20) DEFAULT 'ACTIVE',
  CONSTRAINT uq_training_courses_code UNIQUE (code),
  CONSTRAINT fk_training_courses_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_training_courses_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 培训记录表
CREATE TABLE IF NOT EXISTS training_records (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  created_by INTEGER,
  updated_by INTEGER,
  is_active BOOLEAN DEFAULT TRUE,
  employee_id INTEGER NOT NULL,
  training_course_id INTEGER NOT NULL,
  start_date DATE NOT NULL,
  end_date DATE,
  completion_status VARCHAR(20) DEFAULT 'ENROLLED',
  score DECIMAL(5,2),
  certificate_number VARCHAR(100),
  notes TEXT,
  CONSTRAINT fk_training_records_employee FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE CASCADE,
  CONSTRAINT fk_training_records_training_course FOREIGN KEY (training_course_id) REFERENCES training_courses(id) ON DELETE CASCADE,
  CONSTRAINT fk_training_records_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_training_records_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- ============================================================================
-- 3. 销售管理模块 (11个表)
-- ============================================================================

-- 客户表
CREATE TABLE IF NOT EXISTS customers (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  created_by INTEGER,
  updated_by INTEGER,
  is_active BOOLEAN DEFAULT TRUE,
  code VARCHAR(50) NOT NULL,
  name VARCHAR(200) NOT NULL,
  description TEXT,
  customer_type VARCHAR(20) DEFAULT 'COMPANY',
  industry VARCHAR(100),
  address TEXT,
  phone VARCHAR(50),
  email VARCHAR(100),
  website VARCHAR(200),
  tax_number VARCHAR(50),
  credit_limit DECIMAL(15,2) DEFAULT 0,
  payment_terms VARCHAR(50),
  sales_rep_id INTEGER,
  status VARCHAR(20) DEFAULT 'ACTIVE',
  CONSTRAINT uq_customers_code UNIQUE (code),
  CONSTRAINT fk_customers_sales_rep FOREIGN KEY (sales_rep_id) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_customers_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_customers_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 客户联系人表
CREATE TABLE IF NOT EXISTS customer_contacts (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  created_by INTEGER,
  updated_by INTEGER,
  is_active BOOLEAN DEFAULT TRUE,
  customer_id INTEGER NOT NULL,
  name VARCHAR(100) NOT NULL,
  title VARCHAR(100),
  phone VARCHAR(50),
  email VARCHAR(100),
  is_primary BOOLEAN DEFAULT FALSE,
  notes TEXT,
  CONSTRAINT fk_customer_contacts_customer FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE,
  CONSTRAINT fk_customer_contacts_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_customer_contacts_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 销售机会表
CREATE TABLE IF NOT EXISTS opportunities (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  created_by INTEGER,
  updated_by INTEGER,
  is_active BOOLEAN DEFAULT TRUE,
  code VARCHAR(50) NOT NULL,
  name VARCHAR(200) NOT NULL,
  description TEXT,
  customer_id INTEGER NOT NULL,
  sales_rep_id INTEGER,
  estimated_value DECIMAL(15,2) DEFAULT 0,
  probability INTEGER DEFAULT 0 CHECK (probability >= 0 AND probability <= 100),
  expected_close_date DATE,
  actual_close_date DATE,
  stage VARCHAR(50) DEFAULT 'PROSPECTING',
  source VARCHAR(50),
  status VARCHAR(20) DEFAULT 'OPEN',
  CONSTRAINT uq_opportunities_code UNIQUE (code),
  CONSTRAINT fk_opportunities_customer FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE,
  CONSTRAINT fk_opportunities_sales_rep FOREIGN KEY (sales_rep_id) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_opportunities_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_opportunities_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 报价单表
CREATE TABLE IF NOT EXISTS quotations (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  created_by INTEGER,
  updated_by INTEGER,
  is_active BOOLEAN DEFAULT TRUE,
  code VARCHAR(50) NOT NULL,
  customer_id INTEGER NOT NULL,
  opportunity_id INTEGER,
  quotation_date DATE NOT NULL,
  valid_until DATE,
  subtotal DECIMAL(15,2) DEFAULT 0,
  tax_amount DECIMAL(15,2) DEFAULT 0,
  total_amount DECIMAL(15,2) DEFAULT 0,
  currency VARCHAR(3) DEFAULT 'CNY',
  payment_terms VARCHAR(100),
  delivery_terms VARCHAR(100),
  notes TEXT,
  status VARCHAR(20) DEFAULT 'DRAFT',
  CONSTRAINT uq_purchase_orders_code UNIQUE (code),
  CONSTRAINT fk_purchase_orders_supplier FOREIGN KEY (supplier_id) REFERENCES suppliers(id) ON DELETE CASCADE,
  CONSTRAINT fk_purchase_orders_supplier_quote FOREIGN KEY (supplier_quote_id) REFERENCES supplier_quotes(id) ON DELETE SET NULL,
  CONSTRAINT fk_purchase_orders_buyer FOREIGN KEY (buyer_id) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_purchase_orders_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_purchase_orders_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 采购订单明细表
CREATE TABLE IF NOT EXISTS purchase_order_items (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  purchase_order_id INTEGER NOT NULL,
  line_number INTEGER DEFAULT 0,
  item_id INTEGER,
  item_code VARCHAR(50),
  item_name VARCHAR(200) NOT NULL,
  description TEXT,
  quantity DECIMAL(15,4) NOT NULL,
  received_quantity DECIMAL(15,4) DEFAULT 0,
  unit VARCHAR(20),
  unit_price DECIMAL(15,2) NOT NULL,
  discount_percent DECIMAL(5,2) DEFAULT 0,
  discount_amount DECIMAL(15,2) DEFAULT 0,
  line_total DECIMAL(15,2) DEFAULT 0,
  required_date DATE,
  CONSTRAINT fk_purchase_order_items_purchase_order FOREIGN KEY (purchase_order_id) REFERENCES purchase_orders(id) ON DELETE CASCADE,
  CONSTRAINT fk_purchase_order_items_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE SET NULL
);

-- 收货单表
CREATE TABLE IF NOT EXISTS receipts (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  created_by INTEGER,
  updated_by INTEGER,
  is_active BOOLEAN DEFAULT TRUE,
  code VARCHAR(50) NOT NULL,
  purchase_order_id INTEGER NOT NULL,
  receipt_date DATE NOT NULL,
  warehouse_id INTEGER NOT NULL,
  supplier_delivery_note VARCHAR(100),
  carrier VARCHAR(100),
  notes TEXT,
  status VARCHAR(20) DEFAULT 'DRAFT',
  CONSTRAINT uq_receipts_code UNIQUE (code),
  CONSTRAINT fk_receipts_purchase_order FOREIGN KEY (purchase_order_id) REFERENCES purchase_orders(id) ON DELETE CASCADE,
  CONSTRAINT fk_receipts_warehouse FOREIGN KEY (warehouse_id) REFERENCES warehouses(id) ON DELETE CASCADE,
  CONSTRAINT fk_receipts_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_receipts_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 收货单明细表
CREATE TABLE IF NOT EXISTS receipt_items (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  receipt_id INTEGER NOT NULL,
  purchase_order_item_id INTEGER NOT NULL,
  quantity DECIMAL(15,4) NOT NULL,
  unit VARCHAR(20),
  unit_cost DECIMAL(15,2) DEFAULT 0,
  total_cost DECIMAL(15,2) DEFAULT 0,
  location_id INTEGER,
  batch_number VARCHAR(100),
  expiry_date DATE,
  notes TEXT,
  CONSTRAINT fk_receipt_items_receipt FOREIGN KEY (receipt_id) REFERENCES receipts(id) ON DELETE CASCADE,
  CONSTRAINT fk_receipt_items_purchase_order_item FOREIGN KEY (purchase_order_item_id) REFERENCES purchase_order_items(id) ON DELETE CASCADE,
  CONSTRAINT fk_receipt_items_location FOREIGN KEY (location_id) REFERENCES locations(id) ON DELETE SET NULL
);

-- 采购发票表
CREATE TABLE IF NOT EXISTS purchase_invoices (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  created_by INTEGER,
  updated_by INTEGER,
  is_active BOOLEAN DEFAULT TRUE,
  code VARCHAR(50) NOT NULL,
  supplier_id INTEGER NOT NULL,
  purchase_order_id INTEGER,
  invoice_date DATE NOT NULL,
  due_date DATE,
  supplier_invoice_number VARCHAR(100),
  subtotal DECIMAL(15,2) DEFAULT 0,
  tax_amount DECIMAL(15,2) DEFAULT 0,
  total_amount DECIMAL(15,2) DEFAULT 0,
  paid_amount DECIMAL(15,2) DEFAULT 0,
  currency VARCHAR(3) DEFAULT 'CNY',
  payment_terms VARCHAR(100),
  notes TEXT,
  status VARCHAR(20) DEFAULT 'DRAFT',
  CONSTRAINT uq_purchase_invoices_code UNIQUE (code),
  CONSTRAINT fk_purchase_invoices_supplier FOREIGN KEY (supplier_id) REFERENCES suppliers(id) ON DELETE CASCADE,
  CONSTRAINT fk_purchase_invoices_purchase_order FOREIGN KEY (purchase_order_id) REFERENCES purchase_orders(id) ON DELETE SET NULL,
  CONSTRAINT fk_purchase_invoices_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_purchase_invoices_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 采购发票明细表
CREATE TABLE IF NOT EXISTS purchase_invoice_items (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  purchase_invoice_id INTEGER NOT NULL,
  line_number INTEGER DEFAULT 0,
  item_id INTEGER,
  item_code VARCHAR(50),
  item_name VARCHAR(200) NOT NULL,
  description TEXT,
  quantity DECIMAL(15,4) NOT NULL,
  unit VARCHAR(20),
  unit_price DECIMAL(15,2) NOT NULL,
  discount_percent DECIMAL(5,2) DEFAULT 0,
  discount_amount DECIMAL(15,2) DEFAULT 0,
  line_total DECIMAL(15,2) DEFAULT 0,
  CONSTRAINT fk_purchase_invoice_items_purchase_invoice FOREIGN KEY (purchase_invoice_id) REFERENCES purchase_invoices(id) ON DELETE CASCADE,
  CONSTRAINT fk_purchase_invoice_items_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE SET NULL
);

-- ============================================================================
-- 6. 财务会计模块 (11个表)
-- ============================================================================

-- 会计科目表
CREATE TABLE IF NOT EXISTS accounts (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  created_by INTEGER,
  updated_by INTEGER,
  is_active BOOLEAN DEFAULT TRUE,
  code VARCHAR(50) NOT NULL,
  name VARCHAR(200) NOT NULL,
  description TEXT,
  account_type VARCHAR(20) NOT NULL,
  parent_id INTEGER,
  level INTEGER DEFAULT 1,
  is_leaf BOOLEAN DEFAULT TRUE,
  normal_balance VARCHAR(10) DEFAULT 'DEBIT',
  currency VARCHAR(3) DEFAULT 'CNY',
  CONSTRAINT uq_accounts_code UNIQUE (code),
  CONSTRAINT fk_accounts_parent FOREIGN KEY (parent_id) REFERENCES accounts(id) ON DELETE SET NULL,
  CONSTRAINT fk_accounts_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_accounts_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 会计期间表
CREATE TABLE IF NOT EXISTS fiscal_periods (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  created_by INTEGER,
  updated_by INTEGER,
  is_active BOOLEAN DEFAULT TRUE,
  code VARCHAR(20) NOT NULL,
  name VARCHAR(100) NOT NULL,
  start_date DATE NOT NULL,
  end_date DATE NOT NULL,
  fiscal_year INTEGER NOT NULL,
  period_number INTEGER NOT NULL,
  is_closed BOOLEAN DEFAULT FALSE,
  closed_by INTEGER,
  closed_at TIMESTAMP WITH TIME ZONE,
  CONSTRAINT uq_fiscal_periods_code UNIQUE (code),
  CONSTRAINT fk_fiscal_periods_closed_by FOREIGN KEY (closed_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_fiscal_periods_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_fiscal_periods_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 会计凭证表
CREATE TABLE IF NOT EXISTS journal_entries (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  created_by INTEGER,
  updated_by INTEGER,
  is_active BOOLEAN DEFAULT TRUE,
  code VARCHAR(50) NOT NULL,
  entry_date DATE NOT NULL,
  fiscal_period_id INTEGER NOT NULL,
  reference_type VARCHAR(50),
  reference_id INTEGER,
  description TEXT,
  total_debit DECIMAL(15,2) DEFAULT 0,
  total_credit DECIMAL(15,2) DEFAULT 0,
  currency VARCHAR(3) DEFAULT 'CNY',
  posted_by INTEGER,
  posted_at TIMESTAMP WITH TIME ZONE,
  status VARCHAR(20) DEFAULT 'DRAFT',
  CONSTRAINT uq_journal_entries_code UNIQUE (code),
  CONSTRAINT fk_journal_entries_fiscal_period FOREIGN KEY (fiscal_period_id) REFERENCES fiscal_periods(id) ON DELETE CASCADE,
  CONSTRAINT fk_journal_entries_posted_by FOREIGN KEY (posted_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_journal_entries_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_journal_entries_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 会计凭证明细表
CREATE TABLE IF NOT EXISTS journal_entry_lines (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  journal_entry_id INTEGER NOT NULL,
  line_number INTEGER DEFAULT 0,
  account_id INTEGER NOT NULL,
  description TEXT,
  debit_amount DECIMAL(15,2) DEFAULT 0,
  credit_amount DECIMAL(15,2) DEFAULT 0,
  currency VARCHAR(3) DEFAULT 'CNY',
  CONSTRAINT fk_journal_entry_lines_journal_entry FOREIGN KEY (journal_entry_id) REFERENCES journal_entries(id) ON DELETE CASCADE,
  CONSTRAINT fk_journal_entry_lines_account FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE
);

-- 总账余额表
CREATE TABLE IF NOT EXISTS general_ledger_balances (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  account_id INTEGER NOT NULL,
  fiscal_period_id INTEGER NOT NULL,
  opening_balance DECIMAL(15,2) DEFAULT 0,
  debit_amount DECIMAL(15,2) DEFAULT 0,
  credit_amount DECIMAL(15,2) DEFAULT 0,
  closing_balance DECIMAL(15,2) DEFAULT 0,
  currency VARCHAR(3) DEFAULT 'CNY',
  CONSTRAINT uq_general_ledger_balances_account_period UNIQUE (account_id, fiscal_period_id),
  CONSTRAINT fk_general_ledger_balances_account FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE,
  CONSTRAINT fk_general_ledger_balances_fiscal_period FOREIGN KEY (fiscal_period_id) REFERENCES fiscal_periods(id) ON DELETE CASCADE
);

-- 应收账款表
CREATE TABLE IF NOT EXISTS accounts_receivable (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  created_by INTEGER,
  updated_by INTEGER,
  is_active BOOLEAN DEFAULT TRUE,
  customer_id INTEGER NOT NULL,
  invoice_id INTEGER,
  invoice_number VARCHAR(50),
  invoice_date DATE NOT NULL,
  due_date DATE,
  original_amount DECIMAL(15,2) NOT NULL,
  paid_amount DECIMAL(15,2) DEFAULT 0,
  outstanding_amount DECIMAL(15,2) DEFAULT 0,
  currency VARCHAR(3) DEFAULT 'CNY',
  aging_days INTEGER DEFAULT 0,
  status VARCHAR(20) DEFAULT 'OPEN',
  CONSTRAINT fk_accounts_receivable_customer FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE,
  CONSTRAINT fk_accounts_receivable_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_accounts_receivable_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 应付账款表
CREATE TABLE IF NOT EXISTS accounts_payable (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  created_by INTEGER,
  updated_by INTEGER,
  is_active BOOLEAN DEFAULT TRUE,
  supplier_id INTEGER NOT NULL,
  invoice_id INTEGER,
  invoice_number VARCHAR(50),
  invoice_date DATE NOT NULL,
  due_date DATE,
  original_amount DECIMAL(15,2) NOT NULL,
  paid_amount DECIMAL(15,2) DEFAULT 0,
  outstanding_amount DECIMAL(15,2) DEFAULT 0,
  currency VARCHAR(3) DEFAULT 'CNY',
  aging_days INTEGER DEFAULT 0,
  status VARCHAR(20) DEFAULT 'OPEN',
  CONSTRAINT fk_accounts_payable_supplier FOREIGN KEY (supplier_id) REFERENCES suppliers(id) ON DELETE CASCADE,
  CONSTRAINT fk_accounts_payable_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_accounts_payable_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 付款单表
CREATE TABLE IF NOT EXISTS payments (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  created_by INTEGER,
  updated_by INTEGER,
  is_active BOOLEAN DEFAULT TRUE,
  code VARCHAR(50) NOT NULL,
  payment_date DATE NOT NULL,
  supplier_id INTEGER,
  customer_id INTEGER,
  payment_method VARCHAR(20) DEFAULT 'BANK_TRANSFER',
  bank_account_id INTEGER,
  total_amount DECIMAL(15,2) NOT NULL,
  currency VARCHAR(3) DEFAULT 'CNY',
  reference_number VARCHAR(100),
  notes TEXT,
  status VARCHAR(20) DEFAULT 'DRAFT',
  CONSTRAINT uq_payments_code UNIQUE (code),
  CONSTRAINT fk_payments_supplier FOREIGN KEY (supplier_id) REFERENCES suppliers(id) ON DELETE SET NULL,
  CONSTRAINT fk_payments_customer FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE SET NULL,
  CONSTRAINT fk_payments_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_payments_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 付款单明细表
CREATE TABLE IF NOT EXISTS payment_items (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  payment_id INTEGER NOT NULL,
  line_number INTEGER DEFAULT 0,
  invoice_id INTEGER,
  invoice_number VARCHAR(50),
  original_amount DECIMAL(15,2) DEFAULT 0,
  payment_amount DECIMAL(15,2) NOT NULL,
  discount_amount DECIMAL(15,2) DEFAULT 0,
  currency VARCHAR(3) DEFAULT 'CNY',
  notes TEXT,
  CONSTRAINT fk_payment_items_payment FOREIGN KEY (payment_id) REFERENCES payments(id) ON DELETE CASCADE
);

-- 收款单表
CREATE TABLE IF NOT EXISTS receipts_ar (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  created_by INTEGER,
  updated_by INTEGER,
  is_active BOOLEAN DEFAULT TRUE,
  code VARCHAR(50) NOT NULL,
  receipt_date DATE NOT NULL,
  customer_id INTEGER NOT NULL,
  payment_method VARCHAR(20) DEFAULT 'BANK_TRANSFER',
  bank_account_id INTEGER,
  total_amount DECIMAL(15,2) NOT NULL,
  currency VARCHAR(3) DEFAULT 'CNY',
  reference_number VARCHAR(100),
  notes TEXT,
  status VARCHAR(20) DEFAULT 'DRAFT',
  CONSTRAINT uq_receipts_ar_code UNIQUE (code),
  CONSTRAINT fk_receipts_ar_customer FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE,
  CONSTRAINT fk_receipts_ar_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_receipts_ar_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 收款单明细表
CREATE TABLE IF NOT EXISTS receipt_ar_items (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  receipt_ar_id INTEGER NOT NULL,
  line_number INTEGER DEFAULT 0,
  invoice_id INTEGER,
  invoice_number VARCHAR(50),
  original_amount DECIMAL(15,2) DEFAULT 0,
  receipt_amount DECIMAL(15,2) NOT NULL,
  discount_amount DECIMAL(15,2) DEFAULT 0,
  currency VARCHAR(3) DEFAULT 'CNY',
  notes TEXT,
  CONSTRAINT fk_receipt_ar_items_receipt_ar FOREIGN KEY (receipt_ar_id) REFERENCES receipts_ar(id) ON DELETE CASCADE
);

-- 银行账户表
CREATE TABLE IF NOT EXISTS bank_accounts (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  created_by INTEGER,
  updated_by INTEGER,
  is_active BOOLEAN DEFAULT TRUE,
  code VARCHAR(50) NOT NULL,
  account_name VARCHAR(200) NOT NULL,
  account_number VARCHAR(100) NOT NULL,
  bank_name VARCHAR(200) NOT NULL,
  bank_branch VARCHAR(200),
  swift_code VARCHAR(20),
  iban VARCHAR(50),
  currency VARCHAR(3) DEFAULT 'CNY',
  account_type VARCHAR(20) DEFAULT 'CHECKING',
  opening_balance DECIMAL(15,2) DEFAULT 0,
  current_balance DECIMAL(15,2) DEFAULT 0,
  status VARCHAR(20) DEFAULT 'ACTIVE',
  CONSTRAINT uq_bank_accounts_code UNIQUE (code),
  CONSTRAINT uq_bank_accounts_account_number UNIQUE (account_number),
  CONSTRAINT fk_bank_accounts_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_bank_accounts_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- ================================
-- 生产管理模块表
-- ================================

-- 物料清单表
CREATE TABLE IF NOT EXISTS boms (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  created_by INTEGER,
  updated_by INTEGER,
  is_active BOOLEAN DEFAULT TRUE,
  code VARCHAR(50) NOT NULL,
  name VARCHAR(200) NOT NULL,
  item_id INTEGER NOT NULL,
  version VARCHAR(20) DEFAULT '1.0',
  quantity DECIMAL(15,4) DEFAULT 1,
  unit_id INTEGER,
  bom_type VARCHAR(20) DEFAULT 'PRODUCTION',
  status VARCHAR(20) DEFAULT 'DRAFT',
  effective_date DATE,
  expiry_date DATE,
  notes TEXT,
  CONSTRAINT uq_boms_code UNIQUE (code),
  CONSTRAINT fk_boms_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE,
  CONSTRAINT fk_boms_unit FOREIGN KEY (unit_id) REFERENCES units(id) ON DELETE SET NULL,
  CONSTRAINT fk_boms_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_boms_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 物料清单明细表
CREATE TABLE IF NOT EXISTS bom_items (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  bom_id INTEGER NOT NULL,
  line_number INTEGER DEFAULT 0,
  item_id INTEGER NOT NULL,
  quantity DECIMAL(15,4) NOT NULL,
  unit_id INTEGER,
  scrap_rate DECIMAL(5,2) DEFAULT 0,
  component_type VARCHAR(20) DEFAULT 'MATERIAL',
  is_phantom BOOLEAN DEFAULT FALSE,
  notes TEXT,
  CONSTRAINT fk_bom_items_bom FOREIGN KEY (bom_id) REFERENCES boms(id) ON DELETE CASCADE,
  CONSTRAINT fk_bom_items_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE,
  CONSTRAINT fk_bom_items_unit FOREIGN KEY (unit_id) REFERENCES units(id) ON DELETE SET NULL
);

-- 工作中心表
CREATE TABLE IF NOT EXISTS work_centers (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  created_by INTEGER,
  updated_by INTEGER,
  is_active BOOLEAN DEFAULT TRUE,
  code VARCHAR(50) NOT NULL,
  name VARCHAR(200) NOT NULL,
  description TEXT,
  work_center_type VARCHAR(20) DEFAULT 'MACHINE',
  capacity DECIMAL(10,2) DEFAULT 1,
  efficiency DECIMAL(5,2) DEFAULT 100,
  cost_per_hour DECIMAL(10,2) DEFAULT 0,
  setup_time INTEGER DEFAULT 0,
  queue_time INTEGER DEFAULT 0,
  move_time INTEGER DEFAULT 0,
  status VARCHAR(20) DEFAULT 'ACTIVE',
  CONSTRAINT uq_work_centers_code UNIQUE (code),
  CONSTRAINT fk_work_centers_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_work_centers_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 工艺路线表
CREATE TABLE IF NOT EXISTS routings (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  created_by INTEGER,
  updated_by INTEGER,
  is_active BOOLEAN DEFAULT TRUE,
  code VARCHAR(50) NOT NULL,
  name VARCHAR(200) NOT NULL,
  item_id INTEGER NOT NULL,
  version VARCHAR(20) DEFAULT '1.0',
  routing_type VARCHAR(20) DEFAULT 'PRODUCTION',
  status VARCHAR(20) DEFAULT 'DRAFT',
  effective_date DATE,
  expiry_date DATE,
  notes TEXT,
  CONSTRAINT uq_routings_code UNIQUE (code),
  CONSTRAINT fk_routings_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE,
  CONSTRAINT fk_routings_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_routings_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 工艺路线操作表
CREATE TABLE IF NOT EXISTS routing_operations (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  routing_id INTEGER NOT NULL,
  operation_number INTEGER NOT NULL,
  operation_name VARCHAR(200) NOT NULL,
  work_center_id INTEGER NOT NULL,
  setup_time INTEGER DEFAULT 0,
  run_time INTEGER DEFAULT 0,
  wait_time INTEGER DEFAULT 0,
  move_time INTEGER DEFAULT 0,
  operation_type VARCHAR(20) DEFAULT 'PRODUCTION',
  description TEXT,
  CONSTRAINT fk_routing_operations_routing FOREIGN KEY (routing_id) REFERENCES routings(id) ON DELETE CASCADE,
  CONSTRAINT fk_routing_operations_work_center FOREIGN KEY (work_center_id) REFERENCES work_centers(id) ON DELETE CASCADE
);

-- 生产计划表
CREATE TABLE IF NOT EXISTS production_plans (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  created_by INTEGER,
  updated_by INTEGER,
  is_active BOOLEAN DEFAULT TRUE,
  code VARCHAR(50) NOT NULL,
  name VARCHAR(200) NOT NULL,
  plan_date DATE NOT NULL,
  start_date DATE,
  end_date DATE,
  status VARCHAR(20) DEFAULT 'DRAFT',
  notes TEXT,
  CONSTRAINT uq_production_plans_code UNIQUE (code),
  CONSTRAINT fk_production_plans_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_production_plans_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 生产计划明细表
CREATE TABLE IF NOT EXISTS production_plan_items (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  production_plan_id INTEGER NOT NULL,
  line_number INTEGER DEFAULT 0,
  item_id INTEGER NOT NULL,
  planned_quantity DECIMAL(15,4) NOT NULL,
  unit_id INTEGER,
  required_date DATE,
  priority INTEGER DEFAULT 0,
  notes TEXT,
  CONSTRAINT fk_production_plan_items_plan FOREIGN KEY (production_plan_id) REFERENCES production_plans(id) ON DELETE CASCADE,
  CONSTRAINT fk_production_plan_items_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE,
  CONSTRAINT fk_production_plan_items_unit FOREIGN KEY (unit_id) REFERENCES units(id) ON DELETE SET NULL
);

-- 生产订单表
CREATE TABLE IF NOT EXISTS production_orders (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  created_by INTEGER,
  updated_by INTEGER,
  is_active BOOLEAN DEFAULT TRUE,
  code VARCHAR(50) NOT NULL,
  item_id INTEGER NOT NULL,
  bom_id INTEGER,
  routing_id INTEGER,
  planned_quantity DECIMAL(15,4) NOT NULL,
  produced_quantity DECIMAL(15,4) DEFAULT 0,
  unit_id INTEGER,
  planned_start_date DATE,
  planned_end_date DATE,
  actual_start_date DATE,
  actual_end_date DATE,
  priority INTEGER DEFAULT 0,
  status VARCHAR(20) DEFAULT 'PLANNED',
  notes TEXT,
  CONSTRAINT uq_production_orders_code UNIQUE (code),
  CONSTRAINT fk_production_orders_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE,
  CONSTRAINT fk_production_orders_bom FOREIGN KEY (bom_id) REFERENCES boms(id) ON DELETE SET NULL,
  CONSTRAINT fk_production_orders_routing FOREIGN KEY (routing_id) REFERENCES routings(id) ON DELETE SET NULL,
  CONSTRAINT fk_production_orders_unit FOREIGN KEY (unit_id) REFERENCES units(id) ON DELETE SET NULL,
  CONSTRAINT fk_production_orders_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_production_orders_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 生产订单物料需求表
CREATE TABLE IF NOT EXISTS production_order_materials (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  production_order_id INTEGER NOT NULL,
  item_id INTEGER NOT NULL,
  required_quantity DECIMAL(15,4) NOT NULL,
  issued_quantity DECIMAL(15,4) DEFAULT 0,
  unit_id INTEGER,
  required_date DATE,
  status VARCHAR(20) DEFAULT 'PLANNED',
  CONSTRAINT fk_production_order_materials_order FOREIGN KEY (production_order_id) REFERENCES production_orders(id) ON DELETE CASCADE,
  CONSTRAINT fk_production_order_materials_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE,
  CONSTRAINT fk_production_order_materials_unit FOREIGN KEY (unit_id) REFERENCES units(id) ON DELETE SET NULL
);

-- 生产订单操作表
CREATE TABLE IF NOT EXISTS production_order_operations (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  production_order_id INTEGER NOT NULL,
  operation_number INTEGER NOT NULL,
  operation_name VARCHAR(200) NOT NULL,
  work_center_id INTEGER NOT NULL,
  planned_start_date TIMESTAMP WITH TIME ZONE,
  planned_end_date TIMESTAMP WITH TIME ZONE,
  actual_start_date TIMESTAMP WITH TIME ZONE,
  actual_end_date TIMESTAMP WITH TIME ZONE,
  setup_time INTEGER DEFAULT 0,
  run_time INTEGER DEFAULT 0,
  status VARCHAR(20) DEFAULT 'PLANNED',
  CONSTRAINT fk_production_order_operations_order FOREIGN KEY (production_order_id) REFERENCES production_orders(id) ON DELETE CASCADE,
  CONSTRAINT fk_production_order_operations_work_center FOREIGN KEY (work_center_id) REFERENCES work_centers(id) ON DELETE CASCADE
);

-- 生产报工表
CREATE TABLE IF NOT EXISTS work_orders (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  created_by INTEGER,
  updated_by INTEGER,
  is_active BOOLEAN DEFAULT TRUE,
  code VARCHAR(50) NOT NULL,
  production_order_id INTEGER NOT NULL,
  operation_id INTEGER,
  work_center_id INTEGER NOT NULL,
  employee_id INTEGER,
  work_date DATE NOT NULL,
  start_time TIMESTAMP WITH TIME ZONE,
  end_time TIMESTAMP WITH TIME ZONE,
  quantity_completed DECIMAL(15,4) DEFAULT 0,
  quantity_scrapped DECIMAL(15,4) DEFAULT 0,
  unit_id INTEGER,
  status VARCHAR(20) DEFAULT 'IN_PROGRESS',
  notes TEXT,
  CONSTRAINT uq_work_orders_code UNIQUE (code),
  CONSTRAINT fk_work_orders_production_order FOREIGN KEY (production_order_id) REFERENCES production_orders(id) ON DELETE CASCADE,
  CONSTRAINT fk_work_orders_operation FOREIGN KEY (operation_id) REFERENCES production_order_operations(id) ON DELETE SET NULL,
  CONSTRAINT fk_work_orders_work_center FOREIGN KEY (work_center_id) REFERENCES work_centers(id) ON DELETE CASCADE,
  CONSTRAINT fk_work_orders_employee FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE SET NULL,
  CONSTRAINT fk_work_orders_unit FOREIGN KEY (unit_id) REFERENCES units(id) ON DELETE SET NULL,
  CONSTRAINT fk_work_orders_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_work_orders_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 质量检验表
CREATE TABLE IF NOT EXISTS quality_inspections (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  created_by INTEGER,
  updated_by INTEGER,
  is_active BOOLEAN DEFAULT TRUE,
  code VARCHAR(50) NOT NULL,
  inspection_type VARCHAR(20) DEFAULT 'INCOMING',
  item_id INTEGER NOT NULL,
  batch_number VARCHAR(100),
  quantity_inspected DECIMAL(15,4) NOT NULL,
  quantity_passed DECIMAL(15,4) DEFAULT 0,
  quantity_failed DECIMAL(15,4) DEFAULT 0,
  unit_id INTEGER,
  inspection_date DATE NOT NULL,
  inspector_id INTEGER,
  supplier_id INTEGER,
  production_order_id INTEGER,
  status VARCHAR(20) DEFAULT 'PENDING',
  result VARCHAR(20) DEFAULT 'PENDING',
  notes TEXT,
  CONSTRAINT uq_quality_inspections_code UNIQUE (code),
  CONSTRAINT fk_quality_inspections_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE,
  CONSTRAINT fk_quality_inspections_unit FOREIGN KEY (unit_id) REFERENCES units(id) ON DELETE SET NULL,
  CONSTRAINT fk_quality_inspections_inspector FOREIGN KEY (inspector_id) REFERENCES employees(id) ON DELETE SET NULL,
  CONSTRAINT fk_quality_inspections_supplier FOREIGN KEY (supplier_id) REFERENCES suppliers(id) ON DELETE SET NULL,
  CONSTRAINT fk_quality_inspections_production_order FOREIGN KEY (production_order_id) REFERENCES production_orders(id) ON DELETE SET NULL,
  CONSTRAINT fk_quality_inspections_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_quality_inspections_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- ================================
-- 项目管理模块表
-- ================================

-- 项目表
CREATE TABLE IF NOT EXISTS projects (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  created_by INTEGER,
  updated_by INTEGER,
  is_active BOOLEAN DEFAULT TRUE,
  code VARCHAR(50) NOT NULL,
  name VARCHAR(200) NOT NULL,
  description TEXT,
  customer_id INTEGER,
  project_manager_id INTEGER,
  start_date DATE,
  end_date DATE,
  budget DECIMAL(15,2) DEFAULT 0,
  currency VARCHAR(3) DEFAULT 'CNY',
  status VARCHAR(20) DEFAULT 'PLANNING',
  priority INTEGER DEFAULT 0,
  CONSTRAINT uq_projects_code UNIQUE (code),
  CONSTRAINT fk_projects_customer FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE SET NULL,
  CONSTRAINT fk_projects_manager FOREIGN KEY (project_manager_id) REFERENCES employees(id) ON DELETE SET NULL,
  CONSTRAINT fk_projects_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_projects_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 项目任务表
CREATE TABLE IF NOT EXISTS project_tasks (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  created_by INTEGER,
  updated_by INTEGER,
  is_active BOOLEAN DEFAULT TRUE,
  project_id INTEGER NOT NULL,
  parent_task_id INTEGER,
  task_name VARCHAR(200) NOT NULL,
  description TEXT,
  assigned_to INTEGER,
  start_date DATE,
  end_date DATE,
  estimated_hours DECIMAL(8,2) DEFAULT 0,
  actual_hours DECIMAL(8,2) DEFAULT 0,
  progress DECIMAL(5,2) DEFAULT 0,
  priority INTEGER DEFAULT 0,
  status VARCHAR(20) DEFAULT 'NOT_STARTED',
  CONSTRAINT fk_project_tasks_project FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
  CONSTRAINT fk_project_tasks_parent FOREIGN KEY (parent_task_id) REFERENCES project_tasks(id) ON DELETE SET NULL,
  CONSTRAINT fk_project_tasks_assigned_to FOREIGN KEY (assigned_to) REFERENCES employees(id) ON DELETE SET NULL,
  CONSTRAINT fk_project_tasks_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_project_tasks_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- ================================
-- 系统管理模块表
-- ================================

-- 系统配置表
CREATE TABLE IF NOT EXISTS system_configs (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  created_by INTEGER,
  updated_by INTEGER,
  is_active BOOLEAN DEFAULT TRUE,
  config_key VARCHAR(100) NOT NULL,
  config_value TEXT,
  config_type VARCHAR(20) DEFAULT 'STRING',
  description TEXT,
  is_system BOOLEAN DEFAULT FALSE,
  CONSTRAINT uq_system_configs_key UNIQUE (config_key),
  CONSTRAINT fk_system_configs_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_system_configs_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 审计日志表
CREATE TABLE IF NOT EXISTS audit_logs (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  user_id INTEGER,
  action VARCHAR(50) NOT NULL,
  table_name VARCHAR(100),
  record_id INTEGER,
  old_values JSONB,
  new_values JSONB,
  ip_address INET,
  user_agent TEXT,
  CONSTRAINT fk_audit_logs_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
);

-- 通知模板表
CREATE TABLE IF NOT EXISTS notification_templates (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  created_by INTEGER,
  updated_by INTEGER,
  is_active BOOLEAN DEFAULT TRUE,
  template_code VARCHAR(50) NOT NULL,
  template_name VARCHAR(200) NOT NULL,
  template_type VARCHAR(20) DEFAULT 'EMAIL',
  subject VARCHAR(500),
  content TEXT,
  variables TEXT,
  CONSTRAINT uq_notification_templates_code UNIQUE (template_code),
  CONSTRAINT fk_notification_templates_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_notification_templates_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 通知表
CREATE TABLE IF NOT EXISTS notifications (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  user_id INTEGER NOT NULL,
  notification_type VARCHAR(20) DEFAULT 'SYSTEM',
  title VARCHAR(500) NOT NULL,
  content TEXT,
  is_read BOOLEAN DEFAULT FALSE,
  read_at TIMESTAMP WITH TIME ZONE,
  priority INTEGER DEFAULT 0,
  CONSTRAINT fk_notifications_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- 文件管理表
CREATE TABLE IF NOT EXISTS files (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  created_by INTEGER,
  updated_by INTEGER,
  is_active BOOLEAN DEFAULT TRUE,
  file_name VARCHAR(500) NOT NULL,
  original_name VARCHAR(500) NOT NULL,
  file_path VARCHAR(1000) NOT NULL,
  file_size BIGINT DEFAULT 0,
  file_type VARCHAR(100),
  mime_type VARCHAR(200),
  file_hash VARCHAR(64),
  storage_type VARCHAR(20) DEFAULT 'LOCAL',
  bucket_name VARCHAR(100),
  object_key VARCHAR(500),
  is_public BOOLEAN DEFAULT FALSE,
  CONSTRAINT fk_files_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_files_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 数据字典表
CREATE TABLE IF NOT EXISTS dictionaries (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  created_by INTEGER,
  updated_by INTEGER,
  is_active BOOLEAN DEFAULT TRUE,
  dict_type VARCHAR(50) NOT NULL,
  dict_key VARCHAR(100) NOT NULL,
  dict_value VARCHAR(500) NOT NULL,
  dict_label VARCHAR(200) NOT NULL,
  sort_order INTEGER DEFAULT 0,
  is_system BOOLEAN DEFAULT FALSE,
  description TEXT,
  CONSTRAINT uq_dictionaries_type_key UNIQUE (dict_type, dict_key),
  CONSTRAINT fk_dictionaries_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_dictionaries_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 定时任务表
CREATE TABLE IF NOT EXISTS scheduled_jobs (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  created_by INTEGER,
  updated_by INTEGER,
  is_active BOOLEAN DEFAULT TRUE,
  job_name VARCHAR(200) NOT NULL,
  job_group VARCHAR(100) DEFAULT 'DEFAULT',
  job_class VARCHAR(500) NOT NULL,
  cron_expression VARCHAR(100) NOT NULL,
  description TEXT,
  status VARCHAR(20) DEFAULT 'PAUSED',
  last_run_time TIMESTAMP WITH TIME ZONE,
  next_run_time TIMESTAMP WITH TIME ZONE,
  CONSTRAINT uq_scheduled_jobs_name_group UNIQUE (job_name, job_group),
  CONSTRAINT fk_scheduled_jobs_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_scheduled_jobs_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 任务执行日志表
CREATE TABLE IF NOT EXISTS job_execution_logs (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  job_id INTEGER NOT NULL,
  start_time TIMESTAMP WITH TIME ZONE NOT NULL,
  end_time TIMESTAMP WITH TIME ZONE,
  status VARCHAR(20) DEFAULT 'RUNNING',
  result_message TEXT,
  error_message TEXT,
  execution_time INTEGER DEFAULT 0,
  CONSTRAINT fk_job_execution_logs_job FOREIGN KEY (job_id) REFERENCES scheduled_jobs(id) ON DELETE CASCADE
);

-- 系统监控表
CREATE TABLE IF NOT EXISTS system_monitors (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  monitor_time TIMESTAMP WITH TIME ZONE NOT NULL,
  cpu_usage DECIMAL(5,2) DEFAULT 0,
  memory_usage DECIMAL(5,2) DEFAULT 0,
  disk_usage DECIMAL(5,2) DEFAULT 0,
  network_in BIGINT DEFAULT 0,
  network_out BIGINT DEFAULT 0,
  active_connections INTEGER DEFAULT 0,
  response_time DECIMAL(10,2) DEFAULT 0
);

-- ================================
-- 遗留模块表
-- ================================

-- 数据迁移记录表
CREATE TABLE IF NOT EXISTS data_migrations (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  created_by INTEGER,
  updated_by INTEGER,
  is_active BOOLEAN DEFAULT TRUE,
  migration_name VARCHAR(200) NOT NULL,
  source_system VARCHAR(100) NOT NULL,
  target_table VARCHAR(100) NOT NULL,
  migration_type VARCHAR(20) DEFAULT 'IMPORT',
  start_time TIMESTAMP WITH TIME ZONE,
  end_time TIMESTAMP WITH TIME ZONE,
  total_records INTEGER DEFAULT 0,
  success_records INTEGER DEFAULT 0,
  failed_records INTEGER DEFAULT 0,
  status VARCHAR(20) DEFAULT 'PENDING',
  error_message TEXT,
  CONSTRAINT fk_data_migrations_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_data_migrations_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 数据迁移详细日志表
CREATE TABLE IF NOT EXISTS migration_logs (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  migration_id INTEGER NOT NULL,
  record_id VARCHAR(100),
  operation_type VARCHAR(20) DEFAULT 'INSERT',
  status VARCHAR(20) DEFAULT 'SUCCESS',
  error_message TEXT,
  old_data JSONB,
  new_data JSONB,
  CONSTRAINT fk_migration_logs_migration FOREIGN KEY (migration_id) REFERENCES data_migrations(id) ON DELETE CASCADE
);

-- ================================
-- 初始数据插入
-- ================================

-- 插入默认公司
INSERT INTO companies (id, code, name, tax_number, legal_representative, address, phone, email, website, status) VALUES
(1, 'GALAXY', '银河科技有限公司', '91110000123456789X', '张三', '北京市朝阳区科技园区1号', '010-12345678', 'info@galaxy.com', 'www.galaxy.com', 'ACTIVE')
ON CONFLICT (id) DO NOTHING;

-- 插入默认部门
INSERT INTO departments (id, company_id, code, name, parent_id, manager_id, description, status) VALUES
(1, 1, 'IT', 'IT部门', NULL, NULL, '信息技术部门', 'ACTIVE'),
(2, 1, 'HR', '人力资源部', NULL, NULL, '人力资源部门', 'ACTIVE'),
(3, 1, 'SALES', '销售部', NULL, NULL, '销售部门', 'ACTIVE'),
(4, 1, 'FINANCE', '财务部', NULL, NULL, '财务部门', 'ACTIVE'),
(5, 1, 'PRODUCTION', '生产部', NULL, NULL, '生产部门', 'ACTIVE')
ON CONFLICT (id) DO NOTHING;

-- 插入默认角色
INSERT INTO roles (id, code, name, description, is_system) VALUES
(1, 'ADMIN', '系统管理员', '系统管理员角色', TRUE),
(2, 'USER', '普通用户', '普通用户角色', TRUE),
(3, 'MANAGER', '部门经理', '部门经理角色', FALSE),
(4, 'EMPLOYEE', '员工', '普通员工角色', FALSE)
ON CONFLICT (id) DO NOTHING;

-- 插入默认权限
INSERT INTO permissions (id, code, name, resource, action, description) VALUES
(1, 'USER_READ', '查看用户', 'user', 'read', '查看用户信息'),
(2, 'USER_WRITE', '编辑用户', 'user', 'write', '编辑用户信息'),
(3, 'USER_DELETE', '删除用户', 'user', 'delete', '删除用户'),
(4, 'ROLE_READ', '查看角色', 'role', 'read', '查看角色信息'),
(5, 'ROLE_WRITE', '编辑角色', 'role', 'write', '编辑角色信息'),
(6, 'SYSTEM_CONFIG', '系统配置', 'system', 'config', '系统配置管理')
ON CONFLICT (id) DO NOTHING;

-- 插入角色权限关联
INSERT INTO role_permissions (role_id, permission_id) VALUES
(1, 1), (1, 2), (1, 3), (1, 4), (1, 5), (1, 6),
(2, 1), (2, 4),
(3, 1), (3, 2), (3, 4), (3, 5),
(4, 1), (4, 4)
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- 插入默认用户
INSERT INTO users (id, company_id, department_id, username, email, password_hash, real_name, phone, status, is_admin) VALUES
(1, 1, 1, 'admin', 'admin@galaxy.com', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iKXgwHNDHDES6LKdOaWpjfwIDAWO', '系统管理员', '13800138000', 'ACTIVE', 1),
(2, 1, 2, 'hr_manager', 'hr@galaxy.com', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iKXgwHNDHDES6LKdOaWpjfwIDAWO', '人事经理', '13800138001', 'ACTIVE', 0),
(3, 1, 3, 'sales_manager', 'sales@galaxy.com', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iKXgwHNDHDES6LKdOaWpjfwIDAWO', '销售经理', '13800138002', 'ACTIVE', 0)
ON CONFLICT (id) DO NOTHING;

-- 插入用户角色关联
INSERT INTO user_roles (user_id, role_id, created_at, updated_at, assigned_at, is_active) VALUES
(1, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, true),
(2, 3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, true),
(3, 3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, true)
ON CONFLICT (user_id, role_id) DO NOTHING;

-- 插入默认职位
INSERT INTO positions (id, code, name, department_id, level, description, status) VALUES
(1, 'DEV', '开发工程师', 1, 'JUNIOR', '软件开发工程师', 'ACTIVE'),
(2, 'PM', '项目经理', 1, 'SENIOR', '项目管理', 'ACTIVE'),
(3, 'HR_SPEC', '人事专员', 2, 'JUNIOR', '人力资源专员', 'ACTIVE'),
(4, 'SALES_REP', '销售代表', 3, 'JUNIOR', '销售代表', 'ACTIVE')
ON CONFLICT (id) DO NOTHING;

-- 插入默认薪资等级
INSERT INTO salary_grades (id, grade_code, grade_name, min_salary, max_salary, currency, description) VALUES
(1, 'G1', '初级', 5000.00, 8000.00, 'CNY', '初级员工薪资等级'),
(2, 'G2', '中级', 8000.00, 15000.00, 'CNY', '中级员工薪资等级'),
(3, 'G3', '高级', 15000.00, 25000.00, 'CNY', '高级员工薪资等级'),
(4, 'G4', '专家', 25000.00, 40000.00, 'CNY', '专家级薪资等级')
ON CONFLICT (id) DO NOTHING;

-- 插入默认员工
INSERT INTO employees (id, company_id, department_id, position_id, employee_number, name, gender, birth_date, phone, email, hire_date, salary_grade_id, status) VALUES
(1, 1, 1, 1, 'EMP001', '张三', 'MALE', '1990-01-01', '13800138000', 'zhangsan@galaxy.com', '2023-01-01', 2, 'ACTIVE'),
(2, 1, 2, 3, 'EMP002', '李四', 'FEMALE', '1992-05-15', '13800138001', 'lisi@galaxy.com', '2023-02-01', 2, 'ACTIVE'),
(3, 1, 3, 4, 'EMP003', '王五', 'MALE', '1988-08-20', '13800138002', 'wangwu@galaxy.com', '2023-03-01', 3, 'ACTIVE')
ON CONFLICT (id) DO NOTHING;

-- 插入默认请假类型
INSERT INTO leave_types (id, type_code, type_name, max_days_per_year, is_paid, description) VALUES
(1, 'ANNUAL', '年假', 10, TRUE, '年度带薪假期'),
(2, 'SICK', '病假', 30, TRUE, '病假'),
(3, 'PERSONAL', '事假', 5, FALSE, '个人事假'),
(4, 'MATERNITY', '产假', 98, TRUE, '产假')
ON CONFLICT (id) DO NOTHING;

-- 插入默认计量单位
INSERT INTO units (id, code, name, unit_type, description) VALUES
(1, 'PCS', '个', 'QUANTITY', '计数单位'),
(2, 'KG', '千克', 'WEIGHT', '重量单位'),
(3, 'M', '米', 'LENGTH', '长度单位'),
(4, 'L', '升', 'VOLUME', '体积单位'),
(5, 'SET', '套', 'QUANTITY', '成套单位')
ON CONFLICT (id) DO NOTHING;

-- 插入默认物料分类
INSERT INTO item_categories (id, code, name, parent_id, description) VALUES
(1, 'RAW', '原材料', NULL, '生产用原材料'),
(2, 'SEMI', '半成品', NULL, '半成品物料'),
(3, 'FINISHED', '成品', NULL, '最终产品'),
(4, 'CONSUMABLE', '耗材', NULL, '消耗性物料')
ON CONFLICT (id) DO NOTHING;

-- 插入默认物料
INSERT INTO items (id, code, name, category_id, item_type, unit_id, specification, status) VALUES
(1, 'RAW001', '钢材A', 1, 'RAW_MATERIAL', 2, '规格：10mm*20mm', 'ACTIVE'),
(2, 'SEMI001', '组件A', 2, 'SEMI_FINISHED', 1, '半成品组件', 'ACTIVE'),
(3, 'FIN001', '产品A', 3, 'FINISHED_GOOD', 1, '最终产品A', 'ACTIVE'),
(4, 'CON001', '螺丝', 4, 'CONSUMABLE', 1, 'M6螺丝', 'ACTIVE')
ON CONFLICT (id) DO NOTHING;

-- 插入默认仓库
INSERT INTO warehouses (id, code, name, warehouse_type, address, manager_id, status) VALUES
(1, 'WH001', '主仓库', 'MAIN', '北京市朝阳区仓储区1号', 1, 'ACTIVE'),
(2, 'WH002', '原料仓', 'RAW_MATERIAL', '北京市朝阳区仓储区2号', 1, 'ACTIVE'),
(3, 'WH003', '成品仓', 'FINISHED_GOOD', '北京市朝阳区仓储区3号', 1, 'ACTIVE')
ON CONFLICT (id) DO NOTHING;

-- 插入默认库位
INSERT INTO locations (id, warehouse_id, code, name, location_type, status) VALUES
(1, 1, 'A01', 'A区01位', 'STORAGE', 'ACTIVE'),
(2, 1, 'A02', 'A区02位', 'STORAGE', 'ACTIVE'),
(3, 2, 'B01', 'B区01位', 'STORAGE', 'ACTIVE'),
(4, 3, 'C01', 'C区01位', 'STORAGE', 'ACTIVE')
ON CONFLICT (id) DO NOTHING;

-- 插入默认客户
INSERT INTO customers (id, code, name, customer_type, contact_person, phone, email, address, status) VALUES
(1, 'CUST001', '北京科技公司', 'ENTERPRISE', '张经理', '010-88888888', 'zhang@bjtech.com', '北京市海淀区科技园', 'ACTIVE'),
(2, 'CUST002', '上海贸易公司', 'ENTERPRISE', '李总', '021-99999999', 'li@shtrade.com', '上海市浦东新区贸易区', 'ACTIVE'),
(3, 'CUST003', '个人客户王先生', 'INDIVIDUAL', '王先生', '13900139000', 'wang@email.com', '广州市天河区', 'ACTIVE')
ON CONFLICT (id) DO NOTHING;

-- 插入默认供应商
INSERT INTO suppliers (id, code, name, supplier_type, contact_person, phone, email, address, status) VALUES
(1, 'SUPP001', '钢材供应商', 'MATERIAL', '赵经理', '010-77777777', 'zhao@steel.com', '河北省唐山市钢铁园区', 'ACTIVE'),
(2, 'SUPP002', '电子元件供应商', 'COMPONENT', '钱总', '0755-66666666', 'qian@electronic.com', '深圳市南山区电子城', 'ACTIVE'),
(3, 'SUPP003', '包装材料供应商', 'PACKAGING', '孙经理', '021-55555555', 'sun@package.com', '上海市松江区包装园', 'ACTIVE')
ON CONFLICT (id) DO NOTHING;

-- 插入默认会计科目
INSERT INTO accounts (id, account_code, account_name, account_type, parent_id, is_leaf, description) VALUES
(1, '1001', '库存现金', 'ASSET', NULL, TRUE, '库存现金科目'),
(2, '1002', '银行存款', 'ASSET', NULL, TRUE, '银行存款科目'),
(3, '1122', '应收账款', 'ASSET', NULL, TRUE, '应收账款科目'),
(4, '2202', '应付账款', 'LIABILITY', NULL, TRUE, '应付账款科目'),
(5, '4001', '主营业务收入', 'REVENUE', NULL, TRUE, '主营业务收入科目'),
(6, '5001', '主营业务成本', 'EXPENSE', NULL, TRUE, '主营业务成本科目')
ON CONFLICT (id) DO NOTHING;

-- 插入默认会计期间
INSERT INTO fiscal_periods (id, code, name, start_date, end_date, fiscal_year, period_number) VALUES
(1, '2024-01', '2024年1月', '2024-01-01', '2024-01-31', 2024, 1),
(2, '2024-02', '2024年2月', '2024-02-01', '2024-02-29', 2024, 2),
(3, '2024-03', '2024年3月', '2024-03-01', '2024-03-31', 2024, 3)
ON CONFLICT (id) DO NOTHING;

-- 插入默认银行账户
INSERT INTO bank_accounts (id, code, account_name, account_number, bank_name, currency, account_type, opening_balance, current_balance, status) VALUES
(1, 'BANK001', '基本户', '1234567890123456789', '中国工商银行', 'CNY', 'CHECKING', 1000000.00, 1000000.00, 'ACTIVE'),
(2, 'BANK002', '一般户', '9876543210987654321', '中国建设银行', 'CNY', 'SAVINGS', 500000.00, 500000.00, 'ACTIVE')
ON CONFLICT (id) DO NOTHING;

-- 插入默认工作中心
INSERT INTO work_centers (id, code, name, description, work_center_type, capacity, efficiency, cost_per_hour, status) VALUES
(1, 'WC001', '装配线1', '主装配生产线', 'ASSEMBLY', 8.00, 95.00, 50.00, 'ACTIVE'),
(2, 'WC002', '加工中心1', 'CNC加工中心', 'MACHINE', 24.00, 90.00, 80.00, 'ACTIVE'),
(3, 'WC003', '质检站1', '质量检验工作站', 'INSPECTION', 8.00, 100.00, 30.00, 'ACTIVE'),
(4, 'WC004', '包装线1', '产品包装生产线', 'PACKAGING', 8.00, 98.00, 25.00, 'ACTIVE')
ON CONFLICT (id) DO NOTHING;

-- 插入系统配置
INSERT INTO system_configs (id, config_key, config_value, config_type, description, is_system) VALUES
(1, 'SYSTEM_NAME', '银河ERP系统', 'STRING', '系统名称', TRUE),
(2, 'SYSTEM_VERSION', '1.0.0', 'STRING', '系统版本', TRUE),
(3, 'DEFAULT_CURRENCY', 'CNY', 'STRING', '默认货币', FALSE),
(4, 'DEFAULT_LANGUAGE', 'zh-CN', 'STRING', '默认语言', FALSE),
(5, 'SESSION_TIMEOUT', '3600', 'INTEGER', '会话超时时间(秒)', FALSE),
(6, 'MAX_LOGIN_ATTEMPTS', '5', 'INTEGER', '最大登录尝试次数', FALSE)
ON CONFLICT (id) DO NOTHING;

-- 插入数据字典
INSERT INTO dictionaries (id, dict_type, dict_key, dict_value, dict_label, sort_order, is_system, description) VALUES
(1, 'GENDER', 'MALE', 'MALE', '男', 1, TRUE, '性别-男'),
(2, 'GENDER', 'FEMALE', 'FEMALE', '女', 2, TRUE, '性别-女'),
(3, 'STATUS', 'ACTIVE', 'ACTIVE', '启用', 1, TRUE, '状态-启用'),
(4, 'STATUS', 'INACTIVE', 'INACTIVE', '禁用', 2, TRUE, '状态-禁用'),
(5, 'PRIORITY', 'LOW', 'LOW', '低', 1, FALSE, '优先级-低'),
(6, 'PRIORITY', 'MEDIUM', 'MEDIUM', '中', 2, FALSE, '优先级-中'),
(7, 'PRIORITY', 'HIGH', 'HIGH', '高', 3, FALSE, '优先级-高'),
(8, 'PRIORITY', 'URGENT', 'URGENT', '紧急', 4, FALSE, '优先级-紧急')
ON CONFLICT (id) DO NOTHING;

-- 插入通知模板
INSERT INTO notification_templates (id, template_code, template_name, template_type, subject, content, variables) VALUES
(1, 'USER_WELCOME', '用户欢迎', 'EMAIL', '欢迎加入银河ERP系统', '亲爱的{name}，欢迎您加入银河ERP系统！您的用户名是：{username}', 'name,username'),
(2, 'PASSWORD_RESET', '密码重置', 'EMAIL', '密码重置通知', '您的密码已重置，新密码是：{password}，请及时登录修改。', 'password'),
(3, 'ORDER_CREATED', '订单创建', 'SYSTEM', '新订单创建', '订单{orderNumber}已创建，请及时处理。', 'orderNumber'),
(4, 'INVENTORY_LOW', '库存不足', 'SYSTEM', '库存预警', '物料{itemName}库存不足，当前库存：{quantity}', 'itemName,quantity')
ON CONFLICT (id) DO NOTHING;

-- 插入数据迁移记录示例
INSERT INTO data_migrations (id, migration_name, source_system, target_table, migration_type, total_records, success_records, failed_records, status) VALUES
(1, '客户数据迁移', 'Legacy_CRM', 'customers', 'IMPORT', 1000, 950, 50, 'COMPLETED'),
(2, '供应商数据迁移', 'Legacy_SCM', 'suppliers', 'IMPORT', 500, 480, 20, 'COMPLETED'),
(3, '物料数据迁移', 'Legacy_WMS', 'items', 'IMPORT', 2000, 1980, 20, 'COMPLETED')
ON CONFLICT (id) DO NOTHING;

-- 设置序列值
SELECT setval('companies_id_seq', (SELECT MAX(id) FROM companies));
SELECT setval('departments_id_seq', (SELECT MAX(id) FROM departments));
SELECT setval('roles_id_seq', (SELECT MAX(id) FROM roles));
SELECT setval('permissions_id_seq', (SELECT MAX(id) FROM permissions));
SELECT setval('users_id_seq', (SELECT MAX(id) FROM users));
SELECT setval('positions_id_seq', (SELECT MAX(id) FROM positions));
SELECT setval('salary_grades_id_seq', (SELECT MAX(id) FROM salary_grades));
SELECT setval('employees_id_seq', (SELECT MAX(id) FROM employees));
SELECT setval('leave_types_id_seq', (SELECT MAX(id) FROM leave_types));
SELECT setval('units_id_seq', (SELECT MAX(id) FROM units));
SELECT setval('item_categories_id_seq', (SELECT MAX(id) FROM item_categories));
SELECT setval('items_id_seq', (SELECT MAX(id) FROM items));
SELECT setval('warehouses_id_seq', (SELECT MAX(id) FROM warehouses));
SELECT setval('locations_id_seq', (SELECT MAX(id) FROM locations));
SELECT setval('customers_id_seq', (SELECT MAX(id) FROM customers));
SELECT setval('suppliers_id_seq', (SELECT MAX(id) FROM suppliers));
SELECT setval('accounts_id_seq', (SELECT MAX(id) FROM accounts));
SELECT setval('fiscal_periods_id_seq', (SELECT MAX(id) FROM fiscal_periods));
SELECT setval('bank_accounts_id_seq', (SELECT MAX(id) FROM bank_accounts));
SELECT setval('work_centers_id_seq', (SELECT MAX(id) FROM work_centers));
SELECT setval('system_configs_id_seq', (SELECT MAX(id) FROM system_configs));
SELECT setval('dictionaries_id_seq', (SELECT MAX(id) FROM dictionaries));
SELECT setval('notification_templates_id_seq', (SELECT MAX(id) FROM notification_templates));
SELECT setval('data_migrations_id_seq', (SELECT MAX(id) FROM data_migrations));
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  quotation_id INTEGER NOT NULL,
  line_number INTEGER DEFAULT 0,
  product_code VARCHAR(50),
  product_name VARCHAR(200) NOT NULL,
  description TEXT,
  quantity DECIMAL(15,4) NOT NULL,
  unit VARCHAR(20),
  unit_price DECIMAL(15,2) NOT NULL,
  discount_percent DECIMAL(5,2) DEFAULT 0,
  discount_amount DECIMAL(15,2) DEFAULT 0,
  line_total DECIMAL(15,2) DEFAULT 0,
  CONSTRAINT fk_quotation_items_quotation FOREIGN KEY (quotation_id) REFERENCES quotations(id) ON DELETE CASCADE
);

-- 销售订单表
CREATE TABLE IF NOT EXISTS sales_orders (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  created_by INTEGER,
  updated_by INTEGER,
  is_active BOOLEAN DEFAULT TRUE,
  code VARCHAR(50) NOT NULL,
  customer_id INTEGER NOT NULL,
  quotation_id INTEGER,
  order_date DATE NOT NULL,
  required_date DATE,
  promised_date DATE,
  subtotal DECIMAL(15,2) DEFAULT 0,
  tax_amount DECIMAL(15,2) DEFAULT 0,
  total_amount DECIMAL(15,2) DEFAULT 0,
  currency VARCHAR(3) DEFAULT 'CNY',
  payment_terms VARCHAR(100),
  delivery_address TEXT,
  notes TEXT,
  status VARCHAR(20) DEFAULT 'PENDING',
  CONSTRAINT uq_sales_orders_code UNIQUE (code),
  CONSTRAINT fk_sales_orders_customer FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE,
  CONSTRAINT fk_sales_orders_quotation FOREIGN KEY (quotation_id) REFERENCES quotations(id) ON DELETE SET NULL,
  CONSTRAINT fk_sales_orders_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_sales_orders_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 销售订单明细表
CREATE TABLE IF NOT EXISTS sales_order_items (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  sales_order_id INTEGER NOT NULL,
  line_number INTEGER DEFAULT 0,
  product_code VARCHAR(50),
  product_name VARCHAR(200) NOT NULL,
  description TEXT,
  quantity DECIMAL(15,4) NOT NULL,
  shipped_quantity DECIMAL(15,4) DEFAULT 0,
  unit VARCHAR(20),
  unit_price DECIMAL(15,2) NOT NULL,
  discount_percent DECIMAL(5,2) DEFAULT 0,
  discount_amount DECIMAL(15,2) DEFAULT 0,
  line_total DECIMAL(15,2) DEFAULT 0,
  CONSTRAINT fk_sales_order_items_sales_order FOREIGN KEY (sales_order_id) REFERENCES sales_orders(id) ON DELETE CASCADE
);

-- 发货单表
CREATE TABLE IF NOT EXISTS shipments (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  created_by INTEGER,
  updated_by INTEGER,
  is_active BOOLEAN DEFAULT TRUE,
  code VARCHAR(50) NOT NULL,
  sales_order_id INTEGER NOT NULL,
  shipment_date DATE NOT NULL,
  carrier VARCHAR(100),
  tracking_number VARCHAR(100),
  shipping_address TEXT,
  notes TEXT,
  status VARCHAR(20) DEFAULT 'PREPARING',
  CONSTRAINT uq_shipments_code UNIQUE (code),
  CONSTRAINT fk_shipments_sales_order FOREIGN KEY (sales_order_id) REFERENCES sales_orders(id) ON DELETE CASCADE,
  CONSTRAINT fk_shipments_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_shipments_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 发货单明细表
CREATE TABLE IF NOT EXISTS shipment_items (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  shipment_id INTEGER NOT NULL,
  sales_order_item_id INTEGER NOT NULL,
  quantity DECIMAL(15,4) NOT NULL,
  unit VARCHAR(20),
  notes TEXT,
  CONSTRAINT fk_shipment_items_shipment FOREIGN KEY (shipment_id) REFERENCES shipments(id) ON DELETE CASCADE,
  CONSTRAINT fk_shipment_items_sales_order_item FOREIGN KEY (sales_order_item_id) REFERENCES sales_order_items(id) ON DELETE CASCADE
);

-- 销售发票表
CREATE TABLE IF NOT EXISTS sales_invoices (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  created_by INTEGER,
  updated_by INTEGER,
  is_active BOOLEAN DEFAULT TRUE,
  code VARCHAR(50) NOT NULL,
  customer_id INTEGER NOT NULL,
  sales_order_id INTEGER,
  invoice_date DATE NOT NULL,
  due_date DATE,
  subtotal DECIMAL(15,2) DEFAULT 0,
  tax_amount DECIMAL(15,2) DEFAULT 0,
  total_amount DECIMAL(15,2) DEFAULT 0,
  paid_amount DECIMAL(15,2) DEFAULT 0,
  currency VARCHAR(3) DEFAULT 'CNY',
  payment_terms VARCHAR(100),
  notes TEXT,
  status VARCHAR(20) DEFAULT 'DRAFT',
  CONSTRAINT uq_sales_invoices_code UNIQUE (code),
  CONSTRAINT fk_sales_invoices_customer FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE,
  CONSTRAINT fk_sales_invoices_sales_order FOREIGN KEY (sales_order_id) REFERENCES sales_orders(id) ON DELETE SET NULL,
  CONSTRAINT fk_sales_invoices_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_sales_invoices_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 销售发票明细表
CREATE TABLE IF NOT EXISTS sales_invoice_items (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  sales_invoice_id INTEGER NOT NULL,
  line_number INTEGER DEFAULT 0,
  product_code VARCHAR(50),
  product_name VARCHAR(200) NOT NULL,
  description TEXT,
  quantity DECIMAL(15,4) NOT NULL,
  unit VARCHAR(20),
  unit_price DECIMAL(15,2) NOT NULL,
  discount_percent DECIMAL(5,2) DEFAULT 0,
  discount_amount DECIMAL(15,2) DEFAULT 0,
  line_total DECIMAL(15,2) DEFAULT 0,
  CONSTRAINT fk_sales_invoice_items_sales_invoice FOREIGN KEY (sales_invoice_id) REFERENCES sales_invoices(id) ON DELETE CASCADE
);

-- ============================================================================
-- 4. 库存管理模块 (13个表)
-- ============================================================================

-- 物料分类表
CREATE TABLE IF NOT EXISTS item_categories (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  created_by INTEGER,
  updated_by INTEGER,
  is_active BOOLEAN DEFAULT TRUE,
  code VARCHAR(50) NOT NULL,
  name VARCHAR(100) NOT NULL,
  description TEXT,
  parent_id INTEGER,
  CONSTRAINT uq_item_categories_code UNIQUE (code),
  CONSTRAINT fk_item_categories_parent FOREIGN KEY (parent_id) REFERENCES item_categories(id) ON DELETE SET NULL,
  CONSTRAINT fk_item_categories_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_item_categories_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 计量单位表
CREATE TABLE IF NOT EXISTS units (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  created_by INTEGER,
  updated_by INTEGER,
  is_active BOOLEAN DEFAULT TRUE,
  code VARCHAR(20) NOT NULL,
  name VARCHAR(50) NOT NULL,
  description TEXT,
  unit_type VARCHAR(20) DEFAULT 'COUNT',
  base_unit_id INTEGER,
  conversion_factor DECIMAL(15,6) DEFAULT 1,
  CONSTRAINT uq_units_code UNIQUE (code),
  CONSTRAINT fk_units_base_unit FOREIGN KEY (base_unit_id) REFERENCES units(id) ON DELETE SET NULL,
  CONSTRAINT fk_units_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_units_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 物料表
CREATE TABLE IF NOT EXISTS items (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  created_by INTEGER,
  updated_by INTEGER,
  is_active BOOLEAN DEFAULT TRUE,
  code VARCHAR(50) NOT NULL,
  name VARCHAR(200) NOT NULL,
  description TEXT,
  category_id INTEGER,
  unit_id INTEGER,
  item_type VARCHAR(20) DEFAULT 'PRODUCT',
  specification TEXT,
  brand VARCHAR(100),
  model VARCHAR(100),
  barcode VARCHAR(100),
  weight DECIMAL(15,4),
  dimensions VARCHAR(100),
  cost DECIMAL(15,2) DEFAULT 0,
  price DECIMAL(15,2) DEFAULT 0,
  currency VARCHAR(3) DEFAULT 'CNY',
  min_stock_level DECIMAL(15,4) DEFAULT 0,
  max_stock_level DECIMAL(15,4) DEFAULT 0,
  reorder_point DECIMAL(15,4) DEFAULT 0,
  lead_time_days INTEGER DEFAULT 0,
  status VARCHAR(20) DEFAULT 'ACTIVE',
  CONSTRAINT uq_items_code UNIQUE (code),
  CONSTRAINT uq_items_barcode UNIQUE (barcode),
  CONSTRAINT fk_items_category FOREIGN KEY (category_id) REFERENCES item_categories(id) ON DELETE SET NULL,
  CONSTRAINT fk_items_unit FOREIGN KEY (unit_id) REFERENCES units(id) ON DELETE SET NULL,
  CONSTRAINT fk_items_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_items_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 仓库表
CREATE TABLE IF NOT EXISTS warehouses (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  created_by INTEGER,
  updated_by INTEGER,
  is_active BOOLEAN DEFAULT TRUE,
  code VARCHAR(50) NOT NULL,
  name VARCHAR(100) NOT NULL,
  description TEXT,
  address TEXT,
  manager_id INTEGER,
  warehouse_type VARCHAR(20) DEFAULT 'MAIN',
  status VARCHAR(20) DEFAULT 'ACTIVE',
  CONSTRAINT uq_warehouses_code UNIQUE (code),
  CONSTRAINT fk_warehouses_manager FOREIGN KEY (manager_id) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_warehouses_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_warehouses_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 库位表
CREATE TABLE IF NOT EXISTS locations (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  created_by INTEGER,
  updated_by INTEGER,
  is_active BOOLEAN DEFAULT TRUE,
  warehouse_id INTEGER NOT NULL,
  code VARCHAR(50) NOT NULL,
  name VARCHAR(100) NOT NULL,
  description TEXT,
  location_type VARCHAR(20) DEFAULT 'SHELF',
  capacity DECIMAL(15,4),
  status VARCHAR(20) DEFAULT 'ACTIVE',
  CONSTRAINT uq_locations_warehouse_code UNIQUE (warehouse_id, code),
  CONSTRAINT fk_locations_warehouse FOREIGN KEY (warehouse_id) REFERENCES warehouses(id) ON DELETE CASCADE,
  CONSTRAINT fk_locations_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_locations_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 库存表
CREATE TABLE IF NOT EXISTS inventories (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  created_by INTEGER,
  updated_by INTEGER,
  is_active BOOLEAN DEFAULT TRUE,
  item_id INTEGER NOT NULL,
  warehouse_id INTEGER NOT NULL,
  location_id INTEGER,
  quantity DECIMAL(15,4) DEFAULT 0,
  reserved_quantity DECIMAL(15,4) DEFAULT 0,
  available_quantity DECIMAL(15,4) DEFAULT 0,
  unit_cost DECIMAL(15,2) DEFAULT 0,
  total_cost DECIMAL(15,2) DEFAULT 0,
  last_movement_date TIMESTAMP WITH TIME ZONE,
  CONSTRAINT uq_inventories_item_warehouse_location UNIQUE (item_id, warehouse_id, location_id),
  CONSTRAINT fk_inventories_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE,
  CONSTRAINT fk_inventories_warehouse FOREIGN KEY (warehouse_id) REFERENCES warehouses(id) ON DELETE CASCADE,
  CONSTRAINT fk_inventories_location FOREIGN KEY (location_id) REFERENCES locations(id) ON DELETE SET NULL,
  CONSTRAINT fk_inventories_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_inventories_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 批次表
CREATE TABLE IF NOT EXISTS batches (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  created_by INTEGER,
  updated_by INTEGER,
  is_active BOOLEAN DEFAULT TRUE,
  batch_number VARCHAR(100) NOT NULL,
  item_id INTEGER NOT NULL,
  warehouse_id INTEGER NOT NULL,
  location_id INTEGER,
  quantity DECIMAL(15,4) DEFAULT 0,
  production_date DATE,
  expiry_date DATE,
  supplier_batch_number VARCHAR(100),
  notes TEXT,
  status VARCHAR(20) DEFAULT 'ACTIVE',
  CONSTRAINT uq_batches_batch_number UNIQUE (batch_number),
  CONSTRAINT fk_batches_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE,
  CONSTRAINT fk_batches_warehouse FOREIGN KEY (warehouse_id) REFERENCES warehouses(id) ON DELETE CASCADE,
  CONSTRAINT fk_batches_location FOREIGN KEY (location_id) REFERENCES locations(id) ON DELETE SET NULL,
  CONSTRAINT fk_batches_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_batches_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 序列号表
CREATE TABLE IF NOT EXISTS serial_numbers (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  created_by INTEGER,
  updated_by INTEGER,
  is_active BOOLEAN DEFAULT TRUE,
  serial_number VARCHAR(100) NOT NULL,
  item_id INTEGER NOT NULL,
  warehouse_id INTEGER NOT NULL,
  location_id INTEGER,
  batch_id INTEGER,
  production_date DATE,
  warranty_expiry_date DATE,
  notes TEXT,
  status VARCHAR(20) DEFAULT 'AVAILABLE',
  CONSTRAINT uq_serial_numbers_serial_number UNIQUE (serial_number),
  CONSTRAINT fk_serial_numbers_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE,
  CONSTRAINT fk_serial_numbers_warehouse FOREIGN KEY (warehouse_id) REFERENCES warehouses(id) ON DELETE CASCADE,
  CONSTRAINT fk_serial_numbers_location FOREIGN KEY (location_id) REFERENCES locations(id) ON DELETE SET NULL,
  CONSTRAINT fk_serial_numbers_batch FOREIGN KEY (batch_id) REFERENCES batches(id) ON DELETE SET NULL,
  CONSTRAINT fk_serial_numbers_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_serial_numbers_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 库存移动表
CREATE TABLE IF NOT EXISTS inventory_movements (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  created_by INTEGER,
  updated_by INTEGER,
  is_active BOOLEAN DEFAULT TRUE,
  movement_date TIMESTAMP WITH TIME ZONE NOT NULL,
  movement_type VARCHAR(20) NOT NULL,
  reference_type VARCHAR(50),
  reference_id INTEGER,
  item_id INTEGER NOT NULL,
  warehouse_id INTEGER NOT NULL,
  location_id INTEGER,
  batch_id INTEGER,
  quantity DECIMAL(15,4) NOT NULL,
  unit_cost DECIMAL(15,2) DEFAULT 0,
  total_cost DECIMAL(15,2) DEFAULT 0,
  notes TEXT,
  CONSTRAINT fk_inventory_movements_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE,
  CONSTRAINT fk_inventory_movements_warehouse FOREIGN KEY (warehouse_id) REFERENCES warehouses(id) ON DELETE CASCADE,
  CONSTRAINT fk_inventory_movements_location FOREIGN KEY (location_id) REFERENCES locations(id) ON DELETE SET NULL,
  CONSTRAINT fk_inventory_movements_batch FOREIGN KEY (batch_id) REFERENCES batches(id) ON DELETE SET NULL,
  CONSTRAINT fk_inventory_movements_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_inventory_movements_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 库存调整表
CREATE TABLE IF NOT EXISTS inventory_adjustments (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  created_by INTEGER,
  updated_by INTEGER,
  is_active BOOLEAN DEFAULT TRUE,
  code VARCHAR(50) NOT NULL,
  adjustment_date DATE NOT NULL,
  warehouse_id INTEGER NOT NULL,
  reason VARCHAR(100),
  notes TEXT,
  approved_by INTEGER,
  approved_at TIMESTAMP WITH TIME ZONE,
  status VARCHAR(20) DEFAULT 'DRAFT',
  CONSTRAINT uq_inventory_adjustments_code UNIQUE (code),
  CONSTRAINT fk_inventory_adjustments_warehouse FOREIGN KEY (warehouse_id) REFERENCES warehouses(id) ON DELETE CASCADE,
  CONSTRAINT fk_inventory_adjustments_approved_by FOREIGN KEY (approved_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_inventory_adjustments_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_inventory_adjustments_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 库存调整明细表
CREATE TABLE IF NOT EXISTS inventory_adjustment_items (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  inventory_adjustment_id INTEGER NOT NULL,
  item_id INTEGER NOT NULL,
  location_id INTEGER,
  batch_id INTEGER,
  system_quantity DECIMAL(15,4) DEFAULT 0,
  actual_quantity DECIMAL(15,4) DEFAULT 0,
  adjustment_quantity DECIMAL(15,4) DEFAULT 0,
  unit_cost DECIMAL(15,2) DEFAULT 0,
  total_cost DECIMAL(15,2) DEFAULT 0,
  reason VARCHAR(100),
  notes TEXT,
  CONSTRAINT fk_inventory_adjustment_items_inventory_adjustment FOREIGN KEY (inventory_adjustment_id) REFERENCES inventory_adjustments(id) ON DELETE CASCADE,
  CONSTRAINT fk_inventory_adjustment_items_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE,
  CONSTRAINT fk_inventory_adjustment_items_location FOREIGN KEY (location_id) REFERENCES locations(id) ON DELETE SET NULL,
  CONSTRAINT fk_inventory_adjustment_items_batch FOREIGN KEY (batch_id) REFERENCES batches(id) ON DELETE SET NULL
);

-- 库存盘点表
CREATE TABLE IF NOT EXISTS stock_counts (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  created_by INTEGER,
  updated_by INTEGER,
  is_active BOOLEAN DEFAULT TRUE,
  code VARCHAR(50) NOT NULL,
  count_date DATE NOT NULL,
  warehouse_id INTEGER NOT NULL,
  count_type VARCHAR(20) DEFAULT 'FULL',
  notes TEXT,
  approved_by INTEGER,
  approved_at TIMESTAMP WITH TIME ZONE,
  status VARCHAR(20) DEFAULT 'PLANNING',
  CONSTRAINT uq_stock_counts_code UNIQUE (code),
  CONSTRAINT fk_stock_counts_warehouse FOREIGN KEY (warehouse_id) REFERENCES warehouses(id) ON DELETE CASCADE,
  CONSTRAINT fk_stock_counts_approved_by FOREIGN KEY (approved_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_stock_counts_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_stock_counts_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 库存盘点明细表
CREATE TABLE IF NOT EXISTS stock_count_items (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  stock_count_id INTEGER NOT NULL,
  item_id INTEGER NOT NULL,
  location_id INTEGER,
  batch_id INTEGER,
  system_quantity DECIMAL(15,4) DEFAULT 0,
  counted_quantity DECIMAL(15,4) DEFAULT 0,
  variance_quantity DECIMAL(15,4) DEFAULT 0,
  unit_cost DECIMAL(15,2) DEFAULT 0,
  variance_cost DECIMAL(15,2) DEFAULT 0,
  counter_id INTEGER,
  count_time TIMESTAMP WITH TIME ZONE,
  notes TEXT,
  CONSTRAINT fk_stock_count_items_stock_count FOREIGN KEY (stock_count_id) REFERENCES stock_counts(id) ON DELETE CASCADE,
  CONSTRAINT fk_stock_count_items_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE,
  CONSTRAINT fk_stock_count_items_location FOREIGN KEY (location_id) REFERENCES locations(id) ON DELETE SET NULL,
  CONSTRAINT fk_stock_count_items_batch FOREIGN KEY (batch_id) REFERENCES batches(id) ON DELETE SET NULL,
  CONSTRAINT fk_stock_count_items_counter FOREIGN KEY (counter_id) REFERENCES users(id) ON DELETE SET NULL
);

-- ============================================================================
-- 5. 采购管理模块 (15个表)
-- ============================================================================

-- 供应商表
CREATE TABLE IF NOT EXISTS suppliers (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  created_by INTEGER,
  updated_by INTEGER,
  is_active BOOLEAN DEFAULT TRUE,
  code VARCHAR(50) NOT NULL,
  name VARCHAR(200) NOT NULL,
  description TEXT,
  supplier_type VARCHAR(20) DEFAULT 'COMPANY',
  industry VARCHAR(100),
  address TEXT,
  phone VARCHAR(50),
  email VARCHAR(100),
  website VARCHAR(200),
  tax_number VARCHAR(50),
  credit_rating VARCHAR(10),
  payment_terms VARCHAR(50),
  lead_time_days INTEGER DEFAULT 0,
  buyer_id INTEGER,
  status VARCHAR(20) DEFAULT 'ACTIVE',
  CONSTRAINT uq_suppliers_code UNIQUE (code),
  CONSTRAINT fk_suppliers_buyer FOREIGN KEY (buyer_id) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_suppliers_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_suppliers_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 供应商联系人表
CREATE TABLE IF NOT EXISTS supplier_contacts (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  created_by INTEGER,
  updated_by INTEGER,
  is_active BOOLEAN DEFAULT TRUE,
  supplier_id INTEGER NOT NULL,
  name VARCHAR(100) NOT NULL,
  title VARCHAR(100),
  phone VARCHAR(50),
  email VARCHAR(100),
  is_primary BOOLEAN DEFAULT FALSE,
  notes TEXT,
  CONSTRAINT fk_supplier_contacts_supplier FOREIGN KEY (supplier_id) REFERENCES suppliers(id) ON DELETE CASCADE,
  CONSTRAINT fk_supplier_contacts_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_supplier_contacts_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 采购申请表
CREATE TABLE IF NOT EXISTS purchase_requests (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  created_by INTEGER,
  updated_by INTEGER,
  is_active BOOLEAN DEFAULT TRUE,
  code VARCHAR(50) NOT NULL,
  request_date DATE NOT NULL,
  required_date DATE,
  department_id INTEGER,
  requester_id INTEGER NOT NULL,
  priority VARCHAR(20) DEFAULT 'MEDIUM',
  justification TEXT,
  total_estimated_cost DECIMAL(15,2) DEFAULT 0,
  currency VARCHAR(3) DEFAULT 'CNY',
  approved_by INTEGER,
  approved_at TIMESTAMP WITH TIME ZONE,
  status VARCHAR(20) DEFAULT 'DRAFT',
  CONSTRAINT uq_purchase_requests_code UNIQUE (code),
  CONSTRAINT fk_purchase_requests_department FOREIGN KEY (department_id) REFERENCES departments(id) ON DELETE SET NULL,
  CONSTRAINT fk_purchase_requests_requester FOREIGN KEY (requester_id) REFERENCES users(id) ON DELETE CASCADE,
  CONSTRAINT fk_purchase_requests_approved_by FOREIGN KEY (approved_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_purchase_requests_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_purchase_requests_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 采购申请明细表
CREATE TABLE IF NOT EXISTS purchase_request_items (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  purchase_request_id INTEGER NOT NULL,
  line_number INTEGER DEFAULT 0,
  item_id INTEGER,
  item_code VARCHAR(50),
  item_name VARCHAR(200) NOT NULL,
  description TEXT,
  quantity DECIMAL(15,4) NOT NULL,
  unit VARCHAR(20),
  estimated_unit_cost DECIMAL(15,2) DEFAULT 0,
  estimated_total_cost DECIMAL(15,2) DEFAULT 0,
  required_date DATE,
  notes TEXT,
  CONSTRAINT fk_purchase_request_items_purchase_request FOREIGN KEY (purchase_request_id) REFERENCES purchase_requests(id) ON DELETE CASCADE,
  CONSTRAINT fk_purchase_request_items_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE SET NULL
);

-- 询价单表
CREATE TABLE IF NOT EXISTS rfqs (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  created_by INTEGER,
  updated_by INTEGER,
  is_active BOOLEAN DEFAULT TRUE,
  code VARCHAR(50) NOT NULL,
  rfq_date DATE NOT NULL,
  response_deadline DATE,
  buyer_id INTEGER,
  terms_and_conditions TEXT,
  notes TEXT,
  status VARCHAR(20) DEFAULT 'DRAFT',
  CONSTRAINT uq_rfqs_code UNIQUE (code),
  CONSTRAINT fk_rfqs_buyer FOREIGN KEY (buyer_id) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_rfqs_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_rfqs_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 询价单明细表
CREATE TABLE IF NOT EXISTS rfq_items (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  rfq_id INTEGER NOT NULL,
  line_number INTEGER DEFAULT 0,
  item_id INTEGER,
  item_code VARCHAR(50),
  item_name VARCHAR(200) NOT NULL,
  description TEXT,
  quantity DECIMAL(15,4) NOT NULL,
  unit VARCHAR(20),
  required_date DATE,
  specifications TEXT,
  CONSTRAINT fk_rfq_items_rfq FOREIGN KEY (rfq_id) REFERENCES rfqs(id) ON DELETE CASCADE,
  CONSTRAINT fk_rfq_items_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE SET NULL
);

-- 询价供应商表
CREATE TABLE IF NOT EXISTS rfq_suppliers (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  rfq_id INTEGER NOT NULL,
  supplier_id INTEGER NOT NULL,
  sent_date DATE,
  response_date DATE,
  status VARCHAR(20) DEFAULT 'SENT',
  CONSTRAINT fk_rfq_suppliers_rfq FOREIGN KEY (rfq_id) REFERENCES rfqs(id) ON DELETE CASCADE,
  CONSTRAINT fk_rfq_suppliers_supplier FOREIGN KEY (supplier_id) REFERENCES suppliers(id) ON DELETE CASCADE,
  CONSTRAINT uq_rfq_suppliers_rfq_supplier UNIQUE (rfq_id, supplier_id)
);

-- 供应商报价表
CREATE TABLE IF NOT EXISTS supplier_quotes (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  created_by INTEGER,
  updated_by INTEGER,
  is_active BOOLEAN DEFAULT TRUE,
  code VARCHAR(50) NOT NULL,
  rfq_id INTEGER NOT NULL,
  supplier_id INTEGER NOT NULL,
  quote_date DATE NOT NULL,
  valid_until DATE,
  subtotal DECIMAL(15,2) DEFAULT 0,
  tax_amount DECIMAL(15,2) DEFAULT 0,
  total_amount DECIMAL(15,2) DEFAULT 0,
  currency VARCHAR(3) DEFAULT 'CNY',
  payment_terms VARCHAR(100),
  delivery_terms VARCHAR(100),
  notes TEXT,
  status VARCHAR(20) DEFAULT 'RECEIVED',
  CONSTRAINT uq_supplier_quotes_code UNIQUE (code),
  CONSTRAINT fk_supplier_quotes_rfq FOREIGN KEY (rfq_id) REFERENCES rfqs(id) ON DELETE CASCADE,
  CONSTRAINT fk_supplier_quotes_supplier FOREIGN KEY (supplier_id) REFERENCES suppliers(id) ON DELETE CASCADE,
  CONSTRAINT fk_supplier_quotes_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_supplier_quotes_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 供应商报价明细表
CREATE TABLE IF NOT EXISTS supplier_quote_items (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  supplier_quote_id INTEGER NOT NULL,
  rfq_item_id INTEGER NOT NULL,
  line_number INTEGER DEFAULT 0,
  quantity DECIMAL(15,4) NOT NULL,
  unit VARCHAR(20),
  unit_price DECIMAL(15,2) NOT NULL,
  discount_percent DECIMAL(5,2) DEFAULT 0,
  discount_amount DECIMAL(15,2) DEFAULT 0,
  line_total DECIMAL(15,2) DEFAULT 0,
  lead_time_days INTEGER DEFAULT 0,
  notes TEXT,
  CONSTRAINT fk_supplier_quote_items_supplier_quote FOREIGN KEY (supplier_quote_id) REFERENCES supplier_quotes(id) ON DELETE CASCADE,
  CONSTRAINT fk_supplier_quote_items_rfq_item FOREIGN KEY (rfq_item_id) REFERENCES rfq_items(id) ON DELETE CASCADE
);

-- 采购订单表
CREATE TABLE IF NOT EXISTS purchase_orders (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE,
  created_by INTEGER,
  updated_by INTEGER,
  is_active BOOLEAN DEFAULT TRUE,
  code VARCHAR(50) NOT NULL,
  supplier_id INTEGER NOT NULL,
  supplier_quote_id INTEGER,
  order_date DATE NOT NULL,
  required_date DATE,
  promised_date DATE,
  subtotal DECIMAL(15,2) DEFAULT 0,
  tax_amount DECIMAL(15,2) DEFAULT 0,
  total_amount DECIMAL(15,2) DEFAULT 0,
  currency VARCHAR(3) DEFAULT 'CNY',
  payment_terms VARCHAR(100),
  delivery_address TEXT,
  buyer_id INTEGER,
  notes TEXT,
  status VARCHAR(20) DEFAULT 'DRAFT',
  CONSTRAINT uq