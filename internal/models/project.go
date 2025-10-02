package models

import (
	"time"
)

// Project 项目模型
type Project struct {
	AuditableModel
	ProjectNumber   string     `json:"project_number" gorm:"uniqueIndex;size:100;not null"`
	ProjectName     string     `json:"project_name" gorm:"size:255;not null"`
	Description     string     `json:"description,omitempty" gorm:"type:text"`
	CustomerID      *uint      `json:"customer_id,omitempty" gorm:"index"`
	ManagerID       uint       `json:"manager_id" gorm:"index;not null"`
	StartDate       time.Time  `json:"start_date" gorm:"index;not null"`
	EndDate         time.Time  `json:"end_date" gorm:"index;not null"`
	ActualStartDate *time.Time `json:"actual_start_date,omitempty"`
	ActualEndDate   *time.Time `json:"actual_end_date,omitempty"`
	Budget          float64    `json:"budget" gorm:"default:0"`
	ActualCost      float64    `json:"actual_cost" gorm:"default:0"`
	Progress        float64    `json:"progress" gorm:"default:0"`                      // 进度百分比
	Priority        string     `json:"priority" gorm:"size:20;default:'normal';index"` // low, normal, high, urgent
	Status          string     `json:"status" gorm:"size:50;default:'planning';index"` // planning, active, on_hold, completed, cancelled

	// 关联
	Customer *Customer        `json:"customer,omitempty" gorm:"foreignKey:CustomerID"`
	Manager  Employee         `json:"manager,omitempty" gorm:"foreignKey:ManagerID"`
	Tasks    []Task           `json:"tasks,omitempty" gorm:"foreignKey:ProjectID"`
	Members  []ProjectMember  `json:"members,omitempty" gorm:"foreignKey:ProjectID"`
	Expenses []ProjectExpense `json:"expenses,omitempty" gorm:"foreignKey:ProjectID"`
}

// Task 任务模型
type Task struct {
	AuditableModel
	TaskNumber      string     `json:"task_number" gorm:"uniqueIndex;size:100;not null"`
	TaskName        string     `json:"task_name" gorm:"size:255;not null"`
	Description     string     `json:"description,omitempty" gorm:"type:text"`
	ProjectID       uint       `json:"project_id" gorm:"index;not null"`
	ParentTaskID    *uint      `json:"parent_task_id,omitempty" gorm:"index"`
	AssigneeID      *uint      `json:"assignee_id,omitempty" gorm:"index"`
	StartDate       time.Time  `json:"start_date" gorm:"index;not null"`
	EndDate         time.Time  `json:"end_date" gorm:"index;not null"`
	ActualStartDate *time.Time `json:"actual_start_date,omitempty"`
	ActualEndDate   *time.Time `json:"actual_end_date,omitempty"`
	EstimatedHours  float64    `json:"estimated_hours" gorm:"default:0"`
	ActualHours     float64    `json:"actual_hours" gorm:"default:0"`
	Progress        float64    `json:"progress" gorm:"default:0"`                      // 进度百分比
	Priority        string     `json:"priority" gorm:"size:20;default:'normal';index"` // low, normal, high, urgent
	Status          string     `json:"status" gorm:"size:50;default:'todo';index"`     // todo, in_progress, review, done, cancelled

	// 关联
	Project     Project       `json:"project,omitempty" gorm:"foreignKey:ProjectID"`
	ParentTask  *Task         `json:"parent_task,omitempty" gorm:"foreignKey:ParentTaskID"`
	SubTasks    []Task        `json:"sub_tasks,omitempty" gorm:"foreignKey:ParentTaskID"`
	Assignee    *Employee     `json:"assignee,omitempty" gorm:"foreignKey:AssigneeID"`
	TimeEntries []TimeEntry   `json:"time_entries,omitempty" gorm:"foreignKey:TaskID"`
	Comments    []TaskComment `json:"comments,omitempty" gorm:"foreignKey:TaskID"`
}

// ProjectMember 项目成员模型
type ProjectMember struct {
	BaseModel
	ProjectID  uint       `json:"project_id" gorm:"index;not null"`
	EmployeeID uint       `json:"employee_id" gorm:"index;not null"`
	Role       string     `json:"role" gorm:"size:100;not null"`
	JoinDate   time.Time  `json:"join_date" gorm:"index;not null"`
	LeaveDate  *time.Time `json:"leave_date,omitempty" gorm:"index"`
	IsActive   bool       `json:"is_active" gorm:"default:true;index"`

	// 关联
	Project  Project  `json:"project,omitempty" gorm:"foreignKey:ProjectID"`
	Employee Employee `json:"employee,omitempty" gorm:"foreignKey:EmployeeID"`
}

// TimeEntry 工时记录模型
type TimeEntry struct {
	BaseModel
	EmployeeID  uint      `json:"employee_id" gorm:"index;not null"`
	ProjectID   uint      `json:"project_id" gorm:"index;not null"`
	TaskID      *uint     `json:"task_id,omitempty" gorm:"index"`
	Date        time.Time `json:"date" gorm:"index;not null"`
	StartTime   time.Time `json:"start_time" gorm:"not null"`
	EndTime     time.Time `json:"end_time" gorm:"not null"`
	Hours       float64   `json:"hours" gorm:"not null"`
	Description string    `json:"description,omitempty" gorm:"type:text"`
	IsBillable  bool      `json:"is_billable" gorm:"default:true;index"`
	HourlyRate  float64   `json:"hourly_rate,omitempty" gorm:"default:0"`
	Amount      float64   `json:"amount,omitempty" gorm:"default:0"`
	Status      string    `json:"status" gorm:"size:50;default:'draft';index"` // draft, submitted, approved, billed

	// 关联
	Employee Employee `json:"employee,omitempty" gorm:"foreignKey:EmployeeID"`
	Project  Project  `json:"project,omitempty" gorm:"foreignKey:ProjectID"`
	Task     *Task    `json:"task,omitempty" gorm:"foreignKey:TaskID"`
}

// TaskComment 任务评论模型
type TaskComment struct {
	BaseModel
	TaskID     uint   `json:"task_id" gorm:"index;not null"`
	EmployeeID uint   `json:"employee_id" gorm:"index;not null"`
	Comment    string `json:"comment" gorm:"type:text;not null"`

	// 关联
	Task     Task     `json:"task,omitempty" gorm:"foreignKey:TaskID"`
	Employee Employee `json:"employee,omitempty" gorm:"foreignKey:EmployeeID"`
}

// ProjectExpense 项目费用模型
type ProjectExpense struct {
	AuditableModel
	ProjectID      uint       `json:"project_id" gorm:"index;not null"`
	ExpenseDate    time.Time  `json:"expense_date" gorm:"index;not null"`
	ExpenseType    string     `json:"expense_type" gorm:"size:100;not null;index"` // travel, material, equipment, other
	Amount         float64    `json:"amount" gorm:"not null"`
	Currency       string     `json:"currency" gorm:"size:10;default:'CNY'"`
	Description    string     `json:"description" gorm:"type:text;not null"`
	Receipt        string     `json:"receipt,omitempty" gorm:"size:500"`
	IsReimbursable bool       `json:"is_reimbursable" gorm:"default:true;index"`
	Status         string     `json:"status" gorm:"size:50;default:'pending';index"` // pending, approved, rejected, reimbursed
	ApprovedBy     *uint      `json:"approved_by,omitempty" gorm:"index"`
	ApprovedAt     *time.Time `json:"approved_at,omitempty"`

	// 关联
	Project      Project   `json:"project,omitempty" gorm:"foreignKey:ProjectID"`
	ApprovedUser *Employee `json:"approved_user,omitempty" gorm:"foreignKey:ApprovedBy"`
}

// Milestone 里程碑模型
type Milestone struct {
	BaseModel
	ProjectID     uint       `json:"project_id" gorm:"index;not null"`
	MilestoneName string     `json:"milestone_name" gorm:"size:255;not null"`
	Description   string     `json:"description,omitempty" gorm:"type:text"`
	DueDate       time.Time  `json:"due_date" gorm:"index;not null"`
	CompletedDate *time.Time `json:"completed_date,omitempty"`
	Status        string     `json:"status" gorm:"size:50;default:'pending';index"` // pending, completed, overdue

	// 关联
	Project Project `json:"project,omitempty" gorm:"foreignKey:ProjectID"`
}

// ProjectDocument 项目文档模型
type ProjectDocument struct {
	BaseModel
	ProjectID    uint   `json:"project_id" gorm:"index;not null"`
	DocumentName string `json:"document_name" gorm:"size:255;not null"`
	DocumentType string `json:"document_type" gorm:"size:100;not null;index"` // requirement, design, test, report
	FilePath     string `json:"file_path" gorm:"size:500;not null"`
	FileSize     int64  `json:"file_size" gorm:"default:0"`
	Version      string `json:"version" gorm:"size:50;default:'1.0'"`
	UploadedBy   uint   `json:"uploaded_by" gorm:"index;not null"`
	Description  string `json:"description,omitempty" gorm:"type:text"`

	// 关联
	Project      Project  `json:"project,omitempty" gorm:"foreignKey:ProjectID"`
	UploadedUser Employee `json:"uploaded_user,omitempty" gorm:"foreignKey:UploadedBy"`
}

// ProjectMilestone 项目里程碑模型 (别名)
type ProjectMilestone = Milestone

// ProjectTask 项目任务模型 (别名)
type ProjectTask = Task

// TaskTimeRecord 任务工时记录模型 (别名)
type TaskTimeRecord = TimeEntry

// ProjectResource 项目资源模型
type ProjectResource struct {
	BaseModel
	ProjectID    uint       `json:"project_id" gorm:"index;not null"`
	ResourceType string     `json:"resource_type" gorm:"size:100;not null;index"` // human, equipment, material
	ResourceID   uint       `json:"resource_id" gorm:"index;not null"`
	Quantity     float64    `json:"quantity" gorm:"default:1"`
	Unit         string     `json:"unit" gorm:"size:50"`
	CostPerUnit  float64    `json:"cost_per_unit" gorm:"default:0"`
	TotalCost    float64    `json:"total_cost" gorm:"default:0"`
	AllocatedAt  time.Time  `json:"allocated_at" gorm:"index;not null"`
	ReleasedAt   *time.Time `json:"released_at,omitempty" gorm:"index"`
	Status       string     `json:"status" gorm:"size:50;default:'allocated';index"` // allocated, in_use, released

	// 关联
	Project Project `json:"project,omitempty" gorm:"foreignKey:ProjectID"`
}

// ProjectReport 项目报告模型
type ProjectReport struct {
	AuditableModel
	ProjectID   uint      `json:"project_id" gorm:"index;not null"`
	ReportType  string    `json:"report_type" gorm:"size:100;not null;index"` // progress, financial, risk, final
	ReportDate  time.Time `json:"report_date" gorm:"index;not null"`
	Title       string    `json:"title" gorm:"size:255;not null"`
	Content     string    `json:"content" gorm:"type:text;not null"`
	GeneratedBy uint      `json:"generated_by" gorm:"index;not null"`
	Status      string    `json:"status" gorm:"size:50;default:'draft';index"` // draft, submitted, approved

	// 关联
	Project   Project  `json:"project,omitempty" gorm:"foreignKey:ProjectID"`
	Generator Employee `json:"generator,omitempty" gorm:"foreignKey:GeneratedBy"`
}
