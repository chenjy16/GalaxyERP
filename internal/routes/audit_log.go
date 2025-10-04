package routes

import (
	"github.com/galaxyerp/galaxyErp/internal/handlers"
	"github.com/gin-gonic/gin"
)

// SetupAuditLogRoutes 设置审计日志路由
func SetupAuditLogRoutes(router *gin.RouterGroup, auditLogHandler *handlers.AuditLogHandler) {
	// 审计日志路由组，需要认证
	auditLogGroup := router.Group("/audit-logs")
	{
		// 获取审计日志列表
		auditLogGroup.GET("", auditLogHandler.GetAuditLogs)

		// 根据ID获取审计日志
		auditLogGroup.GET("/:id", auditLogHandler.GetAuditLogByID)

		// 获取用户的审计日志
		auditLogGroup.GET("/user/:userId", auditLogHandler.GetUserAuditLogs)

		// 获取资源的审计日志
		auditLogGroup.GET("/resource/:resource/:resourceId", auditLogHandler.GetResourceAuditLogs)

		// 清理旧日志（管理员权限）
		auditLogGroup.DELETE("/cleanup", auditLogHandler.CleanupOldLogs)
	}
}
