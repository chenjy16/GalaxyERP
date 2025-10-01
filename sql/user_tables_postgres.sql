-- ============================================================================
-- GalaxyERP 用户管理模块 - PostgreSQL 建表脚本
-- 生成时间: 2025-01-01
-- 说明: 基于 /internal/models/user.go 的结构，包含用户、角色、权限管理
-- ============================================================================

-- ============================================================================
-- 基础依赖表（公司）
-- ============================================================================

-- 公司表
CREATE TABLE IF NOT EXISTS companies (
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
  address TEXT NULL,
  phone VARCHAR(50) NULL,
  email VARCHAR(255) NULL,
  website VARCHAR(255) NULL,
  tax_number VARCHAR(100) NULL,
  registration_number VARCHAR(100) NULL,
  CONSTRAINT uq_companies_code UNIQUE (code)
);
CREATE INDEX IF NOT EXISTS idx_companies_deleted_at ON companies (deleted_at);
CREATE INDEX IF NOT EXISTS idx_companies_is_active ON companies (is_active);
CREATE INDEX IF NOT EXISTS idx_companies_code ON companies (code);

-- ============================================================================
-- 用户管理核心表
-- ============================================================================

-- 部门表
CREATE TABLE IF NOT EXISTS departments (
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
  company_id BIGINT NOT NULL,
  parent_id BIGINT NULL,
  CONSTRAINT uq_departments_code UNIQUE (code),
  CONSTRAINT fk_departments_company FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE,
  CONSTRAINT fk_departments_parent FOREIGN KEY (parent_id) REFERENCES departments(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_departments_deleted_at ON departments (deleted_at);
CREATE INDEX IF NOT EXISTS idx_departments_is_active ON departments (is_active);
CREATE INDEX IF NOT EXISTS idx_departments_company_id ON departments (company_id);
CREATE INDEX IF NOT EXISTS idx_departments_parent_id ON departments (parent_id);
CREATE INDEX IF NOT EXISTS idx_departments_code ON departments (code);

-- 用户表
CREATE TABLE IF NOT EXISTS users (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL,
  created_by BIGINT NULL,
  updated_by BIGINT NULL,
  is_active BOOLEAN DEFAULT TRUE,
  username VARCHAR(100) NOT NULL,
  email VARCHAR(255) NOT NULL,
  password_hash VARCHAR(255) NOT NULL,
  first_name VARCHAR(100) NULL,
  last_name VARCHAR(100) NULL,
  phone VARCHAR(50) NULL,
  avatar VARCHAR(500) NULL,
  last_login_at TIMESTAMP NULL,
  email_verified_at TIMESTAMP NULL,
  phone_verified_at TIMESTAMP NULL,
  two_factor_enabled BOOLEAN DEFAULT FALSE,
  two_factor_secret VARCHAR(255) NULL,
  language VARCHAR(10) DEFAULT 'zh-CN',
  timezone VARCHAR(50) DEFAULT 'Asia/Shanghai',
  CONSTRAINT uq_users_username UNIQUE (username),
  CONSTRAINT uq_users_email UNIQUE (email)
);
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users (deleted_at);
CREATE INDEX IF NOT EXISTS idx_users_is_active ON users (is_active);
CREATE INDEX IF NOT EXISTS idx_users_username ON users (username);
CREATE INDEX IF NOT EXISTS idx_users_email ON users (email);
CREATE INDEX IF NOT EXISTS idx_users_last_login_at ON users (last_login_at);

-- 角色表
CREATE TABLE IF NOT EXISTS roles (
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
  is_system BOOLEAN DEFAULT FALSE,
  CONSTRAINT uq_roles_code UNIQUE (code)
);
CREATE INDEX IF NOT EXISTS idx_roles_deleted_at ON roles (deleted_at);
CREATE INDEX IF NOT EXISTS idx_roles_is_active ON roles (is_active);
CREATE INDEX IF NOT EXISTS idx_roles_code ON roles (code);
CREATE INDEX IF NOT EXISTS idx_roles_is_system ON roles (is_system);

-- 权限表
CREATE TABLE IF NOT EXISTS permissions (
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
  resource VARCHAR(100) NOT NULL,
  action VARCHAR(50) NOT NULL,
  CONSTRAINT uq_permissions_code UNIQUE (code),
  CONSTRAINT uq_permissions_resource_action UNIQUE (resource, action)
);
CREATE INDEX IF NOT EXISTS idx_permissions_deleted_at ON permissions (deleted_at);
CREATE INDEX IF NOT EXISTS idx_permissions_is_active ON permissions (is_active);
CREATE INDEX IF NOT EXISTS idx_permissions_code ON permissions (code);
CREATE INDEX IF NOT EXISTS idx_permissions_resource ON permissions (resource);
CREATE INDEX IF NOT EXISTS idx_permissions_action ON permissions (action);

-- ============================================================================
-- 关联表
-- ============================================================================

-- 用户角色关联表
CREATE TABLE IF NOT EXISTS user_roles (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL,
  user_id BIGINT NOT NULL,
  role_id BIGINT NOT NULL,
  assigned_by BIGINT NULL,
  assigned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  expires_at TIMESTAMP NULL,
  CONSTRAINT fk_user_roles_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  CONSTRAINT fk_user_roles_role FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
  CONSTRAINT fk_user_roles_assigned_by FOREIGN KEY (assigned_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT uq_user_roles_user_role UNIQUE (user_id, role_id)
);
CREATE INDEX IF NOT EXISTS idx_user_roles_deleted_at ON user_roles (deleted_at);
CREATE INDEX IF NOT EXISTS idx_user_roles_user_id ON user_roles (user_id);
CREATE INDEX IF NOT EXISTS idx_user_roles_role_id ON user_roles (role_id);
CREATE INDEX IF NOT EXISTS idx_user_roles_assigned_by ON user_roles (assigned_by);
CREATE INDEX IF NOT EXISTS idx_user_roles_expires_at ON user_roles (expires_at);

-- 角色权限关联表
CREATE TABLE IF NOT EXISTS role_permissions (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL,
  role_id BIGINT NOT NULL,
  permission_id BIGINT NOT NULL,
  CONSTRAINT fk_role_permissions_role FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
  CONSTRAINT fk_role_permissions_permission FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE,
  CONSTRAINT uq_role_permissions_role_permission UNIQUE (role_id, permission_id)
);
CREATE INDEX IF NOT EXISTS idx_role_permissions_deleted_at ON role_permissions (deleted_at);
CREATE INDEX IF NOT EXISTS idx_role_permissions_role_id ON role_permissions (role_id);
CREATE INDEX IF NOT EXISTS idx_role_permissions_permission_id ON role_permissions (permission_id);

-- 用户会话表
CREATE TABLE IF NOT EXISTS user_sessions (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL,
  user_id BIGINT NOT NULL,
  token VARCHAR(255) NOT NULL,
  refresh_token VARCHAR(255) NULL,
  expires_at TIMESTAMP NOT NULL,
  ip_address INET NULL,
  user_agent TEXT NULL,
  is_active BOOLEAN DEFAULT TRUE,
  CONSTRAINT fk_user_sessions_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  CONSTRAINT uq_user_sessions_token UNIQUE (token)
);
CREATE INDEX IF NOT EXISTS idx_user_sessions_deleted_at ON user_sessions (deleted_at);
CREATE INDEX IF NOT EXISTS idx_user_sessions_user_id ON user_sessions (user_id);
CREATE INDEX IF NOT EXISTS idx_user_sessions_token ON user_sessions (token);
CREATE INDEX IF NOT EXISTS idx_user_sessions_expires_at ON user_sessions (expires_at);
CREATE INDEX IF NOT EXISTS idx_user_sessions_is_active ON user_sessions (is_active);

-- ============================================================================
-- 初始数据插入
-- ============================================================================

-- 插入默认公司
INSERT INTO companies (code, name, description, is_active) 
VALUES ('DEFAULT', 'GalaxyERP 默认公司', '系统默认公司', TRUE)
ON CONFLICT (code) DO NOTHING;

-- 插入默认部门
INSERT INTO departments (code, name, description, company_id, is_active)
SELECT 'IT', 'IT部门', '信息技术部门', c.id, TRUE
FROM companies c WHERE c.code = 'DEFAULT'
ON CONFLICT (code) DO NOTHING;

-- 插入系统角色
INSERT INTO roles (code, name, description, is_system, is_active) VALUES
('SUPER_ADMIN', '超级管理员', '系统超级管理员，拥有所有权限', TRUE, TRUE),
('ADMIN', '管理员', '系统管理员', TRUE, TRUE),
('USER', '普通用户', '普通用户', TRUE, TRUE)
ON CONFLICT (code) DO NOTHING;

-- 插入基础权限
INSERT INTO permissions (code, name, description, resource, action, is_active) VALUES
('USER_CREATE', '创建用户', '创建新用户', 'user', 'create', TRUE),
('USER_READ', '查看用户', '查看用户信息', 'user', 'read', TRUE),
('USER_UPDATE', '更新用户', '更新用户信息', 'user', 'update', TRUE),
('USER_DELETE', '删除用户', '删除用户', 'user', 'delete', TRUE),
('ROLE_CREATE', '创建角色', '创建新角色', 'role', 'create', TRUE),
('ROLE_READ', '查看角色', '查看角色信息', 'role', 'read', TRUE),
('ROLE_UPDATE', '更新角色', '更新角色信息', 'role', 'update', TRUE),
('ROLE_DELETE', '删除角色', '删除角色', 'role', 'delete', TRUE),
('PERMISSION_READ', '查看权限', '查看权限信息', 'permission', 'read', TRUE),
('SYSTEM_CONFIG', '系统配置', '系统配置管理', 'system', 'config', TRUE)
ON CONFLICT (code) DO NOTHING;

-- 为超级管理员角色分配所有权限
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.code = 'SUPER_ADMIN'
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- ============================================================================
-- 结束
-- ============================================================================