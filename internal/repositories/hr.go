package repositories

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"github.com/galaxyerp/galaxyErp/internal/models"
)

// EmployeeRepository 员工仓储接口
type EmployeeRepository interface {
	Create(ctx context.Context, employee *models.Employee) error
	GetByID(ctx context.Context, id uint) (*models.Employee, error)
	GetByCode(ctx context.Context, code string) (*models.Employee, error)
	Update(ctx context.Context, employee *models.Employee) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, offset, limit int) ([]*models.Employee, int64, error)
	Search(ctx context.Context, query string, offset, limit int) ([]*models.Employee, int64, error)
	GetByDepartmentID(ctx context.Context, departmentID uint, offset, limit int) ([]*models.Employee, int64, error)
}

// EmployeeRepositoryImpl 员工仓储实现
type EmployeeRepositoryImpl struct {
	db *gorm.DB
}

// NewEmployeeRepository 创建员工仓储实例
func NewEmployeeRepository(db *gorm.DB) EmployeeRepository {
	return &EmployeeRepositoryImpl{
		db: db,
	}
}

// Create 创建员工
func (r *EmployeeRepositoryImpl) Create(ctx context.Context, employee *models.Employee) error {
	return r.db.WithContext(ctx).Create(employee).Error
}

// GetByID 根据ID获取员工
func (r *EmployeeRepositoryImpl) GetByID(ctx context.Context, id uint) (*models.Employee, error) {
	var employee models.Employee
	err := r.db.WithContext(ctx).Preload("Department").Preload("Position").Preload("Manager").First(&employee, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &employee, nil
}

// GetByCode 根据员工编号获取员工
func (r *EmployeeRepositoryImpl) GetByCode(ctx context.Context, code string) (*models.Employee, error) {
	var employee models.Employee
	err := r.db.WithContext(ctx).Preload("Department").Preload("Position").Preload("Manager").Where("code = ?", code).First(&employee).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &employee, nil
}

// Update 更新员工
func (r *EmployeeRepositoryImpl) Update(ctx context.Context, employee *models.Employee) error {
	return r.db.WithContext(ctx).Save(employee).Error
}

// Delete 删除员工
func (r *EmployeeRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Employee{}, id).Error
}

// List 获取员工列表
func (r *EmployeeRepositoryImpl) List(ctx context.Context, offset, limit int) ([]*models.Employee, int64, error) {
	var employees []*models.Employee
	var total int64
	
	// 获取总数
	if err := r.db.WithContext(ctx).Model(&models.Employee{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// 获取分页数据
	err := r.db.WithContext(ctx).Preload("Department").Preload("Position").Offset(offset).Limit(limit).Find(&employees).Error
	if err != nil {
		return nil, 0, err
	}
	
	return employees, total, nil
}

// Search 搜索员工
func (r *EmployeeRepositoryImpl) Search(ctx context.Context, query string, offset, limit int) ([]*models.Employee, int64, error) {
	var employees []*models.Employee
	var total int64
	
	searchQuery := "%" + query + "%"
	
	// 获取总数
	if err := r.db.WithContext(ctx).Model(&models.Employee{}).
		Where("code LIKE ? OR first_name LIKE ? OR last_name LIKE ? OR full_name LIKE ? OR email LIKE ?", 
			searchQuery, searchQuery, searchQuery, searchQuery, searchQuery).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// 获取分页数据
	err := r.db.WithContext(ctx).Preload("Department").Preload("Position").
		Where("code LIKE ? OR first_name LIKE ? OR last_name LIKE ? OR full_name LIKE ? OR email LIKE ?", 
			searchQuery, searchQuery, searchQuery, searchQuery, searchQuery).
		Offset(offset).Limit(limit).Find(&employees).Error
	if err != nil {
		return nil, 0, err
	}
	
	return employees, total, nil
}

// GetByDepartmentID 根据部门ID获取员工
func (r *EmployeeRepositoryImpl) GetByDepartmentID(ctx context.Context, departmentID uint, offset, limit int) ([]*models.Employee, int64, error) {
	var employees []*models.Employee
	var total int64
	
	// 获取总数
	if err := r.db.WithContext(ctx).Model(&models.Employee{}).
		Where("department_id = ?", departmentID).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// 获取分页数据
	err := r.db.WithContext(ctx).Preload("Department").Preload("Position").
		Where("department_id = ?", departmentID).
		Offset(offset).Limit(limit).Find(&employees).Error
	if err != nil {
		return nil, 0, err
	}
	
	return employees, total, nil
}