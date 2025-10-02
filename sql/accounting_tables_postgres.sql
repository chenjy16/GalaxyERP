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
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE NULL,
  code VARCHAR(50) NOT NULL,
  name VARCHAR(255) NOT NULL,
  description TEXT NULL,
  account_type VARCHAR(50) NOT NULL,
  balance DECIMAL(15,2) DEFAULT 0,
  is_active BOOLEAN DEFAULT TRUE,
  parent_id INTEGER NULL,
  currency VARCHAR(10) DEFAULT 'USD',
  CONSTRAINT uq_accounts_code UNIQUE (code),
  CONSTRAINT fk_accounts_parent FOREIGN KEY (parent_id) REFERENCES accounts(id) ON DELETE SET NULL,
  CONSTRAINT chk_accounts_account_type CHECK (account_type IN ('ASSET', 'LIABILITY', 'EQUITY', 'REVENUE', 'EXPENSE'))
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
  created_at TIMESTAMP WITH TIME ZONE NULL,
  updated_at TIMESTAMP WITH TIME ZONE NULL,
  deleted_at TIMESTAMP WITH TIME ZONE NULL,
  created_by BIGINT NULL,
  updated_by BIGINT NULL,
  is_active BOOLEAN DEFAULT TRUE,
  code VARCHAR(50) NOT NULL,
  name VARCHAR(255) NOT NULL,
  description TEXT NULL,
  year INTEGER NOT NULL,
  period INTEGER NOT NULL,
  start_date TIMESTAMP WITH TIME ZONE NOT NULL,
  end_date TIMESTAMP WITH TIME ZONE NOT NULL,
  is_closed BOOLEAN DEFAULT FALSE,
  closed_by BIGINT NULL,
  closed_at TIMESTAMP WITH TIME ZONE NULL,
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

-- 会计分录表
CREATE TABLE IF NOT EXISTS journal_entries (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE NULL,
  transaction_id INTEGER NOT NULL,
  account_id INTEGER NOT NULL,
  debit DECIMAL(15,2) DEFAULT 0,
  credit DECIMAL(15,2) DEFAULT 0,
  description TEXT NULL,
  CONSTRAINT fk_journal_entries_account FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_journal_entries_deleted_at ON journal_entries (deleted_at);
CREATE INDEX IF NOT EXISTS idx_journal_entries_transaction_id ON journal_entries (transaction_id);
CREATE INDEX IF NOT EXISTS idx_journal_entries_account_id ON journal_entries (account_id);

-- 应收账款表
CREATE TABLE IF NOT EXISTS receivables (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE NULL,
  customer_id INTEGER NOT NULL,
  invoice_date TIMESTAMP WITH TIME ZONE NOT NULL,
  due_date TIMESTAMP WITH TIME ZONE NOT NULL,
  invoice_number VARCHAR(100) NOT NULL,
  description TEXT NULL,
  amount DECIMAL(15,2) NOT NULL,
  amount_paid DECIMAL(15,2) DEFAULT 0,
  currency VARCHAR(10) DEFAULT 'USD',
  exchange_rate DECIMAL(10,4) DEFAULT 1,
  status VARCHAR(50) DEFAULT 'open',
  CONSTRAINT fk_receivables_customer FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_receivables_deleted_at ON receivables (deleted_at);
CREATE INDEX IF NOT EXISTS idx_receivables_customer_id ON receivables (customer_id);
CREATE INDEX IF NOT EXISTS idx_receivables_invoice_date ON receivables (invoice_date);
CREATE INDEX IF NOT EXISTS idx_receivables_due_date ON receivables (due_date);
CREATE INDEX IF NOT EXISTS idx_receivables_status ON receivables (status);

-- 应付账款表
CREATE TABLE IF NOT EXISTS payables (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE NULL,
  supplier_id INTEGER NOT NULL,
  invoice_date TIMESTAMP WITH TIME ZONE NOT NULL,
  due_date TIMESTAMP WITH TIME ZONE NOT NULL,
  invoice_number VARCHAR(100) NOT NULL,
  description TEXT NULL,
  amount DECIMAL(15,2) NOT NULL,
  amount_paid DECIMAL(15,2) DEFAULT 0,
  currency VARCHAR(10) DEFAULT 'USD',
  exchange_rate DECIMAL(10,4) DEFAULT 1,
  status VARCHAR(50) DEFAULT 'open',
  CONSTRAINT fk_payables_supplier FOREIGN KEY (supplier_id) REFERENCES suppliers(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_payables_deleted_at ON payables (deleted_at);
CREATE INDEX IF NOT EXISTS idx_payables_supplier_id ON payables (supplier_id);
CREATE INDEX IF NOT EXISTS idx_payables_invoice_date ON payables (invoice_date);
CREATE INDEX IF NOT EXISTS idx_payables_due_date ON payables (due_date);
CREATE INDEX IF NOT EXISTS idx_payables_status ON payables (status);

-- 付款表
CREATE TABLE IF NOT EXISTS payments (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE NULL,
  created_by INTEGER NULL,
  updated_by INTEGER NULL,
  payment_number VARCHAR(100) NOT NULL,
  payment_date TIMESTAMP WITH TIME ZONE NOT NULL,
  payment_type VARCHAR(50) NOT NULL,
  amount DECIMAL(15,2) NOT NULL,
  currency VARCHAR(10) DEFAULT 'CNY',
  exchange_rate DECIMAL(10,4) DEFAULT 1,
  payer_type VARCHAR(50) NOT NULL,
  payer_id INTEGER NOT NULL,
  payee_type VARCHAR(50) NOT NULL,
  payee_id INTEGER NOT NULL,
  bank_account_id INTEGER NULL,
  reference_type VARCHAR(50) NULL,
  reference_id INTEGER NULL,
  status VARCHAR(50) DEFAULT 'pending',
  notes TEXT NULL,
  CONSTRAINT uq_payments_payment_number UNIQUE (payment_number),
  CONSTRAINT fk_payments_bank_account FOREIGN KEY (bank_account_id) REFERENCES bank_accounts(id) ON DELETE SET NULL,
  CONSTRAINT fk_payments_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_payments_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_payments_deleted_at ON payments (deleted_at);
CREATE INDEX IF NOT EXISTS idx_payments_payment_number ON payments (payment_number);
CREATE INDEX IF NOT EXISTS idx_payments_payment_date ON payments (payment_date);
CREATE INDEX IF NOT EXISTS idx_payments_payment_type ON payments (payment_type);
CREATE INDEX IF NOT EXISTS idx_payments_payer_type ON payments (payer_type);
CREATE INDEX IF NOT EXISTS idx_payments_payer_id ON payments (payer_id);
CREATE INDEX IF NOT EXISTS idx_payments_payee_type ON payments (payee_type);
CREATE INDEX IF NOT EXISTS idx_payments_payee_id ON payments (payee_id);
CREATE INDEX IF NOT EXISTS idx_payments_status ON payments (status);
CREATE INDEX IF NOT EXISTS idx_payments_bank_account_id ON payments (bank_account_id);

-- 银行账户表
CREATE TABLE IF NOT EXISTS bank_accounts (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE NULL,
  account_name VARCHAR(255) NOT NULL,
  bank_name VARCHAR(255) NOT NULL,
  account_number VARCHAR(100) NOT NULL,
  iban VARCHAR(100) NULL,
  swift_code VARCHAR(50) NULL,
  currency VARCHAR(10) DEFAULT 'USD',
  balance DECIMAL(15,2) DEFAULT 0,
  is_default BOOLEAN DEFAULT FALSE,
  is_active BOOLEAN DEFAULT TRUE,
  account_id INTEGER NULL,
  CONSTRAINT uq_bank_accounts_account_number UNIQUE (account_number),
  CONSTRAINT fk_bank_accounts_account FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_bank_accounts_deleted_at ON bank_accounts (deleted_at);
CREATE INDEX IF NOT EXISTS idx_bank_accounts_account_number ON bank_accounts (account_number);
CREATE INDEX IF NOT EXISTS idx_bank_accounts_currency ON bank_accounts (currency);
CREATE INDEX IF NOT EXISTS idx_bank_accounts_is_default ON bank_accounts (is_default);
CREATE INDEX IF NOT EXISTS idx_bank_accounts_is_active ON bank_accounts (is_active);
CREATE INDEX IF NOT EXISTS idx_bank_accounts_account_name ON bank_accounts (account_name);
CREATE INDEX IF NOT EXISTS idx_bank_accounts_account_id ON bank_accounts (account_id);

-- 预算表
CREATE TABLE IF NOT EXISTS budgets (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE NULL,
  created_by INTEGER NULL,
  updated_by INTEGER NULL,
  budget_name VARCHAR(255) NOT NULL,
  budget_year INTEGER NOT NULL,
  start_date TIMESTAMP WITH TIME ZONE NOT NULL,
  end_date TIMESTAMP WITH TIME ZONE NOT NULL,
  total_amount DECIMAL(15,2) NOT NULL,
  used_amount DECIMAL(15,2) DEFAULT 0,
  remaining_amount DECIMAL(15,2) DEFAULT 0,
  status VARCHAR(50) DEFAULT 'draft',
  CONSTRAINT fk_budgets_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_budgets_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_budgets_deleted_at ON budgets (deleted_at);
CREATE INDEX IF NOT EXISTS idx_budgets_budget_year ON budgets (budget_year);
CREATE INDEX IF NOT EXISTS idx_budgets_start_date ON budgets (start_date);
CREATE INDEX IF NOT EXISTS idx_budgets_end_date ON budgets (end_date);
CREATE INDEX IF NOT EXISTS idx_budgets_status ON budgets (status);

-- 预算明细表
CREATE TABLE IF NOT EXISTS budget_items (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE NULL,
  budget_id INTEGER NOT NULL,
  account_id INTEGER NOT NULL,
  budget_amount DECIMAL(15,2) NOT NULL,
  actual_amount DECIMAL(15,2) DEFAULT 0,
  variance_amount DECIMAL(15,2) DEFAULT 0,
  notes TEXT NULL,
  CONSTRAINT fk_budget_items_budget FOREIGN KEY (budget_id) REFERENCES budgets(id) ON DELETE CASCADE,
  CONSTRAINT fk_budget_items_account FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_budget_items_deleted_at ON budget_items (deleted_at);
CREATE INDEX IF NOT EXISTS idx_budget_items_budget_id ON budget_items (budget_id);
CREATE INDEX IF NOT EXISTS idx_budget_items_account_id ON budget_items (account_id);

-- 成本中心表
CREATE TABLE IF NOT EXISTS cost_centers (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE NULL,
  code VARCHAR(50) NOT NULL,
  name VARCHAR(255) NOT NULL,
  description TEXT NULL,
  is_active BOOLEAN DEFAULT TRUE,
  manager_id INTEGER NULL,
  department_id INTEGER NULL,
  CONSTRAINT uq_cost_centers_code UNIQUE (code),
  CONSTRAINT fk_cost_centers_manager FOREIGN KEY (manager_id) REFERENCES employees(id) ON DELETE SET NULL,
  CONSTRAINT fk_cost_centers_department FOREIGN KEY (department_id) REFERENCES departments(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_cost_centers_deleted_at ON cost_centers (deleted_at);
CREATE INDEX IF NOT EXISTS idx_cost_centers_code ON cost_centers (code);
CREATE INDEX IF NOT EXISTS idx_cost_centers_is_active ON cost_centers (is_active);
CREATE INDEX IF NOT EXISTS idx_cost_centers_manager_id ON cost_centers (manager_id);
CREATE INDEX IF NOT EXISTS idx_cost_centers_department_id ON cost_centers (department_id);

-- 财务报表表
CREATE TABLE IF NOT EXISTS financial_reports (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE NULL,
  created_by INTEGER NULL,
  updated_by INTEGER NULL,
  report_name VARCHAR(255) NOT NULL,
  report_type VARCHAR(50) NOT NULL,
  period_type VARCHAR(50) NOT NULL,
  start_date TIMESTAMP WITH TIME ZONE NOT NULL,
  end_date TIMESTAMP WITH TIME ZONE NOT NULL,
  status VARCHAR(50) DEFAULT 'draft',
  file_path VARCHAR(500) NULL,
  CONSTRAINT fk_financial_reports_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_financial_reports_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_financial_reports_deleted_at ON financial_reports (deleted_at);
CREATE INDEX IF NOT EXISTS idx_financial_reports_report_type ON financial_reports (report_type);
CREATE INDEX IF NOT EXISTS idx_financial_reports_period_type ON financial_reports (period_type);
CREATE INDEX IF NOT EXISTS idx_financial_reports_start_date ON financial_reports (start_date);
CREATE INDEX IF NOT EXISTS idx_financial_reports_end_date ON financial_reports (end_date);
CREATE INDEX IF NOT EXISTS idx_financial_reports_status ON financial_reports (status);

-- 财务报表明细表
CREATE TABLE IF NOT EXISTS financial_report_items (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE NULL,
  report_id INTEGER NOT NULL,
  account_id INTEGER NOT NULL,
  amount DECIMAL(15,2) NOT NULL,
  percentage DECIMAL(5,2) DEFAULT 0,
  CONSTRAINT fk_financial_report_items_report FOREIGN KEY (report_id) REFERENCES financial_reports(id) ON DELETE CASCADE,
  CONSTRAINT fk_financial_report_items_account FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_financial_report_items_deleted_at ON financial_report_items (deleted_at);
CREATE INDEX IF NOT EXISTS idx_financial_report_items_report_id ON financial_report_items (report_id);
CREATE INDEX IF NOT EXISTS idx_financial_report_items_account_id ON financial_report_items (account_id);

-- 交易记录表
CREATE TABLE IF NOT EXISTS transactions (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE NULL,
  created_by INTEGER NULL,
  updated_by INTEGER NULL,
  transaction_number VARCHAR(100) NOT NULL,
  transaction_date TIMESTAMP WITH TIME ZONE NOT NULL,
  transaction_type VARCHAR(50) NOT NULL,
  amount DECIMAL(15,2) NOT NULL,
  currency VARCHAR(10) DEFAULT 'CNY',
  exchange_rate DECIMAL(10,4) DEFAULT 1,
  description TEXT NOT NULL,
  reference_type VARCHAR(50) NULL,
  reference_id INTEGER NULL,
  status VARCHAR(50) DEFAULT 'pending',
  notes TEXT NULL,
  CONSTRAINT uq_transactions_transaction_number UNIQUE (transaction_number),
  CONSTRAINT fk_transactions_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_transactions_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_transactions_deleted_at ON transactions (deleted_at);
CREATE INDEX IF NOT EXISTS idx_transactions_transaction_number ON transactions (transaction_number);
CREATE INDEX IF NOT EXISTS idx_transactions_transaction_date ON transactions (transaction_date);
CREATE INDEX IF NOT EXISTS idx_transactions_transaction_type ON transactions (transaction_type);
CREATE INDEX IF NOT EXISTS idx_transactions_reference_type ON transactions (reference_type);
CREATE INDEX IF NOT EXISTS idx_transactions_reference_id ON transactions (reference_id);
CREATE INDEX IF NOT EXISTS idx_transactions_status ON transactions (status);

-- 固定资产表
CREATE TABLE IF NOT EXISTS fixed_assets (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE NULL,
  created_by INTEGER NULL,
  updated_by INTEGER NULL,
  asset_number VARCHAR(100) NOT NULL,
  asset_name VARCHAR(255) NOT NULL,
  asset_category VARCHAR(100) NOT NULL,
  purchase_date TIMESTAMP WITH TIME ZONE NOT NULL,
  purchase_price DECIMAL(15,2) NOT NULL,
  current_value DECIMAL(15,2) DEFAULT 0,
  depreciation_rate DECIMAL(5,2) DEFAULT 0,
  useful_life INTEGER DEFAULT 0,
  location VARCHAR(255) NULL,
  status VARCHAR(50) DEFAULT 'active',
  notes TEXT NULL,
  CONSTRAINT uq_fixed_assets_asset_number UNIQUE (asset_number),
  CONSTRAINT fk_fixed_assets_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_fixed_assets_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_fixed_assets_deleted_at ON fixed_assets (deleted_at);
CREATE INDEX IF NOT EXISTS idx_fixed_assets_asset_number ON fixed_assets (asset_number);
CREATE INDEX IF NOT EXISTS idx_fixed_assets_asset_category ON fixed_assets (asset_category);
CREATE INDEX IF NOT EXISTS idx_fixed_assets_purchase_date ON fixed_assets (purchase_date);
CREATE INDEX IF NOT EXISTS idx_fixed_assets_status ON fixed_assets (status);

-- 折旧记录表
CREATE TABLE IF NOT EXISTS depreciation_entries (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE NULL,
  created_by INTEGER NULL,
  updated_by INTEGER NULL,
  asset_id INTEGER NOT NULL,
  depreciation_date TIMESTAMP WITH TIME ZONE NOT NULL,
  depreciation_amount DECIMAL(15,2) NOT NULL,
  accumulated_amount DECIMAL(15,2) DEFAULT 0,
  book_value DECIMAL(15,2) DEFAULT 0,
  method VARCHAR(50) NOT NULL,
  notes TEXT NULL,
  CONSTRAINT fk_depreciation_entries_asset FOREIGN KEY (asset_id) REFERENCES fixed_assets(id) ON DELETE CASCADE,
  CONSTRAINT fk_depreciation_entries_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_depreciation_entries_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_depreciation_entries_deleted_at ON depreciation_entries (deleted_at);
CREATE INDEX IF NOT EXISTS idx_depreciation_entries_asset_id ON depreciation_entries (asset_id);
CREATE INDEX IF NOT EXISTS idx_depreciation_entries_depreciation_date ON depreciation_entries (depreciation_date);

-- 税务记录表
CREATE TABLE IF NOT EXISTS tax_entries (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE NULL,
  created_by INTEGER NULL,
  updated_by INTEGER NULL,
  tax_number VARCHAR(100) NOT NULL,
  tax_date TIMESTAMP WITH TIME ZONE NOT NULL,
  tax_type VARCHAR(50) NOT NULL,
  taxable_amount DECIMAL(15,2) NOT NULL,
  tax_rate DECIMAL(5,2) NOT NULL,
  tax_amount DECIMAL(15,2) NOT NULL,
  status VARCHAR(50) DEFAULT 'pending',
  notes TEXT NULL,
  CONSTRAINT uq_tax_entries_tax_number UNIQUE (tax_number),
  CONSTRAINT fk_tax_entries_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_tax_entries_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_tax_entries_deleted_at ON tax_entries (deleted_at);
CREATE INDEX IF NOT EXISTS idx_tax_entries_tax_number ON tax_entries (tax_number);
CREATE INDEX IF NOT EXISTS idx_tax_entries_tax_date ON tax_entries (tax_date);
CREATE INDEX IF NOT EXISTS idx_tax_entries_tax_type ON tax_entries (tax_type);
CREATE INDEX IF NOT EXISTS idx_tax_entries_status ON tax_entries (status);

-- 货币表
CREATE TABLE IF NOT EXISTS currencies (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE NULL,
  created_by INTEGER NULL,
  updated_by INTEGER NULL,
  currency_code VARCHAR(3) NOT NULL,
  currency_name VARCHAR(100) NOT NULL,
  symbol VARCHAR(10) NOT NULL,
  exchange_rate DECIMAL(10,4) DEFAULT 1.0000,
  is_active BOOLEAN DEFAULT true,
  CONSTRAINT uq_currencies_currency_code UNIQUE (currency_code),
  CONSTRAINT fk_currencies_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_currencies_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_currencies_deleted_at ON currencies (deleted_at);
CREATE INDEX IF NOT EXISTS idx_currencies_currency_code ON currencies (currency_code);
CREATE INDEX IF NOT EXISTS idx_currencies_is_active ON currencies (is_active);

-- 会计年度表
CREATE TABLE IF NOT EXISTS fiscal_years (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE NULL,
  created_by INTEGER NULL,
  updated_by INTEGER NULL,
  year_name VARCHAR(100) NOT NULL,
  start_date TIMESTAMP WITH TIME ZONE NOT NULL,
  end_date TIMESTAMP WITH TIME ZONE NOT NULL,
  is_current BOOLEAN DEFAULT false,
  is_closed BOOLEAN DEFAULT false,
  CONSTRAINT uq_fiscal_years_year_name UNIQUE (year_name),
  CONSTRAINT fk_fiscal_years_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_fiscal_years_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_fiscal_years_deleted_at ON fiscal_years (deleted_at);
CREATE INDEX IF NOT EXISTS idx_fiscal_years_start_date ON fiscal_years (start_date);
CREATE INDEX IF NOT EXISTS idx_fiscal_years_end_date ON fiscal_years (end_date);
CREATE INDEX IF NOT EXISTS idx_fiscal_years_is_current ON fiscal_years (is_current);

-- 会计期间表
CREATE TABLE IF NOT EXISTS accounting_periods (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE NULL,
  created_by INTEGER NULL,
  updated_by INTEGER NULL,
  fiscal_year_id INTEGER NOT NULL,
  period_name VARCHAR(100) NOT NULL,
  start_date TIMESTAMP WITH TIME ZONE NOT NULL,
  end_date TIMESTAMP WITH TIME ZONE NOT NULL,
  is_current BOOLEAN DEFAULT false,
  is_closed BOOLEAN DEFAULT false,
  CONSTRAINT fk_accounting_periods_fiscal_year FOREIGN KEY (fiscal_year_id) REFERENCES fiscal_years(id) ON DELETE CASCADE,
  CONSTRAINT fk_accounting_periods_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_accounting_periods_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_accounting_periods_deleted_at ON accounting_periods (deleted_at);
CREATE INDEX IF NOT EXISTS idx_accounting_periods_fiscal_year_id ON accounting_periods (fiscal_year_id);
CREATE INDEX IF NOT EXISTS idx_accounting_periods_start_date ON accounting_periods (start_date);
CREATE INDEX IF NOT EXISTS idx_accounting_periods_end_date ON accounting_periods (end_date);
CREATE INDEX IF NOT EXISTS idx_accounting_periods_is_current ON accounting_periods (is_current);

-- ============================================================================
-- 初始数据插入
-- ============================================================================

-- 插入默认会计科目
INSERT INTO accounts (code, name, description, account_type, balance, is_active, parent_id, currency) VALUES
-- 资产类
('1000', '资产', '资产类科目', 'ASSET', 0, TRUE, NULL, 'USD'),
('1001', '流动资产', '流动资产', 'ASSET', 0, TRUE, NULL, 'USD'),
('100101', '库存现金', '库存现金', 'ASSET', 0, TRUE, NULL, 'USD'),
('100102', '银行存款', '银行存款', 'ASSET', 0, TRUE, NULL, 'USD'),
('100103', '应收账款', '应收账款', 'ASSET', 0, TRUE, NULL, 'USD'),
('100104', '预付账款', '预付账款', 'ASSET', 0, TRUE, NULL, 'USD'),
('100105', '存货', '存货', 'ASSET', 0, TRUE, NULL, 'USD'),
-- 负债类
('2000', '负债', '负债类科目', 'LIABILITY', 0, TRUE, NULL, 'USD'),
('2001', '流动负债', '流动负债', 'LIABILITY', 0, TRUE, NULL, 'USD'),
('200101', '应付账款', '应付账款', 'LIABILITY', 0, TRUE, NULL, 'USD'),
('200102', '预收账款', '预收账款', 'LIABILITY', 0, TRUE, NULL, 'USD'),
('200103', '应付职工薪酬', '应付职工薪酬', 'LIABILITY', 0, TRUE, NULL, 'USD'),
-- 所有者权益类
('3000', '所有者权益', '所有者权益类科目', 'EQUITY', 0, TRUE, NULL, 'USD'),
('300101', '实收资本', '实收资本', 'EQUITY', 0, TRUE, NULL, 'USD'),
('300102', '未分配利润', '未分配利润', 'EQUITY', 0, TRUE, NULL, 'USD'),
-- 收入类
('4000', '收入', '收入类科目', 'REVENUE', 0, TRUE, NULL, 'USD'),
('400101', '主营业务收入', '主营业务收入', 'REVENUE', 0, TRUE, NULL, 'USD'),
('400102', '其他业务收入', '其他业务收入', 'REVENUE', 0, TRUE, NULL, 'USD'),
-- 费用类
('5000', '费用', '费用类科目', 'EXPENSE', 0, TRUE, NULL, 'USD'),
('500101', '主营业务成本', '主营业务成本', 'EXPENSE', 0, TRUE, NULL, 'USD'),
('500102', '销售费用', '销售费用', 'EXPENSE', 0, TRUE, NULL, 'USD'),
('500103', '管理费用', '管理费用', 'EXPENSE', 0, TRUE, NULL, 'USD')
ON CONFLICT (code) DO NOTHING;

-- 更新父子关系
UPDATE accounts SET parent_id = (SELECT id FROM accounts WHERE code = '1000') WHERE code IN ('1001');

-- 插入示例数据
INSERT INTO bank_accounts (account_name, bank_name, account_number, currency, balance, is_default, is_active, created_at, updated_at) VALUES
('主营业务账户', '中国银行', '6217001234567890123', 'CNY', 100000.00, TRUE, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('美元账户', '中国银行', '6217001234567890124', 'USD', 15000.00, FALSE, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('备用账户', '工商银行', '6222001234567890125', 'CNY', 50000.00, FALSE, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT (account_number) DO NOTHING;

INSERT INTO currencies (currency_code, currency_name, symbol, exchange_rate, is_active, created_at, updated_at) VALUES
('CNY', '人民币', '¥', 1.0000, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('USD', '美元', '$', 7.2500, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('EUR', '欧元', '€', 7.8500, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('JPY', '日元', '¥', 0.0650, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT (currency_code) DO NOTHING;

INSERT INTO fiscal_years (year_name, start_date, end_date, is_current, is_closed, created_at, updated_at) VALUES
('2024年度', '2024-01-01 00:00:00+00', '2024-12-31 23:59:59+00', TRUE, FALSE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('2023年度', '2023-01-01 00:00:00+00', '2023-12-31 23:59:59+00', FALSE, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('2025年度', '2025-01-01 00:00:00+00', '2025-12-31 23:59:59+00', FALSE, FALSE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT (year_name) DO NOTHING;
UPDATE accounts SET parent_id = (SELECT id FROM accounts WHERE code = '1001') WHERE code IN ('100101', '100102', '100103', '100104', '100105');
UPDATE accounts SET parent_id = (SELECT id FROM accounts WHERE code = '2000') WHERE code IN ('2001');
UPDATE accounts SET parent_id = (SELECT id FROM accounts WHERE code = '2001') WHERE code IN ('200101', '200102', '200103');

-- 插入当前会计期间
INSERT INTO accounting_periods (fiscal_year_id, period_name, start_date, end_date, is_current, is_closed, created_at, updated_at) VALUES
((SELECT id FROM fiscal_years WHERE year_name = '2024年度'), '2024年1月', '2024-01-01 00:00:00+00', '2024-01-31 23:59:59+00', FALSE, FALSE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
((SELECT id FROM fiscal_years WHERE year_name = '2024年度'), '2024年2月', '2024-02-01 00:00:00+00', '2024-02-29 23:59:59+00', FALSE, FALSE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
((SELECT id FROM fiscal_years WHERE year_name = '2024年度'), '2024年3月', '2024-03-01 00:00:00+00', '2024-03-31 23:59:59+00', FALSE, FALSE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
((SELECT id FROM fiscal_years WHERE year_name = '2024年度'), '2024年4月', '2024-04-01 00:00:00+00', '2024-04-30 23:59:59+00', FALSE, FALSE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
((SELECT id FROM fiscal_years WHERE year_name = '2024年度'), '2024年5月', '2024-05-01 00:00:00+00', '2024-05-31 23:59:59+00', FALSE, FALSE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
((SELECT id FROM fiscal_years WHERE year_name = '2024年度'), '2024年6月', '2024-06-01 00:00:00+00', '2024-06-30 23:59:59+00', FALSE, FALSE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
((SELECT id FROM fiscal_years WHERE year_name = '2024年度'), '2024年7月', '2024-07-01 00:00:00+00', '2024-07-31 23:59:59+00', FALSE, FALSE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
((SELECT id FROM fiscal_years WHERE year_name = '2024年度'), '2024年8月', '2024-08-01 00:00:00+00', '2024-08-31 23:59:59+00', FALSE, FALSE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
((SELECT id FROM fiscal_years WHERE year_name = '2024年度'), '2024年9月', '2024-09-01 00:00:00+00', '2024-09-30 23:59:59+00', FALSE, FALSE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
((SELECT id FROM fiscal_years WHERE year_name = '2024年度'), '2024年10月', '2024-10-01 00:00:00+00', '2024-10-31 23:59:59+00', FALSE, FALSE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
((SELECT id FROM fiscal_years WHERE year_name = '2024年度'), '2024年11月', '2024-11-01 00:00:00+00', '2024-11-30 23:59:59+00', FALSE, FALSE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
((SELECT id FROM fiscal_years WHERE year_name = '2024年度'), '2024年12月', '2024-12-01 00:00:00+00', '2024-12-31 23:59:59+00', TRUE, FALSE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- 插入默认银行账户
INSERT INTO bank_accounts (code, name, description, bank_name, account_number, account_type, currency, opening_balance, current_balance, is_default, status, is_active) VALUES
('BANK_DEFAULT', '默认银行账户', '系统默认银行账户', '中国银行', '1234567890123456789', 'CHECKING', 'CNY', 1000000.00, 1000000.00, TRUE, 'ACTIVE', TRUE)
ON CONFLICT (code) DO NOTHING;

-- ============================================================================
-- 结束
-- ============================================================================