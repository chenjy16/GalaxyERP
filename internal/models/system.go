package models

import (
	"time"
)

// Company 公司模型
type Company struct {
	DescriptionModel
	Address string `json:"address,omitempty" gorm:"type:text"`
	Phone   string `json:"phone,omitempty" gorm:"size:20"`
	Email   string `json:"email,omitempty" gorm:"size:255"`

	// 关联
	Departments []Department `json:"departments,omitempty" gorm:"foreignKey:CompanyID"`
	Users       []User       `json:"users,omitempty" gorm:"foreignKey:CompanyID"`
}

// DataPermission 数据权限模型
type DataPermission struct {
	StatusModel
	Name        string `json:"name" gorm:"size:255;not null"`
	Description string `json:"description,omitempty" gorm:"type:text"`
	Resource    string `json:"resource" gorm:"size:100;not null;index"`
	Scope       string `json:"scope" gorm:"size:50;not null"`         // all, own, department, company
	Constraint  string `json:"constraint,omitempty" gorm:"type:text"` // JSON格式的约束条件
	RoleID      *uint  `json:"role_id,omitempty" gorm:"index"`
	UserID      *uint  `json:"user_id,omitempty" gorm:"index"`

	// 关联
	Role *Role `json:"role,omitempty" gorm:"foreignKey:RoleID"`
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// SystemConfig 系统配置模型
type SystemConfig struct {
	StatusModel
	Key         string `json:"key" gorm:"uniqueIndex;size:100;not null"`
	Value       string `json:"value" gorm:"type:text"`
	Description string `json:"description,omitempty" gorm:"type:text"`
	DataType    string `json:"data_type" gorm:"size:20;default:'string'"` // string, number, boolean, json
	Category    string `json:"category" gorm:"size:50;index"`
	IsEncrypted bool   `json:"is_encrypted" gorm:"default:false"`
}

// ApprovalWorkflow 审批工作流模型
type ApprovalWorkflow struct {
	StatusModel
	Name        string `json:"name" gorm:"size:255;not null"`
	Description string `json:"description,omitempty" gorm:"type:text"`
	Resource    string `json:"resource" gorm:"size:100;not null;index"` // 关联的资源类型

	// 关联
	Steps []ApprovalStep `json:"steps,omitempty" gorm:"foreignKey:WorkflowID"`
}

// ApprovalStep 审批步骤模型
type ApprovalStep struct {
	BaseModel
	WorkflowID uint   `json:"workflow_id" gorm:"index;not null"`
	StepNumber int    `json:"step_number" gorm:"not null"`
	Name       string `json:"name" gorm:"size:255;not null"`
	Approver   string `json:"approver" gorm:"size:100;not null"`    // 审批人类型：user, role, department
	Condition  string `json:"condition,omitempty" gorm:"type:text"` // 审批条件

	// 关联
	Workflow *ApprovalWorkflow `json:"workflow,omitempty" gorm:"foreignKey:WorkflowID"`
}

// AuditLog 审计日志模型
type AuditLog struct {
	BaseModel
	UserID      uint   `json:"user_id" gorm:"index;not null"`
	Action      string `json:"action" gorm:"size:50;not null;index"`
	Resource    string `json:"resource" gorm:"size:100;not null;index"`
	ResourceID  uint   `json:"resource_id,omitempty" gorm:"index"`
	Description string `json:"description,omitempty" gorm:"type:text"`
	IPAddress   string `json:"ip_address,omitempty" gorm:"size:45"`
	UserAgent   string `json:"user_agent,omitempty" gorm:"type:text"`

	// 关联
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// Backup 备份模型
type Backup struct {
	BaseModel
	Name        string     `json:"name" gorm:"size:255;not null"`
	Description string     `json:"description,omitempty" gorm:"type:text"`
	FilePath    string     `json:"file_path" gorm:"size:500;not null"`
	FileSize    int64      `json:"file_size" gorm:"default:0"`
	Status      string     `json:"status" gorm:"size:20;default:'pending'"` // pending, completed, failed
	StartedAt   *time.Time `json:"started_at,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	ErrorMsg    string     `json:"error_msg,omitempty" gorm:"type:text"`
}
