-- ============================================================================
-- GalaxyERP 系统管理模块 - PostgreSQL 建表脚本
-- 生成时间: 2025-01-01
-- 说明: 基于 /internal/models/system.go 的结构，包含系统配置、审计日志、通知管理
-- ============================================================================

-- ============================================================================
-- 系统配置管理
-- ============================================================================

-- 系统配置表
CREATE TABLE IF NOT EXISTS system_configs (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL,
  created_by INTEGER NULL,
  updated_by INTEGER NULL,
  config_key VARCHAR(100) NOT NULL,
  config_value TEXT NOT NULL,
  config_type VARCHAR(20) DEFAULT 'STRING',
  category VARCHAR(50) DEFAULT 'GENERAL',
  description TEXT NULL,
  is_encrypted BOOLEAN DEFAULT FALSE,
  is_system BOOLEAN DEFAULT FALSE,
  default_value TEXT NULL,
  validation_rule VARCHAR(500) NULL,
  display_order INTEGER DEFAULT 0,
  CONSTRAINT fk_system_configs_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_system_configs_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT uq_system_configs_config_key UNIQUE (config_key),
  CONSTRAINT chk_system_configs_config_type CHECK (config_type IN ('STRING', 'INTEGER', 'DECIMAL', 'BOOLEAN', 'JSON', 'TEXT'))
);
CREATE INDEX IF NOT EXISTS idx_system_configs_deleted_at ON system_configs (deleted_at);
CREATE INDEX IF NOT EXISTS idx_system_configs_is_active ON system_configs (is_active);
CREATE INDEX IF NOT EXISTS idx_system_configs_config_key ON system_configs (config_key);
CREATE INDEX IF NOT EXISTS idx_system_configs_category ON system_configs (category);
CREATE INDEX IF NOT EXISTS idx_system_configs_is_system ON system_configs (is_system);
CREATE INDEX IF NOT EXISTS idx_system_configs_display_order ON system_configs (display_order);

-- ============================================================================
-- 审计日志管理
-- ============================================================================

-- 审计日志表
CREATE TABLE IF NOT EXISTS audit_logs (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  user_id INTEGER NULL,
  username VARCHAR(100) NULL,
  action VARCHAR(50) NOT NULL,
  resource_type VARCHAR(50) NOT NULL,
  resource_id INTEGER NULL,
  resource_name VARCHAR(255) NULL,
  old_values JSONB NULL,
  new_values JSONB NULL,
  ip_address VARCHAR(45) NULL,
  user_agent TEXT NULL,
  request_id VARCHAR(100) NULL,
  session_id VARCHAR(100) NULL,
  result VARCHAR(20) DEFAULT 'SUCCESS',
  error_message TEXT NULL,
  duration_ms INTEGER NULL,
  CONSTRAINT fk_audit_logs_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT chk_audit_logs_result CHECK (result IN ('SUCCESS', 'FAILED', 'ERROR'))
);
CREATE INDEX IF NOT EXISTS idx_audit_logs_created_at ON audit_logs (created_at);
CREATE INDEX IF NOT EXISTS idx_audit_logs_user_id ON audit_logs (user_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_session_id ON audit_logs (session_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_action ON audit_logs (action);
CREATE INDEX IF NOT EXISTS idx_audit_logs_resource_type ON audit_logs (resource_type);
CREATE INDEX IF NOT EXISTS idx_audit_logs_resource_id ON audit_logs (resource_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_ip_address ON audit_logs (ip_address);
CREATE INDEX IF NOT EXISTS idx_audit_logs_request_method ON audit_logs (request_method);
CREATE INDEX IF NOT EXISTS idx_audit_logs_response_status ON audit_logs (response_status);

-- ============================================================================
-- 系统通知管理
-- ============================================================================

-- 通知模板表
CREATE TABLE IF NOT EXISTS notification_templates (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL,
  created_by BIGINT NULL,
  updated_by BIGINT NULL,
  is_active BOOLEAN DEFAULT TRUE,
  code VARCHAR(50) NOT NULL,
  name VARCHAR(100) NOT NULL,
  description TEXT NULL,
  notification_type VARCHAR(20) NOT NULL,
  category VARCHAR(50) DEFAULT 'GENERAL',
  subject_template TEXT NULL,
  body_template TEXT NOT NULL,
  variables JSONB NULL,
  is_system BOOLEAN DEFAULT FALSE,
  status VARCHAR(20) DEFAULT 'ACTIVE',
  CONSTRAINT uq_notification_templates_code UNIQUE (code),
  CONSTRAINT chk_notification_templates_notification_type CHECK (notification_type IN ('EMAIL', 'SMS', 'PUSH', 'IN_APP', 'WEBHOOK'))
);
CREATE INDEX IF NOT EXISTS idx_notification_templates_deleted_at ON notification_templates (deleted_at);
CREATE INDEX IF NOT EXISTS idx_notification_templates_is_active ON notification_templates (is_active);
CREATE INDEX IF NOT EXISTS idx_notification_templates_code ON notification_templates (code);
CREATE INDEX IF NOT EXISTS idx_notification_templates_notification_type ON notification_templates (notification_type);
CREATE INDEX IF NOT EXISTS idx_notification_templates_category ON notification_templates (category);
CREATE INDEX IF NOT EXISTS idx_notification_templates_is_system ON notification_templates (is_system);
CREATE INDEX IF NOT EXISTS idx_notification_templates_status ON notification_templates (status);

-- 通知表
CREATE TABLE IF NOT EXISTS notifications (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL,
  notification_template_id BIGINT NULL,
  notification_type VARCHAR(20) NOT NULL,
  recipient_type VARCHAR(20) NOT NULL,
  recipient_id BIGINT NULL,
  recipient_address VARCHAR(255) NULL,
  subject VARCHAR(500) NULL,
  body TEXT NOT NULL,
  variables JSONB NULL,
  priority VARCHAR(10) DEFAULT 'NORMAL',
  scheduled_at TIMESTAMP NULL,
  sent_at TIMESTAMP NULL,
  delivered_at TIMESTAMP NULL,
  read_at TIMESTAMP NULL,
  error_message TEXT NULL,
  retry_count INTEGER DEFAULT 0,
  max_retries INTEGER DEFAULT 3,
  status VARCHAR(20) DEFAULT 'PENDING',
  CONSTRAINT fk_notifications_notification_template FOREIGN KEY (notification_template_id) REFERENCES notification_templates(id) ON DELETE SET NULL,
  CONSTRAINT chk_notifications_notification_type CHECK (notification_type IN ('EMAIL', 'SMS', 'PUSH', 'IN_APP', 'WEBHOOK')),
  CONSTRAINT chk_notifications_recipient_type CHECK (recipient_type IN ('USER', 'EMAIL', 'PHONE', 'DEVICE', 'WEBHOOK')),
  CONSTRAINT chk_notifications_priority CHECK (priority IN ('LOW', 'NORMAL', 'HIGH', 'URGENT')),
  CONSTRAINT chk_notifications_status CHECK (status IN ('PENDING', 'SCHEDULED', 'SENDING', 'SENT', 'DELIVERED', 'READ', 'FAILED', 'CANCELLED'))
);
CREATE INDEX IF NOT EXISTS idx_notifications_deleted_at ON notifications (deleted_at);
CREATE INDEX IF NOT EXISTS idx_notifications_notification_template_id ON notifications (notification_template_id);
CREATE INDEX IF NOT EXISTS idx_notifications_notification_type ON notifications (notification_type);
CREATE INDEX IF NOT EXISTS idx_notifications_recipient_type ON notifications (recipient_type);
CREATE INDEX IF NOT EXISTS idx_notifications_recipient_id ON notifications (recipient_id);
CREATE INDEX IF NOT EXISTS idx_notifications_priority ON notifications (priority);
CREATE INDEX IF NOT EXISTS idx_notifications_scheduled_at ON notifications (scheduled_at);
CREATE INDEX IF NOT EXISTS idx_notifications_sent_at ON notifications (sent_at);
CREATE INDEX IF NOT EXISTS idx_notifications_status ON notifications (status);
CREATE INDEX IF NOT EXISTS idx_notifications_created_at ON notifications (created_at);







-- 数据权限表
CREATE TABLE IF NOT EXISTS data_permissions (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL,
  created_by INTEGER NULL,
  updated_by INTEGER NULL,
  user_id INTEGER NOT NULL,
  resource_type VARCHAR(50) NOT NULL,
  resource_id INTEGER NULL,
  permission_type VARCHAR(20) NOT NULL,
  granted_by INTEGER NULL,
  granted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  expires_at TIMESTAMP NULL,
  CONSTRAINT fk_data_permissions_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  CONSTRAINT fk_data_permissions_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_data_permissions_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_data_permissions_granted_by FOREIGN KEY (granted_by) REFERENCES users(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_data_permissions_deleted_at ON data_permissions (deleted_at);
CREATE INDEX IF NOT EXISTS idx_data_permissions_user_id ON data_permissions (user_id);
CREATE INDEX IF NOT EXISTS idx_data_permissions_resource_type ON data_permissions (resource_type);
CREATE INDEX IF NOT EXISTS idx_data_permissions_resource_id ON data_permissions (resource_id);
CREATE INDEX IF NOT EXISTS idx_data_permissions_permission_type ON data_permissions (permission_type);

-- 审批步骤表
CREATE TABLE IF NOT EXISTS approval_steps (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL,
  created_by INTEGER NULL,
  updated_by INTEGER NULL,
  workflow_type VARCHAR(50) NOT NULL,
  step_order INTEGER NOT NULL,
  step_name VARCHAR(100) NOT NULL,
  approver_type VARCHAR(20) NOT NULL,
  approver_id INTEGER NULL,
  is_required BOOLEAN DEFAULT TRUE,
  conditions JSONB NULL,
  CONSTRAINT fk_approval_steps_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_approval_steps_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_approval_steps_approver FOREIGN KEY (approver_id) REFERENCES users(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_approval_steps_deleted_at ON approval_steps (deleted_at);
CREATE INDEX IF NOT EXISTS idx_approval_steps_workflow_type ON approval_steps (workflow_type);
CREATE INDEX IF NOT EXISTS idx_approval_steps_step_order ON approval_steps (step_order);
CREATE INDEX IF NOT EXISTS idx_approval_steps_approver_type ON approval_steps (approver_type);
CREATE INDEX IF NOT EXISTS idx_approval_steps_approver_id ON approval_steps (approver_id);

-- ============================================================================
-- 初始数据插入
-- ============================================================================

-- 插入系统配置
INSERT INTO system_configs (config_key, config_value, config_type, category, description, is_system, default_value, display_order) VALUES
-- 基础配置
('system.name', 'GalaxyERP', 'STRING', 'BASIC', '系统名称', TRUE, 'GalaxyERP', 1),
('system.version', '1.0.0', 'STRING', 'BASIC', '系统版本', TRUE, '1.0.0', 2),
('system.timezone', 'Asia/Shanghai', 'STRING', 'BASIC', '系统时区', TRUE, 'Asia/Shanghai', 3),
('system.language', 'zh-CN', 'STRING', 'BASIC', '系统语言', TRUE, 'zh-CN', 4),
-- 安全配置
('security.password.min_length', '8', 'INTEGER', 'SECURITY', '密码最小长度', TRUE, '8', 10),
('security.session.timeout', '3600', 'INTEGER', 'SECURITY', '会话超时时间(秒)', TRUE, '3600', 15)
ON CONFLICT (config_key) DO NOTHING;

-- ============================================================================
-- 结束
-- ============================================================================