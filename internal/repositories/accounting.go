package repositories

import (
	"context"
	"github.com/galaxyerp/galaxyErp/internal/models"
	"gorm.io/gorm"
)

// AccountRepository 会计科目仓储接口
type AccountRepository interface {
	BaseRepository[models.Account]
	GetByCode(ctx context.Context, code string) (*models.Account, error)
	GetByType(ctx context.Context, accountType string, offset, limit int) ([]*models.Account, int64, error)
	GetChildren(ctx context.Context, parentID uint) ([]*models.Account, error)
}

// AccountRepositoryImpl 会计科目仓储实现
type AccountRepositoryImpl struct {
	BaseRepository[models.Account]
	db *gorm.DB
}

// NewAccountRepository 创建会计科目仓储实例
func NewAccountRepository(db *gorm.DB) AccountRepository {
	return &AccountRepositoryImpl{
		BaseRepository: NewBaseRepository[models.Account](db),
		db:             db,
	}
}

// GetByCode 根据科目编码获取会计科目
func (r *AccountRepositoryImpl) GetByCode(ctx context.Context, code string) (*models.Account, error) {
	var account models.Account
	err := r.db.WithContext(ctx).Preload("Parent").Preload("Children").Where("code = ?", code).First(&account).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

// GetByType 根据科目类型获取会计科目
func (r *AccountRepositoryImpl) GetByType(ctx context.Context, accountType string, offset, limit int) ([]*models.Account, int64, error) {
	var accounts []*models.Account
	var total int64

	// 获取总数
	if err := r.db.WithContext(ctx).Model(&models.Account{}).
		Where("account_type = ?", accountType).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := r.db.WithContext(ctx).Preload("Parent").
		Where("account_type = ?", accountType).
		Offset(offset).Limit(limit).Find(&accounts).Error
	if err != nil {
		return nil, 0, err
	}

	return accounts, total, nil
}

// GetChildren 获取子科目
func (r *AccountRepositoryImpl) GetChildren(ctx context.Context, parentID uint) ([]*models.Account, error) {
	var accounts []*models.Account
	err := r.db.WithContext(ctx).Where("parent_id = ?", parentID).Find(&accounts).Error
	return accounts, err
}

// JournalEntryRepository 日记账分录仓储接口
type JournalEntryRepository interface {
	BaseRepository[models.JournalEntry]
	GetByDateRange(ctx context.Context, startDate, endDate string, offset, limit int) ([]*models.JournalEntry, int64, error)
}

// JournalEntryRepositoryImpl 日记账分录仓储实现
type JournalEntryRepositoryImpl struct {
	BaseRepository[models.JournalEntry]
	db *gorm.DB
}

// NewJournalEntryRepository 创建日记账分录仓储实例
func NewJournalEntryRepository(db *gorm.DB) JournalEntryRepository {
	return &JournalEntryRepositoryImpl{
		BaseRepository: NewBaseRepository[models.JournalEntry](db),
		db:             db,
	}
}

// GetByDateRange 根据日期范围获取日记账分录
func (r *JournalEntryRepositoryImpl) GetByDateRange(ctx context.Context, startDate, endDate string, offset, limit int) ([]*models.JournalEntry, int64, error) {
	var entries []*models.JournalEntry
	var total int64

	// 获取总数
	if err := r.db.WithContext(ctx).Model(&models.JournalEntry{}).
		Where("posting_date >= ? AND posting_date <= ?", startDate, endDate).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := r.db.WithContext(ctx).Preload("JournalEntryAccounts").
		Where("posting_date >= ? AND posting_date <= ?", startDate, endDate).
		Offset(offset).Limit(limit).Find(&entries).Error
	if err != nil {
		return nil, 0, err
	}

	return entries, total, nil
}

// PaymentEntryRepository 付款分录仓储接口
type PaymentEntryRepository interface {
	BaseRepository[models.PaymentEntry]
	GetByParty(ctx context.Context, partyType string, partyID uint, offset, limit int) ([]*models.PaymentEntry, int64, error)
}

// PaymentEntryRepositoryImpl 付款分录仓储实现
type PaymentEntryRepositoryImpl struct {
	BaseRepository[models.PaymentEntry]
	db *gorm.DB
}

// NewPaymentEntryRepository 创建付款分录仓储实例
func NewPaymentEntryRepository(db *gorm.DB) PaymentEntryRepository {
	return &PaymentEntryRepositoryImpl{
		BaseRepository: NewBaseRepository[models.PaymentEntry](db),
		db:             db,
	}
}

// GetByParty 根据当事方获取付款分录
func (r *PaymentEntryRepositoryImpl) GetByParty(ctx context.Context, partyType string, partyID uint, offset, limit int) ([]*models.PaymentEntry, int64, error) {
	var payments []*models.PaymentEntry
	var total int64

	// 获取总数
	if err := r.db.WithContext(ctx).Model(&models.PaymentEntry{}).
		Where("party_type = ? AND party_id = ?", partyType, partyID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := r.db.WithContext(ctx).Preload("PaymentEntryReferences").
		Where("party_type = ? AND party_id = ?", partyType, partyID).
		Offset(offset).Limit(limit).Find(&payments).Error
	if err != nil {
		return nil, 0, err
	}

	return payments, total, nil
}
