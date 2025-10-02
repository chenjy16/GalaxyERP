package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/galaxyerp/galaxyErp/internal/container"
	"github.com/galaxyerp/galaxyErp/internal/controllers"
)

// RegisterHRRoutes 注册人力资源相关路由
func RegisterHRRoutes(router *gin.RouterGroup, container *container.Container) {
	// 创建HR控制器
	hrController := controllers.NewHRController(
		container.EmployeeService,
		container.AttendanceService,
		container.PayrollService,
		container.LeaveService,
	)

	// HR模块路由组
	hr := router.Group("/hr")

	// 员工管理
	employees := hr.Group("/employees")
	{
		employees.POST("/", hrController.CreateEmployee)
		employees.GET("/", hrController.GetEmployees)
		employees.POST("/search", hrController.SearchEmployees)
		employees.GET("/:id", hrController.GetEmployee)
		employees.PUT("/:id", hrController.UpdateEmployee)
		employees.DELETE("/:id", hrController.DeleteEmployee)
		employees.GET("/:id/leaves", hrController.GetEmployeeLeaves)
	}

	// 考勤管理
	attendance := hr.Group("/attendance")
	{
		attendance.POST("/", hrController.CreateAttendance)
		attendance.GET("/", hrController.GetAttendanceList)
		attendance.GET("/:id", hrController.GetAttendance)
		attendance.PUT("/:id", hrController.UpdateAttendance)
		attendance.DELETE("/:id", hrController.DeleteAttendance)
	}

	// 薪资管理
	payroll := hr.Group("/payroll")
	{
		payroll.POST("/", hrController.CreatePayroll)
		payroll.GET("/", hrController.GetPayrollList)
		payroll.GET("/:id", hrController.GetPayroll)
		payroll.PUT("/:id", hrController.UpdatePayroll)
		payroll.DELETE("/:id", hrController.DeletePayroll)
	}

	// 请假管理
	leaves := hr.Group("/leaves")
	{
		leaves.POST("/", hrController.CreateLeave)
		leaves.GET("/", hrController.GetLeaveList)
		leaves.GET("/pending", hrController.GetPendingLeaves)
		leaves.GET("/:id", hrController.GetLeave)
		leaves.PUT("/:id", hrController.UpdateLeave)
		leaves.DELETE("/:id", hrController.DeleteLeave)
		leaves.POST("/:id/approve", hrController.ApproveLeave)
		leaves.POST("/:id/cancel", hrController.CancelLeave)
	}
}