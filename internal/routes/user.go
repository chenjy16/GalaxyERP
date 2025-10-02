package routes

import (
	"github.com/galaxyerp/galaxyErp/internal/container"
	"github.com/gin-gonic/gin"
)

// RegisterUserRoutes 注册用户个人功能相关路由
// 注意：用户管理功能（CRUD、角色分配等）已移至 /system/users 路由组
func RegisterUserRoutes(router *gin.RouterGroup, container *container.Container) {
	users := router.Group("/users")
	{
		// 用户个人资料管理
		users.GET("/profile", container.UserController.GetProfile)
		users.PUT("/profile", container.UserController.UpdateProfile)
		users.PUT("/password", container.UserController.ChangePassword)

		// 用户个人相关的查询功能
		users.POST("/search", container.UserController.SearchUsers) // 用于选择器等场景
	}
}
