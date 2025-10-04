package dto

import (
	"time"
)

// EmployeeCreateRequest 员工创建请求
type EmployeeCreateRequest struct {
	Code             string     `json:"code" validate:"required,max=50"`
	FirstName        string     `json:"first_name" validate:"required,max=50"`
	LastName         string     `json:"last_name" validate:"required,max=50"`
	Email            string     `json:"email" validate:"required,email"`
	Phone            string     `json:"phone,omitempty" validate:"omitempty,max=20"`
	DateOfBirth      *time.Time `json:"date_of_birth,omitempty"`
	Gender           string     `json:"gender,omitempty" validate:"omitempty,oneof=male female other"`
	HireDate         time.Time  `json:"hire_date" validate:"required"`
	DepartmentID     *uint      `json:"department_id,omitempty"`
	PositionID       *uint      `json:"position_id,omitempty"`
	ManagerID        *uint      `json:"manager_id,omitempty"`
	Status           string     `json:"status" validate:"required,oneof=active inactive terminated"`
	EmergencyContact string     `json:"emergency_contact,omitempty"`
	IDNumber         string     `json:"id_number,omitempty" validate:"omitempty,max=50"`
	Address          string     `json:"address,omitempty"`
	BankAccount      string     `json:"bank_account,omitempty"`
}

// EmployeeUpdateRequest 员工更新请求
type EmployeeUpdateRequest struct {
	FirstName        string     `json:"first_name,omitempty" validate:"omitempty,max=50"`
	LastName         string     `json:"last_name,omitempty" validate:"omitempty,max=50"`
	Email            string     `json:"email,omitempty" validate:"omitempty,email"`
	Phone            string     `json:"phone,omitempty" validate:"omitempty,max=20"`
	DateOfBirth      *time.Time `json:"date_of_birth,omitempty"`
	Gender           string     `json:"gender,omitempty" validate:"omitempty,oneof=male female other"`
	DepartmentID     *uint      `json:"department_id,omitempty"`
	PositionID       *uint      `json:"position_id,omitempty"`
	ManagerID        *uint      `json:"manager_id,omitempty"`
	Status           string     `json:"status,omitempty" validate:"omitempty,oneof=active inactive terminated"`
	EmergencyContact string     `json:"emergency_contact,omitempty"`
	IDNumber         string     `json:"id_number,omitempty" validate:"omitempty,max=50"`
	Address          string     `json:"address,omitempty"`
	BankAccount      string     `json:"bank_account,omitempty"`
}

// EmployeeResponse 员工响应
type EmployeeResponse struct {
	ID               uint                `json:"id"`
	Code             string              `json:"code"`
	FirstName        string              `json:"first_name"`
	LastName         string              `json:"last_name"`
	FullName         string              `json:"full_name"`
	Email            string              `json:"email"`
	Phone            string              `json:"phone,omitempty"`
	DateOfBirth      *time.Time          `json:"date_of_birth,omitempty"`
	Gender           string              `json:"gender,omitempty"`
	HireDate         time.Time           `json:"hire_date"`
	DepartmentID     *uint               `json:"department_id,omitempty"`
	PositionID       *uint               `json:"position_id,omitempty"`
	ManagerID        *uint               `json:"manager_id,omitempty"`
	Status           string              `json:"status"`
	EmergencyContact string              `json:"emergency_contact,omitempty"`
	IDNumber         string              `json:"id_number,omitempty"`
	Address          string              `json:"address,omitempty"`
	BankAccount      string              `json:"bank_account,omitempty"`
	Department       *DepartmentResponse `json:"department,omitempty"`
	Position         *PositionResponse   `json:"position,omitempty"`
	Manager          *EmployeeResponse   `json:"manager,omitempty"`
	CreatedAt        time.Time           `json:"created_at"`
	UpdatedAt        time.Time           `json:"updated_at"`
}

// EmployeeListResponse 员工列表响应
type EmployeeListResponse struct {
	ID         uint   `json:"id"`
	Code       string `json:"code"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	FullName   string `json:"full_name"`
	Email      string `json:"email"`
	Phone      string `json:"phone,omitempty"`
	Department string `json:"department,omitempty"`
	Position   string `json:"position,omitempty"`
	Status     string `json:"status"`
}

// PositionCreateRequest 职位创建请求
type PositionCreateRequest struct {
	Name         string  `json:"name" validate:"required,max=100"`
	Code         string  `json:"code" validate:"required,max=50"`
	Description  string  `json:"description,omitempty"`
	DepartmentID uint    `json:"department_id" validate:"required"`
	Level        int     `json:"level" validate:"required,min=1"`
	MinSalary    float64 `json:"min_salary,omitempty" validate:"omitempty,min=0"`
	MaxSalary    float64 `json:"max_salary,omitempty" validate:"omitempty,min=0"`
}

// PositionUpdateRequest 职位更新请求
type PositionUpdateRequest struct {
	Name         string   `json:"name,omitempty" validate:"omitempty,max=100"`
	Description  string   `json:"description,omitempty"`
	DepartmentID *uint    `json:"department_id,omitempty"`
	Level        *int     `json:"level,omitempty" validate:"omitempty,min=1"`
	MinSalary    *float64 `json:"min_salary,omitempty" validate:"omitempty,min=0"`
	MaxSalary    *float64 `json:"max_salary,omitempty" validate:"omitempty,min=0"`
	IsActive     *bool    `json:"is_active,omitempty"`
}

// PositionResponse 职位响应
type PositionResponse struct {
	ID          uint               `json:"id"`
	Name        string             `json:"name"`
	Code        string             `json:"code"`
	Description string             `json:"description,omitempty"`
	Level       int                `json:"level"`
	MinSalary   float64            `json:"min_salary"`
	MaxSalary   float64            `json:"max_salary"`
	IsActive    bool               `json:"is_active"`
	Department  DepartmentResponse `json:"department"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
}

// EmployeeSearchRequest 员工搜索请求
type EmployeeSearchRequest struct {
	SearchRequest
	DepartmentID *uint  `json:"department_id,omitempty" form:"department_id"`
	PositionID   *uint  `json:"position_id,omitempty" form:"position_id"`
	Status       string `json:"status,omitempty" form:"status"`
	Gender       string `json:"gender,omitempty" form:"gender"`
}

// EmployeeFilter 员工筛选
type EmployeeFilter struct {
	PaginationRequest
	DepartmentID *uint  `json:"department_id,omitempty" form:"department_id"`
	PositionID   *uint  `json:"position_id,omitempty" form:"position_id"`
	Status       string `json:"status,omitempty" form:"status"`
	Name         string `json:"name,omitempty" form:"name"`
	Code         string `json:"code,omitempty" form:"code"`
	Email        string `json:"email,omitempty" form:"email"`
}

// AttendanceCreateRequest 考勤创建请求
type AttendanceCreateRequest struct {
	EmployeeID   uint       `json:"employee_id" validate:"required"`
	Date         time.Time  `json:"date" validate:"required"`
	CheckInTime  *time.Time `json:"check_in_time,omitempty"`
	CheckOutTime *time.Time `json:"check_out_time,omitempty"`
	Status       string     `json:"status" validate:"required,oneof=present absent late early_leave"`
	Notes        string     `json:"notes,omitempty"`
}

// AttendanceUpdateRequest 考勤更新请求
type AttendanceUpdateRequest struct {
	CheckInTime  *time.Time `json:"check_in_time,omitempty"`
	CheckOutTime *time.Time `json:"check_out_time,omitempty"`
	Status       string     `json:"status,omitempty" validate:"omitempty,oneof=present absent late early_leave"`
	Notes        string     `json:"notes,omitempty"`
}

// AttendanceResponse 考勤响应
type AttendanceResponse struct {
	ID           uint                  `json:"id"`
	EmployeeID   uint                  `json:"employee_id"`
	Date         time.Time             `json:"date"`
	CheckInTime  *time.Time            `json:"check_in_time,omitempty"`
	CheckOutTime *time.Time            `json:"check_out_time,omitempty"`
	Status       string                `json:"status"`
	Notes        string                `json:"notes,omitempty"`
	Employee     *EmployeeListResponse `json:"employee,omitempty"`
	CreatedAt    time.Time             `json:"created_at"`
	UpdatedAt    time.Time             `json:"updated_at"`
}

// PayrollCreateRequest 工资单创建请求
type PayrollCreateRequest struct {
	EmployeeID      uint      `json:"employee_id" validate:"required"`
	PayPeriodStart  time.Time `json:"pay_period_start" validate:"required"`
	PayPeriodEnd    time.Time `json:"pay_period_end" validate:"required"`
	BasicSalary     float64   `json:"basic_salary" validate:"required,min=0"`
	OvertimePay     float64   `json:"overtime_pay,omitempty" validate:"omitempty,min=0"`
	Allowance       float64   `json:"allowance,omitempty" validate:"omitempty,min=0"`
	Bonus           float64   `json:"bonus,omitempty" validate:"omitempty,min=0"`
	Deductions      float64   `json:"deductions,omitempty" validate:"omitempty,min=0"`
	SocialInsurance float64   `json:"social_insurance,omitempty" validate:"omitempty,min=0"`
	HousingFund     float64   `json:"housing_fund,omitempty" validate:"omitempty,min=0"`
	Tax             float64   `json:"tax,omitempty" validate:"omitempty,min=0"`
	NetPay          float64   `json:"net_pay" validate:"required,min=0"`
	Status          string    `json:"status" validate:"required,oneof=draft confirmed paid"`
}

// PayrollUpdateRequest 工资单更新请求
type PayrollUpdateRequest struct {
	BasicSalary     *float64 `json:"basic_salary,omitempty" validate:"omitempty,min=0"`
	OvertimePay     *float64 `json:"overtime_pay,omitempty" validate:"omitempty,min=0"`
	Allowance       *float64 `json:"allowance,omitempty" validate:"omitempty,min=0"`
	Bonus           *float64 `json:"bonus,omitempty" validate:"omitempty,min=0"`
	Deductions      *float64 `json:"deductions,omitempty" validate:"omitempty,min=0"`
	SocialInsurance *float64 `json:"social_insurance,omitempty" validate:"omitempty,min=0"`
	HousingFund     *float64 `json:"housing_fund,omitempty" validate:"omitempty,min=0"`
	Tax             *float64 `json:"tax,omitempty" validate:"omitempty,min=0"`
	NetPay          *float64 `json:"net_pay,omitempty" validate:"omitempty,min=0"`
	Status          string   `json:"status,omitempty" validate:"omitempty,oneof=draft confirmed paid"`
}

// PayrollResponse 薪资响应
type PayrollResponse struct {
	ID              uint                  `json:"id"`
	EmployeeID      uint                  `json:"employee_id"`
	PayPeriodStart  time.Time             `json:"pay_period_start"`
	PayPeriodEnd    time.Time             `json:"pay_period_end"`
	BasicSalary     float64               `json:"basic_salary"`
	OvertimePay     float64               `json:"overtime_pay"`
	Allowance       float64               `json:"allowance"`
	Bonus           float64               `json:"bonus"`
	Deductions      float64               `json:"deductions"`
	SocialInsurance float64               `json:"social_insurance"`
	HousingFund     float64               `json:"housing_fund"`
	Tax             float64               `json:"tax"`
	NetPay          float64               `json:"net_pay"`
	Status          string                `json:"status"`
	PaidAt          *time.Time            `json:"paid_at,omitempty"`
	Employee        *EmployeeListResponse `json:"employee,omitempty"`
	CreatedAt       time.Time             `json:"created_at"`
	UpdatedAt       time.Time             `json:"updated_at"`
}

// ===== 请假管理 =====

// LeaveCreateRequest 请假创建请求
type LeaveCreateRequest struct {
	EmployeeID uint      `json:"employee_id" validate:"required"`
	LeaveType  string    `json:"leave_type" validate:"required,oneof=annual sick personal maternity paternity emergency"`
	StartDate  time.Time `json:"start_date" validate:"required"`
	EndDate    time.Time `json:"end_date" validate:"required"`
	Days       float64   `json:"days" validate:"required,min=0.5"`
	Reason     string    `json:"reason" validate:"required,max=500"`
}

// LeaveUpdateRequest 请假更新请求
type LeaveUpdateRequest struct {
	LeaveType *string    `json:"leave_type,omitempty" validate:"omitempty,oneof=annual sick personal maternity paternity emergency"`
	StartDate *time.Time `json:"start_date,omitempty"`
	EndDate   *time.Time `json:"end_date,omitempty"`
	Days      *float64   `json:"days,omitempty" validate:"omitempty,min=0.5"`
	Reason    *string    `json:"reason,omitempty" validate:"omitempty,max=500"`
}

// LeaveApprovalRequest 请假审批请求
type LeaveApprovalRequest struct {
	Status   string `json:"status" validate:"required,oneof=approved rejected"`
	Comments string `json:"comments,omitempty" validate:"omitempty,max=500"`
}

// LeaveResponse 请假申请响应
type LeaveResponse struct {
	ID           uint                  `json:"id"`
	EmployeeID   uint                  `json:"employee_id"`
	LeaveType    string                `json:"leave_type"`
	StartDate    time.Time             `json:"start_date"`
	EndDate      time.Time             `json:"end_date"`
	Days         float64               `json:"days"`
	Reason       string                `json:"reason"`
	Status       string                `json:"status"`
	ApprovedBy   *uint                 `json:"approved_by,omitempty"`
	ApprovedAt   *time.Time            `json:"approved_at,omitempty"`
	Employee     *EmployeeListResponse `json:"employee,omitempty"`
	ApprovedUser *EmployeeListResponse `json:"approved_user,omitempty"`
	CreatedAt    time.Time             `json:"created_at"`
	UpdatedAt    time.Time             `json:"updated_at"`
}

// LeaveListResponse 请假申请列表响应
type LeaveListResponse struct {
	ID           uint      `json:"id"`
	EmployeeID   uint      `json:"employee_id"`
	EmployeeName string    `json:"employee_name"`
	LeaveType    string    `json:"leave_type"`
	StartDate    time.Time `json:"start_date"`
	EndDate      time.Time `json:"end_date"`
	Days         float64   `json:"days"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
}

// LeaveFilter 请假申请筛选条件
type LeaveFilter struct {
	PaginationRequest
	EmployeeID *uint      `json:"employee_id,omitempty" form:"employee_id"`
	LeaveType  string     `json:"leave_type,omitempty" form:"leave_type"`
	Status     string     `json:"status,omitempty" form:"status"`
	StartDate  *time.Time `json:"start_date,omitempty" form:"start_date"`
	EndDate    *time.Time `json:"end_date,omitempty" form:"end_date"`
}
