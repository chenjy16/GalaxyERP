package repositories

import (
	"context"
	"time"
	"github.com/galaxyerp/galaxyErp/internal/dto"
	"github.com/galaxyerp/galaxyErp/internal/models"
	"gorm.io/gorm"
)

// AuditLogRepository 审计日志仓储接口
type AuditLogRepository interface {
	BaseRepository[models.AuditLog]
	GetByUser(ctx context.Context, userID uint, req *dto.PaginationRequest) ([]*models.AuditLog, int64, error)
	GetByResource(ctx context.Context, resource string, resourceID string, req *dto.PaginationRequest) ([]*models.AuditLog, int64, error)
	GetByAction(ctx context.Context, action string, req *dto.PaginationRequest) ([]*models.AuditLog, int64, error)
	Search(ctx context.Context, req *dto.AuditLogSearchRequest) ([]*models.AuditLog, int64, error)
	DeleteOldLogs(ctx context.Context, days int) error
}

// AuditLogRepositoryImpl 审计日志仓储实现
type AuditLogRepositoryImpl struct {
	BaseRepository[models.AuditLog]
	db *gorm.DB
}

// NewAuditLogRepository 创建审计日志仓储实例
func NewAuditLogRepository(db *gorm.DB) AuditLogRepository {
	return &AuditLogRepositoryImpl{
		BaseRepository: NewBaseRepository[models.AuditLog](db),
		db:             db,
	}
}

// GetByUser 根据用户ID获取审计日志
func (r *AuditLogRepositoryImpl) GetByUser(ctx context.Context, userID uint, req *dto.PaginationRequest) ([]*models.AuditLog, int64, error) {
	var auditLogs []*models.AuditLog
	var total int64

	query := r.db.WithContext(ctx).Model(&models.AuditLog{}).Where("user_id = ?", userID)

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	err := query.Preload("User").
		Order("created_at DESC").
		Offset(req.GetOffset()).
		Limit(req.GetLimit()).
		Find(&auditLogs).Error

	return auditLogs, total, err
}

// GetByResource 根据资源获取审计日志
func (r *AuditLogRepositoryImpl) GetByResource(ctx context.Context, resource string, resourceID string, req *dto.PaginationRequest) ([]*models.AuditLog, int64, error) {
	var auditLogs []*models.AuditLog
	var total int64

	query := r.db.WithContext(ctx).Model(&models.AuditLog{}).
		Where("resource = ? AND resource_id = ?", resource, resourceID)

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	err := query.Preload("User").
		Order("created_at DESC").
		Offset(req.GetOffset()).
		Limit(req.GetLimit()).
		Find(&auditLogs).Error

	return auditLogs, total, err
}

// GetByAction 根据操作类型获取审计日志
func (r *AuditLogRepositoryImpl) GetByAction(ctx context.Context, action string, req *dto.PaginationRequest) ([]*models.AuditLog, int64, error) {
	var auditLogs []*models.AuditLog
	var total int64

	query := r.db.WithContext(ctx).Model(&models.AuditLog{}).Where("action = ?", action)

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	err := query.Preload("User").
		Order("created_at DESC").
		Offset(req.GetOffset()).
		Limit(req.GetLimit()).
		Find(&auditLogs).Error

	return auditLogs, total, err
}

// Search 搜索审计日志
func (r *AuditLogRepositoryImpl) Search(ctx context.Context, req *dto.AuditLogSearchRequest) ([]*models.AuditLog, int64, error) {
	var auditLogs []*models.AuditLog
	var total int64

	query := r.db.WithContext(ctx).Model(&models.AuditLog{})

	// 应用搜索条件
	if req.UserID != nil {
		query = query.Where("user_id = ?", *req.UserID)
	}
	if req.Username != "" {
		query = query.Joins("JOIN users ON users.id = audit_logs.user_id").
			Where("users.username LIKE ?", "%"+req.Username+"%")
	}
	if req.Action != "" {
		query = query.Where("action = ?", req.Action)
	}
	if req.Resource != "" {
		query = query.Where("resource = ?", req.Resource)
	}
	if req.ResourceID != "" {
		query = query.Where("resource_id = ?", req.ResourceID)
	}
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}
	if req.IPAddress != "" {
		query = query.Where("ip_address = ?", req.IPAddress)
	}
	if !req.StartTime.IsZero() {
		query = query.Where("created_at >= ?", req.StartTime)
	}
	if !req.EndTime.IsZero() {
		query = query.Where("created_at <= ?", req.EndTime)
	}
	if req.Search != "" {
		query = query.Where("description LIKE ? OR path LIKE ?", "%"+req.Search+"%", "%"+req.Search+"%")
	}

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	err := query.Preload("User").
		Order("created_at DESC").
		Offset(req.GetOffset()).
		Limit(req.GetLimit()).
		Find(&auditLogs).Error

	return auditLogs, total, err
}

// DeleteOldLogs 删除旧的审计日志
func (r *AuditLogRepositoryImpl) DeleteOldLogs(ctx context.Context, days int) error {
	cutoffDate := time.Now().AddDate(0, 0, -days)
	return r.db.WithContext(ctx).
		Where("created_at < ?", cutoffDate).
		Delete(&models.AuditLog{}).Error
}
