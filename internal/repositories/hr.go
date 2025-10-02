package repositories

import (
	"context"
	"errors"
	"github.com/galaxyerp/galaxyErp/internal/models"
	"gorm.io/gorm"
	"time"
)

// EmployeeRepository 员工仓储接口
type EmployeeRepository interface {
	Create(ctx context.Context, employee *models.Employee) error
	GetByID(ctx context.Context, id uint) (*models.Employee, error)
	GetByCode(ctx context.Context, code string) (*models.Employee, error)
	GetByEmail(ctx context.Context, email string) (*models.Employee, error)
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

// GetByEmail 根据邮箱获取员工
func (r *EmployeeRepositoryImpl) GetByEmail(ctx context.Context, email string) (*models.Employee, error) {
	var employee models.Employee
	err := r.db.WithContext(ctx).Preload("Department").Preload("Position").Preload("Manager").Where("email = ?", email).First(&employee).Error
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

// ===== 请假管理 =====

// LeaveRepository 请假仓储接口
type LeaveRepository interface {
	Create(ctx context.Context, leave *models.Leave) error
	GetByID(ctx context.Context, id uint) (*models.Leave, error)
	Update(ctx context.Context, leave *models.Leave) error
	Delete(ctx context.Context, id uint) error
	GetByEmployeeID(ctx context.Context, employeeID uint, offset, limit int) ([]*models.Leave, int64, error)
	GetByEmployeeIDAndDateRange(ctx context.Context, employeeID uint, startDate, endDate time.Time) ([]*models.Leave, error)
	GetByDateRange(ctx context.Context, startDate, endDate time.Time, offset, limit int) ([]*models.Leave, int64, error)
	GetByStatus(ctx context.Context, status string, offset, limit int) ([]*models.Leave, int64, error)
	List(ctx context.Context, offset, limit int) ([]*models.Leave, int64, error)
	ListWithFilters(ctx context.Context, employeeID *uint, leaveType, status string, startDate, endDate *time.Time, offset, limit int) ([]*models.Leave, int64, error)
}

// LeaveRepositoryImpl 请假仓储实现
type LeaveRepositoryImpl struct {
	db *gorm.DB
}

// NewLeaveRepository 创建请假仓储实例
func NewLeaveRepository(db *gorm.DB) LeaveRepository {
	return &LeaveRepositoryImpl{
		db: db,
	}
}

// Create 创建请假申请
func (r *LeaveRepositoryImpl) Create(ctx context.Context, leave *models.Leave) error {
	return r.db.WithContext(ctx).Create(leave).Error
}

// GetByID 根据ID获取请假申请
func (r *LeaveRepositoryImpl) GetByID(ctx context.Context, id uint) (*models.Leave, error) {
	var leave models.Leave
	err := r.db.WithContext(ctx).Preload("Employee").Preload("ApprovedUser").First(&leave, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &leave, nil
}

// Update 更新请假申请
func (r *LeaveRepositoryImpl) Update(ctx context.Context, leave *models.Leave) error {
	return r.db.WithContext(ctx).Save(leave).Error
}

// Delete 删除请假申请
func (r *LeaveRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Leave{}, id).Error
}

// GetByEmployeeID 根据员工ID获取请假申请列表
func (r *LeaveRepositoryImpl) GetByEmployeeID(ctx context.Context, employeeID uint, offset, limit int) ([]*models.Leave, int64, error) {
	var leaves []*models.Leave
	var total int64

	// 获取总数
	if err := r.db.WithContext(ctx).Model(&models.Leave{}).Where("employee_id = ?", employeeID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := r.db.WithContext(ctx).
		Preload("Employee").
		Preload("ApprovedUser").
		Where("employee_id = ?", employeeID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&leaves).Error

	return leaves, total, err
}

// GetByEmployeeIDAndDateRange 根据员工ID和日期范围获取请假申请
func (r *LeaveRepositoryImpl) GetByEmployeeIDAndDateRange(ctx context.Context, employeeID uint, startDate, endDate time.Time) ([]*models.Leave, error) {
	var leaves []*models.Leave
	err := r.db.WithContext(ctx).
		Preload("Employee").
		Preload("ApprovedUser").
		Where("employee_id = ? AND ((start_date BETWEEN ? AND ?) OR (end_date BETWEEN ? AND ?) OR (start_date <= ? AND end_date >= ?))",
			employeeID, startDate, endDate, startDate, endDate, startDate, endDate).
		Order("start_date ASC").
		Find(&leaves).Error
	return leaves, err
}

// GetByDateRange 根据日期范围获取请假申请列表
func (r *LeaveRepositoryImpl) GetByDateRange(ctx context.Context, startDate, endDate time.Time, offset, limit int) ([]*models.Leave, int64, error) {
	var leaves []*models.Leave
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Leave{}).
		Where("(start_date BETWEEN ? AND ?) OR (end_date BETWEEN ? AND ?) OR (start_date <= ? AND end_date >= ?)",
			startDate, endDate, startDate, endDate, startDate, endDate)

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := r.db.WithContext(ctx).
		Preload("Employee").
		Preload("ApprovedUser").
		Where("(start_date BETWEEN ? AND ?) OR (end_date BETWEEN ? AND ?) OR (start_date <= ? AND end_date >= ?)",
			startDate, endDate, startDate, endDate, startDate, endDate).
		Order("start_date ASC").
		Offset(offset).
		Limit(limit).
		Find(&leaves).Error

	return leaves, total, err
}

// GetByStatus 根据状态获取请假申请列表
func (r *LeaveRepositoryImpl) GetByStatus(ctx context.Context, status string, offset, limit int) ([]*models.Leave, int64, error) {
	var leaves []*models.Leave
	var total int64

	// 获取总数
	if err := r.db.WithContext(ctx).Model(&models.Leave{}).Where("status = ?", status).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := r.db.WithContext(ctx).
		Preload("Employee").
		Preload("ApprovedUser").
		Where("status = ?", status).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&leaves).Error

	return leaves, total, err
}

// List 获取请假申请列表
func (r *LeaveRepositoryImpl) List(ctx context.Context, offset, limit int) ([]*models.Leave, int64, error) {
	var leaves []*models.Leave
	var total int64

	// 获取总数
	if err := r.db.WithContext(ctx).Model(&models.Leave{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := r.db.WithContext(ctx).
		Preload("Employee").
		Preload("ApprovedUser").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&leaves).Error

	return leaves, total, err
}

// ListWithFilters 根据筛选条件获取请假申请列表
func (r *LeaveRepositoryImpl) ListWithFilters(ctx context.Context, employeeID *uint, leaveType, status string, startDate, endDate *time.Time, offset, limit int) ([]*models.Leave, int64, error) {
	var leaves []*models.Leave
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Leave{})

	// 添加筛选条件
	if employeeID != nil {
		query = query.Where("employee_id = ?", *employeeID)
	}
	if leaveType != "" {
		query = query.Where("leave_type = ?", leaveType)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if startDate != nil && endDate != nil {
		query = query.Where("(start_date BETWEEN ? AND ?) OR (end_date BETWEEN ? AND ?) OR (start_date <= ? AND end_date >= ?)",
			*startDate, *endDate, *startDate, *endDate, *startDate, *endDate)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := query.
		Preload("Employee").
		Preload("ApprovedUser").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&leaves).Error

	return leaves, total, err
}

// AttendanceRepository 考勤仓储接口
type AttendanceRepository interface {
	Create(ctx context.Context, attendance *models.Attendance) error
	GetByID(ctx context.Context, id uint) (*models.Attendance, error)
	Update(ctx context.Context, attendance *models.Attendance) error
	Delete(ctx context.Context, id uint) error
	GetByEmployeeID(ctx context.Context, employeeID uint, offset, limit int) ([]*models.Attendance, int64, error)
	GetByEmployeeIDAndDateRange(ctx context.Context, employeeID uint, startDate, endDate time.Time) ([]*models.Attendance, error)
	GetByDateRange(ctx context.Context, startDate, endDate time.Time, offset, limit int) ([]*models.Attendance, int64, error)
	List(ctx context.Context, offset, limit int) ([]*models.Attendance, int64, error)
}

// AttendanceRepositoryImpl 考勤仓储实现
type AttendanceRepositoryImpl struct {
	db *gorm.DB
}

// NewAttendanceRepository 创建考勤仓储实例
func NewAttendanceRepository(db *gorm.DB) AttendanceRepository {
	return &AttendanceRepositoryImpl{
		db: db,
	}
}

// Create 创建考勤记录
func (r *AttendanceRepositoryImpl) Create(ctx context.Context, attendance *models.Attendance) error {
	return r.db.WithContext(ctx).Create(attendance).Error
}

// GetByID 根据ID获取考勤记录
func (r *AttendanceRepositoryImpl) GetByID(ctx context.Context, id uint) (*models.Attendance, error) {
	var attendance models.Attendance
	err := r.db.WithContext(ctx).Preload("Employee").First(&attendance, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &attendance, nil
}

// Update 更新考勤记录
func (r *AttendanceRepositoryImpl) Update(ctx context.Context, attendance *models.Attendance) error {
	return r.db.WithContext(ctx).Save(attendance).Error
}

// Delete 删除考勤记录
func (r *AttendanceRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Attendance{}, id).Error
}

// GetByEmployeeID 根据员工ID获取考勤记录
func (r *AttendanceRepositoryImpl) GetByEmployeeID(ctx context.Context, employeeID uint, offset, limit int) ([]*models.Attendance, int64, error) {
	var attendances []*models.Attendance
	var total int64

	// 获取总数
	if err := r.db.WithContext(ctx).Model(&models.Attendance{}).
		Where("employee_id = ?", employeeID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := r.db.WithContext(ctx).Preload("Employee").
		Where("employee_id = ?", employeeID).
		Order("date DESC").
		Offset(offset).Limit(limit).Find(&attendances).Error
	if err != nil {
		return nil, 0, err
	}

	return attendances, total, nil
}

// GetByEmployeeIDAndDateRange 根据员工ID和日期范围获取考勤记录
func (r *AttendanceRepositoryImpl) GetByEmployeeIDAndDateRange(ctx context.Context, employeeID uint, startDate, endDate time.Time) ([]*models.Attendance, error) {
	var attendances []*models.Attendance
	err := r.db.WithContext(ctx).Preload("Employee").
		Where("employee_id = ? AND date >= ? AND date <= ?", employeeID, startDate, endDate).
		Order("date ASC").Find(&attendances).Error
	if err != nil {
		return nil, err
	}
	return attendances, nil
}

// GetByDateRange 根据日期范围获取考勤记录
func (r *AttendanceRepositoryImpl) GetByDateRange(ctx context.Context, startDate, endDate time.Time, offset, limit int) ([]*models.Attendance, int64, error) {
	var attendances []*models.Attendance
	var total int64

	// 获取总数
	if err := r.db.WithContext(ctx).Model(&models.Attendance{}).
		Where("date >= ? AND date <= ?", startDate, endDate).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := r.db.WithContext(ctx).Preload("Employee").
		Where("date >= ? AND date <= ?", startDate, endDate).
		Order("date DESC").
		Offset(offset).Limit(limit).Find(&attendances).Error
	if err != nil {
		return nil, 0, err
	}

	return attendances, total, nil
}

// List 获取考勤记录列表
func (r *AttendanceRepositoryImpl) List(ctx context.Context, offset, limit int) ([]*models.Attendance, int64, error) {
	var attendances []*models.Attendance
	var total int64

	// 获取总数
	if err := r.db.WithContext(ctx).Model(&models.Attendance{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := r.db.WithContext(ctx).Preload("Employee").
		Order("date DESC").
		Offset(offset).Limit(limit).Find(&attendances).Error
	if err != nil {
		return nil, 0, err
	}

	return attendances, total, nil
}

// PayrollRepository 薪资仓储接口
type PayrollRepository interface {
	Create(ctx context.Context, payroll *models.Payroll) error
	GetByID(ctx context.Context, id uint) (*models.Payroll, error)
	Update(ctx context.Context, payroll *models.Payroll) error
	Delete(ctx context.Context, id uint) error
	GetByEmployeeID(ctx context.Context, employeeID uint, offset, limit int) ([]*models.Payroll, int64, error)
	GetByEmployeeIDAndPeriod(ctx context.Context, employeeID uint, startDate, endDate time.Time) (*models.Payroll, error)
	GetByPeriod(ctx context.Context, startDate, endDate time.Time, offset, limit int) ([]*models.Payroll, int64, error)
	List(ctx context.Context, offset, limit int) ([]*models.Payroll, int64, error)
}

// PayrollRepositoryImpl 薪资仓储实现
type PayrollRepositoryImpl struct {
	db *gorm.DB
}

// NewPayrollRepository 创建薪资仓储实例
func NewPayrollRepository(db *gorm.DB) PayrollRepository {
	return &PayrollRepositoryImpl{
		db: db,
	}
}

// Create 创建薪资记录
func (r *PayrollRepositoryImpl) Create(ctx context.Context, payroll *models.Payroll) error {
	return r.db.WithContext(ctx).Create(payroll).Error
}

// GetByID 根据ID获取薪资记录
func (r *PayrollRepositoryImpl) GetByID(ctx context.Context, id uint) (*models.Payroll, error) {
	var payroll models.Payroll
	err := r.db.WithContext(ctx).Preload("Employee").First(&payroll, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &payroll, nil
}

// Update 更新薪资记录
func (r *PayrollRepositoryImpl) Update(ctx context.Context, payroll *models.Payroll) error {
	return r.db.WithContext(ctx).Save(payroll).Error
}

// Delete 删除薪资记录
func (r *PayrollRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Payroll{}, id).Error
}

// GetByEmployeeID 根据员工ID获取薪资记录
func (r *PayrollRepositoryImpl) GetByEmployeeID(ctx context.Context, employeeID uint, offset, limit int) ([]*models.Payroll, int64, error) {
	var payrolls []*models.Payroll
	var total int64

	// 获取总数
	if err := r.db.WithContext(ctx).Model(&models.Payroll{}).
		Where("employee_id = ?", employeeID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := r.db.WithContext(ctx).Preload("Employee").
		Where("employee_id = ?", employeeID).
		Order("pay_period_start DESC").
		Offset(offset).Limit(limit).Find(&payrolls).Error
	if err != nil {
		return nil, 0, err
	}

	return payrolls, total, nil
}

// GetByEmployeeIDAndPeriod 根据员工ID和薪资周期获取薪资记录
func (r *PayrollRepositoryImpl) GetByEmployeeIDAndPeriod(ctx context.Context, employeeID uint, startDate, endDate time.Time) (*models.Payroll, error) {
	var payroll models.Payroll
	err := r.db.WithContext(ctx).Preload("Employee").
		Where("employee_id = ? AND pay_period_start = ? AND pay_period_end = ?", employeeID, startDate, endDate).
		First(&payroll).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &payroll, nil
}

// GetByPeriod 根据薪资周期获取薪资记录
func (r *PayrollRepositoryImpl) GetByPeriod(ctx context.Context, startDate, endDate time.Time, offset, limit int) ([]*models.Payroll, int64, error) {
	var payrolls []*models.Payroll
	var total int64

	// 获取总数
	if err := r.db.WithContext(ctx).Model(&models.Payroll{}).
		Where("pay_period_start = ? AND pay_period_end = ?", startDate, endDate).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := r.db.WithContext(ctx).Preload("Employee").
		Where("pay_period_start = ? AND pay_period_end = ?", startDate, endDate).
		Order("created_at DESC").
		Offset(offset).Limit(limit).Find(&payrolls).Error
	if err != nil {
		return nil, 0, err
	}

	return payrolls, total, nil
}

// List 获取薪资记录列表
func (r *PayrollRepositoryImpl) List(ctx context.Context, offset, limit int) ([]*models.Payroll, int64, error) {
	var payrolls []*models.Payroll
	var total int64

	// 获取总数
	if err := r.db.WithContext(ctx).Model(&models.Payroll{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := r.db.WithContext(ctx).Preload("Employee").
		Order("pay_period_start DESC").
		Offset(offset).Limit(limit).Find(&payrolls).Error
	if err != nil {
		return nil, 0, err
	}

	return payrolls, total, nil
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
