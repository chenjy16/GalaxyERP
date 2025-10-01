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
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL,
  created_by BIGINT NULL,
  updated_by BIGINT NULL,
  is_active BOOLEAN DEFAULT TRUE,
  config_key VARCHAR(100) NOT NULL,
  config_value TEXT NULL,
  config_type VARCHAR(20) DEFAULT 'STRING',
  category VARCHAR(50) DEFAULT 'GENERAL',
  description TEXT NULL,
  is_system BOOLEAN DEFAULT FALSE,
  is_encrypted BOOLEAN DEFAULT FALSE,
  validation_rule TEXT NULL,
  default_value TEXT NULL,
  display_order INTEGER DEFAULT 0,
  CONSTRAINT uq_system_configs_config_key UNIQUE (config_key),
  CONSTRAINT chk_system_configs_config_type CHECK (config_type IN ('STRING', 'INTEGER', 'DECIMAL', 'BOOLEAN', 'JSON', 'DATE', 'DATETIME'))
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
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  user_id BIGINT NULL,
  session_id VARCHAR(255) NULL,
  ip_address VARCHAR(45) NULL,
  user_agent TEXT NULL,
  action VARCHAR(50) NOT NULL,
  resource_type VARCHAR(50) NOT NULL,
  resource_id BIGINT NULL,
  resource_name VARCHAR(255) NULL,
  old_values JSONB NULL,
  new_values JSONB NULL,
  changes JSONB NULL,
  request_method VARCHAR(10) NULL,
  request_url TEXT NULL,
  request_body TEXT NULL,
  response_status INTEGER NULL,
  response_body TEXT NULL,
  execution_time_ms INTEGER NULL,
  error_message TEXT NULL,
  tags JSONB NULL,
  CONSTRAINT fk_audit_logs_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
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
  name VARCHAR(255) NOT NULL,
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

-- ============================================================================
-- 文件管理
-- ============================================================================

-- 文件表
CREATE TABLE IF NOT EXISTS files (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL,
  created_by BIGINT NULL,
  updated_by BIGINT NULL,
  is_active BOOLEAN DEFAULT TRUE,
  filename VARCHAR(255) NOT NULL,
  original_filename VARCHAR(255) NOT NULL,
  file_path VARCHAR(500) NOT NULL,
  file_size BIGINT DEFAULT 0,
  mime_type VARCHAR(100) NULL,
  file_extension VARCHAR(10) NULL,
  file_hash VARCHAR(64) NULL,
  storage_type VARCHAR(20) DEFAULT 'LOCAL',
  storage_config JSONB NULL,
  reference_type VARCHAR(50) NULL,
  reference_id BIGINT NULL,
  category VARCHAR(50) DEFAULT 'GENERAL',
  description TEXT NULL,
  is_public BOOLEAN DEFAULT FALSE,
  download_count INTEGER DEFAULT 0,
  status VARCHAR(20) DEFAULT 'ACTIVE',
  CONSTRAINT fk_files_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT chk_files_storage_type CHECK (storage_type IN ('LOCAL', 'S3', 'OSS', 'COS', 'FTP', 'SFTP'))
);
CREATE INDEX IF NOT EXISTS idx_files_deleted_at ON files (deleted_at);
CREATE INDEX IF NOT EXISTS idx_files_is_active ON files (is_active);
CREATE INDEX IF NOT EXISTS idx_files_filename ON files (filename);
CREATE INDEX IF NOT EXISTS idx_files_file_path ON files (file_path);
CREATE INDEX IF NOT EXISTS idx_files_mime_type ON files (mime_type);
CREATE INDEX IF NOT EXISTS idx_files_file_extension ON files (file_extension);
CREATE INDEX IF NOT EXISTS idx_files_file_hash ON files (file_hash);
CREATE INDEX IF NOT EXISTS idx_files_storage_type ON files (storage_type);
CREATE INDEX IF NOT EXISTS idx_files_reference_type ON files (reference_type);
CREATE INDEX IF NOT EXISTS idx_files_reference_id ON files (reference_id);
CREATE INDEX IF NOT EXISTS idx_files_category ON files (category);
CREATE INDEX IF NOT EXISTS idx_files_created_by ON files (created_by);
CREATE INDEX IF NOT EXISTS idx_files_is_public ON files (is_public);
CREATE INDEX IF NOT EXISTS idx_files_status ON files (status);

-- ============================================================================
-- 数据字典管理
-- ============================================================================

-- 数据字典表
CREATE TABLE IF NOT EXISTS dictionaries (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL,
  created_by BIGINT NULL,
  updated_by BIGINT NULL,
  is_active BOOLEAN DEFAULT TRUE,
  dict_type VARCHAR(50) NOT NULL,
  dict_key VARCHAR(100) NOT NULL,
  dict_value VARCHAR(500) NOT NULL,
  dict_label VARCHAR(255) NOT NULL,
  description TEXT NULL,
  sort_order INTEGER DEFAULT 0,
  parent_id BIGINT NULL,
  is_system BOOLEAN DEFAULT FALSE,
  extra_data JSONB NULL,
  status VARCHAR(20) DEFAULT 'ACTIVE',
  CONSTRAINT fk_dictionaries_parent FOREIGN KEY (parent_id) REFERENCES dictionaries(id) ON DELETE SET NULL,
  CONSTRAINT uq_dictionaries_type_key UNIQUE (dict_type, dict_key)
);
CREATE INDEX IF NOT EXISTS idx_dictionaries_deleted_at ON dictionaries (deleted_at);
CREATE INDEX IF NOT EXISTS idx_dictionaries_is_active ON dictionaries (is_active);
CREATE INDEX IF NOT EXISTS idx_dictionaries_dict_type ON dictionaries (dict_type);
CREATE INDEX IF NOT EXISTS idx_dictionaries_dict_key ON dictionaries (dict_key);
CREATE INDEX IF NOT EXISTS idx_dictionaries_dict_value ON dictionaries (dict_value);
CREATE INDEX IF NOT EXISTS idx_dictionaries_sort_order ON dictionaries (sort_order);
CREATE INDEX IF NOT EXISTS idx_dictionaries_parent_id ON dictionaries (parent_id);
CREATE INDEX IF NOT EXISTS idx_dictionaries_is_system ON dictionaries (is_system);
CREATE INDEX IF NOT EXISTS idx_dictionaries_status ON dictionaries (status);

-- ============================================================================
-- 任务调度管理
-- ============================================================================

-- 定时任务表
CREATE TABLE IF NOT EXISTS scheduled_jobs (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL,
  created_by BIGINT NULL,
  updated_by BIGINT NULL,
  is_active BOOLEAN DEFAULT TRUE,
  job_name VARCHAR(100) NOT NULL,
  job_group VARCHAR(50) DEFAULT 'DEFAULT',
  description TEXT NULL,
  job_class VARCHAR(255) NOT NULL,
  job_method VARCHAR(100) NULL,
  job_params JSONB NULL,
  cron_expression VARCHAR(100) NOT NULL,
  timezone VARCHAR(50) DEFAULT 'Asia/Shanghai',
  start_time TIMESTAMP NULL,
  end_time TIMESTAMP NULL,
  last_run_time TIMESTAMP NULL,
  next_run_time TIMESTAMP NULL,
  run_count INTEGER DEFAULT 0,
  max_runs INTEGER NULL,
  timeout_seconds INTEGER DEFAULT 300,
  retry_count INTEGER DEFAULT 0,
  max_retries INTEGER DEFAULT 3,
  is_concurrent BOOLEAN DEFAULT FALSE,
  priority INTEGER DEFAULT 5,
  status VARCHAR(20) DEFAULT 'ACTIVE',
  CONSTRAINT uq_scheduled_jobs_name_group UNIQUE (job_name, job_group)
);
CREATE INDEX IF NOT EXISTS idx_scheduled_jobs_deleted_at ON scheduled_jobs (deleted_at);
CREATE INDEX IF NOT EXISTS idx_scheduled_jobs_is_active ON scheduled_jobs (is_active);
CREATE INDEX IF NOT EXISTS idx_scheduled_jobs_job_name ON scheduled_jobs (job_name);
CREATE INDEX IF NOT EXISTS idx_scheduled_jobs_job_group ON scheduled_jobs (job_group);
CREATE INDEX IF NOT EXISTS idx_scheduled_jobs_next_run_time ON scheduled_jobs (next_run_time);
CREATE INDEX IF NOT EXISTS idx_scheduled_jobs_status ON scheduled_jobs (status);
CREATE INDEX IF NOT EXISTS idx_scheduled_jobs_priority ON scheduled_jobs (priority);

-- 任务执行日志表
CREATE TABLE IF NOT EXISTS job_execution_logs (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  scheduled_job_id BIGINT NOT NULL,
  job_name VARCHAR(100) NOT NULL,
  job_group VARCHAR(50) NOT NULL,
  trigger_type VARCHAR(20) DEFAULT 'SCHEDULED',
  start_time TIMESTAMP NOT NULL,
  end_time TIMESTAMP NULL,
  execution_time_ms INTEGER NULL,
  status VARCHAR(20) DEFAULT 'RUNNING',
  result_message TEXT NULL,
  error_message TEXT NULL,
  stack_trace TEXT NULL,
  server_name VARCHAR(100) NULL,
  process_id INTEGER NULL,
  thread_id VARCHAR(50) NULL,
  CONSTRAINT fk_job_execution_logs_scheduled_job FOREIGN KEY (scheduled_job_id) REFERENCES scheduled_jobs(id) ON DELETE CASCADE,
  CONSTRAINT chk_job_execution_logs_trigger_type CHECK (trigger_type IN ('SCHEDULED', 'MANUAL', 'RETRY')),
  CONSTRAINT chk_job_execution_logs_status CHECK (status IN ('RUNNING', 'SUCCESS', 'FAILED', 'TIMEOUT', 'CANCELLED'))
);
CREATE INDEX IF NOT EXISTS idx_job_execution_logs_scheduled_job_id ON job_execution_logs (scheduled_job_id);
CREATE INDEX IF NOT EXISTS idx_job_execution_logs_job_name ON job_execution_logs (job_name);
CREATE INDEX IF NOT EXISTS idx_job_execution_logs_job_group ON job_execution_logs (job_group);
CREATE INDEX IF NOT EXISTS idx_job_execution_logs_start_time ON job_execution_logs (start_time);
CREATE INDEX IF NOT EXISTS idx_job_execution_logs_status ON job_execution_logs (status);
CREATE INDEX IF NOT EXISTS idx_job_execution_logs_trigger_type ON job_execution_logs (trigger_type);
CREATE INDEX IF NOT EXISTS idx_job_execution_logs_created_at ON job_execution_logs (created_at);

-- ============================================================================
-- 系统监控管理
-- ============================================================================

-- 系统监控表
CREATE TABLE IF NOT EXISTS system_monitors (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  server_name VARCHAR(100) NOT NULL,
  monitor_type VARCHAR(20) NOT NULL,
  metric_name VARCHAR(100) NOT NULL,
  metric_value DECIMAL(15,4) NOT NULL,
  metric_unit VARCHAR(20) NULL,
  threshold_warning DECIMAL(15,4) NULL,
  threshold_critical DECIMAL(15,4) NULL,
  status VARCHAR(20) DEFAULT 'NORMAL',
  tags JSONB NULL,
  CONSTRAINT chk_system_monitors_monitor_type CHECK (monitor_type IN ('CPU', 'MEMORY', 'DISK', 'NETWORK', 'DATABASE', 'APPLICATION', 'CUSTOM')),
  CONSTRAINT chk_system_monitors_status CHECK (status IN ('NORMAL', 'WARNING', 'CRITICAL', 'UNKNOWN'))
);
CREATE INDEX IF NOT EXISTS idx_system_monitors_created_at ON system_monitors (created_at);
CREATE INDEX IF NOT EXISTS idx_system_monitors_server_name ON system_monitors (server_name);
CREATE INDEX IF NOT EXISTS idx_system_monitors_monitor_type ON system_monitors (monitor_type);
CREATE INDEX IF NOT EXISTS idx_system_monitors_metric_name ON system_monitors (metric_name);
CREATE INDEX IF NOT EXISTS idx_system_monitors_status ON system_monitors (status);

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
('system.currency', 'CNY', 'STRING', 'BASIC', '系统货币', TRUE, 'CNY', 5),
-- 安全配置
('security.password.min_length', '8', 'INTEGER', 'SECURITY', '密码最小长度', TRUE, '8', 10),
('security.password.require_uppercase', 'true', 'BOOLEAN', 'SECURITY', '密码需要大写字母', TRUE, 'true', 11),
('security.password.require_lowercase', 'true', 'BOOLEAN', 'SECURITY', '密码需要小写字母', TRUE, 'true', 12),
('security.password.require_number', 'true', 'BOOLEAN', 'SECURITY', '密码需要数字', TRUE, 'true', 13),
('security.password.require_special', 'false', 'BOOLEAN', 'SECURITY', '密码需要特殊字符', TRUE, 'false', 14),
('security.session.timeout', '3600', 'INTEGER', 'SECURITY', '会话超时时间(秒)', TRUE, '3600', 15),
('security.login.max_attempts', '5', 'INTEGER', 'SECURITY', '最大登录尝试次数', TRUE, '5', 16),
('security.login.lockout_duration', '300', 'INTEGER', 'SECURITY', '账户锁定时间(秒)', TRUE, '300', 17),
-- 业务配置
('business.fiscal_year_start', '01-01', 'STRING', 'BUSINESS', '财年开始日期', TRUE, '01-01', 20),
('business.default_currency', 'CNY', 'STRING', 'BUSINESS', '默认货币', TRUE, 'CNY', 21),
('business.tax_rate', '13.00', 'DECIMAL', 'BUSINESS', '默认税率(%)', TRUE, '13.00', 22),
('business.payment_terms', '30', 'INTEGER', 'BUSINESS', '默认付款条件(天)', TRUE, '30', 23),
-- 通知配置
('notification.email.enabled', 'true', 'BOOLEAN', 'NOTIFICATION', '启用邮件通知', TRUE, 'true', 30),
('notification.sms.enabled', 'false', 'BOOLEAN', 'NOTIFICATION', '启用短信通知', TRUE, 'false', 31),
('notification.push.enabled', 'true', 'BOOLEAN', 'NOTIFICATION', '启用推送通知', TRUE, 'true', 32)
ON CONFLICT (config_key) DO NOTHING;

-- 插入数据字典
INSERT INTO dictionaries (dict_type, dict_key, dict_value, dict_label, description, sort_order, is_system, status) VALUES
-- 性别
('GENDER', 'MALE', 'MALE', '男', '性别-男', 1, TRUE, 'ACTIVE'),
('GENDER', 'FEMALE', 'FEMALE', '女', '性别-女', 2, TRUE, 'ACTIVE'),
-- 状态
('STATUS', 'ACTIVE', 'ACTIVE', '启用', '状态-启用', 1, TRUE, 'ACTIVE'),
('STATUS', 'INACTIVE', 'INACTIVE', '禁用', '状态-禁用', 2, TRUE, 'ACTIVE'),
('STATUS', 'PENDING', 'PENDING', '待处理', '状态-待处理', 3, TRUE, 'ACTIVE'),
('STATUS', 'APPROVED', 'APPROVED', '已批准', '状态-已批准', 4, TRUE, 'ACTIVE'),
('STATUS', 'REJECTED', 'REJECTED', '已拒绝', '状态-已拒绝', 5, TRUE, 'ACTIVE'),
-- 优先级
('PRIORITY', 'LOW', 'LOW', '低', '优先级-低', 1, TRUE, 'ACTIVE'),
('PRIORITY', 'NORMAL', 'NORMAL', '普通', '优先级-普通', 2, TRUE, 'ACTIVE'),
('PRIORITY', 'HIGH', 'HIGH', '高', '优先级-高', 3, TRUE, 'ACTIVE'),
('PRIORITY', 'URGENT', 'URGENT', '紧急', '优先级-紧急', 4, TRUE, 'ACTIVE'),
-- 货币
('CURRENCY', 'CNY', 'CNY', '人民币', '货币-人民币', 1, TRUE, 'ACTIVE'),
('CURRENCY', 'USD', 'USD', '美元', '货币-美元', 2, TRUE, 'ACTIVE'),
('CURRENCY', 'EUR', 'EUR', '欧元', '货币-欧元', 3, TRUE, 'ACTIVE'),
('CURRENCY', 'JPY', 'JPY', '日元', '货币-日元', 4, TRUE, 'ACTIVE')
ON CONFLICT (dict_type, dict_key) DO NOTHING;

-- 插入通知模板
INSERT INTO notification_templates (code, name, description, notification_type, category, subject_template, body_template, is_system, status) VALUES
('USER_WELCOME', '用户欢迎', '新用户注册欢迎通知', 'EMAIL', 'USER', '欢迎加入{{.SystemName}}', '亲爱的{{.UserName}}，欢迎您加入{{.SystemName}}系统！', TRUE, 'ACTIVE'),
('PASSWORD_RESET', '密码重置', '用户密码重置通知', 'EMAIL', 'SECURITY', '密码重置通知', '您的密码重置链接：{{.ResetLink}}，有效期30分钟。', TRUE, 'ACTIVE'),
('ORDER_CREATED', '订单创建', '新订单创建通知', 'IN_APP', 'BUSINESS', '新订单通知', '订单{{.OrderNumber}}已创建，金额：{{.Amount}}', TRUE, 'ACTIVE'),
('PAYMENT_OVERDUE', '付款逾期', '付款逾期提醒', 'EMAIL', 'FINANCE', '付款逾期提醒', '订单{{.OrderNumber}}付款已逾期{{.OverdueDays}}天，请及时处理。', TRUE, 'ACTIVE')
ON CONFLICT (code) DO NOTHING;

-- ============================================================================
-- 结束
-- ============================================================================