package dto

import (
	"time"
)

// EmployeeCreateRequest 员工创建请求
type EmployeeCreateRequest struct {
	Code             string     `json:"code" binding:"required,max=50"`
	FirstName        string     `json:"first_name" binding:"required,max=50"`
	LastName         string     `json:"last_name" binding:"required,max=50"`
	Email            string     `json:"email" binding:"required,email"`
	Phone            string     `json:"phone,omitempty" binding:"omitempty,max=20"`
	DateOfBirth      *time.Time `json:"date_of_birth,omitempty"`
	Gender           string     `json:"gender,omitempty" binding:"omitempty,oneof=male female other"`
	HireDate         time.Time  `json:"hire_date" binding:"required"`
	DepartmentID     *uint      `json:"department_id,omitempty"`
	PositionID       *uint      `json:"position_id,omitempty"`
	ManagerID        *uint      `json:"manager_id,omitempty"`
	Status           string     `json:"status" binding:"required,oneof=active inactive terminated"`
	EmergencyContact string     `json:"emergency_contact,omitempty"`
	IDNumber         string     `json:"id_number,omitempty" binding:"omitempty,max=50"`
	Address          string     `json:"address,omitempty"`
	BankAccount      string     `json:"bank_account,omitempty"`
}

// EmployeeUpdateRequest 员工更新请求
type EmployeeUpdateRequest struct {
	FirstName        string     `json:"first_name,omitempty" binding:"omitempty,max=50"`
	LastName         string     `json:"last_name,omitempty" binding:"omitempty,max=50"`
	Email            string     `json:"email,omitempty" binding:"omitempty,email"`
	Phone            string     `json:"phone,omitempty" binding:"omitempty,max=20"`
	DateOfBirth      *time.Time `json:"date_of_birth,omitempty"`
	Gender           string     `json:"gender,omitempty" binding:"omitempty,oneof=male female other"`
	DepartmentID     *uint      `json:"department_id,omitempty"`
	PositionID       *uint      `json:"position_id,omitempty"`
	ManagerID        *uint      `json:"manager_id,omitempty"`
	Status           string     `json:"status,omitempty" binding:"omitempty,oneof=active inactive terminated"`
	EmergencyContact string     `json:"emergency_contact,omitempty"`
	IDNumber         string     `json:"id_number,omitempty" binding:"omitempty,max=50"`
	Address          string     `json:"address,omitempty"`
	BankAccount      string     `json:"bank_account,omitempty"`
}

// EmployeeResponse 员工响应
type EmployeeResponse struct {
	ID               uint                 `json:"id"`
	Code             string               `json:"code"`
	FirstName        string               `json:"first_name"`
	LastName         string               `json:"last_name"`
	FullName         string               `json:"full_name"`
	Email            string               `json:"email"`
	Phone            string               `json:"phone,omitempty"`
	DateOfBirth      *time.Time           `json:"date_of_birth,omitempty"`
	Gender           string               `json:"gender,omitempty"`
	HireDate         time.Time            `json:"hire_date"`
	DepartmentID     *uint                `json:"department_id,omitempty"`
	PositionID       *uint                `json:"position_id,omitempty"`
	ManagerID        *uint                `json:"manager_id,omitempty"`
	Status           string               `json:"status"`
	EmergencyContact string               `json:"emergency_contact,omitempty"`
	IDNumber         string               `json:"id_number,omitempty"`
	Address          string               `json:"address,omitempty"`
	BankAccount      string               `json:"bank_account,omitempty"`
	Department       *DepartmentResponse  `json:"department,omitempty"`
	Position         *PositionResponse    `json:"position,omitempty"`
	Manager          *EmployeeResponse    `json:"manager,omitempty"`
	CreatedAt        time.Time            `json:"created_at"`
	UpdatedAt        time.Time            `json:"updated_at"`
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
	Name         string `json:"name" binding:"required,max=100"`
	Code         string `json:"code" binding:"required,max=50"`
	Description  string `json:"description,omitempty"`
	DepartmentID uint   `json:"department_id" binding:"required"`
	Level        int    `json:"level" binding:"required,min=1"`
	MinSalary    float64 `json:"min_salary,omitempty" binding:"omitempty,min=0"`
	MaxSalary    float64 `json:"max_salary,omitempty" binding:"omitempty,min=0"`
}

// PositionUpdateRequest 职位更新请求
type PositionUpdateRequest struct {
	Name         string   `json:"name,omitempty" binding:"omitempty,max=100"`
	Description  string   `json:"description,omitempty"`
	DepartmentID *uint    `json:"department_id,omitempty"`
	Level        *int     `json:"level,omitempty" binding:"omitempty,min=1"`
	MinSalary    *float64 `json:"min_salary,omitempty" binding:"omitempty,min=0"`
	MaxSalary    *float64 `json:"max_salary,omitempty" binding:"omitempty,min=0"`
	IsActive     *bool    `json:"is_active,omitempty"`
}

// PositionResponse 职位响应
type PositionResponse struct {
	ID           uint                `json:"id"`
	Name         string              `json:"name"`
	Code         string              `json:"code"`
	Description  string              `json:"description,omitempty"`
	Level        int                 `json:"level"`
	MinSalary    float64             `json:"min_salary"`
	MaxSalary    float64             `json:"max_salary"`
	IsActive     bool                `json:"is_active"`
	Department   DepartmentResponse  `json:"department"`
	CreatedAt    time.Time           `json:"created_at"`
	UpdatedAt    time.Time           `json:"updated_at"`
}

// EmployeeSearchRequest 员工搜索请求
type EmployeeSearchRequest struct {
	SearchRequest
	DepartmentID *uint  `json:"department_id,omitempty" form:"department_id"`
	PositionID   *uint  `json:"position_id,omitempty" form:"position_id"`
	Status       string `json:"status,omitempty" form:"status"`
	Gender       string `json:"gender,omitempty" form:"gender"`
}

// EmployeeFilter 员工过滤器
type EmployeeFilter struct {
	PaginationRequest
	DepartmentID *uint  `json:"department_id,omitempty" form:"department_id"`
	PositionID   *uint  `json:"position_id,omitempty" form:"position_id"`
	Status       string `json:"status,omitempty" form:"status"`
	Name         string `json:"name,omitempty" form:"name"`
	Code         string `json:"code,omitempty" form:"code"`
	Email        string `json:"email,omitempty" form:"email"`
}