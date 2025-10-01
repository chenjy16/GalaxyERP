package dto

import (
	"time"
)

// ProjectCreateRequest 项目创建请求
type ProjectCreateRequest struct {
	ProjectName   string     `json:"project_name" binding:"required,max=255"`
	ProjectCode   string     `json:"project_code" binding:"required,max=50"`
	Description   string     `json:"description,omitempty"`
	ClientID      *uint      `json:"client_id,omitempty"`
	ManagerID     uint       `json:"manager_id" binding:"required"`
	StartDate     time.Time  `json:"start_date" binding:"required"`
	EndDate       time.Time  `json:"end_date" binding:"required"`
	Budget        float64    `json:"budget" binding:"min=0"`
	Priority      string     `json:"priority" binding:"required,oneof=low normal high urgent"`
	Status        string     `json:"status" binding:"required,oneof=planning active on_hold completed cancelled"`
}

// ProjectUpdateRequest 项目更新请求
type ProjectUpdateRequest struct {
	ProjectName   string     `json:"project_name,omitempty" binding:"omitempty,max=255"`
	Description   string     `json:"description,omitempty"`
	ClientID      *uint      `json:"client_id,omitempty"`
	ManagerID     uint       `json:"manager_id,omitempty"`
	StartDate     *time.Time `json:"start_date,omitempty"`
	EndDate       *time.Time `json:"end_date,omitempty"`
	Budget        *float64   `json:"budget,omitempty" binding:"omitempty,min=0"`
	Priority      string     `json:"priority,omitempty" binding:"omitempty,oneof=low normal high urgent"`
	Status        string     `json:"status,omitempty" binding:"omitempty,oneof=planning active on_hold completed cancelled"`
}

// ProjectResponse 项目响应
type ProjectResponse struct {
	BaseModel
	ProjectName   string     `json:"project_name"`
	ProjectCode   string     `json:"project_code"`
	Description   string     `json:"description,omitempty"`
	ClientID      *uint      `json:"client_id,omitempty"`
	ManagerID     uint       `json:"manager_id"`
	StartDate     time.Time  `json:"start_date"`
	EndDate       time.Time  `json:"end_date"`
	ActualStartDate *time.Time `json:"actual_start_date,omitempty"`
	ActualEndDate   *time.Time `json:"actual_end_date,omitempty"`
	Budget        float64    `json:"budget"`
	ActualCost    float64    `json:"actual_cost"`
	Progress      float64    `json:"progress"`
	Priority      string     `json:"priority"`
	Status        string     `json:"status"`
	Client        *CustomerResponse `json:"client,omitempty"`
	Manager       *EmployeeResponse `json:"manager,omitempty"`
}

// TaskCreateRequest 任务创建请求
type TaskCreateRequest struct {
	TaskName       string     `json:"task_name" binding:"required,max=255"`
	TaskNumber     string     `json:"task_number" binding:"required,max=100"`
	Description    string     `json:"description,omitempty"`
	ProjectID      uint       `json:"project_id" binding:"required"`
	ParentTaskID   *uint      `json:"parent_task_id,omitempty"`
	AssigneeID     *uint      `json:"assignee_id,omitempty"`
	StartDate      time.Time  `json:"start_date" binding:"required"`
	EndDate        time.Time  `json:"end_date" binding:"required"`
	EstimatedHours float64    `json:"estimated_hours" binding:"min=0"`
	Priority       string     `json:"priority" binding:"required,oneof=low normal high urgent"`
}

// TaskUpdateRequest 任务更新请求
type TaskUpdateRequest struct {
	TaskName       string     `json:"task_name,omitempty" binding:"omitempty,max=255"`
	Description    string     `json:"description,omitempty"`
	ParentTaskID   *uint      `json:"parent_task_id,omitempty"`
	AssigneeID     *uint      `json:"assignee_id,omitempty"`
	StartDate      *time.Time `json:"start_date,omitempty"`
	EndDate        *time.Time `json:"end_date,omitempty"`
	EstimatedHours *float64   `json:"estimated_hours,omitempty" binding:"omitempty,min=0"`
	ActualHours    *float64   `json:"actual_hours,omitempty" binding:"omitempty,min=0"`
	Progress       *float64   `json:"progress,omitempty" binding:"omitempty,min=0,max=100"`
	Priority       string     `json:"priority,omitempty" binding:"omitempty,oneof=low normal high urgent"`
	Status         string     `json:"status,omitempty" binding:"omitempty,oneof=todo in_progress review done cancelled"`
}

// TaskResponse 任务响应
type TaskResponse struct {
	BaseModel
	TaskNumber      string     `json:"task_number"`
	TaskName        string     `json:"task_name"`
	Description     string     `json:"description,omitempty"`
	ProjectID       uint       `json:"project_id"`
	ParentTaskID    *uint      `json:"parent_task_id,omitempty"`
	AssigneeID      *uint      `json:"assignee_id,omitempty"`
	StartDate       time.Time  `json:"start_date"`
	EndDate         time.Time  `json:"end_date"`
	ActualStartDate *time.Time `json:"actual_start_date,omitempty"`
	ActualEndDate   *time.Time `json:"actual_end_date,omitempty"`
	EstimatedHours  float64    `json:"estimated_hours"`
	ActualHours     float64    `json:"actual_hours"`
	Progress        float64    `json:"progress"`
	Priority        string     `json:"priority"`
	Status          string     `json:"status"`
	Project         *ProjectResponse `json:"project,omitempty"`
	ParentTask      *TaskResponse    `json:"parent_task,omitempty"`
	Assignee        *EmployeeResponse `json:"assignee,omitempty"`
}

// MilestoneCreateRequest 里程碑创建请求
type MilestoneCreateRequest struct {
	ProjectID     uint      `json:"project_id" binding:"required"`
	MilestoneName string    `json:"milestone_name" binding:"required,max=255"`
	Description   string    `json:"description,omitempty"`
	DueDate       time.Time `json:"due_date" binding:"required"`
}

// MilestoneUpdateRequest 里程碑更新请求
type MilestoneUpdateRequest struct {
	MilestoneName string     `json:"milestone_name,omitempty" binding:"omitempty,max=255"`
	Description   string     `json:"description,omitempty"`
	DueDate       *time.Time `json:"due_date,omitempty"`
	Status        string     `json:"status,omitempty" binding:"omitempty,oneof=pending completed overdue"`
}

// MilestoneResponse 里程碑响应
type MilestoneResponse struct {
	BaseModel
	ProjectID     uint       `json:"project_id"`
	MilestoneName string     `json:"milestone_name"`
	Description   string     `json:"description,omitempty"`
	DueDate       time.Time  `json:"due_date"`
	CompletedDate *time.Time `json:"completed_date,omitempty"`
	Status        string     `json:"status"`
	Project       *ProjectResponse `json:"project,omitempty"`
}

// TimeEntryCreateRequest 工时记录创建请求
type TimeEntryCreateRequest struct {
	ProjectID   uint      `json:"project_id" binding:"required"`
	TaskID      *uint     `json:"task_id,omitempty"`
	Date        time.Time `json:"date" binding:"required"`
	StartTime   time.Time `json:"start_time" binding:"required"`
	EndTime     time.Time `json:"end_time" binding:"required"`
	Description string    `json:"description,omitempty"`
	IsBillable  bool      `json:"is_billable"`
	HourlyRate  float64   `json:"hourly_rate,omitempty" binding:"min=0"`
}

// TimeEntryUpdateRequest 工时记录更新请求
type TimeEntryUpdateRequest struct {
	TaskID      *uint      `json:"task_id,omitempty"`
	Date        *time.Time `json:"date,omitempty"`
	StartTime   *time.Time `json:"start_time,omitempty"`
	EndTime     *time.Time `json:"end_time,omitempty"`
	Description string     `json:"description,omitempty"`
	IsBillable  *bool      `json:"is_billable,omitempty"`
	HourlyRate  *float64   `json:"hourly_rate,omitempty" binding:"omitempty,min=0"`
	Status      string     `json:"status,omitempty" binding:"omitempty,oneof=draft submitted approved billed"`
}

// TimeEntryResponse 工时记录响应
type TimeEntryResponse struct {
	BaseModel
	EmployeeID  uint      `json:"employee_id"`
	ProjectID   uint      `json:"project_id"`
	TaskID      *uint     `json:"task_id,omitempty"`
	Date        time.Time `json:"date"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Hours       float64   `json:"hours"`
	Description string    `json:"description,omitempty"`
	IsBillable  bool      `json:"is_billable"`
	HourlyRate  float64   `json:"hourly_rate"`
	Amount      float64   `json:"amount"`
	Status      string    `json:"status"`
	Employee    *EmployeeResponse `json:"employee,omitempty"`
	Project     *ProjectResponse  `json:"project,omitempty"`
	Task        *TaskResponse     `json:"task,omitempty"`
}

// ProjectMemberCreateRequest 项目成员创建请求
type ProjectMemberCreateRequest struct {
	ProjectID  uint      `json:"project_id" binding:"required"`
	EmployeeID uint      `json:"employee_id" binding:"required"`
	Role       string    `json:"role" binding:"required,max=100"`
	JoinDate   time.Time `json:"join_date" binding:"required"`
}

// ProjectMemberUpdateRequest 项目成员更新请求
type ProjectMemberUpdateRequest struct {
	Role      string     `json:"role,omitempty" binding:"omitempty,max=100"`
	LeaveDate *time.Time `json:"leave_date,omitempty"`
	IsActive  *bool      `json:"is_active,omitempty"`
}

// ProjectMemberResponse 项目成员响应
type ProjectMemberResponse struct {
	BaseModel
	ProjectID  uint       `json:"project_id"`
	EmployeeID uint       `json:"employee_id"`
	Role       string     `json:"role"`
	JoinDate   time.Time  `json:"join_date"`
	LeaveDate  *time.Time `json:"leave_date,omitempty"`
	IsActive   bool       `json:"is_active"`
	Project    *ProjectResponse  `json:"project,omitempty"`
	Employee   *EmployeeResponse `json:"employee,omitempty"`
}

// ProjectFilter 项目过滤器
type ProjectFilter struct {
	SearchRequest
	ManagerID *uint      `json:"manager_id,omitempty" form:"manager_id"`
	ClientID  *uint      `json:"client_id,omitempty" form:"client_id"`
	Status    string     `json:"status,omitempty" form:"status"`
	Priority  string     `json:"priority,omitempty" form:"priority"`
	StartDate *time.Time `json:"start_date,omitempty" form:"start_date"`
	EndDate   *time.Time `json:"end_date,omitempty" form:"end_date"`
}

// TaskFilter 任务过滤器
type TaskFilter struct {
	SearchRequest
	ProjectID  *uint      `json:"project_id,omitempty" form:"project_id"`
	AssigneeID *uint      `json:"assignee_id,omitempty" form:"assignee_id"`
	Status     string     `json:"status,omitempty" form:"status"`
	Priority   string     `json:"priority,omitempty" form:"priority"`
	StartDate  *time.Time `json:"start_date,omitempty" form:"start_date"`
	EndDate    *time.Time `json:"end_date,omitempty" form:"end_date"`
}