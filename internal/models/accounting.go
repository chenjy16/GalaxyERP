package models

import (
	"time"
)

// Account 会计科目模型
type Account struct {
	BaseModel
	Code        string  `json:"code" gorm:"uniqueIndex;not null"`
	Name        string  `json:"name" gorm:"not null"`
	Description string  `json:"description,omitempty"`
	AccountType string  `json:"account_type" gorm:"not null"`
	Balance     float64 `json:"balance" gorm:"default:0"`
	IsActive    bool    `json:"is_active" gorm:"default:true"`
	ParentID    *uint   `json:"parent_id,omitempty"`
	Currency    string  `json:"currency" gorm:"default:'USD'"`
	
	// 关联
	Parent   *Account  `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Children []Account `json:"children,omitempty" gorm:"foreignKey:ParentID"`
}

// JournalEntry 会计分录模型
type JournalEntry struct {
	BaseModel
	TransactionID uint    `json:"transaction_id" gorm:"not null"`
	AccountID     uint    `json:"account_id" gorm:"not null"`
	Debit         float64 `json:"debit,omitempty"`
	Credit        float64 `json:"credit,omitempty"`
	Description   string  `json:"description,omitempty"`
	
	// 关联
	Account Account `json:"account,omitempty" gorm:"foreignKey:AccountID"`
}

// Payment 付款模型
type Payment struct {
	AuditableModel
	PaymentNumber string    `json:"payment_number" gorm:"uniqueIndex;size:100;not null"`
	PaymentDate   time.Time `json:"payment_date" gorm:"index;not null"`
	PaymentType   string    `json:"payment_type" gorm:"size:50;not null;index"` // cash, bank, check, card
	Amount        float64   `json:"amount" gorm:"not null"`
	Currency      string    `json:"currency" gorm:"size:10;default:'CNY'"`
	ExchangeRate  float64   `json:"exchange_rate" gorm:"default:1"`
	PayerType     string    `json:"payer_type" gorm:"size:50;not null;index"` // customer, supplier, employee
	PayerID       uint      `json:"payer_id" gorm:"index;not null"`
	PayeeType     string    `json:"payee_type" gorm:"size:50;not null;index"` // customer, supplier, employee
	PayeeID       uint      `json:"payee_id" gorm:"index;not null"`
	BankAccountID *uint     `json:"bank_account_id,omitempty" gorm:"index"`
	ReferenceType string    `json:"reference_type,omitempty" gorm:"size:50;index"` // invoice, purchase_order, salary
	ReferenceID   *uint     `json:"reference_id,omitempty" gorm:"index"`
	Status        string    `json:"status" gorm:"size:50;default:'pending';index"` // pending, completed, cancelled
	Notes         string    `json:"notes,omitempty" gorm:"type:text"`
	
	// 关联
	BankAccount *BankAccount `json:"bank_account,omitempty" gorm:"foreignKey:BankAccountID"`
}

// BankAccount 银行账户模型
type BankAccount struct {
	BaseModel
	AccountName   string  `json:"account_name" gorm:"not null"`
	BankName      string  `json:"bank_name" gorm:"not null"`
	AccountNumber string  `json:"account_number" gorm:"uniqueIndex;not null"`
	IBAN          string  `json:"iban,omitempty"`
	SwiftCode     string  `json:"swift_code,omitempty"`
	Currency      string  `json:"currency" gorm:"default:'USD'"`
	Balance       float64 `json:"balance" gorm:"default:0"`
	IsDefault     bool    `json:"is_default" gorm:"default:false"`
	IsActive      bool    `json:"is_active" gorm:"default:true"`
	AccountID     *uint   `json:"account_id,omitempty"`
	
	// 关联
	Account  *Account  `json:"account,omitempty" gorm:"foreignKey:AccountID"`
	Payments []Payment `json:"payments,omitempty" gorm:"foreignKey:BankAccountID"`
}

// Budget 预算模型
type Budget struct {
	AuditableModel
	BudgetName   string    `json:"budget_name" gorm:"size:255;not null"`
	BudgetYear   int       `json:"budget_year" gorm:"index;not null"`
	StartDate    time.Time `json:"start_date" gorm:"index;not null"`
	EndDate      time.Time `json:"end_date" gorm:"index;not null"`
	TotalAmount  float64   `json:"total_amount" gorm:"not null"`
	UsedAmount   float64   `json:"used_amount" gorm:"default:0"`
	RemainingAmount float64 `json:"remaining_amount" gorm:"default:0"`
	Status       string    `json:"status" gorm:"size:50;default:'draft';index"` // draft, approved, active, closed
	
	// 关联
	Items []BudgetItem `json:"items,omitempty" gorm:"foreignKey:BudgetID"`
}

// BudgetItem 预算明细模型
type BudgetItem struct {
	BaseModel
	BudgetID       uint    `json:"budget_id" gorm:"index;not null"`
	AccountID      uint    `json:"account_id" gorm:"index;not null"`
	BudgetAmount   float64 `json:"budget_amount" gorm:"not null"`
	ActualAmount   float64 `json:"actual_amount" gorm:"default:0"`
	VarianceAmount float64 `json:"variance_amount" gorm:"default:0"`
	Notes          string  `json:"notes,omitempty" gorm:"type:text"`
	
	// 关联
	Budget  Budget  `json:"budget,omitempty" gorm:"foreignKey:BudgetID"`
	Account Account `json:"account,omitempty" gorm:"foreignKey:AccountID"`
}

// TaxRate 税率模型
type TaxRate struct {
	CodeModel
	TaxType    string  `json:"tax_type" gorm:"size:50;not null;index"` // vat, income, sales
	Rate       float64 `json:"rate" gorm:"not null"`
	EffectiveDate time.Time `json:"effective_date" gorm:"index;not null"`
	ExpiryDate    *time.Time `json:"expiry_date,omitempty" gorm:"index"`
}

// CostCenter 成本中心模型
type CostCenter struct {
	CodeModel
	ManagerID    *uint `json:"manager_id,omitempty" gorm:"index"`
	DepartmentID *uint `json:"department_id,omitempty" gorm:"index"`
	
	// 关联
	Manager    *Employee   `json:"manager,omitempty" gorm:"foreignKey:ManagerID"`
	Department *Department `json:"department,omitempty" gorm:"foreignKey:DepartmentID"`
}

// FinancialReport 财务报表模型
type FinancialReport struct {
	AuditableModel
	ReportName   string    `json:"report_name" gorm:"size:255;not null"`
	ReportType   string    `json:"report_type" gorm:"size:50;not null;index"` // balance_sheet, income_statement, cash_flow
	PeriodType   string    `json:"period_type" gorm:"size:50;not null;index"` // monthly, quarterly, yearly
	StartDate    time.Time `json:"start_date" gorm:"index;not null"`
	EndDate      time.Time `json:"end_date" gorm:"index;not null"`
	Status       string    `json:"status" gorm:"size:50;default:'draft';index"` // draft, generated, approved
	FilePath     string    `json:"file_path,omitempty" gorm:"size:500"`
	
	// 关联
	Items []FinancialReportItem `json:"items,omitempty" gorm:"foreignKey:ReportID"`
}

// FinancialReportItem 财务报表明细模型
type FinancialReportItem struct {
	BaseModel
	ReportID    uint    `json:"report_id" gorm:"index;not null"`
	AccountID   uint    `json:"account_id" gorm:"index;not null"`
	Amount      float64 `json:"amount" gorm:"not null"`
	Percentage  float64 `json:"percentage,omitempty" gorm:"default:0"`
	
	// 关联
	Report  FinancialReport `json:"report,omitempty" gorm:"foreignKey:ReportID"`
	Account Account         `json:"account,omitempty" gorm:"foreignKey:AccountID"`
}

// Transaction 交易记录模型
type Transaction struct {
	AuditableModel
	TransactionNumber string    `json:"transaction_number" gorm:"uniqueIndex;size:100;not null"`
	TransactionDate   time.Time `json:"transaction_date" gorm:"index;not null"`
	TransactionType   string    `json:"transaction_type" gorm:"size:50;not null;index"` // income, expense, transfer
	Amount            float64   `json:"amount" gorm:"not null"`
	Currency          string    `json:"currency" gorm:"size:10;default:'CNY'"`
	ExchangeRate      float64   `json:"exchange_rate" gorm:"default:1"`
	Description       string    `json:"description" gorm:"type:text;not null"`
	ReferenceType     string    `json:"reference_type,omitempty" gorm:"size:50;index"` // invoice, payment, adjustment
	ReferenceID       *uint     `json:"reference_id,omitempty" gorm:"index"`
	Status            string    `json:"status" gorm:"size:50;default:'pending';index"` // pending, completed, cancelled
	Notes             string    `json:"notes,omitempty" gorm:"type:text"`
}

// Receivable 应收账款模型
type Receivable struct {
	BaseModel
	CustomerID    uint      `json:"customer_id" gorm:"not null"`
	InvoiceDate   time.Time `json:"invoice_date" gorm:"not null"`
	DueDate       time.Time `json:"due_date" gorm:"not null"`
	InvoiceNumber string    `json:"invoice_number" gorm:"not null"`
	Description   string    `json:"description,omitempty"`
	Amount        float64   `json:"amount" gorm:"not null"`
	AmountPaid    float64   `json:"amount_paid" gorm:"default:0"`
	Currency      string    `json:"currency" gorm:"default:'USD'"`
	ExchangeRate  float64   `json:"exchange_rate" gorm:"default:1"`
	Status        string    `json:"status" gorm:"default:'open'"`
	
	// 关联
	Customer *Customer `json:"customer,omitempty" gorm:"foreignKey:CustomerID"`
}

// Payable 应付账款模型
type Payable struct {
	BaseModel
	SupplierID    uint      `json:"supplier_id" gorm:"not null"`
	InvoiceDate   time.Time `json:"invoice_date" gorm:"not null"`
	DueDate       time.Time `json:"due_date" gorm:"not null"`
	InvoiceNumber string    `json:"invoice_number" gorm:"not null"`
	Description   string    `json:"description,omitempty"`
	Amount        float64   `json:"amount" gorm:"not null"`
	AmountPaid    float64   `json:"amount_paid" gorm:"default:0"`
	Currency      string    `json:"currency" gorm:"default:'USD'"`
	ExchangeRate  float64   `json:"exchange_rate" gorm:"default:1"`
	Status        string    `json:"status" gorm:"default:'open'"`
	
	// 关联
	Supplier *Supplier `json:"supplier,omitempty" gorm:"foreignKey:SupplierID"`
}

// FixedAsset 固定资产模型
type FixedAsset struct {
	AuditableModel
	AssetNumber      string    `json:"asset_number" gorm:"uniqueIndex;size:100;not null"`
	AssetName        string    `json:"asset_name" gorm:"size:255;not null"`
	AssetCategory    string    `json:"asset_category" gorm:"size:100;not null;index"`
	PurchaseDate     time.Time `json:"purchase_date" gorm:"index;not null"`
	PurchasePrice    float64   `json:"purchase_price" gorm:"not null"`
	CurrentValue     float64   `json:"current_value" gorm:"default:0"`
	DepreciationRate float64   `json:"depreciation_rate" gorm:"default:0"` // 年折旧率
	UsefulLife       int       `json:"useful_life" gorm:"default:0"`       // 使用年限
	Location         string    `json:"location" gorm:"size:255"`
	Status           string    `json:"status" gorm:"size:50;default:'active';index"` // active, disposed, sold
	Notes            string    `json:"notes,omitempty" gorm:"type:text"`
}

// DepreciationEntry 折旧记录模型
type DepreciationEntry struct {
	AuditableModel
	AssetID            uint      `json:"asset_id" gorm:"index;not null"`
	DepreciationDate   time.Time `json:"depreciation_date" gorm:"index;not null"`
	DepreciationAmount float64   `json:"depreciation_amount" gorm:"not null"`
	AccumulatedAmount  float64   `json:"accumulated_amount" gorm:"default:0"`
	BookValue          float64   `json:"book_value" gorm:"default:0"`
	Method             string    `json:"method" gorm:"size:50;not null"` // straight_line, declining_balance
	Notes              string    `json:"notes,omitempty" gorm:"type:text"`
	
	// 关联
	Asset FixedAsset `json:"asset,omitempty" gorm:"foreignKey:AssetID"`
}

// TaxEntry 税务记录模型
type TaxEntry struct {
	AuditableModel
	TaxNumber     string    `json:"tax_number" gorm:"uniqueIndex;size:100;not null"`
	TaxDate       time.Time `json:"tax_date" gorm:"index;not null"`
	TaxType       string    `json:"tax_type" gorm:"size:50;not null;index"` // vat, income, sales
	TaxableAmount float64   `json:"taxable_amount" gorm:"not null"`
	TaxRate       float64   `json:"tax_rate" gorm:"not null"`
	TaxAmount     float64   `json:"tax_amount" gorm:"not null"`
	Status        string    `json:"status" gorm:"size:50;default:'pending';index"` // pending, filed, paid
	Notes         string    `json:"notes,omitempty" gorm:"type:text"`
}

// Currency 货币模型
type Currency struct {
	BaseModel
	Code         string  `json:"code" gorm:"uniqueIndex;size:10;not null"`
	Name         string  `json:"name" gorm:"size:100;not null"`
	Symbol       string  `json:"symbol" gorm:"size:10;not null"`
	ExchangeRate float64 `json:"exchange_rate" gorm:"default:1"`
	IsBaseCurrency bool  `json:"is_base_currency" gorm:"default:false"`
	Status       string  `json:"status" gorm:"size:50;default:'active';index"`
}

// PaymentEntry 付款记录模型
type PaymentEntry struct {
	BaseModel
	PaymentType     string    `json:"payment_type" gorm:"not null"`
	PartyType       string    `json:"party_type" gorm:"not null"`
	PartyID         uint      `json:"party_id" gorm:"not null"`
	PostingDate     time.Time `json:"posting_date" gorm:"not null"`
	PaidAmount      float64   `json:"paid_amount" gorm:"not null"`
	ReceivedAmount  float64   `json:"received_amount" gorm:"not null"`
	Currency        string    `json:"currency" gorm:"default:'USD'"`
	ExchangeRate    float64   `json:"exchange_rate" gorm:"default:1"`
	BankAccountID   *uint     `json:"bank_account_id,omitempty"`
	CashAccountID   *uint     `json:"cash_account_id,omitempty"`
	Reference       string    `json:"reference,omitempty"`
	Remarks         string    `json:"remarks,omitempty"`
	Status          string    `json:"status" gorm:"default:'draft'"`
	IsPosted        bool      `json:"is_posted" gorm:"default:false"`
	PostedAt        *time.Time `json:"posted_at,omitempty"`
	CostCenterID    *uint     `json:"cost_center_id,omitempty"`
	ProjectID       *uint     `json:"project_id,omitempty"`
	
	// 关联
	BankAccount *BankAccount `json:"bank_account,omitempty" gorm:"foreignKey:BankAccountID"`
	CashAccount *Account     `json:"cash_account,omitempty" gorm:"foreignKey:CashAccountID"`
	CostCenter  *CostCenter  `json:"cost_center,omitempty" gorm:"foreignKey:CostCenterID"`
}

// ExchangeRateHistory 汇率历史模型
type ExchangeRateHistory struct {
	BaseModel
	CurrencyCode string    `json:"currency_code" gorm:"size:10;not null;index"`
	ExchangeRate float64   `json:"exchange_rate" gorm:"not null"`
	EffectiveDate time.Time `json:"effective_date" gorm:"index;not null"`
	Source       string    `json:"source" gorm:"size:100"`
}

// TaxTemplate 税务模板模型
type TaxTemplate struct {
	CodeModel
	TaxType     string  `json:"tax_type" gorm:"size:50;not null;index"`
	Rate        float64 `json:"rate" gorm:"not null"`
	IsDefault   bool    `json:"is_default" gorm:"default:false"`
	Calculation string  `json:"calculation" gorm:"size:50;not null"` // percentage, fixed
}

// FiscalYear 财政年度模型
type FiscalYear struct {
	BaseModel
	Year      int       `json:"year" gorm:"uniqueIndex;not null"`
	StartDate time.Time `json:"start_date" gorm:"index;not null"`
	EndDate   time.Time `json:"end_date" gorm:"index;not null"`
	IsCurrent bool      `json:"is_current" gorm:"default:false"`
	Status    string    `json:"status" gorm:"size:50;default:'active';index"`
}

// AccountingPeriod 会计期间模型
type AccountingPeriod struct {
	BaseModel
	Name       string     `json:"name" gorm:"not null"`
	StartDate  time.Time  `json:"start_date" gorm:"not null"`
	EndDate    time.Time  `json:"end_date" gorm:"not null"`
	FiscalYear string     `json:"fiscal_year" gorm:"not null"`
	IsClosed   bool       `json:"is_closed" gorm:"default:false"`
	ClosedBy   *uint      `json:"closed_by,omitempty"`
	ClosedAt   *time.Time `json:"closed_at,omitempty"`
	Company    string     `json:"company,omitempty"`
}