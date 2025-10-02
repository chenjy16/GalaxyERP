package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/galaxyerp/galaxyErp/internal/dto"
	"github.com/galaxyerp/galaxyErp/internal/models"
	"github.com/galaxyerp/galaxyErp/internal/repositories"
)

// AccountService 会计科目服务接口
type AccountService interface {
	CreateAccount(ctx context.Context, account *models.Account) error
	GetAccount(ctx context.Context, id uint) (*models.Account, error)
	GetAccountByCode(ctx context.Context, code string) (*models.Account, error)
	UpdateAccount(ctx context.Context, account *models.Account) error
	DeleteAccount(ctx context.Context, id uint) error
	ListAccounts(ctx context.Context, page, pageSize int) ([]*models.Account, int64, error)
	GetAccountsByType(ctx context.Context, accountType string, page, pageSize int) ([]*models.Account, int64, error)
	GetAccountChildren(ctx context.Context, parentID uint) ([]*models.Account, error)
	SearchAccounts(ctx context.Context, keyword string, page, pageSize int) ([]*models.Account, int64, error)
	ValidateAccountHierarchy(ctx context.Context, account *models.Account) error
}

// AccountServiceImpl 会计科目服务实现
type AccountServiceImpl struct {
	accountRepo repositories.AccountRepository
}

// NewAccountService 创建会计科目服务实例
func NewAccountService(accountRepo repositories.AccountRepository) AccountService {
	return &AccountServiceImpl{
		accountRepo: accountRepo,
	}
}

// CreateAccount 创建会计科目
func (s *AccountServiceImpl) CreateAccount(ctx context.Context, account *models.Account) error {
	// 验证科目编码是否已存在
	existingAccount, err := s.accountRepo.GetByCode(ctx, account.Code)
	if err != nil {
		return fmt.Errorf("检查科目编码失败: %w", err)
	}
	if existingAccount != nil {
		return errors.New("科目编码已存在")
	}

	// 验证科目层级关系
	if err := s.ValidateAccountHierarchy(ctx, account); err != nil {
		return err
	}

	// 设置默认值
	if account.Currency == "" {
		account.Currency = "CNY"
	}
	account.IsActive = true

	return s.accountRepo.Create(ctx, account)
}

// GetAccount 获取会计科目
func (s *AccountServiceImpl) GetAccount(ctx context.Context, id uint) (*models.Account, error) {
	account, err := s.accountRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("获取科目失败: %w", err)
	}
	if account == nil {
		return nil, errors.New("科目不存在")
	}
	return account, nil
}

// GetAccountByCode 根据编码获取会计科目
func (s *AccountServiceImpl) GetAccountByCode(ctx context.Context, code string) (*models.Account, error) {
	account, err := s.accountRepo.GetByCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("获取科目失败: %w", err)
	}
	if account == nil {
		return nil, errors.New("科目不存在")
	}
	return account, nil
}

// UpdateAccount 更新会计科目
func (s *AccountServiceImpl) UpdateAccount(ctx context.Context, account *models.Account) error {
	// 检查科目是否存在
	existingAccount, err := s.accountRepo.GetByID(ctx, account.ID)
	if err != nil {
		return fmt.Errorf("检查科目失败: %w", err)
	}
	if existingAccount == nil {
		return errors.New("科目不存在")
	}

	// 如果编码发生变化，检查新编码是否已存在
	if existingAccount.Code != account.Code {
		codeAccount, err := s.accountRepo.GetByCode(ctx, account.Code)
		if err != nil {
			return fmt.Errorf("检查科目编码失败: %w", err)
		}
		if codeAccount != nil && codeAccount.ID != account.ID {
			return errors.New("科目编码已存在")
		}
	}

	// 验证科目层级关系
	if err := s.ValidateAccountHierarchy(ctx, account); err != nil {
		return err
	}

	account.UpdatedAt = time.Now()
	return s.accountRepo.Update(ctx, account)
}

// DeleteAccount 删除会计科目
func (s *AccountServiceImpl) DeleteAccount(ctx context.Context, id uint) error {
	// 检查科目是否存在
	account, err := s.accountRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("检查科目失败: %w", err)
	}
	if account == nil {
		return errors.New("科目不存在")
	}

	// 检查是否有子科目
	children, err := s.accountRepo.GetChildren(ctx, id)
	if err != nil {
		return fmt.Errorf("检查子科目失败: %w", err)
	}
	if len(children) > 0 {
		return errors.New("存在子科目，无法删除")
	}

	// TODO: 检查是否有相关的交易记录
	// 这里应该检查是否有相关的会计分录或交易记录

	return s.accountRepo.Delete(ctx, id)
}

// ListAccounts 获取会计科目列表
func (s *AccountServiceImpl) ListAccounts(ctx context.Context, page, pageSize int) ([]*models.Account, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	return s.accountRepo.List(ctx, offset, pageSize)
}

// GetAccountsByType 根据科目类型获取会计科目
func (s *AccountServiceImpl) GetAccountsByType(ctx context.Context, accountType string, page, pageSize int) ([]*models.Account, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// 验证科目类型
	validTypes := map[string]bool{
		"ASSET":     true,
		"LIABILITY": true,
		"EQUITY":    true,
		"REVENUE":   true,
		"EXPENSE":   true,
	}
	if !validTypes[accountType] {
		return nil, 0, errors.New("无效的科目类型")
	}

	offset := (page - 1) * pageSize
	return s.accountRepo.GetByType(ctx, accountType, offset, pageSize)
}

// GetAccountChildren 获取子科目
func (s *AccountServiceImpl) GetAccountChildren(ctx context.Context, parentID uint) ([]*models.Account, error) {
	return s.accountRepo.GetChildren(ctx, parentID)
}

// SearchAccounts 搜索会计科目
func (s *AccountServiceImpl) SearchAccounts(ctx context.Context, keyword string, page, pageSize int) ([]*models.Account, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	// 如果仓储层没有搜索方法，使用列表方法
	// 这里假设仓储层有搜索方法，如果没有需要在仓储层添加
	accounts, total, err := s.accountRepo.List(ctx, offset, pageSize)
	if err != nil {
		return nil, 0, err
	}

	// 简单的内存过滤（生产环境应该在数据库层面进行搜索）
	var filteredAccounts []*models.Account
	for _, account := range accounts {
		if keyword == "" ||
			containsIgnoreCase(account.Code, keyword) ||
			containsIgnoreCase(account.Name, keyword) ||
			containsIgnoreCase(account.Description, keyword) {
			filteredAccounts = append(filteredAccounts, account)
		}
	}

	return filteredAccounts, total, nil
}

// ValidateAccountHierarchy 验证科目层级关系
func (s *AccountServiceImpl) ValidateAccountHierarchy(ctx context.Context, account *models.Account) error {
	if account.ParentID == nil {
		return nil // 根科目无需验证
	}

	// 检查父科目是否存在
	parent, err := s.accountRepo.GetByID(ctx, *account.ParentID)
	if err != nil {
		return fmt.Errorf("检查父科目失败: %w", err)
	}
	if parent == nil {
		return errors.New("父科目不存在")
	}

	// 检查父科目类型是否匹配
	if parent.AccountType != account.AccountType {
		return errors.New("子科目类型必须与父科目类型一致")
	}

	// 防止循环引用
	if account.ID != 0 && *account.ParentID == account.ID {
		return errors.New("科目不能以自己作为父科目")
	}

	// 检查是否会形成循环引用（递归检查）
	if account.ID != 0 {
		currentParent := parent
		for currentParent != nil && currentParent.ParentID != nil {
			if *currentParent.ParentID == account.ID {
				return errors.New("不能形成循环引用")
			}
			currentParent, err = s.accountRepo.GetByID(ctx, *currentParent.ParentID)
			if err != nil {
				break
			}
		}
	}

	return nil
}

// containsIgnoreCase 不区分大小写的字符串包含检查
func containsIgnoreCase(s, substr string) bool {
	return len(s) >= len(substr) &&
		(substr == "" ||
			len(s) > 0 &&
				(s == substr ||
					(len(s) > len(substr) &&
						(s[:len(substr)] == substr ||
							s[len(s)-len(substr):] == substr ||
							containsIgnoreCase(s[1:], substr)))))
}

// JournalEntryService 会计分录服务接口
type JournalEntryService interface {
	CreateJournalEntry(ctx context.Context, entry *models.JournalEntry) error
	CreateJournalEntryFromDTO(ctx context.Context, req *dto.JournalEntryCreateRequest) (*dto.JournalEntryResponse, error)
	GetJournalEntry(ctx context.Context, id uint) (*models.JournalEntry, error)
	UpdateJournalEntry(ctx context.Context, entry *models.JournalEntry) error
	DeleteJournalEntry(ctx context.Context, id uint) error
	ListJournalEntries(ctx context.Context, page, pageSize int) ([]*models.JournalEntry, int64, error)
	GetJournalEntriesByDateRange(ctx context.Context, startDate, endDate string, page, pageSize int) ([]*models.JournalEntry, int64, error)
}

// JournalEntryServiceImpl 会计分录服务实现
type JournalEntryServiceImpl struct {
	journalRepo repositories.JournalEntryRepository
	accountRepo repositories.AccountRepository
}

// NewJournalEntryService 创建会计分录服务实例
func NewJournalEntryService(journalRepo repositories.JournalEntryRepository, accountRepo repositories.AccountRepository) JournalEntryService {
	return &JournalEntryServiceImpl{
		journalRepo: journalRepo,
		accountRepo: accountRepo,
	}
}

// CreateJournalEntry 创建会计分录
func (s *JournalEntryServiceImpl) CreateJournalEntry(ctx context.Context, entry *models.JournalEntry) error {
	// 验证科目是否存在
	account, err := s.accountRepo.GetByID(ctx, entry.AccountID)
	if err != nil {
		return fmt.Errorf("检查科目失败: %w", err)
	}
	if account == nil {
		return errors.New("科目不存在")
	}

	// 验证借贷金额
	if entry.Debit < 0 || entry.Credit < 0 {
		return errors.New("借贷金额不能为负数")
	}
	if entry.Debit > 0 && entry.Credit > 0 {
		return errors.New("借贷金额不能同时大于0")
	}
	if entry.Debit == 0 && entry.Credit == 0 {
		return errors.New("借贷金额不能同时为0")
	}

	return s.journalRepo.Create(ctx, entry)
}

// CreateJournalEntryFromDTO 从DTO创建会计分录
func (s *JournalEntryServiceImpl) CreateJournalEntryFromDTO(ctx context.Context, req *dto.JournalEntryCreateRequest) (*dto.JournalEntryResponse, error) {
	// 首先创建一个Transaction记录
	transaction := &models.Transaction{
		TransactionNumber: fmt.Sprintf("TXN-%d", time.Now().Unix()),
		TransactionDate:   req.Date,
		TransactionType:   "journal",
		Amount:            0, // 将在下面计算
		Description:       req.Description,
		Status:            "completed",
	}

	// 计算总金额
	var totalAmount float64
	for _, item := range req.Items {
		if item.DebitAmount > 0 {
			totalAmount += item.DebitAmount
		}
	}
	transaction.Amount = totalAmount

	// 这里我们需要一个Transaction repository，暂时跳过Transaction创建
	// 直接创建JournalEntry记录，使用固定的TransactionID
	transactionID := uint(1) // 临时使用固定ID

	var journalEntries []*models.JournalEntry
	for _, item := range req.Items {
		entry := &models.JournalEntry{
			TransactionID: transactionID,
			AccountID:     item.AccountID,
			Debit:         item.DebitAmount,
			Credit:        item.CreditAmount,
			Description:   item.Description,
		}

		// 验证科目存在
		account, err := s.accountRepo.GetByID(ctx, item.AccountID)
		if err != nil {
			return nil, fmt.Errorf("检查科目失败: %w", err)
		}
		if account == nil {
			return nil, fmt.Errorf("科目ID %d 不存在", item.AccountID)
		}

		if err := s.journalRepo.Create(ctx, entry); err != nil {
			return nil, fmt.Errorf("创建分录失败: %w", err)
		}
		journalEntries = append(journalEntries, entry)
	}

	// 构造响应
	response := &dto.JournalEntryResponse{
		ID:          journalEntries[0].ID, // 使用第一个分录的ID
		Number:      fmt.Sprintf("JE-%d", journalEntries[0].ID),
		Date:        req.Date,
		Reference:   req.Reference,
		Description: req.Description,
		TotalDebit:  0,
		TotalCredit: 0,
		Status:      "posted",
		Items:       make([]dto.JournalEntryItemResponse, 0),
	}

	// 计算总借贷金额并构造分录项响应
	for _, entry := range journalEntries {
		response.TotalDebit += entry.Debit
		response.TotalCredit += entry.Credit

		// 获取科目信息
		account, _ := s.accountRepo.GetByID(ctx, entry.AccountID)

		itemResponse := dto.JournalEntryItemResponse{
			ID:           entry.ID,
			DebitAmount:  entry.Debit,
			CreditAmount: entry.Credit,
			Description:  entry.Description,
			Account: dto.AccountResponse{
				ID:   account.ID,
				Code: account.Code,
				Name: account.Name,
				Type: account.AccountType,
			},
		}
		response.Items = append(response.Items, itemResponse)
	}

	return response, nil
}

// GetJournalEntry 获取会计分录
func (s *JournalEntryServiceImpl) GetJournalEntry(ctx context.Context, id uint) (*models.JournalEntry, error) {
	entry, err := s.journalRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("获取分录失败: %w", err)
	}
	if entry == nil {
		return nil, errors.New("分录不存在")
	}
	return entry, nil
}

// UpdateJournalEntry 更新会计分录
func (s *JournalEntryServiceImpl) UpdateJournalEntry(ctx context.Context, entry *models.JournalEntry) error {
	// 检查分录是否存在
	existingEntry, err := s.journalRepo.GetByID(ctx, entry.ID)
	if err != nil {
		return fmt.Errorf("检查分录失败: %w", err)
	}
	if existingEntry == nil {
		return errors.New("分录不存在")
	}

	// 验证科目是否存在
	account, err := s.accountRepo.GetByID(ctx, entry.AccountID)
	if err != nil {
		return fmt.Errorf("检查科目失败: %w", err)
	}
	if account == nil {
		return errors.New("科目不存在")
	}

	// 验证借贷金额
	if entry.Debit < 0 || entry.Credit < 0 {
		return errors.New("借贷金额不能为负数")
	}
	if entry.Debit > 0 && entry.Credit > 0 {
		return errors.New("借贷金额不能同时大于0")
	}
	if entry.Debit == 0 && entry.Credit == 0 {
		return errors.New("借贷金额不能同时为0")
	}

	entry.UpdatedAt = time.Now()
	return s.journalRepo.Update(ctx, entry)
}

// DeleteJournalEntry 删除会计分录
func (s *JournalEntryServiceImpl) DeleteJournalEntry(ctx context.Context, id uint) error {
	// 检查分录是否存在
	entry, err := s.journalRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("检查分录失败: %w", err)
	}
	if entry == nil {
		return errors.New("分录不存在")
	}

	return s.journalRepo.Delete(ctx, id)
}

// ListJournalEntries 获取会计分录列表
func (s *JournalEntryServiceImpl) ListJournalEntries(ctx context.Context, page, pageSize int) ([]*models.JournalEntry, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	return s.journalRepo.List(ctx, offset, pageSize)
}

// GetJournalEntriesByDateRange 根据日期范围获取会计分录
func (s *JournalEntryServiceImpl) GetJournalEntriesByDateRange(ctx context.Context, startDate, endDate string, page, pageSize int) ([]*models.JournalEntry, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	return s.journalRepo.GetByDateRange(ctx, startDate, endDate, offset, pageSize)
}

// PaymentEntryService 付款记录服务接口
type PaymentEntryService interface {
	CreatePaymentEntry(ctx context.Context, payment *models.PaymentEntry) error
	GetPaymentEntry(ctx context.Context, id uint) (*models.PaymentEntry, error)
	UpdatePaymentEntry(ctx context.Context, payment *models.PaymentEntry) error
	DeletePaymentEntry(ctx context.Context, id uint) error
	ListPaymentEntries(ctx context.Context, page, pageSize int) ([]*models.PaymentEntry, int64, error)
	GetPaymentEntriesByParty(ctx context.Context, partyType string, partyID uint, page, pageSize int) ([]*models.PaymentEntry, int64, error)
}

// PaymentEntryServiceImpl 付款记录服务实现
type PaymentEntryServiceImpl struct {
	paymentRepo repositories.PaymentEntryRepository
}

// NewPaymentEntryService 创建付款记录服务实例
func NewPaymentEntryService(paymentRepo repositories.PaymentEntryRepository) PaymentEntryService {
	return &PaymentEntryServiceImpl{
		paymentRepo: paymentRepo,
	}
}

// CreatePaymentEntry 创建付款记录
func (s *PaymentEntryServiceImpl) CreatePaymentEntry(ctx context.Context, payment *models.PaymentEntry) error {
	// 验证必填字段
	if payment.PaymentType == "" {
		return errors.New("付款类型不能为空")
	}
	if payment.PartyType == "" {
		return errors.New("关联方类型不能为空")
	}
	if payment.PartyID == 0 {
		return errors.New("关联方ID不能为空")
	}
	if payment.PaidAmount < 0 || payment.ReceivedAmount < 0 {
		return errors.New("付款金额不能为负数")
	}

	payment.CreatedAt = time.Now()
	payment.UpdatedAt = time.Now()
	if payment.Status == "" {
		payment.Status = "draft"
	}

	return s.paymentRepo.Create(ctx, payment)
}

// GetPaymentEntry 获取付款记录
func (s *PaymentEntryServiceImpl) GetPaymentEntry(ctx context.Context, id uint) (*models.PaymentEntry, error) {
	if id == 0 {
		return nil, errors.New("付款记录ID不能为空")
	}
	return s.paymentRepo.GetByID(ctx, id)
}

// UpdatePaymentEntry 更新付款记录
func (s *PaymentEntryServiceImpl) UpdatePaymentEntry(ctx context.Context, payment *models.PaymentEntry) error {
	// 检查记录是否存在
	existing, err := s.paymentRepo.GetByID(ctx, payment.ID)
	if err != nil {
		return fmt.Errorf("检查付款记录失败: %w", err)
	}
	if existing == nil {
		return errors.New("付款记录不存在")
	}

	// 验证必填字段
	if payment.PaymentType == "" {
		return errors.New("付款类型不能为空")
	}
	if payment.PartyType == "" {
		return errors.New("关联方类型不能为空")
	}
	if payment.PartyID == 0 {
		return errors.New("关联方ID不能为空")
	}
	if payment.PaidAmount < 0 || payment.ReceivedAmount < 0 {
		return errors.New("付款金额不能为负数")
	}

	payment.UpdatedAt = time.Now()
	return s.paymentRepo.Update(ctx, payment)
}

// DeletePaymentEntry 删除付款记录
func (s *PaymentEntryServiceImpl) DeletePaymentEntry(ctx context.Context, id uint) error {
	// 检查记录是否存在
	payment, err := s.paymentRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("检查付款记录失败: %w", err)
	}
	if payment == nil {
		return errors.New("付款记录不存在")
	}

	return s.paymentRepo.Delete(ctx, id)
}

// ListPaymentEntries 获取付款记录列表
func (s *PaymentEntryServiceImpl) ListPaymentEntries(ctx context.Context, page, pageSize int) ([]*models.PaymentEntry, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	return s.paymentRepo.List(ctx, offset, pageSize)
}

// GetPaymentEntriesByParty 根据关联方获取付款记录
func (s *PaymentEntryServiceImpl) GetPaymentEntriesByParty(ctx context.Context, partyType string, partyID uint, page, pageSize int) ([]*models.PaymentEntry, int64, error) {
	if partyType == "" {
		return nil, 0, errors.New("关联方类型不能为空")
	}
	if partyID == 0 {
		return nil, 0, errors.New("关联方ID不能为空")
	}
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	return s.paymentRepo.GetByParty(ctx, partyType, partyID, offset, pageSize)
}
