package dto

import (
	"time"
)

// UserCreateRequest 用户创建请求
type UserCreateRequest struct {
	Username     string `json:"username" binding:"required,min=3,max=50"`
	Email        string `json:"email" binding:"required,email"`
	Password     string `json:"password" binding:"required,min=6"`
	FirstName    string `json:"first_name" binding:"required,max=50"`
	LastName     string `json:"last_name" binding:"required,max=50"`
	Phone        string `json:"phone,omitempty" binding:"omitempty,max=20"`
	DepartmentID uint   `json:"department_id,omitempty"`
	RoleIDs      []uint `json:"role_ids,omitempty"`
}

// UserUpdateRequest 用户更新请求
type UserUpdateRequest struct {
	Username     string `json:"username,omitempty" binding:"omitempty,min=3,max=50"`
	Email        string `json:"email,omitempty" binding:"omitempty,email"`
	FirstName    string `json:"first_name,omitempty" binding:"omitempty,max=50"`
	LastName     string `json:"last_name,omitempty" binding:"omitempty,max=50"`
	Phone        string `json:"phone,omitempty" binding:"omitempty,max=20"`
	DepartmentID *uint  `json:"department_id,omitempty"`
	RoleIDs      []uint `json:"role_ids,omitempty"`
	IsActive     *bool  `json:"is_active,omitempty"`
}

// UserResponse 用户响应
type UserResponse struct {
	ID          uint                `json:"id"`
	Username    string              `json:"username"`
	Email       string              `json:"email"`
	FirstName   string              `json:"first_name"`
	LastName    string              `json:"last_name"`
	Phone       string              `json:"phone,omitempty"`
	IsActive    bool                `json:"is_active"`
	LastLoginAt *time.Time          `json:"last_login_at,omitempty"`
	Department  *DepartmentResponse `json:"department,omitempty"`
	Roles       []RoleResponse      `json:"roles,omitempty"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
}

// UserListResponse 用户列表响应
type UserListResponse struct {
	ID          uint       `json:"id"`
	Username    string     `json:"username"`
	Email       string     `json:"email"`
	FirstName   string     `json:"first_name"`
	LastName    string     `json:"last_name"`
	IsActive    bool       `json:"is_active"`
	LastLoginAt *time.Time `json:"last_login_at,omitempty"`
	Department  string     `json:"department,omitempty"`
	Roles       []string   `json:"roles,omitempty"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username  string `json:"username" binding:"required,min=3,max=50"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	FirstName string `json:"first_name" binding:"required,max=50"`
	LastName  string `json:"last_name" binding:"required,max=50"`
	Phone     string `json:"phone,omitempty" binding:"omitempty,max=20"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token        string       `json:"token"`
	RefreshToken string       `json:"refresh_token"`
	ExpiresAt    time.Time    `json:"expires_at"`
	User         UserResponse `json:"user"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// ResetPasswordRequest 重置密码请求
type ResetPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ResetPasswordConfirmRequest 重置密码确认请求
type ResetPasswordConfirmRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// RoleFilter 角色过滤器
type RoleFilter struct {
	Page     int    `form:"page" binding:"min=1" json:"page"`
	Limit    int    `form:"limit" binding:"min=1,max=100" json:"limit"`
	Name     string `form:"name" json:"name,omitempty"`
	IsActive *bool  `form:"is_active" json:"is_active,omitempty"`
}

// RoleCreateRequest 角色创建请求
type RoleCreateRequest struct {
	Name        string `json:"name" binding:"required,max=50"`
	Code        string `json:"code" binding:"required,max=50"`
	Description string `json:"description,omitempty"`
	IsActive    bool   `json:"is_active"`
	Permissions []uint `json:"permissions,omitempty"`
}

// RoleUpdateRequest 角色更新请求
type RoleUpdateRequest struct {
	Name        string `json:"name,omitempty" binding:"omitempty,max=50"`
	Description string `json:"description,omitempty"`
	Permissions []uint `json:"permissions,omitempty"`
	IsActive    *bool  `json:"is_active,omitempty"`
}

// RoleResponse 角色响应
type RoleResponse struct {
	ID          uint                 `json:"id"`
	Name        string               `json:"name"`
	Code        string               `json:"code"`
	Description string               `json:"description,omitempty"`
	IsActive    bool                 `json:"is_active"`
	Permissions []PermissionResponse `json:"permissions,omitempty"`
	CreatedAt   time.Time            `json:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at"`
}

// PermissionCreateRequest 权限创建请求
type PermissionCreateRequest struct {
	Name        string `json:"name" binding:"required,max=50"`
	Code        string `json:"code" binding:"required,max=50"`
	Resource    string `json:"resource" binding:"required,max=50"`
	Action      string `json:"action" binding:"required,max=50"`
	Description string `json:"description,omitempty"`
}

// PermissionUpdateRequest 权限更新请求
type PermissionUpdateRequest struct {
	Name        string `json:"name,omitempty" binding:"omitempty,max=50"`
	Description string `json:"description,omitempty"`
	IsActive    *bool  `json:"is_active,omitempty"`
}

// PermissionResponse 权限响应
type PermissionResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Resource    string    `json:"resource"`
	Action      string    `json:"action"`
	Description string    `json:"description,omitempty"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// DepartmentCreateRequest 部门创建请求
type DepartmentCreateRequest struct {
	Name        string `json:"name" binding:"required,max=100"`
	Code        string `json:"code" binding:"required,max=50"`
	Description string `json:"description,omitempty"`
	ParentID    *uint  `json:"parent_id,omitempty"`
	ManagerID   *uint  `json:"manager_id,omitempty"`
}

// DepartmentUpdateRequest 部门更新请求
type DepartmentUpdateRequest struct {
	Name        string `json:"name,omitempty" binding:"omitempty,max=100"`
	Description string `json:"description,omitempty"`
	ParentID    *uint  `json:"parent_id,omitempty"`
	ManagerID   *uint  `json:"manager_id,omitempty"`
	IsActive    *bool  `json:"is_active,omitempty"`
}

// DepartmentResponse 部门响应
type DepartmentResponse struct {
	ID          uint                 `json:"id"`
	Name        string               `json:"name"`
	Code        string               `json:"code"`
	Description string               `json:"description,omitempty"`
	IsActive    bool                 `json:"is_active"`
	Parent      *DepartmentResponse  `json:"parent,omitempty"`
	Manager     *UserResponse        `json:"manager,omitempty"`
	Children    []DepartmentResponse `json:"children,omitempty"`
	CreatedAt   time.Time            `json:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at"`
}

// UserSearchRequest 用户搜索请求
type UserSearchRequest struct {
	SearchRequest
	DepartmentID *uint  `json:"department_id,omitempty" form:"department_id"`
	RoleID       *uint  `json:"role_id,omitempty" form:"role_id"`
	IsActive     *bool  `json:"is_active,omitempty" form:"is_active"`
	Email        string `json:"email,omitempty" form:"email"`
}

// UserProfileResponse 用户档案响应
type UserProfileResponse struct {
	UserResponse
	Permissions []string `json:"permissions"`
	MenuItems   []string `json:"menu_items"`
}

// RefreshTokenRequest 刷新令牌请求
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// LogoutRequest 登出请求
type LogoutRequest struct {
	Token string `json:"token" binding:"required"`
}

// CreateUserRequest 创建用户请求（服务层使用）
type CreateUserRequest struct {
	Username     string `json:"username" binding:"required,min=3,max=50"`
	Email        string `json:"email" binding:"required,email"`
	Password     string `json:"password" binding:"required,min=6"`
	FirstName    string `json:"first_name" binding:"required,max=50"`
	LastName     string `json:"last_name" binding:"required,max=50"`
	Phone        string `json:"phone,omitempty" binding:"omitempty,max=20"`
	CompanyID    uint   `json:"company_id,omitempty"`
	DepartmentID uint   `json:"department_id,omitempty"`
	PositionID   uint   `json:"position_id,omitempty"`
}

// UpdateUserRequest 更新用户请求（服务层使用）
type UpdateUserRequest struct {
	Username     string `json:"username,omitempty" binding:"omitempty,min=3,max=50"`
	Email        string `json:"email,omitempty" binding:"omitempty,email"`
	FirstName    string `json:"first_name,omitempty" binding:"omitempty,max=50"`
	LastName     string `json:"last_name,omitempty" binding:"omitempty,max=50"`
	Phone        string `json:"phone,omitempty" binding:"omitempty,max=20"`
	CompanyID    *uint  `json:"company_id,omitempty"`
	DepartmentID *uint  `json:"department_id,omitempty"`
	PositionID   *uint  `json:"position_id,omitempty"`
	IsActive     *bool  `json:"is_active,omitempty"`
}

// UserFilter 用户过滤器
type UserFilter struct {
	PaginationRequest
	Username     string `json:"username,omitempty" form:"username"`
	Email        string `json:"email,omitempty" form:"email"`
	DepartmentID *uint  `json:"department_id,omitempty" form:"department_id"`
	IsActive     *bool  `json:"is_active,omitempty" form:"is_active"`
}

// PermissionFilter 权限过滤器
type PermissionFilter struct {
	Page     int    `form:"page" binding:"min=1" json:"page"`
	Limit    int    `form:"limit" binding:"min=1,max=100" json:"limit"`
	Name     string `form:"name" json:"name,omitempty"`
	Resource string `form:"resource" json:"resource,omitempty"`
	Action   string `form:"action" json:"action,omitempty"`
}

// DepartmentFilter 部门过滤器
type DepartmentFilter struct {
	Page     int    `form:"page" binding:"min=1" json:"page"`
	Limit    int    `form:"limit" binding:"min=1,max=100" json:"limit"`
	Name     string `form:"name" json:"name,omitempty"`
	ParentID *uint  `form:"parent_id" json:"parent_id,omitempty"`
}
