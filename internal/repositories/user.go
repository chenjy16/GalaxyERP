package repositories

import (
	"context"
	"errors"
	"github.com/galaxyerp/galaxyErp/internal/models"
	"gorm.io/gorm"
)

// UserRepository 用户仓储接口
type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id uint) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, offset, limit int) ([]*models.User, int64, error)
	Search(ctx context.Context, query string, offset, limit int) ([]*models.User, int64, error)
}

// UserRepositoryImpl 用户仓储实现
type UserRepositoryImpl struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓储实例
func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{
		db: db,
	}
}

// Create 创建用户
func (r *UserRepositoryImpl) Create(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// GetByID 根据ID获取用户
func (r *UserRepositoryImpl) GetByID(ctx context.Context, id uint) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// GetByEmail 根据邮箱获取用户
func (r *UserRepositoryImpl) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Preload("Roles").Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// GetByUsername 根据用户名获取用户
func (r *UserRepositoryImpl) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Preload("Roles").Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// Update 更新用户
func (r *UserRepositoryImpl) Update(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

// Delete 删除用户
func (r *UserRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.User{}, id).Error
}

// List 获取用户列表
func (r *UserRepositoryImpl) List(ctx context.Context, offset, limit int) ([]*models.User, int64, error) {
	var users []*models.User
	var total int64

	// 获取总数
	if err := r.db.WithContext(ctx).Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := r.db.WithContext(ctx).Preload("Role").Offset(offset).Limit(limit).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// Search 搜索用户
func (r *UserRepositoryImpl) Search(ctx context.Context, query string, offset, limit int) ([]*models.User, int64, error) {
	var users []*models.User
	var total int64

	searchQuery := "%" + query + "%"

	// 获取总数
	if err := r.db.WithContext(ctx).Model(&models.User{}).
		Where("username LIKE ? OR email LIKE ? OR first_name LIKE ? OR last_name LIKE ?",
			searchQuery, searchQuery, searchQuery, searchQuery).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := r.db.WithContext(ctx).Preload("Role").
		Where("username LIKE ? OR email LIKE ? OR first_name LIKE ? OR last_name LIKE ?",
			searchQuery, searchQuery, searchQuery, searchQuery).
		Offset(offset).Limit(limit).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}
