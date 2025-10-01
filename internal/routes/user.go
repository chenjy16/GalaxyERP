package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/galaxyerp/galaxyErp/internal/container"
)

// RegisterUserRoutes 注册用户相关路由
func RegisterUserRoutes(router *gin.RouterGroup, container *container.Container) {
	users := router.Group("/users")
	{
		// 用户注册和登录
		users.POST("/register", container.UserController.Register)
		users.POST("/login", container.UserController.Login)
		
		// 用户个人资料管理
		users.GET("/profile", container.UserController.GetProfile)
		users.PUT("/profile", container.UserController.UpdateProfile)
		users.PUT("/password", container.UserController.ChangePassword)
		
		// 用户管理（管理员功能）
		users.GET("/", container.UserController.GetUsers)
		users.DELETE("/:id", container.UserController.DeleteUser)
		users.POST("/search", container.UserController.SearchUsers)
	}
}