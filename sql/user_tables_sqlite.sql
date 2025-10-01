-- ============================================================================
-- GalaxyERP 用户管理模块 - SQLite 建表脚本
-- 生成时间: 2025-01-01
-- 说明: 基于 /internal/models/user.go 的结构，包含用户、角色、权限管理
-- ============================================================================

-- ============================================================================
-- 组织架构管理
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
  legal_name TEXT,
  tax_number TEXT,
  registration_number TEXT,
  address TEXT,
  phone TEXT,
  email TEXT,
  website TEXT,
  logo_url TEXT,
  status TEXT DEFAULT 'ACTIVE',
  CONSTRAINT uq_companies_code UNIQUE (code)
);
CREATE INDEX IF NOT EXISTS idx_companies_deleted_at ON companies (deleted_at);
CREATE INDEX IF NOT EXISTS idx_companies_is_active ON companies (is_active);
CREATE INDEX IF NOT EXISTS idx_companies_code ON companies (code);
CREATE INDEX IF NOT EXISTS idx_companies_name ON companies (name);
CREATE INDEX IF NOT EXISTS idx_companies_status ON companies (status);

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
  CONSTRAINT fk_departments_parent FOREIGN KEY (parent_id) REFERENCES departments(id) ON DELETE SET NULL,
  CONSTRAINT fk_departments_manager FOREIGN KEY (manager_id) REFERENCES users(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_departments_deleted_at ON departments (deleted_at);
CREATE INDEX IF NOT EXISTS idx_departments_is_active ON departments (is_active);
CREATE INDEX IF NOT EXISTS idx_departments_code ON departments (code);
CREATE INDEX IF NOT EXISTS idx_departments_name ON departments (name);
CREATE INDEX IF NOT EXISTS idx_departments_company_id ON departments (company_id);
CREATE INDEX IF NOT EXISTS idx_departments_parent_id ON departments (parent_id);
CREATE INDEX IF NOT EXISTS idx_departments_manager_id ON departments (manager_id);
CREATE INDEX IF NOT EXISTS idx_departments_level ON departments (level);
CREATE INDEX IF NOT EXISTS idx_departments_status ON departments (status);

-- ============================================================================
-- 用户管理
-- ============================================================================

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
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users (deleted_at);
CREATE INDEX IF NOT EXISTS idx_users_is_active ON users (is_active);
CREATE INDEX IF NOT EXISTS idx_users_username ON users (username);
CREATE INDEX IF NOT EXISTS idx_users_email ON users (email);
CREATE INDEX IF NOT EXISTS idx_users_phone ON users (phone);
CREATE INDEX IF NOT EXISTS idx_users_company_id ON users (company_id);
CREATE INDEX IF NOT EXISTS idx_users_department_id ON users (department_id);
CREATE INDEX IF NOT EXISTS idx_users_employee_number ON users (employee_number);
CREATE INDEX IF NOT EXISTS idx_users_status ON users (status);
CREATE INDEX IF NOT EXISTS idx_users_last_login_at ON users (last_login_at);

-- ============================================================================
-- 权限管理
-- ============================================================================

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
CREATE INDEX IF NOT EXISTS idx_roles_deleted_at ON roles (deleted_at);
CREATE INDEX IF NOT EXISTS idx_roles_is_active ON roles (is_active);
CREATE INDEX IF NOT EXISTS idx_roles_code ON roles (code);
CREATE INDEX IF NOT EXISTS idx_roles_name ON roles (name);
CREATE INDEX IF NOT EXISTS idx_roles_role_type ON roles (role_type);
CREATE INDEX IF NOT EXISTS idx_roles_is_system ON roles (is_system);
CREATE INDEX IF NOT EXISTS idx_roles_status ON roles (status);

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
CREATE INDEX IF NOT EXISTS idx_permissions_deleted_at ON permissions (deleted_at);
CREATE INDEX IF NOT EXISTS idx_permissions_is_active ON permissions (is_active);
CREATE INDEX IF NOT EXISTS idx_permissions_code ON permissions (code);
CREATE INDEX IF NOT EXISTS idx_permissions_name ON permissions (name);
CREATE INDEX IF NOT EXISTS idx_permissions_resource ON permissions (resource);
CREATE INDEX IF NOT EXISTS idx_permissions_action ON permissions (action);
CREATE INDEX IF NOT EXISTS idx_permissions_permission_type ON permissions (permission_type);
CREATE INDEX IF NOT EXISTS idx_permissions_parent_id ON permissions (parent_id);
CREATE INDEX IF NOT EXISTS idx_permissions_is_system ON permissions (is_system);
CREATE INDEX IF NOT EXISTS idx_permissions_status ON permissions (status);

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
CREATE INDEX IF NOT EXISTS idx_user_roles_deleted_at ON user_roles (deleted_at);
CREATE INDEX IF NOT EXISTS idx_user_roles_user_id ON user_roles (user_id);
CREATE INDEX IF NOT EXISTS idx_user_roles_role_id ON user_roles (role_id);
CREATE INDEX IF NOT EXISTS idx_user_roles_assigned_by ON user_roles (assigned_by);
CREATE INDEX IF NOT EXISTS idx_user_roles_is_active ON user_roles (is_active);
CREATE INDEX IF NOT EXISTS idx_user_roles_expires_at ON user_roles (expires_at);

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
CREATE INDEX IF NOT EXISTS idx_role_permissions_deleted_at ON role_permissions (deleted_at);
CREATE INDEX IF NOT EXISTS idx_role_permissions_role_id ON role_permissions (role_id);
CREATE INDEX IF NOT EXISTS idx_role_permissions_permission_id ON role_permissions (permission_id);
CREATE INDEX IF NOT EXISTS idx_role_permissions_assigned_by ON role_permissions (assigned_by);
CREATE INDEX IF NOT EXISTS idx_role_permissions_is_active ON role_permissions (is_active);

-- ============================================================================
-- 会话管理
-- ============================================================================

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
CREATE INDEX IF NOT EXISTS idx_user_sessions_deleted_at ON user_sessions (deleted_at);
CREATE INDEX IF NOT EXISTS idx_user_sessions_user_id ON user_sessions (user_id);
CREATE INDEX IF NOT EXISTS idx_user_sessions_session_id ON user_sessions (session_id);
CREATE INDEX IF NOT EXISTS idx_user_sessions_ip_address ON user_sessions (ip_address);
CREATE INDEX IF NOT EXISTS idx_user_sessions_login_at ON user_sessions (login_at);
CREATE INDEX IF NOT EXISTS idx_user_sessions_last_activity_at ON user_sessions (last_activity_at);
CREATE INDEX IF NOT EXISTS idx_user_sessions_expires_at ON user_sessions (expires_at);
CREATE INDEX IF NOT EXISTS idx_user_sessions_is_active ON user_sessions (is_active);

-- ============================================================================
-- 初始数据插入
-- ============================================================================

-- 插入默认公司
INSERT INTO companies (code, name, description, legal_name, address, phone, email, status) VALUES
('GALAXY', 'Galaxy Corporation', '银河集团总公司', '银河科技有限公司', '北京市朝阳区银河大厦', '010-12345678', 'info@galaxy.com', 'ACTIVE');

-- 插入默认部门
INSERT INTO departments (code, name, description, company_id, level, sort_order, status) VALUES
('ROOT', '总公司', '集团总部', 1, 1, 1, 'ACTIVE'),
('IT', '信息技术部', '负责公司IT系统建设和维护', 1, 2, 10, 'ACTIVE'),
('HR', '人力资源部', '负责人力资源管理', 1, 2, 20, 'ACTIVE'),
('FINANCE', '财务部', '负责财务管理', 1, 2, 30, 'ACTIVE'),
('SALES', '销售部', '负责销售业务', 1, 2, 40, 'ACTIVE'),
('PURCHASE', '采购部', '负责采购业务', 1, 2, 50, 'ACTIVE');

-- 插入默认用户
INSERT INTO users (username, email, password_hash, first_name, last_name, display_name, company_id, department_id, position, employee_number, status) VALUES
('admin', 'admin@galaxy.com', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iKWKQK5OzniLor.VgUxRK.cKIm2y', 'System', 'Admin', '系统管理员', 1, 1, '系统管理员', 'EMP001', 'ACTIVE'),
('demo', 'demo@galaxy.com', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iKWKQK5OzniLor.VgUxRK.cKIm2y', 'Demo', 'User', '演示用户', 1, 2, '开发工程师', 'EMP002', 'ACTIVE');

-- 插入默认角色
INSERT INTO roles (code, name, description, role_type, is_system, sort_order, status) VALUES
('SUPER_ADMIN', '超级管理员', '系统超级管理员，拥有所有权限', 'SYSTEM', 1, 1, 'ACTIVE'),
('ADMIN', '系统管理员', '系统管理员，拥有大部分管理权限', 'SYSTEM', 1, 2, 'ACTIVE'),
('USER', '普通用户', '普通用户，拥有基本操作权限', 'SYSTEM', 1, 10, 'ACTIVE'),
('GUEST', '访客', '访客用户，只有查看权限', 'SYSTEM', 1, 20, 'ACTIVE');

-- 插入默认权限
INSERT INTO permissions (code, name, description, resource, action, permission_type, is_system, sort_order, status) VALUES
-- 系统管理
('SYSTEM_MANAGE', '系统管理', '系统管理权限', 'SYSTEM', 'MANAGE', 'MODULE', 1, 1, 'ACTIVE'),
('USER_MANAGE', '用户管理', '用户管理权限', 'USER', 'MANAGE', 'MODULE', 1, 10, 'ACTIVE'),
('USER_CREATE', '创建用户', '创建用户权限', 'USER', 'CREATE', 'FUNCTION', 1, 11, 'ACTIVE'),
('USER_UPDATE', '修改用户', '修改用户权限', 'USER', 'UPDATE', 'FUNCTION', 1, 12, 'ACTIVE'),
('USER_DELETE', '删除用户', '删除用户权限', 'USER', 'DELETE', 'FUNCTION', 1, 13, 'ACTIVE'),
('USER_VIEW', '查看用户', '查看用户权限', 'USER', 'VIEW', 'FUNCTION', 1, 14, 'ACTIVE'),
-- 角色管理
('ROLE_MANAGE', '角色管理', '角色管理权限', 'ROLE', 'MANAGE', 'MODULE', 1, 20, 'ACTIVE'),
('ROLE_CREATE', '创建角色', '创建角色权限', 'ROLE', 'CREATE', 'FUNCTION', 1, 21, 'ACTIVE'),
('ROLE_UPDATE', '修改角色', '修改角色权限', 'ROLE', 'UPDATE', 'FUNCTION', 1, 22, 'ACTIVE'),
('ROLE_DELETE', '删除角色', '删除角色权限', 'ROLE', 'DELETE', 'FUNCTION', 1, 23, 'ACTIVE'),
('ROLE_VIEW', '查看角色', '查看角色权限', 'ROLE', 'VIEW', 'FUNCTION', 1, 24, 'ACTIVE'),
-- 权限管理
('PERMISSION_MANAGE', '权限管理', '权限管理权限', 'PERMISSION', 'MANAGE', 'MODULE', 1, 30, 'ACTIVE'),
('PERMISSION_VIEW', '查看权限', '查看权限权限', 'PERMISSION', 'VIEW', 'FUNCTION', 1, 31, 'ACTIVE');

-- 插入用户角色关联
INSERT INTO user_roles (user_id, role_id, assigned_by, assigned_at, is_active) VALUES
(1, 1, 1, datetime('now'), 1),  -- admin 用户分配超级管理员角色
(2, 3, 1, datetime('now'), 1);  -- demo 用户分配普通用户角色

-- 插入角色权限关联
INSERT INTO role_permissions (role_id, permission_id, assigned_by, assigned_at, is_active) VALUES
-- 超级管理员拥有所有权限
(1, 1, 1, datetime('now'), 1),
(1, 2, 1, datetime('now'), 1),
(1, 3, 1, datetime('now'), 1),
(1, 4, 1, datetime('now'), 1),
(1, 5, 1, datetime('now'), 1),
(1, 6, 1, datetime('now'), 1),
(1, 7, 1, datetime('now'), 1),
(1, 8, 1, datetime('now'), 1),
(1, 9, 1, datetime('now'), 1),
(1, 10, 1, datetime('now'), 1),
(1, 11, 1, datetime('now'), 1),
(1, 12, 1, datetime('now'), 1),
(1, 13, 1, datetime('now'), 1),
-- 普通用户只有查看权限
(3, 6, 1, datetime('now'), 1),
(3, 11, 1, datetime('now'), 1),
(3, 13, 1, datetime('now'), 1);

-- ============================================================================
-- 结束
-- ============================================================================