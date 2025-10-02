package repositories

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"github.com/galaxyerp/galaxyErp/internal/models"
)

// AccountRepository 会计科目仓储接口
type AccountRepository interface {
	Create(ctx context.Context, account *models.Account) error
	GetByID(ctx context.Context, id uint) (*models.Account, error)
	GetByCode(ctx context.Context, code string) (*models.Account, error)
	Update(ctx context.Context, account *models.Account) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, offset, limit int) ([]*models.Account, int64, error)
	GetByType(ctx context.Context, accountType string, offset, limit int) ([]*models.Account, int64, error)
	GetChildren(ctx context.Context, parentID uint) ([]*models.Account, error)
}

// AccountRepositoryImpl 会计科目仓储实现
type AccountRepositoryImpl struct {
	db *gorm.DB
}

// NewAccountRepository 创建会计科目仓储实例
func NewAccountRepository(db *gorm.DB) AccountRepository {
	return &AccountRepositoryImpl{
		db: db,
	}
}

// Create 创建会计科目
func (r *AccountRepositoryImpl) Create(ctx context.Context, account *models.Account) error {
	return r.db.WithContext(ctx).Create(account).Error
}

// GetByID 根据ID获取会计科目
func (r *AccountRepositoryImpl) GetByID(ctx context.Context, id uint) (*models.Account, error) {
	var account models.Account
	err := r.db.WithContext(ctx).Preload("Parent").Preload("Children").First(&account, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &account, nil
}

// GetByCode 根据科目编码获取会计科目
func (r *AccountRepositoryImpl) GetByCode(ctx context.Context, code string) (*models.Account, error) {
	var account models.Account
	err := r.db.WithContext(ctx).Preload("Parent").Preload("Children").Where("code = ?", code).First(&account).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &account, nil
}

// Update 更新会计科目
func (r *AccountRepositoryImpl) Update(ctx context.Context, account *models.Account) error {
	return r.db.WithContext(ctx).Save(account).Error
}

// Delete 删除会计科目
func (r *AccountRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Account{}, id).Error
}

// List 获取会计科目列表
func (r *AccountRepositoryImpl) List(ctx context.Context, offset, limit int) ([]*models.Account, int64, error) {
	var accounts []*models.Account
	var total int64
	
	// 获取总数
	if err := r.db.WithContext(ctx).Model(&models.Account{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// 获取分页数据
	err := r.db.WithContext(ctx).Preload("Parent").Offset(offset).Limit(limit).Find(&accounts).Error
	if err != nil {
		return nil, 0, err
	}
	
	return accounts, total, nil
}

// GetByType 根据科目类型获取会计科目
func (r *AccountRepositoryImpl) GetByType(ctx context.Context, accountType string, offset, limit int) ([]*models.Account, int64, error) {
	var accounts []*models.Account
	var total int64
	
	// 获取总数
	if err := r.db.WithContext(ctx).Model(&models.Account{}).
		Where("type = ?", accountType).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// 获取分页数据
	err := r.db.WithContext(ctx).Preload("Parent").
		Where("type = ?", accountType).
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
	if err != nil {
		return nil, err
	}
	return accounts, nil
}

// JournalEntryRepository 记账凭证仓储接口
type JournalEntryRepository interface {
	Create(ctx context.Context, entry *models.JournalEntry) error
	GetByID(ctx context.Context, id uint) (*models.JournalEntry, error)
	Update(ctx context.Context, entry *models.JournalEntry) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, offset, limit int) ([]*models.JournalEntry, int64, error)
	GetByDateRange(ctx context.Context, startDate, endDate string, offset, limit int) ([]*models.JournalEntry, int64, error)
}

// JournalEntryRepositoryImpl 记账凭证仓储实现
type JournalEntryRepositoryImpl struct {
	db *gorm.DB
}

// NewJournalEntryRepository 创建记账凭证仓储实例
func NewJournalEntryRepository(db *gorm.DB) JournalEntryRepository {
	return &JournalEntryRepositoryImpl{
		db: db,
	}
}

// Create 创建记账凭证
func (r *JournalEntryRepositoryImpl) Create(ctx context.Context, entry *models.JournalEntry) error {
	return r.db.WithContext(ctx).Create(entry).Error
}

// GetByID 根据ID获取记账凭证
func (r *JournalEntryRepositoryImpl) GetByID(ctx context.Context, id uint) (*models.JournalEntry, error) {
	var entry models.JournalEntry
	err := r.db.WithContext(ctx).Preload("Items").Preload("Items.Account").First(&entry, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &entry, nil
}

// Update 更新记账凭证
func (r *JournalEntryRepositoryImpl) Update(ctx context.Context, entry *models.JournalEntry) error {
	return r.db.WithContext(ctx).Save(entry).Error
}

// Delete 删除记账凭证
func (r *JournalEntryRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.JournalEntry{}, id).Error
}

// List 获取记账凭证列表
func (r *JournalEntryRepositoryImpl) List(ctx context.Context, offset, limit int) ([]*models.JournalEntry, int64, error) {
	var entries []*models.JournalEntry
	var total int64
	
	// 获取总数
	if err := r.db.WithContext(ctx).Model(&models.JournalEntry{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// 获取分页数据，预加载Account关联
	err := r.db.WithContext(ctx).Preload("Account").Offset(offset).Limit(limit).Find(&entries).Error
	if err != nil {
		return nil, 0, err
	}
	
	return entries, total, nil
}

// GetByDateRange 根据日期范围获取记账凭证
func (r *JournalEntryRepositoryImpl) GetByDateRange(ctx context.Context, startDate, endDate string, offset, limit int) ([]*models.JournalEntry, int64, error) {
	var entries []*models.JournalEntry
	var total int64
	
	// 获取总数
	if err := r.db.WithContext(ctx).Model(&models.JournalEntry{}).
		Where("entry_date BETWEEN ? AND ?", startDate, endDate).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// 获取分页数据
	err := r.db.WithContext(ctx).Preload("Items").
		Where("entry_date BETWEEN ? AND ?", startDate, endDate).
		Offset(offset).Limit(limit).Find(&entries).Error
	if err != nil {
		return nil, 0, err
	}
	
	return entries, total, nil
}