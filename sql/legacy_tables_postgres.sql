-- ============================================================================
-- GalaxyERP 遗留模块 - PostgreSQL 建表脚本
-- 生成时间: 2025-01-01
-- 说明: 基于 /internal/models/legacy_item.go 的结构，包含遗留物料数据
-- 注意: LegacyItem 与 Item 共享 items 表，此脚本仅作为参考
-- ============================================================================

-- ============================================================================
-- 遗留物料管理
-- ============================================================================

-- 注意：LegacyItem 结构体与 Item 结构体共享同一个数据表 'items'
-- 这里仅作为文档说明，实际的 items 表已在 inventory_tables_postgres.sql 中定义

-- LegacyItem 结构体字段映射说明：
-- - 继承自 BaseModel: id, created_at, updated_at, deleted_at
-- - 继承自 AuditableModel: created_by, updated_by  
-- - 继承自 StatusModel: is_active
-- - 继承自 CodeModel: code
-- - 继承自 DescriptionModel: name, description
-- - 特有字段: legacy_code, legacy_name, legacy_description, migration_status, migration_notes

-- 如果需要单独的遗留物料表（不与items表共享），可以使用以下结构：

/*
CREATE TABLE IF NOT EXISTS legacy_items (
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
  legacy_code VARCHAR(100) NULL,
  legacy_name VARCHAR(255) NULL,
  legacy_description TEXT NULL,
  migration_status VARCHAR(20) DEFAULT 'PENDING',
  migration_notes TEXT NULL,
  migrated_item_id BIGINT NULL,
  CONSTRAINT uq_legacy_items_code UNIQUE (code),
  CONSTRAINT fk_legacy_items_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_legacy_items_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_legacy_items_migrated_item FOREIGN KEY (migrated_item_id) REFERENCES items(id) ON DELETE SET NULL,
  CONSTRAINT chk_legacy_items_migration_status CHECK (migration_status IN ('PENDING', 'IN_PROGRESS', 'COMPLETED', 'FAILED', 'SKIPPED'))
);

CREATE INDEX IF NOT EXISTS idx_legacy_items_deleted_at ON legacy_items (deleted_at);
CREATE INDEX IF NOT EXISTS idx_legacy_items_is_active ON legacy_items (is_active);
CREATE INDEX IF NOT EXISTS idx_legacy_items_code ON legacy_items (code);
CREATE INDEX IF NOT EXISTS idx_legacy_items_name ON legacy_items (name);
CREATE INDEX IF NOT EXISTS idx_legacy_items_legacy_code ON legacy_items (legacy_code);
CREATE INDEX IF NOT EXISTS idx_legacy_items_migration_status ON legacy_items (migration_status);
CREATE INDEX IF NOT EXISTS idx_legacy_items_migrated_item_id ON legacy_items (migrated_item_id);
CREATE INDEX IF NOT EXISTS idx_legacy_items_created_by ON legacy_items (created_by);
CREATE INDEX IF NOT EXISTS idx_legacy_items_updated_by ON legacy_items (updated_by);
*/

-- ============================================================================
-- 数据迁移相关表
-- ============================================================================

-- 数据迁移记录表
CREATE TABLE IF NOT EXISTS data_migrations (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL,
  migration_name VARCHAR(100) NOT NULL,
  migration_type VARCHAR(50) NOT NULL,
  source_table VARCHAR(100) NULL,
  target_table VARCHAR(100) NULL,
  total_records INTEGER DEFAULT 0,
  processed_records INTEGER DEFAULT 0,
  success_records INTEGER DEFAULT 0,
  failed_records INTEGER DEFAULT 0,
  start_time TIMESTAMP NULL,
  end_time TIMESTAMP NULL,
  duration_seconds INTEGER NULL,
  status VARCHAR(20) DEFAULT 'PENDING',
  error_message TEXT NULL,
  migration_config JSONB NULL,
  created_by BIGINT NULL,
  CONSTRAINT uq_data_migrations_name UNIQUE (migration_name),
  CONSTRAINT fk_data_migrations_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT chk_data_migrations_status CHECK (status IN ('PENDING', 'RUNNING', 'COMPLETED', 'FAILED', 'CANCELLED'))
);

CREATE INDEX IF NOT EXISTS idx_data_migrations_migration_name ON data_migrations (migration_name);
CREATE INDEX IF NOT EXISTS idx_data_migrations_migration_type ON data_migrations (migration_type);
CREATE INDEX IF NOT EXISTS idx_data_migrations_status ON data_migrations (status);
CREATE INDEX IF NOT EXISTS idx_data_migrations_start_time ON data_migrations (start_time);
CREATE INDEX IF NOT EXISTS idx_data_migrations_created_by ON data_migrations (created_by);

-- 数据迁移详细日志表
CREATE TABLE IF NOT EXISTS migration_logs (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  migration_id BIGINT NOT NULL,
  record_id BIGINT NULL,
  record_key VARCHAR(255) NULL,
  operation VARCHAR(20) NOT NULL,
  status VARCHAR(20) NOT NULL,
  old_data JSONB NULL,
  new_data JSONB NULL,
  error_message TEXT NULL,
  processing_time_ms INTEGER NULL,
  CONSTRAINT fk_migration_logs_migration FOREIGN KEY (migration_id) REFERENCES data_migrations(id) ON DELETE CASCADE,
  CONSTRAINT chk_migration_logs_operation CHECK (operation IN ('INSERT', 'UPDATE', 'DELETE', 'VALIDATE', 'TRANSFORM')),
  CONSTRAINT chk_migration_logs_status CHECK (status IN ('SUCCESS', 'FAILED', 'SKIPPED', 'WARNING'))
);

CREATE INDEX IF NOT EXISTS idx_migration_logs_migration_id ON migration_logs (migration_id);
CREATE INDEX IF NOT EXISTS idx_migration_logs_record_id ON migration_logs (record_id);
CREATE INDEX IF NOT EXISTS idx_migration_logs_record_key ON migration_logs (record_key);
CREATE INDEX IF NOT EXISTS idx_migration_logs_operation ON migration_logs (operation);
CREATE INDEX IF NOT EXISTS idx_migration_logs_status ON migration_logs (status);
CREATE INDEX IF NOT EXISTS idx_migration_logs_created_at ON migration_logs (created_at);

-- ============================================================================
-- 初始数据插入
-- ============================================================================

-- 插入示例数据迁移记录
INSERT INTO data_migrations (migration_name, migration_type, source_table, target_table, status, migration_config, created_by) VALUES
('legacy_items_migration', 'DATA_IMPORT', 'legacy_items', 'items', 'PENDING', '{"batch_size": 1000, "validate_data": true, "skip_duplicates": true}', 1),
('customer_data_cleanup', 'DATA_CLEANUP', 'customers', 'customers', 'PENDING', '{"remove_duplicates": true, "normalize_phone": true, "validate_email": true}', 1),
('inventory_reconciliation', 'DATA_RECONCILIATION', 'inventories', 'inventories', 'PENDING', '{"check_negative_stock": true, "validate_locations": true}', 1)
ON CONFLICT (migration_name) DO NOTHING;

-- ============================================================================
-- 结束
-- ============================================================================