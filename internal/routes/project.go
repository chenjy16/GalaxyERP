package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/galaxyerp/galaxyErp/internal/container"
)

// RegisterProjectRoutes 注册项目管理相关路由
func RegisterProjectRoutes(router *gin.RouterGroup, container *container.Container) {
	projectController := container.ProjectController

	// 项目管理
	projects := router.Group("/projects")
	{
		projects.POST("/", projectController.CreateProject)
		projects.GET("/", projectController.ListProjects)
		projects.GET("/:id", projectController.GetProject)
		projects.PUT("/:id", projectController.UpdateProject)
		projects.DELETE("/:id", projectController.DeleteProject)
	}

	// 任务管理
	tasks := router.Group("/tasks")
	{
		tasks.POST("/", projectController.CreateTask)
		tasks.GET("/", projectController.ListTasks)
		tasks.GET("/:id", projectController.GetTask)
		tasks.PUT("/:id", projectController.UpdateTask)
		tasks.DELETE("/:id", projectController.DeleteTask)
	}

	// 里程碑路由
	milestones := router.Group("/milestones")
	{
		milestones.POST("/", projectController.CreateMilestone)
		milestones.GET("/:id", projectController.GetMilestone)
		milestones.PUT("/:id", projectController.UpdateMilestone)
		milestones.DELETE("/:id", projectController.DeleteMilestone)
	}

	// 工时记录路由
	timeEntries := router.Group("/time-entries")
	{
		timeEntries.POST("/", projectController.CreateTimeEntry)
		timeEntries.GET("/:id", projectController.GetTimeEntry)
		timeEntries.PUT("/:id", projectController.UpdateTimeEntry)
		timeEntries.DELETE("/:id", projectController.DeleteTimeEntry)
	}

	// 项目相关的列表路由 - 使用不同的路径避免冲突
	router.GET("/project-milestones/:project_id", projectController.ListMilestones)
	router.GET("/project-time-entries/:project_id", projectController.ListTimeEntries)
}