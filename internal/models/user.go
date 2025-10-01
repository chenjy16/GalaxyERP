package models

import (
	"time"
)

// User 用户模型
type User struct {
	BaseModel
	Username     string     `json:"username" gorm:"uniqueIndex;not null"`
	Email        string     `json:"email" gorm:"uniqueIndex;not null"`
	Password     string     `json:"-" gorm:"not null"`
	FirstName    string     `json:"first_name,omitempty"`
	LastName     string     `json:"last_name,omitempty"`
	IsActive     bool       `json:"is_active" gorm:"default:true"`
	Role         string     `json:"role,omitempty"`
	CompanyID    *uint      `json:"company_id,omitempty"`
	DepartmentID *uint      `json:"department_id,omitempty"`
	PositionID   *uint      `json:"position_id,omitempty"`
	Phone        string     `json:"phone,omitempty"`
	LastLogin    *time.Time `json:"last_login,omitempty"`
	Avatar       string     `json:"avatar,omitempty"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty"`
	
	// 关联
	Roles        []Role       `json:"roles,omitempty" gorm:"many2many:user_roles;"`
	Permissions  []Permission `json:"permissions,omitempty" gorm:"many2many:user_permissions;"`
	Department   *Department  `json:"department,omitempty" gorm:"foreignKey:DepartmentID"`
}

// Role 角色模型
type Role struct {
	BaseModel
	Name        string       `json:"name" gorm:"uniqueIndex;size:100;not null"`
	Description string       `json:"description,omitempty" gorm:"type:text"`
	IsActive    bool         `json:"is_active" gorm:"default:true;index"`
	
	// 关联
	Users       []User       `json:"users,omitempty" gorm:"many2many:user_roles;"`
	Permissions []Permission `json:"permissions,omitempty" gorm:"many2many:role_permissions;"`
}

// Permission 权限模型
type Permission struct {
	BaseModel
	Name        string `json:"name" gorm:"uniqueIndex;size:100;not null"`
	Resource    string `json:"resource" gorm:"size:100;not null;index"`
	Action      string `json:"action" gorm:"size:50;not null;index"`
	Description string `json:"description,omitempty" gorm:"type:text"`
	
	// 关联
	Users []User `json:"users,omitempty" gorm:"many2many:user_permissions;"`
	Roles []Role `json:"roles,omitempty" gorm:"many2many:role_permissions;"`
}

// Department 部门模型
type Department struct {
	BaseModel
	Name        string       `json:"name" gorm:"not null"`
	Code        string       `json:"code" gorm:"uniqueIndex;not null"`
	Description string       `json:"description,omitempty"`
	CompanyID   uint         `json:"company_id" gorm:"not null"`
	ParentID    *uint        `json:"parent_id,omitempty"`
	IsActive    bool         `json:"is_active" gorm:"default:true"`
	CreatedBy   *uint        `json:"created_by,omitempty"`
	UpdatedBy   *uint        `json:"updated_by,omitempty"`
	
	// 关联
	Parent   *Department  `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Children []Department `json:"children,omitempty" gorm:"foreignKey:ParentID"`
	Users    []User       `json:"users,omitempty" gorm:"foreignKey:DepartmentID"`
}

// UserSession 用户会话模型
type UserSession struct {
	BaseModel
	UserID    uint      `json:"user_id" gorm:"index;not null"`
	Token     string    `json:"token" gorm:"uniqueIndex;size:500;not null"`
	ExpiresAt time.Time `json:"expires_at" gorm:"index;not null"`
	IPAddress string    `json:"ip_address,omitempty" gorm:"size:45"`
	UserAgent string    `json:"user_agent,omitempty" gorm:"type:text"`
	IsActive  bool      `json:"is_active" gorm:"default:true;index"`
	
	// 关联
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// UserRole 用户角色关联表
type UserRole struct {
	UserID uint `json:"user_id" gorm:"primaryKey"`
	RoleID uint `json:"role_id" gorm:"primaryKey"`
	
	// 关联
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Role Role `json:"role,omitempty" gorm:"foreignKey:RoleID"`
}

// RolePermission 角色权限关联表
type RolePermission struct {
	RoleID       uint `json:"role_id" gorm:"primaryKey"`
	PermissionID uint `json:"permission_id" gorm:"primaryKey"`
	
	// 关联
	Role       Role       `json:"role,omitempty" gorm:"foreignKey:RoleID"`
	Permission Permission `json:"permission,omitempty" gorm:"foreignKey:PermissionID"`
}