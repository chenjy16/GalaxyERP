-- ============================================================================
-- GalaxyERP 财务会计模块 - PostgreSQL 建表脚本
-- 生成时间: 2025-01-01
-- 说明: 基于 /internal/models/accounting.go 的结构，包含科目、凭证、账簿管理
-- ============================================================================

-- ============================================================================
-- 会计科目管理
-- ============================================================================

-- 会计科目表
CREATE TABLE IF NOT EXISTS accounts (
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
  parent_id BIGINT NULL,
  level INTEGER DEFAULT 1,
  path VARCHAR(500) NULL,
  account_type VARCHAR(20) NOT NULL,
  balance_type VARCHAR(10) NOT NULL,
  is_leaf BOOLEAN DEFAULT TRUE,
  is_system BOOLEAN DEFAULT FALSE,
  currency VARCHAR(10) DEFAULT 'CNY',
  opening_balance DECIMAL(15,2) DEFAULT 0,
  current_balance DECIMAL(15,2) DEFAULT 0,
  status VARCHAR(20) DEFAULT 'ACTIVE',
  CONSTRAINT uq_accounts_code UNIQUE (code),
  CONSTRAINT fk_accounts_parent FOREIGN KEY (parent_id) REFERENCES accounts(id) ON DELETE SET NULL,
  CONSTRAINT chk_accounts_account_type CHECK (account_type IN ('ASSET', 'LIABILITY', 'EQUITY', 'REVENUE', 'EXPENSE')),
  CONSTRAINT chk_accounts_balance_type CHECK (balance_type IN ('DEBIT', 'CREDIT'))
);
CREATE INDEX IF NOT EXISTS idx_accounts_deleted_at ON accounts (deleted_at);
CREATE INDEX IF NOT EXISTS idx_accounts_is_active ON accounts (is_active);
CREATE INDEX IF NOT EXISTS idx_accounts_code ON accounts (code);
CREATE INDEX IF NOT EXISTS idx_accounts_parent_id ON accounts (parent_id);
CREATE INDEX IF NOT EXISTS idx_accounts_account_type ON accounts (account_type);
CREATE INDEX IF NOT EXISTS idx_accounts_balance_type ON accounts (balance_type);
CREATE INDEX IF NOT EXISTS idx_accounts_level ON accounts (level);
CREATE INDEX IF NOT EXISTS idx_accounts_is_leaf ON accounts (is_leaf);
CREATE INDEX IF NOT EXISTS idx_accounts_status ON accounts (status);
CREATE INDEX IF NOT EXISTS idx_accounts_name ON accounts (name);

-- ============================================================================
-- 会计期间管理
-- ============================================================================

-- 会计期间表
CREATE TABLE IF NOT EXISTS fiscal_periods (
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
  year INTEGER NOT NULL,
  period INTEGER NOT NULL,
  start_date DATE NOT NULL,
  end_date DATE NOT NULL,
  is_closed BOOLEAN DEFAULT FALSE,
  closed_by BIGINT NULL,
  closed_at TIMESTAMP NULL,
  status VARCHAR(20) DEFAULT 'OPEN',
  CONSTRAINT uq_fiscal_periods_code UNIQUE (code),
  CONSTRAINT uq_fiscal_periods_year_period UNIQUE (year, period),
  CONSTRAINT fk_fiscal_periods_closed_by FOREIGN KEY (closed_by) REFERENCES users(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_fiscal_periods_deleted_at ON fiscal_periods (deleted_at);
CREATE INDEX IF NOT EXISTS idx_fiscal_periods_is_active ON fiscal_periods (is_active);
CREATE INDEX IF NOT EXISTS idx_fiscal_periods_code ON fiscal_periods (code);
CREATE INDEX IF NOT EXISTS idx_fiscal_periods_year ON fiscal_periods (year);
CREATE INDEX IF NOT EXISTS idx_fiscal_periods_period ON fiscal_periods (period);
CREATE INDEX IF NOT EXISTS idx_fiscal_periods_start_date ON fiscal_periods (start_date);
CREATE INDEX IF NOT EXISTS idx_fiscal_periods_end_date ON fiscal_periods (end_date);
CREATE INDEX IF NOT EXISTS idx_fiscal_periods_is_closed ON fiscal_periods (is_closed);
CREATE INDEX IF NOT EXISTS idx_fiscal_periods_status ON fiscal_periods (status);

-- ============================================================================
-- 会计凭证管理
-- ============================================================================

-- 会计凭证表
CREATE TABLE IF NOT EXISTS journal_entries (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL,
  created_by BIGINT NULL,
  updated_by BIGINT NULL,
  is_active BOOLEAN DEFAULT TRUE,
  code VARCHAR(50) NOT NULL,
  fiscal_period_id BIGINT NOT NULL,
  entry_date DATE NOT NULL,
  reference_type VARCHAR(50) NULL,
  reference_id BIGINT NULL,
  reference_number VARCHAR(100) NULL,
  description TEXT NOT NULL,
  total_debit DECIMAL(15,2) DEFAULT 0,
  total_credit DECIMAL(15,2) DEFAULT 0,
  currency VARCHAR(10) DEFAULT 'CNY',
  exchange_rate DECIMAL(10,4) DEFAULT 1,
  prepared_by BIGINT NULL,
  reviewed_by BIGINT NULL,
  approved_by BIGINT NULL,
  posted_by BIGINT NULL,
  posted_at TIMESTAMP NULL,
  status VARCHAR(20) DEFAULT 'DRAFT',
  CONSTRAINT uq_journal_entries_code UNIQUE (code),
  CONSTRAINT fk_journal_entries_fiscal_period FOREIGN KEY (fiscal_period_id) REFERENCES fiscal_periods(id) ON DELETE CASCADE,
  CONSTRAINT fk_journal_entries_prepared_by FOREIGN KEY (prepared_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_journal_entries_reviewed_by FOREIGN KEY (reviewed_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_journal_entries_approved_by FOREIGN KEY (approved_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_journal_entries_posted_by FOREIGN KEY (posted_by) REFERENCES users(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_journal_entries_deleted_at ON journal_entries (deleted_at);
CREATE INDEX IF NOT EXISTS idx_journal_entries_is_active ON journal_entries (is_active);
CREATE INDEX IF NOT EXISTS idx_journal_entries_code ON journal_entries (code);
CREATE INDEX IF NOT EXISTS idx_journal_entries_fiscal_period_id ON journal_entries (fiscal_period_id);
CREATE INDEX IF NOT EXISTS idx_journal_entries_entry_date ON journal_entries (entry_date);
CREATE INDEX IF NOT EXISTS idx_journal_entries_reference_type ON journal_entries (reference_type);
CREATE INDEX IF NOT EXISTS idx_journal_entries_reference_id ON journal_entries (reference_id);
CREATE INDEX IF NOT EXISTS idx_journal_entries_status ON journal_entries (status);
CREATE INDEX IF NOT EXISTS idx_journal_entries_prepared_by ON journal_entries (prepared_by);
CREATE INDEX IF NOT EXISTS idx_journal_entries_posted_at ON journal_entries (posted_at);

-- 会计凭证明细表
CREATE TABLE IF NOT EXISTS journal_entry_lines (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL,
  journal_entry_id BIGINT NOT NULL,
  line_number INTEGER NOT NULL,
  account_id BIGINT NOT NULL,
  description TEXT NULL,
  debit_amount DECIMAL(15,2) DEFAULT 0,
  credit_amount DECIMAL(15,2) DEFAULT 0,
  currency VARCHAR(10) DEFAULT 'CNY',
  exchange_rate DECIMAL(10,4) DEFAULT 1,
  cost_center_id BIGINT NULL,
  project_id BIGINT NULL,
  CONSTRAINT fk_journal_entry_lines_journal_entry FOREIGN KEY (journal_entry_id) REFERENCES journal_entries(id) ON DELETE CASCADE,
  CONSTRAINT fk_journal_entry_lines_account FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE,
  CONSTRAINT fk_journal_entry_lines_cost_center FOREIGN KEY (cost_center_id) REFERENCES departments(id) ON DELETE SET NULL,
  CONSTRAINT fk_journal_entry_lines_project FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE SET NULL,
  CONSTRAINT chk_journal_entry_lines_amount CHECK ((debit_amount > 0 AND credit_amount = 0) OR (debit_amount = 0 AND credit_amount > 0))
);
CREATE INDEX IF NOT EXISTS idx_journal_entry_lines_deleted_at ON journal_entry_lines (deleted_at);
CREATE INDEX IF NOT EXISTS idx_journal_entry_lines_journal_entry_id ON journal_entry_lines (journal_entry_id);
CREATE INDEX IF NOT EXISTS idx_journal_entry_lines_account_id ON journal_entry_lines (account_id);
CREATE INDEX IF NOT EXISTS idx_journal_entry_lines_cost_center_id ON journal_entry_lines (cost_center_id);
CREATE INDEX IF NOT EXISTS idx_journal_entry_lines_project_id ON journal_entry_lines (project_id);
CREATE INDEX IF NOT EXISTS idx_journal_entry_lines_line_number ON journal_entry_lines (line_number);

-- ============================================================================
-- 总账管理
-- ============================================================================

-- 总账余额表
CREATE TABLE IF NOT EXISTS general_ledger_balances (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL,
  account_id BIGINT NOT NULL,
  fiscal_period_id BIGINT NOT NULL,
  opening_balance DECIMAL(15,2) DEFAULT 0,
  debit_amount DECIMAL(15,2) DEFAULT 0,
  credit_amount DECIMAL(15,2) DEFAULT 0,
  closing_balance DECIMAL(15,2) DEFAULT 0,
  currency VARCHAR(10) DEFAULT 'CNY',
  CONSTRAINT fk_general_ledger_balances_account FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE,
  CONSTRAINT fk_general_ledger_balances_fiscal_period FOREIGN KEY (fiscal_period_id) REFERENCES fiscal_periods(id) ON DELETE CASCADE,
  CONSTRAINT uq_general_ledger_balances_account_period UNIQUE (account_id, fiscal_period_id)
);
CREATE INDEX IF NOT EXISTS idx_general_ledger_balances_deleted_at ON general_ledger_balances (deleted_at);
CREATE INDEX IF NOT EXISTS idx_general_ledger_balances_account_id ON general_ledger_balances (account_id);
CREATE INDEX IF NOT EXISTS idx_general_ledger_balances_fiscal_period_id ON general_ledger_balances (fiscal_period_id);

-- ============================================================================
-- 应收账款管理
-- ============================================================================

-- 应收账款表
CREATE TABLE IF NOT EXISTS accounts_receivable (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL,
  created_by BIGINT NULL,
  updated_by BIGINT NULL,
  is_active BOOLEAN DEFAULT TRUE,
  code VARCHAR(50) NOT NULL,
  customer_id BIGINT NOT NULL,
  sales_invoice_id BIGINT NULL,
  invoice_number VARCHAR(100) NOT NULL,
  invoice_date DATE NOT NULL,
  due_date DATE NOT NULL,
  original_amount DECIMAL(15,2) NOT NULL,
  paid_amount DECIMAL(15,2) DEFAULT 0,
  balance_amount DECIMAL(15,2) NOT NULL,
  currency VARCHAR(10) DEFAULT 'CNY',
  exchange_rate DECIMAL(10,4) DEFAULT 1,
  payment_terms INTEGER DEFAULT 30,
  overdue_days INTEGER DEFAULT 0,
  aging_bucket VARCHAR(20) DEFAULT 'CURRENT',
  notes TEXT NULL,
  status VARCHAR(20) DEFAULT 'OPEN',
  CONSTRAINT uq_accounts_receivable_code UNIQUE (code),
  CONSTRAINT fk_accounts_receivable_customer FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE,
  CONSTRAINT fk_accounts_receivable_sales_invoice FOREIGN KEY (sales_invoice_id) REFERENCES sales_invoices(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_accounts_receivable_deleted_at ON accounts_receivable (deleted_at);
CREATE INDEX IF NOT EXISTS idx_accounts_receivable_is_active ON accounts_receivable (is_active);
CREATE INDEX IF NOT EXISTS idx_accounts_receivable_code ON accounts_receivable (code);
CREATE INDEX IF NOT EXISTS idx_accounts_receivable_customer_id ON accounts_receivable (customer_id);
CREATE INDEX IF NOT EXISTS idx_accounts_receivable_sales_invoice_id ON accounts_receivable (sales_invoice_id);
CREATE INDEX IF NOT EXISTS idx_accounts_receivable_due_date ON accounts_receivable (due_date);
CREATE INDEX IF NOT EXISTS idx_accounts_receivable_aging_bucket ON accounts_receivable (aging_bucket);
CREATE INDEX IF NOT EXISTS idx_accounts_receivable_status ON accounts_receivable (status);
CREATE INDEX IF NOT EXISTS idx_accounts_receivable_invoice_date ON accounts_receivable (invoice_date);

-- ============================================================================
-- 应付账款管理
-- ============================================================================

-- 应付账款表
CREATE TABLE IF NOT EXISTS accounts_payable (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL,
  created_by BIGINT NULL,
  updated_by BIGINT NULL,
  is_active BOOLEAN DEFAULT TRUE,
  code VARCHAR(50) NOT NULL,
  supplier_id BIGINT NOT NULL,
  purchase_invoice_id BIGINT NULL,
  invoice_number VARCHAR(100) NOT NULL,
  invoice_date DATE NOT NULL,
  due_date DATE NOT NULL,
  original_amount DECIMAL(15,2) NOT NULL,
  paid_amount DECIMAL(15,2) DEFAULT 0,
  balance_amount DECIMAL(15,2) NOT NULL,
  currency VARCHAR(10) DEFAULT 'CNY',
  exchange_rate DECIMAL(10,4) DEFAULT 1,
  payment_terms INTEGER DEFAULT 30,
  overdue_days INTEGER DEFAULT 0,
  aging_bucket VARCHAR(20) DEFAULT 'CURRENT',
  notes TEXT NULL,
  status VARCHAR(20) DEFAULT 'OPEN',
  CONSTRAINT uq_accounts_payable_code UNIQUE (code),
  CONSTRAINT fk_accounts_payable_supplier FOREIGN KEY (supplier_id) REFERENCES suppliers(id) ON DELETE CASCADE,
  CONSTRAINT fk_accounts_payable_purchase_invoice FOREIGN KEY (purchase_invoice_id) REFERENCES purchase_invoices(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_accounts_payable_deleted_at ON accounts_payable (deleted_at);
CREATE INDEX IF NOT EXISTS idx_accounts_payable_is_active ON accounts_payable (is_active);
CREATE INDEX IF NOT EXISTS idx_accounts_payable_code ON accounts_payable (code);
CREATE INDEX IF NOT EXISTS idx_accounts_payable_supplier_id ON accounts_payable (supplier_id);
CREATE INDEX IF NOT EXISTS idx_accounts_payable_purchase_invoice_id ON accounts_payable (purchase_invoice_id);
CREATE INDEX IF NOT EXISTS idx_accounts_payable_due_date ON accounts_payable (due_date);
CREATE INDEX IF NOT EXISTS idx_accounts_payable_aging_bucket ON accounts_payable (aging_bucket);
CREATE INDEX IF NOT EXISTS idx_accounts_payable_status ON accounts_payable (status);
CREATE INDEX IF NOT EXISTS idx_accounts_payable_invoice_date ON accounts_payable (invoice_date);

-- ============================================================================
-- 付款管理
-- ============================================================================

-- 付款单表
CREATE TABLE IF NOT EXISTS payments (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL,
  created_by BIGINT NULL,
  updated_by BIGINT NULL,
  is_active BOOLEAN DEFAULT TRUE,
  code VARCHAR(50) NOT NULL,
  payment_type VARCHAR(20) NOT NULL,
  payee_type VARCHAR(20) NOT NULL,
  payee_id BIGINT NOT NULL,
  payment_date DATE NOT NULL,
  amount DECIMAL(15,2) NOT NULL,
  currency VARCHAR(10) DEFAULT 'CNY',
  exchange_rate DECIMAL(10,4) DEFAULT 1,
  payment_method VARCHAR(20) NOT NULL,
  bank_account_id BIGINT NULL,
  reference_number VARCHAR(100) NULL,
  description TEXT NULL,
  notes TEXT NULL,
  status VARCHAR(20) DEFAULT 'PENDING',
  approved_by BIGINT NULL,
  approved_at TIMESTAMP NULL,
  CONSTRAINT uq_payments_code UNIQUE (code),
  CONSTRAINT fk_payments_bank_account FOREIGN KEY (bank_account_id) REFERENCES bank_accounts(id) ON DELETE SET NULL,
  CONSTRAINT fk_payments_approved_by FOREIGN KEY (approved_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT chk_payments_payment_type CHECK (payment_type IN ('SUPPLIER', 'EMPLOYEE', 'OTHER')),
  CONSTRAINT chk_payments_payee_type CHECK (payee_type IN ('SUPPLIER', 'EMPLOYEE', 'OTHER')),
  CONSTRAINT chk_payments_payment_method CHECK (payment_method IN ('CASH', 'BANK_TRANSFER', 'CHECK', 'CREDIT_CARD', 'OTHER'))
);
CREATE INDEX IF NOT EXISTS idx_payments_deleted_at ON payments (deleted_at);
CREATE INDEX IF NOT EXISTS idx_payments_is_active ON payments (is_active);
CREATE INDEX IF NOT EXISTS idx_payments_code ON payments (code);
CREATE INDEX IF NOT EXISTS idx_payments_payment_type ON payments (payment_type);
CREATE INDEX IF NOT EXISTS idx_payments_payee_type ON payments (payee_type);
CREATE INDEX IF NOT EXISTS idx_payments_payee_id ON payments (payee_id);
CREATE INDEX IF NOT EXISTS idx_payments_payment_date ON payments (payment_date);
CREATE INDEX IF NOT EXISTS idx_payments_status ON payments (status);
CREATE INDEX IF NOT EXISTS idx_payments_bank_account_id ON payments (bank_account_id);

-- 付款明细表
CREATE TABLE IF NOT EXISTS payment_items (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL,
  payment_id BIGINT NOT NULL,
  reference_type VARCHAR(50) NOT NULL,
  reference_id BIGINT NOT NULL,
  amount DECIMAL(15,2) NOT NULL,
  discount_amount DECIMAL(15,2) DEFAULT 0,
  notes TEXT NULL,
  CONSTRAINT fk_payment_items_payment FOREIGN KEY (payment_id) REFERENCES payments(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_payment_items_deleted_at ON payment_items (deleted_at);
CREATE INDEX IF NOT EXISTS idx_payment_items_payment_id ON payment_items (payment_id);
CREATE INDEX IF NOT EXISTS idx_payment_items_reference_type ON payment_items (reference_type);
CREATE INDEX IF NOT EXISTS idx_payment_items_reference_id ON payment_items (reference_id);

-- ============================================================================
-- 收款管理
-- ============================================================================

-- 收款单表
CREATE TABLE IF NOT EXISTS receipts_ar (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL,
  created_by BIGINT NULL,
  updated_by BIGINT NULL,
  is_active BOOLEAN DEFAULT TRUE,
  code VARCHAR(50) NOT NULL,
  customer_id BIGINT NOT NULL,
  receipt_date DATE NOT NULL,
  amount DECIMAL(15,2) NOT NULL,
  currency VARCHAR(10) DEFAULT 'CNY',
  exchange_rate DECIMAL(10,4) DEFAULT 1,
  payment_method VARCHAR(20) NOT NULL,
  bank_account_id BIGINT NULL,
  reference_number VARCHAR(100) NULL,
  description TEXT NULL,
  notes TEXT NULL,
  status VARCHAR(20) DEFAULT 'PENDING',
  approved_by BIGINT NULL,
  approved_at TIMESTAMP NULL,
  CONSTRAINT uq_receipts_ar_code UNIQUE (code),
  CONSTRAINT fk_receipts_ar_customer FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE,
  CONSTRAINT fk_receipts_ar_bank_account FOREIGN KEY (bank_account_id) REFERENCES bank_accounts(id) ON DELETE SET NULL,
  CONSTRAINT fk_receipts_ar_approved_by FOREIGN KEY (approved_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT chk_receipts_ar_payment_method CHECK (payment_method IN ('CASH', 'BANK_TRANSFER', 'CHECK', 'CREDIT_CARD', 'OTHER'))
);
CREATE INDEX IF NOT EXISTS idx_receipts_ar_deleted_at ON receipts_ar (deleted_at);
CREATE INDEX IF NOT EXISTS idx_receipts_ar_is_active ON receipts_ar (is_active);
CREATE INDEX IF NOT EXISTS idx_receipts_ar_code ON receipts_ar (code);
CREATE INDEX IF NOT EXISTS idx_receipts_ar_customer_id ON receipts_ar (customer_id);
CREATE INDEX IF NOT EXISTS idx_receipts_ar_receipt_date ON receipts_ar (receipt_date);
CREATE INDEX IF NOT EXISTS idx_receipts_ar_status ON receipts_ar (status);
CREATE INDEX IF NOT EXISTS idx_receipts_ar_bank_account_id ON receipts_ar (bank_account_id);

-- 收款明细表
CREATE TABLE IF NOT EXISTS receipt_ar_items (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL,
  receipt_id BIGINT NOT NULL,
  accounts_receivable_id BIGINT NOT NULL,
  amount DECIMAL(15,2) NOT NULL,
  discount_amount DECIMAL(15,2) DEFAULT 0,
  notes TEXT NULL,
  CONSTRAINT fk_receipt_ar_items_receipt FOREIGN KEY (receipt_id) REFERENCES receipts_ar(id) ON DELETE CASCADE,
  CONSTRAINT fk_receipt_ar_items_accounts_receivable FOREIGN KEY (accounts_receivable_id) REFERENCES accounts_receivable(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_receipt_ar_items_deleted_at ON receipt_ar_items (deleted_at);
CREATE INDEX IF NOT EXISTS idx_receipt_ar_items_receipt_id ON receipt_ar_items (receipt_id);
CREATE INDEX IF NOT EXISTS idx_receipt_ar_items_accounts_receivable_id ON receipt_ar_items (accounts_receivable_id);

-- ============================================================================
-- 银行账户管理
-- ============================================================================

-- 银行账户表
CREATE TABLE IF NOT EXISTS bank_accounts (
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
  bank_name VARCHAR(255) NOT NULL,
  bank_branch VARCHAR(255) NULL,
  account_number VARCHAR(100) NOT NULL,
  account_type VARCHAR(20) DEFAULT 'CHECKING',
  currency VARCHAR(10) DEFAULT 'CNY',
  opening_balance DECIMAL(15,2) DEFAULT 0,
  current_balance DECIMAL(15,2) DEFAULT 0,
  is_default BOOLEAN DEFAULT FALSE,
  status VARCHAR(20) DEFAULT 'ACTIVE',
  CONSTRAINT uq_bank_accounts_code UNIQUE (code),
  CONSTRAINT uq_bank_accounts_account_number UNIQUE (account_number),
  CONSTRAINT chk_bank_accounts_account_type CHECK (account_type IN ('CHECKING', 'SAVINGS', 'CREDIT', 'OTHER'))
);
CREATE INDEX IF NOT EXISTS idx_bank_accounts_deleted_at ON bank_accounts (deleted_at);
CREATE INDEX IF NOT EXISTS idx_bank_accounts_is_active ON bank_accounts (is_active);
CREATE INDEX IF NOT EXISTS idx_bank_accounts_code ON bank_accounts (code);
CREATE INDEX IF NOT EXISTS idx_bank_accounts_account_number ON bank_accounts (account_number);
CREATE INDEX IF NOT EXISTS idx_bank_accounts_currency ON bank_accounts (currency);
CREATE INDEX IF NOT EXISTS idx_bank_accounts_is_default ON bank_accounts (is_default);
CREATE INDEX IF NOT EXISTS idx_bank_accounts_status ON bank_accounts (status);
CREATE INDEX IF NOT EXISTS idx_bank_accounts_name ON bank_accounts (name);

-- ============================================================================
-- 初始数据插入
-- ============================================================================

-- 插入默认会计科目
INSERT INTO accounts (code, name, description, account_type, balance_type, level, is_leaf, is_system, status, is_active) VALUES
-- 资产类
('1000', '资产', '资产类科目', 'ASSET', 'DEBIT', 1, FALSE, TRUE, 'ACTIVE', TRUE),
('1001', '流动资产', '流动资产', 'ASSET', 'DEBIT', 2, FALSE, TRUE, 'ACTIVE', TRUE),
('100101', '库存现金', '库存现金', 'ASSET', 'DEBIT', 3, TRUE, TRUE, 'ACTIVE', TRUE),
('100102', '银行存款', '银行存款', 'ASSET', 'DEBIT', 3, TRUE, TRUE, 'ACTIVE', TRUE),
('100103', '应收账款', '应收账款', 'ASSET', 'DEBIT', 3, TRUE, TRUE, 'ACTIVE', TRUE),
('100104', '预付账款', '预付账款', 'ASSET', 'DEBIT', 3, TRUE, TRUE, 'ACTIVE', TRUE),
('100105', '存货', '存货', 'ASSET', 'DEBIT', 3, TRUE, TRUE, 'ACTIVE', TRUE),
-- 负债类
('2000', '负债', '负债类科目', 'LIABILITY', 'CREDIT', 1, FALSE, TRUE, 'ACTIVE', TRUE),
('2001', '流动负债', '流动负债', 'LIABILITY', 'CREDIT', 2, FALSE, TRUE, 'ACTIVE', TRUE),
('200101', '应付账款', '应付账款', 'LIABILITY', 'CREDIT', 3, TRUE, TRUE, 'ACTIVE', TRUE),
('200102', '预收账款', '预收账款', 'LIABILITY', 'CREDIT', 3, TRUE, TRUE, 'ACTIVE', TRUE),
('200103', '应付职工薪酬', '应付职工薪酬', 'LIABILITY', 'CREDIT', 3, TRUE, TRUE, 'ACTIVE', TRUE),
-- 所有者权益类
('3000', '所有者权益', '所有者权益类科目', 'EQUITY', 'CREDIT', 1, FALSE, TRUE, 'ACTIVE', TRUE),
('300101', '实收资本', '实收资本', 'EQUITY', 'CREDIT', 2, TRUE, TRUE, 'ACTIVE', TRUE),
('300102', '未分配利润', '未分配利润', 'EQUITY', 'CREDIT', 2, TRUE, TRUE, 'ACTIVE', TRUE),
-- 收入类
('4000', '收入', '收入类科目', 'REVENUE', 'CREDIT', 1, FALSE, TRUE, 'ACTIVE', TRUE),
('400101', '主营业务收入', '主营业务收入', 'REVENUE', 'CREDIT', 2, TRUE, TRUE, 'ACTIVE', TRUE),
('400102', '其他业务收入', '其他业务收入', 'REVENUE', 'CREDIT', 2, TRUE, TRUE, 'ACTIVE', TRUE),
-- 费用类
('5000', '费用', '费用类科目', 'EXPENSE', 'DEBIT', 1, FALSE, TRUE, 'ACTIVE', TRUE),
('500101', '主营业务成本', '主营业务成本', 'EXPENSE', 'DEBIT', 2, TRUE, TRUE, 'ACTIVE', TRUE),
('500102', '销售费用', '销售费用', 'EXPENSE', 'DEBIT', 2, TRUE, TRUE, 'ACTIVE', TRUE),
('500103', '管理费用', '管理费用', 'EXPENSE', 'DEBIT', 2, TRUE, TRUE, 'ACTIVE', TRUE)
ON CONFLICT (code) DO NOTHING;

-- 更新父子关系
UPDATE accounts SET parent_id = (SELECT id FROM accounts WHERE code = '1000') WHERE code IN ('1001');
UPDATE accounts SET parent_id = (SELECT id FROM accounts WHERE code = '1001') WHERE code IN ('100101', '100102', '100103', '100104', '100105');
UPDATE accounts SET parent_id = (SELECT id FROM accounts WHERE code = '2000') WHERE code IN ('2001');
UPDATE accounts SET parent_id = (SELECT id FROM accounts WHERE code = '2001') WHERE code IN ('200101', '200102', '200103');

-- 插入当前会计期间
INSERT INTO fiscal_periods (code, name, year, period, start_date, end_date, status, is_active) VALUES
('2024-01', '2024年1月', 2024, 1, '2024-01-01', '2024-01-31', 'OPEN', TRUE),
('2024-02', '2024年2月', 2024, 2, '2024-02-01', '2024-02-29', 'OPEN', TRUE),
('2024-03', '2024年3月', 2024, 3, '2024-03-01', '2024-03-31', 'OPEN', TRUE),
('2024-04', '2024年4月', 2024, 4, '2024-04-01', '2024-04-30', 'OPEN', TRUE),
('2024-05', '2024年5月', 2024, 5, '2024-05-01', '2024-05-31', 'OPEN', TRUE),
('2024-06', '2024年6月', 2024, 6, '2024-06-01', '2024-06-30', 'OPEN', TRUE),
('2024-07', '2024年7月', 2024, 7, '2024-07-01', '2024-07-31', 'OPEN', TRUE),
('2024-08', '2024年8月', 2024, 8, '2024-08-01', '2024-08-31', 'OPEN', TRUE),
('2024-09', '2024年9月', 2024, 9, '2024-09-01', '2024-09-30', 'OPEN', TRUE),
('2024-10', '2024年10月', 2024, 10, '2024-10-01', '2024-10-31', 'OPEN', TRUE),
('2024-11', '2024年11月', 2024, 11, '2024-11-01', '2024-11-30', 'OPEN', TRUE),
('2024-12', '2024年12月', 2024, 12, '2024-12-01', '2024-12-31', 'OPEN', TRUE)
ON CONFLICT (code) DO NOTHING;

-- 插入默认银行账户
INSERT INTO bank_accounts (code, name, description, bank_name, account_number, account_type, currency, opening_balance, current_balance, is_default, status, is_active) VALUES
('BANK_DEFAULT', '默认银行账户', '系统默认银行账户', '中国银行', '1234567890123456789', 'CHECKING', 'CNY', 1000000.00, 1000000.00, TRUE, 'ACTIVE', TRUE)
ON CONFLICT (code) DO NOTHING;

-- ============================================================================
-- 结束
-- ============================================================================