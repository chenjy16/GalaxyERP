package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/galaxyerp/galaxyErp/internal/common"
	"github.com/galaxyerp/galaxyErp/internal/dto"
	"github.com/galaxyerp/galaxyErp/internal/models"
	"github.com/galaxyerp/galaxyErp/internal/repositories"
)

// EmployeeService 员工服务接口
type EmployeeService interface {
	Create(ctx context.Context, req *dto.EmployeeCreateRequest) (*dto.EmployeeResponse, error)
	GetByID(ctx context.Context, id uint) (*dto.EmployeeResponse, error)
	Update(ctx context.Context, id uint, req *dto.EmployeeUpdateRequest) (*dto.EmployeeResponse, error)
	Delete(ctx context.Context, id uint) error
	GetByCode(ctx context.Context, code string) (*dto.EmployeeResponse, error)
	GetByDepartmentID(ctx context.Context, departmentID uint, req *dto.PaginationRequest) (*dto.PaginatedResponse[dto.EmployeeListResponse], error)
	List(ctx context.Context, req *dto.PaginationRequest) (*dto.PaginatedResponse[dto.EmployeeListResponse], error)
	Search(ctx context.Context, req *dto.EmployeeSearchRequest) (*dto.PaginatedResponse[dto.EmployeeListResponse], error)
}

// EmployeeServiceImpl 员工服务实现
type EmployeeServiceImpl struct {
	employeeRepo repositories.EmployeeRepository
}

// NewEmployeeService 创建员工服务
func NewEmployeeService(employeeRepo repositories.EmployeeRepository) EmployeeService {
	return &EmployeeServiceImpl{
		employeeRepo: employeeRepo,
	}
}

// Create 创建员工
func (s *EmployeeServiceImpl) Create(ctx context.Context, req *dto.EmployeeCreateRequest) (*dto.EmployeeResponse, error) {
	// 检查员工编码是否已存在
	existing, err := s.employeeRepo.GetByCode(ctx, req.Code)
	if err != nil {
		return nil, fmt.Errorf("检查员工编码失败: %w", err)
	}
	if existing != nil {
		return nil, errors.New("员工编码已存在")
	}

	// 检查邮箱是否已存在
	existingByEmail, err := s.employeeRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("检查员工邮箱失败: %w", err)
	}
	if existingByEmail != nil {
		return nil, errors.New("员工邮箱已存在")
	}

	// 创建员工
	employee := &models.Employee{
		Code:             req.Code,
		FirstName:        req.FirstName,
		LastName:         req.LastName,
		FullName:         req.FirstName + " " + req.LastName,
		Email:            req.Email,
		Phone:            req.Phone,
		DateOfBirth:      req.DateOfBirth,
		Gender:           req.Gender,
		HireDate:         &req.HireDate,
		DepartmentID:     req.DepartmentID,
		PositionID:       req.PositionID,
		ManagerID:        req.ManagerID,
		Status:           req.Status,
		EmergencyContact: req.EmergencyContact,
		IDNumber:         req.IDNumber,
		Address:          req.Address,
	}

	// 保存到数据库
	if err := s.employeeRepo.Create(ctx, employee); err != nil {
		return nil, fmt.Errorf("创建员工失败: %w", err)
	}

	// 重新获取完整信息
	created, err := s.employeeRepo.GetByID(ctx, employee.ID)
	if err != nil {
		return nil, fmt.Errorf("获取创建的员工失败: %w", err)
	}

	return s.toEmployeeResponse(created), nil
}

// GetByID 根据ID获取员工
func (s *EmployeeServiceImpl) GetByID(ctx context.Context, id uint) (*dto.EmployeeResponse, error) {
	employee, err := s.employeeRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("获取员工失败: %w", err)
	}
	if employee == nil {
		return nil, errors.New("员工不存在")
	}

	return s.toEmployeeResponse(employee), nil
}

// Update 更新员工
func (s *EmployeeServiceImpl) Update(ctx context.Context, id uint, req *dto.EmployeeUpdateRequest) (*dto.EmployeeResponse, error) {
	// 获取现有员工
	employee, err := s.employeeRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("获取员工失败: %w", err)
	}
	if employee == nil {
		return nil, errors.New("员工不存在")
	}

	// 检查邮箱是否已被其他员工使用
	if req.Email != "" && req.Email != employee.Email {
		existingByEmail, err := s.employeeRepo.GetByEmail(ctx, req.Email)
		if err != nil {
			return nil, fmt.Errorf("检查员工邮箱失败: %w", err)
		}
		if existingByEmail != nil && existingByEmail.ID != id {
			return nil, errors.New("员工邮箱已存在")
		}
	}

	// 更新字段
	if req.FirstName != "" {
		employee.FirstName = req.FirstName
	}
	if req.LastName != "" {
		employee.LastName = req.LastName
	}
	if req.FirstName != "" || req.LastName != "" {
		employee.FullName = employee.FirstName + " " + employee.LastName
	}
	if req.Email != "" {
		employee.Email = req.Email
	}
	if req.Phone != "" {
		employee.Phone = req.Phone
	}
	if req.DateOfBirth != nil {
		employee.DateOfBirth = req.DateOfBirth
	}
	if req.Gender != "" {
		employee.Gender = req.Gender
	}
	if req.DepartmentID != nil {
		employee.DepartmentID = req.DepartmentID
	}
	if req.PositionID != nil {
		employee.PositionID = req.PositionID
	}
	if req.ManagerID != nil {
		employee.ManagerID = req.ManagerID
	}
	if req.Status != "" {
		employee.Status = req.Status
	}
	if req.EmergencyContact != "" {
		employee.EmergencyContact = req.EmergencyContact
	}
	if req.IDNumber != "" {
		employee.IDNumber = req.IDNumber
	}
	if req.Address != "" {
		employee.Address = req.Address
	}

	// 保存更新
	if err := s.employeeRepo.Update(ctx, employee); err != nil {
		return nil, fmt.Errorf("更新员工失败: %w", err)
	}

	// 重新获取完整信息
	updated, err := s.employeeRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("获取更新的员工失败: %w", err)
	}

	return s.toEmployeeResponse(updated), nil
}

// Delete 删除员工
func (s *EmployeeServiceImpl) Delete(ctx context.Context, id uint) error {
	// 检查员工是否存在
	employee, err := s.employeeRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("获取员工失败: %w", err)
	}
	if employee == nil {
		return errors.New("员工不存在")
	}

	return s.employeeRepo.Delete(ctx, id)
}

// GetByCode 根据编码获取员工
func (s *EmployeeServiceImpl) GetByCode(ctx context.Context, code string) (*dto.EmployeeResponse, error) {
	employee, err := s.employeeRepo.GetByCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("获取员工失败: %w", err)
	}
	if employee == nil {
		return nil, errors.New("员工不存在")
	}

	return s.toEmployeeResponse(employee), nil
}

// GetByDepartmentID 根据部门ID获取员工列表
func (s *EmployeeServiceImpl) GetByDepartmentID(ctx context.Context, departmentID uint, req *dto.PaginationRequest) (*dto.PaginatedResponse[dto.EmployeeListResponse], error) {
	offset := (req.Page - 1) * req.PageSize
	employees, total, err := s.employeeRepo.GetByDepartmentID(ctx, departmentID, offset, req.PageSize)
	if err != nil {
		return nil, fmt.Errorf("获取部门员工列表失败: %w", err)
	}

	items := make([]dto.EmployeeListResponse, len(employees))
	for i, employee := range employees {
		items[i] = *s.toEmployeeListResponse(employee)
	}

	totalPages := int(total) / req.PageSize
	if int(total)%req.PageSize > 0 {
		totalPages++
	}

	return &dto.PaginatedResponse[dto.EmployeeListResponse]{
		Data:       items,
		Total:      total,
		Page:       req.Page,
		Limit:      req.PageSize,
		TotalPages: totalPages,
	}, nil
}

// List 获取员工列表
func (s *EmployeeServiceImpl) List(ctx context.Context, req *dto.PaginationRequest) (*dto.PaginatedResponse[dto.EmployeeListResponse], error) {
	options := &common.QueryOptions{
		Pagination: req,
	}
	employees, total, err := s.employeeRepo.List(ctx, options)
	if err != nil {
		return nil, fmt.Errorf("获取员工列表失败: %w", err)
	}

	items := make([]dto.EmployeeListResponse, len(employees))
	for i, employee := range employees {
		items[i] = *s.toEmployeeListResponse(employee)
	}

	totalPages := int(total) / req.PageSize
	if int(total)%req.PageSize > 0 {
		totalPages++
	}

	return &dto.PaginatedResponse[dto.EmployeeListResponse]{
		Data:       items,
		Total:      total,
		Page:       req.Page,
		Limit:      req.PageSize,
		TotalPages: totalPages,
	}, nil
}

// Search 搜索员工
func (s *EmployeeServiceImpl) Search(ctx context.Context, req *dto.EmployeeSearchRequest) (*dto.PaginatedResponse[dto.EmployeeListResponse], error) {
	offset := (req.Page - 1) * req.PageSize
	employees, total, err := s.employeeRepo.Search(ctx, req.Keyword, offset, req.PageSize)
	if err != nil {
		return nil, fmt.Errorf("搜索员工失败: %w", err)
	}

	items := make([]dto.EmployeeListResponse, len(employees))
	for i, employee := range employees {
		items[i] = *s.toEmployeeListResponse(employee)
	}

	totalPages := int(total) / req.PageSize
	if int(total)%req.PageSize > 0 {
		totalPages++
	}

	return &dto.PaginatedResponse[dto.EmployeeListResponse]{
		Data:       items,
		Total:      total,
		Page:       req.Page,
		Limit:      req.PageSize,
		TotalPages: totalPages,
	}, nil
}

// toEmployeeResponse 转换为员工响应
func (s *EmployeeServiceImpl) toEmployeeResponse(employee *models.Employee) *dto.EmployeeResponse {
	var hireDate time.Time
	if employee.HireDate != nil {
		hireDate = *employee.HireDate
	}

	resp := &dto.EmployeeResponse{
		ID:               employee.ID,
		Code:             employee.Code,
		FirstName:        employee.FirstName,
		LastName:         employee.LastName,
		FullName:         employee.FullName,
		Email:            employee.Email,
		Phone:            employee.Phone,
		DateOfBirth:      employee.DateOfBirth,
		Gender:           employee.Gender,
		HireDate:         hireDate,
		DepartmentID:     employee.DepartmentID,
		PositionID:       employee.PositionID,
		ManagerID:        employee.ManagerID,
		Status:           employee.Status,
		EmergencyContact: employee.EmergencyContact,
		IDNumber:         employee.IDNumber,
		Address:          employee.Address,
		CreatedAt:        employee.CreatedAt,
		UpdatedAt:        employee.UpdatedAt,
	}

	return resp
}

// toEmployeeListResponse 转换为员工列表响应
func (s *EmployeeServiceImpl) toEmployeeListResponse(employee *models.Employee) *dto.EmployeeListResponse {
	resp := &dto.EmployeeListResponse{
		ID:        employee.ID,
		Code:      employee.Code,
		FirstName: employee.FirstName,
		LastName:  employee.LastName,
		FullName:  employee.FullName,
		Email:     employee.Email,
		Phone:     employee.Phone,
		Status:    employee.Status,
	}

	// 添加部门和职位信息
	if employee.Department != nil {
		resp.Department = employee.Department.Name
	}
	if employee.Position != nil {
		resp.Position = employee.Position.Name
	}

	return resp
}

// AttendanceService 考勤服务接口
type AttendanceService interface {
	Create(ctx context.Context, req *dto.AttendanceCreateRequest) (*dto.AttendanceResponse, error)
	GetByID(ctx context.Context, id uint) (*dto.AttendanceResponse, error)
	Update(ctx context.Context, id uint, req *dto.AttendanceUpdateRequest) (*dto.AttendanceResponse, error)
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, req *dto.PaginationRequest) (*dto.PaginatedResponse[dto.AttendanceResponse], error)
	GetByEmployeeID(ctx context.Context, employeeID uint, req *dto.PaginationRequest) (*dto.PaginatedResponse[dto.AttendanceResponse], error)
	GetByDateRange(ctx context.Context, startDate, endDate time.Time, req *dto.PaginationRequest) (*dto.PaginatedResponse[dto.AttendanceResponse], error)
}

// AttendanceServiceImpl 考勤服务实现
type AttendanceServiceImpl struct {
	attendanceRepo repositories.AttendanceRepository
	employeeRepo   repositories.EmployeeRepository
}

// NewAttendanceService 创建考勤服务
func NewAttendanceService(attendanceRepo repositories.AttendanceRepository, employeeRepo repositories.EmployeeRepository) AttendanceService {
	return &AttendanceServiceImpl{
		attendanceRepo: attendanceRepo,
		employeeRepo:   employeeRepo,
	}
}

// Create 创建考勤记录
func (s *AttendanceServiceImpl) Create(ctx context.Context, req *dto.AttendanceCreateRequest) (*dto.AttendanceResponse, error) {
	// 检查员工是否存在
	employee, err := s.employeeRepo.GetByID(ctx, req.EmployeeID)
	if err != nil {
		return nil, fmt.Errorf("获取员工失败: %w", err)
	}
	if employee == nil {
		return nil, errors.New("员工不存在")
	}

	// 创建考勤记录
	attendance := &models.Attendance{
		EmployeeID:   req.EmployeeID,
		Date:         req.Date,
		CheckInTime:  req.CheckInTime,
		CheckOutTime: req.CheckOutTime,
		Status:       req.Status,
		Notes:        req.Notes,
	}

	// 保存到数据库
	if err := s.attendanceRepo.Create(ctx, attendance); err != nil {
		return nil, fmt.Errorf("创建考勤记录失败: %w", err)
	}

	// 重新获取完整信息
	created, err := s.attendanceRepo.GetByID(ctx, attendance.ID)
	if err != nil {
		return nil, fmt.Errorf("获取创建的考勤记录失败: %w", err)
	}

	return s.toAttendanceResponse(created), nil
}

// GetByID 根据ID获取考勤记录
func (s *AttendanceServiceImpl) GetByID(ctx context.Context, id uint) (*dto.AttendanceResponse, error) {
	attendance, err := s.attendanceRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("获取考勤记录失败: %w", err)
	}
	if attendance == nil {
		return nil, errors.New("考勤记录不存在")
	}

	return s.toAttendanceResponse(attendance), nil
}

// Update 更新考勤记录
func (s *AttendanceServiceImpl) Update(ctx context.Context, id uint, req *dto.AttendanceUpdateRequest) (*dto.AttendanceResponse, error) {
	// 获取现有考勤记录
	attendance, err := s.attendanceRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("获取考勤记录失败: %w", err)
	}
	if attendance == nil {
		return nil, errors.New("考勤记录不存在")
	}

	// 更新字段
	if req.CheckInTime != nil {
		attendance.CheckInTime = req.CheckInTime
	}
	if req.CheckOutTime != nil {
		attendance.CheckOutTime = req.CheckOutTime
	}
	if req.Status != "" {
		attendance.Status = req.Status
	}
	if req.Notes != "" {
		attendance.Notes = req.Notes
	}

	// 保存更新
	if err := s.attendanceRepo.Update(ctx, attendance); err != nil {
		return nil, fmt.Errorf("更新考勤记录失败: %w", err)
	}

	// 重新获取完整信息
	updated, err := s.attendanceRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("获取更新的考勤记录失败: %w", err)
	}

	return s.toAttendanceResponse(updated), nil
}

// Delete 删除考勤记录
func (s *AttendanceServiceImpl) Delete(ctx context.Context, id uint) error {
	// 检查考勤记录是否存在
	attendance, err := s.attendanceRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("获取考勤记录失败: %w", err)
	}
	if attendance == nil {
		return errors.New("考勤记录不存在")
	}

	return s.attendanceRepo.Delete(ctx, id)
}

// List 获取考勤记录列表
func (s *AttendanceServiceImpl) List(ctx context.Context, req *dto.PaginationRequest) (*dto.PaginatedResponse[dto.AttendanceResponse], error) {
	options := &common.QueryOptions{
		Pagination: req,
	}
	attendances, total, err := s.attendanceRepo.List(ctx, options)
	if err != nil {
		return nil, fmt.Errorf("获取考勤记录列表失败: %w", err)
	}

	items := make([]dto.AttendanceResponse, len(attendances))
	for i, attendance := range attendances {
		items[i] = *s.toAttendanceResponse(attendance)
	}

	totalPages := int(total) / req.PageSize
	if int(total)%req.PageSize > 0 {
		totalPages++
	}

	return &dto.PaginatedResponse[dto.AttendanceResponse]{
		Data:       items,
		Total:      total,
		Page:       req.Page,
		Limit:      req.PageSize,
		TotalPages: totalPages,
	}, nil
}

// GetByEmployeeID 根据员工ID获取考勤记录
func (s *AttendanceServiceImpl) GetByEmployeeID(ctx context.Context, employeeID uint, req *dto.PaginationRequest) (*dto.PaginatedResponse[dto.AttendanceResponse], error) {
	offset := (req.Page - 1) * req.PageSize
	attendances, total, err := s.attendanceRepo.GetByEmployeeID(ctx, employeeID, offset, req.PageSize)
	if err != nil {
		return nil, fmt.Errorf("获取员工考勤记录失败: %w", err)
	}

	items := make([]dto.AttendanceResponse, len(attendances))
	for i, attendance := range attendances {
		items[i] = *s.toAttendanceResponse(attendance)
	}

	totalPages := int(total) / req.PageSize
	if int(total)%req.PageSize > 0 {
		totalPages++
	}

	return &dto.PaginatedResponse[dto.AttendanceResponse]{
		Data:       items,
		Total:      total,
		Page:       req.Page,
		Limit:      req.PageSize,
		TotalPages: totalPages,
	}, nil
}

// GetByDateRange 根据日期范围获取考勤记录
func (s *AttendanceServiceImpl) GetByDateRange(ctx context.Context, startDate, endDate time.Time, req *dto.PaginationRequest) (*dto.PaginatedResponse[dto.AttendanceResponse], error) {
	offset := (req.Page - 1) * req.PageSize
	attendances, total, err := s.attendanceRepo.GetByDateRange(ctx, startDate, endDate, offset, req.PageSize)
	if err != nil {
		return nil, fmt.Errorf("获取日期范围考勤记录失败: %w", err)
	}

	items := make([]dto.AttendanceResponse, len(attendances))
	for i, attendance := range attendances {
		items[i] = *s.toAttendanceResponse(attendance)
	}

	totalPages := int(total) / req.PageSize
	if int(total)%req.PageSize > 0 {
		totalPages++
	}

	return &dto.PaginatedResponse[dto.AttendanceResponse]{
		Data:       items,
		Total:      total,
		Page:       req.Page,
		Limit:      req.PageSize,
		TotalPages: totalPages,
	}, nil
}

// toAttendanceResponse 转换为考勤响应
func (s *AttendanceServiceImpl) toAttendanceResponse(attendance *models.Attendance) *dto.AttendanceResponse {
	resp := &dto.AttendanceResponse{
		ID:           attendance.ID,
		EmployeeID:   attendance.EmployeeID,
		Date:         attendance.Date,
		CheckInTime:  attendance.CheckInTime,
		CheckOutTime: attendance.CheckOutTime,
		Status:       attendance.Status,
		Notes:        attendance.Notes,
		CreatedAt:    attendance.CreatedAt,
		UpdatedAt:    attendance.UpdatedAt,
	}

	// 添加员工信息
	if attendance.Employee.ID != 0 {
		resp.Employee = &dto.EmployeeListResponse{
			ID:        attendance.Employee.ID,
			Code:      attendance.Employee.Code,
			FirstName: attendance.Employee.FirstName,
			LastName:  attendance.Employee.LastName,
			FullName:  attendance.Employee.FullName,
		}
	}

	return resp
}

// PayrollService 薪资服务接口
type PayrollService interface {
	Create(ctx context.Context, req *dto.PayrollCreateRequest) (*dto.PayrollResponse, error)
	GetByID(ctx context.Context, id uint) (*dto.PayrollResponse, error)
	Update(ctx context.Context, id uint, req *dto.PayrollUpdateRequest) (*dto.PayrollResponse, error)
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, req *dto.PaginationRequest) (*dto.PaginatedResponse[dto.PayrollResponse], error)
	GetByEmployeeID(ctx context.Context, employeeID uint, req *dto.PaginationRequest) (*dto.PaginatedResponse[dto.PayrollResponse], error)
	GetByPeriod(ctx context.Context, startDate, endDate time.Time, req *dto.PaginationRequest) (*dto.PaginatedResponse[dto.PayrollResponse], error)
}

// PayrollServiceImpl 薪资服务实现
type PayrollServiceImpl struct {
	payrollRepo  repositories.PayrollRepository
	employeeRepo repositories.EmployeeRepository
}

// NewPayrollService 创建薪资服务
func NewPayrollService(payrollRepo repositories.PayrollRepository, employeeRepo repositories.EmployeeRepository) PayrollService {
	return &PayrollServiceImpl{
		payrollRepo:  payrollRepo,
		employeeRepo: employeeRepo,
	}
}

// Create 创建薪资记录
func (s *PayrollServiceImpl) Create(ctx context.Context, req *dto.PayrollCreateRequest) (*dto.PayrollResponse, error) {
	// 检查员工是否存在
	employee, err := s.employeeRepo.GetByID(ctx, req.EmployeeID)
	if err != nil {
		return nil, fmt.Errorf("获取员工失败: %w", err)
	}
	if employee == nil {
		return nil, errors.New("员工不存在")
	}

	// 创建薪资记录
	payroll := &models.Payroll{
		EmployeeID:      req.EmployeeID,
		PayPeriodStart:  req.PayPeriodStart,
		PayPeriodEnd:    req.PayPeriodEnd,
		BasicSalary:     req.BasicSalary,
		OvertimePay:     req.OvertimePay,
		Allowance:       req.Allowance,
		Bonus:           req.Bonus,
		Deductions:      req.Deductions,
		SocialInsurance: req.SocialInsurance,
		HousingFund:     req.HousingFund,
		Tax:             req.Tax,
		NetPay:          req.NetPay,
		Status:          req.Status,
	}

	// 保存到数据库
	if err := s.payrollRepo.Create(ctx, payroll); err != nil {
		return nil, fmt.Errorf("创建薪资记录失败: %w", err)
	}

	// 重新获取完整信息
	created, err := s.payrollRepo.GetByID(ctx, payroll.ID)
	if err != nil {
		return nil, fmt.Errorf("获取创建的薪资记录失败: %w", err)
	}

	return s.toPayrollResponse(created), nil
}

// GetByID 根据ID获取薪资记录
func (s *PayrollServiceImpl) GetByID(ctx context.Context, id uint) (*dto.PayrollResponse, error) {
	payroll, err := s.payrollRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("获取薪资记录失败: %w", err)
	}
	if payroll == nil {
		return nil, errors.New("薪资记录不存在")
	}

	return s.toPayrollResponse(payroll), nil
}

// Update 更新薪资记录
func (s *PayrollServiceImpl) Update(ctx context.Context, id uint, req *dto.PayrollUpdateRequest) (*dto.PayrollResponse, error) {
	// 获取现有薪资记录
	payroll, err := s.payrollRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("获取薪资记录失败: %w", err)
	}
	if payroll == nil {
		return nil, errors.New("薪资记录不存在")
	}

	// 更新字段
	if req.BasicSalary != nil {
		payroll.BasicSalary = *req.BasicSalary
	}
	if req.OvertimePay != nil {
		payroll.OvertimePay = *req.OvertimePay
	}
	if req.Allowance != nil {
		payroll.Allowance = *req.Allowance
	}
	if req.Bonus != nil {
		payroll.Bonus = *req.Bonus
	}
	if req.Deductions != nil {
		payroll.Deductions = *req.Deductions
	}
	if req.SocialInsurance != nil {
		payroll.SocialInsurance = *req.SocialInsurance
	}
	if req.HousingFund != nil {
		payroll.HousingFund = *req.HousingFund
	}
	if req.Tax != nil {
		payroll.Tax = *req.Tax
	}
	if req.NetPay != nil {
		payroll.NetPay = *req.NetPay
	}
	if req.Status != "" {
		payroll.Status = req.Status
	}

	// 保存更新
	if err := s.payrollRepo.Update(ctx, payroll); err != nil {
		return nil, fmt.Errorf("更新薪资记录失败: %w", err)
	}

	// 重新获取完整信息
	updated, err := s.payrollRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("获取更新的薪资记录失败: %w", err)
	}

	return s.toPayrollResponse(updated), nil
}

// Delete 删除薪资记录
func (s *PayrollServiceImpl) Delete(ctx context.Context, id uint) error {
	// 检查薪资记录是否存在
	payroll, err := s.payrollRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("获取薪资记录失败: %w", err)
	}
	if payroll == nil {
		return errors.New("薪资记录不存在")
	}

	return s.payrollRepo.Delete(ctx, id)
}

// List 获取薪资记录列表
func (s *PayrollServiceImpl) List(ctx context.Context, req *dto.PaginationRequest) (*dto.PaginatedResponse[dto.PayrollResponse], error) {
	options := &common.QueryOptions{
		Pagination: req,
	}
	payrolls, total, err := s.payrollRepo.List(ctx, options)
	if err != nil {
		return nil, fmt.Errorf("获取薪资记录列表失败: %w", err)
	}

	items := make([]dto.PayrollResponse, len(payrolls))
	for i, payroll := range payrolls {
		items[i] = *s.toPayrollResponse(payroll)
	}

	totalPages := int(total) / req.PageSize
	if int(total)%req.PageSize > 0 {
		totalPages++
	}

	return &dto.PaginatedResponse[dto.PayrollResponse]{
		Data:       items,
		Total:      total,
		Page:       req.Page,
		Limit:      req.PageSize,
		TotalPages: totalPages,
	}, nil
}

// GetByEmployeeID 根据员工ID获取薪资记录
func (s *PayrollServiceImpl) GetByEmployeeID(ctx context.Context, employeeID uint, req *dto.PaginationRequest) (*dto.PaginatedResponse[dto.PayrollResponse], error) {
	offset := (req.Page - 1) * req.PageSize
	payrolls, total, err := s.payrollRepo.GetByEmployeeID(ctx, employeeID, offset, req.PageSize)
	if err != nil {
		return nil, fmt.Errorf("获取员工薪资记录失败: %w", err)
	}

	items := make([]dto.PayrollResponse, len(payrolls))
	for i, payroll := range payrolls {
		items[i] = *s.toPayrollResponse(payroll)
	}

	totalPages := int(total) / req.PageSize
	if int(total)%req.PageSize > 0 {
		totalPages++
	}

	return &dto.PaginatedResponse[dto.PayrollResponse]{
		Data:       items,
		Total:      total,
		Page:       req.Page,
		Limit:      req.PageSize,
		TotalPages: totalPages,
	}, nil
}

// GetByPeriod 根据薪资期间获取薪资记录
func (s *PayrollServiceImpl) GetByPeriod(ctx context.Context, startDate, endDate time.Time, req *dto.PaginationRequest) (*dto.PaginatedResponse[dto.PayrollResponse], error) {
	offset := (req.Page - 1) * req.PageSize
	payrolls, total, err := s.payrollRepo.GetByPeriod(ctx, startDate, endDate, offset, req.PageSize)
	if err != nil {
		return nil, fmt.Errorf("获取期间薪资记录失败: %w", err)
	}

	items := make([]dto.PayrollResponse, len(payrolls))
	for i, payroll := range payrolls {
		items[i] = *s.toPayrollResponse(payroll)
	}

	totalPages := int(total) / req.PageSize
	if int(total)%req.PageSize > 0 {
		totalPages++
	}

	return &dto.PaginatedResponse[dto.PayrollResponse]{
		Data:       items,
		Total:      total,
		Page:       req.Page,
		Limit:      req.PageSize,
		TotalPages: totalPages,
	}, nil
}

// toPayrollResponse 转换为薪资响应
func (s *PayrollServiceImpl) toPayrollResponse(payroll *models.Payroll) *dto.PayrollResponse {
	resp := &dto.PayrollResponse{
		ID:              payroll.ID,
		EmployeeID:      payroll.EmployeeID,
		PayPeriodStart:  payroll.PayPeriodStart,
		PayPeriodEnd:    payroll.PayPeriodEnd,
		BasicSalary:     payroll.BasicSalary,
		OvertimePay:     payroll.OvertimePay,
		Allowance:       payroll.Allowance,
		Bonus:           payroll.Bonus,
		Deductions:      payroll.Deductions,
		SocialInsurance: payroll.SocialInsurance,
		HousingFund:     payroll.HousingFund,
		Tax:             payroll.Tax,
		NetPay:          payroll.NetPay,
		Status:          payroll.Status,
		PaidAt:          payroll.PaidAt,
		CreatedAt:       payroll.CreatedAt,
		UpdatedAt:       payroll.UpdatedAt,
	}

	// 添加员工信息
	if payroll.Employee.ID != 0 {
		resp.Employee = &dto.EmployeeListResponse{
			ID:        payroll.Employee.ID,
			Code:      payroll.Employee.Code,
			FirstName: payroll.Employee.FirstName,
			LastName:  payroll.Employee.LastName,
			FullName:  payroll.Employee.FullName,
		}
	}

	return resp
}

// ===== 请假管理 =====

// LeaveService 请假服务接口
type LeaveService interface {
	Create(ctx context.Context, req *dto.LeaveCreateRequest) (*dto.LeaveResponse, error)
	GetByID(ctx context.Context, id uint) (*dto.LeaveResponse, error)
	Update(ctx context.Context, id uint, req *dto.LeaveUpdateRequest) (*dto.LeaveResponse, error)
	Delete(ctx context.Context, id uint) error
	Approve(ctx context.Context, id uint, approverID uint, req *dto.LeaveApprovalRequest) (*dto.LeaveResponse, error)
	Cancel(ctx context.Context, id uint, employeeID uint) (*dto.LeaveResponse, error)
	List(ctx context.Context, req *dto.PaginationRequest) (*dto.PaginatedResponse[dto.LeaveListResponse], error)
	ListWithFilters(ctx context.Context, filter *dto.LeaveFilter) (*dto.PaginatedResponse[dto.LeaveListResponse], error)
	GetByEmployeeID(ctx context.Context, employeeID uint, req *dto.PaginationRequest) (*dto.PaginatedResponse[dto.LeaveListResponse], error)
	GetPendingApprovals(ctx context.Context, req *dto.PaginationRequest) (*dto.PaginatedResponse[dto.LeaveListResponse], error)
}

// LeaveServiceImpl 请假服务实现
type LeaveServiceImpl struct {
	leaveRepo    repositories.LeaveRepository
	employeeRepo repositories.EmployeeRepository
}

// NewLeaveService 创建请假服务
func NewLeaveService(leaveRepo repositories.LeaveRepository, employeeRepo repositories.EmployeeRepository) LeaveService {
	return &LeaveServiceImpl{
		leaveRepo:    leaveRepo,
		employeeRepo: employeeRepo,
	}
}

// Create 创建请假申请
func (s *LeaveServiceImpl) Create(ctx context.Context, req *dto.LeaveCreateRequest) (*dto.LeaveResponse, error) {
	// 验证员工是否存在
	employee, err := s.employeeRepo.GetByID(ctx, req.EmployeeID)
	if err != nil {
		return nil, fmt.Errorf("获取员工信息失败: %w", err)
	}
	if employee == nil {
		return nil, errors.New("员工不存在")
	}

	// 验证日期逻辑
	if req.EndDate.Before(req.StartDate) {
		return nil, errors.New("结束日期不能早于开始日期")
	}

	// 验证请假天数
	expectedDays := req.EndDate.Sub(req.StartDate).Hours()/24 + 1
	if req.Days > expectedDays {
		return nil, errors.New("请假天数不能超过日期范围")
	}

	// 检查是否有重叠的请假申请
	overlappingLeaves, err := s.leaveRepo.GetByEmployeeIDAndDateRange(ctx, req.EmployeeID, req.StartDate, req.EndDate)
	if err != nil {
		return nil, fmt.Errorf("检查重叠请假失败: %w", err)
	}

	for _, leave := range overlappingLeaves {
		if leave.Status == "approved" || leave.Status == "pending" {
			return nil, errors.New("该时间段已有请假申请")
		}
	}

	// 创建请假申请
	leave := &models.Leave{
		EmployeeID: req.EmployeeID,
		LeaveType:  req.LeaveType,
		StartDate:  req.StartDate,
		EndDate:    req.EndDate,
		Days:       req.Days,
		Reason:     req.Reason,
		Status:     "pending",
	}

	if err := s.leaveRepo.Create(ctx, leave); err != nil {
		return nil, fmt.Errorf("创建请假申请失败: %w", err)
	}

	// 获取完整信息并返回
	createdLeave, err := s.leaveRepo.GetByID(ctx, leave.ID)
	if err != nil {
		return nil, fmt.Errorf("获取创建的请假申请失败: %w", err)
	}

	return s.toLeaveResponse(createdLeave), nil
}

// GetByID 根据ID获取请假申请
func (s *LeaveServiceImpl) GetByID(ctx context.Context, id uint) (*dto.LeaveResponse, error) {
	leave, err := s.leaveRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("获取请假申请失败: %w", err)
	}
	if leave == nil {
		return nil, errors.New("请假申请不存在")
	}

	return s.toLeaveResponse(leave), nil
}

// Update 更新请假申请
func (s *LeaveServiceImpl) Update(ctx context.Context, id uint, req *dto.LeaveUpdateRequest) (*dto.LeaveResponse, error) {
	leave, err := s.leaveRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("获取请假申请失败: %w", err)
	}
	if leave == nil {
		return nil, errors.New("请假申请不存在")
	}

	// 只有待审批状态的申请才能修改
	if leave.Status != "pending" {
		return nil, errors.New("只有待审批的请假申请才能修改")
	}

	// 更新字段
	if req.LeaveType != nil {
		leave.LeaveType = *req.LeaveType
	}
	if req.StartDate != nil {
		leave.StartDate = *req.StartDate
	}
	if req.EndDate != nil {
		leave.EndDate = *req.EndDate
	}
	if req.Days != nil {
		leave.Days = *req.Days
	}
	if req.Reason != nil {
		leave.Reason = *req.Reason
	}

	// 验证更新后的日期逻辑
	if leave.EndDate.Before(leave.StartDate) {
		return nil, errors.New("结束日期不能早于开始日期")
	}

	// 验证请假天数
	expectedDays := leave.EndDate.Sub(leave.StartDate).Hours()/24 + 1
	if leave.Days > expectedDays {
		return nil, errors.New("请假天数不能超过日期范围")
	}

	if err := s.leaveRepo.Update(ctx, leave); err != nil {
		return nil, fmt.Errorf("更新请假申请失败: %w", err)
	}

	return s.toLeaveResponse(leave), nil
}

// Delete 删除请假申请
func (s *LeaveServiceImpl) Delete(ctx context.Context, id uint) error {
	leave, err := s.leaveRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("获取请假申请失败: %w", err)
	}
	if leave == nil {
		return errors.New("请假申请不存在")
	}

	// 只有待审批状态的申请才能删除
	if leave.Status != "pending" {
		return errors.New("只有待审批的请假申请才能删除")
	}

	return s.leaveRepo.Delete(ctx, id)
}

// Approve 审批请假申请
func (s *LeaveServiceImpl) Approve(ctx context.Context, id uint, approverID uint, req *dto.LeaveApprovalRequest) (*dto.LeaveResponse, error) {
	leave, err := s.leaveRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("获取请假申请失败: %w", err)
	}
	if leave == nil {
		return nil, errors.New("请假申请不存在")
	}

	// 只有待审批状态的申请才能审批
	if leave.Status != "pending" {
		return nil, errors.New("只有待审批的请假申请才能审批")
	}

	// 验证审批人是否存在
	approver, err := s.employeeRepo.GetByID(ctx, approverID)
	if err != nil {
		return nil, fmt.Errorf("获取审批人信息失败: %w", err)
	}
	if approver == nil {
		return nil, errors.New("审批人不存在")
	}

	// 更新审批信息
	now := time.Now()
	leave.Status = req.Status
	leave.ApprovedBy = &approverID
	leave.ApprovedAt = &now

	if err := s.leaveRepo.Update(ctx, leave); err != nil {
		return nil, fmt.Errorf("更新请假申请失败: %w", err)
	}

	return s.toLeaveResponse(leave), nil
}

// Cancel 取消请假申请
func (s *LeaveServiceImpl) Cancel(ctx context.Context, id uint, employeeID uint) (*dto.LeaveResponse, error) {
	leave, err := s.leaveRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("获取请假申请失败: %w", err)
	}
	if leave == nil {
		return nil, errors.New("请假申请不存在")
	}

	// 验证是否为申请人本人
	if leave.EmployeeID != employeeID {
		return nil, errors.New("只能取消自己的请假申请")
	}

	// 只有待审批和已批准的申请才能取消
	if leave.Status != "pending" && leave.Status != "approved" {
		return nil, errors.New("只有待审批或已批准的请假申请才能取消")
	}

	// 更新状态
	leave.Status = "cancelled"

	if err := s.leaveRepo.Update(ctx, leave); err != nil {
		return nil, fmt.Errorf("取消请假申请失败: %w", err)
	}

	return s.toLeaveResponse(leave), nil
}

// List 获取请假申请列表
func (s *LeaveServiceImpl) List(ctx context.Context, req *dto.PaginationRequest) (*dto.PaginatedResponse[dto.LeaveListResponse], error) {
	options := &common.QueryOptions{
		Pagination: req,
	}
	leaves, total, err := s.leaveRepo.List(ctx, options)
	if err != nil {
		return nil, fmt.Errorf("获取请假申请列表失败: %w", err)
	}

	items := make([]dto.LeaveListResponse, len(leaves))
	for i, leave := range leaves {
		items[i] = s.toLeaveListResponse(leave)
	}

	return &dto.PaginatedResponse[dto.LeaveListResponse]{
		Data:       items,
		Total:      total,
		Page:       req.Page,
		Limit:      req.PageSize,
		TotalPages: (int(total) + req.PageSize - 1) / req.PageSize,
	}, nil
}

// ListWithFilters 根据筛选条件获取请假申请列表
func (s *LeaveServiceImpl) ListWithFilters(ctx context.Context, filter *dto.LeaveFilter) (*dto.PaginatedResponse[dto.LeaveListResponse], error) {
	offset := (filter.Page - 1) * filter.PageSize
	leaves, total, err := s.leaveRepo.ListWithFilters(ctx, filter.EmployeeID, filter.LeaveType, filter.Status, filter.StartDate, filter.EndDate, offset, filter.PageSize)
	if err != nil {
		return nil, fmt.Errorf("获取请假申请列表失败: %w", err)
	}

	items := make([]dto.LeaveListResponse, len(leaves))
	for i, leave := range leaves {
		items[i] = s.toLeaveListResponse(leave)
	}

	return &dto.PaginatedResponse[dto.LeaveListResponse]{
		Data:       items,
		Total:      total,
		Page:       filter.Page,
		Limit:      filter.PageSize,
		TotalPages: (int(total) + filter.PageSize - 1) / filter.PageSize,
	}, nil
}

// GetByEmployeeID 根据员工ID获取请假申请列表
func (s *LeaveServiceImpl) GetByEmployeeID(ctx context.Context, employeeID uint, req *dto.PaginationRequest) (*dto.PaginatedResponse[dto.LeaveListResponse], error) {
	offset := (req.Page - 1) * req.PageSize
	leaves, total, err := s.leaveRepo.GetByEmployeeID(ctx, employeeID, offset, req.PageSize)
	if err != nil {
		return nil, fmt.Errorf("获取员工请假申请列表失败: %w", err)
	}

	items := make([]dto.LeaveListResponse, len(leaves))
	for i, leave := range leaves {
		items[i] = s.toLeaveListResponse(leave)
	}

	return &dto.PaginatedResponse[dto.LeaveListResponse]{
		Data:       items,
		Total:      total,
		Page:       req.Page,
		Limit:      req.PageSize,
		TotalPages: (int(total) + req.PageSize - 1) / req.PageSize,
	}, nil
}

// GetPendingApprovals 获取待审批的请假申请列表
func (s *LeaveServiceImpl) GetPendingApprovals(ctx context.Context, req *dto.PaginationRequest) (*dto.PaginatedResponse[dto.LeaveListResponse], error) {
	offset := (req.Page - 1) * req.PageSize
	leaves, total, err := s.leaveRepo.GetByStatus(ctx, "pending", offset, req.PageSize)
	if err != nil {
		return nil, fmt.Errorf("获取待审批请假申请列表失败: %w", err)
	}

	items := make([]dto.LeaveListResponse, len(leaves))
	for i, leave := range leaves {
		items[i] = s.toLeaveListResponse(leave)
	}

	return &dto.PaginatedResponse[dto.LeaveListResponse]{
		Data:       items,
		Total:      total,
		Page:       req.Page,
		Limit:      req.PageSize,
		TotalPages: (int(total) + req.PageSize - 1) / req.PageSize,
	}, nil
}

// toLeaveResponse 转换为请假申请响应
func (s *LeaveServiceImpl) toLeaveResponse(leave *models.Leave) *dto.LeaveResponse {
	response := &dto.LeaveResponse{
		ID:         leave.ID,
		EmployeeID: leave.EmployeeID,
		LeaveType:  leave.LeaveType,
		StartDate:  leave.StartDate,
		EndDate:    leave.EndDate,
		Days:       leave.Days,
		Reason:     leave.Reason,
		Status:     leave.Status,
		ApprovedBy: leave.ApprovedBy,
		ApprovedAt: leave.ApprovedAt,
		CreatedAt:  leave.CreatedAt,
		UpdatedAt:  leave.UpdatedAt,
	}

	// Employee是直接的Employee类型，不是指针，检查ID是否为0来判断是否已加载
	if leave.Employee.ID != 0 {
		response.Employee = &dto.EmployeeListResponse{
			ID:        leave.Employee.ID,
			Code:      leave.Employee.Code,
			FirstName: leave.Employee.FirstName,
			LastName:  leave.Employee.LastName,
			FullName:  leave.Employee.FullName,
			Email:     leave.Employee.Email,
			Phone:     leave.Employee.Phone,
			Status:    leave.Employee.Status,
		}
		if leave.Employee.Department != nil {
			response.Employee.Department = leave.Employee.Department.Name
		}
		if leave.Employee.Position != nil {
			response.Employee.Position = leave.Employee.Position.Name
		}
	}

	if leave.ApprovedUser != nil {
		response.ApprovedUser = &dto.EmployeeListResponse{
			ID:        leave.ApprovedUser.ID,
			Code:      leave.ApprovedUser.Code,
			FirstName: leave.ApprovedUser.FirstName,
			LastName:  leave.ApprovedUser.LastName,
			FullName:  leave.ApprovedUser.FullName,
			Email:     leave.ApprovedUser.Email,
			Phone:     leave.ApprovedUser.Phone,
			Status:    leave.ApprovedUser.Status,
		}
	}

	return response
}

// toLeaveListResponse 转换为请假申请列表响应
func (s *LeaveServiceImpl) toLeaveListResponse(leave *models.Leave) dto.LeaveListResponse {
	response := dto.LeaveListResponse{
		ID:         leave.ID,
		EmployeeID: leave.EmployeeID,
		LeaveType:  leave.LeaveType,
		StartDate:  leave.StartDate,
		EndDate:    leave.EndDate,
		Days:       leave.Days,
		Status:     leave.Status,
		CreatedAt:  leave.CreatedAt,
	}

	if leave.Employee.ID != 0 {
		response.EmployeeName = leave.Employee.FullName
	}

	return response
}
