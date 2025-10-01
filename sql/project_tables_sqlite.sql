-- ============================================================================
-- GalaxyERP 项目管理模块 - SQLite 建表脚本
-- 生成时间: 2025-10-01
-- 说明: 基于当前 /internal/models 下的结构，包含项目相关表及必要依赖
-- 注意: SQLite 不支持部分高级约束，类型为宽松类型（TEXT/INTEGER/REAL）
-- ============================================================================

PRAGMA foreign_keys = ON;

-- ============================================================================
-- 基础依赖表（部门、职位、员工、客户）
-- ============================================================================

-- 部门表
CREATE TABLE IF NOT EXISTS departments (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME NULL,
  updated_at DATETIME NULL,
  deleted_at DATETIME NULL,
  name TEXT NOT NULL,
  code TEXT NOT NULL,
  description TEXT NULL,
  company_id INTEGER NOT NULL,
  parent_id INTEGER NULL,
  is_active INTEGER DEFAULT 1,
  created_by INTEGER NULL,
  updated_by INTEGER NULL,
  UNIQUE (code),
  FOREIGN KEY (parent_id) REFERENCES departments(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_departments_deleted_at ON departments (deleted_at);
CREATE INDEX IF NOT EXISTS idx_departments_company_id ON departments (company_id);
CREATE INDEX IF NOT EXISTS idx_departments_parent_id ON departments (parent_id);
CREATE INDEX IF NOT EXISTS idx_departments_is_active ON departments (is_active);

-- 职位表
CREATE TABLE IF NOT EXISTS positions (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME NULL,
  updated_at DATETIME NULL,
  deleted_at DATETIME NULL,
  created_by INTEGER NULL,
  updated_by INTEGER NULL,
  is_active INTEGER DEFAULT 1,
  code TEXT NOT NULL,
  name TEXT NOT NULL,
  description TEXT NULL,
  department_id INTEGER NOT NULL,
  UNIQUE (code),
  FOREIGN KEY (department_id) REFERENCES departments(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_positions_deleted_at ON positions (deleted_at);
CREATE INDEX IF NOT EXISTS idx_positions_is_active ON positions (is_active);
CREATE INDEX IF NOT EXISTS idx_positions_department_id ON positions (department_id);

-- 员工表
CREATE TABLE IF NOT EXISTS employees (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME NULL,
  updated_at DATETIME NULL,
  deleted_at DATETIME NULL,
  code TEXT NOT NULL,
  first_name TEXT NOT NULL,
  last_name TEXT NOT NULL,
  full_name TEXT NOT NULL,
  email TEXT NOT NULL,
  phone TEXT NULL,
  date_of_birth DATETIME NULL,
  gender TEXT NULL,
  hire_date DATETIME NULL,
  department_id INTEGER NULL,
  position_id INTEGER NULL,
  manager_id INTEGER NULL,
  status TEXT DEFAULT 'active',
  emergency_contact TEXT NULL,
  id_number TEXT NULL,
  address TEXT NULL,
  UNIQUE (code),
  UNIQUE (email),
  FOREIGN KEY (department_id) REFERENCES departments(id) ON DELETE SET NULL,
  FOREIGN KEY (position_id) REFERENCES positions(id) ON DELETE SET NULL,
  FOREIGN KEY (manager_id) REFERENCES employees(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_employees_deleted_at ON employees (deleted_at);
CREATE INDEX IF NOT EXISTS idx_employees_department_id ON employees (department_id);
CREATE INDEX IF NOT EXISTS idx_employees_position_id ON employees (position_id);
CREATE INDEX IF NOT EXISTS idx_employees_manager_id ON employees (manager_id);
CREATE INDEX IF NOT EXISTS idx_employees_status ON employees (status);

-- 客户表
CREATE TABLE IF NOT EXISTS customers (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME NULL,
  updated_at DATETIME NULL,
  deleted_at DATETIME NULL,
  name TEXT NOT NULL,
  code TEXT NOT NULL,
  email TEXT NULL,
  phone TEXT NULL,
  address TEXT NULL,
  city TEXT NULL,
  state TEXT NULL,
  postal_code TEXT NULL,
  country TEXT NULL,
  contact_person TEXT NULL,
  credit_limit REAL DEFAULT 0,
  customer_group TEXT NULL,
  territory TEXT NULL,
  is_active INTEGER DEFAULT 1,
  UNIQUE (code)
);
CREATE INDEX IF NOT EXISTS idx_customers_deleted_at ON customers (deleted_at);
CREATE INDEX IF NOT EXISTS idx_customers_is_active ON customers (is_active);

-- ============================================================================
-- 项目管理核心表
-- ============================================================================

-- 项目表
CREATE TABLE IF NOT EXISTS projects (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME NULL,
  updated_at DATETIME NULL,
  deleted_at DATETIME NULL,
  created_by INTEGER NULL,
  updated_by INTEGER NULL,
  project_number TEXT NOT NULL,
  project_name TEXT NOT NULL,
  description TEXT NULL,
  customer_id INTEGER NULL,
  manager_id INTEGER NOT NULL,
  start_date DATETIME NOT NULL,
  end_date DATETIME NOT NULL,
  actual_start_date DATETIME NULL,
  actual_end_date DATETIME NULL,
  budget REAL DEFAULT 0,
  actual_cost REAL DEFAULT 0,
  progress REAL DEFAULT 0,
  priority TEXT DEFAULT 'normal',
  status TEXT DEFAULT 'planning',
  UNIQUE (project_number),
  FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE SET NULL,
  FOREIGN KEY (manager_id) REFERENCES employees(id) ON DELETE RESTRICT
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
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME NULL,
  updated_at DATETIME NULL,
  deleted_at DATETIME NULL,
  created_by INTEGER NULL,
  updated_by INTEGER NULL,
  task_number TEXT NOT NULL,
  task_name TEXT NOT NULL,
  description TEXT NULL,
  project_id INTEGER NOT NULL,
  parent_task_id INTEGER NULL,
  assignee_id INTEGER NULL,
  start_date DATETIME NOT NULL,
  end_date DATETIME NOT NULL,
  actual_start_date DATETIME NULL,
  actual_end_date DATETIME NULL,
  estimated_hours REAL DEFAULT 0,
  actual_hours REAL DEFAULT 0,
  progress REAL DEFAULT 0,
  priority TEXT DEFAULT 'normal',
  status TEXT DEFAULT 'todo',
  UNIQUE (task_number),
  FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
  FOREIGN KEY (parent_task_id) REFERENCES tasks(id) ON DELETE SET NULL,
  FOREIGN KEY (assignee_id) REFERENCES employees(id) ON DELETE SET NULL
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
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME NULL,
  updated_at DATETIME NULL,
  deleted_at DATETIME NULL,
  project_id INTEGER NOT NULL,
  milestone_name TEXT NOT NULL,
  description TEXT NULL,
  due_date DATETIME NOT NULL,
  completed_date DATETIME NULL,
  status TEXT DEFAULT 'pending',
  FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_milestones_deleted_at ON milestones (deleted_at);
CREATE INDEX IF NOT EXISTS idx_milestones_project_id ON milestones (project_id);
CREATE INDEX IF NOT EXISTS idx_milestones_due_date ON milestones (due_date);
CREATE INDEX IF NOT EXISTS idx_milestones_status ON milestones (status);

-- 工时记录表
CREATE TABLE IF NOT EXISTS time_entries (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME NULL,
  updated_at DATETIME NULL,
  deleted_at DATETIME NULL,
  employee_id INTEGER NOT NULL,
  project_id INTEGER NOT NULL,
  task_id INTEGER NULL,
  date DATETIME NOT NULL,
  start_time DATETIME NOT NULL,
  end_time DATETIME NOT NULL,
  hours REAL NOT NULL,
  description TEXT NULL,
  is_billable INTEGER DEFAULT 1,
  hourly_rate REAL DEFAULT 0,
  amount REAL DEFAULT 0,
  status TEXT DEFAULT 'draft',
  FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE CASCADE,
  FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
  FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE SET NULL
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
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME NULL,
  updated_at DATETIME NULL,
  deleted_at DATETIME NULL,
  project_id INTEGER NOT NULL,
  employee_id INTEGER NOT NULL,
  role TEXT NOT NULL,
  join_date DATETIME NOT NULL,
  leave_date DATETIME NULL,
  is_active INTEGER DEFAULT 1,
  FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
  FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE CASCADE
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
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME NULL,
  updated_at DATETIME NULL,
  deleted_at DATETIME NULL,
  created_by INTEGER NULL,
  updated_by INTEGER NULL,
  project_id INTEGER NOT NULL,
  expense_date DATETIME NOT NULL,
  expense_type TEXT NOT NULL,
  amount REAL NOT NULL,
  currency TEXT DEFAULT 'CNY',
  description TEXT NOT NULL,
  receipt TEXT NULL,
  is_reimbursable INTEGER DEFAULT 1,
  status TEXT DEFAULT 'pending',
  approved_by INTEGER NULL,
  approved_at DATETIME NULL,
  FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
  FOREIGN KEY (approved_by) REFERENCES employees(id) ON DELETE SET NULL
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
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME NULL,
  updated_at DATETIME NULL,
  deleted_at DATETIME NULL,
  task_id INTEGER NOT NULL,
  employee_id INTEGER NOT NULL,
  comment TEXT NOT NULL,
  FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE,
  FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_task_comments_deleted_at ON task_comments (deleted_at);
CREATE INDEX IF NOT EXISTS idx_task_comments_task_id ON task_comments (task_id);
CREATE INDEX IF NOT EXISTS idx_task_comments_employee_id ON task_comments (employee_id);

-- 项目文档表
CREATE TABLE IF NOT EXISTS project_documents (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME NULL,
  updated_at DATETIME NULL,
  deleted_at DATETIME NULL,
  project_id INTEGER NOT NULL,
  document_name TEXT NOT NULL,
  document_type TEXT NOT NULL,
  file_path TEXT NOT NULL,
  file_size INTEGER DEFAULT 0,
  version TEXT DEFAULT '1.0',
  uploaded_by INTEGER NOT NULL,
  description TEXT NULL,
  FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
  FOREIGN KEY (uploaded_by) REFERENCES employees(id) ON DELETE RESTRICT
);
CREATE INDEX IF NOT EXISTS idx_project_documents_deleted_at ON project_documents (deleted_at);
CREATE INDEX IF NOT EXISTS idx_project_documents_project_id ON project_documents (project_id);
CREATE INDEX IF NOT EXISTS idx_project_documents_document_type ON project_documents (document_type);
CREATE INDEX IF NOT EXISTS idx_project_documents_uploaded_by ON project_documents (uploaded_by);

-- 项目资源表
CREATE TABLE IF NOT EXISTS project_resources (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME NULL,
  updated_at DATETIME NULL,
  deleted_at DATETIME NULL,
  project_id INTEGER NOT NULL,
  resource_type TEXT NOT NULL,
  resource_id INTEGER NOT NULL,
  quantity REAL DEFAULT 1,
  unit TEXT NULL,
  cost_per_unit REAL DEFAULT 0,
  total_cost REAL DEFAULT 0,
  allocated_at DATETIME NOT NULL,
  released_at DATETIME NULL,
  status TEXT DEFAULT 'allocated',
  FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
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
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at DATETIME NULL,
  updated_at DATETIME NULL,
  deleted_at DATETIME NULL,
  created_by INTEGER NULL,
  updated_by INTEGER NULL,
  project_id INTEGER NOT NULL,
  report_type TEXT NOT NULL,
  report_date DATETIME NOT NULL,
  title TEXT NOT NULL,
  content TEXT NOT NULL,
  generated_by INTEGER NOT NULL,
  status TEXT DEFAULT 'draft',
  FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
  FOREIGN KEY (generated_by) REFERENCES employees(id) ON DELETE RESTRICT
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