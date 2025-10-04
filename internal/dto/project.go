package dto

import (
	"time"
)

// ProjectCreateRequest 项目创建请求
type ProjectCreateRequest struct {
	ProjectCode string    `json:"project_code,omitempty"`
	ProjectName string    `json:"project_name" validate:"required"`
	Name        string    `json:"name" validate:"required"`
	Description string    `json:"description,omitempty"`
	StartDate   time.Time `json:"start_date" validate:"required"`
	EndDate     time.Time `json:"end_date" validate:"required"`
	Status      string    `json:"status" validate:"required,oneof=Planning Active Completed Cancelled"`
	Priority    string    `json:"priority" validate:"required,oneof=Low Medium High"`
	Budget      float64   `json:"budget" validate:"min=0"`
	ClientID    *uint     `json:"client_id,omitempty"`
	ManagerID   uint      `json:"manager_id" validate:"required"`
}

// ProjectUpdateRequest 项目更新请求
type ProjectUpdateRequest struct {
	ProjectName string     `json:"project_name,omitempty"`
	Name        string     `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	Description string     `json:"description,omitempty"`
	StartDate   *time.Time `json:"start_date,omitempty"`
	EndDate     *time.Time `json:"end_date,omitempty"`
	Status      string     `json:"status,omitempty" validate:"omitempty,oneof=Planning Active Completed Cancelled"`
	Priority    string     `json:"priority,omitempty" validate:"omitempty,oneof=Low Medium High Critical"`
	Budget      *float64   `json:"budget,omitempty" validate:"omitempty,min=0"`
	ClientID    *uint      `json:"client_id,omitempty"`
	ManagerID   *uint      `json:"manager_id,omitempty"`
}

// TaskCreateRequest 任务创建请求
type TaskCreateRequest struct {
	ProjectID       uint       `json:"project_id" validate:"required"`
	TaskNumber      string     `json:"task_number,omitempty"`
	TaskName        string     `json:"task_name" validate:"required"`
	Name            string     `json:"name" validate:"required"`
	Description     string     `json:"description,omitempty"`
	StartDate       *time.Time `json:"start_date,omitempty"`
	EndDate         *time.Time `json:"end_date,omitempty"`
	Status          string     `json:"status" validate:"required,oneof=Todo InProgress Completed"`
	Priority        string     `json:"priority" validate:"required,oneof=Low Medium High"`
	AssigneeID      *uint      `json:"assignee_id,omitempty"`
	ParentTaskID    *uint      `json:"parent_task_id,omitempty"`
	EstimatedHours  float64    `json:"estimated_hours,omitempty"`
}

// TaskUpdateRequest 任务更新请求
type TaskUpdateRequest struct {
	TaskName        string     `json:"task_name,omitempty"`
	Name            string     `json:"name,omitempty"`
	Description     string     `json:"description,omitempty"`
	StartDate       *time.Time `json:"start_date,omitempty"`
	EndDate         *time.Time `json:"end_date,omitempty"`
	Status          string     `json:"status,omitempty"`
	Priority        string     `json:"priority,omitempty"`
	AssigneeID      *uint      `json:"assignee_id,omitempty"`
	ParentTaskID    *uint      `json:"parent_task_id,omitempty"`
	EstimatedHours  *float64   `json:"estimated_hours,omitempty"`
	ActualHours     *float64   `json:"actual_hours,omitempty"`
	Progress        *float64   `json:"progress,omitempty"`
}

// TimeEntryCreateRequest 时间记录创建请求
type TimeEntryCreateRequest struct {
	EmployeeID  uint      `json:"employee_id" validate:"required"`
	ProjectID   uint      `json:"project_id" validate:"required"`
	TaskID      *uint     `json:"task_id,omitempty"`
	Date        time.Time `json:"date" validate:"required"`
	StartTime   time.Time `json:"start_time" validate:"required"`
	EndTime     time.Time `json:"end_time" validate:"required"`
	Hours       float64   `json:"hours" validate:"required,gt=0,lte=24"`
	Description string    `json:"description,omitempty"`
	IsBillable  bool      `json:"is_billable"`
	HourlyRate  float64   `json:"hourly_rate,omitempty"`
	Status      string    `json:"status" validate:"required,oneof=draft submitted approved billed"`
}

// TimeEntryUpdateRequest 时间记录更新请求
type TimeEntryUpdateRequest struct {
	EmployeeID  *uint      `json:"employee_id,omitempty"`
	ProjectID   *uint      `json:"project_id,omitempty"`
	TaskID      *uint      `json:"task_id,omitempty"`
	Date        *time.Time `json:"date,omitempty"`
	StartTime   *time.Time `json:"start_time,omitempty"`
	EndTime     *time.Time `json:"end_time,omitempty"`
	Hours       *float64   `json:"hours,omitempty" validate:"omitempty,gt=0,lte=24"`
	Description *string    `json:"description,omitempty"`
	IsBillable  *bool      `json:"is_billable,omitempty"`
	HourlyRate  *float64   `json:"hourly_rate,omitempty"`
	Status      *string    `json:"status,omitempty" validate:"omitempty,oneof=draft submitted approved billed"`
}

// ProjectMemberCreateRequest 项目成员创建请求
type ProjectMemberCreateRequest struct {
	ProjectID uint   `json:"project_id" validate:"required"`
	UserID    uint   `json:"user_id" validate:"required"`
	Role      string `json:"role" validate:"required,oneof=Manager Member Viewer"`
}

// ProjectMemberUpdateRequest 项目成员更新请求
type ProjectMemberUpdateRequest struct {
	Role string `json:"role" validate:"required,oneof=Manager Member Viewer"`
}

// ProjectReportRequest 项目报告请求
type ProjectReportRequest struct {
	ProjectID uint       `json:"project_id,omitempty"`
	StartDate *time.Time `json:"start_date,omitempty"`
	EndDate   *time.Time `json:"end_date,omitempty"`
	UserID    *uint      `json:"user_id,omitempty"`
	Status    string     `json:"status,omitempty" validate:"omitempty,oneof=Planning Active Completed Cancelled"`
}

// ProjectResponse 项目响应
type ProjectResponse struct {
	ID              uint       `json:"id"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	ProjectCode     string     `json:"project_code"`
	ProjectName     string     `json:"project_name"`
	Name            string     `json:"name"`
	Description     string     `json:"description"`
	StartDate       time.Time  `json:"start_date"`
	EndDate         time.Time  `json:"end_date"`
	ActualStartDate *time.Time `json:"actual_start_date"`
	ActualEndDate   *time.Time `json:"actual_end_date"`
	Status          string     `json:"status"`
	Priority        string     `json:"priority"`
	Budget          float64    `json:"budget"`
	ActualCost      float64    `json:"actual_cost"`
	Progress        float64    `json:"progress"`
	ClientID        uint       `json:"client_id"`
	CustomerID      uint       `json:"customer_id"`
	ManagerID       uint       `json:"manager_id"`
}

// ProjectFilter 项目过滤器
type ProjectFilter struct {
	Name       string     `json:"name,omitempty"`
	Status     string     `json:"status,omitempty"`
	Priority   string     `json:"priority,omitempty"`
	CustomerID *uint      `json:"customer_id,omitempty"`
	ManagerID  *uint      `json:"manager_id,omitempty"`
	StartDate  *time.Time `json:"start_date,omitempty"`
	EndDate    *time.Time `json:"end_date,omitempty"`
	Page       int        `json:"page,omitempty"`
	PageSize   int        `json:"page_size,omitempty"`
}

// TaskResponse 任务响应
type TaskResponse struct {
	ID              uint       `json:"id"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	TaskNumber      string     `json:"task_number"`
	TaskName        string     `json:"task_name"`
	Description     string     `json:"description"`
	ProjectID       uint       `json:"project_id"`
	ParentTaskID    *uint      `json:"parent_task_id"`
	AssigneeID      *uint      `json:"assignee_id"`
	StartDate       time.Time  `json:"start_date"`
	EndDate         time.Time  `json:"end_date"`
	ActualStartDate *time.Time `json:"actual_start_date"`
	ActualEndDate   *time.Time `json:"actual_end_date"`
	EstimatedHours  float64    `json:"estimated_hours"`
	ActualHours     float64    `json:"actual_hours"`
	Progress        float64    `json:"progress"`
	Priority        string     `json:"priority"`
	Status          string     `json:"status"`
}

// TaskFilter 任务过滤器
type TaskFilter struct {
	ProjectID  *uint      `json:"project_id,omitempty"`
	Status     string     `json:"status,omitempty"`
	Priority   string     `json:"priority,omitempty"`
	AssigneeID *uint      `json:"assignee_id,omitempty"`
	StartDate  *time.Time `json:"start_date,omitempty"`
	EndDate    *time.Time `json:"end_date,omitempty"`
	Page       int        `json:"page,omitempty"`
	PageSize   int        `json:"page_size,omitempty"`
}

// MilestoneCreateRequest 里程碑创建请求
type MilestoneCreateRequest struct {
	ProjectID     uint      `json:"project_id" validate:"required"`
	MilestoneName string    `json:"milestone_name" validate:"required"`
	Name          string    `json:"name" validate:"required,min=1,max=100"`
	Description   string    `json:"description,omitempty"`
	DueDate       time.Time `json:"due_date" validate:"required"`
	Status        string    `json:"status" validate:"required,oneof=Pending Completed"`
}

// MilestoneUpdateRequest 里程碑更新请求
type MilestoneUpdateRequest struct {
	MilestoneName *string    `json:"milestone_name,omitempty"`
	Name          *string    `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	Description   *string    `json:"description,omitempty"`
	DueDate       *time.Time `json:"due_date,omitempty"`
	Status        *string    `json:"status,omitempty" validate:"omitempty,oneof=Pending Completed"`
}

// MilestoneResponse 里程碑响应
type MilestoneResponse struct {
	ID            uint      `json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	ProjectID     uint      `json:"project_id"`
	MilestoneName string    `json:"milestone_name"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	DueDate       time.Time `json:"due_date"`
	Status        string    `json:"status"`
}

// TimeEntryResponse 时间条目响应
type TimeEntryResponse struct {
	ID          uint      `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	EmployeeID  uint      `json:"employee_id"`
	ProjectID   uint      `json:"project_id"`
	TaskID      *uint     `json:"task_id,omitempty"`
	Date        time.Time `json:"date"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Hours       float64   `json:"hours"`
	Description string    `json:"description"`
	IsBillable  bool      `json:"is_billable"`
	HourlyRate  float64   `json:"hourly_rate"`
	Amount      float64   `json:"amount"`
	Status      string    `json:"status"`
}
