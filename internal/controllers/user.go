package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/galaxyerp/galaxyErp/internal/dto"
	"github.com/galaxyerp/galaxyErp/internal/services"
)

// UserController 用户控制器
type UserController struct {
	userService services.UserService
	utils      *ControllerUtils
}

// NewUserController 创建用户控制器实例
func NewUserController(userService services.UserService) *UserController {
	return &UserController{
		userService: userService,
		utils:      NewControllerUtils(),
	}
}

// Login 用户登录
// @Summary 用户登录
// @Description 用户登录
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param login body dto.LoginRequest true "登录信息"
// @Success 200 {object} dto.LoginResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/users/login [post]
func (c *UserController) Login(ctx *gin.Context) {
	var req dto.LoginRequest
	if !c.utils.BindJSON(ctx, &req) {
		return
	}

	response, err := c.userService.Login(ctx.Request.Context(), &req)
	if err != nil {
		c.utils.RespondUnauthorized(ctx, "用户名或密码错误")
		return
	}

	c.utils.RespondOK(ctx, response)
}

// Register 用户注册
// @Summary 用户注册
// @Description 用户注册
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param register body dto.RegisterRequest true "注册信息"
// @Success 201 {object} dto.UserResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/users/register [post]
func (c *UserController) Register(ctx *gin.Context) {
	var req dto.RegisterRequest
	if !c.utils.BindJSON(ctx, &req) {
		return
	}

	user, err := c.userService.Register(ctx.Request.Context(), &req)
	if err != nil {
		c.utils.RespondInternalError(ctx, "注册失败")
		return
	}

	c.utils.RespondCreated(ctx, user)
}

// GetProfile 获取用户信息
// @Summary 获取用户信息
// @Description 获取当前用户信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} dto.UserResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/users/profile [get]
func (c *UserController) GetProfile(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		c.utils.RespondUnauthorized(ctx, "未授权访问")
		return
	}

	user, err := c.userService.GetProfile(ctx.Request.Context(), userID.(uint))
	if err != nil {
		c.utils.RespondInternalError(ctx, "获取用户信息失败")
		return
	}

	c.utils.RespondOK(ctx, user)
}

// UpdateProfile 更新用户信息
// @Summary 更新用户信息
// @Description 更新当前用户信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param profile body dto.UserUpdateRequest true "用户信息"
// @Success 200 {object} dto.BaseResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/users/profile [put]
func (c *UserController) UpdateProfile(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		c.utils.RespondUnauthorized(ctx, "未授权访问")
		return
	}

	var req dto.UserUpdateRequest
	if !c.utils.BindJSON(ctx, &req) {
		return
	}

	err := c.userService.UpdateProfile(ctx.Request.Context(), userID.(uint), &req)
	if err != nil {
		c.utils.RespondInternalError(ctx, "更新用户信息失败")
		return
	}

	c.utils.RespondSuccess(ctx, "更新用户信息成功")
}

// ChangePassword 修改密码
// @Summary 修改密码
// @Description 修改当前用户密码
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param password body dto.ChangePasswordRequest true "密码信息"
// @Success 200 {object} dto.BaseResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/users/change-password [post]
func (c *UserController) ChangePassword(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		c.utils.RespondUnauthorized(ctx, "未授权访问")
		return
	}

	var req dto.ChangePasswordRequest
	if !c.utils.BindJSON(ctx, &req) {
		return
	}

	err := c.userService.ChangePassword(ctx.Request.Context(), userID.(uint), &req)
	if err != nil {
		c.utils.RespondInternalError(ctx, "修改密码失败")
		return
	}

	c.utils.RespondSuccess(ctx, "修改密码成功")
}

// CreateUser 创建用户
func (c *UserController) CreateUser(ctx *gin.Context) {
	c.utils.RespondNotImplemented(ctx, "创建用户功能待实现")
}

// GetUser 获取单个用户
func (c *UserController) GetUser(ctx *gin.Context) {
	c.utils.RespondNotImplemented(ctx, "获取用户功能待实现")
}

// UpdateUser 更新用户
func (c *UserController) UpdateUser(ctx *gin.Context) {
	c.utils.RespondNotImplemented(ctx, "更新用户功能待实现")
}

// DeleteUser 删除用户
func (c *UserController) DeleteUser(ctx *gin.Context) {
	c.utils.RespondNotImplemented(ctx, "删除用户功能待实现")
}

// GetUsers 获取用户列表
func (c *UserController) GetUsers(ctx *gin.Context) {
	c.utils.RespondNotImplemented(ctx, "获取用户列表功能待实现")
}

// SearchUsers 搜索用户
func (c *UserController) SearchUsers(ctx *gin.Context) {
	c.utils.RespondNotImplemented(ctx, "搜索用户功能待实现")
}

// AssignRole 分配角色
func (c *UserController) AssignRole(ctx *gin.Context) {
	c.utils.RespondNotImplemented(ctx, "分配角色功能待实现")
}

// RemoveRole 移除角色
func (c *UserController) RemoveRole(ctx *gin.Context) {
	c.utils.RespondNotImplemented(ctx, "移除角色功能待实现")
}

// CreateRole 创建角色
func (c *UserController) CreateRole(ctx *gin.Context) {
	c.utils.RespondNotImplemented(ctx, "创建角色功能待实现")
}

// GetRoles 获取角色列表
func (c *UserController) GetRoles(ctx *gin.Context) {
	c.utils.RespondNotImplemented(ctx, "获取角色列表功能待实现")
}

// GetRole 获取角色
func (c *UserController) GetRole(ctx *gin.Context) {
	c.utils.RespondNotImplemented(ctx, "获取角色功能待实现")
}

// UpdateRole 更新角色
func (c *UserController) UpdateRole(ctx *gin.Context) {
	c.utils.RespondNotImplemented(ctx, "更新角色功能待实现")
}

// DeleteRole 删除角色
func (c *UserController) DeleteRole(ctx *gin.Context) {
	c.utils.RespondNotImplemented(ctx, "删除角色功能待实现")
}

// AssignPermission 分配权限
func (c *UserController) AssignPermission(ctx *gin.Context) {
	c.utils.RespondNotImplemented(ctx, "分配权限功能待实现")
}

// RemovePermission 移除权限
func (c *UserController) RemovePermission(ctx *gin.Context) {
	c.utils.RespondNotImplemented(ctx, "移除权限功能待实现")
}