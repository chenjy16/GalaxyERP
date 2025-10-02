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

-- 项目表
CREATE TABLE IF NOT EXISTS projects (
  id SERIAL PRIMARY KEY,
  project_number VARCHAR(100) NOT NULL UNIQUE,
  project_name VARCHAR(255) NOT NULL,
  description TEXT,
  customer_id INTEGER,
  manager_id INTEGER NOT NULL,
  start_date TIMESTAMP WITH TIME ZONE NOT NULL,
  end_date TIMESTAMP WITH TIME ZONE NOT NULL,
  actual_start_date TIMESTAMP WITH TIME ZONE,
  actual_end_date TIMESTAMP WITH TIME ZONE,
  budget DECIMAL(15,2) DEFAULT 0,
  actual_cost DECIMAL(15,2) DEFAULT 0,
  progress DECIMAL(5,2) DEFAULT 0,
  priority VARCHAR(20) DEFAULT 'normal',
  status VARCHAR(50) DEFAULT 'planning',
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  created_by INTEGER,
  updated_by INTEGER,
  deleted_at TIMESTAMP WITH TIME ZONE,
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
  id SERIAL PRIMARY KEY,
  task_number VARCHAR(100) NOT NULL UNIQUE,
  task_name VARCHAR(255) NOT NULL,
  description TEXT,
  project_id INTEGER NOT NULL,
  parent_task_id INTEGER,
  assignee_id INTEGER,
  start_date TIMESTAMP WITH TIME ZONE NOT NULL,
  end_date TIMESTAMP WITH TIME ZONE NOT NULL,
  actual_start_date TIMESTAMP WITH TIME ZONE,
  actual_end_date TIMESTAMP WITH TIME ZONE,
  estimated_hours DECIMAL(8,2) DEFAULT 0,
  actual_hours DECIMAL(8,2) DEFAULT 0,
  progress DECIMAL(5,2) DEFAULT 0,
  priority VARCHAR(20) DEFAULT 'normal',
  status VARCHAR(50) DEFAULT 'todo',
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  created_by INTEGER,
  updated_by INTEGER,
  deleted_at TIMESTAMP WITH TIME ZONE,
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

-- 工时记录表
CREATE TABLE time_entries (
    id SERIAL PRIMARY KEY,
    employee_id INTEGER NOT NULL REFERENCES employees(id),
    project_id INTEGER NOT NULL REFERENCES projects(id),
    task_id INTEGER REFERENCES tasks(id),
    date TIMESTAMP WITH TIME ZONE NOT NULL,
    start_time TIMESTAMP WITH TIME ZONE NOT NULL,
    end_time TIMESTAMP WITH TIME ZONE NOT NULL,
    hours DECIMAL(8,2) NOT NULL,
    description TEXT,
    is_billable BOOLEAN DEFAULT true,
    hourly_rate DECIMAL(10,2) DEFAULT 0,
    amount DECIMAL(10,2) DEFAULT 0,
    status VARCHAR(50) DEFAULT 'draft',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- 工时记录表索引
CREATE INDEX idx_time_entries_employee_id ON time_entries(employee_id);
CREATE INDEX idx_time_entries_project_id ON time_entries(project_id);
CREATE INDEX idx_time_entries_task_id ON time_entries(task_id);
CREATE INDEX idx_time_entries_date ON time_entries(date);
CREATE INDEX idx_time_entries_is_billable ON time_entries(is_billable);
CREATE INDEX idx_time_entries_status ON time_entries(status);
CREATE INDEX idx_time_entries_deleted_at ON time_entries(deleted_at);

-- ============================================================================
-- 项目相关辅助表
-- ============================================================================

-- 项目成员表
CREATE TABLE project_members (
    id SERIAL PRIMARY KEY,
    project_id INTEGER NOT NULL REFERENCES projects(id),
    employee_id INTEGER NOT NULL REFERENCES employees(id),
    role VARCHAR(100) NOT NULL,
    join_date TIMESTAMP WITH TIME ZONE NOT NULL,
    leave_date TIMESTAMP WITH TIME ZONE,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(project_id, employee_id)
);

-- 项目成员表索引
CREATE INDEX idx_project_members_project_id ON project_members(project_id);
CREATE INDEX idx_project_members_employee_id ON project_members(employee_id);
CREATE INDEX idx_project_members_join_date ON project_members(join_date);
CREATE INDEX idx_project_members_leave_date ON project_members(leave_date);
CREATE INDEX idx_project_members_is_active ON project_members(is_active);
CREATE INDEX idx_project_members_deleted_at ON project_members(deleted_at);

-- 项目费用表
CREATE TABLE project_expenses (
    id SERIAL PRIMARY KEY,
    project_id INTEGER NOT NULL REFERENCES projects(id),
    expense_date TIMESTAMP WITH TIME ZONE NOT NULL,
    expense_type VARCHAR(100) NOT NULL,
    amount DECIMAL(15,2) NOT NULL,
    currency VARCHAR(10) DEFAULT 'CNY',
    description TEXT NOT NULL,
    receipt VARCHAR(500),
    is_reimbursable BOOLEAN DEFAULT true,
    status VARCHAR(50) DEFAULT 'pending',
    approved_by INTEGER REFERENCES employees(id),
    approved_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    created_by INTEGER,
    updated_by INTEGER
);

-- 项目费用表索引
CREATE INDEX idx_project_expenses_project_id ON project_expenses(project_id);
CREATE INDEX idx_project_expenses_expense_date ON project_expenses(expense_date);
CREATE INDEX idx_project_expenses_expense_type ON project_expenses(expense_type);
CREATE INDEX idx_project_expenses_is_reimbursable ON project_expenses(is_reimbursable);
CREATE INDEX idx_project_expenses_status ON project_expenses(status);
CREATE INDEX idx_project_expenses_approved_by ON project_expenses(approved_by);
CREATE INDEX idx_project_expenses_deleted_at ON project_expenses(deleted_at);
CREATE INDEX idx_project_expenses_created_by ON project_expenses(created_by);
CREATE INDEX idx_project_expenses_updated_by ON project_expenses(updated_by);

-- 任务评论表
CREATE TABLE task_comments (
    id SERIAL PRIMARY KEY,
    task_id INTEGER NOT NULL REFERENCES tasks(id),
    employee_id INTEGER NOT NULL REFERENCES employees(id),
    comment TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- 任务评论表索引
CREATE INDEX idx_task_comments_task_id ON task_comments(task_id);
CREATE INDEX idx_task_comments_employee_id ON task_comments(employee_id);
CREATE INDEX idx_task_comments_deleted_at ON task_comments(deleted_at);

-- 项目文档表
CREATE TABLE project_documents (
    id SERIAL PRIMARY KEY,
    project_id INTEGER NOT NULL REFERENCES projects(id),
    document_name VARCHAR(255) NOT NULL,
    document_type VARCHAR(100) NOT NULL,
    file_path VARCHAR(500) NOT NULL,
    file_size BIGINT DEFAULT 0,
    version VARCHAR(50) DEFAULT '1.0',
    uploaded_by INTEGER NOT NULL REFERENCES employees(id),
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- 项目文档表索引
CREATE INDEX idx_project_documents_project_id ON project_documents(project_id);
CREATE INDEX idx_project_documents_document_type ON project_documents(document_type);
CREATE INDEX idx_project_documents_uploaded_by ON project_documents(uploaded_by);
CREATE INDEX idx_project_documents_deleted_at ON project_documents(deleted_at);

-- 项目资源表
CREATE TABLE project_resources (
    id SERIAL PRIMARY KEY,
    project_id INTEGER NOT NULL REFERENCES projects(id),
    resource_type VARCHAR(100) NOT NULL,
    resource_id INTEGER NOT NULL,
    quantity DECIMAL(15,4) DEFAULT 1,
    unit VARCHAR(50),
    cost_per_unit DECIMAL(15,2) DEFAULT 0,
    total_cost DECIMAL(15,2) DEFAULT 0,
    allocated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    released_at TIMESTAMP WITH TIME ZONE,
    status VARCHAR(50) DEFAULT 'allocated',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- 项目资源表索引
CREATE INDEX idx_project_resources_project_id ON project_resources(project_id);
CREATE INDEX idx_project_resources_resource_type ON project_resources(resource_type);
CREATE INDEX idx_project_resources_resource_id ON project_resources(resource_id);
CREATE INDEX idx_project_resources_allocated_at ON project_resources(allocated_at);
CREATE INDEX idx_project_resources_released_at ON project_resources(released_at);
CREATE INDEX idx_project_resources_status ON project_resources(status);
CREATE INDEX idx_project_resources_deleted_at ON project_resources(deleted_at);

-- 项目报告表
CREATE TABLE project_reports (
    id SERIAL PRIMARY KEY,
    project_id INTEGER NOT NULL REFERENCES projects(id),
    report_type VARCHAR(100) NOT NULL,
    report_date TIMESTAMP WITH TIME ZONE NOT NULL,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    generated_by INTEGER NOT NULL REFERENCES employees(id),
    status VARCHAR(50) DEFAULT 'draft',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    created_by INTEGER,
    updated_by INTEGER
);

-- 项目报告表索引
CREATE INDEX idx_project_reports_project_id ON project_reports(project_id);
CREATE INDEX idx_project_reports_report_type ON project_reports(report_type);
CREATE INDEX idx_project_reports_report_date ON project_reports(report_date);
CREATE INDEX idx_project_reports_generated_by ON project_reports(generated_by);
CREATE INDEX idx_project_reports_status ON project_reports(status);
CREATE INDEX idx_project_reports_deleted_at ON project_reports(deleted_at);
CREATE INDEX idx_project_reports_created_by ON project_reports(created_by);
CREATE INDEX idx_project_reports_updated_by ON project_reports(updated_by);

-- 里程碑表
CREATE TABLE milestones (
    id SERIAL PRIMARY KEY,
    project_id INTEGER NOT NULL REFERENCES projects(id),
    milestone_name VARCHAR(255) NOT NULL,
    description TEXT,
    due_date TIMESTAMP WITH TIME ZONE NOT NULL,
    completed_date TIMESTAMP WITH TIME ZONE,
    status VARCHAR(50) DEFAULT 'pending',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- 里程碑表索引
CREATE INDEX idx_milestones_project_id ON milestones(project_id);
CREATE INDEX idx_milestones_due_date ON milestones(due_date);
CREATE INDEX idx_milestones_status ON milestones(status);
CREATE INDEX idx_milestones_deleted_at ON milestones(deleted_at);

-- ============================================================================
-- 示例数据
-- ============================================================================

-- 插入示例项目数据
INSERT INTO projects (project_number, project_name, description, customer_id, manager_id, start_date, end_date, actual_start_date, budget, priority, status, created_by, updated_by) VALUES
('PRJ001', 'ERP系统开发', 'Galaxy ERP系统开发项目', 1, 1, '2024-01-01 00:00:00+00', '2024-12-31 23:59:59+00', '2024-01-01 09:00:00+00', 500000.00, 'high', 'active', 1, 1),
('PRJ002', '移动应用开发', '企业移动应用开发项目', 2, 2, '2024-02-01 00:00:00+00', '2024-08-31 23:59:59+00', NULL, 200000.00, 'normal', 'planning', 2, 2),
('PRJ003', '数据迁移项目', '旧系统数据迁移到新系统', 1, 1, '2024-03-01 00:00:00+00', '2024-06-30 23:59:59+00', NULL, 100000.00, 'high', 'planning', 1, 1);

-- 插入示例任务数据
INSERT INTO tasks (task_number, task_name, description, project_id, assignee_id, start_date, end_date, actual_start_date, actual_end_date, estimated_hours, actual_hours, progress, priority, status, created_by, updated_by) VALUES
('TSK001', '需求分析', '系统需求分析和设计', 1, 1, '2024-01-01 09:00:00+00', '2024-01-31 18:00:00+00', '2024-01-01 09:00:00+00', '2024-01-30 17:00:00+00', 160.00, 155.50, 100.00, 'high', 'done', 1, 1),
('TSK002', '数据库设计', '设计系统数据库结构', 1, 2, '2024-02-01 09:00:00+00', '2024-02-28 18:00:00+00', '2024-02-01 09:00:00+00', NULL, 120.00, 80.00, 75.00, 'high', 'in_progress', 1, 1),
('TSK003', '前端开发', '用户界面开发', 1, 3, '2024-03-01 09:00:00+00', '2024-06-30 18:00:00+00', NULL, NULL, 480.00, 0.00, 0.00, 'normal', 'todo', 1, 1),
('TSK004', '后端开发', '服务端API开发', 1, 1, '2024-03-01 09:00:00+00', '2024-08-31 18:00:00+00', NULL, NULL, 640.00, 0.00, 0.00, 'high', 'todo', 1, 1);

-- 插入示例项目成员数据
INSERT INTO project_members (project_id, employee_id, role, join_date, is_active) VALUES
(1, 1, '项目经理', '2024-01-01 09:00:00+00', true),
(1, 2, '数据库工程师', '2024-01-01 09:00:00+00', true),
(1, 3, '前端工程师', '2024-03-01 09:00:00+00', true),
(2, 2, '项目经理', '2024-02-01 09:00:00+00', true);

-- 插入示例工时记录数据
INSERT INTO time_entries (employee_id, project_id, task_id, date, start_time, end_time, hours, description, is_billable, hourly_rate, amount, status) VALUES
(1, 1, 1, '2024-01-15 00:00:00+00', '2024-01-15 09:00:00+00', '2024-01-15 17:00:00+00', 8.00, '需求分析会议', true, 200.00, 1600.00, 'approved'),
(1, 1, 1, '2024-01-16 00:00:00+00', '2024-01-16 09:00:00+00', '2024-01-16 15:00:00+00', 6.00, '编写需求文档', true, 200.00, 1200.00, 'approved'),
(2, 1, 2, '2024-02-05 00:00:00+00', '2024-02-05 09:00:00+00', '2024-02-05 17:00:00+00', 8.00, '数据库表设计', true, 150.00, 1200.00, 'submitted'),
(2, 1, 2, '2024-02-06 00:00:00+00', '2024-02-06 09:00:00+00', '2024-02-06 13:00:00+00', 4.00, '数据库关系设计', true, 150.00, 600.00, 'submitted');

-- ============================================================================
-- 结束
-- ============================================================================