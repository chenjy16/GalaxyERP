package routes

import (
	"github.com/galaxyerp/galaxyErp/internal/container"
	"github.com/gin-gonic/gin"
)

// RegisterSystemRoutes 注册系统管理相关路由
func RegisterSystemRoutes(router *gin.RouterGroup, container *container.Container) {
	sys := router.Group("/system")

	// 用户和角色管理
	registerUserManagementRoutes(sys, container)

	// 权限管理
	registerPermissionManagementRoutes(sys, container)

	// 组织架构管理
	registerOrganizationRoutes(sys, container)

	// 系统配置管理
	registerSystemConfigRoutes(sys, container)

	// 系统维护功能
	registerMaintenanceRoutes(sys, container)
}

// registerUserManagementRoutes 注册用户和角色管理路由
func registerUserManagementRoutes(sys *gin.RouterGroup, container *container.Container) {
	// 用户管理 - 复用UserController的已有方法
	users := sys.Group("/users")
	{
		users.POST("/", container.UserController.CreateUser)
		users.GET("/", container.UserController.GetUsers)
		users.GET("/:id", container.UserController.GetUser)
		users.PUT("/:id", container.UserController.UpdateUser)
		users.DELETE("/:id", container.UserController.DeleteUser)
		users.POST("/search", container.UserController.SearchUsers)
		users.POST("/:id/assign-role", container.UserController.AssignRole)
		users.POST("/:id/remove-role", container.UserController.RemoveRole)
	}

	// 角色管理 - 复用UserController的已有方法
	roles := sys.Group("/roles")
	{
		roles.POST("/", container.UserController.CreateRole)
		roles.GET("/", container.UserController.GetRoles)
		roles.GET("/:id", container.UserController.GetRole)
		roles.PUT("/:id", container.UserController.UpdateRole)
		roles.DELETE("/:id", container.UserController.DeleteRole)
		roles.POST("/:id/assign-permission", container.UserController.AssignPermission)
		roles.POST("/:id/remove-permission", container.UserController.RemovePermission)
	}
}

// registerPermissionManagementRoutes 注册权限管理路由
func registerPermissionManagementRoutes(sys *gin.RouterGroup, container *container.Container) {
	// 权限管理
	permissions := sys.Group("/permissions")
	{
		permissions.POST("/", container.SystemController.CreatePermission)
		permissions.GET("/", container.SystemController.GetPermissions)
		permissions.GET("/:id", container.SystemController.GetPermission)
		permissions.PUT("/:id", container.SystemController.UpdatePermission)
		permissions.DELETE("/:id", container.SystemController.DeletePermission)
	}

	// 数据权限管理
	dataPermissions := sys.Group("/data-permissions")
	{
		dataPermissions.POST("/", container.SystemController.CreateDataPermission)
		dataPermissions.GET("/", container.SystemController.GetDataPermissions)
		dataPermissions.GET("/:id", container.SystemController.GetDataPermission)
		dataPermissions.PUT("/:id", container.SystemController.UpdateDataPermission)
		dataPermissions.DELETE("/:id", container.SystemController.DeleteDataPermission)
	}
}

// registerOrganizationRoutes 注册组织架构管理路由
func registerOrganizationRoutes(sys *gin.RouterGroup, container *container.Container) {
	// 公司管理
	companies := sys.Group("/companies")
	{
		companies.POST("/", container.SystemController.CreateCompany)
		companies.GET("/", container.SystemController.GetCompanies)
		companies.GET("/:id", container.SystemController.GetCompany)
		companies.PUT("/:id", container.SystemController.UpdateCompany)
		companies.DELETE("/:id", container.SystemController.DeleteCompany)
	}

	// 部门管理
	departments := sys.Group("/departments")
	{
		departments.POST("/", container.SystemController.CreateDepartment)
		departments.GET("/", container.SystemController.GetDepartments)
		departments.GET("/:id", container.SystemController.GetDepartment)
		departments.PUT("/:id", container.SystemController.UpdateDepartment)
		departments.DELETE("/:id", container.SystemController.DeleteDepartment)
	}

	// 职位管理
	positions := sys.Group("/positions")
	{
		positions.POST("/", container.SystemController.CreatePosition)
		positions.GET("/", container.SystemController.GetPositions)
		positions.GET("/:id", container.SystemController.GetPosition)
		positions.PUT("/:id", container.SystemController.UpdatePosition)
		positions.DELETE("/:id", container.SystemController.DeletePosition)
	}
}

// registerSystemConfigRoutes 注册系统配置管理路由
func registerSystemConfigRoutes(sys *gin.RouterGroup, container *container.Container) {
	// 系统配置
	configs := sys.Group("/configs")
	{
		configs.POST("/", container.SystemController.CreateSystemConfig)
		configs.GET("/", container.SystemController.GetSystemConfigs)
		configs.GET("/:id", container.SystemController.GetSystemConfig)
		configs.PUT("/:id", container.SystemController.UpdateSystemConfig)
		configs.DELETE("/:id", container.SystemController.DeleteSystemConfig)
	}
}

// registerMaintenanceRoutes 注册系统维护功能路由
func registerMaintenanceRoutes(sys *gin.RouterGroup, container *container.Container) {

	// 审计日志 - 简化版本
	auditLogs := sys.Group("/audit-logs")
	{
		auditLogs.POST("/", container.SystemController.CreateAuditLog)
		auditLogs.GET("/", container.SystemController.GetAuditLogs)
	}
}
