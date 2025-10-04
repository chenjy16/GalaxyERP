package services

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/galaxyerp/galaxyErp/internal/common"
	"github.com/galaxyerp/galaxyErp/internal/dto"
	"github.com/galaxyerp/galaxyErp/internal/models"
	"github.com/galaxyerp/galaxyErp/internal/repositories"
	"github.com/galaxyerp/galaxyErp/internal/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// AuditLogService 审计日志服务接口
type AuditLogService interface {
	// 记录审计日志
	LogAction(ctx context.Context, userID uint, username, action, resource, resourceID string, description string, oldValues, newValues interface{}) error
	LogActionWithContext(c *gin.Context, userID uint, username, action, resource, resourceID string, description string, oldValues, newValues interface{}) error
	LogError(ctx context.Context, userID uint, username, action, resource, resourceID string, err error, duration int64) error

	// 查询审计日志
	GetAuditLogs(ctx context.Context, req *dto.AuditLogSearchRequest) ([]*dto.AuditLogResponse, int64, error)
	GetAuditLogByID(ctx context.Context, id uint) (*dto.AuditLogResponse, error)
	GetUserAuditLogs(ctx context.Context, userID uint, req *dto.PaginationRequest) ([]*dto.AuditLogResponse, int64, error)
	GetResourceAuditLogs(ctx context.Context, resource string, resourceID string, req *dto.PaginationRequest) ([]*dto.AuditLogResponse, int64, error)

	// 清理旧日志
	CleanupOldLogs(ctx context.Context, days int) error
}

// AuditLogServiceImpl 审计日志服务实现
type AuditLogServiceImpl struct {
	auditLogRepo repositories.AuditLogRepository
	logger       *zap.Logger
}

// NewAuditLogService 创建审计日志服务实例
func NewAuditLogService(auditLogRepo repositories.AuditLogRepository, logger *zap.Logger) AuditLogService {
	return &AuditLogServiceImpl{
		auditLogRepo: auditLogRepo,
		logger:       logger,
	}
}

// LogAction 记录审计日志
func (s *AuditLogServiceImpl) LogAction(ctx context.Context, userID uint, username, action, resource, resourceID string, description string, oldValues, newValues interface{}) error {
	auditLog := &models.AuditLog{
		UserID:       userID,
		Username:     username,
		Action:       action,
		Resource:     resource,
		ResourceType: resource, // 设置ResourceType字段
		ResourceID:   resourceID,
		Description:  description,
		Status:       "success",
	}

	// 序列化旧值和新值
	if oldValues != nil {
		if oldJSON, err := json.Marshal(oldValues); err == nil {
			auditLog.OldValues = string(oldJSON)
		}
	}

	if newValues != nil {
		if newJSON, err := json.Marshal(newValues); err == nil {
			auditLog.NewValues = string(newJSON)
		}
	}

	// 计算变更内容
	if oldValues != nil && newValues != nil {
		changes := s.calculateChanges(oldValues, newValues)
		if changesJSON, err := json.Marshal(changes); err == nil {
			auditLog.Changes = string(changesJSON)
		}
	}

	if err := s.auditLogRepo.Create(ctx, auditLog); err != nil {
		appErr := common.NewAppErrorFromTypeWithCause("database", "AUDIT_LOG_CREATE_FAILED", "创建审计日志失败", err)
		common.LogAppError(appErr, "audit_log_create",
			utils.Uint("user_id", userID),
			utils.String("action", action),
			utils.String("resource", resource))
		return appErr
	}

	return nil
}

// LogActionWithContext 使用Gin上下文记录审计日志
func (s *AuditLogServiceImpl) LogActionWithContext(c *gin.Context, userID uint, username, action, resource, resourceID string, description string, oldValues, newValues interface{}) error {
	auditLog := &models.AuditLog{
		UserID:      userID,
		Username:    username,
		Action:      action,
		Resource:    resource,
		ResourceID:  resourceID,
		Method:      c.Request.Method,
		Path:        c.Request.URL.Path,
		Description: description,
		IPAddress:   c.ClientIP(),
		UserAgent:   c.Request.UserAgent(),
		Status:      "success",
	}

	// 序列化旧值和新值
	if oldValues != nil {
		if oldJSON, err := json.Marshal(oldValues); err == nil {
			auditLog.OldValues = string(oldJSON)
		}
	}

	if newValues != nil {
		if newJSON, err := json.Marshal(newValues); err == nil {
			auditLog.NewValues = string(newJSON)
		}
	}

	// 计算变更内容
	if oldValues != nil && newValues != nil {
		changes := s.calculateChanges(oldValues, newValues)
		if changesJSON, err := json.Marshal(changes); err == nil {
			auditLog.Changes = string(changesJSON)
		}
	}

	if err := s.auditLogRepo.Create(c.Request.Context(), auditLog); err != nil {
		s.logger.Error("Failed to create audit log",
			zap.Error(err),
			zap.Uint("user_id", userID),
			zap.String("action", action),
			zap.String("resource", resource),
		)
		return utils.NewAppErrorWithCause(utils.ErrorTypeDatabase, "AUDIT_LOG_CREATE_FAILED", "创建审计日志失败", err)
	}

	return nil
}

// LogError 记录错误审计日志
func (s *AuditLogServiceImpl) LogError(ctx context.Context, userID uint, username, action, resource, resourceID string, err error, duration int64) error {
	auditLog := &models.AuditLog{
		UserID:     userID,
		Username:   username,
		Action:     action,
		Resource:   resource,
		ResourceID: resourceID,
		Status:     "failed",
		ErrorMsg:   err.Error(),
		Duration:   duration,
	}

	if createErr := s.auditLogRepo.Create(ctx, auditLog); createErr != nil {
		appErr := common.NewAppErrorFromTypeWithCause("database", "AUDIT_LOG_CREATE_FAILED", "创建错误审计日志失败", createErr)
		common.LogAppError(appErr, "audit_log_error_create",
			utils.Uint("user_id", userID),
			utils.String("action", action),
			utils.String("resource", resource))
		return appErr
	}

	return nil
}

// GetAuditLogs 获取审计日志列表
func (s *AuditLogServiceImpl) GetAuditLogs(ctx context.Context, req *dto.AuditLogSearchRequest) ([]*dto.AuditLogResponse, int64, error) {
	auditLogs, total, err := s.auditLogRepo.Search(ctx, req)
	if err != nil {
		appErr := common.NewAppErrorFromTypeWithCause("database", "AUDIT_LOG_GET_FAILED", "获取审计日志失败", err)
		common.LogAppError(appErr, "audit_log_get_list")
		return nil, 0, appErr
	}

	responses := make([]*dto.AuditLogResponse, len(auditLogs))
	for i, log := range auditLogs {
		responses[i] = s.convertToResponse(log)
	}

	return responses, total, nil
}

// GetAuditLogByID 根据ID获取审计日志
func (s *AuditLogServiceImpl) GetAuditLogByID(ctx context.Context, id uint) (*dto.AuditLogResponse, error) {
	auditLog, err := s.auditLogRepo.GetByID(ctx, id)
	if err != nil {
		appErr := common.NewAppErrorFromTypeWithCause("database", "AUDIT_LOG_GET_FAILED", "获取审计日志失败", err)
		common.LogAppError(appErr, "audit_log_get_by_id", utils.Uint("id", id))
		return nil, appErr
	}

	if auditLog == nil {
		appErr := common.NewAppErrorFromType("business", "AUDIT_LOG_NOT_FOUND", "审计日志不存在")
		common.LogAppError(appErr, "audit_log_get_by_id", utils.Uint("id", id))
		return nil, appErr
	}

	return s.convertToResponse(auditLog), nil
}

// GetUserAuditLogs 获取用户审计日志
func (s *AuditLogServiceImpl) GetUserAuditLogs(ctx context.Context, userID uint, req *dto.PaginationRequest) ([]*dto.AuditLogResponse, int64, error) {
	auditLogs, total, err := s.auditLogRepo.GetByUser(ctx, userID, req)
	if err != nil {
		appErr := common.NewAppErrorFromTypeWithCause("database", "AUDIT_LOG_GET_FAILED", "获取用户审计日志失败", err)
		common.LogAppError(appErr, "audit_log_get_by_user", utils.Uint("user_id", userID))
		return nil, 0, appErr
	}

	responses := make([]*dto.AuditLogResponse, len(auditLogs))
	for i, log := range auditLogs {
		responses[i] = s.convertToResponse(log)
	}

	return responses, total, nil
}

// GetResourceAuditLogs 获取资源审计日志
func (s *AuditLogServiceImpl) GetResourceAuditLogs(ctx context.Context, resource string, resourceID string, req *dto.PaginationRequest) ([]*dto.AuditLogResponse, int64, error) {
	auditLogs, total, err := s.auditLogRepo.GetByResource(ctx, resource, resourceID, req)
	if err != nil {
		appErr := common.NewAppErrorFromTypeWithCause("database", "AUDIT_LOG_GET_FAILED", "获取资源审计日志失败", err)
		common.LogAppError(appErr, "audit_log_get_by_resource",
			utils.String("resource", resource),
			utils.String("resource_id", resourceID))
		return nil, 0, appErr
	}

	responses := make([]*dto.AuditLogResponse, len(auditLogs))
	for i, log := range auditLogs {
		responses[i] = s.convertToResponse(log)
	}

	return responses, total, nil
}

// CleanupOldLogs 清理旧日志
func (s *AuditLogServiceImpl) CleanupOldLogs(ctx context.Context, days int) error {
	if err := s.auditLogRepo.DeleteOldLogs(ctx, days); err != nil {
		appErr := common.NewAppErrorFromTypeWithCause("database", "AUDIT_LOG_CLEANUP_FAILED", "清理旧审计日志失败", err)
		common.LogAppError(appErr, "audit_log_cleanup", utils.Int("days", days))
		return appErr
	}

	s.logger.Info("Successfully cleaned up old audit logs", zap.Int("days", days))
	return nil
}

// convertToResponse 转换为响应DTO
func (s *AuditLogServiceImpl) convertToResponse(auditLog *models.AuditLog) *dto.AuditLogResponse {
	return &dto.AuditLogResponse{
		ID:          auditLog.ID,
		UserID:      auditLog.UserID,
		Username:    auditLog.Username,
		Action:      auditLog.Action,
		Resource:    auditLog.Resource,
		ResourceID:  auditLog.ResourceID,
		Method:      auditLog.Method,
		Path:        auditLog.Path,
		Description: auditLog.Description,
		IPAddress:   auditLog.IPAddress,
		UserAgent:   auditLog.UserAgent,
		OldValues:   auditLog.OldValues,
		NewValues:   auditLog.NewValues,
		Changes:     auditLog.Changes,
		Status:      auditLog.Status,
		ErrorMsg:    auditLog.ErrorMsg,
		Duration:    auditLog.Duration,
		CreatedAt:   auditLog.CreatedAt,
	}
}

// calculateChanges 计算变更内容
func (s *AuditLogServiceImpl) calculateChanges(oldValues, newValues interface{}) map[string]interface{} {
	changes := make(map[string]interface{})

	oldVal := reflect.ValueOf(oldValues)
	newVal := reflect.ValueOf(newValues)

	if oldVal.Type() != newVal.Type() {
		return changes
	}

	switch oldVal.Kind() {
	case reflect.Struct:
		s.compareStructs(oldVal, newVal, changes, "")
	case reflect.Map:
		s.compareMaps(oldVal, newVal, changes, "")
	case reflect.Slice, reflect.Array:
		s.compareSlices(oldVal, newVal, changes, "")
	default:
		if !reflect.DeepEqual(oldValues, newValues) {
			changes["value"] = map[string]interface{}{
				"old": oldValues,
				"new": newValues,
			}
		}
	}

	return changes
}

// compareStructs 比较结构体
func (s *AuditLogServiceImpl) compareStructs(oldVal, newVal reflect.Value, changes map[string]interface{}, prefix string) {
	for i := 0; i < oldVal.NumField(); i++ {
		field := oldVal.Type().Field(i)
		if !field.IsExported() {
			continue
		}

		fieldName := field.Name
		if prefix != "" {
			fieldName = prefix + "." + fieldName
		}

		oldFieldVal := oldVal.Field(i)
		newFieldVal := newVal.Field(i)

		if !reflect.DeepEqual(oldFieldVal.Interface(), newFieldVal.Interface()) {
			changes[fieldName] = map[string]interface{}{
				"old": oldFieldVal.Interface(),
				"new": newFieldVal.Interface(),
			}
		}
	}
}

// compareMaps 比较映射
func (s *AuditLogServiceImpl) compareMaps(oldVal, newVal reflect.Value, changes map[string]interface{}, prefix string) {
	for _, key := range oldVal.MapKeys() {
		keyStr := fmt.Sprintf("%v", key.Interface())
		if prefix != "" {
			keyStr = prefix + "." + keyStr
		}

		oldMapVal := oldVal.MapIndex(key)
		newMapVal := newVal.MapIndex(key)

		if !newMapVal.IsValid() {
			changes[keyStr] = map[string]interface{}{
				"old": oldMapVal.Interface(),
				"new": nil,
			}
		} else if !reflect.DeepEqual(oldMapVal.Interface(), newMapVal.Interface()) {
			changes[keyStr] = map[string]interface{}{
				"old": oldMapVal.Interface(),
				"new": newMapVal.Interface(),
			}
		}
	}

	// 检查新增的键
	for _, key := range newVal.MapKeys() {
		oldMapVal := oldVal.MapIndex(key)
		if !oldMapVal.IsValid() {
			keyStr := fmt.Sprintf("%v", key.Interface())
			if prefix != "" {
				keyStr = prefix + "." + keyStr
			}
			changes[keyStr] = map[string]interface{}{
				"old": nil,
				"new": newVal.MapIndex(key).Interface(),
			}
		}
	}
}

// compareSlices 比较切片
func (s *AuditLogServiceImpl) compareSlices(oldVal, newVal reflect.Value, changes map[string]interface{}, prefix string) {
	oldLen := oldVal.Len()
	newLen := newVal.Len()

	maxLen := oldLen
	if newLen > maxLen {
		maxLen = newLen
	}

	for i := 0; i < maxLen; i++ {
		indexStr := fmt.Sprintf("[%d]", i)
		if prefix != "" {
			indexStr = prefix + indexStr
		}

		if i >= oldLen {
			changes[indexStr] = map[string]interface{}{
				"old": nil,
				"new": newVal.Index(i).Interface(),
			}
		} else if i >= newLen {
			changes[indexStr] = map[string]interface{}{
				"old": oldVal.Index(i).Interface(),
				"new": nil,
			}
		} else if !reflect.DeepEqual(oldVal.Index(i).Interface(), newVal.Index(i).Interface()) {
			changes[indexStr] = map[string]interface{}{
				"old": oldVal.Index(i).Interface(),
				"new": newVal.Index(i).Interface(),
			}
		}
	}
}
