package dto

import (
	"time"
)

// AccountCreateRequest 账户创建请求
type AccountCreateRequest struct {
	Code     string  `json:"code" binding:"required,max=50"`
	Name     string  `json:"name" binding:"required,max=100"`
	Type     string  `json:"type" binding:"required,oneof=asset liability equity revenue expense"`
	ParentID *uint   `json:"parent_id,omitempty"`
	Balance  float64 `json:"balance,omitempty"`
	Status   string  `json:"status" binding:"required,oneof=active inactive"`
}

// AccountUpdateRequest 账户更新请求
type AccountUpdateRequest struct {
	Name     string   `json:"name,omitempty" binding:"omitempty,max=100"`
	Type     string   `json:"type,omitempty" binding:"omitempty,oneof=asset liability equity revenue expense"`
	ParentID *uint    `json:"parent_id,omitempty"`
	Balance  *float64 `json:"balance,omitempty"`
	Status   string   `json:"status,omitempty" binding:"omitempty,oneof=active inactive"`
}

// AccountResponse 账户响应
type AccountResponse struct {
	ID        uint              `json:"id"`
	Code      string            `json:"code"`
	Name      string            `json:"name"`
	Type      string            `json:"type"`
	ParentID  *uint             `json:"parent_id,omitempty"`
	Balance   float64           `json:"balance"`
	Status    string            `json:"status"`
	Parent    *AccountResponse  `json:"parent,omitempty"`
	Children  []AccountResponse `json:"children,omitempty"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

// AccountListResponse 账户列表响应
type AccountListResponse struct {
	ID       uint    `json:"id"`
	Code     string  `json:"code"`
	Name     string  `json:"name"`
	Type     string  `json:"type"`
	Balance  float64 `json:"balance"`
	Status   string  `json:"status"`
	ParentID *uint   `json:"parent_id,omitempty"`
}

// JournalEntryCreateRequest 日记账分录创建请求
type JournalEntryCreateRequest struct {
	Date        time.Time                    `json:"date" binding:"required"`
	Reference   string                       `json:"reference,omitempty"`
	Description string                       `json:"description" binding:"required"`
	Items       []JournalEntryItemRequest    `json:"items" binding:"required,min=2"`
}

// JournalEntryItemRequest 日记账分录项请求
type JournalEntryItemRequest struct {
	AccountID   uint    `json:"account_id" binding:"required"`
	DebitAmount float64 `json:"debit_amount,omitempty" binding:"min=0"`
	CreditAmount float64 `json:"credit_amount,omitempty" binding:"min=0"`
	Description string  `json:"description,omitempty"`
}

// JournalEntryResponse 日记账分录响应
type JournalEntryResponse struct {
	ID          uint                        `json:"id"`
	Number      string                      `json:"number"`
	Date        time.Time                   `json:"date"`
	Reference   string                      `json:"reference,omitempty"`
	Description string                      `json:"description"`
	TotalDebit  float64                     `json:"total_debit"`
	TotalCredit float64                     `json:"total_credit"`
	Status      string                      `json:"status"`
	Items       []JournalEntryItemResponse  `json:"items"`
	CreatedBy   UserResponse                `json:"created_by"`
	CreatedAt   time.Time                   `json:"created_at"`
	UpdatedAt   time.Time                   `json:"updated_at"`
}

// JournalEntryItemResponse 日记账分录项响应
type JournalEntryItemResponse struct {
	ID           uint            `json:"id"`
	DebitAmount  float64         `json:"debit_amount"`
	CreditAmount float64         `json:"credit_amount"`
	Description  string          `json:"description,omitempty"`
	Account      AccountResponse `json:"account"`
}

// PaymentCreateRequest 付款创建请求
type PaymentCreateRequest struct {
	Type        string    `json:"type" binding:"required,oneof=payment receipt"`
	Amount      float64   `json:"amount" binding:"required,gt=0"`
	Date        time.Time `json:"date" binding:"required"`
	Reference   string    `json:"reference,omitempty"`
	Description string    `json:"description,omitempty"`
	AccountID   uint      `json:"account_id" binding:"required"`
	PaymentMethod string  `json:"payment_method" binding:"required,oneof=cash bank_transfer check credit_card"`
	CustomerID  *uint     `json:"customer_id,omitempty"`
	SupplierID  *uint     `json:"supplier_id,omitempty"`
}

// PaymentResponse 付款响应
type PaymentResponse struct {
	ID            uint              `json:"id"`
	Number        string            `json:"number"`
	Type          string            `json:"type"`
	Amount        float64           `json:"amount"`
	Date          time.Time         `json:"date"`
	Reference     string            `json:"reference,omitempty"`
	Description   string            `json:"description,omitempty"`
	PaymentMethod string            `json:"payment_method"`
	Status        string            `json:"status"`
	Account       AccountResponse   `json:"account"`
	Customer      *CustomerResponse `json:"customer,omitempty"`
	Supplier      *SupplierResponse `json:"supplier,omitempty"`
	CreatedBy     UserResponse      `json:"created_by"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
}

// AccountSearchRequest 账户搜索请求
type AccountSearchRequest struct {
	SearchRequest
	Type     string `json:"type,omitempty" form:"type"`
	Status   string `json:"status,omitempty" form:"status"`
	ParentID *uint  `json:"parent_id,omitempty" form:"parent_id"`
}

// PaymentSearchRequest 付款搜索请求
type PaymentSearchRequest struct {
	SearchRequest
	Type       string     `json:"type,omitempty" form:"type"`
	AccountID  *uint      `json:"account_id,omitempty" form:"account_id"`
	CustomerID *uint      `json:"customer_id,omitempty" form:"customer_id"`
	SupplierID *uint      `json:"supplier_id,omitempty" form:"supplier_id"`
	StartDate  *time.Time `json:"start_date,omitempty" form:"start_date"`
	EndDate    *time.Time `json:"end_date,omitempty" form:"end_date"`
	MinAmount  *float64   `json:"min_amount,omitempty" form:"min_amount"`
	MaxAmount  *float64   `json:"max_amount,omitempty" form:"max_amount"`
}

// FinancialReportRequest 财务报表请求
type FinancialReportRequest struct {
	Type      string    `json:"type" binding:"required,oneof=balance_sheet income_statement cash_flow"`
	StartDate time.Time `json:"start_date" binding:"required"`
	EndDate   time.Time `json:"end_date" binding:"required"`
	Format    string    `json:"format" form:"format" binding:"required,oneof=excel pdf csv"`
}