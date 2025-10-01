-- ============================================================================
-- GalaxyERP 项目管理模块 - PostgreSQL 建表脚本
-- 生成时间: 2025-10-01
-- 说明: 基于当前 /internal/models 下的结构，包含项目相关表及必要依赖
-- ============================================================================

-- 建议使用一个独立的 schema（可选）
-- CREATE SCHEMA IF NOT EXISTS galaxyerp;
-- SET search_path TO galaxyerp;

-- ============================================================================
-- 基础依赖表（部门、职位、员工、客户）
-- ============================================================================

-- 部门表
CREATE TABLE IF NOT EXISTS departments (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL,
  name VARCHAR(255) NOT NULL,
  code VARCHAR(50) NOT NULL,
  description TEXT NULL,
  company_id BIGINT NOT NULL,
  parent_id BIGINT NULL,
  is_active BOOLEAN DEFAULT TRUE,
  created_by BIGINT NULL,
  updated_by BIGINT NULL,
  CONSTRAINT uq_departments_code UNIQUE (code),
  CONSTRAINT fk_departments_parent FOREIGN KEY (parent_id) REFERENCES departments(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_departments_deleted_at ON departments (deleted_at);
CREATE INDEX IF NOT EXISTS idx_departments_company_id ON departments (company_id);
CREATE INDEX IF NOT EXISTS idx_departments_parent_id ON departments (parent_id);
CREATE INDEX IF NOT EXISTS idx_departments_is_active ON departments (is_active);

-- 职位表
CREATE TABLE IF NOT EXISTS positions (
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
  department_id BIGINT NOT NULL,
  CONSTRAINT uq_positions_code UNIQUE (code),
  CONSTRAINT fk_positions_department FOREIGN KEY (department_id) REFERENCES departments(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_positions_deleted_at ON positions (deleted_at);
CREATE INDEX IF NOT EXISTS idx_positions_is_active ON positions (is_active);
CREATE INDEX IF NOT EXISTS idx_positions_department_id ON positions (department_id);

-- 员工表
CREATE TABLE IF NOT EXISTS employees (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL,
  code VARCHAR(50) NOT NULL,
  first_name VARCHAR(100) NOT NULL,
  last_name VARCHAR(100) NOT NULL,
  full_name VARCHAR(255) NOT NULL,
  email VARCHAR(255) NOT NULL,
  phone VARCHAR(20) NULL,
  date_of_birth TIMESTAMP NULL,
  gender VARCHAR(10) NULL,
  hire_date TIMESTAMP NULL,
  department_id BIGINT NULL,
  position_id BIGINT NULL,
  manager_id BIGINT NULL,
  status VARCHAR(50) DEFAULT 'active',
  emergency_contact VARCHAR(255) NULL,
  id_number VARCHAR(100) NULL,
  address TEXT NULL,
  CONSTRAINT uq_employees_code UNIQUE (code),
  CONSTRAINT uq_employees_email UNIQUE (email),
  CONSTRAINT fk_employees_department FOREIGN KEY (department_id) REFERENCES departments(id) ON DELETE SET NULL,
  CONSTRAINT fk_employees_position FOREIGN KEY (position_id) REFERENCES positions(id) ON DELETE SET NULL,
  CONSTRAINT fk_employees_manager FOREIGN KEY (manager_id) REFERENCES employees(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_employees_deleted_at ON employees (deleted_at);
CREATE INDEX IF NOT EXISTS idx_employees_department_id ON employees (department_id);
CREATE INDEX IF NOT EXISTS idx_employees_position_id ON employees (position_id);
CREATE INDEX IF NOT EXISTS idx_employees_manager_id ON employees (manager_id);
CREATE INDEX IF NOT EXISTS idx_employees_status ON employees (status);

-- 客户表
CREATE TABLE IF NOT EXISTS customers (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL,
  name VARCHAR(255) NOT NULL,
  code VARCHAR(100) NOT NULL,
  email VARCHAR(255) NULL,
  phone VARCHAR(50) NULL,
  address TEXT NULL,
  city VARCHAR(100) NULL,
  state VARCHAR(100) NULL,
  postal_code VARCHAR(20) NULL,
  country VARCHAR(100) NULL,
  contact_person VARCHAR(255) NULL,
  credit_limit NUMERIC(15,2) DEFAULT 0,
  customer_group VARCHAR(100) NULL,
  territory VARCHAR(100) NULL,
  is_active BOOLEAN DEFAULT TRUE,
  CONSTRAINT uq_customers_code UNIQUE (code)
);
CREATE INDEX IF NOT EXISTS idx_customers_deleted_at ON customers (deleted_at);
CREATE INDEX IF NOT EXISTS idx_customers_is_active ON customers (is_active);

-- ============================================================================
-- 项目管理核心表
-- ============================================================================

-- 项目表（参考 models，字段包含编号、名称、客户、经理、计划/实际日期等）
CREATE TABLE IF NOT EXISTS projects (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL,
  created_by BIGINT NULL,
  updated_by BIGINT NULL,
  project_number VARCHAR(100) NOT NULL,
  project_name VARCHAR(255) NOT NULL,
  description TEXT NULL,
  customer_id BIGINT NULL,
  manager_id BIGINT NOT NULL,
  start_date TIMESTAMP NOT NULL,
  end_date TIMESTAMP NOT NULL,
  actual_start_date TIMESTAMP NULL,
  actual_end_date TIMESTAMP NULL,
  budget NUMERIC(15,2) DEFAULT 0,
  actual_cost NUMERIC(15,2) DEFAULT 0,
  progress NUMERIC(5,2) DEFAULT 0,
  priority VARCHAR(20) DEFAULT 'normal',
  status VARCHAR(50) DEFAULT 'planning',
  CONSTRAINT uq_projects_project_number UNIQUE (project_number),
  CONSTRAINT fk_projects_customer FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE SET NULL,
  CONSTRAINT fk_projects_manager FOREIGN KEY (manager_id) REFERENCES employees(id) ON DELETE RESTRICT
);
CREATE INDEX IF NOT EXISTS idx_projects_deleted_at ON projects (deleted_at);
CREATE INDEX IF NOT EXISTS idx_projects_customer_id ON projects (customer_id);
CREATE INDEX IF NOT EXISTS idx_projects_manager_id ON projects (manager_id);
CREATE INDEX IF NOT EXISTS idx_projects_start_date ON projects (start_date);
CREATE INDEX IF NOT EXISTS idx_projects_end_date ON projects (end_date);
CREATE INDEX IF NOT EXISTS idx_projects_priority ON projects (priority);
CREATE INDEX IF NOT EXISTS idx_projects_status ON projects (status);
CREATE INDEX IF NOT EXISTS idx_projects_status_priority ON projects (status, priority);
CREATE INDEX IF NOT EXISTS idx_projects_manager_status ON projects (manager_id, status);
CREATE INDEX IF NOT EXISTS idx_projects_customer_status ON projects (customer_id, status);

-- 任务表
CREATE TABLE IF NOT EXISTS tasks (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL,
  created_by BIGINT NULL,
  updated_by BIGINT NULL,
  task_number VARCHAR(100) NOT NULL,
  task_name VARCHAR(255) NOT NULL,
  description TEXT NULL,
  project_id BIGINT NOT NULL,
  parent_task_id BIGINT NULL,
  assignee_id BIGINT NULL,
  start_date TIMESTAMP NOT NULL,
  end_date TIMESTAMP NOT NULL,
  actual_start_date TIMESTAMP NULL,
  actual_end_date TIMESTAMP NULL,
  estimated_hours NUMERIC(8,2) DEFAULT 0,
  actual_hours NUMERIC(8,2) DEFAULT 0,
  progress NUMERIC(5,2) DEFAULT 0,
  priority VARCHAR(20) DEFAULT 'normal',
  status VARCHAR(50) DEFAULT 'todo',
  CONSTRAINT uq_tasks_task_number UNIQUE (task_number),
  CONSTRAINT fk_tasks_project FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
  CONSTRAINT fk_tasks_parent_task FOREIGN KEY (parent_task_id) REFERENCES tasks(id) ON DELETE SET NULL,
  CONSTRAINT fk_tasks_assignee FOREIGN KEY (assignee_id) REFERENCES employees(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_tasks_deleted_at ON tasks (deleted_at);
CREATE INDEX IF NOT EXISTS idx_tasks_project_id ON tasks (project_id);
CREATE INDEX IF NOT EXISTS idx_tasks_parent_task_id ON tasks (parent_task_id);
CREATE INDEX IF NOT EXISTS idx_tasks_assignee_id ON tasks (assignee_id);
CREATE INDEX IF NOT EXISTS idx_tasks_start_date ON tasks (start_date);
CREATE INDEX IF NOT EXISTS idx_tasks_end_date ON tasks (end_date);
CREATE INDEX IF NOT EXISTS idx_tasks_priority ON tasks (priority);
CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks (status);
CREATE INDEX IF NOT EXISTS idx_tasks_project_status ON tasks (project_id, status);
CREATE INDEX IF NOT EXISTS idx_tasks_assignee_status ON tasks (assignee_id, status);
CREATE INDEX IF NOT EXISTS idx_tasks_project_assignee ON tasks (project_id, assignee_id);

-- 里程碑表
CREATE TABLE IF NOT EXISTS milestones (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL,
  project_id BIGINT NOT NULL,
  milestone_name VARCHAR(255) NOT NULL,
  description TEXT NULL,
  due_date TIMESTAMP NOT NULL,
  completed_date TIMESTAMP NULL,
  status VARCHAR(50) DEFAULT 'pending',
  CONSTRAINT fk_milestones_project FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_milestones_deleted_at ON milestones (deleted_at);
CREATE INDEX IF NOT EXISTS idx_milestones_project_id ON milestones (project_id);
CREATE INDEX IF NOT EXISTS idx_milestones_due_date ON milestones (due_date);
CREATE INDEX IF NOT EXISTS idx_milestones_status ON milestones (status);

-- 工时记录表
CREATE TABLE IF NOT EXISTS time_entries (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL,
  employee_id BIGINT NOT NULL,
  project_id BIGINT NOT NULL,
  task_id BIGINT NULL,
  date TIMESTAMP NOT NULL,
  start_time TIMESTAMP NOT NULL,
  end_time TIMESTAMP NOT NULL,
  hours NUMERIC(8,2) NOT NULL,
  description TEXT NULL,
  is_billable BOOLEAN DEFAULT TRUE,
  hourly_rate NUMERIC(10,2) DEFAULT 0,
  amount NUMERIC(15,2) DEFAULT 0,
  status VARCHAR(50) DEFAULT 'draft',
  CONSTRAINT fk_time_entries_employee FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE CASCADE,
  CONSTRAINT fk_time_entries_project FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
  CONSTRAINT fk_time_entries_task FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_time_entries_deleted_at ON time_entries (deleted_at);
CREATE INDEX IF NOT EXISTS idx_time_entries_employee_id ON time_entries (employee_id);
CREATE INDEX IF NOT EXISTS idx_time_entries_project_id ON time_entries (project_id);
CREATE INDEX IF NOT EXISTS idx_time_entries_task_id ON time_entries (task_id);
CREATE INDEX IF NOT EXISTS idx_time_entries_date ON time_entries (date);
CREATE INDEX IF NOT EXISTS idx_time_entries_is_billable ON time_entries (is_billable);
CREATE INDEX IF NOT EXISTS idx_time_entries_status ON time_entries (status);
CREATE INDEX IF NOT EXISTS idx_time_entries_employee_date ON time_entries (employee_id, date);
CREATE INDEX IF NOT EXISTS idx_time_entries_project_date ON time_entries (project_id, date);
CREATE INDEX IF NOT EXISTS idx_time_entries_task_date ON time_entries (task_id, date);

-- ============================================================================
-- 项目相关辅助表
-- ============================================================================

-- 项目成员表
CREATE TABLE IF NOT EXISTS project_members (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL,
  project_id BIGINT NOT NULL,
  employee_id BIGINT NOT NULL,
  role VARCHAR(100) NOT NULL,
  join_date TIMESTAMP NOT NULL,
  leave_date TIMESTAMP NULL,
  is_active BOOLEAN DEFAULT TRUE,
  CONSTRAINT fk_project_members_project FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
  CONSTRAINT fk_project_members_employee FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_project_members_deleted_at ON project_members (deleted_at);
CREATE INDEX IF NOT EXISTS idx_project_members_project_id ON project_members (project_id);
CREATE INDEX IF NOT EXISTS idx_project_members_employee_id ON project_members (employee_id);
CREATE INDEX IF NOT EXISTS idx_project_members_join_date ON project_members (join_date);
CREATE INDEX IF NOT EXISTS idx_project_members_leave_date ON project_members (leave_date);
CREATE INDEX IF NOT EXISTS idx_project_members_is_active ON project_members (is_active);
CREATE INDEX IF NOT EXISTS idx_project_members_project_active ON project_members (project_id, is_active);
CREATE INDEX IF NOT EXISTS idx_project_members_employee_active ON project_members (employee_id, is_active);

-- 项目费用表
CREATE TABLE IF NOT EXISTS project_expenses (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL,
  created_by BIGINT NULL,
  updated_by BIGINT NULL,
  project_id BIGINT NOT NULL,
  expense_date TIMESTAMP NOT NULL,
  expense_type VARCHAR(100) NOT NULL,
  amount NUMERIC(15,2) NOT NULL,
  currency VARCHAR(10) DEFAULT 'CNY',
  description TEXT NOT NULL,
  receipt VARCHAR(500) NULL,
  is_reimbursable BOOLEAN DEFAULT TRUE,
  status VARCHAR(50) DEFAULT 'pending',
  approved_by BIGINT NULL,
  approved_at TIMESTAMP NULL,
  CONSTRAINT fk_project_expenses_project FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
  CONSTRAINT fk_project_expenses_approved_by FOREIGN KEY (approved_by) REFERENCES employees(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_project_expenses_deleted_at ON project_expenses (deleted_at);
CREATE INDEX IF NOT EXISTS idx_project_expenses_project_id ON project_expenses (project_id);
CREATE INDEX IF NOT EXISTS idx_project_expenses_expense_date ON project_expenses (expense_date);
CREATE INDEX IF NOT EXISTS idx_project_expenses_expense_type ON project_expenses (expense_type);
CREATE INDEX IF NOT EXISTS idx_project_expenses_is_reimbursable ON project_expenses (is_reimbursable);
CREATE INDEX IF NOT EXISTS idx_project_expenses_status ON project_expenses (status);
CREATE INDEX IF NOT EXISTS idx_project_expenses_approved_by ON project_expenses (approved_by);

-- 任务评论表
CREATE TABLE IF NOT EXISTS task_comments (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL,
  task_id BIGINT NOT NULL,
  employee_id BIGINT NOT NULL,
  comment TEXT NOT NULL,
  CONSTRAINT fk_task_comments_task FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE,
  CONSTRAINT fk_task_comments_employee FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_task_comments_deleted_at ON task_comments (deleted_at);
CREATE INDEX IF NOT EXISTS idx_task_comments_task_id ON task_comments (task_id);
CREATE INDEX IF NOT EXISTS idx_task_comments_employee_id ON task_comments (employee_id);

-- 项目文档表
CREATE TABLE IF NOT EXISTS project_documents (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL,
  project_id BIGINT NOT NULL,
  document_name VARCHAR(255) NOT NULL,
  document_type VARCHAR(100) NOT NULL,
  file_path VARCHAR(500) NOT NULL,
  file_size BIGINT DEFAULT 0,
  version VARCHAR(50) DEFAULT '1.0',
  uploaded_by BIGINT NOT NULL,
  description TEXT NULL,
  CONSTRAINT fk_project_documents_project FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
  CONSTRAINT fk_project_documents_uploaded_by FOREIGN KEY (uploaded_by) REFERENCES employees(id) ON DELETE RESTRICT
);
CREATE INDEX IF NOT EXISTS idx_project_documents_deleted_at ON project_documents (deleted_at);
CREATE INDEX IF NOT EXISTS idx_project_documents_project_id ON project_documents (project_id);
CREATE INDEX IF NOT EXISTS idx_project_documents_document_type ON project_documents (document_type);
CREATE INDEX IF NOT EXISTS idx_project_documents_uploaded_by ON project_documents (uploaded_by);

-- 项目资源表
CREATE TABLE IF NOT EXISTS project_resources (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL,
  project_id BIGINT NOT NULL,
  resource_type VARCHAR(100) NOT NULL,
  resource_id BIGINT NOT NULL,
  quantity NUMERIC(10,2) DEFAULT 1,
  unit VARCHAR(50) NULL,
  cost_per_unit NUMERIC(15,2) DEFAULT 0,
  total_cost NUMERIC(15,2) DEFAULT 0,
  allocated_at TIMESTAMP NOT NULL,
  released_at TIMESTAMP NULL,
  status VARCHAR(50) DEFAULT 'allocated',
  CONSTRAINT fk_project_resources_project FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_project_resources_deleted_at ON project_resources (deleted_at);
CREATE INDEX IF NOT EXISTS idx_project_resources_project_id ON project_resources (project_id);
CREATE INDEX IF NOT EXISTS idx_project_resources_resource_type ON project_resources (resource_type);
CREATE INDEX IF NOT EXISTS idx_project_resources_resource_id ON project_resources (resource_id);
CREATE INDEX IF NOT EXISTS idx_project_resources_allocated_at ON project_resources (allocated_at);
CREATE INDEX IF NOT EXISTS idx_project_resources_released_at ON project_resources (released_at);
CREATE INDEX IF NOT EXISTS idx_project_resources_status ON project_resources (status);

-- 项目报告表
CREATE TABLE IF NOT EXISTS project_reports (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL,
  created_by BIGINT NULL,
  updated_by BIGINT NULL,
  project_id BIGINT NOT NULL,
  report_type VARCHAR(100) NOT NULL,
  report_date TIMESTAMP NOT NULL,
  title VARCHAR(255) NOT NULL,
  content TEXT NOT NULL,
  generated_by BIGINT NOT NULL,
  status VARCHAR(50) DEFAULT 'draft',
  CONSTRAINT fk_project_reports_project FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
  CONSTRAINT fk_project_reports_generated_by FOREIGN KEY (generated_by) REFERENCES employees(id) ON DELETE RESTRICT
);
CREATE INDEX IF NOT EXISTS idx_project_reports_deleted_at ON project_reports (deleted_at);
CREATE INDEX IF NOT EXISTS idx_project_reports_project_id ON project_reports (project_id);
CREATE INDEX IF NOT EXISTS idx_project_reports_report_type ON project_reports (report_type);
CREATE INDEX IF NOT EXISTS idx_project_reports_report_date ON project_reports (report_date);
CREATE INDEX IF NOT EXISTS idx_project_reports_generated_by ON project_reports (generated_by);
CREATE INDEX IF NOT EXISTS idx_project_reports_status ON project_reports (status);

-- ============================================================================
-- 结束
-- ============================================================================