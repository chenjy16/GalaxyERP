package repositories

import (
	"context"
	"github.com/galaxyerp/galaxyErp/internal/models"
	"gorm.io/gorm"
	"time"
)

// EmployeeRepository 员工仓储接口
type EmployeeRepository interface {
	BaseRepository[models.Employee]
	GetByCode(ctx context.Context, code string) (*models.Employee, error)
	GetByEmail(ctx context.Context, email string) (*models.Employee, error)
	Search(ctx context.Context, query string, offset, limit int) ([]*models.Employee, int64, error)
	GetByDepartmentID(ctx context.Context, departmentID uint, offset, limit int) ([]*models.Employee, int64, error)
}

// EmployeeRepositoryImpl 员工仓储实现
type EmployeeRepositoryImpl struct {
	BaseRepository[models.Employee]
	db *gorm.DB
}

// NewEmployeeRepository 创建员工仓储实例
func NewEmployeeRepository(db *gorm.DB) EmployeeRepository {
	return &EmployeeRepositoryImpl{
		BaseRepository: NewBaseRepository[models.Employee](db),
		db:             db,
	}
}

// GetByCode 根据员工编号获取员工
func (r *EmployeeRepositoryImpl) GetByCode(ctx context.Context, code string) (*models.Employee, error) {
	var employee models.Employee
	err := r.db.WithContext(ctx).Preload("Department").Preload("Position").Preload("Manager").Where("code = ?", code).First(&employee).Error
	if err != nil {
		return nil, err
	}
	return &employee, nil
}

// GetByEmail 根据邮箱获取员工
func (r *EmployeeRepositoryImpl) GetByEmail(ctx context.Context, email string) (*models.Employee, error) {
	var employee models.Employee
	err := r.db.WithContext(ctx).Preload("Department").Preload("Position").Preload("Manager").Where("email = ?", email).First(&employee).Error
	if err != nil {
		return nil, err
	}
	return &employee, nil
}

// Search 搜索员工
func (r *EmployeeRepositoryImpl) Search(ctx context.Context, query string, offset, limit int) ([]*models.Employee, int64, error) {
	var employees []*models.Employee
	var total int64

	searchQuery := "%" + query + "%"

	// 获取总数
	if err := r.db.WithContext(ctx).Model(&models.Employee{}).
		Where("code LIKE ? OR first_name LIKE ? OR last_name LIKE ? OR email LIKE ?",
			searchQuery, searchQuery, searchQuery, searchQuery).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := r.db.WithContext(ctx).Preload("Department").Preload("Position").Preload("Manager").
		Where("code LIKE ? OR first_name LIKE ? OR last_name LIKE ? OR email LIKE ?",
			searchQuery, searchQuery, searchQuery, searchQuery).
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
	err := r.db.WithContext(ctx).Preload("Department").Preload("Position").Preload("Manager").
		Where("department_id = ?", departmentID).
		Offset(offset).Limit(limit).Find(&employees).Error
	if err != nil {
		return nil, 0, err
	}

	return employees, total, nil
}

// LeaveRepository 请假仓储接口
type LeaveRepository interface {
	BaseRepository[models.Leave]
	GetByEmployeeID(ctx context.Context, employeeID uint, offset, limit int) ([]*models.Leave, int64, error)
	GetByEmployeeIDAndDateRange(ctx context.Context, employeeID uint, startDate, endDate time.Time) ([]*models.Leave, error)
	GetByDateRange(ctx context.Context, startDate, endDate time.Time, offset, limit int) ([]*models.Leave, int64, error)
	GetByStatus(ctx context.Context, status string, offset, limit int) ([]*models.Leave, int64, error)
	ListWithFilters(ctx context.Context, employeeID *uint, leaveType, status string, startDate, endDate *time.Time, offset, limit int) ([]*models.Leave, int64, error)
}

// LeaveRepositoryImpl 请假仓储实现
type LeaveRepositoryImpl struct {
	BaseRepository[models.Leave]
	db *gorm.DB
}

// NewLeaveRepository 创建请假仓储实例
func NewLeaveRepository(db *gorm.DB) LeaveRepository {
	return &LeaveRepositoryImpl{
		BaseRepository: NewBaseRepository[models.Leave](db),
		db:             db,
	}
}

// GetByEmployeeID 根据员工ID获取请假记录
func (r *LeaveRepositoryImpl) GetByEmployeeID(ctx context.Context, employeeID uint, offset, limit int) ([]*models.Leave, int64, error) {
	var leaves []*models.Leave
	var total int64

	// 获取总数
	if err := r.db.WithContext(ctx).Model(&models.Leave{}).
		Where("employee_id = ?", employeeID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := r.db.WithContext(ctx).Preload("Employee").
		Where("employee_id = ?", employeeID).
		Offset(offset).Limit(limit).Find(&leaves).Error
	if err != nil {
		return nil, 0, err
	}

	return leaves, total, nil
}

// GetByEmployeeIDAndDateRange 根据员工ID和日期范围获取请假记录
func (r *LeaveRepositoryImpl) GetByEmployeeIDAndDateRange(ctx context.Context, employeeID uint, startDate, endDate time.Time) ([]*models.Leave, error) {
	var leaves []*models.Leave
	err := r.db.WithContext(ctx).Preload("Employee").
		Where("employee_id = ? AND start_date >= ? AND end_date <= ?", employeeID, startDate, endDate).
		Find(&leaves).Error
	return leaves, err
}

// GetByDateRange 根据日期范围获取请假记录
func (r *LeaveRepositoryImpl) GetByDateRange(ctx context.Context, startDate, endDate time.Time, offset, limit int) ([]*models.Leave, int64, error) {
	var leaves []*models.Leave
	var total int64

	// 获取总数
	if err := r.db.WithContext(ctx).Model(&models.Leave{}).
		Where("start_date >= ? AND end_date <= ?", startDate, endDate).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := r.db.WithContext(ctx).Preload("Employee").
		Where("start_date >= ? AND end_date <= ?", startDate, endDate).
		Offset(offset).Limit(limit).Find(&leaves).Error
	if err != nil {
		return nil, 0, err
	}

	return leaves, total, nil
}

// GetByStatus 根据状态获取请假记录
func (r *LeaveRepositoryImpl) GetByStatus(ctx context.Context, status string, offset, limit int) ([]*models.Leave, int64, error) {
	var leaves []*models.Leave
	var total int64

	// 获取总数
	if err := r.db.WithContext(ctx).Model(&models.Leave{}).
		Where("status = ?", status).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := r.db.WithContext(ctx).Preload("Employee").
		Where("status = ?", status).
		Offset(offset).Limit(limit).Find(&leaves).Error
	if err != nil {
		return nil, 0, err
	}

	return leaves, total, nil
}

// ListWithFilters 根据过滤条件获取请假记录
func (r *LeaveRepositoryImpl) ListWithFilters(ctx context.Context, employeeID *uint, leaveType, status string, startDate, endDate *time.Time, offset, limit int) ([]*models.Leave, int64, error) {
	var leaves []*models.Leave
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Leave{})

	// 应用过滤条件
	if employeeID != nil {
		query = query.Where("employee_id = ?", *employeeID)
	}
	if leaveType != "" {
		query = query.Where("leave_type = ?", leaveType)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if startDate != nil {
		query = query.Where("start_date >= ?", *startDate)
	}
	if endDate != nil {
		query = query.Where("end_date <= ?", *endDate)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := query.Preload("Employee").Offset(offset).Limit(limit).Find(&leaves).Error
	if err != nil {
		return nil, 0, err
	}

	return leaves, total, nil
}

// AttendanceRepository 考勤仓储接口
type AttendanceRepository interface {
	BaseRepository[models.Attendance]
	GetByEmployeeID(ctx context.Context, employeeID uint, offset, limit int) ([]*models.Attendance, int64, error)
	GetByEmployeeIDAndDateRange(ctx context.Context, employeeID uint, startDate, endDate time.Time) ([]*models.Attendance, error)
	GetByDateRange(ctx context.Context, startDate, endDate time.Time, offset, limit int) ([]*models.Attendance, int64, error)
}

// AttendanceRepositoryImpl 考勤仓储实现
type AttendanceRepositoryImpl struct {
	BaseRepository[models.Attendance]
	db *gorm.DB
}

// NewAttendanceRepository 创建考勤仓储实例
func NewAttendanceRepository(db *gorm.DB) AttendanceRepository {
	return &AttendanceRepositoryImpl{
		BaseRepository: NewBaseRepository[models.Attendance](db),
		db:             db,
	}
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
		Find(&attendances).Error
	return attendances, err
}

// GetByDateRange 根据日期范围获取考勤记录
func (r *AttendanceRepositoryImpl) GetByDateRange(ctx context.Context, startDate, endDate time.Time, offset, limit int) ([]*models.Attendance, int64, error) {
	var attendances []*models.Attendance
	var total int64

	// 获取总数
	if err := r.db.WithContext(ctx).Model(&models.Attendance{}).
		Where("date >= ? AND date <= ?", startDate, endDate).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := r.db.WithContext(ctx).Preload("Employee").
		Where("date >= ? AND date <= ?", startDate, endDate).
		Offset(offset).Limit(limit).Find(&attendances).Error
	if err != nil {
		return nil, 0, err
	}

	return attendances, total, nil
}

// PayrollRepository 薪资仓储接口
type PayrollRepository interface {
	BaseRepository[models.Payroll]
	GetByEmployeeID(ctx context.Context, employeeID uint, offset, limit int) ([]*models.Payroll, int64, error)
	GetByEmployeeIDAndPeriod(ctx context.Context, employeeID uint, startDate, endDate time.Time) (*models.Payroll, error)
	GetByPeriod(ctx context.Context, startDate, endDate time.Time, offset, limit int) ([]*models.Payroll, int64, error)
}

// PayrollRepositoryImpl 薪资仓储实现
type PayrollRepositoryImpl struct {
	BaseRepository[models.Payroll]
	db *gorm.DB
}

// NewPayrollRepository 创建薪资仓储实例
func NewPayrollRepository(db *gorm.DB) PayrollRepository {
	return &PayrollRepositoryImpl{
		BaseRepository: NewBaseRepository[models.Payroll](db),
		db:             db,
	}
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
		Offset(offset).Limit(limit).Find(&payrolls).Error
	if err != nil {
		return nil, 0, err
	}

	return payrolls, total, nil
}

// GetByEmployeeIDAndPeriod 根据员工ID和期间获取薪资记录
func (r *PayrollRepositoryImpl) GetByEmployeeIDAndPeriod(ctx context.Context, employeeID uint, startDate, endDate time.Time) (*models.Payroll, error) {
	var payroll models.Payroll
	err := r.db.WithContext(ctx).Preload("Employee").
		Where("employee_id = ? AND period_start = ? AND period_end = ?", employeeID, startDate, endDate).
		First(&payroll).Error
	if err != nil {
		return nil, err
	}
	return &payroll, nil
}

// GetByPeriod 根据期间获取薪资记录
func (r *PayrollRepositoryImpl) GetByPeriod(ctx context.Context, startDate, endDate time.Time, offset, limit int) ([]*models.Payroll, int64, error) {
	var payrolls []*models.Payroll
	var total int64

	// 获取总数
	if err := r.db.WithContext(ctx).Model(&models.Payroll{}).
		Where("period_start = ? AND period_end = ?", startDate, endDate).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := r.db.WithContext(ctx).Preload("Employee").
		Where("period_start = ? AND period_end = ?", startDate, endDate).
		Offset(offset).Limit(limit).Find(&payrolls).Error
	if err != nil {
		return nil, 0, err
	}

	return payrolls, total, nil
}
