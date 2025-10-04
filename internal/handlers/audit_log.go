package handlers

import (
	"net/http"
	"strconv"

	"github.com/galaxyerp/galaxyErp/internal/common"
	"github.com/galaxyerp/galaxyErp/internal/dto"
	"github.com/galaxyerp/galaxyErp/internal/services"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// AuditLogHandler 审计日志处理器
type AuditLogHandler struct {
	auditLogService services.AuditLogService
	logger          *zap.Logger
}

// NewAuditLogHandler 创建审计日志处理器实例
func NewAuditLogHandler(auditLogService services.AuditLogService, logger *zap.Logger) *AuditLogHandler {
	return &AuditLogHandler{
		auditLogService: auditLogService,
		logger:          logger,
	}
}

// GetAuditLogs 获取审计日志列表
// @Summary 获取审计日志列表
// @Description 根据搜索条件获取审计日志列表
// @Tags 审计日志
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param user_id query int false "用户ID"
// @Param username query string false "用户名"
// @Param action query string false "操作类型"
// @Param resource query string false "资源类型"
// @Param resource_id query string false "资源ID"
// @Param status query string false "状态"
// @Param ip_address query string false "IP地址"
// @Param start_time query string false "开始时间"
// @Param end_time query string false "结束时间"
// @Success 200 {object} dto.PaginatedResponse[dto.AuditLogResponse]
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/audit-logs [get]
func (h *AuditLogHandler) GetAuditLogs(c *gin.Context) {
	var req dto.AuditLogSearchRequest

	// 绑定查询参数
	if err := c.ShouldBindQuery(&req); err != nil {
		common.LogAppError(err, "绑定审计日志搜索参数失败")
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success:    false,
			Message:    "参数错误",
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	// 设置默认分页参数
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	auditLogs, total, err := h.auditLogService.GetAuditLogs(c.Request.Context(), &req)
	if err != nil {
		common.LogAppError(err, "获取审计日志列表失败")
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success:    false,
			Message:    "获取审计日志失败",
			StatusCode: http.StatusInternalServerError,
		})
		return
	}

	totalPages := int((total + int64(req.PageSize) - 1) / int64(req.PageSize))

	c.JSON(http.StatusOK, dto.PaginatedResponse[*dto.AuditLogResponse]{
		Data:       auditLogs,
		Total:      total,
		Page:       req.Page,
		Limit:      req.PageSize,
		TotalPages: totalPages,
	})
}

// GetAuditLogByID 根据ID获取审计日志
// @Summary 根据ID获取审计日志
// @Description 根据ID获取单个审计日志的详细信息
// @Tags 审计日志
// @Accept json
// @Produce json
// @Param id path int true "审计日志ID"
// @Success 200 {object} dto.SuccessResponse{data=dto.AuditLogResponse}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/audit-logs/{id} [get]
func (h *AuditLogHandler) GetAuditLogByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success:    false,
			Message:    "无效的审计日志ID",
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	auditLog, err := h.auditLogService.GetAuditLogByID(c.Request.Context(), uint(id))
	if err != nil {
		if appErr, ok := err.(*common.AppError); ok && appErr.Code == "AUDIT_LOG_NOT_FOUND" {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Success:    false,
				Message:    "审计日志不存在",
				StatusCode: http.StatusNotFound,
			})
			return
		}

		common.LogAppError(err, "获取审计日志失败")
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success:    false,
			Message:    "获取审计日志失败",
			StatusCode: http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponseDTO{
		Success: true,
		Message: "获取审计日志成功",
		Data:    auditLog,
	})
}

// GetUserAuditLogs 获取用户审计日志
// @Summary 获取用户审计日志
// @Description 获取指定用户的审计日志列表
// @Tags 审计日志
// @Accept json
// @Produce json
// @Param user_id path int true "用户ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} dto.PaginatedResponse[dto.AuditLogResponse]
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/users/{user_id}/audit-logs [get]
func (h *AuditLogHandler) GetUserAuditLogs(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success:    false,
			Message:    "无效的用户ID",
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	var req dto.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		common.LogAppError(err, "绑定分页参数失败")
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success:    false,
			Message:    "参数错误",
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	// 设置默认分页参数
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	auditLogs, total, err := h.auditLogService.GetUserAuditLogs(c.Request.Context(), uint(userID), &req)
	if err != nil {
		common.LogAppError(err, "获取用户审计日志失败")
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success:    false,
			Message:    "获取用户审计日志失败",
			StatusCode: http.StatusInternalServerError,
		})
		return
	}

	totalPages := int((total + int64(req.PageSize) - 1) / int64(req.PageSize))

	c.JSON(http.StatusOK, dto.PaginatedResponse[*dto.AuditLogResponse]{
		Data:       auditLogs,
		Total:      total,
		Page:       req.Page,
		Limit:      req.PageSize,
		TotalPages: totalPages,
	})
}

// GetResourceAuditLogs 获取资源审计日志
// @Summary 获取资源审计日志
// @Description 获取指定资源的审计日志列表
// @Tags 审计日志
// @Accept json
// @Produce json
// @Param resource path string true "资源类型"
// @Param resource_id path string true "资源ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} dto.PaginatedResponse[dto.AuditLogResponse]
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/resources/{resource}/{resource_id}/audit-logs [get]
func (h *AuditLogHandler) GetResourceAuditLogs(c *gin.Context) {
	resource := c.Param("resource")
	resourceID := c.Param("resource_id")

	if resource == "" || resourceID == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success:    false,
			Message:    "资源类型和资源ID不能为空",
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	var req dto.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		common.LogAppError(err, "绑定分页参数失败")
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success:    false,
			Message:    "参数错误",
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	// 设置默认分页参数
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	auditLogs, total, err := h.auditLogService.GetResourceAuditLogs(c.Request.Context(), resource, resourceID, &req)
	if err != nil {
		common.LogAppError(err, "获取资源审计日志失败")
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success:    false,
			Message:    "获取资源审计日志失败",
			StatusCode: http.StatusInternalServerError,
		})
		return
	}

	totalPages := int((total + int64(req.PageSize) - 1) / int64(req.PageSize))

	c.JSON(http.StatusOK, dto.PaginatedResponse[*dto.AuditLogResponse]{
		Data:       auditLogs,
		Total:      total,
		Page:       req.Page,
		Limit:      req.PageSize,
		TotalPages: totalPages,
	})
}

// CleanupOldLogs 清理旧的审计日志
// @Summary 清理旧的审计日志
// @Description 清理指定天数之前的审计日志
// @Tags 审计日志
// @Accept json
// @Produce json
// @Param days query int true "保留天数" minimum(1)
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/audit-logs/cleanup [post]
func (h *AuditLogHandler) CleanupOldLogs(c *gin.Context) {
	daysStr := c.Query("days")
	if daysStr == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success:    false,
			Message:    "保留天数不能为空",
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	days, err := strconv.Atoi(daysStr)
	if err != nil || days <= 0 {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success:    false,
			Message:    "保留天数必须是正整数",
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	if err := h.auditLogService.CleanupOldLogs(c.Request.Context(), days); err != nil {
		common.LogAppError(err, "清理旧审计日志失败")
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success:    false,
			Message:    "清理旧审计日志失败",
			StatusCode: http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponseDTO{
		Success: true,
		Message: "清理旧审计日志成功",
	})
}
